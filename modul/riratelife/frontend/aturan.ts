// Aturan layar R/I Rate Life - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/services`):
// nama wajib dan tidak kembar, validasi CSV, hak View only.

import type { IsianRate, Rate, Saringan } from './api'
import { RR } from './labels'

/** `pyPageSize` 50 (b12645) - sama dengan backend `models.UkuranHalaman`. */
export const UKURAN_HALAMAN = 50

/** Grid Rate Detail `pyPageSizeOther` 20 (View Detail.xml b10206) - sama dengan backend `models.UkuranHalamanRate`. */
export const UKURAN_HALAMAN_RATE = 20

/** Pilihan GENDER (radio `.GENDER` b1938; nilainya dari isi RATE_LIFE DEV, sama dengan backend `models.Gender`). */
export const GENDER = ['U', 'M', 'F'] as const

/** Form Rate Detail kosong (Clear Field `clearInputFieldRIRate_act` b1142 / Cancel `NewData_DT` b3497). */
export const ISIAN_RATE_KOSONG: IsianRate = { gender: '', contract: '', age: '', rate: '' }

/** Edit (`EditList_DT` b9853): form diisi dari baris grid. */
export function isianDariRate(r: Rate): IsianRate {
  return { gender: r.gender.trim().toUpperCase(), contract: r.contract.trim(), age: r.age.trim(), rate: r.rate.trim() }
}

/** Pemeriksaan awal form Rate Detail: CONTRACT (b2141) dan RATE (b2579) wajib; `null` = boleh dikirim. */
export function periksaIsianRate(isi: IsianRate): string | null {
  if (isi.contract.trim() === '') return RR.galatContract
  if (isi.rate.trim() === '') return RR.galatRate
  return null
}

/** Batas berkas CSV - sama dengan backend `models.MaksBytesCSV`. */
export const MAKS_BYTES_CSV = 4 * 1024 * 1024

/** Kolom yang boleh diurutkan (backend `models.KolomUrut`). */
export const KOLOM_URUT = ['id', 'usedby', 'operatorid', 'modifieddate'] as const
export type KolomUrut = (typeof KOLOM_URUT)[number]

/** Saringan awal: tanpa filter, urutan bawaan backend (ID menurun), halaman 1. */
export const SARINGAN_AWAL: Saringan = { id: '', usedby: '', urut: '', turun: false, halaman: 1 }

/** Klik kepala kolom: kolom baru = menaik; kolom yang sama = balik arah. Halaman kembali ke 1. */
export function gantiUrut(s: Saringan, kolom: KolomUrut): Saringan {
  if (s.urut === kolom) return { ...s, turun: !s.turun, halaman: 1 }
  return { ...s, urut: kolom, turun: false, halaman: 1 }
}

/** Penanda arah di kepala kolom. */
export function tandaUrut(s: Saringan, kolom: KolomUrut): string {
  if (s.urut !== kolom) return ''
  return s.turun ? ' ▼' : ' ▲'
}

/** Jumlah halaman (minimal 1). */
export function jumlahHalaman(total: number, ukuran = UKURAN_HALAMAN): number {
  return Math.max(1, Math.ceil(total / ukuran))
}

/** Pemeriksaan awal form; `null` = boleh dikirim. */
export function periksaNama(nama: string): string | null {
  return nama.trim() === '' ? RR.galatNama : null
}

/** Berkas yang boleh diunggah: berakhiran .csv dan tidak lebih dari 4 MB. */
export function berkasSah(nama: string, ukuran: number): boolean {
  return /\.csv$/i.test(nama.trim()) && ukuran > 0 && ukuran <= MAKS_BYTES_CSV
}

/** Nomor baris pertama satu halaman (No). */
export function nomorAwal(halaman: number, ukuran = UKURAN_HALAMAN): number {
  return (Math.max(1, halaman) - 1) * ukuran + 1
}
