// Penyajian nilai layar realisasi NB Treaty In - K14 (OQ 7, spec AC 86).
//
// Sumber setiap format: setelan sel Section Pega (`pyModes` baris 2 = mode
// baca, baris 1 = mode sunting), dikutip per medan di `medan.ts`:
//
//   pxNumber `pyDecimalPlaces` N     -> tepat N desimal (dibulatkan, dipadankan)
//   pyFormatType number + `pySeparators false` -> tanpa pemisah ribuan
//   `pyShowReadonlyFormatting` true  -> isian tersunting tampil berformat
//                                       selama tidak difokus (`formatSaatSunting`)
//   pxCurrency / pyFormatType number TANPA `pyDecimalPlaces` -> jumlah desimal
//                                       TIDAK terbaca: pola modul lain
//                                       (`inti/frontend/lib/format.ts`)
//
// Yang TIDAK terbaca dari korpus dan karenanya mengikuti pola inti: pemisah
// ribuan/desimal (locale operator tidak terekspor -> titik ribuan, koma
// desimal) dan mode pembulatan (setengah ke atas pada digit, nol float).
// Tanggal: satu format seluruh sistem (AC 33) -> `formatDate` inti.
//
// ⛔ Nilai tidak pernah berubah di sini; yang berubah hanya bentuk bacanya
// (AC 24). Nol perhitungan uang.

import { DESIMAL_TAK_DIBATASI, formatDate, formatNumber } from '../../../inti/frontend/lib/format'

/** Format angka satu sel. */
export interface FormatAngka {
  /** `pyDecimalPlaces` mode baca. Tidak ada = tak terbaca -> pola inti. */
  desimal?: number
  /** `pySeparators`; bawaan kontrol Pega = true. */
  ribuan?: boolean
  /** `pyShowReadonlyFormatting` mode sunting. */
  formatSaatSunting?: boolean
}

/** Sajian satu sel: angka berformat, atau tanggal. */
export type Sajian = FormatAngka | 'tanggal'

const ANGKA = /^[+-]?(\d+\.?\d*|\.\d+)$/

/** Teks tampilan sebuah angka menurut format selnya. Teks bukan angka apa adanya. */
export function sajikanAngka(nilai: string | null | undefined, f: FormatAngka): string {
  const t = (nilai ?? '').trim()
  if (t === '' || !ANGKA.test(t)) return t
  let s = formatNumber(t, f.desimal ?? DESIMAL_TAK_DIBATASI)
  if (f.ribuan === false) s = s.replace(/\./g, '')
  if (f.desimal !== undefined && f.desimal > 0) {
    const [bulat, pecahan = ''] = s.split(',')
    s = `${bulat},${pecahan.padEnd(f.desimal, '0')}`
  }
  // "-0,00": pembulatan menghabiskan seluruh digit berarti - bukan negatif.
  if (/^-[0.,]*$/.test(s)) s = s.slice(1)
  return s
}

/** Teks tampilan satu nilai menurut sajian selnya (tanpa sajian = apa adanya). */
export function sajikan(nilai: string | null | undefined, s: Sajian | undefined): string {
  if (s === undefined) return nilai ?? ''
  if (s === 'tanggal') return formatDate(nilai ?? '')
  return sajikanAngka(nilai, s)
}
