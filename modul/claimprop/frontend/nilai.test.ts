import { describe, expect, it } from 'vitest'

import type { Halaman, Tata } from './api'
import {
  ambil,
  dariInputWaktu,
  hanyaAngka,
  jalurTerbuka,
  keInputWaktu,
  masukan,
  pecahJalur,
  setel,
  tampilAngka,
} from './nilai'

const h: Halaman = {
  nilai: { 'ClaimData.ReporterName': 'UJI-PELAPOR', 'ClaimData.NoClaim': 'UJI-K1' },
  daftar: { 'ClaimData.InterestList': [{ TSIPerObject: '1000' }, { TSIPerObject: '5', IsAdjVal: 'Yes' }] },
}

const tata: Tata[] = [
  { jenis: 'medan', jalur: 'ClaimData.NoClaim', kendali: 'teks', hanyaBaca: true },
  { jenis: 'bagian', anak: [{ jenis: 'medan', jalur: 'ClaimData.ReporterName', kendali: 'teks' }] },
  {
    jenis: 'grid',
    jalur: 'ClaimData.InterestList',
    kolom: [{ jenis: 'medan', jalur: 'TSIPerObject', kendali: 'angka' }],
    baris: [[{ tampil: true }], [{ tampil: true, hanyaBaca: true }]],
  },
]

describe('nilai layar Claim Prop', () => {
  it('jalur baris dipecah pada indeks terakhir', () => {
    expect(pecahJalur('ClaimData.AdjustmentList(2).DataCommitteeTreaty.Remarks')).toEqual({
      daftar: 'ClaimData.AdjustmentList',
      n: 2,
      prop: 'DataCommitteeTreaty.Remarks',
    })
    expect(pecahJalur('ClaimData.NoClaim')).toBeNull()
  })

  it('ambil / setel jalur halaman dan sel baris tanpa mengubah aslinya', () => {
    const b = setel(h, 'ClaimData.InterestList(1).TSIPerObject', '2000')
    expect(ambil(b, 'ClaimData.InterestList(1).TSIPerObject')).toBe('2000')
    expect(ambil(h, 'ClaimData.InterestList(1).TSIPerObject')).toBe('1000')
  })

  it('hanya jalur terbuka yang dikirim: medan read-only dan sel baris terkunci tidak ikut', () => {
    expect(jalurTerbuka(tata)).toEqual(['ClaimData.ReporterName', 'ClaimData.InterestList(1).TSIPerObject'])
    expect(masukan(h, tata)).toEqual({
      'ClaimData.ReporterName': 'UJI-PELAPOR',
      'ClaimData.InterestList(1).TSIPerObject': '1000',
    })
  })

  it('tanggal-waktu halaman <-> datetime-local', () => {
    expect(keInputWaktu('2026-05-01 08:30:00')).toBe('2026-05-01T08:30')
    expect(dariInputWaktu('2026-05-01T08:30')).toBe('2026-05-01 08:30:00')
  })

  // work owner 08-10-2026: "separator indonesia", "4 angka belakang koma"
  it('angka bergaya Indonesia: titik ribuan, koma desimal, maks 4 desimal, nol ekor dibuang', () => {
    expect(tampilAngka('1234567.0123456789')).toBe('1.234.567,0123')
    expect(tampilAngka('1234567.01235')).toBe('1.234.567,0124')
    expect(tampilAngka('2484250.000000')).toBe('2.484.250')
    expect(tampilAngka('0.5')).toBe('0,5')
    expect(tampilAngka('-1000')).toBe('-1.000')
    expect(tampilAngka('')).toBe('')
    expect(tampilAngka('abc')).toBe('abc')
  })
})

describe('nomor telepon (work owner 08-10-2026: tidak boleh diisi huruf)', () => {
  it('hanya angka yang tersisa; nol di depan dipertahankan', () => {
    expect(hanyaAngka('021-5797 81ab00')).toBe('021579781' + '00')
    expect(hanyaAngka('abc')).toBe('')
    expect(hanyaAngka('0812')).toBe('0812')
  })
})
