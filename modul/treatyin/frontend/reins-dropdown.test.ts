// Reinsurer Name (grid Reinsurer / Facultative Reinsurers) = DROPDOWN, bukan
// autocomplete — permintaan pemakai 9 Oktober 2026.

import { describe, expect, it } from 'vitest'

import { nilaiReins, opsiReins, pilihReins } from './components/TabShareNonProp'

const daftar = [
  { id: '10', nama: 'EVEREST REINSURANCE COMPANY', kembar: false },
  { id: '21', nama: 'REASURANSI MAIPARK INDONESIA', kembar: true },
  { id: '22', nama: 'REASURANSI MAIPARK INDONESIA', kembar: true },
]

describe('dropdown reasuradur', () => {
  it('nilai = ReinsID; nama kembar diberi pengenalnya', () => {
    const b = { ReinsID: '10', ReinsName: 'EVEREST REINSURANCE COMPANY' }
    expect(nilaiReins(b, daftar)).toBe('10')
    expect(opsiReins(b, daftar).map((o) => o.label)).toEqual([
      'EVEREST REINSURANCE COMPANY',
      'REASURANSI MAIPARK INDONESIA (21)',
      'REASURANSI MAIPARK INDONESIA (22)',
    ])
  })
  it('memilih mengisi ReinsID DAN ReinsName', () => {
    expect(pilihReins({ ReinsID: '', ReinsName: '', Layer: '1' }, daftar, '22')).toEqual({
      ReinsID: '22', ReinsName: 'REASURANSI MAIPARK INDONESIA', Layer: '1',
    })
    expect(pilihReins({ ReinsID: '10', ReinsName: 'X' }, daftar, '')).toEqual({ ReinsID: '', ReinsName: '' })
  })
  it('baris lama tanpa ReinsID dicocokkan lewat nama; yang tak dikenal tetap tampil', () => {
    expect(nilaiReins({ ReinsID: '', ReinsName: 'EVEREST REINSURANCE COMPANY' }, daftar)).toBe('10')
    const lama = { ReinsID: '99', ReinsName: 'PT LAMA' }
    expect(opsiReins(lama, daftar)[0]).toEqual({ value: '99', label: 'PT LAMA' })
  })
})
