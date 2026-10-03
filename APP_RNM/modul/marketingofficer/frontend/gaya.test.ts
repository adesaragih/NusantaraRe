// Gaya modul Marketing Officer terisolasi: satu berkas CSS milik modul, setiap pemilih di bawah akar
// `.marketingofficer`, kelas modul berawalan `marketingofficer__`, dan nol properti yang memerangkap Modal.

import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'marketingofficer.css'), 'utf8')

function berkas(dir: string, akhiran: string): string[] {
  return readdirSync(dir).flatMap((n) => {
    const p = join(dir, n)
    if (statSync(p).isDirectory()) return berkas(p, akhiran)
    return n.endsWith(akhiran) ? [p] : []
  })
}

describe('gaya modul Marketing Officer', () => {
  it('satu-satunya berkas CSS adalah marketingofficer.css; rute.tsx mengimpornya dan membungkus halaman dengan akarnya', () => {
    expect(berkas(AKAR, '.css').map((p) => p.slice(AKAR.length + 1))).toEqual(['marketingofficer.css'])
    const rute = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
    expect(rute).toContain("import './marketingofficer.css'")
    expect(rute).toContain('<div className="marketingofficer">')
    expect(CSS).toMatch(/^\.marketingofficer \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di bawah akar .marketingofficer - nol pemilih global', () => {
    const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
    const pemilih = [...tanpaKomentar.matchAll(/([^{}]+)\{/g)].flatMap((m) => m[1]!.split(',').map((s) => s.trim()))
    expect(pemilih.length).toBeGreaterThan(5)
    for (const p of pemilih) {
      expect(p, p).toMatch(/^(:root\[data-theme="dark"\] )?\.marketingofficer(\s|$)/)
    }
  })

  it('kelas yang didefinisikan berawalan marketingofficer__, dan setiap yang dipakai layar didefinisikan', () => {
    const didefinisikan = new Set([...CSS.matchAll(/\.(marketingofficer__[a-z0-9-]+)/g)].map((m) => m[1]!))
    const dipakai = new Set(
      berkas(AKAR, '.tsx').flatMap((p) => [...readFileSync(p, 'utf8').matchAll(/\b(marketingofficer__[a-z0-9-]+)\b/g)].map((m) => m[1]!)),
    )
    expect(dipakai.size).toBeGreaterThan(3)
    for (const k of dipakai) expect(didefinisikan, k).toContain(k)
  })

  it('nol properti yang memerangkap Modal tanpa portal', () => {
    expect(CSS).not.toMatch(/(^|[\s;{])(-webkit-)?(backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain)\s*:/m)
  })
})
