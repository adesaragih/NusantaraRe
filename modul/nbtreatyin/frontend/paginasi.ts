// Paginasi grid ber-`pyGridProps/pyPageMode=Numeric` - satu halaman berisi
// `ukuran` baris, halaman berbasis 1.
//
// Dipakai grid `.SuggestList` `Section/ListSuggest` (pyPageSize "Other",
// pyPageSizeOther 5, pyGridPaginator di area aksi atas) - audit silang P3.

/** pyPageSizeOther grid ListSuggest. */
export const BARIS_PER_HALAMAN_USULAN = 5

export function jumlahHalaman(n: number, ukuran: number): number {
  return Math.max(1, Math.ceil(n / ukuran))
}

/** Baris halaman `hal` (berbasis 1, dijepit ke rentang yang ada). */
export function irisan<T>(baris: T[], hal: number, ukuran: number): T[] {
  const h = Math.min(Math.max(1, hal), jumlahHalaman(baris.length, ukuran))
  return baris.slice((h - 1) * ukuran, h * ukuran)
}
