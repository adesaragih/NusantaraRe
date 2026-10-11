import type { MenuModul } from '../../../inti/frontend/modul'

export const NAMA_TREATYINADJUSTMENT = 'treatyinadjustment'
/** Nama modul - nama folder korpus VERBATIM (= M_NAV_MENU.LABEL barisnya). */
export const KELOMPOK_TREATYINADJUSTMENT = 'Treaty In Adjustment'
export const HALAMAN_TREATYINADJUSTMENT = [
  // ⭐ Layar Adjustment — `Section/InputTreatyInAdjustment.xml`, 5 Oktober
  // 2026: daftar penyesuaian warisan, lalu Old Data ‖ New Data, Attachment,
  // deret tombol, History. Data dari `TREATY_IN_EDM` + `M_TREATY_IN_EDM`.
  'treatyinadjustment-penyesuaian',
  'treatyinadjustment-rantai-versi',
  // ⭐ Halaman KEDUA, 4 Oktober 2026 — panel Attachment + History.
  //
  // ⛔ Ia TIDAK melanggar `L-4`. `L-4` melarang MENGARANG layar yang tidak
  // punya spesifikasi; layar ini disalin dari ekspor Pega modul ini sendiri
  // — `Section/WorkAttachments.xml` dan `ShowAttachmentTreaty.xml`, lengkap
  // dengan keempat kolomnya, kedua tombolnya, dan spanduk birunya.
  'treatyinadjustment-lampiran',
] as const
export type HalamanTreatyInAdjustment = (typeof HALAMAN_TREATYINADJUSTMENT)[number]

/**
 * Halaman yang dibuka tombol "Treaty In Adjustment" di sidebar.
 *
 * ⭐ Layar Adjustment sejak 5 Oktober 2026 — layar sistem lama yang
 * sesungguhnya. Rantai Versi membaca `KONTRAK`/`VERSI_KONTRAK` model baru,
 * yang terukur NOL baris; membukanya lebih dulu memperlihatkan layar kosong
 * padahal 280 penyesuaian warisan ada.
 *
 * ⛔ Seluruh halaman BACA-SAJA. Jalur SIMPAN - yang menegakkan materialitas
 * (tiket 02) dan batas tanggal berlaku (tiket 03) - belum dibangun; tombol
 * tulisnya dimatikan, bukan dihilangkan.
 */
export const HALAMAN_AWAL_TREATYINADJUSTMENT: HalamanTreatyInAdjustment = 'treatyinadjustment-penyesuaian'

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
