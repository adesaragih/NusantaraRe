import { describe, expect, it } from 'vitest'

import { tambahModulDipinjam } from './daftarMenu'

// Modul tanpa menu sendiri yang dipakai lewat menu modul LAIN (keputusan work owner 09-10-2026: menu Komite Claim Prop
// dibuang, kasus komite dibuka dari inbox Claim Prop) ikut dipasang bagi pemegang menu peminjamnya. Nama modul tiruan.
describe('tambahModulDipinjam', () => {
  const pinjam = { komite: ['klaim'] } as const

  it('pemegang menu peminjam: modul pinjaman ikut dipasang', () => {
    expect(tambahModulDipinjam(['klaim', 'lain'], ['klaim', 'lain', 'komite'], pinjam)).toEqual([
      'klaim',
      'lain',
      'komite',
    ])
  })

  it('tanpa menu peminjam, atau modul pinjaman tidak aktif di backend: tidak dipasang', () => {
    expect(tambahModulDipinjam(['lain'], ['klaim', 'lain', 'komite'], pinjam)).toEqual(['lain'])
    expect(tambahModulDipinjam(['klaim'], ['klaim'], pinjam)).toEqual(['klaim'])
  })

  it('tanpa saringan akun (stub) tetap tanpa saringan; daftar aktif belum terbaca = semua dianggap aktif', () => {
    expect(tambahModulDipinjam(null, ['klaim'], pinjam)).toBeNull()
    expect(tambahModulDipinjam(['klaim'], null, pinjam)).toEqual(['klaim', 'komite'])
  })

  it('modul yang menunya dipegang sendiri tidak digandakan', () => {
    expect(tambahModulDipinjam(['klaim', 'komite'], null, pinjam)).toEqual(['klaim', 'komite'])
  })
})
