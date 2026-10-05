// Aturan layar Treaty Description - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend
// (`backend/services`): huruf besar, nama tidak kembar, jenis dan status sah, hak View only.

import type { Desc, Isian } from './api'
import { TD } from './labels'

/** `DESCNAME VARCHAR2(100)` - byte. */
export const BATAS_NAMA = 100

/** Isian form dari satu baris (Edit) atau bawaan Add: Non XOL, Active. NULL (baris lama) dibaca Active. */
export function isianDari(d: Desc | null): Isian {
  if (d === null) return { descName: '', isXol: '0', statusAktif: '1' }
  return { descName: d.descName, isXol: d.isXol === '1' ? '1' : '0', statusAktif: d.aktif ? '1' : '0' }
}

/** Isian yang dikirim: spasi tepi dibuang, nama huruf besar. */
export function rapikanIsian(i: Isian): Isian {
  return { ...i, descName: i.descName.trim().toUpperCase() }
}

const panjangByte = (s: string) => new TextEncoder().encode(s).length

/** Pemeriksaan awal sebelum dikirim; `null` = boleh dikirim. Backend tetap memeriksa seluruhnya. */
export function periksaIsian(i: Isian): string | null {
  if (i.descName.trim() === '') return TD.galatNama
  if (panjangByte(i.descName.trim()) > BATAS_NAMA) return TD.galatPanjang
  return null
}

export const labelJenis = (isXol: string): string => (isXol === '1' ? TD.xol : TD.nonXol)
export const labelStatus = (aktif: boolean): string => (aktif ? TD.aktif : TD.nonaktif)
