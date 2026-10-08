// Aksi tombol kerangka panel New — yang dijalankan DI LAYAR, tanpa menulis
// ke basis data.
//
// ⛔ Aturan pemilik proses: isian tidak masuk basis data sebelum Save atau
// Submit ditekan. Tambah dan hapus baris di sini hanya mengubah keadaan
// layar; Save tetap mati sampai jalur simpannya diputuskan.
//
// ⭐ Isi baris baru DISALIN dari Activity korpus `Treaty In Adjustment`,
// bukan dikarang — tiap entri menyebut langkahnya.

import type { SisiPenyesuaian } from '../api'
import type { AksiTombol, Baca, TombolKerangka } from '../ekspor/jenis'
import { syaratTerpenuhi } from '../ekspor/syarat'
import { TOMBOL_IKON } from '../labelsPenyesuaian'
import { teks } from './baris'
import { rantaiKosong, rantaiRumus } from './rumus'

type Halaman = Readonly<Record<string, string>>
type Baris = Record<string, string>

/** Sel/medan ini BACA-SAJA sekarang? `'selalu'`, atau salah satu syaratnya benar. */
export function terkunci(baca: Baca | null | undefined, h: Halaman): boolean {
  if (baca === undefined || baca === null) return false
  if (baca === 'selalu') return true
  return baca.some((s) => syaratTerpenuhi([s], h))
}

/** Tombol tampil (`pyCondition`) dan, bila tampil, tidak dikunci (`pyDisabledWhen`)? */
export function tombolTampil(t: TombolKerangka | null, h: Halaman): t is TombolKerangka {
  return t !== null && syaratTerpenuhi(t.syarat, h)
}
export function tombolMati(t: TombolKerangka, h: Halaman): boolean {
  return (t.nonaktif ?? []).some((s) => syaratTerpenuhi([s], h))
}

/** Teks tombol: `pyLabel`, atau teks bawaan Pega untuk tombol ikon. */
export function labelTombol(t: TombolKerangka): string {
  if (t.label !== '') return t.label
  return TOMBOL_IKON[t.ikon ?? ''] ?? ''
}

interface TambahBaris {
  larik: string
  /** `nomor` = `.pxListSubscript` halaman tempat tombol itu (mulai 1; 0 di akar). */
  baris: (s: SisiPenyesuaian, nomor?: number) => Baris
}

/** `AddValue(type)` — larik per `param.type` (langkah 1–5, masing-masing `ID = ""`). */
const LARIK_ADD_VALUE: Readonly<Record<string, string>> = {
  pla: 'PLAList',
  cashloss: 'CashLossList',
  claimcoop: 'ClaimCoopList',
  epi: 'EPIList',
  reserve: 'ReserveList',
}

/** DT `AddLimitRetentionCession(type)` — `<larik>(<APPEND>).Note = .TreatyType`. */
const LARIK_LIMIT_RETENSI: Readonly<Record<string, string>> = {
  limit: 'IOOLimitList',
  retention: 'RetentionList',
  cession: 'CessionList',
}

const pertama = (s: SisiPenyesuaian, larik: string, kunci: string) => teks(s.larik[larik]?.[0]?.[kunci])

/** `TreatyInAddAccumulation` — `param.keyword` per periode; selain empat ini kosong. */
const AWALAN_AKUMULASI: Readonly<Record<string, string>> = { quarter: 'Q ', half: 'H ', month: 'M ', none: 'T ' }

/**
 * Activity `|` `Type` → baris yang ditambahkan.
 *
 * ⚠️ Yang TIDAK ada di sini sengaja: `TreatyInNonAddItem(share)` dan
 * `(installment)` MENGHITUNG barisnya dari Limits/Total Share — itu rumus,
 * bukan tambah baris kosong, dan dijalankan tombol berumus.
 */
const TAMBAH_BARIS: Readonly<Record<string, TambahBaris>> = {
  // TreatyInNonAddItem langkah 2 — `TreatyIn.Retention(<APPEND>).ID = ""`.
  'TreatyInNonAddItem|retention': { larik: 'Retention', baris: () => ({ ID: '' }) },
  // Langkah 3 — ID kosong, mata uang disalin dari baris Retention PERTAMA.
  'TreatyInNonAddItem|egnpi': {
    larik: 'EGNPI',
    baris: (s) => ({ ID: '', Currency: pertama(s, 'Retention', 'Currency'), CurrencyID: pertama(s, 'Retention', 'CurrencyID') }),
  },
  // Langkah 5 — `ID = @SizeOfPropertyList(TreatyIn.Limits)` sesudah APPEND,
  // lalu KELUAR (transisi `1==1` → 6): `CopyLastLimitNP` tidak tercapai.
  'TreatyInNonAddItem|limits': { larik: 'Limits', baris: (s) => ({ ID: String((s.larik.Limits?.length ?? 0) + 1) }) },
  // Langkah 15 / 16.
  'TreatyInNonAddItem|sharereins': { larik: 'ShareReins', baris: () => ({ ID: '' }) },
  'TreatyInNonAddItem|sharefacname': { larik: 'ShareFacultativeReinsurers', baris: () => ({ ID: '' }) },
  // Langkah 4 (`actualpremium`, lalu KELUAR) — baris Actual GNPI, mata uang
  // dari `TreatyIn.Retention(1)` (cabang Adjust Premium).
  'TreatyInNonAddItem|actualpremium': {
    larik: 'ActualValue.EGNPI',
    baris: (s) => ({ ID: '', Currency: pertama(s, 'Retention', 'Currency'), CurrencyID: pertama(s, 'Retention', 'CurrencyID') }),
  },
  // TreatyInPropAdd langkah 1 / 2.
  'TreatyInPropAdd|portfolio': { larik: 'Portfolio', baris: () => ({ Description: '' }) },
  'TreatyInPropAdd|limits': { larik: 'Limits', baris: () => ({ ID: '' }) },
  // TreatyInAddCurrency langkah 1 — `CurrencyList(<APPEND>).CurrencyID = ""`.
  'TreatyInAddCurrency|': { larik: 'CurrencyList', baris: () => ({ CurrencyID: '' }) },
  // DataTransform `TreatyInAddAccumulation`: hitung baris yang ada, awalan
  // menurut `AccumulationPeriod`, lalu `Period = awalan + (cacah + 1)`.
  'TreatyInAddAccumulation|': {
    larik: 'AccumulationList',
    baris: (s) => ({ Period: `${AWALAN_AKUMULASI[s.medan.AccumulationPeriod ?? ''] ?? ''}${String((s.larik.AccumulationList?.length ?? 0) + 1)}` }),
  },

  // ⭐ 7 Oktober 2026 — Section rincian Limits Prop dan Share Non-Prop.
  // ⚠️ `AddDeduction` korpus ADJUSTMENT: `.DeductionList(<APPEND>).Currency =
  // .Currency` — korpus Treaty In menulis `.CurrencyIOOLimit`. Kedua
  // korpus BERBEDA di sini; layar ini mengikuti korpusnya sendiri.
  'AddDeduction|': { larik: 'DeductionList', baris: (s) => ({ Currency: teks(s.medan.Currency) }) },
  // `AddClassofBusiness`: ID kosong, `ParentID = .pxListSubscript`, `TreatyType = .TreatyType`.
  'AddClassofBusiness|': {
    larik: 'Detail',
    baris: (s, nomor = 0) => ({ ID: '', ParentID: String(nomor), TreatyType: teks(s.medan.TreatyType) }),
  },
  ...Object.fromEntries(
    Object.entries(LARIK_ADD_VALUE).map(([jenis, larik]) => [`AddValue|${jenis}`, { larik, baris: () => ({ ID: '' }) }]),
  ),
  ...Object.fromEntries(
    Object.entries(LARIK_LIMIT_RETENSI).map(([jenis, larik]) => [
      `AddLimitRetentionCession|${jenis}`,
      { larik, baris: (s: SisiPenyesuaian) => ({ Note: teks(s.medan.TreatyType) }) },
    ]),
  ),
}

/** Jenis tambah baris: `Type`/`type` Activity, atau `type` DataTransform — tanpa kutip (`"reserve"`). */
const jenisAksi = (a: AksiTombol) => (a.param?.Type ?? a.param?.type ?? a.paramDT?.type ?? '').replace(/^"(.*)"$/, '$1')

const kunciAksi = (a: AksiTombol) => `${a.aktivitas ?? a.transformasi ?? a.aksi}|${jenisAksi(a)}`

/** Baris yang tombol ini tambahkan — `undefined` bila tombol ini bukan tambah baris. */
export function tambahDari(t: TombolKerangka): TambahBaris | undefined {
  for (const a of t.aksi) {
    const x = TAMBAH_BARIS[kunciAksi(a)]
    if (x !== undefined) return x
  }
  return undefined
}

/**
 * Tombol TULIS tab Information & Submit blok `TreatyMasterInEDM` — dijalankan
 * form (rute `/api/treaty-in/penyesuaian/*`), bukan di layar:
 *
 *   Submit         `refresh` + Activity `TreatyInSubmitEDM`
 *   Decline offer  `localAction` `TreatyInDeclineConfirmationEDM` (konfirmasi
 *                  lalu `TreatyInDeclineConfirmation_postactEDM`)
 *
 * ⛔ Varian non-EDM (`TreatyInSubmit`, `TreatyInDeclineConfirmation`, blok
 * `!TreatyMasterInEDM`) TIDAK termasuk — ia milik layar Treaty In.
 */
export type JenisTulis = 'submit' | 'tolak'
const AKTIVITAS_TULIS: Readonly<Record<string, JenisTulis>> = {
  TreatyInSubmitEDM: 'submit',
  TreatyInDeclineConfirmationEDM: 'tolak',
}
export function tulisanDari(t: TombolKerangka): JenisTulis | undefined {
  for (const a of t.aksi) {
    const j = AKTIVITAS_TULIS[a.aktivitas ?? '']
    if (j !== undefined) return j
  }
  return undefined
}

/**
 * Tombol yang MENGUNDUH (`showHarness`) — `GenerateCSVTreaty` (Generate
 * Excel, Achievement). Nol perubahan keadaan; isinya dari halaman tombol.
 */
export const unduhanDari = (t: TombolKerangka): 'achievement' | undefined =>
  t.aksi.some((a) => a.aksi === 'showHarness' && a.aktivitas === 'GenerateCSVTreaty') ? 'achievement' : undefined

/** `addRow` / `deleteRow` grid — dijalankan atas grid tempat tombol itu duduk. */
export const tambahBarisGrid = (t: TombolKerangka) => t.aksi.some((a) => a.aksi === 'addRow')
export const hapusBarisGrid = (t: TombolKerangka) => t.aksi.some((a) => a.aksi === 'deleteRow')

/**
 * Rantai SESUDAH `deleteRow` / `addRow` — mis. `Remove` grid IOOLimitList
 * (`deleteRow` lalu `refresh LimitCalculation`), `Add Treaty Group` Layers
 * (`addRow` lalu DT `TreatyTypeSetIndex`). `null` = tombol itu HANYA
 * menghapus/menambah; `undefined` = langkah sesudahnya belum punya rumus,
 * dan tombolnya dimatikan: rantai separuh jalan tidak dijalankan.
 */
function rantaiSesudah(t: TombolKerangka, jenis: 'deleteRow' | 'addRow') {
  const sisa = t.aksi.filter((a) => a.aksi !== jenis)
  if (rantaiKosong(sisa)) return null
  return rantaiRumus({ aksi: sisa })
}

export const rantaiSesudahHapus = (t: TombolKerangka) => rantaiSesudah(t, 'deleteRow')
export const rantaiSesudahTambah = (t: TombolKerangka) => rantaiSesudah(t, 'addRow')
