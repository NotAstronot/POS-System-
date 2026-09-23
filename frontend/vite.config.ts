import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': `${import.meta.dirname}/src`,
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5177,
    strictPort: true,
    proxy: {
      '/api': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
      '/uploads': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
      '/exports': {
        target: 'http://localhost:8082',
        changeOrigin: true,
      },
    },
  },
  base: './',
});

