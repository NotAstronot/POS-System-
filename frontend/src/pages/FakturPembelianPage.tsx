import React, { useEffect, useState } from 'react';
import { adminService, PurchaseInvoice, PurchaseOrder, Supplier } from '../services/api';

export default function FakturPembelianPage() {
  const [data, setData] = useState<PurchaseInvoice[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<PurchaseInvoice | null>(null);
  const [form, setForm] = useState({ purchase_order_id: '', supplier_id: '', invoice_date: '', due_date: '', total_amount: '', notes: '' });

  const [purchaseOrders, setPurchaseOrders] = useState<PurchaseOrder[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);

  const load = async () => {
    try {
      const [res, poRes, sRes] = await Promise.all([
        adminService.listPurchaseInvoices(),
        adminService.listPurchaseOrders(),
        adminService.listSuppliers(),
      ]);
      setData(res.data.data ?? []);
      setPurchaseOrders(poRes.data.data ?? []);
      setSuppliers(sRes.data.data ?? []);
    } catch {} finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  const openAdd = () => {
    setEditItem(null);
    setForm({ purchase_order_id: '', supplier_id: '', invoice_date: '', due_date: '', total_amount: '', notes: '' });
    setShowForm(true);
  };

  const openEdit = (d: PurchaseInvoice) => {
    setEditItem(d);
    setForm({
      purchase_order_id: String(d.purchase_order_id),
      supplier_id: String(d.supplier_id),
      invoice_date: d.invoice_date || '',
      due_date: d.due_date || '',
      total_amount: String(d.total_amount),
      notes: d.notes,
    });
    setShowForm(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload: any = {
        purchase_order_id: parseInt(form.purchase_order_id),
        supplier_id: parseInt(form.supplier_id),
        invoice_date: form.invoice_date,
        due_date: form.due_date || null,
        total_amount: parseFloat(form.total_amount),
        notes: form.notes,
      };
      if (editItem) {
        await adminService.updatePurchaseInvoice(editItem.id, payload);
      } else {
        await adminService.createPurchaseInvoice(payload);
      }
      setShowForm(false); load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus faktur pembelian ini?')) return;
    try { await adminService.deletePurchaseInvoice(id); load(); } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Faktur Pembelian</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Faktur Pembelian</button>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Faktur Pembelian' : 'Tambah Faktur Pembelian'}</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <select value={form.purchase_order_id} onChange={e => setForm({ ...form, purchase_order_id: e.target.value })} className="border p-2 rounded" required>
              <option value="">Pilih No PO</option>
              {purchaseOrders.map(po => (
                <option key={po.id} value={po.id}>{po.po_number}</option>
              ))}
            </select>
            <select value={form.supplier_id} onChange={e => setForm({ ...form, supplier_id: e.target.value })} className="border p-2 rounded" required>
              <option value="">Pilih Supplier</option>
              {suppliers.map(s => (
                <option key={s.id} value={s.id}>{s.name}</option>
              ))}
            </select>
            <input type="date" placeholder="Tanggal Invoice" value={form.invoice_date} onChange={e => setForm({ ...form, invoice_date: e.target.value })} className="border p-2 rounded" required />
            <input type="date" placeholder="Jatuh Tempo" value={form.due_date} onChange={e => setForm({ ...form, due_date: e.target.value })} className="border p-2 rounded" />
            <input type="number" placeholder="Total" value={form.total_amount} onChange={e => setForm({ ...form, total_amount: e.target.value })} className="border p-2 rounded" required />
            <input placeholder="Catatan" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} className="border p-2 rounded col-span-2" />
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
                <th className="p-3 text-left">No Faktur</th>
                <th className="p-3 text-left">No PO</th>
                <th className="p-3 text-left">Supplier</th>
                <th className="p-3 text-left">Tanggal Invoice</th>
                <th className="p-3 text-left">Jatuh Tempo</th>
                <th className="p-3 text-right">Total</th>
                <th className="p-3 text-right">Dibayar</th>
                <th className="p-3 text-right">Sisa</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(d => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.invoice_number}</td>
                  <td className="p-3">{d.po_number}</td>
                  <td className="p-3">{d.supplier_name}</td>
                  <td className="p-3">{d.invoice_date}</td>
                  <td className="p-3">{d.due_date || '-'}</td>
                  <td className="p-3 text-right">Rp {d.total_amount.toLocaleString()}</td>
                  <td className="p-3 text-right">Rp {d.paid_amount.toLocaleString()}</td>
                  <td className="p-3 text-right font-bold">Rp {d.remaining_amount.toLocaleString()}</td>
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
              {data.length === 0 && <tr><td colSpan={10} className="p-4 text-center text-gray-400">Belum ada data faktur pembelian</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
