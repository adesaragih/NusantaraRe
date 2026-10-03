// Uji syarat tampil pemilih Source Of Business - VERBATIM XML:
//   tombol `Select Source Of Business` (`Section/DetailPolicyTreatyIn`) pyVisible `.ClaimType = 'XOL Retro'`
//   tombol `Choose` (`Section/SourceHierarki`) pyVisible `.ChildCount = 0`

import { describe, expect, it } from 'vitest'

import type { Halaman } from '../api'
import { tampilChoose, tampilTombolSOB } from './PilihSumberBisnis'

const hal = (nilai: Record<string, string>): Halaman => ({ nilai, daftar: {} })

describe('pemilih Source Of Business', () => {
  it('tombol hanya tampil bila ClaimType persis "XOL Retro"', () => {
    expect(tampilTombolSOB(hal({ 'PolicyTreatyIn.ClaimType': 'XOL Retro' }))).toBe(true)
    expect(tampilTombolSOB(hal({ 'PolicyTreatyIn.ClaimType': 'XOL' }))).toBe(false)
    expect(tampilTombolSOB(hal({}))).toBe(false)
  })

  it('Choose hanya untuk simpul tanpa anak (ChildCount = 0; kosong = 0)', () => {
    expect(tampilChoose({ id: 'UJI-1', clientName: 'UJI', leader0: '', childCount: '0', clientId: '' })).toBe(true)
    expect(tampilChoose({ id: 'UJI-2', clientName: 'UJI', leader0: '', childCount: '', clientId: '' })).toBe(true)
    expect(tampilChoose({ id: 'UJI-3', clientName: 'UJI', leader0: '', childCount: '2', clientId: '' })).toBe(false)
  })
})
