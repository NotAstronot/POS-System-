import React, { useState, useEffect, useCallback, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Search, Barcode, Plus, Minus, Trash2, ShoppingCart, User,
  UtensilsCrossed, Coffee, Cookie, Box, Percent, Printer,
  CreditCard, Smartphone, Banknote, QrCode, X, ChevronRight,
  Receipt, ScanLine, Package, AlertTriangle, Hash, Copy, DollarSign,
  Wifi, WifiOff, RefreshCw, UploadCloud
} from 'lucide-react';
import { useStore, CartItem } from '../store/useStore';
import { productService, categoryService, availabilityService, revenueService, type Product, type Category } from '../services/api';
import {
  isOnline, getQueue, getFailed, syncOrders, dismissFailed,
  saveProductsCache, loadProductsCache, saveCategoriesCache, loadCategoriesCache,
  saveAvailabilityCache, loadAvailabilityCache,
} from '../services/offline';
import toast from 'react-hot-toast';
import { PaymentModal } from '../components/PaymentModal';
import { ReceiptModal } from '../components/ReceiptModal';
import { ShiftModal } from '../components/ShiftModal';
import { SyncStatusModal } from '../components/SyncStatusModal';

const ICON_MAP: Record<string, React.ReactNode> = {
  'utensils-crossed': <UtensilsCrossed className="w-5 h-5" />,
  'coffee': <Coffee className="w-5 h-5" />,
  'cookie': <Cookie className="w-5 h-5" />,
  'box': <Package className="w-5 h-5" />,
};

const DISABLED_PRODUCT_IDS: string[] = [];

const PRODUCT_IMAGES: Record<string, string> = {
  'BRGR-001': '/asset-image/burger-beef.jpg',
  'BRGR-002': '/asset-image/burger-chicken.jpg',
  'NASI-001': '/asset-image/nasi-goreng-spesial.jpg',
  'NASI-002': '/asset-image/nasi-goreng-ayam.jpg',
  'NASI-003': '/asset-image/nasi-goreng-seafood.jpg',
  'MIE-001': '/asset-image/mie-goreng.jpg',
  'MUM-001': '/asset-image/es-teh-manis.jpg',
  'MUM-002': '/asset-image/kopi-susu.jpg',
  'MUM-003': '/asset-image/jus-alpukat.jpg',
  'MUM-004': '/asset-image/jus-mangga.jpg',
  'MUM-005': '/asset-image/air-mineral.jpg',
  'MUM-006': '/asset-image/milkshake-coklat.jpg',
  'SNK-001': '/asset-image/french-fries.jpg',
  'SNK-002': '/asset-image/chicken-wings.jpg',
  'SNK-003': '/asset-image/onion-rings.jpg',
  'SNK-004': '/asset-image/spring-rolls.jpg',
  'SNK-005': '/asset-image/potato-wedges.jpg',
  'MUM-007': '/asset-image/lemon-tea.jpg',
};

export function POSPage() {
  const navigate = useNavigate();
  const user = useStore((s) => s.user);
  const { items, orderType, tableNumber, customerName, customerPhone, discountCode, discountAmount } = useStore();
  const { addItem, removeItem, updateQuantity, clearCart, setOrderType, setTableNumber, setCustomerName, setCustomerPhone, setDiscountCode, setDiscountAmount } = useStore();

  const [categories, setCategories] = useState<Category[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [activeCategory, setActiveCategory] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [showSearch, setShowSearch] = useState(false);
  const [loading, setLoading] = useState(true);
  const [showPayment, setShowPayment] = useState(false);
  const [showReceipt, setShowReceipt] = useState(false);
  const [lastOrder, setLastOrder] = useState<any>(null);
  const [showShift, setShowShift] = useState(false);
  const [barcodeBuffer, setBarcodeBuffer] = useState('');
  const [unavailableIds, setUnavailableIds] = useState<Set<string>>(new Set());
  const [todayRevenue, setTodayRevenue] = useState<number>(0);
  const [online, setOnline] = useState<boolean>(isOnline());
  const [pendingCount, setPendingCount] = useState<number>(0);
  const [failedCount, setFailedCount] = useState<number>(0);
  const [syncing, setSyncing] = useState<boolean>(false);
  const [showFailedList, setShowFailedList] = useState(false);
  const today = new Date().toISOString().slice(0, 10);
  const barcodeTimer = useRef<ReturnType<typeof setTimeout>>();
  const searchRef = useRef<HTMLInputElement>(null);

  const refreshSyncCounts = () => {
    setPendingCount(getQueue().length);
    setFailedCount(getFailed().length);
  };

  const handleSync = async () => {
    if (!isOnline()) {
      toast.error('Tidak ada koneksi internet');
      return;
    }
    setSyncing(true);
    try {
      const summary = await syncOrders();
      refreshSyncCounts();
      if (summary.total > 0) {
        if (summary.failed > 0) {
          toast.error(`${summary.synced} order tersinkron, ${summary.failed} gagal (cek daftar sync)`);
        } else {
          toast.success(`${summary.synced} order offline berhasil disinkronkan`);
        }
        loadProducts(activeCategory);
        refreshRevenue();
      }
    } catch {
      toast.error('Sinkronisasi gagal — periksa koneksi ke server');
    } finally {
      setSyncing(false);
    }
  };

  useEffect(() => {
    refreshSyncCounts();
    const onOnline = () => {
      setOnline(true);
      if (getQueue().length > 0) handleSync();
    };
    const onOffline = () => setOnline(false);
    window.addEventListener('online', onOnline);
    window.addEventListener('offline', onOffline);
    const timer = setInterval(() => {
      if (isOnline() && getQueue().length > 0) handleSync();
    }, 30000);
    return () => {
      window.removeEventListener('online', onOnline);
      window.removeEventListener('offline', onOffline);
      clearInterval(timer);
    };
  }, []);

  const refreshRevenue = async () => {
    try {
      const res = await revenueService.summary('daily');
      setTodayRevenue(res.data.data.total_sales);
    } catch {}
  };

  useEffect(() => {
    refreshRevenue();
    const timer = setInterval(refreshRevenue, 30000);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    if (!user) { navigate('/login'); return; }
    loadData();
  }, [user]);

  useEffect(() => {
    if (showSearch && searchRef.current) searchRef.current.focus();
  }, [showSearch]);

  // Barcode scanner listener
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;
      if (e.key === 'Enter' && barcodeBuffer.length > 3) {
        handleBarcode(barcodeBuffer);
        setBarcodeBuffer('');
        return;
      }
      if (e.key.length === 1 && !e.ctrlKey && !e.metaKey) {
        setBarcodeBuffer((prev) => prev + e.key);
        clearTimeout(barcodeTimer.current);
        barcodeTimer.current = setTimeout(() => setBarcodeBuffer(''), 100);
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [barcodeBuffer]);

  const applyImages = (list: Product[]): Product[] =>
    list.map((p) => ({
      ...p,
      image_url: p.image_url || PRODUCT_IMAGES[p.code] || '',
    }));

  const loadData = async () => {
    try {
      const outletId = user?.outletId || 'b0000000-0000-0000-0000-000000000001';
      const [catRes, prodRes] = await Promise.all([
        categoryService.list(outletId),
        productService.list(),
      ]);
      setCategories(catRes.data.data || []);
      saveCategoriesCache(catRes.data.data || []);
      const list = applyImages(prodRes.data.data || []);
      setProducts(list);
      saveProductsCache(list);
    } catch (err) {
      // Offline / server tidak terjangkau: pakai cache terakhir agar tetap bisa berjualan
      const cached = loadProductsCache();
      if (cached) {
        setCategories(loadCategoriesCache() || MOCK_CATEGORIES);
        setProducts(cached);
      } else {
        setCategories(MOCK_CATEGORIES);
        setProducts(MOCK_PRODUCTS);
      }
    }
    try {
      const availRes = await availabilityService.list(today);
      const ids = new Set(availRes.data.data.filter((a: any) => !a.is_available).map((a: any) => String(a.product_id)));
      setUnavailableIds(ids);
      saveAvailabilityCache(availRes.data.data || []);
    } catch {
      const cachedAvail = loadAvailabilityCache();
      if (cachedAvail) {
        setUnavailableIds(new Set(cachedAvail.filter((a: any) => !a.is_available).map((a: any) => String(a.product_id))));
      }
    }
    setLoading(false);
  };

  const loadProducts = async (categoryId: string) => {
    setLoading(true);
    try {
      const outletId = user?.outletId || 'b0000000-0000-0000-0000-000000000001';
      let list: Product[];
      if (categoryId === 'all') {
        const res = await productService.list();
        list = res.data.data || [];
      } else {
        const res = await productService.byCategory(outletId, categoryId);
        list = res.data.data || [];
      }
      const mapped = applyImages(list);
      setProducts(mapped);
      saveProductsCache(mapped);
    } catch {
      // Offline: pakai cache produk terakhir (stok sesuai saat terakhir online)
      const cached = loadProductsCache();
      setProducts(cached ? cached : MOCK_PRODUCTS);
    }
    setLoading(false);
  };

  useEffect(() => {
    loadProducts(activeCategory);
  }, [activeCategory, categories]);

  const handleBarcode = async (code: string) => {
    try {
      const outletId = user?.outletId || 'b0000000-0000-0000-0000-000000000001';
      const res = await productService.barcode(outletId, code);
      const product = res.data.data;
      if (product) {
        const id = String(product.id);
        if (unavailableIds.has(id)) {
          toast.error(`${product.name} sedang tidak tersedia`);
          return;
        }
        if (product.stock <= 0) {
          toast.error(`${product.name} sudah habis (stok 0)`);
          return;
        }
        const existingItem = items.find((i) => i.productId === product.id);
        if (existingItem && existingItem.quantity >= product.stock) {
          toast.error(`${product.name} stok tidak mencukupi (tersedia ${product.stock})`);
          return;
        }
        addItem({
          productId: product.id,
          name: product.name,
          code: product.code,
          price: product.base_price,
          quantity: 1,
          variantLabel: '',
          note: '',
          customizations: [],
        });
        toast.success(`${product.name} ditambahkan!`);
      }
    } catch {
      toast.error('Produk tidak ditemukan');
    }
  };

  const handleSearch = useCallback(async (q: string) => {
    setSearchQuery(q);
    if (q.length < 2) {
      loadProducts(activeCategory);
      return;
    }
    try {
      const outletId = user?.outletId || 'b0000000-0000-0000-0000-000000000001';
      const res = await productService.search(outletId, q);
      setProducts(applyImages(res.data.data || []));
    } catch {
      // Filter mock
      setProducts(MOCK_PRODUCTS.filter((p) =>
        p.name.toLowerCase().includes(q.toLowerCase()) ||
        p.code.toLowerCase().includes(q.toLowerCase())
      ));
    }
  }, [activeCategory, user]);

  const toggleAvailability = async (product: Product) => {
    const id = String(product.id);
    const wasAvailable = !unavailableIds.has(id);
    try {
      await availabilityService.set(product.id, today, !wasAvailable);
      const next = new Set(unavailableIds);
      if (wasAvailable) next.add(id); else next.delete(id);
      setUnavailableIds(next);
      toast(wasAvailable ? `${product.name} tidak tersedia` : `${product.name} tersedia`, { icon: wasAvailable ? '🔴' : '✅', duration: 1500 });
    } catch { toast.error('Gagal update'); }
  };

  const handleAddProduct = (product: Product) => {
    const id = String(product.id);
    if (unavailableIds.has(id)) {
      toast.error(`${product.name} sedang tidak tersedia`);
      return;
    }
    if (product.stock <= 0) {
      toast.error(`${product.name} sudah habis (stok 0)`);
      return;
    }
    const existingItem = items.find((i) => i.productId === product.id);
    if (existingItem && existingItem.quantity >= product.stock) {
      toast.error(`${product.name} stok tidak mencukupi (tersedia ${product.stock})`);
      return;
    }
    addItem({
      productId: product.id,
      name: product.name,
      code: product.code,
      price: product.base_price,
      quantity: 1,
      variantLabel: '',
      note: '',
      customizations: [],
    });
    toast(`${product.name} +1`, { icon: '✅', duration: 1000 });
  };

  const subtotal = items.reduce((sum, i) => sum + i.subtotal, 0);

  // Valid discount codes
  const validDiscounts: Record<string, { type: 'percent'; value: number }> = {
    'HEMAT10': { type: 'percent', value: 10 },
    'DISKON15': { type: 'percent', value: 15 },
    'PROMO20': { type: 'percent', value: 20 },
    ' Member': { type: 'percent', value: 5 },
  };
  const discount = discountAmount > 0
    ? discountAmount
    : (validDiscounts[discountCode.toUpperCase().trim()] ? subtotal * validDiscounts[discountCode.toUpperCase().trim()].value / 100 : 0);
  const tax = (subtotal - discount) * 0.11;
  const grandTotal = subtotal - discount + tax;

  const handlePaymentSuccess = (orderResult: any) => {
    setLastOrder(orderResult);
    setShowPayment(false);
    setShowReceipt(true);
    refreshSyncCounts();
    // Refresh stok & pendapatan harian secara real-time setelah transaksi
    loadProducts(activeCategory);
    refreshRevenue();
  };

  // Show shift modal on mount if no active shift
  useEffect(() => {
    const shift = useStore.getState().activeShift;
    if (!shift && user) {
      setShowShift(true);
    }
  }, [user]);

  return (
    <div className="flex flex-col lg:flex-row gap-4 h-full">
      {/* Product Grid */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Top bar */}
        <div className="flex items-center gap-3 mb-4 flex-wrap">
          <div className="flex items-center gap-2">
            <button
              onClick={() => setOrderType('dine_in')}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                orderType === 'dine_in'
                  ? 'pos-btn-primary'
                  : 'pos-btn-secondary'
              }`}
            >
              <UtensilsCrossed className="w-4 h-4 inline mr-1.5" />
              Dine In
            </button>
            <button
              onClick={() => setOrderType('takeaway')}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                orderType === 'takeaway'
                  ? 'pos-btn-primary'
                  : 'pos-btn-secondary'
              }`}
            >
              <Package className="w-4 h-4 inline mr-1.5" />
              Takeaway
            </button>
          </div>

          {orderType === 'dine_in' && (
            <div className="flex items-center gap-2">
              <Hash className="w-4 h-4 text-gray-500" />
              <input
                type="text"
                placeholder="No Meja"
                value={tableNumber}
                onChange={(e) => setTableNumber(e.target.value)}
                className="w-20 pos-input text-sm py-1.5"
              />
            </div>
          )}

          <div className="flex-1" />

          <div className="flex items-center gap-2">
            {/* Status online/offline + sync */}
            <div
              className={`hidden md:flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg border text-xs font-medium ${
                online
                  ? 'bg-success/10 border-success/30 text-success'
                  : 'bg-danger/10 border-danger/30 text-danger'
              }`}
            >
              {online ? <Wifi className="w-3.5 h-3.5" /> : <WifiOff className="w-3.5 h-3.5" />}
              {online ? 'Online' : 'Offline'}
            </div>
            {(pendingCount > 0 || failedCount > 0) && (
              <button
                onClick={() => setShowFailedList(true)}
                className="flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg bg-warning/10 border border-warning/30 text-warning-dark text-xs font-medium hover:bg-warning/20"
                title="Lihat antrian sync"
              >
                <UploadCloud className="w-3.5 h-3.5" />
                <span>{pendingCount} antri</span>
                {failedCount > 0 && (
                  <span className="bg-danger text-white rounded-full px-1.5 text-[10px] font-bold">
                    {failedCount} gagal
                  </span>
                )}
              </button>
            )}
            <button
              onClick={handleSync}
              disabled={syncing || !online}
              className="pos-btn-secondary text-xs disabled:opacity-50"
              title="Sinkronkan order offline"
            >
              <RefreshCw className={`w-4 h-4 ${syncing ? 'animate-spin' : ''}`} />
              Sync
            </button>
            <div className="hidden md:flex items-center gap-2 px-3 py-1.5 rounded-lg bg-success/10 border border-success/30">
              <DollarSign className="w-4 h-4 text-success" />
              <span className="text-xs text-success font-medium">
                Hari ini: <strong>Rp {todayRevenue.toLocaleString()}</strong>
              </span>
            </div>
            <button
              onClick={() => setShowShift(true)}
              className="pos-btn-secondary text-xs"
            >
              Shift
            </button>
            <button
              onClick={() => setShowSearch(!showSearch)}
              className={`pos-btn-ghost ${showSearch ? 'text-brand' : ''}`}
            >
              <Search className="w-4 h-4" />
            </button>
            <div className="relative">
              <Barcode className="w-4 h-4 text-gray-500 absolute left-3 top-1/2 -translate-y-1/2" />
              <input
                placeholder="Scan barcode..."
                value={barcodeBuffer}
                onChange={(e) => handleBarcode(e.target.value)}
                className="pos-input text-sm pl-9 w-44 py-1.5"
              />
            </div>
          </div>
        </div>

        {/* Search bar */}
        {showSearch && (
          <div className="mb-4 animate-scan">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500" />
              <input
                ref={searchRef}
                type="text"
                placeholder="Cari produk... (min 2 karakter)"
                value={searchQuery}
                onChange={(e) => handleSearch(e.target.value)}
                className="pos-input pl-10 py-3 text-lg"
              />
            </div>
          </div>
        )}

        {/* Categories */}
        <div className="flex gap-2 mb-4 overflow-x-auto pb-2 scrollbar-none">
          <button
            onClick={() => setActiveCategory('all')}
            className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
              activeCategory === 'all'
                ? 'pos-btn-primary'
                : 'pos-btn-secondary'
            }`}
          >
            <Box className="w-4 h-4 inline mr-1.5" />
            Semua
          </button>
          {categories.map((cat) => (
            <button
              key={cat.id}
              onClick={() => setActiveCategory(cat.id)}
              className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
                activeCategory === cat.id
                  ? 'pos-btn-primary'
                  : 'pos-btn-secondary'
              }`}
              style={activeCategory === cat.id ? {} : { borderLeft: `3px solid ${cat.color}` }}
            >
              {ICON_MAP[cat.icon] || <Box className="w-4 h-4 inline mr-1.5" />}
              {cat.name}
            </button>
          ))}
        </div>

        {/* Products Grid */}
        <div className="flex-1 overflow-y-auto">
          {loading ? (
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
              {[...Array(10)].map((_, i) => (
                <div key={i} className="bg-neutral-200/50 rounded-xl h-32 animate-pulse" />
              ))}
            </div>
          ) : products.length === 0 ? (
            <div className="text-center py-12 text-neutral-500">
              <Package className="w-12 h-12 mx-auto mb-3 opacity-50" />
              <p>Tidak ada produk</p>
            </div>
          ) : (
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
              {products.map((product) => (
                <div
                  key={product.id}
                  className={unavailableIds.has(String(product.id)) || product.stock <= 0
                    ? 'pos-product-card-disabled'
                    : 'pos-product-card-hover'
                  }
                >
                  <button
                    onClick={() => handleAddProduct(product)}
                    disabled={unavailableIds.has(String(product.id)) || product.stock <= 0}
                    className="w-full text-left p-3 disabled:cursor-not-allowed"
                  >
                    <div className="pos-product-image">
                      {product.image_url ? (
                        <img src={product.image_url} alt={product.name} className="w-full h-full object-cover" />
                      ) : (
                        <Package className="w-8 h-8 text-neutral-400" />
                      )}
                    </div>
                    <p className="text-xs text-neutral-500 mb-0.5">{product.code}</p>
                    <p className="text-sm font-medium text-neutral-900 truncate">{product.name}</p>
                    <p className="text-brand font-bold text-sm mt-1">
                      Rp {product.base_price.toLocaleString()}
                    </p>
                    <div className="flex items-center gap-1.5 mt-1">
                      {product.stock > 0 && product.stock <= 5 && (
                        <span className="pos-badge-warning text-[10px]">Sisa {product.stock}</span>
                      )}
                      {product.stock > 5 && (
                        <span className="pos-badge-neutral text-[10px]">Stok {product.stock}</span>
                      )}
                      {product.stock <= 0 && (
                        <span className="pos-badge-danger text-[10px]">Habis</span>
                      )}
                      {product.has_variants && (
                        <span className="pos-badge-warning text-[10px]">Ada Varian</span>
                      )}
                    </div>
                    {unavailableIds.has(String(product.id)) && (
                      <span className="pos-badge-danger text-[10px] mt-1 block w-fit">Tidak Tersedia</span>
                    )}
                  </button>
                  <button
                    onClick={(e) => { e.stopPropagation(); toggleAvailability(product); }}
                    className={`absolute top-2 right-2 w-6 h-6 rounded-full flex items-center justify-center text-[10px] transition-colors ${
                      unavailableIds.has(String(product.id))
                        ? 'bg-danger text-white'
                        : 'bg-neutral-200 text-neutral-500 hover:bg-neutral-300'
                    }`}
                    title={unavailableIds.has(String(product.id)) ? 'Tandai tersedia' : 'Tandai tidak tersedia'}
                  >
                    {unavailableIds.has(String(product.id)) ? '✓' : '✕'}
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Cart */}
      <div className="w-full lg:w-96 flex flex-col pos-card overflow-hidden">
        <div className="px-4 py-3 border-b border-neutral-200 flex items-center justify-between">
          <h2 className="font-semibold text-neutral-900 flex items-center gap-2">
            <ShoppingCart className="w-4 h-4" />
            Pesanan
            {items.length > 0 && (
              <span className="pos-badge-primary">{items.length}</span>
            )}
          </h2>
          {items.length > 0 && (
            <button onClick={clearCart} className="text-xs text-danger hover:text-danger/80">
              <Trash2 className="w-3.5 h-3.5 inline mr-1" />
              Hapus Semua
            </button>
          )}
        </div>

        {/* Cart Items */}
        <div className="flex-1 overflow-y-auto p-4 space-y-2">
          {items.length === 0 ? (
            <div className="text-center py-12 text-neutral-500">
              <ShoppingCart className="w-12 h-12 mx-auto mb-3 opacity-30" />
              <p className="text-sm">Belum ada item</p>
              <p className="text-xs mt-1">Scan barcode atau pilih produk</p>
            </div>
          ) : (
            items.map((item) => (
              <div key={item.id} className="pos-cart-item">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-neutral-900 truncate">{item.name}</p>
                    {item.variantLabel && (
                      <p className="text-xs text-neutral-500">{item.variantLabel}</p>
                    )}
                    <p className="text-xs text-neutral-500 mt-0.5">Rp {item.price.toLocaleString()}</p>
                  </div>
                  <button
                    onClick={() => removeItem(item.id)}
                    className="opacity-0 group-hover:opacity-100 p-1 text-neutral-500 hover:text-danger transition-all"
                  >
                    <X className="w-3.5 h-3.5" />
                  </button>
                </div>

                <div className="flex items-center justify-between mt-2">
                  <div className="flex items-center gap-1">
                    <button
                      onClick={() => updateQuantity(item.id, item.quantity - 1)}
                      className="pos-qty-btn"
                    >
                      <Minus className="w-3 h-3" />
                    </button>
                    <span className="w-8 text-center text-sm font-medium">{item.quantity}</span>
                    <button
                      onClick={() => {
                        const product = products.find((p) => String(p.id) === item.productId);
                        if (product && item.quantity >= product.stock) {
                          toast.error(`${item.name} stok tidak mencukupi (tersedia ${product.stock})`);
                          return;
                        }
                        updateQuantity(item.id, item.quantity + 1);
                      }}
                      className="pos-qty-btn"
                    >
                      <Plus className="w-3 h-3" />
                    </button>
                  </div>
                  <p className="text-sm font-bold text-brand">
                    Rp {item.subtotal.toLocaleString()}
                  </p>
                </div>

                {/* Customer note */}
                {item.note && (
                  <p className="text-xs text-warning mt-1 italic">Catatan: {item.note}</p>
                )}
              </div>
            ))
          )}
        </div>

        {/* Customer info */}
        <div className="px-4 py-3 border-t border-neutral-200 space-y-2">
          <div className="flex gap-2">
            <div className="relative flex-1">
              <User className="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-neutral-500" />
              <input
                placeholder="Nama pelanggan"
                value={customerName}
                onChange={(e) => setCustomerName(e.target.value)}
                className="pos-input text-xs pl-8 py-1.5"
              />
            </div>
            <input
              placeholder="No. HP"
              value={customerPhone}
              onChange={(e) => setCustomerPhone(e.target.value)}
              className="pos-input text-xs w-28 py-1.5"
            />
          </div>
        </div>

        {/* Totals */}
        <div className="px-4 py-3 border-t border-neutral-200 space-y-1.5">
          <div className="flex justify-between text-sm text-neutral-600">
            <span>Subtotal</span>
            <span>Rp {subtotal.toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-sm text-neutral-600">
            <span>Diskon</span>
            <span className="text-danger">-Rp {discount.toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-sm text-neutral-600">
            <span>Pajak (11%)</span>
            <span>Rp {tax.toLocaleString()}</span>
          </div>
          <div className="pos-grand-total">
            <span>Total</span>
            <span className="pos-grand-total-value">Rp {grandTotal.toLocaleString()}</span>
          </div>
        </div>

        {/* Action buttons */}
        <div className="p-4 pt-2 space-y-2">
          <div className="flex gap-2">
            <input
              placeholder="Kode diskon"
              value={discountCode}
              onChange={(e) => setDiscountCode(e.target.value)}
              className={`pos-input text-sm flex-1 py-2 ${discountCode && !validDiscounts[discountCode.toUpperCase().trim()] && discountAmount === 0 ? 'border-danger' : ''}`}
            />
            {discountCode && (
              <button
                onClick={() => { setDiscountCode(''); setDiscountAmount(0); }}
                className="px-3 rounded-lg bg-neutral-200 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-300"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </div>
          {discountCode && discountAmount === 0 && (
            <p className={`text-xs ${validDiscounts[discountCode.toUpperCase().trim()] ? 'text-success' : 'text-danger'}`}>
              {validDiscounts[discountCode.toUpperCase().trim()]
                ? `Diskon ${validDiscounts[discountCode.toUpperCase().trim()].value}% diterapkan`
                : 'Kode diskon tidak valid'}
            </p>
          )}

          <button
            onClick={() => {
              if (!useStore.getState().activeShift) {
                toast.error('Buka shift terlebih dahulu!');
                setShowShift(true);
                return;
              }
              if (items.length === 0) {
                toast.error('Keranjang masih kosong');
                return;
              }
              setShowPayment(true);
            }}
            disabled={items.length === 0}
            className="w-full py-3.5 pos-btn-primary text-lg disabled:opacity-50 flex items-center justify-center gap-2"
          >
            <CreditCard className="w-5 h-5" />
            Bayar Rp {grandTotal.toLocaleString()}
          </button>
        </div>
      </div>

      {/* Modals */}
      {showPayment && (
        <PaymentModal
          grandTotal={grandTotal}
          subtotal={subtotal}
          discount={discount}
          tax={tax}
          onClose={() => setShowPayment(false)}
          onSuccess={handlePaymentSuccess}
        />
      )}
      {showReceipt && lastOrder && (
        <ReceiptModal
          receipt={lastOrder.receipt}
          onClose={() => { setShowReceipt(false); clearCart(); }}
        />
      )}
      {showShift && (
        <ShiftModal onClose={() => setShowShift(false)} />
      )}
      {showFailedList && (
        <SyncStatusModal
          onClose={() => setShowFailedList(false)}
          onSync={handleSync}
          syncing={syncing}
          onCountsChange={refreshSyncCounts}
        />
      )}
    </div>
  );
}

// Mock data for demo/offline
const MOCK_CATEGORIES: Category[] = [
  { id: 'd1', name: 'Makanan', icon: 'utensils-crossed', color: '#EF4444', sort_order: 1 },
  { id: 'd2', name: 'Minuman', icon: 'coffee', color: '#3B82F6', sort_order: 2 },
  { id: 'd3', name: 'Snack', icon: 'cookie', color: '#F59E0B', sort_order: 3 },
];

const MOCK_PRODUCTS = [
  { id: 'e1', code: 'BRGR-001', name: 'Beef Burger', description: '', base_price: 35000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: false, tax_rate: 11, sort_order: 1, stock: 50 },
  { id: 'e2', code: 'BRGR-002', name: 'Chicken Burger', description: '', base_price: 30000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: false, tax_rate: 11, sort_order: 2, stock: 45 },
  { id: 'e3', code: 'NASI-001', name: 'Nasi Goreng Spesial', description: '', base_price: 45000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: true, tax_rate: 11, sort_order: 3, stock: 30 },
  { id: 'e4', code: 'NASI-002', name: 'Nasi Goreng Ayam', description: '', base_price: 38000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: false, tax_rate: 11, sort_order: 4, stock: 25 },
  { id: 'e5', code: 'NASI-003', name: 'Nasi Goreng Seafood', description: '', base_price: 50000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: false, tax_rate: 11, sort_order: 5, stock: 20 },
  { id: 'e6', code: 'MIE-001', name: 'Mie Goreng Spesial', description: '', base_price: 40000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: false, tax_rate: 11, sort_order: 6, stock: 35 },
  { id: 'e7', code: 'MUM-001', name: 'Es Teh Manis', description: '', base_price: 8000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 1, stock: 100 },
  { id: 'e8', code: 'MUM-002', name: 'Kopi Susu', description: '', base_price: 25000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: true, tax_rate: 11, sort_order: 2, stock: 60 },
  { id: 'e9', code: 'MUM-003', name: 'Jus Alpukat', description: '', base_price: 20000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 3, stock: 40 },
  { id: 'e10', code: 'MUM-004', name: 'Jus Mangga', description: '', base_price: 18000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 4, stock: 35 },
  { id: 'e11', code: 'MUM-005', name: 'Air Mineral', description: '', base_price: 5000, unit: 'botol', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 5, stock: 200 },
  { id: 'e12', code: 'MUM-006', name: 'Milkshake Coklat', description: '', base_price: 30000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 6, stock: 25 },
  { id: 'e13', code: 'SNK-001', name: 'French Fries', description: '', base_price: 18000, unit: 'pcs', image_url: '', category_id: 'd3', has_variants: false, tax_rate: 11, sort_order: 1, stock: 4 },
  { id: 'e14', code: 'SNK-002', name: 'Chicken Wings', description: '', base_price: 35000, unit: 'pcs', image_url: '', category_id: 'd3', has_variants: false, tax_rate: 11, sort_order: 2, stock: 30 },
  { id: 'e15', code: 'SNK-003', name: 'Onion Rings', description: '', base_price: 20000, unit: 'pcs', image_url: '', category_id: 'd3', has_variants: false, tax_rate: 11, sort_order: 3, stock: 25 },
  { id: 'e16', code: 'SNK-004', name: 'Spring Rolls', description: '', base_price: 22000, unit: 'pcs', image_url: '', category_id: 'd3', has_variants: false, tax_rate: 11, sort_order: 4, stock: 20 },
  { id: 'e17', code: 'SNK-005', name: 'Potato Wedges', description: '', base_price: 25000, unit: 'pcs', image_url: '', category_id: 'd3', has_variants: false, tax_rate: 11, sort_order: 5, stock: 15 },
  { id: 'e18', code: 'MUM-007', name: 'Lemon Tea', description: '', base_price: 12000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 7, stock: 50 },
].map(p => ({ ...p, barcode: '', purchase_price: 0, min_stock: 0, is_service: false }));
