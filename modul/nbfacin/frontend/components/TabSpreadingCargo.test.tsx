// Tab Spreading + tab lain + Submit MARINE CARGO (tiket 74 P3-P5) - label korpus, fungsi murni, sambungan layar.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { CargoSpreading } from '../api'
import { DEDUCTION_CARGO, SPREADING_CARGO } from '../labels'
import { barisLengkapCargo, barisSahCargo, gantiBarisCargo, kunciCoverageCargo } from './spreadingCargo'

const baca = (...p: string[]) => readFileSync(join(__dirname, ...p), 'utf8').replace(/\r\n/g, '\n')
const INWARD = baca('..', 'pages', 'InwardFacultative.tsx')
const SPREAD = baca('TabSpreadingCargo.tsx')

const baris = (treatyType: string, sharePercentage: string) => ({
  treatyType, treatyName: '', sharePercentage, tsiGrossSpreaded: '', tsiSpreaded: '', claimEstimation: '', premiumSpreaded: '',
}) as unknown as CargoSpreading['coverages'][number]['spreading'][number]

const cargo = (b: ReturnType<typeof baris>[]): CargoSpreading[] => [
  {
    no: '1', tradingNote: '', goodNote: '', packingNote: '', conveyanceNote: '',
    coverages: [{ coverage: 'C1', currency: 'USD', coverageNote: '', rate: '', premium: '', discount: '', tsi: '', tsiLiability: '', limitOfLiability: '', tsiNusantaraRe: '', premiNusantaraRe: '', spreading: b }],
  },
]

describe('Spreading cargo - label korpus', () => {
  it('grid cargo / coverage = pyValue CoverageSpreadingList / InputCoverageSpreadingMarineCargo (tanpa kolom Old, K73-39)', () => {
    expect(SPREADING_CARGO.kolomCargo).toEqual(['No.', 'Trading', 'Goods', 'Packing', 'Conveyance'])
    expect(SPREADING_CARGO.kolomCoverage).toEqual(['Coverage', 'Standard Rate (%)', 'Gross Premium', 'Premium Nusantara Re'])
    expect(DEDUCTION_CARGO.kolomCoverage).toEqual(['Coverage', 'TSI', 'Gross Premium'])
  })

  it('detail coverage = SpreadingItem FIRE (TSI 100% / TSI Liability / TSI NR / Gross Premium / Premium NR)', () => {
    for (const k of ['tsi', 'tsiLiability', 'tsiNusantaraRe', 'premium', 'premiNusantaraRe']) expect(SPREAD).toContain(`S.detailCoverage.${k}`)
  })
})

describe('spreadingCargo.ts', () => {
  it('ganti baris satu coverage tanpa mengubah masukan; kunci per posisi', () => {
    const awal = cargo([])
    const baru = gantiBarisCargo(awal, { c: 0, n: 0 }, [baris('10001', '100')])
    expect(awal[0]!.coverages[0]!.spreading).toEqual([])
    expect(baru[0]!.coverages[0]!.spreading).toHaveLength(1)
    expect(kunciCoverageCargo({ c: 1, n: 2 })).toBe('cc1-2')
  })

  it('baris sah: Type Treaty + % Share wajib, Σ % Share ≤ 100 per coverage', () => {
    expect(barisSahCargo(cargo([baris('10001', '60'), baris('10007', '40')]))).toBe(true)
    expect(barisSahCargo(cargo([baris('10001', '60'), baris('10007', '41')]))).toBe(false)
    expect(barisLengkapCargo(cargo([baris('', '10')]))).toBe(false)
  })
})

describe('sambungan layar Marine Cargo (P3-P5)', () => {
  it('Spreading / Cedant / Deduction / Payment dan Submit untuk kasusCargo; FIRE / Aneka tidak berubah', () => {
    expect(INWARD).toContain("tab === 'Spreading' && kasusCargo ? (\n              <TabSpreadingCargo caseId={kasus.caseId} tahapUw={dariServer?.ubahBarisSpreading === true} />")
    expect(INWARD).toContain('<TabCedant caseId={kasus.caseId} cargo />')
    expect(INWARD).toContain('<TabDeduction caseId={kasus.caseId} cargo />')
    expect(INWARD).toContain("tab === 'Payment' && kasusCargo ? (\n              <TabPayment caseId={kasus.caseId} />")
    expect(INWARD).toContain('const bolehSubmit = kasusFire || kasusAneka || kasusCargo')
  })

  it('tanpa spreading otomatis dan tanpa Copy Spreading per lokasi (Marine tidak punya cabangnya)', () => {
    expect(SPREAD).not.toContain('CopySpreadingLokasi')
    expect(SPREAD).not.toContain('salinLokasi')
  })
})
