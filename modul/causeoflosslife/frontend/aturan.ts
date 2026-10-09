// Aturan layar Cause Of Loss Life - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/models`,
// `backend/services`): wajib, batas kolom, nama kembar, hak View only, ID dari sequence.

import { COL } from './labels'

/** `pyPageSize` 10 (`InboxCauseofLossLife` b4057) - backend `models.UkuranHalaman`. */
export const UKURAN_HALAMAN = 10

/** Lebar `CAUSEOFLOSS_LIFE.CAUSEOFLOSS` VARCHAR2(200) (byte) - backend `models.BatasCauseOfLoss` (di luar XML: batas kolom). */
export const BATAS_CAUSE_OF_LOSS = 200

/** Jumlah halaman (minimal 1). */
export function jumlahHalaman(total: number, ukuran = UKURAN_HALAMAN): number {
  return Math.max(1, Math.ceil(total / ukuran))
}

/** Pemeriksaan awal form: Cause of Loss wajib (`pyRequired` b996 / b1044); panjang paling banyak 200 byte (batas kolom -
 * di luar XML). Huruf TIDAK diubah (XML tanpa pengubah huruf). `null` = boleh dikirim; backend memeriksa ulang. */
export function periksaCauseOfLoss(nilai: string): string | null {
  const t = nilai.trim()
  if (t === '') return COL.galatWajib
  if (new TextEncoder().encode(t).length > BATAS_CAUSE_OF_LOSS) return COL.galatPanjang(BATAS_CAUSE_OF_LOSS)
  return null
}
