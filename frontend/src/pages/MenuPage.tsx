import React, { useState, useEffect } from 'react';
import { Plus, Pencil, Trash2, X, Search, Package, AlertTriangle } from 'lucide-react';
import { productService, categoryService, type Product, type Category } from '../services/api';
import { useStore } from '../store/useStore';
import toast from 'react-hot-toast';

interface ProductForm {
  code: string;
  name: string;
  description: string;
  base_price: number;
  unit: string;
  category_id: string;
  has_variants: boolean;
  tax_rate: number;
  sort_order: number;
  stock: number;
}

const emptyForm: ProductForm = {
  code: '', name: '', description: '', base_price: 0,
  unit: 'pcs', category_id: '', has_variants: false,
  tax_rate: 11, sort_order: 0, stock: 0,
};

export function MenuPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<string | null>(null);
  const [form, setForm] = useState<ProductForm>(emptyForm);
  const [searchQ, setSearchQ] = useState('');
  const user = useStore((s) => s.user);

  const loadData = async () => {
    try {
      const [catRes] = await Promise.all([
        categoryService.list(user?.outletId || ''),
      ]);
      setCategories(catRes.data.data || []);
    } catch { setCategories([]); }
    loadProducts();
  };

  const loadProducts = async () => {
    try {
      const res = await productService.list();
      setProducts(res.data.data || []);
    } catch { setProducts([]); }
    setLoading(false);
  };

  useEffect(() => { loadData(); }, []);

  const filtered = products.filter((p) =>
    !searchQ || p.name.toLowerCase().includes(searchQ.toLowerCase()) ||
    p.code.toLowerCase().includes(searchQ.toLowerCase())
  );

  const openCreate = () => {
    setEditId(null);
    setForm(emptyForm);
    setShowForm(true);
  };

  const openEdit = (p: Product) => {
    setEditId(p.id);
    setForm({
      code: p.code, name: p.name, description: p.description || '',
      base_price: p.base_price, unit: p.unit,
      category_id: p.category_id, has_variants: p.has_variants,
      tax_rate: p.tax_rate, sort_order: p.sort_order, stock: 0,
    });
    setShowForm(true);
  };

  const handleSave = async () => {
    if (!form.name) { toast.error('Nama produk wajib diisi'); return; }
    try {
      const payload = { ...form };
      if (editId) {
        await productService.update(editId, payload);
        toast.success('Produk diupdate!');
      } else {
        await productService.create(payload);
        toast.success('Produk ditambahkan!');
      }
      setShowForm(false);
      loadProducts();
    } catch (err: any) {
      toast.error(err?.response?.data?.error || 'Gagal menyimpan');
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Hapus produk ini?')) return;
    try {
      await productService.delete(id);
      toast.success('Produk dihapus!');
      loadProducts();
    } catch { toast.error('Gagal menghapus'); }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Manajemen Menu</h1>
          <p className="text-gray-400 text-sm mt-1">Kelola daftar produk</p>
        </div>
        <button onClick={openCreate} className="pos-btn-primary flex items-center gap-2">
          <Plus className="w-4 h-4" /> Tambah Produk
        </button>
      </div>

      <div className="relative max-w-md">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
        <input
          value={searchQ} onChange={(e) => setSearchQ(e.target.value)}
          placeholder="Cari produk..." className="pos-input pl-10"
        />
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="bg-gray-800/50 rounded-xl h-32 animate-pulse" />
          ))}
        </div>
      ) : filtered.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          <Package className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>Tidak ada produk</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filtered.map((p) => (
            <div key={p.id} className="bg-gray-900 border border-gray-800 rounded-xl p-4 hover:border-gray-700 transition-colors">
              <div className="flex items-start justify-between gap-3">
                <div className="flex-1 min-w-0">
                  <p className="text-xs text-gray-500 font-mono">{p.code}</p>
                  <p className="text-white font-medium truncate">{p.name}</p>
                  <p className="text-blue-400 font-bold mt-1">
                    Rp {p.base_price.toLocaleString()}
                  </p>
                </div>
                <div className="flex gap-1 shrink-0">
                  <button onClick={() => openEdit(p)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400 hover:text-blue-400">
                    <Pencil className="w-4 h-4" />
                  </button>
                  <button onClick={() => handleDelete(p.id)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400 hover:text-red-400">
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
              <div className="flex items-center gap-3 mt-3 text-xs text-gray-500">
                <span>{p.unit}</span>
                {p.has_variants && <span className="pos-badge-warning">Varian</span>}
                <span className="ml-auto">Pajak {p.tax_rate}%</span>
              </div>
            </div>
          ))}
        </div>
      )}

      {showForm && (
        <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
          <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto shadow-2xl">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
              <h2 className="font-bold text-white">{editId ? 'Edit Produk' : 'Tambah Produk'}</h2>
              <button onClick={() => setShowForm(false)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400">
                <X className="w-5 h-5" />
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm text-gray-300 mb-1">Kode</label>
                  <input value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} className="pos-input" placeholder="BRGR-001" />
                </div>
                <div>
                  <label className="block text-sm text-gray-300 mb-1">Unit</label>
                  <select value={form.unit} onChange={(e) => setForm({ ...form, unit: e.target.value })} className="pos-input">
                    <option value="pcs">pcs</option>
                    <option value="gelas">gelas</option>
                    <option value="botol">botol</option>
                    <option value="porsi">porsi</option>
                    <option value="pack">pack</option>
                  </select>
                </div>
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Nama Produk *</label>
                <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="pos-input" placeholder="Nama produk" />
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Deskripsi</label>
                <textarea value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} className="pos-input" rows={2} />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm text-gray-300 mb-1">Harga</label>
                  <input type="number" value={form.base_price} onChange={(e) => setForm({ ...form, base_price: Number(e.target.value) })} className="pos-input" />
                </div>
                <div>
                  <label className="block text-sm text-gray-300 mb-1">Kategori</label>
                  <select value={form.category_id} onChange={(e) => setForm({ ...form, category_id: e.target.value })} className="pos-input">
                    <option value="">Pilih kategori</option>
                    {categories.map((c: any) => (
                      <option key={c.id} value={c.id}>{c.name}</option>
                    ))}
                  </select>
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm text-gray-300 mb-1">Pajak (%)</label>
                  <input type="number" value={form.tax_rate} onChange={(e) => setForm({ ...form, tax_rate: Number(e.target.value) })} className="pos-input" />
                </div>
                <div>
                  <label className="block text-sm text-gray-300 mb-1">Urutan</label>
                  <input type="number" value={form.sort_order} onChange={(e) => setForm({ ...form, sort_order: Number(e.target.value) })} className="pos-input" />
                </div>
              </div>
              <label className="flex items-center gap-2 text-sm text-gray-300">
                <input type="checkbox" checked={form.has_variants} onChange={(e) => setForm({ ...form, has_variants: e.target.checked })} className="rounded bg-gray-800 border-gray-600" />
                Punya varian (ukuran/toping)
              </label>
              <button onClick={handleSave} className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition-colors">
                {editId ? 'Simpan Perubahan' : 'Tambah Produk'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
