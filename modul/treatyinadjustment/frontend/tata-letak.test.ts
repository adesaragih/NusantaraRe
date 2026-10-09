// Tata letak layout Pega (`pyLayoutOtherFormat`) — permintaan pemakai
// 8 Oktober 2026: posisi tiap tab sama dengan Pega. Contoh yang dipaku:
// `Total All Layers RNM Share` tab Share Non-Prop = lima `Inline grid triple`
// bertumpuk (gambar Pega pemakai): baris 1 dan 5 tiga grid, baris 2–4 SATU
// grid di kolom pertama.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { BlokKerangka, ButirKerangka } from './ekspor/jenis'
import { KERANGKA_RINCIAN, KERANGKA_TAB } from './ekspor/kerangka.gen'

function cariBlok(bs: readonly ButirKerangka[], judul: string): BlokKerangka | undefined {
  for (const b of bs) {
    if (b.t !== 'blok') continue
    if (b.judul === judul) return b
    const x = cariBlok(b.anak, judul)
    if (x !== undefined) return x
  }
  return undefined
}

/** Blok bertata `tata` paling dangkal di bawah `b`. */
function blokTata(bs: readonly ButirKerangka[], tata: string): BlokKerangka[] {
  const out: BlokKerangka[] = []
  for (const b of bs) {
    if (b.t !== 'blok') continue
    if (b.tata === tata) out.push(b)
    else out.push(...blokTata(b.anak, tata))
  }
  return out
}

/** Larik grid di bawah blok (rekursif — `Total Deduction` berada di blok anak berjudul tersembunyi). */
const grid = (b: BlokKerangka): string[] => b.anak.flatMap((x) => (x.t === 'grid' ? [x.larik] : x.t === 'blok' ? grid(x) : []))

describe('tata letak Pega — Total All Layers RNM Share', () => {
  it('lima baris Inline grid triple; kolom sesuai gambar Pega', () => {
    const isi = KERANGKA_TAB['TreatyInTabsNonProportional#Share']?.isi ?? []
    const total = cariBlok(isi, 'Total All Layers RNM Share')
    expect(total).toBeDefined()
    const baris = blokTata(total?.anak ?? [], 'g3')
    expect(baris.map(grid)).toEqual([
      ['TotalShareRnmNP', 'TotalSpreadedRnmProp', 'TotalSpreadedRnmRIProp'],
      ['TotalShareGrossMinNP'],
      ['TotalShareGrossNP'],
      ['TotalShareDeductionNP'],
      ['TotalShareNetNP', 'TotalSpreadedNetPremi', 'TotalSpreadedNetPremiRI'],
    ])
  })

  it('rincian DetailLimits = gambar Pega 05: Treaty Group | slot Spacer, grid CoB baris kedua; tiga baris [label | grid] 30/70', () => {
    const isi = KERANGKA_RINCIAN.DetailLimits ?? []
    const g2 = isi[0]
    expect(g2?.t === 'blok' ? [g2.tata, g2.anak.map((x) => x.t)] : null).toEqual(['g2', ['blok', 'kosong', 'grid']])
    const pasangan = blokTata(isi, 't3070')[0]
    expect(pasangan?.anak.map((x) => (x.t === 'grid' ? x.larik : x.t))).toEqual([
      'blok', 'IOOLimitList', 'blok', 'RetentionList', 'blok', 'CessionList',
    ])
  })

  it('CSS memberi tiap tata letak susunannya, dan renderer memasang kelasnya', () => {
    const css = readFileSync(join(__dirname, 'treatyinadjustment.css'), 'utf8')
    // ⚠️ Yang dikunci adalah NIATNYA — `g3` TIGA kolom — bukan ejaannya.
    // Bentuk `repeat(3, minmax(0, 1fr))` dulu dipakai dan ia MELUBER: sel
    // boleh menyusut sampai nol, medan berlabel kiri di dalamnya tidak, dan
    // yang tidak muat menimpa kolom sebelah (keluhan 8 Oktober 2026).
    // Penggantinya membungkus, dan pembagi `/ 3` itulah yang menjaga
    // kolomnya tidak pernah lebih dari tiga.
    expect(css).toMatch(/\.tria__blok\.tria__tata--g3 \{\s*grid-template-columns:[^;]*\/ 3\)/)
    expect(css).toContain('.tria__blok.tria__tata--kiri > .field {')
    const tsx = readFileSync(join(__dirname, 'komponen', 'KerangkaTab.tsx'), 'utf8')
    expect(tsx).toContain('`tria__blok tria__tata--${b.tata}`')
  })
})

// Judul grid Limits Prop (100% Limit / Retention / Cession to R/I) — bentuk
// sama dengan Treaty In (permintaan pemakai 8 Oktober 2026): sel kiri pasangan
// 30/70 = blok kiri > g2 [judul | alir [persen, %]]; persennya teks SEBARIS.
describe('judul grid Limits sama dengan Treaty In', () => {
  const css = readFileSync(join(__dirname, 'treatyinadjustment.css'), 'utf8')
  const SEL = '.tria__blok.tria__tata--t3070 > .tria__tata--kiri > .tria__blok.tria__tata--g2 {'

  it('kerangka: judul + persen baca-saja di sel kiri pasangan 30/70', () => {
    const teks = readFileSync(join(__dirname, 'ekspor', 'kerangka.gen.ts'), 'utf8').replace(/\s+/g, '')
    expect(teks).toContain('"teks":"Retention"')
    expect(teks).toContain('"kunci":"RetentionPct"')
  })

  it('simetris: judul | kolom persen (angka rata kanan 6ch), persen tanpa kotak', () => {
    const i = css.indexOf(SEL)
    expect(i).toBeGreaterThan(0)
    const aturan = css.slice(i, css.indexOf('}', i))
    expect(aturan).toContain('grid-template-columns: minmax(min-content, 1fr) auto')
    expect(css).toContain('grid-template-columns: 6ch auto')
    const j = css.indexOf('> .tria__tata--alir > .field > .field__input {', i)
    const kotak = css.slice(j, css.indexOf('}', j))
    expect(kotak).toContain('border: 0')
    expect(kotak).toContain('background: transparent')
  })
})
