// Aturan layar Cover Life - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/models`,
// `backend/services`): wajib, batas kolom, Cover kembar, hak View only, ID dari sequence.

import { CVL } from './labels'

/** `pyPageSize` 50 (`InboxCoverLife` b4101) - backend `models.UkuranHalaman`. */
export const UKURAN_HALAMAN = 50

/** Lebar `M_COVER_LIFE.COVER` VARCHAR2(200) (byte) - backend `models.BatasCover` (di luar XML: batas kolom). */
export const BATAS_COVER = 200

/** Lebar `M_COVER_LIFE.NOTE` VARCHAR2(1000) (byte) - backend `models.BatasNote` (di luar XML: batas kolom). */
export const BATAS_NOTE = 1000

/** Jumlah halaman (minimal 1). */
export function jumlahHalaman(total: number, ukuran = UKURAN_HALAMAN): number {
  return Math.max(1, Math.ceil(total / ukuran))
}

function bytes(s: string): number {
  return new TextEncoder().encode(s).length
}

/** Pemeriksaan awal form: Cover wajib (`pyRequired` b843 / b895); Note opsional (b1020); panjang paling banyak lebar
 * kolom (di luar XML). Huruf TIDAK diubah (XML tanpa pengubah huruf). `null` = boleh dikirim; backend memeriksa ulang. */
export function periksaIsian(cover: string, note: string): string | null {
  const t = cover.trim()
  if (t === '') return CVL.galatWajib
  if (bytes(t) > BATAS_COVER) return CVL.galatPanjang(CVL.cover, BATAS_COVER)
  if (note.trim() !== '' && bytes(note) > BATAS_NOTE) return CVL.galatPanjang(CVL.note, BATAS_NOTE)
  return null
}
