import React, { useEffect, useState } from 'react';
import { adminService, DeliveryOrder, SalesOrder } from '../services/api';

interface DeliveryItem {
  product_name: string;
  quantity: string;
  unit: string;
}

const emptyItem: DeliveryItem = { product_name: '', quantity: '', unit: '' };

const statusConfig: Record<string, { label: string; className: string }> = {
  shipped: { label: 'Dikirim', className: 'bg-blue-100 text-blue-700' },
  delivered: { label: 'Terkirim', className: 'bg-green-100 text-green-700' },
  cancelled: { label: 'Dibatalkan', className: 'bg-red-100 text-red-700' },
};

export default function PengirimanPesananPage() {
  const [data, setData] = useState<DeliveryOrder[]>([]);
  const [salesOrders, setSalesOrders] = useState<SalesOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<DeliveryOrder | null>(null);
  const [form, setForm] = useState({
    sales_order_id: '',
    delivery_date: '',
    notes: '',
  });
  const [items, setItems] = useState<DeliveryItem[]>([{ ...emptyItem }]);
  const [error, setError] = useState('');

  const load = async () => {
    setError('');
    try {
      const [dRes, oRes] = await Promise.all([
        adminService.listDeliveryOrders(),
        adminService.listSalesOrders(),
      ]);
      setData(dRes.data.data ?? []);
      setSalesOrders((oRes.data.data ?? []).filter((o: SalesOrder) => o.status === 'draft' || o.status === 'partial'));
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
    setForm({ sales_order_id: '', delivery_date: '', notes: '' });
    setItems([{ ...emptyItem }]);
    setShowForm(true);
  };

  const openEdit = (d: DeliveryOrder) => {
    setEditItem(d);
    setForm({
      sales_order_id: String(d.sales_order_id),
      delivery_date: d.delivery_date,
      notes: d.notes || '',
    });
    setItems(
      d.items.length > 0
        ? d.items.map(i => ({
            product_name: i.product_name,
            quantity: String(i.quantity),
            unit: i.unit,
          }))
        : [{ ...emptyItem }]
    );
    setShowForm(true);
  };

  const handleItemChange = (index: number, field: keyof DeliveryItem, value: string) => {
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
        sales_order_id: parseInt(form.sales_order_id),
        delivery_date: form.delivery_date,
        notes: form.notes,
        items: items
          .filter(i => i.product_name)
          .map(i => ({
            product_name: i.product_name,
            quantity: parseFloat(i.quantity) || 0,
            unit: i.unit,
          })),
      };
      if (editItem) {
        await adminService.updateDeliveryOrder(editItem.id, payload);
      } else {
        await adminService.createDeliveryOrder(payload);
      }
      setShowForm(false);
      load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus pengiriman pesanan ini?')) return;
    try {
      await adminService.deleteDeliveryOrder(id);
      load();
    } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Pengiriman Pesanan</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Pengiriman</button>
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Pengiriman Pesanan' : 'Tambah Pengiriman Pesanan'}</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <select value={form.sales_order_id} onChange={e => setForm({ ...form, sales_order_id: e.target.value })} className="border p-2 rounded" required>
                <option value="">-- Pilih Pesanan Penjualan --</option>
                {salesOrders.map(o => (
                  <option key={o.id} value={o.id}>{o.order_number} - {o.customer_name}</option>
                ))}
              </select>
              <input type="date" placeholder="Tanggal Kirim" value={form.delivery_date} onChange={e => setForm({ ...form, delivery_date: e.target.value })} className="border p-2 rounded" required />
            </div>
            <input placeholder="Catatan" value={form.notes} onChange={e => setForm({ ...form, notes: e.target.value })} className="border p-2 rounded w-full" />

            <div>
              <div className="flex justify-between items-center mb-2">
                <h3 className="font-semibold">Item Kirim</h3>
                <button type="button" onClick={addItem} className="text-blue-600 text-sm hover:underline">+ Tambah Item</button>
              </div>
              <table className="w-full border">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="p-2 text-left border">Nama Produk</th>
                    <th className="p-2 text-left border">Jumlah</th>
                    <th className="p-2 text-left border">Satuan</th>
                    <th className="p-2 text-center border w-16">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {items.map((item, idx) => (
                    <tr key={idx}>
                      <td className="p-1 border"><input value={item.product_name} onChange={e => handleItemChange(idx, 'product_name', e.target.value)} className="border p-1 rounded w-full" /></td>
                      <td className="p-1 border"><input type="number" value={item.quantity} onChange={e => handleItemChange(idx, 'quantity', e.target.value)} className="border p-1 rounded w-full" /></td>
                      <td className="p-1 border"><input value={item.unit} onChange={e => handleItemChange(idx, 'unit', e.target.value)} className="border p-1 rounded w-full" placeholder="pcs/kg/box" /></td>
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
                <th className="p-3 text-left">No Pengiriman</th>
                <th className="p-3 text-left">No Pesanan</th>
                <th className="p-3 text-left">Pelanggan</th>
                <th className="p-3 text-left">Tanggal Kirim</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(d => {
                const st = statusConfig[d.status] || statusConfig.shipped;
                return (
                  <tr key={d.id} className="border-t">
                    <td className="p-3 font-medium">{d.do_number}</td>
                    <td className="p-3">{d.order_number}</td>
                    <td className="p-3">{d.customer_name}</td>
                    <td className="p-3">{d.delivery_date}</td>
                    <td className="p-3">
                      <span className={`px-2 py-1 rounded text-xs ${st.className}`}>{st.label}</span>
                    </td>
                    <td className="p-3 text-center gap-2 flex justify-center">
                      <button onClick={() => openEdit(d)} className="text-blue-500 hover:underline">Edit</button>
                      <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                    </td>
                  </tr>
                );
              })}
              {data.length === 0 && <tr><td colSpan={6} className="p-4 text-center text-gray-400">Belum ada data pengiriman pesanan</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}