// Baris larik layar Adjustment — DATAR untuk grid, BERSARANG untuk rumus.
//
// ⭐ Satu sumber kebenaran: pohon dari backend (`SisiPenyesuaian.pohon`)
// DIGABUNG ke larik saat panel lahir, jadi baris `Limits[]`/`Share[]` membawa
// larik anaknya sendiri (`MDPList`, `TreatyGroupList`, `DeductionList` …).
// Tambah, hapus, dan sunting baris dengan begitu tidak dapat memisahkan
// sebuah baris dari anak-anaknya.

import type { BarisBersarang, NilaiBaris, SisiPenyesuaian } from '../api'

/** Nilai TEKS satu sel — larik anak bukan nilai sel. */
export function teks(v: NilaiBaris | undefined): string {
  return typeof v === 'string' ? v : ''
}

/** Larik anak satu baris, atau larik kosong. */
export function anak(b: BarisBersarang | undefined, kunci: string): BarisBersarang[] {
  const v = b?.[kunci]
  return Array.isArray(v) ? v : []
}

/**
 * Padanan larik TAMPIL → larik sumbernya, dari `TreatyInNonAddItem` langkah
 * 14.3 (dan 12.3 untuk fakultatif): `.RnmLimitListDisplay = .RnmLimitList`,
 * `.RnmGrossPremiDisplay = .GrossPremiumList`. Larik tampil tidak didaratkan;
 * nilainya SAMA dengan sumbernya.
 */
const SUMBER_TAMPIL: Readonly<Record<string, string>> = {
  RnmLimitListDisplay: 'RnmLimitList',
  RnmGrossPremiDisplay: 'GrossPremiumList',
}

const POLA_INDEKS = /^(\w+)\((\d+)\)\.(\w+)$/

/**
 * Nilai satu sel menurut kuncinya — kunci datar (`Layer`) atau jalur
 * berindeks Pega (`RnmLimitListDisplay(1).Value`, indeks mulai 1).
 */
export function nilaiJalur(b: BarisBersarang, kunci: string): string {
  // Kunci rata apa adanya lebih dulu — bentuk yang sudah diratakan.
  const rata = b[kunci]
  if (typeof rata === 'string') return rata
  const m = POLA_INDEKS.exec(kunci)
  if (m === null) return ''
  const [, larik = '', ke = '1', medan = ''] = m
  let daftar = anak(b, larik)
  const sumber = SUMBER_TAMPIL[larik]
  if (daftar.length === 0 && sumber !== undefined) daftar = anak(b, sumber)
  return teks(daftar[Number(ke) - 1]?.[medan])
}

/** Kunci sel yang berupa jalur berindeks — tidak disunting langsung. */
export const kunciBerjalur = (kunci: string) => POLA_INDEKS.test(kunci)

// ---------------------------------------------------------------------
// ⭐ JALUR HALAMAN — 7 Oktober 2026. Section rincian berjalan di halaman
// BARIS (`TreatyIn.Limits(2).Detail(1)`); rumusnya menulis ke baris itu
// ATAU ke akar panel. Jalur menamai baris itu dari akar, sehingga hasil
// rumus diterapkan atas keadaan TERKINI, bukan salinan saat dipicu.
// ---------------------------------------------------------------------

/** Satu langkah jalur dari akar panel ke halaman baris rincian. */
export interface LangkahJalur {
  larik: string
  /** Indeks baris, mulai 0 (`.pxListSubscript` − 1). */
  indeks: number
}

/** Baris → halaman: skalar ke `medan`, larik anak ke `larik`. */
export function barisKeSisi(b: BarisBersarang): SisiPenyesuaian {
  const medan: Record<string, string> = {}
  const larik: Record<string, BarisBersarang[]> = {}
  for (const [k, v] of Object.entries(b)) {
    if (typeof v === 'string') medan[k] = v
    else if (Array.isArray(v)) larik[k] = v
  }
  return { medan, larik }
}

/** Halaman → baris. */
export const sisiKeBaris = (s: SisiPenyesuaian): BarisBersarang => ({ ...s.medan, ...s.larik })

/** Halaman di ujung jalur; baris yang tidak ada = halaman kosong. */
export function ambilHalaman(akar: SisiPenyesuaian, jalur: readonly LangkahJalur[]): SisiPenyesuaian {
  let h = akar
  for (const j of jalur) h = barisKeSisi(h.larik[j.larik]?.[j.indeks] ?? {})
  return h
}

/**
 * Ubah halaman di ujung jalur — tanpa mengubah aslinya. Baris yang sudah
 * tiada (dihapus selama rumus berjalan) dibiarkan: tidak dilahirkan ulang.
 */
export function ubahHalaman(
  akar: SisiPenyesuaian,
  jalur: readonly LangkahJalur[],
  f: (h: SisiPenyesuaian) => SisiPenyesuaian,
): SisiPenyesuaian {
  const [j, ...sisa] = jalur
  if (j === undefined) return f(akar)
  const baris = akar.larik[j.larik] ?? []
  const b = baris[j.indeks]
  if (b === undefined) return akar
  const baru = sisiKeBaris(ubahHalaman(barisKeSisi(b), sisa, f))
  return { ...akar, larik: { ...akar.larik, [j.larik]: baris.map((x, i) => (i === j.indeks ? baru : x)) } }
}

/** Kesamaan nilai bersarang — urutan kunci tidak berarti. */
export function samaNilai(a: unknown, b: unknown): boolean {
  if (a === b) return true
  if (Array.isArray(a) && Array.isArray(b)) return a.length === b.length && a.every((x, i) => samaNilai(x, b[i]))
  if (typeof a === 'object' && typeof b === 'object' && a !== null && b !== null && !Array.isArray(a) && !Array.isArray(b)) {
    const x = a as Record<string, unknown>
    const y = b as Record<string, unknown>
    const kunci = new Set([...Object.keys(x), ...Object.keys(y)])
    for (const k of kunci) if (!samaNilai(x[k], y[k])) return false
    return true
  }
  return false
}

function gabungLarik(dasar: readonly BarisBersarang[] | undefined, kini: readonly BarisBersarang[] | undefined, hasil: BarisBersarang[]): BarisBersarang[] {
  if (dasar === undefined || kini === undefined || samaNilai(kini, dasar)) return hasil
  // Bentuk larik berubah di kedua sisi — hasil rumus yang berlaku.
  if (dasar.length !== hasil.length || kini.length !== hasil.length) return hasil
  return hasil.map((h, i) => gabungBaris(dasar[i] ?? {}, kini[i] ?? {}, h))
}

function gabungBaris(dasar: BarisBersarang, kini: BarisBersarang, hasil: BarisBersarang): BarisBersarang {
  const out: BarisBersarang = { ...kini }
  for (const [k, v] of Object.entries(hasil)) {
    const d = dasar[k]
    if (samaNilai(v, d)) continue
    const x = kini[k]
    out[k] = Array.isArray(v) && Array.isArray(d) && Array.isArray(x) ? gabungLarik(d, x, v) : v
  }
  return out
}

/**
 * GABUNG TIGA ARAH hasil rumus ke keadaan TERKINI: `dasar` = akar saat rumus
 * dipicu, `hasil` = akar sesudah rumus, `kini` = akar sekarang. Yang rumus
 * UBAH (hasil ≠ dasar) menimpa; yang tidak ia sentuh dibiarkan seperti
 * sekarang.
 *
 * ⛔ PERBAIKAN "input tiba-tiba hilang": rumus berjalan ASINKRON. Isian yang
 * diketik selagi rute menjawab tidak lagi ditimpa nilai lamanya — dahulu
 * jawaban rumus menimpakan larik UTUH yang disalin saat rumus dipicu.
 */
export function gabungTigaArah(dasar: SisiPenyesuaian, kini: SisiPenyesuaian, hasil: SisiPenyesuaian): SisiPenyesuaian {
  const medan = { ...kini.medan }
  for (const [k, v] of Object.entries(hasil.medan)) if (v !== dasar.medan[k]) medan[k] = v
  const larik = { ...kini.larik }
  for (const [k, v] of Object.entries(hasil.larik)) {
    const d = dasar.larik[k]
    if (samaNilai(v, d)) continue
    larik[k] = gabungLarik(d, kini.larik[k], v)
  }
  return { ...kini, medan, larik }
}

/** Larik datar digantikan simpul pohonnya — keduanya membawa skalar yang sama. */
export function gabungPohon(s: SisiPenyesuaian): SisiPenyesuaian {
  if (s.pohon === undefined) return s
  return { medan: s.medan, larik: { ...s.larik, ...s.pohon } }
}
