// Tampilan polis NB Marine Cargo (work owner 11-10-2026 "cek tampilan nya di nb facin!"): dua bentuk halaman (Pega
// `.PolicyData` / tabel datar), Grand Total `totalGrossCargo`, spreading `totalCargo` nbfacin.

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import { hitungSpreadingCargoNB, kapalTambahanNB, medanCargo, totalGrossCargoNB } from './polisNBCargo'

const C = 'OfferFacIn.CargoList'

describe('polis Marine Cargo', () => {
  it('medan pelayaran dibaca dari baris datar atau .PolicyData Pega', () => {
    expect(medanCargo({ SailDate: '20260105' }, 'SailDate')).toBe('20260105')
    expect(medanCargo({ 'PolicyData.SailDate': '20260106' }, 'SailDate')).toBe('20260106')
    expect(medanCargo({ 'PolicyData.Ship.NM_SHIP': 'UJI KAPAL' }, 'Ship.NM_SHIP')).toBe('UJI KAPAL')
    const h = {
      nilai: {},
      daftar: { [C + '(1).PolicyData.Ship.AdditionalShip']: [{ NM_SHIP: 'UJI KAPAL 2' }] },
    } as Halaman
    expect(kapalTambahanNB(h, 1)).toEqual([{ NM_SHIP: 'UJI KAPAL 2' }])
  })

  it('Grand Total per mata uang coverage (persen)', () => {
    expect(
      totalGrossCargoNB([
        { 'Currency.Name': 'IDR', TSI: '1000', Premium: '5' },
        { 'Currency.Name': 'USD', TSI: '0', Premium: '1' },
        { 'Currency.Name': 'IDR', TSI: '1000', Premium: '15' },
      ]),
    ).toEqual([
      { currency: 'IDR', tsi: '2000', premium: '20', rate: '1' },
      { currency: 'USD', tsi: '0', premium: '1', rate: '0' },
    ])
  })

  it('spreading: LoL hanya cargo pertama, cargo FlagDelete dilewati, top risk dari cargo Top Risk', () => {
    const sp = (n: number) => `${C}(${n}).CoverageList(1).SpreadingList`
    const h = {
      nilai: { 'OfferFacIn.PercentShare': '50' },
      daftar: {
        [C]: [{ IsTopRisk: 'true' }, {}, { FlagDelete: '1' }],
        [C + '(1).CoverageList']: [{ 'Currency.Name': 'IDR', LimitofLiability: '1000' }],
        [C + '(2).CoverageList']: [{ 'Currency.Name': 'IDR', LimitofLiability: '1000' }],
        [C + '(3).CoverageList']: [{ 'Currency.Name': 'IDR', LimitofLiability: '1000' }],
        [sp(1)]: [
          {
            TreatyType: '10003',
            TreatyName: 'UJI QS',
            SharePercentage: '60',
            TSISpreaded: '600',
            PremiumSpreaded: '6',
          },
        ],
        [sp(2)]: [
          {
            TreatyType: '10003',
            TreatyName: 'UJI QS',
            SharePercentage: '60',
            TSISpreaded: '300',
            PremiumSpreaded: '3',
          },
        ],
        [sp(3)]: [
          {
            TreatyType: '10003',
            TreatyName: 'UJI QS',
            SharePercentage: '60',
            TSISpreaded: '999',
            PremiumSpreaded: '9',
          },
        ],
      },
    } as Halaman
    const s = hitungSpreadingCargoNB(h)
    expect(s.treaty).toHaveLength(1)
    expect(s.treaty[0]).toMatchObject({
      currency: 'IDR',
      treatyName: 'UJI QS',
      sharePercentage: '60',
      tsiSpreaded: '900',
      premiumSpreaded: '9',
      claimSpreaded: '600',
      claimEstimation: '300', // 1000 x 50% x 60%, cargo pertama saja
    })
  })
})
