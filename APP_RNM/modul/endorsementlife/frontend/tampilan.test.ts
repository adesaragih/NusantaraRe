import { describe, expect, it } from 'vitest'

import { OPSI_EDM_TYPE, OPSI_TYPE_CEDING } from './labels'
import { UKURAN_HALAMAN_EDM, labelEdmType, labelTypeCeding, sel, selAngka, selTanggal, selWaktu } from './tampilan'

describe('pembantu tampilan Endorsement Life', () => {
  it('sel kosong ditandai, bukan dibiarkan kosong', () => {
    expect(sel('')).toBe('—')
    expect(sel('  ')).toBe('—')
    expect(sel(undefined)).toBe('—')
    expect(sel(' UJI-1 ')).toBe('UJI-1')
  })

  it('tanggal dan waktu dibalik ke DD-MM-YYYY', () => {
    expect(selTanggal('2026-10-01')).toBe('01-10-2026')
    expect(selTanggal('')).toBe('—')
    expect(selWaktu('2026-10-01 09:05:07')).toBe('01-10-2026 09:05:07')
    expect(selWaktu('')).toBe('—')
  })

  it('angka tetap teks: negatif jurnal balik bertanda, presisi tidak hilang', () => {
    expect(selAngka('-1234567.5')).toContain('-')
    expect(selAngka('-1234567.5')).toContain('1.234.567')
    expect(selAngka('0.00000001')).toContain('00000001')
    expect(selAngka('')).toBe('—')
  })

  it('EDM Type hanya 1 dan 3 berlabel (spec §5)', () => {
    expect(Object.keys(OPSI_EDM_TYPE).sort()).toEqual(['1', '3'])
    expect(labelEdmType('1')).toBe('Perubahan Data')
    expect(labelEdmType('3')).toBe('Batal')
    expect(labelEdmType('2')).toBe('2')
  })

  it('Reinsurance System dari domain TYPE_CEDING 1-4', () => {
    expect(Object.keys(OPSI_TYPE_CEDING).sort()).toEqual(['1', '2', '3', '4'])
    expect(labelTypeCeding('3')).toBe('QS+SURPLUS')
    expect(labelTypeCeding('')).toBe('—')
  })

  it('ukuran halaman = services.UkuranHalaman', () => {
    expect(UKURAN_HALAMAN_EDM).toBe(20)
  })
})
