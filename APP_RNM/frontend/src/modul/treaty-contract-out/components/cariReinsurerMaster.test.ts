import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { opsiReinsurerMaster } from './cariReinsurerMaster'

const KODE = readFileSync(join(__dirname, 'cariReinsurerMaster.ts'), 'utf8')

describe('pencarian master reinsurer (Reinsurer + Security Name)', () => {
  it('butir: nama, ID sebagai keterangan; tanpa nama = ID', () => {
    expect(opsiReinsurerMaster({ id: '10001', clientName: 'UJI REAS SATU', clientId: '' })).toEqual({
      value: '10001', label: 'UJI REAS SATU', keterangan: '10001',
    })
    expect(opsiReinsurerMaster({ id: '10003', clientName: ' ', clientId: '' }).label).toBe('10003')
  })
  it('hanya jawaban TERAKHIR yang dipakai - sukses maupun gagal', () => {
    const then = KODE.slice(KODE.indexOf('.then('), KODE.indexOf('.catch('))
    const tangkap = KODE.slice(KODE.indexOf('.catch('), KODE.indexOf('[onGalat]'))
    expect(then).toContain('if (ke !== urutan.current) return')
    expect(tangkap).toContain('if (ke !== urutan.current) return')
  })
  it('reset membatalkan jawaban yang masih di jalan', () => {
    const reset = KODE.slice(KODE.indexOf('const reset'))
    expect(reset.slice(0, 120)).toContain('urutan.current++')
  })
})
