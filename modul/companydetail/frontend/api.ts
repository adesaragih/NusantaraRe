// Klien Company Detail - `/api/company-detail` (`modul/companydetail/backend/handlers/rute.go`). Rutenya hanya
// terbuka bagi pemegang menu `companydetail`; cookie sesi dikirim peramban sendiri.

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute - SAMA dengan `handlers.Prefix`. */
export const PREFIX_CD = '/api/company-detail'

/** Satu baris daftar - `models.BarisDaftar`. */
export interface BarisDaftar {
  id: string
  idView: string
  nama: string
  title: string
  npwp: string
  countryName: string
  /** Kode `BU_ID`. */
  businessField: string
  parentName: string
}

/** Satu halaman daftar - `services.Halaman`. */
export interface Halaman {
  daftar: BarisDaftar[]
  total: number
  halaman: number
  ukuran: number
}

/** Satu baris `CLIENT` organisasi - `models.Organisasi`. Semua teks; kosong = "". */
export interface Organisasi {
  /** `ASM-SFAGIS-WORK-ORG ORG-n`; tidak pernah berubah. */
  id: string
  /** `ORG-n`. */
  idView: string
  nama: string
  /** LABEL title (`PT.`). */
  title: string
  npwp: string
  /** `NATION.OLDID` (atau `NATION.ID` bila tanpa OLDID). */
  country: string
  countryName: string
  businessField: string
  parentId: string
  parentName: string
  note: string
  createdBy: string
  createdAt: string
  updatedBy: string
  updatedAt: string
}

/** Satu baris `CLIENT_PICLIST` - `models.PIC`. */
export interface PIC {
  /** Kosong = PIC baru (backend memberi `PIC-n`). */
  userIdentifier: string
  nama: string
  position: string
  /** `1` Male, `2` Female. */
  gender: string
  email: string
  /** `DD-MM-YYYY`. */
  dateOfBirth: string
  phone: string
}

/** Satu nomor Phone and Fax - satu baris `CLIENT_ADDRESS`. */
export interface Telfax {
  type: string
  code: string
  no: string
}

/** Satu alamat - `models.Alamat`. */
export interface Alamat {
  /** `ASMADDRESS` sebelum diubah; kosong = alamat baru. */
  asal: string
  type: string
  address: string
  telfax: Telfax[]
}

/** Satu organisasi lengkap - `models.Detail`. */
export interface Detail extends Organisasi {
  pic: PIC[]
  alamat: Alamat[]
}

/** Isian Create dan ubah - `models.Isian`. */
export interface Isian {
  nama: string
  title: string
  npwp: string
  country: string
  businessField: string
  parentId: string
  note: string
  pic: PIC[]
  alamat: Alamat[]
}

/** Satu baris `M_ENUMERASI` - `models.Pilihan`. */
export interface Pilihan {
  kode: string
  label: string
  aktif: boolean
}

/** Satu baris `NATION` - `models.Negara`. */
export interface Negara {
  id: string
  oldId: string
  nama: string
  nationInitial: string
}

/** Akun login aktif `M_LOGIN_GO` - `models.Akun`; pilihan PIC Name, `jabatan` (JOB_POSITION) = PIC Position. */
export interface Akun {
  loginId: string
  nama: string
  jabatan: string
}

/** Pilihan form - `services.PilihanForm`. */
export interface PilihanForm {
  title: Pilihan[]
  businessField: Pilihan[]
  position: Pilihan[]
  addressType: Pilihan[]
  telfax: Pilihan[]
  kodeArea: Pilihan[]
  gender: Pilihan[]
  negara: Negara[]
  /** Akun login aktif - pilihan PIC Name (perintah work owner 05-10-2026). */
  akun: Akun[]
}

/** Satu organisasi bernama sama/mirip - `services.NamaSerupa`. */
export interface NamaSerupa extends BarisDaftar {
  /** `sama` (sesudah title dan tanda baca dibuang) atau `mirip`. */
  jenis: 'sama' | 'mirip'
  skor: number
}

/** Jawaban pemeriksaan nama - `services.PeriksaNama`. */
export interface PeriksaNama {
  /** LABEL title yang ada di nama (mis. `PT.`); kosong = tidak ada. */
  titleDalamNama: string
  serupa: NamaSerupa[]
}

/** Tombol yang boleh tampil bagi akun ini - `services.Hak`. */
export interface Hak {
  /** Copy Old: hanya superadmin (pemegang menu Kelola User). */
  copyOld: boolean
}

/** Satu baris popup Copy Old - `models.OrgLama`: yang AKAN ditulis untuk organisasi dokumen ini. */
export interface OrgLama {
  id: string
  idView: string
  nama: string
  /** Belum punya baris CLIENT. */
  baru: boolean
  /** PIC yang akan ditambah. */
  pic: number
  /** Nomor Phone and Fax yang akan dipindah. */
  nomor: number
  /** Kolom CLIENT kosong yang akan diisi: PARENT_ID, NOTE, TITLE. */
  isi: string[]
}

export type StatusSalinLama = 'disalin' | 'sudahAda' | 'ditolak' | 'gagal'

export interface HasilSalinLama {
  id: string
  status: StatusSalinLama
  pesan: string[]
}

export interface JawabanSalinLama {
  hasil: HasilSalinLama[]
  disalin: number
}

const e = encodeURIComponent

/**
 * Batas waktu Copy Old: isi popup membaca dokumen Pega (`M_CLIENT`) - 17-23 detik di DEV 04-10-2026, terlalu dekat
 * ke batas bawaan 30 detik ("signal is aborted without reason").
 */
export const BATAS_WAKTU_LAMA_MS = 120_000

export function ambilDaftar(q: string, halaman: number, ukuran?: number): Promise<Halaman> {
  return minta<Halaman>(PREFIX_CD, { kueri: { q: q === '' ? undefined : q, halaman, ukuran } })
}

export function ambilPilihan(): Promise<PilihanForm> {
  return minta<PilihanForm>(`${PREFIX_CD}/pilihan`)
}

/** Pilihan Parent organization; `kecuali` = organisasi yang sedang diubah. */
export function cariInduk(q: string, kecuali: string): Promise<{ daftar: BarisDaftar[] }> {
  return minta<{ daftar: BarisDaftar[] }>(`${PREFIX_CD}/induk`, {
    kueri: { q, kecuali: kecuali === '' ? undefined : kecuali },
  })
}

/** Pemeriksaan nama sebelum Create: title di nama dan organisasi bernama sama/mirip. */
export function periksaNama(nama: string, kecuali: string): Promise<PeriksaNama> {
  return minta<PeriksaNama>(`${PREFIX_CD}/periksa-nama`, { kueri: { nama, kecuali: kecuali === '' ? undefined : kecuali } })
}

export function ambilHak(): Promise<Hak> {
  return minta<Hak>(`${PREFIX_CD}/hak`)
}

/** Isi popup Copy Old - superadmin. */
export function ambilLama(): Promise<{ daftar: OrgLama[] }> {
  return minta<{ daftar: OrgLama[] }>(`${PREFIX_CD}/lama`, { batasWaktuMs: BATAS_WAKTU_LAMA_MS })
}

/** Process Copy - superadmin; hasil per ID. */
export function salinLama(ids: string[]): Promise<JawabanSalinLama> {
  return minta<JawabanSalinLama>(`${PREFIX_CD}/lama/salin`, {
    metode: 'POST',
    badan: { ids },
    batasWaktuMs: BATAS_WAKTU_LAMA_MS,
  })
}

export function ambil(id: string): Promise<Detail> {
  return minta<Detail>(`${PREFIX_CD}/${e(id)}`)
}

/** Create - nomor ORG dari backend. */
export function tambah(isi: Isian): Promise<Detail> {
  return minta<Detail>(PREFIX_CD, { metode: 'POST', badan: isi })
}

/** Ubah - seluruh PIC dan alamat ditulis ulang; baris yang dibuang dari grid dihapus. */
export function ubah(id: string, isi: Isian): Promise<Detail> {
  return minta<Detail>(`${PREFIX_CD}/${e(id)}`, { metode: 'PUT', badan: isi })
}
