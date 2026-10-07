// Definisi medan kedua layar realisasi - VERBATIM sel section Pega.
//
//   Admin   `Section/DetailPolicyTreatyIn`        (flow action InboxPolicyTreatyIn)
//   Atasan  `Section/DetailDeptHeadTreatyIn_UW`   (flow action DeptHeadTreatyIn_UW,
//           lewat `GeneralDeptHeadTreatyIn_UW` yang nol sel sendiri)
//
// Setiap medan membawa: jalur halaman, label, jenis kontrol, syarat tampil
// (`pyVisible` sel DAN `pyContainerVisibleWhen` wadahnya - dibaca ulang
// 2026-10-03), terkunci (`pyReadOnly` / `pyDisabled`), aksi refresh BERURUTAN
// (action set sel: `change -> refresh <Activity>` satu atau lebih), mata uang
// pasangan (AC 85), dan sajian angka/tanggal (`sajian.ts`, K14). Medan wajib
// TIDAK ditulis di sini: daftar wajib datang dari backend (`Layar.medanWajib`)
// supaya satu sumber.
//
// ⛔ Tidak ditampilkan, dan sebabnya:
//   - elemen bersyarat tampil `1=2` / `NEVER` (80 elemen mati, AC 53):
//     `.BizName`, `.DueTo` (admin), label pemisah, label bagian, `.Total*`
//     berlabel layar admin (wadah `1=2`);
//   - tombol `Select Source Of Business` (admin, `.ClaimType = 'XOL Retro'`) BUKAN
//     medan: dibangun paket P3 di `components/PilihSumberBisnis.tsx` (RALAT tiket 11);
//   - tombol `Choose Business R` (wadah `.ClaimType = 'XOL Retro'`) dan subsection
//     `DetailPolicyTreatyOutNonProportional` (wadah `.IsNewPolicyNonProp = 1 &&
//     .ClaimType = 'XOL Retro'`): treaty KELUAR, K8 butir 4 tetap tidak dibangun.
//     Subsection `DetailPoliciesNonProportional` DIBANGUN (K8) -
//     `components/DetailNonProp.tsx`, bukan medan di berkas ini;
//   - tombol `Survey Report` (`HistoricalSurveyReport[UW]`) - K7: tidak dibangun,
//     tidak ada tabel di diagram grilling;
//   - grid `Breakdown Spreading` (`.BreakDownSpreadList`) - K9 (KEPUTUSAN-RONDE-12
//     butir 3/3b): breakdown spreading tidak dimigrasi.

import { KLAIM_XOL_RETRO, POLIS, nilai, type Baris, type Halaman, type Pilihan } from './api'
import {
  BAGIAN,
  PILIHAN_CLAIM_PAYMENT_TYPE,
  PILIHAN_CLAIM_TYPE,
  PILIHAN_STATEMENT_TYPE,
  PILIHAN_SURVEY_REPORT,
  PILIHAN_TYPE_TAX,
} from './labels'
import type { Sajian } from './sajian'
import { negatifTeks } from './tanda'

export type JenisMedan =
  | 'tampil'
  | 'teks'
  | 'angka'
  | 'tanggal'
  | 'area'
  | 'centang'
  | 'mataUang'
  | 'mo'
  | 'pilihan'
  | 'radio'

/** Satu pilihan dropdown: nilai standar disimpan, teks prompt ditampilkan. */
export interface OpsiMedan {
  value: string
  label: string
}

/** Satu refresh berhitung (`POST .../hitung`): aksi backend + parameternya. */
export interface Aksi {
  aksi: string
  param?: string
}

export interface Medan {
  jalur: string
  label: string
  jenis: JenisMedan
  /** Syarat tampil yang berlaku (sel DAN wadah); tidak ada = selalu tampil. */
  tampil?: (h: Halaman) => boolean
  /** Terkunci (`pyReadOnly`, atau `pyDisabled` always) - selalu hanya-baca. */
  kunci?: boolean
  /** Refresh sesudah berubah - BERURUTAN seperti action set sel XML. */
  aksi?: Aksi[]
  /** Kode mata uang pasangan angka uang (AC 85). */
  mataUang?: string
  /** Penyajian nilai (K14): format angka sel, atau tanggal (AC 33). */
  sajian?: Sajian
  /** Daftar pilihan pxDropdown / pxRadioButtons ber-pyListSource `associated` (prompt values property). */
  opsi?: OpsiMedan[]
}

/** Teks prompt sebuah nilai standar; nilai di luar daftar tampil apa adanya (tidak dibuang). */
export function teksPilihan(opsi: OpsiMedan[], v: string): string {
  return opsi.find((o) => o.value === v)?.label ?? v
}

const v = (h: Halaman, m: string) => nilai(h, POLIS + m)
const mu = POLIS + 'Currency'

// ------------------------------------------------------------------ syarat

/** `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`. */
export const bukanNonProp = (h: Halaman) => nilai(h, 'Quotation.ProportionalType') !== 'NonProportional'
/** `.IsNewPolicyNonProp != 1`. */
export const bukanNonPropBaru = (h: Halaman) => v(h, 'IsNewPolicyNonProp') !== '1'
/** `.BalanceDueTo < 0` (teks kosong = 0, aritmetika Pega). */
export const saldoNegatif = (h: Halaman) => negatifTeks(v(h, 'BalanceDueTo'))
/** `.ClaimType != 'XOL Retro'` - wadah FlagPPH/TypeTax/Choose Business admin, dan
 *  `pyVisible` FlagRetroTreaty. */
export const bukanXOLRetro = (h: Halaman) => v(h, 'ClaimType') !== KLAIM_XOL_RETRO
/** `.QuotationData.ProportionalType = 'Proportional'` - wadah Quartal / U/Y. */
const proporsionalQD = (h: Halaman) => v(h, 'QuotationData.ProportionalType') === 'Proportional'
/** `.QuotationData.ProportionalType = 'NonProportional'` - wadah Layer*. */
const nonProporsionalQD = (h: Halaman) => v(h, 'QuotationData.ProportionalType') === 'NonProportional'
/** Wadah bagian uang / spreading / angsuran admin:
 *  `.IsNewPolicyNonProp != 1 && .IsNewPolicyListFormat != 1`. */
export const wadahUangAdmin = (h: Halaman) => bukanNonPropBaru(h) && v(h, 'IsNewPolicyListFormat') !== '1'
/** Wadah bagian uang / spreading / angsuran atasan: `.IsNewPolicyNonProp != 1`. */
export const wadahUangAtasan = bukanNonPropBaru
/** `pyVisible NOTBLANK`. */
const tidakKosong = (jalur: string) => (h: Halaman) => nilai(h, jalur).trim() !== ''

/** Medan dalam wadah bersyarat: syarat sel DAN syarat wadah. */
function dalamWadah(wadah: (h: Halaman) => boolean, ms: Medan[]): Medan[] {
  return ms.map((m) => {
    const sel = m.tampil
    return { ...m, tampil: sel ? (h: Halaman) => wadah(h) && sel(h) : wadah }
  })
}

// ------------------------------------------------------------------ sajian (K14)
//
// Dari `pyModes` sel (`scratchpad format-layar-mentah`, INVENTARIS bab 5):
//   pxNumber `pyDecimalPlaces` 2/4      -> { desimal }
//   pxCurrency tanpa `pyDecimalPlaces`  -> {} (tak terbaca: pola inti)
//   pxTextInput Balance* `pyFormatType number` tanpa desimal -> {}
//   `pyShowReadonlyFormatting=true` mode sunting -> formatSaatSunting

const UANG: Sajian = {}
const DUA: Sajian = { desimal: 2 }
/** Bagian uang (Gross, OGP, ONP, klaim, saldo, potongan, pajak) - perintah work owner 06-10-2026: nol / kosong
 *  tampil "0", angka terisi 4 angka di belakang koma; menggantikan K14 (`pyDecimalPlaces` 2 / pola inti). */
const UANG4: Sajian = { desimal: 4, nolPolos: true }
const UANG4_SUNTING: Sajian = { desimal: 4, nolPolos: true, formatSaatSunting: true }
/** pxTextInput `pyFormatType number`, `pyDecimalPlaces 0`, `pySeparators false`. */
const BULAT_POLOS: Sajian = { desimal: 0, ribuan: false }
const TGL: Sajian = 'tanggal'

// ------------------------------------------------------------------ admin

/** Medan umum layar admin - urutan sel `DetailPolicyTreatyIn` (wadah "General"). */
export const MEDAN_ADMIN_UMUM: Medan[] = [
  { jalur: POLIS + 'NoOffer', label: 'Master ID', jenis: 'tampil' },
  { jalur: 'TreatyIn.Commencement', label: 'Commencement', jenis: 'tampil', sajian: TGL },
  // change -> refresh ber-pyPreDataTransform `SystemSetOneYear_DT` (EndDate = StartDate + 1 tahun)
  { jalur: POLIS + 'StartDate', label: 'Statement Period', jenis: 'tanggal', aksi: [{ aksi: 'SystemSetOneYear' }] },
  { jalur: POLIS + 'SOBName', label: 'Source Of Business', jenis: 'tampil' },
  { jalur: POLIS + 'TreatyGroupName', label: 'Treaty Group', jenis: 'tampil', tampil: bukanNonPropBaru },
  // pyReadOnly `IsUW` (workbasket ReasFacIn*) - tidak pernah benar bagi admin treaty
  // pxRadioButtons (prompt values property: Yes / No)
  {
    jalur: POLIS + 'QuotationData.IsSurveyReport',
    label: 'Survey Report',
    jenis: 'radio',
    opsi: PILIHAN_SURVEY_REPORT,
    tampil: bukanNonProp,
  },
  // pxDropdown pyListSource `associated`, pyHasNoSelection true (prompt values property)
  { jalur: POLIS + 'StatementType', label: 'Statement Type', jenis: 'pilihan', opsi: PILIHAN_STATEMENT_TYPE },
  { jalur: POLIS + 'QuotationData.NoOfferSlip', label: 'No Offer Slip', jenis: 'area' },
  // pyCheckboxCaption sel `.FlagRetroTreaty` / `.FlagPPH` (label sel = teks bawaan "Checkbox")
  { jalur: POLIS + 'FlagRetroTreaty', label: 'Overiding Commision', jenis: 'centang', tampil: bukanXOLRetro },
  // change -> postValue -> runActivity RemoveTypeTax_ACT
  { jalur: POLIS + 'FlagPPH', label: 'With Tax', jenis: 'centang', tampil: bukanXOLRetro, aksi: [{ aksi: 'RemoveTypeTax' }] },
  // pxRadioButtons pyListSource `associated` (prompt values property: Inclusive / Exclusive). XML hanya
  // postValue; keputusan work owner 06-10-2026: pilihan langsung menghitung ulang pajak (server, pemicu Pajak)
  {
    jalur: POLIS + 'TypeTax',
    label: 'Type Tax',
    jenis: 'radio',
    opsi: PILIHAN_TYPE_TAX,
    tampil: (h) => bukanXOLRetro(h) && v(h, 'FlagPPH') === 'true',
    aksi: [{ aksi: 'HitungPajak' }],
  },
  { jalur: POLIS + 'ShareCurrency', label: 'RNM Share', jenis: 'tampil' },
  { jalur: POLIS + 'ShareValue', label: 'ShareValue', jenis: 'tampil', mataUang: POLIS + 'ShareCurrency', sajian: UANG },
  { jalur: POLIS + 'StatementDate', label: 'Statement Date', jenis: 'tampil', sajian: TGL },
  { jalur: 'TreatyIn.Termination', label: 'Termination', jenis: 'tampil', sajian: TGL },
  { jalur: POLIS + 'EndDate', label: 'To', jenis: 'tanggal', aksi: [{ aksi: 'ProtectDate' }] },
  { jalur: POLIS + 'CedingCoName', label: 'Ceding Company', jenis: 'tampil' },
  // sel tanpa label: teks yang TAMPIL di layar Pega (screenshot work owner 06-10-2026)
  { jalur: POLIS + 'InsuredName', label: 'Insured Name', jenis: 'tampil' },
  { jalur: POLIS + 'TreatyType', label: 'Treaty Type', jenis: 'tampil' },
  { jalur: POLIS + 'TreatyYear', label: 'UW Year', jenis: 'tampil', tampil: bukanNonPropBaru },
  // pyFormatType number tanpa desimal (mode baca); pyShowReadonlyFormatting false
  // sel tanpa label - LABEL tampil di depannya: "Q" [Quartal] "/" [YearOfQuartal] "U/Y" [TreatyYear]
  { jalur: POLIS + 'Quartal', label: 'Q', jenis: 'teks', tampil: proporsionalQD, sajian: UANG },
  { jalur: POLIS + 'YearOfQuartal', label: '/', jenis: 'teks', tampil: proporsionalQD },
  // sel `.TreatyYear` kedua, sesudah label "U/Y" (tampil ALWAYS)
  { jalur: POLIS + 'TreatyYear', label: 'U/Y', jenis: 'tampil', tampil: proporsionalQD },
  // change -> postValue -> refresh CheckDataMkt
  { jalur: POLIS + 'QuotationData.MOID', label: 'Marketing Officer', jenis: 'mo', aksi: [{ aksi: 'CheckDataMkt' }] },
  { jalur: POLIS + 'QuotationData.ProportionalType', label: 'Proportional Type', jenis: 'tampil', tampil: bukanNonPropBaru },
  // change -> refresh SetCurrency_act(CURR=.IDCurrency)
  {
    jalur: POLIS + 'IDCurrency',
    label: 'Currency',
    jenis: 'mataUang',
    tampil: bukanNonPropBaru,
    aksi: [{ aksi: 'SetCurrency' }],
  },
  // pxDropdown pyListSource `associated`, pyHasNoSelection true (prompt values property)
  { jalur: POLIS + 'ClaimType', label: 'Claim Type', jenis: 'pilihan', opsi: PILIHAN_CLAIM_TYPE },
  { jalur: POLIS + 'ClaimPaymentType', label: 'Payment Type', jenis: 'pilihan', opsi: PILIHAN_CLAIM_PAYMENT_TYPE },
  ...dalamWadah(nonProporsionalQD, [
    { jalur: POLIS + 'LayerType', label: 'LayerType', jenis: 'tampil' },
    { jalur: POLIS + 'Layer', label: 'Layer', jenis: 'tampil' },
    // LABEL "Of" di depan sel (pyVisible ALWAYS)
    { jalur: POLIS + 'LayerPartType', label: 'Of', jenis: 'tampil' },
    { jalur: POLIS + 'LayerPart', label: 'LayerPart', jenis: 'tampil' },
  ]),
  { jalur: POLIS + 'Remark', label: 'Remark', jenis: 'area' },
]

/** Satu medan uang tersunting layar admin. */
const uangAdmin = (m: string, label: string, aksi?: Aksi[], sajian: Sajian = UANG4_SUNTING, uang = true): Medan => ({
  jalur: POLIS + m,
  label,
  jenis: 'angka',
  aksi,
  mataUang: uang ? mu : undefined,
  sajian,
})
/** Persen/hasil OGP-ONP: `refresh <Act>(Data=..)` LALU `refresh CountOGPONP_Act`. */
const lalu = (aksi: string, param: string): Aksi[] => [{ aksi, param }, { aksi: 'CountOGPONP' }]
const OGPONP: Aksi[] = [{ aksi: 'CountOGPONP' }]

/** Medan uang layar admin - bagian "Old Soa Input Format" (wadah
 *  `.IsNewPolicyNonProp != 1 && .IsNewPolicyListFormat != 1`). */
export const MEDAN_ADMIN_UANG: Medan[] = dalamWadah(wadahUangAdmin, [
  uangAdmin('GrossPremium', 'Gross Premium 100%', [{ aksi: 'CalculatePremi', param: 'PREMIUM' }], UANG4_SUNTING),
  uangAdmin('GrossClaim', 'Claim 100%', [{ aksi: 'CalculatePremi', param: 'CLAIM' }], UANG4_SUNTING),
  uangAdmin('PremiOgp', 'Premi Ogp', OGPONP),
  uangAdmin('RiCommOgp', '(%) Deduction In A (OGP)', lalu('CountResult1', 'Pct'), UANG4_SUNTING, false),
  uangAdmin('ResultOgp1', 'Deduction In A (OGP)', lalu('CountResult1', 'Amount')),
  uangAdmin('OveriddingCommOgp', '(%) Deduction In B (OGP)', lalu('CountResult2Ogp', 'Pct'), UANG4_SUNTING, false),
  uangAdmin('ResultOgp2', 'Deduction In B (OGP)', lalu('CountResult2Ogp', 'Amount')),
  uangAdmin('PremiOnp', 'Premi Onp', OGPONP),
  uangAdmin('RiCommOnp', '(%) Deduction In A (ONP)', lalu('CountResult1Onp', 'Pct'), UANG4_SUNTING, false),
  uangAdmin('ResultOnp1', 'Deduction In A (ONP)', lalu('CountResult1Onp', 'Amount')),
  uangAdmin('OveriddingCommOnp', '(%) Deduction In B (ONP)', lalu('CountResult2Onp', 'Pct'), UANG4_SUNTING, false),
  uangAdmin('ResultOnp2', 'Deduction In B (ONP)', lalu('CountResult2Onp', 'Amount')),
  uangAdmin('Claim', 'Claim', OGPONP),
  uangAdmin('OutstandingClaim', 'Outstanding Claim'),
  uangAdmin('SalvageValue', 'Salvage', OGPONP),
  uangAdmin('ExcessLoss', 'Excess Loss', OGPONP),
  { jalur: POLIS + 'NetPremium', label: 'Total Premium Before Claim', jenis: 'tampil', mataUang: mu, sajian: UANG4 },
  { jalur: POLIS + 'BalanceDueTo', label: 'Balance Due To You', jenis: 'tampil', tampil: saldoNegatif, mataUang: mu, sajian: UANG4 },
  {
    jalur: POLIS + 'BalanceBeforeTax',
    label: 'Balance Before Tax',
    jenis: 'tampil',
    tampil: (h) => !saldoNegatif(h),
    mataUang: mu,
    sajian: UANG4,
  },
  {
    jalur: POLIS + 'BalanceBeforePPH',
    label: 'Balance Before Withholding Tax (PPH 2.2)',
    jenis: 'tampil',
    tampil: (h) => !saldoNegatif(h),
    mataUang: mu,
    sajian: UANG4,
  },
  {
    jalur: POLIS + 'BalanceDueTo',
    label: 'Balance Due To Us',
    jenis: 'tampil',
    tampil: (h) => !saldoNegatif(h),
    mataUang: mu,
    sajian: UANG4,
  },
  // K3: pxCurrency, dipakai rumus sebagai jumlah uang (CountNetPremi_act langkah 4)
  uangAdmin('Deduction1', 'Deduction1', OGPONP),
  uangAdmin('Deduction2', 'Deduction2', OGPONP),
  { jalur: POLIS + 'PPHValue', label: 'PPH 2%', jenis: 'tampil', mataUang: mu, sajian: UANG4 },
  { jalur: POLIS + 'PPNValue', label: 'PPN 2.2%', jenis: 'tampil', mataUang: mu, sajian: UANG4 },
])

// ------------------------------------------------------------------ atasan

/** Medan layar atasan (`DetailDeptHeadTreatyIn_UW`, wadah "General") - SELURUHNYA
 *  hanya-baca: `pyReadOnly`, atau `pyDisabled=always` (DueTo, FlagPPH, No Offer
 *  Slip). Yang dapat diisi atasan hanya `ListSuggest` (AC 52). */
export const MEDAN_ATASAN_UMUM: Medan[] = [
  { jalur: POLIS + 'NoOffer', label: 'Master ID', jenis: 'tampil' },
  { jalur: 'TreatyIn.Commencement', label: 'Commencement', jenis: 'tampil', sajian: TGL },
  { jalur: POLIS + 'StartDate', label: 'Statement Period', jenis: 'tampil', sajian: TGL },
  { jalur: POLIS + 'SOBName', label: 'Source Of Business', jenis: 'tampil' },
  { jalur: POLIS + 'TreatyGroupName', label: 'Treaty Group', jenis: 'tampil', tampil: bukanNonPropBaru },
  // wadah `pyWorkPage.FlagViewPolicy = 1` - FlagViewPolicy tidak diisi rule
  // terjangkau mana pun (hanya dibaca `SumTSIPremiSpreadedRNM_Act`, tak terjangkau)
  ...dalamWadah((h) => nilai(h, 'FlagViewPolicy') === '1', [
    { jalur: POLIS + 'PolicyNo', label: 'No Polis', jenis: 'tampil' },
    { jalur: POLIS + 'ProductionDate', label: 'Production Date', jenis: 'tampil', sajian: TGL },
  ]),
  { jalur: POLIS + 'DueTo', label: 'Due To Us / You', jenis: 'tampil' },
  // radio / dropdown Read-only: teks prompt nilai tersimpan
  {
    jalur: POLIS + 'QuotationData.IsSurveyReport',
    label: 'Survey Report',
    jenis: 'tampil',
    opsi: PILIHAN_SURVEY_REPORT,
    tampil: bukanNonProp,
  },
  { jalur: POLIS + 'StatementType', label: 'Statement Type', jenis: 'tampil', opsi: PILIHAN_STATEMENT_TYPE },
  { jalur: POLIS + 'FlagPPH', label: 'Include Tax', jenis: 'centang', kunci: true },
  { jalur: POLIS + 'TypeTax', label: 'Type Tax', jenis: 'tampil', tampil: tidakKosong(POLIS + 'TypeTax') },
  { jalur: POLIS + 'QuotationData.NoOfferSlip', label: 'No Offer Slip', jenis: 'tampil' },
  { jalur: POLIS + 'ShareCurrency', label: 'RNM Share', jenis: 'tampil' },
  { jalur: POLIS + 'ShareValue', label: 'ShareValue', jenis: 'tampil', mataUang: POLIS + 'ShareCurrency', sajian: UANG },
  { jalur: POLIS + 'StatementDate', label: 'Statement Date', jenis: 'tampil', sajian: TGL },
  { jalur: 'TreatyIn.Termination', label: 'Termination', jenis: 'tampil', sajian: TGL },
  { jalur: POLIS + 'EndDate', label: 'To', jenis: 'tampil', sajian: TGL },
  { jalur: POLIS + 'CedingCoName', label: 'Ceding Company', jenis: 'tampil', tampil: tidakKosong(POLIS + 'CedingCoName') },
  { jalur: POLIS + 'InsuredName', label: 'Insured Name', jenis: 'tampil' },
  { jalur: POLIS + 'TreatyType', label: 'Treaty Type', jenis: 'tampil' },
  { jalur: POLIS + 'TreatyYear', label: 'UW Year', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: POLIS + 'Quartal', label: 'Q', jenis: 'tampil', tampil: proporsionalQD, sajian: BULAT_POLOS },
  { jalur: POLIS + 'YearOfQuartal', label: '/', jenis: 'tampil', tampil: proporsionalQD, sajian: BULAT_POLOS },
  { jalur: POLIS + 'TreatyYear', label: 'U/Y', jenis: 'tampil', tampil: proporsionalQD },
  {
    jalur: POLIS + 'MarketingOfficer',
    label: 'Marketing Officer',
    jenis: 'tampil',
    tampil: tidakKosong(POLIS + 'MarketingOfficer'),
  },
  { jalur: POLIS + 'QuotationData.ProportionalType', label: 'Proportional Type', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: POLIS + 'Currency', label: 'Currency', jenis: 'tampil', tampil: bukanNonPropBaru },
  // dropdown Read-only: teks prompt nilai tersimpan
  { jalur: POLIS + 'ClaimType', label: 'Claim Type', jenis: 'tampil', opsi: PILIHAN_CLAIM_TYPE },
  { jalur: POLIS + 'ClaimPaymentType', label: 'Claim Payment Type', jenis: 'tampil', opsi: PILIHAN_CLAIM_PAYMENT_TYPE },
  ...dalamWadah(nonProporsionalQD, [
    { jalur: POLIS + 'LayerType', label: 'LayerType', jenis: 'tampil' },
    { jalur: POLIS + 'Layer', label: 'Layer', jenis: 'tampil' },
    // LABEL "Of" di depan sel (pyVisible ALWAYS)
    { jalur: POLIS + 'LayerPartType', label: 'Of', jenis: 'tampil' },
    { jalur: POLIS + 'LayerPart', label: 'LayerPart', jenis: 'tampil' },
  ]),
  { jalur: POLIS + 'Remark', label: 'Remark', jenis: 'tampil' },
]

/** Bagian uang layar atasan - label VERBATIM `DetailDeptHeadTreatyIn_UW`
 *  (berbeda dari layar admin: tanpa "(%)", tanpa Claim 100%). Wadah
 *  `.IsNewPolicyNonProp != 1`. */
export const MEDAN_ATASAN_UANG: Medan[] = dalamWadah(wadahUangAtasan, medanAtasanUang())

function medanAtasanUang(): Medan[] {
  const t = (m: string, label: string, sajian: Sajian = UANG4, uang = true): Medan => ({
    jalur: POLIS + m,
    label,
    jenis: 'tampil',
    mataUang: uang ? mu : undefined,
    sajian,
  })
  return [
    t('GrossPremium', 'Gross Premium 100%', UANG4),
    t('PremiOgp', 'Premi Ogp'),
    t('RiCommOgp', 'Deduction In A (OGP)', UANG4, false),
    t('ResultOgp1', 'ResultOgp1'),
    t('OveriddingCommOgp', 'Deduction In B (OGP)', UANG4, false),
    t('ResultOgp2', 'ResultOgp2'),
    t('Claim', 'Claim'),
    t('OutstandingClaim', 'Outstanding Claim'),
    t('SalvageValue', 'Salvage'),
    t('ExcessLoss', 'Excess Loss'),
    t('NetPremium', 'Net Premium'),
    { ...t('BalanceDueTo', 'Balance Due To You'), tampil: saldoNegatif },
    { ...t('BalanceBeforeTax', 'Balance Before Tax'), tampil: (h) => !saldoNegatif(h) },
    { ...t('BalanceBeforePPH', 'Balance Before Withholding Tax (PPH 2.2)'), tampil: (h) => !saldoNegatif(h) },
    { ...t('BalanceDueTo', 'Balance Due To Us'), tampil: (h) => !saldoNegatif(h) },
    t('PremiOnp', 'Premi Onp'),
    t('RiCommOnp', 'Deduction In A (ONP)', UANG4, false),
    t('ResultOnp1', 'ResultOnp1'),
    t('OveriddingCommOnp', 'Deduction In B (ONP)', UANG4, false),
    t('ResultOnp2', 'ResultOnp2'),
    // K3: pxCurrency, jumlah uang
    t('Deduction1', 'Deduction1'),
    t('Deduction2', 'Deduction2'),
    t('PPHValue', 'PPH 2%'),
    t('PPNValue', 'PPN 2.2%'),
  ]
}

// ------------------------------------------------------------------ kelompok wadah uang

/** Satu wadah sel bagian uang; `judul` = LABEL `pyIncludeLabel=true` (Heading 4) bila ada. */
export interface Kelompok {
  judul?: string
  medan: Medan[]
}

/** Potong `ms` menurut jalur: dari medan `dari` (inklusif) sampai sebelum `sampai`. */
function antara(ms: Medan[], dari: string, sampai?: string): Medan[] {
  const i = ms.findIndex((m) => m.jalur === POLIS + dari)
  const j = sampai === undefined ? ms.length : ms.findIndex((m) => m.jalur === POLIS + sampai)
  return ms.slice(i, j)
}

/** Wadah bagian uang admin (`DetailPolicyTreatyIn` S19 - W6 audit silang P3, tanpa judul
 *  buatan): S20 Gross, S24 LABEL "OGP", S25 LABEL "ONP", S26-S28 klaim/saldo/potongan/pajak. */
export const KELOMPOK_UANG_ADMIN: Kelompok[] = [
  { medan: antara(MEDAN_ADMIN_UANG, 'GrossPremium', 'PremiOgp') },
  { judul: BAGIAN.ogp, medan: antara(MEDAN_ADMIN_UANG, 'PremiOgp', 'PremiOnp') },
  { judul: BAGIAN.onp, medan: antara(MEDAN_ADMIN_UANG, 'PremiOnp', 'Claim') },
  { medan: antara(MEDAN_ADMIN_UANG, 'Claim') },
]

/** Wadah bagian uang atasan (`DetailDeptHeadTreatyIn_UW` S90): S91 Gross, S93 LABEL "OGP"
 *  (sampai Balance*), S94 LABEL "ONP" (sampai PPN). */
export const KELOMPOK_UANG_ATASAN: Kelompok[] = [
  { medan: antara(MEDAN_ATASAN_UANG, 'GrossPremium', 'PremiOgp') },
  { judul: BAGIAN.ogp, medan: antara(MEDAN_ATASAN_UANG, 'PremiOgp', 'PremiOnp') },
  { judul: BAGIAN.onp, medan: antara(MEDAN_ATASAN_UANG, 'PremiOnp') },
]

/** Total berlabel di bawah grid spreading layar atasan (`.TotalSharePercentagePremium`
 *  dst., pxNumber 2 desimal). Di layar admin sel yang sama berwadah `1=2` (mati). */
export const MEDAN_ATASAN_TOTAL: Medan[] = dalamWadah(wadahUangAtasan, [
  { jalur: POLIS + 'TotalSharePercentagePremium', label: 'Total %Share', jenis: 'tampil', sajian: DUA },
  { jalur: POLIS + 'TotalPremium', label: 'Total Premium', jenis: 'tampil', mataUang: mu, sajian: DUA },
  { jalur: POLIS + 'TotalSharePercentageClaim', label: 'Total %Share Claim', jenis: 'tampil', sajian: DUA },
  { jalur: POLIS + 'TotalClaim', label: 'Total Claim', jenis: 'tampil', mataUang: mu, sajian: DUA },
])

// ------------------------------------------------------------------ grid (K14)

/** Sajian sel grid spreading `.SpreadingRiskList` kedua layar: pxNumber 4 desimal
 *  (persen tersunting `pyShowReadonlyFormatting=true`), total footer 4 desimal. */
/** Teks sel `.TreatyType` grid spreading yang tidak dapat disunting: dropdown ro - label
 *  pilihan daftar acuannya (`.ID` -> `.Note` / `.TreatyName`), lalu `.TreatyName` baris,
 *  lalu kode apa adanya. Satu pencari untuk layar admin, layar atasan, dan grid NonProp. */
export function labelTreatyType(b: Baris, opsi: Pilihan[]): string {
  return opsi.find((o) => o.nilai === b.TreatyType)?.label || b.TreatyName || b.TreatyType || ''
}

export const SAJIAN_SPREADING = {
  persen: { desimal: 4, formatSaatSunting: true } as Sajian,
  uang: { desimal: 4 } as Sajian,
  total: { desimal: 4 } as Sajian,
}

/** Sajian sel grid `.ListInstallment`: InstallmentNo pxInteger (apa adanya),
 *  DueDate tanggal, InstallmentPercentage & Premium 4 desimal, PaymentTotal 2. */
export const SAJIAN_ANGSURAN = {
  dueDate: TGL,
  persen: { desimal: 4, formatSaatSunting: true } as Sajian,
  premium: { desimal: 4, formatSaatSunting: true } as Sajian,
  total: DUA,
}

/** Medan tampil bagi halaman ini. */
export function medanTampil(daftar: Medan[], h: Halaman): Medan[] {
  return daftar.filter((m) => !m.tampil || m.tampil(h))
}
