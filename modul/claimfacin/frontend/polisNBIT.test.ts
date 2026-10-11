// Tab Loss Record / Correspondence tampilan polis NB (hanya akun IT): total per Class Of Business, tanggal email, nama
// penerima, dan pembersih HTML pesan (salinan nbfacin - hanya keluarannya yang masuk innerHTML).

import { describe, expect, it } from 'vitest'

import {
  bersihkanHtmlNB,
  daftarNamaPenerimaNB,
  tanggalKorespondensiNB,
  totalCobNB,
  type BarisLossRecordNB,
} from './polisNBIT'

const baris = (b: Partial<BarisLossRecordNB>): BarisLossRecordNB => ({
  uwYear: '',
  classOfBusiness: '',
  dateOfLoss: '',
  causeOfLoss: '',
  currency: '',
  lossAmount: '',
  lossAmountIdr: '',
  ...b,
})

describe('Loss Record', () => {
  it('Total Loss per (Class Of Business, Currency) urut kemunculan; IDR kosong dilewati (totalCobLossRecord)', () => {
    expect(
      totalCobNB([
        baris({ classOfBusiness: 'UJI FIRE', currency: 'IDR', lossAmount: '100.5', lossAmountIdr: '100.5' }),
        baris({ classOfBusiness: 'UJI MB', currency: 'USD', lossAmount: '10' }),
        baris({ classOfBusiness: 'UJI FIRE', currency: 'IDR', lossAmount: '0.25', lossAmountIdr: '0.25' }),
        baris({ classOfBusiness: 'UJI FIRE', currency: 'USD', lossAmount: '1', lossAmountIdr: '16000' }),
      ]),
    ).toEqual([
      { classOfBusiness: 'UJI FIRE', currency: 'IDR', lossAmount: '100.75', lossAmountIdr: '100.75' },
      { classOfBusiness: 'UJI MB', currency: 'USD', lossAmount: '10', lossAmountIdr: '' },
      { classOfBusiness: 'UJI FIRE', currency: 'USD', lossAmount: '1', lossAmountIdr: '16000' },
    ])
    expect(totalCobNB([])).toEqual([])
  })
})

describe('Correspondence', () => {
  it('tanggal: stempel Pega GMT -> WIB DD/MM/YY HH:MM; kabel DD-MM-YYYY HH:MM dipendekkan; lain apa adanya', () => {
    expect(tanggalKorespondensiNB('20261006T102800.000 GMT')).toBe('06/10/26 17:28')
    expect(tanggalKorespondensiNB('20261231T200000.000 GMT')).toBe('01/01/27 03:00')
    expect(tanggalKorespondensiNB('06-10-2026 17:28')).toBe('06/10/26 17:28')
    expect(tanggalKorespondensiNB('UJI')).toBe('UJI')
  })

  it('kolom Recipient: nama atau alamat, dipisah koma', () => {
    expect(
      daftarNamaPenerimaNB([
        { recipient: 'UJI A', recipientMail: 'a@contoh.invalid' },
        { recipient: ' ', recipientMail: 'b@contoh.invalid' },
        { recipient: '', recipientMail: '' },
      ]),
    ).toBe('UJI A, b@contoh.invalid')
  })

  it('pembersih HTML: tag polos yang diizinkan saja, tanpa atribut; script / style dibuang; teks di-escape', () => {
    expect(bersihkanHtmlNB('<p onclick="x()">UJI <b style="c">tebal</b></p><script>alert(1)</script>')).toBe(
      '<p>UJI <b>tebal</b></p>',
    )
    expect(bersihkanHtmlNB('<img src=x onerror=alert(1)>a < b & c &amp; <a href="j">t</a><br/>')).toBe(
      'a &lt; b &amp; c &amp; t<br>',
    )
  })
})
