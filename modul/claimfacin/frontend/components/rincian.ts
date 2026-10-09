// Panel rinci baris grid masterDetail (expandPane) - UMUM. Claim Non Prop menyebut tiga grid ber-panel satu per satu
// (`modul/claimnonprop/frontend/components/rincian.ts`); Claim Fac In memakai sarang tiga tingkat (objek -> item ->
// estimasi / adjustment) sehingga panel dibaca dari tata itu sendiri: grid ber-`rincian` = P membuka, untuk baris ke-n,
// panel `Layar.panel[P + ':' + jalurGrid + '(' + n + ')']` (`models.KunciPanel`, `backend/models/tata.go`). Panel tidak
// dikirim server = baris itu tidak dapat dibuka. Grid di dalam panel boleh ber-`rincian` sendiri - kuncinya memuat jalur
// baris induknya, jadi sarang berapa pun dalamnya beralamat tanpa protokol khusus:
//
//   est:ClaimData.ObjectList(1)                                   panel objek layar Input Estimasi
//   estitem:ClaimData.ObjectList(1).ObjectItemList(2)             panel item di dalamnya
//   adjdtl:ClaimData.ObjectList(1).ObjectItemList(2).Adjustment(3) panel adjustment layar Input Adjustment
//
// Aksi unsur di dalam panel dikirim ber-`konteks` = kunci panel itu dan `indeks` = baris grid DI DALAM panel itu.

import type { Tata } from '../api'

/** Kunci panel baris ke-n grid berjalur `jalur` ber-`rincian` = `prefiks` (`models.KunciPanel`). */
export function kunciPanel(prefiks: string, jalur: string, n: number): string {
  return `${prefiks}:${jalur}(${n})`
}

/** Kunci panel baris ke-n grid `t`; null bila grid tanpa `rincian` atau server tidak mengirim panel baris itu. */
export function panelBaris(
  t: Pick<Tata, 'rincian' | 'jalur'>,
  n: number,
  panel: Readonly<Record<string, readonly Tata[] | undefined>> | undefined,
): string | null {
  if (!t.rincian || !t.jalur || n < 1) return null
  const kunci = kunciPanel(t.rincian, t.jalur, n)
  return panel?.[kunci] !== undefined ? kunci : null
}

/**
 * Baris terbuka sesudah jumlah baris berubah: baris baru (Add item, "+" adjustment) langsung terbuka supaya isiannya
 * terlihat; baris yang terhapus ditutup.
 */
export function barisTerbuka(buka: number | null, nLalu: number, nKini: number): number | null {
  if (nKini > nLalu) return nKini
  if (buka !== null && buka > nKini) return null
  return buka
}
