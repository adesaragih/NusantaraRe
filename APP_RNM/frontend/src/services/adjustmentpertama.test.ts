import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { DETAIL, TAHAP } from '../assets/labels.claimlife'
import { bolehAddAdjustment, STATUS_WORK_SELESAI, type Klaim, type Peserta } from './api'

// Butir bo (GILIRAN-13) — `Add` b17937 pada grid adjustment KOSONG.

function klaimUji(tahap: string, statusWork = ''): Klaim {
  return { tahap, statusWork } as unknown as Klaim
}

function pesertaUji(kodeStatus = '', cacahBaris = 0): Peserta {
  return {
    kodeStatus,
    baris: Array.from({ length: cacahBaris }, (_, i) => ({ id: `UJI-A${i}` })),
  } as unknown as Peserta
}

describe('Add baris adjustment pertama', () => {
  it('ditawarkan di Claim Analis pada grid kosong (b18160)', () => {
    expect(bolehAddAdjustment(klaimUji(TAHAP.claimAnalis), pesertaUji())).toBe(true)
  })

  it('TIDAK ditawarkan di tahap lain — syarat tampilnya pyPosition ReasLifeSPV', () => {
    for (const t of [TAHAP.inputRegister, TAHAP.outstandingClaim, TAHAP.medicalCheck, '']) {
      expect(bolehAddAdjustment(klaimUji(t), pesertaUji())).toBe(false)
    }
  })

  it('TIDAK ditawarkan pada grid berbaris — di sana Add adalah jalur putaran', () => {
    expect(bolehAddAdjustment(klaimUji(TAHAP.claimAnalis), pesertaUji('', 1))).toBe(false)
  })

  it('TIDAK ditawarkan bagi peserta yang sudah diputus — tujuh gerbang STS_REJECT', () => {
    for (const kode of ['1', '2']) {
      expect(bolehAddAdjustment(klaimUji(TAHAP.claimAnalis), pesertaUji(kode))).toBe(false)
    }
    expect(bolehAddAdjustment(klaimUji(TAHAP.claimAnalis), pesertaUji('0'))).toBe(true)
  })

  it('TIDAK ditawarkan pada kasus tertutup', () => {
    expect(
      bolehAddAdjustment(klaimUji(TAHAP.claimAnalis, STATUS_WORK_SELESAI), pesertaUji()),
    ).toBe(false)
  })
})

describe('layar Detail memasang Add dan Delete', () => {
  const BERKAS = readFileSync(join(__dirname, '..', 'pages', 'claimlife', 'KlaimLife.tsx'), 'utf8')
  // Sumber TANPA komentar — prosa yang menerangkan tidak boleh dituduh kode.
  const SUMBER = BERKAS.split('\n')
    .filter((b) => {
      const t = b.trimStart()
      return !t.startsWith('//') && !t.startsWith('*') && !t.startsWith('{/*')
    })
    .join('\n')

  it('Add memakai rute putaran yang sama, bergerbang bolehAddAdjustment', () => {
    expect(SUMBER).toContain('bolehAddAdjustment(klaim, p)')
    expect(SUMBER).toContain('{DETAIL.tambahAdjustment}')
    expect(SUMBER).toContain('void putaranBaru(p.id)')
  })

  it('Delete BERDIRI tetapi mati — ADR-U-0031, nol hapus fisik', () => {
    expect(SUMBER).toContain('{DETAIL.hapusAdjustment}')
    // ⛔ Tidak ada pemanggil hapus baris adjustment di layar ini.
    expect(SUMBER).not.toMatch(/hapusAdjustment\s*\(/)
    expect(SUMBER).not.toMatch(/\/adjustment\/[^'"`]*['"`],\s*\{\s*metode:\s*'DELETE'/)
  })

  it('label VERBATIM b17937 dan b19120', () => {
    expect(DETAIL.tambahAdjustment).toBe('Add')
    expect(DETAIL.hapusAdjustment).toBe('Delete')
  })
})
