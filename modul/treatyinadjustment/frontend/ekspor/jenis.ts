// Bentuk kerangka bangkitan (`kerangka.gen.ts`) — ditulis tangan, dibaca layar.
//
// ⛔ Kerangka adalah SALINAN Section ekspor, bukan rancangan: urutan butir =
// urutan di Section, judul tampil hanya bila kepala blok `BAR`/`TABBED`.

/** `sisi` = halaman yang panel itu tampilkan; `akar` = halaman `TreatyIn` (New). */
export type Dari = 'sisi' | 'akar' | null

export interface GridKerangka {
  t: 'grid'
  at: number
  prop: string
  dari: Dari
  /** Jalur larik relatif halaman — boleh bertitik (`ValueDifference.Share`). */
  larik: string
  syarat: readonly string[]
  kolom: readonly string[]
  /** Properti tiap sel — boleh berindeks (`RnmLimitListDisplay(1).Value`). */
  kunci: readonly string[]
  lebar: readonly number[]
  /** `pyDecimalPlaces` sel; `null` = TIDAK dinyatakan ekspor. */
  desimal: readonly (number | null)[]
  format: readonly string[]
  syaratSel: readonly (string | null)[]
  /** Offset bita sel badan tiap kolom. */
  atSel: readonly number[]
  /** Rule rincian baris (`pyGridTemplateName`) — TIDAK ada di ekspor. */
  templatBaris?: string
}

export interface MedanKerangka {
  t: 'medan'
  at: number
  label: string
  dari: Dari
  kunci: string
  format: string
  desimal: number | null
  syarat: readonly string[]
  caption?: string
}

export interface BlokKerangka {
  t: 'blok'
  at: number
  judul: string
  syarat: readonly string[]
  anak: readonly ButirKerangka[]
}

export type ButirKerangka =
  | BlokKerangka
  | GridKerangka
  | MedanKerangka
  | { t: 'teks'; at: number; teks: string; syarat: readonly string[] }
  | { t: 'tombol'; at: number; label: string; syarat: readonly string[] }
  | { t: 'include'; at: number; nama: string; syarat: readonly string[] }

export interface Kerangka {
  at: number
  syarat: readonly string[]
  isi: readonly ButirKerangka[]
}

export interface Terbuang {
  berkas: string
  at: number
  jenis: string
  alasan: string
}
