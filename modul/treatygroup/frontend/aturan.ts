// Aturan layar Treaty Group - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/services`):
// huruf besar, nama tidak kembar, salinan OJK, COAID opsi A, hak View only.

import type { Grup, Isian, Ojk } from './api'
import { TG } from './labels'

/** Isian form dari satu baris (Edit) atau kosong (Add). */
export function isianDari(g: Grup | null): Isian {
  return { ojkId: g?.ojkId ?? '', name: g?.name ?? '', soaName: g?.soaName ?? '' }
}

/** Isian yang dikirim: spasi tepi dibuang (huruf besar dikerjakan backend). */
export function rapikanIsian(i: Isian): Isian {
  return { ojkId: i.ojkId.trim(), name: i.name.trim(), soaName: i.soaName.trim() }
}

/** Pemeriksaan awal sebelum dikirim; `null` = boleh dikirim. Backend tetap memeriksa seluruhnya. */
export function periksaIsian(i: Isian): string | null {
  if (i.ojkId.trim() === '') return TG.galatOjk
  if (i.name.trim() === '') return TG.galatNama
  return null
}

/** Kode dan nama: `01 - PROPERTY`; nama kosong = kode saja; kode kosong = kosong. */
export function kodeNama(kode: string, nama: string): string {
  if (kode === '') return ''
  return nama === '' ? kode : `${kode} - ${nama}`
}

/**
 * Opsi OJK Business, urut seperti dari backend (Order No). OJK tersimpan yang tidak ada lagi di daftar tetap tampil
 * sebagai nilai lama.
 */
export function opsiOjk(ojk: Ojk[], terpasang: string, namaTerpasang = ''): { value: string; label: string }[] {
  const opsi = ojk.map((o) => ({ value: o.id, label: kodeNama(o.id, o.name) }))
  if (terpasang !== '' && !opsi.some((o) => o.value === terpasang)) {
    opsi.push({ value: terpasang, label: TG.nilaiLama(kodeNama(terpasang, namaTerpasang)) })
  }
  return opsi
}

/**
 * COA yang tampil di form (perintah work owner 05-10-2026: "ubah pas pilih OJK Business"): OJK yang sama dengan baris
 * tersimpan = COA baris itu (dibiarkan); OJK lain = COA OJK itu (sama dengan yang disimpan backend); kosong = belum ada.
 */
export function coaTampil(ojk: Ojk[], ojkId: string, lama: Grup | null): string {
  if (lama !== null && lama.ojkId === ojkId) return kodeNama(lama.coaId, lama.coaName)
  const o = ojk.find((x) => x.id === ojkId)
  return o === undefined ? '' : kodeNama(o.coaId, o.coaName)
}
