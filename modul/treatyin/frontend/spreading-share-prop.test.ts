// Tab Share Prop — rincian Detail (`Section/DetailShare.xml`) disamakan
// dengan layar Pega pemakai.
//
// ⭐ KEPUTUSAN PEMILIK PROSES 8 Oktober 2026: Spreading DIPILIH dari master —
// *"pctnya di ambil dari table PROPORTIONALARRG ... select nya menggunakan
// rd"*. Dropdown `Spreading Type` (RD `BrowseTreatyArrangement_ParentReins
// MasterTrt`) selalu tampil; grid Reins Type · Pct baca-saja; grid manual
// (`SetSpreadName`, `AddDelSpreadingTreatyin`) tidak dirender lagi.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { selAngka } from './components/angka'
import { DETAIL_SHARE } from './labelsShareProp'

const SRC = readFileSync(join(__dirname, 'components', 'TabShareProp.tsx'), 'utf8')
const i = SRC.indexOf('function RincianDetailShare')
const RINCIAN = SRC.slice(i, SRC.indexOf('export default function TabShareProp', i))

describe('⭐ Spreading dipilih dari RD, Pct dari PROPORTIONALARRG', () => {
  it('dropdown Spreading Type tampil di mode Edit TANPA syarat Spreading Type terisi', () => {
    expect(RINCIAN).toContain("useIndukSpreading(teksDari(d, 'TreatyGroupID'), commencement, !terkunci)")
    expect(RINCIAN).toMatch(/\{terkunci \? \(\s*<p className="trin__teks-sel">\s*\{DETAIL_SHARE\.jenisSpreading\}/)
    expect(RINCIAN).not.toContain("tipe !== '' ? (")
  })

  it('pilihan → `FetchQSfromMaster` (aksi `spreading`); dikosongkan → nama ikut kosong', () => {
    expect(RINCIAN).toContain("onRumus('spreading', { ...d, SpreadingTypeID: v, ...(v === '' ? { SpreadingType: '' } : {}) })")
  })

  it('grid Reins Type · Pct BACA-SAJA — nol isian, nol Add/Delete, nol `SetSpreadName`', () => {
    expect([...DETAIL_SHARE.kolomSpreading]).toEqual(['Reins Type', 'Pct'])
    expect(RINCIAN).toContain("<td>{teksDari(r, 'ReinsTypeName')}</td>")
    expect(RINCIAN).toContain("<td>{teksDari(r, 'Pct')}</td>")
    for (const dicabut of ['sebar-nama', 'TombolTambah', 'TombolHapus', 'inputMode="decimal"', '<select']) {
      expect(RINCIAN).not.toContain(dicabut)
    }
  })

  it('`Total Spreading Pct :` lalu Value Spreading OR / R/I', () => {
    expect(DETAIL_SHARE.totalSpreadingPct).toBe('Total Spreading Pct :')
    const t = RINCIAN.indexOf('DETAIL_SHARE.totalSpreadingPct')
    expect(RINCIAN.indexOf('DETAIL_SHARE.sebaranOR', t)).toBeGreaterThan(t)
    expect(RINCIAN.indexOf('DETAIL_SHARE.sebaranRI', t)).toBeGreaterThan(t)
  })
})

// Perbandingan layar Pega 7 Oktober 2026 — format dari ekspor.
describe('⭐ tampilan tab Share Prop = layar Pega', () => {
  it('`% RNM Share` Treaty Group: `pxNumber` 2 desimal TANPA %, lalu ShareNote — `25,00 of 80% of 100%`', () => {
    expect(selAngka(['uang', 2], '25') + ' of 80% of 100%').toBe('25,00 of 80% of 100%')
    expect(SRC).toContain("{selAngka(['uang', 2], teksDari(d, 'RNMShare'))}")
    expect(SRC).not.toContain("selAngka(['persen', 2], teksDari(d, 'RNMShare'))")
  })

  it('grid rincian Detail (RNM Share, Value Spreading OR/R/I) berkepala kolom kedua KOSONG', () => {
    expect((SRC.match(/judulNilai=""/g) ?? []).length).toBe(3)
  })

  it('`% RNM Share` / `% Brokerage` kepala: tanpa simbol, placeholder `%`', () => {
    expect((SRC.match(/desimal=\{2\} placeholder="%"/g) ?? []).length).toBe(2)
    expect(SRC).not.toMatch(/persenRnmShare\} value=\{rnmShare\} desimal=\{2\} persen/)
  })
})
