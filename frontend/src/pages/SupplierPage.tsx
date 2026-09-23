import React, { useEffect, useState } from 'react';
import { adminService, Supplier } from '../services/api';

export default function SupplierPage() {
  const [data, setData] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<Supplier | null>(null);
  const [form, setForm] = useState({ name: '', contact_name: '', phone: '', email: '', address: '' });
  const [error, setError] = useState('');

  const load = async () => {
    try { const res = await adminService.listSuppliers(); setData(res.data.data ?? []); } catch {} finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  const openAdd = () => { setEditItem(null); setForm({ name: '', contact_name: '', phone: '', email: '', address: '' }); setError(''); setShowForm(true); };
  const openEdit = (d: Supplier) => { setEditItem(d); setForm({ name: d.name, contact_name: d.contact_name, phone: d.phone, email: d.email, address: d.address }); setError(''); setShowForm(true); };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      if (editItem) {
        await adminService.updateSupplier(editItem.id, form);
      } else {
        await adminService.createSupplier(form);
      }
      setShowForm(false); load();
    } catch (err: any) {
      const msg = err.response?.data?.error 
        || err.response?.data?.message 
        || err.message 
        || 'Gagal menyimpan data';
      setError(msg);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus supplier ini?')) return;
    try { await adminService.deleteSupplier(id); load(); } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Supplier</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Supplier</button>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Supplier' : 'Tambah Supplier'}</h2>
          {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <input placeholder="Nama Supplier *" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} className="border p-2 rounded" required />
            <input placeholder="Nama Kontak" value={form.contact_name} onChange={e => setForm({ ...form, contact_name: e.target.value })} className="border p-2 rounded" />
            <input placeholder="Telepon" value={form.phone} onChange={e => setForm({ ...form, phone: e.target.value })} className="border p-2 rounded" />
            <input type="email" placeholder="Email" value={form.email} onChange={e => setForm({ ...form, email: e.target.value })} className="border p-2 rounded" />
            <input placeholder="Alamat" value={form.address} onChange={e => setForm({ ...form, address: e.target.value })} className="border p-2 rounded col-span-2" />
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
                <th className="p-3 text-left">Kontak</th>
                <th className="p-3 text-left">Telepon</th>
                <th className="p-3 text-left">Email</th>
                <th className="p-3 text-left">Alamat</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(d => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.name}</td>
                  <td className="p-3">{d.contact_name}</td>
                  <td className="p-3">{d.phone}</td>
                  <td className="p-3">{d.email}</td>
                  <td className="p-3">{d.address}</td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs ${d.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                      {d.is_active ? 'Aktif' : 'Nonaktif'}
                    </span>
                  </td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => openEdit(d)} className="text-blue-500 hover:underline">Edit</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && <tr><td colSpan={7} className="p-4 text-center text-gray-400">Belum ada data supplier</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
