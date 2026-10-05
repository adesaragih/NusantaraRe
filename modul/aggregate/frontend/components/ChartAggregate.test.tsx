// Chart Aggregate DIRENDER (react-dom/server): batang bertumpuk tanpa pie, tingkat Ceding dan Treaty Type dapat
// dibuka (tombol), tingkat Coverage tidak; total format Indonesia; jejak tingkat dapat diklik kembali.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { Ringkasan } from '../api'
import { PanelChart } from './ChartAggregate'

const DATA: Ringkasan = {
  irisan: [
    { cedingCode: 'UJI-C1', cedingName: 'UJI CEDING SATU', treatyType: 'QS', coverage: 'EQVET', rnmValueInUsd: '100.10' },
    { cedingCode: 'UJI-C1', cedingName: 'UJI CEDING SATU', treatyType: 'OR', coverage: 'FLOOD', rnmValueInUsd: '49.85' },
    { cedingCode: 'UJI-C2', cedingName: 'UJI CEDING DUA', treatyType: 'SURPLUS', coverage: 'EQVET', rnmValueInUsd: '1300' },
  ],
}

const render = (jalur: { ceding?: string; treatyType?: string }) =>
  renderToStaticMarkup(<PanelChart data={DATA} galat={null} jalur={jalur} onJalur={() => {}} />)

describe('chart Aggregate', () => {
  it('tingkat Ceding: batang bertumpuk dapat dibuka, tanpa SVG pie, total format Indonesia', () => {
    const html = render({})
    expect(html).not.toContain('<svg')
    expect(html).not.toContain('<select')
    expect(html).not.toContain('As At')
    expect(html).toContain('>1.449,95</strong>')
    expect(html).toContain('<button type="button" class="aggregate__batang aggregate__batang--buka">')
    expect(html.indexOf('UJI CEDING DUA')).toBeLessThan(html.indexOf('UJI CEDING SATU'))
    expect(html).toContain('title="QS: 100,10"')
    expect(html).toContain('<button type="button" class="aggregate__remah-tombol" disabled="">All Cedings</button>')
  })

  it('tingkat Coverage: batang tidak dapat dibuka; jejak menampilkan ceding dan treaty type', () => {
    const html = render({ ceding: 'UJI-C1', treatyType: 'QS' })
    expect(html).not.toContain('aggregate__batang--buka')
    expect(html).toContain('<div class="aggregate__batang">')
    expect(html).toContain('<button type="button" class="aggregate__remah-tombol">All Cedings</button>')
    expect(html).toContain('<button type="button" class="aggregate__remah-tombol">UJI CEDING SATU</button>')
    expect(html).toContain('>100,10</strong>')
  })
})
