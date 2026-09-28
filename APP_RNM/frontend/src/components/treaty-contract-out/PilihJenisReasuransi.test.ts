// Uji pemilih jenis reasuransi — tiket 02 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { opsiJenisReasuransi } from './PilihJenisReasuransi'

const SUMBER = readFileSync(join(__dirname, 'PilihJenisReasuransi.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//'))
  .join('\n')

describe('opsi dari daftar server', () => {
  it('nilai = ID sebagai TEKS, label = Note', () => {
    const opsi = opsiJenisReasuransi([
      { id: '10003', note: 'UJI QS', tipe: '1' },
      { id: '00007', note: 'UJI SURPLUS', tipe: '2' },
    ])
    expect(opsi).toEqual([
      { value: '10003', label: 'UJI QS' },
      // ⛔ Nol di depan TETAP: ID adalah kode, bukan bilangan (ADR-U-0022).
      { value: '00007', label: 'UJI SURPLUS' },
    ])
  })
  it('daftar kosong menjadi opsi kosong, bukan karangan', () => {
    expect(opsiJenisReasuransi([])).toEqual([])
  })
})

describe('layar tidak menyaring dan tidak mengarang', () => {
  it('nol .filter( — saringannya di server', () => {
    expect(SUMBER).not.toContain('.filter(')
  })
  it('nol literal ID blacklist di React', () => {
    for (const id of ['10004', '10011', '10012', '10021', '10022', '10025', '10026', '10028', '10248', '10249', '10018', '10217']) {
      expect(SUMBER).not.toContain(id)
    }
  })
  it('memakai Pilih dari ui/dasar dan label dari berkas label', () => {
    expect(SUMBER).toContain("from '../ui/dasar'")
    expect(SUMBER).toContain("from '../../assets/labels.treaty-contract-out'")
  })
  it('ID tidak lewat Number()', () => {
    expect(SUMBER).not.toMatch(/Number\(|parseInt|parseFloat/)
  })
})
