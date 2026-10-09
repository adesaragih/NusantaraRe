// Berkas Resolve Complete / Decline = layar BACA — laporan pemakai 9 Oktober
// 2026: sesudah approve, isian masih dapat diubah.

import { describe, expect, it } from 'vitest'

import { modeEfektif } from './pages/PenyesuaianKontrak'

describe('modeEfektif — status tuntas mengunci', () => {
  it('Resolve Complete / Decline → baca walau dibuka lewat Edit', () => {
    expect(modeEfektif('0', false, { StatusAkseptasi: 'Resolve Complete' })).toBe('1')
    expect(modeEfektif('0', false, { StatusAkseptasi: 'Decline' })).toBe('1')
  })
  it('status lain (termasuk kosong sesudah Revision) mengikuti tombol', () => {
    expect(modeEfektif('0', false, { StatusAkseptasi: 'Accept' })).toBe('0')
    expect(modeEfektif('0', false, { StatusAkseptasi: '', RevisionState: '1' })).toBe('0')
    expect(modeEfektif('1', false, {})).toBe('1')
  })
})
