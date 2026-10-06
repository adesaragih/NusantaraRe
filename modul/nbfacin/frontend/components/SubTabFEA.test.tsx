// Paritas sub-tab FEA (tiket 41) - grid Pega `FEAList`, form padanan `InputFEA_IsUW`. Data uji sintetis.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { GRID_FEA, TEKS_INWARD } from '../labels'
import SubTabFEA, { adaGalatFEA, feaBaru, unitSah } from './SubTabFEA'

describe('SubTabFEA', () => {
  it('grid kosong: kolom berurutan + Tambah, "No items"', () => {
    const html = renderToStaticMarkup(<SubTabFEA fea={[]} ubah={() => {}} />)
    const kolom = [...html.matchAll(/<th[^>]*>([\s\S]*?)<\/th>|<th[^>]*\/>/g)].map((m) => (m[1] ?? '').replace(/<[^>]+>/g, ''))
    expect(kolom).toEqual(['', ...GRID_FEA.map((k) => k.label.replace('&', '&amp;')), 'Tambah'])
    expect(html).toContain(TEKS_INWARD.kosong)
  })

  it('baris: lima jumlah unit lalu Others Info', () => {
    const html = renderToStaticMarkup(
      <SubTabFEA fea={[{ ...feaBaru(), apar: '3', sprinkler: '0', smokeDetector: '12', hydrant: '1', privateTruckBrigade: '2', info: 'UJI' }]} ubah={() => {}} />,
    )
    expect(html).toContain('<td class="nbf-angka">3</td><td class="nbf-angka">0</td><td class="nbf-angka">12</td><td class="nbf-angka">1</td><td class="nbf-angka">2</td><td>UJI</td>')
  })

  it('jumlah unit: kosong / bilangan bulat >= 0 sah; minus / desimal / huruf ditolak', () => {
    expect(['', '0', '25'].every(unitSah)).toBe(true)
    expect(['-1', '1.5', 'a'].some(unitSah)).toBe(false)
    expect(adaGalatFEA([{ ...feaBaru(), hydrant: '-2' }])).toBe(true)
    expect(adaGalatFEA([feaBaru()])).toBe(false)
  })
})
