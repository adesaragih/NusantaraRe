import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { PROPORSI_TCO, labelProporsi } from './proporsi'

// Keputusan work owner 30-09-2026 (OQ-TCO-04): Reinsurance Type tahun treaty
// = dua nilai, bukan master jenis reasuransi.

const baca = (relatif: string) => readFileSync(join(__dirname, relatif), 'utf8')

describe('Reinsurance Type tahun treaty', () => {
  it('tepat dua pilihan: nilai tersimpan dan labelnya', () => {
    expect(PROPORSI_TCO.map((p) => [p.value, p.label])).toEqual([
      ['Proportional', 'Proportional'],
      ['NonProportional', 'Non Proportional'],
    ])
  })
  it('label tampil; nilai lain (data lama) apa adanya', () => {
    expect(labelProporsi('NonProportional')).toBe('Non Proportional')
    expect(labelProporsi('Proportional')).toBe('Proportional')
    expect(labelProporsi('1000003')).toBe('1000003')
  })
  it('form tahun memakai daftar ini, bukan master jenis reasuransi; grid dan kepala klausul memakai labelnya', () => {
    const tahun = baca('pages/InboxTreatyContract.tsx')
    expect(tahun).toContain('opsi={PROPORSI_TCO.map(')
    expect(tahun).not.toContain('PilihJenisReasuransi')
    expect(tahun).toContain('selTahun(labelProporsi(b.proportion))')
    expect(baca('components/PanelKlausulTahun.tsx')).toContain('value={labelProporsi(tahun.proportion)}')
  })
})
