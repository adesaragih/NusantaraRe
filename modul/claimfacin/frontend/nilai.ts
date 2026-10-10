// Disalin dari `modul/claimnonprop/frontend/nilai.ts` (pola, bukan impor; asal Claim Prop): keputusan work
// owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Nilai layar Claim Fac In - fungsi murni atas halaman dan pohon tata (diuji tanpa DOM).
//
// Jalur medan: `ClaimData.X` (halaman), `daftar(n).prop` (sel baris, n berbasis 1; daftar bersarang
// `ClaimData.ObjectList(1).ObjectItemList(2).EstimationList(1).Prop` dipecah pada indeks TERAKHIR). Layar mengirim nilai
// SEMUA jalur terbuka bersama setiap aksi; server menggabungkan hanya yang terbuka menurut tata miliknya sendiri.

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

/** Gabungan semua pohon tata layar (utama, panel baris masterDetail, modal). */
export function semuaTata(utama: readonly Tata[], lain: Record<string, Tata[] | undefined>[]): Tata[] {
  const out = [...utama]
  for (const peta of lain) for (const t of Object.values(peta)) out.push(...(t ?? []))
  return out
}

/**
 * Penanda mode layar (`models.ModeLayar`, `modul/claimfacin/backend/models/layar_dasar.go`): hidup di halaman tetapi
 * tidak disimpan - Edit Catastrophe, lokasi sementara, Search Type / Search Name dan detail polis pop-up Choose Polis,
 * Send Email / Message pop-up Print PLA / DLA, penanda Send to Committe, medan pop-up Comittee / Reject Claim, Close
 * Without Payment.
 */
export const MODE_LAYAR: readonly string[] = [
  'ClaimData.EditCatastrope',
  'TempLocaion.CARI1',
  'SearchType',
  'SearchName',
  'InputParam.CARI4',
  'ProtectPrint.CARI50',
  'Email.Message',
  'Protect.CARI1',
  'Protect.CARI2',
  'TempCommiteClaim.DateOfComitee',
  'TreatyExchangeYearly.UserName',
  'TempCommiteClaim.AllocationShareSalvage',
  'TempCommiteClaim.Initial',
  'TempCommiteClaim.TypeComentAnalysis',
]

/**
 * `mode` setiap aksi: `Layar.mode` terakhir ditimpa nilai lokal penanda mode (diubah layar sesudah layar terakhir
 * diterima). Server memasangnya SEBELUM tata dihitung (`models.PasangMode`, hanya nilai yang sah).
 */
export function modeKirim(mode: Record<string, string> | undefined, h: Halaman): Record<string, string> {
  const out: Record<string, string> = { ...(mode ?? {}) }
  for (const k of MODE_LAYAR) {
    if (Object.hasOwn(h.nilai, k)) out[k] = h.nilai[k] ?? ''
  }
  return out
}

/** Sumber pilihan yang dibaca dari acuan statis (`GET /acuan`) - tidak pernah ditanyakan per sel. */
const SUMBER_ACUAN = new Set(['mataUang', 'jenisReas'])

/** Sumber pilihan yang dibaca dari `GET .../pilihan/{sumber}` (bukan kode `associated`, bukan acuan statis). */
export function sumberServer(sumber: string): boolean {
  return sumber !== '' && !sumber.startsWith('kode:') && !SUMBER_ACUAN.has(sumber)
}

/** Kunci simpanan pilihan server: sumber tingkat item / adjustment bergantung pada panel (konteks) dan baris itemnya. */
export function kunciOpsi(sumber: string, konteks: string, indeks: number): string {
  return `${sumber}|${konteks}|${indeks}`
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
 * Tampilan angka: titik ribuan, koma desimal, paling banyak 4 desimal (dibulatkan setengah ke atas pada digit, bukan
 * float), nol ekor dibuang - inti `formatNumber`. Bukan angka = apa adanya. Nilai tersimpan tidak disentuh.
 */
export function tampilAngka(v: string): string {
  return formatNumber(v, DESIMAL_TAMPIL)
}
