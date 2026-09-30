// Penjaga gaya modul (brief paket 9): kelas CSS khusus modul BERAWALAN `mcrl` dan tinggal di berkas CSS
// modul sendiri (`mcrl.css`); kelas bersama (`panel`, `btn`, `inbox__tabel`, ...) dari `inti` tidak
// didefinisikan ulang di sini.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname

/** Setiap berkas di bawah folder frontend modul (rekursif). */
function berkas(d = AKAR): string[] {
  return readdirSync(d, { withFileTypes: true }).flatMap((e) => (e.isDirectory() ? berkas(join(d, e.name)) : [join(d, e.name)]))
}

const CSS = readFileSync(join(AKAR, 'mcrl.css'), 'utf8')

/** Kelas yang didefinisikan pemilih di `mcrl.css`. */
function kelasCSS(): Set<string> {
  const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
  const pemilih = [...tanpaKomentar.matchAll(/([^{}]+)\{/g)].map((m) => m[1] ?? '')
  return new Set(pemilih.flatMap((p) => [...p.matchAll(/\.([a-zA-Z0-9_-]+)/g)].map((m) => m[1] ?? '')))
}

const POLA_KELAS = /\bmcrl(?:-[a-z0-9]+)+(?:__[a-z0-9]+)?(?:--[a-z0-9]+)?\b/g

/**
 * Kelas berawalan `mcrl` yang dipakai layar: isi atribut `className="…"` dan literal di badan fungsi
 * perakit kelas `kelas*` (mis. `kelasTotal`). Nama halaman (`mcrl-tahun`) bukan kelas.
 */
function kelasTSX(): Set<string> {
  const isi = berkas()
    .filter((f) => /\.tsx?$/.test(f) && !f.endsWith('.test.ts'))
    .map((f) => readFileSync(f, 'utf8'))
  const potongan = isi.flatMap((s) => [
    ...[...s.matchAll(/className="([^"]*)"/g)].map((m) => m[1] ?? ''),
    ...[...s.matchAll(/function kelas\w*\([^)]*\)[^{]*\{([\s\S]*?)\n\}/g)].map((m) => m[1] ?? ''),
  ])
  return new Set(potongan.flatMap((p) => [...p.matchAll(POLA_KELAS)].map((m) => m[0])))
}

describe('kelas CSS modul', () => {
  it('satu-satunya berkas CSS modul adalah mcrl.css, dan halaman awal mengimpornya', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['mcrl.css'])
    expect(readFileSync(join(AKAR, 'pages', 'MasterContractRetroLife.tsx'), 'utf8')).toContain("import '../mcrl.css'")
  })

  it('setiap kelas di mcrl.css berawalan mcrl', () => {
    const kelas = [...kelasCSS()]
    expect(kelas.length).toBeGreaterThan(5)
    expect(kelas.filter((k) => !k.startsWith('mcrl'))).toEqual([])
  })

  it('setiap kelas mcrl yang dipakai layar didefinisikan di mcrl.css', () => {
    const didefinisikan = kelasCSS()
    const dipakai = [...kelasTSX()]
    expect(dipakai.length).toBeGreaterThan(5)
    expect(dipakai.filter((k) => !didefinisikan.has(k))).toEqual([])
  })
})
