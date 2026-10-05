// Kepala grid dua baris (perintah work owner 05-10-2026 "grouping sama seperti bordereaux"): pasangan NoR/IA per
// kategori, Total, dan RNM digabung di baris pertama; kolom lain menempati dua baris.

import { describe, expect, it } from 'vitest'

import { KOLOM_GRID, susunKepala } from './aturan'
import { GRUP_KOLOM, LABEL_KOLOM } from './labels'

describe('kepala grid Aggregate', () => {
  const { baris1, baris2 } = susunKepala(KOLOM_GRID, GRUP_KOLOM, LABEL_KOLOM)

  it('setiap kolom tampil tepat sekali: tanpa grup di baris pertama, bergrup di baris kedua', () => {
    const tunggal = baris1.filter((s) => !s.grup).map((s) => s.kunci)
    expect([...tunggal, ...baris2.map((s) => s.kunci)].sort()).toEqual(KOLOM_GRID.map((k) => k.nama).sort())
    expect(baris1.filter((s) => !s.grup).every((s) => s.rowSpan === 2)).toBe(true)
  })

  it('grup berurutan: 11 kategori NoR/IA, lalu Total dan RNM', () => {
    const grup = baris1.filter((s) => s.grup)
    expect(grup.map((s) => `${s.teks}:${s.colSpan}`)).toEqual([
      'Buildings:2',
      'Stocks:2',
      'Machinery:2',
      'Other Contents:2',
      'Consequential Loss:2',
      'Residential:2',
      'Commercial:2',
      'Industrial:2',
      'Agriculture:2',
      'Miscellaneous:2',
      'Utilities:2',
      'Total:3',
      'RNM:3',
    ])
    expect(baris2.slice(0, 2).map((s) => [s.teks, s.judul])).toEqual([
      ['NoR', 'NoR Buildings'],
      ['IA', 'IA Buildings'],
    ])
    expect(baris1.at(-1)?.teks).toBe('Remark')
  })

  it('setiap kolom bergrup memang kolom grid', () => {
    for (const nama of Object.keys(GRUP_KOLOM))
      expect(
        KOLOM_GRID.some((k) => k.nama === nama),
        nama,
      ).toBe(true)
  })
})
