// Panel rinci baris grid (expand pane). Di XML Claim Non Prop tiga grid ber-`pyRowEditing=masterDetail` dengan
// `pyEditingMode=expandPane`: Acceptation List (`ClaimData.AdjustmentList` -> flow action AdjustmentDetailNP), Insured
// Interests Outstanding (`ClaimData.InterestList` -> InputDtlInterest), dan XOL Allocation panel akseptasi
// (`ClaimData.AdjustmentList(n).SpreadingRisk` -> ShowDetailXOL). Grid Limit Layer (ViewClaimLayerDetail) ber-`tampil[1=2]`
// - tidak pernah tampil. Grid lain tidak dapat dibuka. Pola panel disalin dari `modul/claimprop/frontend/components/
// rincian.ts` (bukan impor).

import type { ReactNode } from 'react'

/** Grid Acceptation List (expand pane AdjustmentDetailNP). */
export const DAFTAR_ADJ = 'ClaimData.AdjustmentList'
/** Grid Insured Interests 100 % Outstanding (expand pane InputDtlInterest). */
export const DAFTAR_INTEREST = 'ClaimData.InterestList'
/** Anak XOL Allocation panel akseptasi (expand pane ShowDetailXOL). */
export const ANAK_XOL = 'SpreadingRisk'

export interface RincianGrid {
  /** Jalur daftar grid yang barisnya dapat dibuka. */
  daftar: string
  /** Isi panel baris ke-n (1..n). */
  isi: (n: number) => ReactNode
  /** Baris terbaru terbuka tanpa klik (Acceptation List - pola Claim Prop, work owner 08-10-2026). */
  bukaAwal?: boolean
  /** Panel membawa nomor akseptasi (`BarisAdjustment`) - hanya panel AdjustmentDetailNP. */
  nomorAkseptasi?: boolean
}

/** Jalur grid XOL Allocation panel akseptasi ke-n. */
export function jalurXOL(n: number): string {
  return `${DAFTAR_ADJ}(${n}).${ANAK_XOL}`
}

/** Nomor akseptasi dan baris XOL dari jalur grid XOL Allocation panel; null bila bukan jalur itu. */
export function pecahJalurXOL(jalur: string): number | null {
  const m = /^ClaimData\.AdjustmentList\((\d+)\)\.SpreadingRisk$/.exec(jalur)
  return m ? Number(m[1]) : null
}

/**
 * Nomor baris adjustment dari kunci modal (`komite:n` = harness KomiteCNP); 0 bila modal halaman (`pla`, `tutupKlaim`,
 * `cwp`). Tombol modal baris (Send Claim to Committee) adalah aksi panel akseptasi.
 */
export function barisModal(kunci: string): number {
  const m = /^komite:(\d+)$/.exec(kunci)
  return m ? Number(m[1]) : 0
}

/** Panel rinci untuk grid berjalur `jalurGrid`; `undefined` = baris grid itu tidak dapat dibuka. */
export function rincianUntuk(
  r: readonly RincianGrid[] | undefined,
  jalurGrid: string | undefined,
): RincianGrid | undefined {
  return r?.find((x) => x.daftar === jalurGrid)
}

/** Baris yang terbuka sejak awal: baris terbaru bila panelnya `bukaAwal`. */
export function bukaAwal(nBaris: number, awal = true): number | null {
  return awal && nBaris > 0 ? nBaris : null
}

/** Baris terbuka sesudah jumlah baris berubah: baris baru (Add) langsung terbuka; baris yang terhapus ditutup. */
export function barisTerbuka(buka: number | null, nLalu: number, nKini: number, awal = true): number | null {
  if (nKini > nLalu) return nKini
  if (buka !== null && buka > nKini) return bukaAwal(nKini, awal)
  return buka
}
