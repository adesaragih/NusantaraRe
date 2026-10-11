// Desimal teks isian Claim RNM (11-10-2026) — nol float.

import { describe, expect, it } from 'vitest'

import { bandingDesimal, normalisasiDesimal, persenKlaimTeks, saringKetikanAngka } from './desimal'

describe('desimal teks', () => {
  it('ketikan: hanya digit dan satu pemisah desimal', () => {
    expect(saringKetikanAngka('12a3,4.5x')).toBe('123,45')
    expect(saringKetikanAngka('1.000.000')).toBe('1.000000')
    expect(saringKetikanAngka('abc')).toBe('')
  })
  it('normalisasi: koma = titik, nol di depan/ekor dibuang; bukan angka = null', () => {
    expect(normalisasiDesimal('0050,500')).toBe('50.5')
    expect(normalisasiDesimal('7.')).toBe('7')
    expect(normalisasiDesimal('.5')).toBe('0.5')
    expect(normalisasiDesimal('')).toBeNull()
    expect(normalisasiDesimal('1e3')).toBeNull()
    expect(normalisasiDesimal('-1')).toBeNull()
  })
  it('bandingkan lewat digit, presisi penuh (tanpa float)', () => {
    expect(bandingDesimal('5077429.5', '5077429.50')).toBe(0)
    expect(bandingDesimal('5077429.51', '5077429.5')).toBe(1)
    expect(bandingDesimal('999', '1000')).toBe(-1)
    // Di luar presisi float64: 2^53 + 1 vs 2^53.
    expect(bandingDesimal('9007199254740993', '9007199254740992')).toBe(1)
    expect(bandingDesimal('0.1000000000000000001', '0.1')).toBe(1)
  })
})

// Kesamaan dengan backend models.HitungPersenKlaim — kasus yang SAMA ada di persenklaim_test.go (hasil Go diperiksa
// 11-10-2026: pecahan dibulatkan setengah-ke-atas pada 5 desimal, lalu × 100).
describe('persen klaim identik dengan backend', () => {
  it.each([
    ['50', '100', '50'],
    ['1', '3', '33.333'],
    ['0', '100', '0'],
    ['2', '3', '66.667'],
    ['665', '5850038.75', '0.011'],
    ['1', '200000', '0.001'],
    ['1', '400000', '0'],
    ['94695.705', '100000', '94.696'],
  ])('%s / %s = %s', (a, b, mau) => {
    expect(persenKlaimTeks(a, b)).toBe(mau)
  })
  it('share nol = kosong', () => {
    expect(persenKlaimTeks('5', '0')).toBe('')
  })
})
