// Penjaga gaya modul Treaty In Adjustment (02-10-2026): CSS modul TERISOLASI dan TIDAK menumpang di inti.
//
//   - satu-satunya berkas CSS modul adalah `treatyinadjustment.css`, diimpor `rute.tsx`;
//   - `rute.tsx` membungkus seluruh halaman modul dengan akar `.treatyinadjustment` (`display: contents`,
//     jadi tata letak kerangka tidak berubah);
//   - SETIAP pemilih di `treatyinadjustment.css` diawali kelas akar (atau `:where(.treatyinadjustment ...)`
//     yang menjaga kekhususan nol) - nol pemilih global;
//   - kelas khusus modul ini tidak didefinisikan di `inti/frontend/styles.css` (work owner 02-10-2026:
//     "untuk css style per modul tidak ada lagi menumpang per inti");
//   - tanpa properti yang mengurung popup `position: fixed`.
//
// ⚠️ Modul ini TIDAK pernah menumpang di inti - ia memang belum punya gaya sama sekali sampai hari ini.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'treatyinadjustment.css'), 'utf8')
const RUTE = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
const INTI = readFileSync(join(AKAR, '..', '..', '..', 'inti', 'frontend', 'styles.css'), 'utf8')

/** Kelas khusus modul ini: wajib hidup di berkas modul, dan nol jejak di inti. */
const KELAS_MODUL = [
  'tria__aturan',
  'tria__redup',
  'tria__tabel',
  'tria__spanduk',
  'tria__catatan',
  'tria__aksi',
  // Layar Adjustment — 5 Oktober 2026.
  'tria__bandingan',
  'tria__sisi',
  'tria__dwikolom',
  'tria__kolom',
  'tria__kepala',
  'tria__radio',
  'tria__centang',
  'tria__tak-ada',
  'tria__belum',
  'tria__teks',
  'tria__grid',
  'tria__subjudul',
  'tria__angka',
  'tria__nilai',
  'tria__prorata',
  'tria__blok',
  'tria__teks-sel',
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
  // ⭐ 8 Oktober 2026 — tema Treaty Exchange Yearly: varian GELAP ditulis
  // `:root[data-theme="dark"] .treatyinadjustment …`, pengecualian SEMPIT yang
  // sama dengan penjaga `treatyexchangeyearly/frontend/gaya.test.ts`.
  const p2 = p.startsWith(':root[data-theme="dark"] ') ? p.slice(':root[data-theme="dark"] '.length) : p
  return p2 === '.treatyinadjustment' || p2.startsWith('.treatyinadjustment ') || p2.startsWith(':where(.treatyinadjustment ')
}

/** Kelas dari `daftar` yang masih disebut pemilih di `css`. */
function masihDi(css: string, daftar: readonly string[]): string[] {
  const pemilih = pemilihCSS(css)
  return daftar.filter((k) => pemilih.some((p) => new RegExp(`\\.${k.replace(/[-]/g, '\\-')}(?![\\w-])`).test(p)))
}

const PENAMPUNG_FIXED = /(?:^|[;{\s])((?:-webkit-)?(?:backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain))\s*:/g

describe('gaya modul Treaty In Adjustment', () => {
  it('satu-satunya berkas CSS modul adalah treatyinadjustment.css; rute.tsx mengimpornya dan membungkus halaman dengan akar', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['treatyinadjustment.css'])
    expect(RUTE).toContain("import './treatyinadjustment.css'")
    expect(RUTE).toContain('<div className="treatyinadjustment">')
    expect(CSS).toMatch(/^\.treatyinadjustment \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di treatyinadjustment.css diawali kelas akar .treatyinadjustment - nol pemilih global', () => {
    const pemilih = pemilihCSS(CSS)
    expect(pemilih.length).toBeGreaterThan(1)
    expect(pemilih.filter((p) => !terisolasi(p))).toEqual([])
  })

  it('aturan isolasi menggigit', () => {
    const contoh = pemilihCSS('table { x: 1 } .treatyinadjustment .a, .b { y: 2 } .treatyinadjustment :is(.c, .d):disabled { z: 3 } @media (max-width: 1px) { .e { w: 4 } }')
    expect(contoh.filter((p) => !terisolasi(p))).toEqual(['table', '.b', '.e'])
    // Tema gelap berakar modul diterima; tema gelap GLOBAL tetap ditolak.
    const gelap = pemilihCSS(':root[data-theme="dark"] .treatyinadjustment > .inbox { a: 1 } :root[data-theme="dark"] .x { b: 2 }')
    expect(gelap.filter((p) => !terisolasi(p))).toEqual([':root[data-theme="dark"] .x'])
  })

  it('kelas khusus modul ini tidak menumpang di inti/frontend/styles.css', () => {
    expect(masihDi(INTI, KELAS_MODUL)).toEqual([])
    expect(masihDi(CSS, KELAS_MODUL)).toEqual([...KELAS_MODUL])
  })

  it('aturan "tidak menumpang" menggigit', () => {
    expect(masihDi('.x, .tria__aturan:hover { a: 1 } .tria__aturan-lain { b: 2 }', ['tria__aturan'])).toEqual(['tria__aturan'])
    expect(masihDi('.tria__aturan-lain { b: 2 }', ['tria__aturan'])).toEqual([])
  })

  it('tanpa properti yang mengurung popup position: fixed', () => {
    const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
    expect([...tanpaKomentar.matchAll(PENAMPUNG_FIXED)].map((m) => m[1])).toEqual([])
    // `container-type` juga mengurung popup fixed (layout containment).
    expect(tanpaKomentar).not.toMatch(/container(-type)?\s*:/)
  })
})

// ⭐ 8 Oktober 2026 — "perbaiki design nya seperti [Pega] … kalau zoom in
// zoom out … jadi bagus".
describe('bentuk Pega + responsif saat zoom', () => {
  const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')

  it('Old ‖ New bertumpuk menurut lebar wadahnya, bukan lebar layar', () => {
    expect(tanpaKomentar).toMatch(
      /\.treatyinadjustment \.tria__bandingan \{[^}]*grid-template-columns: repeat\(auto-fit, minmax\(min\(100%, 520px\), 1fr\)\);/,
    )
    expect(tanpaKomentar).toMatch(
      /\.treatyinadjustment \.tria__dwikolom \{[^}]*grid-template-columns: repeat\(auto-fit, minmax\(min\(100%, 340px\), 1fr\)\);/,
    )
    expect(tanpaKomentar).not.toMatch(/@media \(max-width: 1100px\)/)
  })

  it('strip tab SATU baris, digulir mendatar', () => {
    expect(tanpaKomentar).toMatch(/\.treatyinadjustment \.tabs \{\s*flex-wrap: nowrap;\s*overflow-x: auto;/)
    expect(tanpaKomentar).toMatch(/\.treatyinadjustment \.tabs__item \{\s*flex: 0 0 auto;/)
  })

  it('label di KIRI kotak kepala (132px), hanya dari 900px ke atas', () => {
    const i = tanpaKomentar.indexOf('@media (min-width: 900px)')
    expect(i).toBeGreaterThan(0)
    const blok = tanpaKomentar.slice(i, i + 700)
    expect(blok).toContain('grid-template-columns: 132px minmax(0, 1fr);')
    expect(blok).toMatch(/\.tria__kolom \.field > :not\(\.field__label\) \{\s*grid-column: 2;/)
  })

  it('kerapatan sama dengan Treaty In: kotak 34px, label 13px', () => {
    expect(tanpaKomentar).toMatch(/\.treatyinadjustment \.field__input \{\s*height: 34px;\s*font-size: 14px;/)
    expect(tanpaKomentar).toMatch(/\.treatyinadjustment \.field__label \{\s*font-size: 13px;/)
  })

  it('kotak di sel grid tidak terpotong; tabel ber-rincian (fixed) tetap boleh menyusut', () => {
    expect(tanpaKomentar).toMatch(/\.treatyinadjustment \.tria__tabel td \.field__input \{\s*min-width: 96px;/)
    expect(tanpaKomentar).toMatch(/\.treatyinadjustment \.tria__tabel td \.tria__tgl \{\s*min-width: 150px;/)
    expect(tanpaKomentar).toMatch(
      /\.treatyinadjustment \.tria__tabel:has\(> tbody > \.tria__rincian\) td \.field__input,\s*\.treatyinadjustment \.tria__tabel:has\(> tbody > \.tria__rincian\) td \.tria__tgl \{\s*min-width: 0;/,
    )
  })
})

// ⭐ 8 Oktober 2026 — "ubah tema … menjadi seperti tema Treaty Exchange Yearly".
describe('tema Treaty Exchange Yearly', () => {
  const css = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
  const akar = css.slice(css.indexOf('.treatyinadjustment > .inbox {'))

  it('token terang di AKAR MODUL (juga layar tanpa `.inbox`), kartu akar halaman: `--tria-*` bernilai `--tey-*`, gradasi lembut, sudut 26px', () => {
    expect(css).toMatch(/^\.treatyinadjustment \{[^}]*--tria-latar: #eef1f6;/m)
    expect(akar.slice(0, 3000)).toContain('border-radius: 26px;')
    expect(akar.slice(0, 3000)).toContain('background: linear-gradient(180deg, var(--tria-latar-atas) 0%, var(--tria-latar-bawah) 100%);')
  })

  it('varian gelap ditulis berakar modul', () => {
    expect(css).toMatch(/^:root\[data-theme="dark"\] \.treatyinadjustment \{[^}]*--tria-isi: #1f2638;/m)
  })

  it('tombol kapsul: utama merah bergradasi, kartu panel putih timbul, isian cekung', () => {
    expect(css).toMatch(/\.treatyinadjustment \.btn \{\s*border-radius: 999px;/)
    expect(css).toMatch(/\.treatyinadjustment \.btn--primary,\s*\.treatyinadjustment \.tl-tambah \{\s*background: linear-gradient\(180deg, var\(--tria-tombol\), var\(--tria-tombol-ujung\)\);/)
    expect(css).toMatch(/\.treatyinadjustment \.panel \{\s*background: var\(--tria-isi\);[^}]*border-radius: 18px;[^}]*box-shadow: var\(--tria-timbul\);/)
    expect(css).toMatch(/\.treatyinadjustment \.field__input \{[^}]*box-shadow: var\(--tria-cekung-kecil\);/)
  })

  it('strip tab = kontrol segmen kapsul; kepala tabel navy muda dan TIDAK tembus pandang', () => {
    expect(css).toMatch(/\.treatyinadjustment \.tabs__item--aktif,\s*\.treatyinadjustment \.tabs__item--aktif:hover \{\s*background: var\(--tria-isi\);/)
    expect(css).toMatch(/\.treatyinadjustment \.table-wrap \.tria__tabel thead th \{\s*background: var\(--tria-kepala-tabel\);/)
  })

  it('lapisan tema TIDAK menimpa ukuran kontrol — kerapatan tetap milik bloknya', () => {
    const tema = css.slice(css.indexOf('.treatyinadjustment > .inbox {'))
    expect(tema).not.toMatch(/\.(btn|field__input)[^{]*\{[^}]*height:/)
  })
})
