// Syarat tampil subsection NonProp - VERBATIM kontainer section Pega
// (`DetailPolicyTreatyIn`, `DetailPolicyTreatyInNonProportional`, `SpreadingRiskList`).

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import {
  TOTAL_LIMIT,
  TOTAL_SHARE,
  TOTAL_SPREADED,
  jalurRnmShare,
  kolomTotal,
  polisNonPropBaru,
  spreadingTerbuka,
  tampilFakultatif,
  tampilNonProp,
  type Total,
} from './nonprop'

const hal = (nilai: Record<string, string>): Halaman => ({ nilai, daftar: {} })

describe('subsection DetailPoliciesNonProportional', () => {
  it('tampil bila .IsNewPolicyNonProp = 1 && .ClaimType != XOL Retro', () => {
    expect(tampilNonProp(hal({ 'PolicyTreatyIn.IsNewPolicyNonProp': '1' }))).toBe(true)
    expect(tampilNonProp(hal({ 'PolicyTreatyIn.IsNewPolicyNonProp': '0' }))).toBe(false)
    // XOL Retro -> DetailPolicyTreatyOutNonProportional (treaty keluar, tidak dibangun)
    expect(tampilNonProp(hal({ 'PolicyTreatyIn.IsNewPolicyNonProp': '1', 'PolicyTreatyIn.ClaimType': 'XOL Retro' }))).toBe(false)
    expect(polisNonPropBaru(hal({ 'PolicyTreatyIn.IsNewPolicyNonProp': '1', 'PolicyTreatyIn.ClaimType': 'XOL Retro' }))).toBe(true)
  })

  it('Share Facultative tampil bila FacultativeShare > 0; % RNM Share berganti medan', () => {
    expect(tampilFakultatif(hal({ 'TreatyIn.FacultativeShare': '5' }))).toBe(true)
    expect(tampilFakultatif(hal({ 'TreatyIn.FacultativeShare': '0' }))).toBe(false)
    expect(jalurRnmShare(hal({}))).toBe('TreatyIn.RNMShare')
    expect(jalurRnmShare(hal({ 'TreatyIn.FacultativeShare': '5' }))).toBe('TreatyIn.RnmShareDeducted')
  })

  it('tanda FacultativeShare dibaca dari TEKS, tanpa Number/float (spec §5.6)', () => {
    // angka lebih panjang dari presisi double tetap benar tandanya
    expect(tampilFakultatif(hal({ 'TreatyIn.FacultativeShare': '0.000000000000000000001' }))).toBe(true)
    expect(spreadingTerbuka(hal({ 'TreatyIn.FacultativeShare': '0.000000000000000000001' }))).toBe(false)
    expect(spreadingTerbuka(hal({ 'TreatyIn.FacultativeShare': '0.0000' }))).toBe(true)
    expect(tampilFakultatif(hal({ 'TreatyIn.FacultativeShare': '-1' }))).toBe(false)
    // bukan angka: bukan nol, bukan positif (sama dengan NaN sebelumnya)
    expect(tampilFakultatif(hal({ 'TreatyIn.FacultativeShare': 'UJI-X' }))).toBe(false)
    expect(spreadingTerbuka(hal({ 'TreatyIn.FacultativeShare': 'UJI-X' }))).toBe(false)
  })

  it('SpreadingRiskList terbuka (Add/Delete, %Share) hanya bila FacultativeShare = 0 atau kosong', () => {
    expect(spreadingTerbuka(hal({}))).toBe(true)
    expect(spreadingTerbuka(hal({ 'TreatyIn.FacultativeShare': '0' }))).toBe(true)
    expect(spreadingTerbuka(hal({ 'TreatyIn.FacultativeShare': '2.5' }))).toBe(false)
  })
})

// Audit silang P3: judul grid total `DetailPolicyTreatyInNonProportional` - LABEL
// judul berada di atas kolom `.Value` dan kolom `.Currency` tanpa judul
// (TotalLimitDeductblNP: [ "", "Total Deductible" ]; TotalShareDeductionNP:
// [ "", "Total Brokerage", "Total PPN", "Total PPH" ]), KECUALI grid Total Spreaded:
// [ "Total OR Net Premium", "Value" ].
describe('grid total NonProp: judul kolom = LABEL XML', () => {
  const cari = (ts: Total[], d: string): Total => {
    const t = ts.find((x) => x.daftar === d)
    if (!t) throw new Error(`grid total ${d} tidak ada`)
    return t
  }
  const judul = (t: Total) => kolomTotal(t).map((k) => k.label)
  it('label di atas kolom Value, Currency tanpa judul', () => {
    expect(judul(cari(TOTAL_LIMIT, 'TotalLimitDeductblNP'))).toEqual(['', 'Total Deductible'])
    const potongan = cari(TOTAL_SHARE, 'TotalShareDeductionNP')
    expect(judul(potongan)).toEqual(['', 'Total Brokerage', 'Total PPN', 'Total PPH'])
    expect(kolomTotal(potongan).map((k) => k.m)).toEqual(['Currency', 'Value', 'TotalPPNValue', 'TotalPPHValue'])
  })
  it('Total Spreaded: label di atas Currency, "Value" di atas Value', () => {
    expect(judul(cari(TOTAL_SPREADED, 'TotalSpreadedNetPremi'))).toEqual(['Total OR Net Premium', 'Value'])
    expect(judul(cari(TOTAL_SPREADED, 'TotalSpreadedNetPremiRI'))).toEqual(['Total R/I Net Premium', 'Value'])
  })
})
