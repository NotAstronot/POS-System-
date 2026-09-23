import React, { useEffect, useState } from 'react';
import { adminService, ChartOfAccount, FixedAsset } from '../services/api';

const assetTypes = [
  { value: 'peralatan', label: 'Peralatan Toko' },
  { value: 'komputer', label: 'Komputer Kasir' },
  { value: 'kendaraan', label: 'Kendaraan Operasional' },
  { value: 'bangunan', label: 'Bangunan' },
  { value: 'lainnya', label: 'Lainnya' },
];

const assetTypeLabels: Record<string, string> = Object.fromEntries(assetTypes.map((t) => [t.value, t.label]));

const assetTypeColors: Record<string, string> = {
  peralatan: 'bg-blue-100 text-blue-700',
  komputer: 'bg-purple-100 text-purple-700',
  kendaraan: 'bg-orange-100 text-orange-700',
  bangunan: 'bg-green-100 text-green-700',
  lainnya: 'bg-gray-100 text-gray-700',
};

const methodLabels: Record<string, string> = {
  garis_lurus: 'Garis Lurus',
  saldo_menurun: 'Saldo Menurun (Pajak)',
};

export default function PencatatanAsetPage() {
  const [data, setData] = useState<FixedAsset[]>([]);
  const [accounts, setAccounts] = useState<ChartOfAccount[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<FixedAsset | null>(null);
  const [form, setForm] = useState({
    asset_code: '',
    name: '',
    asset_type: 'peralatan',
    purchase_date: '',
    cost: '',
    salvage_value: '0',
    useful_life_months: '12',
    depreciation_method: 'garis_lurus',
    depreciation_account_id: 0,
    accumulated_account_id: 0,
    is_active: true,
  });

  const load = async () => {
    try {
      const [assetsRes, accountsRes] = await Promise.all([
        adminService.listFixedAssets(),
        adminService.listChartOfAccounts(),
      ]);
      setData(assetsRes.data.data);
      setAccounts(accountsRes.data.data.filter((a) => a.is_active));
    } catch {} finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openAdd = () => {
    setEditItem(null);
    setForm({
      asset_code: '',
      name: '',
      asset_type: 'peralatan',
      purchase_date: '',
      cost: '',
      salvage_value: '0',
      useful_life_months: '12',
      depreciation_method: 'garis_lurus',
      depreciation_account_id: 0,
      accumulated_account_id: 0,
      is_active: true,
    });
    setShowForm(true);
  };

  const openEdit = (item: FixedAsset) => {
    setEditItem(item);
    setForm({
      asset_code: item.asset_code,
      name: item.name,
      asset_type: item.asset_type,
      purchase_date: item.purchase_date,
      cost: String(item.cost),
      salvage_value: String(item.salvage_value),
      useful_life_months: String(item.useful_life_months),
      depreciation_method: item.depreciation_method,
      depreciation_account_id: item.depreciation_account_id,
      accumulated_account_id: item.accumulated_account_id,
      is_active: item.is_active,
    });
    setShowForm(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload = {
        ...form,
        cost: parseFloat(form.cost),
        salvage_value: parseFloat(form.salvage_value),
        useful_life_months: parseInt(form.useful_life_months),
      };
      if (editItem) {
        await adminService.updateFixedAsset(editItem.id, payload);
      } else {
        await adminService.createFixedAsset(payload);
      }
      setShowForm(false);
      load();
    } catch {}
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus aset tetap ini?')) return;
    try {
      await adminService.deleteFixedAsset(id);
      load();
    } catch {}
  };

  const totals = assetTypes.map((t) => ({
    ...t,
    count: data.filter((d) => d.asset_type === t.value).length,
    total: data.filter((d) => d.asset_type === t.value).reduce((s, d) => s + d.cost, 0),
  }));

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold">Pencatatan Aset</h1>
          <p className="text-sm text-gray-500">
            Mencatat peralatan toko, komputer kasir, kendaraan operasional, hingga bangunan
          </p>
        </div>
        <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Aset Baru</button>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-5 gap-4 mb-6">
        {totals.map((t) => (
          <div key={t.value} className={`p-4 rounded ${assetTypeColors[t.value]}`}>
            <div className="text-sm font-medium">{t.label}</div>
            <div className="text-xl font-bold mt-1">Rp {t.total.toLocaleString()}</div>
            <div className="text-xs mt-1 opacity-70">{t.count} aset</div>
          </div>
        ))}
      </div>

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Aset' : 'Aset Baru'}</h2>
          <form onSubmit={handleSubmit} className="grid grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Kode Aset</label>
              <input
                placeholder="contoh: AST-001"
                value={form.asset_code}
                onChange={(e) => setForm({ ...form, asset_code: e.target.value })}
                className="border p-2 rounded w-full"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Nama Aset</label>
              <input
                placeholder="contoh: Mesin Kasir"
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                className="border p-2 rounded w-full"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Tipe Aset</label>
              <select
                value={form.asset_type}
                onChange={(e) => setForm({ ...form, asset_type: e.target.value })}
                className="border p-2 rounded w-full"
              >
                {assetTypes.map((t) => (
                  <option key={t.value} value={t.value}>{t.label}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Tanggal Perolehan</label>
              <input
                type="date"
                value={form.purchase_date}
                onChange={(e) => setForm({ ...form, purchase_date: e.target.value })}
                className="border p-2 rounded w-full"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Harga Perolehan (Rp)</label>
              <input
                type="number"
                placeholder="Harga perolehan"
                value={form.cost}
                onChange={(e) => setForm({ ...form, cost: e.target.value })}
                className="border p-2 rounded w-full"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Nilai Sisa (Rp)</label>
              <input
                type="number"
                placeholder="Nilai sisa"
                value={form.salvage_value}
                onChange={(e) => setForm({ ...form, salvage_value: e.target.value })}
                className="border p-2 rounded w-full"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Umur Manfaat (bulan)</label>
              <input
                type="number"
                value={form.useful_life_months}
                onChange={(e) => setForm({ ...form, useful_life_months: e.target.value })}
                className="border p-2 rounded w-full"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Metode Penyusutan</label>
              <select
                value={form.depreciation_method}
                onChange={(e) => setForm({ ...form, depreciation_method: e.target.value })}
                className="border p-2 rounded w-full"
              >
                <option value="garis_lurus">Garis Lurus</option>
                <option value="saldo_menurun">Saldo Menurun (Pajak)</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Akun Beban Depresiasi</label>
              <select
                value={form.depreciation_account_id}
                onChange={(e) => setForm({ ...form, depreciation_account_id: Number(e.target.value) })}
                className="border p-2 rounded w-full"
                required
              >
                <option value={0}>-- Pilih Akun --</option>
                {accounts.map((a) => (
                  <option key={a.id} value={a.id}>{a.code} - {a.name}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Akun Akumulasi Depresiasi</label>
              <select
                value={form.accumulated_account_id}
                onChange={(e) => setForm({ ...form, accumulated_account_id: Number(e.target.value) })}
                className="border p-2 rounded w-full"
                required
              >
                <option value={0}>-- Pilih Akun --</option>
                {accounts.map((a) => (
                  <option key={a.id} value={a.id}>{a.code} - {a.name}</option>
                ))}
              </select>
            </div>
            <div className="flex items-end">
              <label className="flex items-center gap-2 text-sm font-medium">
                <input
                  type="checkbox"
                  checked={form.is_active}
                  onChange={(e) => setForm({ ...form, is_active: e.target.checked })}
                  className="w-4 h-4"
                />
                Aset Aktif
              </label>
            </div>
            <div className="col-span-3 flex gap-2">
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
                <th className="p-3 text-left">Nama Aset</th>
                <th className="p-3 text-left">Tipe</th>
                <th className="p-3 text-left">Metode</th>
                <th className="p-3 text-right">Harga Perolehan</th>
                <th className="p-3 text-right">Nilai Sisa</th>
                <th className="p-3 text-right">Akumulasi</th>
                <th className="p-3 text-right">Nilai Buku</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map((d) => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.asset_code}</td>
                  <td className="p-3">{d.name}</td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs ${assetTypeColors[d.asset_type] || assetTypeColors.lainnya}`}>
                      {assetTypeLabels[d.asset_type] || d.asset_type}
                    </span>
                  </td>
                  <td className="p-3 text-sm">{methodLabels[d.depreciation_method] || d.depreciation_method}</td>
                  <td className="p-3 text-right">Rp {d.cost.toLocaleString()}</td>
                  <td className="p-3 text-right">Rp {d.salvage_value.toLocaleString()}</td>
                  <td className="p-3 text-right">Rp {d.accumulated_depreciation.toLocaleString()}</td>
                  <td className="p-3 text-right font-medium">Rp {(d.cost - d.accumulated_depreciation).toLocaleString()}</td>
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
                  <td colSpan={10} className="p-4 text-center text-gray-400">Belum ada aset tetap. Klik "+ Aset Baru" untuk mencatat.</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
