// Subsection `DetailPoliciesNonProportional` -> `DetailPolicyTreatyInNonProportional`
// (polis NB NonProporsional / XOL, `[keputusan work owner]` K8): label VERBATIM sel
// section, syarat tampil, dan aturan sunting grid. Isinya halaman master `TreatyIn`
// (dibaca backend, baca-saja) dan daftar polis `SpreadingRiskList` /
// `ListInstallment().InstallmentList` (dihitung backend saat pilih bisnis).
//
// ⛔ Tidak ditampilkan, dan sebabnya (bukti XML):
//   - `DetailPolicyTreatyInNonProportionalEDM`: tampil bila `TreatyMasterInEDM &&
//     IsEDMInputOnNB != true`; InputPolicyTreatyInDetail_NonProp langkah 10 mengisi
//     `IsEDMInputOnNB = true` tepat ketika `TreatyMasterInEDM` - tidak terjangkau;
//   - tombol `.pyTemplateButton` (refresh `TreatyInNonSetTotal`): tampil `1=2`;
//   - (label "NON EDM" / "EDM" DIGERBANG tempat berperan tiket 05 - `tempat.ts`,
//     `components/DetailNonProp.tsx`; tertunda selama pemetaan IAM kosong.)

import { KLAIM_XOL_RETRO, MASTER, POLIS, nilai, type Baris, type Halaman } from './api'
import type { Sajian } from './sajian'
import { nolTeks, positifTeks } from './tanda'

/** `.IsNewPolicyNonProp = 1` - kontainer proporsional tersembunyi, subsection NonProp tampil. */
export const polisNonPropBaru = (h: Halaman) => nilai(h, POLIS + 'IsNewPolicyNonProp') === '1'

/** Kontainer `DetailPoliciesNonProportional`: `.IsNewPolicyNonProp = 1 && .ClaimType != 'XOL Retro'`. */
export const tampilNonProp = (h: Halaman) => polisNonPropBaru(h) && nilai(h, POLIS + 'ClaimType') !== KLAIM_XOL_RETRO

/** `pyWorkPage.TreatyIn.FacultativeShare` = 0 (kosong dihitung 0; tanda teks, nol `Number`). */
export const fakultatifNol = (h: Halaman) => nolTeks(nilai(h, MASTER + 'FacultativeShare'))

/** Kontainer "Share Facultative": `pyWorkPage.TreatyIn.FacultativeShare>0`. */
export const tampilFakultatif = (h: Halaman) => positifTeks(nilai(h, MASTER + 'FacultativeShare'))

/**
 * Section `SpreadingRiskList`: tombol Add/Delete tampil bila FacultativeShare = 0 atau '';
 * TreatyType / %Share / %Share Claim terkunci bila FacultativeShare > 0.
 */
export const spreadingTerbuka = (h: Halaman) => fakultatifNol(h)

/** Satu kolom grid: anggota baris, label sel VERBATIM, dan format sel bila BUKAN
 *  pxNumber `pyDecimalPlaces` 2 (bawaan grid master, K14): `{}` = pxNumber tanpa
 *  `pyDecimalPlaces` (pola inti), `'mentah'` = sel tanpa Format (apa adanya). */
export interface Kolom {
  m: string
  label: string
  format?: Sajian | 'mentah'
}

const k = (m: string, label: string, format?: Sajian | 'mentah'): Kolom =>
  format === undefined ? { m, label } : { m, label, format }
/** pxNumber tanpa `pyDecimalPlaces`. */
const POLA_INTI: Sajian = {}

/** Grid "Limits" - `pyWorkPage.TreatyIn.LimitSummaryList`. */
export const KOLOM_LIMIT: Kolom[] = [
  k('Note', 'Note'),
  k('Limit', '100% Limit (IDR)'),
  k('Limit2', '100% Limit (USD)'),
  k('Deductible', 'Deductible (IDR)'),
  k('Deductible2', 'Deductible (USD)'),
  k('MDP', 'MDP (IDR)'),
  k('MDP2', 'MDP (USD)'),
]

/** Grid "Share" - `pyWorkPage.TreatyIn.LimitShareSummaryList` (S20). Sel `.Limit` pxNumber
 *  TANPA `pyDecimalPlaces`, sel angka lain 2 (audit silang P3 W6). */
export const KOLOM_SHARE: Kolom[] = [
  k('Note', 'Note'),
  k('Limit', '100% Limit (IDR)', POLA_INTI),
  k('Limit2', '100% Limit (USD)'),
  k('MDP', 'MDP (IDR)'),
  k('MDP2', 'MDP (USD)'),
  k('Deductible', 'Brokerage (IDR)'),
  k('Deductible2', 'Brokerage (USD)'),
  k('NetPremi', 'Net Premi (IDR)'),
  k('NetPremiAfterPPN', 'Net After PPN(IDR)'),
  k('NetPremiAfterPPH', 'Net After PPH(IDR)'),
  k('NetPremi2', 'Net Premi (USD)'),
  k('NetPremiAfterPPN2', 'Net After PPN(USD)'),
  k('NetPremiAfterPPH2', 'Net After PPH(USD)'),
]

/** Grid "Share Facultative" - `pyWorkPage.TreatyIn.LimitFacShareSummaryList` (S52):
 *  seluruh sel FIELD TANPA Format (bukan pxNumber) - nilai apa adanya (audit silang P3 W6). */
export const KOLOM_FAKULTATIF: Kolom[] = [
  k('Note', 'Note'),
  k('Limit', '100% Limit (IDR)', 'mentah'),
  k('Limit2', '100% Limit (USD)', 'mentah'),
  k('MDP', 'MDP (IDR)', 'mentah'),
  k('MDP2', 'MDP (USD)', 'mentah'),
  k('Deductible', 'Brokerage (IDR)', 'mentah'),
  k('Deductible2', 'Brokerage (USD)', 'mentah'),
  k('NetPremi', 'Net Premi (IDR)', 'mentah'),
  k('NetPremi2', 'Net Premi (USD)', 'mentah'),
]

/** Satu grid total per mata uang: daftar master, label, kolom tambahan, dan kolom
 *  pertama FIELD tanpa properti (`kosongAwal`: S7 TotalLimitIOONP dan S57
 *  TotalFacShareRnmNP C[2.1]; judul C[1.3] di atas `.Value`). */
export interface Total {
  daftar: string
  label: string
  tambahan?: Kolom[]
  kosongAwal?: boolean
}

export const TOTAL_LIMIT: Total[] = [
  { daftar: 'TotalLimitIOONP', label: 'Total Limit', kosongAwal: true },
  { daftar: 'TotalLimitDeductblNP', label: 'Total Deductible' },
  { daftar: 'TotalLimitMDPNP', label: 'Total MDP' },
]

export const TOTAL_SHARE: Total[] = [
  { daftar: 'TotalShareRnmNP', label: 'Total Limit' },
  { daftar: 'TotalShareGrossNP', label: 'Total MDP' },
  {
    daftar: 'TotalShareDeductionNP',
    label: 'Total Brokerage',
    tambahan: [k('TotalPPNValue', 'Total PPN'), k('TotalPPHValue', 'Total PPH')],
  },
  {
    daftar: 'TotalShareNetNP',
    label: 'Total Net Premi Before PPn PPh',
    tambahan: [k('TotalNetPremiAfterPPN', 'Total Net Premi After PPn'), k('TotalNetPremiAfterTax', 'Total Net Premi After PPn PPh')],
  },
]

export const TOTAL_SPREADED: Total[] = [
  { daftar: 'TotalSpreadedNetPremi', label: 'Total OR Net Premium' },
  { daftar: 'TotalSpreadedNetPremiRI', label: 'Total R/I Net Premium' },
]

export const TOTAL_FAKULTATIF: Total[] = [
  { daftar: 'TotalFacShareRnmNP', label: 'Total Limit', kosongAwal: true },
  { daftar: 'TotalFacShareGrossNP', label: 'Total MDP' },
  { daftar: 'TotalFacShareDeductionNP', label: 'Total Brokerage' },
  { daftar: 'TotalFacShareNetNP', label: 'Total Net Premi' },
]

/** Kolom grid total: LABEL judul di atas kolom `.Value`, kolom `.Currency` tanpa
 *  judul - kecuali grid Total Spreaded (label di atas `.Currency`, "Value" di atas
 *  `.Value`). LABEL judul `DetailPolicyTreatyInNonProportional` (audit silang P3). */
export function kolomTotal(t: Total): Kolom[] {
  const dasar = t.daftar.startsWith('TotalSpreaded')
    ? [k('Currency', t.label), k('Value', 'Value')]
    : [k('Currency', ''), k('Value', t.label)]
  return [...(t.kosongAwal ? [k('', '')] : []), ...dasar, ...(t.tambahan ?? [])]
}

/** Section `InstallmentList` (flow action `InstallmentList`, hanya-baca) - rincian satu angsuran. */
export const KOLOM_RINCI: Kolom[] = [
  k('DueDate', 'Payment Date'),
  k('InstallmentPercentage', 'Percentage'),
  // LABEL judul kolom ke-3 Section InstallmentList kosong
  k('Currency', ''),
  k('Premium', 'Balance Before Tax'),
  k('PremiumAfterPPN', 'Balance Before Withholding Tax (PPH 2.2)'),
  k('PremiumAfterTax', 'Balance Due To'),
]

export const JUDUL_NONPROP = {
  /** LABEL sel `DetailPoliciesNonProportional` - tampil menurut tempat berperan (tiket 05). */
  nonEdm: 'NON EDM',
  edm: 'EDM',
  limits: 'Limits',
  share: 'Share',
  totalSpreaded: 'Total Spreaded',
  fakultatif: 'Share Facultative',
  rnmShare: '% RNM Share',
  fakultatifPersen: '% Share Facultative',
  currency: 'Currency',
  value: 'Value',
} as const

/** `% RNM Share`: RNMShare bila FacultativeShare = 0, RnmShareDeducted bila != 0. */
export const jalurRnmShare = (h: Halaman) => MASTER + (fakultatifNol(h) ? 'RNMShare' : 'RnmShareDeducted')

/** Baris daftar master (kosong bila tidak ada). */
export function daftarMaster(h: Halaman, nama: string): Baris[] {
  return h.daftar?.[MASTER + nama] ?? []
}
