// Gaya modul Claim Prop terisolasi: satu berkas CSS milik modul, setiap pemilih di bawah akar `.claimprop`, kelas modul
// berawalan `claimprop__`, dan nol properti yang memerangkap Modal.

import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'claimprop.css'), 'utf8')
const ATURAN = CSS.replace(/\/\*[\s\S]*?\*\//g, '')

function berkas(dir: string, akhiran: string): string[] {
  return readdirSync(dir).flatMap((n) => {
    const p = join(dir, n)
    if (statSync(p).isDirectory()) return berkas(p, akhiran)
    return n.endsWith(akhiran) ? [p] : []
  })
}

describe('gaya modul Claim Prop', () => {
  it('satu-satunya berkas CSS adalah claimprop.css; rute.tsx mengimpornya dan membungkus halaman dengan akarnya', () => {
    expect(berkas(AKAR, '.css').map((p) => p.slice(AKAR.length + 1))).toEqual(['claimprop.css'])
    const rute = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
    expect(rute).toContain("import './claimprop.css'")
    expect(rute).toContain('<div className="claimprop">')
    expect(CSS).toMatch(/^\.claimprop \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di bawah akar .claimprop - nol pemilih global', () => {
    const pemilih = [...ATURAN.matchAll(/([^{}]+)\{/g)]
      .flatMap((m) => m[1]!.split(',').map((s) => s.trim()))
      .filter((p) => !p.startsWith('@'))
    expect(pemilih.length).toBeGreaterThan(5)
    for (const p of pemilih) expect(p, p).toMatch(/^(:root\[data-theme="dark"\] )?\.claimprop(\s|$)/)
  })

  it('kelas yang dipakai layar didefinisikan, berawalan claimprop__', () => {
    const didefinisikan = new Set([...CSS.matchAll(/\.(claimprop__[a-z0-9-]+)/g)].map((m) => m[1]!))
    const dipakai = new Set(
      berkas(AKAR, '.tsx').flatMap((p) =>
        [...readFileSync(p, 'utf8').matchAll(/\b(claimprop__[a-z0-9-]+)\b/g)].map((m) => m[1]!),
      ),
    )
    expect(dipakai.size).toBeGreaterThan(5)
    for (const k of dipakai) expect(didefinisikan, k).toContain(k)
  })

  it('kulit Kelola User: akar .claimprop__akar bertoken --cp-* terang dan gelap; warna hanya di blok token', () => {
    expect(CSS).toMatch(/^\.claimprop \.claimprop__akar \{[^}]*--cp-latar: #eef1f6;/m)
    expect(CSS).toMatch(/^:root\[data-theme="dark"\] \.claimprop \.claimprop__akar \{[^}]*--cp-latar: #1b2130;/m)
    const tanpaToken = ATURAN.replace(/(^|\n)[^{}\n]*\.claimprop__akar \{[^}]*\}/g, '')
    expect(tanpaToken).not.toMatch(/#[0-9a-fA-F]{3,8}\b|rgba?\(|hsla?\(/)
    for (const [berkas, kelas] of [
      [['pages', 'ClaimProp.tsx'], 'className="inbox claimprop__akar"'],
      [['components', 'LayarKasus.tsx'], 'className="inbox claimprop__akar"'],
      [['components', 'TataView.tsx'], 'className="claimprop__tabel"'],
      [['components', 'Popup.tsx'], 'className="claimprop__tabel"'],
    ] as const) {
      expect(readFileSync(join(AKAR, ...berkas), 'utf8'), kelas).toContain(kelas)
    }
  })

  it('halaman awal memakai kelas inti Kelola User: toolbar, field__input, inbox__tabel, btn--primary', () => {
    const awal = readFileSync(join(AKAR, 'pages', 'ClaimProp.tsx'), 'utf8')
    for (const k of [
      'className="toolbar"',
      'field__input',
      'className="inbox__tabel claimprop__tabel"',
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
      'claimprop__baris',
      'claimprop__label-medan',
      'claimprop__dua',
      'tabs__item',
      'field__error',
    ]) {
      expect(tata, k).toContain(k)
    }
    expect(CSS).toMatch(/\.claimprop \.claimprop__baris \{[^}]*grid-template-columns: 150px minmax\(0, 1fr\)/)
    // tata letak tanpa flex kolom (basis flex medan pernah menjadi tinggi - 12.204px, 07-10-2026)
    expect(ATURAN).not.toMatch(/flex-direction:\s*column/)
  })

  it('huruf modul diperkecil (work owner 08-10-2026 "font juga kecilin"): akar 13px, isian dan tombol 32px', () => {
    expect(CSS).toMatch(/\.claimprop \.claimprop__akar \{\s*font-size: 13px;\s*\}/)
    expect(CSS).toMatch(/\.claimprop \.claimprop__akar \.field__input \{[^}]*height: 32px;[^}]*font-size: 13px;/)
    expect(CSS).toMatch(/\.claimprop \.claimprop__akar \.btn \{[^}]*height: 32px;[^}]*font-size: 13px;/)
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
