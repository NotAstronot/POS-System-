import React, { useEffect, useState } from 'react';
import { adminService, AssetDepreciationSchedule, PeriodClosing } from '../services/api';

const assetTypeLabels: Record<string, string> = {
  peralatan: 'Peralatan Toko',
  komputer: 'Komputer Kasir',
  kendaraan: 'Kendaraan Operasional',
  bangunan: 'Bangunan',
  lainnya: 'Lainnya',
};

const methodLabels: Record<string, string> = {
  garis_lurus: 'Garis Lurus',
  saldo_menurun: 'Saldo Menurun (Pajak)',
};

export default function PenyusutanOtomatisPage() {
  const [schedule, setSchedule] = useState<AssetDepreciationSchedule[]>([]);
  const [periods, setPeriods] = useState<PeriodClosing[]>([]);
  const [loading, setLoading] = useState(true);
  const [running, setRunning] = useState(false);

  const activePeriod = periods.find((p) => p.status === 'open') || null;

  const load = async () => {
    try {
      const [scheduleRes, periodsRes] = await Promise.all([
        adminService.listDepreciationSchedule(),
        adminService.listPeriodEnd(),
      ]);
      setSchedule(scheduleRes.data.data);
      setPeriods(periodsRes.data.data);
    } catch {} finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const currentMonth = () => {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
  };

  const handleRunDepreciation = async () => {
    if (running) return;
    setRunning(true);
    try {
      let period = activePeriod;
      if (!period) {
        await adminService.openPeriod(currentMonth());
        const periodsRes = await adminService.listPeriodEnd();
        setPeriods(periodsRes.data.data);
        period = periodsRes.data.data.find((p) => p.status === 'open') || null;
      }
      if (!period) {
        alert('Tidak ada periode yang terbuka. Buka periode lewat menu Proses Akhir Bulan.');
        return;
      }
      const res = await adminService.runDepreciation(period.id);
      alert(
        `Penyusutan tercatat untuk periode ${res.data.data.period}: ` +
        `${res.data.data.assets_count} aset, total Rp ${res.data.data.depreciation_total.toLocaleString()}. ` +
        `Jurnal otomatis masuk ke Jurnal Umum.`
      );
      await load();
    } catch {} finally {
      setRunning(false);
    }
  };

  const active = schedule.filter((s) => s.is_active && s.remaining_value > 0);
  const totalBookValue = active.reduce((s, a) => s + a.remaining_value, 0);
  const totalMonthly = active.reduce((s, a) => s + a.first_depreciation, 0);

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold">Penyusutan Otomatis</h1>
          <p className="text-sm text-gray-500">
            Perhitungan nilai penyusutan (depresiasi) aset secara otomatis sesuai metode akuntansi/pajak
          </p>
        </div>
        <button
          onClick={handleRunDepreciation}
          disabled={running}
          className="bg-blue-600 text-white px-4 py-2 rounded disabled:opacity-50"
        >
          {running ? 'Memproses...' : 'Jalankan Penyusutan Bulan Ini'}
        </button>
      </div>

      <div className="bg-blue-50 p-4 rounded mb-6 flex justify-between items-center">
        <div className="text-sm text-blue-700">
          Periode aktif: <strong>{activePeriod ? activePeriod.period : 'Belum ada (akan otomatis dibuka bulan ini)'}</strong>
        </div>
        <div className="text-sm text-blue-700">
          Total penyusutan bulan ini: <strong>Rp {totalMonthly.toLocaleString()}</strong>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="bg-white p-4 rounded shadow">
          <div className="text-sm text-gray-500">Aset Aktif (belum habis disusutkan)</div>
          <div className="text-xl font-bold mt-1">{active.length} aset</div>
        </div>
        <div className="bg-white p-4 rounded shadow">
          <div className="text-sm text-gray-500">Total Nilai Buku</div>
          <div className="text-xl font-bold mt-1">Rp {totalBookValue.toLocaleString()}</div>
        </div>
        <div className="bg-white p-4 rounded shadow">
          <div className="text-sm text-gray-500">Depresiasi Bulan Ini</div>
          <div className="text-xl font-bold mt-1">Rp {totalMonthly.toLocaleString()}</div>
        </div>
      </div>

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
                <th className="p-3 text-right">Akumulasi</th>
                <th className="p-3 text-right">Nilai Buku</th>
                <th className="p-3 text-right">Depresiasi Bulan Ini</th>
                <th className="p-3 text-right">Sisa Bulan</th>
              </tr>
            </thead>
            <tbody>
              {schedule.map((a) => (
                <tr key={a.id} className={`border-t ${!a.is_active || a.remaining_value <= 0 ? 'opacity-50' : ''}`}>
                  <td className="p-3 font-medium">{a.asset_code}</td>
                  <td className="p-3">{a.name}</td>
                  <td className="p-3 text-sm">{assetTypeLabels[a.asset_type] || a.asset_type}</td>
                  <td className="p-3 text-sm">{methodLabels[a.depreciation_method] || a.depreciation_method}</td>
                  <td className="p-3 text-right">Rp {a.cost.toLocaleString()}</td>
                  <td className="p-3 text-right">Rp {a.accumulated_depreciation.toLocaleString()}</td>
                  <td className="p-3 text-right font-medium">Rp {a.remaining_value.toLocaleString()}</td>
                  <td className="p-3 text-right font-medium text-blue-700">Rp {a.first_depreciation.toLocaleString()}</td>
                  <td className="p-3 text-right">{a.remaining_value > 0 && a.monthly_depreciation > 0 ? a.months_remaining.toFixed(1) : '-'}</td>
                </tr>
              ))}
              {schedule.length === 0 && (
                <tr>
                  <td colSpan={9} className="p-4 text-center text-gray-400">
                    Belum ada aset tetap. Catat aset lewat menu Pencatatan Aset.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      <div className="bg-yellow-50 border border-yellow-200 p-4 rounded mt-6 text-sm text-yellow-800">
        <strong>Metode Penyusutan:</strong> Garis Lurus = (harga perolehan - nilai sisa) / umur manfaat. Saldo Menurun (pajak) = 2 x nilai buku / umur manfaat.
        Hasil penyusutan otomatis dicatat sebagai jurnal (Beban Depresiasi ↔ Akumulasi Depresiasi) di <strong>Jurnal Umum</strong> dan masuk ke ringkasan <strong>Proses Akhir Bulan</strong>.
      </div>
    </div>
  );
}
