// modul/nbtreatyin/frontend/api.ts - panggilan backend modul NB Treaty In, satu
// fungsi per rute (`backend/handlers/rute.go`). Klien HTTP-nya
// `inti/frontend/klien.ts`.
//
// ⛔ Nol perhitungan uang di sini. Setiap rumus Pega dijalankan backend
// (`POST .../hitung`, models/hitung.go) - layar hanya mengirim halaman dan
// menampilkan hasilnya (AC 25, 79).

import { minta } from '../../../inti/frontend/klien'

/** Prefix rute API modul ini - SAMA dengan `handlers.Prefix`. */
export const PREFIX_NBTREATYIN = '/api/nb-treaty-in'

/** Satu anggota PageList - nilai teks per nama properti. */
export type Baris = Record<string, string>

/** Halaman kerja - SAMA dengan `models.Halaman`. Jalur relatif `pyWorkPage`. */
export interface Halaman {
  nilai: Record<string, string>
  daftar: Record<string, Baris[] | null>
  pesan?: Record<string, string[]>
}

/** Keadaan kerja satu kasus - `models.Kasus`. */
export interface Kasus {
  id: string
  position: string
  statusWork: string
  positionNote: string
  noPolis: string
  generasiTertutup: boolean
  createOp: string
  tglCreate: string
}

/** Satu baris daftar portal - `models.RingkasanKasus`. */
export interface RingkasanKasus {
  id: string
  businessName: string
  insuredName: string
  marketingName: string
  nbStatus: string
  statusWork: string
  positionNote: string
  noPolis: string
  tglCreate: string
}

/** Tombol submit yang tampil - `models.TombolKirim`. */
export type TombolKirim = '' | 'kirim' | 'konfirmasi-tolak' | 'nomor-polis'

/** Satu kasus siap tampil - `services.Layar`. */
export interface Layar {
  kasus: Kasus
  halaman: Halaman
  bolehKerja: boolean
  tombol: TombolKirim
  medanWajib: string[] | null
  tempat: Record<string, boolean>
  pesan?: string[]
}

export interface Pilihan {
  nilai: string
  label: string
}

/** Daftar pilihan layar - `services.Acuan`. */
export interface Acuan {
  mataUang: Pilihan[] | null
  mo: Pilihan[] | null
  spreading: Pilihan[] | null
  jenisReas: Pilihan[] | null
}

/** Satu baris view kontrak (nama kolom view) - `models.BarisKontrak`. */
export type BarisKontrak = Record<string, string>

/** Satu baris RD `BrowseAgentHierarkiList_RD` (pemilih SOB) - `models.BarisAgen`. */
export interface BarisAgen {
  id: string
  clientName: string
  leader0: string
  childCount: string
  clientId: string
}

export interface NomorPolis {
  id: string
  policyNo: string
}

export interface Riwayat {
  idPega: string
  status: string
  username: string
  workbasket: string
  operatorId: string
  tglTransfer: string
}

/** Action set satu sel - `services.PermintaanHitung`, SATU bentuk: `urutan`
 *  berisi satu refresh atau lebih, dijalankan berurutan atas halaman yang sama. */
export interface PermintaanHitung {
  urutan: { aksi: string; param?: string }[]
  indeks?: number
  halaman: Halaman
}

const kasus = (id: string) => `${PREFIX_NBTREATYIN}/kasus/${encodeURIComponent(id)}`

export function daftarKasus(cari: string, posisi = ''): Promise<RingkasanKasus[]> {
  return minta<RingkasanKasus[]>(`${PREFIX_NBTREATYIN}/kasus`, {
    kueri: { cari: cari || undefined, posisi: posisi || undefined },
  })
}

export function buatKasus(): Promise<Kasus> {
  return minta<Kasus>(`${PREFIX_NBTREATYIN}/kasus`, { metode: 'POST' })
}

export function bukaKasus(id: string): Promise<Layar> {
  return minta<Layar>(kasus(id))
}

export function simpanKasus(id: string, halaman: Halaman): Promise<Layar> {
  return minta<Layar>(kasus(id), { metode: 'PUT', badan: { halaman } })
}

export function hitung(id: string, p: PermintaanHitung): Promise<Layar> {
  return minta<Layar>(`${kasus(id)}/hitung`, { metode: 'POST', badan: p })
}

export function pilihBisnis(id: string, idDetail: string, halaman: Halaman): Promise<Layar> {
  return minta<Layar>(`${kasus(id)}/pilih-bisnis`, { metode: 'POST', badan: { idDetail, halaman } })
}

/** Jawaban klik satu baris popup `SOB` - `services.HasilSumberBisnis`: nilai yang
 *  ditulis `SearchHierarkiSourceBizAgent_PostDT`, per jalur halaman (`Quotation.*`). */
export interface HasilSumberBisnis {
  nilai: Record<string, string>
}

/** Klik satu baris popup `SOB` - pra-proses `SearchHierarkiSourceBizAgent_PostDT`,
 *  TANPA simpan (F4): hasilnya dipegang layar (`pegangSumberBisnis`) dan ikut
 *  terkirim pada Save/Submit/refresh; server menerimanya hanya bila cocok dengan
 *  RD `BrowseAgentHierarkiList_RD` yang dijalankan ulang. */
export function pilihSumberBisnis(id: string, idAgen: string, halaman: Halaman): Promise<HasilSumberBisnis> {
  return minta<HasilSumberBisnis>(`${kasus(id)}/pilih-sumber-bisnis`, { metode: 'POST', badan: { idAgen, halaman } })
}

export function terbitkanNomor(id: string, halaman: Halaman): Promise<NomorPolis> {
  return minta<NomorPolis>(`${kasus(id)}/nomor-polis`, { metode: 'POST', badan: { halaman } })
}

/** Akibat satu submit - `services.HasilKirim`. `pesanKonversi` = `FlagErrorKonversi`
 *  bila konversi Arasapas sesudah selesai gagal (penyimpanan tetap berhasil). */
export interface HasilKirim {
  kasus: Kasus
  pesanKonversi?: string
}

export function kirimKasus(id: string, halaman: Halaman): Promise<HasilKirim> {
  return minta<HasilKirim>(`${kasus(id)}/kirim`, { metode: 'POST', badan: { halaman } })
}

export function riwayatKasus(id: string): Promise<Riwayat[]> {
  return minta<Riwayat[]>(`${kasus(id)}/riwayat`)
}

export function daftarBisnis(cari: string): Promise<BarisKontrak[]> {
  return minta<BarisKontrak[]>(`${PREFIX_NBTREATYIN}/bisnis`, { kueri: { cari: cari || undefined } })
}

/** Isi TreeGrid popup `SOB` (`Section/SourceHierarki`). */
export function daftarSumberBisnis(): Promise<BarisAgen[]> {
  return minta<BarisAgen[]>(`${PREFIX_NBTREATYIN}/sumber-bisnis`)
}

export function ambilAcuan(): Promise<Acuan> {
  return minta<Acuan>(`${PREFIX_NBTREATYIN}/acuan`)
}

// ------------------------------------------------------------------ halaman

/** Awalan jalur halaman polis - SATU-SATUNYA salinan `models.HalamanPolis` + ".". */
export const POLIS = 'PolicyTreatyIn.'

/** Awalan jalur halaman master kontrak - `models.HalamanMaster` + ".". */
export const MASTER = 'TreatyIn.'

/** Nilai `.ClaimType` XOL Retro - `models.KlaimXOLRetro`. */
export const KLAIM_XOL_RETRO = 'XOL Retro'

/** Nilai satu jalur halaman ("" bila tidak ada). */
export function nilai(h: Halaman, jalur: string): string {
  return h.nilai?.[jalur] ?? ''
}

/** Salinan halaman dengan satu nilai diganti. */
export function setel(h: Halaman, jalur: string, v: string): Halaman {
  return { ...h, nilai: { ...h.nilai, [jalur]: v } }
}

/** Baris satu daftar (kosong bila tidak ada). */
export function daftar(h: Halaman, jalur: string): Baris[] {
  return h.daftar?.[jalur] ?? []
}

/** Salinan halaman dengan satu daftar diganti. */
export function setelDaftar(h: Halaman, jalur: string, b: Baris[]): Halaman {
  return { ...h, daftar: { ...h.daftar, [jalur]: b } }
}
