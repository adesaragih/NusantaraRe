// Klien Disease Life - `/api/disease-life` (`modul/diseaselife/backend/handlers/rute.go`). Rutenya hanya terbuka bagi
// pemegang menu `diseaselife`; tulis ditolak bagi akses View only. Tidak ada hapus (XML tanpa Delete). 97.586 baris
// DEV: saring dan halaman SELALU dikirim ke server - klien tidak pernah memuat seluruh tabel.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_DSL = '/api/disease-life'

/** Satu baris DISEASE_LIFE - `models.Penyakit` (grid: ID = `.Number` b4321, ICD Code = `.ICD_Code` b4484, Disease =
 * `.Disease` b4629). */
export interface Penyakit {
  id: string
  icdCode: string
  disease: string
}

/** Satu halaman grid - `models.Halaman`. */
export interface Halaman {
  daftar: Penyakit[]
  total: number
  halaman: number
  ukuran: number
}

/** Kolom urut pilihan: ID (`pyColumnSorting` true b5014) atau ICD Code (b5036). */
export type KolomUrut = 'id' | 'icd'

export interface Saringan {
  icdCode: string
  disease: string
  urut: KolomUrut
  /** Menaik; bawaan false = ID menurun (`pySortType` DESC b5012). */
  naik: boolean
  halaman: number
}

export function ambilDaftar(s: Saringan): Promise<Halaman> {
  return minta(PREFIX_DSL, {
    kueri: {
      icd: s.icdCode.trim() === '' ? undefined : s.icdCode.trim(),
      disease: s.disease.trim() === '' ? undefined : s.disease.trim(),
      urut: s.urut === 'icd' ? 'icd' : undefined,
      arah: s.naik ? 'asc' : undefined,
      halaman: s.halaman > 1 ? s.halaman : undefined,
    },
  })
}

/** Save - Add (`AddToList_Act` b2121; ID = TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL), dibentuk server). */
export function tambah(icdCode: string, disease: string): Promise<Penyakit> {
  return minta(PREFIX_DSL, { metode: 'POST', badan: { icdCode, disease } })
}

/** Save sesudah Edit (`EditList_DT` b4852). */
export function ubah(id: string, icdCode: string, disease: string): Promise<Penyakit> {
  return minta(`${PREFIX_DSL}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: { icdCode, disease } })
}
