// Rumus kepala form panel New — `TreatyInSetTreatyYear` dan
// `TreatyCalculateProratePct`, disalin dari DataTransform ekspor.

import { describe, expect, it } from 'vitest'

import { ikutBerubah, proRata, tahunTreaty } from './komponen/rumusKepala'

describe('TreatyInSetTreatyYear — Commencement berubah', () => {
  it('Treaty Year = 4 digit tahun; Termination = +1 tahun, bentuk tersimpan', () => {
    expect(tahunTreaty({ Commencement: '20180716' })).toEqual({ TreatyYear: '2018', Termination: '20190716' })
    // Masukan kotak tanggal (`YYYY-MM-DD`) juga dikenali.
    expect(tahunTreaty({ Commencement: '2021-03-01' })).toEqual({ TreatyYear: '2021', Termination: '20220301' })
  })

  it('⛔ 29 Februari + 1 tahun = 28 Februari (`@addCalendar`), bukan 1 Maret', () => {
    expect(tahunTreaty({ Commencement: '20240229' })).toEqual({ TreatyYear: '2024', Termination: '20250228' })
  })

  it('tanggal tak terbaca → nol perubahan', () => {
    expect(tahunTreaty({ Commencement: '' })).toEqual({})
  })
})

describe('TreatyCalculateProratePct — Effective / Is Pro Rate berubah', () => {
  it('Days = Effective → Termination; Total = Commencement → Termination; % = @divide dua desimal × 100', () => {
    expect(proRata({ IsProRate: 'true', Commencement: '20200301', EDMEffective: '20210301', Termination: '20210331' })).toEqual({
      ProRateDays: '30',
      ProRateTotalDays: '395',
      // 30/395 = 0,0759… → @divide 0,08 → 8,00
      ProRatePercent: '8.00',
    })
  })

  it('⭐ sama dengan penyesuaian terbaru yang tersimpan: 134/364 → 37,00 dan 90/455 → 20,00', () => {
    expect(proRata({ IsProRate: 'true', Commencement: '20230101', EDMEffective: '20230820', Termination: '20231231' })).toMatchObject({
      ProRateDays: '133',
      ProRateTotalDays: '364',
    })
    const p = proRata({ IsProRate: 'true', Commencement: '20230101', EDMEffective: '20230819', Termination: '20231231' })
    expect(p).toEqual({ ProRateDays: '134', ProRateTotalDays: '364', ProRatePercent: '37.00' })
  })

  it('Is Pro Rate bukan true → ketiganya 0', () => {
    expect(proRata({ IsProRate: 'false', Commencement: '20200301', EDMEffective: '20210301', Termination: '20210331' })).toEqual({
      ProRateDays: '0',
      ProRateTotalDays: '0',
      ProRatePercent: '0',
    })
  })

  it('⛔ Termination TIDAK punya event — mengubahnya tidak menghitung ulang apa pun', () => {
    expect(ikutBerubah('Termination', { IsProRate: 'true', Commencement: '20200301', EDMEffective: '20210301', Termination: '20210331' })).toEqual({})
    expect(Object.keys(ikutBerubah('EDMEffective', { IsProRate: 'true' }))).toEqual(['ProRateDays', 'ProRateTotalDays', 'ProRatePercent'])
  })
})
