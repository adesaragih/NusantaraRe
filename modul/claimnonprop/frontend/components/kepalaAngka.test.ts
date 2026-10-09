// Disalin dari `modul/claimprop/frontend/components/kepalaAngka.test.ts` (pola, bukan impor): keputusan work owner yang disebut di bawah
// adalah keputusan layar Claim Prop yang ditiru Claim Non Prop (prompt Claim Non Prop tahap 1).
// Judul kolom angka sejajar dengan nilainya (work owner 08-10-2026: "ga sejajar header sama nilai" - Claim Amount, Loss
// Allocation, Estimation List, Total Original Currency Estimation, tab Interest dan Spreading): nilai angka rata kanan
// (`claimnonprop__angka`), maka judul kolom berkendali angka ikut rata kanan. Satu komponen Grid untuk semua tabel.

import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Halaman, Tata } from '../api'
import TataView from './TataView'

const grid: Tata = {
  jenis: 'grid',
  jalur: 'ClaimData.ListClaimAmount',
  bernomor: true,
  kolom: [
    { jenis: 'medan', jalur: 'Currency', label: 'Currency', kendali: 'tampil' },
    { jenis: 'medan', jalur: 'Value', label: 'Claim Amount Ceding', kendali: 'angka' },
  ],
  baris: [
    [
      { tampil: true, hanyaBaca: true },
      { tampil: true, hanyaBaca: true },
    ],
  ],
}
const h: Halaman = { nilai: {}, daftar: { 'ClaimData.ListClaimAmount': [{ Currency: 'IDR', Value: '35000000000' }] } }

describe('judul kolom grid', () => {
  const html = renderToStaticMarkup(
    createElement(TataView, {
      tata: [grid],
      k: { h, ubah: () => {}, aksi: () => {}, opsi: () => [], pesanMedan: {}, sibuk: false },
    }),
  )

  it('kolom angka: judul rata kanan seperti nilainya', () => {
    expect(html).toContain('<th><span class="claimnonprop__angka">Claim Amount Ceding</span></th>')
    expect(html).toContain('35.000.000.000')
  })

  it('kolom bukan angka: judul tetap biasa', () => {
    expect(html).toContain('<th>Currency</th>')
  })
})
