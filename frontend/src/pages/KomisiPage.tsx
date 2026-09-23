import React, { useState, useEffect } from 'react';
import { Plus, Pencil, Trash2, X, Percent } from 'lucide-react';
import { adminService, type Commission, type Employee } from '../services/api';
import toast from 'react-hot-toast';

interface CommForm { employee_id: number; commission_type: string; rate: number; description: string }
const emptyForm: CommForm = { employee_id: 0, commission_type: 'percentage', rate: 0, description: '' };

export function KomisiPage() {
  const [items, setItems] = useState<Commission[]>([]);
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);
  const [form, setForm] = useState<CommForm>(emptyForm);

  const load = async () => {
    try {
      const [commRes, empRes] = await Promise.all([adminService.listCommissions(), adminService.listEmployees()]);
      setItems(commRes.data.data || []);
      setEmployees(empRes.data.data || []);
    } catch { setItems([]); }
    setLoading(false);
  };

  useEffect(() => { load(); }, []);

  const openCreate = () => { setEditId(null); setForm(emptyForm); setShowForm(true); };
  const openEdit = (c: Commission) => {
    setEditId(c.id);
    setForm({ employee_id: c.employee_id, commission_type: c.commission_type, rate: c.rate, description: c.description });
    setShowForm(true);
  };

  const handleSave = async () => {
    if (!form.employee_id) { toast.error('Karyawan wajib dipilih'); return; }
    if (form.rate <= 0) { toast.error('Rate harus lebih dari 0'); return; }
    try {
      if (editId) {
        await adminService.updateCommission(editId, form);
        toast.success('Komisi diupdate!');
      } else {
        await adminService.createCommission(form);
        toast.success('Komisi ditambahkan!');
      }
      setShowForm(false);
      load();
    } catch (err: any) { toast.error(err?.response?.data?.error || 'Gagal menyimpan'); }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus komisi ini?')) return;
    try { await adminService.deleteCommission(id); toast.success('Komisi dihapus!'); load(); }
    catch { toast.error('Gagal menghapus'); }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Komisi Karyawan</h1>
          <p className="text-gray-500 text-sm mt-1">Atur komisi per karyawan</p>
        </div>
        <button onClick={openCreate} className="pos-btn-primary flex items-center gap-2"><Plus className="w-4 h-4" /> Tambah Komisi</button>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(3)].map((_, i) => <div key={i} className="bg-gray-800/50 rounded-xl h-28 animate-pulse" />)}
        </div>
      ) : items.length === 0 ? (
        <div className="text-center py-12 text-gray-500"><Percent className="w-12 h-12 mx-auto mb-3 opacity-50" /><p>Belum ada komisi</p></div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {items.map((c) => (
            <div key={c.id} className="bg-gray-900 border border-gray-800 rounded-xl p-4">
              <div className="flex items-start justify-between">
                <div>
                  <p className="text-white font-medium">{c.employee_name}</p>
                  <p className="text-xs text-gray-400 mt-1 capitalize">{c.commission_type === 'percentage' ? 'Persentase' : 'Flat'}</p>
                </div>
                <span className="text-blue-400 font-bold text-lg">{c.commission_type === 'percentage' ? `${c.rate}%` : `Rp ${c.rate.toLocaleString()}`}</span>
              </div>
              {c.description && <p className="text-sm text-gray-400 mt-2">{c.description}</p>}
              <div className="flex gap-1 mt-4">
                <button onClick={() => openEdit(c)} className="px-3 py-1.5 text-xs bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg flex items-center gap-1"><Pencil className="w-3 h-3" /> Edit</button>
                <button onClick={() => handleDelete(c.id)} className="px-3 py-1.5 text-xs bg-red-900/30 hover:bg-red-900/50 text-red-400 rounded-lg flex items-center gap-1"><Trash2 className="w-3 h-3" /> Hapus</button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showForm && (
        <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
          <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-md shadow-2xl">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
              <h2 className="font-bold text-white">{editId ? 'Edit Komisi' : 'Tambah Komisi'}</h2>
              <button onClick={() => setShowForm(false)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400"><X className="w-5 h-5" /></button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm text-gray-300 mb-1">Karyawan *</label>
                <select value={form.employee_id || ''} onChange={(e) => setForm({ ...form, employee_id: Number(e.target.value) })} className="pos-input" disabled={!!editId}>
                  <option value="">-- Pilih Karyawan --</option>
                  {employees.map((e) => <option key={e.id} value={e.id}>{e.name} ({e.nik})</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Tipe Komisi</label>
                <select value={form.commission_type} onChange={(e) => setForm({ ...form, commission_type: e.target.value })} className="pos-input">
                  <option value="percentage">Persentase (%)</option>
                  <option value="flat">Flat (Rp)</option>
                </select>
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Rate *</label>
                <input type="number" value={form.rate} onChange={(e) => setForm({ ...form, rate: Number(e.target.value) })} className="pos-input" placeholder={form.commission_type === 'percentage' ? 'Contoh: 5' : 'Contoh: 10000'} />
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Deskripsi</label>
                <textarea value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} className="pos-input" rows={2} />
              </div>
              <button onClick={handleSave} className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition-colors">
                {editId ? 'Simpan Perubahan' : 'Tambah Komisi'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
