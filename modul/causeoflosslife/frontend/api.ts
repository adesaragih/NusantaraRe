// Klien Cause Of Loss Life - `/api/cause-of-loss-life` (`modul/causeoflosslife/backend/handlers/rute.go`). Rutenya
// hanya terbuka bagi pemegang menu `causeoflosslife`; tulis ditolak bagi akses View only. Tidak ada hapus (XML tanpa
// Delete), saring, maupun urut pilihan.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_COL = '/api/cause-of-loss-life'

/** Satu baris CAUSEOFLOSS_LIFE - `models.CauseOfLoss` (grid: ID = `.ID` b3415, Cause of Loss = `.CauseofLoss` b3562). */
export interface CauseOfLoss {
  id: string
  causeOfLoss: string
}

/** Satu halaman grid - `models.Halaman`. */
export interface Halaman {
  daftar: CauseOfLoss[]
  total: number
  halaman: number
  ukuran: number
}

/** Grid ID menaik tetap (b3923) - satu-satunya masukan adalah halaman. */
export function ambilDaftar(halaman: number): Promise<Halaman> {
  return minta(PREFIX_COL, { kueri: { halaman: halaman > 1 ? halaman : undefined } })
}

/** Save - Add (`AddToList_Act` b5385; ID = '1' || LPAD(M_CAUSEOFLOSS_LIFE_SEQ, 5), dibentuk server). */
export function tambah(causeOfLoss: string): Promise<CauseOfLoss> {
  return minta(PREFIX_COL, { metode: 'POST', badan: { causeOfLoss } })
}

/** Save sesudah Edit (`EditList_DT` b3770). */
export function ubah(id: string, causeOfLoss: string): Promise<CauseOfLoss> {
  return minta(`${PREFIX_COL}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: { causeOfLoss } })
}
