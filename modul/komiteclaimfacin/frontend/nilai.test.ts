// Nilai layar Komite Claim Fac In (pola `modul/komiteclaimnonprop/frontend/nilai.test.ts`): format angka / tanggal,
// validasi isian ShowTransfer LS45, dan isian yang dikirim (dua Propose hanya TT2 tingkat 1).

import { describe, expect, it } from 'vitest'

import type { Keputusan } from './api'
import { isianKirim, periksa, PESAN_KOSONG, tampil, tercentang, usulTerbuka } from './nilai'

const kosong: Keputusan = { acceptStatus: '', comment: '', usulTutup: false, usulCadang: false }

describe('tampil - angka paling banyak 4 desimal pemisah Indonesia, tanggal dd-mm-yyyy', () => {
  it('angka: pemisah ribuan, 4 desimal maksimum, nol ekor dibuang', () => {
    expect(tampil('1234567.5', 'angka')).toBe('1.234.567,5')
    expect(tampil('0.123456', 'angka')).toBe('0,1235')
    expect(tampil('57750000.0000', 'angka')).toBe('57.750.000')
    expect(tampil('-1500.25', 'angka')).toBe('-1.500,25')
    expect(tampil('', 'angka')).toBe('')
  })
  it('tanggal dan tanggal-jam', () => {
    expect(tampil('2026-09-15', 'tanggal')).toBe('15-09-2026')
    expect(tampil('2026-10-10 10:05:00', 'tanggalJam')).toBe('10-10-2026 10:05')
    expect(tampil('2026-10-10T10:05:00+07:00', 'tanggalJam')).toBe('10-10-2026 10:05')
    expect(tampil('UJI', 'tanggal')).toBe('UJI')
  })
  it('teks dan tautan apa adanya', () => {
    expect(tampil(' UJI ', 'teks')).toBe('UJI')
    expect(tampil('retro:1:1:2', 'tautan')).toBe('retro:1:1:2')
  })
  it('centang hanya-baca (pxCheckbox): "true" / "1" dicentang', () => {
    expect(tercentang('true')).toBe(true)
    expect(tercentang('1')).toBe(true)
    expect(tercentang('false')).toBe(false)
    expect(tercentang('')).toBe(false)
    expect(tampil('1', 'centang')).toBe('✓')
    expect(tampil('0', 'centang')).toBe('')
  })
})

describe('isian ShowTransfer', () => {
  it('wajib isi: AcceptStatus (1 / 2) dan Note', () => {
    expect(periksa(kosong)).toEqual({ acceptStatus: PESAN_KOSONG, comment: PESAN_KOSONG })
    expect(periksa({ ...kosong, acceptStatus: '3', comment: '  ' })).toEqual({
      acceptStatus: PESAN_KOSONG,
      comment: PESAN_KOSONG,
    })
    expect(periksa({ ...kosong, acceptStatus: '2', comment: 'UJI' })).toEqual({})
  })

  it('dua Propose terbuka hanya bila tampil (TT2) dan terbuka (tingkat 1)', () => {
    expect(usulTerbuka({ tampilUsul: true, terbuka: true })).toBe(true)
    expect(usulTerbuka({ tampilUsul: true, terbuka: false })).toBe(false)
    expect(usulTerbuka({ tampilUsul: false, terbuka: true })).toBe(false)
  })

  it('dua Propose tersembunyi / nonaktif tidak membawa nilai baru', () => {
    const k: Keputusan = { acceptStatus: '1', comment: 'UJI', usulTutup: true, usulCadang: true }
    expect(isianKirim(k, { tampilUsul: true, terbuka: true })).toEqual(k)
    expect(isianKirim(k, { tampilUsul: true, terbuka: false })).toEqual({ ...k, usulTutup: false, usulCadang: false })
    expect(isianKirim(k, { tampilUsul: false, terbuka: false })).toEqual({ ...k, usulTutup: false, usulCadang: false })
  })
})
