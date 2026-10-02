// Penjaga gaya modul Komite Claim Life (02-10-2026): CSS modul TERISOLASI dan TIDAK LAGI menumpang di inti.
//
//   - satu-satunya berkas CSS modul adalah `komiteclaimlife.css`, diimpor `rute.tsx`;
//   - `rute.tsx` membungkus seluruh halaman modul dengan akar `.komiteclaimlife` (`display: contents`, jadi tata
//     letak kerangka tidak berubah);
//   - SETIAP pemilih di `komiteclaimlife.css` diawali kelas akar (atau `:where(.komiteclaimlife ...)` yang menjaga kekhususan
//     nol) - nol pemilih global;
//   - kelas khusus modul ini tidak lagi didefinisikan di `inti/frontend/styles.css` (work owner 02-10-2026:
//     "untuk css style per modul tidak ada lagi menumpang per inti");
//   - tanpa properti yang mengurung popup `position: fixed`.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'komiteclaimlife.css'), 'utf8')
const RUTE = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
const INTI = readFileSync(join(AKAR, '..', '..', '..', 'inti', 'frontend', 'styles.css'), 'utf8')

/** Kelas khusus modul ini yang DULU didefinisikan di inti. */
const KELAS_MODUL = ['inbox__aturan'] as const

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
  return p === '.komiteclaimlife' || p.startsWith('.komiteclaimlife ') || p.startsWith(':where(.komiteclaimlife ')
}

/** Kelas dari `daftar` yang masih disebut pemilih di `css`. */
function masihDi(css: string, daftar: readonly string[]): string[] {
  const pemilih = pemilihCSS(css)
  return daftar.filter((k) => pemilih.some((p) => new RegExp(`\\.${k.replace(/[-]/g, '\\-')}(?![\\w-])`).test(p)))
}

const PENAMPUNG_FIXED = /(?:^|[;{\s])((?:-webkit-)?(?:backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain))\s*:/g

describe('gaya modul Komite Claim Life', () => {
  it('satu-satunya berkas CSS modul adalah komiteclaimlife.css; rute.tsx mengimpornya dan membungkus halaman dengan akar', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['komiteclaimlife.css'])
    expect(RUTE).toContain("import './komiteclaimlife.css'")
    expect(RUTE).toContain('<div className="komiteclaimlife">')
    expect(CSS).toMatch(/^\.komiteclaimlife \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di komiteclaimlife.css diawali kelas akar .komiteclaimlife - nol pemilih global', () => {
    const pemilih = pemilihCSS(CSS)
    expect(pemilih.length).toBeGreaterThan(1)
    expect(pemilih.filter((p) => !terisolasi(p))).toEqual([])
  })

  it('aturan isolasi menggigit', () => {
    const contoh = pemilihCSS('table { x: 1 } .komiteclaimlife .a, .b { y: 2 } .komiteclaimlife :is(.c, .d):disabled { z: 3 } @media (max-width: 1px) { .e { w: 4 } }')
    expect(contoh.filter((p) => !terisolasi(p))).toEqual(['table', '.b', '.e'])
  })

  it('kelas khusus modul ini tidak lagi menumpang di inti/frontend/styles.css', () => {
    expect(masihDi(INTI, KELAS_MODUL)).toEqual([])
    expect(masihDi(CSS, KELAS_MODUL)).toEqual([...KELAS_MODUL])
  })

  it('aturan "tidak menumpang" menggigit', () => {
    expect(masihDi('.x, .inbox__aturan:hover { a: 1 } .inbox__aturan-lain { b: 2 }', ['inbox__aturan'])).toEqual(['inbox__aturan'])
    expect(masihDi('.inbox__aturan-lain { b: 2 }', ['inbox__aturan'])).toEqual([])
  })

  it('tanpa properti yang mengurung popup position: fixed', () => {
    const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
    expect([...tanpaKomentar.matchAll(PENAMPUNG_FIXED)].map((m) => m[1])).toEqual([])
  })
})
