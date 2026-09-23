import React, { useEffect, useState } from 'react';
import { adminService, Debt } from '../services/api';

export default function HutangPage() {
  const [data, setData] = useState<Debt[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<Debt | null>(null);
  const [form, setForm] = useState({ creditor_name: '', amount: '', description: '', due_date: '', paid_amount: '0' });

  const load = async () => {
    try { const res = await adminService.listDebts(); setData(res.data.data); } catch {} finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  const openAdd = () => { setEditItem(null); setForm({ creditor_name: '', amount: '', description: '', due_date: '', paid_amount: '0' }); setShowForm(true); };
  const openEdit = (d: Debt) => { setEditItem(d); setForm({ creditor_name: d.creditor_name, amount: String(d.amount), description: d.description, due_date: d.due_date || '', paid_amount: String(d.paid_amount) }); setShowForm(true); };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload: any = { ...form, amount: parseFloat(form.amount) };
      if (!form.due_date) payload.due_date = null;
      if (editItem) {
        payload.paid_amount = parseFloat(form.paid_amount);
        await adminService.updateDebt(editItem.id, payload);
      } else {
        delete payload.paid_amount;
        await adminService.createDebt(payload);
      }
      setShowForm(false); load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus hutang ini?')) return;
    try { await adminService.deleteDebt(id); load(); } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Hutang</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Hutang</button>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Hutang' : 'Tambah Hutang'}</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <input placeholder="Nama Kreditur" value={form.creditor_name} onChange={e => setForm({ ...form, creditor_name: e.target.value })} className="border p-2 rounded" required />
            <input type="number" placeholder="Jumlah" value={form.amount} onChange={e => setForm({ ...form, amount: e.target.value })} className="border p-2 rounded" required />
            <input type="date" placeholder="Jatuh Tempo" value={form.due_date} onChange={e => setForm({ ...form, due_date: e.target.value })} className="border p-2 rounded" />
            {editItem && <input type="number" placeholder="Sudah Dibayar" value={form.paid_amount} onChange={e => setForm({ ...form, paid_amount: e.target.value })} className="border p-2 rounded" />}
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
                <th className="p-3 text-left">Kreditur</th>
                <th className="p-3 text-right">Jumlah</th>
                <th className="p-3 text-right">Dibayar</th>
                <th className="p-3 text-right">Sisa</th>
                <th className="p-3 text-left">Jatuh Tempo</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data?.map(d => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.creditor_name}</td>
                  <td className="p-3 text-right">Rp {d.amount.toLocaleString()}</td>
                  <td className="p-3 text-right">Rp {d.paid_amount.toLocaleString()}</td>
                  <td className="p-3 text-right font-bold">Rp {d.remaining_amount.toLocaleString()}</td>
                  <td className="p-3">{d.due_date || '-'}</td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs ${d.status === 'paid' ? 'bg-green-100 text-green-700' : d.status === 'partial' ? 'bg-yellow-100 text-yellow-700' : 'bg-red-100 text-red-700'}`}>
                      {d.status === 'paid' ? 'Lunas' : d.status === 'partial' ? 'Sebagian' : 'Belum Bayar'}
                    </span>
                  </td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => openEdit(d)} className="text-blue-500 hover:underline">Edit</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data?.length === 0 && <tr><td colSpan={7} className="p-4 text-center text-gray-400">Belum ada data hutang</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
