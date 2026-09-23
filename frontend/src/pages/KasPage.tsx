import React, { useEffect, useState } from 'react';
import { adminService, CashTransaction } from '../services/api';

export default function KasPage() {
  const [data, setData] = useState<CashTransaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ type: 'in', amount: '', description: '', category: '', reference: '' });

  const load = async () => {
    try {
      const res = await adminService.listCash();
      setData(res.data.data);
    } catch {} finally { setLoading(false); }
  };

  useEffect(() => { load(); }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await adminService.createCash({ ...form, amount: parseFloat(form.amount) });
      setShowForm(false);
      setForm({ type: 'in', amount: '', description: '', category: '', reference: '' });
      load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus data kas ini?')) return;
    try { await adminService.deleteCash(id); load(); } catch {}
  };

  const totalIn = data?.filter(d => d.type === 'in')?.reduce((s, d) => s + d.amount, 0) || 0;
  const totalOut = data?.filter(d => d.type === 'out')?.reduce((s, d) => s + d.amount, 0) || 0;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Kas</h1>
        <button onClick={() => setShowForm(true)} className="bg-blue-600 text-white px-4 py-2 rounded">+ Transaksi</button>
      </div>

      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-green-100 p-4 rounded">
          <div className="text-sm text-green-700">Total Pemasukan</div>
          <div className="text-xl font-bold text-green-800">Rp {totalIn.toLocaleString()}</div>
        </div>
        <div className="bg-red-100 p-4 rounded">
          <div className="text-sm text-red-700">Total Pengeluaran</div>
          <div className="text-xl font-bold text-red-800">Rp {totalOut.toLocaleString()}</div>
        </div>
        <div className="bg-blue-100 p-4 rounded">
          <div className="text-sm text-blue-700">Saldo</div>
          <div className="text-xl font-bold text-blue-800">Rp {(totalIn - totalOut).toLocaleString()}</div>
        </div>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">Transaksi Kas</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <select value={form.type} onChange={e => setForm({ ...form, type: e.target.value })} className="border p-2 rounded">
              <option value="in">Pemasukan</option>
              <option value="out">Pengeluaran</option>
            </select>
            <input type="number" placeholder="Jumlah" value={form.amount} onChange={e => setForm({ ...form, amount: e.target.value })} className="border p-2 rounded" required />
            <input placeholder="Kategori" value={form.category} onChange={e => setForm({ ...form, category: e.target.value })} className="border p-2 rounded" />
            <input placeholder="Referensi" value={form.reference} onChange={e => setForm({ ...form, reference: e.target.value })} className="border p-2 rounded" />
            <input placeholder="Deskripsi" value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} className="border p-2 rounded col-span-2" />
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
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-left">Tipe</th>
                <th className="p-3 text-left">Kategori</th>
                <th className="p-3 text-left">Deskripsi</th>
                <th className="p-3 text-left">Referensi</th>
                <th className="p-3 text-right">Jumlah</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data?.map(d => (
                <tr key={d.id} className="border-t">
                  <td className="p-3">{new Date(d.created_at).toLocaleDateString('id-ID')}</td>
                  <td className="p-3"><span className={d.type === 'in' ? 'text-green-600' : 'text-red-600'}>{d.type === 'in' ? 'Masuk' : 'Keluar'}</span></td>
                  <td className="p-3">{d.category || '-'}</td>
                  <td className="p-3">{d.description || '-'}</td>
                  <td className="p-3">{d.reference || '-'}</td>
                  <td className="p-3 text-right">Rp {d.amount.toLocaleString()}</td>
                  <td className="p-3 text-center">
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data?.length === 0 && <tr><td colSpan={7} className="p-4 text-center text-gray-400">Belum ada transaksi</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
