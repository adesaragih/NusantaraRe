// Tempat berperan tiket 05 di layar: kode sama dengan backend, tertunda = tidak tampil.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import { TEMPAT_LABEL_EDM, TEMPAT_LABEL_NON_EDM, TEMPAT_PRODUKSI_TAMPIL, tampilTanggalProduksi, tempatTerbuka } from './tempat'

const GO = readFileSync(join(__dirname, '..', 'backend', 'models', 'peran_tempat.go'), 'utf8')
const hal = (nilai: Record<string, string>): Halaman => ({ nilai, daftar: {} })

describe('tempat berperan (tiket 05)', () => {
  it('kode layar terdaftar di models/peran_tempat.go', () => {
    for (const k of [...TEMPAT_PRODUKSI_TAMPIL, TEMPAT_LABEL_NON_EDM, TEMPAT_LABEL_EDM]) {
      expect(GO).toContain(`"${k}"`)
    }
  })

  it('tertunda (pemetaan kosong) = tidak tampil; tampil hanya bila backend membuka tempatnya', () => {
    expect(tempatTerbuka(undefined, TEMPAT_LABEL_NON_EDM)).toBe(false)
    expect(tempatTerbuka({ [TEMPAT_LABEL_NON_EDM]: false }, TEMPAT_LABEL_NON_EDM)).toBe(false)
    expect(tempatTerbuka({ [TEMPAT_LABEL_NON_EDM]: true }, TEMPAT_LABEL_NON_EDM)).toBe(true)
  })

  it('Production Date: IsApproved 1 DAN salah satu dari dua tempat tampil', () => {
    const disetujui = hal({ 'PolicyTreatyIn.IsApproved': '1' })
    expect(tampilTanggalProduksi(disetujui, {})).toBe(false)
    expect(tampilTanggalProduksi(disetujui, { [TEMPAT_PRODUKSI_TAMPIL[1]]: true })).toBe(true)
    expect(tampilTanggalProduksi(hal({ 'PolicyTreatyIn.IsApproved': '0' }), { [TEMPAT_PRODUKSI_TAMPIL[0]]: true })).toBe(false)
  })
})
