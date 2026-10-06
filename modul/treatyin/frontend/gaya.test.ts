// Penjaga gaya modul Treaty In (02-10-2026): CSS modul TERISOLASI dan TIDAK menumpang di inti.
//
//   - satu-satunya berkas CSS modul adalah `treatyin.css`, diimpor `rute.tsx`;
//   - `rute.tsx` membungkus seluruh halaman modul dengan akar `.treatyin` (`display: contents`, jadi tata
//     letak kerangka tidak berubah);
//   - SETIAP pemilih di `treatyin.css` diawali kelas akar (atau `:where(.treatyin ...)` yang menjaga kekhususan
//     nol) - nol pemilih global;
//   - kelas khusus modul ini tidak didefinisikan di `inti/frontend/styles.css` (work owner 02-10-2026:
//     "untuk css style per modul tidak ada lagi menumpang per inti");
//   - tanpa properti yang mengurung popup `position: fixed`.
//
// ⚠️ Modul ini TIDAK pernah menumpang di inti - ia memang belum punya gaya sama sekali sampai hari ini.
// Penjaganya tetap sama bentuknya dengan tujuh modul yang dipindah, supaya yang menambah aturan besok
// tidak perlu tahu modul mana yang lahir dari pemindahan dan mana yang lahir baru.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'treatyin.css'), 'utf8')
const RUTE = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
const INTI = readFileSync(join(AKAR, '..', '..', '..', 'inti', 'frontend', 'styles.css'), 'utf8')

/** Kelas khusus modul ini: wajib hidup di berkas modul, dan nol jejak di inti. */
const KELAS_MODUL = [
  'trin__aturan',
  'trin__lencana',
  // Ronde layar 1 — kelas kedua layar baru.
  'trin__redup',
  'trin__catatan',
  'trin__galat',
  'trin__kepala',
  'trin__id',
  'trin__radio',
  'trin__pilih-luar',
  'trin__centang',
  'trin__aksi',
  // Ronde layar 2 — tata letak.
  'trin__dwikolom',
  'trin__kolom',
  'trin__panel-kepala',
  'trin__kaki',
  'trin__tabel',
  'trin__tabel-kurs',
  'trin__kol-mata-uang',
  'trin__kol-nilai',
  'trin__kol-tanggal',
  'trin__kepala-kolom',
  'trin__kepala-kanan',
  'trin__ikon-saring',
  'trin__baris-saring',
  'trin__belum',
  'trin__belum-judul',
  'trin__belum-petunjuk',
  'trin__teks',
  'trin__teks-asal',
  'trin__teks-lain',
  'trin__spanduk',
] as const

/** Setiap berkas di bawah folder frontend modul (rekursif). */
function berkas(d = AKAR): string[] {
  return readdirSync(d, { withFileTypes: true }).flatMap((e) => (e.isDirectory() ? berkas(join(d, e.name)) : [join(d, e.name)]))
}

/** Pisah daftar pemilih pada koma tingkat atas (koma di dalam `:is(...)`/`:where(...)` tidak memisah). */
function pisahKoma(p: string): string[] {
  const hasil: string[] = []
  let dalam = 0
  let awal = 0
  for (let i = 0; i < p.length; i++) {
    if (p[i] === '(') dalam++
    else if (p[i] === ')') dalam--
    else if (p[i] === ',' && dalam === 0) {
      hasil.push(p.slice(awal, i))
      awal = i + 1
    }
  }
  hasil.push(p.slice(awal))
  return hasil.map((s) => s.trim())
}

/** Pemilih CSS tingkat atas (tanpa komentar; isi @media ikut, kepala @media tidak). */
function pemilihCSS(css: string): string[] {
  const tanpaKomentar = css.replace(/\/\*[\s\S]*?\*\//g, '')
  return [...tanpaKomentar.matchAll(/([^{};]+)\{/g)]
    .map((m) => (m[1] ?? '').trim())
    .filter((p) => p !== '' && !p.startsWith('@'))
    .flatMap(pisahKoma)
}

function terisolasi(p: string): boolean {
  return p === '.treatyin' || p.startsWith('.treatyin ') || p.startsWith(':where(.treatyin ')
}

/** Kelas dari `daftar` yang masih disebut pemilih di `css`. */
function masihDi(css: string, daftar: readonly string[]): string[] {
  const pemilih = pemilihCSS(css)
  return daftar.filter((k) => pemilih.some((p) => new RegExp(`\\.${k.replace(/[-]/g, '\\-')}(?![\\w-])`).test(p)))
}

const PENAMPUNG_FIXED = /(?:^|[;{\s])((?:-webkit-)?(?:backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain))\s*:/g

describe('gaya modul Treaty In', () => {
  it('satu-satunya berkas CSS modul adalah treatyin.css; rute.tsx mengimpornya dan membungkus halaman dengan akar', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['treatyin.css'])
    expect(RUTE).toContain("import './treatyin.css'")
    expect(RUTE).toContain('<div className="treatyin">')
    expect(CSS).toMatch(/^\.treatyin \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di treatyin.css diawali kelas akar .treatyin - nol pemilih global', () => {
    const pemilih = pemilihCSS(CSS)
    expect(pemilih.length).toBeGreaterThan(1)
    expect(pemilih.filter((p) => !terisolasi(p))).toEqual([])
  })

  it('aturan isolasi menggigit', () => {
    const contoh = pemilihCSS('table { x: 1 } .treatyin .a, .b { y: 2 } .treatyin :is(.c, .d):disabled { z: 3 } @media (max-width: 1px) { .e { w: 4 } }')
    expect(contoh.filter((p) => !terisolasi(p))).toEqual(['table', '.b', '.e'])
  })

  it('kelas khusus modul ini tidak menumpang di inti/frontend/styles.css', () => {
    expect(masihDi(INTI, KELAS_MODUL)).toEqual([])
    expect(masihDi(CSS, KELAS_MODUL)).toEqual([...KELAS_MODUL])
  })

  it('aturan "tidak menumpang" menggigit', () => {
    expect(masihDi('.x, .trin__aturan:hover { a: 1 } .trin__aturan-lain { b: 2 }', ['trin__aturan'])).toEqual(['trin__aturan'])
    expect(masihDi('.trin__aturan-lain { b: 2 }', ['trin__aturan'])).toEqual([])
  })

  it('tanpa properti yang mengurung popup position: fixed', () => {
    const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
    expect([...tanpaKomentar.matchAll(PENAMPUNG_FIXED)].map((m) => m[1])).toEqual([])
  })
})
