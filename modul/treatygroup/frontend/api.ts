// Klien Treaty Group - `/api/treaty-group` (`modul/treatygroup/backend/handlers/rute.go`). Rutenya hanya terbuka bagi
// pemegang menu `treatygroup`; tulis ditolak bagi akses View only. Tanpa hapus.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_TG = '/api/treaty-group'

/** Satu baris - `models.Grup`. */
export interface Grup {
  id: string
  oldId: string
  ojkId: string
  ojkName: string
  ojkNameIdn: string
  orderNo: string
  name: string
  soaName: string
  /** Teks format Pega `YYYYMMDDTHHMMSS.mmm GMT`. */
  tglUpdate: string
  /** Tampilan WIB `DD-MM-YYYY HH:MM`. */
  diubah: string
  userId: string
  coaId: string
  coaName: string
}

/** Grup bisnis anak (`BUSINESSGROUP.TOPID`), tanpa SYARIAH. */
export interface BisnisGrup {
  id: string
  name: string
  alias: string
}

/** Satu grup dan anaknya - `models.Detail`. */
export interface Detail extends Grup {
  anak: BisnisGrup[]
}

/** Satu OJK Business - `models.Ojk` - dan COA-nya (yang ikut tersimpan bila OJK ini dipilih; kosong = belum ada). */
export interface Ojk {
  id: string
  name: string
  nameIdn: string
  orderNo: string
  coaId: string
  coaName: string
}

/** Isian form Add / Edit - `models.Isian` (ID di jalur). COAID bukan isian (opsi A). */
export interface Isian {
  ojkId: string
  name: string
  soaName: string
}

export function ambilDaftar(q: string, ojk: string): Promise<{ daftar: Grup[] }> {
  return minta(PREFIX_TG, { kueri: { q: q.trim() === '' ? undefined : q.trim(), ojk: ojk === '' ? undefined : ojk } })
}

export function ambilPilihan(): Promise<{ ojk: Ojk[] }> {
  return minta(`${PREFIX_TG}/pilihan`)
}

export function ambil(id: string): Promise<Detail> {
  return minta(`${PREFIX_TG}/${encodeURIComponent(id)}`)
}

/** Add - ID dibuat backend (situs aktif + TREATYGROUP_SEQ). */
export function tambah(isi: Isian): Promise<Detail> {
  return minta(PREFIX_TG, { metode: 'POST', badan: isi })
}

export function ubah(id: string, isi: Isian): Promise<Detail> {
  return minta(`${PREFIX_TG}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: isi })
}
