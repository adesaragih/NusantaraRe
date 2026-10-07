// Tab Payment kasus FIRE (tiket 51) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

import { PAYMENT as P, TEKS_PAYMENT as T } from '../labels'
import { installmentSah, pesanPersenCicilan } from './TabPayment'

const NB = 'D:\\migrasi\\RNM\\NB FacIn\\'
const EDM = 'D:\\migrasi\\RNM\\Endorsment Fac In\\'
const SUMBER = readFileSync(join(__dirname, 'TabPayment.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe.skipIf(!existsSync(NB + 'Section\\PaymentCurrencyList.xml'))('label = korpus', () => {
  it('Payment Information, Installment, grid mata uang (Dtl + PaymentCurrencyList)', () => {
    const dtl = readFileSync(NB + 'Section\\InputInwardFacultativeDtl.xml', 'utf-8')
    expect(dtl).toContain(`<pyValue>&lt;b&gt;${P.judul}&lt;/b&gt;</pyValue>`)
    const pcl = readFileSync(NB + 'Section\\PaymentCurrencyList.xml', 'utf-8')
    expect(pcl).toContain(`<pyLabelFieldValue>${P.installment}</pyLabelFieldValue>`)
    for (const k of P.kolomMataUang) expect(pcl).toContain(`<pyValue>${k}</pyValue>`)
  })
  it('grid cicilan (InputDtlPayment_FacIn)', () => {
    const xml = readFileSync(NB + 'Section\\InputDtlPayment_FacIn.xml', 'utf-8')
    expect(xml).toContain(`>${P.judulCicilan}<`)
    for (const k of [P.no, P.dueDate, P.persen, P.premium, P.deduction, P.deduction2, P.brokerage, P.discount, P.totalPayment]) {
      expect(xml).toContain(`<pyValue>${k}</pyValue>`)
    }
  })
  it.skipIf(!existsSync(EDM + 'Activity\\SetValidateInstallment_Act.xml'))('pesan Σ % (SetValidateInstallment_Act kelas Currency)', () => {
    const xml = readFileSync(EDM + 'Activity\\SetValidateInstallment_Act.xml', 'utf-8')
    expect(xml).toContain('<pxInsName>ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY!SETVALIDATEINSTALLMENT_ACT</pxInsName>')
    for (const t of [T.persenLebih, T.persenKurang]) expect(xml).toContain(`<PropertiesValue>"${t}"</PropertiesValue>`)
  })
})

describe('TabPayment - aturan', () => {
  const c = (installmentPercentage: string) => ({
    installmentNo: '1', dueDate: '', installmentPercentage, premium: '', riCommision: '', deduction2: '', brokerageFee: '', pph: '', ppn: '',
    discount: '', paymentTotal: '',
  })
  it('Σ % Installment: tepat 100 -> tanpa pesan; > 100 / < 100 -> pesan Pega', () => {
    expect(pesanPersenCicilan([c('33.33333333333333333333'), c('33.33333333333333333333'), c('33.33333333333333333334')])).toBe('')
    expect(pesanPersenCicilan([c('60'), c('40.0001')])).toBe(T.persenLebih)
    expect(pesanPersenCicilan([c('60'), c('39.9999')])).toBe(T.persenKurang)
    expect(pesanPersenCicilan([])).toBe('')
  })
  it('Installment bilangan bulat 1..99', () => {
    expect(['1', '12', '99'].map(installmentSah)).toEqual([true, true, true])
    expect(['', '0', '100', '1.5', 'a'].map(installmentSah)).toEqual([false, false, false, false, false])
  })
  it('hitungan di backend; tanpa float', () => {
    expect(SUMBER).toContain("void hitung('installment', '', 0, v)")
    expect(SUMBER).toContain("void hitung('persen', m.currency, n, v)")
    expect(SUMBER).toContain('setV(await simpanPayment(caseId, v))')
    expect(SUMBER).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|parseInt/)
  })
})
