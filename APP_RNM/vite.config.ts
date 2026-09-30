/// <reference types="vitest/config" />
import { fileURLToPath } from 'node:url'
import { join } from 'node:path'

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
// Keduanya ditulis di frontend/.env (contohnya di frontend/.env.example).
//
// # Struktur tim satu folder per modul (30-09-2026)
//
// Berkas ini, package.json, dan tsconfig.json tinggal di APP_RNM/, bukan di
// frontend/: kode frontend kini tersebar di tiga tempat - `frontend/`
// (perakit), `inti/frontend/` (kerangka bersama), dan `modul/<nama>/frontend/`
// (satu folder per modul). Berkas di luar `frontend/` mencari `react` ke atas
// dari foldernya sendiri, jadi `node_modules` harus berada di induk ketiganya.
//
// - `root` tetap `frontend/`: index.html, public/, dist/, dan .env di sana.
// - `loadEnv` membaca folder yang SAMA dengan `root` - bukan folder kerja
//   proses. `npm run dev` dijalankan dari APP_RNM/, dan APP_RNM/.env adalah
//   env BACKEND (ORACLE_DSN dan kawan-kawan) yang tidak boleh terbaca di sini.
//   Awalan kosong membuat DEV_PROXY_TARGET (tanpa awalan VITE_) ikut terbaca.
//   Env var yang sudah ada di proses (mis. disetel di terminal) menang atas
//   isi .env - itu perilaku bawaan loadEnv.
// - `server.fs.allow` menyebut ketiga tempat kode itu: server dev menolak
//   menyajikan berkas di luar daftar ini.
// - `test.dir` + `test.include`: Vitest mencari berkas uji di ketiganya juga;
//   bawaannya hanya di bawah `root`.

const AKAR = fileURLToPath(new URL('.', import.meta.url))
const AKAR_FRONTEND = join(AKAR, 'frontend')

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, AKAR_FRONTEND, '')
  const proxyTarget = env.DEV_PROXY_TARGET

  return {
    root: AKAR_FRONTEND,
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
      fs: {
        allow: [AKAR_FRONTEND, join(AKAR, 'inti', 'frontend'), join(AKAR, 'modul'), join(AKAR, 'node_modules')],
      },
    },
    build: {
      // Relatif `root`: frontend/dist, seperti sebelumnya.
      outDir: 'dist',
    },
    test: {
      dir: AKAR,
      // Pola bawaan Vitest, di ketiga tempat - bukan hanya `*.test.ts`.
      include: ['frontend', 'inti/frontend', 'modul/*/frontend'].map((d) => `${d}/**/*.{test,spec}.?(c|m)[jt]s?(x)`),
    },
  }
})
