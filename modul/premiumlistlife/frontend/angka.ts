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

/**
 * Menyisipkan pemisah ribuan ke teks desimal. Teks yang bukan desimal
 * (kosong, `—`, teks biasa) DIKEMBALIKAN APA ADANYA.
 */
export function pemisahRibuan(teks: string): string {
  const m = POLA_DESIMAL.exec(teks.trim())
  if (m === null) return teks
  const bulat = (m[2] ?? '').replace(/\B(?=(\d{3})+(?!\d))/g, '.')
  // Paling banyak 4 digit di belakang koma, DIPOTONG — bukan dibulatkan
  // (keputusan work owner 03-10-2026: "ambil saja 4 digit"). Hanya tampilan:
  // nilai tersimpan tetap utuh.
  const desimal = (m[3] ?? '').slice(0, 1 + DIGIT_DESIMAL_TAMPIL).replace('.', ',')
  return `${m[1] ?? ''}${bulat}${desimal}`
}
