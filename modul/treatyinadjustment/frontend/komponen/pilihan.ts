// Isi dropdown/autocomplete panel New — dari sumber yang ekspor nyatakan.
//
// ⛔ Satu entri per SUMBER yang tercatat di kerangka (`SumberPilihan`), bukan
// per medan. RD yang sama di Pega (`BrowseCurrency_RD` …) memberi daftar yang
// sama di sini, di medan mana pun ia dipakai.
//
// ⚠️ Daftar `associated` hidup di rule Property yang TIDAK diekspor. Yang
// dipakai di sini hanya yang punya bukti, dan buktinya disebut per entri;
// yang tanpa bukti TIDAK diisi — medannya tetap kotak teks.

import type { BarisBersarang, OpsiKepalaTreatyIn, OpsiLimitsTreatyIn, PilihanTreatyIn, SisiPenyesuaian } from '../api'
import type { SumberPilihan } from '../ekspor/jenis'

export interface Opsi {
  value: string
  label: string
  /**
   * Pengenal baris master (`.ID`, `.BizCode`, `.ReinsTypeID`) — nilai medan
   * yang IKUT diisi (`SumberPilihan.setel`) saat pilihan ini dipilih.
   */
  id?: string
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
  /**
   * ⭐ Daftar BERPARAMETER yang sudah dimuat (Class of Business per Treaty
   * Group, Spreading Type per Treaty Group + Commencement), per kunci muat.
   */
  muatan?: Readonly<Record<string, readonly Opsi[]>>
  /** Minta satu daftar berparameter dimuat; hasilnya masuk `muatan[kunci]`. */
  muat?: (kunci: string, ambil: () => Promise<Opsi[]>) => void
}

/**
 * Halaman tempat sel/medan itu duduk — untuk parameter RD. `halaman` =
 * skalar akar (`Commencement`) + skalar halaman/baris BERTITIK
 * (`.TreatyGroupID`); `sisi` = halaman itu sendiri (larik `.TreatyGroupList`).
 */
export interface TempatOpsi {
  halaman: Readonly<Record<string, string>>
  sisi?: SisiPenyesuaian
}

/** Rute pemuat daftar berparameter — disuntik `SisiPenyesuaian`, bukan diimpor di sini. */
export interface PemuatOpsi {
  kelasBisnis: (treatyGroupId: string) => Promise<readonly PilihanTreatyIn[]>
  indukSpreading: (treatyGroupId: string, mulai: string) => Promise<readonly { reinsTypeId: string; reinsTypeName: string }[]>
}

const nama = (xs: readonly PilihanTreatyIn[] | undefined): Opsi[] | undefined =>
  xs?.map((x) => ({ value: x.nama, label: x.nama, id: x.id }))

/**
 * Nilai menurut medan NILAI RD-nya: `ID` → pengenal (tampil nama), selain
 * itu nama. ⛔ RALAT 7 Oktober 2026: Treaty Group (`.TreatyGroupID`, rincian
 * Limits Prop) dan Treaty Type (`.TreatyTypeID`) bernilai `ID` di ekspor,
 * tetapi dahulu diisi NAMA — nama masuk ke medan pengenal, dan
 * `SetTreatyGroupName_Act` / `SetTreatyTypeName_Act` tidak menemukannya.
 */
const menurutNilai = (sp: SumberPilihan, xs: readonly PilihanTreatyIn[] | undefined): Opsi[] | undefined =>
  sp.nilai === 'ID' ? xs?.map((x) => ({ value: x.id, label: x.nama, id: x.id })) : nama(xs)

/**
 * Nilai satu parameter RD atas halaman sel/medan itu — ungkapan Pega apa
 * adanya: `"10001"` (harfiah), `TreatyIn.X` (akar), `.L(n).X` (baris ke-n
 * larik halaman ini), `.X` (skalar halaman/baris), `Halaman.X` (halaman sesi).
 */
export function nilaiParam(ekspr: string, t: TempatOpsi): string {
  const e = ekspr.trim()
  if (/^".*"$/.test(e)) return e.slice(1, -1)
  if (e.startsWith('TreatyIn.')) return t.halaman[e.slice('TreatyIn.'.length)] ?? ''
  const larik = /^\.(\w+)\((\d+)\)\.(\w+)$/.exec(e)
  if (larik !== null) {
    const v = t.sisi?.larik[larik[1] ?? '']?.[Number(larik[2]) - 1]?.[larik[3] ?? '']
    return typeof v === 'string' ? v : ''
  }
  if (e.startsWith('.')) return t.halaman[e] ?? t.sisi?.medan[e.slice(1)] ?? ''
  return t.halaman[e] ?? ''
}

/**
 * Daftar BERPARAMETER: dari `muatan` bila sudah dimuat; bila belum, minta
 * dimuat dan kembalikan daftar KOSONG (dropdown, bukan kotak teks).
 */
function berparameter(d: DataOpsi, kunci: string, ambil: () => Promise<Opsi[]>): Opsi[] {
  const ada = d.muatan?.[kunci]
  if (ada !== undefined) return [...ada]
  d.muat?.(kunci, ambil)
  return []
}

/**
 * ⭐ RD yang daftarnya BERGANTUNG halaman (8 Oktober 2026) — rute Treaty In
 * yang SAMA dengan layar Treaty In:
 *   `BrowseTreatyBusinessWOType_RD` → `kelas-bisnis?treatyGroupId=`
 *     (`pTreatyGroupId`), nilai `.BIZNAME`, `id` `.BizCode`;
 *   `BrowseTreatyArrangement_ParentReinsMasterTrt` → `spreading-induk`
 *     (`TreatyGroupID`, `StartDate`), nilai `.ReinsTypeID` atau
 *     `.ReinsTypeName` menurut sel, tampil `.ReinsTypeName`.
 */
function opsiBerparameter(sp: SumberPilihan, d: DataOpsi, t: TempatOpsi | undefined, p: PemuatOpsi | undefined): Opsi[] | undefined {
  if (t === undefined || p === undefined) return undefined
  const param = sp.param ?? {}
  switch (sp.rd) {
    case 'BrowseTreatyBusinessWOType_RD': {
      const grup = nilaiParam(param.pTreatyGroupId ?? '.TreatyGroupID', t)
      return berparameter(d, `kelas-bisnis|${grup}`, () =>
        p.kelasBisnis(grup).then((xs) => xs.map((x) => ({ value: x.nama, label: x.nama, id: x.id }))),
      )
    }
    case 'BrowseTreatyArrangement_ParentReinsMasterTrt': {
      const grup = nilaiParam(param.TreatyGroupID ?? '.TreatyGroupID', t)
      const mulai = nilaiParam(param.StartDate ?? 'TreatyIn.Commencement', t)
      const xs = berparameter(d, `spreading|${grup}|${mulai}`, () =>
        p.indukSpreading(grup, mulai).then((ys) => ys.map((y) => ({ value: y.reinsTypeId, label: y.reinsTypeName, id: y.reinsTypeId }))),
      )
      // Nilai `.ReinsTypeName` (Share Non-Prop) — yang tersimpan namanya.
      return sp.nilai === 'ReinsTypeName' ? xs.map((x) => ({ ...x, value: x.label })) : xs
    }
    default:
      return undefined
  }
}

/**
 * Medan yang IKUT diisi saat `nilai` dipilih (`SumberPilihan.setel`): tiap
 * target ← `id` pilihan itu. Nilai di luar daftar (diketik bebas) mengosongkan
 * targetnya — pasangan basi lebih buruk daripada kosong.
 */
export function ikutTerpilih(sp: SumberPilihan | null | undefined, opsi: readonly Opsi[] | undefined, nilai: string): Record<string, string> {
  if (sp?.setel === undefined || opsi === undefined) return {}
  const o = opsi.find((x) => x.value === nilai)
  const out: Record<string, string> = {}
  for (const s of sp.setel) out[s.target] = o?.id ?? ''
  return out
}

/**
 * `Option` tab Share Prop — nilainya `1`/`2`.
 *
 * ⭐ *(8 Okt, E)* Label `2` kini DIKETAHUI: rule Property `OptionLimit`
 * (kelas `ASM-FW-GISFW-Int-TREATY_IN`, Table type `Prompt List`) dikirim
 * pemilik proses 7 Oktober 2026 — 1 → `Of Cession to R/I`, 2 → `Of 100%
 * Limit`. Modul Treaty In sudah memakainya (`TabShareProp.tsx` `OPSI_SHARE`,
 * `backend/services/prompt_value.go`); di sini disamakan, bukan ditebak.
 */
const OPSI_OPTION_LIMIT: readonly Opsi[] = [
  { value: '1', label: 'Of Cession to R/I' },
  { value: '2', label: 'Of 100% Limit' },
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
export function opsiUntuk(
  sp: SumberPilihan,
  kunci: string,
  d: DataOpsi,
  tempat?: TempatOpsi,
  pemuat?: PemuatOpsi,
): Opsi[] | undefined {
  if (sp.sumber === 'reportdefinition') {
    const bp = opsiBerparameter(sp, d, tempat, pemuat)
    if (bp !== undefined) return bp
    switch (sp.rd) {
      case 'BrowseCurrency_RD':
        // Nilai `.ID` (grid kurs New) atau `.Currency`; tampil `.Currency`.
        return sp.nilai === 'ID'
          ? d.limits?.mataUang.map((x) => ({ value: x.id, label: x.nama, id: x.id }))
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
      // ⭐ 8 Oktober 2026 — rincian Layers Non-Prop: daftar yang SAMA dengan
      // tab Limits Non-Prop Treaty In (`TabLimitsNonProp.tsx` memakai
      // `jenisLayer` untuk `Layer` DAN `Part Of`).
      case 'LayerType':
      case 'LayerPartType':
        return d.limits?.jenisLayer
      case 'Cover':
        return d.limits?.cover
      case 'CurrencyRelation':
        return d.limits?.relasiMataUang
      case 'ReinstatementNote':
        return d.limits?.catatanReinstatement
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
