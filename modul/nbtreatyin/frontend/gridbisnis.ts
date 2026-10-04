// Saring, urut, dan paging grid AKTIF popup `Section/BusinessAndSOBList` (RD `BrowseTreatyJoinEDM`,
// audit silang P3 W1) - dikerjakan di klien atas baris jawaban server (<= 500, RD pyMaxRecords).
//
// Sumber XML (`pyGridProps` dan `pyColumnDefaults` grid itu):
//   pyGridSorting true, pyGridFiltering true      -> urut dan saring per kolom; nol kotak cari global
//   pyPageMode Numeric, pyPageSize 50             -> 50 baris per halaman, tombol nomor halaman
//   kolom 1 (tombol Choose) pyColumnSorting false, pyColumnFilteringDropDown false
//   kolom 2-26 pyColumnSorting true, pyColumnFilteringDropDown true -> saringan berupa daftar nilai
//   kolom 2 `.TREATYID` pySortType ASC, pySortOrder 1 (sama dengan urutan RD)
//
// Yang TIDAK terbaca dari XML dan karenanya dipilih di sini: urutan teks = urutan kode karakter
// (sama dengan ORDER BY TREATYID di Oracle berurutan biner, sehingga urutan awal klien = urutan
// server), dan nilai kosong di depan pada urutan naik. Kolom `pxCurrency` diurut menurut nilai
// desimalnya tanpa `Number` (spec §5.6: uang tidak pernah float).

import { desimalSah, jumlahDesimal } from '../../../inti/frontend/lib/desimal'

import type { BarisKontrak } from './api'
import { tandaTeks } from './tanda'

/** `pyPageSize` grid. */
export const BARIS_PER_HALAMAN_BISNIS = 50

export type ArahUrut = 'naik' | 'turun'

export interface UrutBisnis {
  kolom: string
  arah: ArahUrut
}

/** Kolom 2 `.TREATYID` pySortType ASC, pySortOrder 1. */
export const URUT_AWAL_BISNIS: UrutBisnis = { kolom: 'TREATYID', arah: 'naik' }

/** Sel `pyFormat` pxCurrency (sel lain pxDisplayText/pxTextInput = teks). */
const KOLOM_UANG = new Set(['LIMITVALUE', 'RETENTIONVALUE', 'EPIVALUE', 'MDPVALUE', 'NETPREMIVALUE', 'SHAREVALUE'])

/** Saringan aktif per kolom; "" atau tidak ada = kolom itu tidak disaring. */
export type SaringanKolom = Record<string, string>

/** Lawan tanda teks desimal sah (`12` -> `-12`, `-1.5` -> `1.5`). */
function lawan(s: string): string {
  const t = s.trim()
  return t.startsWith('-') ? t.slice(1) : '-' + t.replace(/^\+/, '')
}

/** Bandingkan dua angka desimal teks (`-12.50`, `.5`) tanpa float - tanda selisih eksak
 *  `jumlahDesimal` (inti/frontend/lib/desimal.ts); `null` bila salah satunya bukan angka. */
function bandingAngka(a: string, b: string): number | null {
  if (!desimalSah(a) || !desimalSah(b)) return null
  return tandaTeks(jumlahDesimal([a.trim(), lawan(b)]).total)
}

function bandingTeks(a: string, b: string): number {
  return a === b ? 0 : a < b ? -1 : 1
}

/** Pembanding naik satu kolom: kosong di depan; kolom uang menurut nilai, selainnya teks. */
function bandingKolom(kolom: string, a: string, b: string): number {
  if (a === '' || b === '') return bandingTeks(a, b)
  if (KOLOM_UANG.has(kolom)) {
    const n = bandingAngka(a, b)
    if (n !== null) return n
  }
  return bandingTeks(a, b)
}

/** Baris yang lolos seluruh saringan kolom, diurut stabil menurut `urut` (null = urutan datang). */
export function saringUrutBisnis<T extends BarisKontrak>(baris: T[], saringan: SaringanKolom, urut: UrutBisnis | null): T[] {
  const aktif = Object.entries(saringan).filter(([, v]) => v !== '')
  const hasil = baris.filter((r) => aktif.every(([k, v]) => (r[k] ?? '') === v))
  if (urut === null) return hasil
  const arah = urut.arah === 'naik' ? 1 : -1
  return hasil
    .map((r, i) => ({ r, i }))
    .sort((x, y) => arah * bandingKolom(urut.kolom, x.r[urut.kolom] ?? '', y.r[urut.kolom] ?? '') || x.i - y.i)
    .map((x) => x.r)
}

/** Pilihan daftar saringan satu kolom: nilai berbeda selain kosong, berurut naik seperti kolomnya. */
export function nilaiSaringBisnis(baris: BarisKontrak[], kolom: string): string[] {
  const ada = new Set<string>()
  for (const r of baris) {
    const v = r[kolom] ?? ''
    if (v !== '') ada.add(v)
  }
  return [...ada].sort((a, b) => bandingKolom(kolom, a, b))
}

/** Urutan sesudah judul `kolom` diklik: kolom lain mulai naik; kolom yang sama berganti arah. */
export function urutBerikut(kini: UrutBisnis | null, kolom: string): UrutBisnis {
  if (kini !== null && kini.kolom === kolom) return { kolom, arah: kini.arah === 'naik' ? 'turun' : 'naik' }
  return { kolom, arah: 'naik' }
}

/** Kolom uang (pxCurrency) - disajikan berformat angka. */
export function kolomUangBisnis(kolom: string): boolean {
  return KOLOM_UANG.has(kolom)
}
