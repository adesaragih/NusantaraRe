import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { kolomPanelSummary } from './PanelSummary'

// Panel Summary di Premium List Detail — keputusan work owner 03-10-2026.

const PANEL = readFileSync(join(__dirname, 'PanelSummary.tsx'), 'utf8')
const DETAIL = readFileSync(join(__dirname, '..', 'pages', 'PremiumListDetail.tsx'), 'utf8')

describe('panel Summary', () => {
  it('membaca rekap TERSIMPAN dari server — layar tidak menghitung', () => {
    expect(PANEL).toContain('ambilRekapPolis(polisID)')
    expect(PANEL).not.toMatch(/\bNumber\(|parseFloat|toFixed/)
  })

  it('COB dan PL NUMBER tidak ditampilkan sebelum Confirm', () => {
    const kolom = kolomPanelSummary('QR').map((k) => k.kolom)
    expect(kolom).not.toContain('COB')
    expect(kolom).not.toContain('PL_NUMBER')
    expect(kolom).toContain('CURRENCY')
    expect(kolom).toContain('PREMIUM')
  })

  it('TAX, PROF COMM, dan CLAIM disembunyikan di semua Type (03-10-2026)', () => {
    for (const tipe of ['QR', 'QP', 'TP', 'TR']) {
      const kolom = kolomPanelSummary(tipe).map((k) => k.kolom)
      for (const k of ['TAX', 'PROF_COMM', 'CLAIM']) {
        expect(kolom).not.toContain(k)
      }
      expect(kolom).toContain('BALANCE')
    }
  })

  it('Type asing → tanpa kolom (bukan grid QR)', () => {
    expect(kolomPanelSummary('')).toEqual([])
  })

  it('menjadi TAB Summary di panel Participants, dibaca ulang sesudah setiap Save', () => {
    expect(DETAIL).toContain("{tab === 'summary' && <PanelSummary key={versiSummary} polisID={polisID} />}")
    expect(DETAIL).toContain("{tab === 'rincian' && (")
    expect(DETAIL).toContain('<StripTab tab={TAB_PESERTA}')
    expect(DETAIL.match(/<PanelSummary/g)?.length).toBe(1)
    expect(DETAIL.match(/setVersiSummary\(\(v\) => v \+ 1\)/g)?.length).toBe(2)
  })

  it('isi tab tanpa bingkai panel sendiri — satu panel, bukan dua tabel bertumpuk', () => {
    expect(PANEL).not.toContain('<section className="panel">')
  })
})
