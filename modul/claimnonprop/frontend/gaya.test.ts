// Disalin dari `modul/claimprop/frontend/gaya.test.ts` (pola, bukan impor): keputusan work owner yang disebut di bawah
// adalah keputusan layar Claim Prop yang ditiru Claim Non Prop (prompt Claim Non Prop tahap 1).
// Gaya modul Claim Prop terisolasi: satu berkas CSS milik modul, setiap pemilih di bawah akar `.claimnonprop`, kelas modul
// berawalan `claimnonprop__`, dan nol properti yang memerangkap Modal.

import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'claimnonprop.css'), 'utf8')
const ATURAN = CSS.replace(/\/\*[\s\S]*?\*\//g, '')

function berkas(dir: string, akhiran: string): string[] {
  return readdirSync(dir).flatMap((n) => {
    const p = join(dir, n)
    if (statSync(p).isDirectory()) return berkas(p, akhiran)
    return n.endsWith(akhiran) ? [p] : []
  })
}

describe('gaya modul Claim Non Prop', () => {
  it('satu-satunya berkas CSS adalah claimnonprop.css; rute.tsx mengimpornya dan membungkus halaman dengan akarnya', () => {
    expect(berkas(AKAR, '.css').map((p) => p.slice(AKAR.length + 1))).toEqual(['claimnonprop.css'])
    const rute = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
    expect(rute).toContain("import './claimnonprop.css'")
    expect(rute).toContain('<div className="claimnonprop">')
    expect(CSS).toMatch(/^\.claimnonprop \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di bawah akar .claimnonprop - nol pemilih global', () => {
    const pemilih = [...ATURAN.matchAll(/([^{}]+)\{/g)]
      .flatMap((m) => m[1]!.split(',').map((s) => s.trim()))
      .filter((p) => !p.startsWith('@'))
    expect(pemilih.length).toBeGreaterThan(5)
    for (const p of pemilih) expect(p, p).toMatch(/^(:root\[data-theme="dark"\] )?\.claimnonprop(\s|$)/)
  })

  it('kelas yang dipakai layar didefinisikan, berawalan claimnonprop__', () => {
    const didefinisikan = new Set([...CSS.matchAll(/\.(claimnonprop__[a-z0-9-]+)/g)].map((m) => m[1]!))
    const dipakai = new Set(
      berkas(AKAR, '.tsx').flatMap((p) =>
        [...readFileSync(p, 'utf8').matchAll(/\b(claimnonprop__[a-z0-9-]+)\b/g)].map((m) => m[1]!),
      ),
    )
    expect(dipakai.size).toBeGreaterThan(5)
    for (const k of dipakai) expect(didefinisikan, k).toContain(k)
  })

  it('kulit Kelola User: akar .claimnonprop__akar bertoken --cnp-* terang dan gelap; warna hanya di blok token', () => {
    expect(CSS).toMatch(/^\.claimnonprop \.claimnonprop__akar \{[^}]*--cnp-latar: #eef1f6;/m)
    expect(CSS).toMatch(/^:root\[data-theme="dark"\] \.claimnonprop \.claimnonprop__akar \{[^}]*--cnp-latar: #1b2130;/m)
    const tanpaToken = ATURAN.replace(/(^|\n)[^{}\n]*\.claimnonprop__akar \{[^}]*\}/g, '')
    expect(tanpaToken).not.toMatch(/#[0-9a-fA-F]{3,8}\b|rgba?\(|hsla?\(/)
    for (const [berkas, kelas] of [
      [['pages', 'ClaimNonProp.tsx'], 'className="inbox claimnonprop__akar"'],
      [['components', 'LayarKasus.tsx'], 'className="inbox claimnonprop__akar"'],
      [['components', 'TataView.tsx'], 'className="claimnonprop__tabel"'],
      [['components', 'Popup.tsx'], 'className="claimnonprop__tabel"'],
    ] as const) {
      expect(readFileSync(join(AKAR, ...berkas), 'utf8'), kelas).toContain(kelas)
    }
  })

  it('halaman awal memakai kelas inti Kelola User: toolbar, field__input, inbox__tabel, btn--primary', () => {
    const awal = readFileSync(join(AKAR, 'pages', 'ClaimNonProp.tsx'), 'utf8')
    for (const k of [
      'className="toolbar"',
      'field__input',
      'className="inbox__tabel claimnonprop__tabel"',
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
      'claimnonprop__baris',
      'claimnonprop__label-medan',
      'claimnonprop__dua',
      'tabs__item',
      'field__error',
    ]) {
      expect(tata, k).toContain(k)
    }
    expect(CSS).toMatch(/\.claimnonprop \.claimnonprop__baris \{[^}]*grid-template-columns: 150px minmax\(0, 1fr\)/)
    // tata letak tanpa flex kolom (basis flex medan pernah menjadi tinggi - 12.204px, 07-10-2026)
    expect(ATURAN).not.toMatch(/flex-direction:\s*column/)
  })

  it('huruf modul diperkecil (work owner 08-10-2026 "font juga kecilin"): akar 13px, isian dan tombol 32px', () => {
    expect(CSS).toMatch(/\.claimnonprop \.claimnonprop__akar \{\s*font-size: 13px;\s*\}/)
    expect(CSS).toMatch(/\.claimnonprop \.claimnonprop__akar \.field__input \{[^}]*height: 32px;[^}]*font-size: 13px;/)
    expect(CSS).toMatch(/\.claimnonprop \.claimnonprop__akar \.btn \{[^}]*height: 32px;[^}]*font-size: 13px;/)
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
