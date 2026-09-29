// Uji aturan status simpan utuh — tiket 09 Treaty Contract Out.

import { describe, expect, it } from 'vitest'

import { statusSimpanSukses } from './api'

describe('status simpan (AC 39)', () => {
  it('HANYA "1" sukses; kosong, null, undefined, dan nilai lain gagal', () => {
    expect(statusSimpanSukses('1')).toBe(true)
    for (const s of ['', '0', ' 1', '1 ', '01', '2', 'sukses', null, undefined]) {
      expect(statusSimpanSukses(s), String(s)).toBe(false)
    }
  })
})

describe('hapus kontrak mengirim cacah kontrak lain (OQ-TCO-21)', () => {
  it('kueri DELETE memuat bersama', async () => {
    const { readFileSync } = await import('node:fs')
    const { join } = await import('node:path')
    const kode = readFileSync(join(__dirname, 'api.ts'), 'utf8')
    expect(kode).toContain('bersama: String(d.bersama)')
  })
})
