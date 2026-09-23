import React, { useState, useEffect } from 'react';
import { Plus, Pencil, Trash2, X, MapPin, Phone, Building2, Package } from 'lucide-react';
import { adminService, inventoryService, type Branch, type Warehouse } from '../services/api';
import toast from 'react-hot-toast';

interface BranchForm {
  name: string;
  address: string;
  phone: string;
}

const emptyForm: BranchForm = { name: '', address: '', phone: '' };

export function BranchPage() {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);
  const [form, setForm] = useState<BranchForm>(emptyForm);

  const loadData = async () => {
    setLoading(true);
    try {
      const [brRes, whRes] = await Promise.all([
        adminService.listBranches(),
        inventoryService.listWarehouses()
      ]);
      setBranches(brRes.data.data || []);
      setWarehouses(whRes.data.data || []);
    } catch { 
        setBranches([]); 
        setWarehouses([]);
    }
    setLoading(false);
  };

  useEffect(() => { loadData(); }, []);

  const getWarehouseCount = (branchId: number) => {
      return warehouses.filter(w => w.branch_id === branchId).length;
  };

  const openCreate = () => {
    setEditId(null);
    setForm(emptyForm);
    setShowForm(true);
  };

  const openEdit = (b: Branch) => {
    setEditId(b.id);
    setForm({ name: b.name, address: b.address, phone: b.phone });
    setShowForm(true);
  };

  const handleSave = async () => {
    if (!form.name) { toast.error('Nama cabang wajib diisi'); return; }
    try {
      if (editId) {
        await adminService.updateBranch(editId, form);
        toast.success('Cabang diupdate!');
      } else {
        await adminService.createBranch(form);
        toast.success('Cabang ditambahkan!');
      }
      setShowForm(false);
      loadData();
    } catch (err: any) {
      toast.error(err?.response?.data?.error || 'Gagal menyimpan');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus cabang ini?')) return;
    try {
      await adminService.deleteBranch(id);
      toast.success('Cabang dihapus!');
      loadData();
    } catch { toast.error('Gagal menghapus'); }
  };

  const handleToggleActive = async (b: Branch) => {
    try {
      await adminService.updateBranch(b.id, { is_active: !b.is_active });
      toast.success(b.is_active ? 'Cabang dinonaktifkan' : 'Cabang diaktifkan');
      loadData();
    } catch { toast.error('Gagal update status'); }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Manajemen Cabang</h1>
          <p className="text-gray-500 text-sm mt-1">Kelola daftar cabang outlet</p>
        </div>
        <button onClick={openCreate} className="pos-btn-primary flex items-center gap-2">
          <Plus className="w-4 h-4" /> Tambah Cabang
        </button>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(3)].map((_, i) => (
            <div key={i} className="bg-gray-800/50 rounded-xl h-32 animate-pulse" />
          ))}
        </div>
      ) : branches.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          <Building2 className="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>Belum ada cabang</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {branches.map((b) => (
            <div key={b.id} className="bg-gray-900 border border-gray-800 rounded-xl p-4 hover:border-gray-700 transition-colors">
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-white">
                    <Building2 className="w-5 h-5" />
                  </div>
                  <div>
                    <p className="text-white font-medium">{b.name}</p>
                    <div className="flex items-center gap-1 mt-0.5">
                      <span className={`w-2 h-2 rounded-full ${b.is_active ? 'bg-emerald-500' : 'bg-red-500'}`} />
                      <span className="text-xs text-gray-400">{b.is_active ? 'Aktif' : 'Nonaktif'}</span>
                    </div>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-xs text-gray-500">Gudang</p>
                  <p className="text-xl font-bold text-blue-400">{getWarehouseCount(b.id)}</p>
                </div>
              </div>

              <div className="mt-3 space-y-1.5 text-sm text-gray-400">
                {b.address && (
                  <div className="flex items-center gap-2">
                    <MapPin className="w-3.5 h-3.5 shrink-0" />
                    <span className="truncate">{b.address}</span>
                  </div>
                )}
                {b.phone && (
                  <div className="flex items-center gap-2">
                    <Phone className="w-3.5 h-3.5 shrink-0" />
                    <span>{b.phone}</span>
                  </div>
                )}
                <div className="flex items-center gap-2 text-xs text-gray-500">
                  <Package className="w-3.5 h-3.5" />
                  <span>{getWarehouseCount(b.id)} Gudang</span>
                </div>
              </div>

              <div className="flex gap-1 mt-4 flex-wrap">
                <button onClick={() => openEdit(b)} className="px-3 py-1.5 text-xs bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg flex items-center gap-1">
                  <Pencil className="w-3 h-3" /> Edit
                </button>
                <button onClick={() => handleToggleActive(b)} className={`px-3 py-1.5 text-xs rounded-lg flex items-center gap-1 ${b.is_active ? 'bg-red-900/30 hover:bg-red-900/50 text-red-400' : 'bg-emerald-900/30 hover:bg-emerald-900/50 text-emerald-400'}`}>
                  {b.is_active ? 'Nonaktif' : 'Aktifkan'}
                </button>
                <button onClick={() => handleDelete(b.id)} className="px-3 py-1.5 text-xs bg-red-900/30 hover:bg-red-900/50 text-red-400 rounded-lg flex items-center gap-1">
                  <Trash2 className="w-3 h-3" /> Hapus
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showForm && (
        <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
          <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-md shadow-2xl">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
              <h2 className="font-bold text-white">{editId ? 'Edit Cabang' : 'Tambah Cabang'}</h2>
              <button onClick={() => setShowForm(false)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400">
                <X className="w-5 h-5" />
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm text-gray-300 mb-1">Nama Cabang *</label>
                <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="pos-input" placeholder="Contoh: Cabang Pusat" />
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Alamat</label>
                <textarea value={form.address} onChange={(e) => setForm({ ...form, address: e.target.value })} className="pos-input" rows={2} placeholder="Alamat cabang" />
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Telepon</label>
                <input value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} className="pos-input" placeholder="08123456789" />
              </div>
              <button onClick={handleSave} className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition-colors">
                {editId ? 'Simpan Perubahan' : 'Tambah Cabang'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
