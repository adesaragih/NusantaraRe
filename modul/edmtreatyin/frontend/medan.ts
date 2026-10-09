// Definisi medan layar endorsemen - VERBATIM sel section Pega korpus `EDM Treaty In`. Asal pola:
// `modul/nbtreatyin/frontend/medan.ts` (06-10-2026: tipe `Medan`, `dalamWadah`, sajian bagian uang 4 desimal).
//
//   Header   `Section/DetailPolicyTreatyInAddendum` S3 (kolom kiri) / S7 (kolom kanan) / S10 (Remark)
//   Tab      `DetailPolicyTreatyInPropOldData` / `...OldData2` / `...NewData` / `...NewData2` / `...ValueDifference`
//            (layout group Tab `DetailPolicyTreatyInAddGeneralEditable` S2108) - SATU kerangka, beda awalan
//            properti, label, dan status sunting (`tataTab`)
//
// Setiap medan membawa jalur halaman, label, jenis kontrol, syarat tampil (`pyVisible` sel DAN
// `pyContainerVisibleWhen` wadahnya), terkunci, dan action set BERURUTAN. Medan WAJIB datang dari backend
// (`Layar.medanWajib`). ⛔ Nol perhitungan: semua refresh berhitung dikirim ke `POST .../hitung`.
//
// ⛔ Tidak ditampilkan (bukti XML): "No Statement" (`Never`), "Due To Us / You" (`NEVER`), "Survey Report" (`1=2`),
// `.InsuredName` (`1=2`), "Claim Type" / "Claim Payment Type" (`NEVER`), label pemisah `1=2`, wadah total berlabel
// `1=2` (OldData/OldData2/ValueDifference/NewData2), `.Installment` PropNewData (tampil `PN = ADM` di section yang
// hanya tampil `PN != ADM`), tombol Delete grid spreading (keputusan work owner: baris endorsemen tidak dapat
// dihapus, spec-penyimpanan ID-16).

import { POLIS, nilai, type Baris, type Halaman, type Pilihan } from './api'
import { BAGIAN, PILIHAN_TYPE_TAX } from './labels'
import type { Sajian } from './sajian'
import { negatifTeks } from './tanda'

export type JenisMedan = 'tampil' | 'angka' | 'area' | 'centang' | 'radio' | 'mo'

/** Satu pilihan dropdown / radio: nilai standar disimpan, teks prompt ditampilkan. */
export interface OpsiMedan {
  value: string
  label: string
}

/** Daftar acuan yang memberi teks nilai dropdown hanya-baca (`GET /acuan`). */
export type SumberAcuan = 'mataUang' | 'mo' | 'jenisEdm'

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
  /** Terkunci (`pyReadOnly`, `pyDisabled`) - selalu hanya-baca. */
  kunci?: boolean
  /** Terbuka HANYA di layar Admin - sel header yang XML tidak kunci per posisi, keputusan work owner 07-10-2026:
   *  atasan tidak boleh mengubahnya (backend `medanHeader` hanya diterima dari Admin). */
  hanyaAdmin?: boolean
  /** Refresh sesudah berubah - BERURUTAN seperti action set sel XML. */
  aksi?: Aksi[]
  /** Penyajian nilai: format angka sel, atau tanggal. */
  sajian?: Sajian
  /** Pilihan pxRadioButtons ber-pyListSource `associated` (prompt values property). */
  opsi?: OpsiMedan[]
  /** Dropdown hanya-baca ber-RD / pageList: teks nilai dari daftar acuan ini. */
  acuan?: SumberAcuan
  /** pyMaxLength isian teks. */
  panjangMaks?: number
}

/** Teks prompt sebuah nilai standar; nilai di luar daftar tampil apa adanya (tidak dibuang). */
export function teksPilihan(opsi: OpsiMedan[], v: string): string {
  return opsi.find((o) => o.value === v)?.label ?? v
}

const v = (h: Halaman, m: string) => nilai(h, POLIS + m)

// ------------------------------------------------------------------ syarat

/** `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`. */
export const bukanNonPropQ = (h: Halaman) => nilai(h, 'Quotation.ProportionalType') !== 'NonProportional'
/** `.QuotationData.ProportionalType = 'NonProportional'` (UW Year). */
const nonPropQD = (h: Halaman) => v(h, 'QuotationData.ProportionalType') === 'NonProportional'
/** `pyWorkPage.PolicyTreatyIn.IsNewPolicyNonProp != '1'` (Currency). */
export const bukanNonPropBaru = (h: Halaman) => v(h, 'IsNewPolicyNonProp') !== '1'
/** S11 `.IsNewPolicyNonProp = 1` - cabang AddPremi. */
export const nonPropBaru = (h: Halaman) => v(h, 'IsNewPolicyNonProp') === '1'
/** S17 `.IsNewPolicyNonProp = 0 || .IsNewPolicyNonProp = ''` (wadah `DetailPolicyTreatyInAddendum`); syarat S22
 *  `DetailPolicyTreatyInAddGeneral` `.IsNewPolicyNonProp != 1` sudah tercakup - tab Old / New / Diff. */
export const tampilTabData = (h: Halaman) => {
  const x = v(h, 'IsNewPolicyNonProp')
  return x === '0' || x === ''
}
/** S9 `QuotationData.ProportionalType = 'NonProportional' && IsNewPolicyNonProp != 1` - baris Layer. */
const wadahLayer = (h: Halaman) => nonPropQD(h) && bukanNonPropBaru(h)
/** `.FlagPPH = true` - Type Tax. */
const denganPajak = (h: Halaman) => v(h, 'FlagPPH') === 'true'
/** `.OldData.EDMNo = ''` - tab Old Data = `PropOldData`; selain itu `PropOldData2`. */
export const oldDataTanpaEDM = (h: Halaman) => v(h, 'OldData.EDMNo') === ''

/** Medan dalam wadah bersyarat: syarat sel DAN syarat wadah. */
function dalamWadah(wadah: (h: Halaman) => boolean, ms: Medan[]): Medan[] {
  return ms.map((m) => {
    const sel = m.tampil
    return { ...m, tampil: sel ? (h: Halaman) => wadah(h) && sel(h) : wadah }
  })
}

/** Medan tampil bagi halaman ini. */
export function medanTampil(daftar: Medan[], h: Halaman): Medan[] {
  return daftar.filter((m) => !m.tampil || m.tampil(h))
}

// ------------------------------------------------------------------ sajian

/** Bagian uang - perintah work owner 06-10-2026 (NB): nol / kosong tampil "0", angka terisi 4 desimal. */
const UANG4: Sajian = { desimal: 4, nolPolos: true }
const UANG4_SUNTING: Sajian = { desimal: 4, nolPolos: true, formatSaatSunting: true }
/** pxNumber `pyDecimalPlaces` 2 (total berlabel PropNewData). */
const DUA: Sajian = { desimal: 2 }
/** `.Quartal` `/decnone`: tanpa desimal, tanpa pemisah ribuan (pola NB). */
const BULAT_POLOS: Sajian = { desimal: 0, ribuan: false }
const TGL: Sajian = 'tanggal'

// ------------------------------------------------------------------ header (DetailPolicyTreatyInAddendum)

/** Kolom kiri (wadah S3 + S4 + S5). Tombol "Choose Business" (S6) bukan medan - `pages/LayarKasus.tsx`. */
export const MEDAN_KIRI: Medan[] = [
  { jalur: POLIS + 'OldData.NoOffer', label: 'Old Master ID', jenis: 'tampil' },
  { jalur: POLIS + 'NoOffer', label: 'Current Master ID', jenis: 'tampil' },
  { jalur: POLIS + 'TreatyGroupName', label: 'Treaty Group', jenis: 'tampil', tampil: bukanNonPropQ },
  // ⛔ Keputusan work owner 07-10-2026 ("Class Of Business GANTI JADII Source Of Business, SAMAIN DENGAN NB NYA!"): sel
  // XML `.BizName` "Class Of Business" (tampil `Quotation.ProportionalType != 'NonProportional'`) diganti `.SOBName`
  // tanpa syarat seperti NB `MEDAN_ADMIN_UMUM`; sel `.SOBName` wadah S7 kolom kanan dipindah ke sini (tidak dobel)
  { jalur: POLIS + 'SOBName', label: 'Source Of Business', jenis: 'tampil' },
  // S4: `pyWorkPage.PolicyTreatyIn.PolicyNo` / `.EDMNo`
  { jalur: POLIS + 'PolicyNo', label: 'No Polis', jenis: 'tampil' },
  { jalur: POLIS + 'EDMNo', label: 'No EDM', jenis: 'tampil' },
  // pxDropdown Read-only; prompt values `associated` tidak ada di korpus - teks dari acuan jenisEdm (DT
  // TreatyEDMListType, screenshot work owner 07-10-2026)
  { jalur: POLIS + 'EDMType', label: 'EDM Type', jenis: 'tampil', acuan: 'jenisEdm' },
  // pyCheckboxCaption (label sel = "Checkbox"); Auto, click -> postValue. Hanya Admin (WO 07-10-2026).
  { jalur: POLIS + 'FlagRetroTreaty', label: 'Overiding Commision', jenis: 'centang', hanyaAdmin: true },
  // S5: click -> postValue -> runActivity RemoveTypeTax_ACT. Hanya Admin (WO 07-10-2026).
  {
    jalur: POLIS + 'FlagPPH',
    label: 'With Tax',
    jenis: 'centang',
    aksi: [{ aksi: 'RemoveTypeTax' }],
    hanyaAdmin: true,
  },
  // Auto, wajib bila `.FlagPPH = true` (medanWajib backend). XML hanya postValue; pajak dihitung ulang seperti
  // sel NB yang sama (keputusan work owner 06-10-2026, aksi `HitungPajak` backend)
  {
    jalur: POLIS + 'TypeTax',
    label: 'Type Tax',
    jenis: 'radio',
    opsi: PILIHAN_TYPE_TAX,
    tampil: denganPajak,
    aksi: [{ aksi: 'HitungPajak' }],
    hanyaAdmin: true,
  },
]

/** Kolom kanan (wadah S7, S8 deret Q, S9 deret Layer). */
export const MEDAN_KANAN: Medan[] = [
  // RO{ALWAYS || ...} + disabled
  { jalur: POLIS + 'StatementDate', label: 'Statement Date', jenis: 'tampil', sajian: TGL },
  { jalur: POLIS + 'StartDate', label: 'Statement Period', jenis: 'tampil', sajian: TGL },
  { jalur: POLIS + 'EndDate', label: 'To', jenis: 'tampil', sajian: TGL },
  { jalur: POLIS + 'CedingCoName', label: 'Ceding Company', jenis: 'tampil' },
  { jalur: POLIS + 'TreatyType', label: 'Treaty Type', jenis: 'tampil' },
  { jalur: POLIS + 'TreatyYear', label: 'UW Year', jenis: 'tampil', tampil: nonPropQD },
  // S8 `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`: LABEL "Q" [Quartal] "/" [YearOfQuartal] "U/Y" [TreatyYear]
  ...dalamWadah(bukanNonPropQ, [
    { jalur: POLIS + 'Quartal', label: 'Q', jenis: 'tampil', sajian: BULAT_POLOS },
    { jalur: POLIS + 'YearOfQuartal', label: '/', jenis: 'tampil' },
    { jalur: POLIS + 'TreatyYear', label: 'U/Y', jenis: 'tampil' },
  ]),
  // pxDropdown RD BrowseMarketingOfficer_RD, Auto: change -> postValue -> refresh CheckDataMkt. Hanya Admin (WO
  // 07-10-2026).
  {
    jalur: POLIS + 'QuotationData.MOID',
    label: 'Marketing Officer',
    jenis: 'mo',
    aksi: [{ aksi: 'CheckDataMkt' }],
    hanyaAdmin: true,
  },
  { jalur: POLIS + 'QuotationData.ProportionalType', label: 'Proportional Type', jenis: 'tampil' },
  // pxDropdown RD BrowseCurrencyTreatyIn_RD, RO{ALWAYS || ...} + disabled
  { jalur: POLIS + 'IDCurrency', label: 'Currency', jenis: 'tampil', acuan: 'mataUang', tampil: bukanNonPropBaru },
  // S9: [LayerType] [Layer] "Of" [LayerPartType] [LayerPart] - sel tanpa label
  ...dalamWadah(wadahLayer, [
    { jalur: POLIS + 'LayerType', label: 'LayerType', jenis: 'tampil' },
    { jalur: POLIS + 'Layer', label: 'Layer', jenis: 'tampil' },
    { jalur: POLIS + 'LayerPartType', label: 'Of', jenis: 'tampil' },
    { jalur: POLIS + 'LayerPart', label: 'LayerPart', jenis: 'tampil' },
  ]),
]

/** S10 `.Remark` pxTextArea (label "Text Area") - disabled bila `pyWorkPage.PositionNote != 'ReasTreatyInAdmin'`. */
export function medanRemark(admin: boolean): Medan {
  return { jalur: POLIS + 'Remark', label: 'Remark', jenis: 'area', kunci: !admin }
}

/** S12 `.Installment` cabang NonProp baru (label "Text Input", pyMaxLength 2): change -> FillPaymentInstallmentEDMT. */
export const MEDAN_ANGSURAN_NP: Medan = {
  jalur: POLIS + 'Installment',
  label: 'Installment',
  jenis: 'tampil',
  panjangMaks: 2,
  aksi: [{ aksi: 'FillPaymentInstallmentEDMT' }],
}

// ------------------------------------------------------------------ deret satu baris (pola NB tataletak.ts)

const DERET_Q = [POLIS + 'Quartal', POLIS + 'YearOfQuartal', POLIS + 'TreatyYear']
const DERET_LAYER = [POLIS + 'LayerType', POLIS + 'Layer', POLIS + 'LayerPartType', POLIS + 'LayerPart']
const TANPA_LABEL = new Set([POLIS + 'LayerType', POLIS + 'Layer', POLIS + 'LayerPart'])

/** Deret layer NonProp ("layer 1 Of layer 1"). */
export const deretLayer = (d: Medan[]) => d[0]?.jalur === DERET_LAYER[0]

/** Satukan "Q / U/Y" (tiga medan berurutan) dan Layer (empat medan berurutan, label nama properti dibuang) menjadi
 *  satu deret; medan lain tetap sendiri. Salinan `deretQ` NB. */
export function deretQ(ms: Medan[]): (Medan | Medan[])[] {
  const hasil: (Medan | Medan[])[] = []
  for (let i = 0; i < ms.length; i++) {
    const pola = [DERET_Q, DERET_LAYER].find((p) => {
      const potong = ms.slice(i, i + p.length)
      return potong.length === p.length && potong.every((m, n) => m.jalur === p[n])
    })
    if (pola === undefined) {
      hasil.push(ms[i]!)
      continue
    }
    hasil.push(ms.slice(i, i + pola.length).map((m) => (TANPA_LABEL.has(m.jalur) ? { ...m, label: '' } : m)))
    i += pola.length - 1
  }
  return hasil
}

// ------------------------------------------------------------------ tab Old Data / New Data / Value Difference

/**
 * Varian section tab:
 *   lama      `PropOldData`          (`.OldData.EDMNo = ''`)      awalan `.OldData.`
 *   lama2     `PropOldData2`         (`.OldData.EDMNo != ''`)     awalan `.OldData.TreatyDifference.`
 *   baru      `PropNewData`          (`PN != 'ReasTreatyInAdmin'`) awalan `.` - hanya-baca
 *   baruAdmin `PropNewData2`         (`PN = 'ReasTreatyInAdmin'`)  awalan `.` - dapat diisi
 *   selisih   `PropValueDifference`  (selalu)                      awalan `.TreatyDifference.`
 */
export type VarianTab = 'lama' | 'lama2' | 'baru' | 'baruAdmin' | 'selisih'

export const AWALAN_TAB: Record<VarianTab, string> = {
  lama: POLIS + 'OldData.',
  lama2: POLIS + 'OldData.TreatyDifference.',
  baru: POLIS,
  baruAdmin: POLIS,
  selisih: POLIS + 'TreatyDifference.',
}

/** Satu wadah sel bagian uang; `judul` = LABEL Heading 4 ("OGP" / "ONP"). Tombol Save `PropNewData2` S7 di dalam
 *  wadah DIBUANG (keputusan work owner 07-10-2026 "HAPUS AJA") - Save kaki halaman (S19) tetap. */
export interface Kelompok {
  judul?: string
  medan: Medan[]
}

export interface TataTab {
  varian: VarianTab
  awalan: string
  /** Wadah S2 (Gross Premium 100%). */
  gross: Medan[]
  /** Kolom kiri (OGP ...) dan kanan (ONP ...). */
  kiri: Kelompok[]
  kanan: Kelompok[]
  /** Grid spreading dapat diisi (Type Treaty, % Share, % Share klaim, tombol Add). */
  /** Total berlabel di bawah grid spreading - hanya `PropNewData` S11 (selainnya `1=2`). */
  totalBerlabel: Medan[]
  /** Sel `.Installment` tampil (`PropNewData`: mati). */
  installment: Medan | null
}

/** Action set: `refresh <Act>(Data=..)` LALU `refresh CountOGPONP_Act`, lalu (bila sel me-refresh section
 *  `DetailPolicyTreatyInPropValueDifference` ber-defer-load) `EDMTCalculateTreatyDifference`. */
const SELISIH: Aksi = { aksi: 'EDMTCalculateTreatyDifference' }

/** Sel persen spreading `PropNewData2`: CountSpreading_Act(Index=.pxListSubscript) lalu hitung ulang selisih.
 *  ⛔ PENYIMPANGAN SADAR (WO 08-10-2026 "COBA CEK TOMBOLITU, APAKAH MASIH DIPERLUKAN? KALAU SUDAH TIDAK DIHAPUS AJA!"): tombol "Calculate Value Difference"
 *  (`PropNewData2` S24) dibuang - sel spreading dan `.Installment` ikut menghitung ulang tab Value Difference
 *  (sel uang sudah sejak awal), server menghitung ulang lagi saat Save / Submit (`HitungSelisihGenerasi`). */
export const AKSI_SPREADING: Aksi[] = [{ aksi: 'CountSpreading' }, SELISIH]
const OGPONP: Aksi = { aksi: 'CountOGPONP' }

/** Kerangka satu tab - urutan sel VERBATIM section varian itu. */
export function tataTab(varian: VarianTab): TataTab {
  const a = AWALAN_TAB[varian]
  const admin = varian === 'baruAdmin'
  /** `% Deduction ...` di empat section hanya-baca; `(%) Deduction ...` di PropNewData2. */
  const persen = (s: string) => (admin ? `(%) ${s}` : `% ${s}`)
  const m = (prop: string, label: string, aksi?: Aksi[]): Medan => ({
    jalur: a + prop,
    label,
    jenis: admin ? 'angka' : 'tampil',
    aksi: admin ? aksi : undefined,
    sajian: admin ? UANG4_SUNTING : UANG4,
  })
  const t = (prop: string, label: string, tampil?: (h: Halaman) => boolean): Medan => ({
    jalur: a + prop,
    label,
    jenis: 'tampil',
    sajian: UANG4,
    ...(tampil ? { tampil } : {}),
  })
  // Balance Due To You / Us: `<awalan>.BalanceDueTo < 0` / `>= 0`. Balance Before Tax / Before Withholding Tax:
  // `.BalanceDueTo >= 0` TANPA awalan di KELIMA section (kejanggalan XML B2 bab 7 butir 1, ditiru VERBATIM).
  const dueToYou = (h: Halaman) => negatifTeks(nilai(h, a + 'BalanceDueTo'))
  const dueToUs = (h: Halaman) => !negatifTeks(nilai(h, a + 'BalanceDueTo'))
  const sebelumPajak = (h: Halaman) => !negatifTeks(nilai(h, POLIS + 'BalanceDueTo'))

  const ogp = [
    m('PremiOgp', 'Premi Ogp', [OGPONP, SELISIH]),
    // XML: RiCommOgp TANPA refresh section Value Difference (beda dengan RiCommOnp)
    m('RiCommOgp', persen('Deduction In A (OGP)'), [{ aksi: 'CountResult1', param: 'Pct' }, OGPONP]),
    m('ResultOgp1', 'Deduction In A (OGP)', [{ aksi: 'CountResult1', param: 'Amount' }, OGPONP, SELISIH]),
    m('OveriddingCommOgp', persen('Deduction In B (OGP)'), [
      { aksi: 'CountResult2Ogp', param: 'Pct' },
      OGPONP,
      SELISIH,
    ]),
    m('ResultOgp2', 'Deduction In B (OGP)', [{ aksi: 'CountResult2Ogp', param: 'Amount' }, OGPONP, SELISIH]),
  ]
  const onp = [
    m('PremiOnp', 'Premi Onp', [OGPONP, SELISIH]),
    m('RiCommOnp', persen('Deduction In A (ONP)'), [{ aksi: 'CountResult1Onp', param: 'Pct' }, OGPONP, SELISIH]),
    m('ResultOnp1', 'Deduction In A (ONP)', [{ aksi: 'CountResult1Onp', param: 'Amount' }, OGPONP, SELISIH]),
    m('OveriddingCommOnp', persen('Deduction In B (ONP)'), [
      { aksi: 'CountResult2Onp', param: 'Pct' },
      OGPONP,
      SELISIH,
    ]),
    m('ResultOnp2', 'Deduction In B (ONP)', [{ aksi: 'CountResult2Onp', param: 'Amount' }, OGPONP, SELISIH]),
  ]
  const klaim = m('Claim', 'Claim', [OGPONP, SELISIH])
  const salvage = m('SalvageValue', 'Salvage', [OGPONP, SELISIH])
  const excess = m('ExcessLoss', 'Excess Loss', [OGPONP, SELISIH])
  const potongan = [m('Deduction1', 'Deduction1', [OGPONP, SELISIH]), m('Deduction2', 'Deduction2', [OGPONP, SELISIH])]
  const pajak = [t('PPHValue', 'PPH 2%'), t('PPNValue', 'PPN 2.2%')]
  const balYou = t('BalanceDueTo', 'Balance Due To You', dueToYou)
  const balUs = t('BalanceDueTo', 'Balance Due To Us', dueToUs)
  const balTax = t('BalanceBeforeTax', 'Balance Before Tax', sebelumPajak)
  const balPPH = t('BalanceBeforePPH', 'Balance Before Withholding Tax (PPH 2.2)', sebelumPajak)
  // Gross Premium 100%: PropNewData2 Auto TANPA action set (NB memicu CalculatePremi_Act; EDM tidak)
  const gross = [m('GrossPremium', 'Gross Premium 100%')]
  // sel `.Installment` (label "Text Input"): dirender `components/TabData.tsx` sebagai isian polos (pola NB)
  const installment: Medan = {
    jalur: a + 'Installment',
    label: 'Installment',
    jenis: 'tampil',
    // PropNewData2 S19: change -> refresh FillPaymentInstallment(Installment=.Installment); + selisih (tombol S24
    // dibuang, lihat AKSI_SPREADING)
    aksi: admin ? [{ aksi: 'FillPaymentInstallment' }, SELISIH] : undefined,
  }

  if (admin) {
    return {
      varian,
      awalan: a,
      gross,
      // S4 OGP lalu S6 klaim/saldo | S5 ONP lalu S7 potongan (+ Save) / S8 pajak
      kiri: [
        { judul: BAGIAN.ogp, medan: ogp },
        {
          medan: [
            klaim,
            // OutstandingClaim: Auto, wajib, TANPA action set
            m('OutstandingClaim', 'Outstanding Claim'),
            salvage,
            excess,
            t('NetPremium', 'Total Premium Before Claim'),
            balYou,
            balTax,
            balPPH,
            balUs,
          ],
        },
      ],
      kanan: [{ judul: BAGIAN.onp, medan: onp }, { medan: potongan }, { medan: pajak }],
      totalBerlabel: [],
      installment,
    }
  }
  if (varian === 'baru') {
    return {
      varian,
      awalan: a,
      gross,
      // S4: OGP + klaim + Net Premium + Balance Due To You/Us; S5: ONP + Deduction1/2 (tanpa PPH/PPN)
      kiri: [
        { judul: BAGIAN.ogp, medan: [...ogp, klaim, salvage, excess, t('NetPremium', 'Net Premium'), balYou, balUs] },
      ],
      kanan: [{ judul: BAGIAN.onp, medan: [...onp, ...potongan] }],
      // S11 (tanpa syarat): pxNumber 2 desimal
      totalBerlabel: [
        { jalur: a + 'TotalSharePercentagePremium', label: 'Total %Share', jenis: 'tampil', sajian: DUA },
        { jalur: a + 'TotalPremium', label: 'Total Premium', jenis: 'tampil', sajian: DUA },
        { jalur: a + 'TotalSharePercentageClaim', label: 'Total %Share Claim', jenis: 'tampil', sajian: DUA },
        { jalur: a + 'TotalClaim', label: 'Total Claim', jenis: 'tampil', sajian: DUA },
      ],
      installment: null,
    }
  }
  // lama / lama2 / selisih: S4 OGP + klaim + Net Premium + empat Balance; S5 ONP + Deduction1/2 + S6 PPH/PPN
  return {
    varian,
    awalan: a,
    gross,
    kiri: [
      {
        judul: BAGIAN.ogp,
        medan: [...ogp, klaim, salvage, excess, t('NetPremium', 'Net Premium'), balYou, balTax, balPPH, balUs],
      },
    ],
    kanan: [{ judul: BAGIAN.onp, medan: [...onp, ...potongan] }, { medan: pajak }],
    totalBerlabel: [],
    installment,
  }
}

/** Varian tab Old Data dan New Data menurut halaman dan posisi. */
export function varianTab(tab: 'Old Data' | 'New Data' | 'Value Difference', h: Halaman, admin: boolean): VarianTab {
  if (tab === 'Old Data') return oldDataTanpaEDM(h) ? 'lama' : 'lama2'
  if (tab === 'New Data') return admin ? 'baruAdmin' : 'baru'
  return 'selisih'
}

// ------------------------------------------------------------------ grid

/** Teks sel `.TreatyType` grid spreading hanya-baca: label acuan (`.ID` -> `.Note`), lalu `.TreatyName`, lalu kode.
 *  Salinan `labelTreatyType` NB. */
export function labelTreatyType(b: Baris, opsi: Pilihan[]): string {
  return opsi.find((o) => o.nilai === b.TreatyType)?.label || b.TreatyName || b.TreatyType || ''
}

/** Grid spreading kelima section Prop: pxNumber 4 desimal (persen tersunting `pyShowReadonlyFormatting`). */
export const SAJIAN_SPREADING = {
  persen: { desimal: 4, formatSaatSunting: true } as Sajian,
  uang: { desimal: 4 } as Sajian,
  total: { desimal: 4 } as Sajian,
}

/** Grid spreading `DetailPolicyTreatyInAddPremi` S17: pxNumber 2 desimal, hanya-baca. */
export const SAJIAN_SPREADING_NP: Sajian = { desimal: 2 }

/** Grid `.ListInstallment` section Prop: InstallmentNo pxInteger, DueDate tanggal, % dan Premium 4 desimal,
 *  PaymentTotal 2. */
export const SAJIAN_ANGSURAN = {
  dueDate: TGL,
  persen: { desimal: 4 } as Sajian,
  premium: { desimal: 4 } as Sajian,
  total: DUA,
}
