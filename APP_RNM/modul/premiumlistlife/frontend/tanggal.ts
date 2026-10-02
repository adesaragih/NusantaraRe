// Tanggal di layar modul PremiumList Life — SELALU dd/mm/yyyy (permintaan work
// owner 02-10-2026: "tampilkan semua format tanggal dd/mm/yyyy").
//
// ⛔ Hanya BENTUK TAMPIL. Kabel ke server tetap `YYYY-MM-DD` (masukan) dan ISO
// (jawaban): mengubah kontraknya demi tampilan akan membuat dua bentuk tanggal
// berselisih di jalur yang sama.
//
// ⛔ Tanpa `Date` lokal dan tanpa zona: tanggal tanpa jam yang melewati zona
// setempat dapat bergeser satu hari. Semua dikerjakan sebagai TEKS.

/** `YYYY-MM-DD` di awal teks, dengan atau tanpa jam (`T`/spasi) sesudahnya. */
const POLA_ISO = /^(\d{4})-(\d{2})-(\d{2})(?:[T ](\d{2}):(\d{2}))?/

/**
 * Teks tanggal apa pun dari server → `dd/mm/yyyy`.
 *
 * Menerima `2026-10-01`, `2026-10-01T00:00:00Z`, dan `2026-10-01 00:00:00`.
 * Teks yang bukan tanggal ISO (mis. STNC `31/10/2026`, atau nilai bukan
 * tanggal) DIKEMBALIKAN APA ADANYA. Kosong / `null` → `''`. Tanggal nol Go
 * (`0001-01-01`) dianggap kosong.
 */
export function tanggalTampil(nilai: string | null | undefined): string {
  const s = (nilai ?? '').trim()
  if (s === '') return ''
  const m = POLA_ISO.exec(s)
  if (m === null) return s
  if (m[1] === '0001') return ''
  return `${m[3] ?? ''}/${m[2] ?? ''}/${m[1] ?? ''}`
}

/** Seperti `tanggalTampil`, dengan jam `HH:MM` bila ada: `dd/mm/yyyy HH:MM`. */
export function tanggalJamTampil(nilai: string | null | undefined): string {
  const s = (nilai ?? '').trim()
  const m = POLA_ISO.exec(s)
  const tgl = tanggalTampil(s)
  if (m === null || tgl === '' || m[4] === undefined) return tgl
  return `${tgl} ${m[4]}:${m[5] ?? '00'}`
}

/**
 * `dd/mm/yyyy` → `YYYY-MM-DD`, atau `null` bila bentuknya salah ATAU
 * tanggalnya tidak ada di kalender (31/02/2026).
 */
export function isoDariTampil(teks: string): string | null {
  const m = /^(\d{2})\/(\d{2})\/(\d{4})$/.exec(teks.trim())
  if (m === null) return null
  const [h, b, t] = [Number(m[1]), Number(m[2]), Number(m[3])]
  if (b < 1 || b > 12 || h < 1) return null
  // Hari terakhir bulan itu, dihitung dalam UTC — tanpa zona setempat.
  const akhir = new Date(Date.UTC(t, b, 0)).getUTCDate()
  if (h > akhir) return null
  return `${m[3] ?? ''}-${m[2] ?? ''}-${m[1] ?? ''}`
}

/**
 * Topeng ketik dd/mm/yyyy: hanya angka, garis miring disisipkan otomatis
 * (`13122024` → `13/12/2024`), paling banyak delapan angka.
 */
export function topengTanggal(teks: string): string {
  const d = teks.replace(/\D/g, '').slice(0, 8)
  if (d.length <= 2) return d
  if (d.length <= 4) return `${d.slice(0, 2)}/${d.slice(2)}`
  return `${d.slice(0, 2)}/${d.slice(2, 4)}/${d.slice(4)}`
}
