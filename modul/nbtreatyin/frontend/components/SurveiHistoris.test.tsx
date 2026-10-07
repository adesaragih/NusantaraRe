// Popup Historical Survey Report DIRENDER (react-dom/server). Harapan dari XML:
//   admin  `Section/HistoricalSurveyReportDtl`: Insured Name hanya-baca, tombol Add / Delete / Submit,
//          kolom Date of Survey, Surveyed by (Ceding Company), Loss Prevention, Remarks
//   atasan `Section/HistoricalSurveyReportDtlUW`: hanya-baca, tanpa Add / Delete / Submit, tanpa Loss Prevention

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Halaman } from '../api'
import SurveiHistoris from './SurveiHistoris'

const h: Halaman = {
  nilai: { 'PolicyTreatyIn.InsuredName': 'UJI-TERTANGGUNG' },
  daftar: { 'PolicyTreatyIn.QuotationData.SurveyReportList': [{ DateofSurvey: '2026-09-01', SurveyedBy: 'UJI-SURVEYOR', LossPrevention: '12.5' }] },
}
const render = (admin: boolean, sunting: boolean) =>
  renderToStaticMarkup(<SurveiHistoris halaman={h} admin={admin} sunting={sunting} onSetel={() => undefined} onTutup={() => undefined} />)
const tombol = (html: string) => [...html.matchAll(/<button[^>]*>([^<]*)<\/button>/g)].map((m) => m[1])

describe('popup Historical Survey Report', () => {
  it('admin: Insured Name, tombol Add / Delete / Submit, empat kolom', () => {
    const html = render(true, true)
    expect(html).toContain('Historical Survey Report')
    expect(html).toContain('UJI-TERTANGGUNG')
    expect(tombol(html)).toEqual(expect.arrayContaining(['Add', 'Delete', 'Submit']))
    for (const k of ['Date of Survey', 'Surveyed by (Ceding Company)', 'Loss Prevention', 'Remarks']) expect(html).toContain(k)
    expect(html).toContain('value="UJI-SURVEYOR"')
  })

  it('atasan: hanya-baca, tanpa Add / Delete / Submit dan tanpa Loss Prevention', () => {
    const html = render(false, false)
    expect(tombol(html)).not.toEqual(expect.arrayContaining(['Add']))
    expect(tombol(html)).not.toContain('Delete')
    expect(tombol(html)).not.toContain('Submit')
    expect(html).not.toContain('Loss Prevention')
    expect(html).toContain('UJI-SURVEYOR')
    expect(html).not.toContain('<input')
  })
})

// Prompt values property `.Remarks` (screenshot work owner 06-10-2026): 1 = Satisfied, 0 = Unsatisfied.
describe('Remarks survei', () => {
  const dengan = (remarks: string): Halaman => ({
    nilai: {},
    daftar: { 'PolicyTreatyIn.QuotationData.SurveyReportList': [{ DateofSurvey: '2026-09-01', Remarks: remarks }] },
  })
  const tampil = (hh: Halaman, admin: boolean, sunting: boolean) =>
    renderToStaticMarkup(<SurveiHistoris halaman={hh} admin={admin} sunting={sunting} onSetel={() => undefined} onTutup={() => undefined} />)

  it('admin: dropdown Satisfied / Unsatisfied (nilai 1 / 0), nilai tersimpan terpilih', () => {
    const html = tampil(dengan('0'), true, true)
    expect(html).toContain('<option value="1">Satisfied</option>')
    expect(html).toMatch(/<option value="0" selected="">Unsatisfied<\/option>/)
  })

  it('atasan: teks prompt, bukan nilai standar', () => {
    const html = tampil(dengan('1'), false, false)
    expect(html).toContain('Satisfied')
    expect(html).not.toContain('<select')
  })
})
