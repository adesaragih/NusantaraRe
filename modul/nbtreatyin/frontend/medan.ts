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
//   - tombol `Select Source Of Business` (admin, `.ClaimType = 'XOL Retro'`) -
//     dibangun paket P3 (pemilih SOB);
//   - tombol `Choose Business R` (wadah `.ClaimType = 'XOL Retro'`) dan subsection
//     `DetailPoliciesNonProportional` / `DetailPolicyTreatyOutNonProportional`
//     (wadah `.IsNewPolicyNonProp = 1`) - jalur NonProp/XOL, paket P5;
//   - tombol `Survey Report` (`HistoricalSurveyReport[UW]`) - K7: tidak dibangun,
//     tidak ada tabel di diagram grilling;
//   - grid `Breakdown Spreading` (`.BreakDownSpreadList`) - K9 (KEPUTUSAN-RONDE-12
//     butir 3/3b): breakdown spreading tidak dimigrasi.

import { nilai, type Halaman } from './api'
import type { Sajian } from './sajian'

export type JenisMedan = 'tampil' | 'teks' | 'angka' | 'tanggal' | 'area' | 'centang' | 'mataUang' | 'mo'

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
}

const P = 'PolicyTreatyIn.'
const v = (h: Halaman, m: string) => nilai(h, P + m)
const mu = P + 'Currency'

// ------------------------------------------------------------------ syarat

/** `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`. */
export const bukanNonProp = (h: Halaman) => nilai(h, 'Quotation.ProportionalType') !== 'NonProportional'
/** `.IsNewPolicyNonProp != 1`. */
export const bukanNonPropBaru = (h: Halaman) => v(h, 'IsNewPolicyNonProp') !== '1'
/** `.BalanceDueTo < 0` (teks kosong = 0, aritmetika Pega). */
export function saldoNegatif(h: Halaman): boolean {
  const s = v(h, 'BalanceDueTo').trim()
  return s !== '' && s.startsWith('-') && /[1-9]/.test(s)
}
/** `.ClaimType != 'XOL Retro'` - wadah FlagPPH/TypeTax/Choose Business admin, dan
 *  `pyVisible` FlagRetroTreaty. */
export const bukanXOLRetro = (h: Halaman) => v(h, 'ClaimType') !== 'XOL Retro'
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
const UANG_SUNTING: Sajian = { formatSaatSunting: true }
const DUA: Sajian = { desimal: 2 }
const DUA_SUNTING: Sajian = { desimal: 2, formatSaatSunting: true }
/** pxTextInput `pyFormatType number`, `pyDecimalPlaces 0`, `pySeparators false`. */
const BULAT_POLOS: Sajian = { desimal: 0, ribuan: false }
const TGL: Sajian = 'tanggal'

// ------------------------------------------------------------------ admin

/** Medan umum layar admin - urutan sel `DetailPolicyTreatyIn` (wadah "General"). */
export const MEDAN_ADMIN_UMUM: Medan[] = [
  { jalur: P + 'NoOffer', label: 'Master ID', jenis: 'tampil' },
  { jalur: 'TreatyIn.Commencement', label: 'Commencement', jenis: 'tampil', sajian: TGL },
  // change -> refresh (tanpa activity)
  { jalur: P + 'StartDate', label: 'Statement Period', jenis: 'tanggal' },
  { jalur: P + 'SOBName', label: 'Source Of Business', jenis: 'tampil' },
  { jalur: P + 'TreatyGroupName', label: 'Treaty Group', jenis: 'tampil', tampil: bukanNonPropBaru },
  // pyReadOnly `IsUW` (workbasket ReasFacIn*) - tidak pernah benar bagi admin treaty
  { jalur: P + 'QuotationData.IsSurveyReport', label: 'Survey Report', jenis: 'teks', tampil: bukanNonProp },
  { jalur: P + 'StatementType', label: 'Statement Type', jenis: 'teks' },
  { jalur: P + 'QuotationData.NoOfferSlip', label: 'No Offer Slip', jenis: 'area' },
  { jalur: P + 'FlagRetroTreaty', label: 'FlagRetroTreaty', jenis: 'centang', tampil: bukanXOLRetro },
  // change -> postValue -> runActivity RemoveTypeTax_ACT
  { jalur: P + 'FlagPPH', label: 'FlagPPH', jenis: 'centang', tampil: bukanXOLRetro, aksi: [{ aksi: 'RemoveTypeTax' }] },
  {
    jalur: P + 'TypeTax',
    label: 'Type Tax',
    jenis: 'teks',
    tampil: (h) => bukanXOLRetro(h) && v(h, 'FlagPPH') === 'true',
  },
  { jalur: P + 'ShareCurrency', label: 'RNM Share', jenis: 'tampil' },
  { jalur: P + 'ShareValue', label: 'ShareValue', jenis: 'tampil', mataUang: P + 'ShareCurrency', sajian: UANG },
  { jalur: P + 'StatementDate', label: 'Statement Date', jenis: 'tampil', sajian: TGL },
  { jalur: 'TreatyIn.Termination', label: 'Termination', jenis: 'tampil', sajian: TGL },
  { jalur: P + 'EndDate', label: 'To', jenis: 'tanggal', aksi: [{ aksi: 'ProtectDate' }] },
  { jalur: P + 'CedingCoName', label: 'Ceding Company', jenis: 'tampil' },
  { jalur: P + 'InsuredName', label: 'InsuredName', jenis: 'tampil' },
  { jalur: P + 'TreatyType', label: 'Treaty Type', jenis: 'tampil' },
  { jalur: P + 'TreatyYear', label: 'UW Year', jenis: 'tampil', tampil: bukanNonPropBaru },
  // pyFormatType number tanpa desimal (mode baca); pyShowReadonlyFormatting false
  { jalur: P + 'Quartal', label: '.Quartal', jenis: 'teks', tampil: proporsionalQD, sajian: UANG },
  { jalur: P + 'YearOfQuartal', label: 'YearOfQuartal', jenis: 'teks', tampil: proporsionalQD },
  // sel `.TreatyYear` kedua, sesudah label "U/Y" (tampil ALWAYS)
  { jalur: P + 'TreatyYear', label: 'U/Y', jenis: 'tampil', tampil: proporsionalQD },
  // change -> postValue -> refresh CheckDataMkt
  { jalur: P + 'QuotationData.MOID', label: 'Marketing Officer', jenis: 'mo', aksi: [{ aksi: 'CheckDataMkt' }] },
  { jalur: P + 'QuotationData.ProportionalType', label: 'Proportional Type', jenis: 'tampil', tampil: bukanNonPropBaru },
  // change -> refresh SetCurrency_act(CURR=.IDCurrency)
  {
    jalur: P + 'IDCurrency',
    label: 'Currency',
    jenis: 'mataUang',
    tampil: bukanNonPropBaru,
    aksi: [{ aksi: 'SetCurrency' }],
  },
  { jalur: P + 'ClaimType', label: 'Claim Type', jenis: 'teks' },
  { jalur: P + 'ClaimPaymentType', label: 'Payment Type', jenis: 'teks' },
  ...dalamWadah(nonProporsionalQD, [
    { jalur: P + 'LayerType', label: 'LayerType', jenis: 'tampil' },
    { jalur: P + 'Layer', label: 'Layer', jenis: 'tampil' },
    { jalur: P + 'LayerPartType', label: 'LayerPartType', jenis: 'tampil' },
    { jalur: P + 'LayerPart', label: 'LayerPart', jenis: 'tampil' },
  ]),
  { jalur: P + 'Remark', label: 'Remark', jenis: 'area' },
]

/** Satu medan uang tersunting layar admin. */
const uangAdmin = (m: string, label: string, aksi?: Aksi[], sajian: Sajian = UANG_SUNTING, uang = true): Medan => ({
  jalur: P + m,
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
  uangAdmin('GrossPremium', 'Gross Premium 100%', [{ aksi: 'CalculatePremi', param: 'PREMIUM' }], DUA_SUNTING),
  uangAdmin('GrossClaim', 'Claim 100%', [{ aksi: 'CalculatePremi', param: 'CLAIM' }], DUA_SUNTING),
  uangAdmin('PremiOgp', 'Premi Ogp', OGPONP),
  uangAdmin('RiCommOgp', '(%) Deduction In A (OGP)', lalu('CountResult1', 'Pct'), DUA_SUNTING, false),
  uangAdmin('ResultOgp1', 'Deduction In A (OGP)', lalu('CountResult1', 'Amount')),
  uangAdmin('OveriddingCommOgp', '(%) Deduction In B (OGP)', lalu('CountResult2Ogp', 'Pct'), DUA_SUNTING, false),
  uangAdmin('ResultOgp2', 'Deduction In B (OGP)', lalu('CountResult2Ogp', 'Amount')),
  uangAdmin('PremiOnp', 'Premi Onp', OGPONP),
  uangAdmin('RiCommOnp', '(%) Deduction In A (ONP)', lalu('CountResult1Onp', 'Pct'), DUA_SUNTING, false),
  uangAdmin('ResultOnp1', 'Deduction In A (ONP)', lalu('CountResult1Onp', 'Amount')),
  uangAdmin('OveriddingCommOnp', '(%) Deduction In B (ONP)', lalu('CountResult2Onp', 'Pct'), DUA_SUNTING, false),
  uangAdmin('ResultOnp2', 'Deduction In B (ONP)', lalu('CountResult2Onp', 'Amount')),
  uangAdmin('Claim', 'Claim', OGPONP),
  uangAdmin('OutstandingClaim', 'Outstanding Claim'),
  uangAdmin('SalvageValue', 'Salvage', OGPONP),
  uangAdmin('ExcessLoss', 'Excess Loss', OGPONP),
  { jalur: P + 'NetPremium', label: 'Total Premium Before Claim', jenis: 'tampil', mataUang: mu, sajian: UANG },
  { jalur: P + 'BalanceDueTo', label: 'Balance Due To You', jenis: 'tampil', tampil: saldoNegatif, mataUang: mu, sajian: UANG },
  {
    jalur: P + 'BalanceBeforeTax',
    label: 'Balance Before Tax',
    jenis: 'tampil',
    tampil: (h) => !saldoNegatif(h),
    mataUang: mu,
    sajian: UANG,
  },
  {
    jalur: P + 'BalanceBeforePPH',
    label: 'Balance Before Withholding Tax (PPH 2.2)',
    jenis: 'tampil',
    tampil: (h) => !saldoNegatif(h),
    mataUang: mu,
    sajian: UANG,
  },
  {
    jalur: P + 'BalanceDueTo',
    label: 'Balance Due To Us',
    jenis: 'tampil',
    tampil: (h) => !saldoNegatif(h),
    mataUang: mu,
    sajian: UANG,
  },
  // K3: pxCurrency, dipakai rumus sebagai jumlah uang (CountNetPremi_act langkah 4)
  uangAdmin('Deduction1', 'Deduction1', OGPONP),
  uangAdmin('Deduction2', 'Deduction2', OGPONP),
  { jalur: P + 'PPHValue', label: 'PPH 2%', jenis: 'tampil', mataUang: mu, sajian: UANG },
  { jalur: P + 'PPNValue', label: 'PPN 2.2%', jenis: 'tampil', mataUang: mu, sajian: UANG },
])

// ------------------------------------------------------------------ atasan

/** Medan layar atasan (`DetailDeptHeadTreatyIn_UW`, wadah "General") - SELURUHNYA
 *  hanya-baca: `pyReadOnly`, atau `pyDisabled=always` (DueTo, FlagPPH, No Offer
 *  Slip). Yang dapat diisi atasan hanya `ListSuggest` (AC 52). */
export const MEDAN_ATASAN_UMUM: Medan[] = [
  { jalur: P + 'NoOffer', label: 'Master ID', jenis: 'tampil' },
  { jalur: 'TreatyIn.Commencement', label: 'Commencement', jenis: 'tampil', sajian: TGL },
  { jalur: P + 'StartDate', label: 'Statement Period', jenis: 'tampil', sajian: TGL },
  { jalur: P + 'SOBName', label: 'Source Of Business', jenis: 'tampil' },
  { jalur: P + 'TreatyGroupName', label: 'Treaty Group', jenis: 'tampil', tampil: bukanNonPropBaru },
  // wadah `pyWorkPage.FlagViewPolicy = 1` - FlagViewPolicy tidak diisi rule
  // terjangkau mana pun (hanya dibaca `SumTSIPremiSpreadedRNM_Act`, tak terjangkau)
  ...dalamWadah((h) => nilai(h, 'FlagViewPolicy') === '1', [
    { jalur: P + 'PolicyNo', label: 'No Polis', jenis: 'tampil' },
    { jalur: P + 'ProductionDate', label: 'Production Date', jenis: 'tampil', sajian: TGL },
  ]),
  { jalur: P + 'DueTo', label: 'Due To Us / You', jenis: 'tampil' },
  { jalur: P + 'QuotationData.IsSurveyReport', label: 'Survey Report', jenis: 'tampil', tampil: bukanNonProp },
  { jalur: P + 'StatementType', label: 'Statement Type', jenis: 'tampil' },
  { jalur: P + 'FlagPPH', label: 'FlagPPH', jenis: 'centang', kunci: true },
  { jalur: P + 'TypeTax', label: 'Type Tax', jenis: 'tampil', tampil: tidakKosong(P + 'TypeTax') },
  { jalur: P + 'QuotationData.NoOfferSlip', label: 'No Offer Slip', jenis: 'tampil' },
  { jalur: P + 'ShareCurrency', label: 'RNM Share', jenis: 'tampil' },
  { jalur: P + 'ShareValue', label: 'ShareValue', jenis: 'tampil', mataUang: P + 'ShareCurrency', sajian: UANG },
  { jalur: P + 'StatementDate', label: 'Statement Date', jenis: 'tampil', sajian: TGL },
  { jalur: 'TreatyIn.Termination', label: 'Termination', jenis: 'tampil', sajian: TGL },
  { jalur: P + 'EndDate', label: 'To', jenis: 'tampil', sajian: TGL },
  { jalur: P + 'CedingCoName', label: 'Ceding Company', jenis: 'tampil', tampil: tidakKosong(P + 'CedingCoName') },
  { jalur: P + 'InsuredName', label: 'InsuredName', jenis: 'tampil' },
  { jalur: P + 'TreatyType', label: 'Treaty Type', jenis: 'tampil' },
  { jalur: P + 'TreatyYear', label: 'UW Year', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: P + 'Quartal', label: '.Quartal', jenis: 'tampil', tampil: proporsionalQD, sajian: BULAT_POLOS },
  { jalur: P + 'YearOfQuartal', label: 'YearOfQuartal', jenis: 'tampil', tampil: proporsionalQD, sajian: BULAT_POLOS },
  { jalur: P + 'TreatyYear', label: 'U/Y', jenis: 'tampil', tampil: proporsionalQD },
  {
    jalur: P + 'MarketingOfficer',
    label: 'Marketing Officer',
    jenis: 'tampil',
    tampil: tidakKosong(P + 'MarketingOfficer'),
  },
  { jalur: P + 'QuotationData.ProportionalType', label: 'Proportional Type', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: P + 'Currency', label: 'Currency', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: P + 'ClaimType', label: 'Claim Type', jenis: 'tampil' },
  { jalur: P + 'ClaimPaymentType', label: 'Claim Payment Type', jenis: 'tampil' },
  ...dalamWadah(nonProporsionalQD, [
    { jalur: P + 'LayerType', label: 'LayerType', jenis: 'tampil' },
    { jalur: P + 'Layer', label: 'Layer', jenis: 'tampil' },
    { jalur: P + 'LayerPartType', label: 'LayerPartType', jenis: 'tampil' },
    { jalur: P + 'LayerPart', label: 'LayerPart', jenis: 'tampil' },
  ]),
  { jalur: P + 'Remark', label: 'Remark', jenis: 'tampil' },
]

/** Bagian uang layar atasan - label VERBATIM `DetailDeptHeadTreatyIn_UW`
 *  (berbeda dari layar admin: tanpa "(%)", tanpa Claim 100%). Wadah
 *  `.IsNewPolicyNonProp != 1`. */
export const MEDAN_ATASAN_UANG: Medan[] = dalamWadah(wadahUangAtasan, medanAtasanUang())

function medanAtasanUang(): Medan[] {
  const t = (m: string, label: string, sajian: Sajian = UANG, uang = true): Medan => ({
    jalur: P + m,
    label,
    jenis: 'tampil',
    mataUang: uang ? mu : undefined,
    sajian,
  })
  return [
    t('GrossPremium', 'Gross Premium 100%', DUA),
    t('PremiOgp', 'Premi Ogp'),
    t('RiCommOgp', 'Deduction In A (OGP)', DUA, false),
    t('ResultOgp1', 'ResultOgp1'),
    t('OveriddingCommOgp', 'Deduction In B (OGP)', DUA, false),
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
    t('RiCommOnp', 'Deduction In A (ONP)', DUA, false),
    t('ResultOnp1', 'ResultOnp1'),
    t('OveriddingCommOnp', 'Deduction In B (ONP)', DUA, false),
    t('ResultOnp2', 'ResultOnp2'),
    // K3: pxCurrency, jumlah uang
    t('Deduction1', 'Deduction1'),
    t('Deduction2', 'Deduction2'),
    t('PPHValue', 'PPH 2%'),
    t('PPNValue', 'PPN 2.2%'),
  ]
}

/** Total berlabel di bawah grid spreading layar atasan (`.TotalSharePercentagePremium`
 *  dst., pxNumber 2 desimal). Di layar admin sel yang sama berwadah `1=2` (mati). */
export const MEDAN_ATASAN_TOTAL: Medan[] = dalamWadah(wadahUangAtasan, [
  { jalur: P + 'TotalSharePercentagePremium', label: 'Total %Share', jenis: 'tampil', sajian: DUA },
  { jalur: P + 'TotalPremium', label: 'Total Premium', jenis: 'tampil', mataUang: mu, sajian: DUA },
  { jalur: P + 'TotalSharePercentageClaim', label: 'Total %Share Claim', jenis: 'tampil', sajian: DUA },
  { jalur: P + 'TotalClaim', label: 'Total Claim', jenis: 'tampil', mataUang: mu, sajian: DUA },
])

// ------------------------------------------------------------------ grid (K14)

/** Sajian sel grid spreading `.SpreadingRiskList` kedua layar: pxNumber 4 desimal
 *  (persen tersunting `pyShowReadonlyFormatting=true`), total footer 4 desimal. */
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
