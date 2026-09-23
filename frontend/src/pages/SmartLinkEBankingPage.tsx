import React, { useEffect, useState } from 'react';
import toast from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';
import { smartlinkService, BankStatementImport, BankStatementLine } from '../services/api';

export default function SmartLinkEBankingPage() {
  const navigate = useNavigate();
  const [data, setData] = useState<BankStatementImport[]>([]);
  const [loading, setLoading] = useState(true);

  const [form, setForm] = useState({
    account_name: '', account_type: 'bank', statement_date: '',
    source: 'csv', file_name: '',
  });
  const [lines, setLines] = useState<BankStatementLine[]>([]);

  const load = async () => {
    setLoading(true);
    try {
      const res = await smartlinkService.listBankImports();
      setData(res.data.data);
    } catch {
      toast.error('Gagal memuat mutasi bank');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const addLine = () => {
    setLines([...lines, { transaction_date: form.statement_date, description: '', reference: '', amount: 0, type: 'debit' }]);
  };

  const updateLine = (index: number, field: keyof BankStatementLine, value: any) => {
    const next = [...lines];
    (next[index] as any)[field] = value;
    setLines(next);
  };

  const removeLine = (index: number) => {
    setLines(lines.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (lines.length === 0) return toast.error('Tambah minimal 1 baris mutasi');
    try {
      const res = await smartlinkService.importBankLines({
        ...form,
        lines: lines.map((l) => ({
          ...l,
          amount: parseFloat(String(l.amount)) || 0,
        })),
      });
      toast.success(`Mutasi ${res.data.data.total_rows} baris berhasil diimpor`);
      setForm({ account_name: '', account_type: 'bank', statement_date: '', source: 'csv', file_name: '' });
      setLines([]);
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal impor mutasi');
    }
  };

  const handleReconcile = async (item: BankStatementImport) => {
    const bookBalance = parseFloat(prompt(`Saldo buku (kas/bank) untuk ${item.account_name}?`, '0') || '0');
    try {
      await smartlinkService.reconcileBankImport(item.id, bookBalance);
      toast.success('Rekonsiliasi Bank berhasil dibuat dari mutasi');
      load();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Gagal membuat rekonsiliasi');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus impor mutasi ini?')) return;
    try {
      await smartlinkService.deleteBankImport(id);
      load();
    } catch {}
  };

  const balanceOf = (item: BankStatementImport) => {
    let bal = 0;
    for (const l of item.lines) {
      if (l.type === 'credit' || l.type === 'kredit') bal += l.amount;
      else bal -= l.amount;
    }
    return bal;
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-2">
        <h1 className="text-2xl font-bold">SmartLink e-Banking</h1>
      </div>
      <p className="text-sm text-gray-500 mb-6">
        Impor mutasi rekening bank secara otomatis untuk mempercepat proses rekonsiliasi.
        Mutasi yang diimpor dapat langsung dibuatkan Rekonsiliasi Bank di menu <b>Keuangan</b>.
      </p>

      <div className="bg-white p-4 rounded shadow mb-6">
        <h2 className="font-bold mb-4">Impor Mutasi Rekening</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium mb-1">Nama Rekening</label>
              <input placeholder="cth: BCA - 1234567890" value={form.account_name} onChange={(e) => setForm({ ...form, account_name: e.target.value })} className="border p-2 rounded w-full" required />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Tipe Rekening</label>
              <select value={form.account_type} onChange={(e) => setForm({ ...form, account_type: e.target.value })} className="border p-2 rounded w-full">
                <option value="bank">Bank</option>
                <option value="kas">Kas</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Tanggal Mutasi</label>
              <input type="date" value={form.statement_date} onChange={(e) => setForm({ ...form, statement_date: e.target.value })} className="border p-2 rounded w-full" />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Sumber</label>
              <select value={form.source} onChange={(e) => setForm({ ...form, source: e.target.value })} className="border p-2 rounded w-full">
                <option value="csv">CSV / Excel</option>
                <option value="api">API Bank</option>
                <option value="manual">Manual</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Nama File</label>
              <input placeholder="cth: mutasi_bca_202608.csv" value={form.file_name} onChange={(e) => setForm({ ...form, file_name: e.target.value })} className="border p-2 rounded w-full" />
            </div>
          </div>

          <div>
            <div className="flex justify-between items-center mb-2">
              <h3 className="font-medium">Baris Mutasi</h3>
              <button type="button" onClick={addLine} className="text-blue-600 text-sm hover:underline">+ Tambah Baris</button>
            </div>
            {lines.length > 0 && (
              <div className="border rounded overflow-hidden">
                <table className="w-full text-sm">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="p-2 text-left">Tanggal</th>
                      <th className="p-2 text-left">Deskripsi</th>
                      <th className="p-2 text-left">Referensi</th>
                      <th className="p-2 text-right">Jumlah</th>
                      <th className="p-2 text-left">Tipe</th>
                      <th className="p-2 text-center">Aksi</th>
                    </tr>
                  </thead>
                  <tbody>
                    {lines.map((l, idx) => (
                      <tr key={idx} className="border-t">
                        <td className="p-1">
                          <input type="date" value={l.transaction_date} onChange={(e) => updateLine(idx, 'transaction_date', e.target.value)} className="border p-1 rounded w-full text-sm" />
                        </td>
                        <td className="p-1">
                          <input placeholder="Deskripsi mutasi" value={l.description} onChange={(e) => updateLine(idx, 'description', e.target.value)} className="border p-1 rounded w-full text-sm" />
                        </td>
                        <td className="p-1">
                          <input placeholder="Referensi" value={l.reference} onChange={(e) => updateLine(idx, 'reference', e.target.value)} className="border p-1 rounded w-full text-sm" />
                        </td>
                        <td className="p-1">
                          <input type="number" placeholder="0" value={l.amount || ''} onChange={(e) => updateLine(idx, 'amount', e.target.value)} className="border p-1 rounded w-full text-sm text-right" />
                        </td>
                        <td className="p-1">
                          <select value={l.type} onChange={(e) => updateLine(idx, 'type', e.target.value)} className="border p-1 rounded w-full text-sm">
                            <option value="debit">Debit (Keluar)</option>
                            <option value="credit">Kredit (Masuk)</option>
                          </select>
                        </td>
                        <td className="p-1 text-center">
                          <button type="button" onClick={() => removeLine(idx)} className="text-red-500 hover:underline text-sm">Hapus</button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>

          <div className="flex gap-2">
            <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded">Impor Mutasi</button>
            <button type="button" onClick={() => navigate('/bank-reconciliations')} className="bg-gray-100 text-gray-700 px-4 py-2 rounded hover:bg-gray-200">
              Buka Rekonsiliasi Bank →
            </button>
          </div>
        </form>
      </div>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <div className="bg-white rounded shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-3 text-left">No. Impor</th>
                <th className="p-3 text-left">Rekening</th>
                <th className="p-3 text-left">Tanggal</th>
                <th className="p-3 text-left">Sumber</th>
                <th className="p-3 text-right">Jumlah Baris</th>
                <th className="p-3 text-right">Selisih Bersih</th>
                <th className="p-3 text-left">Status</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map((d) => (
                <tr key={d.id} className="border-t">
                  <td className="p-3 font-medium">{d.import_number}</td>
                  <td className="p-3">{d.account_name} <span className="text-xs text-gray-500">({d.account_type})</span></td>
                  <td className="p-3 text-sm">{d.statement_date}</td>
                  <td className="p-3 capitalize text-sm">{d.source}{d.file_name && <p className="text-xs text-gray-500">{d.file_name}</p>}</td>
                  <td className="p-3 text-right">{d.total_rows}</td>
                  <td className={`p-3 text-right font-medium ${balanceOf(d) >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                    {balanceOf(d) >= 0 ? '+' : ''}{balanceOf(d).toLocaleString()}
                  </td>
                  <td className="p-3">
                    <span className={`px-2 py-1 rounded text-xs ${d.lines.some((l) => l.is_reconciled) ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'}`}>
                      {d.lines.some((l) => l.is_reconciled) ? 'Sebagian Direkonsiliasi' : 'Terimpor'}
                    </span>
                  </td>
                  <td className="p-3 text-center gap-2 flex justify-center">
                    <button onClick={() => handleReconcile(d)} className="text-green-600 hover:underline font-medium">Buat Rekonsiliasi</button>
                    <button onClick={() => handleDelete(d.id)} className="text-red-500 hover:underline">Hapus</button>
                  </td>
                </tr>
              ))}
              {data.length === 0 && (
                <tr><td colSpan={8} className="p-4 text-center text-gray-400">Belum ada impor mutasi bank</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}