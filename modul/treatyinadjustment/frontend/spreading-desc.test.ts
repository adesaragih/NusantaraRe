// Dropdown Reins Type / Spreading Type Adjustment mengirim `TreatyDescID`
// seperti ekspor (`"10001"`) — tanpa itu `ORS` dan susunan non-QS ikut
// muncul (laporan pemakai 9 Oktober 2026).

import { describe, expect, it, vi } from 'vitest'

import type { SumberPilihan } from './ekspor/jenis'
import { opsiUntuk, type PemuatOpsi } from './komponen/pilihan'

const REINS: SumberPilihan = {
  sumber: 'reportdefinition',
  rd: 'BrowseTreatyArrangement_ParentReinsMasterTrt',
  nilai: 'ReinsTypeID',
  tampil: 'ReinsTypeName',
  param: { TreatyYear: 'TreatyIn.TreatyYear', TreatyGroupID: 'TempSprd.TreatyGroupID', TreatyDescID: '"10001"', StartDate: 'TreatyIn.Commencement', ReinsTypeID: '"10246"' },
}

describe('spreading-induk Adjustment = parameter ekspor', () => {
  it('TreatyDescID 10001 ikut dikirim', async () => {
    const muat = vi.fn()
    const pemuat: PemuatOpsi = {
      kelasBisnis: vi.fn(() => Promise.resolve([])),
      indukSpreading: vi.fn(() => Promise.resolve([{ reinsTypeId: '10300', reinsTypeName: '2026 QS 212M TRT HR' }])),
    }
    const tempat = { halaman: { Commencement: '20260101' }, sisi: { medan: {}, larik: {} } }
    expect(opsiUntuk(REINS, 'ReinsTypeID', { muat } as never, tempat, pemuat)).toEqual([])
    expect(muat).toHaveBeenCalledWith('spreading||20260101|10001', expect.any(Function))
    const [, ambil] = muat.mock.calls[0] as [string, () => Promise<unknown>]
    await ambil()
    expect(pemuat.indukSpreading).toHaveBeenCalledWith('', '20260101', '10001')
  })
})

describe('TempSprd.TreatyGroupID = Treaty Group Detail (AddDelSpreadingTreatyin [1])', () => {
  it('dropdown Reins Type disaring grup Detail — ORS grup lain tidak ikut', () => {
    const muat = vi.fn()
    const pemuat: PemuatOpsi = {
      kelasBisnis: vi.fn(() => Promise.resolve([])),
      indukSpreading: vi.fn(() => Promise.resolve([])),
    }
    // Sel grid SpreadingList: `sisi` = Detail (Treaty Group 10015).
    const tempat = { halaman: { Commencement: '20250701' }, sisi: { medan: { TreatyGroupID: '10015' }, larik: {} } }
    opsiUntuk(REINS, 'ReinsTypeID', { muat } as never, tempat, pemuat)
    expect(muat).toHaveBeenCalledWith('spreading|10015|20250701|10001', expect.any(Function))
  })
})
