// Pembantu tampilan MURNI modul Endorsement Life.
//
// ⛔ Angka tidak pernah menjadi `Number`: sel angka lewat `formatNumber` (digit, bukan float).

import { DESIMAL_TAK_DIBATASI, formatDate, formatNumber } from '../../../inti/frontend/lib/format'
import { OPSI_EDM_TYPE, OPSI_TYPE_CEDING } from './labels'

/** Baris per halaman grid - SAMA dengan `services.UkuranHalaman` (`pyGridPaginator`). */
export const UKURAN_HALAMAN_EDM = 20

/** Sel kosong ditandai, bukan dibiarkan kosong (ADR-U-0027). */
export function sel(v: string | undefined): string {
  const t = (v ?? '').trim()
  return t === '' ? '—' : t
}

/** Sel tanggal `YYYY-MM-DD` → `DD-MM-YYYY`. */
export function selTanggal(v: string | undefined): string {
  return sel(formatDate(v ?? ''))
}

/** Sel waktu `YYYY-MM-DD HH:MM:SS` → `DD-MM-YYYY HH:MM:SS`; bentuk lain apa adanya. */
export function selWaktu(v: string | undefined): string {
  const m = /^(\d{4})-(\d{2})-(\d{2})[ T](\d{2}:\d{2}:\d{2})/.exec((v ?? '').trim())
  return m ? `${m[3]}-${m[2]}-${m[1]} ${m[4]}` : sel(v)
}

/** Sel angka - digit apa adanya, pemisah ribuan, nol ekor dipangkas; negatif tetap bertanda. */
export function selAngka(v: string | undefined): string {
  const t = (v ?? '').trim()
  return t === '' ? '—' : formatNumber(t, DESIMAL_TAK_DIBATASI)
}

/** Label `EDM Type` (`1` Perubahan Data, `3` Batal); nilai lain apa adanya. */
export function labelEdmType(v: string): string {
  return OPSI_EDM_TYPE[v] ?? sel(v)
}

/** Label `Reinsurance System` dari kode `TYPE_CEDING`; kode tak dikenal apa adanya. */
export function labelTypeCeding(v: string): string {
  return OPSI_TYPE_CEDING[v] ?? sel(v)
}
