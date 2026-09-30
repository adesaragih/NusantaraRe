// Uji panel business — tiket 07 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { BUSINESS_TCO } from '../labels'
import { formBusinessDari, formBusinessKosong, keMasukBusiness, labelAktif } from './PanelBusinessKombinasi'

const KODE = readFileSync(join(__dirname, 'PanelBusinessKombinasi.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
  .join('\n')

describe('form business', () => {
  it('baris baru kosong; Active wajib dipilih', () => {
    expect(formBusinessKosong()).toEqual({ id: '', bizCode: '', isActive: '' })
  })
  it('Edit menyalin ID, kode, status', () => {
    expect(formBusinessDari({
      id: '1000002', isActive: '0', aktif: false, treatyYear: '2026', treatyYearId: '1000001', treatyGroupId: '10001',
      treatyGroupName: 'UJI GRUP', reinsTypeId: '10003', reinsTypeName: 'UJI QS', bizCode: 'UJI-B1',
      bizName: 'UJI BISNIS', userId: 'UJI-ADMIN', tglUpdate: '',
    })).toEqual({ id: '1000002', bizCode: 'UJI-B1', isActive: '0' })
  })
  it('nama bisnis TIDAK dikirim — server mengambilnya dari master', () => {
    expect(keMasukBusiness({ id: '', bizCode: ' UJI-B1 ', isActive: '1' })).toEqual({ id: '', bizCode: 'UJI-B1', isActive: '1' })
  })
  it('status nonaktif tetap terbaca sebagai nonaktif (AC 22)', () => {
    expect(labelAktif({ aktif: true })).toBe(BUSINESS_TCO.aktif)
    expect(labelAktif({ aktif: false })).toBe(BUSINESS_TCO.nonaktif)
  })
})

describe('paritas layar business', () => {
  it('seluruh label dari BUSINESS_TCO', () => {
    for (const k of ['judul', 'add', 'kolomTreatyGroup', 'kolomBusinessId', 'kolomBusinessName', 'edit', 'delete',
      'formBusinessName', 'formActive', 'formBusinessCode', 'save', 'information', 'closeList'] as const) {
      expect(KODE).toContain(`BUSINESS_TCO.${k}`)
    }
  })
  it('grid tidak menyaring baris nonaktif', () => {
    expect(KODE).not.toMatch(/daftar\.daftar\.filter/)
    expect(KODE).toContain('labelAktif(b)')
  })
  it('kolom Business ID menampilkan ID BARIS (b4228), bukan kode bisnis', () => {
    expect(KODE).toMatch(/<td>\{b\.id\}<\/td>/)
  })
})
