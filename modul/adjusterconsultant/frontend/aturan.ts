// Aturan layar Adjuster Consultant - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend
// (`backend/services`): huruf besar, nama ganda, panjang kolom, hak View only.

import type { Adjuster, Isian, Status } from './api'

/** Urutan saringan status di layar: bawaan Active (keputusan work owner 05-10-2026). */
export const URUTAN_STATUS: readonly Status[] = ['active', 'inactive', '']

export const STATUS_AWAL: Status = 'active'

/** Isian form dari satu baris (Edit / View) atau kosong (Add). */
export function isianDari(a: Adjuster | null): Isian {
  return { name: a?.name ?? '', address: a?.address ?? '', telpNo: a?.telpNo ?? '' }
}

/** Pemeriksaan awal sebelum dikirim; `null` = boleh dikirim. Backend tetap memeriksa seluruhnya. */
export function periksaIsian(i: Isian, pesanNama: string): string | null {
  return i.name.trim() === '' ? pesanNama : null
}

/** Isian yang dikirim: spasi tepi dibuang (huruf besar dikerjakan backend, seperti Pega). */
export function rapikanIsian(i: Isian): Isian {
  return { name: i.name.trim(), address: i.address.trim(), telpNo: i.telpNo.trim() }
}

/** Nama tampilan di kalimat konfirmasi. */
export function namaTampil(a: Adjuster): string {
  return a.name.trim() === '' ? a.id : `${a.name} (${a.id})`
}
