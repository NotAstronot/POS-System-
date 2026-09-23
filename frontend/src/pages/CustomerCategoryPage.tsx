import React, { useEffect, useState } from 'react';
import { adminService, CustomerCategory } from '../services/api';

export default function CustomerCategoryPage() {
  const [data, setData] = useState<CustomerCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<CustomerCategory | null>(null);
  const [form, setForm] = useState({ name: '', description: '', price_scheme: 'retail' });
  const [error, setError] = useState('');

  const load = async () => {
    setError('');
    try { const res = await adminService.listCustomerCategories(); setData(res.data.data ?? []); }
    catch (err: any) { setError(err.response?.data?.error || err.message || 'Gagal memuat data'); }
    finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  const openAdd = () => { setEditItem(null); setForm({ name: '', description: '', price_scheme: 'retail' }); setError(''); setShowForm(true); };
  const openEdit = (d: CustomerCategory) => { setEditItem(d); setForm({ name: d.name, description: d.description, price_scheme: d.price_scheme }); setError(''); setShowForm(true); };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      if (editItem) { await adminService.updateCustomerCategory(editItem.id, form); }
      else { await adminService.createCustomerCategory(form); }
      setShowForm(false); load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.response?.data?.message || err.message || 'Gagal menyimpan data');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus kategori pelanggan ini?')) return;
    try { await adminService.deleteCustomerCategory(id); load(); } catch {}
  };

  const schemeLabel = (s: string) =>
    s === 'wholesale' ? 'Grosir' : s === 'reseller' ? 'Reseller' : 'Retail';

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Kategori Pelanggan</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Kategori</button>
      </div>

      {error && <div className="mb-4 p-3 text-red-500 bg-red-50 border border-red-200 rounded">{error}</div>}

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Kategori Pelanggan' : 'Tambah Kategori Pelanggan'}</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Nama Kategori</label>
              <input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} className="border p-2 rounded w-full" required />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Skema Harga</label>
              <select value={form.price_scheme} onChange={e => setForm({ ...form, price_scheme: e.target.value })} className="border p-2 rounded w-full">
                <option value="retail">Retail</option>
                <option value="wholesale">Grosir</option>
                <option value="reseller">Reseller</option>
              </select>
            </div>
            <div className="col-span-2">
              <label className="block text-sm font-medium mb-1">Deskripsi</label>
              <input placeholder="Deskripsi" value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} className="border p-2 rounded w-full" />
            </div>
            <div className="col-span-2 flex gap-2">
              <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Simpan</button>
              <button type="button" onClick={() => setShowForm(false)} className="bg-gray-400 text-white px-4 py-2 rounded">Batal</button>
            </div>
          </form>
        </div>
      )}

      {loading ? <p>Loading...</p> : (
        <div className="bg-white rounded shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-3 text-left">Nama</th>
                <th className="p-3 text-left">Skema Harga</th>
                <th className="p-3 text-left">Deskripsi</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(d => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.name}</td>
                  <td className="p-3">
                    <span className="px-2 py-1 rounded text-xs bg-blue-100 text-blue-700">{schemeLabel(d.price_scheme)}</span>
                  </td>
                  <td className="p-3">{d.description || '-'}</td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => openEdit(d)} className="text-blue-500 hover:underline">Edit</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && <tr><td colSpan={4} className="p-4 text-center text-gray-400">Belum ada data kategori pelanggan</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}