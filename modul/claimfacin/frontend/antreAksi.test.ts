// Disalin dari `modul/claimnonprop/frontend/antreAksi.test.ts` (pola, bukan impor; asal Claim Prop): keputusan work
// owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Laporan work owner 09-10-2026: sesudah sebuah medan diisi, klik tombol tidak membuka popup - harus klik dua kali. Klik
// itu memicu blur medan beraksi lebih dulu; aksi kedua kini diantre, bukan hilang.

import { describe, expect, it } from 'vitest'

import { putuskanAksi } from './antreAksi'

describe('putuskanAksi', () => {
  const panel = 'adjdtl:ClaimData.ObjectList(1).ObjectItemList(1).Adjustment(1)'
  const medan = { aksi: 'SetGrossAdjustment', indeks: 0, konteks: panel }
  const tombol = { aksi: 'SendPICProtect', indeks: 0, konteks: panel }

  it('tanpa aksi berjalan: langsung jalan', () => {
    expect(putuskanAksi(null, tombol)).toBe('jalan')
  })
  it('klik tombol saat aksi medan berjalan: diantre', () => {
    expect(putuskanAksi(medan, tombol)).toBe('antre')
  })
  it('klik ganda aksi yang sama pada baris dan konteks yang sama: diabaikan', () => {
    expect(putuskanAksi(tombol, { ...tombol })).toBe('abaikan')
  })
  it('aksi sama di baris lain: diantre', () => {
    const hapus = { aksi: 'HapusItem', indeks: 1, konteks: 'est:ClaimData.ObjectList(1)' }
    expect(putuskanAksi(hapus, { ...hapus, indeks: 2 })).toBe('antre')
  })
  it('aksi dan baris sama di panel objek lain: diantre, bukan klik ganda', () => {
    const hapus = { aksi: 'HapusItem', indeks: 1, konteks: 'est:ClaimData.ObjectList(1)' }
    expect(putuskanAksi(hapus, { ...hapus, konteks: 'est:ClaimData.ObjectList(2)' })).toBe('antre')
  })
})
