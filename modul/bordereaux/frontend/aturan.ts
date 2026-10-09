// Aturan layar Bordereaux - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/services`).

import type { BerkasLama, Filter, HasilSalinLama, IrisanChart, Kolom, StatusSalinLama } from './api'
import { BDX } from './labels'

export const FILTER_KOSONG: Filter = {
  id: '',
  type: '',
  business: '',
  reffSoa: '',
  reffBdx: '',
  ceding: '',
  treaty: '',
  start: '',
  end: '',
  position: '',
  status: '',
}

/** Kunci kombinasi `TYPE|BUSINESS` (kunci `Pilihan.kolom` dan `Pilihan.templat`). */
export function kunciKombinasi(type: string, business: string): string {
  return `${type}|${business}`
}

/** SUBROGATION hanya BONDING (`SetSubrogation_ACT`). */
export function businessUntuk(type: string, business: string): string {
  return type === 'SUBROGATION' ? 'BONDING' : business
}

/** Pemisah ribuan format Indonesia. */
function ribuan(bulat: string): string {
  return bulat.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
}

function tambahSatu(bulat: string): string {
  const d = bulat.split('')
  let i = d.length - 1
  while (i >= 0 && d[i] === '9') d[i--] = '0'
  if (i < 0) return `1${d.join('')}`
  d[i] = String(Number(d[i]) + 1)
  return d.join('')
}

/**
 * Angka kabel (`1234567.5`) -> format Indonesia dua desimal (`1.234.567,50`), dibulatkan setengah ke atas pada teks
 * (tanpa float). Kosong tetap kosong; bukan angka dikembalikan apa adanya. Sama dengan tampilan Aggregate.
 */
export function formatAngka(s: string | undefined): string {
  if (s === undefined || s === '') return ''
  const m = /^(-?)(\d*)(?:\.(\d*))?$/.exec(s)
  if (m === null || (m[2] === '' && (m[3] ?? '') === '')) return s
  let bulat = (m[2] === '' ? '0' : m[2]!).replace(/^0+(?=\d)/, '')
  const pecahan = (m[3] ?? '').padEnd(3, '0')
  let sen = Number(pecahan.slice(0, 2))
  if (pecahan[2]! >= '5') sen += 1
  if (sen === 100) {
    sen = 0
    bulat = tambahSatu(bulat)
  }
  const tanda = m[1] === '-' && (bulat !== '0' || sen !== 0) ? '-' : ''
  return `${tanda}${ribuan(bulat)},${String(sen).padStart(2, '0')}`
}

/** Tanggal input HTML (`YYYY-MM-DD`) <-> kabel (`DD-MM-YYYY`). */
export function keKabel(iso: string): string {
  const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(iso)
  return m === null ? '' : `${m[3]}-${m[2]}-${m[1]}`
}

export function keIso(kabel: string): string {
  const m = /^(\d{2})-(\d{2})-(\d{4})$/.exec(kabel)
  return m === null ? '' : `${m[3]}-${m[2]}-${m[1]}`
}

/** Satu sel kepala grid detail. */
export interface SelKepala {
  kunci: string
  teks: string
  colSpan: number
  rowSpan: number
  /** Kolom angka tanpa grup (judulnya menempati dua baris). */
  angka: boolean
}

/** Kepala grid detail: kolom yang tampil, baris kepala pertama, dan baris kedua (kolom bergrup). */
export interface SusunanKepala {
  tampil: Kolom[]
  baris1: SelKepala[]
  baris2: Kolom[]
  duaBaris: boolean
  /** Kombinasi memakai kepala format Excel bordereaux 2025. */
  excel: boolean
}

/** Judul kepala satu kolom: menurut format Excel bila ada, selain itu judul CSV. */
export function judulKepala(k: Kolom): string {
  return k.kepala !== undefined && k.kepala !== '' ? k.kepala : k.judul
}

/**
 * Susun kepala grid detail seperti sheet format bordereaux 2025 (perintah work owner 05-10-2026): kolom bergrup
 * bersebelahan digabung di baris pertama (colSpan), judulnya di baris kedua; kolom tanpa grup menempati dua baris
 * (rowSpan 2). Kolom `sembunyi` tidak tampil. Tanpa grup sama sekali = satu baris kepala.
 */
export function susunKepala(kolom: readonly Kolom[]): SusunanKepala {
  const tampil = kolom.filter((k) => k.sembunyi !== true)
  const duaBaris = tampil.some((k) => (k.grup ?? '') !== '')
  const excel = kolom.some((k) => (k.kepala ?? '') !== '')
  const baris1: SelKepala[] = []
  const baris2: Kolom[] = []
  for (const k of tampil) {
    const grup = k.grup ?? ''
    if (grup === '') {
      baris1.push({ kunci: k.kolom, teks: judulKepala(k), colSpan: 1, rowSpan: duaBaris ? 2 : 1, angka: k.jenis === 'angka' })
      continue
    }
    baris2.push(k)
    const akhir = baris1[baris1.length - 1]
    if (akhir !== undefined && akhir.rowSpan === 1 && duaBaris && akhir.teks === grup && akhir.kunci.startsWith('grup:')) {
      akhir.colSpan++
    } else {
      baris1.push({ kunci: `grup:${k.kolom}`, teks: grup, colSpan: 1, rowSpan: 1, angka: false })
    }
  }
  return { tampil, baris1, baris2, duaBaris, excel }
}

/** Jumlah halaman (sekurangnya 1). */
export function jumlahHalaman(total: number, ukuran: number): number {
  return Math.max(1, Math.ceil(total / Math.max(1, ukuran)))
}

/** Pesan galat backend -> baris layar (pesan validasi dipisah baris baru). */
export function barisPesan(pesan: string): string[] {
  return pesan
    .split('\n')
    .map((s) => s.trim())
    .filter((s) => s !== '')
}

/** Urutan seri (warna) chart: Type Pega. */
export const SERI_TYPE = ['PREMIUM', 'CLAIM', 'SUBROGATION'] as const

/** Jalur buka chart: kosong = semua Business; `business` = Ceding satu Business. */
export interface JalurChart {
  business?: string
}

export interface RuasChart {
  kunci: string
  jumlah: number
  seri: number
  lebar: number
}

export interface BatangChart {
  kunci: string
  label: string
  jumlah: number
  persen: number
  lebar: number
  bagian: RuasChart[]
}

export interface TampilanChart {
  tingkat: 'business' | 'ceding'
  total: number
  batang: BatangChart[]
  bisaBuka: boolean
}

/**
 * Isi chart daftar (keputusan work owner 04-10-2026: "group by bisnis type, lalu diklik membuka ceding nya"): jumlah
 * berkas per Business, ditumpuk per Type; dibuka = per Ceding Business itu. Terbanyak dulu.
 */
export function tampilanChart(irisan: readonly IrisanChart[], jalur: JalurChart): TampilanChart {
  const tingkat = jalur.business === undefined ? 'business' : 'ceding'
  const daun = tingkat === 'business' ? irisan : irisan.filter((i) => i.business === jalur.business)
  const per = new Map<string, { label: string; ruas: Map<string, number> }>()
  for (const i of daun) {
    const kunci = tingkat === 'business' ? i.business : i.cedingId || i.cedingName
    const label = tingkat === 'business' ? i.business : i.cedingName || i.cedingId
    const g = per.get(kunci) ?? { label, ruas: new Map<string, number>() }
    g.ruas.set(i.type, (g.ruas.get(i.type) ?? 0) + i.jumlah)
    per.set(kunci, g)
  }
  const total = daun.reduce((n, i) => n + i.jumlah, 0)
  const mentah = [...per].map(([kunci, g]) => ({
    kunci,
    label: g.label,
    jumlah: [...g.ruas.values()].reduce((a, b) => a + b, 0),
    ruas: g.ruas,
  }))
  mentah.sort((a, b) => b.jumlah - a.jumlah || a.label.localeCompare(b.label))
  const terbesar = mentah.reduce((m, b) => Math.max(m, b.jumlah), 0)
  const persen = (n: number, dari: number) => (dari > 0 ? (n / dari) * 100 : 0)
  const urutSeri = (t: string) => {
    const i = (SERI_TYPE as readonly string[]).indexOf(t)
    return i < 0 ? SERI_TYPE.length : i
  }
  return {
    tingkat,
    total,
    bisaBuka: tingkat === 'business',
    batang: mentah.map((b) => ({
      kunci: b.kunci,
      label: b.label,
      jumlah: b.jumlah,
      persen: persen(b.jumlah, total),
      lebar: persen(b.jumlah, terbesar),
      bagian: [...b.ruas]
        .map(([t, n]) => ({ kunci: t, jumlah: n, seri: urutSeri(t), lebar: persen(n, terbesar) }))
        .sort((x, y) => x.seri - y.seri),
    })),
  }
}

/** Persen untuk dibaca: `75,0%`. */
export function teksPersen(p: number): string {
  return `${p.toFixed(1).replace('.', ',')}%`
}

/** Saring popup Copy Old Data: ID, ceding, treaty, Type, atau Business memuat kata (tanpa beda huruf besar). */
export function saringLama(daftar: BerkasLama[], kata: string): BerkasLama[] {
  const k = kata.trim().toUpperCase()
  if (k === '') return daftar
  return daftar.filter((b) => [b.id, b.ceding, b.treaty, b.type, b.business].some((s) => s.toUpperCase().includes(k)))
}

/** Catatan satu baris Copy Old Data: apa yang akan disalin (sama dengan aturan backend). */
export function catatanLama(b: BerkasLama): string {
  const bagian: string[] = []
  if (b.tanpaHeader) bagian.push(BDX.lamaHeader)
  if (b.barisJson > 0 && b.barisTabel === 0) bagian.push(BDX.lamaDetail(b.barisJson))
  if (b.komentar > 0 && b.riwayat === 0) bagian.push(BDX.lamaRiwayat(b.komentar))
  return bagian.length === 0 ? BDX.lamaTakTerbaca : bagian.join(' · ')
}

/** Jumlah hasil Process Copy per status. */
export function ringkasSalin(hasil: HasilSalinLama[]): Record<StatusSalinLama, number> {
  const r: Record<StatusSalinLama, number> = { disalin: 0, sudahAda: 0, ditolak: 0, gagal: 0 }
  for (const h of hasil) r[h.status]++
  return r
}

/** Ukuran satu permintaan Process Copy: satu berkas bisa ratusan baris detail, jadi kelompoknya kecil. */
export const UKURAN_SALIN_LAMA = 10

/** Pecah ID menjadi kelompok berurutan berukuran `ukuran`. */
export function potong<T>(daftar: T[], ukuran: number): T[][] {
  const n = Math.max(1, ukuran)
  const out: T[][] = []
  for (let i = 0; i < daftar.length; i += n) out.push(daftar.slice(i, i + n))
  return out
}

/**
 * Upload File lampiran: berkas dari pilihan dan seret-lepas DIGABUNG ke daftar terpilih; berkas yang sama (nama tanpa
 * beda huruf besar dan ukuran) tidak digandakan, urutan pilihan lama dipertahankan.
 */
export function gabungBerkas(lama: readonly File[], baru: readonly File[]): File[] {
  const kunci = (f: File) => `${f.name.toLowerCase()}|${f.size}`
  const ada = new Set(lama.map(kunci))
  const hasil = [...lama]
  for (const f of baru) {
    if (ada.has(kunci(f))) continue
    ada.add(kunci(f))
    hasil.push(f)
  }
  return hasil
}
