import React, { useEffect, useState } from 'react';
import { adminService, Customer, CustomerCategory } from '../services/api';

export default function CustomerPage() {
  const [data, setData] = useState<Customer[]>([]);
  const [categories, setCategories] = useState<CustomerCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<Customer | null>(null);
  const [form, setForm] = useState({ name: '', phone: '', email: '', address: '', customer_category_id: '' });
  const [error, setError] = useState('');

  const load = async () => {
    setError('');
    try {
      const [res, resCat] = await Promise.all([adminService.listCustomers(), adminService.listCustomerCategories()]);
      setData(res.data.data ?? []);
      setCategories(resCat.data.data ?? []);
    } catch (err: any) { setError(err.response?.data?.error || err.message || 'Gagal memuat data'); }
    finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  const openAdd = () => { setEditItem(null); setForm({ name: '', phone: '', email: '', address: '', customer_category_id: '' }); setError(''); setShowForm(true); };
  const openEdit = (d: Customer) => { setEditItem(d); setForm({ name: d.name, phone: d.phone, email: d.email, address: d.address, customer_category_id: d.customer_category_id ? String(d.customer_category_id) : '' }); setError(''); setShowForm(true); };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    const payload = {
      name: form.name,
      phone: form.phone,
      email: form.email,
      address: form.address,
      customer_category_id: form.customer_category_id ? parseInt(form.customer_category_id) : null,
    };
    try {
      if (editItem) { await adminService.updateCustomer(editItem.id, payload); }
      else { await adminService.createCustomer(payload); }
      setShowForm(false); load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.response?.data?.message || err.message || 'Gagal menyimpan data');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus pelanggan ini?')) return;
    try { await adminService.deleteCustomer(id); load(); } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Pelanggan</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Pelanggan</button>
      </div>

      {error && <div className="mb-4 p-3 text-red-500 bg-red-50 border border-red-200 rounded">{error}</div>}

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Pelanggan' : 'Tambah Pelanggan'}</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Nama</label>
              <input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} className="border p-2 rounded w-full" required />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">No HP</label>
              <input value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })} className="border p-2 rounded w-full" />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Email</label>
              <input type="email" value={form.email} onChange={e => setForm({ ...form, email: e.target.value })} className="border p-2 rounded w-full" />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Kategori</label>
              <select value={form.customer_category_id} onChange={e => setForm({ ...form, customer_category_id: e.target.value })} className="border p-2 rounded w-full">
                <option value="">-- Tanpa Kategori --</option>
                {categories.filter(c => c.is_active).map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
            </div>
            <div className="col-span-2">
              <label className="block text-sm font-medium mb-1">Alamat</label>
              <input placeholder="Alamat" value={form.address} onChange={e => setForm({ ...form, address: e.target.value })} className="border p-2 rounded w-full" />
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
                <th className="p-3 text-left">No HP</th>
                <th className="p-3 text-left">Email</th>
                <th className="p-3 text-left">Kategori</th>
                <th className="p-3 text-left">Alamat</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(d => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.name}</td>
                  <td className="p-3">{d.phone || '-'}</td>
                  <td className="p-3">{d.email || '-'}</td>
                  <td className="p-3">{d.customer_category_name || '-'}</td>
                  <td className="p-3">{d.address || '-'}</td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => openEdit(d)} className="text-blue-500 hover:underline">Edit</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && <tr><td colSpan={6} className="p-4 text-center text-gray-400">Belum ada data pelanggan</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}