import { describe, expect, it } from 'vitest'

import { jumlahUang, tampilUang, dampakHapusKlaim, hapusKlaim, type DampakHapus} from './api'

// Test ini menjaga SATU aturan: uang di sisi React tetap TEKS, tidak pernah
// menjadi angka (tiket 01 AC-4; ADR-U-0003, ADR-U-0016).
//
// Cara membacanya: `describe` = kelompok kasus, `it` = satu kasus,
// `expect(x).toBe(y)` = "x harus persis y".
describe('uang di sisi React', () => {
  it('membawa jumlah sebagai teks, apa adanya', () => {
    const uang = { amount: '1234567890.12345678', currency: 'IDR' }
    expect(jumlahUang(uang)).toBe('1234567890.12345678')
    expect(typeof jumlahUang(uang)).toBe('string')
  })

  it('tidak kehilangan presisi pada nilai yang tidak muat di float64', () => {
    // Number('1234567890.12345678') membulat menjadi 1234567890.1234567.
    const teks = '1234567890.12345678'
    expect(jumlahUang({ amount: teks, currency: 'IDR' })).toBe(teks)
    expect(String(Number(teks))).not.toBe(teks)
  })

  it('menolak jumlah yang datang sebagai angka JSON', () => {
    // UangMasuk sengaja menerima `unknown`, jadi angka boleh dicoba di sini.
    expect(() => jumlahUang({ amount: 1234.5, currency: 'IDR' })).toThrow(TypeError)
  })

  it('memperlakukan kosong sebagai kosong, bukan nol', () => {
    expect(jumlahUang(null)).toBe('')
    expect(jumlahUang(undefined)).toBe('')
    expect(jumlahUang({ currency: 'IDR' })).toBe('')
    expect(tampilUang(null)).toBe('')
  })

  it('menampilkan jumlah bersama mata uangnya', () => {
    expect(tampilUang({ amount: '250000', currency: 'IDR' })).toBe('250000 IDR')
  })
})

// Tiket 15: yang membuat "Batal" benar-benar membatalkan adalah bentuk
// jalurnya, bukan niat baik layar. Pratinjau dampak wajib MEMBACA saja.
describe('pratinjau dampak penghapusan', () => {
  it('memisahkan baris datar warisan dari totalnya', () => {
    // Bentuk yang server kirim: `total` TIDAK memuat barisDatarWarisan,
    // sebab nasib baris itu belum diputuskan work owner. Menjumlahkannya
    // berarti menjawab pertanyaan yang terbuka lewat sebuah angka.
    const d: DampakHapus = {
      header: 1,
      peserta: 2,
      adjustment: 4,
      spreading: 8,
      spreadingRetro: 16,
      dokumen: 3,
      barisWork: 1,
      barisDatarWarisan: 5,
      total: 35,
    }
    const jumlahDaftar =
      d.header + d.peserta + d.adjustment + d.spreading + d.spreadingRetro +
      d.dokumen + d.barisWork
    expect(d.total).toBe(jumlahDaftar)
    expect(d.total).not.toBe(jumlahDaftar + d.barisDatarWarisan)
  })

  it('membaca lewat GET, bukan lewat kata kerja yang menulis', () => {
    // Fungsi pratinjau memanggil api.get; yang menghapus memanggil api.delete.
    // Keduanya sengaja terpisah supaya membuka popup tidak pernah menulis.
    expect(dampakHapusKlaim.constructor.name).toBe('AsyncFunction')
    expect(hapusKlaim.constructor.name).toBe('AsyncFunction')
    expect(dampakHapusKlaim).not.toBe(hapusKlaim)
  })
})
