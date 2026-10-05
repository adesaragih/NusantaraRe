// Klien Business Group - `/api/business-group` (`modul/businessgroup/backend/handlers/rute.go`). Rutenya hanya terbuka
// bagi pemegang menu `businessgroup`; tulis ditolak bagi akses View only. Tanpa hapus, tanpa SYARIAH.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_BG = '/api/business-group'

/** Satu baris - `models.BisnisGrup`. */
export interface BisnisGrup {
  id: string
  /** `NOTE`. */
  name: string
  alias: string
  /** `TOPID` = `TREATYGROUP.ID`. */
  topId: string
  /** Salinan nama Treaty Group saat disimpan. */
  treatyName: string
}

/** Satu Treaty Group - `models.TreatyGroup`. */
export interface TreatyGroup {
  id: string
  name: string
}

/** Isian form Add / Edit - `models.Isian` (ID di jalur). */
export interface Isian {
  topId: string
  name: string
  alias: string
}

export function ambilDaftar(q: string, treatyGroup: string): Promise<{ daftar: BisnisGrup[] }> {
  return minta(PREFIX_BG, {
    kueri: { q: q.trim() === '' ? undefined : q.trim(), treatyGroup: treatyGroup === '' ? undefined : treatyGroup },
  })
}

export function ambilPilihan(): Promise<{ treatyGroup: TreatyGroup[] }> {
  return minta(`${PREFIX_BG}/pilihan`)
}

export function ambil(id: string): Promise<BisnisGrup> {
  return minta(`${PREFIX_BG}/${encodeURIComponent(id)}`)
}

/** Add - ID dibuat backend (situs aktif + BUSINESSGROUP_SEQ). */
export function tambah(isi: Isian): Promise<BisnisGrup> {
  return minta(PREFIX_BG, { metode: 'POST', badan: isi })
}

export function ubah(id: string, isi: Isian): Promise<BisnisGrup> {
  return minta(`${PREFIX_BG}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: isi })
}
