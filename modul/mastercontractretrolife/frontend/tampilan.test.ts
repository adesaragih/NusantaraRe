import { describe, expect, it } from 'vitest'

import { UKURAN_HALAMAN_MCRL, jepitHalaman, keTanggalKabel, potongHalaman, sel, selAngka, selTanggal, selWaktu, waktuKini } from './tampilan'

describe('penomoran 10 baris (pyRDLPageSize 10)', () => {
  const baris = Array.from({ length: 23 }, (_, i) => i + 1)

  it('ukuran halaman = 10', () => {
    expect(UKURAN_HALAMAN_MCRL).toBe(10)
  })

  it('memotong per halaman, halaman terakhir sisanya', () => {
    expect(potongHalaman(baris, 1)).toEqual([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    expect(potongHalaman(baris, 3)).toEqual([21, 22, 23])
    expect(potongHalaman(baris, 0)).toEqual(potongHalaman(baris, 1))
  })

  it('halaman dijepit sesudah daftar menyusut (baris terakhir halaman terhapus)', () => {
    expect(jepitHalaman(3, 20)).toBe(2)
    expect(jepitHalaman(3, 0)).toBe(1)
    expect(jepitHalaman(2, 23)).toBe(2)
  })
})

describe('sel', () => {
  it('kosong ditandai (ADR-U-0027)', () => {
    expect(sel('')).toBe('—')
    expect(sel('  ')).toBe('—')
    expect(sel('UJI-1')).toBe('UJI-1')
  })

  it('tanggal YYYY-MM-DD → DD/MM/YYYY', () => {
    expect(selTanggal('2026-01-31')).toBe('31/01/2026')
    expect(selTanggal('')).toBe('—')
  })

  it('waktu YYYY-MM-DD HH:MM:SS → DD/MM/YYYY HH:MM:SS', () => {
    expect(selWaktu('2026-09-30 10:05:07')).toBe('30/09/2026 10:05:07')
    expect(selWaktu('')).toBe('—')
  })

  it('angka lewat digit, bukan float - 17 digit tetap utuh', () => {
    // Number('12345678901234567.25') = 12345678901234568 - float kehilangan digitnya.
    expect(selAngka('12345678901234567.25')).toBe('12.345.678.901.234.567,25')
    expect(selAngka('33.333333')).toBe('33,333333')
    expect(selAngka('')).toBe('—')
  })
})

describe('tanggal isian ke kabel', () => {
  it('DD-MM-YYYY dan YYYY-MM-DD → YYYY-MM-DD', () => {
    expect(keTanggalKabel('31-01-2026')).toBe('2026-01-31')
    expect(keTanggalKabel('2026-01-31')).toBe('2026-01-31')
  })

  it('bukan tanggal dikirim apa adanya - server yang menolak dengan kalimatnya', () => {
    expect(keTanggalKabel('besok')).toBe('besok')
    expect(keTanggalKabel('  ')).toBe('')
  })
})

describe('@CurrentDateTime()', () => {
  it('YYYY-MM-DD HH:MM:SS waktu setempat', () => {
    expect(waktuKini(new Date(2026, 8, 30, 7, 4, 9))).toBe('2026-09-30 07:04:09')
  })
})
