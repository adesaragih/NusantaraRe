// CSV objek MARINE CARGO (tiket 74 P6, keputusan work owner 11-10-2026 K74-14..17). Template = berkas Pega kiriman work owner
// `D:\migrasi\RNM\DDL\templatecsvmarinecargo.csv` (K74-14; judul VERBATIM, urutan apa adanya). Pega mengikat kolom per NAMA judul
// (`UploadCSVDataPolis_Act` `pxUploadCSVResults`), jadi urutan berkas unggahan bebas. Wajib = kolom template yang DIBACA impor Pega
// `InsertUploadMarineCargo_act` / `ProtectUploadMarineCargo_act` [terverifikasi]: INVOICEDATE ada di template tetapi tidak dibaca
// (grep 0) -> opsional; COVERAGE_NAME dibaca Insert tetapi tidak ada di template -> opsional.

import { pecahBarisCsv } from './csvLossRecord'
import { PESAN_CSV_OBJEK } from './csvObjek'

const BOM = String.fromCharCode(0xfeff)
const tanpaBom = (t: string) => (t.startsWith(BOM) ? t.slice(1) : t)

/** Judul template Pega, urutan VERBATIM (`templatecsvmarinecargo.csv` baris 1). */
export const KOLOM_CSV_CARGO = [
  'TRADING', 'ETD', 'DEPARTURE', 'DESTINATION', 'INVOICE', 'INVOICEDATE', 'BLNUMBER', 'CONVEYENCE', 'CONVEYENCE_NAME', 'DATASHIP',
  'GOODS', 'GOODS_NAME', 'PACKING', 'AGENT', 'COVERAGE_IDX', 'CURRENCY', 'TSI', 'TOP_RISK',
] as const

/** Kolom template yang tidak dibaca impor Pega - boleh tidak ada. */
export const KOLOM_CSV_CARGO_OPSIONAL = ['INVOICEDATE'] as const

/** Kolom wajib = template tanpa yang opsional (17). COVERAGE_NAME (dibaca Pega, di luar template) juga opsional. */
export const KOLOM_CSV_CARGO_WAJIB = KOLOM_CSV_CARGO.filter((k) => !(KOLOM_CSV_CARGO_OPSIONAL as readonly string[]).includes(k))

/** Isi tombol Download Template CSV: baris judul template Pega saja (contoh barisnya tidak disebar). */
export const TEMPLATE_CSV_CARGO = KOLOM_CSV_CARGO.join(';') + '\r\n'

/** Judul memuat semua kolom wajib (tak peka huruf / spasi tepi; urutan bebas, kolom tambahan diabaikan seperti Pega). */
export function judulCsvCargoSah(barisJudul: string): boolean {
  const ada = new Set(pecahBarisCsv(tanpaBom(barisJudul), ';').map((x) => x.trim().toUpperCase()))
  return KOLOM_CSV_CARGO_WAJIB.every((k) => ada.has(k))
}

/** Judul berkas dari potongan awal saja (tiket 75); '' = sah. */
export async function periksaJudulBerkasCargo(berkas: Blob): Promise<string> {
  if (berkas.size === 0) return PESAN_CSV_OBJEK.kosong
  const awal = await berkas.slice(0, 64 * 1024).text()
  const pertama = tanpaBom(awal)
    .split(/\r?\n/)
    .find((b) => b.trim() !== '')
  if (pertama === undefined) return PESAN_CSV_OBJEK.takTerbaca
  return judulCsvCargoSah(pertama) ? '' : PESAN_CSV_CARGO.judul
}

/** Teks sistem baru (bukan korpus). */
export const PESAN_CSV_CARGO = {
  judul: `Baris judul CSV Marine Cargo harus memuat kolom: ${KOLOM_CSV_CARGO_WAJIB.join(', ')} (Download Template CSV).`,
} as const
