import { orderService } from './api';

export interface OfflineOrder {
  client_order_id: string;
  payload: any;
  created_at: string;
}

export interface FailedOrder extends OfflineOrder {
  reason: string;
  synced_at: string;
}

const QUEUE_KEY = 'pos-offline-queue';
const FAILED_KEY = 'pos-offline-failed';
const CACHE_PRODUCTS_KEY = 'pos-cache-products';
const CACHE_CATEGORIES_KEY = 'pos-cache-categories';
const CACHE_AVAIL_KEY = 'pos-cache-availability';

export function isOnline(): boolean {
  return typeof navigator !== 'undefined' ? navigator.onLine : true;
}

function read<T>(key: string, fallback: T): T {
  try {
    const raw = localStorage.getItem(key);
    return raw ? JSON.parse(raw) : fallback;
  } catch {
    return fallback;
  }
}

function write(key: string, value: unknown) {
  try {
    localStorage.setItem(key, JSON.stringify(value));
  } catch {
    /* storage penuh / tidak tersedia — abaikan */
  }
}

export function getQueue(): OfflineOrder[] {
  return read<OfflineOrder[]>(QUEUE_KEY, []);
}

export function getFailed(): FailedOrder[] {
  return read<FailedOrder[]>(FAILED_KEY, []);
}

export function enqueueOrder(payload: any): OfflineOrder {
  const entry: OfflineOrder = {
    client_order_id: `LCL-${Date.now()}-${crypto.randomUUID().slice(0, 8)}`,
    payload,
    created_at: new Date().toISOString(),
  };
  write(QUEUE_KEY, [...getQueue(), entry]);
  return entry;
}

export interface SyncSummary {
  total: number;
  synced: number;
  failed: number;
  failedItems: FailedOrder[];
}

export async function syncOrders(): Promise<SyncSummary> {
  const queue = getQueue();
  if (queue.length === 0) {
    return { total: 0, synced: 0, failed: 0, failedItems: [] };
  }
  if (!isOnline()) {
    throw new Error('offline');
  }

  const res = await orderService.sync({
    orders: queue.map((q) => ({ client_order_id: q.client_order_id, ...q.payload })),
  });

  const data = res.data?.data;
  const failedItems: FailedOrder[] = [];
  const processedIds = new Set<string>();

  if (data?.results) {
    for (const r of data.results) {
      processedIds.add(r.client_order_id);
      if (!r.success) {
        const original = queue.find((q) => q.client_order_id === r.client_order_id);
        if (original) {
          failedItems.push({
            ...original,
            reason: r.error || 'Gagal disinkronkan',
            synced_at: new Date().toISOString(),
          });
        }
      }
    }
  }

  write(
    QUEUE_KEY,
    queue.filter((q) => !processedIds.has(q.client_order_id))
  );

  const prevFailed = getFailed();
  write(FAILED_KEY, [...prevFailed, ...failedItems]);

  return {
    total: queue.length,
    synced: data?.synced ?? queue.length - failedItems.length,
    failed: failedItems.length,
    failedItems,
  };
}

export function dismissFailed(id: string) {
  write(
    FAILED_KEY,
    getFailed().filter((f) => f.client_order_id !== id)
  );
}

// --- Cache data (produk + availability) agar bisa dipakai saat offline ---
export function saveProductsCache(products: any[]) {
  write(CACHE_PRODUCTS_KEY, products);
}

export function loadProductsCache(): any[] | null {
  return read<any[] | null>(CACHE_PRODUCTS_KEY, null);
}

export function saveCategoriesCache(categories: any[]) {
  write(CACHE_CATEGORIES_KEY, categories);
}

export function loadCategoriesCache(): any[] | null {
  return read<any[] | null>(CACHE_CATEGORIES_KEY, null);
}

export function saveAvailabilityCache(list: any[]) {
  write(CACHE_AVAIL_KEY, list);
}

export function loadAvailabilityCache(): any[] | null {
  return read<any[] | null>(CACHE_AVAIL_KEY, null);
}

// --- Struk lokal untuk transaksi offline ---
export function buildLocalReceipt(
  entry: OfflineOrder,
  items: { name: string; variantLabel?: string; quantity: number; price: number; subtotal: number }[],
  subtotal: number,
  discount: number,
  tax: number,
  grandTotal: number,
  method: string,
  cashierName: string,
  storeName: string
): any {
  const now = new Date();
  const createdAt = now.toLocaleString('id-ID', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
  return {
    order: {
      id: 0,
      order_number: entry.client_order_id,
      order_type: entry.payload.order_type || 'dine_in',
      table_number: entry.payload.table_number || '',
      total: grandTotal,
      status: 'pending_sync',
      created_at: entry.created_at,
    },
    receipt: {
      store_name: storeName,
      store_address: '',
      store_phone: '',
      order_number: entry.client_order_id,
      order_type: entry.payload.order_type || 'dine_in',
      table_number: entry.payload.table_number || '',
      cashier_name: cashierName,
      items: items.map((i) => ({
        name: i.name,
        variant: i.variantLabel || '',
        qty: i.quantity,
        price: i.price,
        sub_total: i.subtotal,
      })),
      sub_total: subtotal,
      discount,
      tax,
      grand_total: grandTotal,
      amount_paid: grandTotal,
      change_amount: 0,
      payment_method: method,
      created_at: createdAt,
      offline: true,
    },
  };
}