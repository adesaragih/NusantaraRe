// Pendaftaran layar modul Komite Claim Life - modul TANPA MENU (perintah work owner 09-10-2026: menu komite dihapus,
// "anggap menu itu tidak pernah ada" - komite digabung ke menu Claim Life). `PENDAFTARAN_LAYAR` dibaca perakit
// `frontend/daftar.ts` lewat `import.meta.glob`; modul ini dipasang bagi pemegang menu Claim Life (`MODUL_DIPINJAM`
// frontend/App.tsx) dan kasusnya dibuka DI TEMPAT dari tabel komite inbox Claim Life (`PropsRute.onBukaModul`).

import type { LayarModul } from '../../../inti/frontend/modul'

/** Nama modul - SAMA dengan `const Nama` di `modul/komiteclaimlife/backend/modul.go`. */
export const NAMA_KOMITE = 'komiteclaimlife'

/** Nama modul - nama folder korpus VERBATIM (`D:\XML\RNM_BRD\Komite Claim Life`). Modul ini TANPA menu. */
export const KELOMPOK_KOMITE = 'Komite Claim Life'

export const HALAMAN_KOMITE = ['komiteclaimlife-kasus'] as const
export type HalamanKomite = (typeof HALAMAN_KOMITE)[number]

export const PENDAFTARAN_LAYAR: LayarModul<HalamanKomite> = {
  nama: NAMA_KOMITE,
  kelompok: KELOMPOK_KOMITE,
  halaman: HALAMAN_KOMITE,
  halamanAwal: 'komiteclaimlife-kasus',
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    komiteclaimlife: HalamanKomite
  }
}
