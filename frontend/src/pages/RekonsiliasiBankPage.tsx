import React, { useEffect, useState } from 'react';
import { adminService, BankReconciliation, BankReconciliationItem } from '../services/api';

export default function RekonsiliasiBankPage() {
  const [data, setData] = useState<BankReconciliation[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [detailItem, setDetailItem] = useState<BankReconciliation | null>(null);
  const [form, setForm] = useState({
    account_name: '',
    account_type: 'bank',
    period_start: '',
    period_end: '',
    book_balance: '',
    statement_balance: '',
    notes: '',
  });
  const [items, setItems] = useState<BankReconciliationItem[]>([]);

  const load = async () => {
    try {
      const res = await adminService.listBankReconciliations();
      setData(res.data.data);
    } catch {} finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openAdd = () => {
    setDetailItem(null);
    setForm({
      account_name: '',
      account_type: 'bank',
      period_start: '',
      period_end: '',
      book_balance: '',
      statement_balance: '',
      notes: '',
    });
    setItems([]);
    setShowForm(true);
  };

  const openDetail = (item: BankReconciliation) => {
    setDetailItem(item);
    setShowForm(false);
  };

  const addItem = () => {
    setItems([
      ...items,
      { description: '', transaction_date: '', amount: 0, type: 'debit', is_matched: false },
    ]);
  };

  const updateItem = (index: number, field: keyof BankReconciliationItem, value: any) => {
    const newItems = [...items];
    (newItems[index] as any)[field] = value;
    setItems(newItems);
  };

  const removeItem = (index: number) => {
    setItems(items.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const bookBal = parseFloat(form.book_balance) || 0;
      const stmtBal = parseFloat(form.statement_balance) || 0;
      await adminService.createBankReconciliation({
        ...form,
        book_balance: bookBal,
        statement_balance: stmtBal,
        items: items.map((item) => ({
          ...item,
          amount: parseFloat(String(item.amount)) || 0,
        })),
      });
      setShowForm(false);
      load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus rekonsiliasi ini?')) return;
    try {
      await adminService.deleteBankReconciliation(id);
      load();
    } catch {}
  };

  const getDiff = (book: number, stmt: number) => stmt - book;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Rekonsiliasi Bank</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Rekonsiliasi</button>
      </div>

      <div className="bg-purple-50 p-4 rounded mb-6">
        <p className="text-sm text-purple-700">Menyelaraskan catatan kas/bank di Accurate dengan mutasi bank riil</p>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">Rekonsiliasi Baru</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">Nama Rekening</label>
                <input
                  placeholder="Nama Rekening"
                  value={form.account_name}
                  onChange={(e) => setForm({ ...form, account_name: e.target.value })}
                  className="border p-2 rounded w-full"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Tipe Rekening</label>
                <select
                  value={form.account_type}
                  onChange={(e) => setForm({ ...form, account_type: e.target.value })}
                  className="border p-2 rounded w-full"
                >
                  <option value="bank">Bank</option>
                  <option value="kas">Kas</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Periode Awal</label>
                <input
                  type="date"
                  value={form.period_start}
                  onChange={(e) => setForm({ ...form, period_start: e.target.value })}
                  className="border p-2 rounded w-full"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Periode Akhir</label>
                <input
                  type="date"
                  value={form.period_end}
                  onChange={(e) => setForm({ ...form, period_end: e.target.value })}
                  className="border p-2 rounded w-full"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Saldo Buku (Rp)</label>
                <input
                  type="number"
                  placeholder="Saldo Buku"
                  value={form.book_balance}
                  onChange={(e) => setForm({ ...form, book_balance: e.target.value })}
                  className="border p-2 rounded w-full"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Saldo Mutasi Bank (Rp)</label>
                <input
                  type="number"
                  placeholder="Saldo Mutasi Bank"
                  value={form.statement_balance}
                  onChange={(e) => setForm({ ...form, statement_balance: e.target.value })}
                  className="border p-2 rounded w-full"
                />
              </div>
            </div>

            {form.book_balance && form.statement_balance && (
              <div className={`p-3 rounded text-sm ${
                getDiff(parseFloat(form.book_balance) || 0, parseFloat(form.statement_balance) || 0) === 0
                  ? 'bg-green-50 text-green-700'
                  : 'bg-red-50 text-red-700'
              }`}>
                Selisih: Rp {getDiff(parseFloat(form.book_balance) || 0, parseFloat(form.statement_balance) || 0).toLocaleString()}
                {getDiff(parseFloat(form.book_balance) || 0, parseFloat(form.statement_balance) || 0) === 0
                  ? ' (Sesuai)'
                  : ' (Belum Sesuai)'}
              </div>
            )}

            <div>
              <div className="flex justify-between items-center mb-2">
                <h3 className="font-medium">Item Rekonsiliasi</h3>
                <button type="button" onClick={addItem} className="text-blue-600 text-sm hover:underline">+ Tambah Item</button>
              </div>
              {items.length > 0 && (
                <div className="border rounded overflow-hidden">
                  <table className="w-full text-sm">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="p-2 text-left">Deskripsi</th>
                        <th className="p-2 text-left">Tanggal</th>
                        <th className="p-2 text-right">Jumlah</th>
                        <th className="p-2 text-left">Tipe</th>
                        <th className="p-2 text-center">Cocok</th>
                        <th className="p-2 text-center">Aksi</th>
                      </tr>
                    </thead>
                    <tbody>
                      {items.map((item, idx) => (
                        <tr key={idx} className="border-t">
                          <td className="p-1">
                            <input
                              placeholder="Deskripsi"
                              value={item.description}
                              onChange={(e) => updateItem(idx, 'description', e.target.value)}
                              className="border p-1 rounded w-full text-sm"
                            />
                          </td>
                          <td className="p-1">
                            <input
                              type="date"
                              value={item.transaction_date}
                              onChange={(e) => updateItem(idx, 'transaction_date', e.target.value)}
                              className="border p-1 rounded w-full text-sm"
                            />
                          </td>
                          <td className="p-1">
                            <input
                              type="number"
                              placeholder="0"
                              value={item.amount || ''}
                              onChange={(e) => updateItem(idx, 'amount', e.target.value)}
                              className="border p-1 rounded w-full text-sm text-right"
                            />
                          </td>
                          <td className="p-1">
                            <select
                              value={item.type}
                              onChange={(e) => updateItem(idx, 'type', e.target.value)}
                              className="border p-1 rounded w-full text-sm"
                            >
                              <option value="debit">Debit</option>
                              <option value="kredit">Kredit</option>
                            </select>
                          </td>
                          <td className="p-1 text-center">
                            <input
                              type="checkbox"
                              checked={item.is_matched}
                              onChange={(e) => updateItem(idx, 'is_matched', e.target.checked)}
                              className="rounded"
                            />
                          </td>
                          <td className="p-1 text-center">
                            <button type="button" onClick={() => removeItem(idx)} className="text-red-500 hover:underline text-sm">Hapus</button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">Catatan</label>
              <input
                placeholder="Catatan rekonsiliasi"
                value={form.notes}
                onChange={(e) => setForm({ ...form, notes: e.target.value })}
                className="border p-2 rounded w-full"
              />
            </div>

            <div className="flex gap-2">
              <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Simpan</button>
              <button type="button" onClick={() => setShowForm(false)} className="bg-gray-400 text-white px-4 py-2 rounded">Batal</button>
            </div>
          </form>
        </div>
      )}

      {detailItem && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="font-bold">Detail: {detailItem.reconciliation_number}</h2>
            <button onClick={() => setDetailItem(null)} className="text-gray-500 hover:text-gray-700">Tutup</button>
          </div>
          <div className="grid grid-cols-2 gap-4 mb-4 text-sm">
            <div><span className="text-gray-500">Rekening:</span> {detailItem.account_name} ({detailItem.account_type})</div>
            <div><span className="text-gray-500">Status:</span> <span className={`px-2 py-1 rounded text-xs ${detailItem.status === 'completed' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'}`}>{detailItem.status === 'completed' ? 'Selesai' : 'Pending'}</span></div>
            <div><span className="text-gray-500">Periode:</span> {detailItem.period_start} - {detailItem.period_end}</div>
            <div><span className="text-gray-500">Selisih:</span> Rp {detailItem.difference.toLocaleString()}</div>
            <div><span className="text-gray-500">Saldo Buku:</span> Rp {detailItem.book_balance.toLocaleString()}</div>
            <div><span className="text-gray-500">Saldo Mutasi:</span> Rp {detailItem.statement_balance.toLocaleString()}</div>
          </div>
          {detailItem.items && detailItem.items.length > 0 && (
            <div>
              <h3 className="font-medium mb-2">Item Rekonsiliasi</h3>
              <table className="w-full text-sm">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="p-2 text-left">Deskripsi</th>
                    <th className="p-2 text-left">Tanggal</th>
                    <th className="p-2 text-right">Jumlah</th>
                    <th className="p-2 text-left">Tipe</th>
                    <th className="p-2 text-center">Cocok</th>
                  </tr>
                </thead>
                <tbody>
                  {detailItem.items.map((item, idx) => (
                    <tr key={idx} className="border-t">
                      <td className="p-2">{item.description}</td>
                      <td className="p-2">{item.transaction_date}</td>
                      <td className="p-2 text-right">Rp {item.amount.toLocaleString()}</td>
                      <td className="p-2 capitalize">{item.type}</td>
                      <td className="p-2 text-center">{item.is_matched ? '✓' : '✗'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {detailItem.notes && <p className="mt-3 text-sm text-gray-600"><span className="font-medium">Catatan:</span> {detailItem.notes}</p>}
        </div>
      )}

      {loading ? (
        <p>Loading...</p>
      ) : (
        <div className="bg-white rounded shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-3 text-left">No. Rekonsiliasi</th>
                <th className="p-3 text-left">Rekening</th>
                <th className="p-3 text-left">Periode</th>
                <th className="p-3 text-right">Saldo Buku</th>
                <th className="p-3 text-right">Saldo Mutasi</th>
                <th className="p-3 text-right">Selisih</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map((d) => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.reconciliation_number}</td>
                  <td className="p-3">{d.account_name} <span className="text-xs text-gray-500">({d.account_type})</span></td>
                  <td className="p-3 text-sm">{d.period_start} - {d.period_end}</td>
                  <td className="p-3 text-right">Rp {d.book_balance.toLocaleString()}</td>
                  <td className="p-3 text-right">Rp {d.statement_balance.toLocaleString()}</td>
                  <td className={`p-3 text-right font-bold ${d.difference === 0 ? 'text-green-600' : 'text-red-600'}`}>
                    Rp {d.difference.toLocaleString()}
                  </td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs ${d.status === 'completed' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'}`}>
                      {d.status === 'completed' ? 'Selesai' : 'Pending'}
                    </span>
                  </td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => openDetail(d)} className="text-blue-500 hover:underline">Detail</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && (
                <tr>
                  <td colSpan={8} className="p-4 text-center text-gray-400">Belum ada rekonsiliasi bank</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
