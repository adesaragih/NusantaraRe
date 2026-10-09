// Aturan layar Disease Life - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/models`,
// `backend/services`): wajib, huruf besar, batas kolom, ICD kembar, hak View only, ID dari sequence.

import type { KolomUrut, Saringan } from './api'
import { DSL } from './labels'

/** `pyPageSize` 10 (`InboxDisease` b5169) - backend `models.UkuranHalaman`. */
export const UKURAN_HALAMAN = 10

/** Lebar `DISEASE_LIFE.ICD_CODE` VARCHAR2(100) (byte) - backend `models.BatasICDCode` (di luar XML: batas kolom). */
export const BATAS_ICD = 100

/** Lebar `DISEASE_LIFE.DISEASE` VARCHAR2(1000) (byte) - backend `models.BatasDisease` (di luar XML: batas kolom). */
export const BATAS_DISEASE = 1000

/** Saringan awal: tanpa filter, ID MENURUN (`pySortType` DESC b5012, `pySortOrder` 1 b5017), halaman 1. */
export const SARINGAN_AWAL: Saringan = { icdCode: '', disease: '', urut: 'id', naik: false, halaman: 1 }

/** Klik kepala kolom: kolom yang sama membalik arah; kolom lain dipilih dengan arah bawaannya (ID menurun, ICD Code
 * menaik); halaman kembali 1. Kolom Disease tidak dapat diurutkan (`pyColumnSorting` false b5058). */
export function pilihUrut(s: Saringan, kolom: KolomUrut): Saringan {
  if (s.urut === kolom) return { ...s, naik: !s.naik, halaman: 1 }
  return { ...s, urut: kolom, naik: kolom === 'icd', halaman: 1 }
}

/** Penanda arah di kepala kolom yang sedang dipakai mengurut; kolom lain tanpa tanda. */
export function tandaArah(s: Saringan, kolom: KolomUrut): string {
  if (s.urut !== kolom) return ''
  return s.naik ? ' ▲' : ' ▼'
}

/** Jumlah halaman (minimal 1). */
export function jumlahHalaman(total: number, ukuran = UKURAN_HALAMAN): number {
  return Math.max(1, Math.ceil(total / ukuran))
}

/** `SetUpperCase_DT` (ICD Code b1265, Disease b1538): perubahan medan menjadikannya huruf besar - dipanggil saat medan
 * ditinggalkan (peristiwa `change` Pega b1253 / b1526). Backend menerapkan hal yang sama. */
export function hurufBesar(s: string): string {
  return s.toUpperCase()
}

function bytes(s: string): number {
  return new TextEncoder().encode(s).length
}

/** Pemeriksaan awal form: ICD Code wajib (keputusan WO D3 - XML `pyRequired` false b1244) dan Disease wajib
 * (`pyRequired` b1465); panjang paling banyak lebar kolom (di luar XML). `null` = boleh dikirim; backend memeriksa
 * ulang. */
export function periksaIsian(icdCode: string, disease: string): string | null {
  const icd = hurufBesar(icdCode.trim())
  const nama = hurufBesar(disease.trim())
  if (icd === '') return DSL.galatICD
  if (nama === '') return DSL.galatDisease
  if (bytes(icd) > BATAS_ICD) return DSL.galatPanjang(DSL.icdCode, BATAS_ICD)
  if (bytes(nama) > BATAS_DISEASE) return DSL.galatPanjang(DSL.disease, BATAS_DISEASE)
  return null
}
