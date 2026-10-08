// Isi dropdown/autocomplete panel New — dari sumber yang ekspor nyatakan.
//
// ⛔ Satu entri per SUMBER yang tercatat di kerangka (`SumberPilihan`), bukan
// per medan. RD yang sama di Pega (`BrowseCurrency_RD` …) memberi daftar yang
// sama di sini, di medan mana pun ia dipakai.
//
// ⚠️ Daftar `associated` hidup di rule Property yang TIDAK diekspor. Yang
// dipakai di sini hanya yang punya bukti, dan buktinya disebut per entri;
// yang tanpa bukti TIDAK diisi — medannya tetap kotak teks.

import type { BarisBersarang, OpsiKepalaTreatyIn, OpsiLimitsTreatyIn, PilihanTreatyIn } from '../api'
import type { SumberPilihan } from '../ekspor/jenis'

export interface Opsi {
  value: string
  label: string
}

/** Daftar mentah yang sudah dimuat dari rute Treaty In. */
export interface DataOpsi {
  limits?: OpsiLimitsTreatyIn
  kepala?: OpsiKepalaTreatyIn
  agen?: PilihanTreatyIn[]
  /**
   * Larik AKAR panel — sumber `pageList` halaman SESI (`TempQuarter.pxResults`,
   * `TempQuarterYear.pxResults`, nilai `.CARI6`) yang `GetAchievement` isi.
   */
  halaman?: Readonly<Record<string, BarisBersarang[]>>
}

const nama = (xs: readonly PilihanTreatyIn[] | undefined): Opsi[] | undefined => xs?.map((x) => ({ value: x.nama, label: x.nama }))

/**
 * Nilai menurut medan NILAI RD-nya: `ID` → pengenal (tampil nama), selain
 * itu nama. ⛔ RALAT 7 Oktober 2026: Treaty Group (`.TreatyGroupID`, rincian
 * Limits Prop) dan Treaty Type (`.TreatyTypeID`) bernilai `ID` di ekspor,
 * tetapi dahulu diisi NAMA — nama masuk ke medan pengenal, dan
 * `SetTreatyGroupName_Act` / `SetTreatyTypeName_Act` tidak menemukannya.
 */
const menurutNilai = (sp: SumberPilihan, xs: readonly PilihanTreatyIn[] | undefined): Opsi[] | undefined =>
  sp.nilai === 'ID' ? xs?.map((x) => ({ value: x.id, label: x.nama })) : nama(xs)

/**
 * `Option` tab Share Prop — nilainya `1`/`2`. Hanya label `1` yang terbaca
 * (tangkapan layar Pega pemilik proses); label `2` TIDAK ditebak — sama
 * persis dengan modul Treaty In (`TabShareProp.tsx`).
 */
const OPSI_OPTION_LIMIT: readonly Opsi[] = [
  { value: '1', label: 'Of Cession to R/I' },
  { value: '2', label: '2' },
]

/**
 * `AccumulationPeriod` — keempat nilai yang DataTransform
 * `TreatyInAddAccumulation` periksa (`=="quarter"`, `"half"`, `"month"`,
 * `"none"`). Labelnya tidak diekspor: nilai tampil apa adanya.
 */
const OPSI_PERIODE_AKUMULASI: readonly Opsi[] = ['quarter', 'half', 'month', 'none'].map((v) => ({ value: v, label: v }))

/**
 * Daftar untuk satu sel, atau `undefined` bila sumbernya belum dikenal /
 * belum dimuat — pemanggil lalu memakai kotak teks.
 */
export function opsiUntuk(sp: SumberPilihan, kunci: string, d: DataOpsi): Opsi[] | undefined {
  if (sp.sumber === 'reportdefinition') {
    switch (sp.rd) {
      case 'BrowseCurrency_RD':
        // Nilai `.ID` (grid kurs New) atau `.Currency`; tampil `.Currency`.
        return sp.nilai === 'ID'
          ? d.limits?.mataUang.map((x) => ({ value: x.id, label: x.nama }))
          : nama(d.limits?.mataUang)
      case 'BrowseCurrencyTreatyIn_RD':
        return nama(d.limits?.mataUang)
      case 'BrowseTreatyGroup_RD':
        return menurutNilai(sp, d.limits?.kelompokTreaty)
      case 'BrowseReinsuranceType_RD':
        // Nilai `.Note` — kolom `NOTE` REINSURANCETYPE, `nama` di rute Treaty
        // In; atau `.ID` (Treaty Type rincian Limits Prop).
        return menurutNilai(sp, d.limits?.jenisTreaty)
      case 'BrowseAgentNusaRe_RD':
        return nama(d.agen)
      default:
        return undefined
    }
  }
  if (sp.sumber === 'pageList' && sp.halaman?.startsWith('TempQuarter') === true) {
    const daftar = d.halaman?.[sp.halaman.replace(/\.pxResults$/, '')]
    return daftar?.map((b) => {
      const v = typeof b.CARI6 === 'string' ? b.CARI6 : ''
      return { value: v, label: v }
    })
  }
  if (sp.sumber === 'associated') {
    switch (kunci) {
      case 'ReportingPeriod':
        return d.kepala?.periodePelaporan
      case 'AccountingMode':
        return d.kepala?.caraPembukuan
      case 'AccountingModeNonProp':
        return d.kepala?.caraPembukuanNonProp
      case 'Bordeaux':
        return d.kepala?.bordereaux
      case 'TreatyGroup':
        return nama(d.limits?.kelompokTreaty)
      case 'OptionLimit':
        return [...OPSI_OPTION_LIMIT]
      case 'AccumulationPeriod':
        return [...OPSI_PERIODE_AKUMULASI]
      default:
        return undefined
    }
  }
  return undefined
}

/**
 * Daftar yang DITAMPILKAN: nilai tersimpan yang tidak ada di daftar tetap
 * ditawarkan (ditandai), bukan dijatuhkan — dropdown yang diam-diam
 * mengosongkan nilai lama terbaca sebagai data yang hilang.
 */
export function denganNilaiKini(opsi: readonly Opsi[], nilai: string, tandaAsing: string): Opsi[] {
  if (nilai === '' || opsi.some((o) => o.value === nilai)) return [...opsi]
  return [{ value: nilai, label: `${nilai} ${tandaAsing}` }, ...opsi]
}
