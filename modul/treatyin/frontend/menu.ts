import type { MenuModul } from '../../../inti/frontend/modul'

export const NAMA_TREATYIN = 'treatyin'
/** Nama modul - nama folder korpus VERBATIM (= M_NAV_MENU.LABEL barisnya). */
export const KELOMPOK_TREATYIN = 'Treaty In'
export const HALAMAN_TREATYIN = ['treatyin-kontrak', 'treatyin-acuan'] as const
export type HalamanTreatyIn = (typeof HALAMAN_TREATYIN)[number]

/**
 * Halaman yang dibuka tombol "Treaty In" di sidebar.
 *
 * ⚠️ DUA halaman sejak 3 Oktober 2026, dan sebabnya layak dibaca. Berkas ini
 * dulu menyatakan SATU halaman dengan alasan `L-4` — *"tidak ada spesifikasi
 * layar di mana pun"*. **`L-4` DICABUT**: spesifikasinya ada dan selalu ada,
 * 50 berkas `Section/` dan 3 `Harness/` di ekspor Pega 2026-09 untuk modul
 * ini (68 + 6 untuk Adjustment). Yang tidak ada bukan spesifikasinya,
 * melainkan pembacaannya.
 *
 * Larangan "jangan mengarang layar" TETAP berlaku — dan kini tidak ada alasan
 * untuk mengarang: tiap medan di kedua layar membawa nama rule dan posisi
 * bitanya di `labels.ts`.
 */
export const HALAMAN_AWAL_TREATYIN: HalamanTreatyIn = 'treatyin-kontrak'

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
