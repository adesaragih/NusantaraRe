// Asal: salinan `modul/nbtreatyin/frontend/gaya.test.ts` (06-10-2026), disesuaikan akar `.edmtreatyin`.
//
// Penjaga gaya modul EDM Treaty In: CSS modul TERISOLASI dan TIDAK menumpang di inti.
//
//   - satu-satunya berkas CSS modul adalah `edmtreatyin.css`, diimpor `rute.tsx`;
//   - `rute.tsx` membungkus seluruh halaman modul dengan akar `.edmtreatyin` (`display: contents`);
//   - SETIAP pemilih di `edmtreatyin.css` diawali kelas akar - nol pemilih global;
//   - kelas khusus modul ini tidak didefinisikan di `inti/frontend/styles.css`;
//   - tanpa properti yang mengurung popup `position: fixed`.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'edmtreatyin.css'), 'utf8')
const RUTE = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
const INTI = readFileSync(join(AKAR, '..', '..', '..', 'inti', 'frontend', 'styles.css'), 'utf8')

/** Kelas khusus modul ini: wajib hidup di berkas modul, dan nol jejak di inti. */
const KELAS_MODUL = ['edmt__saring', 'edmt__tautan', 'edmt__kaki', 'edmt__lipat', 'edmt__rinci-baris'] as const

/** Setiap berkas di bawah folder frontend modul (rekursif). */
function berkas(d = AKAR): string[] {
  return readdirSync(d, { withFileTypes: true }).flatMap((e) =>
    e.isDirectory() ? berkas(join(d, e.name)) : [join(d, e.name)],
  )
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
  // Tema gelap: `:root[data-theme="dark"] .edmtreatyin …` tetap di bawah akar modul.
  const q = p.startsWith(':root[data-theme="dark"] ') ? p.slice(':root[data-theme="dark"] '.length) : p
  return q === '.edmtreatyin' || q.startsWith('.edmtreatyin ') || q.startsWith(':where(.edmtreatyin ')
}

/** Kelas dari `daftar` yang masih disebut pemilih di `css`. */
function masihDi(css: string, daftar: readonly string[]): string[] {
  const pemilih = pemilihCSS(css)
  return daftar.filter((k) => pemilih.some((p) => new RegExp(`\\.${k.replace(/[-]/g, '\\-')}(?![\\w-])`).test(p)))
}

const PENAMPUNG_FIXED =
  /(?:^|[;{\s])((?:-webkit-)?(?:backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain))\s*:/g

describe('gaya modul EDM Treaty In', () => {
  it('satu-satunya berkas CSS modul adalah edmtreatyin.css; rute.tsx mengimpornya dan membungkus halaman dengan akar', () => {
    expect(
      berkas()
        .filter((f) => f.endsWith('.css'))
        .map((f) => f.slice(AKAR.length + 1)),
    ).toEqual(['edmtreatyin.css'])
    expect(RUTE).toContain("import './edmtreatyin.css'")
    expect(RUTE).toContain('<div className="edmtreatyin">')
    expect(CSS).toMatch(/^\.edmtreatyin \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di edmtreatyin.css diawali kelas akar .edmtreatyin - nol pemilih global', () => {
    const pemilih = pemilihCSS(CSS)
    expect(pemilih.length).toBeGreaterThan(1)
    expect(pemilih.filter((p) => !terisolasi(p))).toEqual([])
  })

  it('aturan isolasi menggigit', () => {
    const contoh = pemilihCSS(
      'table { x: 1 } .edmtreatyin .a, .b { y: 2 } .edmtreatyin :is(.c, .d):disabled { z: 3 } @media (max-width: 1px) { .e { w: 4 } }',
    )
    expect(contoh.filter((p) => !terisolasi(p))).toEqual(['table', '.b', '.e'])
    const gelap = pemilihCSS(':root[data-theme="dark"] .edmtreatyin .a { x: 1 } :root[data-theme="dark"] .a { y: 2 }')
    expect(gelap.filter((p) => !terisolasi(p))).toEqual([':root[data-theme="dark"] .a'])
  })

  it('kelas khusus modul ini tidak menumpang di inti/frontend/styles.css', () => {
    expect(masihDi(INTI, KELAS_MODUL)).toEqual([])
    expect(masihDi(CSS, KELAS_MODUL)).toEqual([...KELAS_MODUL])
  })

  it('aturan "tidak menumpang" menggigit', () => {
    expect(masihDi('.x, .edmt__kaki:hover { a: 1 } .edmt__kaki-lain { b: 2 }', ['edmt__kaki'])).toEqual(['edmt__kaki'])
    expect(masihDi('.edmt__kaki-lain { b: 2 }', ['edmt__kaki'])).toEqual([])
  })

  it('tanpa properti yang mengurung popup position: fixed', () => {
    const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
    expect([...tanpaKomentar.matchAll(PENAMPUNG_FIXED)].map((m) => m[1])).toEqual([])
  })

  it('tidak ada sisa akar / kelas NB di gaya maupun kode modul', () => {
    expect(CSS.replace(/\/\*[\s\S]*?\*\//g, '')).not.toMatch(/nbtreatyin|nbti__|--nbti-/)
    const kode = berkas().filter((f) => /\.tsx?$/.test(f) && !f.endsWith('.test.ts') && !f.endsWith('.test.tsx'))
    const sisa = kode.filter((f) => /className=[^>]*nbti__|['"]nbti__/.test(readFileSync(f, 'utf8')))
    expect(sisa).toEqual([])
  })
})
