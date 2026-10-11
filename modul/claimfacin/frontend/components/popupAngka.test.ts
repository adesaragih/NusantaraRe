// Tabel pop-up (Outstanding Claim, Retro List): work owner 11-10-2026 "PERBAIKI TAMPILAN INI" - nilai TSI RNM /
// Estimation / Kurs / Convert bertumpuk di kolom TSI RNM. Sel angka ber-kelas `claimfacin__angka` (display: block untuk
// nilai medan) menjadikan <td> sebuah blok, sehingga sel-sel angka sebaris menumpuk di satu kolom. Angka dibungkus
// <span> di dalam <td> seperti grid lain, dan judul kolom angka rata kanan seperti nilainya.

import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { Angka, Tabel } from './Popup'

describe('tabel pop-up', () => {
  const html = renderToStaticMarkup(
    createElement(Tabel, {
      kolom: ['Currency', 'TSI RNM', 'Kurs Value'],
      angka: [1, 2],
      baris: [
        createElement(
          'tr',
          { key: 1 },
          createElement('td', null, 'IDR'),
          createElement(Angka, { v: '50000000000000' }),
          createElement(Angka, { v: '1' }),
        ),
      ],
    }),
  )

  it('sel angka tetap sel tabel (tidak bertumpuk): angka di span di dalam td', () => {
    expect(html).not.toContain('<td class="claimfacin__angka"')
    expect(html).toContain('<td><span class="claimfacin__angka">50.000.000.000.000,0000</span></td>')
    expect(html).toContain('<td><span class="claimfacin__angka">1,0000</span></td>')
  })

  it('judul kolom angka rata kanan, kolom lain biasa', () => {
    expect(html).toContain('<th><span class="claimfacin__angka">TSI RNM</span></th>')
    expect(html).toContain('<th><span class="claimfacin__angka">Kurs Value</span></th>')
    expect(html).toContain('<th>Currency</th>')
  })
})
