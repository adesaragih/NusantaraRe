import { create } from 'zustand'

// [usulan] Zustand belum pernah diputuskan lewat ADR; ia dipakai sebagai
// usulan dan boleh diganti lewat keputusan tertulis.
//
// Fase 0 - scaffold: nol keadaan dagang. Yang ada hanya keadaan sesi paling
// dasar, supaya tiket pertama yang punya layar menemukan bentuknya.

export const useAppStore = create((set) => ({
  // Keadaan backend sebagaimana dilaporkan GET /healthz.
  kesehatan: null,
  setKesehatan: (kesehatan) => set({ kesehatan }),

  // Galat terakhir yang perlu ditampilkan ke pengguna.
  galat: null,
  setGalat: (galat) => set({ galat }),
  bersihkanGalat: () => set({ galat: null }),
}))
