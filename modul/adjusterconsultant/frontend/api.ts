// Klien Adjuster Consultant - `/api/adjuster-consultant` (`modul/adjusterconsultant/backend/handlers/rute.go`). Rutenya
// hanya terbuka bagi pemegang menu `adjusterconsultant`; tulis ditolak bagi akses View only.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_ADJ = '/api/adjuster-consultant'

/** Satu baris - `models.Adjuster`. */
export interface Adjuster {
  id: string
  name: string
  address: string
  telpNo: string
  /** Akun pengubah terakhir. */
  username: string
  /** `DD-MM-YYYY HH24:MI`. */
  editDate: string
  active: boolean
}

/** Isian form Add / Edit - `models.Isian` (ID di jalur). */
export interface Isian {
  name: string
  address: string
  telpNo: string
}

/** Saringan status daftar. */
export type Status = '' | 'active' | 'inactive'

export function ambilDaftar(q: string, status: Status): Promise<{ daftar: Adjuster[] }> {
  return minta(PREFIX_ADJ, { kueri: { q: q.trim() === '' ? undefined : q.trim(), status: status === '' ? undefined : status } })
}

export function ambil(id: string): Promise<Adjuster> {
  return minta(`${PREFIX_ADJ}/${encodeURIComponent(id)}`)
}

/** Add - ID dibuat backend (situs aktif + sequence). */
export function tambah(isi: Isian): Promise<Adjuster> {
  return minta(PREFIX_ADJ, { metode: 'POST', badan: isi })
}

export function ubah(id: string, isi: Isian): Promise<Adjuster> {
  return minta(`${PREFIX_ADJ}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: isi })
}

/** Activate / Deactivate - pengganti hapus. */
export function setelAktif(id: string, active: boolean): Promise<Adjuster> {
  return minta(`${PREFIX_ADJ}/${encodeURIComponent(id)}/aktif`, { metode: 'POST', badan: { active } })
}
