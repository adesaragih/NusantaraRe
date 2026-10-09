// Laporan work owner 09-10-2026: sesudah Occupation diisi, klik "Send to Committe" tidak membuka popup - harus klik dua
// kali. Klik itu memicu blur medan beraksi (MakeLowercase) lebih dulu; aksi kedua kini diantre, bukan hilang.

import { describe, expect, it } from 'vitest'

import { putuskanAksi } from './antreAksi'

describe('putuskanAksi', () => {
  const medan = { aksi: 'MakeLowercase', indeks: 0 }
  const tombol = { aksi: 'BukaKomite', indeks: 1 }

  it('tanpa aksi berjalan: langsung jalan', () => {
    expect(putuskanAksi(null, tombol)).toBe('jalan')
  })
  it('klik tombol saat aksi medan berjalan: diantre', () => {
    expect(putuskanAksi(medan, tombol)).toBe('antre')
  })
  it('klik ganda aksi yang sama pada baris yang sama: diabaikan', () => {
    expect(putuskanAksi(tombol, { aksi: 'BukaKomite', indeks: 1 })).toBe('abaikan')
  })
  it('aksi sama di baris lain: diantre', () => {
    expect(putuskanAksi(tombol, { aksi: 'BukaKomite', indeks: 2 })).toBe('antre')
  })
})
