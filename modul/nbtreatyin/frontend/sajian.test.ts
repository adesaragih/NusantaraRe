// Uji pemformat penyajian layar realisasi (K14, OQ 7, AC 86). Nilai harapan
// dihitung tangan dari setelan sel Section Pega (`pyModes` mode baca):
//   pxNumber `pyDecimalPlaces` 2  - GrossPremium, RiComm*, OveriddingComm*, PaymentTotal
//   pxNumber `pyDecimalPlaces` 4  - grid spreading, InstallmentPercentage, Premium angsuran
//   pxTextInput `pyFormatType number` `pyDecimalPlaces 0` `pySeparators false` - Quartal atasan
//   pxCurrency tanpa `pyDecimalPlaces`, Balance* `pyFormatType number` tanpa desimal
//     -> tak terbaca: pola `inti/frontend/lib/format.ts` (titik ribuan, koma desimal,
//        tanpa batas desimal, nol ekor dibuang)

import { describe, expect, it } from 'vitest'

import { sajikan, sajikanAngka } from './sajian'

describe('sajikanAngka - pyDecimalPlaces terbaca: tepat N desimal', () => {
  it('2 desimal: dipadankan dan dibulatkan setengah ke atas, titik ribuan koma desimal', () => {
    expect(sajikanAngka('1234567.5', { desimal: 2 })).toBe('1.234.567,50')
    expect(sajikanAngka('1000', { desimal: 2 })).toBe('1.000,00')
    expect(sajikanAngka('0.005', { desimal: 2 })).toBe('0,01')
    expect(sajikanAngka('-1234.567', { desimal: 2 })).toBe('-1.234,57')
  })

  it('4 desimal: grid spreading dan angsuran', () => {
    expect(sajikanAngka('33.3333333333', { desimal: 4 })).toBe('33,3333')
    expect(sajikanAngka('100', { desimal: 4 })).toBe('100,0000')
  })

  it('0 desimal tanpa pemisah ribuan (pySeparators false): tahun tidak menjadi "2.026"', () => {
    expect(sajikanAngka('2026', { desimal: 0, ribuan: false })).toBe('2026')
    expect(sajikanAngka('4', { desimal: 0, ribuan: false })).toBe('4')
  })

  it('nol negatif hasil pembulatan tampil tanpa tanda', () => {
    expect(sajikanAngka('-0.001', { desimal: 2 })).toBe('0,00')
  })
})

describe('sajikanAngka - jumlah desimal tak terbaca: pola inti', () => {
  it('tanpa batas desimal, nol ekor dibuang, minus di depan (pyNegativeFormat minusStyle)', () => {
    expect(sajikanAngka('830.82191780804', {})).toBe('830,82191780804')
    expect(sajikanAngka('2484250.000000', {})).toBe('2.484.250')
    expect(sajikanAngka('-1500', {})).toBe('-1.500')
  })

  it('teks bukan angka dan kosong ditampilkan apa adanya', () => {
    expect(sajikanAngka('UJI-TEKS', { desimal: 2 })).toBe('UJI-TEKS')
    expect(sajikanAngka('', { desimal: 2 })).toBe('')
  })
})

describe('sajikan - per jenis sajian sel', () => {
  it('tanggal: satu format seluruh sistem (AC 33, formatDate inti)', () => {
    expect(sajikan('2026-10-03 09:00:00', 'tanggal')).toBe('03-10-2026')
    expect(sajikan('', 'tanggal')).toBe('')
  })

  it('tanpa sajian: nilai apa adanya (teks)', () => {
    expect(sajikan('2026', undefined)).toBe('2026')
  })

  it('angka: diteruskan ke sajikanAngka', () => {
    expect(sajikan('1000', { desimal: 2 })).toBe('1.000,00')
  })
})
