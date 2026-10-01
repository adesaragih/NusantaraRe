// Penjaga gaya modul (brief bab 1): kelas CSS khusus modul BERAWALAN `mpnl` dan tinggal di berkas CSS modul
// sendiri (`masterproductnamelife.css`, diimpor `rute.tsx`); kelas bersama (`panel`, `btn`, `inbox__tabel`, ...)
// dari `inti` tidak didefinisikan ulang di sini.

import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const AKAR = __dirname

/** Setiap berkas di bawah folder frontend modul (rekursif). */
function berkas(d = AKAR): string[] {
  return readdirSync(d, { withFileTypes: true }).flatMap((e) => (e.isDirectory() ? berkas(join(d, e.name)) : [join(d, e.name)]))
}

const CSS = readFileSync(join(AKAR, 'masterproductnamelife.css'), 'utf8')

/** Kelas yang didefinisikan pemilih di CSS modul. */
function kelasCSS(): Set<string> {
  const tanpaKomentar = CSS.replace(/\/\*[\s\S]*?\*\//g, '')
  const pemilih = [...tanpaKomentar.matchAll(/([^{}]+)\{/g)].map((m) => m[1] ?? '')
  return new Set(pemilih.flatMap((p) => [...p.matchAll(/\.([a-zA-Z0-9_-]+)/g)].map((m) => m[1] ?? '')))
}

const POLA_KELAS = /\bmpnl(?:-[a-z0-9]+)+(?:__[a-z0-9]+)?(?:--[a-z0-9]+)?\b/g

/** Kelas berawalan `mpnl` yang dipakai layar (isi `className`); nama halaman (`mpnl-produk`) bukan kelas. */
function kelasTSX(): Set<string> {
  const isi = berkas()
    .filter((f) => /\.tsx$/.test(f))
    .map((f) => readFileSync(f, 'utf8'))
  const potongan = isi.flatMap((s) => [
    ...[...s.matchAll(/className=\{?["'`]([^"'`]*)["'`]/g)].map((m) => m[1] ?? ''),
    ...[...s.matchAll(/'(mpnl-[a-z0-9_-]+(?: [a-z0-9_ -]+)?)'/g)].map((m) => m[1] ?? ''),
    ...[...s.matchAll(/' (mpnl-[a-z0-9_-]+)'/g)].map((m) => m[1] ?? ''),
  ])
  return new Set(potongan.flatMap((p) => [...p.matchAll(POLA_KELAS)].map((m) => m[0])))
}

describe('kelas CSS modul', () => {
  it('satu-satunya berkas CSS modul adalah masterproductnamelife.css, dan rute.tsx mengimpornya', () => {
    expect(berkas().filter((f) => f.endsWith('.css')).map((f) => f.slice(AKAR.length + 1))).toEqual(['masterproductnamelife.css'])
    expect(readFileSync(join(AKAR, 'rute.tsx'), 'utf8')).toContain("import './masterproductnamelife.css'")
  })

  it('setiap kelas di CSS modul berawalan mpnl', () => {
    const kelas = [...kelasCSS()]
    expect(kelas.length).toBeGreaterThan(5)
    expect(kelas.filter((k) => !k.startsWith('mpnl'))).toEqual([])
  })

  it('setiap kelas mpnl yang dipakai layar didefinisikan di CSS modul', () => {
    const didefinisikan = kelasCSS()
    const dipakai = [...kelasTSX()].filter((k) => k !== 'mpnl-produk')
    expect(dipakai.length).toBeGreaterThan(5)
    expect(dipakai.filter((k) => !didefinisikan.has(k))).toEqual([])
  })
})
