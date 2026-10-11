// Paritas form coverage Aneka (tiket 73 P2) - `NB FacIn\Section\InputCoverageAneka_FacIn.xml` + tangkapan layar Pega work owner
// 09-10-2026 (layout 3 kolom). Data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { TampilanAneka } from '../api'
import { FORM_COVERAGE_ANEKA as F, OPSI_METODE_ANEKA, OPSI_UNIT_INDEMNITY } from '../labels'
import { coverageAnekaBaru } from './coverageAneka'
import FormCoverageAneka from './FormCoverageAneka'

const SUMBER = readFileSync(join(__dirname, 'FormCoverageAneka.tsx'), 'utf8').replace(/\r\n/g, '\n')

const T = (p: Partial<TampilanAneka>): TampilanAneka => ({
  kontainer: ['C1'], kolomTahunJumlah: false, he: false, golf: false, tanamanHasil: false, engineering: false, bonding: false,
  briguna: false, topRisk: false, sectionLol: false, jumlahSatuan: true, currencyCoverage: false, mbd: false, ...p,
})

const render = (tampilan: TampilanAneka, metode = '1') =>
  renderToStaticMarkup(
    <FormCoverageAneka
      caseId="NB-1"
      idDasar="u"
      c={{ ...coverageAnekaBaru({ currency: 'IDR', tsi: '500' }, metode) }}
      ubah={() => {}}
      hitung={() => {}}
      tampilan={tampilan}
      mataUang={['IDR']}
      proRate="50"
      ubahProRate={() => {}}
      ubahDeductible={() => {}}
    />,
  )

/** Label sel berurutan sesuai kemunculan di HTML. */
const urutan = (html: string, label: string[]) => label.map((l) => html.indexOf(`>${l}<`))

describe('FormCoverageAneka', () => {
  it('opsi = properti Pega: Calculate Method 1/2/3, Indemnity Unit 0/1/2', () => {
    expect(OPSI_METODE_ANEKA.map((o) => o.label)).toEqual(['PRO RATA', 'SHORT PERIOD', 'FIX RATE'])
    expect(OPSI_UNIT_INDEMNITY.map((o) => o.label)).toEqual(['Day (s)', 'Month (s)', 'Year (s)'])
  })

  it('nilai awal sel (pyDefaultValue): % Short Period 0, Day 365, % Loading 0, % Scale 100, Indemnity 0, % Indemnity 100, % Discount 0', () => {
    const c = coverageAnekaBaru({ currency: '', tsi: '' }, '')
    expect([c.calculateMethodFacIn, c.pctShortPeriod, c.day, c.loading, c.lostLimit, c.indemnity, c.indemnityPercentage, c.discountPercentage]).toEqual(
      ['1', '0', '365', '0', '100', '0', '100', '0'],
    )
  })

  it('layout Pega: satu kisi, urutan sel Coverage Note → Calculate Method → % Pro Rate → Day → TSI → % Loading → % Scale → % Rate → Premium → % Discount → Discount → LoL → % LoL → Sub Limit → Deductible', () => {
    const h = render(T({}))
    const p = urutan(h, [
      F.coverageNote, F.calculateMethod, F.proRate, F.day, F.tsi, F.loading, F.scale, F.rate, F.premium, F.pctDiscount,
      F.discount, F.lol, F.pctLol, F.subLimit,
    ])
    expect(p.every((x) => x >= 0)).toBe(true)
    expect([...p].sort((a, b) => a - b)).toEqual(p)
    expect(h.indexOf('Deductible')).toBeGreaterThan(p[p.length - 1]!)
    expect(h).toContain('class="nbf-aneka-cov__grid"')
  })

  it('% Pro Rate hanya metode 1 (nilai case), % Short Period hanya metode 2', () => {
    expect(render(T({}), '1')).toContain(F.proRate)
    expect(render(T({}), '1')).not.toContain(F.shortPeriod)
    expect(render(T({}), '2')).toContain(F.shortPeriod)
    expect(render(T({}), '2')).not.toContain(F.proRate)
    expect(render(T({}), '3')).not.toContain(F.proRate)
  })

  it('baca-saja tanpa syarat (Currency, TSI, % Indemnity) = teks, bukan isian', () => {
    const h = render(T({ currencyCoverage: true, mbd: true }))
    expect(h).not.toContain('id="u-tsi"')
    expect(h).not.toContain('id="u-indemnityPercentage"')
    expect(h).toContain('id="u-indemnity"')
    expect(h).toContain(`>${F.currency}<`)
    expect(SUMBER).toContain('<Tampil label={F.tsi} nilai={angka4(c.tsi)} angka />')
    expect(SUMBER).toContain('{tampilan.currencyCoverage && <Tampil label={F.currency} nilai={c.currency} />}')
  })

  it('Currency hanya bila currencyCoverage; Indemnity hanya MBD; Coverage Note tersembunyi untuk Yield Shortfall', () => {
    expect(render(T({}))).not.toContain(`>${F.currency}<`)
    expect(render(T({}))).not.toContain(F.pctIndemnity)
    expect(render(T({ mbd: true }))).toContain(F.pctIndemnity)
    expect(render(T({}))).toContain(F.coverageNote)
    expect(render(T({ tanamanHasil: true }))).not.toContain(F.coverageNote)
  })

  it('pemicu = rantai on-change sel Pega yang dapat diisi (nama kabel c3)', () => {
    for (const s of [
      "pilihLalu('calculateMethodFacIn', 'metode')",
      "onSelesai={() => hitung(c, 'proRate')}",
      "angka('pctShortPeriod', F.shortPeriod, 'shortPeriod')",
      "pilihLalu('day', 'day')",
      "angka('loading', F.loading, 'loading')",
      "angka('lostLimit', F.scale, 'scale')",
      "angka('rate', F.rate, 'rate')",
      "angka('indemnity', F.indemnity, 'indemnity')",
      "pilihLalu('unit', 'indemnity')",
      "angka('premium', F.premium, 'premium')",
      "angka('discountPercentage', F.pctDiscount, 'diskonPersen')",
      "angka('discount', F.discount, 'diskonJumlah')",
      "angka('limitOfLiability', F.lol, 'lol')",
      "angka('pctLol', F.pctLol, 'pctLol')",
    ])
      expect(SUMBER).toContain(s)
  })

  it('Coverage Note = saran sambil mengetik (work owner 09-10-2026 "balikin"): pilih `.ID` → coverage, `.NamaCoverage` → coverageNote; tanpa tombol Choose Coverage', () => {
    expect(SUMBER).toContain('ubah({ ...c, coverage: s.id, coverageNote: s.label, oldId: b?.oldId ?? c.oldId })')
    expect(SUMBER).toContain('<IsianSaran')
    // Ketikan di luar daftar dihapus (work owner 09-10-2026); mengetik melepas kode coverage lama.
    expect(SUMBER).toMatch(/<IsianSaran[\s\S]*?wajibDariDaftar[\s\S]*?onKetik=\{\(v\) => ubah\(\{ \.\.\.c, coverageNote: v, coverage: '', oldId: '' \}\)\}/)
    expect(render(T({}))).not.toContain('Choose Coverage')
  })
})
