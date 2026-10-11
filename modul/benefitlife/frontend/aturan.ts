// Aturan layar Benefit - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/models`,
// `backend/services`): wajib, huruf besar, batas kolom, hak View only, ID dari sequence.

import type { Saringan } from './api'
import { BN } from './labels'

/** `pyPageSize` 10 (`InboxBenefit` b4526) - backend `models.UkuranHalaman`. */
export const UKURAN_HALAMAN = 10

/** Lebar `BENEFIT_LIFE.BENEFIT` VARCHAR2(200) (byte) - backend `models.BatasBenefit` (di luar XML: batas kolom). */
export const BATAS_BENEFIT = 200

/** Saringan awal: tanpa filter, ID MENURUN (`pySortType` DESC b4391, `pySortOrder` 1 b4397), halaman 1. */
export const SARINGAN_AWAL: Saringan = { id: '', benefit: '', naik: false, halaman: 1 }

/** Klik kepala kolom ID: balik arah, halaman kembali 1. Kolom Benefit tidak dapat diurutkan (`pyColumnSorting` false b4415). */
export function balikArah(s: Saringan): Saringan {
  return { ...s, naik: !s.naik, halaman: 1 }
}

/** Penanda arah di kepala kolom ID. */
export function tandaArah(s: Saringan): string {
  return s.naik ? ' ▲' : ' ▼'
}

/** Jumlah halaman (minimal 1). */
export function jumlahHalaman(total: number, ukuran = UKURAN_HALAMAN): number {
  return Math.max(1, Math.ceil(total / ukuran))
}

/** `SetUpperCase_DT` (b1211): perubahan textarea Benefit menjadikannya huruf besar - dipanggil saat medan ditinggalkan
 * (peristiwa `change` Pega b1199). Backend menerapkan hal yang sama. */
export function hurufBesar(s: string): string {
  return s.toUpperCase()
}

/** Pemeriksaan awal form: Benefit wajib (`pyRequired` b1137); panjang paling banyak 200 byte (batas kolom - di luar
 * XML). `null` = boleh dikirim; backend memeriksa ulang. */
export function periksaBenefit(benefit: string): string | null {
  const t = benefit.trim()
  if (t === '') return BN.galatBenefit
  if (new TextEncoder().encode(t.toUpperCase()).length > BATAS_BENEFIT) return BN.galatPanjang(BATAS_BENEFIT)
  return null
}
