import React, { useState, useEffect } from 'react';
import { Users, Plus, Pencil, Trash2, X, Shield, Lock, Eye, EyeOff, UserCheck, UserX, Search, FileText } from 'lucide-react';
import { adminService, type User, type Permission, type TransactionLog } from '../services/api';
import toast from 'react-hot-toast';

export function AdminPage() {
  const [activeTab, setActiveTab] = useState<'users' | 'logs'>('users');
  const [users, setUsers] = useState<User[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);
  const [showUserForm, setShowUserForm] = useState(false);
  const [showPermModal, setShowPermModal] = useState(false);
  const [showPassModal, setShowPassModal] = useState(false);
  const [editUser, setEditUser] = useState<User | null>(null);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [showPass, setShowPass] = useState(false);

  const [form, setForm] = useState({ username: '', name: '', password: '', role: 'cashier' });
  const [newPass, setNewPass] = useState('');
  const [selectedPerms, setSelectedPerms] = useState<number[]>([]);

  const [logSearchUsername, setLogSearchUsername] = useState<number | ''>('');
  const [logData, setLogData] = useState<TransactionLog[]>([]);
  const [logLoading, setLogLoading] = useState(false);
  const [logError, setLogError] = useState('');

  const loadUsers = async () => {
    try {
      const res = await adminService.listUsers();
      setUsers(res.data.data || []);
    } catch { setUsers([]); }
    setLoading(false);
  };

  const loadPermissions = async () => {
    try {
      const res = await adminService.getAllPermissions();
      setPermissions(res.data.data || []);
    } catch { setPermissions([]); }
  };

  useEffect(() => { loadUsers(); loadPermissions(); }, []);

  const openCreate = () => {
    setEditUser(null);
    setForm({ username: '', name: '', password: '', role: 'cashier' });
    setShowUserForm(true);
  };

  const openEdit = (u: User) => {
    setEditUser(u);
    setForm({ username: u.username, name: u.name, password: '', role: u.role });
    setShowUserForm(true);
  };

  const handleSave = async () => {
    if (!form.username || !form.name) { toast.error('Username dan nama wajib diisi'); return; }
    try {
      if (editUser) {
        await adminService.updateUser(editUser.id, { name: form.name, role: form.role });
        toast.success('User diupdate!');
      } else {
        if (!form.password) { toast.error('Password wajib diisi'); return; }
        if (!/^\d{6}$/.test(form.password)) { toast.error('PIN harus tepat 6 digit angka'); return; }
        await adminService.createUser(form);
        toast.success('User dibuat!');
      }
      setShowUserForm(false);
      loadUsers();
    } catch (err: any) {
      toast.error(err?.response?.data?.error || 'Gagal menyimpan');
    }
  };

  const handleDelete = async (u: User) => {
    if (!confirm(`Hapus user ${u.username}?`)) return;
    try {
      await adminService.deleteUser(u.id);
      toast.success('User dihapus!');
      loadUsers();
    } catch { toast.error('Gagal menghapus'); }
  };

  const handleToggleActive = async (u: User) => {
    try {
      await adminService.updateUser(u.id, { is_active: !u.is_active });
      toast.success(u.is_active ? 'User dinonaktifkan' : 'User diaktifkan');
      loadUsers();
    } catch { toast.error('Gagal update status'); }
  };

  const openPermissions = async (u: User) => {
    setSelectedUser(u);
    try {
      const res = await adminService.getUserPermissions(u.id);
      setSelectedPerms((res.data.data || []).map((p) => p.id));
    } catch { setSelectedPerms([]); }
    setShowPermModal(true);
  };

  const handleSavePerms = async () => {
    if (!selectedUser) return;
    try {
      await adminService.setUserPermissions(selectedUser.id, selectedPerms);
      toast.success('Permission disimpan!');
      setShowPermModal(false);
    } catch { toast.error('Gagal menyimpan permission'); }
  };

  const togglePerm = (permId: number) => {
    setSelectedPerms((prev) =>
      prev.includes(permId) ? prev.filter((id) => id !== permId) : [...prev, permId]
    );
  };

  const openPassword = (u: User) => {
    setSelectedUser(u);
    setNewPass('');
    setShowPassModal(true);
  };

  const handleSearchLog = async () => {
    if (logSearchUsername === '') {
      toast.error('Pilih kasir terlebih dahulu');
      return;
    }
    const user = users.find((u) => u.id === logSearchUsername);
    if (!user) { toast.error('User tidak ditemukan'); return; }
    setLogLoading(true);
    setLogError('');
    setLogData([]);
    try {
      const res = await adminService.getTransactionLogsByUsername(user.username);
      const logs = res.data.data || [];
      setLogData(logs);
      if (logs.length === 0) {
        setLogError('Tidak ada transaksi ditemukan untuk user "' + user.username + '"');
      }
    } catch {
      setLogError('Gagal mencari transaksi');
    }
    setLogLoading(false);
  };

  const handleChangePass = async () => {
    if (!selectedUser || !newPass) { toast.error('Password wajib diisi'); return; }
    if (!/^\d{6}$/.test(newPass)) { toast.error('PIN harus tepat 6 digit angka'); return; }
    try {
      await adminService.updatePassword(selectedUser.id, newPass);
      toast.success('Password diupdate!');
      setShowPassModal(false);
    } catch { toast.error('Gagal update password'); }
  };

  const roleLabel = (r: string) => r === 'admin' ? 'Admin' : 'Kasir';

  return (
    <div className="space-y-6">
      {/* Tabs */}
      <div className="flex gap-1 border-b border-gray-200">
        <button
          onClick={() => setActiveTab('users')}
          className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
            activeTab === 'users'
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-500 hover:text-gray-700'
          }`}
        >
          <Users className="w-4 h-4 inline mr-1.5" />
          Kelola User
        </button>
        <button
          onClick={() => setActiveTab('logs')}
          className={`px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
            activeTab === 'logs'
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-500 hover:text-gray-700'
          }`}
        >
          <FileText className="w-4 h-4 inline mr-1.5" />
          Log Transaksi
        </button>
      </div>

      {/* Tab: Users */}
      {activeTab === 'users' && (
        <>
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Kelola User</h1>
              <p className="text-gray-500 text-sm mt-1">Buat, edit, dan atur hak akses kasir</p>
            </div>
            <button onClick={openCreate} className="pos-btn-primary flex items-center gap-2">
              <Plus className="w-4 h-4" /> Tambah User
            </button>
          </div>

          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {[...Array(3)].map((_, i) => (
                <div key={i} className="bg-gray-800/50 rounded-xl h-32 animate-pulse" />
              ))}
            </div>
          ) : users.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <Users className="w-12 h-12 mx-auto mb-3 opacity-50" />
              <p>Belum ada user</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {users.map((u) => (
                <div key={u.id} className="bg-gray-900 border border-gray-800 rounded-xl p-4">
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-3">
                      <div className={`w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold ${u.role === 'admin' ? 'bg-gradient-to-br from-purple-500 to-blue-600 text-white' : 'bg-gradient-to-br from-emerald-500 to-teal-600 text-white'}`}>
                        {u.name.charAt(0).toUpperCase()}
                      </div>
                      <div>
                        <p className="text-white font-medium">{u.name}</p>
                        <p className="text-xs text-gray-400">@{u.username}</p>
                      </div>
                    </div>
                    <span className={`text-xs px-2 py-1 rounded-full ${u.role === 'admin' ? 'bg-purple-500/20 text-purple-400' : 'bg-emerald-500/20 text-emerald-400'}`}>
                      {roleLabel(u.role)}
                    </span>
                  </div>

                  <div className="flex items-center gap-2 mt-3">
                    <span className={`w-2 h-2 rounded-full ${u.is_active ? 'bg-emerald-500' : 'bg-red-500'}`} />
                    <span className="text-xs text-gray-400">{u.is_active ? 'Aktif' : 'Nonaktif'}</span>
                  </div>

                  <div className="flex gap-1 mt-4 flex-wrap">
                    <button onClick={() => openEdit(u)} className="px-3 py-1.5 text-xs bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg flex items-center gap-1">
                      <Pencil className="w-3 h-3" /> Edit
                    </button>
                    <button onClick={() => openPassword(u)} className="px-3 py-1.5 text-xs bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg flex items-center gap-1">
                      <Lock className="w-3 h-3" /> Password
                    </button>
                    <button onClick={() => openPermissions(u)} className="px-3 py-1.5 text-xs bg-gray-800 hover:bg-gray-700 text-gray-300 rounded-lg flex items-center gap-1">
                      <Shield className="w-3 h-3" /> Akses
                    </button>
                    <button onClick={() => handleToggleActive(u)} className={`px-3 py-1.5 text-xs rounded-lg flex items-center gap-1 ${u.is_active ? 'bg-red-900/30 hover:bg-red-900/50 text-red-400' : 'bg-emerald-900/30 hover:bg-emerald-900/50 text-emerald-400'}`}>
                      {u.is_active ? <><UserX className="w-3 h-3" /> Nonaktif</> : <><UserCheck className="w-3 h-3" /> Aktifkan</>}
                    </button>
                    {u.username !== 'admin' && (
                      <button onClick={() => handleDelete(u)} className="px-3 py-1.5 text-xs bg-red-900/30 hover:bg-red-900/50 text-red-400 rounded-lg flex items-center gap-1">
                        <Trash2 className="w-3 h-3" /> Hapus
                      </button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </>
      )}

      {/* Tab: Log Transaksi */}
      {activeTab === 'logs' && (
        <>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Log Transaksi</h1>
            <p className="text-gray-500 text-sm mt-1">Cari riwayat detail transaksi berdasarkan ID</p>
          </div>

          <div className="flex gap-3 items-end">
            <div className="flex-1 max-w-xs">
              <label className="block text-sm text-gray-600 mb-1">Pilih Kasir</label>
              <select
                value={logSearchUsername}
                onChange={(e) => setLogSearchUsername(Number(e.target.value))}
                className="w-full px-4 py-2.5 border border-gray-300 rounded-xl text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="">-- Pilih Kasir --</option>
                {users.map((u) => (
                  <option key={u.id} value={u.id}>
                    {u.name} (@{u.username})
                  </option>
                ))}
              </select>
            </div>
            <button
              onClick={handleSearchLog}
              disabled={logLoading || logSearchUsername === ''}
              className="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition-colors flex items-center gap-2 disabled:opacity-50"
            >
              <Search className="w-4 h-4" />
              {logLoading ? 'Mencari...' : 'Cari'}
            </button>
          </div>

          {logError && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-xl">
              {logError}
            </div>
          )}

          {logData.length > 0 && (
            <div className="space-y-4">
              <p className="text-sm text-gray-500">Ditemukan {logData.length} transaksi</p>
              {logData.map((log) => (
                <div key={log.id} className="bg-white border border-gray-200 rounded-2xl overflow-hidden">
                  <div className="bg-gray-900 px-6 py-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <h2 className="text-lg font-bold text-white">Transaksi #{log.id}</h2>
                        <p className="text-gray-400 text-sm">Order #{log.order_id}</p>
                      </div>
                      <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                        log.status === 'completed' ? 'bg-emerald-500/20 text-emerald-400' : 'bg-yellow-500/20 text-yellow-400'
                      }`}>
                        {log.status === 'completed' ? 'Selesai' : log.status}
                      </span>
                    </div>
                  </div>

                  <div className="p-6 space-y-4">
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                      <div className="bg-gray-50 rounded-xl p-3">
                        <p className="text-xs text-gray-500 mb-1">Total Bayar</p>
                        <p className="text-lg font-bold text-gray-900">Rp {log.amount.toLocaleString()}</p>
                      </div>
                      <div className="bg-gray-50 rounded-xl p-3">
                        <p className="text-xs text-gray-500 mb-1">Metode Bayar</p>
                        <p className="text-lg font-bold text-gray-900 capitalize">{log.payment_method}</p>
                      </div>
                      <div className="bg-gray-50 rounded-xl p-3">
                        <p className="text-xs text-gray-500 mb-1">Kasir</p>
                        <p className="text-lg font-bold text-gray-900">{log.user_name}</p>
                        <p className="text-xs text-gray-500 capitalize">{log.user_role}</p>
                      </div>
                      <div className="bg-gray-50 rounded-xl p-3">
                        <p className="text-xs text-gray-500 mb-1">Shift</p>
                        <p className="text-lg font-bold text-gray-900">#{log.shift_number}</p>
                        <p className="text-xs text-gray-500">{new Date(log.shift_opened).toLocaleDateString('id-ID')}</p>
                      </div>
                    </div>

                    <div className="flex gap-6 text-sm text-gray-500">
                      <span>Dibuat: {new Date(log.created_at).toLocaleString('id-ID')}</span>
                    </div>

                    {log.items && log.items.length > 0 && (
                      <div>
                        <h3 className="text-sm font-semibold text-gray-900 mb-3">Item Pesanan</h3>
                        <div className="border border-gray-200 rounded-xl overflow-hidden">
                          <table className="w-full text-sm">
                            <thead className="bg-gray-50">
                              <tr>
                                <th className="text-left px-4 py-2.5 text-gray-600 font-medium">Produk</th>
                                <th className="text-right px-4 py-2.5 text-gray-600 font-medium">Qty</th>
                                <th className="text-right px-4 py-2.5 text-gray-600 font-medium">Harga</th>
                                <th className="text-right px-4 py-2.5 text-gray-600 font-medium">Subtotal</th>
                              </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200">
                              {log.items.map((item, idx) => (
                                <tr key={idx}>
                                  <td className="px-4 py-2.5 text-gray-900">{item.product_name}</td>
                                  <td className="px-4 py-2.5 text-gray-600 text-right">{item.quantity}</td>
                                  <td className="px-4 py-2.5 text-gray-600 text-right">Rp {item.price.toLocaleString()}</td>
                                  <td className="px-4 py-2.5 text-gray-900 font-medium text-right">Rp {item.subtotal.toLocaleString()}</td>
                                </tr>
                              ))}
                            </tbody>
                            <tfoot className="bg-gray-50">
                              <tr>
                                <td colSpan={3} className="px-4 py-2.5 text-gray-600 font-medium text-right">Total Order</td>
                                <td className="px-4 py-2.5 text-gray-900 font-bold text-right">Rp {log.order_total.toLocaleString()}</td>
                              </tr>
                            </tfoot>
                          </table>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </>
      )}

      {/* User Form Modal */}
      {showUserForm && (
        <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
          <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-md shadow-2xl">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
              <h2 className="font-bold text-white">{editUser ? 'Edit User' : 'Tambah User'}</h2>
              <button onClick={() => setShowUserForm(false)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400">
                <X className="w-5 h-5" />
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm text-gray-300 mb-1">Username *</label>
                <input value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })} className="pos-input" disabled={!!editUser} placeholder="username_kasir" />
              </div>
              <div>
                <label className="block text-sm text-gray-300 mb-1">Nama Lengkap *</label>
                <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="pos-input" placeholder="Nama Kasir" />
              </div>
              {!editUser && (
                <div>
                  <label className="block text-sm text-gray-300 mb-1">PIN (6 digit angka) *</label>
                  <input type="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} className="pos-input" placeholder="Contoh: 123456" maxLength={6} inputMode="numeric" />
                </div>
              )}
              <div>
                <label className="block text-sm text-gray-300 mb-1">Role</label>
                <select value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value })} className="pos-input">
                  <option value="cashier">Kasir</option>
                  <option value="admin">Admin</option>
                </select>
              </div>
              <button onClick={handleSave} className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition-colors">
                {editUser ? 'Simpan Perubahan' : 'Buat User'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Permissions Modal */}
      {showPermModal && selectedUser && (
        <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
          <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-md shadow-2xl">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
              <h2 className="font-bold text-white">Hak Akses - {selectedUser.name}</h2>
              <button onClick={() => setShowPermModal(false)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400">
                <X className="w-5 h-5" />
              </button>
            </div>
            <div className="p-6 space-y-3">
              {permissions.map((p) => (
                <label key={p.id} className="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-800 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={selectedPerms.includes(p.id)}
                    onChange={() => togglePerm(p.id)}
                    className="w-4 h-4 rounded bg-gray-800 border-gray-600 text-blue-600 focus:ring-blue-500"
                  />
                  <div>
                    <p className="text-sm text-white">{p.description}</p>
                    <p className="text-xs text-gray-500">{p.name}</p>
                  </div>
                </label>
              ))}
              <button onClick={handleSavePerms} className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition-colors mt-4">
                Simpan Hak Akses
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Password Modal */}
      {showPassModal && selectedUser && (
        <div className="fixed inset-0 bg-black/70 z-50 flex items-center justify-center p-4">
          <div className="bg-gray-900 border border-gray-800 rounded-2xl w-full max-w-md shadow-2xl">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-800">
              <h2 className="font-bold text-white">Ubah Password - {selectedUser.name}</h2>
              <button onClick={() => setShowPassModal(false)} className="p-2 rounded-lg hover:bg-gray-800 text-gray-400">
                <X className="w-5 h-5" />
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div>
                <label className="block text-sm text-gray-300 mb-1">PIN Baru (6 digit angka)</label>
                <div className="relative">
                  <input
                    type={showPass ? 'text' : 'password'}
                    value={newPass}
                    onChange={(e) => setNewPass(e.target.value)}
                    className="pos-input pr-10"
                    placeholder="Contoh: 123456"
                    maxLength={6}
                    inputMode="numeric"
                  />
                  <button onClick={() => setShowPass(!showPass)} className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400">
                    {showPass ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
              </div>
              <button onClick={handleChangePass} className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-xl transition-colors">
                Simpan Password
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
