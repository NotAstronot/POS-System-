import React, { useState, useEffect, useCallback, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Search, Barcode, Plus, Minus, Trash2, ShoppingCart, User,
  UtensilsCrossed, Coffee, Cookie, Box, Percent, Printer,
  CreditCard, Smartphone, Banknote, QrCode, X, ChevronRight,
  Receipt, ScanLine, Package, AlertTriangle, Hash, Copy
} from 'lucide-react';
import { useStore, CartItem } from '../store/useStore';
import { productService, categoryService, availabilityService, type Product, type Category } from '../services/api';
import toast from 'react-hot-toast';
import { PaymentModal } from '../components/PaymentModal';
import { ReceiptModal } from '../components/ReceiptModal';
import { ShiftModal } from '../components/ShiftModal';

const ICON_MAP: Record<string, React.ReactNode> = {
  'utensils-crossed': <UtensilsCrossed className="w-5 h-5" />,
  'coffee': <Coffee className="w-5 h-5" />,
  'cookie': <Cookie className="w-5 h-5" />,
  'box': <Package className="w-5 h-5" />,
};

const DISABLED_PRODUCT_IDS: string[] = [];

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
  const today = new Date().toISOString().slice(0, 10);
  const barcodeTimer = useRef<ReturnType<typeof setTimeout>>();
  const searchRef = useRef<HTMLInputElement>(null);

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

  const loadData = async () => {
    try {
      const outletId = user?.outletId || 'b0000000-0000-0000-0000-000000000001';
      const [catRes, prodRes] = await Promise.all([
        categoryService.list(outletId),
        // We'll load all products via categories
        Promise.resolve({ data: { data: [] } }),
      ]);
      setCategories(catRes.data.data || []);
    } catch (err) {
      setCategories(MOCK_CATEGORIES);
    }
    try {
      const availRes = await availabilityService.list(today);
      const ids = new Set(availRes.data.data.filter((a: any) => !a.is_available).map((a: any) => String(a.product_id)));
      setUnavailableIds(ids);
    } catch {}
    setLoading(false);
  };

  const loadProducts = async (categoryId: string) => {
    setLoading(true);
    try {
      const outletId = user?.outletId || 'b0000000-0000-0000-0000-000000000001';
      if (categoryId === 'all') {
        // Load all categories
        const allProducts: Product[] = [];
        for (const cat of categories) {
          const res = await productService.byCategory(outletId, cat.id);
          allProducts.push(...(res.data.data || []));
        }
        setProducts(allProducts);
      } else {
        const res = await productService.byCategory(outletId, categoryId);
        setProducts(res.data.data || []);
      }
    } catch {
      setProducts(MOCK_PRODUCTS);
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
      setProducts(res.data.data || []);
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
  const discount = discountAmount > 0 ? discountAmount : (discountCode ? subtotal * 0.1 : 0);
  const tax = (subtotal - discount) * 0.11;
  const grandTotal = subtotal - discount + tax;

  const handlePaymentSuccess = (orderResult: any) => {
    setLastOrder(orderResult);
    setShowPayment(false);
    setShowReceipt(true);
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
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-800 text-gray-400 hover:text-white'
              }`}
            >
              <UtensilsCrossed className="w-4 h-4 inline mr-1.5" />
              Dine In
            </button>
            <button
              onClick={() => setOrderType('takeaway')}
              className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                orderType === 'takeaway'
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-800 text-gray-400 hover:text-white'
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
            <button
              onClick={() => setShowShift(true)}
              className="pos-btn-secondary text-xs"
            >
              Shift
            </button>
            <button
              onClick={() => setShowSearch(!showSearch)}
              className={`pos-btn-ghost ${showSearch ? 'text-blue-400' : ''}`}
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
                ? 'bg-blue-600 text-white'
                : 'bg-gray-800 text-gray-400 hover:text-white'
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
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-800 text-gray-400 hover:text-white'
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
                <div key={i} className="bg-gray-800/50 rounded-xl h-32 animate-pulse" />
              ))}
            </div>
          ) : products.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <Package className="w-12 h-12 mx-auto mb-3 opacity-50" />
              <p>Tidak ada produk</p>
            </div>
          ) : (
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-3">
              {products.map((product) => (
                <div
                  key={product.id}
                  className={`relative bg-gray-900 border rounded-xl transition-all group ${unavailableIds.has(String(product.id)) ? 'border-red-900/50 opacity-60' : 'border-gray-800 hover:border-blue-600/50 hover:bg-gray-800'}`}
                >
                  <button onClick={() => handleAddProduct(product)} className="w-full text-left p-3">
                    <div className="w-full aspect-square rounded-lg bg-gradient-to-br from-gray-800 to-gray-700 mb-2 flex items-center justify-center overflow-hidden">
                      {product.image_url ? (
                        <img src={product.image_url} alt={product.name} className="w-full h-full object-cover" />
                      ) : (
                        <Package className="w-8 h-8 text-gray-600" />
                      )}
                    </div>
                    <p className="text-xs text-gray-500 mb-0.5">{product.code}</p>
                    <p className="text-sm font-medium text-white truncate">{product.name}</p>
                    <p className="text-blue-400 font-bold text-sm mt-1">
                      Rp {product.base_price.toLocaleString()}
                    </p>
                    {product.has_variants && (
                      <span className="pos-badge-warning text-[10px] mt-1">Ada Varian</span>
                    )}
                    {unavailableIds.has(String(product.id)) && (
                      <span className="pos-badge-danger text-[10px] mt-1 block w-fit">Tidak Tersedia</span>
                    )}
                  </button>
                  <button
                    onClick={(e) => { e.stopPropagation(); toggleAvailability(product); }}
                    className={`absolute top-2 right-2 w-6 h-6 rounded-full flex items-center justify-center text-[10px] transition-colors ${unavailableIds.has(String(product.id)) ? 'bg-red-600 text-white' : 'bg-gray-700 text-gray-400 hover:bg-gray-600'}`}
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
      <div className="w-full lg:w-96 flex flex-col bg-gray-900 border border-gray-800 rounded-2xl overflow-hidden">
        <div className="px-4 py-3 border-b border-gray-800 flex items-center justify-between">
          <h2 className="font-semibold text-white flex items-center gap-2">
            <ShoppingCart className="w-4 h-4" />
            Pesanan
            {items.length > 0 && (
              <span className="pos-badge bg-blue-600 text-white">{items.length}</span>
            )}
          </h2>
          {items.length > 0 && (
            <button onClick={clearCart} className="text-xs text-red-400 hover:text-red-300">
              <Trash2 className="w-3.5 h-3.5 inline mr-1" />
              Hapus Semua
            </button>
          )}
        </div>

        {/* Cart Items */}
        <div className="flex-1 overflow-y-auto p-4 space-y-2">
          {items.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <ShoppingCart className="w-12 h-12 mx-auto mb-3 opacity-30" />
              <p className="text-sm">Belum ada item</p>
              <p className="text-xs mt-1">Scan barcode atau pilih produk</p>
            </div>
          ) : (
            items.map((item) => (
              <div key={item.id} className="bg-gray-800/50 rounded-xl p-3 group hover:bg-gray-800 transition-colors">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-white truncate">{item.name}</p>
                    {item.variantLabel && (
                      <p className="text-xs text-gray-400">{item.variantLabel}</p>
                    )}
                    <p className="text-xs text-gray-500 mt-0.5">Rp {item.price.toLocaleString()}</p>
                  </div>
                  <button
                    onClick={() => removeItem(item.id)}
                    className="opacity-0 group-hover:opacity-100 p-1 text-gray-500 hover:text-red-400 transition-all"
                  >
                    <X className="w-3.5 h-3.5" />
                  </button>
                </div>

                <div className="flex items-center justify-between mt-2">
                  <div className="flex items-center gap-1">
                    <button
                      onClick={() => updateQuantity(item.id, item.quantity - 1)}
                      className="w-7 h-7 rounded-lg bg-gray-700 hover:bg-gray-600 flex items-center justify-center text-white text-sm"
                    >
                      <Minus className="w-3 h-3" />
                    </button>
                    <span className="w-8 text-center text-sm font-medium">{item.quantity}</span>
                    <button
                      onClick={() => updateQuantity(item.id, item.quantity + 1)}
                      className="w-7 h-7 rounded-lg bg-gray-700 hover:bg-gray-600 flex items-center justify-center text-white text-sm"
                    >
                      <Plus className="w-3 h-3" />
                    </button>
                  </div>
                  <p className="text-sm font-bold text-blue-400">
                    Rp {item.subtotal.toLocaleString()}
                  </p>
                </div>

                {/* Customer note */}
                {item.note && (
                  <p className="text-xs text-amber-400 mt-1 italic">Catatan: {item.note}</p>
                )}
              </div>
            ))
          )}
        </div>

        {/* Customer info */}
        <div className="px-4 py-3 border-t border-gray-800 space-y-2">
          <div className="flex gap-2">
            <div className="relative flex-1">
              <User className="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-gray-500" />
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
        <div className="px-4 py-3 border-t border-gray-800 space-y-1.5">
          <div className="flex justify-between text-sm text-gray-400">
            <span>Subtotal</span>
            <span>Rp {subtotal.toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-sm text-gray-400">
            <span>Diskon</span>
            <span className="text-red-400">-Rp {discount.toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-sm text-gray-400">
            <span>Pajak (11%)</span>
            <span>Rp {tax.toLocaleString()}</span>
          </div>
          <div className="flex justify-between text-lg font-bold text-white pt-2 border-t border-gray-700">
            <span>Total</span>
            <span className="text-blue-400">Rp {grandTotal.toLocaleString()}</span>
          </div>
        </div>

        {/* Action buttons */}
        <div className="p-4 pt-2 space-y-2">
          <div className="flex gap-2">
            <input
              placeholder="Kode diskon"
              value={discountCode}
              onChange={(e) => setDiscountCode(e.target.value)}
              className="pos-input text-sm flex-1 py-2"
            />
            {discountCode && (
              <button
                onClick={() => setDiscountCode('')}
                className="px-3 rounded-lg bg-gray-700 text-gray-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </div>

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
            className="w-full py-3.5 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-700 disabled:text-gray-500 text-white font-bold rounded-xl transition-colors text-lg flex items-center justify-center gap-2"
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
    </div>
  );
}

// Mock data for demo/offline
const MOCK_CATEGORIES: Category[] = [
  { id: 'd1', name: 'Makanan', icon: 'utensils-crossed', color: '#EF4444', sort_order: 1 },
  { id: 'd2', name: 'Minuman', icon: 'coffee', color: '#3B82F6', sort_order: 2 },
  { id: 'd3', name: 'Snack', icon: 'cookie', color: '#F59E0B', sort_order: 3 },
];

const MOCK_PRODUCTS: Product[] = [
  { id: 'e1', code: 'BRGR-001', name: 'Beef Burger', description: '', base_price: 35000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: false, tax_rate: 11, sort_order: 1 },
  { id: 'e2', code: 'BRGR-002', name: 'Chicken Burger', description: '', base_price: 30000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: false, tax_rate: 11, sort_order: 2 },
  { id: 'e3', code: 'NASI-001', name: 'Nasi Goreng Spesial', description: '', base_price: 45000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: true, tax_rate: 11, sort_order: 3 },
  { id: 'e4', code: 'NASI-002', name: 'Nasi Goreng Ayam', description: '', base_price: 38000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: false, tax_rate: 11, sort_order: 4 },
  { id: 'e5', code: 'NASI-003', name: 'Nasi Goreng Seafood', description: '', base_price: 50000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: false, tax_rate: 11, sort_order: 5 },
  { id: 'e6', code: 'MIE-001', name: 'Mie Goreng Spesial', description: '', base_price: 40000, unit: 'pcs', image_url: '', category_id: 'd1', has_variants: false, tax_rate: 11, sort_order: 6 },
  { id: 'e7', code: 'MUM-001', name: 'Es Teh Manis', description: '', base_price: 8000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 1 },
  { id: 'e8', code: 'MUM-002', name: 'Kopi Susu', description: '', base_price: 25000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: true, tax_rate: 11, sort_order: 2 },
  { id: 'e9', code: 'MUM-003', name: 'Jus Alpukat', description: '', base_price: 20000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 3 },
  { id: 'e10', code: 'MUM-004', name: 'Jus Mangga', description: '', base_price: 18000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 4 },
  { id: 'e11', code: 'MUM-005', name: 'Air Mineral', description: '', base_price: 5000, unit: 'botol', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 5 },
  { id: 'e12', code: 'MUM-006', name: 'Milkshake Coklat', description: '', base_price: 30000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 6 },
  { id: 'e13', code: 'SNK-001', name: 'French Fries', description: '', base_price: 18000, unit: 'pcs', image_url: '', category_id: 'd3', has_variants: false, tax_rate: 11, sort_order: 1 },
  { id: 'e14', code: 'SNK-002', name: 'Chicken Wings', description: '', base_price: 35000, unit: 'pcs', image_url: '', category_id: 'd3', has_variants: false, tax_rate: 11, sort_order: 2 },
  { id: 'e15', code: 'SNK-003', name: 'Onion Rings', description: '', base_price: 20000, unit: 'pcs', image_url: '', category_id: 'd3', has_variants: false, tax_rate: 11, sort_order: 3 },
  { id: 'e16', code: 'SNK-004', name: 'Spring Rolls', description: '', base_price: 22000, unit: 'pcs', image_url: '', category_id: 'd3', has_variants: false, tax_rate: 11, sort_order: 4 },
  { id: 'e17', code: 'SNK-005', name: 'Potato Wedges', description: '', base_price: 25000, unit: 'pcs', image_url: '', category_id: 'd3', has_variants: false, tax_rate: 11, sort_order: 5 },
  { id: 'e18', code: 'MUM-007', name: 'Lemon Tea', description: '', base_price: 12000, unit: 'gelas', image_url: '', category_id: 'd2', has_variants: false, tax_rate: 11, sort_order: 7 },
];
