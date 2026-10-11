// Klien Plan - `/api/plan-life` (`modul/planlife/backend/handlers/rute.go`). Rutenya hanya terbuka bagi pemegang menu
// `planlife`; tulis ditolak bagi akses View only. Tidak ada hapus (XML tanpa Delete).

import { minta } from '../../../inti/frontend/klien'

export const PREFIX_PL = '/api/plan-life'

/** Satu baris PRODUCT_TYPE_LIFE - `models.Plan`. */
export interface Plan {
  id: string
  coverName: string
  business: string
  businessId: string
  benefit: string
  benefitId: string
}

/** Satu halaman grid - `models.Halaman`. */
export interface Halaman {
  daftar: Plan[]
  total: number
  halaman: number
  ukuran: number
}

/** Urut grid: '' = bawaan (ID menaik), 'covername' / 'benefit' (Business tidak dapat diurutkan b4269). */
export interface Saringan {
  urut: '' | 'covername' | 'benefit'
  turun: boolean
  halaman: number
}

/** Pilihan Business - `models.PilihanBusiness` (`BrowseBusinessLife_RD`: OLDID b1200, Note b1231, ID b1265). */
export interface PilihanBusiness {
  id: string
  oldId: string
  note: string
}

/** Pilihan Benefit - `models.PilihanBenefit` (`BrowseBenefitLife_RD`: Benefit b1552, Number b1563). */
export interface PilihanBenefit {
  id: string
  benefit: string
}

/** Isian form - `models.Isian` (teks saja; ID business / benefit diisi server dari master, K4). */
export interface Isian {
  coverName: string
  business: string
  benefit: string
}

export function ambilDaftar(s: Saringan): Promise<Halaman> {
  return minta(PREFIX_PL, {
    kueri: {
      urut: s.urut === '' ? undefined : s.urut,
      arah: s.urut === '' ? undefined : s.turun ? 'desc' : 'asc',
      halaman: s.halaman > 1 ? s.halaman : undefined,
    },
  })
}

export function ambilPilihanBusiness(): Promise<PilihanBusiness[]> {
  return minta(`${PREFIX_PL}/pilihan/business`)
}

export function ambilPilihanBenefit(): Promise<PilihanBenefit[]> {
  return minta(`${PREFIX_PL}/pilihan/benefit`)
}

/** Save - Add (`SaveProductTypeLife_Act` b6422; ID dibentuk server). */
export function tambah(isi: Isian): Promise<Plan> {
  return minta(PREFIX_PL, { metode: 'POST', badan: isi })
}

/** Save sesudah Edit (`EditProductTypeLife_Act` b4045). */
export function ubah(id: string, isi: Isian): Promise<Plan> {
  return minta(`${PREFIX_PL}/${encodeURIComponent(id)}`, { metode: 'PUT', badan: isi })
}
