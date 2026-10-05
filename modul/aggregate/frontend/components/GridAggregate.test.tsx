// Grid TempCSV DIRENDER (react-dom/server): hanya dibaca - nol kotak isian (permintaan work owner 04-10-2026
// "upload csv nya read only") - dan angka tampil format Indonesia dua desimal, nilai penuh tetap di title.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import GridAggregate from './GridAggregate'

const HTML = renderToStaticMarkup(
  <GridAggregate
    baris={[
      { ASSESMENT_ZONE: '1.1', TREATY_TYPE: 'QS', CURRENCY: 'IDR', TO_USD: '0.00006060606060606061', RNM_VALUE: '35747721.9', REMARK: 'UJI' },
      { ASSESMENT_ZONE: 'Total :', CURRENCY: 'IDR', RNM_VALUE: '35747721.9' },
    ]}
  />,
)

describe('grid Aggregate', () => {
  it('hanya dibaca: tidak ada kotak isian, baris data maupun baris Total', () => {
    expect(HTML).not.toContain('<input')
    expect(HTML).toContain('>QS</td>')
    expect(HTML).toContain('>UJI</td>')
  })

  it('angka format Indonesia 2 desimal; nilai penuh di title', () => {
    expect(HTML).toContain('title="35747721.9">35.747.721,90</td>')
    expect(HTML).toContain('title="0.00006060606060606061">0,00</td>')
  })
})
