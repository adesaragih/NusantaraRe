// Penjaga gaya modul NB FacIn (tiket 25/26) - pola `treatycontractout/frontend/gaya.test.ts`:
// satu berkas CSS modul (`nbfacin.css`, diimpor `rute.tsx`), setiap pemilih di bawah kelas akar
// `.nbfacin`, kelas modul berawalan `nbf-`.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname
const CSS = readFileSync(join(AKAR, 'nbfacin.css'), 'utf8')

function berkas(d = AKAR): string[] {
  return readdirSync(d, { withFileTypes: true }).flatMap((e) => (e.isDirectory() ? berkas(join(d, e.name)) : [join(d, e.name)]))
}

/** Pemilih CSS (tanpa komentar; isi @media ikut, kepala @media tidak). */
function pemilih(css: string): string[] {
  return [...css.replace(/\/\*[\s\S]*?\*\//g, '').matchAll(/([^{};]+)\{/g)]
    .map((m) => (m[1] ?? '').trim())
    .filter((p) => p !== '' && !p.startsWith('@'))
    .flatMap((p) => p.split(',').map((s) => s.trim()))
}

describe('gaya NB FacIn terisolasi', () => {
  it('satu-satunya berkas CSS modul, diimpor rute.tsx', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['nbfacin.css'])
    expect(readFileSync(join(AKAR, 'rute.tsx'), 'utf8')).toContain("import './nbfacin.css'")
  })

  it('setiap pemilih di bawah .nbfacin', () => {
    const p = pemilih(CSS)
    expect(p.length).toBeGreaterThan(0)
    expect(p.filter((s) => s !== '.nbfacin' && !s.startsWith('.nbfacin '))).toEqual([])
  })

  it('kelas yang didefinisikan berawalan nbf- (selain akar)', () => {
    const kelas = [...new Set(pemilih(CSS).flatMap((s) => [...s.matchAll(/\.([\w-]+)/g)].map((m) => m[1])))]
    expect(kelas.filter((k) => k !== 'nbfacin' && !k!.startsWith('nbf-'))).toEqual([])
  })

  it('uji ini menggigit: pemilih global tertangkap', () => {
    expect(pemilih('.btn { color: red }').filter((s) => !s.startsWith('.nbfacin'))).toEqual(['.btn'])
  })
})
