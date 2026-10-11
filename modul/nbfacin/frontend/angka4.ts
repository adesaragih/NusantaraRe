// Tampilan angka NB FacIn: SELALU tepat 4 angka di belakang koma (work owner 06-10-2026: "yang nilai uang, di tampilan,
// buatin yang muncul 4 belakang koma dong, biar seragam semua, walaupun dihitungannya lebih dari 4 belakang koma, namun
// yang tampil hanya 4 saja") - seperti pxNumber 4 desimal Pega ("20.000.000,0000").
//
// `formatNumber` inti membuang nol di ekor dan dipakai modul lain, jadi tidak diubah; di sini hasilnya hanya dilengkapi
// nol sampai 4 desimal. Pembulatan = `formatNumber` (setengah ke atas, eksak di teks - tanpa float). Hanya tampilan:
// nilai yang dikirim ke backend tetap teks lengkap.

import { desimalSah } from '../../../inti/frontend/lib/desimal'
import { formatNumber } from '../../../inti/frontend/lib/format'

/** Jumlah angka di belakang koma pada tampilan. */
export const DESIMAL_TAMPIL = 4

/** "20000000" -> "20.000.000,0000"; "0.123456" -> "0,1235"; kosong -> ""; bukan angka -> apa adanya. */
export function angka4(nilai: string | null | undefined): string {
  const teks = (nilai ?? '').trim()
  if (teks === '' || !desimalSah(teks)) return formatNumber(teks, DESIMAL_TAMPIL)
  const hasil = formatNumber(teks, DESIMAL_TAMPIL)
  const [kiri = '', kanan = ''] = hasil.split(',')
  return `${kiri},${kanan.padEnd(DESIMAL_TAMPIL, '0')}`
}
