// Tab Premium Deduction / Payment / Cedant tampilan polis NB (work owner 11-10-2026 "cek tampilan nya di nb facin!"):
// nilai tersimpan (keputusan work owner), rincian Cedant per mata uang = rumus nbfacin `ringkasanRNM`.

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import {
  cicilanNB,
  dasarCedantNB,
  itemDasarFireNB,
  mataUangCedantNB,
  pesanPersenCicilanNB,
  shareOfCedingNB,
} from './polisNBBayar'

describe('Payment', () => {
  it('pesan Σ % Installment (SetValidateInstallment_Act)', () => {
    expect(pesanPersenCicilanNB([])).toBe('')
    expect(pesanPersenCicilanNB([{ InstallmentPercentage: '60' }, { InstallmentPercentage: '40.0' }])).toBe('')
    expect(pesanPersenCicilanNB([{ InstallmentPercentage: '60' }, { InstallmentPercentage: '50' }])).toBe(
      "Total % Installment Can't Be More Than 100%",
    )
    expect(pesanPersenCicilanNB([{ InstallmentPercentage: '60' }])).toBe('Total % Installment Must Be 100%')
  })

  it('cicilan dibaca dari bentuk tabel datar atau .Policy.Payment Pega', () => {
    const datar = {
      nilai: {},
      daftar: { 'OfferFacIn.CurrencyList(1).ListInstallment': [{ InstallmentNo: '1' }] },
    } as Halaman
    const pega = {
      nilai: {},
      daftar: { 'OfferFacIn.CurrencyList(1).Policy.Payment.ListInstallment': [{ InstallmentNo: '2' }] },
    } as Halaman
    expect(cicilanNB(datar, 1)).toEqual([{ InstallmentNo: '1' }])
    expect(cicilanNB(pega, 1)).toEqual([{ InstallmentNo: '2' }])
  })
})

describe('Cedant', () => {
  it('"of" ceding: Gross 100%, Share RNM = % Share RNM polis', () => {
    expect(shareOfCedingNB('0', '25', '')).toBe('100%')
    expect(shareOfCedingNB('1', '25.5000', '')).toBe('25.5%')
    expect(shareOfCedingNB('', '25', 'UJI')).toBe('UJI')
  })

  it('dasar per mata uang item FIRE lalu x ShareCeding / 100', () => {
    const L = 'OfferFacIn.LocationList'
    const h = {
      nilai: { 'OfferFacIn.PercentShare': '50' },
      daftar: {
        [L]: [{}],
        [L + '(1).Property.PropertyItemList']: [
          { Currency: 'IDR', TSIObjectItem: '1000' },
          { Currency: 'USD', TSIObjectItem: '10' },
          { Currency: '', TSIObjectItem: '5' },
        ],
        [L + '(1).Property.PropertyItemList(1).CoverageList']: [{ Premium: '100', Discount: '20' }],
        [L + '(1).Property.PropertyItemList(2).CoverageList']: [{ Premium: '4' }],
      },
    } as Halaman
    const dasar = dasarCedantNB(h, itemDasarFireNB(h))
    expect(dasar).toEqual([
      { currency: 'IDR', tsiRnm: '500', premiumRnm: '40', tsiGross: '1000', premiumGross: '80' },
      { currency: 'USD', tsiRnm: '5', premiumRnm: '2', tsiGross: '10', premiumGross: '4' },
    ])
    // Share RNM: dasar RNM; Gross: dasar gross
    expect(mataUangCedantNB('10', '1', dasar)[0]).toEqual({ currency: 'IDR', tsi: '50', premium: '4' })
    expect(mataUangCedantNB('10', '0', dasar)[0]).toEqual({ currency: 'IDR', tsi: '100', premium: '8' })
    expect(mataUangCedantNB('x', '1', dasar)).toEqual([])
  })

  it('tanpa % Share RNM: tanpa dasar', () => {
    expect(dasarCedantNB({ nilai: {}, daftar: {} } as Halaman, [{ currency: 'IDR', tsi: '1', coverage: [] }])).toEqual(
      [],
    )
  })
})
