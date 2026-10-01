// Penjaga gaya modul Treaty Contract Out (UI 02-10-2026): CSS modul TERISOLASI.
//
//   - satu-satunya berkas CSS modul adalah `tco.css`, diimpor `rute.tsx`;
//   - SETIAP pemilih di `tco.css` diawali kelas akar `.tco` - nol pemilih global;
//   - kelas di `tco.css` berawalan `tco`, kecuali kelas kerangka inti yang DITIMPA di bawah `.tco`;
//   - ketiga halaman memasang kelas akar `tco`;
//   - setiap tabel modul berada di pembungkus gulir `tco-tabel` (tabel lebar menggulir di dalam panel,
//     tidak keluar dari panel - contoh bug tangkapan layar 02-10-2026).

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname

/** Setiap berkas di bawah folder frontend modul (rekursif). */
function berkas(d = AKAR): string[] {
  return readdirSync(d, { withFileTypes: true }).flatMap((e) => (e.isDirectory() ? berkas(join(d, e.name)) : [join(d, e.name)]))
}

const CSS = readFileSync(join(AKAR, 'tco.css'), 'utf8')

/** Pemilih CSS tingkat atas (tanpa komentar; isi @media ikut, kepala @media tidak). */
function pemilihCSS(css: string): string[] {
  const tanpaKomentar = css.replace(/\/\*[\s\S]*?\*\//g, '')
  return [...tanpaKomentar.matchAll(/([^{};]+)\{/g)]
    .map((m) => (m[1] ?? '').trim())
    .filter((p) => p !== '' && !p.startsWith('@'))
    .flatMap((p) => p.split(',').map((s) => s.trim()))
}

function tidakTerisolasi(css: string): string[] {
  return pemilihCSS(css).filter((p) => {
    // Tema gelap boleh, HANYA di bawah kelas akar.
    const tanpaTema = p.replace(/^:root\[data-theme="dark"\]\s+/, '')
    return tanpaTema !== '.tco' && !tanpaTema.startsWith('.tco ')
  })
}

/** Kelas yang disebut pemilih di `tco.css`. */
function kelasCSS(): Set<string> {
  return new Set(pemilihCSS(CSS).flatMap((p) => [...p.matchAll(/\.([a-zA-Z0-9_-]+)/g)].map((m) => m[1] ?? '')))
}

const KELAS_BERSAMA = new Set(['inbox__kepala', 'panel', 'inbox__tabel', 'inbox__rinci', 'table__actions', 'btn', 'aksi-baris', 'modal', 'modal__actions', 'inbox__judul', 'panel__title', 'btn--primary', 'btn--ghost', 'field__input'])

describe('isolasi CSS modul Treaty Contract Out', () => {
  it('satu-satunya berkas CSS modul adalah tco.css, dan rute.tsx mengimpornya', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['tco.css'])
    expect(readFileSync(join(AKAR, 'rute.tsx'), 'utf8')).toContain("import './tco.css'")
  })

  it('SETIAP pemilih di tco.css diawali kelas akar .tco - nol pemilih global', () => {
    expect(pemilihCSS(CSS).length).toBeGreaterThan(20)
    expect(tidakTerisolasi(CSS)).toEqual([])
  })

  it('aturan isolasi menggigit: pemilih global tertangkap', () => {
    expect(tidakTerisolasi('table { x: 1 } .tco .a, .btn { y: 2 } @media (max-width: 640px) { .tco-b { z: 3 } }')).toEqual([
      'table',
      '.btn',
      '.tco-b',
    ])
    expect(tidakTerisolasi(':root[data-theme="dark"] .tco { a: 1 } :root[data-theme="dark"] .panel { b: 2 }')).toEqual([
      ':root[data-theme="dark"] .panel',
    ])
  })

  it('kelas di tco.css berawalan tco, kecuali kelas kerangka inti yang ditimpa', () => {
    expect([...kelasCSS()].filter((k) => !k.startsWith('tco') && !KELAS_BERSAMA.has(k))).toEqual([])
  })

  it('ketiga halaman memasang kelas akar tco', () => {
    const tahun = readFileSync(join(AKAR, 'pages', 'InboxTreatyContract.tsx'), 'utf8')
    expect(tahun.match(/<section className="inbox tco">/g)?.length).toBe(2)
    expect(tahun).not.toContain('<section className="inbox">')
    for (const h of ['InboxTreatyContractDescription.tsx', 'InboxTreatyContractReinsType.tsx']) {
      const kode = readFileSync(join(AKAR, 'pages', h), 'utf8')
      expect(kode).toContain('<div className="inbox tco">')
      expect(kode).not.toContain('<div className="inbox">')
    }
  })

  it('setiap tabel modul berada di pembungkus gulir tco-tabel', () => {
    for (const f of berkas().filter((x) => x.endsWith('.tsx'))) {
      const kode = readFileSync(f, 'utf8')
      const tabel = kode.match(/<table\b/g)?.length ?? 0
      const bungkus = kode.match(/<div className="tco-tabel">\s*<table\b/g)?.length ?? 0
      expect(`${f.slice(AKAR.length + 1)}: ${bungkus}/${tabel}`).toBe(`${f.slice(AKAR.length + 1)}: ${tabel}/${tabel}`)
    }
  })
})

/**
 * Properti yang menjadikan elemen blok penampung bagi keturunan `position: fixed`. Popup inti
 * (`.modal__backdrop`, `position: fixed; inset: 0`) dirender DI DALAM panel modul, bukan lewat portal;
 * bila satu leluhurnya memakai salah satu properti ini, popup menempel ke panel dan tidak lagi di tengah
 * layar (bug 02-10-2026: `backdrop-filter` pada panel tema login). `container-type: inline-size` aman
 * (diukur di Edge: popup di dalam pembungkus tabel tetap menutup layar penuh).
 */
const PENAMPUNG_FIXED = /(?:^|[;{\s])((?:-webkit-)?(?:backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain))\s*:/g

function penampungFixed(css: string): string[] {
  const tanpaKomentar = css.replace(/\/\*[\s\S]*?\*\//g, '')
  return [...tanpaKomentar.matchAll(PENAMPUNG_FIXED)].map((m) => m[1] ?? '')
}

describe('popup tetap di tengah layar', () => {
  it('tco.css tidak memakai properti yang mengurung popup position: fixed', () => {
    expect(penampungFixed(CSS)).toEqual([])
  })

  it('aturan penampung menggigit', () => {
    expect(
      penampungFixed('.tco .panel { -webkit-backdrop-filter: blur(4px); backdrop-filter: blur(4px) } .tco .a{transform: none; will-change:x} /* filter: x */ .tco .b { container-type: inline-size }'),
    ).toEqual(['-webkit-backdrop-filter', 'backdrop-filter', 'transform', 'will-change'])
  })
})
