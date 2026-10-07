// Grid popup `BusinessAndSOBListEDM`: urut awal kolom 2 `.ID` DESC (pySortOrder 1), saringan daftar nilai
// (`pyColumnFilteringDropDown`), urut per kolom. Fixture UJI-.

import { describe, expect, it } from 'vitest'

import { URUT_AWAL_BISNIS, nilaiSaringBisnis, saringUrutBisnis, urutBerikut } from './gridbisnis'

const baris = [
  { ID: 'UJI-0001', OLDID: '', Ceding: 'UJI-B' },
  { ID: 'UJI-0003', OLDID: 'UJI-0001', Ceding: 'UJI-A' },
  { ID: 'UJI-0002', OLDID: '', Ceding: 'UJI-B' },
]

describe('grid BusinessAndSOBListEDM', () => {
  it('urut awal ID Revision turun', () => {
    expect(URUT_AWAL_BISNIS).toEqual({ kolom: 'ID', arah: 'turun' })
    expect(saringUrutBisnis(baris, {}, URUT_AWAL_BISNIS).map((b) => b.ID)).toEqual(['UJI-0003', 'UJI-0002', 'UJI-0001'])
  })

  it('saringan daftar nilai: nilai sama persis; kosong = tanpa saringan', () => {
    expect(saringUrutBisnis(baris, { Ceding: 'UJI-B' }, null).map((b) => b.ID)).toEqual(['UJI-0001', 'UJI-0002'])
    expect(saringUrutBisnis(baris, { Ceding: '' }, null)).toHaveLength(3)
    expect(nilaiSaringBisnis(baris, 'Ceding')).toEqual(['UJI-A', 'UJI-B'])
    expect(nilaiSaringBisnis(baris, 'OLDID')).toEqual(['UJI-0001'])
  })

  it('klik judul: kolom lain mulai naik, kolom sama berganti arah', () => {
    expect(urutBerikut(URUT_AWAL_BISNIS, 'Ceding')).toEqual({ kolom: 'Ceding', arah: 'naik' })
    expect(urutBerikut(URUT_AWAL_BISNIS, 'ID')).toEqual({ kolom: 'ID', arah: 'naik' })
    expect(urutBerikut({ kolom: 'ID', arah: 'naik' }, 'ID')).toEqual({ kolom: 'ID', arah: 'turun' })
  })
})
