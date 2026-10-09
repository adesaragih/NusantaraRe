// Kotak angka mentah Treaty In menolak huruf/simbol — permintaan pemakai
// 8 Oktober 2026. Lihat `components/saringAngka.ts`.

import { describe, expect, it } from 'vitest'

import { saringAngka } from './components/saringAngka'

describe('saringAngka', () => {
  it('huruf dan simbol dibuang, angka utuh', () => {
    expect(saringAngka('12a3b')).toBe('123')
    expect(saringAngka('Rp 500%')).toBe('500')
    expect(saringAngka('abc')).toBe('')
    expect(saringAngka('')).toBe('')
  })
  it('satu pemisah desimal; koma menjadi titik (bentuk kabel)', () => {
    expect(saringAngka('12,5')).toBe('12.5')
    expect(saringAngka('0.125')).toBe('0.125')
    expect(saringAngka('1.5.')).toBe('1.5')
    expect(saringAngka('12,')).toBe('12.')
  })
  it('minus hanya di depan', () => {
    expect(saringAngka('-7.5')).toBe('-7.5')
    expect(saringAngka('7-5')).toBe('75')
  })
  it('bulat: digit saja, ekor desimal dibuang', () => {
    expect(saringAngka('30 hari', true)).toBe('30')
    expect(saringAngka('10,5', true)).toBe('10')
    expect(saringAngka('-3', true)).toBe('3')
  })
})
