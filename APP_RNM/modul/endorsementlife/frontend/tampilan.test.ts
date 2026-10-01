import { describe, expect, it } from 'vitest'

import { OPSI_EDM_TYPE, OPSI_TYPE_CEDING } from './labels'
import {
  PILIHAN_KOSONG,
  UKURAN_HALAMAN_EDM,
  alihBaris,
  alihSemua,
  bolehCentang,
  labelEdmType,
  labelTypeCeding,
  sel,
  selAngka,
  selTanggal,
  selWaktu,
  tercentang,
} from './tampilan'

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

describe('kotak centang .EdmBatal dan DELETE ALL (SelectAllEdmLife_act)', () => {
  it('DELETE ALL sakelar: kosong → semua, semua → kosong (b259)', () => {
    const semua = alihSemua(PILIHAN_KOSONG)
    expect(semua).toEqual({ pilih: [], semua: true, kecuali: [] })
    expect(tercentang(semua, 'UJI-D9')).toBe(true)
    expect(alihSemua(alihBaris(semua, 'UJI-D1'))).toEqual(PILIHAN_KOSONG)
    // Centang satu per satu lalu DELETE ALL: seluruh baris tercentang, pilihan lama dibuang.
    expect(alihSemua(alihBaris(PILIHAN_KOSONG, 'UJI-D1'))).toEqual({ pilih: [], semua: true, kecuali: [] })
  })

  it('centang baris: daftar pilihan, atau pengecualian sesudah DELETE ALL', () => {
    const satu = alihBaris(PILIHAN_KOSONG, 'UJI-D1')
    expect(satu.pilih).toEqual(['UJI-D1'])
    expect(tercentang(satu, 'UJI-D1')).toBe(true)
    expect(tercentang(satu, 'UJI-D2')).toBe(false)
    expect(alihBaris(satu, 'UJI-D1')).toEqual(PILIHAN_KOSONG)
    const kecuali = alihBaris(alihSemua(PILIHAN_KOSONG), 'UJI-D2')
    expect(kecuali).toEqual({ pilih: [], semua: true, kecuali: ['UJI-D2'] })
    expect(tercentang(kecuali, 'UJI-D2')).toBe(false)
    expect(tercentang(alihBaris(kecuali, 'UJI-D2'), 'UJI-D2')).toBe(true)
  })

  it('hanya peserta Old pada kasus Perubahan Data terbuka yang belum disimpan (R30, b15763, b37200)', () => {
    const kasus = { status: '', sudahSimpan: false, edmType: '1' }
    expect(bolehCentang({ edmStatus: 'Old' }, kasus)).toBe(true)
    for (const s of ['New', 'Delete', 'Batal', '']) expect(bolehCentang({ edmStatus: s }, kasus), s).toBe(false)
    expect(bolehCentang({ edmStatus: 'Old' }, { ...kasus, sudahSimpan: true })).toBe(false)
    expect(bolehCentang({ edmStatus: 'Old' }, { ...kasus, edmType: '3' })).toBe(false)
    expect(bolehCentang({ edmStatus: 'Old' }, { ...kasus, status: 'Resolved-Completed' })).toBe(false)
  })
})
