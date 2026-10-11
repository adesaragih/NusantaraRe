// Pendaftaran layar modul Komite Claim Prop - modul TANPA MENU (perintah work owner 09-10-2026: menu komite dihapus,
// "anggap menu itu tidak pernah ada" - komite digabung ke menu Claim Prop). `PENDAFTARAN_LAYAR` dibaca perakit
// `frontend/daftar.ts` lewat `import.meta.glob`; modul ini dipasang bagi pemegang menu Claim Prop (`MODUL_DIPINJAM`
// frontend/App.tsx) dan kasusnya dibuka DI TEMPAT dari tabel komite inbox Claim Prop (`PropsRute.onBukaModul`).

import type { LayarModul } from '../../../inti/frontend/modul'
import { KORPUS_KCP } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_KCP = 'komiteclaimprop'
export const HALAMAN_KCP = ['komiteclaimprop-kasus'] as const
export type HalamanKCP = (typeof HALAMAN_KCP)[number]

export const PENDAFTARAN_LAYAR: LayarModul<HalamanKCP> = {
  nama: NAMA_KCP,
  kelompok: KORPUS_KCP.kelompok,
  halaman: HALAMAN_KCP,
  halamanAwal: 'komiteclaimprop-kasus',
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    komiteclaimprop: HalamanKCP
  }
}
