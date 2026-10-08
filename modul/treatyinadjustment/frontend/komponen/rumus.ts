// Rumus tombol panel New — dijalankan lewat rute `/hitung/*` modul Treaty In.
//
// ---------------------------------------------------------------------
// ⭐ MENGAPA RUTE TREATY IN, BUKAN SALINAN RUMUS
// ---------------------------------------------------------------------
// Panel New Adjustment di Pega ADALAH form Treaty In: Section yang sama, dan
// Activity rumusnya IDENTIK di kedua korpus — diadu langkah demi langkah
// 7 Oktober 2026: TreatyInNPSetTotal, TreatyInEGNPIListValue,
// TreatyInSetValueInstallment, TreatyInSummaryLimit/MDP,
// TreatyInLimitsListValue, TreatyInXOLAddSpreading, TreatyInSetBrokerage,
// TreatyInSummaryLimitShare, TreatyInShareListValue, TreatyInSetReport,
// TreatyInNonAddItem, TreatyInPropAdd, TreatyInPropshare — 14 dari 14 sama.
// Menyalin rumusnya ke modul ini membuat dua salinan yang suatu hari
// berselisih. Modul tidak boleh saling IMPOR; memanggil rute HTTP modul lain
// bukan impor (preseden: claimlife → `/api/polis-life`).
//
// ⛔ Rute `/hitung/*` MURNI menghitung — nol tulisan ke basis data. Hasilnya
// masuk keadaan layar saja, sama seperti isian.
//
// ⚠️ Bentuk masukan rute itu kini KONTRAK BERSAMA dengan modul Treaty In;
// perubahan bentuknya di sana harus diberitahukan ke modul ini.
//
// ---------------------------------------------------------------------
// ⭐ 7 Oktober 2026 — KETERGANTUNGAN ANTARTAB DAN SECTION RINCIAN
// ---------------------------------------------------------------------
// Rumus Section rincian (`Layers`, `DetailLimits`, `DetailShare`, `Share`,
// `DetailEGNPI`) berjalan di halaman BARIS, dan banyak yang menulis ke AKAR
// panel (`TreatyInPropshareDetail` menyusun ulang `Limits`). Satu rantai
// kini berjalan atas AKAR + JALUR baris pemicunya (`baris.ts`), dan
// hasilnya — akar sebelum dan sesudah — diterapkan GABUNG TIGA ARAH di
// `SisiPenyesuaian`: isian yang diketik selama rute menjawab tidak hilang.
//
// Syarat aksi (`pyActionConditions`, `AksiTombol.syarat`) dinilai atas
// halaman PEMICU; aksi yang syaratnya tidak terpenuhi dilewati.

import { minta } from '../../../../inti/frontend/klien'
import type { BarisBersarang, OpsiLimitsTreatyIn, PilihanTreatyIn, SisiPenyesuaian } from '../api'
import type { AksiTombol, TombolKerangka } from '../ekspor/jenis'
import { syaratTerpenuhi } from '../ekspor/syarat'
import { PENYESUAIAN } from '../labelsPenyesuaian'
import { ambilHalaman, anak, sisiKeBaris, teks, ubahHalaman, type LangkahJalur } from './baris'
import { tambah } from './desimal'

/** Prefix rute rumus modul Treaty In. */
export const PREFIX_HITUNG_TREATYIN = '/api/treaty-in/hitung'

type Baris = BarisBersarang
type NilaiMataUang = { Currency: string; CurrencyID: string; Value: string }

/** Perubahan satu halaman — larik dan/atau skalar. */
export interface Ubahan {
  larik?: Record<string, Baris[]>
  medan?: Record<string, string>
}

/** Perubahan yang satu rumus kembalikan — ditimpakan ke keadaan sisi New. */
export interface HasilRumus {
  /** Perubahan halaman TEMPAT rumus berjalan (akar, atau baris rincian). */
  larik: Record<string, Baris[]>
  medan: Record<string, string>
  /** Pesan Activity (`Property-Set-Messages`), apa adanya. */
  pesan: string[]
  /**
   * Perubahan di AKAR panel (`TreatyIn.*`) — rumus Section rincian yang
   * menulis ke luar barisnya (`TreatyInPropshareDetail` → `Limits`,
   * `FetchQSfromMasterXOL` → seluruh `Share`).
   */
  akar?: Ubahan
}

/** Hasil satu RANTAI — perubahan halaman pemicu, plus akar sebelum dan sesudah. */
export interface HasilRantai extends HasilRumus {
  /** Akar panel saat rantai dipicu — halaman pemicu sudah ditaruh di jalurnya. */
  awal: SisiPenyesuaian
  /** Akar panel sesudah seluruh langkah; diterapkan GABUNG TIGA ARAH. */
  akhir: SisiPenyesuaian
}

/** Yang `KonteksKerangka.terapkan` terima — hasil rumus, dengan atau tanpa akar. */
export type HasilTerapan = HasilRumus & Partial<Pick<HasilRantai, 'awal' | 'akhir'>>

/** Sel grid pemicu — perilaku `change` sel, atau baris yang dihapus. */
export interface SelRumus {
  larik: string
  /** Indeks baris, mulai 0 (`.pxListSubscript` − 1). */
  indeks: number
  /** Kunci sel; kosong untuk tombol baris (Remove). */
  kunci: string
}

/**
 * Halaman lain yang SEBAGIAN rumus baca, di luar halaman tempat ia berjalan.
 *
 * ⭐ Section rincian baris menjalankan Activity di halaman BARIS
 * (`SetTotalInstallment` atas satu `Installment`), sementara masukannya
 * sebagian di akar (`TreatyIn.TotalShareNetNP`). `TreatyInSetValueInstallment`
 * langkah 10 membaca `TreatyIn.OLDDATA` — sisi Old.
 */
export interface LingkupRumus {
  /** Akar panel (`TreatyIn`) — sama dengan halaman rumus bila ia berjalan di akar. */
  panel?: SisiPenyesuaian
  /** `TreatyIn.OLDDATA` — sisi Old penyesuaian. */
  lama?: SisiPenyesuaian
  /**
   * Jalur halaman rumus dari akar — `[]` di akar. TANPA jalur, halaman
   * rumus diperlakukan sebagai akarnya sendiri dan `panel` dibaca apa adanya.
   */
  jalur?: readonly LangkahJalur[]
  /** Sel grid pemicu — `.pxListSubscript` dan kuncinya. */
  sel?: SelRumus
  /**
   * Halaman PEMICU — akar + skalar bertitik halaman/baris pemicu
   * (`.Note`, `.Layer`). Syarat aksi dan parameter `.X` dibaca di sini.
   */
  halaman?: Readonly<Record<string, string>>
  /** Daftar master yang dropdown panel ini muat — pasangan id ↔ nama. */
  master?: OpsiLimitsTreatyIn
}

type Rumus = (s: SisiPenyesuaian, a: AksiTombol, l: LingkupRumus) => Promise<HasilRumus>

const kirim = <T>(rute: string, badan: unknown) => minta<T>(`${PREFIX_HITUNG_TREATYIN}/${rute}`, { metode: 'POST', badan })

const kosong = (): HasilRumus => ({ larik: {}, medan: {}, pesan: [] })

/**
 * Baris hasil ditimpakan ke baris asal per indeks — kunci yang rute itu
 * tidak kenal (`pyTemplateRichTextEditor`, …) TIDAK hilang.
 */
const timpa = (asal: readonly Baris[] | undefined, hasil: readonly Baris[]) =>
  hasil.map((b, i) => ({ ...(asal?.[i] ?? {}), ...b }))

/** Nilai kosong Go (`""`, `[]`) setara dengan kunci yang tidak ada. */
const kosongGo = (v: unknown) => v === undefined || v === '' || (Array.isArray(v) && v.length === 0)

/**
 * Seperti `timpa`, tetapi kunci yang KOSONG di jawaban dan TIDAK ADA di asal
 * tidak dilahirkan — rute bertipe Go menjawab setiap medan strukturnya.
 */
const timpaIsi = (asal: readonly Baris[] | undefined, hasil: readonly Baris[]): Baris[] =>
  hasil.map((h, i) => {
    const a = asal?.[i] ?? {}
    const out: Baris = { ...a }
    for (const [k, v] of Object.entries(h)) if (!(kosongGo(v) && kosongGo(a[k]))) out[k] = v
    return out
  })

const kurs = (s: SisiPenyesuaian) =>
  (s.larik.CurrencyList ?? []).map((k) => ({ Currency: teks(k.Currency), Conversion: teks(k.Conversion) }))

const keBaris = (xs: readonly NilaiMataUang[]): Baris[] => xs.map((x) => ({ ...x }))

/** Akar panel tempat rumus berjalan. */
const akarDari = (s: SisiPenyesuaian, l: LingkupRumus) => l.panel ?? s

/** `.pxListSubscript` halaman rumus (mulai 1) — 0 di akar. */
const nomorHalaman = (l: LingkupRumus) => {
  const j = l.jalur ?? []
  const t = j[j.length - 1]
  return t === undefined ? 0 : t.indeks + 1
}

/** Indeks baris akar `larik` yang memuat halaman rincian ini (`Limits(n)`, `Share(n)`), atau 0. */
const indeksAkar = (l: LingkupRumus, larik: string) => {
  const j = (l.jalur ?? [])[0]
  return j !== undefined && j.larik === larik ? j.indeks : 0
}

/**
 * Nilai satu parameter Activity sebagaimana Pega membacanya: `.X` dari
 * halaman PEMICU (baris sel, atau halaman rincian), `TreatyIn.X` dari akar,
 * `"x"` harfiah, `.pxListSubscript` = nomor baris pemicu (mulai 1).
 */
export function nilaiParam(v: string | undefined, s: SisiPenyesuaian, l: LingkupRumus): string {
  if (v === undefined) return ''
  if (v === '.pxListSubscript') return String(l.sel !== undefined ? l.sel.indeks + 1 : nomorHalaman(l))
  if (v.startsWith('.')) return l.halaman?.[v] ?? s.medan[v.slice(1)] ?? ''
  if (v.startsWith('TreatyIn.')) return akarDari(s, l).medan[v.slice('TreatyIn.'.length)] ?? ''
  return v.replace(/^"(.*)"$/, '$1')
}

/**
 * Tanggal ke bentuk SIMPAN `YYYYMMDD`. Kotak tanggal panel New mengirim
 * `YYYY-MM-DD` (sel) atau `DD-MM-YYYY` (medan); rute Treaty In membaca
 * `YYYYMMDD` dan `DD-MM-YYYY` saja. Bentuk lain dikirim apa adanya.
 */
export function keSimpanTanggal(v: string | undefined): string {
  const t = (v ?? '').trim()
  let m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(t)
  if (m !== null) return `${m[1] ?? ''}${m[2] ?? ''}${m[3] ?? ''}`
  m = /^(\d{2})-(\d{2})-(\d{4})$/.exec(t)
  if (m !== null) return `${m[3] ?? ''}${m[2] ?? ''}${m[1] ?? ''}`
  return t
}

/** `TreatyInNPSetTotal(type=retention)` — Σ Amount per mata uang. */
const totalRetensi: Rumus = async (s) => {
  const h = await kirim<{ retensi: Baris[]; TotalRetentionAmountNP: NilaiMataUang[]; pesan: string[] }>('retensi', {
    aksi: 'total',
    retensi: s.larik.Retention ?? [],
    indeks: 0,
  })
  return {
    larik: { Retention: timpa(s.larik.Retention, h.retensi), TotalRetentionAmountNP: keBaris(h.TotalRetentionAmountNP) },
    medan: {},
    pesan: h.pesan,
  }
}

/** Baris Retention yang rute EGNPI baca (`TreatyInNonAddItem(egnpi)` langkah 3). */
const retensiUntukEgnpi = (s: SisiPenyesuaian) =>
  (s.larik.Retention ?? []).map((r) => ({ Currency: teks(r.Currency), CurrencyID: teks(r.CurrencyID) }))

/** `TreatyInNPSetTotal(type=egnpi)` dan `TreatyInEGNPIListValue` — satu rute, dua aksi. */
const egnpi =
  (aksi: 'total' | 'nilai'): Rumus =>
  async (s) => {
    const h = await kirim<{
      egnpi: Baris[]
      TotalEgnpiAmount: string
      TotalEgnpiProportion: string
      TotalEgnpiAmountNP: NilaiMataUang[]
      pesan: string[]
    }>('egnpi', {
      aksi,
      egnpi: s.larik.EGNPI ?? [],
      kurs: kurs(s),
      retensi: retensiUntukEgnpi(s),
      indeks: 0,
    })
    return {
      larik: { EGNPI: timpa(s.larik.EGNPI, h.egnpi), TotalEgnpiAmountNP: keBaris(h.TotalEgnpiAmountNP) },
      medan: { TotalEgnpiAmount: h.TotalEgnpiAmount, TotalEgnpiProportion: h.TotalEgnpiProportion },
      pesan: h.pesan,
    }
  }

/**
 * `SetAmountConversion` — Section rincian `DetailEGNPI` (`Amount`,
 * `Currency`), atas baris EGNPI ke-n. Hanya larik `EGNPI` yang diambil:
 * aksi `konversi` tidak menghitung total.
 */
const konversiEgnpi: Rumus = async (s, _a, l) => {
  const akar = akarDari(s, l)
  const h = await kirim<{ egnpi: Baris[]; pesan: string[] }>('egnpi', {
    aksi: 'konversi',
    egnpi: akar.larik.EGNPI ?? [],
    kurs: kurs(akar),
    retensi: retensiUntukEgnpi(akar),
    indeks: indeksAkar(l, 'EGNPI'),
  })
  return { ...kosong(), pesan: h.pesan, akar: { larik: { EGNPI: timpaIsi(akar.larik.EGNPI, h.egnpi) } } }
}

/**
 * `TreatyInSetReport(startdate, autocalculate)` — tombol Apply tab
 * Reporting Period (`TreatyIn.ReportingStart`, `true`) dan sel Initial Date
 * (`.InitialDate`, `.AutoCalculate` barisnya). Langkah 1 KELUAR bila
 * `autocalculate` bukan `true`. Rute Treaty In menjawab dalam bentuk layar
 * Treaty In; dipetakan balik ke kunci Pega larik `ReportingPeriodList`.
 */
const periodePelaporan: Rumus = async (s, a, l): Promise<HasilRumus> => {
  if (nilaiParam(a.param?.autocalculate, s, l) !== 'true') return kosong()
  const m = s.medan
  const h = await kirim<{
    baris: {
      periode: string
      hitungOtomatis: string
      tanggalAwalAsli: string
      jatuhTempoKirimAsli: string
      jatuhTempoKonfirmasiAsli: string
      jatuhTempoBayarAsli: string
    }[]
    galat: Record<string, string>
  }>('periode-pelaporan', {
    mulai: keSimpanTanggal(m.ReportingStart),
    akhir: keSimpanTanggal(m.ReportingEnd),
    periode: m.ReportingPeriod ?? '',
    interval: m.ReportingInterval ?? '',
    penyerahan: m.ReportingSubmission ?? '',
    konfirmasi: m.ReportingConfirmation ?? '',
    pelunasan: m.ReportingSettlement ?? '',
    awal: keSimpanTanggal(nilaiParam(a.param?.startdate, s, l)),
  })
  const pesan = Object.values(h.galat)
  if (pesan.length > 0) return { larik: {}, medan: {}, pesan }
  return {
    larik: {
      ReportingPeriodList: h.baris.map((b) => ({
        Period: b.periode,
        AutoCalculate: b.hitungOtomatis,
        InitialDate: b.tanggalAwalAsli,
        SubmissionDue: b.jatuhTempoKirimAsli,
        ConfirmationDue: b.jatuhTempoKonfirmasiAsli,
        SettlementDue: b.jatuhTempoBayarAsli,
      })),
    },
    medan: {},
    pesan: [],
  }
}

/**
 * `setValue` — pasangan `TreatyIn.X = nilai` (ke AKAR) atau `.X = nilai`
 * (ke halaman pemicu), dijalankan di layar.
 */
const setel: Rumus = (_s, a) => {
  const akar: Record<string, string> = {}
  const halaman: Record<string, string> = {}
  for (const [k, v] of Object.entries(a.param ?? {})) {
    const nilai = v.replace(/^"(.*)"$/, '$1')
    if (k.startsWith('TreatyIn.')) akar[k.slice('TreatyIn.'.length)] = nilai
    else halaman[k.replace(/^\./, '')] = nilai
  }
  return Promise.resolve({ larik: {}, medan: halaman, pesan: [], akar: { medan: akar } })
}

// ---------------------------------------------------------------------
// Accumulation (Prop) — rute `akumulasi` (`TreatyInSetAccountReport`,
// `TreatyInAccumulationSetSubDue`; `treatyin/backend/services/hitung_akumulasi.go`).
// ---------------------------------------------------------------------

/**
 * `periode` = isian Period (didahului DT `TreatyInDeleteAccumulationLists`,
 * jadi daftar yang dikirim sudah KOSONG); `jatuh-tempo` = sel Reporting
 * Date / Submission Days. Tanggal dikirim bentuk simpan.
 */
const akumulasi =
  (aksi: 'periode' | 'jatuh-tempo'): Rumus =>
  async (s) => {
    const daftar = s.larik.AccumulationList ?? []
    const h = await kirim<{ AccumulationList: Baris[] }>('akumulasi', {
      aksi,
      AccumulationPeriod: s.medan.AccumulationPeriod ?? '',
      ReportingStart: keSimpanTanggal(s.medan.ReportingStart),
      ReportingEnd: keSimpanTanggal(s.medan.ReportingEnd),
      AccumulationList: daftar.map((b) => ({
        Period: teks(b.Period),
        ReportDate: keSimpanTanggal(teks(b.ReportDate)),
        SubDays: teks(b.SubDays),
        SubDueDate: keSimpanTanggal(teks(b.SubDueDate)),
      })),
    })
    return { larik: { AccumulationList: timpa(daftar, h.AccumulationList) }, medan: {}, pesan: [] }
  }

// ---------------------------------------------------------------------
// Share Prop — rute `share-prop` (TreatyInPropshare, TreatyInPropshareDetail,
// FetchQSfromMaster, SetSpreadName; `hitung_share_prop.go`).
// ---------------------------------------------------------------------

interface HasilShareProp {
  Limits: Baris[]
  TotalShareRnmProp: NilaiMataUang[]
  TotalSpreadedRnmProp: NilaiMataUang[]
  TotalSpreadedRnmRIProp: NilaiMataUang[]
  pesan: string[]
}

/**
 * Satu aksi tab Share Prop. Detail sasaran = jalur rincian
 * (`Limits(i).Detail(j)` — Section `DetailShare` dan `DetailLimits`).
 * `spreading` membawa `ParentReinsTypeID` dari parameternya:
 * `.SpreadingTypeID` (isian Spreading Type) atau `""` (rantai Treaty Group
 * rincian Limits). Hasilnya ke AKAR: `Limits` dan ketiga total.
 */
const shareProp =
  (aksi: 'share' | 'detail' | 'spreading' | 'sebar-nama'): Rumus =>
  async (s, a, l) => {
    const akar = akarDari(s, l)
    const m = akar.medan
    const j = l.jalur ?? []
    const badan: Record<string, unknown> = {
      aksi,
      Limits: akar.larik.Limits ?? [],
      RNMShareP: m.RNMShareP ?? '',
      BrokeragePercentP: m.BrokeragePercentP ?? '',
      OptionLimit: m.OptionLimit ?? '',
      RNMShareAcrossTheBoard: m.RNMShareAcrossTheBoard ?? '',
      Commencement: keSimpanTanggal(m.Commencement),
      TotalShareRnmProp: akar.larik.TotalShareRnmProp ?? [],
      TotalSpreadedRnmProp: akar.larik.TotalSpreadedRnmProp ?? [],
      TotalSpreadedRnmRIProp: akar.larik.TotalSpreadedRnmRIProp ?? [],
      indeksLimit: j[0]?.indeks ?? 0,
      indeksDetail: j[1]?.indeks ?? 0,
    }
    if (aksi === 'spreading') badan.ParentReinsTypeID = nilaiParam(a.param?.ParentReinsTypeID, s, l)
    const h = await kirim<HasilShareProp>('share-prop', badan)
    return {
      ...kosong(),
      pesan: h.pesan,
      akar: {
        larik: {
          Limits: h.Limits,
          TotalShareRnmProp: keBaris(h.TotalShareRnmProp),
          TotalSpreadedRnmProp: keBaris(h.TotalSpreadedRnmProp),
          TotalSpreadedRnmRIProp: keBaris(h.TotalSpreadedRnmRIProp),
        },
      },
    }
  }

// ---------------------------------------------------------------------
// Limits Prop — rincian `DetailLimits` (rute `limit`, `limit-deduksi`,
// `limit-cadangan`) dan `LimitProportional`.
// ---------------------------------------------------------------------

/**
 * Medan yang `LimitCalculation` TULIS, per mode — hanya ini yang digabung
 * kembali ke Detail. SALINAN `DITULIS_LIMIT_CALCULATION` Treaty In (nol
 * impor antarmodul).
 */
const DITULIS_LIMIT_CALCULATION = {
  // [2]–[4]
  qs: ['RetentionPct', 'CessionPct', 'RetentionList', 'CessionList'],
  // [6]–[14]
  surplus: ['IOOPct', 'CessionPct', 'IOOLimitList', 'RetentionList', 'CessionList', 'COBList'],
} as const

/**
 * `LimitCalculation(kindoftreaty, add, autocalculate)` atas Detail ini.
 * Langkah 1 KELUAR bila `autocalculate` bukan `true` — parameternya
 * `true`, atau `.Layer` baris sel pemicu (kotak centang hitung otomatis).
 *
 * ⚠️ Sel 100% Limit mengirim `QS` HURUF BESAR ke perbandingan `=="qs"`;
 * dibaca `qs` — keputusan yang sama dengan Treaty In (`TabLimitsProp`).
 * Pohon = SELURUH Detail di seluruh Limits (langkah 9), dengan Detail ini
 * dalam keadaan TERKINI.
 */
const limitCalculation: Rumus = async (s, a, l) => {
  if (nilaiParam(a.param?.autocalculate, s, l) !== 'true') return kosong()
  const mode = nilaiParam(a.param?.kindoftreaty, s, l).toLowerCase() === 'qs' ? 'qs' : 'surplus'
  const akar = akarDari(s, l)
  const h = await kirim<Baris>('limit', {
    jenis: mode,
    tambah: nilaiParam(a.param?.add, s, l) === 'man' ? 'man' : '',
    otomatis: true,
    detail: sisiKeBaris(s),
    pohon: (akar.larik.Limits ?? []).flatMap((x) => anak(x, 'Detail')),
  })
  const larik: Record<string, Baris[]> = {}
  const medan: Record<string, string> = {}
  for (const k of DITULIS_LIMIT_CALCULATION[mode]) {
    const v = h[k]
    if (typeof v === 'string') medan[k] = v
    else if (Array.isArray(v)) larik[k] = v
  }
  return { larik, medan, pesan: [] }
}

/** `CalculateDeduction(sts, index)` di Detail Limits — rute `limit-deduksi`. */
const deduksiLimit: Rumus = async (s, a, l) => {
  const nomor = Number(nilaiParam(a.param?.index, s, l))
  const h = await kirim<{ DeductionList: Baris[]; DeductionTotalList: Baris[]; pesan: string[] }>('limit-deduksi', {
    sts: nilaiParam(a.param?.sts, s, l),
    indeks: Number.isFinite(nomor) && nomor > 0 ? nomor - 1 : 0,
    DeductionList: s.larik.DeductionList ?? [],
    GrossPremiumList: s.larik.GrossPremiumList ?? [],
  })
  return { larik: { DeductionList: h.DeductionList, DeductionTotalList: h.DeductionTotalList }, medan: {}, pesan: h.pesan }
}

/** `PremiumReserveCalculate` — `% Premium Reserve` × Cession per mata uang. */
const cadanganPremi: Rumus = async (s) => {
  const h = await kirim<{ ReserveList: Baris[]; pesan: string[] }>('limit-cadangan', {
    PremiumReservePct: s.medan.PremiumReservePct ?? '',
    CessionList: s.larik.CessionList ?? [],
  })
  return { larik: { ReserveList: h.ReserveList }, medan: {}, pesan: h.pesan }
}

const masterBelum = (): HasilRumus => ({ ...kosong(), pesan: [PENYESUAIAN.masterBelum] })

/**
 * `SetCurrName_Act` — pasangan mata uang baris sel pemicu.
 *
 * ⚠️ KEPUTUSAN. Activity-nya (korpus Adjustment) menulis
 * `MstBankAccount.CURRENCY` bila `Param.CURRID == .ID` — halaman LAIN, dan
 * `CURRID` tidak pernah dikirim sel mana pun: atas baris itu sendiri ia tidak
 * mengubah apa pun. Data tersimpan Pega tetap membawa `Currency` DAN
 * `CurrencyID` berpasangan, jadi pasangannya diisi dropdown lewat medan
 * tambahan yang tidak diekspor. Di sini pasangan itu diisi dari daftar
 * `BrowseCurrency…_RD` — sama dengan Treaty In (`PilihKode`).
 */
const pasangMataUang: Rumus = (s, _a, l) => {
  const sel = l.sel
  if (sel === undefined) return Promise.resolve(kosong())
  const daftar = l.master?.mataUang
  if (daftar === undefined) return Promise.resolve(masterBelum())
  const baris = s.larik[sel.larik] ?? []
  const b = baris[sel.indeks]
  if (b === undefined) return Promise.resolve(kosong())
  const baru: Baris =
    sel.kunci === 'CurrencyID'
      ? { ...b, Currency: daftar.find((x) => x.id === teks(b.CurrencyID))?.nama ?? '' }
      : { ...b, CurrencyID: daftar.find((x) => x.nama === teks(b.Currency))?.id ?? '' }
  return Promise.resolve({ larik: { [sel.larik]: baris.map((y, i) => (i === sel.indeks ? baru : y)) }, medan: {}, pesan: [] })
}

/**
 * `SetCurrNameMasterTreaty_Act` — perilaku `change` sel `.CurrencyID` grid
 * Rate of Exchange panel New (`GRID_KURS.baru`: postValue → runActivity →
 * refresh). *(8 Okt, E)*
 *
 * Korpus `Treaty In Adjustment/Activity/SetCurrNameMasterTreaty_Act.xml`
 * (langkahnya identik dengan korpus Treaty In), kelas
 * `ASM-FW-GISFW-Data-TreatyInCurrencyList` — halaman pemicu = BARIS kurs:
 *
 *   [1] Page-Remove `CURRENCY`
 *   [2] Obj-Browse `CURRENCY` kelas `ASM-FW-GISFW-Int-CURRENCY`,
 *       `.ID = .CurrencyID` (@4723/@4763), pilih `.Currency` (@9916),
 *       MaxRecords 1 (@11461)
 *   [3] Property-Set `.Currency = CURRENCY.pxResults(1).Currency` (@18421/@18463)
 *
 * Hanya `.Currency` baris itu yang ditulis; pengenal yang tidak ditemukan
 * memberi `pxResults(1)` kosong → `.Currency` kosong. Berbeda dari
 * `SetCurrName_Act` (pasangan dua arah) — Activity ini SATU arah, ID → nama.
 *
 * ⚠️ Daftar master = rute Treaty In (`BrowseCurrency_RD`, tanpa `ITL`),
 * sedang Obj-Browse tidak menyaring `ITL`. Dropdown sel ini memakai daftar
 * yang sama, jadi pengenal `ITL` tidak dapat DIPILIH; bedanya hanya muncul
 * pada nilai tersimpan yang tidak pernah diubah — dan itu tidak memicu
 * `change`.
 */
const namaMataUangMaster: Rumus = (s, _a, l) => {
  const sel = l.sel
  if (sel === undefined) return Promise.resolve(kosong())
  const daftar = l.master?.mataUang
  if (daftar === undefined) return Promise.resolve(masterBelum())
  const baris = s.larik[sel.larik] ?? []
  const b = baris[sel.indeks]
  if (b === undefined) return Promise.resolve(kosong())
  const baru: Baris = { ...b, Currency: daftar.find((x) => x.id === teks(b.CurrencyID))?.nama ?? '' }
  return Promise.resolve({ larik: { [sel.larik]: baris.map((y, i) => (i === sel.indeks ? baru : y)) }, medan: {}, pesan: [] })
}

/**
 * `SetTreatyGroupName_Act` — `.TreatyGroup = TREATYGROUP(.TreatyGroupID).TreatyGroupName`.
 * Langkah 2 (pencarian) dilewati bila pengenalnya kosong, tetapi langkah 3
 * tetap menulis — hasil kosong.
 */
const namaKelompokTreaty: Rumus = (s, _a, l) => {
  const id = s.medan.TreatyGroupID ?? ''
  if (id === '') return Promise.resolve({ larik: {}, medan: { TreatyGroup: '' }, pesan: [] })
  const daftar = l.master?.kelompokTreaty
  if (daftar === undefined) return Promise.resolve(masterBelum())
  return Promise.resolve({ larik: {}, medan: { TreatyGroup: daftar.find((x) => x.id === id)?.nama ?? '' }, pesan: [] })
}

/**
 * Kind of Treaty dari baris REINSURANCETYPE — kolom SOA Name, atau Name
 * bila kosong. ⭐ KEPUTUSAN PEMAKAI 6 Oktober 2026 (Treaty In
 * `kindOfTreatyDari`, `PERTANYAAN-TERBUKA-LAYAR-PEGA.md` §20); menyimpang
 * dari langkah 3 `SetTreatyTypeName_Act` (`.Note`). Disalin, nol impor.
 */
const kindOfTreaty = (x: PilihanTreatyIn) => {
  const soa = (x.namaSoa ?? '').trim()
  return soa !== '' ? soa : x.nama
}

/**
 * DT `TreatyTypeSetIndex` atas halaman Limits ini:
 *   1    `.ID = .pxListSubscript`
 *   2    FOR_EACH `.Detail`: `.TreatyType = .TreatyType` (jenis INDUK)
 */
const setIndeksJenis = (s: SisiPenyesuaian, l: LingkupRumus): HasilRumus => {
  const jenis = s.medan.TreatyType ?? ''
  const detail = s.larik.Detail
  return {
    larik: detail === undefined ? {} : { Detail: detail.map((d) => ({ ...d, TreatyType: jenis })) },
    medan: { ID: String(nomorHalaman(l)) },
    pesan: [],
  }
}

/** `SetTreatyTypeName_Act` — Kind of Treaty dari Treaty Type, lalu DT `TreatyTypeSetIndex` (langkah 4). */
const namaJenisTreaty: Rumus = (s, _a, l) => {
  const id = s.medan.TreatyTypeID ?? ''
  let jenis = ''
  if (id !== '') {
    const daftar = l.master?.jenisTreaty
    if (daftar === undefined) return Promise.resolve(masterBelum())
    const x = daftar.find((o) => o.id === id)
    jenis = x === undefined ? '' : kindOfTreaty(x)
  }
  const h = setIndeksJenis({ ...s, medan: { ...s.medan, TreatyType: jenis } }, l)
  return Promise.resolve({ ...h, medan: { ...h.medan, TreatyType: jenis } })
}

// ---------------------------------------------------------------------
// Limits Non-Prop — rute `limit-np` (Activity TreatyInNPSetTotal(limits),
// TreatyInSummaryLimit, TreatyInLimitsListValue; identik di kedua korpus).
// ---------------------------------------------------------------------

interface HasilLimitNP {
  layers: Baris[]
  Total?: Record<string, Baris[]>
  TotalLimitsROL?: string
  LimitSummaryList?: Baris[]
  pesan: string[]
}

/** Baris EGNPI yang `TotalEgnpi` baca. */
const egnpiUntukLimit = (s: SisiPenyesuaian) =>
  (s.larik.EGNPI ?? []).map((e) => ({
    TreatyGroup: teks(e.TreatyGroup),
    Currency: teks(e.Currency),
    CurrencyID: teks(e.CurrencyID),
    Amount: teks(e.Amount),
  }))

/**
 * `total` = `NPSetTotal(limits)` + `SummaryLimit` (rute Treaty In,
 * `case AksiNPTotal`); `nilai-list` = `TreatyInLimitsListValue`. Layer
 * (beserta larik anaknya) DITIMPAKAN per indeks; total, ringkasan, dan
 * `TotalLimitsROL` diganti.
 */
const limitNP =
  (aksi: 'total' | 'nilai-list'): Rumus =>
  async (s): Promise<HasilRumus> => {
    const h = await kirim<HasilLimitNP>('limit-np', {
      aksi,
      layers: s.larik.Limits ?? [],
      egnpi: egnpiUntukLimit(s),
      kurs: kurs(s),
      indeks: 0,
      baris: 0,
    })
    const larik: Record<string, Baris[]> = { Limits: timpa(s.larik.Limits, h.layers) }
    if (h.LimitSummaryList !== undefined) larik.LimitSummaryList = h.LimitSummaryList
    for (const [nama, baris] of Object.entries(h.Total ?? {})) larik[nama] = baris
    return {
      larik,
      medan: h.TotalLimitsROL === undefined ? {} : { TotalLimitsROL: h.TotalLimitsROL },
      pesan: h.pesan,
    }
  }

/** Aksi rute `limit-np` atas SATU layer (atau, `egnpi`, seluruh layer). */
type AksiLayerNP = 'egnpi' | 'reinstatement' | 'reinst-jumlah' | 'reinst-persen' | 'reinst-tambahan' | 'adj' | 'mdp'

/**
 * Medan layer yang tiap aksi TULIS — hanya ini yang digabung kembali.
 * SALINAN `DITULIS_LIMIT_NP` Treaty In (nol impor antarmodul).
 */
const DITULIS_LIMIT_NP: Readonly<Record<AksiLayerNP, readonly string[]>> = {
  egnpi: ['EgnpiTotalList'],
  reinstatement: ['Reinstatement_List'],
  'reinst-jumlah': ['Reinstatement_List'],
  'reinst-persen': ['Reinstatement_List'],
  'reinst-tambahan': ['Reinstatement_List'],
  adj: ['PremiumEarnedList', 'ROLPct'],
  // ⭐ 8 Oktober 2026 — sama dengan Treaty In: sesudah MDP berubah, rute
  // `limit-np` menyegarkan Reinstatement Premium Amount dengan rumus
  // `ReCalculateReinstatement` yang SAMA (`segarkanReinstatement`).
  mdp: ['MDPList', 'MDPMinList', 'ROLPct', 'Reinstatement_List'],
}

/**
 * Rumus Section rincian `Layers` (layer = `Limits(n)` jalur): DetailCalculation
 * (`adj`/`mdp`), SetReinstatementPct, DT reinstatement per baris (`baris` =
 * sel pemicu), dan `TotalEgnpi` — yang terakhir menapaki SEMUA layer.
 */
const limitNPLayer =
  (aksi: AksiLayerNP): Rumus =>
  async (s, _a, l) => {
    const akar = akarDari(s, l)
    // Rincian `LayersEDM` (cabang Adjust Premium) berjalan di baris
    // `ActualValue.Limits(n)`; DetailCalculation/SetReinstatementPct menulis
    // halaman itu. `TotalEgnpi` selalu menapaki `TreatyIn.Limits`.
    const larikJalur = (l.jalur ?? [])[0]?.larik ?? ''
    const larikLayer = aksi !== 'egnpi' && larikJalur.endsWith('Limits') ? larikJalur : 'Limits'
    const layers = akar.larik[larikLayer] ?? []
    const indeks = indeksAkar(l, larikLayer)
    const h = await kirim<HasilLimitNP>('limit-np', {
      aksi,
      layers,
      egnpi: egnpiUntukLimit(akar),
      kurs: kurs(akar),
      indeks,
      baris: l.sel?.indeks ?? 0,
    })
    const baru = layers.map((x, i) => {
      const hb = h.layers[i]
      if (hb === undefined || (aksi !== 'egnpi' && i !== indeks)) return x
      const y: Baris = { ...x }
      for (const k of DITULIS_LIMIT_NP[aksi]) {
        const v = hb[k]
        if (v !== undefined) y[k] = v
      }
      return y
    })
    return { ...kosong(), pesan: h.pesan, akar: { larik: { [larikLayer]: baru } } }
  }

/**
 * `Activity/TreatyInSummaryMDP` — disalin langkah demi langkah:
 *
 *   per layer (kunci LayerType·Layer·LayerPartType·LayerPart):
 *     3.2  flag = 1 bila TempMDP sudah punya baris berkunci sama
 *     per MDPList:
 *       3.3.2  `local.currency != .Currency` atas baris TempMDP — baris itu
 *              TIDAK punya `.Currency`, jadi flag hanya berubah bila mata
 *              uang MDP-nya KOSONG (dan TempMDP berisi). Disalin apa adanya.
 *       3.3.3  flag 0 → baris baru: IDR/USD = nilai bila mata uangnya itu, 0
 *       3.3.4  flag 1 → baris berkunci sama: IDR/USD += nilai
 *       3.3.5  flag = 1
 *   4/5  MDPSummaryList = TempMDP; Note = LayerType+Layer+" of "+LayerPartType+LayerPart
 *
 * Mata uang selain IDR/USD tidak dijumlahkan ke mana pun — bunyi Activity-nya.
 */
export function ringkasanMDP(layers: readonly Baris[]): Baris[] {
  const kunci = ['LayerType', 'Layer', 'LayerPartType', 'LayerPart'] as const
  const sama = (a: Baris, b: Baris) => kunci.every((k) => teks(a[k]) === teks(b[k]))
  const temp: Baris[] = []
  for (const l of layers) {
    let flag = temp.some((r) => sama(r, l)) ? 1 : 0
    for (const m of anak(l, 'MDPList')) {
      const cur = teks(m.Currency)
      const val = teks(m.Value)
      if (cur === '' && temp.length > 0) flag = 1
      if (flag === 0) {
        temp.push({
          LayerType: teks(l.LayerType),
          Layer: teks(l.Layer),
          LayerPartType: teks(l.LayerPartType),
          LayerPart: teks(l.LayerPart),
          IDR: cur === 'IDR' ? val : '0',
          USD: cur === 'USD' ? val : '0',
        })
      } else {
        for (const r of temp) {
          if (!sama(r, l)) continue
          if (cur === 'IDR') r.IDR = tambah(teks(r.IDR), val)
          if (cur === 'USD') r.USD = tambah(teks(r.USD), val)
        }
      }
      flag = 1
    }
  }
  return temp.map((r) => ({
    ...r,
    Note: `${teks(r.LayerType)}${teks(r.Layer)} of ${teks(r.LayerPartType)}${teks(r.LayerPart)}`,
  }))
}

const summaryMDP: Rumus = (s) =>
  Promise.resolve({ larik: { MDPSummaryList: ringkasanMDP(s.larik.Limits ?? []) }, medan: {}, pesan: [] })

// ---------------------------------------------------------------------
// Share Non-Prop — rute `share-np` (Activity TreatyInNPSetTotal(share),
// TreatyInSummaryLimitShare/FacShare, TreatyInShareListValue,
// TreatyInXOLAddSpreading, TreatyInSetBrokerage, TreatyInNonAddItem(share),
// dan rincian `Share`: TreatyInXOLAddSpreadingDetail, FetchQSfromMasterXOL,
// SetSpreadingXOL, CalculateDeduction).
// ---------------------------------------------------------------------

/** Skalar akar tab Share. */
const SKALAR_SHARE = [
  'RNMShare', 'BrokeragePercent', 'RNMShareAcrossTheBoard', 'FacultativeShare',
  'FacultativeShareBrokerage', 'RnmShareDeducted', 'IsProRate',
] as const

/** Larik tab Share yang ikut dikirim dan diterima kembali. */
const LARIK_SHARE = [
  'ShareReins', 'ShareFacultativeReinsurers', 'Share', 'FacultativeShareList',
  'LimitShareSummaryList', 'LimitFacShareSummaryList',
] as const

/** Kesembilan grid `Total All Layers RNM Share`. */
const TOTAL_SHARE = [
  'TotalShareRnmNP', 'TotalShareGrossNP', 'TotalShareGrossMinNP', 'TotalShareDeductionNP', 'TotalShareNetNP',
  'TotalFacShareRnmNP', 'TotalFacShareGrossNP', 'TotalFacShareDeductionNP', 'TotalFacShareNetNP',
] as const

type ShareNP = Record<string, string | Baris[] | Record<string, Baris[]>>

/** Aksi rute `share-np` — satu per Activity (`AksiShare…` di services Treaty In). */
type AksiShareNP =
  | 'total'
  | 'nilai-share'
  | 'rnm'
  | 'set-brokerage'
  | 'summary'
  | 'spreading-type'
  | 'rnm-baris'
  | 'deduksi'
  | 'spreading-pct'

/**
 * Satu aksi rute `share-np` atas AKAR. Baris Share sasaran = jalur rincian
 * (`Share(n)`); baris deduksi = sel pemicu; `sts` dari parameter. Larik Share
 * DIGANTI utuh: rute itu menyusunnya ulang (bukan menimpa per indeks).
 * Pesan yang Pega tempelkan pada baris Share ini ikut ditampilkan.
 */
const shareNP =
  (aksi: AksiShareNP): Rumus =>
  async (s, a, l): Promise<HasilRumus> => {
    const akar = akarDari(s, l)
    const share: ShareNP = {}
    for (const k of SKALAR_SHARE) share[k] = akar.medan[k] ?? ''
    for (const k of LARIK_SHARE) share[k] = akar.larik[k] ?? []
    share.Total = Object.fromEntries(TOTAL_SHARE.map((k) => [k, akar.larik[k] ?? []]))
    // Baris Share = `Param.idx` / `indexShare` (`.pxListSubscript`) bila ada —
    // rincian yang sama juga dibuka dari `ActualValue.Share(n)`, dan Activity
    // tetap menulis `TreatyIn.Share(idx)`; selain itu jalur `Share(n)`.
    const nomor = Number(nilaiParam(a.param?.idx ?? a.param?.indexShare, s, l))
    const indeks = Number.isFinite(nomor) && nomor > 0 ? nomor - 1 : indeksAkar(l, 'Share')
    const h = await kirim<{ share: ShareNP; pesan: string[]; pesanBaris?: { indeks: number; pesan: string }[] }>('share-np', {
      aksi,
      share,
      layers: akar.larik.Limits ?? [],
      idKontrak: akar.medan.ID ?? '',
      commencement: keSimpanTanggal(akar.medan.Commencement),
      indeks,
      baris: l.sel?.indeks ?? 0,
      sts: nilaiParam(a.param?.sts, s, l),
    })
    const larik: Record<string, Baris[]> = {}
    const medan: Record<string, string> = {}
    for (const k of SKALAR_SHARE) {
      const v = h.share[k]
      if (typeof v === 'string') medan[k] = v
    }
    for (const k of LARIK_SHARE) {
      const v = h.share[k]
      if (Array.isArray(v)) larik[k] = v
    }
    const total = h.share.Total
    if (total !== undefined && !Array.isArray(total) && typeof total !== 'string') {
      for (const [k, v] of Object.entries(total)) larik[k] = v
    }
    const diRincian = (l.jalur ?? []).length > 0
    const pesanBaris = (h.pesanBaris ?? []).filter((p) => !diRincian || p.indeks === indeks).map((p) => p.pesan)
    return { ...kosong(), pesan: [...h.pesan, ...pesanBaris], akar: { larik, medan } }
  }

/**
 * `CalculateDeduction(sts, index)` — SATU Activity, dua tempat: rincian
 * `Share` Non-Prop (rute `share-np` `deduksi`, baris Share = jalur) dan
 * rincian `DetailLimits` Prop (rute `limit-deduksi`).
 */
const calculateDeduction: Rumus = (s, a, l) => {
  const j = l.jalur ?? []
  return j[j.length - 1]?.larik === 'Share' ? shareNP('deduksi')(s, a, l) : deduksiLimit(s, a, l)
}

// ---------------------------------------------------------------------
// Installment — rute `angsuran` (TreatyInSetValueInstallment,
// SetTotalInstallment, TreatyInNPSetTotal(installment); identik di kedua
// korpus). Lihat `treatyin/backend/services/hitung_angsuran.go`.
// ---------------------------------------------------------------------

interface AngsuranRute {
  Currency: string
  AmountTotal: string
  PctTotal: string
  InstallmentList: Record<string, string>[]
}

interface HasilAngsuran {
  angsuran: AngsuranRute[]
  /** `null` bila aksinya tidak menyentuh total (`total-baris`). */
  TotalInstallmentNP: NilaiMataUang[] | null
  InstallmentNo: string
  pesan: string[]
}

const KUNCI_BARIS_ANGSURAN = ['Installment', 'DueDate', 'WPC', 'PaymentDate', 'Currency', 'InstallmentPct', 'Amount'] as const

/** Satu halaman `Installment` dalam bentuk rute. */
const angsuranRute = (a: Baris): AngsuranRute => ({
  Currency: teks(a.Currency),
  AmountTotal: teks(a.AmountTotal),
  PctTotal: teks(a.PctTotal),
  InstallmentList: anak(a, 'InstallmentList').map((r) => Object.fromEntries(KUNCI_BARIS_ANGSURAN.map((k) => [k, teks(r[k])]))),
})

/** `TreatyIn.TotalShareNetNP` — sumber nilai angsuran. */
const netPremium = (s: SisiPenyesuaian) =>
  (s.larik.TotalShareNetNP ?? []).map((n) => ({ Currency: teks(n.Currency), CurrencyID: teks(n.CurrencyID), Value: teks(n.Value) }))

/**
 * `TreatyInSetValueInstallment` — isian `Installment` (perilaku `change`,
 * tanpa `status`) dan tombol `Update Value` (`status=update`).
 * `Param.Installment` = `TreatyIn.InstallmentNo`, dibaca dari halaman ini.
 * ⛔ Larik `Installment` DIGANTI utuh — langkah 2 membuangnya.
 */
const nilaiAngsuran: Rumus = async (s, a, l) => {
  const h = await kirim<HasilAngsuran>('angsuran', {
    aksi: 'nilai',
    status: a.param?.status ?? '',
    installmentNo: s.medan.InstallmentNo ?? '',
    edmState: s.medan.EDMState ?? '',
    angsuran: (s.larik.Installment ?? []).map(angsuranRute),
    angsuranLama: (l.lama?.larik.Installment ?? []).map(angsuranRute),
    netPremium: netPremium(s),
    indeks: 0,
  })
  return {
    larik: {
      Installment: h.angsuran.map((x) => ({ ...x, InstallmentList: x.InstallmentList.map((r) => ({ ...r })) })),
      TotalInstallmentNP: keBaris(h.TotalInstallmentNP ?? []),
    },
    medan: { InstallmentNo: h.InstallmentNo },
    pesan: h.pesan,
  }
}

/**
 * `SetTotalInstallment` — berjalan di halaman SATU `Installment` (Section
 * rincian `Installments`): sel `% Installment` (`status=editpercentage`,
 * Amount dihitung ulang) dan sel `Amount` (hanya total). Net Premium dibaca
 * dari akar panel. ⛔ Yang diubah hanya `Amount` tiap baris, `AmountTotal`,
 * dan `PctTotal` — kunci lain baris tetap.
 */
const totalHalamanAngsuran: Rumus = async (s, a, l): Promise<HasilRumus> => {
  const halaman: Baris = { ...s.medan, InstallmentList: s.larik.InstallmentList ?? [] }
  const h = await kirim<HasilAngsuran>('angsuran', {
    aksi: 'total-baris',
    status: a.param?.status ?? '',
    angsuran: [angsuranRute(halaman)],
    netPremium: netPremium(l.panel ?? s),
    indeks: 0,
  })
  const hasil = h.angsuran[0]
  if (hasil === undefined) return { larik: {}, medan: {}, pesan: h.pesan }
  return {
    larik: {
      InstallmentList: (s.larik.InstallmentList ?? []).map((r, i) => ({ ...r, Amount: hasil.InstallmentList[i]?.Amount ?? teks(r.Amount) })),
    },
    medan: { AmountTotal: hasil.AmountTotal, PctTotal: hasil.PctTotal },
    pesan: h.pesan,
  }
}

/**
 * `TreatyInUpdatePaymentDate` (sel Due Date, halaman Installment ini) dan
 * `_Act` (sel WPC, SEMUA halaman `TreatyIn.Installment`) — rute `angsuran`.
 * Hanya `PaymentDate` yang diambil; tanggal dikirim bentuk simpan.
 */
const tanggalBayar =
  (semua: boolean): Rumus =>
  async (s, _a, l) => {
    const akar = akarDari(s, l)
    const bentuk = (a: Baris) => {
      const r = angsuranRute(a)
      return { ...r, InstallmentList: r.InstallmentList.map((b) => ({ ...b, DueDate: keSimpanTanggal(b.DueDate) })) }
    }
    if (!semua) {
      const halaman: Baris = { ...s.medan, InstallmentList: s.larik.InstallmentList ?? [] }
      const h = await kirim<HasilAngsuran>('angsuran', { aksi: 'tanggal-bayar', angsuran: [bentuk(halaman)], indeks: 0 })
      const hasil = h.angsuran[0]
      if (hasil === undefined) return { ...kosong(), pesan: h.pesan }
      return {
        larik: { InstallmentList: (s.larik.InstallmentList ?? []).map((r, i) => ({ ...r, PaymentDate: hasil.InstallmentList[i]?.PaymentDate ?? teks(r.PaymentDate) })) },
        medan: {},
        pesan: h.pesan,
      }
    }
    const daftar = akar.larik.Installment ?? []
    const h = await kirim<HasilAngsuran>('angsuran', { aksi: 'tanggal-bayar-semua', angsuran: daftar.map(bentuk), indeks: 0 })
    const baru = daftar.map((a, i) => {
      const hasil = h.angsuran[i]
      if (hasil === undefined) return a
      return { ...a, InstallmentList: anak(a, 'InstallmentList').map((r, j) => ({ ...r, PaymentDate: hasil.InstallmentList[j]?.PaymentDate ?? teks(r.PaymentDate) })) }
    })
    return { ...kosong(), pesan: h.pesan, akar: { larik: { Installment: baru } } }
  }

/** Tombol `Update Total` tab Installment — `TreatyInNPSetTotal(type=installment)` langkah 27–28. */
const totalAngsuran: Rumus = async (s) => {
  const h = await kirim<HasilAngsuran>('angsuran', { aksi: 'total', angsuran: (s.larik.Installment ?? []).map(angsuranRute), indeks: 0 })
  return { larik: { TotalInstallmentNP: keBaris(h.TotalInstallmentNP ?? []) }, medan: {}, pesan: h.pesan }
}

// ---------------------------------------------------------------------
// Value Difference — rute `selisih` (`TreatyEDMCalculateDifference` dan
// kesembilan Activity panggilannya, identik di kedua korpus;
// `treatyin/backend/services/hitung_selisih.go`).
// ---------------------------------------------------------------------

/** Halaman bertitik di akar panel: `ValueDifference.X`, `ValueBeforeProrate.X`, `ActualValue.X`. */
const SELISIH = 'ValueDifference.'
const SEBELUM_PRORATA = 'ValueBeforeProrate.'
const AKTUAL = 'ActualValue.'

interface HalamanPohon {
  medan: Record<string, string>
  larik: Record<string, Baris[]>
}

/** Halaman bertitik → pohon (awalan dibuang). */
function halamanBertitik(akar: SisiPenyesuaian, awalan: string): HalamanPohon {
  const medan: Record<string, string> = {}
  const larik: Record<string, Baris[]> = {}
  for (const [k, v] of Object.entries(akar.medan)) if (k.startsWith(awalan)) medan[k.slice(awalan.length)] = v
  for (const [k, v] of Object.entries(akar.larik)) if (k.startsWith(awalan)) larik[k.slice(awalan.length)] = v
  return { medan, larik }
}

/**
 * Pohon → kunci bertitik, MENGGANTI halaman itu utuh: kunci bertitik lama
 * yang tidak ada di jawaban dikosongkan (Property-Remove langkah 1).
 */
function gantiBertitik(akar: SisiPenyesuaian, awalan: string, h: HalamanPohon): Required<Ubahan> {
  const medan: Record<string, string> = {}
  const larik: Record<string, Baris[]> = {}
  for (const k of Object.keys(akar.medan)) if (k.startsWith(awalan)) medan[k] = ''
  for (const k of Object.keys(akar.larik)) if (k.startsWith(awalan)) larik[k] = []
  for (const [k, v] of Object.entries(h.medan)) medan[awalan + k] = v
  for (const [k, v] of Object.entries(h.larik)) larik[awalan + k] = v
  return { medan, larik }
}

/** Skalar dan larik yang rantai Value Difference baca — badan rute tetap kecil. */
const MEDAN_SELISIH = ['EDMState', 'RNMShare', 'BrokeragePercent', 'FacultativeShareBrokerage', 'IsProRate', 'ProRatePercent', 'TotalLimitsROL'] as const
const LARIK_SELISIH = [
  'EGNPI', 'Limits', 'LimitSummaryList', 'TotalLimitIOONP', 'TotalLimitDeductblNP', 'TotalLimitPremiEarnNP', 'TotalLimitMDPNP',
  'Share', 'LimitShareSummaryList', 'FacultativeShareList', 'Installment',
] as const

const pilihSelisih = (h: SisiPenyesuaian): HalamanPohon => ({
  medan: Object.fromEntries(MEDAN_SELISIH.map((k) => [k, h.medan[k] ?? ''])),
  larik: Object.fromEntries(LARIK_SELISIH.filter((k) => h.larik[k] !== undefined).map((k) => [k, h.larik[k] ?? []])),
})

/**
 * Tombol `Update Value` tab Value Difference. Selisih = New (akar) − Old
 * (`TreatyIn.OLDDATA`, sisi Old), per indeks. `ValueDifference` disusun
 * ULANG; `ValueBeforeProrate` hanya bila Pro Rate; `ActualValue` ringkasan
 * Share-nya; `FacultativeShareList` akar dihitung ulang (Deduction [5]).
 */
const selisihEDM: Rumus = async (s, _a, l) => {
  const akar = akarDari(s, l)
  const h = await kirim<{
    selisih: HalamanPohon
    sebelumProrata: HalamanPohon | null
    actual: HalamanPohon
    FacultativeShareList: Baris[] | null
    pesan: string[]
  }>('selisih', {
    aksi: 'edm',
    akar: pilihSelisih(akar),
    lama: pilihSelisih(l.lama ?? { medan: {}, larik: {} }),
    actual: halamanBertitik(akar, AKTUAL),
  })
  const vd = gantiBertitik(akar, SELISIH, h.selisih)
  const larik: Record<string, Baris[]> = { ...vd.larik }
  const medan: Record<string, string> = { ...vd.medan }
  if (h.sebelumProrata !== null) {
    const sb = gantiBertitik(akar, SEBELUM_PRORATA, h.sebelumProrata)
    Object.assign(larik, sb.larik)
    Object.assign(medan, sb.medan)
  }
  const ringkas = h.actual.larik.LimitShareSummaryList
  if (ringkas !== undefined) larik[`${AKTUAL}LimitShareSummaryList`] = ringkas
  if (h.FacultativeShareList !== null && akar.larik.FacultativeShareList !== undefined) {
    larik.FacultativeShareList = timpa(akar.larik.FacultativeShareList, h.FacultativeShareList)
  }
  return { ...kosong(), pesan: h.pesan, akar: { larik, medan } }
}

// ---------------------------------------------------------------------
// Achievement rincian Limits Prop — rute `achievement` (`GetAchievement`).
// Halaman SESI (`SearchData.CARI1/2`, `FlagExcel.CARI1`, `TempQuarter*`)
// hidup di akar panel berkunci utuh.
// ---------------------------------------------------------------------

interface HasilAchievement {
  limits: Baris[][]
  kuartal: string[]
  tahunKuartal: string[]
  flagExcel: boolean
}

/**
 * `GetAchievement` — tombol Refresh (tanpa parameter) dan isian Quarter Year
 * (`search = "search"`). Atas SELURUH `Limits`; hanya medan yang Activity
 * tulis yang digabung ke tiap Treaty Group — sama dengan Treaty In.
 */
const achievement: Rumus = async (s, a, l) => {
  const akar = akarDari(s, l)
  const limits = akar.larik.Limits ?? []
  const h = await kirim<HasilAchievement>('achievement', {
    idKontrak: akar.medan.ID ?? '',
    cari: nilaiParam(a.param?.search, s, l) === 'search',
    asAt: akar.medan['SearchData.CARI1'] ?? '',
    tahun: akar.medan['SearchData.CARI2'] ?? '',
    rnmShareP: '',
    limits: limits.map((x) => ({
      TreatyType: teks(x.TreatyType),
      Detail: anak(x, 'Detail').map((d) => ({ TreatyGroup: teks(d.TreatyGroup), EPIList: anak(d, 'EPIList') })),
    })),
  })
  const baru = limits.map((x, i) => ({ ...x, Detail: anak(x, 'Detail').map((d, k) => ({ ...d, ...(h.limits[i]?.[k] ?? {}) })) }))
  return {
    ...kosong(),
    akar: {
      larik: {
        Limits: baru,
        // Daftar dropdown As At / Quarter Year (`pageList`, nilai `.CARI6`).
        TempQuarter: h.kuartal.map((q) => ({ CARI6: q })),
        TempQuarterYear: h.tahunKuartal.map((q) => ({ CARI6: q })),
      },
      medan: { 'FlagExcel.CARI1': h.flagExcel ? '1' : '' },
    },
  }
}

// ---------------------------------------------------------------------
// Cabang ADJUST PREMIUM (`EDMState = 3`) — rute `aktual`
// (`treatyin/backend/services/hitung_aktual.go`): halaman `ActualValue`,
// selisihnya `ValueDifference`, dan kunci akar yang Activity tulis.
// ---------------------------------------------------------------------

/** Kunci akar yang rantai Adjust Premium baca. */
const MEDAN_AKTUAL = ['ID', 'TotalEgnpiAmount', 'RNMShare', 'BrokeragePercent', 'FacultativeShare', 'FacultativeShareBrokerage'] as const
const LARIK_AKTUAL = [
  'Limits', 'Share', 'EGNPI', 'FacultativeShareList', 'LimitShareSummaryList', 'LimitFacShareSummaryList',
  'TotalLimitPremiEarnNP', 'TotalLimitMDPNP',
  'TotalSpreadedRnmProp', 'TotalSpreadedRnmRIProp', 'TotalSpreadedNetPremi', 'TotalSpreadedNetPremiRI',
  'TotalFacShareRnmNP', 'TotalFacShareGrossNP', 'TotalFacShareNetNP', 'TotalFacShareDeductionNP',
] as const

/**
 * Satu aksi rute `aktual`. `ActualValue` dan (bila disentuh) `ValueDifference`
 * DIGANTI utuh; kunci akar hanya yang Activity tulis. `rnm-baris` = rincian
 * Share baris `Param.idx`.
 */
const aktual =
  (aksi: 'nilai' | 'share' | 'limits' | 'ringkas-limit' | 'ringkas-share' | 'rnm-baris'): Rumus =>
  async (s, a, l) => {
    const akar = akarDari(s, l)
    const nomor = Number(nilaiParam(a.param?.idx, s, l))
    const h = await kirim<{ actual: HalamanPohon; selisih: HalamanPohon | null; akar: HalamanPohon; pesan: string[] }>('aktual', {
      aksi,
      akar: {
        medan: Object.fromEntries(MEDAN_AKTUAL.map((k) => [k, akar.medan[k] ?? ''])),
        larik: Object.fromEntries(LARIK_AKTUAL.filter((k) => akar.larik[k] !== undefined).map((k) => [k, akar.larik[k] ?? []])),
      },
      actual: halamanBertitik(akar, AKTUAL),
      indeks: Number.isFinite(nomor) && nomor > 0 ? nomor - 1 : 0,
    })
    const av = gantiBertitik(akar, AKTUAL, h.actual)
    const larik: Record<string, Baris[]> = { ...av.larik, ...h.akar.larik }
    const medan: Record<string, string> = { ...av.medan, ...h.akar.medan }
    if (h.selisih !== null) {
      const vd = gantiBertitik(akar, SELISIH, h.selisih)
      Object.assign(larik, vd.larik)
      Object.assign(medan, vd.medan)
    }
    return { ...kosong(), pesan: h.pesan, akar: { larik, medan } }
  }

/** Activity `|` `Type` (huruf kecil) → rumus. */
const RUMUS: Readonly<Record<string, Rumus>> = {
  'TreatyInNPSetTotal|retention': totalRetensi,
  'TreatyInNPSetTotal|egnpi': egnpi('total'),
  'TreatyInEGNPIListValue|': egnpi('nilai'),
  'SetAmountConversion|': konversiEgnpi,
  'TreatyInSetReport|': periodePelaporan,
  'TreatyInNPSetTotal|limits': limitNP('total'),
  'TreatyInSummaryLimit|': limitNP('total'),
  'TreatyInSummaryMDP|': summaryMDP,
  'TreatyInLimitsListValue|': limitNP('nilai-list'),
  'DetailCalculation|adj': limitNPLayer('adj'),
  'DetailCalculation|mdp': limitNPLayer('mdp'),
  'SetReinstatementPct|': limitNPLayer('reinstatement'),
  'TotalEgnpi|': limitNPLayer('egnpi'),
  'TreatyInNPSetTotal|share': shareNP('total'),
  'TreatyInSummaryLimitShare|': shareNP('total'),
  'TreatyInSummaryLimitFacShare|': shareNP('total'),
  'TreatyInShareListValue|': shareNP('nilai-share'),
  'TreatyInXOLAddSpreading|': shareNP('rnm'),
  'TreatyInSetBrokerage|': shareNP('set-brokerage'),
  'TreatyInNonAddItem|share': shareNP('summary'),
  'TreatyInXOLAddSpreadingDetail|': shareNP('rnm-baris'),
  'FetchQSfromMasterXOL|': shareNP('spreading-type'),
  'SetSpreadingXOL|': shareNP('spreading-pct'),
  'CalculateDeduction|': calculateDeduction,
  'TreatyInSetValueInstallment|': nilaiAngsuran,
  'SetTotalInstallment|': totalHalamanAngsuran,
  'TreatyInNPSetTotal|installment': totalAngsuran,
  'TreatyEDMCalculateDifference|': selisihEDM,
  'GetAchievement|': achievement,
  'TreatyInActualUpdateValue|': aktual('nilai'),
  'TreatyInActualUpdateValueShare|share': aktual('share'),
  'TreatyInNPSetTotalActual|limits': aktual('limits'),
  'TreatyInSummaryLimitActual|': aktual('ringkas-limit'),
  'TreatyInSummaryLimitShareActual|': aktual('ringkas-share'),
  'TreatyInXOLAddSpreadingDetailActual|': aktual('rnm-baris'),
  'TreatyInUpdatePaymentDate|': tanggalBayar(false),
  'TreatyInUpdatePaymentDate_Act|': tanggalBayar(true),
  'TreatyInSetAccountReport|': akumulasi('periode'),
  'TreatyInAccumulationSetSubDue|': akumulasi('jatuh-tempo'),
  'TreatyInPropshare|': shareProp('share'),
  'TreatyInPropshareDetail|': shareProp('detail'),
  'FetchQSfromMaster|': shareProp('spreading'),
  'SetSpreadName|': shareProp('sebar-nama'),
  'LimitCalculation|': limitCalculation,
  'PremiumReserveCalculate|': cadanganPremi,
  'SetCurrName_Act|': pasangMataUang,
  'SetCurrNameMasterTreaty_Act|': namaMataUangMaster,
  'SetTreatyGroupName_Act|': namaKelompokTreaty,
  'SetTreatyTypeName_Act|': namaJenisTreaty,
}

/**
 * DataTransform (pra-refresh `transformasi`, atau `runDataTransform`) →
 * rumus. Isinya disalin dari korpus `Treaty In Adjustment/DataTransform`.
 */
const RUMUS_DT: Readonly<Record<string, Rumus>> = {
  // `TreatyIn.AccumulationList` REMOVE — pra-refresh isian Period.
  TreatyInDeleteAccumulationLists: () => Promise.resolve({ ...kosong(), akar: { larik: { AccumulationList: [] } } }),
  // `.ID = .pxListSubscript` — pra-refresh `FetchQSfromMaster` rincian Limits.
  SetDetailsID: (_s, _a, l) => Promise.resolve({ larik: {}, medan: { ID: String(nomorHalaman(l)) }, pesan: [] }),
  TreatyTypeSetIndex: (s, _a, l) => Promise.resolve(setIndeksJenis(s, l)),
  // FOR_EACH `TreatyIn.Limits` (idx = subscript): FOR_EACH `.LayerList`: `.Layer = idx`.
  SetIndexLayer_DT: (s, _a, l) => {
    const limits = akarDari(s, l).larik.Limits ?? []
    if (!limits.some((x) => Array.isArray(x.LayerList))) return Promise.resolve(kosong())
    const baru = limits.map((x, i) =>
      Array.isArray(x.LayerList) ? { ...x, LayerList: anak(x, 'LayerList').map((y) => ({ ...y, Layer: String(i + 1) })) } : x,
    )
    return Promise.resolve({ ...kosong(), akar: { larik: { Limits: baru } } })
  },
  // Sel `Reinstatement_List` rincian Layers — `Subscript = .pxListSubscript`.
  ReCalculateReinstatement: limitNPLayer('reinst-tambahan'),
  CalculateReinstatement: limitNPLayer('reinst-jumlah'),
  CalculateReinstatementPct: limitNPLayer('reinst-persen'),
  // `SearchData.CARI2 = ""` — isian As At (Achievement) berubah.
  Reset_DT: () => Promise.resolve({ ...kosong(), akar: { medan: { 'SearchData.CARI2': '' } } }),
  // `TreatyIn.Exclusions = TreatyIn.ExclusionsP`, `SpecialConditions = SpecialConditionsP`.
  TreatyInCopyConditions: (s, _a, l) => {
    const m = akarDari(s, l).medan
    return Promise.resolve({ ...kosong(), akar: { medan: { Exclusions: m.ExclusionsP ?? '', SpecialConditions: m.SpecialConditionsP ?? '' } } })
  },
}

/**
 * DataTransform yang menulis nilai halaman ke DIRINYA sendiri — parameternya
 * properti halaman yang sama (`SetCoB(business = .ClassOfBusiness, …)`,
 * `SetCurrencyID(Currency = .Currency, ID = .CurrencyID)`). Nol perubahan.
 */
const DT_IDENTITAS: ReadonlySet<string> = new Set(['SetCoB', 'SetCoBID', 'SetCurrencyID', 'SetTreatyGroupID'])

/**
 * Activity yang SUDAH dijalankan rute bersama langkah lain — dilewati bila
 * langkah pencakupnya sudah berjalan dalam rantai yang sama, supaya satu
 * klik tidak memanggil rute yang sama tiga kali. Sendirian, ia tetap
 * dijalankan lewat entrinya di `RUMUS`.
 */
const CAKUPAN: Readonly<Record<string, readonly string[]>> = {
  'TreatyInNPSetTotal|limits': ['TreatyInSummaryLimit|'],
  'TreatyInNPSetTotal|share': ['TreatyInSummaryLimitShare|', 'TreatyInSummaryLimitFacShare|'],
  // Rute `summary` = NonAddItem(share) + XOLAddSpreading + SetBrokerage.
  'TreatyInNonAddItem|share': ['TreatyInXOLAddSpreading|', 'TreatyInSetBrokerage|'],
  // `TotalEgnpi` menapaki SEMUA layer — `idxLimit` tidak dibacanya.
  'TotalEgnpi|': ['TotalEgnpi|'],
  // Rute `aktual` `limits` = NPSetTotalActual + kedua Summary Actual.
  'TreatyInNPSetTotalActual|limits': ['TreatyInSummaryLimitActual|', 'TreatyInSummaryLimitShareActual|'],
}

const kunciRumus = (a: AksiTombol) => `${a.aktivitas ?? ''}|${(a.param?.Type ?? a.param?.type ?? '').toLowerCase()}`

/** Satu langkah rantai sesudah dipecah — pra-DT dan Activity menjadi dua. */
interface Langkah {
  a: AksiTombol
  kunci: string
  /** Rumusnya; tidak ada = langkah tanpa perubahan (`postValue`, refresh kosong). */
  r?: Rumus
  /** Langkah yang benar-benar menghitung — rantai tanpa satu pun bukan rumus. */
  nyata: boolean
  /** Langkah BERSYARAT tanpa rumus — rantai batal bila syaratnya terpenuhi. */
  belum?: boolean
}

/** Pecah satu aksi; `undefined` = aksi tanpa rumus (rantai mati). */
function pecah(a: AksiTombol): Langkah[] | undefined {
  const kunci = kunciRumus(a)
  switch (a.aksi) {
    case 'setValue':
      return [{ a, kunci, r: setel, nyata: false }]
    // `postValue` hanya mengirim isian ke clipboard Pega — di sini isian SUDAH
    // di keadaan layar.
    case 'postValue':
      return [{ a, kunci, nyata: false }]
    case 'refresh':
    case 'runActivity':
    case 'runDataTransform': {
      const out: Langkah[] = []
      const dt = a.transformasi
      if (dt !== undefined && !DT_IDENTITAS.has(dt)) {
        const r = RUMUS_DT[dt]
        if (r !== undefined) out.push({ a, kunci: `DT:${dt}`, r, nyata: true })
        else if (a.syarat !== undefined) out.push({ a, kunci: `DT:${dt}`, nyata: false, belum: true })
        else return undefined
      }
      if (a.aktivitas !== undefined && a.aksi !== 'runDataTransform') {
        const r = RUMUS[kunci]
        if (r !== undefined) out.push({ a, kunci, r, nyata: true })
        else if (a.syarat !== undefined) out.push({ a, kunci, nyata: false, belum: true })
        else return undefined
      }
      if (out.length === 0) out.push({ a, kunci, nyata: false })
      return out
    }
    default:
      return undefined
  }
}

function susun(aksi: readonly AksiTombol[]): Langkah[] | undefined {
  const out: Langkah[] = []
  for (const a of aksi) {
    const x = pecah(a)
    if (x === undefined) return undefined
    out.push(...x)
  }
  return out
}

/**
 * Rantai yang SELURUH langkahnya tanpa perubahan (`postValue`, refresh
 * kosong, DT identitas) — bukan rumus, dan juga tidak mati.
 */
export function rantaiKosong(aksi: readonly AksiTombol[]): boolean {
  const xs = susun(aksi)
  return xs !== undefined && !xs.some((x) => x.nyata || x.belum === true)
}

/** Jalankan langkah-langkah atas akar + jalur; lihat kepala berkas. */
const jalankan =
  (xs: readonly Langkah[]) =>
  async (awal: SisiPenyesuaian, l: LingkupRumus = {}): Promise<HasilRantai> => {
    const berjalur = l.jalur !== undefined
    const jalur = l.jalur ?? []
    // Halaman pemicu DITARUH di akar pada jalurnya — penghapusan baris memberi
    // halaman yang barisnya sudah dibuang.
    let akar = jalur.length === 0 ? awal : ubahHalaman(l.panel ?? awal, jalur, () => awal)
    const akarAwal = akar
    const halaman = l.halaman ?? {}
    const pesan: string[] = []
    const larik: Record<string, Baris[]> = {}
    const medan: Record<string, string> = {}
    const tercakup = new Set<string>()
    for (const x of xs) {
      if (x.a.syarat !== undefined && !syaratTerpenuhi([x.a.syarat], halaman)) continue
      if (x.belum === true) {
        const nama = x.a.aktivitas ?? x.a.transformasi ?? x.a.aksi
        return { larik: {}, medan: {}, pesan: [`${PENYESUAIAN.langkahBelum} ${nama}`], awal: akarAwal, akhir: akarAwal }
      }
      if (x.r === undefined || tercakup.has(x.kunci)) continue
      for (const c of CAKUPAN[x.kunci] ?? []) tercakup.add(c)
      // Tanpa jalur (pemakaian lama/uji): halaman rumus = akarnya sendiri, `panel` apa adanya.
      const h = berjalur ? await x.r(ambilHalaman(akar, jalur), x.a, { ...l, panel: akar, jalur }) : await x.r(akar, x.a, l)
      if (h.akar !== undefined) {
        akar = { ...akar, medan: { ...akar.medan, ...h.akar.medan }, larik: { ...akar.larik, ...h.akar.larik } }
        if (jalur.length === 0) {
          Object.assign(larik, h.akar.larik)
          Object.assign(medan, h.akar.medan)
        }
      }
      if (Object.keys(h.larik).length > 0 || Object.keys(h.medan).length > 0) {
        akar = ubahHalaman(akar, jalur, (p) => ({ ...p, medan: { ...p.medan, ...h.medan }, larik: { ...p.larik, ...h.larik } }))
        Object.assign(larik, h.larik)
        Object.assign(medan, h.medan)
      }
      pesan.push(...h.pesan)
    }
    return { larik, medan, pesan, awal: akarAwal, akhir: akar }
  }

/**
 * Rantai rumus SATU tombol / perilaku `change` — atau `undefined` bila salah
 * satu Activity-nya belum punya rumus, atau rantainya tanpa satu pun langkah
 * yang menghitung. ⛔ Rantai yang separuh jalan TIDAK dijalankan: tombol
 * Pega menjalankan semua langkahnya atau tidak sama sekali. Langkah BERSYARAT
 * tanpa rumus diterima; bila syaratnya terpenuhi, rantai batal dengan pesan.
 */
export function rantaiRumus(t: Pick<TombolKerangka, 'aksi'>): ((s: SisiPenyesuaian, l?: LingkupRumus) => Promise<HasilRantai>) | undefined {
  const xs = susun(t.aksi)
  if (xs === undefined || !xs.some((x) => x.nyata)) return undefined
  return jalankan(xs)
}
