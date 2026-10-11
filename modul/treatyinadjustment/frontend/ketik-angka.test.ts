// Kotak tanggal & kotak angka-berbentuk-teks menolak huruf — permintaan
// pemakai 9 Oktober 2026 ("cek keseluruhan").

import { describe, expect, it } from 'vitest'

import { angkaBerformat, penyaringKetik, saringAngka, saringTanggal } from './komponen/ketikAngka'

describe('ketikan kotak tanggal', () => {
  it('hanya digit dan pemisah / - . spasi, maks 10 aksara', () => {
    expect(saringTanggal('18/01/2024')).toBe('18/01/2024')
    expect(saringTanggal('18a/0b1/2024x')).toBe('18/01/2024')
    expect(saringTanggal('2024-01-18')).toBe('2024-01-18')
    expect(saringTanggal('abc')).toBe('')
    expect(saringTanggal('18/01/20245678')).toBe('18/01/2024')
  })
})

describe('pxTextInput berisi angka', () => {
  it('pecahan: Amount, % RNM Share, Limit, Layer (Limits)', () => {
    for (const k of ['Amount', 'RNMShare', 'Limit', 'AgregateLimit2', 'Layer', 'RIOGR']) {
      expect(penyaringKetik(k)?.('12a,5%')).toBe('12.5')
    }
  })
  it('bulat: WPC, Installment, Reporting, Submission Days', () => {
    for (const k of ['WPC', 'InstallmentNo', 'ReportingSubmission', 'SubDays']) {
      expect(penyaringKetik(k)?.('30 hari')).toBe('30')
    }
  })
  it('teks sungguhan tetap bebas', () => {
    expect(penyaringKetik('Comment')).toBeUndefined()
    expect(penyaringKetik('CoInShare')).toBeUndefined()
    expect(penyaringKetik('Layer', 'ShareReins')).toBeUndefined()
    expect(penyaringKetik('Layer', 'ShareFacultativeReinsurers')).toBeUndefined()
    expect(penyaringKetik('TreatyContractName')).toBeUndefined()
  })
  it('saringAngka sama dengan salinan Treaty In', () => {
    expect(saringAngka('-7,5x')).toBe('-7.5')
    expect(saringAngka('10,5', true)).toBe('10')
  })
})

describe('angka pecahan ber-pxTextInput tampil berformat', () => {
  it('Amount (rincian Maximum Retention), % RNM Share, Limit: berformat; Layer/Part & bulat: tidak', () => {
    expect(angkaBerformat('Amount')).toBe(true)
    expect(angkaBerformat('RNMShare')).toBe(true)
    expect(angkaBerformat('Limit', 'Limits')).toBe(true)
    expect(angkaBerformat('Conversion', 'CurrencyList')).toBe(true)
    expect(angkaBerformat('Layer')).toBe(false)
    expect(angkaBerformat('LayerPart')).toBe(false)
    expect(angkaBerformat('WPC')).toBe(false)
    expect(angkaBerformat('Comment')).toBe(false)
  })
})
