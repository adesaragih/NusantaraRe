// Mode baris grid (`pyGridProps/pyRowEditing`) — laporan pemakai 8 Oktober
// 2026: sel `Kind of Treaty` tab Achievement In IDR panel New tampil sebagai
// dropdown, padahal di Pega tidak dapat disunting.
//
// Bukti ekspor (`Treaty In Adjustment/Section`): grid `TreatyIn.Limits` tab
// Achievement (L69) dan tab Limits (L28) ber-`pyRowEditing = masterDetail`,
// `pyEditingMode = expandPane` — baris TAMPIL saja, penyuntingan di panel
// rincian. Grid Reporting Period (L17) ber-`row`/`readWrite` — disunting di
// tempat. Pembangkit kini membaca mode itu (`alat/ekstrak_kerangka.py`).

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { SisiPenyesuaian } from './api'
import type { ButirKerangka, GridKerangka, MedanKerangka } from './ekspor/jenis'
import { KERANGKA_RINCIAN, KERANGKA_TAB } from './ekspor/kerangka.gen'
import { RenderKerangka } from './komponen/KerangkaTab'
import { terkunci } from './komponen/aksiTombol'

function semuaGrid(isi: readonly ButirKerangka[]): GridKerangka[] {
  const out: GridKerangka[] = []
  const jalan = (x: unknown): void => {
    if (Array.isArray(x)) {
      for (const e of x) jalan(e)
      return
    }
    if (x === null || typeof x !== 'object') return
    const o = x as { t?: string; anak?: unknown; isi?: unknown }
    if (o.t === 'grid') out.push(o as GridKerangka)
    jalan(o.anak)
    jalan(o.isi)
  }
  jalan(isi)
  return out
}

function grid(tab: string, prop: string): GridKerangka {
  const g = semuaGrid(KERANGKA_TAB[tab]?.isi ?? []).find((x) => x.prop === prop)
  if (g === undefined) throw new Error(`${tab} ${prop} tidak ada`)
  return g
}

const sisi: SisiPenyesuaian = {
  medan: { ProportionType: 'Proportional' },
  larik: { Limits: [{ TreatyType: 'QUOTA SHARE' }], ReportingPeriodList: [{ Period: 'Q 1' }] },
}
const kEdit = { sisi, akar: sisi, halaman: { ViewState: '0' }, ubah: true, ubahLarik: () => undefined }

describe('mode baris grid dibaca dari ekspor', () => {
  it('Achievement In IDR dan Limits: masterDetail → Kind of Treaty TIDAK dapat disunting', () => {
    for (const tab of ['TreatyInTabsProportional#Achievement In IDR', 'TreatyInTabsProportional#Limits']) {
      const g = grid(tab, 'TreatyIn.Limits')
      expect(g.modeBaris).toBe('masterDetail')
      expect(g.baca[g.kunci.indexOf('TreatyType')]).toBe('selalu')
    }
  })

  it('Reporting Period: row → sel tetap disunting di tempat', () => {
    const g = grid('TreatyInTabsProportional#Reporting Period', 'TreatyIn.ReportingPeriodList')
    expect(g.modeBaris).toBe('row')
    expect(g.baca.some((b) => b !== 'selalu')).toBe(true)
  })

  it('Kind of Treaty baris baru tetap dapat dipilih — di panel rincian LimitProportional', () => {
    const medan = (KERANGKA_RINCIAN.LimitProportional ?? []).filter((b): b is MedanKerangka => b.t === 'medan')
    const tipe = medan.find((m) => m.kunci === 'TreatyTypeID')
    expect(tipe).toBeDefined()
    // Dikunci HANYA di mode lihat (`TreatyIn.ViewState = 1`); mode Edit terbuka.
    expect(tipe?.baca).not.toBe('selalu')
    expect(terkunci(tipe?.baca, { ViewState: '0' })).toBe(false)
    expect(terkunci(tipe?.baca, { ViewState: '1' })).toBe(true)
  })

  it('mode Edit panel New: sel Achievement tampil sebagai TEKS, bukan dropdown', () => {
    const html = renderToStaticMarkup(
      <RenderKerangka isi={KERANGKA_TAB['TreatyInTabsProportional#Achievement In IDR']?.isi ?? []} k={kEdit} />,
    )
    expect(html).toContain('QUOTA SHARE')
    expect(html).not.toContain('<select')
    expect(html).not.toMatch(/<input[^>]*value="QUOTA SHARE"/)
  })

  it('mode Edit panel New: Reporting Period tetap berisi kotak isian', () => {
    const html = renderToStaticMarkup(
      <RenderKerangka isi={KERANGKA_TAB['TreatyInTabsProportional#Reporting Period']?.isi ?? []} k={kEdit} />,
    )
    expect(html).toMatch(/<input/)
  })
})
