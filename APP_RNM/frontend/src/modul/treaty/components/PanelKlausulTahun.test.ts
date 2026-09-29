// Uji panel klausul tahun — tiket 08 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { alihJenis } from './PanelKlausulTahun'

const KODE = readFileSync(join(__dirname, 'PanelKlausulTahun.tsx'), 'utf8')

describe('panel klausul tahun', () => {
  it('Show membuka/menutup satu jenis tanpa menyentuh jenis lain (AC 29)', () => {
    expect(alihJenis([], '10009')).toEqual(['10009'])
    expect(alihJenis(['10009', '10013'], '10009')).toEqual(['10013'])
    expect(alihJenis(['10013'], '10001')).toEqual(['10013', '10001'])
  })
  it('dua grid: For Non XOL (isXol 0) dan For XOL (isXol 1)', () => {
    expect(KODE).toContain("ambilJenisKlausul('0')")
    expect(KODE).toContain("ambilJenisKlausul('1')")
    expect(KODE).toContain('KLAUSUL_TCO.gridNonXol')
    expect(KODE).toContain('KLAUSUL_TCO.gridXol')
  })
  it('kepala tahun hanya dibaca — tujuh medan', () => {
    expect((KODE.match(/readOnly/g) ?? []).length).toBe(7)
  })
  it('label bersilang OQ-TCO-05 dibawa apa adanya', () => {
    expect(KODE).toContain('label={KLAUSUL_TCO.headerUnderwritingYear} value={tahun.treatyYear}')
    expect(KODE).toContain('label={KLAUSUL_TCO.headerTransactionYear} value={tahun.underwritingYear}')
  })
})
