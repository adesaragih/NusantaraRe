import { describe, expect, it } from 'vitest'
import type { PilihanWarisan } from './api'
import { saringTerdekat, peringkatKedekatan } from './saring'

const D = [
  { id: 'G0000006', nama: 'ASURANSI ADIRA DINAMIKA', kembar: true },
  { id: 'G0000010', nama: 'ASURANSI DIGITAL BERSAMA (D/H SARANA LINDUNG UPAYA)', kembar: false },
  { id: 'G0000038', nama: 'ASURANSI KRESNA MITRA TBK', kembar: true },
  { id: 'G0000041', nama: 'ASURANSI MEGA PRATAMA', kembar: false },
  { id: 'G0000050', nama: 'ASURANSI SIMAS INSURTECH', kembar: false },
]

describe('saringan terdekat', () => {
  it('"asuransi adira" hanya memunculkan yang itu', () => {
    const h = saringTerdekat(D, 'asuransi adira')
    expect(h.map((x) => x.nama)).toEqual(['ASURANSI ADIRA DINAMIKA'])
  })
  it('"adira" menemukan lewat awal KATA', () => {
    expect(saringTerdekat(D, 'adira').map((x) => x.nama)).toEqual(['ASURANSI ADIRA DINAMIKA'])
  })
  it('"asuransi" memang memunculkan semuanya — ketikannya memang seluas itu', () => {
    expect(saringTerdekat(D, 'asuransi')).toHaveLength(5)
  })
  it('kode pengenal tetap bisa dicari', () => {
    expect(saringTerdekat(D, 'G0000038').map((x) => x.id)).toEqual(['G0000038'])
  })
  it('nol cocok mengembalikan kosong', () => {
    expect(saringTerdekat(D, 'zzz')).toHaveLength(0)
  })
  it('peringkat: sama persis < awalan < awal kata < mengandung', () => {
    expect(peringkatKedekatan('ASURANSI ADIRA DINAMIKA', 'G1', 'asuransi adira dinamika')).toBe(0)
    expect(peringkatKedekatan('ASURANSI ADIRA DINAMIKA', 'G1', 'asuransi')).toBe(1)
    expect(peringkatKedekatan('ASURANSI ADIRA DINAMIKA', 'G1', 'adira')).toBe(2)
    expect(peringkatKedekatan('ASURANSI ADIRA DINAMIKA', 'G1', 'dir')).toBe(3)
    expect(peringkatKedekatan('ASURANSI ADIRA DINAMIKA', 'G1', 'zzz')).toBeNull()
  })
})

// ===========================================================================
// Keluhan pemilik proses 6 Oktober 2026 — tangkapan layar `Ceding`
// ===========================================================================
//
// Diketik `zurich`, yang tampil TUJUH baris dan hanya satu yang cocok:
//
//   ASURANSI PERISAI LISTRIK NASIONAL · ASURANSI SIMAS INSURTECH ·
//   ASURANSI SINAR MAS · ASURANSI SOMPO JAPAN NIPPONKOA INDONESIA ·
//   BESS CENTRAL INSURANCE · BRI ASURANSI INDONESIA · ZURICH ASURANSI INDONESIA
//
// ⭐ Diuji dengan nama-nama PERSIS itu, bukan dengan contoh yang dikarang:
// contoh yang dikarang membuktikan bahwa contoh itu lewat, bukan bahwa
// keluhannya tertutup.
describe('keluhan nyata — mengetik `zurich` pada daftar cedant', () => {
  const daftarNyata: PilihanWarisan[] = [
    'ASURANSI PERISAI LISTRIK NASIONAL',
    'ASURANSI SIMAS INSURTECH',
    'ASURANSI SINAR MAS',
    'ASURANSI SOMPO JAPAN NIPPONKOA INDONESIA',
    'BESS CENTRAL INSURANCE',
    'BRI ASURANSI INDONESIA (D/H BRINS GENERAL INSURANCE)',
    'ZURICH ASURANSI INDONESIA',
  ].map((nama, i) => ({ id: `G000000${String(i)}`, nama, kembar: false }))

  it('⛔ hanya ZURICH yang tampil — bukan tujuh baris', () => {
    expect(saringTerdekat(daftarNyata, 'zurich').map((x) => x.nama)).toEqual([
      'ZURICH ASURANSI INDONESIA',
    ])
  })

  it('huruf besar-kecil tidak membedakan', () => {
    for (const q of ['ZURICH', 'Zurich', '  zUrIcH  ']) {
      expect(saringTerdekat(daftarNyata, q)).toHaveLength(1)
    }
  })

  // ⚠️ Ketikan KOSONG memang menampilkan seluruhnya — itu keadaan "belum
  // mengetik apa-apa", bukan saringan yang gagal. Dibedakan di sini supaya
  // tidak ada yang kelak "memperbaikinya" menjadi daftar kosong.
  it('ketikan kosong menampilkan seluruhnya — itu memang maksudnya', () => {
    expect(saringTerdekat(daftarNyata, '')).toHaveLength(7)
  })

  it('potongan di TENGAH kata tidak menarik baris yang tidak relevan', () => {
    // `sinar` hanya menemukan SINAR MAS, bukan setiap nama berawalan ASURANSI.
    expect(saringTerdekat(daftarNyata, 'sinar').map((x) => x.nama)).toEqual([
      'ASURANSI SINAR MAS',
    ])
  })
})
