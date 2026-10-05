// Aturan layar Business Group - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/services`):
// huruf besar, nama tidak kembar, tanpa SYARIAH, salinan nama Treaty Group, hak View only.

import type { BisnisGrup, Isian, TreatyGroup } from './api'
import { BG } from './labels'

/** Isian form dari satu baris (Edit) atau kosong (Add). */
export function isianDari(b: BisnisGrup | null): Isian {
  return { topId: b?.topId ?? '', name: b?.name ?? '', alias: b?.alias ?? '' }
}

/** Isian yang dikirim: spasi tepi dibuang (huruf besar dan Alias kosong = Name dikerjakan backend). */
export function rapikanIsian(i: Isian): Isian {
  return { topId: i.topId.trim(), name: i.name.trim(), alias: i.alias.trim() }
}

/** Nama berakhiran SYARIAH (tanpa beda huruf dan spasi tepi) - tidak dikelola modul ini. */
export function syariah(nama: string): boolean {
  return nama.trim().toUpperCase().endsWith('SYARIAH')
}

/** Pemeriksaan awal sebelum dikirim; `null` = boleh dikirim. Backend tetap memeriksa seluruhnya. */
export function periksaIsian(i: Isian): string | null {
  if (i.topId.trim() === '') return BG.galatTreaty
  if (i.name.trim() === '') return BG.galatNama
  if (syariah(i.name)) return BG.galatSyariah
  return null
}

/** Kode dan nama: `10007 - PROPERTY`; nama kosong = kode saja; kode kosong = kosong. */
export function kodeNama(kode: string, nama: string): string {
  if (kode === '') return ''
  return nama === '' ? kode : `${kode} - ${nama}`
}

/** Opsi Treaty Group; TOPID tersimpan yang tidak ada lagi di TREATYGROUP tetap tampil sebagai nilai lama. */
export function opsiTreaty(
  treaty: TreatyGroup[],
  terpasang: string,
  namaTerpasang = '',
): { value: string; label: string }[] {
  const opsi = treaty.map((t) => ({ value: t.id, label: kodeNama(t.id, t.name) }))
  if (terpasang !== '' && !opsi.some((o) => o.value === terpasang)) {
    opsi.push({ value: terpasang, label: BG.nilaiLama(kodeNama(terpasang, namaTerpasang)) })
  }
  return opsi
}
