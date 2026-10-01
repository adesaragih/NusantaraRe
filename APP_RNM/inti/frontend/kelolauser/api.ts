// Klien Kelola User — `/api/admin/*` (`inti/backend/login/kelola_rute.go`,
// keputusan work owner 01-10-2026).
//
// ⛔ Rute ini menuntut SESI LOGIN pemegang menu Kelola User; cookie HttpOnly
// dikirim peramban sendiri. Sandi awal hanya hidup di badan `buatPengguna`
// dan tidak pernah kembali di jawaban mana pun. Reset sandi DITUNDA.

import { minta } from '../klien'

/** Satu baris daftar — `login.RingkasAkun`. */
export interface RingkasAkun {
  akunId: string
  nama: string
  organisasi: string
  divisi: string
  unit: string
  aktif: boolean
  terkunci: boolean
  wajibGantiSandi: boolean
  /** `YYYY-MM-DD HH:MI` jam server; kosong = belum pernah login. */
  loginTerakhir: string
}

/** Satu akun beserta workbasket dan menunya — `login.RinciAkun`. */
export interface RinciAkun extends RingkasAkun {
  workbasket: string[]
  menu: string[]
}

/** Satu pilihan dropdown master; `induk` = CODE organisasi/divisi pemiliknya. */
export interface OpsiMaster {
  kode: string
  nama: string
  induk?: string
}

/** Satu kotak centang menu. */
export interface OpsiMenu {
  kode: string
  label: string
  /** GROUPMENU (TREATY, …) atau ADMIN untuk menu aplikasi. */
  golongan: string
  dimigrasi: boolean
}

/** Isi pilihan form — `login.PilihanKelola`. Master hanya yang aktif. */
export interface PilihanKelola {
  organisasi: OpsiMaster[]
  divisi: OpsiMaster[]
  unit: OpsiMaster[]
  workbasket: OpsiMaster[]
  menu: OpsiMenu[]
}

/** Badan ubah. `akunId` tidak dapat diubah. */
export interface BadanUbah {
  nama: string
  organisasi: string
  divisi: string
  unit: string
  workbasket: string[]
  menu: string[]
}

/** Badan buat: ubah + username + sandi awal yang diketik admin. */
export interface BadanBaru extends BadanUbah {
  akunId: string
  sandi: string
}

const DASAR = '/api/admin/pengguna'
const id = (akunId: string) => `${DASAR}/${encodeURIComponent(akunId)}`

export function ambilDaftarPengguna(): Promise<RingkasAkun[]> {
  return minta<RingkasAkun[]>(DASAR)
}

export function ambilPengguna(akunId: string): Promise<RinciAkun> {
  return minta<RinciAkun>(id(akunId))
}

export function ambilPilihanPengguna(): Promise<PilihanKelola> {
  return minta<PilihanKelola>('/api/admin/pilihan-pengguna')
}

export function buatPengguna(badan: BadanBaru): Promise<RinciAkun> {
  return minta<RinciAkun>(DASAR, { metode: 'POST', badan })
}

export function ubahPengguna(akunId: string, badan: BadanUbah): Promise<RinciAkun> {
  return minta<RinciAkun>(id(akunId), { metode: 'PUT', badan })
}

export function setelAktifPengguna(akunId: string, aktif: boolean): Promise<RinciAkun> {
  return minta<RinciAkun>(`${id(akunId)}/aktif`, { metode: 'POST', badan: { aktif } })
}

export function bukaKunciPengguna(akunId: string): Promise<RinciAkun> {
  return minta<RinciAkun>(`${id(akunId)}/buka-kunci`, { metode: 'POST' })
}

/** HAPUS PERMANEN — akun, workbasket, dan menunya. */
export async function hapusPengguna(akunId: string): Promise<void> {
  await minta<undefined>(id(akunId), { metode: 'DELETE' })
}
