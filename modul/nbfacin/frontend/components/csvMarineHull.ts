// CSV objek MARINE HULL (tiket 73 P6, keputusan work owner 11-10-2026 K73-52..62). Template = berkas Pega kiriman work owner
// `D:\migrasi\RNM\DDL\templatecsvmarinehull.csv` (`SetDownloadMasterCSV` IsMarineHull -> TemplateCSVMarineHull; judul VERBATIM,
// urutan apa adanya). Impor Pega = `InsertUploadAneka_act` (`ChooseActivityMappingCSV_Act` IsAneka), mengikat kolom per NAMA judul,
// jadi urutan berkas unggahan bebas. Wajib = seluruh 31 judul template (K73-60); kolom lain diabaikan (pola K74-18).

import { pecahBarisCsv } from './csvLossRecord'
import { PESAN_CSV_OBJEK } from './csvObjek'

const BOM = String.fromCharCode(0xfeff)
const tanpaBom = (t: string) => (t.startsWith(BOM) ? t.slice(1) : t)

/** Judul template Pega, urutan VERBATIM (`templatecsvmarinehull.csv` baris 1). */
export const KOLOM_CSV_MARINE_HULL = [
  'NO', 'TYPE_ADDRESS', 'ADDRESS', 'BUILDING_NO', 'POSTAL_CODE', 'COUNTRY', 'OCCUPATION_NAME', 'OBJECT_DESCRIPTION', 'SECTION',
  'CURRENCY', 'TSI', 'YEAR', 'PLACE_OF_MANUFACTURE', 'OBJECT_TYPE', 'NRT', 'GRT', 'DWT', 'FLAGNAME', 'CONSTRUCTION',
  'CLASSIFICATION', 'COVERAGE_NAME', 'CALCULATE_METHOD', 'DAY', 'PCT_LOADING', 'PCT_SCALE', 'PCT_RATE', 'PREMIUM', 'PCT_DISCOUNT',
  'DISCOUNT', 'PCT_LOL', 'LOL',
] as const

/** Isi tombol Download Template CSV: baris judul template Pega saja (contoh barisnya tidak disebar). */
export const TEMPLATE_CSV_MARINE_HULL = KOLOM_CSV_MARINE_HULL.join(';') + '\r\n'

/** Judul memuat semua 31 kolom (tak peka huruf / spasi tepi; urutan bebas, kolom tambahan diabaikan). */
export function judulCsvMarineHullSah(barisJudul: string): boolean {
  const ada = new Set(pecahBarisCsv(tanpaBom(barisJudul), ';').map((x) => x.trim().toUpperCase()))
  return KOLOM_CSV_MARINE_HULL.every((k) => ada.has(k))
}

/** Judul berkas dari potongan awal saja (tiket 75); '' = sah. */
export async function periksaJudulBerkasMarineHull(berkas: Blob): Promise<string> {
  if (berkas.size === 0) return PESAN_CSV_OBJEK.kosong
  const awal = await berkas.slice(0, 64 * 1024).text()
  const pertama = tanpaBom(awal)
    .split(/\r?\n/)
    .find((b) => b.trim() !== '')
  if (pertama === undefined) return PESAN_CSV_OBJEK.takTerbaca
  return judulCsvMarineHullSah(pertama) ? '' : PESAN_CSV_MARINE_HULL.judul
}

/** Teks sistem baru (bukan korpus). */
export const PESAN_CSV_MARINE_HULL = {
  judul: `Baris judul CSV Marine Hull harus memuat kolom: ${KOLOM_CSV_MARINE_HULL.join(', ')} (Download Template CSV).`,
} as const
