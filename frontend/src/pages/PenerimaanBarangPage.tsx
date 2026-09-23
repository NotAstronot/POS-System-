import React, { useEffect, useState } from 'react';
import { adminService, inventoryService, PurchaseOrder, PurchaseOrderItem, Warehouse } from '../services/api';

export default function PenerimaanBarangPage() {
  const [data, setData] = useState<PurchaseOrder[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState<PurchaseOrder | null>(null);
  const [selectedWarehouse, setSelectedWarehouse] = useState('');
  const [receivedQty, setReceivedQty] = useState<Record<number, number>>({});
  const [submitting, setSubmitting] = useState(false);

  const load = async () => {
    try {
      const [res, wRes] = await Promise.all([
        adminService.listPurchaseOrders(),
        inventoryService.listWarehouses(),
      ]);
      setData(res.data.data ?? []);
      setWarehouses((wRes.data.data ?? []).filter((w: Warehouse) => w.is_active));
    } catch {} finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openReceive = (order: PurchaseOrder) => {
    setSelectedOrder(order);
    setSelectedWarehouse('');
    const qtyMap: Record<number, number> = {};
    order.items.forEach((item) => {
      qtyMap[item.id] = 0;
    });
    setReceivedQty(qtyMap);
    setShowForm(true);
  };

  const handleSave = async () => {
    if (!selectedOrder) return;
    setSubmitting(true);
    try {
      const items = Object.entries(receivedQty)
        .filter(([_, qty]) => qty > 0)
        .map(([id, qty]) => ({ id: parseInt(id), quantity: qty }));
      await adminService.receivePurchaseOrderItems(selectedOrder.id, items, selectedWarehouse ? parseInt(selectedWarehouse) : undefined);
      setShowForm(false);
      load();
    } catch {} finally {
      setSubmitting(false);
    }
  };

  const formatRupiah = (amount: number) =>
    `Rp ${amount.toLocaleString('id-ID')}`;

  const statusLabel = (status: string) => {
    const map: Record<string, string> = {
      draft: 'Draft',
      ordered: 'Dipesan',
      partial: 'Sebagian Diterima',
      received: 'Diterima',
      cancelled: 'Dibatalkan',
    };
    return map[status] || status;
  };

  const statusColor = (status: string) => {
    const map: Record<string, string> = {
      draft: 'bg-gray-100 text-gray-700',
      ordered: 'bg-blue-100 text-blue-700',
      partial: 'bg-yellow-100 text-yellow-700',
      received: 'bg-green-100 text-green-700',
      cancelled: 'bg-red-100 text-red-700',
    };
    return map[status] || 'bg-gray-100 text-gray-700';
  };

  const canReceive = (order: PurchaseOrder) =>
    order.status === 'ordered' || order.status === 'partial';

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Penerimaan Barang</h1>
      </div>

      {showForm && selectedOrder && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">
            Terima Barang - PO {selectedOrder.po_number}
          </h2>
          <p className="text-sm text-gray-600 mb-4">
            Supplier: {selectedOrder.supplier_name}
          </p>
          <div className="mb-4">
            <label className="block text-sm font-medium mb-1">Tujuan Gudang</label>
            <select
              value={selectedWarehouse}
              onChange={(e) => setSelectedWarehouse(e.target.value)}
              className="border p-2 rounded w-full max-w-sm"
            >
              <option value="">-- Gudang Utama (default) --</option>
              {warehouses.map((w) => (
                <option key={w.id} value={w.id}>{w.name}</option>
              ))}
            </select>
          </div>
          <table className="w-full mb-4">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-2 text-left">Produk</th>
                <th className="p-2 text-right">Dipesan</th>
                <th className="p-2 text-right">Sudah Diterima</th>
                <th className="p-2 text-right">Sisa</th>
                <th className="p-2 text-right">Diterima Sekarang</th>
              </tr>
            </thead>
            <tbody>
              {selectedOrder.items.map((item) => {
                const remaining = item.quantity - item.received_qty;
                return (
                  <tr key={item.id} className="border-t">
                    <td className="p-2">{item.product_name}</td>
                    <td className="p-2 text-right">
                      {item.quantity} {item.unit}
                    </td>
                    <td className="p-2 text-right">
                      {item.received_qty} {item.unit}
                    </td>
                    <td className="p-2 text-right font-medium">
                      {remaining} {item.unit}
                    </td>
                    <td className="p-2 text-right">
                      <input
                        type="number"
                        min={0}
                        max={remaining}
                        value={receivedQty[item.id] || 0}
                        onChange={(e) => {
                          const val = Math.min(
                            Math.max(0, parseInt(e.target.value) || 0),
                            remaining
                          );
                          setReceivedQty({ ...receivedQty, [item.id]: val });
                        }}
                        className="border p-1 rounded w-24 text-right"
                      />
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
          <div className="flex gap-2">
            <button
              onClick={handleSave}
              disabled={submitting}
              className="bg-green-600 text-white px-4 py-2 rounded disabled:opacity-50"
            >
              {submitting ? 'Menyimpan...' : 'Simpan'}
            </button>
            <button
              onClick={() => setShowForm(false)}
              className="bg-gray-400 text-white px-4 py-2 rounded"
            >
              Batal
            </button>
          </div>
        </div>
      )}

      {loading ? (
        <p>Loading...</p>
      ) : (
        <div className="bg-white rounded shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-3 text-left">No PO</th>
                <th className="p-3 text-left">Supplier</th>
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-right">Total</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map((po) => (
                <tr key={po.id} className="border-t">
                  <td className="p-3 font-medium">{po.po_number}</td>
                  <td className="p-3">{po.supplier_name}</td>
                  <td className="p-3">{po.order_date}</td>
                  <td className="p-3">
                    <span
                      className={`px-2 py-1 rounded text-xs ${statusColor(po.status)}`}
                    >
                      {statusLabel(po.status)}
                    </span>
                  </td>
                  <td className="p-3 text-right">
                    {formatRupiah(po.total_amount)}
                  </td>
                  <td className="p-3 text-center">
                    {canReceive(po) && (
                      <button
                        onClick={() => openReceive(po)}
                        className="bg-blue-600 text-white px-3 py-1 rounded text-sm hover:bg-blue-700"
                      >
                        Terima Barang
                      </button>
                    )}
                  </td>
                </tr>
              ))}
              {data.length === 0 && (
                <tr>
                  <td colSpan={6} className="p-4 text-center text-gray-400">
                    Belum ada data purchase order
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
