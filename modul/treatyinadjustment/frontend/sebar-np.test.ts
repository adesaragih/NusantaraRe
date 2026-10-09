// ⭐ DUA spreading Share Non-Prop layar Adjustment — persis Pega `Share.xml`
// (keputusan pemakai 9 Oktober 2026, *"di pega ada terdapat 2 spreading"*):
//
//   `.SpreadingTypeXOL != ''` → spreading LAMA: dropdown Spreading Type + grid
//                               `readOnly` hasil `FetchQSfromMasterXOL`, nol tombol.
//   `.SpreadingTypeXOL == ''` → spreading BARU: grid manual `AddSpreadingXOL`
//                               (Add/Delete), sel `SetSpreadingXOL`, Total Share Pct.
//
// Dulu grid manual NP terbangkit `readOnly` tanpa tombol (kolom ketiga tertulis
// `.pyTemplateInputBox`) — pembangkit hanya mengenali grid manual Prop.

import { describe, expect, it } from 'vitest'

import type { ButirKerangka, GridKerangka } from './ekspor/jenis'
import { KERANGKA_RINCIAN } from './ekspor/kerangka.gen'
import { hapusBarisGrid, rantaiSesudahHapus, tambahDari } from './komponen/aksiTombol'

function gridDalam(isi: readonly ButirKerangka[], syarat: string): GridKerangka[] {
  const out: GridKerangka[] = []
  const jalan = (xs: readonly ButirKerangka[], dalam: boolean) => {
    for (const b of xs) {
      const kena = dalam || b.syarat.some((s) => s.replace(/\s/g, '') === syarat)
      if (b.t === 'grid' && kena && b.larik === 'SpreadingListXOL') out.push(b)
      if (b.t === 'blok') jalan(b.anak, kena)
    }
  }
  jalan(isi, false)
  return out
}

describe('Share NP Adjustment — spreading lama vs baru', () => {
  const share = KERANGKA_RINCIAN.Share ?? []

  it("spreading LAMA (`.SpreadingTypeXOL!=''`): grid baca, nol tombol", () => {
    const [g] = gridDalam(share, ".SpreadingTypeXOL!=''")
    expect(g).toBeDefined()
    expect(g?.tombol.every((t) => t === null)).toBe(true)
    expect(g?.tombolKepala.every((t) => t === null)).toBe(true)
    expect(g?.baca.every((b) => b === 'selalu')).toBe(true)
  })

  it("spreading BARU (`.SpreadingTypeXOL==''`): Add/Delete `AddSpreadingXOL`, sel dapat diubah", () => {
    const [g] = gridDalam(share, ".SpreadingTypeXOL==''")
    expect(g).toBeDefined()
    const tambah = g?.tombolKepala.find((t) => t !== null)
    const hapus = g?.tombol.find((t) => t !== null)
    expect(tambah?.label).toBe('Add')
    expect(hapus?.label).toBe('Delete')
    expect(tambah && tambahDari(tambah)).toMatchObject({ larik: 'SpreadingListXOL' })
    expect(hapus && hapusBarisGrid(hapus)).toBe(true)
    expect(g?.kolom).toEqual(['Reins Type', 'Pct Share', ''])
    expect(g?.baca.slice(0, 2)).toEqual([["TreatyIn.ViewState = '1'"], ["TreatyIn.ViewState = '1'"]])
  })

  it('Delete lalu CountTotalPctSpreadXOL — Σ .Pct baris tersisa', async () => {
    const [g] = gridDalam(share, ".SpreadingTypeXOL==''")
    const hapus = g?.tombol.find((t) => t !== null)
    const rantai = hapus ? rantaiSesudahHapus(hapus) : undefined
    expect(typeof rantai).toBe('function')
    const h = await rantai?.({ medan: {}, larik: { SpreadingListXOL: [{ Pct: '40' }, { Pct: '35.5' }] } } as never, {} as never)
    expect(h?.medan).toMatchObject({ SpreadingTotalPctXOL: '75.5' })
  })
})
