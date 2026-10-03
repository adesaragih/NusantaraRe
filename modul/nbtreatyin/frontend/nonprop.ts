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
//   - label "NON EDM" / "EDM": tampil hanya untuk SATU identitas orang (tiket 05).

import { nilai, type Baris, type Halaman } from './api'

const P = 'PolicyTreatyIn.'
const M = 'TreatyIn.'

/** `.IsNewPolicyNonProp = 1` - kontainer proporsional tersembunyi, subsection NonProp tampil. */
export const polisNonPropBaru = (h: Halaman) => nilai(h, P + 'IsNewPolicyNonProp') === '1'

/** Kontainer `DetailPoliciesNonProportional`: `.IsNewPolicyNonProp = 1 && .ClaimType != 'XOL Retro'`. */
export const tampilNonProp = (h: Halaman) => polisNonPropBaru(h) && nilai(h, P + 'ClaimType') !== 'XOL Retro'

/** Angka master sebagai bilangan (kosong = 0) - hanya untuk SYARAT tampil, bukan perhitungan. */
function angka(s: string): number {
  const n = Number(s.trim() === '' ? '0' : s)
  return Number.isFinite(n) ? n : NaN
}

/** `pyWorkPage.TreatyIn.FacultativeShare` = 0 (kosong dihitung 0). */
export const fakultatifNol = (h: Halaman) => angka(nilai(h, M + 'FacultativeShare')) === 0

/** Kontainer "Share Facultative": `pyWorkPage.TreatyIn.FacultativeShare>0`. */
export const tampilFakultatif = (h: Halaman) => angka(nilai(h, M + 'FacultativeShare')) > 0

/**
 * Section `SpreadingRiskList`: tombol Add/Delete tampil bila FacultativeShare = 0 atau '';
 * TreatyType / %Share / %Share Claim terkunci bila FacultativeShare > 0.
 */
export const spreadingTerbuka = (h: Halaman) => fakultatifNol(h)

/** Satu kolom grid: anggota baris dan label sel VERBATIM. */
export interface Kolom {
  m: string
  label: string
}

const k = (m: string, label: string): Kolom => ({ m, label })

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

/** Grid "Share" - `pyWorkPage.TreatyIn.LimitShareSummaryList`. */
export const KOLOM_SHARE: Kolom[] = [
  k('Note', 'Note'),
  k('Limit', '100% Limit (IDR)'),
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

/** Grid "Share Facultative" - `pyWorkPage.TreatyIn.LimitFacShareSummaryList`. */
export const KOLOM_FAKULTATIF: Kolom[] = [
  k('Note', 'Note'),
  k('Limit', '100% Limit (IDR)'),
  k('Limit2', '100% Limit (USD)'),
  k('MDP', 'MDP (IDR)'),
  k('MDP2', 'MDP (USD)'),
  k('Deductible', 'Brokerage (IDR)'),
  k('Deductible2', 'Brokerage (USD)'),
  k('NetPremi', 'Net Premi (IDR)'),
  k('NetPremi2', 'Net Premi (USD)'),
]

/** Satu grid total per mata uang: daftar master, label, dan kolom tambahan. */
export interface Total {
  daftar: string
  label: string
  tambahan?: Kolom[]
}

export const TOTAL_LIMIT: Total[] = [
  { daftar: 'TotalLimitIOONP', label: 'Total Limit' },
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
  { daftar: 'TotalFacShareRnmNP', label: 'Total Limit' },
  { daftar: 'TotalFacShareGrossNP', label: 'Total MDP' },
  { daftar: 'TotalFacShareDeductionNP', label: 'Total Brokerage' },
  { daftar: 'TotalFacShareNetNP', label: 'Total Net Premi' },
]

/** Section `InstallmentList` (flow action `InstallmentList`, hanya-baca) - rincian satu angsuran. */
export const KOLOM_RINCI: Kolom[] = [
  k('DueDate', 'Payment Date'),
  k('InstallmentPercentage', 'Percentage'),
  k('Currency', 'Currency'),
  k('Premium', 'Balance Before Tax'),
  k('PremiumAfterPPN', 'Balance Before Withholding Tax (PPH 2.2)'),
  k('PremiumAfterTax', 'Balance Due To'),
]

export const JUDUL_NONPROP = {
  limits: 'Limits',
  share: 'Share',
  totalSpreaded: 'Total Spreaded',
  fakultatif: 'Share Facultative',
  rnmShare: '% RNM Share',
  fakultatifPersen: '% Share Facultative',
  installment: 'Installment',
  currency: 'Currency',
  value: 'Value',
} as const

/** `% RNM Share`: RNMShare bila FacultativeShare = 0, RnmShareDeducted bila != 0. */
export const jalurRnmShare = (h: Halaman) => M + (fakultatifNol(h) ? 'RNMShare' : 'RnmShareDeducted')

/** Baris daftar master (kosong bila tidak ada). */
export function daftarMaster(h: Halaman, nama: string): Baris[] {
  return h.daftar?.[M + nama] ?? []
}
