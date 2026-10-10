// Disalin dari `modul/claimnonprop/frontend/nilai.test.ts` (pola, bukan impor; asal Claim Prop): keputusan work
// owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
import { describe, expect, it } from 'vitest'

import type { Halaman, Tata } from './api'
import {
  ambil,
  dariInputWaktu,
  hanyaAngka,
  jalurTerbuka,
  keInputWaktu,
  kunciOpsi,
  masukan,
  modeKirim,
  pecahJalur,
  semuaTata,
  setel,
  sumberServer,
  tampilAngka,
} from './nilai'

const ITEM = 'ClaimData.ObjectList(1).ObjectItemList'

const h: Halaman = {
  nilai: { 'ClaimData.ReporterName': 'UJI-PELAPOR', 'ClaimData.NoClaim': 'UJI-K1' },
  daftar: {
    'ClaimData.ObjectList': [{ ObjectName: 'UJI-OBJEK' }],
    [ITEM]: [{ TSIPerObject: '1000' }, { TSIPerObject: '5', Amount: '2' }],
  },
}

const tata: Tata[] = [
  { jenis: 'medan', jalur: 'ClaimData.NoClaim', kendali: 'teks', hanyaBaca: true },
  { jenis: 'bagian', anak: [{ jenis: 'medan', jalur: 'ClaimData.ReporterName', kendali: 'teks' }] },
  {
    jenis: 'grid',
    jalur: ITEM,
    kolom: [{ jenis: 'medan', jalur: 'TSIPerObject', kendali: 'angka' }],
    baris: [[{ tampil: true }], [{ tampil: true, hanyaBaca: true }]],
  },
]

describe('nilai layar Claim Fac In', () => {
  it('jalur baris dipecah pada indeks terakhir - juga daftar bersarang', () => {
    expect(pecahJalur('ClaimData.ObjectList(1).ObjectItemList(2).Adjustment(3).DataCommitteFacin.Remarks')).toEqual({
      daftar: 'ClaimData.ObjectList(1).ObjectItemList(2).Adjustment',
      n: 3,
      prop: 'DataCommitteFacin.Remarks',
    })
    expect(pecahJalur(`${ITEM}(2).Amount`)).toEqual({ daftar: ITEM, n: 2, prop: 'Amount' })
    expect(pecahJalur('ClaimData.NoClaim')).toBeNull()
  })

  it('ambil / setel jalur halaman dan sel baris bersarang tanpa mengubah aslinya', () => {
    const b = setel(h, `${ITEM}(1).TSIPerObject`, '2000')
    expect(ambil(b, `${ITEM}(1).TSIPerObject`)).toBe('2000')
    expect(ambil(h, `${ITEM}(1).TSIPerObject`)).toBe('1000')
    expect(ambil(h, `${ITEM}(2).Amount`)).toBe('2')
  })

  it('hanya jalur terbuka yang dikirim: medan read-only dan sel baris terkunci tidak ikut', () => {
    expect(jalurTerbuka(tata)).toEqual(['ClaimData.ReporterName', `${ITEM}(1).TSIPerObject`])
    expect(masukan(h, tata)).toEqual({
      'ClaimData.ReporterName': 'UJI-PELAPOR',
      [`${ITEM}(1).TSIPerObject`]: '1000',
    })
  })

  it('masukan mencakup medan panel baris dan modal', () => {
    const panel: Tata[] = [{ jenis: 'medan', jalur: `${ITEM}(2).Amount`, kendali: 'angka' }]
    const modal: Tata[] = [{ jenis: 'medan', jalur: 'SearchName', kendali: 'teks' }]
    const semua = semuaTata(tata, [{ [`estitem:${ITEM}(2)`]: panel }, { pilihPolis: modal }])
    expect(Object.keys(masukan(h, semua))).toEqual([
      'ClaimData.ReporterName',
      `${ITEM}(1).TSIPerObject`,
      `${ITEM}(2).Amount`,
      'SearchName',
    ])
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

describe('mode layar dikembalikan setiap aksi', () => {
  it('Layar.mode terakhir ditimpa nilai lokal penanda mode; jalur lain tidak ikut', () => {
    const lokal: Halaman = {
      nilai: { SearchType: '3', SearchName: 'UJI-TERTANGGUNG', 'ClaimData.ReporterName': 'UJI-PELAPOR' },
      daftar: {},
    }
    expect(modeKirim({ 'ClaimData.EditCatastrope': 'true', SearchType: '1' }, lokal)).toEqual({
      'ClaimData.EditCatastrope': 'true',
      SearchType: '3',
      SearchName: 'UJI-TERTANGGUNG',
    })
  })

  it('tanpa Layar.mode dan tanpa penanda lokal: kosong', () => {
    expect(modeKirim(undefined, { nilai: {}, daftar: {} })).toEqual({})
  })
})

describe('sumber pilihan', () => {
  it('kode associated dan acuan statis tidak ditanyakan ke server; sumber lain ya', () => {
    expect(sumberServer('kode:PaymentType')).toBe(false)
    expect(sumberServer('mataUang')).toBe(false)
    expect(sumberServer('jenisReas')).toBe(false)
    expect(sumberServer('')).toBe(false)
    for (const s of ['mataUangAdj', 'rekening', 'okupasi', 'covFire', 'adjuster', 'provinsi']) {
      expect(sumberServer(s), s).toBe(true)
    }
  })

  it('kunci simpanan pilihan membedakan panel dan baris item', () => {
    const a = kunciOpsi('okupasi', 'est:ClaimData.ObjectList(1)', 2)
    expect(a).not.toBe(kunciOpsi('okupasi', 'est:ClaimData.ObjectList(2)', 2))
    expect(a).not.toBe(kunciOpsi('okupasi', 'est:ClaimData.ObjectList(1)', 3))
  })
})

describe('nomor telepon (work owner 08-10-2026: tidak boleh diisi huruf)', () => {
  it('hanya angka yang tersisa; nol di depan dipertahankan', () => {
    expect(hanyaAngka('021-5797 81ab00')).toBe('021579781' + '00')
    expect(hanyaAngka('abc')).toBe('')
    expect(hanyaAngka('0812')).toBe('0812')
  })
})
