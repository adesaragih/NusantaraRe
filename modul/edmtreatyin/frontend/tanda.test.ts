// Asal: salinan `modul/nbtreatyin/frontend/tanda.test.ts` (06-10-2026), kelas `nbti__` -> `edmt__`.
// Tanda angka teks tanpa Number (spec §5.6: uang dan persentase tidak pernah float).

import { describe, expect, it } from 'vitest'

import { negatifTeks, tandaTeks } from './tanda'

describe('tandaTeks', () => {
  it('membaca tanda dari teks desimal apa adanya', () => {
    expect(tandaTeks('')).toBe(0)
    expect(tandaTeks(' 0 ')).toBe(0)
    expect(tandaTeks('-0.00')).toBe(0)
    expect(tandaTeks('.5')).toBe(1)
    expect(tandaTeks('+12.50')).toBe(1)
    expect(tandaTeks('-0.000000000000000000001')).toBe(-1)
    expect(tandaTeks('99999999999999999999999.99')).toBe(1)
    expect(tandaTeks('1,5')).toBeNull()
    expect(tandaTeks('UJI')).toBeNull()
  })

  it('negatif', () => {
    expect(negatifTeks('-0.01')).toBe(true)
    expect(negatifTeks('UJI-1')).toBe(false)
  })
})
