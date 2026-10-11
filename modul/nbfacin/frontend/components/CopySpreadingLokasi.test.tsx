// Copy Spreading per lokasi (tiket 73 P5, K73-40/50/51) - label & opsi Pega, aturan Type Facultative, pemasangan di tab Spreading.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { OPSI_FLAG_SPREADING, SPREADING_LOKASI as L } from '../labels'
import CopySpreadingLokasi, { fakultatifKurang, TREATY_WAJIB_FAKULTATIF } from './CopySpreadingLokasi'

const SUMBER = readFileSync(join(__dirname, 'CopySpreadingLokasi.tsx'), 'utf8').replace(/\r\n/g, '\n')
const FIRE = readFileSync(join(__dirname, 'TabSpreading.tsx'), 'utf8').replace(/\r\n/g, '\n')
const ANEKA = readFileSync(join(__dirname, 'TabSpreadingAneka.tsx'), 'utf8').replace(/\r\n/g, '\n')
const API = readFileSync(join(__dirname, '..', 'api.ts'), 'utf8').replace(/\r\n/g, '\n')

const baris = (treatyType: string, flagSpreading = '') => ({ treatyType, treatyName: '', sharePercentage: '10', flagSpreading })

describe('CopySpreadingLokasi', () => {
  it('label kolom / tombol = sel Pega; opsi Type Facultative = DDL\\FlagSpreading.xml', () => {
    expect(L.kolom).toEqual(['Type Treaty', '% Share', 'Type Facultative'])
    expect([L.tambah, L.hapus, L.salin]).toEqual(['Add', 'Delete', 'Copy To All Spreading'])
    expect(OPSI_FLAG_SPREADING).toEqual(['SF-MR 2023', 'Open Fac Market'])
    const h = renderToStaticMarkup(<CopySpreadingLokasi treaty={[]} menghitung={false} onSalin={() => {}} />)
    for (const k of L.kolom) expect(h).toContain(`>${k}<`)
    expect(h).toContain('>Copy To All Spreading<')
  })

  it('Type Facultative wajib hanya untuk treaty 10015 (L3014 / L3020)', () => {
    expect(TREATY_WAJIB_FAKULTATIF).toBe('10015')
    expect(fakultatifKurang([baris('10015')])).toBe(true)
    expect(fakultatifKurang([baris('10015', 'SF-MR 2023')])).toBe(false)
    expect(fakultatifKurang([baris('10007')])).toBe(false)
  })

  it('ResetFlagSpreading_ACT: Type Treaty berubah mengosongkan Type Facultative semua baris', () => {
    expect(SUMBER).toMatch(/template\.map\(\(x, k\) => \(\{\s*\.\.\.x,\s*flagSpreading: '',/)
  })

  it('terpasang di detail lokasi FIRE dan Aneka, tahap UW; template case di layar dipertahankan (jawaban template: [])', () => {
    expect(FIRE).toContain('salinLokasiSpreading(caseId, share.trim(), o + 1, t).then((h) => ({ ...h, template: v.template }))')
    expect(ANEKA).toContain('salinLokasiSpreadingAneka(caseId, share.trim(), a + 1, t).then((h) => ({ ...h, template: v.template }))')
    expect(FIRE).toMatch(/\{tahapUw && \(\s*<CopySpreadingLokasi/)
    expect(ANEKA).toMatch(/\{tahapUw && \(\s*<CopySpreadingLokasi/)
    expect(API).toContain('`${jalurSpreading(caseId)}/salin-lokasi`, { metode: \'POST\', badan: { percentShare, idxLokasi, template } }')
  })
})
