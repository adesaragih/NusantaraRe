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
  return pemilihCSS(css).filter((p) => p !== '.tco' && !p.startsWith('.tco '))
}

/** Kelas yang disebut pemilih di `tco.css`. */
function kelasCSS(): Set<string> {
  return new Set(pemilihCSS(CSS).flatMap((p) => [...p.matchAll(/\.([a-zA-Z0-9_-]+)/g)].map((m) => m[1] ?? '')))
}

const KELAS_BERSAMA = new Set(['inbox__kepala', 'panel', 'inbox__tabel', 'inbox__rinci', 'table__actions', 'btn', 'aksi-baris', 'modal', 'modal__actions'])

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
