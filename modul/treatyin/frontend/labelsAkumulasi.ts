// Label tab ACCUMULATION (cabang Proporsional) — `Section/TreatyInTabsProportional.xml`
// tab `Accumulation`, dan gambar Pega `19-prop-menu-accumulation.png`.

export const AKUMULASI = {
  /** Judul blok @1250071 (`pyIncludeHeader` true) — gambar 19. */
  judul: 'Accumulation Control',
  /** Medan `TreatyIn.AccumulationPeriod` — gambar 19. */
  periode: 'Period',
  tanpaBaris: 'No items',
  tambah: 'Add',
  hapus: 'Delete',
} as const

/**
 * Pilihan `TreatyIn.AccumulationPeriod` — NILAI yang
 * `Activity/TreatyInSetAccountReport.xml` periksa (langkah 2–6: `none`,
 * `quarter`, `half`, `month`, `other`).
 *
 * ⛔ Label: hanya `None` yang terlihat (gambar 19). Rule Property-nya tidak
 * diekspor, jadi sisanya tampil APA ADANYA — konvensi modul ini untuk prompt
 * value yang belum diketahui (`backend/services/prompt_value.go`).
 */
export const OPSI_PERIODE_AKUMULASI: readonly { value: string; label: string }[] = [
  { value: 'none', label: 'None' },
  { value: 'quarter', label: 'quarter' },
  { value: 'half', label: 'half' },
  { value: 'month', label: 'month' },
  { value: 'other', label: 'other' },
]

/**
 * Awalan `Period` baris baru — DataTransform `TreatyInAddAccumulation`
 * langkah 4–7 (`param.keyword`): quarter `Q `, half `H `, month `M `,
 * none `T `. Periode lain: tanpa awalan.
 */
export const AWALAN_AKUMULASI: Readonly<Record<string, string>> = {
  quarter: 'Q ',
  half: 'H ',
  month: 'M ',
  none: 'T ',
}
