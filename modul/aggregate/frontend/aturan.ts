// Aturan layar Aggregate - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend (`backend/services`).

import { jumlahDesimal } from '../../../inti/frontend/lib/desimal'

import type { Baris, IrisanRingkasan, Kelompok, Kunci } from './api'

/**
 * Satu kolom grid TempCSV - cermin `models.KolomGrid` (urutan, angka; `aturan.test.ts` mencocokkan). Grid hanya
 * dibaca (work owner 04-10-2026: "upload csv nya read only"), jadi tanda kolom yang dapat diubah di Pega tidak dibawa.
 */
export interface KolomGrid {
  nama: string
  angka: boolean
}

const teks = (nama: string): KolomGrid => ({ nama, angka: false })
const angka = (nama: string): KolomGrid => ({ nama, angka: true })

export const KOLOM_GRID: readonly KolomGrid[] = [
  teks('ASSESMENT_ZONE'),
  teks('M_TREATY_ID'),
  teks('TREATY_TYPE'),
  teks('COVERAGE'),
  teks('AS_AT'),
  teks('UW_YEAR'),
  teks('TREATYYEAR'),
  teks('CEDING_CODE'),
  teks('CEDING_NAME'),
  teks('CURRENCY'),
  angka('TO_USD'),
  angka('NOR_BUILDINGS'),
  angka('BUILDINGS'),
  angka('NOR_STOCKS'),
  angka('STOCKS'),
  angka('NOR_MACHINERY'),
  angka('MACHINERY'),
  angka('NOR_OTHER_CONTENTS'),
  angka('OTHER_CONTENTS'),
  angka('NOR_CONSEQUENTIAL_LOSS'),
  angka('CONSEQUENTIAL_LOSS'),
  angka('NOR_RESIDENTIAL'),
  angka('RESIDENTIAL'),
  angka('NOR_COMMERCIAL'),
  angka('COMMERCIAL'),
  angka('NOR_INDUSTRIAL'),
  angka('INDUSTRIAL'),
  angka('NOR_AGRICULTURE'),
  angka('AGRICULTURE'),
  angka('NOR_MISCELLANEOUS'),
  angka('MISCELLANEOUS'),
  angka('NOR_UTILITIES'),
  angka('UTILITIES'),
  angka('TOTAL_NO_OF_RISK'),
  angka('TOTAL_IN_AMOUNT'),
  angka('TOTAL_IN_AMOUNT_IN_USD'),
  angka('RNM_SHARE'),
  angka('RNM_VALUE'),
  angka('RNM_VALUE_IN_USD'),
  teks('REMARK'),
]

/** Satu sel kepala grid. */
export interface SelKepala {
  kunci: string
  teks: string
  /** Caption Pega lengkap (title). */
  judul: string
  colSpan: number
  rowSpan: number
  /** Sel judul grup (baris pertama, di atas kolom-kolomnya). */
  grup: boolean
  angka: boolean
}

/**
 * Susun kepala grid dua baris seperti grid Bordereaux (perintah work owner 05-10-2026): kolom bergrup bersebelahan
 * digabung di baris pertama, judul bawahnya di baris kedua; kolom tanpa grup menempati dua baris. `grup` dan `label`
 * diberikan pemanggil (`GRUP_KOLOM`, `LABEL_KOLOM`).
 */
export function susunKepala(
  kolom: readonly KolomGrid[],
  grup: Readonly<Record<string, readonly [string, string]>>,
  label: Readonly<Record<string, string>>,
): { baris1: SelKepala[]; baris2: SelKepala[] } {
  const baris1: SelKepala[] = []
  const baris2: SelKepala[] = []
  for (const k of kolom) {
    const judul = label[k.nama] ?? k.nama
    const g = grup[k.nama]
    if (g === undefined) {
      baris1.push({ kunci: k.nama, teks: judul, judul, colSpan: 1, rowSpan: 2, grup: false, angka: k.angka })
      continue
    }
    baris2.push({ kunci: k.nama, teks: g[1], judul, colSpan: 1, rowSpan: 1, grup: false, angka: k.angka })
    const akhir = baris1[baris1.length - 1]
    if (akhir !== undefined && akhir.grup && akhir.teks === g[0]) akhir.colSpan++
    else baris1.push({ kunci: `grup:${k.nama}`, teks: g[0], judul: g[0], colSpan: 1, rowSpan: 1, grup: true, angka: false })
  }
  return { baris1, baris2 }
}

/** Assessment Zone baris jumlah per mata uang - tidak disimpan. */
export const ZONA_TOTAL = 'Total :'
/**
 * Templat Upload CSV dikelola Template Manager (keputusan work owner 04-10-2026): kode slotnya = `KodeTemplat` di
 * `backend/modul.go`; berkas bawaannya `backend/templat/aggregate.csv` (header 38 kolom, pemisah `;`).
 */
export const KODE_TEMPLAT = 'aggregate.upload'

/** Nama berkas unduhan - sama dengan `NamaUnduhan` slot di backend. */
export const NAMA_TEMPLATE = 'aggregate.csv'

export function barisTotal(b: Baris): boolean {
  return b['ASSESMENT_ZONE'] === ZONA_TOTAL
}

/** Pemisah ribuan format Indonesia (titik) untuk deret angka bulat. */
function ribuan(bulat: string): string {
  return bulat.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
}

/** Tambah satu pada deret angka bulat tanpa tanda (`999` -> `1000`), tanpa float. */
function tambahSatu(bulat: string): string {
  const d = bulat.split('')
  let i = d.length - 1
  while (i >= 0 && d[i] === '9') d[i--] = '0'
  if (i < 0) return `1${d.join('')}`
  d[i] = String(Number(d[i]) + 1)
  return d.join('')
}

/**
 * Angka kabel (`1234567.5`) -> tampilan format Indonesia dua desimal (`1.234.567,50`), dibulatkan setengah ke atas
 * (work owner 04-10-2026: "format indo dan 2 belakang koma untuk tampilannya"). Hanya tampilan - nilai yang dikirim
 * dan disimpan tetap teks backend utuh. Dihitung pada teks, bukan float, supaya `2166.5286000000001408` tidak
 * bergeser. Kosong tetap kosong; teks bukan angka dikembalikan apa adanya.
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

/** Bilangan bulat (jumlah baris) format Indonesia tanpa desimal: `13450` -> `13.450`. */
export function formatBulat(n: number): string {
  const bulat = String(Math.trunc(Math.abs(n)))
  return `${n < 0 ? '-' : ''}${ribuan(bulat)}`
}

/** Persen chart format Indonesia dua desimal: `75` -> `75,00%`. */
export function formatPersen(persen: number): string {
  return `${formatAngka(persen.toFixed(4))}%`
}

/** Kunci baris daftar. */
export function kunciDari(k: Kelompok): Kunci {
  return {
    tanggalInput: k.tanggalInput,
    cedingCode: k.cedingCode,
    cedingName: k.cedingName,
    treatyType: k.treatyType,
    asAt: k.asAt,
    uwYear: k.uwYear,
  }
}

/** Kunci teks satu baris daftar (kunci React). */
export function kunciTeks(k: Kunci): string {
  return [k.tanggalInput, k.cedingCode, k.cedingName, k.treatyType, k.asAt, k.uwYear].join('|')
}

/** Jumlah halaman daftar (sekurangnya 1). */
export function jumlahHalaman(total: number, ukuran: number): number {
  return Math.max(1, Math.ceil(total / Math.max(1, ukuran)))
}

/** Pesan galat backend -> baris layar (Save memisah pesan Pega dengan baris baru). */
export function barisPesan(pesan: string): string[] {
  return pesan
    .split('\n')
    .map((s) => s.trim())
    .filter((s) => s !== '')
}


/** Tingkat chart: Ceding, dibuka ke Treaty Type, dibuka ke Coverage (work owner 04-10-2026). */
export type TingkatChart = 'ceding' | 'treatyType' | 'coverage'

/** Jalur buka chart: kosong = seluruh ceding; `ceding` = CEDING_CODE; `treatyType` sesudah ceding. */
export interface JalurChart {
  ceding?: string
  treatyType?: string
}

/** Satu ruas batang bertumpuk. `seri` = indeks warna; `lebar` = persen lebar lintasan. */
export interface RuasChart {
  kunci: string
  nilai: string
  seri: number
  lebar: number
}

/** Satu batang: `lebar` relatif terhadap batang terbesar, `persen` terhadap total tingkat ini. */
export interface BatangChart {
  kunci: string
  label: string
  nilai: string
  persen: number
  lebar: number
  bagian: RuasChart[]
}

export interface TampilanChart {
  tingkat: TingkatChart
  total: string
  /** Nama ceding pada jalur (tingkat Treaty Type dan Coverage). */
  namaCeding: string
  /** Seri tumpukan berurutan (legenda); kosong di tingkat Coverage. */
  seri: string[]
  /** Batang dapat dibuka ke tingkat berikutnya. */
  bisaBuka: boolean
  batang: BatangChart[]
}

/** Nilai kabel -> angka untuk lebar dan urutan saja; nilai tampil tetap teks eksak. */
function angkaLebar(nilai: string): number {
  const n = Number(nilai)
  return Number.isFinite(n) && n > 0 ? n : 0
}

/**
 * Isi chart satu tingkat dari daun ringkasan backend (Ceding x Treaty Type x Coverage). Jumlah dihitung eksak
 * (`jumlahDesimal`, BigInt); batang terbesar dulu. Tingkat Ceding ditumpuk per Treaty Type, tingkat Treaty Type per
 * Coverage; warna Coverage di tingkat terakhir sama dengan warnanya di tumpukan tingkat Treaty Type.
 */
export function tampilanRingkasan(irisan: readonly IrisanRingkasan[], jalur: JalurChart): TampilanChart {
  const tingkat: TingkatChart = jalur.ceding === undefined ? 'ceding' : jalur.treatyType === undefined ? 'treatyType' : 'coverage'
  const milikCeding = jalur.ceding === undefined ? irisan : irisan.filter((i) => i.cedingCode === jalur.ceding)
  const daun = tingkat === 'coverage' ? milikCeding.filter((i) => i.treatyType === jalur.treatyType) : milikCeding
  const kunciBatang = (i: IrisanRingkasan) => (tingkat === 'ceding' ? i.cedingCode : tingkat === 'treatyType' ? i.treatyType : i.coverage)
  const kunciRuas = (i: IrisanRingkasan) => (tingkat === 'ceding' ? i.treatyType : i.coverage)
  const urutanWarna = [...new Set(milikCeding.map(kunciRuas))].sort()

  const kelompok = new Map<string, { label: string; ruas: Map<string, string[]> }>()
  for (const i of daun) {
    const k = kunciBatang(i)
    const g = kelompok.get(k) ?? { label: tingkat === 'ceding' ? i.cedingName || i.cedingCode : k, ruas: new Map<string, string[]>() }
    const r = tingkat === 'coverage' ? k : kunciRuas(i)
    g.ruas.set(r, [...(g.ruas.get(r) ?? []), i.rnmValueInUsd])
    kelompok.set(k, g)
  }

  const total = jumlahDesimal(daun.map((i) => i.rnmValueInUsd)).total
  const mentah = [...kelompok].map(([kunci, g]) => {
    const bagian = [...g.ruas]
      .map(([r, nilai]) => ({ kunci: r, nilai: jumlahDesimal(nilai).total, seri: Math.max(0, urutanWarna.indexOf(r)) }))
      .sort((a, b) => a.seri - b.seri)
    return { kunci, label: g.label, nilai: jumlahDesimal(bagian.map((b) => b.nilai)).total, bagian }
  })
  mentah.sort((a, b) => angkaLebar(b.nilai) - angkaLebar(a.nilai) || a.label.localeCompare(b.label))
  const terbesar = mentah.reduce((m, b) => Math.max(m, angkaLebar(b.nilai)), 0)
  const totalAngka = angkaLebar(total)
  const persenDari = (n: number, dari: number) => (dari > 0 ? (n / dari) * 100 : 0)

  return {
    tingkat,
    total,
    namaCeding: jalur.ceding === undefined ? '' : (milikCeding[0]?.cedingName || jalur.ceding),
    seri: tingkat === 'coverage' ? [] : urutanWarna,
    bisaBuka: tingkat !== 'coverage',
    batang: mentah.map((b) => ({
      kunci: b.kunci,
      label: b.label,
      nilai: b.nilai,
      persen: persenDari(angkaLebar(b.nilai), totalAngka),
      lebar: persenDari(angkaLebar(b.nilai), terbesar),
      bagian: b.bagian.map((s) => ({ ...s, lebar: persenDari(angkaLebar(s.nilai), terbesar) })),
    })),
  }
}
