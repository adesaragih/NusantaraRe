// Klien Treaty Description - `/api/treaty-description` (`modul/treatydescription/backend/handlers/rute.go`). Rutenya
// hanya terbuka bagi pemegang menu `treatydescription`; tulis ditolak bagi akses View only. Tanpa hapus.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_TD = '/api/treaty-description'

/** `ISXOL`: "0" Non XOL, "1" XOL. */
export type Jenis = '0' | '1'
/** `STATUSAKTIF` yang ditulis: "1" Active, "0" Inactive. */
export type Status = '1' | '0'

/** Satu baris - `models.Desc`. */
export interface Desc {
  id: string
  descName: string
  isXol: string
  /** Teks apa adanya: "1", "0", atau "" (NULL - baris lama Pega, dibaca Active). */
  statusAktif: string
  aktif: boolean
}

/** Isian form Add / Edit - `models.Isian` (ID di jalur). */
export interface Isian {
  descName: string
  isXol: Jenis
  statusAktif: Status
}

/** Saringan daftar; "" = semua. */
export interface Saringan {
  q: string
  xol: '' | Jenis
  status: '' | Status
}

const kosongJadiHilang = (s: string) => (s.trim() === '' ? undefined : s.trim())

export function ambilDaftar(s: Saringan): Promise<{ daftar: Desc[] }> {
  return minta(PREFIX_TD, {
    kueri: { q: kosongJadiHilang(s.q), xol: kosongJadiHilang(s.xol), status: kosongJadiHilang(s.status) },
  })
}

export function ambil(id: string): Promise<Desc> {
  return minta(`${PREFIX_TD}/${encodeURIComponent(id)}`)
}

/** Add - ID dibuat backend ('1' + TREATY_DESCRIPTION_SEQ 4 digit). */
export function tambah(isi: Isian): Promise<Desc> {
  return minta(PREFIX_TD, { metode: 'POST', badan: isi })
}

export function ubah(id: string, isi: Isian): Promise<Desc> {
  return minta(`${PREFIX_TD}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: isi })
}
