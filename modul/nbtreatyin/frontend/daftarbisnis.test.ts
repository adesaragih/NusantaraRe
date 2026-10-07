// Popup Choose Business: saringan kolom dikirim ke SERVER bersama isian layar (keputusan work owner 06-10-2026),
// supaya kontrak di luar 500 baris pertama (RD pyMaxRecords 500) dapat ditemukan.

import { beforeEach, describe, expect, it, vi } from 'vitest'

const minta = vi.fn((..._a: unknown[]) => Promise.resolve([]))
vi.mock('../../../inti/frontend/klien', () => ({ minta: (...a: unknown[]) => minta(...a) }))

import { daftarBisnis, type Halaman } from './api'

const h: Halaman = { nilai: { 'PolicyTreatyIn.QuotationData.ProportionalType': 'Proportional' }, daftar: {} }

describe('daftarBisnis', () => {
  beforeEach(() => minta.mockClear())

  it('mengirim saringan kolom bersama halaman', async () => {
    await daftarBisnis('UJI-K1', h, { TREATYID: '1002059' })
    expect(minta).toHaveBeenCalledTimes(1)
    const [jalur, opsi] = minta.mock.calls[0] as [string, { metode: string; badan: unknown }]
    expect(jalur).toMatch(/\/kasus\/UJI-K1\/bisnis$/)
    expect(opsi.metode).toBe('POST')
    expect(opsi.badan).toEqual({ halaman: h, saringan: { TREATYID: '1002059' } })
  })

  it('tanpa saringan: objek kosong', async () => {
    await daftarBisnis('UJI-K1', h)
    const [, opsi] = minta.mock.calls[0] as [string, { badan: unknown }]
    expect(opsi.badan).toEqual({ halaman: h, saringan: {} })
  })
})
