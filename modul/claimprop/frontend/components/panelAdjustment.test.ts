// Panel AdjustmentDetail (work owner 08-10-2026 "perbaiki tampilan" + "tidak harus klik angka sebelah kiri untuk membuka
// adjustment nya"): tabel bebas berkolom XML tampil sebagai tabel, dan baris Acceptation List terbaru terbuka tanpa klik.

import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Halaman, Tata } from '../api'
import { DAFTAR_RINCI } from './rincian'
import TataView, { type KonteksTata } from './TataView'

const dasar: Omit<KonteksTata, 'h'> = { ubah: () => {}, aksi: () => {}, opsi: () => [], pesanMedan: {}, sibuk: false }

describe('LetakTabel', () => {
  const j = (p: string) => `ClaimData.AdjustmentList(1).${p}`
  const tabel: Tata = {
    jenis: 'bagian',
    letak: 'tabel',
    anak: [
      {
        jenis: 'bagian',
        anak: [
          { jenis: 'label', label: 'Currency' },
          { jenis: 'label', label: 'Treaty Gross (100%)' },
        ],
      },
      {
        jenis: 'bagian',
        anak: [
          { jenis: 'medan', jalur: j('Currency'), kendali: 'teks', hanyaBaca: true },
          { jenis: 'medan', jalur: j('GrossAdjustment'), kendali: 'angka', hanyaBaca: true },
        ],
      },
      {
        jenis: 'bagian',
        anak: [
          { jenis: 'label', label: 'Total' },
          { jenis: 'medan', jalur: j('AdjustmentValue'), kendali: 'angka', hanyaBaca: true },
        ],
      },
    ],
  }
  const h: Halaman = {
    nilai: {},
    daftar: { 'ClaimData.AdjustmentList': [{ Currency: 'UJA', GrossAdjustment: '1500', AdjustmentValue: '375' }] },
  }
  const html = renderToStaticMarkup(createElement(TataView, { tata: [tabel], k: { ...dasar, h } }))

  it('baris pertama = judul kolom, judul kolom angka rata kanan', () => {
    expect(html).toContain('<table class="claimprop__tabel claimprop__tabel--tetap">')
    expect(html).toContain('<th>Currency</th>')
    expect(html).toContain('<th><span class="claimprop__angka">Treaty Gross (100%)</span></th>')
  })

  it('sel medan menampilkan nilai baris adjustment tanpa label, sel label apa adanya', () => {
    expect(html).toContain('UJA')
    expect(html).toContain('1.500')
    expect(html).toContain('<td>Total</td>')
    expect(html).toContain('375')
  })
})

describe('Acceptation List: panel baris terbaru terbuka tanpa klik', () => {
  const grid: Tata = {
    jenis: 'grid',
    jalur: DAFTAR_RINCI,
    bernomor: true,
    kolom: [{ jenis: 'medan', jalur: 'Type', label: 'Type', kendali: 'tampil' }],
    baris: [[{ tampil: true, hanyaBaca: true }], [{ tampil: true, hanyaBaca: true }]],
  }
  const h: Halaman = { nilai: {}, daftar: { [DAFTAR_RINCI]: [{ Type: 'UJI-1' }, { Type: 'UJI-2' }] } }
  const html = renderToStaticMarkup(
    createElement(TataView, {
      tata: [grid],
      k: { ...dasar, h, rincian: { daftar: DAFTAR_RINCI, isi: (n: number) => `PANEL-${n}` } },
    }),
  )

  it('hanya baris terbaru yang panelnya terbuka', () => {
    expect(html).toContain('PANEL-2')
    expect(html).not.toContain('PANEL-1')
    expect(html).toContain('aria-expanded="true"')
    expect(html).toContain('aria-expanded="false"')
  })
})
