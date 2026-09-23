import React, { useEffect, useState } from 'react';
import { adminService, SalesQuotation, Customer, SalesCategory } from '../services/api';

interface QuotationItem {
  product_name: string;
  quantity: string;
  unit: string;
  unit_price: string;
}

const emptyItem: QuotationItem = { product_name: '', quantity: '', unit: '', unit_price: '' };

const statusConfig: Record<string, { label: string; className: string }> = {
  draft: { label: 'Draft', className: 'bg-blue-100 text-blue-700' },
  accepted: { label: 'Diterima', className: 'bg-green-100 text-green-700' },
  rejected: { label: 'Ditolak', className: 'bg-red-100 text-red-700' },
  expired: { label: 'Kadaluarsa', className: 'bg-gray-100 text-gray-700' },
};

export default function PenawaranPenjualanPage() {
  const [data, setData] = useState<SalesQuotation[]>([]);
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [salesCategories, setSalesCategories] = useState<SalesCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<SalesQuotation | null>(null);
  const [form, setForm] = useState({
    customer_id: '',
    sales_category_id: '',
    quotation_date: '',
    valid_until: '',
    notes: '',
  });
  const [items, setItems] = useState<QuotationItem[]>([{ ...emptyItem }]);
  const [error, setError] = useState('');

  const load = async () => {
    setError('');
    try {
      const [qRes, cRes, scRes] = await Promise.all([
        adminService.listSalesQuotations(),
        adminService.listCustomers(),
        adminService.listSalesCategories(),
      ]);
      setData(qRes.data.data ?? []);
      setCustomers((cRes.data.data ?? []).filter((c: Customer) => c.is_active));
      setSalesCategories((scRes.data.data ?? []).filter((c: SalesCategory) => c.is_active));
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
    setForm({ customer_id: '', sales_category_id: '', quotation_date: '', valid_until: '', notes: '' });
    setItems([{ ...emptyItem }]);
    setShowForm(true);
  };

  const openEdit = (q: SalesQuotation) => {
    setEditItem(q);
    setForm({
      customer_id: String(q.customer_id),
      sales_category_id: q.sales_category_id ? String(q.sales_category_id) : '',
      quotation_date: q.quotation_date,
      valid_until: q.valid_until || '',
      notes: q.notes || '',
    });
    setItems(
      q.items.length > 0
        ? q.items.map(i => ({
            product_name: i.product_name,
            quantity: String(i.quantity),
            unit: i.unit,
            unit_price: String(i.unit_price),
          }))
        : [{ ...emptyItem }]
    );
    setShowForm(true);
  };

  const handleItemChange = (index: number, field: keyof QuotationItem, value: string) => {
    const updated = [...items];
    updated[index][field] = value;
    setItems(updated);
  };

  const addItem = () => setItems([...items, { ...emptyItem }]);
  const removeItem = (index: number) => setItems(items.filter((_, i) => i !== index));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload: any = {
        customer_id: parseInt(form.customer_id),
        sales_category_id: form.sales_category_id ? parseInt(form.sales_category_id) : null,
        quotation_date: form.quotation_date,
        valid_until: form.valid_until || null,
        notes: form.notes,
        items: items
          .filter(i => i.product_name)
          .map(i => ({
            product_name: i.product_name,
            quantity: parseFloat(i.quantity) || 0,
            unit: i.unit,
            unit_price: parseFloat(i.unit_price) || 0,
          })),
      };
      if (editItem) {
        await adminService.updateSalesQuotation(editItem.id, payload);
      } else {
        await adminService.createSalesQuotation(payload);
      }
      setShowForm(false);
      load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus penawaran penjualan ini?')) return;
    try {
      await adminService.deleteSalesQuotation(id);
      load();
    } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Penawaran Penjualan</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Penawaran</button>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Penawaran Penjualan' : 'Tambah Penawaran Penjualan'}</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-3 gap-4">
              <select value={form.customer_id} onChange={e => setForm({ ...form, customer_id: e.target.value })} className="border p-2 rounded" required>
                <option value="">-- Pilih Pelanggan --</option>
                {customers.map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
              <select value={form.sales_category_id} onChange={e => setForm({ ...form, sales_category_id: e.target.value })} className="border p-2 rounded">
                <option value="">-- Pilih Kategori Penjualan --</option>
                {salesCategories.map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
              <input type="date" placeholder="Tanggal Penawaran" value={form.quotation_date} onChange={e => setForm({ ...form, quotation_date: e.target.value })} className="border p-2 rounded" required />
              <input type="date" placeholder="Berlaku Hingga" value={form.valid_until} onChange={e => setForm({ ...form, valid_until: e.target.value })} className="border p-2 rounded" />
            </div>
            <input placeholder="Catatan" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} className="border p-2 rounded w-full" />

            <div>
              <div className="flex justify-between items-center mb-2">
                <h3 className="font-semibold">Item Penawaran</h3>
                <button type="button" onClick={addItem} className="text-blue-600 text-sm hover:underline">+ Tambah Item</button>
              </div>
              <table className="w-full border">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="p-2 text-left border">Nama Produk</th>
                    <th className="p-2 text-left border">Jumlah</th>
                    <th className="p-2 text-left border">Satuan</th>
                    <th className="p-2 text-left border">Harga Satuan</th>
                    <th className="p-2 text-center border w-16">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map((item, idx) => (
                    <tr key={idx}>
                      <td className="p-1 border"><input value={item.product_name} onChange={e => handleItemChange(idx, 'product_name', e.target.value)} className="border p-1 rounded w-full" /></td>
                      <td className="p-1 border"><input type="number" value={item.quantity} onChange={e => handleItemChange(idx, 'quantity', e.target.value)} className="border p-1 rounded w-full" /></td>
                      <td className="p-1 border"><input value={item.unit} onChange={e => handleItemChange(idx, 'unit', e.target.value)} className="border p-1 rounded w-full" placeholder="pcs/kg/box" /></td>
                      <td className="p-1 border"><input type="number" value={item.unit_price} onChange={e => handleItemChange(idx, 'unit_price', e.target.value)} className="border p-1 rounded w-full" /></td>
                      <td className="p-1 border text-center">
                        {items.length > 1 && (
                          <button type="button" onClick={() => removeItem(idx)} className="text-red-500 hover:underline text-sm">Hapus</button>
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
          {error && <div className="p-4 text-red-500 bg-red-50 border-b">{error}</div>}
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-3 text-left">No Penawaran</th>
                <th className="p-3 text-left">Pelanggan</th>
                <th className="p-3 text-left">Kategori</th>
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-left">Berlaku Hingga</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-right">Total</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(q => {
                const st = statusConfig[q.status] || statusConfig.draft;
                return (
                  <tr key={q.id} className="border-t">
                    <td className="p-3 font-medium">{q.quotation_number}</td>
                    <td className="p-3">{q.customer_name}</td>
                    <td className="p-3">{q.sales_category_name || '-'}</td>
                    <td className="p-3">{q.quotation_date}</td>
                    <td className="p-3">{q.valid_until || '-'}</td>
                    <td className="p-3">
                      <span className={`px-2 py-1 rounded text-xs ${st.className}`}>{st.label}</span>
                    </td>
                    <td className="p-3 text-right">Rp {q.total_amount.toLocaleString()}</td>
                    <td className="p-3 text-center gap-2 flex justify-center">
                      <button onClick={() => openEdit(q)} className="text-blue-500 hover:underline">Edit</button>
                      <button onClick={() => handleDelete(q.id)} className="text-red-500 hover:underline">Hapus</button>
                    </td>
                  </tr>
                );
              })}
              {data.length === 0 && <tr><td colSpan={8} className="p-4 text-center text-gray-400">Belum ada data penawaran penjualan</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}