// Tanda angka TEKS halaman (uang, persen) tanpa `Number`/float - nilai halaman
// selalu teks desimal dari backend (apd, `models.TeksSkalarJSON`); spec §5.6:
// uang dan persentase tidak pernah float. Hanya untuk SYARAT tampil/kunci -
// setiap perhitungan dijalankan backend. Teks kosong = 0 (aritmetika Pega).

/** Angka desimal bertitik dengan tanda opsional (`-12.50`, `.5`, `0`). */
const ANGKA = /^[+-]?(\d+\.?\d*|\.\d+)$/

/** Tanda angka teks: -1, 0, 1; `null` bila bukan angka. */
export function tandaTeks(s: string): -1 | 0 | 1 | null {
  const t = s.trim()
  if (t === '') return 0
  if (!ANGKA.test(t)) return null
  if (!/[1-9]/.test(t)) return 0
  return t.startsWith('-') ? -1 : 1
}

/** `= 0` (kosong dihitung 0). */
export const nolTeks = (s: string) => tandaTeks(s) === 0

/** `> 0`. */
export const positifTeks = (s: string) => tandaTeks(s) === 1

/** `< 0`. */
export const negatifTeks = (s: string) => tandaTeks(s) === -1
