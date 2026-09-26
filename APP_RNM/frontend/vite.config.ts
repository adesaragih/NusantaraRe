import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

// Pengaturan Vite (alat yang menjalankan dan membangun frontend).
//
// Nol alamat host di berkas ini (ADR-U-0004). Dua alamat datang dari env:
//   - VITE_API_BASE_URL  : alamat backend yang dipanggil browser; kosong berarti
//                          sama-asal. Dibaca kode lewat import.meta.env.
//   - DEV_PROXY_TARGET   : tujuan proxy SAAT DEV saja. Permintaan /api/* dan
//                          /healthz yang tiba di server dev Vite diteruskan ke
//                          alamat ini, sehingga browser tetap sama-asal dan
//                          backend tidak perlu header CORS. Tanpa awalan VITE_
//                          supaya tidak pernah ikut ke bundel browser.
// Keduanya ditulis di frontend/.env (contohnya di .env.example).
//
// `loadEnv(mode, '.', '')` membaca SELURUH variabel .env di folder frontend,
// termasuk yang tanpa awalan VITE_; awalan kosong itulah yang membuat
// DEV_PROXY_TARGET terbaca di sini. Env var yang sudah ada di proses (mis.
// disetel di terminal) menang atas isi .env - itu perilaku bawaan loadEnv.
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, '.', '')
  const proxyTarget = env.DEV_PROXY_TARGET

  return {
    plugins: [react()],
    server: {
      port: 5173,
      // Bila DEV_PROXY_TARGET kosong, tidak ada proxy: perilaku lama.
      proxy: proxyTarget
        ? {
            '/api': { target: proxyTarget, changeOrigin: true },
            '/healthz': { target: proxyTarget, changeOrigin: true },
          }
        : undefined,
    },
    build: {
      outDir: 'dist',
    },
  }
})
