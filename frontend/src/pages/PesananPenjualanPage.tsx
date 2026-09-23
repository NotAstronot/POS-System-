import React, { useEffect, useState } from 'react';
import { adminService, SalesOrder, Customer, SalesCategory, SalesQuotation } from '../services/api';

interface OrderItem {
  product_name: string;
  quantity: string;
  unit: string;
  unit_price: string;
}

const emptyItem: OrderItem = { product_name: '', quantity: '', unit: '', unit_price: '' };

const statusConfig: Record<string, { label: string; className: string }> = {
  draft: { label: 'Draft', className: 'bg-blue-100 text-blue-700' },
  partial: { label: 'Sebagian', className: 'bg-yellow-100 text-yellow-700' },
  delivered: { label: 'Terkirim', className: 'bg-green-100 text-green-700' },
  cancelled: { label: 'Dibatalkan', className: 'bg-red-100 text-red-700' },
};

export default function PesananPenjualanPage() {
  const [data, setData] = useState<SalesOrder[]>([]);
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [salesCategories, setSalesCategories] = useState<SalesCategory[]>([]);
  const [quotations, setQuotations] = useState<SalesQuotation[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<SalesOrder | null>(null);
  const [form, setForm] = useState({
    customer_id: '',
    quotation_id: '',
    sales_category_id: '',
    order_date: '',
    expected_date: '',
    notes: '',
  });
  const [items, setItems] = useState<OrderItem[]>([{ ...emptyItem }]);
  const [error, setError] = useState('');

  const load = async () => {
    setError('');
    try {
      const [oRes, cRes, scRes, qRes] = await Promise.all([
        adminService.listSalesOrders(),
        adminService.listCustomers(),
        adminService.listSalesCategories(),
        adminService.listSalesQuotations(),
      ]);
      setData(oRes.data.data ?? []);
      setCustomers((cRes.data.data ?? []).filter((c: Customer) => c.is_active));
      setSalesCategories((scRes.data.data ?? []).filter((c: SalesCategory) => c.is_active));
      setQuotations(qRes.data.data ?? []);
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
    setForm({ customer_id: '', quotation_id: '', sales_category_id: '', order_date: '', expected_date: '', notes: '' });
    setItems([{ ...emptyItem }]);
    setShowForm(true);
  };

  const openEdit = (o: SalesOrder) => {
    setEditItem(o);
    setForm({
      customer_id: String(o.customer_id),
      quotation_id: o.quotation_id ? String(o.quotation_id) : '',
      sales_category_id: o.sales_category_id ? String(o.sales_category_id) : '',
      order_date: o.order_date,
      expected_date: o.expected_date || '',
      notes: o.notes || '',
    });
    setItems(
      o.items.length > 0
        ? o.items.map(i => ({
            product_name: i.product_name,
            quantity: String(i.quantity),
            unit: i.unit,
            unit_price: String(i.unit_price),
          }))
        : [{ ...emptyItem }]
    );
    setShowForm(true);
  };

  const handleItemChange = (index: number, field: keyof OrderItem, value: string) => {
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
        quotation_id: form.quotation_id ? parseInt(form.quotation_id) : null,
        sales_category_id: form.sales_category_id ? parseInt(form.sales_category_id) : null,
        order_date: form.order_date,
        expected_date: form.expected_date || null,
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
        await adminService.updateSalesOrder(editItem.id, payload);
      } else {
        await adminService.createSalesOrder(payload);
      }
      setShowForm(false);
      load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus pesanan penjualan ini?')) return;
    try {
      await adminService.deleteSalesOrder(id);
      load();
    } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Pesanan Penjualan</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Pesanan</button>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Pesanan Penjualan' : 'Tambah Pesanan Penjualan'}</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-3 gap-4">
              <select value={form.customer_id} onChange={e => setForm({ ...form, customer_id: e.target.value })} className="border p-2 rounded" required>
                <option value="">-- Pilih Pelanggan --</option>
                {customers.map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
              <select value={form.quotation_id} onChange={e => setForm({ ...form, quotation_id: e.target.value })} className="border p-2 rounded">
                <option value="">-- Pilih Penawaran --</option>
                {quotations.map(q => (
                  <option key={q.id} value={q.id}>{q.quotation_number}</option>
                ))}
              </select>
              <select value={form.sales_category_id} onChange={e => setForm({ ...form, sales_category_id: e.target.value })} className="border p-2 rounded">
                <option value="">-- Pilih Kategori Penjualan --</option>
                {salesCategories.map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
              <input type="date" placeholder="Tanggal Pesanan" value={form.order_date} onChange={e => setForm({ ...form, order_date: e.target.value })} className="border p-2 rounded" required />
              <input type="date" placeholder="Jatuh Tempo" value={form.expected_date} onChange={e => setForm({ ...form, expected_date: e.target.value })} className="border p-2 rounded" />
            </div>
            <input placeholder="Catatan" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} className="border p-2 rounded w-full" />

            <div>
              <div className="flex justify-between items-center mb-2">
                <h3 className="font-semibold">Item Pesanan</h3>
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
                <th className="p-3 text-left">No Pesanan</th>
                <th className="p-3 text-left">Pelanggan</th>
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-left">Jatuh Tempo</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-right">Total</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(o => {
                const st = statusConfig[o.status] || statusConfig.draft;
                return (
                  <tr key={o.id} className="border-t">
                    <td className="p-3 font-medium">{o.order_number}</td>
                    <td className="p-3">{o.customer_name}</td>
                    <td className="p-3">{o.order_date}</td>
                    <td className="p-3">{o.expected_date || '-'}</td>
                    <td className="p-3">
                      <span className={`px-2 py-1 rounded text-xs ${st.className}`}>{st.label}</span>
                    </td>
                    <td className="p-3 text-right">Rp {o.total_amount.toLocaleString()}</td>
                    <td className="p-3 text-center gap-2 flex justify-center">
                      <button onClick={() => openEdit(o)} className="text-blue-500 hover:underline">Edit</button>
                      <button onClick={() => handleDelete(o.id)} className="text-red-500 hover:underline">Hapus</button>
                    </td>
                  </tr>
                );
              })}
              {data.length === 0 && <tr><td colSpan={7} className="p-4 text-center text-gray-400">Belum ada data pesanan penjualan</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}