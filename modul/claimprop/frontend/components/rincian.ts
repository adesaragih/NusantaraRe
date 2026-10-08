// Panel rinci baris grid (expand pane). Di XML hanya grid Adjustment (`ClaimData.AdjustmentList`) yang punya expand
// pane - Section `AdjustmentDetail`. Grid lain (Insured Interests, Total in Original Currency, ...) tidak bisa dibuka;
// dulu panel itu menempel di semua grid sehingga klik baris mana pun membuka "AdjustmentDetail" kosong (laporan work
// owner 08-10-2026).

import type { ReactNode } from 'react'

/** Satu-satunya grid ber-expand pane di XML. */
export const DAFTAR_RINCI = 'ClaimData.AdjustmentList'

export interface RincianGrid {
  /** Jalur daftar grid yang barisnya dapat dibuka. */
  daftar: string
  /** Isi panel baris ke-n (1..n). */
  isi: (n: number) => ReactNode
}

/** Panel rinci untuk grid berjalur `jalurGrid`; `undefined` = baris grid itu tidak dapat dibuka. */
export function rincianUntuk(r: RincianGrid | undefined, jalurGrid: string | undefined): RincianGrid | undefined {
  return r !== undefined && r.daftar === jalurGrid ? r : undefined
}

/**
 * Baris grid ber-panel yang terbuka tanpa klik (work owner 08-10-2026 "tidak harus klik angka sebelah kiri untuk
 * membuka adjustment nya"): baris terbaru terbuka sejak awal.
 */
export function bukaAwal(nBaris: number): number | null {
  return nBaris > 0 ? nBaris : null
}

/** Baris terbuka sesudah jumlah baris berubah: baris baru (Add) langsung terbuka; baris yang terhapus ditutup. */
export function barisTerbuka(buka: number | null, nLalu: number, nKini: number): number | null {
  if (nKini > nLalu) return nKini
  if (buka !== null && buka > nKini) return bukaAwal(nKini)
  return buka
}
