import React, { useEffect, useState } from 'react';
import { adminService, ChartOfAccount } from '../services/api';

const categoryLabels: Record<string, string> = {
  aset: 'Aset',
  kewajiban: 'Kewajiban',
  ekuitas: 'Ekuitas',
  pendapatan: 'Pendapatan',
  beban: 'Beban',
};

const categoryColors: Record<string, string> = {
  aset: 'bg-blue-100 text-blue-700',
  kewajiban: 'bg-orange-100 text-orange-700',
  ekuitas: 'bg-green-100 text-green-700',
  pendapatan: 'bg-purple-100 text-purple-700',
  beban: 'bg-red-100 text-red-700',
};

export default function ChartOfAccountsPage() {
  const [data, setData] = useState<ChartOfAccount[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<ChartOfAccount | null>(null);
  const [form, setForm] = useState({
    code: '',
    name: '',
    category: 'aset',
    normal_balance: 'debit',
    description: '',
    is_active: true,
  });

  const load = async () => {
    try {
      const res = await adminService.listChartOfAccounts();
      setData(res.data.data);
    } catch {} finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openAdd = () => {
    setEditItem(null);
    setForm({ code: '', name: '', category: 'aset', normal_balance: 'debit', description: '', is_active: true });
    setShowForm(true);
  };

  const openEdit = (item: ChartOfAccount) => {
    setEditItem(item);
    setForm({
      code: item.code,
      name: item.name,
      category: item.category,
      normal_balance: item.normal_balance,
      description: item.description,
      is_active: item.is_active,
    });
    setShowForm(true);
  };

  const changeCategory = (category: string) => {
    const normalBalance = category === 'aset' || category === 'beban' ? 'debit' : 'kredit';
    setForm({ ...form, category, normal_balance: normalBalance });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editItem) {
        await adminService.updateChartOfAccount(editItem.id, form);
      } else {
        await adminService.createChartOfAccount(form);
      }
      setShowForm(false);
      load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus akun perkiraan ini?')) return;
    try {
      await adminService.deleteChartOfAccount(id);
      load();
    } catch {}
  };

  const categories = Object.keys(categoryLabels);

  const totals = categories.map((cat) => ({
    category: cat,
    label: categoryLabels[cat],
    count: data.filter((d) => d.category === cat).length,
    total: data.filter((d) => d.category === cat).reduce((s, d) => s + d.balance, 0),
  }));

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold">Akun Perkiraan</h1>
          <p className="text-sm text-gray-500">Bagan struktur akun keuangan (Aset, Kewajiban, Ekuitas, Pendapatan, Beban)</p>
        </div>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Akun Baru</button>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-6">
        {totals.map((t) => (
          <div key={t.category} className={`p-4 rounded ${categoryColors[t.category]}`}>
            <div className="text-sm font-medium">{t.label}</div>
            <div className="text-xl font-bold mt-1">Rp {t.total.toLocaleString()}</div>
            <div className="text-xs mt-1 opacity-70">{t.count} akun</div>
          </div>
        ))}
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Akun' : 'Akun Baru'}</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Kode Akun</label>
              <input
                placeholder="contoh: 1-1000"
                value={form.code}
                onChange={(e) => setForm({ ...form, code: e.target.value })}
                className="border p-2 rounded w-full"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Nama Akun</label>
              <input
                placeholder="contoh: Kas"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                className="border p-2 rounded w-full"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Kategori</label>
              <select
                value={form.category}
                onChange={(e) => changeCategory(e.target.value)}
                className="border p-2 rounded w-full"
              >
                {categories.map((cat) => (
                  <option key={cat} value={cat}>{categoryLabels[cat]}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Saldo Normal</label>
              <select
                value={form.normal_balance}
                onChange={(e) => setForm({ ...form, normal_balance: e.target.value })}
                className="border p-2 rounded w-full"
              >
                <option value="debit">Debit</option>
                <option value="kredit">Kredit</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Deskripsi</label>
              <input
                placeholder="Keterangan akun"
                value={form.description}
                onChange={(e) => setForm({ ...form, description: e.target.value })}
                className="border p-2 rounded w-full"
              />
            </div>
            <div className="flex items-end">
              <label className="flex items-center gap-2 text-sm font-medium">
                <input
                  type="checkbox"
                  checked={form.is_active}
                  onChange={(e) => setForm({ ...form, is_active: e.target.checked })}
                  className="w-4 h-4"
                />
                Akun Aktif
              </label>
            </div>
            <div className="col-span-2 flex gap-2">
              <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Simpan</button>
              <button type="button" onClick={() => setShowForm(false)} className="bg-gray-400 text-white px-4 py-2 rounded">Batal</button>
            </div>
          </form>
        </div>
      )}

      {loading ? (
        <p>Loading...</p>
      ) : (
        <div className="bg-white rounded shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-3 text-left">Kode</th>
                <th className="p-3 text-left">Nama Akun</th>
                <th className="p-3 text-left">Kategori</th>
                <th className="p-3 text-left">Saldo Normal</th>
                <th className="p-3 text-right">Saldo</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map((d) => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.code}</td>
                  <td className="p-3">
                    {d.name}
                    {d.description && <div className="text-xs text-gray-400">{d.description}</div>}
                  </td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs ${categoryColors[d.category]}`}>
                      {categoryLabels[d.category]}
                    </span>
                  </td>
                  <td className="p-3 capitalize">{d.normal_balance}</td>
                  <td className="p-3 text-right font-medium">Rp {d.balance.toLocaleString()}</td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs ${d.is_active ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
                      {d.is_active ? 'Aktif' : 'Nonaktif'}
                    </span>
                  </td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => openEdit(d)} className="text-blue-500 hover:underline">Edit</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && (
                <tr>
                  <td colSpan={7} className="p-4 text-center text-gray-400">Belum ada akun perkiraan</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
