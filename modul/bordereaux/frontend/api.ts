// Klien Bordereaux - `/api/bordereaux` (`modul/bordereaux/backend/handlers/rute.go`). Rutenya hanya terbuka bagi
// pemegang menu `bordereaux`.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_BDX = '/api/bordereaux'

/** Batas waktu Upload CSV dan Save: sampai 5.000 baris. */
export const BATAS_WAKTU_BESAR_MS = 120_000

/** Satu baris detail: kolom tabel -> teks (angka bertitik desimal tanpa pemisah ribuan, tanggal DD-MM-YYYY). */
export type Baris = Record<string, string>

export interface Header {
  bdxId: string
  tanggal: string
  userInput: string
  type: string
  typeBusiness: string
  masterId: string
  cedingId: string
  cedingName: string
  sobId: string
  sobName: string
  treatyName: string
  reportStart: string
  reportEnd: string
  reffNoSoa: string
  reffNoBdx: string
  position: string
  statusAksep: string
}

export interface Hak {
  ubah: boolean
  hapus: boolean
  submit: boolean
  putuskan: boolean
}

export interface BarisDaftar extends Header {
  hak: Hak
}

export interface Halaman {
  daftar: BarisDaftar[]
  total: number
  halaman: number
  ukuran: number
}

export interface Ringkasan {
  currency: string
  reinsurer: string
  rnm: string
}

export interface Riwayat {
  tanggal: string
  pic: string
  isApproved: boolean
  komentar: string
}

export interface Rincian {
  header: Header
  baris: Baris[]
  ringkasan: Ringkasan[]
  riwayat: Riwayat[]
  hak: Hak
}

export interface Kolom {
  kolom: string
  /** Judul kepala templat CSV (pesan validasi). */
  judul: string
  jenis: 'teks' | 'angka' | 'tanggal'
  /** Judul kepala grid menurut format Excel bordereaux 2025 (boleh ber-baris baru); tidak ada = `judul`. */
  kepala?: string
  /** Judul sel gabungan di atas kolom ini (mis. `*PERIOD OF INSURANCE`); tidak ada = tanpa grup. */
  grup?: string
  /** Kolom tidak ada di format Excel: tidak ditampilkan grid. */
  sembunyi?: boolean
}

export interface Pilihan {
  type: string[]
  business: Record<string, string[]>
  status: string[]
  /** Kolom grid detail per `TYPE|BUSINESS`. */
  kolom: Record<string, Kolom[]>
  /** Slot Template Manager per `TYPE|BUSINESS`. */
  templat: Record<string, { kode: string; nama: string }>
  bolehBuat: boolean
  /** Tombol Copy Old Data - superadmin. */
  copyOld: boolean
}

export interface Filter {
  id: string
  type: string
  business: string
  reffSoa: string
  reffBdx: string
  ceding: string
  treaty: string
  /** DD-MM-YYYY. */
  start: string
  end: string
  position: string
  status: string
}

export interface Cedant {
  id: string
  nama: string
}

export interface MasterTreaty {
  id: string
  contractName: string
  reinsType: string
  sobId: string
  sobName: string
  cedingId: string
  cedingName: string
}

export interface PermintaanSimpan {
  bdxId: string
  type: string
  business: string
  masterId: string
  reportStart: string
  reportEnd: string
  reffNoSoa: string
  reffNoBdx: string
  baris: Baris[]
}

export function ambilPilihan(): Promise<Pilihan> {
  return minta<Pilihan>(`${PREFIX_BDX}/pilihan`)
}

export function ambilDaftar(f: Filter, halaman: number): Promise<Halaman> {
  const kueri: Record<string, string | number | undefined> = { halaman }
  for (const [k, v] of Object.entries(f)) kueri[k] = v === '' ? undefined : v
  return minta<Halaman>(PREFIX_BDX, { kueri })
}

export function buka(id: string): Promise<Rincian> {
  return minta<Rincian>(`${PREFIX_BDX}/berkas/${encodeURIComponent(id)}`)
}

export function unggahCsv(type: string, business: string, csv: string): Promise<{ baris: Baris[]; ringkasan: Ringkasan[] }> {
  return minta(`${PREFIX_BDX}/unggah-csv`, { metode: 'POST', badan: { type, business, csv }, batasWaktuMs: BATAS_WAKTU_BESAR_MS })
}

export function simpan(p: PermintaanSimpan): Promise<{ bdxId: string }> {
  return minta(`${PREFIX_BDX}/simpan`, { metode: 'POST', badan: p, batasWaktuMs: BATAS_WAKTU_BESAR_MS })
}

export function submit(id: string, setuju: boolean, komentar: string): Promise<{ ok: boolean }> {
  return minta(`${PREFIX_BDX}/berkas/${encodeURIComponent(id)}/submit`, { metode: 'POST', badan: { setuju, komentar } })
}

export function hapus(id: string): Promise<{ ok: boolean }> {
  return minta(`${PREFIX_BDX}/berkas/${encodeURIComponent(id)}/hapus`, { metode: 'POST' })
}

/** Satu irisan chart daftar - `models.IrisanChart`. */
export interface IrisanChart {
  business: string
  type: string
  cedingId: string
  cedingName: string
  jumlah: number
}

export function ambilChart(): Promise<{ irisan: IrisanChart[] }> {
  return minta(`${PREFIX_BDX}/chart`)
}

export function cariCedant(q: string): Promise<{ daftar: Cedant[] }> {
  return minta(`${PREFIX_BDX}/cedant`, { kueri: { q } })
}

export function cariTreaty(ceding: string): Promise<{ daftar: MasterTreaty[] }> {
  return minta(`${PREFIX_BDX}/master-treaty`, { kueri: { ceding } })
}

/** Batas waktu Copy Old Data: satu kelompok berkas bisa membawa ratusan baris detail. */
export const BATAS_WAKTU_LAMA_MS = 120_000

/** Satu baris popup Copy Old Data - `models.BerkasLama`: berkas yang isinya masih tertinggal di JSON lama. */
export interface BerkasLama {
  id: string
  type: string
  business: string
  ceding: string
  treaty: string
  status: string
  /** Baris detail di JSON lama dan di tabel detail kombinasinya. */
  barisJson: number
  barisTabel: number
  /** Baris CommentList JSON lama dan baris BORDEREAUX_HISTORY. */
  komentar: number
  riwayat: number
  /** Baris BORDEREAUX belum ada; dibuat dari JSON. */
  tanpaHeader: boolean
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

/** Isi popup Copy Old Data - superadmin. */
export function ambilLama(): Promise<{ daftar: BerkasLama[] }> {
  return minta(`${PREFIX_BDX}/lama`, { batasWaktuMs: BATAS_WAKTU_LAMA_MS })
}

/** Process Copy - superadmin; hasil per ID. */
export function salinLama(ids: string[]): Promise<JawabanSalinLama> {
  return minta(`${PREFIX_BDX}/lama/salin`, { metode: 'POST', badan: { ids }, batasWaktuMs: BATAS_WAKTU_LAMA_MS })
}
