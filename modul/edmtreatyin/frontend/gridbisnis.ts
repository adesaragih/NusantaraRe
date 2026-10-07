// Saring, urut, dan paging grid popup `Section/BusinessAndSOBListEDM` - dikerjakan di klien atas baris jawaban
// `POST /kasus/{id}/bisnis`. Asal pola: `modul/nbtreatyin/frontend/gridbisnis.ts` (06-10-2026), versi saringan daftar
// nilai (`pyColumnFilteringDropDown`) - backend EDM tidak menerima saringan, jadi nilai disaring di klien.
//
// Sumber XML (`pyGridProps` / `pyColumnDefaults` kedua grid S1 dan S6, identik):
//   pyGridSorting true, pyGridFiltering true   -> urut dan saring per kolom; nol kotak cari global
//   pyPageMode Numeric, pyPageSize 50          -> 50 baris per halaman
//   kolom 1 (tombol Choose) tanpa urut / saring
//   kolom 2-9 pyColumnSorting true, pyColumnFilteringDropDown true -> saringan berupa daftar nilai
//   kolom 2 `.ID` pySortType DESC, pySortOrder 1 -> urutan awal ID Revision turun
//
// Yang TIDAK terbaca dari XML dan dipilih di sini: urutan teks = urutan kode karakter; nilai kosong di depan pada
// urutan naik. Seluruh kolom grid ini teks / tanggal (tanpa pxCurrency).

import type { BarisBisnis } from './api'

export type ArahUrut = 'naik' | 'turun'

export interface UrutBisnis {
  kolom: string
  arah: ArahUrut
}

/** Kolom 2 `.ID` pySortType DESC, pySortOrder 1. */
export const URUT_AWAL_BISNIS: UrutBisnis = { kolom: 'ID', arah: 'turun' }

/** Saringan aktif per kolom; "" atau tidak ada = kolom itu tidak disaring. */
export type SaringanKolom = Record<string, string>

function banding(a: string, b: string): number {
  return a === b ? 0 : a < b ? -1 : 1
}

/** Baris yang lolos seluruh saringan kolom (nilai sama persis), diurut stabil menurut `urut`. */
export function saringUrutBisnis<T extends BarisBisnis>(
  baris: T[],
  saringan: SaringanKolom,
  urut: UrutBisnis | null,
): T[] {
  const aktif = Object.entries(saringan).filter(([, v]) => v !== '')
  const hasil = baris.filter((r) => aktif.every(([k, v]) => (r[k] ?? '') === v))
  if (urut === null) return hasil
  const arah = urut.arah === 'naik' ? 1 : -1
  return hasil
    .map((r, i) => ({ r, i }))
    .sort((x, y) => arah * banding(x.r[urut.kolom] ?? '', y.r[urut.kolom] ?? '') || x.i - y.i)
    .map((x) => x.r)
}

/** Pilihan daftar saringan satu kolom: nilai berbeda selain kosong, berurut naik. */
export function nilaiSaringBisnis(baris: BarisBisnis[], kolom: string): string[] {
  const ada = new Set<string>()
  for (const r of baris) {
    const v = r[kolom] ?? ''
    if (v !== '') ada.add(v)
  }
  return [...ada].sort(banding)
}

/** Urutan sesudah judul `kolom` diklik: kolom lain mulai naik; kolom yang sama berganti arah. */
export function urutBerikut(kini: UrutBisnis | null, kolom: string): UrutBisnis {
  if (kini !== null && kini.kolom === kolom) return { kolom, arah: kini.arah === 'naik' ? 'turun' : 'naik' }
  return { kolom, arah: 'naik' }
}
