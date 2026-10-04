// Aturan murni form Kelola User — tanpa React, diuji `kelolauser.test.ts`.
//
// Pemeriksaan di sini hanya AWAL, supaya admin tidak menunggu server untuk
// salah yang jelas; backend tetap menegakkan semuanya (jenjang organisasi,
// master aktif, KODE menu, penjaga diri sendiri dan admin terakhir).

import { KELOLA_USER, PANJANG_MIN_SANDI } from '../labels'
import { KODE_MENU_KELOLA_USER } from '../lib/daftarMenu'
import type { BadanBaru, BadanSandi, BadanUbah, OpsiMaster, OpsiMenu, PilihanKelola, RingkasAkun, RinciAkun } from './api'

/** Bentuk username — sama dengan `login.polaAkun` di backend. */
export const POLA_AKUN = /^[A-Za-z0-9._@-]{1,64}$/

/** Username sah — `login.AkunSah`: pola di atas, dan bukan titik saja. */
export function akunSah(akun: string): boolean {
  return POLA_AKUN.test(akun) && akun.replace(/\./g, '') !== ''
}

/** Isian form — satu bentuk untuk buat dan ubah. */
export interface IsianForm {
  akunId: string
  nama: string
  sandi: string
  ulangiSandi: string
  organisasi: string
  divisi: string
  unit: string
  workbasket: string[]
  menu: string[]
  /** Centang "Change Password Next Login" (tab Security). */
  wajibGanti: boolean
  /** Kontak akun (Kelola User 03-10-2026) — opsional. */
  email: string
  telepon: string
  nik: string
  jabatan: string
}

/** Tab form tambah/ubah user. */
export type TabForm = 'profil' | 'security'

export function isianKosong(): IsianForm {
  return {
    akunId: '',
    nama: '',
    sandi: '',
    ulangiSandi: '',
    organisasi: '',
    divisi: '',
    unit: '',
    workbasket: [],
    menu: [],
    // User baru bawaannya WAJIB mengganti password yang diketik admin.
    wajibGanti: true,
    email: '',
    telepon: '',
    nik: '',
    jabatan: '',
  }
}

export function isianDari(r: RinciAkun): IsianForm {
  return {
    akunId: r.akunId,
    nama: r.nama,
    sandi: '',
    ulangiSandi: '',
    organisasi: r.organisasi,
    divisi: r.divisi,
    unit: r.unit,
    workbasket: [...r.workbasket],
    menu: [...r.menu],
    wajibGanti: r.wajibGantiSandi,
    email: r.email,
    telepon: r.telepon,
    nik: r.nik,
    jabatan: r.jabatan,
  }
}

/** Bentuk kontak — sama dengan `login.PeriksaKontak` di backend (kosong selalu sah). */
export const POLA_EMAIL = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
export const POLA_TELEPON = /^\+?[0-9][0-9 -]*[0-9]$/
export const POLA_NIK = /^[A-Za-z0-9./-]+$/

/** Panjang BYTE UTF-8, seperti kolom VARCHAR2 dan `len()` Go. */
function panjangByte(s: string): number {
  return new TextEncoder().encode(s).length
}

/** Galat kontak PERTAMA (urut medan di form); `null` = sah. */
export function periksaKontak(isi: Pick<IsianForm, 'email' | 'telepon' | 'nik' | 'jabatan'>): string | null {
  const email = isi.email.trim()
  if (email !== '' && (panjangByte(email) > 254 || !POLA_EMAIL.test(email))) return KELOLA_USER.galatEmail
  const telepon = isi.telepon.trim()
  if (telepon !== '') {
    const digit = telepon.replace(/[^0-9]/g, '').length
    if (telepon.length > 30 || !POLA_TELEPON.test(telepon) || digit < 8 || digit > 15) return KELOLA_USER.galatTelepon
  }
  const nik = isi.nik.trim()
  if (nik !== '' && (nik.length > 30 || !POLA_NIK.test(nik))) return KELOLA_USER.galatNIK
  if (panjangByte(isi.jabatan.trim()) > 150) return KELOLA_USER.galatJabatan
  return null
}

/** Divisi milik organisasi itu. */
export function divisiUntuk(p: PilihanKelola, organisasi: string): OpsiMaster[] {
  return organisasi === '' ? [] : p.divisi.filter((d) => d.induk === organisasi)
}

/** Unit milik divisi itu. */
export function unitUntuk(p: PilihanKelola, divisi: string): OpsiMaster[] {
  return divisi === '' ? [] : p.unit.filter((u) => u.induk === divisi)
}

/** Mengganti organisasi; divisi dan unit yang bukan miliknya dikosongkan. */
export function setelOrganisasi(isi: IsianForm, organisasi: string, p: PilihanKelola): IsianForm {
  const divisiTetap = divisiUntuk(p, organisasi).some((d) => d.kode === isi.divisi)
  return setelDivisi({ ...isi, organisasi }, divisiTetap ? isi.divisi : '', p)
}

/** Mengganti divisi; unit yang bukan miliknya dikosongkan. */
export function setelDivisi(isi: IsianForm, divisi: string, p: PilihanKelola): IsianForm {
  const unitTetap = unitUntuk(p, divisi).some((u) => u.kode === isi.unit)
  return { ...isi, divisi, unit: unitTetap ? isi.unit : '' }
}

/** Menyalakan atau mematikan satu kode; hasilnya urut dan tanpa ganda. */
export function alihkan(daftar: readonly string[], kode: string, nyala: boolean): string[] {
  const s = new Set(daftar)
  if (nyala) s.add(kode)
  else s.delete(kode)
  return [...s].sort()
}

/**
 * Pemeriksaan awal; `null` = boleh dikirim. Panjang sandi dihitung KARAKTER,
 * seperti backend. `akunSaya` = akun yang sedang login. Password WAJIB saat
 * membuat user; saat mengubah, kosong = tidak diganti.
 */
export function periksaIsian(isi: IsianForm, baru: boolean, akunSaya: string): string | null {
  if (baru && !akunSah(isi.akunId.trim())) return KELOLA_USER.galatAkun
  const nama = isi.nama.trim()
  if (nama === '' || [...nama].length > 150) return KELOLA_USER.galatNama
  const kontak = periksaKontak(isi)
  if (kontak !== null) return kontak
  if (baru || isi.sandi !== '' || isi.ulangiSandi !== '') {
    if ([...isi.sandi].length < PANJANG_MIN_SANDI) return KELOLA_USER.galatSandi
    if (isi.sandi !== isi.ulangiSandi) return KELOLA_USER.galatUlangi
  }
  if (!baru && isi.akunId === akunSaya && !isi.menu.includes(KODE_MENU_KELOLA_USER)) return KELOLA_USER.menuDiriSendiri
  return null
}

/**
 * Username atau email sudah terdaftar (migrasi 905) — dicek di layar dari daftar yang sudah termuat supaya pesannya
 * langsung muncul; backend tetap penjaga sebenarnya (409). Aturannya sama dengan `login.PemakaiUsername` dan
 * `login.PemakaiEmail`: username dengan username, email dengan email, masing-masing tanpa beda huruf. `akunId` = akun
 * yang diubah (`null` = user baru); `daftar` `null` = belum termuat, tidak menebak.
 */
export function periksaGanda(
  isi: Pick<IsianForm, 'akunId' | 'email'>,
  akunId: string | null,
  daftar: readonly RingkasAkun[] | null,
): string | null {
  if (daftar === null) return null
  const lain = daftar.filter((a) => a.akunId !== akunId)
  const sama = (a: string, b: string) => a !== '' && a.toLowerCase() === b.toLowerCase()
  if (akunId === null) {
    const id = isi.akunId.trim()
    if (lain.some((a) => sama(a.akunId, id))) return KELOLA_USER.galatUsernameTerdaftar
  }
  const email = isi.email.trim()
  if (email !== '' && lain.some((a) => sama(a.email, email))) return KELOLA_USER.galatEmailTerdaftar
  return null
}

export function badanUbah(isi: IsianForm): BadanUbah {
  return {
    nama: isi.nama.trim(),
    organisasi: isi.organisasi,
    divisi: isi.divisi,
    unit: isi.unit,
    workbasket: [...isi.workbasket],
    menu: [...isi.menu],
    email: isi.email.trim(),
    telepon: isi.telepon.trim(),
    nik: isi.nik.trim(),
    jabatan: isi.jabatan.trim(),
  }
}

export function badanBaru(isi: IsianForm): BadanBaru {
  return { ...badanUbah(isi), akunId: isi.akunId.trim(), sandi: isi.sandi, wajibGanti: isi.wajibGanti }
}

/**
 * Badan tab Security saat MENGUBAH user; `null` = tidak ada yang berubah
 * (password kosong dan centangnya sama dengan `wajibAwal`).
 */
export function badanSandi(isi: IsianForm, wajibAwal: boolean): BadanSandi | null {
  if (isi.sandi === '' && isi.wajibGanti === wajibAwal) return null
  return { sandi: isi.sandi, wajibGanti: isi.wajibGanti }
}

/** Tab tempat galat itu harus diperbaiki. */
export function tabGalat(pesan: string): TabForm {
  return pesan === KELOLA_USER.galatSandi || pesan === KELOLA_USER.galatUlangi ? 'security' : 'profil'
}

/** Kotak centang menu per golongan — urutan server (urutan sidebar). */
export function kelompokMenu(menu: readonly OpsiMenu[]): { golongan: string; menu: OpsiMenu[] }[] {
  const hasil: { golongan: string; menu: OpsiMenu[] }[] = []
  for (const m of menu) {
    const g = hasil.find((x) => x.golongan === m.golongan)
    if (g === undefined) hasil.push({ golongan: m.golongan, menu: [m] })
    else g.menu.push(m)
  }
  return hasil
}

/** Saringan daftar: setiap kata harus ada di username atau nama. */
export function saringDaftar(daftar: readonly RingkasAkun[], kueri: string): RingkasAkun[] {
  const kata = kueri.trim().toLowerCase().split(/\s+/).filter(Boolean)
  if (kata.length === 0) return [...daftar]
  return daftar.filter((a) => {
    const jerami = [a.akunId, a.contactId, a.nama, a.email].join(' ').toLowerCase()
    return kata.every((k) => jerami.includes(k))
  })
}

/** Teks satu kolom Organisasi / Divisi / Unit. */
export function teksJenjang(a: Pick<RingkasAkun, 'organisasi' | 'divisi' | 'unit'>): string {
  return [a.organisasi, a.divisi, a.unit].filter((x) => x !== '').join(' / ') || '—'
}
