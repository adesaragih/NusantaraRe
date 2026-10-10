// Klien R/I Comm Life - `/api/ri-comm-life` (`modul/ricommlife/backend/handlers/rute.go`). Rutenya hanya terbuka bagi
// pemegang menu `ricommlife`; tulis (termasuk View Upload) ditolak bagi akses View only.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_RC = '/api/ri-comm-life'

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
  jumlahKomisi?: number
}

/** Satu baris R/I COMM DETAIL - `models.Komisi` (kolom M_RICOMM_LIFE; angka bertitik desimal kanonik). */
export interface Komisi {
  id: string
  idUsedBy: string
  usedby: string
  contract: string
  year: string
  comm: string
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
  comm: string
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
  return minta(PREFIX_RC, {
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
  return minta(`${PREFIX_RC}/${encodeURIComponent(id)}`)
}

/** Save - Add (ID = site + sequence warisan, dibentuk server). */
export function tambah(usedby: string): Promise<Ringkasan> {
  return minta(PREFIX_RC, { metode: 'POST', badan: { usedby } })
}

/** Save - Edit (nama baru ikut ke rinciannya). */
export function ubah(id: string, usedby: string): Promise<Ringkasan> {
  return minta(`${PREFIX_RC}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: { usedby } })
}

/** Delete - ringkasan beserta rinciannya. */
export function hapus(id: string): Promise<{ id: string; komisiTerhapus: number }> {
  return minta(`${PREFIX_RC}/${encodeURIComponent(id)}`, { metode: 'DELETE' })
}

/** R/I COMM DETAIL - satu halaman. */
export function ambilDetail(id: string, halaman: number): Promise<Halaman<Komisi>> {
  return minta(`${PREFIX_RC}/${encodeURIComponent(id)}/detail`, {
    kueri: { halaman: halaman > 1 ? halaman : undefined },
  })
}

/** Isian form R/I COMM DETAIL - `models.IsianKomisi` (USEDBY dan IDUSEDBY diambil server dari ringkasan). */
export interface IsianKomisi {
  contract: string
  year: string
  comm: string
}

/** R/I COMM DETAIL Save - tambah baris (`AddToList_Act`). */
export function tambahDetail(id: string, isi: IsianKomisi): Promise<Komisi> {
  return minta(`${PREFIX_RC}/${encodeURIComponent(id)}/detail`, { metode: 'POST', badan: isi })
}

/** R/I COMM DETAIL Save - ubah baris (`EditList_DT`). */
export function ubahDetail(id: string, idDetail: string, isi: IsianKomisi): Promise<Komisi> {
  return minta(`${PREFIX_RC}/${encodeURIComponent(id)}/detail/${encodeURIComponent(idDetail)}`, { metode: 'PUT', badan: isi })
}

/** View Upload - tanpa menulis. */
export function pratinjauUnggah(csv: string): Promise<HasilUnggah> {
  return minta(`${PREFIX_RC}/unggah/pratinjau`, { metode: 'POST', badan: { csv }, batasWaktuMs: BATAS_BERAT_MS })
}

/** Simpan Upload. */
export function simpanUnggah(csv: string): Promise<HasilSimpanUnggah> {
  return minta(`${PREFIX_RC}/unggah`, { metode: 'POST', badan: { csv }, batasWaktuMs: BATAS_BERAT_MS })
}
