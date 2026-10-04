// Uji pemilih form MB Capacity (10017) — Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { opsiOccupationMB, saringGrupTreaty } from './PilihMBCapacity'

const KODE = readFileSync(join(__dirname, 'PilihMBCapacity.tsx'), 'utf8')

describe('Occupation MB Capacity', () => {
  it('opsi dari server (SetOccupationLimitMB): ID nilai, nama tampil', () => {
    expect(opsiOccupationMB([{ id: '01', nama: 'RESIDENTIAL RISK' }, { id: '02', nama: 'INDUSTRIAL RISK' }])).toEqual([
      { value: '01', label: 'RESIDENTIAL RISK' },
      { value: '02', label: 'INDUSTRIAL RISK' },
    ])
  })
  it('dropdown inti berteks kosong Choose, dimuat dari klausul-pilihan occupation-limitmb', () => {
    expect(KODE).toContain("cariPilihanKlausul('occupation-limitmb', '')")
    expect(KODE).toContain('kosong={KLAUSUL_TCO.choose}')
  })
})

describe('TerritorialLimit MB Capacity (autocomplete BrowseTreatyGroup_RD)', () => {
  const grup = [
    { id: '10012', treatyGroupName: 'UJI FIRE QS' },
    { id: '10011', treatyGroupName: 'UJI MARINE CARGO' },
    { id: '10010', treatyGroupName: 'UJI FIRE QS' },
  ]
  it('nilai = TreatyGroupName (yang disimpan), ID sebagai keterangan; nama kembar tampil sekali', () => {
    expect(saringGrupTreaty(grup, '')).toEqual([
      { value: 'UJI FIRE QS', label: 'UJI FIRE QS', keterangan: '10012' },
      { value: 'UJI MARINE CARGO', label: 'UJI MARINE CARGO', keterangan: '10011' },
    ])
  })
  it('ketikan menyaring nama (tanpa beda huruf besar-kecil) atau ID', () => {
    expect(saringGrupTreaty(grup, 'marine').map((o) => o.value)).toEqual(['UJI MARINE CARGO'])
    expect(saringGrupTreaty(grup, ' 10012 ').map((o) => o.value)).toEqual(['UJI FIRE QS'])
    expect(saringGrupTreaty(grup, 'tidak ada')).toEqual([])
  })
  it('daftar dari master grup treaty, sekali muat; pemilih yang dapat difilter', () => {
    expect(KODE).toContain('ambilGrupTreaty()')
    expect(KODE).toContain('<PilihSaring')
  })
})
