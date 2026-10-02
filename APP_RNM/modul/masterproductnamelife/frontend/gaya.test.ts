// Penjaga gaya modul (brief bab 1): kelas CSS khusus modul BERAWALAN `mpnl` dan tinggal di berkas CSS modul
// sendiri (`masterproductnamelife.css`, diimpor `rute.tsx`); kelas bersama (`panel`, `btn`, `inbox__tabel`, ...)
// dari `inti` tidak didefinisikan ulang di sini.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname

/** Setiap berkas di bawah folder frontend modul (rekursif). */
function berkas(d = AKAR): string[] {
  return readdirSync(d, { withFileTypes: true }).flatMap((e) => (e.isDirectory() ? berkas(join(d, e.name)) : [join(d, e.name)]))
}

const CSS = readFileSync(join(AKAR, 'masterproductnamelife.css'), 'utf8')

/** Kelas yang didefinisikan pemilih di CSS modul. */
function kelasCSS(): Set<string> {
  const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
  const pemilih = [...tanpaKomentar.matchAll(/([^{}]+)\{/g)].map((m) => m[1] ?? '')
  return new Set(pemilih.flatMap((p) => [...p.matchAll(/\.([a-zA-Z0-9_-]+)/g)].map((m) => m[1] ?? '')))
}

const POLA_KELAS = /\bmpnl(?:-[a-z0-9]+)+(?:__[a-z0-9]+)?(?:--[a-z0-9]+)?\b/g

/** Kelas berawalan `mpnl` yang dipakai layar (isi `className`); nama halaman (`mpnl-produk`) bukan kelas. */
function kelasTSX(): Set<string> {
  const isi = berkas()
    .filter((f) => /\.tsx$/.test(f))
    .map((f) => readFileSync(f, 'utf8'))
  const potongan = isi.flatMap((s) => [
    ...[...s.matchAll(/className=\{?["'`]([^"'`]*)["'`]/g)].map((m) => m[1] ?? ''),
    ...[...s.matchAll(/'(mpnl-[a-z0-9_-]+(?: [a-z0-9_ -]+)?)'/g)].map((m) => m[1] ?? ''),
    ...[...s.matchAll(/' (mpnl-[a-z0-9_-]+)'/g)].map((m) => m[1] ?? ''),
  ])
  return new Set(potongan.flatMap((p) => [...p.matchAll(POLA_KELAS)].map((m) => m[0])))
}

describe('kelas CSS modul', () => {
  it('satu-satunya berkas CSS modul adalah masterproductnamelife.css, dan rute.tsx mengimpornya', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['masterproductnamelife.css'])
    expect(readFileSync(join(AKAR, 'rute.tsx'), 'utf8')).toContain("import './masterproductnamelife.css'")
  })

  it('setiap kelas di CSS modul berawalan mpnl, kecuali kelas bersama inti yang DITIMPA di bawah .mpnl', () => {
    // UI 02-10-2026 (gaya Treaty Contract Out): kelas bersama boleh ditimpa HANYA di bawah kelas akar
    // `.mpnl` (dijaga "isolasi CSS modul" di bawah) dan hanya kelas kerangka inti yang dipakai layar ini.
    const BERSAMA = new Set(['inbox__kepala', 'panel', 'inbox__tabel', 'table__actions', 'btn', 'aksi-baris', 'modal', 'modal__actions', 'inbox__judul', 'panel__title', 'btn--primary', 'btn--ghost', 'field__input', 'modal__title', 'modal__head', 'muted', 'alert'])
    const kelas = [...kelasCSS()]
    expect(kelas.length).toBeGreaterThan(5)
    expect(kelas.filter((k) => !k.startsWith('mpnl') && !BERSAMA.has(k))).toEqual([])
  })

  it('setiap kelas mpnl yang dipakai layar didefinisikan di CSS modul', () => {
    const didefinisikan = kelasCSS()
    const dipakai = [...kelasTSX()].filter((k) => k !== 'mpnl-produk')
    expect(dipakai.length).toBeGreaterThan(5)
    expect(dipakai.filter((k) => !didefinisikan.has(k))).toEqual([])
  })
})

/** Pemilih CSS tingkat atas (tanpa komentar; isi @media ikut, kepala @media tidak). */
function pemilihCSS(css: string): string[] {
  const tanpaKomentar = css.replace(/\/\*[\s\S]*?\*\//g, '')
  return [...tanpaKomentar.matchAll(/([^{};]+)\{/g)]
    .map((m) => (m[1] ?? '').trim())
    .filter((p) => p !== '' && !p.startsWith('@'))
    .flatMap((p) => p.split(',').map((s) => s.trim()))
}

/** Pemilih terisolasi: kelas akar .mpnl, keturunannya, atau varian tema gelap di bawah kelas akar. */
function terisolasi(p: string): boolean {
  const tanpaTema = p.replace(/^:root\[data-theme="dark"\]\s+/, '')
  return tanpaTema === '.mpnl' || tanpaTema.startsWith('.mpnl ')
}

describe('isolasi CSS modul (UI 02-10-2026, gaya Treaty Contract Out)', () => {
  it('SETIAP pemilih di masterproductnamelife.css diawali kelas akar .mpnl - nol pemilih global', () => {
    const pemilih = pemilihCSS(CSS)
    expect(pemilih.length).toBeGreaterThan(20)
    expect(pemilih.filter((p) => !terisolasi(p))).toEqual([])
  })

  it('aturan isolasi menggigit: pemilih global tertangkap', () => {
    const pemilih = pemilihCSS('table { x: 1 } .mpnl .a, .btn { y: 2 } @media (max-width: 640px) { .mpnl-b { z: 3 } }')
    expect(pemilih.filter((p) => !terisolasi(p))).toEqual(['table', '.btn', '.mpnl-b'])
    // Tema gelap boleh, HANYA di bawah kelas akar; :root tanpa .mpnl tetap ditolak.
    expect(terisolasi(':root[data-theme="dark"] .mpnl')).toBe(true)
    expect(terisolasi(':root[data-theme="dark"] .panel')).toBe(false)
  })

  it('halaman awal memasang kelas akar mpnl', () => {
    const kode = readFileSync(join(AKAR, 'pages', 'MasterProductNameLife.tsx'), 'utf8')
    expect(kode.match(/<section className="inbox mpnl">/g)?.length).toBe(1)
    expect(kode).not.toContain('<section className="inbox">')
  })

  it('setiap tabel modul berada di pembungkus gulir mpnl-tabel', () => {
    for (const f of berkas().filter((x) => x.endsWith('.tsx'))) {
      const kode = readFileSync(f, 'utf8')
      const tabel = kode.match(/<table\b/g)?.length ?? 0
      const bungkus = kode.match(/<div className="mpnl-tabel">\s*<table\b/g)?.length ?? 0
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
  it('masterproductnamelife.css tidak memakai properti yang mengurung popup position: fixed', () => {
    expect(penampungFixed(CSS)).toEqual([])
  })

  it('aturan penampung menggigit', () => {
    expect(
      penampungFixed('.mpnl .panel { -webkit-backdrop-filter: blur(4px); backdrop-filter: blur(4px) } .mpnl .a{transform: none; will-change:x} /* filter: x */ .mpnl .b { container-type: inline-size }'),
    ).toEqual(['-webkit-backdrop-filter', 'backdrop-filter', 'transform', 'will-change'])
  })
})

/** Rasio kontras WCAG 2 dua warna hex `#rrggbb`. */
function kontras(a: string, b: string): number {
  const terang = (hex: string): number => {
    const [r, g, bl] = [1, 3, 5].map((i) => {
      const v = parseInt(hex.slice(i, i + 2), 16) / 255
      return v <= 0.03928 ? v / 12.92 : ((v + 0.055) / 1.055) ** 2.4
    })
    return 0.2126 * (r ?? 0) + 0.7152 * (g ?? 0) + 0.0722 * (bl ?? 0)
  }
  const [t, g] = [terang(a), terang(b)].sort((x, y) => y - x)
  return ((t ?? 0) + 0.05) / ((g ?? 0) + 0.05)
}

/** Token `--mt-*` bernilai hex dari aturan yang pemilihnya PERSIS `pemilih`. */
function tokenMt(css: string, pemilih: string): Record<string, string> {
  const tanpaKomentar = css.replace(/\/\*[\s\S]*?\*\//g, '')
  return Object.fromEntries(
    [...tanpaKomentar.matchAll(/([^{};]+)\{([^{}]*)\}/g)]
      .filter((m) => (m[1] ?? '').trim() === pemilih)
      .flatMap((m) => [...(m[2] ?? '').matchAll(/(--mt-[a-z-]+)\s*:\s*(#[0-9a-f]{6})\s*;/gi)].map((x) => [x[1] ?? '', x[2] ?? ''])),
  )
}

describe('soft UI: teks tetap terbaca (02-10-2026)', () => {
  it('kontras token teks --mt-* terhadap latar, kartu, kepala tabel, dan sorot baris minimal 4,5:1, terang dan gelap', () => {
    const terang = tokenMt(CSS, '.mpnl')
    const gelap = { ...terang, ...tokenMt(CSS, ':root[data-theme="dark"] .mpnl') }
    const pasangan = [
      ['--mt-teks', '--mt-latar'],
      ['--mt-teks', '--mt-kartu'],
      ['--mt-teks', '--mt-baris-hover'],
      ['--mt-teks-redup', '--mt-latar'],
      ['--mt-teks-redup', '--mt-kartu'],
      ['--mt-teks-redup', '--mt-kepala-tabel'],
    ] as const
    const kurang = (t: Record<string, string>, nama: string): string[] =>
      pasangan.flatMap(([a, b]) => {
        const x = t[a]
        const y = t[b]
        if (x === undefined || y === undefined) return [`${nama} ${a}/${b} token hilang`]
        const r = kontras(x, y)
        return r < 4.5 ? [`${nama} ${a}/${b} ${r.toFixed(2)}`] : []
      })
    expect(Object.keys(terang).length).toBeGreaterThan(6)
    expect([...kurang(terang, 'terang'), ...kurang(gelap, 'gelap')]).toEqual([])
  })

  it('rumus kontras menggigit: hitam/putih 21:1, #777777/putih di bawah 4,5:1', () => {
    expect(kontras('#000000', '#ffffff')).toBeCloseTo(21, 5)
    expect(kontras('#777777', '#ffffff')).toBeLessThan(4.5)
  })
})
