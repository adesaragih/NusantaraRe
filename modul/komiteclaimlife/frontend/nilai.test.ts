import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { BarisInboxKomite } from './api'
import { selKomite, tingkatKomite } from './nilai'

// Uji nilai kasus komite — tiket 01 Komite Claim Life (dipindah dari uji Inbox Komite).

const SUMBER = readFileSync(join(__dirname, 'pages', 'KasusKomite.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//'))
  .join('\n')

const BARIS: BarisInboxKomite = {
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

describe('nilai kasus komite', () => {
  it('tingkat berjalan dari seluruh tingkat', () => {
    expect(tingkatKomite(BARIS)).toBe('2 / 3')
  })
  it('sel kosong ditandai', () => {
    expect(selKomite('')).toBe('—')
    expect(selKomite('1500000.25')).toBe('1500000.25')
  })
  it('uang tidak lewat float; layar tidak menyaring ulang', () => {
    expect(SUMBER).not.toMatch(/Number\(|parseFloat|parseInt/)
    expect(SUMBER).not.toContain('.filter(')
  })
})
