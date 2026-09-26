import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Store, Eye, EyeOff } from 'lucide-react';
import { useStore } from '../store/useStore';
import { authService } from '../services/api';
import toast from 'react-hot-toast';

export function RegisterPage() {
  const navigate = useNavigate();
  const setUser = useStore((s) => s.setUser);
  const setToken = useStore((s) => s.setToken);
  const setOutletName = useStore((s) => s.setOutletName);

  const [merchantName, setMerchantName] = useState('');
  const [slug, setSlug] = useState('');
  const [username, setUsername] = useState('');
  const [ownerName, setOwnerName] = useState('');
  const [showPin, setShowPin] = useState(false);
  const [pin, setPin] = useState(['', '', '', '', '', '']);
  const [loading, setLoading] = useState(false);

  const handlePinChange = (index: number, value: string) => {
    if (value.length > 1) return;
    if (value && !/^\d$/.test(value)) return;
    const newPin = [...pin];
    newPin[index] = value;
    setPin(newPin);

    if (value && index < 5) {
      const next = document.getElementById(`reg-pin-${index + 1}`);
      next?.focus();
    }
  };

  const handleKeyDown = (index: number, e: React.KeyboardEvent) => {
    if (e.key === 'Backspace' && !pin[index] && index > 0) {
      const prev = document.getElementById(`reg-pin-${index - 1}`);
      prev?.focus();
    }
  };

  const handleRegister = async () => {
    if (!merchantName.trim()) { toast.error('Nama bisnis wajib diisi'); return; }
    if (!username.trim()) { toast.error('Username owner wajib diisi'); return; }
    const password = pin.join('');
    if (!/^\d{6}$/.test(password)) { toast.error('PIN harus tepat 6 digit angka'); return; }

    setLoading(true);
    try {
      const res = await authService.registerMerchant({
        merchant_name: merchantName.trim(),
        ...(slug.trim() ? { slug: slug.trim().toLowerCase() } : {}),
        username: username.trim(),
        owner_name: ownerName.trim() || username.trim(),
        password,
      });

      const { token, claims, tenant, user } = res.data.data;
      setToken(token);
      setUser({
        id: String(user?.id ?? claims?.user_id ?? ''),
        name: user?.name || ownerName || username,
        role: user?.role || claims?.role || 'admin',
        outletId: String(user?.outletId || user?.outlet_id || claims?.outlet_id || ''),
        tenantId: String(user?.tenantId ?? user?.tenant_id ?? claims?.tenant_id ?? ''),
        permissions: [],
      });
      setOutletName(tenant?.name || merchantName);
      toast.success('Bisnis berhasil didaftarkan. Selamat datang!');
      navigate('/dashboard');
    } catch (err: any) {
      toast.error(err?.response?.data?.error || 'Pendaftaran gagal');
    }
    setLoading(false);
  };

  return (
    <div className="min-h-screen bg-neutral-50 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <div className="inline-flex items-center justify-center w-20 h-20 rounded-2xl bg-brand mb-4">
            <Store className="w-10 h-10 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-neutral-900">Daftarkan Bisnis Anda</h1>
          <p className="text-neutral-500 mt-2">Buat akun merchant baru dan langsung kelola toko</p>
        </div>

        <div className="pos-card-elevated p-8">
          <div className="mb-4">
            <label className="block text-sm font-medium text-neutral-700 mb-2">Nama Bisnis *</label>
            <input
              type="text"
              value={merchantName}
              onChange={(e) => setMerchantName(e.target.value)}
              className="pos-input"
              placeholder="cth: Toko Maju Jaya"
            />
          </div>

          <div className="mb-4">
            <label className="block text-sm font-medium text-neutral-700 mb-2">
              Slug <span className="text-neutral-400 font-normal">(opsional, dibuat otomatis bila kosong)</span>
            </label>
            <input
              type="text"
              value={slug}
              onChange={(e) => setSlug(e.target.value)}
              className="pos-input"
              placeholder="cth: toko-maju-jaya"
            />
          </div>

          <div className="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-neutral-700 mb-2">Username Owner *</label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="pos-input"
                placeholder="cth: budi"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-neutral-700 mb-2">Nama Owner</label>
              <input
                type="text"
                value={ownerName}
                onChange={(e) => setOwnerName(e.target.value)}
                className="pos-input"
                placeholder="cth: Budi Santoso"
              />
            </div>
          </div>

          <div className="mb-6">
            <label className="block text-sm font-medium text-neutral-700 mb-2">PIN Kasir (6 digit angka) *</label>
            <div className="flex justify-center gap-3">
              {pin.map((digit, i) => (
                <input
                  key={i}
                  id={`reg-pin-${i}`}
                  type={showPin ? 'text' : 'password'}
                  maxLength={1}
                  value={digit}
                  onChange={(e) => handlePinChange(i, e.target.value)}
                  onKeyDown={(e) => handleKeyDown(i, e)}
                  className="w-12 h-14 text-center text-xl font-bold bg-neutral-100 border border-neutral-300 rounded-xl text-neutral-900 focus:outline-none focus:ring-2 focus:ring-brand/50 focus:border-brand"
                  inputMode="numeric"
                />
              ))}
            </div>
            <button
              onClick={() => setShowPin(!showPin)}
              className="mt-2 text-xs text-neutral-500 hover:text-neutral-700 flex items-center justify-center gap-1 w-full"
            >
              {showPin ? <EyeOff className="w-3 h-3" /> : <Eye className="w-3 h-3" />}
              {showPin ? 'Sembunyikan' : 'Tampilkan'}
            </button>
          </div>

          <button
            onClick={handleRegister}
            disabled={loading}
            className="w-full py-3 bg-brand hover:bg-brand/90 text-white font-medium rounded-xl transition-colors disabled:opacity-50 focus:ring-2 focus:ring-brand/50 focus:ring-offset-2"
          >
            {loading ? 'Mendaftarkan...' : 'Daftar & Mulai Berjualan'}
          </button>

          <p className="text-center text-sm text-neutral-500 mt-4">
            Sudah punya akun?{' '}
            <Link to="/login" className="text-brand font-medium hover:underline">
              Masuk di sini
            </Link>
          </p>
        </div>

        <p className="text-center text-neutral-500 text-xs mt-6">
          POS System v1.0 &mdash; High Performance Point of Sale
        </p>
      </div>
    </div>
  );
}
