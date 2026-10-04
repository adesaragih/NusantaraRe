// Klien API modul Master Data (`/api/masterdata`, kontrak `docs/issues/02-api-master-data.md`).
//
// Seluruh nilai baris = teks (kolom tabel apa adanya); status aktif = boolean `aktif`.

import { minta } from '../../../inti/frontend/klien'

/** Satu kolom master (`GET /api/masterdata`). */
export interface KolomMaster {
  /** Kunci JSON baris. */
  kunci: string
  /** Nama kolom tabel. */
  kolom: string
  /** Lebar kolom (byte). */
  lebar: number
  wajib: boolean
  /** Diisi backend dari tabel rujukan - tidak dikirim. */
  turunan: boolean
}

/** Satu master: tabel, judul, dan kolomnya. */
export interface MetaMaster {
  kunci: string
  judul: string
  /** ID dibuat backend (Accumulation). */
  idOtomatis: boolean
  kolom: KolomMaster[]
}

/** Satu baris master: kunci kolom -> teks, plus status. */
export type BarisMaster = { aktif: boolean } & Record<string, string | boolean>

export interface HalamanMaster {
  baris: BarisMaster[]
  total: number
  halaman: number
  ukuran: number
}

/** Saringan status daftar. */
export type StatusSaring = '' | 'aktif' | 'nonaktif'

/** `GET /api/masterdata` - delapan master beserta kolomnya. */
export function daftarMaster(): Promise<{ tabel: MetaMaster[] }> {
  return minta<{ tabel: MetaMaster[] }>('/api/masterdata')
}

/** `GET /api/masterdata/{tabel}?q=&status=&halaman=`. */
export function cariMaster(tabel: string, q: string, status: StatusSaring, halaman: number): Promise<HalamanMaster> {
  const kueri: Record<string, string | number> = { halaman }
  if (q.trim() !== '') kueri.q = q.trim()
  if (status !== '') kueri.status = status
  return minta<HalamanMaster>(`/api/masterdata/${tabel}`, { kueri })
}

/** `POST /api/masterdata/{tabel}` -> ID baris baru. */
export function tambahMaster(tabel: string, isi: Record<string, string>): Promise<{ id: string }> {
  return minta<{ id: string }>(`/api/masterdata/${tabel}`, { metode: 'POST', badan: isi })
}

/** `PUT /api/masterdata/{tabel}/{id}`. */
export function ubahMaster(tabel: string, id: string, isi: Record<string, string>): Promise<{ id: string }> {
  return minta<{ id: string }>(`/api/masterdata/${tabel}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: isi })
}

/** `PUT /api/masterdata/{tabel}/{id}/status`. */
export function ubahStatusMaster(tabel: string, id: string, aktif: boolean): Promise<{ id: string; aktif: boolean }> {
  return minta<{ id: string; aktif: boolean }>(`/api/masterdata/${tabel}/${encodeURIComponent(id)}/status`, {
    metode: 'PUT',
    badan: { aktif },
  })
}
