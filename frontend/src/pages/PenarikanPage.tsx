import React, { useEffect, useState } from 'react';
import { adminService, Withdrawal } from '../services/api';

export default function PenarikanPage() {
  const [data, setData] = useState<Withdrawal[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<Withdrawal | null>(null);
  const [form, setForm] = useState({ amount: '', description: '', method: 'cash' });

  const load = async () => {
    try { const res = await adminService.listWithdrawals(); setData(res.data.data); } catch {} finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  const openAdd = () => { setEditItem(null); setForm({ amount: '', description: '', method: 'cash' }); setShowForm(true); };
  const openEdit = (w: Withdrawal) => { setEditItem(w); setForm({ amount: String(w.amount), description: w.description, method: w.method }); setShowForm(true); };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload = { ...form, amount: parseFloat(form.amount) };
      if (editItem) {
        await adminService.updateWithdrawal(editItem.id, payload);
      } else {
        await adminService.createWithdrawal(payload);
      }
      setShowForm(false); load();
    } catch {}
  };

  const handleApprove = async (id: number) => {
    try { await adminService.approveWithdrawal(id); load(); } catch {}
  };
  const handleReject = async (id: number) => {
    if (!confirm('Tolak penarikan ini?')) return;
    try { await adminService.rejectWithdrawal(id); load(); } catch {}
  };
  const handleDelete = async (id: number) => {
    if (!confirm('Hapus penarikan ini?')) return;
    try { await adminService.deleteWithdrawal(id); load(); } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Penarikan Dana</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Penarikan</button>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Penarikan' : 'Ajukan Penarikan'}</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <input type="number" placeholder="Jumlah" value={form.amount} onChange={e => setForm({ ...form, amount: e.target.value })} className="border p-2 rounded" required />
            <select value={form.method} onChange={e => setForm({ ...form, method: e.target.value })} className="border p-2 rounded">
              <option value="cash">Cash</option>
              <option value="transfer">Transfer</option>
            </select>
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
                <th className="p-3 text-right">Jumlah</th>
                <th className="p-3 text-left">Metode</th>
                <th className="p-3 text-left">Deskripsi</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data?.map(w => (
                <tr key={w.id} className="border-t">
                  <td className="p-3">{new Date(w.created_at).toLocaleDateString('id-ID')}</td>
                  <td className="p-3 text-right">Rp {w.amount.toLocaleString()}</td>
                  <td className="p-3 capitalize">{w.method}</td>
                  <td className="p-3">{w.description || '-'}</td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs ${w.status === 'approved' ? 'bg-green-100 text-green-700' : w.status === 'rejected' ? 'bg-red-100 text-red-700' : 'bg-yellow-100 text-yellow-700'}`}>
                      {w.status === 'approved' ? 'Disetujui' : w.status === 'rejected' ? 'Ditolak' : 'Pending'}
                    </span>
                  </td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    {w.status === 'pending' && (
                      <>
                        <button onClick={() => handleApprove(w.id)} className="text-green-600 hover:underline">Setuju</button>
                        <button onClick={() => handleReject(w.id)} className="text-red-500 hover:underline">Tolak</button>
                      </>
                    )}
                    <button onClick={() => openEdit(w)} className="text-blue-500 hover:underline">Edit</button>
                    <button onClick={() => handleDelete(w.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data?.length === 0 && <tr><td colSpan={6} className="p-4 text-center text-gray-400">Belum ada data penarikan</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
