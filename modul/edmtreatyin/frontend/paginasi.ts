// Asal: salinan `modul/nbtreatyin/frontend/paginasi.ts` (06-10-2026), kelas `nbti__` -> `edmt__`.
// Paginasi grid ber-`pyGridProps/pyPageMode=Numeric` - satu halaman berisi
// `ukuran` baris, halaman berbasis 1.
//
// EDM Treaty In memakainya untuk grid `.SuggestList` `Section/ListSuggestEDM` (pyPageSizeOther 5), grid portal
// `Section/SFAPortal_Endorsement_Treaty` (pyPageSize 50) dan grid popup `Section/BusinessAndSOBListEDM`
// (pyPageSize 50). Tombolnya satu komponen (`components/Paginasi.tsx`).

/** pyPageSizeOther grid ListSuggestEDM. */
export const BARIS_PER_HALAMAN_USULAN = 5

/** pyPageSize grid portal `SFAPortal_Endorsement_Treaty` dan popup `BusinessAndSOBListEDM`. */
export const BARIS_PER_HALAMAN_GRID = 50

export function jumlahHalaman(n: number, ukuran: number): number {
  return Math.max(1, Math.ceil(n / ukuran))
}

/** Halaman `hal` (berbasis 1) dijepit ke rentang 1..`jumlah`. */
export function halamanTerjepit(hal: number, jumlah: number): number {
  return Math.min(Math.max(1, hal), jumlah)
}

/** Baris halaman `hal` (berbasis 1, dijepit ke rentang yang ada). */
export function irisan<T>(baris: T[], hal: number, ukuran: number): T[] {
  const h = halamanTerjepit(hal, jumlahHalaman(baris.length, ukuran))
  return baris.slice((h - 1) * ukuran, h * ukuran)
}
