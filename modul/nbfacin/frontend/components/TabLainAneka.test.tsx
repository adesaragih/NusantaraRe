// Tab Inw Fac Cedant Panels / Premium Deduction / Payment untuk keluarga Aneka (tiket 73 P3). Data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { LokasiSpreadingAneka } from '../api'
import { DEDUCTION_ANEKA as DA, SPREADING as S } from '../labels'
import { DaftarKomisiAneka } from './TabDeduction'
import { RingkasanSpreading } from './TabSpreading'

const HALAMAN = readFileSync(join(__dirname, '..', 'pages', 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
const CEDANT = readFileSync(join(__dirname, 'TabCedant.tsx'), 'utf8').replace(/\r\n/g, '\n')
const DEDUCT = readFileSync(join(__dirname, 'TabDeduction.tsx'), 'utf8').replace(/\r\n/g, '\n')

const v = {
  ringkasanTreaty: [],
  ringkasanMataUang: [{ currency: 'IDR', tsiTopRisk: '9', tsi: '100', lol: '5', premium: '1' }],
}

describe('tab lain untuk Aneka (tiket 73 P3)', () => {
  it('ringkasan per mata uang Cedant !IsFire tanpa kolom Top Risk (L71388-72464); FIRE tetap 5 kolom', () => {
    const fire = renderToStaticMarkup(<RingkasanSpreading v={v} hanyaMataUang />)
    const aneka = renderToStaticMarkup(<RingkasanSpreading v={v} hanyaMataUang tanpaTopRisk />)
    expect(fire).toContain(S.kolomRingkasanMataUang[1])
    expect(aneka).not.toContain(S.kolomRingkasanMataUang[1])
    expect(aneka).toContain('Total TSI Spreaded')
    expect(aneka).toContain('Total Limit of Liability')
    expect((aneka.match(/<td/g) ?? []).length).toBe(4)
    expect((fire.match(/<td/g) ?? []).length).toBe(5)
  })

  it('Cedant / Deduction mengambil ringkasan dari spreading-aneka bila aneka; FIRE tetap ambilSpreading', () => {
    expect(CEDANT).toContain(';(aneka ? ambilSpreadingAneka : ambilSpreading)(caseId).then(')
    expect(DEDUCT).toContain(';(aneka ? ambilSpreadingAneka : ambilSpreading)(caseId).then(')
    expect(DEDUCT).toContain('aneka ? <DaftarKomisiAneka')
  })

  it('daftar komisi Aneka: label rantai CommisionObjectAneka / CommisionCoverageAneka (tertutup = grid lokasi)', () => {
    const lokasi: LokasiSpreadingAneka[] = [{ objectNo: '1', location: 'UJI LOKASI', okupasi: [], objekGolf: [], total: [] }]
    const h = renderToStaticMarkup(<DaftarKomisiAneka lokasi={lokasi} />)
    for (const k of DA.kolomLokasi) expect(h).toContain(`>${k}<`)
    expect(h).toContain('UJI LOKASI')
    expect(DA.kolomObjek).toEqual(['Object Name', 'Currency', 'TSI']) // TSI Old hanya EDM (work owner 09-10-2026)
    expect(DA.kolomCoverage).toEqual(['Coverage', 'Gross Premium'])
    expect(DA.detailCoverage.tsi).toBe('TSI Liability')
  })

  it('halaman case: Cedant / Premium Deduction / Payment dibuka untuk case Aneka (di belakang FIRE)', () => {
    expect(HALAMAN).toContain("tab === 'Inw Fac Cedant Panels' && kasusAneka ? (\n              <TabCedant caseId={kasus.caseId} aneka />")
    expect(HALAMAN).toContain("tab === 'Premium Deduction / Brokerage Fee' && kasusAneka ? (\n              <TabDeduction caseId={kasus.caseId} aneka />")
    expect(HALAMAN).toContain("tab === 'Payment' && kasusAneka ? (\n              <TabPayment caseId={kasus.caseId} />")
    expect(HALAMAN).toMatch(/: tab === 'Inw Fac Cedant Panels' && kasusFire \? \(\s*<TabCedant caseId=\{kasus\.caseId\} \/>/)
  })
})
