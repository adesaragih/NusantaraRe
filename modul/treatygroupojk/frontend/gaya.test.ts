// Gaya modul Treaty Group OJK terisolasi: satu berkas CSS milik modul, setiap pemilih di bawah akar
// `.treatygroupojk`, kelas modul berawalan `treatygroupojk__`, dan nol properti yang memerangkap Modal.

import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'treatygroupojk.css'), 'utf8')
/** CSS tanpa komentar - komentar boleh menyebut kelas inti atau properti terlarang sebagai penjelasan. */
const ATURAN = CSS.replace(/\/\*[\s\S]*?\*\//g, '')

function berkas(dir: string, akhiran: string): string[] {
  return readdirSync(dir).flatMap((n) => {
    const p = join(dir, n)
    if (statSync(p).isDirectory()) return berkas(p, akhiran)
    return n.endsWith(akhiran) ? [p] : []
  })
}

describe('gaya modul Treaty Group OJK', () => {
  it('satu-satunya berkas CSS adalah treatygroupojk.css; rute.tsx mengimpornya dan membungkus halaman dengan akarnya', () => {
    expect(berkas(AKAR, '.css').map((p) => p.slice(AKAR.length + 1))).toEqual(['treatygroupojk.css'])
    const rute = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
    expect(rute).toContain("import './treatygroupojk.css'")
    expect(rute).toContain('<div className="treatygroupojk">')
    expect(CSS).toMatch(/^\.treatygroupojk \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di bawah akar .treatygroupojk - nol pemilih global', () => {
    // Aturan `@media` bukan pemilih; pemilih DI DALAMNYA tetap diperiksa.
    const pemilih = [...ATURAN.matchAll(/([^{}]+)\{/g)]
      .flatMap((m) => m[1]!.split(',').map((s) => s.trim()))
      .filter((p) => !p.startsWith('@'))
    expect(pemilih.length).toBeGreaterThan(5)
    for (const p of pemilih) {
      expect(p, p).toMatch(/^(:root\[data-theme="dark"\] )?\.treatygroupojk(\s|$)/)
    }
  })

  it('kelas yang didefinisikan berawalan treatygroupojk__, dan setiap yang dipakai layar didefinisikan', () => {
    const didefinisikan = new Set([...CSS.matchAll(/\.(treatygroupojk__[a-z0-9-]+)/g)].map((m) => m[1]!))
    const dipakai = new Set(
      berkas(AKAR, '.tsx').flatMap((p) =>
        [...readFileSync(p, 'utf8').matchAll(/\b(treatygroupojk__[a-z0-9-]+)\b/g)].map((m) => m[1]!),
      ),
    )
    expect(dipakai.size).toBeGreaterThan(3)
    for (const k of dipakai) expect(didefinisikan, k).toContain(k)
  })

  it('tema modul: halaman berakar .treatygroupojk__akar dengan token --tgo-* terang dan gelap', () => {
    const halaman = readFileSync(join(AKAR, 'pages', 'TreatyGroupOjk.tsx'), 'utf8')
    expect(halaman).toContain('<section className="inbox treatygroupojk__akar">')
    expect(CSS).toMatch(/^\.treatygroupojk \.treatygroupojk__akar \{[^}]*--tgo-latar: #eef1f6;/m)
    expect(CSS).toMatch(
      /^:root\[data-theme="dark"\] \.treatygroupojk \.treatygroupojk__akar \{[^}]*--tgo-latar: #1b2130;/m,
    )
    expect(ATURAN).not.toMatch(/kelola-user|marketingofficer/)
  })

  it('tombol aksi sebaris ke samping, tidak turun ke bawah', () => {
    const aksi = /\.treatygroupojk \.treatygroupojk__aksi \{([^}]*)\}/.exec(ATURAN)?.[1] ?? ''
    expect(aksi).toContain('flex-wrap: nowrap;')
    expect(aksi).toContain('white-space: nowrap;')
  })

  it('nol properti yang memerangkap Modal tanpa portal', () => {
    expect(ATURAN).not.toMatch(
      /(^|[\s;{])(-webkit-)?(backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain)\s*:/m,
    )
  })
})
