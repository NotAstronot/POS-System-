import React, { useEffect, useState } from 'react';
import toast from 'react-hot-toast';
import { smartlinkService, EFakturExport } from '../services/api';

const EXPORT_TYPES = [
  { value: 'ppn_keluaran', label: 'e-Faktur PPN Keluaran (Faktur Penjualan)' },
  { value: 'ppn_masukan', label: 'e-Faktur PPN Masukan (Faktur Pembelian)' },
  { value: 'pph23', label: 'Laporan PPh 23 (Pembelian / Vendor)' },
  { value: 'pph21', label: 'Laporan PPh 21 (Gaji Karyawan)' },
];

export default function SmartLinkTaxPage() {
  const [data, setData] = useState<EFakturExport[]>([]);
  const [loading, setLoading] = useState(true);
  const [detailItem, setDetailItem] = useState<EFakturExport | null>(null);
  const [form, setForm] = useState({ export_type: 'ppn_keluaran', period: new Date().toISOString().slice(0, 7) });

  const load = async () => {
    setLoading(true);
    try {
      const res = await smartlinkService.listTaxExports();
      setData(res.data.data);
    } catch {
      toast.error('Gagal memuat ekspor pajak');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const generate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await smartlinkService.generateTaxExport(form.export_type, form.period);
      toast.success(`Ekspor ${EXPORT_TYPES.find((t) => t.value === form.export_type)?.label} berhasil dibuat (${res.data.data.total_rows} baris)`);
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal membuat ekspor pajak');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus ekspor pajak ini?')) return;
    try {
      await smartlinkService.deleteTaxExport(id);
      load();
    } catch {}
  };

  const typeLabel = (t: string) => EXPORT_TYPES.find((x) => x.value === t)?.label || t;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-2">
        <h1 className="text-2xl font-bold">SmartLink Tax / e-Faktur</h1>
      </div>
      <p className="text-sm text-gray-500 mb-6">
        Integrasi perpajakan untuk mengekspor e-Faktur PPN atau laporan PPh secara praktis.
        Data diambil otomatis dari <b>Faktur Penjualan</b>, <b>Faktur Pembelian</b>, dan <b>Gaji</b>.
      </p>

      <div className="bg-white p-4 rounded shadow mb-6">
        <h2 className="font-bold mb-4">Buat Ekspor Pajak</h2>
        <form onSubmit={generate} className="grid grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium mb-1">Jenis Ekspor</label>
            <select value={form.export_type} onChange={(e) => setForm({ ...form, export_type: e.target.value })} className="border p-2 rounded w-full">
              {EXPORT_TYPES.map((t) => <option key={t.value} value={t.value}>{t.label}</option>)}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Periode (YYYY-MM)</label>
            <input type="month" value={form.period} onChange={(e) => setForm({ ...form, period: e.target.value })} className="border p-2 rounded w-full" required />
          </div>
          <div className="flex items-end">
            <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded w-full">Generate Ekspor</button>
          </div>
        </form>
      </div>

      {detailItem && (
        <div className="bg-white p-4 rounded shadow mb-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="font-bold">Detail: {detailItem.export_number} - {typeLabel(detailItem.export_type)}</h2>
            <button onClick={() => setDetailItem(null)} className="text-gray-500 hover:text-gray-700">Tutup</button>
          </div>
          <div className="grid grid-cols-3 gap-4 mb-4 text-sm">
            <div><span className="text-gray-500">Periode:</span> {detailItem.period}</div>
            <div><span className="text-gray-500">Total Baris:</span> {detailItem.total_rows}</div>
            <div><span className="text-gray-500">Total DPP:</span> Rp {detailItem.total_dpp.toLocaleString()}</div>
            <div><span className="text-gray-500">Total Pajak:</span> Rp {detailItem.total_tax.toLocaleString()}</div>
            {detailItem.file_path && (
              <div>
                <a href={detailItem.file_path} target="_blank" rel="noreferrer" className="text-blue-600 hover:underline">Download CSV</a>
              </div>
            )}
          </div>
          {detailItem.lines.length > 0 && (
            <table className="w-full text-sm">
              <thead className="bg-gray-50">
                <tr>
                  <th className="p-2 text-left">Referensi</th>
                  <th className="p-2 text-left">Nama</th>
                  <th className="p-2 text-left">NPWP</th>
                  <th className="p-2 text-left">Tanggal</th>
                  <th className="p-2 text-right">DPP</th>
                  <th className="p-2 text-right">Pajak</th>
                </tr>
              </thead>
              <tbody>
                {detailItem.lines.map((l, idx) => (
                  <tr key={idx} className="border-t">
                    <td className="p-2">{l.reference_number}</td>
                    <td className="p-2">{l.party_name}</td>
                    <td className="p-2">{l.npwp || '-'}</td>
                    <td className="p-2">{l.line_date || '-'}</td>
                    <td className="p-2 text-right">Rp {l.taxable_base.toLocaleString()}</td>
                    <td className="p-2 text-right">Rp {l.tax_amount.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}

      {loading ? (
        <p>Loading...</p>
      ) : (
        <div className="bg-white rounded shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-3 text-left">No. Ekspor</th>
                <th className="p-3 text-left">Jenis</th>
                <th className="p-3 text-left">Periode</th>
                <th className="p-3 text-right">Baris</th>
                <th className="p-3 text-right">Total DPP</th>
                <th className="p-3 text-right">Total Pajak</th>
                <th className="p-3 text-left">File</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map((d) => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.export_number}</td>
                  <td className="p-3 text-sm">{typeLabel(d.export_type)}</td>
                  <td className="p-3">{d.period}</td>
                  <td className="p-3 text-right">{d.total_rows}</td>
                  <td className="p-3 text-right">Rp {d.total_dpp.toLocaleString()}</td>
                  <td className="p-3 text-right font-medium">Rp {d.total_tax.toLocaleString()}</td>
                  <td className="p-3">
                    {d.file_path ? (
                      <a href={d.file_path} target="_blank" rel="noreferrer" className="text-blue-600 hover:underline text-sm">Download</a>
                    ) : '-'}
                  </td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => setDetailItem(d)} className="text-blue-500 hover:underline">Detail</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && (
                <tr><td colSpan={8} className="p-4 text-center text-gray-400">Belum ada ekspor pajak. Generate dari form di atas.</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}