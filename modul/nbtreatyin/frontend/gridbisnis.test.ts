// Grid AKTIF popup `Section/BusinessAndSOBList` (RD `BrowseTreatyJoinEDM`) - audit silang P3 W1.
// Harapan diketik dari XML Section (pyGridProps, pyColumnDefaults) dan dihitung tangan:
//
//   pyGridSorting true, pyGridFiltering true, pyPageMode Numeric, pyPageSize 50
//   kolom 1 (tombol Choose): pyColumnSorting false, pyColumnFilteringDropDown false
//   kolom 2-26: pyColumnSorting true, pyColumnFilteringDropDown true
//   kolom 2 (.TREATYID): pySortType ASC, pySortOrder 1
//   sel .LIMITVALUE/.RETENTIONVALUE/.EPIVALUE/.MDPVALUE/.NETPREMIVALUE/.SHAREVALUE: pxCurrency

import { describe, expect, it } from 'vitest'

import { BARIS_PER_HALAMAN_BISNIS, URUT_AWAL_BISNIS, nilaiSaringBisnis, saringUrutBisnis, urutBerikut } from './gridbisnis'

const b = (ID: string, isi: Record<string, string>) => ({ ID, ...isi })

describe('grid popup BusinessAndSOBList', () => {
  it('paging 50 baris (pyPageSize), urut awal Treaty Offer ID naik (kolom 2 pySortType ASC)', () => {
    expect(BARIS_PER_HALAMAN_BISNIS).toBe(50)
    expect(URUT_AWAL_BISNIS).toEqual({ kolom: 'TREATYID', arah: 'naik' })
  })

  it('saring per kolom: nilai sama persis, beberapa kolom = DAN', () => {
    const baris = [
      b('1', { TREATYTYPE: 'XOL', LIMITCURRENCY: 'IDR' }),
      b('2', { TREATYTYPE: 'XOL', LIMITCURRENCY: 'USD' }),
      b('3', { TREATYTYPE: 'Quota Share', LIMITCURRENCY: 'IDR' }),
      b('4', { TREATYTYPE: 'XOL Retro', LIMITCURRENCY: 'IDR' }),
    ]
    const id = (x: { ID: string }[]) => x.map((r) => r.ID)
    expect(id(saringUrutBisnis(baris, { TREATYTYPE: 'XOL' }, null))).toEqual(['1', '2'])
    expect(id(saringUrutBisnis(baris, { TREATYTYPE: 'XOL', LIMITCURRENCY: 'IDR' }, null))).toEqual(['1'])
    // saringan kosong = kolom itu tidak disaring
    expect(id(saringUrutBisnis(baris, { TREATYTYPE: '' }, null))).toEqual(['1', '2', '3', '4'])
  })

  it('urut teks: urutan kode karakter (sama dengan ORDER BY RD), stabil, naik/turun', () => {
    const baris = [b('1', { TREATYID: 'UJI-T2' }), b('2', { TREATYID: 'UJI-T10' }), b('3', { TREATYID: 'UJI-T2' }), b('4', { TREATYID: '' })]
    const id = (x: { ID: string }[]) => x.map((r) => r.ID)
    // '' < 'UJI-T10' < 'UJI-T2' (kode '1' < '2'); seri 1 dan 3 tetap urutan datang
    expect(id(saringUrutBisnis(baris, {}, { kolom: 'TREATYID', arah: 'naik' }))).toEqual(['4', '2', '1', '3'])
    expect(id(saringUrutBisnis(baris, {}, { kolom: 'TREATYID', arah: 'turun' }))).toEqual(['1', '3', '2', '4'])
  })

  it('urut kolom pxCurrency menurut nilai desimal, tanpa float', () => {
    const baris = [
      b('a', { LIMITVALUE: '1000' }),
      b('b', { LIMITVALUE: '200' }),
      b('c', { LIMITVALUE: '-5' }),
      b('d', { LIMITVALUE: '9007199254740993.5' }),
      b('e', { LIMITVALUE: '9007199254740993.25' }),
      b('f', { LIMITVALUE: '' }),
      b('g', { LIMITVALUE: '.5' }),
    ]
    const id = (x: { ID: string }[]) => x.map((r) => r.ID)
    // kosong di depan; -5 < .5 < 200 < 1000 < ...993.25 < ...993.5
    expect(id(saringUrutBisnis(baris, {}, { kolom: 'LIMITVALUE', arah: 'naik' }))).toEqual(['f', 'c', 'g', 'b', 'a', 'e', 'd'])
    // kolom teks bukan uang tetap diurut sebagai teks: '1000' < '200'
    const teks = [b('a', { LAYER: '1000' }), b('b', { LAYER: '200' })]
    expect(id(saringUrutBisnis(teks, {}, { kolom: 'LAYER', arah: 'naik' }))).toEqual(['a', 'b'])
  })

  it('pilihan saringan satu kolom: nilai berbeda selain kosong, berurut seperti kolomnya', () => {
    // "" adalah pilihan "tanpa saringan", bukan nilai yang disaring
    const baris = [b('1', { LIMITVALUE: '1000' }), b('2', { LIMITVALUE: '200' }), b('3', { LIMITVALUE: '1000' }), b('4', {})]
    expect(nilaiSaringBisnis(baris, 'LIMITVALUE')).toEqual(['200', '1000'])
    expect(nilaiSaringBisnis([b('1', { SOB: 'UJI B' }), b('2', { SOB: 'UJI A' })], 'SOB')).toEqual(['UJI A', 'UJI B'])
  })

  it('klik judul kolom: kolom lain -> naik; kolom yang sama berganti naik/turun', () => {
    expect(urutBerikut({ kolom: 'TREATYID', arah: 'naik' }, 'SOB')).toEqual({ kolom: 'SOB', arah: 'naik' })
    expect(urutBerikut({ kolom: 'SOB', arah: 'naik' }, 'SOB')).toEqual({ kolom: 'SOB', arah: 'turun' })
    expect(urutBerikut({ kolom: 'SOB', arah: 'turun' }, 'SOB')).toEqual({ kolom: 'SOB', arah: 'naik' })
    expect(urutBerikut(null, 'SOB')).toEqual({ kolom: 'SOB', arah: 'naik' })
  })
})
