import React, { useEffect, useState } from 'react';
import { inventoryService, PriceChange, Item } from '../services/api';

interface PriceRow {
  product_id: string;
  product_name: string;
  old_price: number;
  new_price: string;
}

const statusConfig: Record<string, { label: string; className: string }> = {
  applied: { label: 'Diterapkan', className: 'bg-green-100 text-green-700' },
  scheduled: { label: 'Terjadwal', className: 'bg-yellow-100 text-yellow-700' },
  cancelled: { label: 'Dibatalkan', className: 'bg-red-100 text-red-700' },
};

export default function PriceChangePage() {
  const [data, setData] = useState<PriceChange[]>([]);
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ change_date: '', notes: '' });
  const [rows, setRows] = useState<PriceRow[]>([]);
  const [error, setError] = useState('');

  const load = async () => {
    setError('');
    try {
      const [pRes, iRes] = await Promise.all([
        inventoryService.listPriceChanges(),
        inventoryService.listItems(),
      ]);
      setData(pRes.data.data ?? []);
      setItems(iRes.data.data ?? []);
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal memuat data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openAdd = () => {
    setForm({ change_date: new Date().toISOString().slice(0, 10), notes: '' });
    setRows([]);
    setShowForm(true);
  };

  const addRow = () => {
    setRows([...rows, { product_id: '', product_name: '', old_price: 0, new_price: '' }]);
  };

  const handleRowChange = (idx: number, field: keyof PriceRow, value: string) => {
    const updated = [...rows];
    if (field === 'product_id') {
      const it = items.find(x => String(x.id) === value);
      updated[idx].product_id = value;
      updated[idx].product_name = it ? it.name : '';
      updated[idx].old_price = it ? it.base_price : 0;
      updated[idx].new_price = it ? String(it.base_price) : '';
    } else {
      (updated[idx] as any)[field] = value;
    }
    setRows(updated);
  };

  const removeRow = (idx: number) => setRows(rows.filter((_, i) => i !== idx));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload: any = {
        change_date: form.change_date,
        notes: form.notes,
        items: rows
          .filter(r => r.product_id)
          .map(r => ({
            product_id: parseInt(r.product_id),
            product_name: r.product_name,
            old_price: r.old_price,
            new_price: parseFloat(r.new_price) || 0,
          })),
      };
      await inventoryService.createPriceChange(payload);
      setShowForm(false);
      load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal menyimpan data');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus penyesuaian harga ini? Harga yang sudah diterapkan tetap berlaku.')) return;
    try {
      await inventoryService.deletePriceChange(id);
      load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal menghapus data');
    }
  };

  const formatRupiah = (n: number) => `Rp ${(n || 0).toLocaleString('id-ID')}`;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Penyesuaian Harga Jual</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Penyesuaian Harga</button>
      </div>

      {error && <div className="mb-4 p-4 text-red-500 bg-red-50 border border-red-200 rounded">{error}</div>}

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6 border border-gray-200">
          <h2 className="font-bold mb-4">Tambah Penyesuaian Harga</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <input type="date" value={form.change_date} onChange={e => setForm({ ...form, change_date: e.target.value })} className="border p-2 rounded" required />
              <input placeholder="Catatan (mis. Kenaikan harga bahan baku)" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} className="border p-2 rounded" />
            </div>

            {rows.length === 0 && (
              <div className="p-4 bg-yellow-50 border border-yellow-200 rounded text-sm text-yellow-800">
                Tambahkan barang yang harganya ingin disesuaikan.
              </div>
            )}

            {rows.length > 0 && (
              <div>
                <div className="flex justify-between items-center mb-2">
                  <h3 className="font-semibold">Item Penyesuaian</h3>
                  <button type="button" onClick={addRow} className="text-blue-600 text-sm hover:underline">+ Tambah Item</button>
                </div>
                <table className="w-full border">
                  <thead className="bg-gray-100">
                    <tr>
                      <th className="p-2 text-left border">Barang</th>
                      <th className="p-2 text-right border">Harga Lama</th>
                      <th className="p-2 text-right border">Harga Baru</th>
                      <th className="p-2 text-right border">Selisih</th>
                      <th className="p-2 text-center border w-16">Aksi</th>
                    </tr>
                  </thead>
                  <tbody>
                    {rows.map((r, idx) => {
                      const diff = (parseFloat(r.new_price) || 0) - r.old_price;
                      return (
                        <tr key={idx}>
                          <td className="p-1 border">
                            <select value={r.product_id} onChange={e => handleRowChange(idx, 'product_id', e.target.value)} className="border p-1 rounded w-full">
                              <option value="">-- Pilih Barang --</option>
                              {items.map(it => (
                                <option key={it.id} value={it.id}>{it.name}</option>
                              ))}
                            </select>
                          </td>
                          <td className="p-1 border text-right">{formatRupiah(r.old_price)}</td>
                          <td className="p-1 border"><input type="number" min="0" value={r.new_price} onChange={e => handleRowChange(idx, 'new_price', e.target.value)} className="border p-1 rounded w-full text-right" /></td>
                          <td className={`p-1 border text-right font-bold ${diff !== 0 ? (diff > 0 ? 'text-green-600' : 'text-red-600') : ''}`}>
                            {diff > 0 ? `+${formatRupiah(diff)}` : diff < 0 ? `-${formatRupiah(Math.abs(diff))}` : '-'}
                          </td>
                          <td className="p-1 border text-center">
                            {rows.length > 1 && (
                              <button type="button" onClick={() => removeRow(idx)} className="text-red-500 hover:underline text-sm">Hapus</button>
                            )}
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}

            <div className="flex gap-2">
              <button type="button" onClick={addRow} className="bg-blue-600 text-white px-4 py-2 rounded">+ Item</button>
              <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Terapkan Harga</button>
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
                <th className="p-3 text-left">No Penyesuaian</th>
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-left">Catatan</th>
                <th className="p-3 text-left">Item</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(p => {
                const st = statusConfig[p.status] || statusConfig.applied;
                return (
                  <tr key={p.id} className="border-t">
                    <td className="p-3 font-medium">{p.price_change_number}</td>
                    <td className="p-3">{p.change_date}</td>
                    <td className="p-3">
                      <span className={`px-2 py-1 rounded text-xs ${st.className}`}>{st.label}</span>
                    </td>
                    <td className="p-3">{p.notes || '-'}</td>
                    <td className="p-3 text-sm text-gray-600">
                      {(p.items && p.items.length > 0) ? p.items.map(i => `${i.product_name} (${formatRupiah(i.old_price)} → ${formatRupiah(i.new_price)})`).join('; ') : '-'}
                    </td>
                    <td className="p-3 text-center">
                      <button onClick={() => handleDelete(p.id)} className="text-red-500 hover:underline">Hapus</button>
                    </td>
                  </tr>
                );
              })}
              {data.length === 0 && <tr><td colSpan={6} className="p-4 text-center text-gray-400">Belum ada penyesuaian harga</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}