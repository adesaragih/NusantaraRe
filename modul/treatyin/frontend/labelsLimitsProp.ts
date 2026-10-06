// TAB LIMITS CABANG PROPORSIONAL — TIGA TINGKAT, dibaca dari ekspor
// `D:\XML_NURE\Treaty In\Section`.
//
//   tingkat 1  `TreatyInTabsProportional.xml` @524276 — grid `Kind of Treaty`
//              (`.TreatyType` W1315), `pyEditingMode=expandPane`,
//              `pyEditAction=LimitProportional` @604564
//   tingkat 2  `LimitProportional.xml` (kelas TreatyInLimits) — `Treaty Type`
//              + grid `.Detail` (Treaty Group), expandPane,
//              `pyEditAction=DetailLimits` @87354
//   tingkat 3  `DetailLimits.xml` (kelas TreatyInLimitsDetail) — medan,
//              empat grid, sebelas tab `TABBED`
//
// Offset `@n` = posisi bita sesudah `<pyIncludedRuleXML>` dibuang dengan
// MENGHITUNG SARANG. Kolom tanpa kepala memang tanpa kepala di ekspor.

import type { GolonganAngka } from './labels'

/** Satu kolom grid daftar Limits. `desimal` = `pyDecimalPlaces` ekspor; null = tidak dinyatakan. */
export interface KolomLimit {
  label: string
  kunci: string
  lebar: number
  golongan: GolonganAngka
  desimal: number | null
  /** Syarat tampil SEL — dinilai per baris. */
  syarat?: (baris: Readonly<Record<string, unknown>>) => boolean
  /**
   * Sel dropdown mata uang (CURRENCY). `nama`: sel mengikat `.Currency`
   * langsung (`BrowseCurrencyTreatyIn_RD`, nilai `.Currency`). `id`: sel
   * mengikat `.CurrencyID` (`BrowseCurrency_RD`, nilai `.ID`) dan
   * `SetCurrName_Act` mengisi `.Currency` dengan namanya.
   */
  mataUang?: 'nama' | 'id'
}

const JENIS_SURPLUS = ['SURPLUS', '2ND SURPLUS', '3RD SURPLUS', 'SPECIAL SURPLUS']

/** Grid `Currency · Value` yang berulang di DetailLimits — kepalanya kosong di ekspor. */
function kolomNilai(lebarMataUang: number, lebarNilai: number): readonly KolomLimit[] {
  return [
    { label: '', kunci: 'Currency', lebar: lebarMataUang, golongan: 'teks', desimal: null, mataUang: 'nama' },
    { label: '', kunci: 'Value', lebar: lebarNilai, golongan: 'uang', desimal: 2 },
  ]
}

export const LIMITS_PROP = {
  judul: 'Limits', // pyTitle @534799
  treatyType: 'Treaty Type', // LimitProportional @17821, `.TreatyTypeID` pxDropdown
  treatyGroup: 'Treaty Group', // DetailLimits @33452, `.TreatyGroupID` pxDropdown
  tambah: 'Add', // tombol kepala grid — `TreatyIn.ViewState !='1'`
  hapus: 'Delete', // baris Kind of Treaty / Treaty Group — syarat sama
  hapusBaris: 'Remove', // baris grid DetailLimits (100% Limit, Retention, …)
  /**
   * `.TreatyTypeID` adalah KODE dropdown (`10042` 13 kali, `10035` 5,
   * `10037` 1 — terisi 19 dari 1.360). Pilihannya dari `REINSURANCETYPE`
   * (`BrowseReinsuranceType_RD`): nilai `.ID`, label `.Note`.
   */
  pilihKosong: 'Choose', // pyNoSelectionText LimitProportional
  namaKembar: 'nama dipakai lebih dari satu kode',
  qs: 'QS %', // @185962 `.QSPct` d2 — bila `.TreatyType = 'QUOTA SHARE'` @170234
  lines: 'Lines', // @228088 `.Surplus` d2 — bila jenis surplus @212267
  limit100: '100% Limit', // @283497
  seratus: '100', // @297544 — bila QUOTA SHARE
  persen: '%',
  retensi: 'Retention', // @460405; `.RetentionPct` @474451 bila QUOTA SHARE
  cession: 'Cession to R/I', // @642766; `.CessionPct` @656819 bila QUOTA SHARE
  tanpaBaris: 'No items',
  tidakAda: 'Tidak ada di dokumen sistem lama',
} as const

/** `.TreatyType` QUOTA SHARE / jenis surplus — syarat medan DetailLimits. */
export const syaratQS = (jenis: string) => jenis === 'QUOTA SHARE'
export const syaratSurplus = (jenis: string) => JENIS_SURPLUS.includes(jenis)

/** Kolom grid tingkat 1 dan 2. */
export const KOLOM_KIND_OF_TREATY: readonly KolomLimit[] = [
  { label: 'Kind of Treaty', kunci: 'TreatyType', lebar: 1315, golongan: 'teks', desimal: null },
]
export const KOLOM_TREATY_GROUP: readonly KolomLimit[] = [
  { label: 'Treaty Group', kunci: 'TreatyGroup', lebar: 1368, golongan: 'teks', desimal: null },
]

/** Grid DetailLimits — kolom, lebar, desimal, syarat sel, dari ekspor. */
export const GRID_DETAIL: Readonly<Record<string, { kolom: readonly KolomLimit[]; tambah: boolean; jejak: string }>> = {
  COBList: {
    kolom: [{ label: 'Class of Business', kunci: 'ClassOfBusiness', lebar: 272, golongan: 'teks', desimal: null }],
    tambah: false, // tombolnya `TreatyIn.ViewState !='1' && 1=2` — MATI di ekspor
    jejak: 'DetailLimits @86887',
  },
  IOOLimitList: {
    kolom: [
      ...kolomNilai(194, 361),
      // `.Layer` tampil bila `.Note = 'QUOTA SHARE'` — baris itu sendiri.
      { label: '', kunci: 'Layer', lebar: 128, golongan: 'teks', desimal: null, syarat: (b) => b.Note === 'QUOTA SHARE' },
    ],
    tambah: true,
    jejak: 'DetailLimits @345241',
  },
  RetentionList: {
    kolom: [
      ...kolomNilai(196, 363),
      { label: '', kunci: 'Layer', lebar: 128, golongan: 'teks', desimal: null, syarat: (b) => JENIS_SURPLUS.includes(String(b.Note ?? '')) },
    ],
    tambah: true,
    jejak: 'DetailLimits @524406',
  },
  CessionList: { kolom: kolomNilai(193, 352), tambah: true, jejak: 'DetailLimits @706791' },
  DeductionList: {
    kolom: [
      { label: 'Description', kunci: 'Comment', lebar: 206, golongan: 'teks', desimal: null },
      // Sel mengikat `.CurrencyID` (dropdown); yang tampil namanya, `.Currency`.
      { label: 'Currency', kunci: 'Currency', lebar: 197, golongan: 'teks', desimal: null, mataUang: 'id' },
      { label: 'Deduction', kunci: 'Deduction', lebar: 201, golongan: 'uang', desimal: 2 },
      { label: 'or', kunci: '', lebar: 31, golongan: 'teks', desimal: null },
      { label: 'Deduction %', kunci: 'DeductionPct', lebar: 213, golongan: 'persen', desimal: 2 },
    ],
    tambah: true,
    jejak: 'DetailLimits @1127846',
  },
  DeductionTotalList: {
    kolom: [
      { label: 'Deductions', kunci: 'Currency', lebar: 196, golongan: 'teks', desimal: null },
      { label: '', kunci: 'Value', lebar: 356, golongan: 'uang', desimal: 2 },
    ],
    tambah: false, // tanpa kolom tombol di ekspor
    jejak: 'DetailLimits @1259471',
  },
  ReserveList: { kolom: kolomNilai(198, 363), tambah: true, jejak: 'DetailLimits @1427153' },
  PLAList: { kolom: kolomNilai(197, 361), tambah: true, jejak: 'DetailLimits @1673566' },
  CashLossList: { kolom: kolomNilai(197, 361), tambah: true, jejak: 'DetailLimits @1851745' },
  ClaimCoopList: { kolom: kolomNilai(197, 361), tambah: true, jejak: 'DetailLimits @2029946' },
  EPIList: { kolom: kolomNilai(197, 361), tambah: true, jejak: 'DetailLimits @2284378' },
  AchievementLists: {
    kolom: [
      { label: 'Quarter', kunci: 'Quarter', lebar: 83, golongan: 'teks', desimal: null },
      { label: 'Quarter Year', kunci: 'QUARTERYEAR', lebar: 75, golongan: 'teks', desimal: null },
      { label: 'Currency', kunci: 'Currency', lebar: 93, golongan: 'teks', desimal: null },
      { label: 'Premium', kunci: 'PREMIUM', lebar: 138, golongan: 'uang', desimal: 2 },
      { label: 'R/I COMM', kunci: 'RICOMM', lebar: 140, golongan: 'uang', desimal: 2 },
      { label: 'Brokerage', kunci: 'BROKERAGE', lebar: 159, golongan: 'uang', desimal: 2 },
      { label: 'Net Premium Before Claim', kunci: 'NETPREMIUM', lebar: 124, golongan: 'uang', desimal: 2 },
      { label: 'Paid Claim', kunci: 'PaidClaim', lebar: 131, golongan: 'uang', desimal: 2 },
      { label: 'Cash Call Claim', kunci: 'CASHCALL', lebar: 130, golongan: 'uang', desimal: 2 },
      { label: 'Outstanding Claim', kunci: 'OutstandingClaim', lebar: 163, golongan: 'uang', desimal: 2 },
      // Tiga kolom terakhir TIDAK menyatakan desimal — aturan lama.
      { label: 'Incured Claim', kunci: 'IncuredClaim', lebar: 161, golongan: 'uang', desimal: null },
      { label: 'Total', kunci: 'Total', lebar: 113, golongan: 'uang', desimal: null },
      { label: 'Net Loss Ratio', kunci: 'LossRatio', lebar: 73, golongan: 'persen', desimal: null },
    ],
    tambah: false, // tanpa kolom tombol di ekspor
    jejak: 'DetailLimits @2413255',
  },
}

/** Satu medan DetailLimits. `desimal`: ekspor; `null` tidak dinyatakan; `-1` tak dibatasi. */
export interface MedanLimit {
  label: string
  kunci: string
  golongan: GolonganAngka
  desimal: number | null
  /** Pasangan mata uang di sebelah kiri (Event Limits). */
  mataUang?: string
}

/** Butir satu tab DetailLimits: medan, grid (kunci `GRID_DETAIL`), atau teks label ekspor. */
export type ButirTabLimit =
  | { t: 'medan'; m: MedanLimit }
  | { t: 'grid'; larik: string; judul?: string }
  | { t: 'teks'; teks: string }

/** Sebelas tab DetailLimits — `pyHeaderType=TABBED`, urutan ekspor. */
export const TAB_DETAIL_LIMIT: readonly { judul: string; at: number; isi: readonly ButirTabLimit[] }[] = [
  {
    judul: 'Event Limits',
    at: 788648,
    isi: [
      { t: 'medan', m: { label: 'RSMD Limit', mataUang: 'CurrencyRSMD', kunci: 'RSMDLimit', golongan: 'uang', desimal: null } },
      { t: 'medan', m: { label: 'Earthquake Limit', mataUang: 'CurrencyEarthquake', kunci: 'Earthquake', golongan: 'uang', desimal: null } },
      { t: 'medan', m: { label: 'Flood Limit (Jabodetabek)', mataUang: 'CurrencyFloodJab', kunci: 'FloodJab', golongan: 'uang', desimal: null } },
      { t: 'medan', m: { label: 'Flood Limit (Nationwide)', mataUang: 'CurrencyFloodNat', kunci: 'FloodNation', golongan: 'uang', desimal: null } },
    ],
  },
  {
    judul: 'Deduction In A',
    at: 1054344,
    isi: [
      { t: 'medan', m: { label: '% OGR', kunci: 'RIOGR', golongan: 'persen', desimal: null } },
      { t: 'medan', m: { label: '% ONR', kunci: 'RIONR', golongan: 'persen', desimal: null } },
    ],
  },
  {
    judul: 'Deduction',
    at: 1100591,
    isi: [
      { t: 'grid', larik: 'DeductionList', judul: 'Deduction Details' },
      { t: 'grid', larik: 'DeductionTotalList' },
    ],
  },
  {
    judul: 'Reserve',
    at: 1311728,
    isi: [
      { t: 'medan', m: { label: '% Premium Reserve', kunci: 'PremiumReservePct', golongan: 'persen', desimal: null } },
      { t: 'teks', teks: 'Premium Reserve' },
      { t: 'grid', larik: 'ReserveList' },
    ],
  },
  {
    judul: 'Experience Premium Refund',
    at: 1515314,
    isi: [
      { t: 'medan', m: { label: '% Commision', kunci: 'ProfitCommision', golongan: 'persen', desimal: null } },
      { t: 'medan', m: { label: '% ME', kunci: 'ProfitME', golongan: 'persen', desimal: null } },
      // `pyDecimalPlaces = -999` @1549354 — padanan "tak dibatasi" Pega.
      { t: 'medan', m: { label: 'YDCF', kunci: 'ProfitYDCF', golongan: 'uang', desimal: -1 } },
      { t: 'teks', teks: 'Note:' },
      { t: 'teks', teks: 'ME = Management Expense' },
      { t: 'teks', teks: 'YDCF = Years Deficit Carried Forward' },
    ],
  },
  { judul: 'PLA', at: 1583182, isi: [{ t: 'teks', teks: 'PLA' }, { t: 'grid', larik: 'PLAList' }] },
  { judul: 'Cash Loss Limit', at: 1761325, isi: [{ t: 'teks', teks: 'Cash Loss Limit' }, { t: 'grid', larik: 'CashLossList' }] },
  { judul: 'Claim Cooperation', at: 1939519, isi: [{ t: 'teks', teks: 'Claim Cooperation' }, { t: 'grid', larik: 'ClaimCoopList' }] },
  {
    judul: 'LPC',
    at: 2117723,
    isi: [
      { t: 'teks', teks: 'Loss Participation Clause (LPC)' },
      { t: 'medan', m: { label: 'Lower Band', kunci: 'LowerBand', golongan: 'uang', desimal: null } },
      { t: 'medan', m: { label: 'Upper Band', kunci: 'UpperBand', golongan: 'uang', desimal: null } },
      { t: 'medan', m: { label: 'Reisured Participant', kunci: 'ReisuredParticipant', golongan: 'uang', desimal: null } },
      { t: 'medan', m: { label: 'Period (Month)', kunci: 'Periode', golongan: 'teks', desimal: null } },
    ],
  },
  { judul: 'EPI', at: 2194135, isi: [{ t: 'teks', teks: 'EPI' }, { t: 'grid', larik: 'EPIList' }] },
  // Achievement: grid saja. `As At Quarter`/`Quarter Year` mengikat halaman
  // SESI `SearchData`, `Refresh` menjalankan Activity — tidak tersimpan, tidak
  // dibangun. Tiga blok `1=2` di sekitarnya mati.
  { judul: 'Achievement', at: 2372416, isi: [{ t: 'grid', larik: 'AchievementLists' }] },
]
