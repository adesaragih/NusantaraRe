// Tampilan rincian `DetailLimits` panel Adjustment — laporan pemakai
// 8 Oktober 2026: isi meluber ke kanan, label "100% Limit100%" menempel, dan
// sebelas blok Event Limits … Achievement bertumpuk ke bawah.
//
// Bukti ekspor: layout group `DetailLimits` ber-`pyHeaderType = TABBED`
// (sebelas tab, urutan sama dengan `modul/treatyin/frontend/labelsLimitsProp.ts`
// `TAB_DETAIL_LIMIT`). Pembangkit menandainya `tab: true`, dan
// `RenderKerangka` menghimpun blok bertanda yang berurutan menjadi SATU strip.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { SisiPenyesuaian } from './api'
import type { BlokKerangka } from './ekspor/jenis'
import { KERANGKA_RINCIAN } from './ekspor/kerangka.gen'
import { RenderKerangka } from './komponen/KerangkaTab'

const TAB_EKSPOR = [
  'Event Limits', 'Deduction In A', 'Deduction', 'Reserve', 'Experience Premium Refund',
  'PLA', 'Cash Loss Limit', 'Claim Cooperation', 'LPC', 'EPI', 'Achievement',
]

const sisi: SisiPenyesuaian = { medan: { ProportionType: 'Proportional' }, larik: {} }

describe('rincian DetailLimits — layout group TABBED', () => {
  it('kerangka membawa sebelas blok tab BERURUTAN, judul dan urutan ekspor', () => {
    const isi = KERANGKA_RINCIAN.DetailLimits ?? []
    const tab = isi.filter((b): b is BlokKerangka => b.t === 'blok' && b.tab === true)
    expect(tab.map((b) => b.judul)).toEqual(TAB_EKSPOR)
    const awal = isi.indexOf(tab[0] as BlokKerangka)
    expect(isi.slice(awal, awal + tab.length)).toEqual(tab)
  })

  it('digambar sebagai SATU strip tab; hanya isi tab pertama yang dirender', () => {
    const html = renderToStaticMarkup(
      <RenderKerangka isi={KERANGKA_RINCIAN.DetailLimits ?? []} k={{ sisi, akar: sisi, halaman: { ViewState: '1' } }} />,
    )
    for (const j of TAB_EKSPOR) expect(html).toContain(j)
    expect(html.match(/tria__grup-tab/g)).toHaveLength(1)
    // Isi tab pertama (Event Limits) tampil; isi tab lain (Premium Reserve) tidak.
    expect(html).toContain('RSMD Limit')
    expect(html).not.toContain('% Premium Reserve')
  })

  it('isi rincian baris dibungkus blok mengalir, dan tabelnya dikunci selebar wadah', () => {
    const tsx = readFileSync(join(__dirname, 'komponen', 'KerangkaTab.tsx'), 'utf8')
    const i = tsx.indexOf('<tr className="tria__rincian">')
    expect(tsx.slice(i, i + 600)).toContain('<div className="tria__blok">')
    const css = readFileSync(join(__dirname, 'treatyinadjustment.css'), 'utf8')
    expect(css).toContain('.treatyinadjustment .tria__rincian .tria__tabel {')
    expect(css).toMatch(/table-layout: fixed;/)
  })
})
