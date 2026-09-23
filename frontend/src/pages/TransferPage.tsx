import React, { useEffect, useState } from 'react';
import { inventoryService, StockTransfer, Warehouse, Item } from '../services/api';

interface TransferItemRow {
  product_id: string;
  product_name: string;
  quantity: string;
  unit: string;
}

const emptyRow: TransferItemRow = { product_id: '', product_name: '', quantity: '', unit: 'pcs' };

const statusConfig: Record<string, { label: string; className: string }> = {
  draft: { label: 'Draft', className: 'bg-gray-100 text-gray-700' },
  sent: { label: 'Dikirim', className: 'bg-blue-100 text-blue-700' },
  received: { label: 'Diterima', className: 'bg-green-100 text-green-700' },
  cancelled: { label: 'Dibatalkan', className: 'bg-red-100 text-red-700' },
};

export default function TransferPage() {
  const [data, setData] = useState<StockTransfer[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<StockTransfer | null>(null);
  const [form, setForm] = useState({ from_warehouse_id: '', to_warehouse_id: '', transfer_date: '', notes: '' });
  const [rows, setRows] = useState<TransferItemRow[]>([{ ...emptyRow }]);
  const [error, setError] = useState('');

  const load = async () => {
    setError('');
    try {
      const [tRes, wRes, iRes] = await Promise.all([
        inventoryService.listTransfers(),
        inventoryService.listWarehouses(),
        inventoryService.listItems(),
      ]);
      setData(tRes.data.data ?? []);
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
    setEditItem(null);
    setForm({ from_warehouse_id: '', to_warehouse_id: '', transfer_date: new Date().toISOString().slice(0, 10), notes: '' });
    setRows([{ ...emptyRow }]);
    setShowForm(true);
  };

  const openEdit = (t: StockTransfer) => {
    setEditItem(t);
    setForm({
      from_warehouse_id: String(t.from_warehouse_id),
      to_warehouse_id: String(t.to_warehouse_id),
      transfer_date: t.transfer_date,
      notes: t.notes || '',
    });
    setRows(
      (t.items && t.items.length > 0)
        ? t.items.map(i => ({ product_id: String(i.product_id), product_name: i.product_name, quantity: String(i.quantity), unit: i.unit }))
        : [{ ...emptyRow }]
    );
    setShowForm(true);
  };

  const handleItemChange = (idx: number, field: keyof TransferItemRow, value: string) => {
    const updated = [...rows];
    if (field === 'product_id') {
      const it = items.find(x => String(x.id) === value);
      updated[idx].product_id = value;
      updated[idx].product_name = it ? it.name : '';
      updated[idx].unit = it ? it.unit : 'pcs';
    } else {
      (updated[idx] as any)[field] = value;
    }
    setRows(updated);
  };

  const addRow = () => setRows([...rows, { ...emptyRow }]);
  const removeRow = (idx: number) => setRows(rows.filter((_, i) => i !== idx));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload: any = {
        from_warehouse_id: parseInt(form.from_warehouse_id),
        to_warehouse_id: parseInt(form.to_warehouse_id),
        transfer_date: form.transfer_date,
        notes: form.notes,
        items: rows
          .filter(r => r.product_id)
          .map(r => ({
            product_id: parseInt(r.product_id),
            product_name: r.product_name,
            quantity: parseFloat(r.quantity) || 0,
            unit: r.unit,
          })),
      };
      if (editItem) {
        await inventoryService.updateTransfer(editItem.id, payload);
      } else {
        await inventoryService.createTransfer(payload);
      }
      setShowForm(false);
      load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal menyimpan data');
    }
  };

  const doAction = async (id: number, action: 'send' | 'receive' | 'cancel' | 'delete') => {
    const labels: Record<string, string> = {
      send: 'Kirim transfer ini?',
      receive: 'Terima transfer ini?',
      cancel: 'Batalkan transfer ini?',
      delete: 'Hapus transfer ini?',
    };
    if (!confirm(labels[action])) return;
    try {
      if (action === 'send') await inventoryService.sendTransfer(id);
      if (action === 'receive') await inventoryService.receiveTransfer(id);
      if (action === 'cancel') await inventoryService.cancelTransfer(id);
      if (action === 'delete') await inventoryService.deleteTransfer(id);
      load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Aksi gagal');
    }
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Transfer Barang</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Transfer</button>
      </div>

      {error && <div className="mb-4 p-4 text-red-500 bg-red-50 border border-red-200 rounded">{error}</div>}

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6 border border-gray-200">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Transfer Barang' : 'Tambah Transfer Barang'}</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-3 gap-4">
              <select value={form.from_warehouse_id} onChange={e => setForm({ ...form, from_warehouse_id: e.target.value })} className="border p-2 rounded" required>
                <option value="">-- Dari Gudang --</option>
                {warehouses.map(w => (
                  <option key={w.id} value={w.id}>{w.name}</option>
                ))}
              </select>
              <select value={form.to_warehouse_id} onChange={e => setForm({ ...form, to_warehouse_id: e.target.value })} className="border p-2 rounded" required>
                <option value="">-- Ke Gudang --</option>
                {warehouses.filter(w => String(w.id) !== form.from_warehouse_id).map(w => (
                  <option key={w.id} value={w.id}>{w.name}</option>
                ))}
              </select>
              <input type="date" value={form.transfer_date} onChange={e => setForm({ ...form, transfer_date: e.target.value })} className="border p-2 rounded" required />
            </div>
            <input placeholder="Catatan" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} className="border p-2 rounded w-full" />

            <div>
              <div className="flex justify-between items-center mb-2">
                <h3 className="font-semibold">Item Transfer</h3>
                <button type="button" onClick={addRow} className="text-blue-600 text-sm hover:underline">+ Tambah Item</button>
              </div>
              <table className="w-full border">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="p-2 text-left border">Barang</th>
                    <th className="p-2 text-left border">Jumlah</th>
                    <th className="p-2 text-left border">Satuan</th>
                    <th className="p-2 text-center border w-16">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {rows.map((r, idx) => (
                    <tr key={idx}>
                      <td className="p-1 border">
                        <select value={r.product_id} onChange={e => handleItemChange(idx, 'product_id', e.target.value)} className="border p-1 rounded w-full">
                          <option value="">-- Pilih Barang --</option>
                          {items.map(it => (
                            <option key={it.id} value={it.id}>{it.name} (stok {it.stock})</option>
                          ))}
                        </select>
                      </td>
                      <td className="p-1 border"><input type="number" min="0" value={r.quantity} onChange={e => handleItemChange(idx, 'quantity', e.target.value)} className="border p-1 rounded w-full" /></td>
                      <td className="p-1 border"><input value={r.unit} onChange={e => handleItemChange(idx, 'unit', e.target.value)} className="border p-1 rounded w-full" /></td>
                      <td className="p-1 border text-center">
                        {rows.length > 1 && (
                          <button type="button" onClick={() => removeRow(idx)} className="text-red-500 hover:underline text-sm">Hapus</button>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="flex gap-2">
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
                <th className="p-3 text-left">No Transfer</th>
                <th className="p-3 text-left">Dari</th>
                <th className="p-3 text-left">Ke</th>
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(t => {
                const st = statusConfig[t.status] || statusConfig.draft;
                return (
                  <tr key={t.id} className="border-t">
                    <td className="p-3 font-medium">{t.transfer_number}</td>
                    <td className="p-3">{t.from_warehouse_name}</td>
                    <td className="p-3">{t.to_warehouse_name}</td>
                    <td className="p-3">{t.transfer_date}</td>
                    <td className="p-3">
                      <span className={`px-2 py-1 rounded text-xs ${st.className}`}>{st.label}</span>
                    </td>
                    <td className="p-3 text-center gap-2 flex justify-center">
                      {t.status === 'draft' && (
                        <>
                          <button onClick={() => openEdit(t)} className="text-blue-500 hover:underline">Edit</button>
                          <button onClick={() => doAction(t.id, 'send')} className="text-green-600 hover:underline">Kirim</button>
                        </>
                      )}
                      {t.status === 'sent' && (
                        <button onClick={() => doAction(t.id, 'receive')} className="text-green-600 hover:underline">Terima</button>
                      )}
                      {(t.status === 'draft' || t.status === 'sent') && (
                        <button onClick={() => doAction(t.id, 'cancel')} className="text-red-500 hover:underline">Batalkan</button>
                      )}
                      {t.status === 'draft' && (
                        <button onClick={() => doAction(t.id, 'delete')} className="text-red-500 hover:underline">Hapus</button>
                      )}
                    </td>
                  </tr>
                );
              })}
              {data.length === 0 && <tr><td colSpan={6} className="p-4 text-center text-gray-400">Belum ada transfer barang</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}