// Pendaftaran layar modul Komite Claim Non Prop - modul TANPA MENU (perintah work owner 09-10-2026: menu komite dihapus,
// "anggap menu itu tidak pernah ada" - komite digabung ke menu Claim Non Prop). `PENDAFTARAN_LAYAR` dibaca perakit
// `frontend/daftar.ts` lewat `import.meta.glob`; modul ini dipasang bagi pemegang menu Claim Non Prop (`MODUL_DIPINJAM`
// frontend/App.tsx) dan kasusnya dibuka DI TEMPAT dari tabel komite inbox Claim Non Prop (`PropsRute.onBukaModul`).

import type { LayarModul } from '../../../inti/frontend/modul'
import { KORPUS_KCNP } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_KCNP = 'komiteclaimnonprop'
export const HALAMAN_KCNP = ['komiteclaimnonprop-kasus'] as const
export type HalamanKCNP = (typeof HALAMAN_KCNP)[number]

export const PENDAFTARAN_LAYAR: LayarModul<HalamanKCNP> = {
  nama: NAMA_KCNP,
  kelompok: KORPUS_KCNP.kelompok,
  halaman: HALAMAN_KCNP,
  halamanAwal: 'komiteclaimnonprop-kasus',
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    komiteclaimnonprop: HalamanKCNP
  }
}
