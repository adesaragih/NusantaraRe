import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Pengaturan Vite (alat yang menjalankan dan membangun frontend).
// Nol alamat host di berkas ini. Alamat backend datang dari env var
// VITE_API_BASE_URL (ADR-U-0004); lihat .env.example dan src/vite-env.d.ts.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
  },
  build: {
    outDir: 'dist',
  },
})
