import React, { useState, useEffect } from 'react';
import { Plus, Pencil, Trash2, X, Users, Search } from 'lucide-react';
import { adminService, type Employee, type Department } from '../services/api';
import toast from 'react-hot-toast';

interface EmpForm {
  nik: string; name: string; email: string; phone: string;
  department_id: number | null; position: string; hire_date: string; salary_base: number;
}
const emptyForm: EmpForm = { nik: '', name: '', email: '', phone: '', department_id: null, position: '', hire_date: '', salary_base: 0 };

export function DaftarKaryawanPage() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);
  const [form, setForm] = useState<EmpForm>(emptyForm);
  const [searchQ, setSearchQ] = useState('');

  const load = async () => {
    try {
      const [empRes, deptRes] = await Promise.all([adminService.listEmployees(), adminService.listDepartments()]);
      setEmployees(empRes.data.data || []);
      setDepartments(deptRes.data.data || []);
    } catch { setEmployees([]); }
    setLoading(false);
  };

  useEffect(() => { load(); }, []);

  const filtered = employees.filter((e) => !searchQ || e.name.toLowerCase().includes(searchQ.toLowerCase()) || e.nik.toLowerCase().includes(searchQ.toLowerCase()));

  const openCreate = () => { setEditId(null); setForm(emptyForm); setShowForm(true); };
  const openEdit = (e: Employee) => {
    setEditId(e.id);
    setForm({ nik: e.nik, name: e.name, email: e.email, phone: e.phone, department_id: e.department_id, position: e.position, hire_date: e.hire_date?.split('T')[0] || '', salary_base: e.salary_base });
    setShowForm(true);
  };

  const handleSave = async () => {
    if (!form.nik || !form.name) { toast.error('NIK dan nama wajib diisi'); return; }
    try {
      if (editId) {
        await adminService.updateEmployee(editId, form);
        toast.success('Karyawan diupdate!');
      } else {
        await adminService.createEmployee(form);
        toast.success('Karyawan ditambahkan!');
      }
      setShowForm(false);
      load();
    } catch (err: any) { toast.error(err?.response?.data?.error || 'Gagal menyimpan'); }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus karyawan ini?')) return;
    try { await adminService.deleteEmployee(id); toast.success('Karyawan dihapus!'); load(); }
    catch { toast.error('Gagal menghapus'); }
  };

  const handleToggleActive = async (e: Employee) => {
    try { await adminService.updateEmployee(e.id, { is_active: !e.is_active }); toast.success(e.is_active ? 'Dinonaktifkan' : 'Diaktifkan'); load(); }
    catch { toast.error('Gagal update status'); }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Daftar Karyawan</h1>
          <p className="text-gray-500 text-sm mt-1">Kelola data karyawan</p>
        </div>
        <button onClick={openCreate} className="pos-btn-primary flex items-center gap-2"><Plus className="w-4 h-4" /> Tambah Karyawan</button>
      </div>

      <div className="relative max-w-md">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
        <input value={searchQ} onChange={(e) => setSearchQ(e.target.value)} placeholder="Cari karyawan..." className="pos-input pl-10" />
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[...Array(3)].map((_, i) => <div key={i} className="bg-gray-800/50 rounded-xl h-32 animate-pulse" />)}
        </div>
      ) : filtered.length === 0 ? (
        <div className="text-center py-12 text-gray-500"><Users className="w-12 h-12 mx-auto mb-3 opacity-50" /><p>Belum ada karyawan</p></div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filtered.map((e) => (
            <div key={e.id} className="bg-gray-900 border border-gray-800 rounded-xl p-4">
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center text-white text-sm font-bold">{e.name.charAt(0).toUpperCase()}</div>
                  <div>
                    <p className="text-white font-medium">{e.name}</p>
                    <p className="text-xs text-gray-400">NIK: {e.nik}</p>
                  </div>
                </div>
                <div className="flex items-center gap-1.5">
                  <span className={`w-2 h-2 rounded-full ${e.is_active ? 'bg-emerald-500' : 'bg-red-500'}`} />
                  <span className="text-xs text-gray-400">{e.is_active ? 'Aktif' : 'Nonaktif'}</span>
                </div>
              </div>
              <div className="mt-3 space-y-1 text-sm text-gray-400">
                <p>Posisi: <span className="text-gray-300">{e.position || '-'}</span></p>
                <p>Dept: <span className="text-gray-300">{e.dept_name || '-'}</span></p>
                <p>Gaji: <span className="text-blue-400 font-medium">Rp {e.salary_base.toLocaleString()}</span></p>
              </div>
              <div className="flex gap-1 mt-4 flex-wrap">
                <button onClick={() => openEdit(e)} className="px-3 py-1.5 text-xs bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg flex items-center gap-1"><Pencil className="w-3 h-3" /> Edit</button>
                <button onClick={() => handleToggleActive(e)} className={`px-3 py-1.5 text-xs rounded-lg ${e.is_active ? 'bg-red-900/30 hover:bg-red-900/50 text-red-400' : 'bg-emerald-900/30 hover:bg-emerald-900/50 text-emerald-400'}`}>
                  {e.is_active ? 'Nonaktif' : 'Aktifkan'}
                </button>
                <button onClick={() => handleDelete(e.id)} className="px-3 py-1.5 text-xs bg-red-900/30 hover:bg-red-900/50 text-red-400 rounded-lg flex items-center gap-1"><Trash2 className="w-3 h-3" /> Hapus</button>
              </div>
            </div>
          ))}
        </div>
      )}

      {showForm && (
        <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
          <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto shadow-2xl">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
              <h2 className="font-bold text-white">{editId ? 'Edit Karyawan' : 'Tambah Karyawan'}</h2>
              <button onClick={() => setShowForm(false)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400"><X className="w-5 h-5" /></button>
            </div>
            <div className="p-6 space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div><label className="block text-sm text-gray-300 mb-1">NIK *</label><input value={form.nik} onChange={(e) => setForm({ ...form, nik: e.target.value })} className="pos-input" disabled={!!editId} placeholder="NIK Karyawan" /></div>
                <div><label className="block text-sm text-gray-300 mb-1">Nama *</label><input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="pos-input" placeholder="Nama Lengkap" /></div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div><label className="block text-sm text-gray-300 mb-1">Email</label><input type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} className="pos-input" placeholder="email@contoh.com" /></div>
                <div><label className="block text-sm text-gray-300 mb-1">Telepon</label><input value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} className="pos-input" placeholder="08123456789" /></div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm text-gray-300 mb-1">Departemen</label>
                  <select value={form.department_id || ''} onChange={(e) => setForm({ ...form, department_id: e.target.value ? Number(e.target.value) : null })} className="pos-input">
                    <option value="">-- Pilih --</option>
                    {departments.map((d) => <option key={d.id} value={d.id}>{d.name}</option>)}
                  </select>
                </div>
                <div><label className="block text-sm text-gray-300 mb-1">Posisi</label><input value={form.position} onChange={(e) => setForm({ ...form, position: e.target.value })} className="pos-input" placeholder="Kasir, Server, dll" /></div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div><label className="block text-sm text-gray-300 mb-1">Tanggal Masuk</label><input type="date" value={form.hire_date} onChange={(e) => setForm({ ...form, hire_date: e.target.value })} className="pos-input" /></div>
                <div><label className="block text-sm text-gray-300 mb-1">Gaji Pokok</label><input type="number" value={form.salary_base} onChange={(e) => setForm({ ...form, salary_base: Number(e.target.value) })} className="pos-input" /></div>
              </div>
              <button onClick={handleSave} className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition-colors">
                {editId ? 'Simpan Perubahan' : 'Tambah Karyawan'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
