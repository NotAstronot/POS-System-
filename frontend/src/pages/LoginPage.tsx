import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Store, Eye, EyeOff } from 'lucide-react';
import { useStore } from '../store/useStore';
import { authService } from '../services/api';
import toast from 'react-hot-toast';

export function LoginPage() {
  const navigate = useNavigate();
  const setUser = useStore((s) => s.setUser);
  const setToken = useStore((s) => s.setToken);
  const [showPin, setShowPin] = useState(false);
  const [pin, setPin] = useState(['', '', '', '', '', '']);
  const [username, setUsername] = useState('admin');
  const [loading, setLoading] = useState(false);

  const handlePinChange = (index: number, value: string) => {
    if (value.length > 1) return;
    const newPin = [...pin];
    newPin[index] = value;
    setPin(newPin);

    if (value && index < 5) {
      const next = document.getElementById(`pin-${index + 1}`);
      next?.focus();
    }
  };

  const handleKeyDown = (index: number, e: React.KeyboardEvent) => {
    if (e.key === 'Backspace' && !pin[index] && index > 0) {
      const prev = document.getElementById(`pin-${index - 1}`);
      prev?.focus();
    }
  };

  const handleLogin = async () => {
    setLoading(true);

    const password = pin.join('');

    try {
      const res = await authService.login({
        username,
        password,
        tenant_id: 'a0000000-0000-0000-0000-000000000001',
        outlet_id: 'b0000000-0000-0000-0000-000000000001',
      });

      const { token, user } = res.data.data;
      setToken(token);
      setUser({ ...user, permissions: user.permissions || [] });
      useStore.getState().setOutletName('Outlet Pusat');
      toast.success('Login berhasil!');
      navigate('/pos');
    } catch (err: any) {
      toast.error(err?.response?.data?.error || 'Login gagal');
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
          <h1 className="text-3xl font-bold text-neutral-900">POS System</h1>
          <p className="text-neutral-500 mt-2">Masuk ke sistem kasir</p>
        </div>

        <div className="pos-card-elevated p-8">
          <div className="mb-6">
            <label className="block text-sm font-medium text-neutral-700 mb-2">Username</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="pos-input"
              placeholder="Masukkan username"
            />
          </div>

          <div className="mb-6">
            <label className="block text-sm font-medium text-neutral-700 mb-2">PIN Kasir</label>
            <div className="flex justify-center gap-3">
              {pin.map((digit, i) => (
                <input
                  key={i}
                  id={`pin-${i}`}
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
            onClick={handleLogin}
            disabled={loading}
            className="w-full py-3 bg-brand hover:bg-brand/90 text-white font-medium rounded-xl transition-colors disabled:opacity-50 focus:ring-2 focus:ring-brand/50 focus:ring-offset-2"
          >
            {loading ? 'Memproses...' : 'Masuk'}
          </button>
        </div>

        <p className="text-center text-neutral-500 text-xs mt-6">
          POS System v1.0 &mdash; High Performance Point of Sale
        </p>
      </div>
    </div>
  );
}