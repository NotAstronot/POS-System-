import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { adminService, FixedAsset, FxDifference, PeriodClosing } from '../services/api';

export default function PeriodEndPage() {
  const [periods, setPeriods] = useState<PeriodClosing[]>([]);
  const [assets, setAssets] = useState<FixedAsset[]>([]);
  const [fxList, setFxList] = useState<FxDifference[]>([]);
  const [loading, setLoading] = useState(true);
  const [newPeriod, setNewPeriod] = useState('');
  const [fxForm, setFxForm] = useState({ currency: 'USD', description: '', gain_amount: '', loss_amount: '' });
  const [running, setRunning] = useState(false);

  const activePeriod = periods.find((p) => p.status === 'open') || null;

  const load = async () => {
    try {
      const [periodsRes, assetsRes] = await Promise.all([
        adminService.listPeriodEnd(),
        adminService.listFixedAssets(),
      ]);
      setPeriods(periodsRes.data.data);
      setAssets(assetsRes.data.data);
      const open = periodsRes.data.data.find((p) => p.status === 'open');
      if (open) {
        const fxRes = await adminService.listFxDifferences(open.id);
        setFxList(fxRes.data.data);
      } else {
        setFxList([]);
      }
    } catch {} finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const handleOpenPeriod = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await adminService.openPeriod(newPeriod);
      setNewPeriod('');
      setFxForm({ currency: 'USD', description: '', gain_amount: '', loss_amount: '' });
      await load();
    } catch {}
  };

  const handleRunDepreciation = async () => {
    if (!activePeriod) return;
    if (!confirm('Hitung dan catat depresiasi aset tetap untuk periode ini ke jurnal?')) return;
    setRunning(true);
    try {
      const res = await adminService.runDepreciation(activePeriod.id);
      alert(`Depresiasi tercatat: ${res.data.data.assets_count} aset, total Rp ${res.data.data.depreciation_total.toLocaleString()}`);
      await load();
    } catch {} finally {
      setRunning(false);
    }
  };

  const handleAddFx = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!activePeriod) return;
    try {
      await adminService.createFxDifference({
        period: activePeriod.period,
        currency: fxForm.currency,
        description: fxForm.description,
        gain_amount: parseFloat(fxForm.gain_amount) || 0,
        loss_amount: parseFloat(fxForm.loss_amount) || 0,
      });
      setFxForm({ currency: 'USD', description: '', gain_amount: '', loss_amount: '' });
      await load();
    } catch {}
  };

  const handleClosePeriod = async () => {
    if (!activePeriod) return;
    if (!confirm('Tutup buku periode ini? Jurnal penutupan akan dibuat otomatis.')) return;
    setRunning(true);
    try {
      await adminService.closePeriod(activePeriod.id);
      await load();
    } catch {} finally {
      setRunning(false);
    }
  };

  const monthlyDepreciation = (a: FixedAsset) =>
    a.useful_life_months > 0 ? (a.cost - a.salvage_value) / a.useful_life_months : 0;

  const fxGainTotal = fxList.reduce((s, f) => s + f.gain_amount, 0);
  const fxLossTotal = fxList.reduce((s, f) => s + f.loss_amount, 0);

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold">Proses Akhir Bulan</h1>
          <p className="text-sm text-gray-500">Depresiasi aset tetap, selisih kurs valas, dan penutupan buku bulanan</p>
        </div>
      </div>

      {!activePeriod && (
        <form onSubmit={handleOpenPeriod} className="bg-white p-4 rounded shadow mb-6 flex items-end gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">Buka Periode Baru (bulan)</label>
            <input
              type="month"
              value={newPeriod}
              onChange={(e) => setNewPeriod(e.target.value)}
              className="border p-2 rounded"
              required
            />
          </div>
          <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">Buka Periode</button>
          <p className="text-sm text-gray-500 pb-2">Periode dibuka otomatis untuk bulan yang dipilih. Hanya satu periode aktif.</p>
        </form>
      )}

      {loading ? (
        <p>Loading...</p>
      ) : activePeriod ? (
        <div className="space-y-6">
          <div className="bg-blue-50 p-4 rounded flex justify-between items-center">
            <div>
              <div className="text-sm text-blue-700">Periode Aktif</div>
              <div className="text-xl font-bold text-blue-800">{activePeriod.period}</div>
            </div>
            <div className="flex gap-4 text-sm text-blue-700">
              <div>Depresiasi: <strong>Rp {activePeriod.depreciation_total.toLocaleString()}</strong></div>
              <div>Selisih Kurs (untung): <strong>Rp {activePeriod.fx_gain.toLocaleString()}</strong></div>
              <div>Selisih Kurs (rugi): <strong>Rp {activePeriod.fx_loss.toLocaleString()}</strong></div>
            </div>
          </div>

          <div className="bg-white rounded shadow overflow-hidden">
            <div className="flex justify-between items-center px-4 py-3 bg-gray-100">
              <h2 className="font-bold">Depresiasi Aset Tetap</h2>
              <div className="flex gap-2">
                <Link to="/fixed-assets" className="text-blue-600 text-sm hover:underline self-center">Kelola Aset</Link>
                <Link to="/depreciation" className="text-blue-600 text-sm hover:underline self-center">Jadwal Penyusutan</Link>
                <button
                  onClick={handleRunDepreciation}
                  disabled={running}
                  className="bg-blue-600 text-white px-4 py-2 rounded disabled:opacity-50"
                >
                  {running ? 'Memproses...' : 'Hitung & Catat Depresiasi'}
                </button>
              </div>
            </div>
            <div className="px-4 py-2 text-sm text-gray-500">
              Dihitung otomatis dari aset tetap aktif sesuai metode masing-masing (garis lurus / saldo menurun). Hasil dicatat ke Jurnal Umum.
            </div>
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="p-3 text-left">Kode</th>
                  <th className="p-3 text-left">Nama Aset</th>
                  <th className="p-3 text-right">Harga Perolehan</th>
                  <th className="p-3 text-right">Nilai Sisa</th>
                  <th className="p-3 text-right">Umur (bulan)</th>
                  <th className="p-3 text-right">Akumulasi Depresiasi</th>
                  <th className="p-3 text-right">Depresiasi / Bulan</th>
                </tr>
              </thead>
              <tbody>
                {assets.filter((a) => a.is_active).map((a) => (
                  <tr key={a.id} className="border-t">
                    <td className="p-3 font-medium">{a.asset_code}</td>
                    <td className="p-3">{a.name}</td>
                    <td className="p-3 text-right">Rp {a.cost.toLocaleString()}</td>
                    <td className="p-3 text-right">Rp {a.salvage_value.toLocaleString()}</td>
                    <td className="p-3 text-right">{a.useful_life_months}</td>
                    <td className="p-3 text-right">Rp {a.accumulated_depreciation.toLocaleString()}</td>
                    <td className="p-3 text-right font-medium">Rp {monthlyDepreciation(a).toLocaleString()}</td>
                  </tr>
                ))}
                {assets.filter((a) => a.is_active).length === 0 && (
                  <tr>
                    <td colSpan={7} className="p-4 text-center text-gray-400">
                      Belum ada aset tetap.{' '}
                      <Link to="/fixed-assets" className="text-blue-600 hover:underline">Catat aset lewat menu Pencatatan Aset</Link>
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          <div className="bg-white rounded shadow overflow-hidden">
            <div className="px-4 py-3 bg-gray-100">
              <h2 className="font-bold">Selisih Kurs Valas</h2>
            </div>
            <form onSubmit={handleAddFx} className="grid grid-cols-5 gap-4 p-4 bg-gray-50">
              <div>
                <label className="block text-sm font-medium mb-1">Mata Uang</label>
                <input
                  value={fxForm.currency}
                  onChange={(e) => setFxForm({ ...fxForm, currency: e.target.value })}
                  className="border p-2 rounded w-full"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Keterangan</label>
                <input
                  placeholder="mis. Revaluasi saldo USD"
                  value={fxForm.description}
                  onChange={(e) => setFxForm({ ...fxForm, description: e.target.value })}
                  className="border p-2 rounded w-full"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Keuntungan Kurs (Rp)</label>
                <input
                  type="number"
                  value={fxForm.gain_amount}
                  onChange={(e) => setFxForm({ ...fxForm, gain_amount: e.target.value })}
                  className="border p-2 rounded w-full"
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Kerugian Kurs (Rp)</label>
                <input
                  type="number"
                  value={fxForm.loss_amount}
                  onChange={(e) => setFxForm({ ...fxForm, loss_amount: e.target.value })}
                  className="border p-2 rounded w-full"
                />
              </div>
              <div className="flex items-end">
                <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">+ Tambah</button>
              </div>
            </form>
            <table className="w-full">
              <thead className="bg-gray-50">
                <tr>
                  <th className="p-3 text-left">Mata Uang</th>
                  <th className="p-3 text-left">Keterangan</th>
                  <th className="p-3 text-right">Keuntungan</th>
                  <th className="p-3 text-right">Kerugian</th>
                </tr>
              </thead>
              <tbody>
                {fxList.map((f) => (
                  <tr key={f.id} className="border-t">
                    <td className="p-3 font-medium">{f.currency}</td>
                    <td className="p-3">{f.description || '-'}</td>
                    <td className="p-3 text-right text-green-600">{f.gain_amount > 0 ? `Rp ${f.gain_amount.toLocaleString()}` : '-'}</td>
                    <td className="p-3 text-right text-red-600">{f.loss_amount > 0 ? `Rp ${f.loss_amount.toLocaleString()}` : '-'}</td>
                  </tr>
                ))}
                {fxList.length === 0 && (
                  <tr>
                    <td colSpan={4} className="p-4 text-center text-gray-400">Belum ada selisih kurs pada periode ini</td>
                  </tr>
                )}
                {fxList.length > 0 && (
                  <tr className="bg-gray-50 font-semibold">
                    <td className="p-3">Total</td>
                    <td></td>
                    <td className="p-3 text-right text-green-600">Rp {fxGainTotal.toLocaleString()}</td>
                    <td className="p-3 text-right text-red-600">Rp {fxLossTotal.toLocaleString()}</td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          <div className="bg-white rounded shadow p-4 flex justify-between items-center">
            <div>
              <h2 className="font-bold">Penutupan Buku Bulanan</h2>
              <p className="text-sm text-gray-500">
                Jurnal penutupan dibuat otomatis: menutup akun Pendapatan dan Beban ke akun Ekuitas (laba ditahan).
              </p>
            </div>
            <button
              onClick={handleClosePeriod}
              disabled={running}
              className="bg-green-600 text-white px-4 py-2 rounded disabled:opacity-50"
            >
              {running ? 'Memproses...' : 'Tutup Buku Periode Ini'}
            </button>
          </div>
        </div>
      ) : (
        <div className="bg-white rounded shadow p-8 text-center text-gray-400">
          Buka periode bulanan terlebih dahulu untuk menjalankan proses akhir bulan.
        </div>
      )}

      <div className="bg-white rounded shadow overflow-hidden mt-6">
        <div className="px-4 py-3 bg-gray-100">
          <h2 className="font-bold">Riwayat Periode</h2>
        </div>
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="p-3 text-left">Periode</th>
              <th className="p-3 text-right">Depresiasi</th>
              <th className="p-3 text-right">Selisih Kurs (untung)</th>
              <th className="p-3 text-right">Selisih Kurs (rugi)</th>
              <th className="p-3 text-left">Status</th>
              <th className="p-3 text-left">Ditutup Pada</th>
            </tr>
          </thead>
          <tbody>
            {periods.map((p) => (
              <tr key={p.id} className="border-t">
                <td className="p-3 font-medium">{p.period}</td>
                <td className="p-3 text-right">Rp {p.depreciation_total.toLocaleString()}</td>
                <td className="p-3 text-right text-green-600">Rp {p.fx_gain.toLocaleString()}</td>
                <td className="p-3 text-right text-red-600">Rp {p.fx_loss.toLocaleString()}</td>
                <td className="p-3">
                  <span className={`px-2 py-1 rounded text-xs ${p.status === 'open' ? 'bg-yellow-100 text-yellow-700' : 'bg-green-100 text-green-700'}`}>
                    {p.status === 'open' ? 'Terbuka' : 'Ditutup'}
                  </span>
                </td>
                <td className="p-3 text-sm">{p.closed_at ? new Date(p.closed_at).toLocaleString('id-ID') : '-'}</td>
              </tr>
            ))}
            {periods.length === 0 && (
              <tr>
                <td colSpan={6} className="p-4 text-center text-gray-400">Belum ada periode</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
