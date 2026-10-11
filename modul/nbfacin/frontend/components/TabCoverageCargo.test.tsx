// Tab Coverage MARINE CARGO (tiket 74 P2) - label korpus, rantai hitung, Deskripsi deductible aditif, gerbang layar.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { FORM_COVERAGE_CARGO as F, GRID_CARGO_COVERAGE, GRID_COVERAGE_CARGO } from '../labels'
import { coverageCargoKosong, terapkanHitungCargo } from './cargo'
import GridDeductible, { deductibleBaru } from './GridDeductible'
import { FormCoverageCargo } from './TabCoverageCargo'

const SUMBER = readFileSync(join(__dirname, 'TabCoverageCargo.tsx'), 'utf8').replace(/\r\n/g, '\n')
const INWARD = readFileSync(join(__dirname, '..', 'pages', 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
const COV_ANEKA = readFileSync(join(__dirname, 'FormCoverageAneka.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('TabCoverageCargo - label korpus (pyLabelPreview / pyLabelFieldValue)', () => {
  it('grid cargo, grid coverage, grid total, popup Select Coverage', () => {
    expect(GRID_CARGO_COVERAGE).toEqual(['No.', 'Trading', 'Goods', 'Packing', 'Conveyance'])
    expect(GRID_COVERAGE_CARGO.kolom).toEqual(['Coverage Note', 'Coverage Initial'])
    expect(GRID_COVERAGE_CARGO.total).toEqual(['Currency', 'Grand Total Gross TSI', 'Grand Total Gross Premium', 'Grand Total Gross Rate'])
  })

  it('form: bukan pyLabelFor ("Keterangan Jaminan", "Premi", "Diskon (%)")', () => {
    expect([F.coverageNote, F.currency, F.rate, F.limitOfLiability, F.discountPercentage, F.tsi, F.premium, F.minPremium, F.discount]).toEqual([
      'Coverage Note', 'Currency', 'Rate (%)', 'Limit of Liability', 'Discount (%)', 'TSI', 'Premium', 'Min Premium', 'Discount',
    ])
  })
})

describe('TabCoverageCargo - rantai hitung', () => {
  it('pemicu = sel ber-aksi Pega: currency / rate / diskonPersen / tsi / premium', () => {
    for (const p of ["hitung(baru, 'currency')", "hitung(c, 'rate')", "hitung(c, 'diskonPersen')", "hitung(c, 'tsi')", "hitung(c, 'premium')"]) {
      expect(SUMBER).toContain(p)
    }
    // Limit of Liability dan Min Premium tanpa aksi on-change di Pega.
    expect(SUMBER).toMatch(/label=\{F\.limitOfLiability\} value=\{c\.limitOfLiability\} onChange=\{set\('limitOfLiability'\)\} \/>/)
    expect(SUMBER).toMatch(/label=\{F\.minPremium\} value=\{c\.minPremium\} onChange=\{set\('minPremium'\)\} \/>/)
  })

  it('hasil hitung mempertahankan Warranty Extension dan Deductible di layar', () => {
    const layar = { ...coverageCargoKosong(), warrantyExtension: [{ coverage: 'W1', coverageNote: 'N' }], deductibles: [deductibleBaru('USD')] }
    const hasil = { ...coverageCargoKosong(), tsi: '100', rate: '1', premium: '1', discount: '0', warrantyExtension: [], deductibles: [] }
    const r = terapkanHitungCargo(layar, hasil)
    expect(r.premium).toBe('1')
    expect(r.warrantyExtension).toEqual(layar.warrantyExtension)
    expect(r.deductibles).toEqual(layar.deductibles)
  })
})

describe('Deskripsi deductible - aditif, Marine Cargo saja', () => {
  it('tanpa prop deskripsi (FIRE / Aneka) tidak ada isian Deskripsi; FormCoverageAneka tidak mengirimnya', () => {
    expect(COV_ANEKA).not.toContain('deskripsi=')
    const fire = renderToStaticMarkup(<GridDeductible daftar={[deductibleBaru('IDR')]} currency="IDR" mataUang={['IDR']} ubah={() => {}} />)
    expect(fire).not.toContain(F.deskripsi)
  })

  it('form coverage cargo memasang Deskripsi; Select Coverage + Warranty Extension tampil', () => {
    expect(SUMBER).toContain('deskripsi={F.deskripsi}')
    const h = renderToStaticMarkup(<FormCoverageCargo c={coverageCargoKosong()} idx="0-0" ubah={() => {}} hitung={() => {}} mataUang={['IDR']} />)
    for (const t of [F.coverageNote, F.coverageInitial, F.currency, F.rate, F.tsi, F.premium, F.grupWarranty]) expect(h).toContain(t)
    // Work owner 11-10-2026: tombol Select Coverage diganti isian saran (wajib dari daftar); urutan = tangkapan layar Pega.
    expect(h).not.toContain(`>${F.selectCoverage}<`)
    expect(SUMBER).toMatch(/<IsianSaran\s+label=\{F\.coverageNote\}[\s\S]*?wajibDariDaftar[\s\S]*?onKetik=\{\(v\) => ubah\(\{ \.\.\.c, coverageNote: v, coverage: '' \}\)\}/)
    const urut = [F.coverageNote, F.coverageInitial, F.currency, F.rate, F.limitOfLiability, F.discountPercentage, F.tsi, F.premium, F.minPremium]
    const pos = urut.map((t) => h.indexOf(`>${t}<`))
    expect(pos.every((p) => p >= 0)).toBe(true)
    expect([...pos].sort((a, b) => a - b)).toEqual(pos)
  })
})

describe('TabCoverageCargo - gerbang layar', () => {
  it('tab Coverage Marine Cargo hanya cabang tambahan', () => {
    expect(INWARD).toContain("tab === 'Coverage' && kasusCargo ? (\n              <TabCoverageCargo caseId={kasus.caseId} />")
  })
})
