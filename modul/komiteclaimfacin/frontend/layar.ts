// Pendaftaran layar modul Komite Claim Fac In - modul TANPA MENU (perintah work owner 09-10-2026: menu komite dihapus,
// "anggap menu itu tidak pernah ada" - komite digabung ke menu Claim Fac In; prompt tahap 2 §2 butir 2). Pola
// `modul/komiteclaimnonprop/frontend/layar.ts`. `PENDAFTARAN_LAYAR` dibaca perakit `frontend/daftar.ts` lewat
// `import.meta.glob`; modul ini dipasang bagi pemegang menu Claim Fac In (`MODUL_DIPINJAM` frontend/App.tsx) dan
// kasusnya dibuka DI TEMPAT dari tabel komite inbox Claim Fac In (`PropsRute.onBukaModul`).

import type { LayarModul } from '../../../inti/frontend/modul'
import { KORPUS_KCFI } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_KCFI = 'komiteclaimfacin'
export const HALAMAN_KCFI = ['komiteclaimfacin-kasus'] as const
export type HalamanKCFI = (typeof HALAMAN_KCFI)[number]

export const PENDAFTARAN_LAYAR: LayarModul<HalamanKCFI> = {
  nama: NAMA_KCFI,
  kelompok: KORPUS_KCFI.kelompok,
  halaman: HALAMAN_KCFI,
  halamanAwal: 'komiteclaimfacin-kasus',
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    komiteclaimfacin: HalamanKCFI
  }
}
