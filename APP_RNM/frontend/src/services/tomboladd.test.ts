import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { DETAIL, TAHAP } from '../assets/labels.claimlife'
import { bolehAddAdjustment, STATUS_WORK_SELESAI, type Klaim, type Peserta } from './api'

// Tombol `Add` grid adjustment — `ClaimLifeDetailGCNM.xml` b17937.
//
// ⛔ Sejak GILIRAN-14 butir bp `Add` = jalur PUTARAN saja: baris pertama lahir
// saat Submit Register (`SavePesertaClaim` 7.8), bukan lewat `Add` (ralat bo).

function klaimUji(tahap: string, statusWork = ''): Klaim {
  return { tahap, statusWork } as unknown as Klaim
}

function pesertaUji(...kode: string[]): Peserta {
  return {
    kodeStatus: '',
    baris: kode.map((k, i) => ({ id: `UJI-A${i}`, kodeStatus: k })),
  } as unknown as Peserta
}

describe('Add = putaran berikutnya', () => {
  it('ditawarkan di Claim Analis bila baris terakhir ditolak (b18160)', () => {
    expect(bolehAddAdjustment(klaimUji(TAHAP.claimAnalis), pesertaUji('2'))).toBe(true)
    expect(bolehAddAdjustment(klaimUji(TAHAP.claimAnalis), pesertaUji('1', '2'))).toBe(true)
  })

  it('TIDAK ditawarkan di tahap lain — syarat tampilnya pyPosition ReasLifeSPV', () => {
    for (const t of [TAHAP.inputRegister, TAHAP.outstandingClaim, TAHAP.medicalCheck, '']) {
      expect(bolehAddAdjustment(klaimUji(t), pesertaUji('2'))).toBe(false)
    }
  })

  it('TIDAK ditawarkan bila baris terakhir belum ditolak — layanan menjawab 409', () => {
    expect(bolehAddAdjustment(klaimUji(TAHAP.claimAnalis), pesertaUji('0'))).toBe(false)
    expect(bolehAddAdjustment(klaimUji(TAHAP.claimAnalis), pesertaUji('2', '0'))).toBe(false)
  })

  it('TIDAK ditawarkan pada grid kosong — baris pertama lahir saat Register (bp)', () => {
    expect(bolehAddAdjustment(klaimUji(TAHAP.claimAnalis), pesertaUji())).toBe(false)
  })

  it('TIDAK ditawarkan pada kasus tertutup', () => {
    expect(
      bolehAddAdjustment(klaimUji(TAHAP.claimAnalis, STATUS_WORK_SELESAI), pesertaUji('2')),
    ).toBe(false)
  })
})

describe('layar Detail memasang Add dan Delete', () => {
  const BERKAS = readFileSync(join(__dirname, '..', 'pages', 'claimlife', 'KlaimLife.tsx'), 'utf8')
  // Sumber TANPA komentar — prosa yang menerangkan tidak boleh dituduh kode.
  // Blok komentar JSX berbaris banyak dibuang UTUH lebih dulu: baris
  // lanjutannya tidak berawalan penanda apa pun.
  const SUMBER = BERKAS.replace(/\{\/\*[\s\S]*?\*\/\}/g, '')
    .split('\n')
    .filter((b) => {
      const t = b.trimStart()
      return !t.startsWith('//') && !t.startsWith('*') && !t.startsWith('{/*')
    })
    .join('\n')

  it('SATU tombol Add untuk putaran — label VERBATIM, bukan "Putaran berikutnya"', () => {
    expect(SUMBER).toContain('bolehAddAdjustment(klaim, p)')
    expect(SUMBER).toContain('{DETAIL.tambahAdjustment}')
    expect(SUMBER).toContain('void putaranBaru(p.id)')
    expect(SUMBER).not.toContain('Putaran berikutnya')
    expect(SUMBER.match(/void putaranBaru\(p\.id\)/g)).toHaveLength(1)
  })

  it('Delete BERDIRI tetapi mati — ADR-U-0031, nol hapus fisik', () => {
    expect(SUMBER).toContain('{DETAIL.hapusAdjustment}')
    expect(SUMBER).not.toMatch(/hapusAdjustment\s*\(/)
  })

  it('label VERBATIM b17937 dan b19120', () => {
    expect(DETAIL.tambahAdjustment).toBe('Add')
    expect(DETAIL.hapusAdjustment).toBe('Delete')
  })
})
