// Tab Premium Deduction / Brokerage Fee kasus FIRE (tiket 50) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

import { DEDUCTION as D, TEKS_DEDUCTION as T } from '../labels'
import { deduction2Lebih, kepalaPajak, lebihBesar } from './TabDeduction'

const NB = 'D:\\migrasi\\RNM\\NB FacIn\\'
const SUMBER = readFileSync(join(__dirname, 'TabDeduction.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe.skipIf(!existsSync(NB + 'Section\\InputInwardFacultativeDtl.xml'))('label = korpus', () => {
  it('isian, centang, kepala grid (InputInwardFacultativeDtl tab Premium Deduction / Brokerage Fee)', () => {
    const xml = readFileSync(NB + 'Section\\InputInwardFacultativeDtl.xml', 'utf-8')
    expect(xml).toContain('<pyTitle>Premium Deduction / Brokerage Fee</pyTitle>')
    for (const t of [D.pctPremiumDeduction, D.pctBrokerageFee]) expect(xml).toContain(`<pyLabelFieldValue>${t}</pyLabelFieldValue>`)
    expect(xml).toContain(`<pyCheckboxCaption>${D.excludePpn}</pyCheckboxCaption>`)
    for (const t of [D.kolomMataUang, D.totalPremiumDeduction, D.totalPremiumDeduction2, D.totalBrokerageFee]) {
      expect(xml).toContain(`<pyValue>${t}</pyValue>`)
    }
  })
  it('daftar komisi (CoverageCommisionList → … → CommisionItem)', () => {
    const baca = (f: string) => readFileSync(NB + 'Section\\' + f, 'utf-8')
    for (const t of D.kolomLokasi) expect(baca('CoverageCommisionList.xml')).toContain(`<pyValue>${t}</pyValue>`)
    for (const t of D.kolomItem) expect(baca('PropertyItemListCoverageCommision.xml')).toContain(`<pyValue>${t}</pyValue>`)
    for (const t of D.kolomCoverage) expect(baca('InputCoverageCommisionFire.xml')).toContain(`<pyValue>${t}</pyValue>`)
    for (const t of Object.values(D.detailCoverage)) expect(baca('CommisionItem.xml')).toContain(`<pyLabelFieldValue>${t}</pyLabelFieldValue>`)
    for (const t of D.kolomSpreading) expect(baca('InputDtlSpreadingCommisionCoverage_FacIn.xml')).toContain(`<pyValue>${t}</pyValue>`)
  })
  it('pesan Deduction 2 dan kepala pajak (ProtekDeduction2_Act, SetTAXNote)', () => {
    expect(readFileSync(NB + 'Activity\\ProtekDeduction2_Act.xml', 'utf-8')).toContain(`<PropertiesValue>"${T.deduction2Lebih}"</PropertiesValue>`)
    const tax = readFileSync(NB + 'Activity\\SetTAXNote.xml', 'utf-8')
    expect(tax).toContain('<PropertiesValue>"PPH " +  Tax.pxResults(1).PPH + " %"</PropertiesValue>')
    expect(tax).toContain('<PropertiesValue>"PPN " +  Tax.pxResults(1).PPN + " %"</PropertiesValue>')
  })
})

describe('TabDeduction - aturan', () => {
  it('kepala pajak "PPH 2 %" / "PPN 2.2 %" (gambar layar Pega); tarif kosong -> label saja', () => {
    expect(kepalaPajak(D.pph, '2')).toBe('PPH 2 %')
    expect(kepalaPajak(D.ppn, '2.20000000')).toBe('PPN 2.2 %')
    expect(kepalaPajak(D.pph, '')).toBe('PPH')
  })

  it('lebihBesar eksak', () => {
    expect([lebihBesar('10', '9.9999'), lebihBesar('10', '10.0000'), lebihBesar('9', '10'), lebihBesar('', '1')]).toEqual([true, false, false, false])
  })

  it('Deduction 2 > Premium ditolak (ProtekDeduction2_Act)', () => {
    const b = (deduction2: string) => ({ currency: 'IDR', premium: '469550', commision: '', deduction2, brokerageFee: '', pph: '', ppn: '' })
    expect(deduction2Lebih(b('469550.0001'))).toBe(true)
    expect(deduction2Lebih(b('469550'))).toBe(false)
    expect(deduction2Lebih(b(''))).toBe(false)
  })

  it('hitungan di backend: persen / sel grid -> hitung; Save -> simpan; tanpa float', () => {
    expect(SUMBER).toContain("void hitung('persen', 0, v)")
    expect(SUMBER).toContain("void hitung('commision', n, v)")
    expect(SUMBER).toContain("void hitung('brokerage', n, v)")
    expect(SUMBER).toContain('if (!deduction2Lebih(b)) void hitung(\'deduction2\', n, v)')
    expect(SUMBER).toContain('setV(await simpanDeduction(caseId, v))')
    expect(SUMBER).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|parseInt/)
  })
})
