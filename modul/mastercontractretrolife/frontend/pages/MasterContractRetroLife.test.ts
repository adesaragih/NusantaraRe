// Halaman awal Master Contract Retro Life: penyimpangan tampilan atas keputusan work owner.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { TAHUN_MCRL } from '../labels'

const KODE = readFileSync(join(__dirname, 'MasterContractRetroLife.tsx'), 'utf8')

describe('halaman awal Master Contract Retro Life', () => {
  it('label sel End Period TIDAK ditampilkan (keputusan work owner 02-10-2026); tombol Add tetap', () => {
    // Labelnya tetap disimpan VERBATIM di labels.ts sebagai bukti XML (b8927, dijaga labels.test).
    expect(TAHUN_MCRL.labelSelAdd).toBe('End Period')
    expect(KODE).not.toContain('{TAHUN_MCRL.labelSelAdd}')
    expect(KODE).not.toContain('className="mcrl-label-sel"')
    expect(KODE).toContain('{TAHUN_MCRL.add}')
    expect(KODE).toContain('title={TAHUN_MCRL.tooltipAdd}')
  })
})
