import React, { useState, useEffect } from 'react';
import { Plus, Pencil, Trash2, X, Building } from 'lucide-react';
import { adminService, type Department } from '../services/api';
import toast from 'react-hot-toast';

interface DeptForm { name: string; description: string }
const emptyForm: DeptForm = { name: '', description: '' };

export function DepartmentPage() {
  const [items, setItems] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);
  const [form, setForm] = useState<DeptForm>(emptyForm);

  const load = async () => {
    try {
      const res = await adminService.listDepartments();
      setItems(res.data.data || []);
    } catch { setItems([]); }
    setLoading(false);
  };

  useEffect(() => { load(); }, []);

  const openCreate = () => { setEditId(null); setForm(emptyForm); setShowForm(true); };
  const openEdit = (d: Department) => { setEditId(d.id); setForm({ name: d.name, description: d.description }); setShowForm(true); };

  const handleSave = async () => {
    if (!form.name) { toast.error('Nama departemen wajib diisi'); return; }
    try {
      if (editId) {
        await adminService.updateDepartment(editId, form);
        toast.success('Departemen diupdate!');
      } else {
        await adminService.createDepartment(form);
        toast.success('Departemen ditambahkan!');
      }
      setShowForm(false);
      load();
    } catch (err: any) { toast.error(err?.response?.data?.error || 'Gagal menyimpan'); }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus departemen ini?')) return;
    try { await adminService.deleteDepartment(id); toast.success('Departemen dihapus!'); load(); }
    catch { toast.error('Gagal menghapus'); }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Departemen</h1>
          <p className="text-gray-500 text-sm mt-1">Kelola daftar departemen</p>
        </div>
        <button onClick={openCreate} className="pos-btn-primary flex items-center gap-2">
          <Plus className="w-4 h-4" /> Tambah Departemen
        </button>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(3)].map((_, i) => <div key={i} className="bg-gray-800/50 rounded-xl h-24 animate-pulse" />)}
        </div>
      ) : items.length === 0 ? (
        <div className="text-center py-12 text-gray-500"><Building className="w-12 h-12 mx-auto mb-3 opacity-50" /><p>Belum ada departemen</p></div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {items.map((d) => (
            <div key={d.id} className="bg-gray-900 border border-gray-800 rounded-xl p-4">
              <div className="flex items-start justify-between">
                <div>
                  <p className="text-white font-medium">{d.name}</p>
                  <p className="text-sm text-gray-400 mt-1">{d.description || '-'}</p>
                  <div className="flex items-center gap-1.5 mt-2">
                    <span className={`w-2 h-2 rounded-full ${d.is_active ? 'bg-emerald-500' : 'bg-red-500'}`} />
                    <span className="text-xs text-gray-400">{d.is_active ? 'Aktif' : 'Nonaktif'}</span>
                  </div>
                </div>
              </div>
              <div className="flex gap-1 mt-4">
                <button onClick={() => openEdit(d)} className="px-3 py-1.5 text-xs bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg flex items-center gap-1"><Pencil className="w-3 h-3" /> Edit</button>
                <button onClick={() => handleDelete(d.id)} className="px-3 py-1.5 text-xs bg-red-900/30 hover:bg-red-900/50 text-red-400 rounded-lg flex items-center gap-1"><Trash2 className="w-3 h-3" /> Hapus</button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showForm && (
        <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
          <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-md shadow-2xl">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
              <h2 className="font-bold text-white">{editId ? 'Edit Departemen' : 'Tambah Departemen'}</h2>
              <button onClick={() => setShowForm(false)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400"><X className="w-5 h-5" /></button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm text-gray-300 mb-1">Nama Departemen *</label>
                <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="pos-input" placeholder="Contoh: Marketing" />
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Deskripsi</label>
                <textarea value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} className="pos-input" rows={2} />
              </div>
              <button onClick={handleSave} className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition-colors">
                {editId ? 'Simpan Perubahan' : 'Tambah Departemen'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
