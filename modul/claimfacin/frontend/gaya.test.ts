// Disalin dari `modul/claimnonprop/frontend/gaya.test.ts` (pola, bukan impor; asal Claim Prop): keputusan work
// owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Gaya modul Claim Fac In terisolasi: satu berkas CSS milik modul, setiap pemilih di bawah akar `.claimfacin`, kelas modul
// berawalan `claimfacin__`, dan nol properti yang memerangkap Modal.

import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'claimfacin.css'), 'utf8')
const ATURAN = CSS.replace(/\/\*[\s\S]*?\*\//g, '')

function berkas(dir: string, akhiran: string): string[] {
  return readdirSync(dir).flatMap((n) => {
    const p = join(dir, n)
    if (statSync(p).isDirectory()) return berkas(p, akhiran)
    return n.endsWith(akhiran) ? [p] : []
  })
}

describe('gaya modul Claim Fac In', () => {
  it('satu-satunya berkas CSS adalah claimfacin.css; rute.tsx mengimpornya dan membungkus halaman dengan akarnya', () => {
    expect(berkas(AKAR, '.css').map((p) => p.slice(AKAR.length + 1))).toEqual(['claimfacin.css'])
    const rute = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
    expect(rute).toContain("import './claimfacin.css'")
    expect(rute).toContain('<div className="claimfacin">')
    expect(CSS).toMatch(/^\.claimfacin \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di bawah akar .claimfacin - nol pemilih global', () => {
    const pemilih = [...ATURAN.matchAll(/([^{}]+)\{/g)]
      .flatMap((m) => m[1]!.split(',').map((s) => s.trim()))
      .filter((p) => !p.startsWith('@'))
    expect(pemilih.length).toBeGreaterThan(5)
    for (const p of pemilih) expect(p, p).toMatch(/^(:root\[data-theme="dark"\] )?\.claimfacin(\s|$)/)
  })

  it('kelas yang dipakai layar didefinisikan, berawalan claimfacin__', () => {
    const didefinisikan = new Set([...CSS.matchAll(/\.(claimfacin__[a-z0-9-]+)/g)].map((m) => m[1]!))
    const dipakai = new Set(
      berkas(AKAR, '.tsx').flatMap((p) =>
        [...readFileSync(p, 'utf8').matchAll(/\b(claimfacin__[a-z0-9-]+)\b/g)].map((m) => m[1]!),
      ),
    )
    expect(dipakai.size).toBeGreaterThan(5)
    for (const k of dipakai) expect(didefinisikan, k).toContain(k)
  })

  it('kulit Kelola User: akar .claimfacin__akar bertoken --cfi-* terang dan gelap; warna hanya di blok token', () => {
    expect(CSS).toMatch(/^\.claimfacin \.claimfacin__akar \{[^}]*--cfi-latar: #eef1f6;/m)
    expect(CSS).toMatch(/^:root\[data-theme="dark"\] \.claimfacin \.claimfacin__akar \{[^}]*--cfi-latar: #1b2130;/m)
    const tanpaToken = ATURAN.replace(/(^|\n)[^{}\n]*\.claimfacin__akar \{[^}]*\}/g, '')
    expect(tanpaToken).not.toMatch(/#[0-9a-fA-F]{3,8}\b|rgba?\(|hsla?\(/)
    for (const [berkas, kelas] of [
      [['pages', 'ClaimFacIn.tsx'], 'className="inbox claimfacin__akar"'],
      [['components', 'LayarKasus.tsx'], 'className="inbox claimfacin__akar"'],
      [['components', 'TataView.tsx'], 'className="claimfacin__tabel"'],
      [['components', 'Popup.tsx'], 'className="claimfacin__tabel"'],
    ] as const) {
      expect(readFileSync(join(AKAR, ...berkas), 'utf8'), kelas).toContain(kelas)
    }
  })

  it('halaman awal memakai kelas inti Kelola User: toolbar, field__input, inbox__tabel, btn--primary', () => {
    const awal = readFileSync(join(AKAR, 'pages', 'ClaimFacIn.tsx'), 'utf8')
    for (const k of [
      'className="toolbar"',
      'field__input',
      'className="inbox__tabel claimfacin__tabel"',
      'btn btn--primary',
    ]) {
      expect(awal, k).toContain(k)
    }
  })

  it('layar kasus = layout lama Pega dirapikan: label kiri, dua kolom, tab, kartu panel Kelola User', () => {
    const tata = readFileSync(join(AKAR, 'components', 'TataView.tsx'), 'utf8')
    for (const k of [
      'className="panel"',
      'panel__title',
      'claimfacin__baris',
      'claimfacin__label-medan',
      'claimfacin__dua',
      'tabs__item',
      'field__error',
    ]) {
      expect(tata, k).toContain(k)
    }
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__baris \{[^}]*grid-template-columns: 150px minmax\(0, 1fr\)/)
    // tata letak tanpa flex kolom (basis flex medan pernah menjadi tinggi - 12.204px, 07-10-2026)
    expect(ATURAN).not.toMatch(/flex-direction:\s*column/)
  })

  it('huruf modul diperkecil (work owner 08-10-2026 "font juga kecilin"): akar 13px, isian dan tombol 32px', () => {
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__akar \{\s*font-size: 13px;\s*\}/)
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__akar \.field__input \{[^}]*height: 32px;[^}]*font-size: 13px;/)
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__akar \.btn \{[^}]*height: 32px;[^}]*font-size: 13px;/)
  })

  it('nol properti yang memerangkap Modal tanpa portal', () => {
    expect(ATURAN).not.toMatch(
      /(^|[\s;{])(-webkit-)?(backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain)\s*:/m,
    )
  })

  it('tanpa kotak masuk Beranda (bukan bagian Pega)', () => {
    expect(readFileSync(join(AKAR, 'menu.ts'), 'utf8')).not.toMatch(/antreanBeranda\s*:|daftarBeranda\s*:/)
  })
})

describe('gaya panel bersarang Claim Fac In', () => {
  it('panel baris di dalam panel baris tetap bergaris aksen dan berlatar tenang', () => {
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__rinci \.claimfacin__rinci \{/)
  })

  it('tab layout group layar kasus bergaris bawah selebar teks (label tab Pega panjang), bukan kapsul 110 px', () => {
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__tab > \.tabs \{[^}]*width: 100%;/)
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__tab > \.tabs::before \{\s*display: none;/)
  })

  it('grid lebar menggulir di kartunya; daftar autocomplete sel yang terbuka diberi ruang, tidak terpotong', () => {
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__grid \{[^}]*overflow-x: auto;/)
    expect(CSS).toMatch(/\.claimfacin \.claimfacin__grid:has\(.*\.pilih-saring__daftar\) \{\s*padding-bottom: 270px;/)
  })
})
