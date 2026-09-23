import React, { useEffect, useState } from 'react';
import toast from 'react-hot-toast';
import {
  adminService, inventoryService, smartlinkService,
  MarketplaceConnection, MarketplaceOrder, Customer, SalesCategory, Item,
} from '../services/api';

const PLATFORMS = ['Tokopedia', 'Shopee', 'Lazada', 'TikTok Shop', 'Bukalapak', 'Blibli'];

export default function SmartLinkECommercePage() {
  const [tab, setTab] = useState<'connections' | 'orders'>('connections');
  const [connections, setConnections] = useState<MarketplaceConnection[]>([]);
  const [orders, setOrders] = useState<MarketplaceOrder[]>([]);
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [categories, setCategories] = useState<SalesCategory[]>([]);
  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);

  const [connForm, setConnForm] = useState({
    platform: 'Tokopedia', shop_name: '', api_token: '', account_name: '',
    customer_id: '' as string, sales_category_id: '' as string,
    shipping_fee_account: 'Biaya Ongkir', commission_account: 'Komisi Platform',
  });
  const [importForm, setImportForm] = useState({
    connection_id: '' as string,
    marketplace_order_id: '', order_date: '', customer_name: '',
    product_id: '' as string, product_name: '', quantity: '', unit_price: '',
    shipping_fee: '', platform_fee: '',
  });

  const load = async () => {
    setLoading(true);
    try {
      const [c, o, cu, ca, it] = await Promise.all([
        smartlinkService.listMarketplaceConnections(),
        smartlinkService.listMarketplaceOrders(),
        adminService.listCustomers(),
        adminService.listSalesCategories(),
        inventoryService.listItems(),
      ]);
      setConnections(c.data.data);
      setOrders(o.data.data);
      setCustomers(cu.data.data);
      setCategories(ca.data.data);
      setItems(it.data.data);
    } catch {
      toast.error('Gagal memuat data SmartLink e-Commerce');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const saveConnection = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await smartlinkService.createMarketplaceConnection({
        ...connForm,
        customer_id: connForm.customer_id ? Number(connForm.customer_id) : null,
        sales_category_id: connForm.sales_category_id ? Number(connForm.sales_category_id) : null,
      });
      toast.success('Koneksi marketplace tersimpan');
      setConnForm({ ...connForm, shop_name: '', api_token: '', account_name: '' });
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal menyimpan koneksi');
    }
  };

  const toggleConnection = async (c: MarketplaceConnection) => {
    try {
      await smartlinkService.updateMarketplaceConnection(c.id, { is_active: !c.is_active });
      load();
    } catch {}
  };

  const deleteConnection = async (id: number) => {
    if (!confirm('Hapus koneksi marketplace ini?')) return;
    try {
      await smartlinkService.deleteMarketplaceConnection(id);
      toast.success('Koneksi dihapus');
      load();
    } catch {}
  };

  const runImport = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!importForm.connection_id) return toast.error('Pilih koneksi marketplace terlebih dahulu');
    const product = items.find((i) => i.id === Number(importForm.product_id));
    try {
      const res = await smartlinkService.importMarketplaceOrders(Number(importForm.connection_id), [{
        marketplace_order_id: importForm.marketplace_order_id,
        order_date: importForm.order_date,
        customer_name: importForm.customer_name,
        product_id: importForm.product_id ? Number(importForm.product_id) : null,
        product_name: importForm.product_name || product?.name || '',
        quantity: parseFloat(importForm.quantity) || 0,
        unit_price: parseFloat(importForm.unit_price) || 0,
        shipping_fee: parseFloat(importForm.shipping_fee) || 0,
        platform_fee: parseFloat(importForm.platform_fee) || 0,
      }]);
      toast.success(`${res.data.data.imported} pesanan berhasil diimpor`);
      setImportForm({
        connection_id: importForm.connection_id, marketplace_order_id: '', order_date: '',
        customer_name: '', product_id: '', product_name: '', quantity: '', unit_price: '',
        shipping_fee: '', platform_fee: '',
      });
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal impor pesanan');
    }
  };

  const postOrder = async (id: number) => {
    try {
      await smartlinkService.postMarketplaceOrder(id);
      toast.success('Pesanan diposting ke Pesanan Penjualan');
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal memposting pesanan');
    }
  };

  const deleteOrder = async (id: number) => {
    if (!confirm('Hapus pesanan marketplace ini?')) return;
    try {
      await smartlinkService.deleteMarketplaceOrder(id);
      load();
    } catch {}
  };

  const statusBadge = (status: string) => {
    const map: Record<string, string> = {
      imported: 'bg-blue-100 text-blue-700',
      posted: 'bg-green-100 text-green-700',
    };
    return <span className={`px-2 py-1 rounded text-xs ${map[status] || 'bg-yellow-100 text-yellow-700'}`}>{status === 'posted' ? 'Terposting' : 'Terimpor'}</span>;
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-2">
        <h1 className="text-2xl font-bold">SmartLink e-Commerce</h1>
      </div>
      <p className="text-sm text-gray-500 mb-6">
        Integrasi marketplace (Tokopedia, Shopee, Lazada, TikTok Shop, dll). Impor transaksi penjualan,
        ongkir, dan komisi platform otomatis ke modul Penjualan.
      </p>

      <div className="flex gap-2 mb-6">
        {([['connections', 'Koneksi Marketplace'], ['orders', 'Pesanan Marketplace']] as const).map(([key, label]) => (
          <button
            key={key}
            onClick={() => setTab(key)}
            className={`px-4 py-2 rounded font-medium ${tab === key ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
          >
            {label}
          </button>
        ))}
      </div>

      {tab === 'connections' && (
        <>
          <div className="bg-white p-4 rounded shadow mb-6">
            <h2 className="font-bold mb-4">Tambah Koneksi Marketplace</h2>
            <form onSubmit={saveConnection} className="grid grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">Platform</label>
                <select value={connForm.platform} onChange={(e) => setConnForm({ ...connForm, platform: e.target.value })} className="border p-2 rounded w-full" required>
                  {PLATFORMS.map((p) => <option key={p} value={p}>{p}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Nama Toko</label>
                <input placeholder="Nama Toko Marketplace" value={connForm.shop_name} onChange={(e) => setConnForm({ ...connForm, shop_name: e.target.value })} className="border p-2 rounded w-full" required />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">API Token</label>
                <input placeholder="Token / API Key" value={connForm.api_token} onChange={(e) => setConnForm({ ...connForm, api_token: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Rekening Penampung</label>
                <input placeholder="cth: BCA - Toko Utama" value={connForm.account_name} onChange={(e) => setConnForm({ ...connForm, account_name: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Pelanggan (pemetaan SO)</label>
                <select value={connForm.customer_id} onChange={(e) => setConnForm({ ...connForm, customer_id: e.target.value })} className="border p-2 rounded w-full">
                  <option value="">-- Pilih Pelanggan --</option>
                  {customers.map((cu) => <option key={cu.id} value={cu.id}>{cu.name}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Kategori Penjualan</label>
                <select value={connForm.sales_category_id} onChange={(e) => setConnForm({ ...connForm, sales_category_id: e.target.value })} className="border p-2 rounded w-full">
                  <option value="">-- Pilih Kategori --</option>
                  {categories.map((ca) => <option key={ca.id} value={ca.id}>{ca.name}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Akun Biaya Ongkir</label>
                <input value={connForm.shipping_fee_account} onChange={(e) => setConnForm({ ...connForm, shipping_fee_account: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Akun Komisi Platform</label>
                <input value={connForm.commission_account} onChange={(e) => setConnForm({ ...connForm, commission_account: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              <div className="flex items-end">
                <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded w-full">Simpan Koneksi</button>
              </div>
            </form>
          </div>

          {loading ? <p>Loading...</p> : (
            <div className="bg-white rounded shadow overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="p-3 text-left">Platform</th>
                    <th className="p-3 text-left">Toko</th>
                    <th className="p-3 text-left">Pelanggan Terpetakan</th>
                    <th className="p-3 text-left">Rekening</th>
                    <th className="p-3 text-left">Sync Terakhir</th>
                    <th className="p-3 text-center">Status</th>
                    <th className="p-3 text-center">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {connections.map((c) => (
                    <tr key={c.id} className="border-t">
                      <td className="p-3 font-medium">{c.platform}</td>
                      <td className="p-3">{c.shop_name}</td>
                      <td className="p-3">{c.customer_name || <span className="text-gray-400">Belum dipetakan</span>}</td>
                      <td className="p-3 text-sm">{c.account_name || '-'}</td>
                      <td className="p-3 text-sm">{c.last_sync_at ? new Date(c.last_sync_at).toLocaleString('id-ID') : 'Belum pernah sync'}</td>
                      <td className="p-3 text-center">
                        <span className={`px-2 py-1 rounded text-xs ${c.is_active ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-500'}`}>{c.is_active ? 'Aktif' : 'Nonaktif'}</span>
                      </td>
                      <td className="p-3 text-center gap-2 flex justify-center">
                        <button onClick={() => toggleConnection(c)} className="text-blue-500 hover:underline">{c.is_active ? 'Nonaktifkan' : 'Aktifkan'}</button>
                        <button onClick={() => deleteConnection(c.id)} className="text-red-500 hover:underline">Hapus</button>
                      </td>
                    </tr>
                  ))}
                  {connections.length === 0 && (
                    <tr><td colSpan={7} className="p-4 text-center text-gray-400">Belum ada koneksi marketplace</td></tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}

      {tab === 'orders' && (
        <>
          <div className="bg-white p-4 rounded shadow mb-6">
            <h2 className="font-bold mb-4">Impor Pesanan Marketplace</h2>
            <form onSubmit={runImport} className="grid grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">Koneksi Marketplace</label>
                <select value={importForm.connection_id} onChange={(e) => setImportForm({ ...importForm, connection_id: e.target.value })} className="border p-2 rounded w-full" required>
                  <option value="">-- Pilih Koneksi --</option>
                  {connections.map((c) => <option key={c.id} value={c.id}>{c.platform} - {c.shop_name}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">ID Pesanan Marketplace</label>
                <input placeholder="cth: INV/2026/0813-123" value={importForm.marketplace_order_id} onChange={(e) => setImportForm({ ...importForm, marketplace_order_id: e.target.value })} className="border p-2 rounded w-full" required />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Tanggal Pesanan</label>
                <input type="date" value={importForm.order_date} onChange={(e) => setImportForm({ ...importForm, order_date: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Nama Pembeli</label>
                <input placeholder="Nama pembeli marketplace" value={importForm.customer_name} onChange={(e) => setImportForm({ ...importForm, customer_name: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Produk (dari Barang & Jasa)</label>
                <select value={importForm.product_id} onChange={(e) => { const it = items.find((i) => i.id === Number(e.target.value)); setImportForm({ ...importForm, product_id: e.target.value, product_name: it?.name || '' }); }} className="border p-2 rounded w-full">
                  <option value="">-- Pilih Produk --</option>
                  {items.map((it) => <option key={it.id} value={it.id}>{it.name}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Nama Produk (tampilan)</label>
                <input value={importForm.product_name} onChange={(e) => setImportForm({ ...importForm, product_name: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Qty</label>
                <input type="number" placeholder="1" value={importForm.quantity} onChange={(e) => setImportForm({ ...importForm, quantity: e.target.value })} className="border p-2 rounded w-full" required />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Harga Satuan (Rp)</label>
                <input type="number" placeholder="0" value={importForm.unit_price} onChange={(e) => setImportForm({ ...importForm, unit_price: e.target.value })} className="border p-2 rounded w-full" required />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Ongkir (Rp)</label>
                <input type="number" placeholder="0" value={importForm.shipping_fee} onChange={(e) => setImportForm({ ...importForm, shipping_fee: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Komisi Platform (Rp)</label>
                <input type="number" placeholder="0" value={importForm.platform_fee} onChange={(e) => setImportForm({ ...importForm, platform_fee: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              <div className="flex items-end">
                <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded w-full">Impor Pesanan</button>
              </div>
            </form>
          </div>

          {loading ? <p>Loading...</p> : (
            <div className="bg-white rounded shadow overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="p-3 text-left">Pesanan</th>
                    <th className="p-3 text-left">Marketplace</th>
                    <th className="p-3 text-left">Produk</th>
                    <th className="p-3 text-right">Qty</th>
                    <th className="p-3 text-right">Subtotal</th>
                    <th className="p-3 text-right">Ongkir</th>
                    <th className="p-3 text-right">Komisi</th>
                    <th className="p-3 text-right">Total</th>
                    <th className="p-3 text-left">Status</th>
                    <th className="p-3 text-center">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {orders.map((o) => (
                    <tr key={o.id} className="border-t">
                      <td className="p-3">
                        <p className="font-medium">{o.marketplace_order_id}</p>
                        <p className="text-xs text-gray-500">{o.order_date}{o.customer_name ? ` • ${o.customer_name}` : ''}</p>
                      </td>
                      <td className="p-3">{o.platform} <span className="text-xs text-gray-500">({o.shop_name})</span></td>
                      <td className="p-3">{o.product_name}</td>
                      <td className="p-3 text-right">{o.quantity}</td>
                      <td className="p-3 text-right">Rp {o.subtotal.toLocaleString()}</td>
                      <td className="p-3 text-right">Rp {o.shipping_fee.toLocaleString()}</td>
                      <td className="p-3 text-right">Rp {o.platform_fee.toLocaleString()}</td>
                      <td className="p-3 text-right font-bold">Rp {o.grand_total.toLocaleString()}</td>
                      <td className="p-3">{statusBadge(o.status)}
                        {o.sales_order_number && <p className="text-xs text-gray-500 mt-1">SO: {o.sales_order_number}</p>}
                      </td>
                      <td className="p-3 text-center gap-2 flex justify-center">
                        {o.status === 'imported' && (
                          <button onClick={() => postOrder(o.id)} className="text-green-600 hover:underline font-medium">Posting ke Penjualan</button>
                        )}
                        <button onClick={() => deleteOrder(o.id)} className="text-red-500 hover:underline">Hapus</button>
                      </td>
                    </tr>
                  ))}
                  {orders.length === 0 && (
                    <tr><td colSpan={10} className="p-4 text-center text-gray-400">Belum ada pesanan marketplace. Gunakan form impor di atas.</td></tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}
    </div>
  );
}