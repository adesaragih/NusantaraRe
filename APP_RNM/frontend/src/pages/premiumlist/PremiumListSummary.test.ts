import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { GRID_REKAP } from '../../assets/labels.premiumlist'
import type { RekapMataUangPolis } from '../../services/api'
import { kolomGridRekap, selRekap } from './PremiumListSummary'

// Uji layar rekap premium list — tiket 05a bagian 2.

const BERKAS = readFileSync(join(__dirname, 'PremiumListSummary.tsx'), 'utf8')
const SUMBER = BERKAS.split('\n')
  .filter((b) => {
    const t = b.trimStart()
    return !t.startsWith('//') && !t.startsWith('*')
  })
  .join('\n')

const REKAP: RekapMataUangPolis = {
  currency: 'IDR',
  premium: '1000.0000',
  commission: '0.0000',
  balance: '975.0000',
  jumlah: { DEDUCTION: '10.0000', TAX: '' },
  cacahBaris: 1,
}

describe('grid rekap per Type', () => {
  it('empat grid, satu per Type, VERBATIM urut section', () => {
    expect(Object.keys(GRID_REKAP).sort()).toEqual(['QP', 'QR', 'TP', 'TR'])
    expect(kolomGridRekap('QR').map((k) => k.kolom)).toEqual([
      'COB', 'PL_NUMBER', 'CURRENCY', 'PREMIUM', 'DEDUCTION', 'BROKERAGE_FEE',
      'RI_ADMIN_FEE', 'TAX', 'PROF_COMM', 'CLAIM', 'BALANCE',
    ])
    expect(kolomGridRekap('QP')).toHaveLength(16)
    expect(kolomGridRekap('TP')).toHaveLength(13)
    expect(kolomGridRekap('TR')).toHaveLength(13)
  })

  it('keanehan warisan: PREMIUM DEDUCTION berdiri di atas COMMISSION pada QP/TP/TR', () => {
    for (const t of ['QP', 'TP', 'TR']) {
      const k = kolomGridRekap(t).find((x) => x.judul === 'PREMIUM DEDUCTION')
      expect(k?.kolom).toBe('COMMISSION')
    }
    expect(kolomGridRekap('QR').some((x) => x.kolom === 'COMMISSION')).toBe(false)
  })

  it('Type asing tidak meminjam grid QR', () => {
    expect(kolomGridRekap('FAC')).toEqual([])
    expect(kolomGridRekap('')).toEqual([])
  })
})

describe('sel rekap', () => {
  it('uang tampil apa adanya, ekor nol tidak dibuang', () => {
    expect(selRekap(REKAP, 'BALANCE', 'PL1')).toBe('975.0000')
    expect(selRekap(REKAP, 'PREMIUM', 'PL1')).toBe('1000.0000')
    expect(selRekap(REKAP, 'DEDUCTION', 'PL1')).toBe('10.0000')
  })

  it('kosong ditandai, COB tidak ditebak', () => {
    expect(selRekap(REKAP, 'TAX', 'PL1')).toBe('—')
    expect(selRekap(REKAP, 'COB', 'PL1')).toBe('—')
    expect(selRekap(REKAP, 'PL_NUMBER', '')).toBe('—')
    expect(selRekap(REKAP, 'PL_NUMBER', 'PL1')).toBe('PL1')
  })
})

describe('penjaga statik', () => {
  it('nol Number(...) — uang tidak lewat float', () => {
    expect(SUMBER).not.toMatch(/Number\(|parseFloat|parseInt/)
  })

  it('sesudah Submit kasus tertutup dan tombolnya mati (05b)', () => {
    expect(SUMBER).toContain('setSelesai(true)')
    expect(SUMBER).toMatch(/disabled=\{[^}]*selesai/)
  })

  it('rekap tidak dihitung di layar', () => {
    expect(SUMBER).not.toContain('GROSS_PREMIUM_RETRO +')
    expect(SUMBER).toContain('ambilRekapPolis')
    expect(SUMBER).toContain('submitRekapPolis')
  })
})
