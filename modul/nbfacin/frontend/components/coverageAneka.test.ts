// Fungsi murni tab Coverage Aneka (tiket 73 P2). Data uji sintetis.

import { describe, expect, it } from 'vitest'

import type { CoverageAneka, Deductible, LokasiAneka, ObjekAneka } from '../api'
import { lokasiAnekaBaru, objekAnekaBaru } from './aneka'
import {
  coverageAnekaBaru,
  dorongTsiObjek,
  metodeAwal,
  salinCoverageSeksi,
  salinDeductibleAneka,
  setelCurrencyDeductibleTerakhir,
  tambahCoverage,
  ubahCoverageDi,
} from './coverageAneka'

const ded = (currency: string, amount = '1'): Deductible => ({
  typeDeductible: '1', pctDeductible: '', minMax: '', currency, typeDeductible2: '', pctDeductible2: '', condition: '', amount,
  inputCondition: '', timeExcess: '',
})

const cov = (p: Partial<CoverageAneka>): CoverageAneka => ({ ...coverageAnekaBaru({ currency: 'IDR', tsi: '100' }, '1'), ...p })
const obj = (p: Partial<ObjekAneka>): ObjekAneka => ({ ...objekAnekaBaru(), ...p })

function pohon(): LokasiAneka[] {
  return [
    {
      ...lokasiAnekaBaru('1'),
      okupasi: [
        {
          occupationId: '1', occupationName: 'A',
          objek: [
            obj({ section: '1', tsi: '1000', currency: 'IDR', coverages: [cov({ coverage: 'C1', calculateMethodFacIn: '2', tsi: '1000', deductibles: [ded('IDR')] })] }),
            obj({ section: '2', tsi: '2000', currency: 'USD', coverages: [] }),
          ],
        },
      ],
    },
    {
      ...lokasiAnekaBaru('2'),
      okupasi: [{ occupationId: '2', occupationName: 'B', objek: [obj({ section: '1', tsi: '3000', coverages: [cov({ coverage: 'C1' }), cov({ coverage: 'C9' })] })] }],
      objekGolf: [obj({ section: '1', tsi: '9', coverages: [cov({ coverage: 'C1' })] })],
    },
  ]
}

describe('coverageAneka', () => {
  it('coverage baru: nilai awal sel Pega (Day 365, Loading 0, Scale 100, % Indemnity 100); metode kosong → 1', () => {
    const c = coverageAnekaBaru({ currency: 'USD', tsi: '5' }, '')
    expect([c.calculateMethodFacIn, c.day, c.loading, c.lostLimit, c.indemnityPercentage, c.currency, c.tsi]).toEqual(['1', '365', '0', '100', '100', 'USD', '5'])
  })

  it('Add: metode dari coverage pertama objek pertama okupasi pertama lokasi pertama (indeks keras L294-295); TSI/Currency objek', () => {
    const l = pohon()
    expect(metodeAwal(l)).toBe('2')
    const h = tambahCoverage(l, { l: 0, k: 0, j: 1 })
    const baru = h[0]!.okupasi[0]!.objek[1]!.coverages![0]!
    expect([baru.calculateMethodFacIn, baru.currency, baru.tsi]).toEqual(['2', 'USD', '2000'])
  })

  it('setTSIAneka_preDT: TSI tiap objek didorong ke semua coverage objek itu', () => {
    const l = pohon()
    l[1]!.okupasi[0]!.objek[0]!.coverages![1]!.tsi = '7'
    const h = dorongTsiObjek(l)
    expect(h[1]!.okupasi[0]!.objek[0]!.coverages!.map((c) => c.tsi)).toEqual(['3000', '3000'])
    expect(h[1]!.objekGolf[0]!.coverages![0]!.tsi).toBe('9')
  })

  it('SetTSI_DT / SetCurrency_DT: TSI dan Currency coverage yang diubah ikut ke objek', () => {
    const p = { l: 0, k: 0, j: 0 }
    const l = pohon()
    const c = l[0]!.okupasi[0]!.objek[0]!.coverages![0]!
    const h = ubahCoverageDi(l, p, 0, { ...c, tsi: '1500', currency: 'USD' })
    expect([h[0]!.okupasi[0]!.objek[0]!.tsi, h[0]!.okupasi[0]!.objek[0]!.currency]).toEqual(['1500', 'USD'])
    const tetap = ubahCoverageDi(l, p, 0, { ...c, rate: '1' })
    expect(tetap[0]!.okupasi[0]!.objek[0]!.tsi).toBe('1000')
  })

  it('Copy Coverage From This Section: objek okupasi ber-Section sama (termasuk sumber), TSI sasaran; Golf tidak ikut', () => {
    const h = salinCoverageSeksi(pohon(), { l: 0, k: 0, j: 0 })
    expect(h.sasaran).toEqual([{ l: 0, k: 0, j: 0 }, { l: 1, k: 0, j: 0 }])
    const sasaran = h.lokasi[1]!.okupasi[0]!.objek[0]!.coverages!
    expect(sasaran.map((c) => [c.coverage, c.tsi])).toEqual([['C1', '3000']])
    expect(sasaran[0]!.deductibles).toEqual([ded('IDR')])
    expect(h.lokasi[0]!.okupasi[0]!.objek[1]!.coverages).toEqual([]) // Section 2
    expect(h.lokasi[1]!.objekGolf[0]!.coverages![0]!.tsi).toBe('100') // Golf tidak disentuh
  })

  it('Copy Deductible: coverage berkode sama di semua lokasi okupasi diganti deductible sumber; kode lain tetap', () => {
    const h = salinDeductibleAneka(pohon(), { l: 0, k: 0, j: 0 })
    expect(h.jumlah).toBe(2) // sumber sendiri + C1 lokasi 2
    const l2 = h.lokasi[1]!.okupasi[0]!.objek[0]!.coverages!
    expect(l2[0]!.deductibles).toEqual([ded('IDR')])
    expect(l2[1]!.deductibles).toEqual([])
    expect(h.lokasi[1]!.objekGolf[0]!.coverages![0]!.deductibles).toEqual([])
    expect(salinDeductibleAneka(pohon(), { l: 0, k: 0, j: 1 }).jumlah).toBe(0)
  })

  it('addDeductibleAneka_ACT: Currency deductible TERAKHIR setiap coverage = Currency coverage itu', () => {
    const l = pohon()
    l[1]!.okupasi[0]!.objek[0]!.coverages![1] = cov({ coverage: 'C9', currency: 'SGD', deductibles: [ded('IDR'), ded('IDR')] })
    const h = setelCurrencyDeductibleTerakhir(l)
    expect(h[1]!.okupasi[0]!.objek[0]!.coverages![1]!.deductibles.map((d) => d.currency)).toEqual(['IDR', 'SGD'])
  })
})
