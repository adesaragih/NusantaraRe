// Klien Benefit - `/api/benefit-life` (`modul/benefitlife/backend/handlers/rute.go`). Rutenya hanya terbuka bagi
// pemegang menu `benefitlife`; tulis ditolak bagi akses View only. Tidak ada hapus (XML tanpa Delete).

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_BN = '/api/benefit-life'

/** Satu baris BENEFIT_LIFE - `models.Benefit` (grid: ID = `.Number` b3858, Benefit = `.Benefit` b4020). */
export interface Benefit {
  id: string
  benefit: string
}

/** Satu halaman grid - `models.Halaman`. */
export interface Halaman {
  daftar: Benefit[]
  total: number
  halaman: number
  ukuran: number
}

export interface Saringan {
  id: string
  benefit: string
  /** ID menaik; bawaan false = ID menurun (`pySortType` DESC b4391). */
  naik: boolean
  halaman: number
}

export function ambilDaftar(s: Saringan): Promise<Halaman> {
  return minta(PREFIX_BN, {
    kueri: {
      id: s.id.trim() === '' ? undefined : s.id.trim(),
      benefit: s.benefit.trim() === '' ? undefined : s.benefit.trim(),
      arah: s.naik ? 'asc' : undefined,
      halaman: s.halaman > 1 ? s.halaman : undefined,
    },
  })
}

/** Save - Add (`AddToList_Act` b1791; ID = '1' || LPAD(M_BENEFIT_LIFE_SEQ, 5), dibentuk server). */
export function tambah(benefit: string): Promise<Benefit> {
  return minta(PREFIX_BN, { metode: 'POST', badan: { benefit } })
}

/** Save sesudah Edit (`EditList_DT` b4243). */
export function ubah(id: string, benefit: string): Promise<Benefit> {
  return minta(`${PREFIX_BN}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: { benefit } })
}
