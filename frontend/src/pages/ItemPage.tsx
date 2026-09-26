import React, { useEffect, useState } from 'react';
import { Upload, Download, FileSpreadsheet, X } from 'lucide-react';
import * as XLSX from 'xlsx';
import { inventoryService, Item, ItemVariant, categoryService } from '../services/api';

interface CategoryOpt {
  id: string;
  name: string;
}

interface ItemForm {
  code: string;
  name: string;
  barcode: string;
  description: string;
  category_id: string;
  purchase_price: string;
  base_price: string;
  unit: string;
  min_stock: string;
  is_service: boolean;
  has_variants: boolean;
  tax_rate: string;
  initial_stock: string;
}

const emptyForm: ItemForm = {
  code: '', name: '', barcode: '', description: '', category_id: '',
  purchase_price: '', base_price: '', unit: 'pcs', min_stock: '0',
  is_service: false, has_variants: false, tax_rate: '11', initial_stock: '0',
};

const emptyVariant: ItemVariant = { name: '', barcode: '', additional_price: 0, stock: 0, is_active: true };

export default function ItemPage() {
  const [data, setData] = useState<Item[]>([]);
  const [categories, setCategories] = useState<CategoryOpt[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<Item | null>(null);
  const [form, setForm] = useState<ItemForm>({ ...emptyForm });
  const [variants, setVariants] = useState<ItemVariant[]>([{ ...emptyVariant }]);
  const [expanded, setExpanded] = useState<number | null>(null);
  const [details, setDetails] = useState<Record<number, Item>>({});
  const [error, setError] = useState('');

  // ---- Bulk Upload Excel/CSV ----
  interface BulkRow {
    data: Record<string, any>;
    rowNum: number;
    error?: string;
  }
  const [showBulk, setShowBulk] = useState(false);
  const [bulkRows, setBulkRows] = useState<BulkRow[]>([]);
  const [bulkFileName, setBulkFileName] = useState('');
  const [importing, setImporting] = useState(false);
  const [importProgress, setImportProgress] = useState({ done: 0, total: 0 });
  const [importResult, setImportResult] = useState<{ ok: number; fail: BulkRow[] } | null>(null);

  const normKey = (h: any) => String(h ?? '').toLowerCase().trim().replace(/[\s_]+/g, '');
  const colMap: Record<string, string> = {
    code: 'code', kode: 'code',
    name: 'name', nama: 'name', namabarang: 'name',
    barcode: 'barcode',
    description: 'description', deskripsi: 'description',
    category: 'category', kategori: 'category',
    unit: 'unit', satuan: 'unit',
    purchaseprice: 'purchase_price', hpp: 'purchase_price', hargabeli: 'purchase_price',
    baseprice: 'base_price', hargajual: 'base_price', harga: 'base_price',
    minstock: 'min_stock', stokmin: 'min_stock', stokminimal: 'min_stock',
    initialstock: 'initial_stock', stokawal: 'initial_stock', stok: 'initial_stock',
    taxrate: 'tax_rate', pajak: 'tax_rate',
    isservice: 'is_service', jasa: 'is_service',
  };
  const numVal = (v: any) => {
    const n = parseFloat(String(v ?? '').replace(/[^0-9.\-]/g, ''));
    return isNaN(n) ? 0 : n;
  };
  const boolVal = (v: any) => /^(true|1|ya|y|jasa)$/i.test(String(v ?? '').trim());

  const downloadTemplate = () => {
    const ws = XLSX.utils.json_to_sheet([
      { code: 'BRG-001', name: 'Kopi Susu 250ml', barcode: '8991234567890', description: 'Kopi susu gula aren', category: 'Minuman', unit: 'pcs', purchase_price: 8000, base_price: 15000, min_stock: 10, initial_stock: 100, tax_rate: 11, is_service: 'FALSE' },
      { code: 'JSA-001', name: 'Cuci Motor', barcode: '', description: 'Jasa cuci motor', category: 'Jasa', unit: 'pcs', purchase_price: 0, base_price: 25000, min_stock: 0, initial_stock: 0, tax_rate: 0, is_service: 'TRUE' },
    ]);
    ws['!cols'] = [{ wch: 12 }, { wch: 24 }, { wch: 18 }, { wch: 24 }, { wch: 14 }, { wch: 8 }, { wch: 14 }, { wch: 12 }, { wch: 10 }, { wch: 13 }, { wch: 10 }, { wch: 11 }];
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Produk');
    XLSX.writeFile(wb, 'template-bulk-produk.xlsx');
  };

  const handleBulkFile = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    e.target.value = '';
    if (!file) return;
    setBulkFileName(file.name);
    setImportResult(null);
    const reader = new FileReader();
    reader.onload = () => {
      try {
        const wb = XLSX.read(reader.result, { type: 'array' });
        const ws = wb.Sheets[wb.SheetNames[0]];
        const aoa: any[][] = XLSX.utils.sheet_to_json(ws, { header: 1, defval: '' });
        if (aoa.length < 2) { setBulkRows([]); setError('File kosong (butuh baris header + minimal 1 data)'); return; }
        const headers = aoa[0].map(normKey).map((h) => colMap[h] || '');
        const rows: BulkRow[] = [];
        for (let i = 1; i < aoa.length; i++) {
          const data: Record<string, any> = {};
          let empty = true;
          headers.forEach((key, j) => {
            if (!key) return;
            const v = aoa[i][j];
            data[key] = typeof v === 'string' ? v.trim() : v;
            if (data[key] !== '' && data[key] !== undefined) empty = false;
          });
          if (empty) continue;
          const row: BulkRow = { data, rowNum: i + 1 };
          if (!data.name) row.error = 'Nama produk wajib diisi';
          rows.push(row);
        }
        setBulkRows(rows);
      } catch {
        setBulkRows([]);
        setError('Gagal membaca file. Gunakan format .xlsx atau .csv');
      }
    };
    reader.readAsArrayBuffer(file);
  };

  const handleBulkImport = async () => {
    const valid = bulkRows.filter((r) => !r.error);
    if (valid.length === 0) { setError('Tidak ada baris valid untuk diimport'); return; }
    setImporting(true);
    setImportProgress({ done: 0, total: valid.length });
    const catCache = new Map<string, number | null>();
    categories.forEach((c) => catCache.set(c.name.toLowerCase(), Number(c.id)));
    let ok = 0;
    const fail: BulkRow[] = [];
    for (let i = 0; i < valid.length; i++) {
      const row = valid[i];
      try {
        const d = row.data;
        let categoryId: number | null = null;
        const catName = String(d.category ?? '').trim();
        if (catName) {
          const key = catName.toLowerCase();
          if (catCache.has(key)) {
            categoryId = catCache.get(key) ?? null;
          } else {
            const res = await categoryService.create(catName);
            categoryId = Number((res.data.data as any).id);
            catCache.set(key, categoryId);
          }
        }
        await inventoryService.createItem({
          code: String(d.code ?? ''),
          name: String(d.name ?? ''),
          barcode: String(d.barcode ?? ''),
          description: String(d.description ?? ''),
          category_id: categoryId,
          purchase_price: numVal(d.purchase_price),
          base_price: numVal(d.base_price),
          unit: String(d.unit ?? 'pcs') || 'pcs',
          min_stock: Math.round(numVal(d.min_stock)),
          is_service: boolVal(d.is_service),
          has_variants: false,
          tax_rate: numVal(d.tax_rate),
          initial_stock: Math.round(numVal(d.initial_stock)),
          variants: [],
        });
        ok++;
      } catch (err: any) {
        fail.push({ ...row, error: err.response?.data?.error || err.message || 'Gagal import' });
      }
      setImportProgress({ done: i + 1, total: valid.length });
    }
    setImporting(false);
    setImportResult({ ok, fail });
    load();
  };

  const load = async () => {
    setError('');
    try {
      const [iRes, cRes] = await Promise.all([
        inventoryService.listItems(),
        categoryService.list(''),
      ]);
      setData(iRes.data.data ?? []);
      setCategories(cRes.data.data ?? []);
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal memuat data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    load();
  }, []);

  const openAdd = () => {
    setEditItem(null);
    setForm({ ...emptyForm });
    setVariants([{ ...emptyVariant }]);
    setShowForm(true);
  };

  const openEdit = (it: Item) => {
    setEditItem(it);
    setForm({
      code: it.code,
      name: it.name,
      barcode: it.barcode,
      description: it.description,
      category_id: it.category_id ? String(it.category_id) : '',
      purchase_price: String(it.purchase_price || ''),
      base_price: String(it.base_price || ''),
      unit: it.unit,
      min_stock: String(it.min_stock || 0),
      is_service: it.is_service,
      has_variants: it.has_variants,
      tax_rate: String(it.tax_rate || 0),
      initial_stock: '0',
    });
    setVariants(it.variants.length > 0 ? it.variants.map(v => ({ ...v })) : [{ ...emptyVariant }]);
    setShowForm(true);
  };

  const handleVariantChange = (index: number, field: keyof ItemVariant, value: any) => {
    const updated = [...variants];
    (updated[index] as any)[field] = value;
    setVariants(updated);
  };

  const addVariant = () => setVariants([...variants, { ...emptyVariant }]);
  const removeVariant = (index: number) => setVariants(variants.filter((_, i) => i !== index));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const payload: any = {
        code: form.code,
        name: form.name,
        barcode: form.barcode,
        description: form.description,
        category_id: form.category_id ? parseInt(form.category_id) : null,
        purchase_price: parseFloat(form.purchase_price) || 0,
        base_price: parseFloat(form.base_price) || 0,
        unit: form.unit,
        min_stock: parseInt(form.min_stock) || 0,
        is_service: form.is_service,
        has_variants: form.has_variants,
        tax_rate: parseFloat(form.tax_rate) || 0,
        initial_stock: parseFloat(form.initial_stock) || 0,
        variants: form.has_variants
          ? variants.filter(v => v.name).map(v => ({
              name: v.name,
              barcode: v.barcode,
              additional_price: v.additional_price || 0,
              stock: v.stock || 0,
              is_active: v.is_active !== false,
            }))
          : [],
      };
      if (editItem) {
        await inventoryService.updateItem(editItem.id, payload);
      } else {
        await inventoryService.createItem(payload);
      }
      setShowForm(false);
      load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal menyimpan data');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Hapus barang ini?')) return;
    try {
      await inventoryService.deleteItem(id);
      load();
    } catch (err: any) {
      setError(err.response?.data?.error || err.message || 'Gagal menghapus data');
    }
  };

  const toggleDetail = async (id: number) => {
    if (expanded === id) {
      setExpanded(null);
      return;
    }
    if (!details[id]) {
      try {
        const res = await inventoryService.getItem(id);
        setDetails(prev => ({ ...prev, [id]: res.data.data }));
      } catch {}
    }
    setExpanded(id);
  };

  const formatRupiah = (n: number) => `Rp ${(n || 0).toLocaleString('id-ID')}`;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Barang & Jasa</h1>
        <div className="flex gap-2">
          <button onClick={downloadTemplate} className="bg-white border border-gray-300 text-gray-700 px-4 py-2 rounded flex items-center gap-2 hover:bg-gray-50">
            <Download className="w-4 h-4" /> Template
          </button>
          <button onClick={() => { setShowBulk(true); setBulkRows([]); setImportResult(null); setBulkFileName(''); }} className="bg-white border border-gray-300 text-gray-700 px-4 py-2 rounded flex items-center gap-2 hover:bg-gray-50">
            <FileSpreadsheet className="w-4 h-4" /> Bulk Upload
          </button>
          <button onClick={openAdd} className="bg-blue-600 text-white px-4 py-2 rounded">+ Barang</button>
        </div>
      </div>

      {error && <div className="mb-4 p-4 text-red-500 bg-red-50 border border-red-200 rounded">{error}</div>}

      {showForm && (
        <div className="bg-white p-4 rounded shadow mb-6 border border-gray-200">
          <h2 className="font-bold mb-4">{editItem ? 'Edit Barang & Jasa' : 'Tambah Barang & Jasa'}</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-3 gap-4">
              <input value={form.code} onChange={e => setForm({ ...form, code: e.target.value })} className="border p-2 rounded" placeholder="Kode (mis. BRG-001)" />
              <input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} className="border p-2 rounded" placeholder="Nama Barang *" required />
              <input value={form.barcode} onChange={e => setForm({ ...form, barcode: e.target.value })} className="border p-2 rounded" placeholder="Barcode" />
              <select value={form.category_id} onChange={e => setForm({ ...form, category_id: e.target.value })} className="border p-2 rounded">
                <option value="">-- Kategori --</option>
                {categories.map(c => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
              <input value={form.unit} onChange={e => setForm({ ...form, unit: e.target.value })} className="border p-2 rounded" placeholder="Satuan (pcs/kg/box)" />
              <input value={form.tax_rate} onChange={e => setForm({ ...form, tax_rate: e.target.value })} className="border p-2 rounded" placeholder="Pajak (%)" />
            </div>
            <input value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} className="border p-2 rounded w-full" placeholder="Deskripsi" />

            <div className="grid grid-cols-4 gap-4">
              <div>
                <label className="block text-sm font-medium mb-1">Harga Beli</label>
                <input type="number" value={form.purchase_price} onChange={e => setForm({ ...form, purchase_price: e.target.value })} className="border p-2 rounded w-full" placeholder="0" />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Harga Jual</label>
                <input type="number" value={form.base_price} onChange={e => setForm({ ...form, base_price: e.target.value })} className="border p-2 rounded w-full" placeholder="0" required />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Stok Minimal</label>
                <input type="number" value={form.min_stock} onChange={e => setForm({ ...form, min_stock: e.target.value })} className="border p-2 rounded w-full" placeholder="0" />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">Stok Awal {editItem ? '' : '(gudang utama)'}</label>
                <input type="number" value={form.initial_stock} disabled={!!editItem} onChange={e => setForm({ ...form, initial_stock: e.target.value })} className="border p-2 rounded w-full" placeholder="0" />
              </div>
            </div>

            <div className="flex gap-6">
              <label className="flex items-center gap-2">
                <input type="checkbox" checked={form.has_variants} onChange={e => setForm({ ...form, has_variants: e.target.checked })} />
                <span className="text-sm">Barang ber-varian (warna/ukuran)</span>
              </label>
              <label className="flex items-center gap-2">
                <input type="checkbox" checked={form.is_service} onChange={e => setForm({ ...form, is_service: e.target.checked })} />
                <span className="text-sm">Jasa (tidak dihitung stok)</span>
              </label>
            </div>

            {form.has_variants && (
              <div>
                <div className="flex justify-between items-center mb-2">
                  <h3 className="font-semibold text-sm">Varian</h3>
                  <button type="button" onClick={addVariant} className="text-blue-600 text-sm hover:underline">+ Tambah Varian</button>
                </div>
                <table className="w-full border">
                  <thead className="bg-gray-100">
                    <tr>
                      <th className="p-2 text-left border">Nama Varian</th>
                      <th className="p-2 text-left border">Barcode</th>
                      <th className="p-2 text-left border">Harga Tambahan</th>
                      <th className="p-2 text-left border">Stok</th>
                      <th className="p-2 text-center border w-16">Aksi</th>
                    </tr>
                  </thead>
                  <tbody>
                    {variants.map((v, idx) => (
                      <tr key={idx}>
                        <td className="p-1 border"><input value={v.name} onChange={e => handleVariantChange(idx, 'name', e.target.value)} className="border p-1 rounded w-full" /></td>
                        <td className="p-1 border"><input value={v.barcode} onChange={e => handleVariantChange(idx, 'barcode', e.target.value)} className="border p-1 rounded w-full" /></td>
                        <td className="p-1 border"><input type="number" value={v.additional_price} onChange={e => handleVariantChange(idx, 'additional_price', parseFloat(e.target.value) || 0)} className="border p-1 rounded w-full" /></td>
                        <td className="p-1 border"><input type="number" value={v.stock} onChange={e => handleVariantChange(idx, 'stock', parseInt(e.target.value) || 0)} className="border p-1 rounded w-full" /></td>
                        <td className="p-1 border text-center">
                          {variants.length > 1 && (
                            <button type="button" onClick={() => removeVariant(idx)} className="text-red-500 hover:underline text-sm">Hapus</button>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}

            <div className="flex gap-2">
              <button type="submit" className="bg-green-600 text-white px-4 py-2 rounded">Simpan</button>
              <button type="button" onClick={() => setShowForm(false)} className="bg-gray-400 text-white px-4 py-2 rounded">Batal</button>
            </div>
          </form>
        </div>
      )}

      {loading ? <p>Loading...</p> : (
        <div className="bg-white rounded shadow overflow-hidden">
          <table className="w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="p-3 text-left">Kode</th>
                <th className="p-3 text-left">Nama</th>
                <th className="p-3 text-left">Barcode</th>
                <th className="p-3 text-left">Kategori</th>
                <th className="p-3 text-right">Harga Beli</th>
                <th className="p-3 text-right">Harga Jual</th>
                <th className="p-3 text-right">Stok</th>
                <th className="p-3 text-left">Tipe</th>
                <th className="p-3 text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              {data.map(it => {
                const detail = details[it.id];
                const showDetail = detail || it;
                return (
                <React.Fragment key={it.id}>
                  <tr className="border-t">
                    <td className="p-3 font-medium">{it.code || '-'}</td>
                    <td className="p-3">{it.name}</td>
                    <td className="p-3">{it.barcode || '-'}</td>
                    <td className="p-3">{it.category_name || '-'}</td>
                    <td className="p-3 text-right">{formatRupiah(it.purchase_price)}</td>
                    <td className="p-3 text-right">{formatRupiah(it.base_price)}</td>
                    <td className="p-3 text-right font-bold">
                      {it.is_service ? '-' : it.stock}
                      {!it.is_service && it.min_stock > 0 && it.stock <= it.min_stock && (
                        <span className="ml-1 text-xs text-red-600">(min {it.min_stock})</span>
                      )}
                    </td>
                    <td className="p-3">
                      <span className={`px-2 py-1 rounded text-xs ${it.is_service ? 'bg-purple-100 text-purple-700' : it.has_variants ? 'bg-yellow-100 text-yellow-700' : 'bg-blue-100 text-blue-700'}`}>
                        {it.is_service ? 'Jasa' : it.has_variants ? 'Varian' : 'Barang'}
                      </span>
                    </td>
                    <td className="p-3 text-center gap-2 flex justify-center">
                      <button onClick={() => toggleDetail(it.id)} className="text-gray-500 hover:underline">Detail</button>
                      <button onClick={() => openEdit(it)} className="text-blue-500 hover:underline">Edit</button>
                      <button onClick={() => handleDelete(it.id)} className="text-red-500 hover:underline">Hapus</button>
                    </td>
                  </tr>
                  {expanded === it.id && (
                    <tr className="bg-gray-50">
                      <td colSpan={9} className="p-4">
                        <div className="grid grid-cols-2 gap-6">
                          <div>
                            <h4 className="font-semibold text-sm mb-2">Stok per Gudang</h4>
                            {showDetail.stock_by_warehouse.length === 0 && <p className="text-sm text-gray-400">Tidak ada data stok</p>}
                            <table className="w-full text-sm">
                              <tbody>
                                {showDetail.stock_by_warehouse.map(s => (
                                  <tr key={s.warehouse_id} className="border-t">
                                    <td className="py-1">{s.warehouse_name}</td>
                                    <td className="py-1 text-right font-medium">{s.quantity}</td>
                                  </tr>
                                ))}
                              </tbody>
                            </table>
                          </div>
                          {showDetail.has_variants && (
                            <div>
                              <h4 className="font-semibold text-sm mb-2">Varian</h4>
                              {showDetail.variants.length === 0 && <p className="text-sm text-gray-400">Tidak ada varian</p>}
                              <table className="w-full text-sm">
                                <tbody>
                                  {showDetail.variants.map(v => (
                                    <tr key={v.id} className="border-t">
                                      <td className="py-1">{v.name}</td>
                                      <td className="py-1 text-gray-500">{v.barcode || '-'}</td>
                                      <td className="py-1 text-right">+{formatRupiah(v.additional_price)}</td>
                                      <td className="py-1 text-right font-medium">{v.stock}</td>
                                    </tr>
                                  ))}
                                </tbody>
                              </table>
                            </div>
                          )}
                        </div>
                      </td>
                    </tr>
                  )}
                </React.Fragment>
                );
              })}
              {data.length === 0 && <tr><td colSpan={9} className="p-4 text-center text-gray-400">Belum ada barang</td></tr>}
            </tbody>
          </table>
        </div>
      )}

      {showBulk && (
        <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-xl w-full max-w-3xl max-h-[90vh] overflow-y-auto shadow-2xl">
            <div className="flex items-center justify-between px-6 py-4 border-b">
              <h2 className="font-bold flex items-center gap-2">
                <FileSpreadsheet className="w-5 h-5 text-green-600" /> Bulk Upload Produk (Excel/CSV)
              </h2>
              <button onClick={() => setShowBulk(false)} disabled={importing} className="p-2 rounded-lg hover:bg-gray-100 text-gray-500">
                <X className="w-5 h-5" />
              </button>
            </div>
            <div className="p-6 space-y-4">
              <div className="text-sm text-gray-600 bg-gray-50 border rounded-lg p-3">
                <p className="font-medium mb-1">Langkah:</p>
                <ol className="list-decimal ml-5 space-y-0.5">
                  <li>Unduh <button onClick={downloadTemplate} className="text-blue-600 hover:underline">template Excel</button> bila belum punya.</li>
                  <li>Isi data (kolom <b>name</b> wajib; harga & stok berupa angka polos).</li>
                  <li>Pilih file .xlsx/.csv, periksa preview, lalu klik Import.</li>
                </ol>
              </div>

              <label className="flex items-center justify-center gap-2 border-2 border-dashed border-gray-300 rounded-xl py-8 cursor-pointer hover:border-blue-500 transition-colors">
                <Upload className="w-5 h-5 text-gray-500" />
                <span className="text-sm text-gray-600">{bulkFileName || 'Klik untuk pilih file Excel/CSV'}</span>
                <input type="file" accept=".xlsx,.xls,.csv" onChange={handleBulkFile} className="hidden" disabled={importing} />
              </label>

              {bulkRows.length > 0 && (
                <div>
                  <p className="text-sm font-medium mb-2">
                    Preview: {bulkRows.length} baris
                    <span className="text-green-700"> ({bulkRows.filter((r) => !r.error).length} valid</span>
                    {bulkRows.some((r) => r.error) && <span className="text-red-600">, {bulkRows.filter((r) => r.error).length} bermasalah</span>})
                  </p>
                  <div className="overflow-x-auto border rounded-lg max-h-64 overflow-y-auto">
                    <table className="w-full text-sm">
                      <thead className="bg-gray-100 sticky top-0">
                        <tr>
                          <th className="p-2 text-left">Baris</th>
                          <th className="p-2 text-left">Nama</th>
                          <th className="p-2 text-right">Harga Jual</th>
                          <th className="p-2 text-right">HPP</th>
                          <th className="p-2 text-right">Stok Awal</th>
                          <th className="p-2 text-left">Status</th>
                        </tr>
                      </thead>
                      <tbody>
                        {bulkRows.slice(0, 50).map((r, i) => (
                          <tr key={i} className="border-t">
                            <td className="p-2">{r.rowNum}</td>
                            <td className="p-2">{String(r.data.name || '-')}</td>
                            <td className="p-2 text-right">{numVal(r.data.base_price).toLocaleString('id-ID')}</td>
                            <td className="p-2 text-right">{numVal(r.data.purchase_price).toLocaleString('id-ID')}</td>
                            <td className="p-2 text-right">{Math.round(numVal(r.data.initial_stock))}</td>
                            <td className="p-2">{r.error ? <span className="text-red-600 text-xs">{r.error}</span> : <span className="text-green-600 text-xs">Siap</span>}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                    {bulkRows.length > 50 && <p className="text-xs text-gray-400 p-2">... dan {bulkRows.length - 50} baris lainnya</p>}
                  </div>
                </div>
              )}

              {importing && (
                <div>
                  <div className="flex justify-between text-sm mb-1">
                    <span>Mengimport...</span>
                    <span>{importProgress.done}/{importProgress.total}</span>
                  </div>
                  <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
                    <div className="h-full bg-blue-600 transition-all" style={{ width: `${importProgress.total > 0 ? (importProgress.done / importProgress.total) * 100 : 0}%` }} />
                  </div>
                </div>
              )}

              {importResult && (
                <div className={`border rounded-lg p-3 text-sm ${importResult.fail.length === 0 ? 'bg-green-50 border-green-200 text-green-800' : 'bg-amber-50 border-amber-200 text-amber-800'}`}>
                  <p className="font-medium">Import selesai: {importResult.ok} berhasil{importResult.fail.length > 0 && `, ${importResult.fail.length} gagal`}.</p>
                  {importResult.fail.slice(0, 10).map((f, i) => (
                    <p key={i} className="text-xs mt-1">Baris {f.rowNum} ({String(f.data.name || '-')}) : {f.error}</p>
                  ))}
                </div>
              )}

              <div className="flex gap-2 justify-end">
                <button onClick={() => setShowBulk(false)} disabled={importing} className="bg-gray-200 text-gray-700 px-4 py-2 rounded disabled:opacity-50">Tutup</button>
                <button onClick={handleBulkImport} disabled={importing || bulkRows.filter((r) => !r.error).length === 0} className="bg-green-600 text-white px-4 py-2 rounded disabled:opacity-50">
                  {importing ? 'Mengimport...' : `Import ${bulkRows.filter((r) => !r.error).length} Produk`}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}