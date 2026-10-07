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

  it('tema terang dan gelap lewat token --cp-*', () => {
    expect(CSS).toMatch(/^\.claimprop \.claimprop__akar \{[^}]*--cp-latar: #f4f6fa;/m)
    expect(CSS).toMatch(/^:root\[data-theme="dark"\] \.claimprop \.claimprop__akar \{[^}]*--cp-latar: #1b2130;/m)
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
