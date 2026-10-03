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

/** Satu refresh berhitung - `services.PermintaanHitung`. */
export interface PermintaanHitung {
  aksi: string
  param?: string
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

export function ambilAcuan(): Promise<Acuan> {
  return minta<Acuan>(`${PREFIX_NBTREATYIN}/acuan`)
}

// ------------------------------------------------------------------ halaman

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
