// Tombol Survey Report (`Section/DetailPolicyTreatyIn` sel 22, `DetailDeptHeadTreatyIn_UW`): pyVisible
// `pyWorkPage.Quotation.ProportionalType != 'NonProportional'`; pyDisabledWhen
// `.QuotationData.IsSurveyReport=='No' || .QuotationData.IsSurveyReport==''`. Popup admin
// `Section/HistoricalSurveyReportDtl` mewajibkan `.DateofSurvey` (`InputHistoricalSurveyReportDtl`).

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import { aktifTombolSurvei, barisTanpaTanggal, DAFTAR_SURVEI, tampilTombolSurvei } from './survei'

const hal = (nilai: Record<string, string>): Halaman => ({ nilai, daftar: {} })
const SR = 'PolicyTreatyIn.QuotationData.IsSurveyReport'

describe('tombol Survey Report', () => {
  it('tampil kecuali kontrak NonProportional', () => {
    expect(tampilTombolSurvei(hal({}))).toBe(true)
    expect(tampilTombolSurvei(hal({ 'Quotation.ProportionalType': 'Proportional' }))).toBe(true)
    expect(tampilTombolSurvei(hal({ 'Quotation.ProportionalType': 'NonProportional' }))).toBe(false)
  })

  it('aktif hanya bila Survey Report bukan No / kosong', () => {
    expect(aktifTombolSurvei(hal({ [SR]: 'Yes' }))).toBe(true)
    expect(aktifTombolSurvei(hal({ [SR]: 'No' }))).toBe(false)
    expect(aktifTombolSurvei(hal({ [SR]: '' }))).toBe(false)
    expect(aktifTombolSurvei(hal({}))).toBe(false)
  })

  it('daftar grid = PolicyTreatyIn.QuotationData.SurveyReportList', () => {
    expect(DAFTAR_SURVEI).toBe('PolicyTreatyIn.QuotationData.SurveyReportList')
  })
})

describe('Date of Survey wajib', () => {
  it('nomor baris (1..n) yang belum berisi Date of Survey', () => {
    expect(barisTanpaTanggal([{ DateofSurvey: '2026-09-01' }, { SurveyedBy: 'UJI' }, { DateofSurvey: ' ' }])).toEqual([2, 3])
    expect(barisTanpaTanggal([])).toEqual([])
  })
})
