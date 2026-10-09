// Klien R/I Risk - `/api/ri-risk-life` (`modul/ririsklife/backend/handlers/rute.go`). Rutenya hanya terbuka bagi
// pemegang menu `ririsklife`; tulis (termasuk View Upload) ditolak bagi akses View only.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_RK = '/api/ri-risk-life'

/** Batas waktu simpan unggahan sampai 10.000 baris. */
const BATAS_BERAT_MS = 180_000

/** Satu ringkasan - `models.Ringkasan`. */
export interface Ringkasan {
  id: string
  usedby: string
  operatorId: string
  /** Teks apa adanya (format Pega `YYYYMMDDTHHMMSS.mmm GMT`). */
  modifiedDate: string
  /** Tampilan WIB `DD-MM-YYYY`. */
  diubah: string
  /** Hanya dari `ambil` (dialog Delete). */
  jumlahRincian?: number
}

/** Satu baris R/I RISK DETAIL - `models.Rincian` (kolom RIRISK_LIFE; RISK bertitik desimal kanonik; YEAR / MONTH boleh
 * kosong). */
export interface Rincian {
  id: string
  idUsedBy: string
  usedby: string
  contract: string
  year: string
  month: string
  risk: string
}

/** Satu halaman grid - `models.Halaman`. */
export interface Halaman<T> {
  daftar: T[]
  total: number
  halaman: number
  ukuran: number
}

export interface Saringan {
  id: string
  usedby: string
  urut: string
  turun: boolean
  halaman: number
}

/** Baris CSV sah - `models.BarisCSV`. */
export interface BarisCSV {
  baris: number
  usedby: string
  contract: string
  year: string
  month: string
  risk: string
}

export interface GalatBaris {
  baris: number
  pesan: string
}

export interface RingkasanUnggah {
  usedby: string
  id: string
  baru: boolean
  jumlah: number
}

/** View Upload - `services.HasilUnggah`. */
export interface HasilUnggah {
  baris: BarisCSV[]
  galat: GalatBaris[]
  ringkasan: RingkasanUnggah[]
  sah: boolean
}

/** Simpan Upload - `services.HasilSimpanUnggah`. */
export interface HasilSimpanUnggah {
  disimpan: number
  ringkasanBaru: number
  ringkasan: RingkasanUnggah[]
}

export function ambilDaftar(s: Saringan): Promise<Halaman<Ringkasan>> {
  return minta(PREFIX_RK, {
    kueri: {
      id: s.id.trim() === '' ? undefined : s.id.trim(),
      usedby: s.usedby.trim() === '' ? undefined : s.usedby.trim(),
      urut: s.urut === '' ? undefined : s.urut,
      arah: s.urut === '' ? undefined : s.turun ? 'desc' : 'asc',
      halaman: s.halaman > 1 ? s.halaman : undefined,
    },
  })
}

/** Satu ringkasan beserta jumlah rinciannya. */
export function ambil(id: string): Promise<Ringkasan> {
  return minta(`${PREFIX_RK}/${encodeURIComponent(id)}`)
}

/** Save - Add (ID = site || LPAD(M_RIRISK_LIFE_SUMMARY_SEQ, 6), dibentuk server). Keputusan work owner 08-10-2026:
 * Save ditampilkan walau wadahnya `1=2` di XML (b1511). */
export function tambah(usedby: string): Promise<Ringkasan> {
  return minta(PREFIX_RK, { metode: 'POST', badan: { usedby } })
}

/** Save - Edit (nama baru ikut ke rinciannya, satu transaksi). */
export function ubah(id: string, usedby: string): Promise<Ringkasan> {
  return minta(`${PREFIX_RK}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: { usedby } })
}

/** Delete - ringkasan beserta rinciannya. */
export function hapus(id: string): Promise<{ id: string; rincianTerhapus: number }> {
  return minta(`${PREFIX_RK}/${encodeURIComponent(id)}`, { metode: 'DELETE' })
}

/** R/I RISK DETAIL - satu halaman (200 baris). */
export function ambilDetail(id: string, halaman: number): Promise<Halaman<Rincian>> {
  return minta(`${PREFIX_RK}/${encodeURIComponent(id)}/detail`, {
    kueri: { halaman: halaman > 1 ? halaman : undefined },
  })
}

/** Isian form R/I RISK DETAIL - `models.IsianRincian` (USEDBY dan IDUSEDBY diambil server dari ringkasan). */
export interface IsianRincian {
  contract: string
  year: string
  month: string
  risk: string
}

/** R/I RISK DETAIL Save - tambah baris (`AddToList_Act`). */
export function tambahDetail(id: string, isi: IsianRincian): Promise<Rincian> {
  return minta(`${PREFIX_RK}/${encodeURIComponent(id)}/detail`, { metode: 'POST', badan: isi })
}

/** R/I RISK DETAIL Save - ubah baris yang diisi EDIT (`EditRIRiskLife_Act`). */
export function ubahDetail(id: string, idDetail: string, isi: IsianRincian): Promise<Rincian> {
  return minta(`${PREFIX_RK}/${encodeURIComponent(id)}/detail/${encodeURIComponent(idDetail)}`, { metode: 'PUT', badan: isi })
}

/** View Upload - tanpa menulis. */
export function pratinjauUnggah(csv: string): Promise<HasilUnggah> {
  return minta(`${PREFIX_RK}/unggah/pratinjau`, { metode: 'POST', badan: { csv }, batasWaktuMs: BATAS_BERAT_MS })
}

/** Simpan Upload. */
export function simpanUnggah(csv: string): Promise<HasilSimpanUnggah> {
  return minta(`${PREFIX_RK}/unggah`, { metode: 'POST', badan: { csv }, batasWaktuMs: BATAS_BERAT_MS })
}
