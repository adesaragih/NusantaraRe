// Bentuk kerangka bangkitan (`kerangka.gen.ts`) — ditulis tangan, dibaca layar.
//
// ⛔ Kerangka adalah SALINAN Section ekspor, bukan rancangan: urutan butir =
// urutan di Section, judul tampil hanya bila kepala blok `BAR`/`TABBED`.

/**
 * `sisi` = halaman yang panel itu tampilkan; `akar` = halaman `TreatyIn`
 * (New). Di Section RINCIAN (`KERANGKA_RINCIAN`) `sisi` adalah halaman BARIS
 * yang dibuka — ikatannya `.X` di ekspor.
 */
export type Dari = 'sisi' | 'akar' | 'sesi' | null
// `sesi` = halaman SESI Pega (`SearchData.CARI1`, `FlagExcel.CARI1` —
// Achievement rincian Limits): kuncinya UTUH dan hidup di akar panel New,
// bukan bagian dokumen `TreatyIn`.

/**
 * Kapan sebuah sel/medan BACA-SAJA di mode Edit panel New: `'selalu'`, atau
 * daftar syarat — salah satu BENAR → baca-saja (`pyReadOnlyCondition` dan
 * `pyDisabledWhen`, termasuk kunci per Material Type). Tidak ada → dapat
 * disunting. Dibangkitkan `kunci_baca` di `alat/ekstrak_kerangka.py`.
 */
export type Baca = 'selalu' | readonly string[]

/**
 * Sumber daftar sel `pxDropdown`/`pxAutoComplete` (`pyListDataSource`):
 * `reportdefinition` membawa nama RD beserta medan NILAI dan TAMPIL-nya;
 * `associated` = daftar milik rule Property yang TIDAK diekspor;
 * `pageList` = halaman sesi.
 */
export interface SumberPilihan {
  sumber: string
  rd?: string
  nilai?: string
  tampil?: string
  halaman?: string
}

/** Satu aksi klik tombol, berurutan seperti di `pyBehaviors` ekspor. */
export interface AksiTombol {
  /** `refresh` · `addRow` · `deleteRow` · `setValue` · `localAction` · `runActivity` … */
  aksi: string
  /** Activity / local action / flow action yang dijalankan. */
  aktivitas?: string
  /** Parameter Activity (`Type=retention`) atau pasangan `setValue`. */
  param?: Readonly<Record<string, string>>
  /** DataTransform PRA-refresh (`pyPreDataTransform`) — Add Accumulation. */
  transformasi?: string
  /** Parameter DataTransform itu (`pyDataTransformParams`) — `type = limit`, `Subscript = .pxListSubscript`. */
  paramDT?: Readonly<Record<string, string>>
  /**
   * Syarat aksi ini (`pyActionConditions`) — aksinya DILEWATI bila tidak
   * terpenuhi. Teks dari `syarat_aksi` pembangkit; dinilai atas halaman
   * pemicu (akar + skalar bertitik baris sel pemicu).
   */
  syarat?: string
}

/** Tombol ekspor — juga tombol ikon tanpa `pyLabel`. */
export interface TombolKerangka {
  t: 'tombol'
  at: number
  label: string
  syarat: readonly string[]
  aksi: readonly AksiTombol[]
  /** `pyImage` — `IconAdd.png`, `IconTrash.png`, `rpadd.gif` … */
  ikon?: string
  /** `pyDisabledWhen` — tombol tampil tetapi mati bila salah satu benar. */
  nonaktif?: readonly string[]
}

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
  /** Baca-saja per kolom; `null` = dapat disunting di mode Edit. */
  baca: readonly (Baca | null)[]
  /** Tombol baris per kolom (Delete/Hapus) — `null` untuk kolom data. */
  tombol: readonly (TombolKerangka | null)[]
  /** Tombol di sel KEPALA kolom (Tambah Co-Ins) — `null` bila tidak ada. */
  tombolKepala: readonly (TombolKerangka | null)[]
  /** Sumber daftar sel dropdown/autocomplete per kolom; `null` untuk yang lain. */
  pilihan: readonly (SumberPilihan | null)[]
  /** Aksi `change` sel per kolom (Activity sesudah isian berubah); `null` = tidak ada. */
  aksiUbah: readonly (readonly AksiTombol[] | null)[]
  /** Templat bingkai baris (`pyGridTemplateName`) — bingkai saja, bukan isinya. */
  templatBaris?: string
  /**
   * Section RINCIAN BARIS (`pyEditingMode = expandPane` → `pyEditAction` →
   * FlowAction → `pySectionReference`), kunci `KERANGKA_RINCIAN`.
   */
  rincian?: string
  /** Flow action rincian yang dirujuk tetapi TIDAK ada di ekspor. */
  rincianHilang?: string
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
  /** Baca-saja di mode Edit panel New; tidak ada = dapat disunting. */
  baca?: Baca
  /** Sumber daftar bila medannya dropdown/autocomplete. */
  pilihan?: SumberPilihan
  /** Aksi `change` — Activity yang Pega jalankan sesudah isian berubah. */
  aksiUbah?: readonly AksiTombol[]
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
  | TombolKerangka
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
