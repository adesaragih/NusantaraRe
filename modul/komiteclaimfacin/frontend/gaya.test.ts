// Gaya modul Komite Claim Fac In terisolasi (pola `modul/komiteclaimnonprop/frontend/gaya.test.ts`): satu berkas CSS milik
// modul, setiap pemilih di bawah akar `.komiteclaimfacin`, kelas modul berawalan `komiteclaimfacin__`, warna hanya di
// blok token akar (`--kcfi-*`, terang dan gelap).

import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'komiteclaimfacin.css'), 'utf8')
const ATURAN = CSS.replace(/\/\*[\s\S]*?\*\//g, '')

function berkas(dir: string, akhiran: string): string[] {
  return readdirSync(dir).flatMap((n) => {
    const p = join(dir, n)
    if (statSync(p).isDirectory()) return n === 'node_modules' ? [] : berkas(p, akhiran)
    return n.endsWith(akhiran) ? [p] : []
  })
}

describe('gaya modul Komite Claim Fac In', () => {
  it('satu-satunya berkas CSS adalah komiteclaimfacin.css; rute.tsx mengimpornya dan membungkus halaman', () => {
    expect(berkas(AKAR, '.css').map((p) => p.slice(AKAR.length + 1))).toEqual(['komiteclaimfacin.css'])
    const rute = readFileSync(join(AKAR, 'rute.tsx'), 'utf8')
    expect(rute).toContain("import './komiteclaimfacin.css'")
    expect(rute).toContain('<div className="komiteclaimfacin">')
    expect(CSS).toMatch(/^\.komiteclaimfacin \{\s*display: contents;\s*\}/m)
  })

  it('SETIAP pemilih di bawah akar .komiteclaimfacin - nol pemilih global', () => {
    const pemilih = [...ATURAN.matchAll(/([^{}]+)\{/g)]
      .flatMap((m) => m[1]!.split(',').map((s) => s.trim()))
      .filter((p) => !p.startsWith('@'))
    expect(pemilih.length).toBeGreaterThan(5)
    for (const p of pemilih) expect(p, p).toMatch(/^(:root\[data-theme="dark"\] )?\.komiteclaimfacin(\s|$)/)
  })

  it('kelas yang dipakai layar didefinisikan, berawalan komiteclaimfacin__', () => {
    const didefinisikan = new Set([...CSS.matchAll(/\.(komiteclaimfacin__[a-z0-9-]+)/g)].map((m) => m[1]!))
    const dipakai = new Set(
      berkas(AKAR, '.tsx').flatMap((p) =>
        [...readFileSync(p, 'utf8').matchAll(/\b(komiteclaimfacin__[a-z0-9-]+)\b/g)].map((m) => m[1]!),
      ),
    )
    expect(dipakai.size).toBeGreaterThan(5)
    for (const k of dipakai) expect(didefinisikan, k).toContain(k)
  })

  it('kulit bertoken --kcfi-* terang dan gelap di akar .komiteclaimfacin__akar', () => {
    expect(CSS).toMatch(/^\.komiteclaimfacin \.komiteclaimfacin__akar \{[^}]*--kcfi-teks: #101828;/m)
    expect(CSS).toMatch(
      /^:root\[data-theme="dark"\] \.komiteclaimfacin \.komiteclaimfacin__akar \{[^}]*--kcfi-teks: #eef0f6;/m,
    )
    expect(ATURAN).not.toMatch(/--kcnp-|--kcp-|--cfi-/)
  })

  // Laporan work owner 09-10-2026 "table nya sejajarin header sama value" (Komite Prop): judul kolom angka rata kanan
  // seperti nilainya - aturan `th` tabel (text-align: left, lebih spesifik) tidak boleh menimpanya.
  it('judul kolom angka rata kanan seperti nilainya', () => {
    expect(ATURAN).toMatch(
      /\.komiteclaimfacin \.komiteclaimfacin__tabel th\.komiteclaimfacin__angka \{[^}]*text-align: right;/,
    )
    const kasus = readFileSync(join(AKAR, 'pages', 'KasusKomite.tsx'), 'utf8')
    expect(kasus).toContain("<th key={i} className={k.jenis === 'angka' ? 'komiteclaimfacin__angka' : undefined}>")
  })

  it('baris rincian (expand pane) tidak ikut zebra / sorot: aturannya sesudah keduanya', () => {
    const rinci = ATURAN.indexOf('tr.komiteclaimfacin__baris-rinci > td {')
    expect(rinci).toBeGreaterThan(ATURAN.indexOf('tbody tr:nth-child(even) > td {'))
    expect(rinci).toBeGreaterThan(ATURAN.indexOf('tbody tr:hover > td {'))
  })

  it('warna hanya di blok token akar; nol backdrop-filter / filter / transform (Modal View Retro tanpa portal)', () => {
    const tanpaToken = ATURAN.replace(/(^|\n)[^{}\n]*\.komiteclaimfacin__akar \{[^}]*\}/g, '')
    expect(tanpaToken).not.toMatch(/#[0-9a-fA-F]{3,8}\b|rgba?\(|hsla?\(/)
    expect(ATURAN).not.toMatch(
      /(^|[\s;{])(-webkit-)?(backdrop-filter|filter|transform|translate|rotate|scale|perspective|will-change|contain)\s*:/m,
    )
  })
})
