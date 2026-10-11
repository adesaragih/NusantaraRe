// Tampilan polis NB keluarga Aneka (work owner 11-10-2026 "tampilan view polis yang belum tergenerate tolong diperbaiki.
// cek tampilan nya di nb facin!"): lini / varian tampilan = When nbfacin, total = rumus nbfacin.

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import {
  bisnisNB,
  hitungSpreadingAnekaNB,
  kelompokMedanAneka,
  kolomObjekAnekaNB,
  liniNB,
  tampilanAnekaNB,
  tampilOkupasiNB,
  totalCoverageAnekaNB,
} from './polisNBAneka'

const Q = 'OfferFacIn.QuotationData.'
const halaman = (
  type: string,
  code = '',
  name = '',
  daftar: Halaman['daftar'] = {},
  nilai: Record<string, string> = {},
) =>
  ({
    nilai: { [Q + 'BusinessType']: type, [Q + 'BusinessCode']: code, [Q + 'BusinessName']: name, ...nilai },
    daftar,
  }) as Halaman

describe('lini dan varian tampilan Aneka (lini.go nbfacin)', () => {
  it('gerbang FIRE > Aneka > Marine Cargo', () => {
    expect(liniNB(halaman('Liability'), false)).toMatchObject({ aneka: true, cargo: false })
    expect(liniNB(halaman('Liability'), true)).toMatchObject({ aneka: false, cargo: false })
    expect(liniNB(halaman('MarineCargo'), false)).toMatchObject({ aneka: false, cargo: true })
    expect(liniNB(halaman('GolfInsurance'), false)).toMatchObject({ aneka: true, anekaWhen: false })
    expect(liniNB(halaman('CustomBond'), false)).toMatchObject({ aneka: true, anekaWhen: true })
    // IsElectronicEquipment nbfacin berspasi (" ElectronicEquipment ") - bukan Aneka (NB pun "belum tersedia")
    expect(liniNB(halaman('ElectronicEquipment'), false)).toMatchObject({ aneka: false, cargo: false })
    expect(liniNB(halaman('MBUCar'), false)).toMatchObject({ aneka: false, cargo: false })
  })

  it('kontainer dan sakelar per BusinessType', () => {
    expect(tampilanAnekaNB(bisnisNB(halaman('Liability', '10047')))).toMatchObject({
      kontainer: ['C1'],
      topRisk: true,
      jumlahSatuan: true,
    })
    expect(tampilanAnekaNB(bisnisNB(halaman('HE')))).toMatchObject({
      kontainer: ['C5', 'C6'],
      he: true,
      topRisk: true,
      currencyCoverage: true,
    })
    expect(tampilanAnekaNB(bisnisNB(halaman('MarineHull'))).kontainer).toEqual(['C2'])
    expect(tampilanAnekaNB(bisnisNB(halaman('MBD')))).toMatchObject({
      kontainer: ['C4', 'C5'],
      kolomTahunJumlah: true,
      engineering: true,
      mbd: true,
    })
    expect(tampilanAnekaNB(bisnisNB(halaman('Boiler'))).kontainer).toEqual(['C1', 'C4', 'C5'])
    expect(tampilanAnekaNB(bisnisNB(halaman('Bonding')))).toMatchObject({ kontainer: ['C1'], bonding: true })
    expect(tampilanAnekaNB(bisnisNB(halaman('Aneka', '10167')))).toMatchObject({
      kontainer: ['C8'],
      tanamanHasil: true,
    })
    expect(tampilanAnekaNB(bisnisNB(halaman('Aneka', '10235'))).briguna).toBe(true)
    expect(tampilanAnekaNB(bisnisNB(halaman('MBUCar'))).kontainer).toEqual([])
  })

  it('grid okupasi: hanya Golf murni yang tanpa okupasi', () => {
    const golf = tampilanAnekaNB(bisnisNB(halaman('GolfInsurance')))
    expect(tampilOkupasiNB(false, golf)).toBe(false)
    expect(tampilOkupasiNB(true, tampilanAnekaNB(bisnisNB(halaman('Liability'))))).toBe(true)
  })

  it('kolom grid objek K2 / K3 / K4', () => {
    expect(kolomObjekAnekaNB(tampilanAnekaNB(bisnisNB(halaman('HE'))))).toEqual([
      'objectName',
      'year',
      'chassisNumber',
      'engineNumber',
    ])
    expect(kolomObjekAnekaNB(tampilanAnekaNB(bisnisNB(halaman('MBD'))))).toEqual([
      'objectName',
      'year',
      'quantity',
      'unitQuantity',
    ])
    expect(kolomObjekAnekaNB(tampilanAnekaNB(bisnisNB(halaman('Aneka', '10167'))))).toEqual([
      'objectName',
      'currency',
      'tsi',
    ])
    expect(kolomObjekAnekaNB(tampilanAnekaNB(bisnisNB(halaman('Liability'))))).toEqual(['objectName'])
  })
})

describe('form objek Aneka (FormObjekAneka nbfacin)', () => {
  it('kartu Liability: Object (Description pertama), Sum Insured, Identity, Period', () => {
    const g = kelompokMedanAneka(tampilanAnekaNB(bisnisNB(halaman('Liability'))))
    expect(g.objek.map((m) => m.label)).toEqual(['Object Description', 'Section', 'Contract No'])
    expect(g.nilai.map((m) => m.label)).toEqual(['Currency', 'TSI'])
    expect(g.identitas.map((m) => m.label)).toEqual(['KTP No', 'NPWP No'])
    expect(g.periode.map((m) => m.label)).toEqual(['Start Date', 'End Date'])
    expect(g.kapal).toEqual([])
  })

  it('Marine Hull: Vessel Details memuat data kapal + Production Year / Place of Manufacture / Flag', () => {
    const g = kelompokMedanAneka(tampilanAnekaNB(bisnisNB(halaman('MarineHull'))))
    expect(g.objek.map((m) => m.label)).toEqual(['Vessel Name'])
    expect(g.kapal.map((m) => m.kunci)).toEqual([
      'Year',
      'BuiltIn',
      'Flag',
      'MarineHull.TypeOfVessel',
      'MarineHull.Classification',
      'MarineHull.Construction',
      'MarineHull.NRTRT',
      'MarineHull.GRTRT',
      'MarineHull.DWTRT',
      'MarineHull.LENGTHM',
      'MarineHull.BREADTHM',
      'MarineHull.DEPTHM',
    ])
  })

  it('Boiler (C1 + C4 + C5): medan kembar tampil sekali, kontainer pertama menang', () => {
    const g = kelompokMedanAneka(tampilanAnekaNB(bisnisNB(halaman('Boiler'))))
    const semua = Object.values(g).flat()
    expect(semua.filter((m) => m.kunci === 'ObjectName')).toHaveLength(1)
    expect(semua.find((m) => m.kunci === 'ObjectName')?.label).toBe('Object Description')
    expect(g.kendaraan.map((m) => m.label)).toEqual(['Chassis Number', 'Engine Number'])
  })
})

describe('total Coverage dan Spreading Aneka (totalaneka.go / spreadinganeka.go nbfacin)', () => {
  const L = 'OfferFacIn.LocationList'
  const ok = L + '(1).Property.RiskLocation.OccupationList'
  const an = ok + '(1).AnekaList'
  const daftar = (cov2Mu: string): Halaman['daftar'] => ({
    [L]: [{ 'Property.ObjectNo': '1' }],
    [ok]: [{ OccupationName: 'UJI OKUPASI' }],
    [an]: [
      { ObjectName: 'UJI A', 'Currency.Name': 'IDR', TSI: '1000' },
      { ObjectName: 'UJI B', 'Currency.Name': 'USD', TSI: '500' },
    ],
    [an + '(1).CoverageList']: [
      { 'Currency.Name': 'IDR', TSI: '1000', Premium: '10', LimitofLiability: '400' },
      { 'Currency.Name': 'IDR', TSI: '1000', Premium: '5' },
    ],
    [an + '(2).CoverageList']: [{ 'Currency.Name': cov2Mu, TSI: '500', Premium: '20' }],
    [an + '(1).CoverageList(1).SpreadingList']: [
      { TreatyType: '10003', TreatyName: 'UJI QS', SharePercentage: '60', TSISpreaded: '600', PremiumSpreaded: '6' },
    ],
    [an + '(1).CoverageList(2).SpreadingList']: [
      { TreatyType: '10003', TreatyName: 'UJI QS', SharePercentage: '60', TSISpreaded: '600', PremiumSpreaded: '3' },
    ],
  })

  it('Grand Total per objek (mata uang objek) untuk lini umum; per coverage untuk 10 When khusus', () => {
    expect(totalCoverageAnekaNB(halaman('Liability', '', '', daftar('IDR')), 1)).toEqual([
      { currency: 'IDR', tsi: '1000', premium: '15', rate: '1.5' },
      { currency: 'USD', tsi: '500', premium: '20', rate: '4' },
    ])
    // HE (khusus): TSI dan premi setiap coverage, mata uang coverage
    expect(totalCoverageAnekaNB(halaman('HE', '', '', daftar('IDR')), 1)).toEqual([
      { currency: 'IDR', tsi: '2500', premium: '35', rate: '1.4' },
    ])
  })

  it('spreading: TSI hanya coverage pertama, premi semua coverage, LoL x PS x share', () => {
    const h = halaman('Liability', '', '', daftar('IDR'), { 'OfferFacIn.PercentShare': '50' })
    const s = hitungSpreadingAnekaNB(h)
    expect(s.lokasi[0]).toHaveLength(1)
    expect(s.lokasi[0]![0]).toMatchObject({
      currency: 'IDR',
      treatyName: 'UJI QS',
      sharePercentage: '60',
      tsiSpreaded: '600',
      premiumSpreaded: '9',
      claimEstimation: '120', // 400 x 50% x 60%
    })
    expect(s.mataUang).toEqual([{ currency: 'IDR', tsiTopRisk: '0', tsi: '600', lol: '120', premium: '9' }])
  })
})
