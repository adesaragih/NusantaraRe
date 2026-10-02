import { describe, expect, it } from 'vitest'

import { badanPremiCargo } from './api'

describe('badanPremiCargo', () => {
  it('angka tetap teks, lini MARINE CARGO, bukan master policy', () => {
    const b = badanPremiCargo({ mataUang: ' IDR ', rate: '0.125', tsi: '1000000.10' })
    expect(b).toEqual({ liniBisnis: 'MARINE CARGO', mataUang: 'IDR', tsi: '1000000.10', rate: '0.125', masterPolicy: false })
    expect(typeof b.tsi).toBe('string')
    expect(typeof b.rate).toBe('string')
  })

  it('tidak membulatkan atau menormalkan angka', () => {
    // 0.1 + 0.2 di float = 0.30000000000000004; teks harus utuh.
    expect(badanPremiCargo({ mataUang: 'IDR', rate: '0.30', tsi: '0000123' }).rate).toBe('0.30')
    expect(badanPremiCargo({ mataUang: 'IDR', rate: '1', tsi: '0000123' }).tsi).toBe('0000123')
  })
})
