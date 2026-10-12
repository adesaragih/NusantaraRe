import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { pemisahRibuan } from '../angka'
import type { SpreadingSatuPeserta } from '../api'
import { LABEL_TAB_PESERTA, PANEL_SPREADING, TAB_PESERTA } from '../labels'
import { barisTabelSpreading, judulModalSpreading, KOLOM_SPREADING, treatyYearTampil } from './ModalSpreadingPeserta'

// Popup spreading satu peserta (keputusan work owner 07-10-2026) dan rekap retro di Summary.

const MODAL = readFileSync(join(__dirname, 'ModalSpreadingPeserta.tsx'), 'utf8')
const SUMMARY = readFileSync(join(__dirname, 'PanelSummary.tsx'), 'utf8')
const API = readFileSync(join(__dirname, '..', 'api.ts'), 'utf8')

function peserta(lapisan: SpreadingSatuPeserta['lapisan']): SpreadingSatuPeserta {
  return {
    id: 'UJI-D1', nama: 'UJI-A', sertifikat: 'UJI-C1', currency: 'IDR', retrocededShare: '4500000000',
    treatyYear: '2026', treatyMulai: '2026-01-01', treatyAkhir: '2026-12-31', lapisan,
  }
}

describe('tab Spreading dihapus (07-10-2026)', () => {
  it('panel Participants hanya Participant Details dan Summary', () => {
    expect([...TAB_PESERTA]).toEqual(['rincian', 'summary'])
    expect(Object.keys(LABEL_TAB_PESERTA)).toEqual(['rincian', 'summary'])
  })

  it('endpoint semua peserta diganti endpoint satu peserta', () => {
    expect(API).not.toContain('/spreading${q}')
    expect(API).not.toContain('ambilSpreadingPolis')
    expect(API).toContain('/peserta/${encodeURIComponent(detailID)}/spreading')
    expect(MODAL).toContain('ambilSpreadingPeserta(polisID, detailID)')
  })
})

describe('isi popup spreading', () => {
  it('lapisan → reinsurer; Reins Type dan bagian hanya di baris pertama lapisan', () => {
    const baris = barisTabelSpreading(peserta([
      { reinsType: 'QS', treatyYear: 'T2', bagian: '3000000000.0000', reinsurer: [
        { nama: 'UJI-RE-A', pctShare: '60', amount: '1800000000.0000', rate: '2.5000', komisi: '20', ovrComm: '10', gross: '1', net: '1' },
        { nama: 'UJI-RE-B', pctShare: '40', amount: '1200000000.0000', rate: '2.5000', komisi: '20', ovrComm: '10', gross: '1', net: '1' },
      ] },
      { reinsType: 'SURPLUS', treatyYear: 'T2', bagian: '1500000000.0000', reinsurer: [
        { nama: 'UJI-RE-C', pctShare: '100', amount: '1500000000.0000', rate: '4.0000', komisi: '0', ovrComm: '5', gross: '1', net: '1' },
      ] },
    ]))
    expect(baris.map((b) => [b.reinsType, b.bagian, b.reinsurer, b.awalLapisan])).toEqual([
      ['QS', '3.000.000.000,0000', 'UJI-RE-A', true],
      ['', '', 'UJI-RE-B', false],
      ['SURPLUS', '1.500.000.000,0000', 'UJI-RE-C', true],
    ])
    expect(baris[0]?.amount).toBe('1.800.000.000,0000')
    expect(KOLOM_SPREADING.map((c) => PANEL_SPREADING.kolom[c.k])).toEqual([
      'Reins Type', 'Retroceded Share', 'Reinsurer', '% Share', 'Amount', 'Rate', 'Commision %', 'Ovr Comm %', 'Gross', 'Net',
    ])
  })

  it('SATU formatter: Rate selalu berkoma desimal, apa pun bentuk teks server', () => {
    const lap = (rate: string) => [{ reinsType: 'QS', treatyYear: 'T2', bagian: '1', reinsurer: [
      { nama: 'UJI-RE', pctShare: '100', amount: '1', rate, komisi: '0', ovrComm: '0', gross: '1', net: '1' }] }]
    for (const rate of ['0.72', '.72', '+0.72', '7.2E-01', ' 0.72 ']) {
      expect(barisTabelSpreading(peserta(lap(rate)))[0]?.rate).toBe('0,72')
    }
    // Sama persis dengan formatter grid peserta.
    expect(MODAL).toContain("import { pemisahRibuan } from '../angka'")
    expect(MODAL).not.toMatch(/\bNumber\(|parseFloat|toFixed|toLocaleString/)
    expect(MODAL).toContain('value={pemisahRibuan(isi.retrocededShare)}')
    expect(pemisahRibuan('4500000000')).toBe('4.500.000.000')
  })

  it('tanpa spreading → pesan kosong, bukan galat', () => {
    expect(barisTabelSpreading(peserta([]))).toEqual([])
    expect(barisTabelSpreading(peserta(null))).toEqual([])
    expect(PANEL_SPREADING.kosong).toBe('No retro spreading for this participant.')
    expect(MODAL).toContain('<Kosong pesan={PANEL_SPREADING.kosong} />')
  })

  it('judul "Spreading - <Name Of Insured>", kepala dan treaty year', () => {
    expect(judulModalSpreading('UJI-A')).toBe('Spreading - UJI-A')
    expect(judulModalSpreading('  ')).toBe('Spreading')
    expect(treatyYearTampil(peserta([]))).toBe('2026 (01/01/2026 - 31/12/2026)')
    expect(treatyYearTampil({ treatyYear: '', treatyMulai: '', treatyAkhir: '' })).toBe('')
    for (const k of ['nama', 'sertifikat', 'currency', 'retrocededShare', 'treatyYear'] as const) {
      expect(MODAL).toContain(`label={k.${k}}`)
    }
  })

  it('Close dan Esc menutup: Modal inti dengan onTutup pemanggil', () => {
    expect(PANEL_SPREADING.tutup).toBe('Close')
    expect(MODAL).toContain('onTutup={onTutup} labelBatal={PANEL_SPREADING.tutup}')
    expect(MODAL).toContain('<div className="pl-spreading">')
  })
})

describe('rekap retro di Summary', () => {
  it('per reinsurer per mata uang dari data, tanpa daftar nama tetap', () => {
    expect(SUMMARY).toContain('ambilRekapRetroPolis(polisID)')
    expect(SUMMARY).toContain('{PANEL_SPREADING.judulRekap}')
    for (const nama of ['NUSARE', 'HANNOVER', 'TOA', 'GIBRALTAR']) {
      expect(SUMMARY.toUpperCase()).not.toContain(nama)
    }
    expect(PANEL_SPREADING.kolomRekap.gross).toBe('Premium Spreaded Gross')
    expect(API).toContain('/spreading/rekap')
  })
})
