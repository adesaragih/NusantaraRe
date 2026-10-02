import type { MenuModul } from '../../../inti/frontend/modul'

export const NAMA_TREATYINADJUSTMENT = 'treatyinadjustment'
/** Nama modul - nama folder korpus VERBATIM (= M_NAV_MENU.LABEL barisnya). */
export const KELOMPOK_TREATYINADJUSTMENT = 'Treaty In Adjustment'
export const HALAMAN_TREATYINADJUSTMENT = ['treatyinadjustment-rantai-versi'] as const
export type HalamanTreatyInAdjustment = (typeof HALAMAN_TREATYINADJUSTMENT)[number]

/**
 * Halaman yang dibuka tombol "Treaty In Adjustment" di sidebar.
 *
 * ⛔ SATU halaman, baca-saja. Papan tiket modul ini menyatakan `L-4`: "tidak ada
 * spesifikasi layar di mana pun". Jalur SIMPAN - yang menegakkan materialitas
 * (tiket 02) dan batas tanggal berlaku (tiket 03) - lahir bersama
 * spesifikasinya.
 */
export const HALAMAN_AWAL_TREATYINADJUSTMENT: HalamanTreatyInAdjustment = 'treatyinadjustment-rantai-versi'

export const PENDAFTARAN_MENU: MenuModul<HalamanTreatyInAdjustment> = {
  nama: NAMA_TREATYINADJUSTMENT,
  kelompok: KELOMPOK_TREATYINADJUSTMENT,
  halaman: HALAMAN_TREATYINADJUSTMENT,
  halamanAwal: HALAMAN_AWAL_TREATYINADJUSTMENT,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    treatyinadjustment: HalamanTreatyInAdjustment
  }
}
