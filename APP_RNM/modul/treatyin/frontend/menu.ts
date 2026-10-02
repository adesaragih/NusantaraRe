import type { MenuModul } from '../../../inti/frontend/modul'

export const NAMA_TREATYIN = 'treatyin'
/** Nama modul - nama folder korpus VERBATIM (= M_NAV_MENU.LABEL barisnya). */
export const KELOMPOK_TREATYIN = 'Treaty In'
export const HALAMAN_TREATYIN = ['treatyin-acuan'] as const
export type HalamanTreatyIn = (typeof HALAMAN_TREATYIN)[number]

/**
 * Halaman yang dibuka tombol "Treaty In" di sidebar.
 *
 * ⛔ SATU halaman, dan itu disengaja. Papan tiket modul ini menyatakan `L-4`:
 * "tidak ada spesifikasi layar di mana pun". Layar kontrak, versi, dan
 * persetujuan lahir bersama spesifikasinya - halaman ini hanya memperlihatkan
 * keenam tabel acuan tiket 15.
 */
export const HALAMAN_AWAL_TREATYIN: HalamanTreatyIn = 'treatyin-acuan'

export const PENDAFTARAN_MENU: MenuModul<HalamanTreatyIn> = {
  nama: NAMA_TREATYIN,
  kelompok: KELOMPOK_TREATYIN,
  halaman: HALAMAN_TREATYIN,
  halamanAwal: HALAMAN_AWAL_TREATYIN,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    treatyin: HalamanTreatyIn
  }
}
