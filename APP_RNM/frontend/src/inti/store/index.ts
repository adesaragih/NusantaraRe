import { create } from 'zustand'

import type { Kesehatan } from '../klien'

// [usulan] Zustand = penyimpan keadaan BERSAMA antar-komponen: nilai yang
// perlu dibaca banyak halaman sekaligus tanpa dioper-oper lewat props.
// Belum pernah diputuskan lewat ADR; boleh diganti lewat keputusan tertulis.
//
// Fase 0: nol keadaan dagang. Yang ada hanya keadaan sesi paling dasar.
// Keadaan yang hanya dipakai satu halaman (mis. teks yang sedang diketik)
// cukup memakai useState di halaman itu, tidak perlu masuk ke sini.

/** Bentuk seluruh keadaan bersama. Tambahkan medan di sini bila tiket memerlukannya. */
export interface AppState {
  /** Keadaan backend sebagaimana dilaporkan GET /healthz; null = belum dicek. */
  kesehatan: Kesehatan | null
  setKesehatan: (kesehatan: Kesehatan | null) => void

  /** Galat terakhir yang perlu ditampilkan ke pengguna. */
  galat: string | null
  setGalat: (galat: string | null) => void
  bersihkanGalat: () => void
}

// Cara pakai di komponen:   const galat = useAppStore((s) => s.galat)
export const useAppStore = create<AppState>()((set) => ({
  kesehatan: null,
  setKesehatan: (kesehatan) => set({ kesehatan }),

  galat: null,
  setGalat: (galat) => set({ galat }),
  bersihkanGalat: () => set({ galat: null }),
}))
