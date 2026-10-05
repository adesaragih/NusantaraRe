// `Section/ListSuggest` grid `.SuggestList`: pyPageMode Numeric, pyPageSize Other,
// pyPageSizeOther 5 (audit silang P3).

import { describe, expect, it } from 'vitest'

import { BARIS_PER_HALAMAN_USULAN, halamanTerjepit, irisan, jumlahHalaman } from './paginasi'

describe('paginasi grid ListSuggest', () => {
  const baris = Array.from({ length: 12 }, (_, i) => i + 1)
  it('lima baris per halaman (pyPageSizeOther)', () => {
    expect(BARIS_PER_HALAMAN_USULAN).toBe(5)
    expect(jumlahHalaman(12, 5)).toBe(3)
    expect(jumlahHalaman(0, 5)).toBe(1)
    expect(irisan(baris, 1, 5)).toEqual([1, 2, 3, 4, 5])
    expect(irisan(baris, 3, 5)).toEqual([11, 12])
  })
  it('halaman di luar rentang dijepit', () => {
    expect(irisan(baris, 9, 5)).toEqual([11, 12])
    expect(irisan(baris, 0, 5)).toEqual([1, 2, 3, 4, 5])
    expect(halamanTerjepit(9, 3)).toBe(3)
    expect(halamanTerjepit(0, 3)).toBe(1)
  })
})
