// Asal: salinan `modul/nbtreatyin/frontend/sajian.ts` (06-10-2026), kelas `nbti__` -> `edmt__`.
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
  /** Nilai nol atau kosong tampil "0" (perintah work owner 06-10-2026, bagian uang OGP / ONP). */
  nolPolos?: boolean
}

/** Sajian satu sel: angka berformat, atau tanggal. */
export type Sajian = FormatAngka | 'tanggal'

/** pxNumber tanpa `pyDecimalPlaces` - pola angka inti (dipakai grid NonProp dan `InstallmentList`). */
export const POLA_INTI: Sajian = {}

const ANGKA = /^[+-]?(\d+\.?\d*|\.\d+)$/
/** Angka bernilai nol (`0`, `-0`, `0.000`, `.0`). */
const NOL = /^[+-]?(0+\.?0*|\.0+)$/

/** Kelas sel tabel untuk nilai angka mentah: rata kanan, digit sama lebar (perintah work owner 06-10-2026). */
export function kelasAngka(v: string | null | undefined): string | undefined {
  return ANGKA.test((v ?? '').trim()) ? 'edmt__angka' : undefined
}

/** Kosong atau bernilai nol - isian angka menampilkannya sebagai placeholder "0" (perintah work owner 06-10-2026). */
export function nilaiNol(v: string | null | undefined): boolean {
  const t = (v ?? '').trim()
  return t === '' || NOL.test(t)
}

/** Teks tampilan sebuah angka menurut format selnya. Teks bukan angka apa adanya. */
export function sajikanAngka(nilai: string | null | undefined, f: FormatAngka): string {
  const t = (nilai ?? '').trim()
  if (f.nolPolos && (t === '' || NOL.test(t))) return '0'
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
/**
 * Tanggal + jam kolom Date grid SuggestList (perintah work owner 06-10-2026: "HARUSNYA PAKE YG ATAS AJA, TAPI
 * TAMBAHIN JAM NYA" - panel History dibuang). Tanggal tetap format seluruh sistem (`formatDate` inti), jam 24 apa
 * adanya dari cap waktu `YYYY-MM-DD HH:MI[:SS]`; nilai tanpa jam = tanggal saja.
 */
export function sajikanTanggalJam(nilai: string | null | undefined): string {
  const t = (nilai ?? '').trim()
  const tgl = formatDate(t)
  const jam = /^\d{4}-\d{2}-\d{2}[ T](\d{2}):(\d{2})(?::(\d{2}))?/.exec(t)
  if (tgl === '' || jam === null) return tgl
  return `${tgl} ${jam[1]}:${jam[2]}:${jam[3] ?? '00'}`
}

export function sajikan(nilai: string | null | undefined, s: Sajian | undefined): string {
  if (s === undefined) return nilai ?? ''
  if (s === 'tanggal') return formatDate(nilai ?? '')
  return sajikanAngka(nilai, s)
}
