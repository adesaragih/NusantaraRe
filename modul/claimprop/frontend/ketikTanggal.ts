// Isian tanggal Claim Prop (work owner 08-10-2026: "format tanggalnya kenapa susah banget di ketik ya? perbaiki") - pola
// `TanggalDMY` nbfacin, disalin (bukan diimpor, batas modul). `<input type="date">` bawaan mengikuti bahasa browser
// (mm/dd/yyyy) dan diketik per segmen; di sini pengguna cukup mengetik angka, pemisah "-" (dan spasi / ":" untuk jam)
// disisipkan otomatis. Tampilan dd-mm-yyyy = format tanggal inti (NFR-14). Nilai halaman TETAP `YYYY-MM-DD` (tanggal)
// dan `YYYY-MM-DD HH:MM:SS` (tanggal-waktu), sama dengan isian bawaan sebelumnya. Semua konversi manipulasi teks tanpa
// `Date`, supaya zona waktu tidak ikut campur.

import { keInputTanggal } from '../../../inti/frontend/lib/tanggalInput'

export type JenisTanggal = 'tanggal' | 'tanggal-waktu'

/** Jumlah digit isian lengkap: ddmmyyyy / ddmmyyyyhhmm. */
const DIGIT = { tanggal: 8, 'tanggal-waktu': 12 } as const

/** Panjang teks isian lengkap: `dd-mm-yyyy` / `dd-mm-yyyy hh:mm`. */
export const PANJANG = { tanggal: 10, 'tanggal-waktu': 16 } as const

/** Contoh bentuk isian (placeholder). */
export const CONTOH = { tanggal: 'dd-mm-yyyy', 'tanggal-waktu': 'dd-mm-yyyy hh:mm' } as const

/** Nilai halaman -> tampilan `dd-mm-yyyy` (+ ` hh:mm`). Tidak terbaca sebagai tanggal = apa adanya. */
export function tampilTanggal(v: string, jenis: JenisTanggal): string {
  const s = v.trim()
  if (s === '') return ''
  const iso = keInputTanggal(s)
  if (iso === '') return s
  const dmy = `${iso.slice(8, 10)}-${iso.slice(5, 7)}-${iso.slice(0, 4)}`
  if (jenis === 'tanggal') return dmy
  const jam = /[T ](\d{2}):(\d{2})/.exec(s)
  return jam ? `${dmy} ${jam[1]}:${jam[2]}` : dmy
}

/** Ketikan bebas -> angka saja, pemisah disisipkan, paling panjang isian lengkap. */
export function rapikanTanggal(s: string, jenis: JenisTanggal): string {
  const d = s.replace(/\D/g, '').slice(0, DIGIT[jenis])
  let out = d.slice(0, 2)
  if (d.length > 2) out += '-' + d.slice(2, 4)
  if (d.length > 4) out += '-' + d.slice(4, 8)
  if (d.length > 8) out += ' ' + d.slice(8, 10)
  if (d.length > 10) out += ':' + d.slice(10, 12)
  return out
}

/** Isian lengkap dan tanggalnya ada (jam 00-23, menit 00-59) -> nilai halaman; selain itu kosong. */
export function halamanTanggal(teks: string, jenis: JenisTanggal): string {
  const m = /^(\d{2})-(\d{2})-(\d{4})(?: (\d{2}):(\d{2}))?$/.exec(teks)
  if (!m || teks.length !== PANJANG[jenis]) return ''
  const iso = keInputTanggal(`${m[1]}-${m[2]}-${m[3]}`)
  if (iso === '') return ''
  if (jenis === 'tanggal') return iso
  if (Number(m[4]) > 23 || Number(m[5]) > 59) return ''
  return `${iso} ${m[4]}:${m[5]}:00`
}
