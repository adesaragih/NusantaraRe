// Penyaring ketikan kotak ANGKA yang menampilkan nilai MENTAH (bentuk kabel:
// titik desimal, tanpa pemisah ribuan) — permintaan pemakai 8 Oktober 2026:
// "inputan yg khusus angka tidak boleh input karakter lain".
//
// ⛔ Hanya untuk kotak yang ekspornya kontrol Number (`pxNumber`/`pxInteger`)
// dan yang menampilkan nilai mentah saat diketik. Kotak `FieldAngka` (inti)
// sudah menyaring sendiri; kotak teks (`pxTextInput`, mis. `CoInShare`
// `>=30% up to < 50%`) TIDAK boleh disaring.

/**
 * Pecahan: digit, SATU pemisah desimal, dan `-` di depan. Koma yang diketik
 * menjadi titik (bentuk yang backend terima); pemisah kedua dan seterusnya
 * dibuang, begitu pula huruf dan simbol.
 *
 * Bulat (`bulat = true`): digit saja; yang diketik sesudah pemisah desimal
 * dibuang, bukan disambung (`10,5` hari → `10`, bukan `105`).
 */
export function saringAngka(teks: string, bulat = false): string {
  const t = teks.trim()
  if (bulat) {
    const batas = t.search(/[.,]/)
    return (batas >= 0 ? t.slice(0, batas) : t).replace(/\D/g, '')
  }
  const minus = t.startsWith('-') ? '-' : ''
  const bersih = t.replace(/,/g, '.').replace(/[^\d.]/g, '')
  const i = bersih.indexOf('.')
  return minus + (i < 0 ? bersih : bersih.slice(0, i + 1) + bersih.slice(i + 1).replace(/\./g, ''))
}
