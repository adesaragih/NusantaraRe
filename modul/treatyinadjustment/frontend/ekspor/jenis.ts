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
  /**
   * Parameter RD (`pyReportDefParams`) — nama → ungkapan Pega apa adanya:
   * `.TreatyGroupID`, `.TreatyGroupList(1).TreatyGroupID`,
   * `TreatyIn.Commencement`, `"10001"`.
   */
  param?: Readonly<Record<string, string>>
  /**
   * Medan yang IKUT diisi saat satu pilihan dipilih (`pyAdditionalFields`,
   * `pySetValueOnSelect`): `CurrencyID` ← `.ID`, `ClassOfBusinessID` ← `.BizCode`.
   */
  setel?: readonly { target: string; dari: string }[]
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
  /** `showHarness` — harness yang dibuka (`pyHarnessName`); isinya `KERANGKA_HARNESS`. */
  harness?: string
  /** Judul jendela harness (`pyWindowName`), mis. `Facultative Calculation`. */
  jendela?: string
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
  /**
   * Mode baris grid (`pyGridProps/pyRowEditing`): `row` = sel disunting di
   * tempat; `readOnly` = grid baca-saja; `masterDetail` = baris TAMPIL saja,
   * penyuntingannya di panel rincian. Untuk dua yang terakhir pembangkit
   * sudah mengisi `baca` setiap kolom data dengan `selalu`.
   */
  modeBaris?: string
  /** `pyRODetails = true` — panel rincian pun baca-saja. */
  rincianBaca?: boolean
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
  /**
   * Layout group `pyHeaderType = TABBED` — blok bertanda ini yang BERURUTAN
   * digambar sebagai SATU strip tab (`KerangkaTab.tsx`, `GrupTab`).
   */
  tab?: true
  /**
   * Tata letak layout dinamis Pega (`pyLayoutOtherFormat`), 8 Oktober 2026:
   * `g2`/`g3`/`g4` = `Inline grid double/triple/quadruple` (kolom sama lebar),
   * `t3070` = `Inline 30 70 table`, `alir` = `Inline*` (mengalir sebaris),
   * `kiri` = `Stacked with labels left`. Tanpa = bertumpuk.
   */
  tata?: 'g2' | 'g3' | 'g4' | 't3070' | 'alir' | 'kiri'
}

export type ButirKerangka =
  | BlokKerangka
  | GridKerangka
  | MedanKerangka
  | { t: 'teks'; at: number; teks: string; syarat: readonly string[] }
  /**
   * Kotak SATUAN baca-saja berlabel — sel `.pyTemplateRichTextEditor`
   * (rincian EGNPI: `Amount in IDR` [IDR]). Teksnya dari tangkapan Pega,
   * `SATUAN_TEMPLAT` pembangkit (9 Oktober 2026).
   */
  | { t: 'satuan'; at: number; label: string; teks: string; syarat: readonly string[] }
  | TombolKerangka
  | { t: 'include'; at: number; nama: string; syarat: readonly string[] }
  /**
   * Slot kosong layout BERKOLOM (`g2`/`g3`/`g4`/`t3070`) — sel tersembunyi
   * (`1=2`, `Spacer`) atau tanpa isi TETAP memakan slot di Pega (gambar 05).
   */
  | { t: 'kosong'; at: number; syarat: readonly string[] }

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
