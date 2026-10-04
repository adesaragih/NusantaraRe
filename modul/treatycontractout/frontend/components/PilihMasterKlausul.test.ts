import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { opsiMasterKlausul } from './PilihMasterKlausul'

const KODE = readFileSync(join(__dirname, 'PilihMasterKlausul.tsx'), 'utf8')

// 10013 Exclusion Treaty [keputusan work owner 02-10-2026]: dropdown ID Occupation / ID Clause berfungsi
// tanpa kotak Search - terisi begitu dibuka, disaring di server.
describe('PilihMasterKlausul', () => {
  it('opsi: nama tampil, ID sebagai keterangan', () => {
    expect(opsiMasterKlausul([{ id: 'UJI-1', nama: 'UJI OCCUPATION' }])).toEqual([
      { value: 'UJI-1', label: 'UJI OCCUPATION', keterangan: 'UJI-1' },
    ])
  })
  it('dimuat tanpa syarat panjang ketikan: kata kosong = daftar awal', () => {
    expect(KODE).toContain('cariPilihanKlausul(master, kata.trim())')
    expect(KODE).not.toMatch(/length\s*<\s*2/)
  })
  it('memilih mengisi ID dan nama', () => {
    expect(KODE).toContain('onPilih(o.value, o.label)')
  })
})
