import { describe, expect, it } from 'vitest'

import { barisTerbuka, bukaAwal, DAFTAR_RINCI, rincianUntuk, type RincianGrid } from './rincian'

const panel: RincianGrid = { daftar: DAFTAR_RINCI, isi: () => null }

describe('panel rinci baris grid', () => {
  it('hanya grid Adjustment yang punya expand pane (Section AdjustmentDetail)', () => {
    expect(DAFTAR_RINCI).toBe('ClaimData.AdjustmentList')
    expect(rincianUntuk(panel, 'ClaimData.AdjustmentList')).toBe(panel)
  })

  it('grid lain tidak dapat dibuka - laporan work owner 08-10-2026', () => {
    expect(rincianUntuk(panel, 'ClaimData.InterestList')).toBeUndefined()
    expect(rincianUntuk(panel, 'ClaimData.TotalInterestInsured')).toBeUndefined()
    expect(rincianUntuk(panel, undefined)).toBeUndefined()
    expect(rincianUntuk(undefined, DAFTAR_RINCI)).toBeUndefined()
  })
})

describe('baris ber-panel terbuka tanpa klik - work owner 08-10-2026', () => {
  it('baris terbaru terbuka sejak awal; tanpa baris tidak ada yang terbuka', () => {
    expect(bukaAwal(3)).toBe(3)
    expect(bukaAwal(1)).toBe(1)
    expect(bukaAwal(0)).toBeNull()
  })

  it('Add membuka baris baru; hapus menutup baris yang hilang; lainnya tetap', () => {
    expect(barisTerbuka(1, 1, 2)).toBe(2)
    expect(barisTerbuka(null, 0, 1)).toBe(1)
    expect(barisTerbuka(2, 2, 1)).toBe(1)
    expect(barisTerbuka(1, 1, 0)).toBeNull()
    expect(barisTerbuka(1, 3, 2)).toBe(1)
    expect(barisTerbuka(null, 2, 2)).toBeNull()
  })
})
