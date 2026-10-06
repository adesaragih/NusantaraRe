// Klien Treaty Group OJK - `/api/treaty-group-ojk` (`modul/treatygroupojk/backend/handlers/rute.go`). Rutenya hanya
// terbuka bagi pemegang menu `treatygroupojk`; tulis ditolak bagi akses View only. Tanpa hapus.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_TGO = '/api/treaty-group-ojk'

/** Satu baris - `models.Ojk`. */
export interface Ojk {
  id: string
  name: string
  nameIdn: string
  orderNo: string
}

/**
 * Isian form Add / Edit - `models.Isian` (ID di jalur). Order No bukan isian (perintah work owner 05-10-2026: "ORDERNO
 * hide aja, isi sesuai max dari order no") - backend mengisinya saat Add.
 */
export interface Isian {
  name: string
  nameIdn: string
}

export function ambilDaftar(q: string): Promise<{ daftar: Ojk[] }> {
  return minta(PREFIX_TGO, { kueri: { q: q.trim() === '' ? undefined : q.trim() } })
}

export function ambil(id: string): Promise<Ojk> {
  return minta(`${PREFIX_TGO}/${encodeURIComponent(id)}`)
}

/** Add - ID dan Order No dibuat backend (masing-masing nomor tertinggi + 1). */
export function tambah(isi: Isian): Promise<Ojk> {
  return minta(PREFIX_TGO, { metode: 'POST', badan: isi })
}

export function ubah(id: string, isi: Isian): Promise<Ojk> {
  return minta(`${PREFIX_TGO}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: isi })
}
