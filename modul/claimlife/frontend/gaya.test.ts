// Penjaga gaya modul Claim Life (02-10-2026): CSS modul TERISOLASI dan TIDAK LAGI menumpang di inti.
//
//   - satu-satunya berkas CSS modul adalah `claimlife.css`, diimpor `rute.tsx`;
//   - `rute.tsx` membungkus seluruh halaman modul dengan akar `.claimlife` (`display: contents`, jadi tata
//     letak kerangka tidak berubah);
//   - SETIAP pemilih di `claimlife.css` diawali kelas akar (atau `:where(.claimlife ...)` yang menjaga kekhususan
//     nol) - nol pemilih global;
//   - kelas khusus modul ini tidak lagi didefinisikan di `inti/frontend/styles.css` (work owner 02-10-2026:
//     "untuk css style per modul tidak ada lagi menumpang per inti");
//   - tanpa properti yang mengurung popup `position: fixed`.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'claimlife.css'), 'utf8')
const RUTE = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
const INTI = readFileSync(join(AKAR, '..', '..', '..', 'inti', 'frontend', 'styles.css'), 'utf8')

/** Kelas khusus modul ini yang DULU didefinisikan di inti. */
const KELAS_MODUL = ['inbox__halaman', 'inbox__lencana', 'inbox__register', 'inbox__tab', 'inbox__tab-butir', 'inbox__tab-butir--aktif', 'inbox__tab-id', 'os__aksi', 'os__detail', 'os__judul', 'os__kepala', 'os__tombol', 'polis', 'polis__daftar', 'polis__judul', 'polis__medan', 'polis__medan--belum'] as const

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
  return p === '.claimlife' || p.startsWith('.claimlife ') || p.startsWith(':where(.claimlife ')
}

/** Kelas dari `daftar` yang masih disebut pemilih di `css`. */
function masihDi(css: string, daftar: readonly string[]): string[] {
  const pemilih = pemilihCSS(css)
  return daftar.filter((k) => pemilih.some((p) => new RegExp(`\\.${k.replace(/[-]/g, '\\-')}(?![\\w-])`).test(p)))
}

const PENAMPUNG_FIXED = /(?:^|[;{\s])((?:-webkit-)?(?:backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain))\s*:/g

describe('gaya modul Claim Life', () => {
  it('satu-satunya berkas CSS modul adalah claimlife.css; rute.tsx mengimpornya dan membungkus halaman dengan akar', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['claimlife.css'])
    expect(RUTE).toContain("import './claimlife.css'")
    expect(RUTE).toContain('<div className="claimlife">')
    expect(CSS).toMatch(/^\.claimlife \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di claimlife.css diawali kelas akar .claimlife - nol pemilih global', () => {
    const pemilih = pemilihCSS(CSS)
    expect(pemilih.length).toBeGreaterThan(25)
    expect(pemilih.filter((p) => !terisolasi(p))).toEqual([])
  })

  it('aturan isolasi menggigit', () => {
    const contoh = pemilihCSS('table { x: 1 } .claimlife .a, .b { y: 2 } .claimlife :is(.c, .d):disabled { z: 3 } @media (max-width: 1px) { .e { w: 4 } }')
    expect(contoh.filter((p) => !terisolasi(p))).toEqual(['table', '.b', '.e'])
  })

  it('kelas khusus modul ini tidak lagi menumpang di inti/frontend/styles.css', () => {
    expect(masihDi(INTI, KELAS_MODUL)).toEqual([])
    expect(masihDi(CSS, KELAS_MODUL)).toEqual([...KELAS_MODUL])
  })

  it('aturan "tidak menumpang" menggigit', () => {
    expect(masihDi('.x, .inbox__halaman:hover { a: 1 } .inbox__halaman-lain { b: 2 }', ['inbox__halaman'])).toEqual(['inbox__halaman'])
    expect(masihDi('.inbox__halaman-lain { b: 2 }', ['inbox__halaman'])).toEqual([])
  })

  it('tanpa properti yang mengurung popup position: fixed', () => {
    const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
    expect([...tanpaKomentar.matchAll(PENAMPUNG_FIXED)].map((m) => m[1])).toEqual([])
  })
})

// Kerapian tata letak (11-10-2026, "jangan ada kolom yang miring2 atau terpotong"): aturan pentingnya dikunci.
describe('kerapian tata letak Claim Life', () => {
  const satu = CSS.replace(/\s+/g, ' ')
  const aturan = (pemilih: string): string => {
    const i = satu.lastIndexOf(pemilih + ' {')
    return i < 0 ? '' : satu.slice(i, satu.indexOf('}', i))
  }
  it('sel rincian di grid kembali membungkus (inti: .inbox__tabel td nowrap)', () => {
    expect(aturan('.claimlife .os__rincian > td')).toContain('white-space: normal')
  })
  it('anak grid form boleh menyusut; wadah tidak melebihi halaman', () => {
    expect(aturan('.claimlife .os__form > *, .claimlife .os__kolom > *')).toContain('min-width: 0')
    expect(aturan('.claimlife .os__form, .claimlife .os__kolom, .claimlife .rincian')).toContain('max-width: 100%')
  })
  it('label tetap lebar dan nowrap; nilai kode nowrap; teks panjang membungkus', () => {
    expect(satu).toContain('.claimlife .os__baris-form { grid-template-columns: 24ch minmax(0, 1fr);')
    expect(aturan('.claimlife .os__baris-form > dt, .claimlife .os__baris-form > .os__label')).toContain('white-space: nowrap')
    // 11-10-2026: nilai kode TIDAK nowrap (meluber ke kolom sebelah di Register Claim) - aturan TERAKHIR membungkus.
    expect(aturan('.claimlife .os__nilai-kode')).toContain('white-space: normal')
    expect(aturan('.claimlife .os__nilai-kode')).toContain('overflow-wrap: anywhere')
    expect(aturan('.claimlife .os__baris-form > dd, .claimlife .os__baris-form > span:not(.os__label)')).toContain('overflow-wrap: anywhere')
  })
  it('isian tanggal/angka berlebar tetap (ch), radio satu baris', () => {
    expect(aturan('.claimlife .os__baris-form input[type="date"]')).toMatch(/width: \d+ch/)
    expect(aturan('.claimlife .os__isian-angka')).toMatch(/width: \d+ch/)
    expect(aturan('.claimlife .os__radio')).toContain('white-space: nowrap')
    expect(aturan('.claimlife .os__angka')).toContain('text-align: right')
  })
})

describe('konsistensi isian Outstanding (11-10-2026)', () => {
  const satu = CSS.replace(/\s+/g, ' ')
  it('lebar isian diatur menimpa aturan lama; kolom nilai angka selebar tetap; textarea inti tinggi otomatis', () => {
    expect(satu).toContain('.claimlife .os__baris-form input.os__isian-angka { width: 100%; text-align: right; }')
    expect(satu).toContain('.claimlife .os__form--angka .os__baris-form { grid-template-columns: minmax(0, 22ch) minmax(0, 20ch); }')
    expect(satu).toContain('.claimlife textarea.field__input { height: auto;')
    expect(satu).toMatch(/\.claimlife \.os__baris-form input\.os__isian-tanggal \{ width: \d+ch;/)
  })
})

describe('inbox berkulit Claim Fac In (work owner 12-10-2026) - tampilan saja', () => {
  const INBOX = readFileSync(join(AKAR, 'pages', 'InboxClaimLife.tsx'), 'utf8')
  const KOMITE = readFileSync(join(AKAR, 'components', 'TabelKomite.tsx'), 'utf8')
  it('kelas kulit dipasang pada elemen yang sudah ada; tanpa impor dari modul claimfacin', () => {
    expect(INBOX).toContain('<section className="inbox claimlife__akar">')
    expect(INBOX).toContain('<table className="inbox__tabel claimlife__tabel">')
    expect(KOMITE).toContain('<table className="inbox__tabel claimlife__tabel">')
    for (const s of [INBOX, KOMITE]) expect(s).not.toMatch(/from '[^']*claimfacin/)
  })
  it('kulit: token --cl-* di akar, tab bersegmen, tabel kartu berkepala, tema gelap lewat :where', () => {
    for (const p of ['.claimlife .claimlife__akar {', '.claimlife .claimlife__akar .inbox__tab {', '.claimlife .claimlife__tabel thead th {',
      '.claimlife .claimlife__tabel tbody tr:hover > td {', '.claimlife .claimlife__akar:where(:root[data-theme="dark"] *) {'])
      expect(CSS).toContain(p)
    expect(CSS).toContain('--cl-tombol: #d63a42;')
  })
})
