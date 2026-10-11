// Angka uang di layar modul PremiumList Life — DENGAN pemisah ribuan
// (permintaan work owner 03-10-2026: "tambahkan separator, karena ini untuk
// nominal uang").
//
// ⛔ TEKS, BUKAN `Number` (ADR-U-0003). Pemisah disisipkan ke teks desimal
// apa adanya: tidak ada pembulatan, tidak ada angka di belakang koma yang
// hilang, dan premi delapan angka desimal tetap delapan angka desimal.
//
// Bentuk Indonesia (keputusan work owner 03-10-2026): TITIK pemisah ribuan,
// KOMA desimal — `2200000.50` → `2.200.000,50`, `23400155` → `23.400.155`.
// ⚠️ Hanya tampilan; teks dari server tetap ber-titik desimal.

/** Digit desimal paling banyak yang ditampilkan (sisanya dipotong). */
export const DIGIT_DESIMAL_TAMPIL = 4

/** Teks desimal bertanda opsional: `-1234.5`, `0`, `123`. */
const POLA_DESIMAL = /^(-?)(\d+)(\.\d+)?$/

/** Bentuk desimal lain yang dapat dikirim Oracle/Go: `+0.72`, `.72`, `-.5`, `7.2E-01`. */
const POLA_DESIMAL_LONGGAR = /^([+-]?)(\d*)(?:\.(\d*))?(?:[eE]([+-]?\d+))?$/

/**
 * Teks desimal → bentuk baku `-?\d+(\.\d+)?` lewat olah TEKS (tanpa `Number`), atau
 * null bila bukan desimal. Satu bentuk untuk semua sumber, supaya angka yang sama
 * tidak tampil `0,72` di satu sel dan `0.72` di sel lain (keputusan work owner 07-10-2026).
 */
export function desimalBaku(teks: string): string | null {
  const m = POLA_DESIMAL_LONGGAR.exec(teks.trim())
  if (m === null) return null
  const bulat = m[2] ?? ''
  const pecahan = m[3] ?? ''
  if (bulat === '' && pecahan === '') return null
  let digit = bulat + pecahan
  let titik = bulat.length + Number.parseInt(m[4] ?? '0', 10)
  if (titik < 0) {
    digit = '0'.repeat(-titik) + digit
    titik = 0
  }
  if (titik > digit.length) digit = digit + '0'.repeat(titik - digit.length)
  const kiri = digit.slice(0, titik).replace(/^0+(?=\d)/, '') || '0'
  const kanan = digit.slice(titik)
  const tanda = m[1] === '-' ? '-' : ''
  return `${tanda}${kiri}${kanan === '' ? '' : `.${kanan}`}`
}

/**
 * Menyisipkan pemisah ribuan ke teks desimal. Teks yang bukan desimal
 * (kosong, `—`, teks biasa) DIKEMBALIKAN APA ADANYA.
 */
export function pemisahRibuan(teks: string): string {
  const m = POLA_DESIMAL.exec(desimalBaku(teks) ?? teks.trim())
  if (m === null) return teks
  const bulat = (m[2] ?? '').replace(/\B(?=(\d{3})+(?!\d))/g, '.')
  // Paling banyak 4 digit di belakang koma, DIPOTONG — bukan dibulatkan
  // (keputusan work owner 03-10-2026: "ambil saja 4 digit"). Hanya tampilan:
  // nilai tersimpan tetap utuh.
  const desimal = (m[3] ?? '').slice(0, 1 + DIGIT_DESIMAL_TAMPIL).replace('.', ',')
  return `${m[1] ?? ''}${bulat}${desimal}`
}
