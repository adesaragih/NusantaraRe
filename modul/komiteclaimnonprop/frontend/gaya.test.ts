// Gaya modul Komite Claim Non Prop terisolasi: satu berkas CSS milik modul, setiap pemilih di bawah akar
// `.komiteclaimnonprop`, kelas modul berawalan `komiteclaimnonprop__`, warna hanya di blok token akar.

import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'komiteclaimnonprop.css'), 'utf8')
const ATURAN = CSS.replace(/\/\*[\s\S]*?\*\//g, '')

function berkas(dir: string, akhiran: string): string[] {
  return readdirSync(dir).flatMap((n) => {
    const p = join(dir, n)
    if (statSync(p).isDirectory()) return berkas(p, akhiran)
    return n.endsWith(akhiran) ? [p] : []
  })
}

describe('gaya modul Komite Claim Non Prop', () => {
  it('satu-satunya berkas CSS adalah komiteclaimnonprop.css; rute.tsx mengimpornya dan membungkus halaman', () => {
    expect(berkas(AKAR, '.css').map((p) => p.slice(AKAR.length + 1))).toEqual(['komiteclaimnonprop.css'])
    const rute = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
    expect(rute).toContain("import './komiteclaimnonprop.css'")
    expect(rute).toContain('<div className="komiteclaimnonprop">')
    expect(CSS).toMatch(/^\.komiteclaimnonprop \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di bawah akar .komiteclaimnonprop - nol pemilih global', () => {
    const pemilih = [...ATURAN.matchAll(/([^{}]+)\{/g)]
      .flatMap((m) => m[1]!.split(',').map((s) => s.trim()))
      .filter((p) => !p.startsWith('@'))
    expect(pemilih.length).toBeGreaterThan(5)
    for (const p of pemilih) expect(p, p).toMatch(/^(:root\[data-theme="dark"\] )?\.komiteclaimnonprop(\s|$)/)
  })

  it('kelas yang dipakai layar didefinisikan, berawalan komiteclaimnonprop__', () => {
    const didefinisikan = new Set([...CSS.matchAll(/\.(komiteclaimnonprop__[a-z0-9-]+)/g)].map((m) => m[1]!))
    const dipakai = new Set(
      berkas(AKAR, '.tsx').flatMap((p) =>
        [...readFileSync(p, 'utf8').matchAll(/\b(komiteclaimnonprop__[a-z0-9-]+)\b/g)].map((m) => m[1]!),
      ),
    )
    expect(dipakai.size).toBeGreaterThan(5)
    for (const k of dipakai) expect(didefinisikan, k).toContain(k)
  })

  // Laporan work owner 09-10-2026 "table nya sejajarin header sama value": judul kolom angka rata kanan seperti
  // nilainya - aturan `th` tabel (text-align: left, lebih spesifik) tidak boleh menimpanya.
  it('judul kolom angka rata kanan seperti nilainya', () => {
    expect(ATURAN).toMatch(
      /\.komiteclaimnonprop \.komiteclaimnonprop__tabel th\.komiteclaimnonprop__angka \{[^}]*text-align: right;/,
    )
    const kasus = readFileSync(join(AKAR, 'pages', 'KasusKomite.tsx'), 'utf8')
    expect(kasus).toContain("<th key={i} className={k.jenis === 'angka' ? 'komiteclaimnonprop__angka' : undefined}>")
  })

  it('warna hanya di blok token akar; nol backdrop-filter / filter / transform', () => {
    const tanpaToken = ATURAN.replace(/(^|\n)[^{}\n]*\.komiteclaimnonprop__akar \{[^}]*\}/g, '')
    expect(tanpaToken).not.toMatch(/#[0-9a-fA-F]{3,8}\b|rgba?\(|hsla?\(/)
    expect(ATURAN).not.toMatch(/backdrop-filter|(^|[;\s])filter:|transform:/)
  })
})
