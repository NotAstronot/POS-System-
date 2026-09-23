import React, { useEffect, useState } from 'react';
import toast from 'react-hot-toast';
import {
  smartlinkService, inventoryService,
  BillOfMaterials, BomItem, WorkOrder, Item, Warehouse,
} from '../services/api';

export default function SmartLinkManufacturingPage() {
  const [tab, setTab] = useState<'bom' | 'workorders'>('bom');
  const [boms, setBoms] = useState<BillOfMaterials[]>([]);
  const [workOrders, setWorkOrders] = useState<WorkOrder[]>([]);
  const [items, setItems] = useState<Item[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [loading, setLoading] = useState(true);

  const [bomForm, setBomForm] = useState({
    product_id: '' as string, quantity_output: '1', unit: 'pcs',
    overhead_cost: '0', notes: '',
  });
  const [bomItems, setBomItems] = useState<BomItem[]>([]);
  const [editingBom, setEditingBom] = useState<number | null>(null);

  const [woForm, setWoForm] = useState({
    bom_id: '' as string, quantity: '', warehouse_id: '' as string,
    scheduled_date: '', notes: '',
  });
  const [woPreview, setWoPreview] = useState<BillOfMaterials | null>(null);

  const load = async () => {
    setLoading(true);
    try {
      const [b, w, i, wh] = await Promise.all([
        smartlinkService.listBoms(),
        smartlinkService.listWorkOrders(),
        inventoryService.listItems(),
        inventoryService.listWarehouses(),
      ]);
      setBoms(b.data.data);
      setWorkOrders(w.data.data);
      setItems(i.data.data);
      setWarehouses(wh.data.data);
    } catch {
      toast.error('Gagal memuat data manufaktur');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const addBomItem = () => {
    setBomItems([...bomItems, { product_id: 0, product_name: '', quantity_required: 1, unit: 'pcs', cost_per_unit: 0, estimated_cost: 0 }]);
  };

  const updateBomItem = (index: number, field: keyof BomItem, value: any) => {
    const next = [...bomItems];
    if (field === 'product_id') {
      const it = items.find((x) => x.id === Number(value));
      next[index] = { ...next[index], product_id: Number(value), product_name: it?.name || '', unit: it?.unit || 'pcs', cost_per_unit: it?.purchase_price || 0 };
    } else {
      (next[index] as any)[field] = value;
    }
    next[index].estimated_cost = (parseFloat(String(next[index].quantity_required)) || 0) * (parseFloat(String(next[index].cost_per_unit)) || 0);
    setBomItems(next);
  };

  const removeBomItem = (index: number) => setBomItems(bomItems.filter((_, i) => i !== index));

  const resetBomForm = () => {
    setBomForm({ product_id: '', quantity_output: '1', unit: 'pcs', overhead_cost: '0', notes: '' });
    setBomItems([]);
    setEditingBom(null);
  };

  const openEditBom = (bom: BillOfMaterials) => {
    setEditingBom(bom.id);
    setBomForm({
      product_id: String(bom.product_id), quantity_output: String(bom.quantity_output),
      unit: bom.unit, overhead_cost: String(bom.overhead_cost), notes: bom.notes,
    });
    setBomItems(bom.items.map((i) => ({ ...i })));
    setTab('bom');
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const saveBom = async (e: React.FormEvent) => {
    e.preventDefault();
    if (bomItems.length === 0) return toast.error('Tambah minimal 1 bahan baku');
    const payload = {
      product_id: Number(bomForm.product_id),
      quantity_output: parseFloat(bomForm.quantity_output) || 1,
      unit: bomForm.unit,
      overhead_cost: parseFloat(bomForm.overhead_cost) || 0,
      notes: bomForm.notes,
      items: bomItems.map((i) => ({
        product_id: i.product_id,
        quantity_required: parseFloat(String(i.quantity_required)) || 0,
        unit: i.unit,
        cost_per_unit: parseFloat(String(i.cost_per_unit)) || 0,
      })),
    };
    try {
      if (editingBom) {
        await smartlinkService.updateBom(editingBom, payload);
        toast.success('BOM berhasil diperbarui');
      } else {
        const res = await smartlinkService.createBom(payload);
        toast.success(`BOM ${res.data.data.bom_number} dibuat, HPP: Rp ${res.data.data.cost_per_unit.toLocaleString()}`);
      }
      resetBomForm();
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal menyimpan BOM');
    }
  };

  const recalculateBom = async (id: number) => {
    try {
      const res = await smartlinkService.recalculateBom(id);
      toast.success(`HPP dihitung ulang dari harga bahan baku terbaru: Rp ${res.data.data.cost_per_unit.toLocaleString()}`);
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal menghitung ulang HPP');
    }
  };

  const deleteBom = async (id: number) => {
    if (!confirm('Hapus BOM ini?')) return;
    try {
      await smartlinkService.deleteBom(id);
      load();
    } catch {}
  };

  const selectBomForWo = async (bomId: string) => {
    setWoForm({ ...woForm, bom_id: bomId });
    if (!bomId) return setWoPreview(null);
    try {
      const res = await smartlinkService.getBom(Number(bomId));
      setWoPreview(res.data.data);
    } catch {
      setWoPreview(null);
    }
  };

  const createWorkOrder = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await smartlinkService.createWorkOrder({
        bom_id: Number(woForm.bom_id),
        quantity: parseFloat(woForm.quantity) || 0,
        warehouse_id: woForm.warehouse_id ? Number(woForm.warehouse_id) : null,
        scheduled_date: woForm.scheduled_date,
        notes: woForm.notes,
      });
      toast.success(`Perintah Kerja ${res.data.data.work_order_number} dibuat`);
      setWoForm({ bom_id: '', quantity: '', warehouse_id: '', scheduled_date: '', notes: '' });
      setWoPreview(null);
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal membuat perintah kerja');
    }
  };

  const startWo = async (id: number) => {
    try {
      await smartlinkService.startWorkOrder(id);
      toast.success('Perintah kerja dimulai');
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal memulai perintah kerja');
    }
  };

  const completeWo = async (id: number) => {
    if (!confirm('Selesaikan produksi? Bahan baku akan dikurangi dan barang jadi ditambahkan ke gudang.')) return;
    try {
      await smartlinkService.completeWorkOrder(id);
      toast.success('Produksi selesai - stok diperbarui otomatis');
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal menyelesaikan produksi');
    }
  };

  const cancelWo = async (id: number) => {
    if (!confirm('Batalkan perintah kerja ini?')) return;
    try {
      await smartlinkService.cancelWorkOrder(id);
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal membatalkan');
    }
  };

  const deleteWo = async (id: number) => {
    if (!confirm('Hapus perintah kerja ini?')) return;
    try {
      await smartlinkService.deleteWorkOrder(id);
      load();
    } catch {}
  };

  const woStatusBadge = (s: string) => {
    const map: Record<string, string> = {
      planned: 'bg-gray-100 text-gray-600',
      in_progress: 'bg-blue-100 text-blue-700',
      completed: 'bg-green-100 text-green-700',
      cancelled: 'bg-red-100 text-red-600',
    };
    const label: Record<string, string> = {
      planned: 'Direncanakan', in_progress: 'Berjalan', completed: 'Selesai', cancelled: 'Dibatalkan',
    };
    return <span className={`px-2 py-1 rounded text-xs ${map[s] || 'bg-gray-100 text-gray-600'}`}>{label[s] || s}</span>;
  };

  const bomItemCost = (items: BomItem[]) => items.reduce((sum, i) => sum + i.estimated_cost, 0);

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-2">
        <h1 className="text-2xl font-bold">Modul Manufaktur / Produksi</h1>
      </div>
      <p className="text-sm text-gray-500 mb-6">
        Untuk bisnis non-retail / pabrikasi. Kelola <b>Bill of Materials (BOM)</b>, <b>Perintah Kerja (Work Order)</b>,
        dan perhitungan <b>HPP barang jadi</b> yang terhubung dengan stok gudang.
      </p>

      <div className="flex gap-2 mb-6">
        {([['bom', 'Bill of Materials'], ['workorders', 'Perintah Kerja']] as const).map(([key, label]) => (
          <button
            key={key}
            onClick={() => setTab(key)}
            className={`px-4 py-2 rounded font-medium ${tab === key ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-600 hover:bg-gray-200'}`}
          >
            {label}
          </button>
        ))}
      </div>

      {tab === 'bom' && (
        <>
          <div className="bg-white p-4 rounded shadow mb-6">
            <h2 className="font-bold mb-4">{editingBom ? `Edit BOM #${editingBom}` : 'Tambah BOM (Bill of Materials)'}</h2>
            <form onSubmit={saveBom} className="space-y-4">
              <div className="grid grid-cols-4 gap-4">
                <div>
                  <label className="block text-sm font-medium mb-1">Barang Jadi</label>
                  <select value={bomForm.product_id} onChange={(e) => setBomForm({ ...bomForm, product_id: e.target.value })} className="border p-2 rounded w-full" required>
                    <option value="">-- Pilih Barang Jadi --</option>
                    {items.map((it) => <option key={it.id} value={it.id}>{it.name}</option>)}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">Jumlah Output</label>
                  <input type="number" step="0.01" value={bomForm.quantity_output} onChange={(e) => setBomForm({ ...bomForm, quantity_output: e.target.value })} className="border p-2 rounded w-full" />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">Satuan</label>
                  <input value={bomForm.unit} onChange={(e) => setBomForm({ ...bomForm, unit: e.target.value })} className="border p-2 rounded w-full" />
                </div>
                <div>
                  <label className="block text-sm font-medium mb-1">Biaya Overhead (Rp)</label>
                  <input type="number" step="0.01" value={bomForm.overhead_cost} onChange={(e) => setBomForm({ ...bomForm, overhead_cost: e.target.value })} className="border p-2 rounded w-full" />
                </div>
              </div>

              <div>
                <div className="flex justify-between items-center mb-2">
                  <h3 className="font-medium">Bahan Baku (dari Barang & Jasa)</h3>
                  <button type="button" onClick={addBomItem} className="text-blue-600 text-sm hover:underline">+ Tambah Bahan Baku</button>
                </div>
                {bomItems.length > 0 && (
                  <div className="border rounded overflow-hidden">
                    <table className="w-full text-sm">
                      <thead className="bg-gray-50">
                        <tr>
                          <th className="p-2 text-left">Bahan Baku</th>
                          <th className="p-2 text-left">Satuan</th>
                          <th className="p-2 text-right">Jumlah Kebutuhan</th>
                          <th className="p-2 text-right">Biaya per Unit (Rp)</th>
                          <th className="p-2 text-right">Estimasi Biaya</th>
                          <th className="p-2 text-center">Aksi</th>
                        </tr>
                      </thead>
                      <tbody>
                        {bomItems.map((bi, idx) => (
                          <tr key={idx} className="border-t">
                            <td className="p-1">
                              <select value={bi.product_id || ''} onChange={(e) => updateBomItem(idx, 'product_id', e.target.value)} className="border p-1 rounded w-full text-sm" required>
                                <option value="">-- Pilih Bahan Baku --</option>
                                {items.map((it) => <option key={it.id} value={it.id}>{it.name}</option>)}
                              </select>
                            </td>
                            <td className="p-1">
                              <input value={bi.unit} onChange={(e) => updateBomItem(idx, 'unit', e.target.value)} className="border p-1 rounded w-full text-sm" />
                            </td>
                            <td className="p-1">
                              <input type="number" step="0.01" value={bi.quantity_required} onChange={(e) => updateBomItem(idx, 'quantity_required', e.target.value)} className="border p-1 rounded w-full text-sm text-right" />
                            </td>
                            <td className="p-1">
                              <input type="number" step="0.01" value={bi.cost_per_unit} onChange={(e) => updateBomItem(idx, 'cost_per_unit', e.target.value)} className="border p-1 rounded w-full text-sm text-right" />
                            </td>
                            <td className="p-1 text-right font-medium">Rp {bi.estimated_cost.toLocaleString()}</td>
                            <td className="p-1 text-center">
                              <button type="button" onClick={() => removeBomItem(idx)} className="text-red-500 hover:underline text-sm">Hapus</button>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                    <div className="p-2 bg-gray-50 text-right text-sm">
                      <span className="font-medium">Total Bahan Baku:</span> Rp {bomItemCost(bomItems).toLocaleString()}
                      {parseFloat(bomForm.overhead_cost) > 0 && (
                        <span className="ml-3 text-gray-500">+ Overhead Rp {(parseFloat(bomForm.overhead_cost) || 0).toLocaleString()}</span>
                      )}
                      {parseFloat(bomForm.quantity_output) > 0 && (
                        <span className="ml-3 font-bold text-blue-700">
                          HPP per unit: Rp {((bomItemCost(bomItems) + (parseFloat(bomForm.overhead_cost) || 0)) / (parseFloat(bomForm.quantity_output) || 1)).toLocaleString()}
                        </span>
                      )}
                    </div>
                  </div>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium mb-1">Catatan</label>
                <input value={bomForm.notes} onChange={(e) => setBomForm({ ...bomForm, notes: e.target.value })} className="border p-2 rounded w-full" />
              </div>

              <div className="flex gap-2">
                <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">{editingBom ? 'Simpan Perubahan' : 'Simpan BOM'}</button>
                {editingBom && <button type="button" onClick={resetBomForm} className="bg-gray-400 text-white px-4 py-2 rounded">Batal</button>}
              </div>
            </form>
          </div>

          {loading ? <p>Loading...</p> : (
            <div className="bg-white rounded shadow overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="p-3 text-left">No. BOM</th>
                    <th className="p-3 text-left">Barang Jadi</th>
                    <th className="p-3 text-left">Output</th>
                    <th className="p-3 text-right">Total Bahan Baku</th>
                    <th className="p-3 text-right">HPP / Unit</th>
                    <th className="p-3 text-left">Status</th>
                    <th className="p-3 text-center">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {boms.map((b) => (
                    <tr key={b.id} className="border-t">
                      <td className="p-3 font-medium">{b.bom_number}</td>
                      <td className="p-3">
                        {b.product_name}
                        <div className="text-xs text-gray-500 mt-1">{b.items.map((i) => i.product_name).join(', ')}</div>
                      </td>
                      <td className="p-3">{b.quantity_output} {b.unit}</td>
                      <td className="p-3 text-right">Rp {bomItemCost(b.items).toLocaleString()}</td>
                      <td className="p-3 text-right font-bold">Rp {b.cost_per_unit.toLocaleString()}</td>
                      <td className="p-3">
                        <span className="px-2 py-1 rounded text-xs bg-green-100 text-green-700 capitalize">{b.status}</span>
                      </td>
                      <td className="p-3 text-center gap-2 flex justify-center">
                        <button onClick={() => openEditBom(b)} className="text-blue-500 hover:underline">Edit</button>
                        <button onClick={() => recalculateBom(b.id)} className="text-indigo-500 hover:underline">Hitung Ulang HPP</button>
                        <button onClick={() => deleteBom(b.id)} className="text-red-500 hover:underline">Hapus</button>
                      </td>
                    </tr>
                  ))}
                  {boms.length === 0 && (
                    <tr><td colSpan={7} className="p-4 text-center text-gray-400">Belum ada BOM</td></tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}

      {tab === 'workorders' && (
        <>
          <div className="bg-white p-4 rounded shadow mb-6">
            <h2 className="font-bold mb-4">Buat Perintah Kerja (Work Order)</h2>
            <form onSubmit={createWorkOrder} className="grid grid-cols-4 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">BOM</label>
                <select value={woForm.bom_id} onChange={(e) => selectBomForWo(e.target.value)} className="border p-2 rounded w-full" required>
                  <option value="">-- Pilih BOM --</option>
                  {boms.map((b) => <option key={b.id} value={b.id}>{b.bom_number} - {b.product_name}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Jumlah Produksi</label>
                <input type="number" step="0.01" placeholder="cth: 100" value={woForm.quantity} onChange={(e) => setWoForm({ ...woForm, quantity: e.target.value })} className="border p-2 rounded w-full" required />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Gudang Produksi</label>
                <select value={woForm.warehouse_id} onChange={(e) => setWoForm({ ...woForm, warehouse_id: e.target.value })} className="border p-2 rounded w-full">
                  <option value="">-- Pilih Gudang --</option>
                  {warehouses.map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Tanggal Jadwal</label>
                <input type="date" value={woForm.scheduled_date} onChange={(e) => setWoForm({ ...woForm, scheduled_date: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              <div className="col-span-4">
                <label className="block text-sm font-medium mb-1">Catatan</label>
                <input value={woForm.notes} onChange={(e) => setWoForm({ ...woForm, notes: e.target.value })} className="border p-2 rounded w-full" />
              </div>
              {woPreview && (
                <div className="col-span-4 bg-blue-50 border border-blue-200 rounded p-3 text-sm">
                  <p className="font-medium text-blue-800 mb-2">Kebutuhan bahan baku untuk {woPreview.product_name} ({woForm.quantity || 0} unit):</p>
                  <ul className="grid grid-cols-2 gap-1 text-blue-700">
                    {woPreview.items.map((i, idx) => (
                      <li key={idx}>• {i.product_name}: {(i.quantity_required * (parseFloat(woForm.quantity) || 0) / (woPreview.quantity_output || 1)).toFixed(2)} {i.unit}</li>
                    ))}
                  </ul>
                </div>
              )}
              <div className="col-span-4 flex gap-2">
                <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Buat Perintah Kerja</button>
              </div>
            </form>
          </div>

          {loading ? <p>Loading...</p> : (
            <div className="bg-white rounded shadow overflow-hidden">
              <table className="w-full">
                <thead className="bg-gray-100">
                  <tr>
                    <th className="p-3 text-left">No. WO</th>
                    <th className="p-3 text-left">Produk</th>
                    <th className="p-3 text-right">Qty</th>
                    <th className="p-3 text-left">Gudang</th>
                    <th className="p-3 text-left">Jadwal</th>
                    <th className="p-3 text-right">Biaya Aktual (HPP)</th>
                    <th className="p-3 text-left">Status</th>
                    <th className="p-3 text-center">Aksi</th>
                  </tr>
                </thead>
                <tbody>
                  {workOrders.map((w) => (
                    <tr key={w.id} className="border-t">
                      <td className="p-3 font-medium">{w.work_order_number}</td>
                      <td className="p-3">
                        {w.product_name}
                        <div className="text-xs text-gray-500 mt-1">BOM: {w.bom_number}</div>
                      </td>
                      <td className="p-3 text-right">{w.quantity}</td>
                      <td className="p-3">{w.warehouse_name || '-'}</td>
                      <td className="p-3 text-sm">{w.scheduled_date}</td>
                      <td className="p-3 text-right font-medium">Rp {w.actual_cost.toLocaleString()}</td>
                      <td className="p-3">{woStatusBadge(w.status)}</td>
                      <td className="p-3 text-center gap-2 flex justify-center">
                        {w.status === 'planned' && (
                          <button onClick={() => startWo(w.id)} className="text-blue-600 hover:underline font-medium">Mulai</button>
                        )}
                        {(w.status === 'planned' || w.status === 'in_progress') && (
                          <button onClick={() => completeWo(w.id)} className="text-green-600 hover:underline font-medium">Selesaikan</button>
                        )}
                        {(w.status === 'planned' || w.status === 'in_progress') && (
                          <button onClick={() => cancelWo(w.id)} className="text-yellow-600 hover:underline">Batalkan</button>
                        )}
                        <button onClick={() => deleteWo(w.id)} className="text-red-500 hover:underline">Hapus</button>
                      </td>
                    </tr>
                  ))}
                  {workOrders.length === 0 && (
                    <tr><td colSpan={8} className="p-4 text-center text-gray-400">Belum ada perintah kerja</td></tr>
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