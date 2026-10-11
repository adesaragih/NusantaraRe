// Fungsi murni tab Spreading Marine Cargo (tiket 74 P3). Aturan baris (Σ % Share ≤ 100, Type Treaty + % Share wajib di tahap UW)
// = FIRE (`TabSpreading` - `cekSpreadingFactIn`); pohonnya Cargo → Coverage.

import { desimalSah } from '../../../../inti/frontend/lib/desimal'
import type { BarisSpreading, CargoSpreading, CoverageSpreadingCargo } from '../api'
import { barisLebih } from './TabSpreading'

/** Posisi satu coverage: cargo, coverage (0-based). */
export interface PosisiCoverageCargo {
  c: number
  n: number
}

export const kunciCoverageCargo = (p: PosisiCoverageCargo) => `cc${p.c}-${p.n}`

export const semuaCoverageCargo = (cargo: CargoSpreading[]): CoverageSpreadingCargo[] => cargo.flatMap((x) => x.coverages)

/** Ganti baris spreading satu coverage tanpa mengubah masukan. */
export function gantiBarisCargo(cargo: CargoSpreading[], p: PosisiCoverageCargo, baris: BarisSpreading[]): CargoSpreading[] {
  return cargo.map((x, c) => (c !== p.c ? x : { ...x, coverages: x.coverages.map((v, n) => (n !== p.n ? v : { ...v, spreading: baris })) }))
}

/** Setiap baris spreading ber-Type Treaty dan % Share angka sah. */
export function barisLengkapCargo(cargo: CargoSpreading[]): boolean {
  return semuaCoverageCargo(cargo).every((v) =>
    v.spreading.every((b) => b.treatyType.trim() !== '' && b.sharePercentage.trim() !== '' && desimalSah(b.sharePercentage.trim())),
  )
}

/** Baris sah untuk dihitung / disimpan: lengkap dan tidak ada coverage dengan Σ % Share > 100. */
export const barisSahCargo = (cargo: CargoSpreading[]) => barisLengkapCargo(cargo) && semuaCoverageCargo(cargo).every((v) => !barisLebih(v.spreading))
