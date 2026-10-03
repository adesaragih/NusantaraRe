// Penjaga "CSS modul tidak menumpang di inti" (work owner 02-10-2026: "untuk css style per modul tidak
// ada lagi menumpang per inti").
//
// Kelas yang didefinisikan `inti/frontend/styles.css` tetapi HANYA dipakai layar satu modul (dan tidak
// dipakai inti maupun perakit `frontend/`) adalah gaya milik modul itu: tempatnya `modul/<nama>/frontend/
// <nama>.css` berawalan kelas akar modul, bukan inti. Kelas yang dipakai dua modul atau lebih (`btn`,
// `panel`, `inbox__tabel`, `polis__catatan`, ...) tetap komponen bersama inti.
//
// Pemakaian dibaca dari atribut `className` (teks, templat, atau ekspresi tanpa kurung kurawal bersarang).

import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const APLIKASI = join(__dirname, '..', '..')
const CSS_INTI = readFileSync(join(__dirname, 'styles.css'), 'utf8')

/** Kelas yang disebut pemilih sebuah lembar gaya (tanpa komentar; isi @media ikut). */
function kelasDidefinisikan(css: string): Set<string> {
  const tanpaKomentar = css.replace(/\/\*[\s\S]*?\*\//g, '')
  const hasil = new Set<string>()
  for (const m of tanpaKomentar.matchAll(/([^{};]+)\{/g)) {
    const p = (m[1] ?? '').trim()
    if (p.startsWith('@')) continue
    for (const k of p.matchAll(/\.([a-zA-Z][\w-]*)/g)) hasil.add(k[1] ?? '')
  }
  return hasil
}

/** Kelas yang disebut atribut `className` di kode sumber (bukan berkas uji). */
function kelasDipakai(kode: string): Set<string> {
  const hasil = new Set<string>()
  const isi = [
    ...[...kode.matchAll(/className="([^"]*)"/g)].map((m) => m[1] ?? ''),
    ...[...kode.matchAll(/className=\{`([^`]*)`\}/g)].map((m) => m[1] ?? ''),
    ...[...kode.matchAll(/className=\{([^{}]*)\}/g)].map((m) => m[1] ?? ''),
  ]
  for (const s of isi) for (const k of s.matchAll(/[a-zA-Z][\w-]*/g)) hasil.add(k[0])
  return hasil
}

function berkasKode(d: string): string[] {
  if (!existsSync(d)) return []
  return readdirSync(d).flatMap((n) => {
    const f = join(d, n)
    if (statSync(f).isDirectory()) return n === 'node_modules' || n === 'dist' ? [] : berkasKode(f)
    return /\.tsx?$/.test(n) && !/\.test\.tsx?$/.test(n) ? [f] : []
  })
}

function pemakaiDari(d: string): Set<string> {
  const hasil = new Set<string>()
  for (const f of berkasKode(d)) for (const k of kelasDipakai(readFileSync(f, 'utf8'))) hasil.add(k)
  return hasil
}

/** Kelas inti yang hanya dipakai SATU modul: { modul: [kelas] }. */
function menumpang(didefinisikan: Set<string>, bersama: Set<string>, perModul: Record<string, Set<string>>): Record<string, string[]> {
  const hasil: Record<string, string[]> = {}
  for (const k of didefinisikan) {
    if (bersama.has(k)) continue
    const siapa = Object.keys(perModul).filter((m) => perModul[m]?.has(k))
    if (siapa.length === 1) (hasil[siapa[0] ?? ''] ??= []).push(k)
  }
  for (const m of Object.keys(hasil)) hasil[m]?.sort()
  return hasil
}

describe('CSS modul tidak menumpang di inti', () => {
  it('tidak ada kelas di styles.css inti yang hanya dipakai satu modul', () => {
    const bersama = new Set([...pemakaiDari(__dirname), ...pemakaiDari(join(APLIKASI, 'frontend'))])
    const perModul: Record<string, Set<string>> = {}
    for (const m of readdirSync(join(APLIKASI, 'modul'))) {
      const fe = join(APLIKASI, 'modul', m, 'frontend')
      if (existsSync(fe)) perModul[m] = pemakaiDari(fe)
    }
    expect(Object.keys(perModul).length).toBeGreaterThan(5)
    expect(menumpang(kelasDidefinisikan(CSS_INTI), bersama, perModul)).toEqual({})
  })

  it('aturan menumpang menggigit: kelas satu modul tertangkap, kelas bersama dan kelas dua modul tidak', () => {
    const css = '.btn { a: 1 } .alfa-khusus, .dipakai-dua { b: 2 } @media (max-width: 1px) { .alfa-lagi { c: 3 } }'
    const bersama = kelasDipakai('<button className="btn">')
    const perModul = {
      alfa: kelasDipakai('<div className={`alfa-khusus ${x ? "alfa-lagi" : ""}`}><p className="dipakai-dua btn" />'),
      beta: kelasDipakai('<p className={ok ? "dipakai-dua" : "lain"} />'),
    }
    expect(menumpang(kelasDidefinisikan(css), bersama, perModul)).toEqual({ alfa: ['alfa-khusus', 'alfa-lagi'] })
  })
})
