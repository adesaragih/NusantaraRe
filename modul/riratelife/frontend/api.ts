// Klien R/I Rate Life - `/api/ri-rate-life` (`modul/riratelife/backend/handlers/rute.go`). Rutenya hanya terbuka bagi
// pemegang menu `riratelife`; tulis (termasuk View Upload) ditolak bagi akses View only.

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_RR = '/api/ri-rate-life'

/** Batas waktu baca rincian M_RATE_LIFE dan simpan unggahan sampai 10.000 baris. */
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
  jumlahRate?: number
}

/** Satu baris rate - `models.Rate`. */
export interface Rate {
  id: string
  idUsedBy: string
  usedby: string
  gender: string
  contract: string
  age: string
  rate: string
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
  gender: string
  age: string
  rate: string
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
  return minta(PREFIX_RR, {
    kueri: {
      id: s.id.trim() === '' ? undefined : s.id.trim(),
      usedby: s.usedby.trim() === '' ? undefined : s.usedby.trim(),
      urut: s.urut === '' ? undefined : s.urut,
      arah: s.urut === '' ? undefined : s.turun ? 'desc' : 'asc',
      halaman: s.halaman > 1 ? s.halaman : undefined,
    },
  })
}

/** Satu ringkasan beserta jumlah rate-nya. */
export function ambil(id: string): Promise<Ringkasan> {
  return minta(`${PREFIX_RR}/${encodeURIComponent(id)}`, { batasWaktuMs: BATAS_BERAT_MS })
}

/** Save - Add (ID dari SEQ_M_RATE_LIFE_SUMMARY). */
export function tambah(usedby: string): Promise<Ringkasan> {
  return minta(PREFIX_RR, { metode: 'POST', badan: { usedby } })
}

/** Save - Edit. */
export function ubah(id: string, usedby: string): Promise<Ringkasan> {
  return minta(`${PREFIX_RR}/${encodeURIComponent(id)}`, {
    metode: 'PUT',
    badan: { usedby },
    batasWaktuMs: BATAS_BERAT_MS,
  })
}

/** Delete - ringkasan beserta rate-nya. */
export function hapus(id: string): Promise<{ id: string; rateTerhapus: number }> {
  return minta(`${PREFIX_RR}/${encodeURIComponent(id)}`, { metode: 'DELETE', batasWaktuMs: BATAS_BERAT_MS })
}

/** Rate Detail - satu halaman. */
export function ambilRate(id: string, halaman: number): Promise<Halaman<Rate>> {
  return minta(`${PREFIX_RR}/${encodeURIComponent(id)}/rate`, {
    kueri: { halaman: halaman > 1 ? halaman : undefined },
    batasWaktuMs: BATAS_BERAT_MS,
  })
}

/** Isian form Rate Detail - `models.IsianRate` (R/I RATE NAME dan IDUSEDBY diambil server dari ringkasan). */
export interface IsianRate {
  gender: string
  contract: string
  age: string
  rate: string
}

/** Rate Detail Save - tambah baris (`AddToList_Act`). */
export function tambahRate(id: string, isi: IsianRate): Promise<Rate> {
  return minta(`${PREFIX_RR}/${encodeURIComponent(id)}/rate`, { metode: 'POST', badan: isi, batasWaktuMs: BATAS_BERAT_MS })
}

/** Rate Detail Save - ubah baris (`EditList_DT`). */
export function ubahRate(id: string, idRate: string, isi: IsianRate): Promise<Rate> {
  return minta(`${PREFIX_RR}/${encodeURIComponent(id)}/rate/${encodeURIComponent(idRate)}`, {
    metode: 'PUT',
    badan: isi,
    batasWaktuMs: BATAS_BERAT_MS,
  })
}

/** View Upload - tanpa menulis. */
export function pratinjauUnggah(csv: string): Promise<HasilUnggah> {
  return minta(`${PREFIX_RR}/unggah/pratinjau`, { metode: 'POST', badan: { csv }, batasWaktuMs: BATAS_BERAT_MS })
}

/** Simpan Upload. */
export function simpanUnggah(csv: string): Promise<HasilSimpanUnggah> {
  return minta(`${PREFIX_RR}/unggah`, { metode: 'POST', badan: { csv }, batasWaktuMs: BATAS_BERAT_MS })
}
