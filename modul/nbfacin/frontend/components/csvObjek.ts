// Upload CSV objek FIRE (tiket 66) - template `D:\migrasi\RNM\DDL\templatecsvfire.csv` dari work owner 08-10-2026 (berkas asli
// Rule-File-Binary `TemplateCSVFire` Pega 55 kolom + 6 kolom wilayah yang ditambah work owner 08-10-2026 16:42 = 61 kolom,
// pemisah ";"). Unggahan berjudul 55 kolom (template Pega lama) tetap diterima - 6 kolom baru kosong. Layar hanya MEMBACA berkas jadi baris bernama-kolom (padanan
// buffer Pega `TempWorkPage.ListB2BHostUpload`, `UploadCSVDataPolis_Act`); pemetaan ke objek / item / coverage, normalisasi angka
// ("." ribuan, "," desimal) dan pencarian master dikerjakan backend saat Save Data (`InsertUploadFire_act`).

import { pecahBarisCsv } from './csvLossRecord'

/** Urutan kolom = baris judul `templatecsvfire.csv` (diuji terhadap berkasnya). */
export const KOLOM_CSV_OBJEK = [
  'NO',
  'NO_OFFER_SLIP',
  'CLASS_OF_BUSINESS',
  'BEGIN_DATE',
  'END_DATE',
  'INSURED_NAME',
  'SOB_NAME',
  'CEDING_CO_NAME',
  'MARKETING_NAME',
  'OCCUPATION_CODE',
  'OCCUPATION_NAME',
  'CLASS_OF_CONSTRUCTION',
  'TYPE_ADDRESS',
  'ADDRESS',
  'BUILDING_NO',
  'POSTAL_CODE',
  'KELURAHAN',
  'OBJECT_TYPE',
  'INTEREST',
  'TOP_RISK',
  'ITEM_TYPE',
  'ITEM_TYPE_1',
  'ITEM_TYPE_2',
  'ITEM_TYPE_3',
  'ITEM_TYPE_4',
  'ITEM_TYPE_5',
  'ITEM_TYPE_6',
  'TSI_OBJECT_ITEM',
  'TSI_OBJECT_ITEM_1',
  'TSI_OBJECT_ITEM_2',
  'TSI_OBJECT_ITEM_3',
  'TSI_OBJECT_ITEM_4',
  'TSI_OBJECT_ITEM_5',
  'TSI_OBJECT_ITEM_6',
  'ACCUMULATION_CODE',
  'ACCUMULATION_NAME',
  'ADJUSTABLE_PCT',
  'CURRENCY',
  'INDEMNITY_MONTH',
  'LOL',
  'RATE_FLEXAS',
  'NET_RATE_FLEXAS',
  'DISCOUNT_RATE_FLEXAS',
  'RATE_RSMDCC',
  'NET_RATE_RSMDCC',
  'DISCOUNT_RATE_RSMDCC',
  'RATE_FWTWD',
  'NET_RATE_FWTWD',
  'DISCOUNT_RATE_FWTWD',
  'RATE_EQ',
  'NET_RATE_EQ',
  'DISCOUNT_RATE_EQ',
  'RATE_OTHERS',
  'NET_RATE_OTHERS',
  'DISCOUNT_RATE_OTHERS',
  // Kolom tambahan work owner (bukan Pega): "Bagian template, dipakai" - POSTAL_CODE dulu, KODE_POS cadangan; KECAMATAN /
  // KABUPATEN_KOTA / PROVINSI membantu mencari RW; SUMBER_KODE_POS / TINGKAT_KEYAKINAN hanya ditampilkan.
  'KODE_POS',
  'KECAMATAN',
  'KABUPATEN_KOTA',
  'PROVINSI',
  'SUMBER_KODE_POS',
  'TINGKAT_KEYAKINAN',
] as const

/** Jumlah kolom template Pega lama (sebelum 6 kolom wilayah) - unggahan berjudul ini tetap diterima. */
export const KOLOM_CSV_PEGA = 55

export type KolomCsvObjek = (typeof KOLOM_CSV_OBJEK)[number]

/** Satu baris CSV: nama kolom -> teks apa adanya (dipangkas). */
export type BarisCsvObjek = Record<KolomCsvObjek, string>

/**
 * Isi Download Template CSV: baris judul saja. Baris contoh di berkas Pega berisi nama tertanggung / marketing / ceding,
 * jadi tidak disalin ke aplikasi `[keputusan agent]`.
 */
export const TEMPLATE_CSV_OBJEK = KOLOM_CSV_OBJEK.join(';') + '\r\n'

export interface HasilCsvObjek {
  baris: BarisCsvObjek[]
  galat: string
}

/** Pesan teks sistem baru (bukan korpus). */
export const PESAN_CSV_OBJEK = {
  kosong: 'Berkas CSV kosong.',
  judul: 'Baris judul CSV tidak sama dengan template (Download Template CSV).',
  tanpaBaris: 'Berkas CSV tidak berisi baris data.',
  terlaluBanyak: (n: number) => `Berkas CSV berisi ${n} baris; paling banyak ${BATAS_BARIS_CSV_OBJEK} baris per unggahan.`,
} as const

/** Batas baris per unggahan (work owner 08-10-2026: "jangan dibatasin 500. 3000 aja buat batas nya"); backend tetap penentu. */
export const BATAS_BARIS_CSV_OBJEK = 3000

/** Baca teks CSV template FIRE: judul wajib sama persis (tak peka huruf / spasi tepi); baris kosong seluruhnya dilewati. */
export function bacaCsvObjek(teks: string): HasilCsvObjek {
  const baris = teks
    .replace(/^\uFEFF/, '')
    .split(/\r?\n/)
    .filter((b) => b.trim() !== '')
  if (baris.length === 0) return { baris: [], galat: PESAN_CSV_OBJEK.kosong }
  const judul = pecahBarisCsv(baris[0] ?? '', ';').map((x) => x.trim().toUpperCase())
  while (judul.length > 0 && judul[judul.length - 1] === '') judul.pop()
  const cocok = (n: number) => judul.length === n && judul.every((j, i) => j === KOLOM_CSV_OBJEK[i])
  if (!cocok(KOLOM_CSV_OBJEK.length) && !cocok(KOLOM_CSV_PEGA)) {
    return { baris: [], galat: PESAN_CSV_OBJEK.judul }
  }
  const hasil: BarisCsvObjek[] = []
  for (const b of baris.slice(1)) {
    const sel = pecahBarisCsv(b, ';').map((x) => x.trim())
    if (sel.every((x) => x === '')) continue
    const r = {} as BarisCsvObjek
    KOLOM_CSV_OBJEK.forEach((k, i) => {
      r[k] = sel[i] ?? ''
    })
    hasil.push(r)
  }
  if (hasil.length === 0) return { baris: [], galat: PESAN_CSV_OBJEK.tanpaBaris }
  if (hasil.length > BATAS_BARIS_CSV_OBJEK) return { baris: [], galat: PESAN_CSV_OBJEK.terlaluBanyak(hasil.length) }
  return { baris: hasil, galat: '' }
}
