// Uji panel klausul tahun — tiket 08 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { alihJenisTunggal } from './PanelKlausulTahun'

const KODE = readFileSync(join(__dirname, 'PanelKlausulTahun.tsx'), 'utf8')

describe('panel klausul tahun', () => {
  it('Show membuka SATU jenis; jenis lain menggantinya, jenis yang sama menutupnya (keputusan work owner 30-09-2026)', () => {
    expect(alihJenisTunggal(null, '10001')).toBe('10001')
    // Treaty Limit terbuka, lalu Cash Loss Limit ditekan: yang tampil HANYA Cash Loss Limit.
    expect(alihJenisTunggal('10001', '10004')).toBe('10004')
    expect(alihJenisTunggal('10004', '10004')).toBeNull()
  })
  it('jenis terbuka tampil di POPUP, bukan di bawah grid', () => {
    expect(KODE).toContain('<Modal')
    expect(KODE).toMatch(/<Modal[\s\S]*<PanelJenisKlausul[\s\S]*<\/Modal>/)
    expect(KODE).not.toContain('terbuka.map(')
  })
  it('satu grid "Treaty Desc" (isXol 0); grid For XOL dibuang', () => {
    expect(KODE).toContain("ambilJenisKlausul('0')")
    expect(KODE).not.toContain("ambilJenisKlausul('1')")
    expect(KODE).toContain('KLAUSUL_TCO.gridNonXol')
    expect(KODE).not.toContain('gridXol')
  })
  it('kepala tahun hanya dibaca — tujuh medan', () => {
    expect((KODE.match(/readOnly/g) ?? []).length).toBe(7)
  })
  it('label bersilang OQ-TCO-05 dibawa apa adanya', () => {
    expect(KODE).toContain('label={KLAUSUL_TCO.headerUnderwritingYear} value={tahun.treatyYear}')
    expect(KODE).toContain('label={KLAUSUL_TCO.headerTransactionYear} value={tahun.underwritingYear}')
  })
})
