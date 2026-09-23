import React, { useEffect, useState } from 'react';
import { inventoryService, Warehouse, adminService, Branch } from '../services/api';

export default function WarehousePage() {
  const [data, setData] = useState<Warehouse[]>([]);
  const [branches, setBranches] = useState<Branch[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<Warehouse | null>(null);
  const [form, setForm] = useState({ code: '', name: '', branch_id: '', address: '' });
  const [error, setError] = useState('');

  const load = async () => {
    setError('');
    try {
      const [wRes, bRes] = await Promise.all([
        inventoryService.listWarehouses(),
        adminService.listBranches(),
      ]);
      setData(wRes.data.data ?? []);
      setBranches(bRes.data.data ?? []);
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
    setForm({ code: '', name: '', branch_id: '', address: '' });
    setShowForm(true);
  };

  const openEdit = (w: Warehouse) => {
    setEditItem(w);
    setForm({
      code: w.code,
      name: w.name,
      branch_id: w.branch_id ? String(w.branch_id) : '',
      address: w.address,
    });
    setShowForm(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload = {
        code: form.code,
        name: form.name,
        branch_id: form.branch_id ? parseInt(form.branch_id) : null,
        address: form.address,
      };
      if (editItem) {
        await inventoryService.updateWarehouse(editItem.id, payload);
      } else {
        await inventoryService.createWarehouse(payload);
      }
      setShowForm(false);
      load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal menyimpan data');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus gudang ini? Stok di gudang ini juga akan terhapus.')) return;
    try {
      await inventoryService.deleteWarehouse(id);
      load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal menghapus data');
    }
  };

  const handleToggleActive = async (w: Warehouse) => {
    try {
      await inventoryService.updateWarehouse(w.id, { is_active: !w.is_active });
      load();
    } catch {}
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Gudang</h1>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Gudang</button>
      </div>

      {error && <div className="mb-4 p-4 text-red-500 bg-red-50 border border-red-200 rounded">{error}</div>}

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6 border border-gray-200">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Gudang' : 'Tambah Gudang'}</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Kode Gudang</label>
              <input value={form.code} onChange={e => setForm({ ...form, code: e.target.value })} className="border p-2 rounded w-full" placeholder="WH-002" required />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Nama Gudang</label>
              <input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} className="border p-2 rounded w-full" placeholder="Gudang Utama" required />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Cabang</label>
              <select value={form.branch_id} onChange={e => setForm({ ...form, branch_id: e.target.value })} className="border p-2 rounded w-full">
                <option value="">-- Tanpa Cabang --</option>
                {branches.filter(b => b.is_active).map(b => (
                  <option key={b.id} value={b.id}>{b.name}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Alamat</label>
              <input value={form.address} onChange={e => setForm({ ...form, address: e.target.value })} className="border p-2 rounded w-full" placeholder="Alamat gudang" />
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
                <th className="p-3 text-left">Kode</th>
                <th className="p-3 text-left">Nama Gudang</th>
                <th className="p-3 text-left">Cabang</th>
                <th className="p-3 text-left">Alamat</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(w => (
                <tr key={w.id} className="border-t">
                  <td className="p-3 font-medium">{w.code}</td>
                  <td className="p-3">{w.name}</td>
                  <td className="p-3">{w.branch_name || '-'}</td>
                  <td className="p-3">{w.address || '-'}</td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs ${w.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                      {w.is_active ? 'Aktif' : 'Nonaktif'}
                    </span>
                  </td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => openEdit(w)} className="text-blue-500 hover:underline">Edit</button>
                    <button onClick={() => handleToggleActive(w)} className="text-yellow-600 hover:underline">{w.is_active ? 'Nonaktifkan' : 'Aktifkan'}</button>
                    <button onClick={() => handleDelete(w.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && <tr><td colSpan={6} className="p-4 text-center text-gray-400">Belum ada gudang</td></tr>}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}