// GOLONGAN nilai per properti — ditulis tangan, dibaca bersama kerangka.
//
// ⛔ Golongan menentukan DUA hal saja: apakah nilainya diformat sebagai angka
// (titik ribuan, koma desimal), dan desimal CADANGAN bila ekspor tidak
// menyatakan `pyDecimalPlaces` sel itu. Desimal yang DINYATAKAN ekspor selalu
// menang (keputusan pemilik proses 5 Oktober 2026: desimal PER KOLOM).
//
// Kunci berindeks (`RnmLimitListDisplay(1).Value`) dan bertitik
// (`ValueDifference.RNMShare`) digolongkan menurut SEGMEN TERAKHIRNYA.

import type { JenisAngka } from '../labelsPenyesuaian'

const TANGGAL = [
  'AsDate', 'InitialDate', 'SubmissionDue', 'ConfirmationDue', 'SettlementDue', 'ReportDate', 'SubDueDate',
  'PeriodStart', 'PeriodEnd', 'ReportingStart', 'ReportingEnd', 'EDMEffective', 'Commencement', 'Termination',
]
const UANG = [
  'Amount', 'AmountIDR', 'Limit', 'Limit2', 'Deductible', 'Deductible2', 'MDP', 'MDP2', 'AggregateLimit',
  'AggregateLimit2', 'NetPremi', 'NetPremi2', 'Value', 'Conversion', 'RSMDLimit', 'Earthquake', 'FloodJab',
  'FloodNation', 'TotalEgnpiAmount',
]
const PERSEN = [
  'Proportion', 'SharePct', 'RNMShare', 'RNMShareP', 'BrokeragePercent', 'BrokeragePercentP', 'FacultativeShare',
  'FacultativeShareBrokerage', 'FacShare', 'FacShareBrokerage', 'RnmShareDeducted', 'ProRatePercent',
  'TotalEgnpiProportion', 'TotalLimitsROL',
]
/**
 * ⛔ TIDAK diformat walau selnya `pxNumber`: `Currency` (kode mata uang),
 * `Note` (teks layer), dan dua medan yang golongannya TIDAK diketahui —
 * `MaxCoGroup`/`MaxCoNonGroup` (terukur hanya bernilai `10`). Didaftarkan di
 * `docs/PERTANYAAN-TERBUKA-LAYAR-ADJUSTMENT.md`.
 */
const TEKS = ['Currency', 'Note', 'MaxCoGroup', 'MaxCoNonGroup']

const PETA: Readonly<Record<string, JenisAngka>> = Object.fromEntries([
  ...TANGGAL.map((k) => [k, 'tanggal'] as const),
  ...UANG.map((k) => [k, 'uang'] as const),
  ...PERSEN.map((k) => [k, 'persen'] as const),
  ...TEKS.map((k) => [k, 'teks'] as const),
])

/** Segmen terakhir kunci: `RnmLimitListDisplay(1).Value` → `Value`. */
export function segmenAkhir(kunci: string): string {
  const s = kunci.split('.')
  return s[s.length - 1] ?? kunci
}

/** Golongan sebuah kunci; `pxDateTime` tanpa entri peta → tanggal. */
export function golongan(kunci: string, format = ''): JenisAngka {
  const g = PETA[segmenAkhir(kunci)]
  if (g !== undefined) return g
  if (format === 'pxDateTime') return 'tanggal'
  return 'teks'
}

/** Kunci yang masuk peta — dipakai uji untuk mendaftar yang belum. */
export const KUNCI_BERGOLONGAN: ReadonlySet<string> = new Set(Object.keys(PETA))
