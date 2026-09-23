import React, { useEffect, useState } from 'react';
import { adminService, SalesCategory } from '../services/api';

export default function SalesCategoryPage() {
  const [data, setData] = useState<SalesCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<SalesCategory | null>(null);
  const [form, setForm] = useState({ name: '', description: '' });
  const [error, setError] = useState('');

  const load = async () => {
    setError('');
    try { const res = await adminService.listSalesCategories(); setData(res.data.data ?? []); }
    catch (err: any) { setError(err.response?.data?.error || err.message || 'Gagal memuat data'); }
    finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  const openAdd = () => { setEditItem(null); setForm({ name: '', description: '' }); setError(''); setShowForm(true); };
  const openEdit = (d: SalesCategory) => { setEditItem(d); setForm({ name: d.name, description: d.description }); setError(''); setShowForm(true); };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      if (editItem) { await adminService.updateSalesCategory(editItem.id, form); }
      else { await adminService.createSalesCategory(form); }
      setShowForm(false); load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.response?.data?.message || err.message || 'Gagal menyimpan data');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus kategori penjualan ini?')) return;
    try { await adminService.deleteSalesCategory(id); load(); } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Kategori Penjualan</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Kategori</button>
      </div>

      {error && <div className="mb-4 p-3 text-red-500 bg-red-50 border border-red-200 rounded">{error}</div>}

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Kategori Penjualan' : 'Tambah Kategori Penjualan'}</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Nama Kategori</label>
              <input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} className="border p-2 rounded w-full" required />
            </div>
            <div>
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
                <th className="p-3 text-left">Deskripsi</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(d => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.name}</td>
                  <td className="p-3">{d.description || '-'}</td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => openEdit(d)} className="text-blue-500 hover:underline">Edit</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && <tr><td colSpan={3} className="p-4 text-center text-gray-400">Belum ada data kategori penjualan</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}