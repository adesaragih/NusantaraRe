// Definisi medan kedua layar realisasi - VERBATIM sel section Pega.
//
//   Admin   `Section/DetailPolicyTreatyIn`        (flow action InboxPolicyTreatyIn)
//   Atasan  `Section/DetailDeptHeadTreatyIn_UW`   (flow action DeptHeadTreatyIn_UW)
//
// Setiap medan membawa: jalur halaman, label, jenis kontrol, syarat tampil
// (`pyVisible` yang BERLAKU - INVENTARIS bab 5), terkunci, dan aksi refresh
// (`change -> refresh <Activity>`). Medan wajib TIDAK ditulis di sini: daftar
// wajib datang dari backend (`Layar.medanWajib`) supaya satu sumber.
//
// ⛔ Tidak ditampilkan, dan sebabnya:
//   - elemen bersyarat tampil `1=2` / `NEVER` (80 elemen mati, AC 53):
//     `.BizName`, `.DueTo` (admin), label pemisah, label bagian;
//   - `Select Source Of Business` (tampil hanya bila ClaimType 'XOL Retro') dan
//     `Choose Business R` (`SetValueRetro_Act` -> `InputPolicyTreatyOutDetail_preACT`):
//     jalur retro/treaty keluar - pembongkar JSON master (P29, AC 58, 62);
//   - tombol `Survey Report` (`HistoricalSurveyReport`): penyimpanan survei tidak
//     dirancang di tiket 00-23 - `[terbuka]`;
//   - subsection `DetailPoliciesNonProportional` / `DetailPolicyTreatyOutNonProportional`:
//     isinya halaman master JSON `TreatyIn.Limits/Share/...` (P29).

import { nilai, type Halaman } from './api'

export type JenisMedan = 'tampil' | 'teks' | 'angka' | 'tanggal' | 'area' | 'centang' | 'mataUang' | 'mo'

export interface Medan {
  jalur: string
  label: string
  jenis: JenisMedan
  /** `pyVisible` yang berlaku; tidak ada = selalu tampil. */
  tampil?: (h: Halaman) => boolean
  /** Terkunci (`pyReadOnly`) - selalu ditampilkan hanya-baca. */
  kunci?: boolean
  /** Refresh sesudah berubah: aksi backend + parameternya. */
  aksi?: { aksi: string; param?: string }
  /** Kode mata uang pasangan angka uang (AC 85). */
  mataUang?: string
}

const P = 'PolicyTreatyIn.'
const v = (h: Halaman, m: string) => nilai(h, P + m)

/** `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`. */
export const bukanNonProp = (h: Halaman) => nilai(h, 'Quotation.ProportionalType') !== 'NonProportional'
/** `.IsNewPolicyNonProp != 1`. */
export const bukanNonPropBaru = (h: Halaman) => v(h, 'IsNewPolicyNonProp') !== '1'
/** `.BalanceDueTo < 0` (teks kosong = 0, aritmetika Pega). */
export function saldoNegatif(h: Halaman): boolean {
  const s = v(h, 'BalanceDueTo').trim()
  return s !== '' && s.startsWith('-') && /[1-9]/.test(s)
}

const mu = P + 'Currency'

/** Medan umum layar admin - urutan sel `DetailPolicyTreatyIn`. */
export const MEDAN_ADMIN_UMUM: Medan[] = [
  { jalur: P + 'NoOffer', label: 'Master ID', jenis: 'tampil' },
  { jalur: 'TreatyIn.Commencement', label: 'Commencement', jenis: 'tampil' },
  // refresh `change` ber-pyPreDataTransform `SystemSetOneYear_DT` (EndDate = StartDate + 1 tahun)
  { jalur: P + 'StartDate', label: 'Statement Period', jenis: 'tanggal', aksi: { aksi: 'SystemSetOneYear' } },
  { jalur: P + 'EndDate', label: 'To', jenis: 'tanggal', aksi: { aksi: 'ProtectDate' } },
  { jalur: P + 'SOBName', label: 'Source Of Business', jenis: 'tampil' },
  { jalur: P + 'TreatyGroupName', label: 'Treaty Group', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: P + 'QuotationData.IsSurveyReport', label: 'Survey Report', jenis: 'teks', tampil: bukanNonProp },
  { jalur: P + 'StatementType', label: 'Statement Type', jenis: 'teks' },
  { jalur: P + 'QuotationData.NoOfferSlip', label: 'No Offer Slip', jenis: 'area' },
  {
    jalur: P + 'FlagRetroTreaty',
    label: 'FlagRetroTreaty',
    jenis: 'centang',
    tampil: (h) => v(h, 'ClaimType') === 'XOL Retro',
  },
  { jalur: P + 'FlagPPH', label: 'FlagPPH', jenis: 'centang', aksi: { aksi: 'RemoveTypeTax' } },
  { jalur: P + 'TypeTax', label: 'Type Tax', jenis: 'teks', tampil: (h) => v(h, 'FlagPPH') === 'true' },
  { jalur: P + 'ShareCurrency', label: 'RNM Share', jenis: 'tampil' },
  { jalur: P + 'ShareValue', label: 'ShareValue', jenis: 'tampil', mataUang: P + 'ShareCurrency' },
  { jalur: P + 'StatementDate', label: 'Statement Date', jenis: 'tampil' },
  { jalur: 'TreatyIn.Termination', label: 'Termination', jenis: 'tampil' },
  { jalur: P + 'CedingCoName', label: 'Ceding Company', jenis: 'tampil' },
  { jalur: P + 'InsuredName', label: 'InsuredName', jenis: 'tampil' },
  { jalur: P + 'TreatyType', label: 'Treaty Type', jenis: 'tampil' },
  { jalur: P + 'TreatyYear', label: 'UW Year', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: P + 'Quartal', label: '.Quartal', jenis: 'teks' },
  { jalur: P + 'YearOfQuartal', label: 'YearOfQuartal', jenis: 'teks' },
  { jalur: P + 'QuotationData.MOID', label: 'Marketing Officer', jenis: 'mo', aksi: { aksi: 'CheckDataMkt' } },
  { jalur: P + 'QuotationData.ProportionalType', label: 'Proportional Type', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: P + 'IDCurrency', label: 'Currency', jenis: 'mataUang', tampil: bukanNonPropBaru, aksi: { aksi: 'SetCurrency' } },
  { jalur: P + 'ClaimType', label: 'Claim Type', jenis: 'teks' },
  { jalur: P + 'ClaimPaymentType', label: 'Payment Type', jenis: 'teks' },
  { jalur: P + 'LayerType', label: 'LayerType', jenis: 'tampil' },
  { jalur: P + 'Layer', label: 'Layer', jenis: 'tampil' },
  { jalur: P + 'LayerPartType', label: 'LayerPartType', jenis: 'tampil' },
  { jalur: P + 'LayerPart', label: 'LayerPart', jenis: 'tampil' },
  { jalur: P + 'Remark', label: 'Remark', jenis: 'area' },
]

/** Medan uang layar admin - bagian "Old Soa Input Format". */
export const MEDAN_ADMIN_UANG: Medan[] = [
  { jalur: P + 'GrossPremium', label: 'Gross Premium 100%', jenis: 'angka', aksi: { aksi: 'CalculatePremi', param: 'PREMIUM' }, mataUang: mu },
  { jalur: P + 'GrossClaim', label: 'Claim 100%', jenis: 'angka', aksi: { aksi: 'CalculatePremi', param: 'CLAIM' }, mataUang: mu },
  { jalur: P + 'PremiOgp', label: 'Premi Ogp', jenis: 'angka', aksi: { aksi: 'CountOGPONP' }, mataUang: mu },
  { jalur: P + 'RiCommOgp', label: '(%) Deduction In A (OGP)', jenis: 'angka', aksi: { aksi: 'CountResult1', param: 'Pct' } },
  { jalur: P + 'ResultOgp1', label: 'Deduction In A (OGP)', jenis: 'angka', aksi: { aksi: 'CountResult1', param: 'Amount' }, mataUang: mu },
  { jalur: P + 'OveriddingCommOgp', label: '(%) Deduction In B (OGP)', jenis: 'angka', aksi: { aksi: 'CountResult2Ogp', param: 'Pct' } },
  { jalur: P + 'ResultOgp2', label: 'Deduction In B (OGP)', jenis: 'angka', aksi: { aksi: 'CountResult2Ogp', param: 'Amount' }, mataUang: mu },
  { jalur: P + 'PremiOnp', label: 'Premi Onp', jenis: 'angka', aksi: { aksi: 'CountOGPONP' }, mataUang: mu },
  { jalur: P + 'RiCommOnp', label: '(%) Deduction In A (ONP)', jenis: 'angka', aksi: { aksi: 'CountResult1Onp', param: 'Pct' } },
  { jalur: P + 'ResultOnp1', label: 'Deduction In A (ONP)', jenis: 'angka', aksi: { aksi: 'CountResult1Onp', param: 'Amount' }, mataUang: mu },
  { jalur: P + 'OveriddingCommOnp', label: '(%) Deduction In B (ONP)', jenis: 'angka', aksi: { aksi: 'CountResult2Onp', param: 'Pct' } },
  { jalur: P + 'ResultOnp2', label: 'Deduction In B (ONP)', jenis: 'angka', aksi: { aksi: 'CountResult2Onp', param: 'Amount' }, mataUang: mu },
  { jalur: P + 'Claim', label: 'Claim', jenis: 'angka', aksi: { aksi: 'CountOGPONP' }, mataUang: mu },
  { jalur: P + 'OutstandingClaim', label: 'Outstanding Claim', jenis: 'angka', mataUang: mu },
  { jalur: P + 'SalvageValue', label: 'Salvage', jenis: 'angka', aksi: { aksi: 'CountOGPONP' }, mataUang: mu },
  { jalur: P + 'ExcessLoss', label: 'Excess Loss', jenis: 'angka', aksi: { aksi: 'CountOGPONP' }, mataUang: mu },
  { jalur: P + 'Deduction1', label: 'Deduction1', jenis: 'angka', aksi: { aksi: 'CountOGPONP' } }, // persen (P29)
  { jalur: P + 'Deduction2', label: 'Deduction2', jenis: 'angka', aksi: { aksi: 'CountOGPONP' } }, // persen (P29)
  { jalur: P + 'NetPremium', label: 'Total Premium Before Claim', jenis: 'tampil', mataUang: mu },
  { jalur: P + 'BalanceDueTo', label: 'Balance Due To You', jenis: 'tampil', tampil: saldoNegatif, mataUang: mu },
  { jalur: P + 'BalanceBeforeTax', label: 'Balance Before Tax', jenis: 'tampil', tampil: (h) => !saldoNegatif(h), mataUang: mu },
  {
    jalur: P + 'BalanceBeforePPH',
    label: 'Balance Before Withholding Tax (PPH 2.2)',
    jenis: 'tampil',
    tampil: (h) => !saldoNegatif(h),
    mataUang: mu,
  },
  { jalur: P + 'BalanceDueTo', label: 'Balance Due To Us', jenis: 'tampil', tampil: (h) => !saldoNegatif(h), mataUang: mu },
  { jalur: P + 'PPHValue', label: 'PPH 2%', jenis: 'tampil', mataUang: mu },
  { jalur: P + 'PPNValue', label: 'PPN 2.2%', jenis: 'tampil', mataUang: mu },
]

/** Medan layar atasan (`DetailDeptHeadTreatyIn_UW`) - seluruhnya terkunci
 *  kecuali FlagPPH dan No Offer Slip (kolom "Kunci" kosong). */
export const MEDAN_ATASAN_UMUM: Medan[] = [
  { jalur: P + 'NoOffer', label: 'Master ID', jenis: 'tampil' },
  { jalur: 'TreatyIn.Commencement', label: 'Commencement', jenis: 'tampil' },
  { jalur: P + 'StartDate', label: 'Statement Period', jenis: 'tampil' },
  { jalur: P + 'EndDate', label: 'To', jenis: 'tampil' },
  { jalur: P + 'SOBName', label: 'Source Of Business', jenis: 'tampil' },
  { jalur: P + 'TreatyGroupName', label: 'Treaty Group', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: P + 'PolicyNo', label: 'No Polis', jenis: 'tampil' },
  { jalur: P + 'ProductionDate', label: 'Production Date', jenis: 'tampil' },
  { jalur: P + 'DueTo', label: 'Due To Us / You', jenis: 'tampil' },
  { jalur: P + 'QuotationData.IsSurveyReport', label: 'Survey Report', jenis: 'tampil', tampil: bukanNonProp },
  { jalur: P + 'StatementType', label: 'Statement Type', jenis: 'tampil' },
  { jalur: P + 'FlagPPH', label: 'FlagPPH', jenis: 'centang' },
  { jalur: P + 'TypeTax', label: 'Type Tax', jenis: 'tampil', tampil: (h) => v(h, 'TypeTax') !== '' },
  { jalur: P + 'QuotationData.NoOfferSlip', label: 'No Offer Slip', jenis: 'area' },
  { jalur: P + 'ShareCurrency', label: 'RNM Share', jenis: 'tampil' },
  { jalur: P + 'ShareValue', label: 'ShareValue', jenis: 'tampil', mataUang: P + 'ShareCurrency' },
  { jalur: P + 'StatementDate', label: 'Statement Date', jenis: 'tampil' },
  { jalur: 'TreatyIn.Termination', label: 'Termination', jenis: 'tampil' },
  { jalur: P + 'CedingCoName', label: 'Ceding Company', jenis: 'tampil' },
  { jalur: P + 'InsuredName', label: 'InsuredName', jenis: 'tampil' },
  { jalur: P + 'TreatyType', label: 'Treaty Type', jenis: 'tampil' },
  { jalur: P + 'TreatyYear', label: 'UW Year', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: P + 'Quartal', label: '.Quartal', jenis: 'tampil' },
  { jalur: P + 'YearOfQuartal', label: 'YearOfQuartal', jenis: 'tampil' },
  { jalur: P + 'MarketingOfficer', label: 'Marketing Officer', jenis: 'tampil' },
  { jalur: P + 'QuotationData.ProportionalType', label: 'Proportional Type', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: P + 'Currency', label: 'Currency', jenis: 'tampil', tampil: bukanNonPropBaru },
  { jalur: P + 'ClaimType', label: 'Claim Type', jenis: 'tampil' },
  { jalur: P + 'ClaimPaymentType', label: 'Claim Payment Type', jenis: 'tampil' },
  { jalur: P + 'LayerType', label: 'LayerType', jenis: 'tampil' },
  { jalur: P + 'Layer', label: 'Layer', jenis: 'tampil' },
  { jalur: P + 'LayerPartType', label: 'LayerPartType', jenis: 'tampil' },
  { jalur: P + 'LayerPart', label: 'LayerPart', jenis: 'tampil' },
  { jalur: P + 'Remark', label: 'Remark', jenis: 'tampil' },
]

/** Bagian uang layar atasan. */
export const MEDAN_ATASAN_UANG: Medan[] = medanAtasanUang()

/** Bagian uang layar atasan - label VERBATIM `DetailDeptHeadTreatyIn_UW`
 *  (berbeda dari layar admin: tanpa "(%)", tanpa Claim 100%). */
function medanAtasanUang(): Medan[] {
  const t = (m: string, label: string, uang = true): Medan => ({ jalur: P + m, label, jenis: 'tampil', mataUang: uang ? mu : undefined })
  return [
    t('GrossPremium', 'Gross Premium 100%'),
    t('PremiOgp', 'Premi Ogp'),
    t('RiCommOgp', 'Deduction In A (OGP)', false),
    t('ResultOgp1', 'ResultOgp1'),
    t('OveriddingCommOgp', 'Deduction In B (OGP)', false),
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
    t('RiCommOnp', 'Deduction In A (ONP)', false),
    t('ResultOnp1', 'ResultOnp1'),
    t('OveriddingCommOnp', 'Deduction In B (ONP)', false),
    t('ResultOnp2', 'ResultOnp2'),
    t('Deduction1', 'Deduction1', false), // persen (P29)
    t('Deduction2', 'Deduction2', false),
    t('PPHValue', 'PPH 2%'),
    t('PPNValue', 'PPN 2.2%'),
  ]
}

/** Medan tampil bagi halaman ini. */
export function medanTampil(daftar: Medan[], h: Halaman): Medan[] {
  return daftar.filter((m) => !m.tampil || m.tampil(h))
}
