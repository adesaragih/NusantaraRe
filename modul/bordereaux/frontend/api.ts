// Klien Bordereaux - `/api/bordereaux` (`modul/bordereaux/backend/handlers/rute.go`). Rutenya hanya terbuka bagi
// pemegang menu `bordereaux`.

import { BATAS_WAKTU_MS, kegagalanDari, minta, mintaFormulir, rakitURL, unduhBerkasBeridentitas } from '../../../inti/frontend/klien'
import { headerIdentitas } from '../../../inti/frontend/store/sesi'

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
  /** Upload File / Delete lampiran (`services.BolehLampiran`); terisi saat berkas dibuka. */
  lampiran: boolean
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

// ---------------------------------------------------------------------------- lampiran
// `AttachmentsBdx` / `AttachmentDetailBdx` - penyimpanan bersama `inti/backend/penyimpanan` (keputusan work owner
// 08-10-2026).

/** Satu kategori grid lampiran (`GetKategoryDocBDX_SQL`: kategori master + jumlah lampiran berkas ini). */
export interface KategoriLampiran {
  id: string
  nama: string
  cacah: number
}

/** Satu lampiran (`AttachDocumentBdx_SQL`). Kunci objek penyimpanan TIDAK dikirim backend. */
export interface Lampiran {
  id: string
  kategoriId: string
  /** Kolom Type popup (`CATEGORY`). */
  kategori: string
  fileName: string
  /** `FILEMIMETYPE` - ekstensi huruf kecil, mis. `xlsx`. */
  ekstensi: string
  username: string
}

function jalurLampiran(id: string, kategori?: string): string {
  const dasar = `${PREFIX_BDX}/berkas/${encodeURIComponent(id)}/lampiran`
  return kategori === undefined ? dasar : `${dasar}/${encodeURIComponent(kategori)}`
}

function jalurSatuLampiran(id: string, l: Lampiran): string {
  return `${jalurLampiran(id, l.kategoriId)}/${encodeURIComponent(l.id)}`
}

/** Grid kategori (`GetKategotyDocBdx`). */
export function ambilKategoriLampiran(id: string): Promise<{ daftar: KategoriLampiran[] }> {
  return minta(jalurLampiran(id))
}

/** Popup View File (`getAttcachmentList`). */
export function ambilLampiran(id: string, kategori: string): Promise<{ daftar: Lampiran[] }> {
  return minta(jalurLampiran(id, kategori))
}

/** Upload File - satu berkas per permintaan (`AttachDocBdx_Post` mengulang setiap berkas). */
export function unggahLampiran(id: string, kategori: string, berkas: File): Promise<Lampiran> {
  const isi = new FormData()
  isi.append('berkas', berkas)
  return mintaFormulir(jalurLampiran(id, kategori), isi)
}

/** Tautan nama berkas (`DownloadAttachmentBdx`) - fetch beridentitas, bukan pranala. */
export function unduhLampiran(id: string, l: Lampiran): Promise<void> {
  return unduhBerkasBeridentitas(`${jalurSatuLampiran(id, l)}/isi`, l.fileName)
}

/**
 * Isi satu lampiran sebagai Blob - rute unduh yang ADA (`DownloadAttachmentBdx`), berheader identitas - untuk `View`
 * pdf / gambar di popup penampil (permintaan work owner 08-10-2026, seperti Product Name Life).
 */
export async function ambilIsiLampiran(id: string, l: Lampiran): Promise<Blob> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  try {
    const jawab = await fetch(rakitURL(`${jalurSatuLampiran(id, l)}/isi`), {
      method: 'GET',
      headers: { ...headerIdentitas() },
      signal: kendali.signal,
    })
    if (!jawab.ok) throw kegagalanDari(jawab.status, await jawab.text())
    return await jawab.blob()
  } finally {
    clearTimeout(jam)
  }
}

/** View Office Online - URL bertanda tangan untuk penampil (`DownloadAttachmentBdx` ViewOffice=true). */
export function tautanOffice(id: string, l: Lampiran): Promise<{ url: string }> {
  return minta(`${jalurSatuLampiran(id, l)}/office`)
}

/** Delete (`DeleteAttachmentBdx`). */
export function hapusLampiran(id: string, l: Lampiran): Promise<{ ok: boolean }> {
  return minta(`${jalurSatuLampiran(id, l)}/hapus`, { metode: 'POST' })
}
