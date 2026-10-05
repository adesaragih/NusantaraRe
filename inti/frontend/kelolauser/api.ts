// Klien Kelola User — `/api/admin/*` (`inti/backend/login/kelola_rute.go`,
// keputusan work owner 01-10-2026).
//
// ⛔ Rute ini menuntut SESI LOGIN pemegang menu Kelola User; cookie HttpOnly
// dikirim peramban sendiri. Password hanya hidup di badan `buatPengguna` dan
// `aturSandiPengguna` (tab Security) dan tidak pernah kembali di jawaban.

import { minta } from '../klien'

/** Satu baris daftar — `login.RingkasAkun`. */
export interface RingkasAkun {
  akunId: string
  /** `CONTACT_ID` `CON-n` (migrasi 905) — diberi saat user dibuat, tidak pernah berubah. */
  contactId: string
  nama: string
  organisasi: string
  divisi: string
  unit: string
  aktif: boolean
  terkunci: boolean
  wajibGantiSandi: boolean
  /** `YYYY-MM-DD HH:MI` jam server; kosong = belum pernah login. */
  loginTerakhir: string
  /** Kontak akun (migrasi 904, Kelola User 03-10-2026) — semuanya opsional; kosong = tidak diisi. */
  email: string
  telepon: string
  /** NIK = Nomor Induk Karyawan (Employee ID). */
  nik: string
  jabatan: string
}

/** Satu akun beserta workbasket dan menunya — `login.RinciAkun`. */
export interface RinciAkun extends RingkasAkun {
  workbasket: string[]
  menu: string[]
  /** Bagian dari `menu` yang View only (`M_LOGIN_GO_MENU.HAK = 'LIHAT'`, migrasi 914). Tidak ada = backend lama. */
  menuLihat?: string[]
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
  /** Modul yang mendukung akses View only (04-10-2026): pilihan Full / View only tampil di sampingnya. */
  bisaLihat?: boolean
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
  /** Bagian dari `menu` yang View only. */
  menuLihat: string[]
  email: string
  telepon: string
  nik: string
  jabatan: string
}

/** Badan buat: ubah + username + sandi awal yang diketik admin. */
export interface BadanBaru extends BadanUbah {
  akunId: string
  sandi: string
  /** Centang "Change Password Next Login". */
  wajibGanti: boolean
}

/** Badan tab Security: `sandi` kosong = hanya centang yang berubah. */
export interface BadanSandi {
  sandi: string
  wajibGanti: boolean
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

/**
 * Tab Security — password baru (sesi akun itu berakhir, kuncinya dibuka) dan/atau
 * centang "Change Password Next Login". Password akun SENDIRI: backend
 * menerbitkan ulang cookie, sesi ini tetap berjalan.
 */
export function aturSandiPengguna(akunId: string, badan: BadanSandi): Promise<RinciAkun> {
  return minta<RinciAkun>(`${id(akunId)}/sandi`, { metode: 'POST', badan })
}

/** HAPUS PERMANEN — akun, workbasket, dan menunya. */
export async function hapusPengguna(akunId: string): Promise<void> {
  await minta<undefined>(id(akunId), { metode: 'DELETE' })
}
