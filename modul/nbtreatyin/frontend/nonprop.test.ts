// Syarat tampil subsection NonProp - VERBATIM kontainer section Pega
// (`DetailPolicyTreatyIn`, `DetailPolicyTreatyInNonProportional`, `SpreadingRiskList`).

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import { jalurRnmShare, polisNonPropBaru, spreadingTerbuka, tampilFakultatif, tampilNonProp } from './nonprop'

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

  it('SpreadingRiskList terbuka (Add/Delete, %Share) hanya bila FacultativeShare = 0 atau kosong', () => {
    expect(spreadingTerbuka(hal({}))).toBe(true)
    expect(spreadingTerbuka(hal({ 'TreatyIn.FacultativeShare': '0' }))).toBe(true)
    expect(spreadingTerbuka(hal({ 'TreatyIn.FacultativeShare': '2.5' }))).toBe(false)
  })
})
