// Tampilan angka 4 desimal seragam (work owner 06-10-2026) - data uji sintetis.

import { describe, expect, it } from 'vitest'

import { angka4 } from '../angka4'

describe('angka4', () => {
  it('selalu 4 angka di belakang koma; dibulatkan setengah ke atas; ribuan bertitik', () => {
    expect(angka4('20000000')).toBe('20.000.000,0000')
    expect(angka4('0.123456')).toBe('0,1235')
    expect(angka4('1500000.12345')).toBe('1.500.000,1235')
    expect(angka4('33.33333333333333333333')).toBe('33,3333')
    expect(angka4('-2.5')).toBe('-2,5000')
  })
  it('kosong -> kosong; bukan angka -> apa adanya', () => {
    expect(angka4('')).toBe('')
    expect(angka4(undefined)).toBe('')
    expect(angka4('abc')).toBe('abc')
  })
})
