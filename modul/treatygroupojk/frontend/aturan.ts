// Aturan layar Treaty Group OJK - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/services`):
// huruf besar, nama tidak kembar, panjang kolom, Order No otomatis, hak View only.

import type { Isian, Ojk } from './api'
import { TGO } from './labels'

/** Isian form dari satu baris (Edit / View) atau kosong (Add). */
export function isianDari(o: Ojk | null): Isian {
  return { name: o?.name ?? '', nameIdn: o?.nameIdn ?? '' }
}

/** Isian yang dikirim: spasi tepi dibuang (huruf besar dikerjakan backend). */
export function rapikanIsian(i: Isian): Isian {
  return { name: i.name.trim(), nameIdn: i.nameIdn.trim() }
}

/** Pemeriksaan awal sebelum dikirim; `null` = boleh dikirim. Backend tetap memeriksa seluruhnya. */
export function periksaIsian(i: Isian): string | null {
  if (i.name.trim() === '') return TGO.galatNama
  if (i.nameIdn.trim() === '') return TGO.galatNamaIdn
  return null
}
