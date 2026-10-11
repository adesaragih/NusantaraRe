// Paritas tab Spreading Aneka (tiket 73 P3) - label rantai Pega + fungsi pohon. Data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { BarisSpreading, LokasiSpreadingAneka, ObjekSpreadingAneka } from '../api'
import { SPREADING_ANEKA as A } from '../labels'
import { barisLengkapAneka, barisSahAneka, gantiBarisAneka, kunciCoverageSpreading, semuaCoverageSpreading } from './spreadingAneka'

const SUMBER = readFileSync(join(__dirname, 'TabSpreadingAneka.tsx'), 'utf8').replace(/\r\n/g, '\n')
const HALAMAN = readFileSync(join(__dirname, '..', 'pages', 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')

const baris = (treatyType: string, sharePercentage: string): BarisSpreading => ({
  treatyType, treatyName: '', sharePercentage, tsiGrossSpreaded: '', tsiSpreaded: '', claimEstimation: '', premiumSpreaded: '',
})
const objek = (spreading: BarisSpreading[][]): ObjekSpreadingAneka => ({
  objectName: 'UJI', currency: 'IDR', tsiOld: '', tsi: '1',
  coverages: spreading.map((s) => ({
    oldId: '', coverageNote: 'C', loading: '', rate: '', tsi: '1', premium: '1', discount: '', limitOfLiability: '', tsiNusantaraRe: '',
    premiNusantaraRe: '', spreading: s,
  })),
})
const pohon = (): LokasiSpreadingAneka[] => [
  { objectNo: '1', location: 'L1', okupasi: [{ occupationName: 'A', objek: [objek([[baris('1', '60')], []])] }], objekGolf: [], total: [] },
  { objectNo: '2', location: 'L2', okupasi: [], objekGolf: [objek([[baris('1', '100')]])], total: [] },
]

describe('TabSpreadingAneka', () => {
  it('label rantai = label sel Pega (CoverageSpreadingList / SpreadingObjectAneka / SpreadingCoverageAneka / InputCoverageSpreadingAneka / SpreadingItemAneka)', () => {
    expect(A.kolomLokasi).toEqual(['No', 'Location'])
    expect(A.kolomOkupasi).toEqual(['Occupation Name'])
    // TSI Old hanya EDM (work owner 09-10-2026; menyimpang dari XML yang ALWAYS).
    expect(A.kolomObjek).toEqual(['Object Name', 'Currency', 'TSI'])
    expect(A.kolomCoverage).toEqual(['Coverage', 'Gross Premium'])
    expect(A.kolomObjekGolf).toEqual(['Object Name'])
    expect(A.kolomCoverageGolf).toEqual(['Coverage Note', '% Loading', '% Rate', 'TSI', 'Premi'])
    expect(Object.values(A.detailCoverage)).toEqual(['Gross TSI', 'TSI Nusantara Re', 'Gross Premium', 'Premium Nusantara Re'])
  })

  it('ganti baris satu coverage (okupasi dan Golf) tanpa mengubah masukan', () => {
    const l = pohon()
    const h = gantiBarisAneka(l, { l: 0, k: 0, j: 0, n: 1 }, [baris('2', '40')])
    expect(h[0]!.okupasi[0]!.objek[0]!.coverages[1]!.spreading).toEqual([baris('2', '40')])
    expect(l[0]!.okupasi[0]!.objek[0]!.coverages[1]!.spreading).toEqual([])
    const g = gantiBarisAneka(l, { l: 1, k: null, j: 0, n: 0 }, [])
    expect(g[1]!.objekGolf[0]!.coverages[0]!.spreading).toEqual([])
    expect(semuaCoverageSpreading(l)).toHaveLength(3)
    expect(kunciCoverageSpreading({ l: 1, k: null, j: 0, n: 0 })).toBe('c1-g-0-0')
  })

  it('baris sah = lengkap (Type Treaty + % Share angka) dan Σ % Share per coverage ≤ 100 (cekSpreadingFactIn)', () => {
    expect(barisSahAneka(pohon())).toBe(true)
    const lebih = gantiBarisAneka(pohon(), { l: 0, k: 0, j: 0, n: 0 }, [baris('1', '60'), baris('2', '50')])
    expect(barisSahAneka(lebih)).toBe(false)
    const kosong = gantiBarisAneka(pohon(), { l: 0, k: 0, j: 0, n: 0 }, [baris('', '10')])
    expect(barisLengkapAneka(kosong)).toBe(false)
  })

  it('K73-19: total per lokasi hanya bila tampilTotalLokasi; Copy Spreading dan edit baris hanya tahap UW', () => {
    expect(SUMBER).toContain('{v.tampilTotalLokasi && (')
    expect(SUMBER).toContain('{tahapUw && (\n          <div className="nbf-spread__copy">')
    expect(SUMBER).toContain('onClick={() => hitung(() => salinSpreadingAneka(caseId, share.trim(), v.template))}')
  })

  it('halaman case: tab Spreading Aneka di belakang FIRE, tahap UW sama dengan FIRE', () => {
    expect(HALAMAN).toContain("tab === 'Spreading' && kasusAneka ? (")
    expect(HALAMAN).toContain('<TabSpreadingAneka caseId={kasus.caseId} tahapUw={dariServer?.ubahBarisSpreading === true} />')
    expect(HALAMAN).toMatch(/: tab === 'Spreading' && kasusFire \? \(\s*<TabSpreading /)
  })

  it('tidak ada komponen yang didefinisikan di dalam komponen', () => {
    const dalam = SUMBER.slice(SUMBER.indexOf('export default function TabSpreadingAneka'))
    expect(dalam).not.toMatch(/^[ \t]+function [A-Z]\w*\(/m)
  })
})
