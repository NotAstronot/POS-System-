import React, { useEffect, useState } from 'react';
import { adminService, BankTransfer, BankReconciliation } from '../services/api';

export default function TransferBankPage() {
  const [data, setData] = useState<BankTransfer[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<BankTransfer | null>(null);
  const [form, setForm] = useState({
    from_account_name: '',
    from_account_type: 'bank',
    to_account_name: '',
    to_account_type: 'bank',
    amount: '',
    transfer_date: '',
    notes: '',
  });

  const load = async () => {
    try {
      const res = await adminService.listBankTransfers();
      setData(res.data.data);
    } catch {} finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openAdd = () => {
    setEditItem(null);
    setForm({
      from_account_name: '',
      from_account_type: 'bank',
      to_account_name: '',
      to_account_type: 'bank',
      amount: '',
      transfer_date: '',
      notes: '',
    });
    setShowForm(true);
  };

  const openEdit = (item: BankTransfer) => {
    setEditItem(item);
    setForm({
      from_account_name: item.from_account_name,
      from_account_type: item.from_account_type,
      to_account_name: item.to_account_name,
      to_account_type: item.to_account_type,
      amount: String(item.amount),
      transfer_date: item.transfer_date,
      notes: item.notes,
    });
    setShowForm(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload = {
        ...form,
        amount: parseFloat(form.amount),
      };
      if (editItem) {
        await adminService.updateBankTransfer(editItem.id, payload);
      } else {
        await adminService.createBankTransfer(payload);
      }
      setShowForm(false);
      load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus transfer bank ini?')) return;
    try {
      await adminService.deleteBankTransfer(id);
      load();
    } catch {}
  };

  const totalTransfer = data.reduce((s, d) => s + d.amount, 0);

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'bank': return 'Bank';
      case 'kas': return 'Kas';
      default: return type;
    }
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Transfer Bank</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Transfer</button>
      </div>

      <div className="bg-blue-50 p-4 rounded mb-6">
        <div className="text-sm text-blue-700">Total Semua Transfer</div>
        <div className="text-xl font-bold text-blue-800">Rp {totalTransfer.toLocaleString()}</div>
        <p className="text-xs text-blue-600 mt-1">Pencatatan perpindahan dana antar rekening bank atau dari kasir ke bank</p>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Transfer' : 'Transfer Baru'}</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Rekening Asal</label>
              <input
                placeholder="Nama Rekening Asal"
                value={form.from_account_name}
                onChange={(e) => setForm({ ...form, from_account_name: e.target.value })}
                className="border p-2 rounded w-full"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Tipe Asal</label>
              <select
                value={form.from_account_type}
                onChange={(e) => setForm({ ...form, from_account_type: e.target.value })}
                className="border p-2 rounded w-full"
              >
                <option value="bank">Bank</option>
                <option value="kas">Kas (Kasir)</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Rekening Tujuan</label>
              <input
                placeholder="Nama Rekening Tujuan"
                value={form.to_account_name}
                onChange={(e) => setForm({ ...form, to_account_name: e.target.value })}
                className="border p-2 rounded w-full"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Tipe Tujuan</label>
              <select
                value={form.to_account_type}
                onChange={(e) => setForm({ ...form, to_account_type: e.target.value })}
                className="border p-2 rounded w-full"
              >
                <option value="bank">Bank</option>
                <option value="kas">Kas (Kasir)</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Jumlah (Rp)</label>
              <input
                type="number"
                placeholder="Jumlah"
                value={form.amount}
                onChange={(e) => setForm({ ...form, amount: e.target.value })}
                className="border p-2 rounded w-full"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Tanggal Transfer</label>
              <input
                type="date"
                value={form.transfer_date}
                onChange={(e) => setForm({ ...form, transfer_date: e.target.value })}
                className="border p-2 rounded w-full"
              />
            </div>
            <div className="col-span-2">
              <label className="block text-sm font-medium mb-1">Catatan</label>
              <input
                placeholder="Catatan transfer"
                value={form.notes}
                onChange={(e) => setForm({ ...form, notes: e.target.value })}
                className="border p-2 rounded w-full"
              />
            </div>
            {(form.from_account_type === 'kas' || form.to_account_type === 'kas') && (
              <div className="col-span-2 bg-yellow-50 border border-yellow-200 p-3 rounded text-sm text-yellow-800">
                Transfer ini akan otomatis tercatat di menu <strong>Kas</strong> sebagai transaksi masuk/keluar.
              </div>
            )}
            <div className="col-span-2 flex gap-2">
              <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Simpan</button>
              <button type="button" onClick={() => setShowForm(false)} className="bg-gray-400 text-white px-4 py-2 rounded">Batal</button>
            </div>
          </form>
        </div>
      )}

      {loading ? (
        <p>Loading...</p>
      ) : (
        <div className="bg-white rounded shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-3 text-left">No. Transfer</th>
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-left">Dari</th>
                <th className="p-3 text-left">Ke</th>
                <th className="p-3 text-right">Jumlah</th>
                <th className="p-3 text-left">Catatan</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map((d) => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.transfer_number}</td>
                  <td className="p-3">{d.transfer_date ? new Date(d.transfer_date).toLocaleDateString('id-ID') : '-'}</td>
                  <td className="p-3">
                    <span>{d.from_account_name}</span>
                    <span className="ml-1 text-xs text-gray-500">({getTypeLabel(d.from_account_type)})</span>
                  </td>
                  <td className="p-3">
                    <span>{d.to_account_name}</span>
                    <span className="ml-1 text-xs text-gray-500">({getTypeLabel(d.to_account_type)})</span>
                  </td>
                  <td className="p-3 text-right font-medium">Rp {d.amount.toLocaleString()}</td>
                  <td className="p-3 text-sm text-gray-600">{d.notes || '-'}</td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs ${
                      d.status === 'completed' ? 'bg-green-100 text-green-700' :
                      d.status === 'pending' ? 'bg-yellow-100 text-yellow-700' :
                      'bg-red-100 text-red-700'
                    }`}>
                      {d.status === 'completed' ? 'Selesai' : d.status === 'pending' ? 'Pending' : 'Dibatalkan'}
                    </span>
                  </td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => openEdit(d)} className="text-blue-500 hover:underline">Edit</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && (
                <tr>
                  <td colSpan={8} className="p-4 text-center text-gray-400">Belum ada transfer bank</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
