// Klien Accounts - `/api/accounts` (`modul/accounts/backend/handlers/rute.go`). Rutenya hanya terbuka bagi pemegang
// menu `accounts`; cookie sesi dikirim peramban sendiri. Akun hanya diinput sekali: nol ubah, nol hapus.

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute - SAMA dengan `handlers.Prefix`. */
export const PREFIX_ACC = '/api/accounts'

/** Satu baris `POOLDATA.T_M_ACCOUNT` - `models.Account`. Semua teks; kosong = "". */
export interface Account {
  /** `ASM-SFAGIS-WORK-ACCOUNT ACC-n`. */
  id: string
  /** `ACC-n` - kosong bila ID tidak berpola itu (baris lama Pega). */
  idView: string
  groupBusinessId: string
  groupBusiness: string
  /** `CLIENT.ID` organisasi. */
  insuredId: string
  /** `ORG-n`. */
  orgId: string
  insuredName: string
  description: string
  /** "Owner" - akun pelaku saat dibuat; kosong untuk baris lama. */
  createOp: string
  /** `YYYY-MM-DD HH:MI` jam Oracle; kosong untuk baris lama. */
  createDate: string
}

/** Satu halaman daftar - `services.HalamanDaftar`. */
export interface HalamanDaftar {
  daftar: Account[]
  total: number
  halaman: number
  ukuran: number
}

/** Pilihan "Insured Name" - `models.Organisasi`. */
export interface Organisasi {
  id: string
  idView: string
  nama: string
}

/** Pilihan "Group Business" - `models.GroupBusiness`. */
export interface GroupBusiness {
  id: string
  note: string
}

/** Isian form tambah - `models.Isian`. */
export interface Isian {
  insuredId: string
  groupBusinessId: string
  description: string
}

export function ambilDaftar(cari: string, halaman: number, ukuran: number): Promise<HalamanDaftar> {
  return minta<HalamanDaftar>(PREFIX_ACC, { kueri: { cari: cari === '' ? undefined : cari, halaman, ukuran } })
}

export function ambilPilihan(): Promise<{ groupBusiness: GroupBusiness[] }> {
  return minta<{ groupBusiness: GroupBusiness[] }>(`${PREFIX_ACC}/pilihan`)
}

export function cariOrganisasi(cari: string): Promise<{ daftar: Organisasi[]; lebih: boolean }> {
  return minta<{ daftar: Organisasi[]; lebih: boolean }>(`${PREFIX_ACC}/organisasi`, {
    kueri: { cari: cari === '' ? undefined : cari },
  })
}

export function tambah(isi: Isian): Promise<Account> {
  return minta<Account>(PREFIX_ACC, { metode: 'POST', badan: isi })
}
