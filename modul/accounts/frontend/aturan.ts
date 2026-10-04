// Aturan murni layar Accounts - diuji `aturan.test.ts`. Otoritas tetap di backend; di sini hanya pemeriksaan wajib
// isi supaya pesan "Value cannot be blank" tampil di bawah medannya sebelum dikirim.

import type { Isian } from './api'
import { ACC } from './labels'

/** Pilihan Insured Name yang sedang terpilih di form. */
export interface InsuredTerpilih {
  id: string
  orgId: string
  nama: string
}

/** Galat per medan; kosong = lolos. */
export interface GalatIsian {
  insured?: string
  groupBusiness?: string
}

export function isianKosong(): Isian {
  return { insuredId: '', groupBusinessId: '', description: '' }
}

/** Wajib isi Insured Name dan Group Business - pesan VERBATIM Pega. */
export function periksa(isi: Isian): GalatIsian {
  const g: GalatIsian = {}
  if (isi.insuredId.trim() === '') g.insured = ACC.wajib
  if (isi.groupBusinessId.trim() === '') g.groupBusiness = ACC.wajib
  return g
}

export function adaGalat(g: GalatIsian): boolean {
  return g.insured !== undefined || g.groupBusiness !== undefined
}

/** Teks sel kosong. */
export function atauStrip(s: string): string {
  return s.trim() === '' ? '—' : s
}
