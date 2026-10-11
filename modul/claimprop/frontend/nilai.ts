// Nilai layar Claim Prop - fungsi murni atas halaman dan pohon tata (diuji tanpa DOM).
//
// Jalur medan: `ClaimData.X` (halaman), `daftar(n).prop` (sel baris, n berbasis 1). Layar mengirim nilai SEMUA jalur
// terbuka bersama setiap aksi; server menggabungkan hanya yang terbuka menurut tata miliknya sendiri.

import { formatNumber } from '../../../inti/frontend/lib/format'
import type { Baris, Halaman, Tata } from './api'

/** Pecah `daftar(n).prop` (indeks terakhir); null bila bukan jalur baris. */
export function pecahJalur(j: string): { daftar: string; n: number; prop: string } | null {
  const i = j.lastIndexOf(').')
  if (i < 0) return null
  const awal = j.slice(0, i)
  const b = awal.lastIndexOf('(')
  if (b < 0) return null
  const n = Number(awal.slice(b + 1))
  if (!Number.isInteger(n)) return null
  return { daftar: awal.slice(0, b), n, prop: j.slice(i + 2) }
}

export function jalurBaris(daftar: string, n: number, prop: string): string {
  return `${daftar}(${n}).${prop}`
}

/** Nilai satu jalur. */
export function ambil(h: Halaman, j: string): string {
  const p = pecahJalur(j)
  if (p) return h.daftar[p.daftar]?.[p.n - 1]?.[p.prop] ?? ''
  return h.nilai[j] ?? ''
}

/** Halaman baru dengan satu jalur diubah (imutabel). */
export function setel(h: Halaman, j: string, v: string): Halaman {
  const p = pecahJalur(j)
  if (!p) return { ...h, nilai: { ...h.nilai, [j]: v } }
  const rows = [...(h.daftar[p.daftar] ?? [])]
  const lama: Baris = rows[p.n - 1] ?? {}
  rows[p.n - 1] = { ...lama, [p.prop]: v }
  return { ...h, daftar: { ...h.daftar, [p.daftar]: rows } }
}

/** Jalur medan terbuka (tampil, tidak hanya-baca, tidak nonaktif, bukan tampilan) - padanan `models.MedanTerbuka`. */
export function jalurTerbuka(ts: readonly Tata[]): string[] {
  const out: string[] = []
  const jalan = (xs: readonly Tata[]) => {
    for (const t of xs) {
      if (t.jenis === 'medan' && !t.hanyaBaca && !t.nonaktif && t.kendali !== 'tampil' && t.jalur) out.push(t.jalur)
      if (t.jenis === 'bagian') jalan(t.anak ?? [])
      if (t.jenis === 'grid') {
        ;(t.baris ?? []).forEach((sel, i) => {
          sel.forEach((s, k) => {
            const kol = t.kolom?.[k]
            if (
              kol &&
              kol.jenis === 'medan' &&
              s.tampil &&
              !s.hanyaBaca &&
              !s.nonaktif &&
              kol.kendali !== 'tampil' &&
              kol.jalur
            ) {
              out.push(jalurBaris(t.jalur ?? '', i + 1, kol.jalur))
            }
          })
        })
        jalan(t.kaki ?? [])
      }
    }
  }
  jalan(ts)
  return out
}

/** Masukan aksi: nilai setiap jalur terbuka. */
export function masukan(h: Halaman, ts: readonly Tata[]): Record<string, string> {
  const out: Record<string, string> = {}
  for (const j of jalurTerbuka(ts)) out[j] = ambil(h, j)
  return out
}

/** Gabungan semua pohon tata layar (utama, baris adjustment, modal). */
export function semuaTata(utama: readonly Tata[], lain: Record<string, Tata[] | undefined>[]): Tata[] {
  const out = [...utama]
  for (const peta of lain) for (const t of Object.values(peta)) out.push(...(t ?? []))
  return out
}

/** Tanggal-waktu halaman ("2006-01-02 15:04:05") <-> nilai `datetime-local` ("2006-01-02T15:04"). */
export function keInputWaktu(v: string): string {
  if (v.length < 16) return v.replace(' ', 'T')
  return v.slice(0, 16).replace(' ', 'T')
}

export function dariInputWaktu(v: string): string {
  if (v === '') return ''
  const s = v.replace('T', ' ')
  return s.length === 16 ? `${s}:00` : s
}

/** Nomor telepon: hanya angka (keputusan work owner 08-10-2026); nol di depan dipertahankan. */
export function hanyaAngka(v: string): string {
  return v.replace(/\D/g, '')
}

/** Desimal tampilan angka (work owner 08-10-2026: "separator indonesia, 4 angka belakang koma"). */
export const DESIMAL_TAMPIL = 4

/**
 * Tampilan angka: titik ribuan, koma desimal, SELALU 4 desimal (work owner 10-10-2026: "untuk tampilannya munculkan 4
 * angka belakang koma" - nol ekor tidak lagi dibuang), dibulatkan setengah ke atas pada digit, bukan float (inti
 * `formatNumber`). Bukan angka = apa adanya. Nilai tersimpan tidak disentuh.
 */
export function tampilAngka(v: string): string {
  const teks = formatNumber(v, DESIMAL_TAMPIL)
  if (!/^[+-]?(\d+\.?\d*|\.\d+)$/.test((v ?? '').trim())) return teks // bukan angka / kosong: apa adanya
  const [bulat, pecahan = ''] = teks.split(',')
  return `${bulat},${pecahan.padEnd(DESIMAL_TAMPIL, '0')}`
}
