// Uji definisi medan layar - syarat tampil VERBATIM section Pega (INVENTARIS bab 5).

import { describe, expect, it } from 'vitest'

import type { Halaman } from './api'
import { MEDAN_ADMIN_UANG, MEDAN_ADMIN_UMUM, MEDAN_ATASAN_UANG, MEDAN_ATASAN_UMUM, medanTampil, saldoNegatif } from './medan'

const hal = (nilai: Record<string, string>): Halaman => ({ nilai, daftar: {} })
const label = (ms: { label: string }[]) => ms.map((m) => m.label)

describe('medan layar NB Treaty In', () => {
  it('saldo: "Balance Due To You" bila BalanceDueTo < 0, "Due To Us" bila >= 0 (kosong = 0)', () => {
    expect(saldoNegatif(hal({ 'PolicyTreatyIn.BalanceDueTo': '-0.01' }))).toBe(true)
    expect(saldoNegatif(hal({ 'PolicyTreatyIn.BalanceDueTo': '-0' }))).toBe(false)
    expect(saldoNegatif(hal({}))).toBe(false)
    const neg = label(medanTampil(MEDAN_ADMIN_UANG, hal({ 'PolicyTreatyIn.BalanceDueTo': '-5' })))
    expect(neg).toContain('Balance Due To You')
    expect(neg).not.toContain('Balance Due To Us')
    const pos = label(medanTampil(MEDAN_ADMIN_UANG, hal({ 'PolicyTreatyIn.BalanceDueTo': '5' })))
    expect(pos).toContain('Balance Due To Us')
    expect(pos).toContain('Balance Before Tax')
  })

  it('Type Tax hanya tampil bila FlagPPH true; Survey Report disembunyikan untuk NonProportional', () => {
    expect(label(medanTampil(MEDAN_ADMIN_UMUM, hal({})))).not.toContain('Type Tax')
    expect(label(medanTampil(MEDAN_ADMIN_UMUM, hal({ 'PolicyTreatyIn.FlagPPH': 'true' })))).toContain('Type Tax')
    expect(label(medanTampil(MEDAN_ADMIN_UMUM, hal({ 'Quotation.ProportionalType': 'NonProportional' })))).not.toContain(
      'Survey Report',
    )
  })

  it('elemen mati (1=2 / NEVER) tidak dibangun: Class Of Business dan Due To tidak ada di layar admin', () => {
    const semua = label(MEDAN_ADMIN_UMUM)
    expect(semua).not.toContain('Class Of Business')
    expect(semua).not.toContain('Due To Us / You')
  })

  it('layar atasan: seluruh medan uang terkunci, label VERBATIM (tanpa "(%)")', () => {
    expect(MEDAN_ATASAN_UANG.every((m) => m.jenis === 'tampil')).toBe(true)
    expect(label(MEDAN_ATASAN_UANG)).toContain('Deduction In A (OGP)')
    expect(label(MEDAN_ATASAN_UANG)).not.toContain('(%) Deduction In A (OGP)')
    // hanya FlagPPH dan No Offer Slip yang dapat diisi di bagian umum
    expect(MEDAN_ATASAN_UMUM.filter((m) => m.jenis !== 'tampil').map((m) => m.label)).toEqual(['FlagPPH', 'No Offer Slip'])
  })

  it('angka uang membawa kode mata uang pasangannya (AC 85)', () => {
    const uang = MEDAN_ADMIN_UANG.filter((m) => !m.label.includes('(%)') && !m.label.startsWith('Deduction'))
    expect(uang.every((m) => m.mataUang !== undefined)).toBe(true)
  })
})
