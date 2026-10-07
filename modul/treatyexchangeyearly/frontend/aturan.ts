// Aturan layar Treaty Exchange Yearly - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend
// (`backend/services`): format tanggal Pega, kembar tahun + mata uang + quarter, hak View only.

import type { Isian, Kurs, MataUang } from './api'
import { TEY } from './labels'

/** Quarter yang ditawarkan: 0 tahunan, 1-4 triwulan (backend `models.Quarter`). */
export const QUARTER = ['0', '1', '2', '3', '4'] as const

const ANGKA = /^[0-9]+(\.[0-9]+)?$/

/** Teks tanggal Pega (`YYYYMMDDT...`, juga bentuk warisan rusak) ke isian `YYYY-MM-DD`; selainnya kosong. */
export function keInputTanggal(pega: string): string {
  const m = /^([0-9]{4})([0-9]{2})([0-9]{2})/.exec(pega.trim())
  return m === null ? '' : `${m[1]}-${m[2]}-${m[3]}`
}

/** Periode bawaan tahun treaty: 1 Juli tahun itu - 30 Juni tahun berikutnya (pola data sejak 2020). */
export function periodeBawaan(tahun: string): { startDate: string; endDate: string } | null {
  if (!/^[0-9]{4}$/.test(tahun.trim())) return null
  const t = Number(tahun.trim())
  return { startDate: `${t}-07-01`, endDate: `${t + 1}-06-30` }
}

/** Isian form dari satu baris (Edit) atau kosong (Add, Quarter 0). */
export function isianDari(k: Kurs | null): Isian {
  if (k === null)
    return { kunci: '', treatyYear: '', idCurrency: '', startDate: '', endDate: '', toIdr: '', toUsd: '', quarter: '0' }
  return {
    kunci: k.kunci,
    treatyYear: k.treatyYear,
    idCurrency: k.idCurrency,
    startDate: keInputTanggal(k.startDate),
    endDate: keInputTanggal(k.endDate),
    toIdr: k.toIdr,
    toUsd: k.toUsd,
    quarter: k.quarter === '' ? '0' : k.quarter,
  }
}

/** Ganti Treaty Year: Add yang tanggalnya masih kosong atau masih periode bawaan tahun lama ikut berganti periode. */
export function gantiTahun(i: Isian, tahun: string, add: boolean): Isian {
  const lama = periodeBawaan(i.treatyYear)
  const baru = periodeBawaan(tahun)
  const masihBawaan =
    (i.startDate === '' && i.endDate === '') ||
    (lama !== null && i.startDate === lama.startDate && i.endDate === lama.endDate)
  if (add && baru !== null && masihBawaan) return { ...i, treatyYear: tahun, ...baru }
  return { ...i, treatyYear: tahun }
}

/** Isian yang dikirim: spasi tepi dibuang. */
export function rapikanIsian(i: Isian): Isian {
  return {
    kunci: i.kunci,
    treatyYear: i.treatyYear.trim(),
    idCurrency: i.idCurrency.trim(),
    startDate: i.startDate.trim(),
    endDate: i.endDate.trim(),
    toIdr: i.toIdr.trim(),
    toUsd: i.toUsd.trim(),
    quarter: i.quarter.trim() === '' ? '0' : i.quarter.trim(),
  }
}

/** Pemeriksaan awal sebelum dikirim; `null` = boleh dikirim. Backend tetap memeriksa seluruhnya. */
export function periksaIsian(i: Isian): string | null {
  if (!/^[0-9]{4}$/.test(i.treatyYear.trim())) return TEY.galatTahun
  if (i.idCurrency.trim() === '') return TEY.galatMataUang
  if (i.startDate === '' || i.endDate === '') return TEY.galatTanggal
  if (i.endDate < i.startDate) return TEY.galatUrutan
  if (!ANGKA.test(i.toIdr.trim())) return TEY.galatIdr
  if (i.toUsd.trim() !== '' && !ANGKA.test(i.toUsd.trim())) return TEY.galatUsd
  return null
}

/**
 * Kurs untuk dibaca, format Indonesia (perintah work owner 05-10-2026: "perbaiki tampilan uang", "format angkanya
 * versi indo lah"): ribuan titik, desimal koma, minimal dua desimal, desimal tersimpan yang lebih panjang tidak dipotong (To USD sampai 7 desimal di DEV, mis. `0.0067809`).
 * Teks yang bukan angka bertitik apa adanya; yang tersimpan tidak berubah.
 */
export function tampilKurs(s: string): string {
  const m = /^([0-9]+)(?:\.([0-9]+))?$/.exec(s.trim())
  if (m === null) return s
  const bulat = m[1]!.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
  return `${bulat},${(m[2] ?? '').padEnd(2, '0')}`
}

/** Kode dan nama: `USD - UNITED STATE DOLLAR`; nama kosong = kode saja. */
export function kodeNama(kode: string, nama: string): string {
  if (kode === '') return ''
  return nama === '' ? kode : `${kode} - ${nama}`
}

/** Opsi Currency; mata uang tersimpan yang tidak ada lagi di view CURRENCY tetap tampil sebagai nilai lama. */
export function opsiMataUang(
  mataUang: MataUang[],
  terpasang: string,
  kodeTerpasang = '',
): { value: string; label: string }[] {
  const opsi = mataUang.map((m) => ({ value: m.id, label: kodeNama(m.kode, m.nama) }))
  if (terpasang !== '' && !opsi.some((o) => o.value === terpasang)) {
    opsi.push({ value: terpasang, label: TEY.nilaiLama(kodeTerpasang === '' ? terpasang : kodeTerpasang) })
  }
  return opsi
}
