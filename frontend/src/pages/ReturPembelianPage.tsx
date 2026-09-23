import React, { useEffect, useState } from 'react';
import { adminService, PurchaseReturn, PurchaseOrder, Supplier } from '../services/api';

export default function ReturPembelianPage() {
  const [data, setData] = useState<PurchaseReturn[]>([]);
  const [purchaseOrders, setPurchaseOrders] = useState<PurchaseOrder[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<PurchaseReturn | null>(null);
  const [form, setForm] = useState({
    purchase_order_id: '',
    supplier_id: '',
    return_date: '',
    reason: '',
    notes: ''
  });
  const [items, setItems] = useState<{ product_name: string; quantity: string; unit_price: string; reason: string }[]>([
    { product_name: '', quantity: '1', unit_price: '', reason: '' }
  ]);

  const load = async () => {
    try {
      const [resData, resPO, resSup] = await Promise.all([
        adminService.listPurchaseReturns(),
        adminService.listPurchaseOrders(),
        adminService.listSuppliers()
      ]);
      setData(resData.data.data ?? []);
      setPurchaseOrders(resPO.data.data ?? []);
      setSuppliers(resSup.data.data ?? []);
    } catch {} finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  const openAdd = () => {
    setEditItem(null);
    setForm({ purchase_order_id: '', supplier_id: '', return_date: '', reason: '', notes: '' });
    setItems([{ product_name: '', quantity: '1', unit_price: '', reason: '' }]);
    setShowForm(true);
  };

  const openEdit = (d: PurchaseReturn) => {
    setEditItem(d);
    setForm({
      purchase_order_id: String(d.purchase_order_id),
      supplier_id: String(d.supplier_id),
      return_date: d.return_date,
      reason: d.reason,
      notes: d.notes
    });
    setItems(d.items.map(it => ({
      product_name: it.product_name,
      quantity: String(it.quantity),
      unit_price: String(it.unit_price),
      reason: it.reason
    })));
    setShowForm(true);
  };

  const addItem = () => {
    setItems([...items, { product_name: '', quantity: '1', unit_price: '', reason: '' }]);
  };

  const removeItem = (index: number) => {
    if (items.length <= 1) return;
    setItems(items.filter((_, i) => i !== index));
  };

  const updateItem = (index: number, field: string, value: string) => {
    const newItems = [...items];
    (newItems[index] as any)[field] = value;
    setItems(newItems);
  };

  const handlePOChange = (poId: string) => {
    const po = purchaseOrders.find(p => p.id === parseInt(poId));
    setForm({ ...form, purchase_order_id: poId, supplier_id: po ? String(po.supplier_id) : '' });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload: any = {
        purchase_order_id: parseInt(form.purchase_order_id),
        supplier_id: parseInt(form.supplier_id),
        return_date: form.return_date,
        reason: form.reason,
        notes: form.notes,
        items: items.map(it => ({
          product_name: it.product_name,
          quantity: parseInt(it.quantity),
          unit_price: parseFloat(it.unit_price),
          reason: it.reason
        }))
      };
      if (editItem) {
        await adminService.updatePurchaseReturn(editItem.id, payload);
      } else {
        await adminService.createPurchaseReturn(payload);
      }
      setShowForm(false); load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus retur pembelian ini?')) return;
    try { await adminService.deletePurchaseReturn(id); load(); } catch {}
  };

  const formatRupiah = (amount: number) => `Rp ${amount.toLocaleString('id-ID')}`;

  const statusBadge = (status: string) => {
    const cls = status === 'approved' ? 'bg-green-100 text-green-700'
      : status === 'rejected' ? 'bg-red-100 text-red-700'
      : 'bg-yellow-100 text-yellow-700';
    const label = status === 'approved' ? 'Disetujui' : status === 'rejected' ? 'Ditolak' : 'Menunggu';
    return <span className={`px-2 py-1 rounded text-xs ${cls}`}>{label}</span>;
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Retur Pembelian</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Retur</button>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Retur Pembelian' : 'Tambah Retur Pembelian'}</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">No PO</label>
                <select value={form.purchase_order_id} onChange={e => handlePOChange(e.target.value)} className="border p-2 rounded w-full" required>
                  <option value="">Pilih PO</option>
                  {purchaseOrders.map(po => (
                    <option key={po.id} value={po.id}>{po.po_number} - {po.supplier_name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Supplier</label>
                <select value={form.supplier_id} onChange={e => setForm({ ...form, supplier_id: e.target.value })} className="border p-2 rounded w-full" required>
                  <option value="">Pilih Supplier</option>
                  {suppliers.map(s => (
                    <option key={s.id} value={s.id}>{s.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Tanggal Retur</label>
                <input type="date" value={form.return_date} onChange={e => setForm({ ...form, return_date: e.target.value })} className="border p-2 rounded w-full" required />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Alasan</label>
                <input placeholder="Alasan retur" value={form.reason} onChange={e => setForm({ ...form, reason: e.target.value })} className="border p-2 rounded w-full" required />
              </div>
              <div className="col-span-2">
                <label className="block text-sm font-medium mb-1">Catatan</label>
                <input placeholder="Catatan" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} className="border p-2 rounded w-full" />
              </div>
            </div>

            <div>
              <div className="flex justify-between items-center mb-2">
                <label className="font-medium">Item Retur</label>
                <button type="button" onClick={addItem} className="text-blue-500 text-sm hover:underline">+ Tambah Item</button>
              </div>
              {items.map((item, idx) => (
                <div key={idx} className="flex gap-2 mb-2 items-end">
                  <input placeholder="Nama Produk" value={item.product_name} onChange={e => updateItem(idx, 'product_name', e.target.value)} className="border p-2 rounded flex-1" required />
                  <input type="number" placeholder="Jumlah" value={item.quantity} onChange={e => updateItem(idx, 'quantity', e.target.value)} className="border p-2 rounded w-24" min="1" required />
                  <input type="number" placeholder="Harga Satuan" value={item.unit_price} onChange={e => updateItem(idx, 'unit_price', e.target.value)} className="border p-2 rounded w-36" min="0" required />
                  <input placeholder="Alasan Item" value={item.reason} onChange={e => updateItem(idx, 'reason', e.target.value)} className="border p-2 rounded flex-1" required />
                  {items.length > 1 && (
                    <button type="button" onClick={() => removeItem(idx)} className="text-red-500 hover:underline text-sm px-2">Hapus</button>
                  )}
                </div>
              ))}
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
                <th className="p-3 text-left">No Retur</th>
                <th className="p-3 text-left">No PO</th>
                <th className="p-3 text-left">Supplier</th>
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-left">Alasan</th>
                <th className="p-3 text-center">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(d => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.return_number}</td>
                  <td className="p-3">{d.po_number}</td>
                  <td className="p-3">{d.supplier_name}</td>
                  <td className="p-3">{d.return_date}</td>
                  <td className="p-3">{d.reason}</td>
                  <td className="p-3 text-center">{statusBadge(d.status)}</td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => openEdit(d)} className="text-blue-500 hover:underline">Edit</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && <tr><td colSpan={7} className="p-4 text-center text-gray-400">Belum ada data retur pembelian</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
