import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { BarisKomite } from '../api'
import { selKomite, tingkatKomite } from './TabelKomite'

// Tabel komite di bawah inbox Claim Life (perintah work owner 09-10-2026: menu Komite Claim Life dihapus) - perilaku
// Inbox Komite tiket 01 yang dipindah ke sini.

const SUMBER = readFileSync(join(__dirname, 'TabelKomite.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//'))
  .join('\n')

const BARIS: BarisKomite = {
  kasusId: 'KMTLF-UJI1',
  tglUpdate: '',
  statusWork: '',
  klaimId: 'UJI-K',
  nomorKlaim: 'UJI-CLM',
  tingkatBerjalan: 2,
  komiteLoop: 3,
  nilaiKlaim: '1500000.25',
  mataUang: 'IDR',
  statusBaris: 'Outstanding',
}

describe('tabel komite inbox Claim Life', () => {
  it('tingkat berjalan dari seluruh tingkat', () => {
    expect(tingkatKomite(BARIS)).toBe('2 / 3')
  })
  it('sel kosong ditandai', () => {
    expect(selKomite('')).toBe('—')
    expect(selKomite('1500000.25')).toBe('1500000.25')
  })
  it('uang tidak lewat float; layar tidak menyaring ulang; kasus dibuka di tempat modul Komite Claim Life', () => {
    expect(SUMBER).not.toMatch(/Number\(|parseFloat|parseInt/)
    expect(SUMBER).not.toContain('.filter(')
    expect(SUMBER).toContain('onBukaModul?.(MODUL_KOMITE, id)')
  })
})
