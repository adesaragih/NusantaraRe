// Cabang NonProp baru (`DetailPolicyTreatyInAddendum` S11, `.IsNewPolicyNonProp = 1`): label VERBATIM dan format sel
// `Section/DetailPolicyTreatyInAddPremi` (3 grid XOL per mata uang), `Section/DetailPolicyAddPremiDetail` (rincian
// per layer, expand pane), `DetailPolicyTreatyInAddendum` S13 (angsuran per mata uang) dan `Section/InstallmentList`
// (rincian angsuran - section bersama NB, isinya sama persis). Pola kolom: `modul/nbtreatyin/frontend/nonprop.ts`.
//
// ⛔ SEMUA hanya-baca: sel grid XOL ber-`Auto` tetapi edit mode expandPane dan panel rinciannya (`DetailPolicyAdd
// PremiDetail`, memo "made everything read only") seluruhnya RO; backend tidak menerima daftar XOL dari layar.

import { POLIS } from './api'
import { BAGIAN } from './labels'
import { POLA_INTI, type Sajian } from './sajian'

/** Satu kolom grid: anggota baris (`''` = sel LABEL), judul kolom VERBATIM, teks sel LABEL, format sel. */
export interface Kolom {
  m: string
  judul: string
  /** Teks sel LABEL baris data (mis. "Total For Currency", "Part of"). */
  teks?: string
  /** Format sel; tidak ada = teks apa adanya. */
  format?: Sajian
  /** `pyVisible NOTBLANK`. */
  bilaAda?: boolean
}

const DUA: Sajian = { desimal: 2 }
/** ⛔ Perintah work owner 07-10-2026 (screenshot Current Premium): "JANGAN NULL TAPI 0" + "BUAT 4 ANGKA BELAKANG KOMA" -
 *  sel uang grid XOL dan rincian layernya 4 desimal, kosong / nol tampil "0" (pola bagian uang OGP / ONP). XML: pxNumber
 *  tanpa `pyDecimalPlaces` (induk) dan 2 desimal (rincian) - penyimpangan sadar. Nilai tersimpan tidak diubah. */
const UANG4: Sajian = { desimal: 4, nolPolos: true }

/** Tiga grid `AddPremi`: wadah berjudul (pyIncludeHeader=true) dan daftarnya. */
export const GRID_XOL = [
  { judul: BAGIAN.premiLama, daftar: POLIS + 'OldData.TreatyXOLList' },
  { judul: BAGIAN.premiBaru, daftar: POLIS + 'TreatyXOLList' },
  { judul: BAGIAN.premiSelisih, daftar: POLIS + 'TreatyXOLDifferenceList' },
] as const

/** Kolom ketiga grid XOL (identik): sel 1 LABEL "Total For Currency" di bawah judul kosong; angka `UANG4`. */
export const KOLOM_XOL: Kolom[] = [
  { m: '', judul: '', teks: 'Total For Currency' },
  { m: 'Currency', judul: 'Currency' },
  { m: 'GrossPremi', judul: 'Gross Premium', format: UANG4 },
  { m: 'Deduction', judul: 'Deduction', format: UANG4 },
  { m: 'PPNValue', judul: 'PPN', format: UANG4 },
  { m: 'PPHValue', judul: 'PPh', format: UANG4 },
  { m: 'NetPremi', judul: 'Net Premium', format: UANG4 },
  { m: 'NetPremiAfterPPN', judul: 'Net Premium After PPN', format: UANG4 },
  { m: 'NetPremiAfterPPH', judul: 'Net Premium After PPh', format: UANG4 },
  { m: 'NetPremiAfterTax', judul: 'Net Premium After Tax', format: UANG4 },
]

/** Anak tiap baris XOL - grid `.ValueList` `DetailPolicyAddPremiDetail`. */
export const ANAK_XOL = 'ValueList'

/** Kolom `DetailPolicyAddPremiDetail` (16 sel, judul baris 1 VERBATIM - lima judul pertama kosong): Layer* pxTextInput
 *  number 0 desimal, `.Currency` tampil bila tidak kosong, angka `UANG4`. */
export const KOLOM_XOL_RINCI: Kolom[] = [
  { m: 'LayerType', judul: '', format: { desimal: 0 } },
  { m: 'Layer', judul: '', format: { desimal: 0 } },
  { m: '', judul: '', teks: 'Part of' },
  { m: 'LayerPartType', judul: '', format: { desimal: 0 } },
  { m: 'LayerPart', judul: '', format: { desimal: 0 } },
  { m: 'Currency', judul: '', bilaAda: true },
  { m: 'GrossPremi', judul: 'Gross Premium', format: UANG4 },
  { m: 'Currency', judul: '', bilaAda: true },
  { m: 'Deduction', judul: 'Deduction', format: UANG4 },
  { m: 'PPNValue', judul: 'PPN', format: UANG4 },
  { m: 'PPHValue', judul: 'PPh', format: UANG4 },
  { m: 'Currency', judul: '', bilaAda: true },
  { m: 'NetPremi', judul: 'Net Premium', format: UANG4 },
  { m: 'NetPremiAfterPPN', judul: 'Net Premium After PPN', format: UANG4 },
  { m: 'NetPremiAfterPPH', judul: 'Net Premium After PPh', format: UANG4 },
  { m: 'NetPremiAfterTax', judul: 'Net Premium After Tax', format: UANG4 },
]

/** Grid "Installment Data Information" S13 (`.ListInstallment`, expand pane `InstallmentList`): Currency / Total. */
export const DAFTAR_ANGSURAN = POLIS + 'ListInstallment'
export const ANAK_ANGSURAN = 'InstallmentList'
export const KOLOM_ANGSURAN_NP: Kolom[] = [
  { m: 'Currency', judul: 'Currency' },
  { m: 'Premium', judul: 'Total', format: POLA_INTI },
]

/** `Section/InstallmentList` - sama dengan NB `KOLOM_RINCI` (format menurut screenshot layar Pega NB 06-10-2026):
 *  saldo dua desimal, Percentage pola inti, judul kolom ke-3 kosong. */
export const KOLOM_RINCI_ANGSURAN: Kolom[] = [
  { m: 'DueDate', judul: 'Payment Date', format: 'tanggal' },
  { m: 'InstallmentPercentage', judul: 'Percentage', format: POLA_INTI },
  { m: 'Currency', judul: '' },
  { m: 'Premium', judul: 'Balance Before Tax', format: DUA },
  { m: 'PremiumAfterPPN', judul: 'Balance Before Withholding Tax (PPH 2.2)', format: DUA },
  { m: 'PremiumAfterTax', judul: 'Balance Due To', format: DUA },
]
