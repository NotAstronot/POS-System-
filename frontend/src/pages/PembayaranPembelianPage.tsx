import React, { useEffect, useState } from 'react';
import { adminService, PurchasePayment, PurchaseInvoice, Supplier } from '../services/api';

export default function PembayaranPembelianPage() {
  const [data, setData] = useState<PurchasePayment[]>([]);
  const [invoices, setInvoices] = useState<PurchaseInvoice[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ purchase_invoice_id: '', supplier_id: '', payment_date: '', amount: '', payment_method: 'cash', notes: '' });

  const load = async () => {
    try {
      const [payRes, invRes, supRes] = await Promise.all([
        adminService.listPurchasePayments(),
        adminService.listPurchaseInvoices(),
        adminService.listSuppliers(),
      ]);
      setData(payRes.data.data ?? []);
      setInvoices(invRes.data.data ?? []);
      setSuppliers(supRes.data.data ?? []);
    } catch {} finally { setLoading(false); }
  };
  useEffect(() => { load(); }, []);

  const openAdd = () => {
    setForm({ purchase_invoice_id: '', supplier_id: '', payment_date: '', amount: '', payment_method: 'cash', notes: '' });
    setShowForm(true);
  };

  const handleInvoiceChange = (invoiceId: string) => {
    const inv = invoices.find(i => i.id === Number(invoiceId));
    setForm({ ...form, purchase_invoice_id: invoiceId, supplier_id: inv ? String(inv.supplier_id) : '' });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await adminService.createPurchasePayment({
        purchase_invoice_id: Number(form.purchase_invoice_id),
        supplier_id: Number(form.supplier_id),
        payment_date: form.payment_date,
        amount: parseFloat(form.amount),
        payment_method: form.payment_method,
        notes: form.notes,
      });
      setShowForm(false); load();
    } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Pembayaran Pembelian</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Pembayaran</button>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">Tambah Pembayaran</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">No Faktur</label>
              <select value={form.purchase_invoice_id} onChange={e => handleInvoiceChange(e.target.value)} className="border p-2 rounded w-full" required>
                <option value="">-- Pilih Faktur --</option>
                {invoices.filter(i => i.status !== 'paid').map(inv => (
                  <option key={inv.id} value={inv.id}>{inv.invoice_number} - Rp {inv.remaining_amount.toLocaleString()}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Supplier</label>
              <select value={form.supplier_id} onChange={e => setForm({ ...form, supplier_id: e.target.value })} className="border p-2 rounded w-full" required>
                <option value="">-- Pilih Supplier --</option>
                {suppliers.filter(s => s.is_active).map(s => (
                  <option key={s.id} value={s.id}>{s.name}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Tanggal Bayar</label>
              <input type="date" value={form.payment_date} onChange={e => setForm({ ...form, payment_date: e.target.value })} className="border p-2 rounded w-full" required />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Jumlah (Rp)</label>
              <input type="number" placeholder="Jumlah" value={form.amount} onChange={e => setForm({ ...form, amount: e.target.value })} className="border p-2 rounded w-full" required />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Metode Pembayaran</label>
              <select value={form.payment_method} onChange={e => setForm({ ...form, payment_method: e.target.value })} className="border p-2 rounded w-full" required>
                <option value="cash">Tunai</option>
                <option value="transfer">Transfer</option>
                <option value="qris">QRIS</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Keterangan</label>
              <input placeholder="Keterangan" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} className="border p-2 rounded w-full" />
            </div>
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
                <th className="p-3 text-left">No Pembayaran</th>
                <th className="p-3 text-left">No Faktur</th>
                <th className="p-3 text-left">Supplier</th>
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-right">Jumlah</th>
                <th className="p-3 text-left">Metode</th>
                <th className="p-3 text-left">Keterangan</th>
              </tr>
            </thead>
            <tbody>
              {data.map(d => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.payment_number}</td>
                  <td className="p-3">{d.invoice_number}</td>
                  <td className="p-3">{d.supplier_name}</td>
                  <td className="p-3">{d.payment_date}</td>
                  <td className="p-3 text-right">Rp {d.amount.toLocaleString()}</td>
                  <td className="p-3">
                    <span className="px-2 py-1 rounded text-xs bg-blue-100 text-blue-700">
                      {d.payment_method === 'cash' ? 'Tunai' : d.payment_method === 'transfer' ? 'Transfer' : 'QRIS'}
                    </span>
                  </td>
                  <td className="p-3">{d.notes || '-'}</td>
                </tr>
              ))}
              {data.length === 0 && <tr><td colSpan={7} className="p-4 text-center text-gray-400">Belum ada data pembayaran</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
