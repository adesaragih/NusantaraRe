// Penjaga gaya modul: kelas CSS khusus modul BERAWALAN `edm` dan tinggal di berkas CSS modul sendiri
// (`endorsementlife.css`); kelas bersama (`panel`, `btn`, `inbox__tabel`, ...) dari `inti` tidak
// didefinisikan ulang di sini.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname

/** Setiap berkas di bawah folder frontend modul (rekursif). */
function berkas(d = AKAR): string[] {
  return readdirSync(d, { withFileTypes: true }).flatMap((e) => (e.isDirectory() ? berkas(join(d, e.name)) : [join(d, e.name)]))
}

const CSS = readFileSync(join(AKAR, 'endorsementlife.css'), 'utf8')

/** Kelas yang didefinisikan pemilih di `endorsementlife.css`. */
function kelasCSS(): Set<string> {
  const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
  const pemilih = [...tanpaKomentar.matchAll(/([^{}]+)\{/g)].map((m) => m[1] ?? '')
  return new Set(pemilih.flatMap((p) => [...p.matchAll(/\.([a-zA-Z0-9_-]+)/g)].map((m) => m[1] ?? '')))
}

const POLA_KELAS = /\bedm(?:-[a-z0-9]+)+(?:__[a-z0-9]+)?(?:--[a-z0-9]+)?\b/g

/** Kelas berawalan `edm` yang dipakai layar: isi atribut `className` (string atau templat). */
function kelasTSX(): Set<string> {
  const isi = berkas()
    .filter((f) => /\.tsx?$/.test(f) && !f.endsWith('.test.ts'))
    .map((f) => readFileSync(f, 'utf8'))
  const potongan = isi.flatMap((s) => [
    ...[...s.matchAll(/className="([^"]*)"/g)].map((m) => m[1] ?? ''),
    ...[...s.matchAll(/className=\{`([^`]*)`\}/g)].map((m) => m[1] ?? ''),
  ])
  return new Set(potongan.flatMap((p) => [...p.matchAll(POLA_KELAS)].map((m) => m[0])))
}

describe('kelas CSS modul', () => {
  it('satu-satunya berkas CSS modul adalah endorsementlife.css, dan halaman mengimpornya', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['endorsementlife.css'])
    expect(readFileSync(join(AKAR, 'pages', 'InboxEndorsementLife.tsx'), 'utf8')).toContain("import '../endorsementlife.css'")
  })

  it('setiap kelas di endorsementlife.css berawalan edm', () => {
    const kelas = [...kelasCSS()]
    expect(kelas.length).toBeGreaterThan(3)
    expect(kelas.filter((k) => !k.startsWith('edm'))).toEqual([])
  })

  it('setiap kelas edm yang dipakai layar didefinisikan di endorsementlife.css', () => {
    const didefinisikan = kelasCSS()
    const dipakai = [...kelasTSX()]
    expect(dipakai.length).toBeGreaterThan(3)
    expect(dipakai.filter((k) => !didefinisikan.has(k))).toEqual([])
  })
})
