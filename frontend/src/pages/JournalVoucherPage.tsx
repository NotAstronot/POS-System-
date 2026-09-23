import React, { useEffect, useState } from 'react';
import { adminService, ChartOfAccount, JournalEntry } from '../services/api';

interface FormItem {
  account_id: number;
  debit: string;
  credit: string;
}

const referenceOptions = [
  'Umum',
  'Kas',
  'Hutang',
  'Piutang',
  'Transfer Bank',
  'Rekonsiliasi Bank',
  'Pembelian',
  'Penjualan',
  'Depresiasi Aset',
  'Penutupan Buku',
];

export default function JournalVoucherPage() {
  const [data, setData] = useState<JournalEntry[]>([]);
  const [accounts, setAccounts] = useState<ChartOfAccount[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({
    entry_date: '',
    description: '',
    reference: 'Umum',
  });
  const [items, setItems] = useState<FormItem[]>([
    { account_id: 0, debit: '', credit: '' },
    { account_id: 0, debit: '', credit: '' },
  ]);

  const load = async () => {
    try {
      const [entriesRes, accountsRes] = await Promise.all([
        adminService.listJournalEntries(),
        adminService.listChartOfAccounts(),
      ]);
      setData(entriesRes.data.data);
      setAccounts(accountsRes.data.data.filter((a) => a.is_active));
    } catch {} finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openAdd = () => {
    setForm({ entry_date: '', description: '', reference: 'Umum' });
    setItems([
      { account_id: 0, debit: '', credit: '' },
      { account_id: 0, debit: '', credit: '' },
    ]);
    setShowForm(true);
  };

  const updateItem = (index: number, field: keyof FormItem, value: string | number) => {
    setItems((prev) => prev.map((it, i) => (i === index ? { ...it, [field]: value } : it)));
  };

  const addItem = () => setItems((prev) => [...prev, { account_id: 0, debit: '', credit: '' }]);

  const removeItem = (index: number) => {
    if (items.length <= 2) return;
    setItems((prev) => prev.filter((_, i) => i !== index));
  };

  const totalDebit = items.reduce((s, it) => s + (parseFloat(it.debit) || 0), 0);
  const totalCredit = items.reduce((s, it) => s + (parseFloat(it.credit) || 0), 0);
  const diff = Math.abs(totalDebit - totalCredit);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload = {
        ...form,
        items: items.map((it) => ({
          account_id: Number(it.account_id),
          debit: parseFloat(it.debit) || 0,
          credit: parseFloat(it.credit) || 0,
        })),
      };
      await adminService.createJournalEntry(payload);
      setShowForm(false);
      load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus jurnal ini?')) return;
    try {
      await adminService.deleteJournalEntry(id);
      load();
    } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold">Jurnal Umum</h1>
          <p className="text-sm text-gray-500">Input jurnal manual untuk penyesuaian khusus</p>
        </div>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Jurnal Baru</button>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">Jurnal Baru</h2>
          <form onSubmit={handleSubmit}>
            <div className="grid grid-cols-3 gap-4 mb-4">
              <div>
                <label className="block text-sm font-medium mb-1">Tanggal</label>
                <input
                  type="date"
                  value={form.entry_date}
                  onChange={(e) => setForm({ ...form, entry_date: e.target.value })}
                  className="border p-2 rounded w-full"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Referensi Modul</label>
                <select
                  value={form.reference}
                  onChange={(e) => setForm({ ...form, reference: e.target.value })}
                  className="border p-2 rounded w-full"
                >
                  {referenceOptions.map((r) => (
                    <option key={r} value={r}>{r}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Deskripsi</label>
                <input
                  placeholder="Keterangan jurnal"
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  className="border p-2 rounded w-full"
                  required
                />
              </div>
            </div>
            <table className="w-full mb-2">
              <thead className="bg-gray-100">
                <tr>
                  <th className="p-2 text-left w-1/2">Akun (dari Akun Perkiraan)</th>
                  <th className="p-2 text-right">Debit</th>
                  <th className="p-2 text-right">Kredit</th>
                  <th className="p-2 w-16"></th>
                </tr>
              </thead>
              <tbody>
                {items.map((it, i) => (
                  <tr key={i} className="border-t">
                    <td className="p-2">
                      <select
                        value={it.account_id}
                        onChange={(e) => updateItem(i, 'account_id', e.target.value)}
                        className="border p-2 rounded w-full"
                        required
                      >
                        <option value={0}>-- Pilih Akun --</option>
                        {accounts.map((a) => (
                          <option key={a.id} value={a.id}>{a.code} - {a.name}</option>
                        ))}
                      </select>
                    </td>
                    <td className="p-2">
                      <input
                        type="number"
                        placeholder="0"
                        value={it.debit}
                        onChange={(e) => updateItem(i, 'debit', e.target.value)}
                        className="border p-2 rounded w-full text-right"
                      />
                    </td>
                    <td className="p-2">
                      <input
                        type="number"
                        placeholder="0"
                        value={it.credit}
                        onChange={(e) => updateItem(i, 'credit', e.target.value)}
                        className="border p-2 rounded w-full text-right"
                      />
                    </td>
                    <td className="p-2 text-center">
                      <button type="button" onClick={() => removeItem(i)} className="text-red-500 hover:underline text-sm">Hapus</button>
                    </td>
                  </tr>
                ))}
              </tbody>
              <tfoot className="bg-gray-50">
                <tr>
                  <td className="p-2 font-semibold text-right">Total</td>
                  <td className="p-2 font-semibold text-right">Rp {totalDebit.toLocaleString()}</td>
                  <td className="p-2 font-semibold text-right">Rp {totalCredit.toLocaleString()}</td>
                  <td></td>
                </tr>
              </tfoot>
            </table>
            <div className="flex items-center justify-between mb-4">
              <button type="button" onClick={addItem} className="text-blue-600 text-sm font-medium hover:underline">+ Tambah Baris</button>
              {diff > 0 && (
                <span className="text-sm text-red-600 bg-red-50 px-3 py-1 rounded">
                  Jurnal belum seimbang (selisih Rp {diff.toLocaleString()})
                </span>
              )}
            </div>
            <div className="flex gap-2">
              <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Simpan Jurnal</button>
              <button type="button" onClick={() => setShowForm(false)} className="bg-gray-400 text-white px-4 py-2 rounded">Batal</button>
            </div>
          </form>
        </div>
      )}

      {loading ? (
        <p>Loading...</p>
      ) : (
        <div className="space-y-4">
          {data.map((e) => (
            <div key={e.id} className="bg-white rounded shadow overflow-hidden">
              <div className="flex justify-between items-center bg-gray-100 px-4 py-3">
                <div className="flex items-center gap-3">
                  <span className="font-bold">{e.entry_number}</span>
                  <span className="text-sm text-gray-600">
                    {e.entry_date ? new Date(e.entry_date).toLocaleDateString('id-ID') : '-'}
                  </span>
                  <span className="text-xs px-2 py-1 bg-blue-100 text-blue-700 rounded">{e.reference || 'Umum'}</span>
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-xs px-2 py-1 bg-green-100 text-green-700 rounded">{e.status}</span>
                  <button onClick={() => handleDelete(e.id)} className="text-red-500 hover:underline text-sm">Hapus</button>
                </div>
              </div>
              <div className="px-4 py-2 text-sm text-gray-600">{e.description || '-'}</div>
              <table className="w-full">
                <thead>
                  <tr className="text-left text-xs text-gray-500 border-t">
                    <th className="p-2 pl-4">Akun</th>
                    <th className="p-2 text-right">Debit</th>
                    <th className="p-2 pr-4 text-right">Kredit</th>
                  </tr>
                </thead>
                <tbody>
                  {e.items?.map((it) => (
                    <tr key={it.id} className="border-t">
                      <td className="p-2 pl-4">{it.account_code} - {it.account_name}</td>
                      <td className="p-2 text-right">{it.debit > 0 ? `Rp ${it.debit.toLocaleString()}` : '-'}</td>
                      <td className="p-2 pr-4 text-right">{it.credit > 0 ? `Rp ${it.credit.toLocaleString()}` : '-'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ))}
          {data.length === 0 && (
            <div className="bg-white rounded shadow p-8 text-center text-gray-400">
              Belum ada jurnal umum. Klik "+ Jurnal Baru" untuk membuat penyesuaian.
            </div>
          )}
        </div>
      )}
    </div>
  );
}
