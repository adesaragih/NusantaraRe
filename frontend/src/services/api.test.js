import { describe, expect, it } from 'vitest'

import { jumlahUang, tampilUang } from './api.js'

// Tiket 01 AC-4: tidak ada nilai uang sebagai binary floating point di lapisan
// mana pun - termasuk di React. Backend mengirimnya sebagai teks desimal;
// penjaga di bawah memastikan batas itu tidak pernah diseberangi diam-diam.
describe('uang di sisi React', () => {
  it('membawa jumlah sebagai teks, apa adanya', () => {
    const uang = { amount: '1234567890.12345678', currency: 'IDR' }
    expect(jumlahUang(uang)).toBe('1234567890.12345678')
    expect(typeof jumlahUang(uang)).toBe('string')
  })

  it('tidak kehilangan presisi pada nilai yang tidak muat di float64', () => {
    // Number('1234567890.12345678') membulat menjadi 1234567890.1234567.
    const teks = '1234567890.12345678'
    expect(jumlahUang({ amount: teks, currency: 'IDR' })).toBe(teks)
    expect(String(Number(teks))).not.toBe(teks)
  })

  it('menolak jumlah yang datang sebagai angka JSON', () => {
    expect(() => jumlahUang({ amount: 1234.5, currency: 'IDR' })).toThrow(TypeError)
  })

  it('memperlakukan kosong sebagai kosong, bukan nol', () => {
    expect(jumlahUang(null)).toBe('')
    expect(jumlahUang(undefined)).toBe('')
    expect(jumlahUang({ currency: 'IDR' })).toBe('')
    expect(tampilUang(null)).toBe('')
  })

  it('menampilkan jumlah bersama mata uangnya', () => {
    expect(tampilUang({ amount: '250000', currency: 'IDR' })).toBe('250000 IDR')
  })
})
