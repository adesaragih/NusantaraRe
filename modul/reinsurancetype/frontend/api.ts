// Klien Reinsurance Type - `/api/reinsurance-type` (`modul/reinsurancetype/backend/handlers/rute.go`). Rutenya hanya
// terbuka bagi pemegang menu `reinsurancetype`; tulis ditolak bagi akses View only. Tanpa hapus.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_RT = '/api/reinsurance-type'

/** Satu baris - `models.Jenis`. */
export interface Jenis {
  id: string
  /** `NOTE`. */
  name: string
  type: string
  soaName: string
  code: string
  flag: string
  noUrut: string
  groupType: string
  userId: string
  /** Teks format Pega `YYYYMMDDTHHMMSS.mmm GMT`. */
  tglUpdate: string
  /** Tampilan WIB `DD-MM-YYYY HH:MM`. */
  diubah: string
}

/** Isian form Add / Edit - `models.Isian` (ID di jalur). */
export interface Isian {
  name: string
  type: string
  soaName: string
  code: string
  flag: string
  noUrut: string
  groupType: string
}

export function ambilDaftar(q: string, type: string, flag: string): Promise<{ daftar: Jenis[] }> {
  return minta(PREFIX_RT, {
    kueri: {
      q: q.trim() === '' ? undefined : q.trim(),
      type: type === '' ? undefined : type,
      flag: flag === '' ? undefined : flag,
    },
  })
}

export function ambil(id: string): Promise<Jenis> {
  return minta(`${PREFIX_RT}/${encodeURIComponent(id)}`)
}

/** Add - ID dibuat backend (1 + M_REINSURANCETYPE_SEQ). */
export function tambah(isi: Isian): Promise<Jenis> {
  return minta(PREFIX_RT, { metode: 'POST', badan: isi })
}

export function ubah(id: string, isi: Isian): Promise<Jenis> {
  return minta(`${PREFIX_RT}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: isi })
}
