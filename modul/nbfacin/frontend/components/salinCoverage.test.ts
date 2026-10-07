// Tombol salin coverage / deductible (tiket 53) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { describe, expect, it } from 'vitest'

import type { CoverageObjek, ObjekFire } from '../api'
import { SALIN_COVERAGE as S } from '../labels'
import { coverageBaru } from './FormCoverage'
import { JENIS_ITEM_BI, salinCoverageDari, salinCoverageLokasi, salinCoverageSemua, salinDeductible } from './salinCoverage'

const NB = 'D:\\migrasi\\RNM\\NB FacIn\\'

describe.skipIf(!existsSync(NB + 'Section\\InputCoverageFire.xml'))('label = korpus', () => {
  it('tombol InputCoverageFire dan popup CopyCoverageFrom', () => {
    const grid = readFileSync(NB + 'Section\\InputCoverageFire.xml', 'utf-8')
    for (const t of [S.coverage, S.semua, S.dari, S.deductibleSemua, S.deductibleBI, S.keItemLain]) expect(grid).toContain(`<pyLabel>${t}</pyLabel>`)
    const pop = readFileSync(NB + 'Section\\CopyCoverageFrom.xml', 'utf-8')
    for (const t of [S.lokasi, S.item]) expect(pop).toContain(`<pyLabelFieldValue>${t}</pyLabelFieldValue>`)
    expect(pop).toContain(`<pyLabel>${S.keItemLain}</pyLabel>`)
  })
  it('BI = 10000992 (CopyDeductibleToAllLocation_Act)', () => {
    const akt = readFileSync(NB + 'Activity\\CopyDeductibleToAllLocation_Act.xml', 'utf-8')
    expect(akt).toContain(`<pyStepsPreCondParamsWhen>.ItemTypeID=="${JENIS_ITEM_BI}"</pyStepsPreCondParamsWhen>`)
    expect(akt).toContain('<pyStepsPreCondParamsWhen>Param.idcoverage==.Coverage</pyStepsPreCondParamsWhen>')
  })
})

const cov = (coverage: string, premium = '', deductibles: CoverageObjek['deductibles'] = []): CoverageObjek => ({
  ...coverageBaru(),
  coverage,
  oldId: coverage,
  premium,
  deductibles,
})
const it1 = (itemTypeId: string, tsi: string, coverages: CoverageObjek[], totalNetRate = '') =>
  ({ itemTypeId, itemType: itemTypeId, tsi, coverages, totalNetRate }) as unknown as ObjekFire['items'][number]
const lok = (items: ObjekFire['items']) => ({ items, riskLocation: 'UJI', objectNo: '1' }) as unknown as ObjekFire
const ded = (amount: string) => ({ amount }) as unknown as NonNullable<CoverageObjek['deductibles']>[number]

describe('salinCoverage - aturan', () => {
  const data = (): ObjekFire[] => [
    lok([it1('A', '100', [cov('C1', '5'), cov('C2', '6')], '1.5'), it1('B', '200', [cov('X', '9')])]),
    lok([it1('A', '300', [cov('Y', '1')])]),
  ]

  it('Copy Coverage: item 1 -> semua item lokasi itu; diganti; TSI tujuan; TotalNetRate disalin', () => {
    const h = salinCoverageLokasi(data(), 0)
    expect(h.tujuan).toEqual([{ o: 0, i: 0 }, { o: 0, i: 1 }])
    const b = h.objek[0]!.items[1]!
    expect(b.coverages!.map((c) => c.coverage)).toEqual(['C1', 'C2'])
    expect(b.coverages!.every((c) => c.tsi === '200')).toBe(true)
    expect(b.totalNetRate).toBe('1.5')
    expect(h.objek[1]!.items[0]!.coverages!.map((c) => c.coverage)).toEqual(['Y'])
  })

  it('Copy Coverage To All: lokasi 1 / item 1 -> semua item semua lokasi', () => {
    const h = salinCoverageSemua(data())
    expect(h.tujuan).toHaveLength(3)
    expect(h.objek[1]!.items[0]!.coverages!.map((c) => [c.coverage, c.tsi])).toEqual([
      ['C1', '300'],
      ['C2', '300'],
    ])
  })

  it('Copy Coverage From: item terpilih -> item yang dibuka saja; TotalNetRate tidak disalin', () => {
    const h = salinCoverageDari(data(), { o: 1, i: 0 }, { o: 0, i: 1 })
    expect(h.tujuan).toEqual([{ o: 0, i: 1 }])
    expect(h.objek[0]!.items[1]!.coverages!.map((c) => [c.coverage, c.tsi])).toEqual([['Y', '200']])
    expect(h.objek[0]!.items[1]!.totalNetRate).toBe('')
    expect(h.objek[0]!.items[0]!.coverages!.map((c) => c.coverage)).toEqual(['C1', 'C2'])
  })

  it('salinan berdiri sendiri (mengubah tujuan tidak mengubah sumber)', () => {
    const awal = data()
    const h = salinCoverageLokasi(awal, 0)
    h.objek[0]!.items[1]!.coverages![0]!.rate = '9'
    expect(h.objek[0]!.items[0]!.coverages![0]!.rate).toBe('')
  })

  it('Copy Deductible To All: per kode coverage; tanpa padanan -> kosong; item BI dilewati', () => {
    const d: ObjekFire[] = [
      lok([
        it1('A', '1', [cov('C1', '', [ded('10')]), cov('C2', '', [ded('20')])]),
        it1('B', '1', [cov('C2', '', [ded('99')]), cov('Z', '', [ded('5')])]),
        it1(JENIS_ITEM_BI, '1', [cov('C1', '', [ded('7')])]),
      ]),
    ]
    const h = salinDeductible(d, { o: 0, i: 0 }, 'ALL')
    expect(h[0]!.items[1]!.coverages!.map((c) => (c.deductibles ?? []).map((x) => x.amount))).toEqual([['20'], []])
    expect(h[0]!.items[2]!.coverages![0]!.deductibles!.map((x) => x.amount)).toEqual(['7'])
  })

  it('Copy Deductible For BI: hanya item BI', () => {
    const d: ObjekFire[] = [
      lok([it1(JENIS_ITEM_BI, '1', [cov('C1', '', [ded('10')])]), it1('B', '1', [cov('C1', '', [ded('99')])])]),
      lok([it1(JENIS_ITEM_BI, '1', [cov('C1', '', [])])]),
    ]
    const h = salinDeductible(d, { o: 0, i: 0 }, 'BI')
    expect(h[1]!.items[0]!.coverages![0]!.deductibles!.map((x) => x.amount)).toEqual(['10'])
    expect(h[0]!.items[1]!.coverages![0]!.deductibles!.map((x) => x.amount)).toEqual(['99'])
  })
})
