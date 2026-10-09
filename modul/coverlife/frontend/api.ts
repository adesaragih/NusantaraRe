// Klien Cover Life - `/api/cover-life` (`modul/coverlife/backend/handlers/rute.go`). Rutenya hanya terbuka bagi pemegang
// menu `coverlife`; tulis ditolak bagi akses View only. Tidak ada hapus (XML tanpa Delete), saring, maupun urut pilihan.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_CVL = '/api/cover-life'

/** Satu baris M_COVER_LIFE - `models.Cover` (grid: ID = `.ID` b3449, Cover = `.Cover` b3596; Note tidak tampil di grid
 * tetapi diisi Edit - `EditList_DT` b3823). */
export interface Cover {
  id: string
  cover: string
  note: string
}

/** Satu halaman grid - `models.Halaman`. */
export interface Halaman {
  daftar: Cover[]
  total: number
  halaman: number
  ukuran: number
}

/** Grid ID menaik tetap (b3967) - satu-satunya masukan adalah halaman. */
export function ambilDaftar(halaman: number): Promise<Halaman> {
  return minta(PREFIX_CVL, { kueri: { halaman: halaman > 1 ? halaman : undefined } })
}

/** Save - Add (`AddToList_Act` b5426; ID = '1' || LPAD(M_COVER_LIFE_SEQ, 5), dibentuk server). */
export function tambah(cover: string, note: string): Promise<Cover> {
  return minta(PREFIX_CVL, { metode: 'POST', badan: { cover, note } })
}

/** Save sesudah Edit (`EditList_DT` b3807). */
export function ubah(id: string, cover: string, note: string): Promise<Cover> {
  return minta(`${PREFIX_CVL}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: { cover, note } })
}
