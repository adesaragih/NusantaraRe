// Fungsi murni tab Spreading keluarga Aneka (tiket 73 P3). Aturan baris (Σ % Share ≤ 100, Type Treaty + % Share wajib di tahap UW)
// = FIRE (`TabSpreading` - `cekSpreadingFactIn`), hanya pohonnya yang berbeda: lokasi → okupasi | Golf → objek → coverage.

import { desimalSah } from '../../../../inti/frontend/lib/desimal'
import type { BarisSpreading, CoverageSpreadingAneka, LokasiSpreadingAneka, ObjekSpreadingAneka } from '../api'
import { barisLebih } from './TabSpreading'

/** Posisi satu coverage: lokasi, okupasi (null = Golf), objek, coverage. Indeks 0-based. */
export interface PosisiCoverageSpreading {
  l: number
  k: number | null
  j: number
  n: number
}

export const kunciCoverageSpreading = (p: PosisiCoverageSpreading) => `c${p.l}-${p.k ?? 'g'}-${p.j}-${p.n}`

/** Setiap coverage di pohon. */
export function semuaCoverageSpreading(lokasi: LokasiSpreadingAneka[]): CoverageSpreadingAneka[] {
  const dariObjek = (o: ObjekSpreadingAneka[]) => o.flatMap((x) => x.coverages)
  return lokasi.flatMap((l) => [...l.okupasi.flatMap((k) => dariObjek(k.objek)), ...dariObjek(l.objekGolf)])
}

/** Ganti baris spreading satu coverage tanpa mengubah masukan. */
export function gantiBarisAneka(lokasi: LokasiSpreadingAneka[], p: PosisiCoverageSpreading, baris: BarisSpreading[]): LokasiSpreadingAneka[] {
  const objek = (daftar: ObjekSpreadingAneka[]) =>
    daftar.map((o, j) => (j !== p.j ? o : { ...o, coverages: o.coverages.map((c, n) => (n !== p.n ? c : { ...c, spreading: baris })) }))
  return lokasi.map((l, a) => {
    if (a !== p.l) return l
    if (p.k === null) return { ...l, objekGolf: objek(l.objekGolf) }
    return { ...l, okupasi: l.okupasi.map((k, b) => (b !== p.k ? k : { ...k, objek: objek(k.objek) })) }
  })
}

/** Setiap baris spreading ber-Type Treaty dan % Share angka sah. */
export function barisLengkapAneka(lokasi: LokasiSpreadingAneka[]): boolean {
  return semuaCoverageSpreading(lokasi).every((c) =>
    c.spreading.every((b) => b.treatyType.trim() !== '' && b.sharePercentage.trim() !== '' && desimalSah(b.sharePercentage.trim())),
  )
}

/** Baris sah untuk dihitung / disimpan: lengkap dan tidak ada coverage dengan Σ % Share > 100. */
export const barisSahAneka = (lokasi: LokasiSpreadingAneka[]) =>
  barisLengkapAneka(lokasi) && semuaCoverageSpreading(lokasi).every((c) => !barisLebih(c.spreading))
