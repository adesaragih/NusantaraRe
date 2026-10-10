// Isian tanggal NB Treaty In (work owner 10-10-2026: "supaya bisa di copy paste dan di ketik lancar") - pola
// `ketikTanggal.ts` Claim Prop / `TanggalDMY` nbfacin, disalin (bukan diimpor, batas modul). `<input type="date">`
// bawaan mengikuti bahasa browser (mm/dd/yyyy), diketik per segmen, dan tidak menerima teks tempelan; di sini pengguna
// cukup mengetik angka (pemisah "-" disisipkan otomatis) atau menempel tanggal utuh. Tampilan dd-mm-yyyy = format
// tanggal inti (NFR-14). Nilai halaman TETAP `YYYY-MM-DD`, sama dengan isian bawaan sebelumnya. Semua konversi
// manipulasi teks tanpa `Date`, supaya zona waktu tidak ikut campur.

import { keInputTanggal } from '../../../inti/frontend/lib/tanggalInput'

/** Panjang teks isian lengkap `dd-mm-yyyy`. */
export const PANJANG_TANGGAL = 10

/** Contoh bentuk isian (placeholder). */
export const CONTOH_TANGGAL = 'dd-mm-yyyy'

const p2 = (s: string) => s.padStart(2, '0')

/** Nama bulan tempelan (Inggris / Indonesia, tiga huruf pertama) -> nomor bulan. */
const BULAN: Readonly<Record<string, string>> = {
  jan: '01',
  feb: '02',
  mar: '03',
  apr: '04',
  may: '05',
  mei: '05',
  jun: '06',
  jul: '07',
  aug: '08',
  agu: '08',
  ags: '08',
  sep: '09',
  oct: '10',
  okt: '10',
  nov: '11',
  dec: '12',
  des: '12',
}

/** Nilai halaman -> tampilan `dd-mm-yyyy`. Tidak terbaca sebagai tanggal = apa adanya. */
export function tampilTanggal(v: string): string {
  const s = v.trim()
  if (s === '') return ''
  const iso = keInputTanggal(s)
  return iso === '' ? s : `${iso.slice(8, 10)}-${iso.slice(5, 7)}-${iso.slice(0, 4)}`
}

/**
 * Tanggal UTUH berpemisah -> `dd-mm-yyyy`; selain itu null. Dipakai untuk tempelan dan ketikan lengkap:
 * `yyyy-mm-dd` (ISO, boleh berjam), `d/m/yyyy` / `d-m-yyyy` / `d.m.yyyy` (hari dulu - format aplikasi), dan
 * `1 Jul 2024` / `01-Jul-2024`. Tanggal yang tidak ada (31-02) tetap dikembalikan; `halamanTanggal` menolaknya.
 */
export function tanggalUtuh(s: string): string | null {
  const t = s.trim()
  const iso = /^(\d{4})[-/.](\d{1,2})[-/.](\d{1,2})(?:$|[T\s])/.exec(t)
  if (iso) return `${p2(iso[3]!)}-${p2(iso[2]!)}-${iso[1]}`
  const dmy = /^(\d{1,2})[-/.\s](\d{1,2})[-/.\s](\d{4})$/.exec(t)
  if (dmy) return `${p2(dmy[1]!)}-${p2(dmy[2]!)}-${dmy[3]}`
  const nama = /^(\d{1,2})[-/.\s]+([A-Za-z]{3,})\.?[-/.,\s]+(\d{4})$/.exec(t)
  const bulan = nama ? BULAN[nama[2]!.slice(0, 3).toLowerCase()] : undefined
  if (nama && bulan) return `${p2(nama[1]!)}-${bulan}-${nama[3]}`
  return null
}

/** Ketikan bebas -> tanggal utuh bila terbaca, selain itu angka saja dengan pemisah disisipkan (paling panjang 8). */
export function rapikanTanggal(s: string): string {
  const utuh = tanggalUtuh(s)
  if (utuh !== null) return utuh
  const d = s.replace(/\D/g, '').slice(0, 8)
  let out = d.slice(0, 2)
  if (d.length > 2) out += '-' + d.slice(2, 4)
  if (d.length > 4) out += '-' + d.slice(4, 8)
  return out
}

/** Isian lengkap dan tanggalnya ada -> nilai halaman `YYYY-MM-DD`; selain itu kosong. */
export function halamanTanggal(teks: string): string {
  if (!/^\d{2}-\d{2}-\d{4}$/.test(teks)) return ''
  return keInputTanggal(teks)
}
