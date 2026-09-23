import React, { useState, useEffect } from 'react';
import { Plus, Pencil, Trash2, X, Banknote, CheckCircle } from 'lucide-react';
import { adminService, type Salary, type Employee } from '../services/api';
import toast from 'react-hot-toast';

interface SalaryForm { employee_id: number; pay_period: string; base_salary: number; commission_total: number; deduction: number }
const emptyForm: SalaryForm = { employee_id: 0, pay_period: '', base_salary: 0, commission_total: 0, deduction: 0 };

export function GajiPage() {
  const [items, setItems] = useState<Salary[]>([]);
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);
  const [form, setForm] = useState<SalaryForm>(emptyForm);

  const load = async () => {
    try {
      const [salRes, empRes] = await Promise.all([adminService.listSalaries(), adminService.listEmployees()]);
      setItems(salRes.data.data || []);
      setEmployees(empRes.data.data || []);
    } catch { setItems([]); }
    setLoading(false);
  };

  useEffect(() => { load(); }, []);

  const openCreate = () => { setEditId(null); setForm(emptyForm); setShowForm(true); };
  const openEdit = (s: Salary) => {
    setEditId(s.id);
    setForm({ employee_id: s.employee_id, pay_period: s.pay_period, base_salary: s.base_salary, commission_total: s.commission_total, deduction: s.deduction });
    setShowForm(true);
  };

  const handleSave = async () => {
    if (!form.employee_id) { toast.error('Karyawan wajib dipilih'); return; }
    if (!form.pay_period) { toast.error('Periode wajib diisi'); return; }
    if (form.base_salary <= 0) { toast.error('Gaji pokok harus lebih dari 0'); return; }
    try {
      if (editId) {
        await adminService.updateSalary(editId, form);
        toast.success('Gaji diupdate!');
      } else {
        await adminService.createSalary(form);
        toast.success('Gaji ditambahkan!');
      }
      setShowForm(false);
      load();
    } catch (err: any) { toast.error(err?.response?.data?.error || 'Gagal menyimpan'); }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus data gaji ini?')) return;
    try { await adminService.deleteSalary(id); toast.success('Gaji dihapus!'); load(); }
    catch { toast.error('Gagal menghapus'); }
  };

  const handleMarkPaid = async (id: number) => {
    if (!confirm('Tandai gaji ini sudah dibayar?')) return;
    try { await adminService.markSalaryPaid(id); toast.success('Gaji ditandai sudah dibayar!'); load(); }
    catch { toast.error('Gagal update status'); }
  };

  const statusBadge = (s: string) => {
    if (s === 'paid') return <span className="px-2 py-0.5 text-xs rounded-full bg-emerald-500/20 text-emerald-400">Dibayar</span>;
    return <span className="px-2 py-0.5 text-xs rounded-full bg-yellow-500/20 text-yellow-400">Pending</span>;
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Gaji Karyawan</h1>
          <p className="text-gray-500 text-sm mt-1">Kelola pembayaran gaji karyawan</p>
        </div>
        <button onClick={openCreate} className="pos-btn-primary flex items-center gap-2"><Plus className="w-4 h-4" /> Tambah Gaji</button>
      </div>

      {loading ? (
        <div className="space-y-3">
          {[...Array(3)].map((_, i) => <div key={i} className="bg-gray-800/50 rounded-xl h-20 animate-pulse" />)}
        </div>
      ) : items.length === 0 ? (
        <div className="text-center py-12 text-gray-500"><Banknote className="w-12 h-12 mx-auto mb-3 opacity-50" /><p>Belum ada data gaji</p></div>
      ) : (
        <div className="bg-white border border-gray-200 rounded-2xl overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50">
              <tr>
                <th className="text-left px-4 py-3 text-gray-600 font-medium">Karyawan</th>
                <th className="text-left px-4 py-3 text-gray-600 font-medium">Periode</th>
                <th className="text-right px-4 py-3 text-gray-600 font-medium">Gaji Pokok</th>
                <th className="text-right px-4 py-3 text-gray-600 font-medium">Komisi</th>
                <th className="text-right px-4 py-3 text-gray-600 font-medium">Potongan</th>
                <th className="text-right px-4 py-3 text-gray-600 font-medium">Total</th>
                <th className="text-center px-4 py-3 text-gray-600 font-medium">Status</th>
                <th className="text-center px-4 py-3 text-gray-600 font-medium">Aksi</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {items.map((s) => (
                <tr key={s.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3 text-gray-900 font-medium">{s.employee_name}</td>
                  <td className="px-4 py-3 text-gray-600">{s.pay_period}</td>
                  <td className="px-4 py-3 text-gray-600 text-right">Rp {s.base_salary.toLocaleString()}</td>
                  <td className="px-4 py-3 text-emerald-600 text-right">Rp {s.commission_total.toLocaleString()}</td>
                  <td className="px-4 py-3 text-red-600 text-right">Rp {s.deduction.toLocaleString()}</td>
                  <td className="px-4 py-3 text-gray-900 font-bold text-right">Rp {s.net_salary.toLocaleString()}</td>
                  <td className="px-4 py-3 text-center">{statusBadge(s.status)}</td>
                  <td className="px-4 py-3 text-center">
                    <div className="flex gap-1 justify-center">
                      {s.status !== 'paid' && (
                        <button onClick={() => handleMarkPaid(s.id)} className="p-1.5 rounded-lg hover:bg-emerald-50 text-emerald-500" title="Tandai Dibayar"><CheckCircle className="w-4 h-4" /></button>
                      )}
                      <button onClick={() => openEdit(s)} className="p-1.5 rounded-lg hover:bg-blue-50 text-blue-500"><Pencil className="w-4 h-4" /></button>
                      <button onClick={() => handleDelete(s.id)} className="p-1.5 rounded-lg hover:bg-red-50 text-red-500"><Trash2 className="w-4 h-4" /></button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showForm && (
        <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
          <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-md shadow-2xl">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
              <h2 className="font-bold text-white">{editId ? 'Edit Gaji' : 'Tambah Gaji'}</h2>
              <button onClick={() => setShowForm(false)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400"><X className="w-5 h-5" /></button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm text-gray-300 mb-1">Karyawan *</label>
                <select value={form.employee_id || ''} onChange={(e) => {
                  const empId = Number(e.target.value);
                  const emp = employees.find((x) => x.id === empId);
                  setForm({ ...form, employee_id: empId, base_salary: emp?.salary_base || 0 });
                }} className="pos-input" disabled={!!editId}>
                  <option value="">-- Pilih Karyawan --</option>
                  {employees.map((e) => <option key={e.id} value={e.id}>{e.name} ({e.nik})</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Periode *</label>
                <input value={form.pay_period} onChange={(e) => setForm({ ...form, pay_period: e.target.value })} className="pos-input" placeholder="Contoh: 2026-08" />
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Gaji Pokok</label>
                <input type="number" value={form.base_salary} onChange={(e) => setForm({ ...form, base_salary: Number(e.target.value) })} className="pos-input" />
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Total Komisi</label>
                <input type="number" value={form.commission_total} onChange={(e) => setForm({ ...form, commission_total: Number(e.target.value) })} className="pos-input" />
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Potongan</label>
                <input type="number" value={form.deduction} onChange={(e) => setForm({ ...form, deduction: Number(e.target.value) })} className="pos-input" />
              </div>
              <div className="bg-gray-800 rounded-xl p-3 text-center">
                <p className="text-xs text-gray-400">Total Diterima</p>
                <p className="text-xl font-bold text-blue-400">Rp {(form.base_salary + form.commission_total - form.deduction).toLocaleString()}</p>
              </div>
              <button onClick={handleSave} className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition-colors">
                {editId ? 'Simpan Perubahan' : 'Tambah Gaji'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
