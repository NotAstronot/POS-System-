import React, { useEffect, useState } from 'react';
import { inventoryService, StockAdjustment, Warehouse, Item } from '../services/api';

interface AdjRow {
  product_id: string;
  product_name: string;
  system_qty: number;
  actual_qty: string;
  reason: string;
}

const statusConfig: Record<string, { label: string; className: string }> = {
  draft: { label: 'Draft', className: 'bg-gray-100 text-gray-700' },
  applied: { label: 'Diterapkan', className: 'bg-green-100 text-green-700' },
};

export default function StockAdjustmentPage() {
  const [data, setData] = useState<StockAdjustment[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ warehouse_id: '', adjustment_date: '', notes: '' });
  const [rows, setRows] = useState<AdjRow[]>([]);
  const [selectedWh, setSelectedWh] = useState(0);
  const [error, setError] = useState('');

  const load = async () => {
    setError('');
    try {
      const [aRes, wRes, iRes] = await Promise.all([
        inventoryService.listAdjustments(),
        inventoryService.listWarehouses(),
        inventoryService.listItems(),
      ]);
      setData(aRes.data.data ?? []);
      setWarehouses((wRes.data.data ?? []).filter((w: Warehouse) => w.is_active));
      setItems((iRes.data.data ?? []).filter((i: Item) => !i.is_service));
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
    setForm({ warehouse_id: '', adjustment_date: new Date().toISOString().slice(0, 10), notes: '' });
    setRows([]);
    setSelectedWh(0);
    setShowForm(true);
  };

  const handleWarehouseChange = (whId: string) => {
    setForm({ ...form, warehouse_id: whId });
    setRows([]);
    setSelectedWh(parseInt(whId) || 0);
  };

  const addRow = () => {
    setRows([...rows, { product_id: '', product_name: '', system_qty: 0, actual_qty: '', reason: '' }]);
  };

  const handleRowChange = async (idx: number, field: keyof AdjRow, value: string) => {
    const updated = [...rows];
    if (field === 'product_id') {
      const it = items.find(x => String(x.id) === value);
      updated[idx].product_id = value;
      updated[idx].product_name = it ? it.name : '';
      if (it && selectedWh) {
        try {
          const detail = await inventoryService.getItem(it.id);
          const whStock = detail.data.data.stock_by_warehouse.find(s => s.warehouse_id === selectedWh);
          updated[idx].system_qty = whStock ? whStock.quantity : 0;
        } catch {
          updated[idx].system_qty = it.stock;
        }
      }
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
        warehouse_id: parseInt(form.warehouse_id),
        adjustment_date: form.adjustment_date,
        notes: form.notes,
        items: rows
          .filter(r => r.product_id)
          .map(r => ({
            product_id: parseInt(r.product_id),
            product_name: r.product_name,
            system_qty: r.system_qty,
            actual_qty: parseFloat(r.actual_qty) || 0,
            reason: r.reason,
          })),
      };
      await inventoryService.createAdjustment(payload);
      setShowForm(false);
      load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal menyimpan data');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus penyesuaian persediaan ini?')) return;
    try {
      await inventoryService.deleteAdjustment(id);
      load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal menghapus data');
    }
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Penyesuaian Persediaan</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Opname</button>
      </div>

      {error && <div className="mb-4 p-4 text-red-500 bg-red-50 border border-red-200 rounded">{error}</div>}

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6 border border-gray-200">
          <h2 className="font-bold mb-4">Tambah Penyesuaian (Stock Opname)</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-3 gap-4">
              <select value={form.warehouse_id} onChange={e => handleWarehouseChange(e.target.value)} className="border p-2 rounded" required>
                <option value="">-- Pilih Gudang --</option>
                {warehouses.map(w => (
                  <option key={w.id} value={w.id}>{w.name}</option>
                ))}
              </select>
              <input type="date" value={form.adjustment_date} onChange={e => setForm({ ...form, adjustment_date: e.target.value })} className="border p-2 rounded" required />
              <input placeholder="Catatan" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} className="border p-2 rounded" />
            </div>

            {rows.length === 0 && (
              <div className="p-4 bg-yellow-50 border border-yellow-200 rounded text-sm text-yellow-800">
                Pilih gudang terlebih dahulu, lalu tambahkan item yang akan diopname.
              </div>
            )}

            {rows.length > 0 && (
              <div>
                <div className="flex justify-between items-center mb-2">
                  <h3 className="font-semibold">Item Opname</h3>
                  <button type="button" onClick={addRow} className="text-blue-600 text-sm hover:underline">+ Tambah Item</button>
                </div>
                <table className="w-full border">
                  <thead className="bg-gray-100">
                    <tr>
                      <th className="p-2 text-left border">Barang</th>
                      <th className="p-2 text-right border">Stok Sistem</th>
                      <th className="p-2 text-right border">Stok Fisik</th>
                      <th className="p-2 text-left border">Selisih</th>
                      <th className="p-2 text-left border">Alasan</th>
                      <th className="p-2 text-center border w-16">Aksi</th>
                    </tr>
                  </thead>
                  <tbody>
                    {rows.map((r, idx) => {
                      const diff = (parseFloat(r.actual_qty) || 0) - r.system_qty;
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
                          <td className="p-1 border text-right">{r.system_qty}</td>
                          <td className="p-1 border"><input type="number" min="0" value={r.actual_qty} onChange={e => handleRowChange(idx, 'actual_qty', e.target.value)} className="border p-1 rounded w-full text-right" /></td>
                          <td className={`p-1 border text-right font-bold ${diff !== 0 ? (diff > 0 ? 'text-green-600' : 'text-red-600') : ''}`}>
                            {diff > 0 ? `+${diff}` : diff}
                          </td>
                          <td className="p-1 border"><input value={r.reason} onChange={e => handleRowChange(idx, 'reason', e.target.value)} className="border p-1 rounded w-full" placeholder="Rusak/hilang/lebih" /></td>
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
              <button type="button" onClick={addRow} className="bg-blue-600 text-white px-4 py-2 rounded">+ Item Opname</button>
              <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Simpan & Terapkan</button>
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
                <th className="p-3 text-left">No Opname</th>
                <th className="p-3 text-left">Gudang</th>
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-left">Catatan</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(a => {
                const st = statusConfig[a.status] || statusConfig.draft;
                return (
                  <tr key={a.id} className="border-t">
                    <td className="p-3 font-medium">{a.adjustment_number}</td>
                    <td className="p-3">{a.warehouse_name}</td>
                    <td className="p-3">{a.adjustment_date}</td>
                    <td className="p-3">
                      <span className={`px-2 py-1 rounded text-xs ${st.className}`}>{st.label}</span>
                    </td>
                    <td className="p-3">{a.notes || '-'}</td>
                    <td className="p-3 text-center gap-2 flex justify-center">
                      <button onClick={() => handleDelete(a.id)} className="text-red-500 hover:underline">Hapus</button>
                    </td>
                  </tr>
                );
              })}
              {data.length === 0 && <tr><td colSpan={6} className="p-4 text-center text-gray-400">Belum ada penyesuaian persediaan</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}