// Aturan layar R/I Comm Life - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/services`,
// `backend/models`): nama wajib dan tidak kembar, angka muat kolom flat, kembar (CONTRACT, YEAR), hak View only.

import type { IsianKomisi, Komisi, Saringan } from './api'
import { RC } from './labels'

/** `pyPageSize` 50 (ringkasan b10081, detail b9716) - sama dengan backend `models.UkuranHalaman`. */
export const UKURAN_HALAMAN = 50

/** YEAR `pyMax` 4 (`InboxRIComm` b2206); TEPAT 4 angka [penyimpangan sadar - menunggu WO] - backend `models.NormalYear`. */
export const DIGIT_YEAR = 4

/** Form R/I COMM DETAIL kosong (Clear Field `clearInputFieldRIComm_act` b1123 / Cancel `NewData_DT` b3258). */
export const ISIAN_KOMISI_KOSONG: IsianKomisi = { contract: '', year: '', comm: '' }

/** Angka API (titik desimal kanonik) ditampilkan berkoma desimal: `12.5` -> `12,5`. */
export function tampilDesimal(s: string): string {
  return s.replace('.', ',')
}

/** Edit (`EditList_DT` b9323): form diisi dari baris grid. */
export function isianDariKomisi(k: Komisi): IsianKomisi {
  return { contract: k.contract.trim(), year: k.year.trim(), comm: tampilDesimal(k.comm.trim()) }
}

/** Pemeriksaan awal form detail: CONTRACT (b1923), YEAR (b2201), COMM (b2412) wajib; YEAR tepat 4 angka, COMM tidak
 * negatif (backend memeriksa ulang semuanya); `null` = boleh dikirim. */
export function periksaIsianKomisi(isi: IsianKomisi): string | null {
  if (isi.contract.trim() === '') return RC.galatContract
  if (isi.year.trim() === '') return RC.galatYear
  if (isi.comm.trim() === '') return RC.galatComm
  if (!/^[1-9][0-9]{3}$/.test(isi.year.trim())) return RC.galatYearEmpat
  if (isi.comm.trim().startsWith('-')) return RC.galatCommNegatif
  return null
}

/** Batas berkas CSV - sama dengan backend `models.MaksBytesCSV`. */
export const MAKS_BYTES_CSV = 4 * 1024 * 1024

/** Kolom yang boleh diurutkan (backend `models.KolomUrut`). */
export const KOLOM_URUT = ['id', 'usedby', 'operatorid', 'modifieddate'] as const
export type KolomUrut = (typeof KOLOM_URUT)[number]

/** Saringan awal: tanpa filter, urutan bawaan backend (ID menaik, sort XML b9857), halaman 1. */
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
  return nama.trim() === '' ? RC.galatNama : null
}

/** Berkas yang boleh diunggah: berakhiran .csv dan tidak lebih dari 4 MB. */
export function berkasSah(nama: string, ukuran: number): boolean {
  return /\.csv$/i.test(nama.trim()) && ukuran > 0 && ukuran <= MAKS_BYTES_CSV
}

/** Nomor baris pertama satu halaman (No). */
export function nomorAwal(halaman: number, ukuran = UKURAN_HALAMAN): number {
  return (Math.max(1, halaman) - 1) * ukuran + 1
}
