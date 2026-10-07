// Baca CSV Loss Record (tiket 55) - padanan `UploadCSV_DataPolis` (`pxUploadCSVResults`) + langkah 7 `SaveLossRecord`
// (`NB FacIn\Activity\SaveLossRecord.xml`): kolom dicari menurut NAMA kepala (CLASS_OF_BUSINESS, DATE_OF_LOSS,
// CAUSE_OF_LOSS, CURRENCY, LOSS_AMOUNT - urutan di template work owner berbeda dari urutan baca Pega), kolom lain diabaikan
// (template memuat catatan di kolom ke-10). Pemisah `;` (template) atau `,`. Loss Amount format Indonesia
// (`@replaceAll(@replaceAll(.LOSS_AMOUNT,".",""),",",".")`); Cause of Loss huruf besar (`@toUpperCase`); Date of Loss
// DD/MM/YYYY -> kabel DD-MM-YYYY, kosong dibiarkan (backend mengisi Begin date - 1 tahun, langkah 7.6).
// Uw Year dan Loss Amount (IDR) dihitung backend. Tanpa float.

import { desimalSah } from '../../../../inti/frontend/lib/desimal'
import type { BarisLossRecord, SaranAkumulasi } from '../api'
import { TEKS_LOSS_RECORD as T } from '../labels'
import { saranSama } from './IsianSaran'

export const KOLOM_CSV_LOSS_RECORD = ['CLASS_OF_BUSINESS', 'DATE_OF_LOSS', 'CAUSE_OF_LOSS', 'CURRENCY', 'LOSS_AMOUNT'] as const

export interface HasilCsvLossRecord {
  baris: BarisLossRecord[]
  /** Nomor baris berkas (kepala = 1) untuk tiap `baris` - wadah yang sama, indeks sama. */
  nomor: number[]
  galat: string[]
}

/** Pecah satu baris CSV; tanda kutip ganda membungkus sel (kutip ganda di dalamnya ditulis dua kali). */
export function pecahBarisCsv(baris: string, pemisah: string): string[] {
  const sel: string[] = []
  let kini = ''
  let dalamKutip = false
  for (let i = 0; i < baris.length; i++) {
    const c = baris[i]
    if (dalamKutip) {
      if (c === '"' && baris[i + 1] === '"') {
        kini += '"'
        i++
      } else if (c === '"') dalamKutip = false
      else kini += c
    } else if (c === '"') dalamKutip = true
    else if (c === pemisah) {
      sel.push(kini)
      kini = ''
    } else kini += c
  }
  sel.push(kini)
  return sel
}

/** "100.000.000" -> "100000000"; "10.000,50" -> "10000.50"; bukan angka -> null. */
export function angkaIndonesia(teks: string): string | null {
  const t = teks.trim().replace(/\./g, '').replace(/,/g, '.')
  return desimalSah(t) ? t : null
}

/** "03/01/2025" -> "03-01-2025"; kosong -> ""; bukan tanggal yang ada -> null. */
export function tanggalCsv(teks: string): string | null {
  const t = teks.trim()
  if (t === '') return ''
  const m = /^(\d{2})\/(\d{2})\/(\d{4})$/.exec(t)
  if (!m) return null
  const [, dd = '', mm = '', yyyy = ''] = m
  const hari = Number.parseInt(dd, 10)
  const bulan = Number.parseInt(mm, 10)
  const tahun = Number.parseInt(yyyy, 10)
  const d = new Date(Date.UTC(tahun, bulan - 1, hari))
  if (d.getUTCFullYear() !== tahun || d.getUTCMonth() !== bulan - 1 || d.getUTCDate() !== hari) return null
  return `${dd}-${mm}-${yyyy}`
}

export function bacaCsvLossRecord(teks: string): HasilCsvLossRecord {
  const galat: string[] = []
  const baris = teks.replace(/^﻿/, '').split(/\r?\n/)
  const kepala = baris[0] ?? ''
  const pemisah = kepala.includes(';') ? ';' : ','
  const nama = pecahBarisCsv(kepala, pemisah).map((s) => s.trim().toUpperCase())
  const indeks = KOLOM_CSV_LOSS_RECORD.map((k) => nama.indexOf(k))
  const hilang = KOLOM_CSV_LOSS_RECORD.filter((_, i) => indeks[i] === -1)
  if (hilang.length > 0) return { baris: [], nomor: [], galat: [T.csvKolom + hilang.join(', ')] }
  const [iCob = 0, iTgl = 0, iSebab = 0, iMu = 0, iJumlah = 0] = indeks
  const hasil: BarisLossRecord[] = []
  const nomorBaris: number[] = []
  baris.slice(1).forEach((b, n) => {
    const sel = pecahBarisCsv(b, pemisah)
    const ambil = (i: number) => (sel[i] ?? '').trim()
    if ([iCob, iTgl, iSebab, iMu, iJumlah].every((i) => ambil(i) === '')) return
    const nomor = n + 2
    const tanggal = tanggalCsv(ambil(iTgl))
    const jumlah = angkaIndonesia(ambil(iJumlah))
    if (tanggal === null) galat.push(T.csvBaris + nomor + T.csvTanggal)
    if (jumlah === null) galat.push(T.csvBaris + nomor + T.csvJumlah)
    nomorBaris.push(nomor)
    hasil.push({
      uwYear: '',
      classOfBusiness: ambil(iCob),
      dateOfLoss: tanggal ?? '',
      causeOfLoss: ambil(iSebab).toUpperCase(),
      currency: ambil(iMu).toUpperCase(),
      lossAmount: jumlah ?? '',
      lossAmountIdr: '',
    })
  })
  if (galat.length === 0 && hasil.length === 0) galat.push(T.csvKosong)
  return galat.length > 0 ? { baris: [], nomor: [], galat } : { baris: hasil, nomor: nomorBaris, galat }
}

/**
 * Class Of Business CSV wajib dari daftar `BUSINESS.NOTE` (work owner 06-10-2026, menyusul aturan ketikan L55-7): tiap nilai
 * berbeda dicari sekali lewat `cari` (saran Class Of Business); sama persis tanpa beda huruf besar/kecil -> diganti ejaan
 * daftar; kosong / tidak ada -> galat per baris, dan seluruh hasil ditolak (grid tidak diganti, L55-3).
 */
export async function cocokkanKelasBisnis(
  h: HasilCsvLossRecord,
  cari: (q: string) => Promise<{ baris: SaranAkumulasi[] }>,
): Promise<HasilCsvLossRecord> {
  if (h.galat.length > 0) return h
  const unik = [...new Set(h.baris.map((b) => b.classOfBusiness.trim().toUpperCase()).filter((x) => x !== ''))]
  const peta = new Map<string, string>()
  await Promise.all(
    unik.map(async (k) => {
      const s = saranSama(k, (await cari(k)).baris)
      if (s) peta.set(k, s.label)
    }),
  )
  const galat: string[] = []
  const baris = h.baris.map((b, i) => {
    const k = b.classOfBusiness.trim().toUpperCase()
    const nomor = h.nomor[i] ?? i + 2
    if (k === '') galat.push(T.csvBaris + nomor + T.csvKelasKosong)
    else if (!peta.has(k)) galat.push(T.csvBaris + nomor + T.csvKelasTidakAda + b.classOfBusiness.trim())
    return { ...b, classOfBusiness: peta.get(k) ?? b.classOfBusiness }
  })
  return galat.length > 0 ? { baris: [], nomor: [], galat } : { baris, nomor: h.nomor, galat }
}
