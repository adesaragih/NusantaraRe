// Aturan layar R/I Risk - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/services`,
// `backend/models`): angka muat kolom RIRISK_LIFE, kembar (CONTRACT, YEAR, MONTH), hak View only.

import type { IsianRincian, Rincian, Saringan } from './api'
import { RK } from './labels'

/** `pyPageSize` 50 grid ringkasan (`InboxSummaryRIRisk` b10030) - backend `models.UkuranHalaman`. */
export const UKURAN_HALAMAN = 50

/** `pyPageSize` Other = 200 grid R/I RISK DETAIL (`InboxRIRisk` b10169/b10243) - backend `models.UkuranHalamanRincian`. */
export const UKURAN_HALAMAN_RINCIAN = 200

/** YEAR / MONTH `pyMax` 4 (`InboxRIRisk` b2220 / b2407), keduanya TIDAK wajib - backend `models.NormalYear/NormalMonth`. */
export const DIGIT_YEAR_MONTH = 4

/** CONTRACT - lebar kolom VARCHAR2(10) warisan (backend `models.DigitTeksAngka`). */
export const DIGIT_CONTRACT = 10

/** Form R/I RISK DETAIL kosong (Clear Field `clearInputFieldRIRisk_act` b1148 / Cancel `NewRIRiskLife_Act` b3423). */
export const ISIAN_RINCIAN_KOSONG: IsianRincian = { contract: '', year: '', month: '', risk: '' }

/** Angka API (titik desimal kanonik) ditampilkan berkoma desimal: `921.9` -> `921,9` (aturan COMM ricommlife). */
export function tampilDesimal(s: string): string {
  return s.replace('.', ',')
}

/** EDIT (`EditRIRiskLife_Act` b9808: ID, Contract, Year, risk, UsedTo, Month): form diisi dari baris grid. */
export function isianDariRincian(k: Rincian): IsianRincian {
  return { contract: k.contract.trim(), year: k.year.trim(), month: k.month.trim(), risk: tampilDesimal(k.risk.trim()) }
}

/** Pemeriksaan awal form detail: CONTRACT (b1933) dan RISK (b2595) wajib; YEAR / MONTH kosong atau bulat <= 4 angka;
 * RISK tidak negatif (backend memeriksa ulang semuanya); `null` = boleh dikirim. */
export function periksaIsianRincian(isi: IsianRincian): string | null {
  if (isi.contract.trim() === '') return RK.galatContract
  if (isi.risk.trim() === '') return RK.galatRisk
  if (!/^[0-9]+$/.test(isi.contract.trim()) || isi.contract.trim().length > DIGIT_CONTRACT) return RK.galatAngka(RK.contract, DIGIT_CONTRACT)
  for (const [medan, nilai] of [
    [RK.year, isi.year],
    [RK.month, isi.month],
  ] as const) {
    const t = nilai.trim()
    if (t !== '' && !new RegExp(`^[0-9]{1,${DIGIT_YEAR_MONTH}}$`).test(t)) return RK.galatAngka(medan, DIGIT_YEAR_MONTH)
  }
  if (isi.risk.trim().startsWith('-')) return RK.galatRiskNegatif
  return null
}

/** Batas berkas CSV - sama dengan backend `models.MaksBytesCSV`. */
export const MAKS_BYTES_CSV = 4 * 1024 * 1024

/** Kolom yang boleh diurutkan (backend `models.KolomUrut`). */
export const KOLOM_URUT = ['id', 'usedby', 'operatorid', 'modifieddate'] as const
export type KolomUrut = (typeof KOLOM_URUT)[number]

/** Saringan awal: tanpa filter, urutan bawaan backend (ID menaik, sort XML b9806), halaman 1. */
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

/** Pemeriksaan awal form ringkasan (R/I RISK NAME wajib b1274); `null` = boleh dikirim. */
export function periksaNama(nama: string): string | null {
  return nama.trim() === '' ? RK.galatNama : null
}

/** Berkas yang boleh diunggah: berakhiran .csv dan tidak lebih dari 4 MB. */
export function berkasSah(nama: string, ukuran: number): boolean {
  return /\.csv$/i.test(nama.trim()) && ukuran > 0 && ukuran <= MAKS_BYTES_CSV
}

/** Nomor baris pertama satu halaman (No). */
export function nomorAwal(halaman: number, ukuran = UKURAN_HALAMAN): number {
  return (Math.max(1, halaman) - 1) * ukuran + 1
}
