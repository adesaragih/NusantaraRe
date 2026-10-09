// Lebar kolom grid layar Adjustment — `pyWidth` ekspor sebagai PERBANDINGAN.
//
// ⛔ Satu tempat untuk seluruh grid layar ini (daftar penyesuaian dan picker
// tombol Add): tata letak kita responsif, Pega tidak, jadi pikselnya menjadi
// persen.

/** Lebar ekspor → persen kolom ke-`i`. */
export function persenLebar(lebar: readonly number[], i: number): string {
  const total = lebar.reduce((a, b) => a + b, 0)
  return `${(((lebar[i] ?? 0) / total) * 100).toFixed(2)}%`
}
