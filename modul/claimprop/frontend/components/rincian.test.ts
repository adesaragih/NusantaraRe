import { describe, expect, it } from 'vitest'

import { DAFTAR_RINCI, rincianUntuk, type RincianGrid } from './rincian'

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
