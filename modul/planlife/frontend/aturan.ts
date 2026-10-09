// Aturan layar Plan - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/models`,
// `backend/services`): wajib, unik, pencocokan ulang master (K4), hak View only, ID dari sequence.

import type { Isian, PilihanBenefit, PilihanBusiness, Saringan } from './api'
import { PL } from './labels'

/** `pyPageSize` 10 (`InboxProductType` b4401) - backend `models.UkuranHalaman`. */
export const UKURAN_HALAMAN = 10

/** Lebar kolom nama (byte) - backend `models.BatasNama` (di luar XML: batas kolom 947). */
export const BATAS_NAMA = 200

/** Batas pilihan yang ditampilkan dropdown sekali ketik. */
export const BATAS_TAMPIL = 50

/** Saringan awal: urutan bawaan backend (ID menaik), halaman 1. */
export const SARINGAN_AWAL: Saringan = { urut: '', turun: false, halaman: 1 }

export const ISIAN_KOSONG: Isian = { coverName: '', business: '', benefit: '' }

/** Klik kepala kolom yang dapat diurutkan (Plan Name b4247, Benefit b4291): baru = menaik, sama = balik arah. */
export function gantiUrut(s: Saringan, kolom: 'covername' | 'benefit'): Saringan {
  if (s.urut === kolom) return { ...s, turun: !s.turun, halaman: 1 }
  return { ...s, urut: kolom, turun: false, halaman: 1 }
}

export function tandaUrut(s: Saringan, kolom: 'covername' | 'benefit'): string {
  if (s.urut !== kolom) return ''
  return s.turun ? ' ▼' : ' ▲'
}

export function jumlahHalaman(total: number, ukuran = UKURAN_HALAMAN): number {
  return Math.max(1, Math.ceil(total / ukuran))
}

const kunci = (s: string) => s.trim().toUpperCase()

/** Autocomplete Business: cari "memuat" pada Note (`pyUseForSearch` b1230), tanpa beda huruf. */
export function saringBusiness(daftar: PilihanBusiness[], ketik: string): PilihanBusiness[] {
  const k = kunci(ketik)
  return daftar.filter((b) => k === '' || kunci(b.note).includes(k)).slice(0, BATAS_TAMPIL)
}

/** Autocomplete Benefit: cari "memuat" pada Benefit (`pyUseForSearch` b1550), tanpa beda huruf. */
export function saringBenefit(daftar: PilihanBenefit[], ketik: string): PilihanBenefit[] {
  const k = kunci(ketik)
  return daftar.filter((b) => k === '' || kunci(b.benefit).includes(k)).slice(0, BATAS_TAMPIL)
}

/** Pemeriksaan awal form (K5): ketiga medan wajib dan <= 200 byte; `null` = boleh dikirim. Pencocokan master (K4) dan
 * Plan Name unik diperiksa backend (422 berkalimat). */
export function periksaIsian(isi: Isian): string | null {
  for (const [medan, nilai] of [
    [PL.planName, isi.coverName],
    [PL.business, isi.business],
    [PL.benefit, isi.benefit],
  ] as const) {
    const t = nilai.trim()
    if (t === '') return PL.galatWajib(medan)
    if (new TextEncoder().encode(t).length > BATAS_NAMA) return PL.galatPanjang(medan, BATAS_NAMA)
  }
  return null
}

/** Edit (`EditProductTypeLife_Act` b4045: ID, Business, BusinessID, Benefit, BenefitID, CoverName): form dari baris grid. */
export function isianDariPlan(p: { coverName: string; business: string; benefit: string }): Isian {
  return { coverName: p.coverName, business: p.business, benefit: p.benefit }
}
