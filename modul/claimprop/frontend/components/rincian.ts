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
