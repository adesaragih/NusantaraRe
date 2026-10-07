import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { pemisahRibuan } from './angka'
import { selAngkaAtauTanggal } from './pages/PremiumListDetail'

// Pemisah ribuan untuk nominal uang — keputusan work owner 03-10-2026.

describe('pemisahRibuan', () => {
  it('titik pemisah ribuan, koma desimal, desimal tidak diubah (03-10-2026)', () => {
    expect(pemisahRibuan('2200000.50')).toBe('2.200.000,50')
    expect(pemisahRibuan('23400155')).toBe('23.400.155')
    expect(pemisahRibuan('16380108.5')).toBe('16.380.108,5')
    expect(pemisahRibuan('558.1988625')).toBe('558,1988')
    expect(pemisahRibuan('130000867.12345678')).toBe('130.000.867,1234')
    expect(pemisahRibuan('-1234567.89')).toBe('-1.234.567,89')
    expect(pemisahRibuan('0')).toBe('0')
  })

  it('teks bukan desimal dibiarkan apa adanya', () => {
    for (const t of ['', '—', 'IDR', '31/10/2026', 'CIDCLSPPP001000', '1,234']) {
      expect(pemisahRibuan(t)).toBe(t)
    }
  })

  it('tanpa Number() — uang tetap teks (ADR-U-0003)', () => {
    const sumber = readFileSync(join(__dirname, 'angka.ts'), 'utf8')
    expect(sumber).not.toMatch(/\bNumber\(|parseFloat|toFixed|toLocaleString/)
  })
})

describe('sel grid peserta', () => {
  it('hanya kolom ANGKA (dari server) yang diberi pemisah', () => {
    expect(selAngkaAtauTanggal('23400155', true)).toBe('23.400.155')
    // Nomor sertifikat yang kebetulan berupa angka TIDAK ikut berpemisah.
    expect(selAngkaAtauTanggal('12345678', false)).toBe('12345678')
    expect(selAngkaAtauTanggal('2025-03-09 00:00:00', false)).toBe('09/03/2025')
    expect(selAngkaAtauTanggal(undefined, true)).toBeUndefined()
  })
})

describe('paling banyak 4 digit desimal (03-10-2026)', () => {
  it('DIPOTONG, bukan dibulatkan; desimal pendek tidak ditambah nol', () => {
    expect(pemisahRibuan('3101.105663')).toBe('3.101,1056')
    expect(pemisahRibuan('1.99999')).toBe('1,9999')
    expect(pemisahRibuan('20.1234')).toBe('20,1234')
    expect(pemisahRibuan('20.5')).toBe('20,5')
  })
})
