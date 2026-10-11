// Zone / 4.2 Construction di form coverage (tiket 52) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { COVERAGE_KONSTRUKSI, COVERAGE_ZONA, ZONA_COVERAGE as Z } from '../labels'
import FormCoverage, { coverageBaru } from './FormCoverage'

const KORPUS = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\CoverageItem.xml'
const SUMBER = readFileSync(join(__dirname, 'FormCoverage.tsx'), 'utf8').replace(/\r\n/g, '\n')
const form = (coverage: string, zone = '', constructionEarthquake = '') =>
  renderToStaticMarkup(
    <FormCoverage caseId="NB-1" c={{ ...coverageBaru(), coverage, oldId: 'UJI', zone, constructionEarthquake }} tsiItem="1000000" ubah={() => {}} />,
  )

describe.skipIf(!existsSync(KORPUS))('label dan syarat = korpus (CoverageItem)', () => {
  const xml = existsSync(KORPUS) ? readFileSync(KORPUS, 'utf-8') : ''
  it('label Zone, 4.2 Construction, tooltip ikon', () => {
    expect(xml).toContain(`<pyLabelFieldValue>${Z.zone}</pyLabelFieldValue>`)
    expect(xml).toContain(`<pyLabelFieldValue>${Z.konstruksi}</pyLabelFieldValue>`)
    expect(xml).toContain(`<pyTooltip>${Z.ubahZona}</pyTooltip>`)
  })
  it('kode coverage = syarat tampil sel 6 dan 13', () => {
    for (const k of COVERAGE_ZONA) expect(xml).toContain(`.Coverage=='${k}'`)
    expect(xml).toContain(".Coverage=='100821' || .Coverage=='100829'")
    expect(COVERAGE_KONSTRUKSI.every((k) => COVERAGE_ZONA.includes(k))).toBe(true)
  })
})

describe('FormCoverage - Zone / 4.2 Construction', () => {
  it('4.2 PRGBI (100829): Zone (autocomplete) dan 4.2 Construction tampil', () => {
    const html = form('100829')
    expect(html).toContain(`>${Z.zone}<`)
    expect(html).toContain(`>${Z.konstruksi}<`)
  })
  it('RSMD 4.1A CC (100828): Zone tampil, 4.2 Construction tidak', () => {
    const html = form('100828')
    expect(html).toContain(`>${Z.zone}<`)
    expect(html).not.toContain(`>${Z.konstruksi}<`)
  })
  it('coverage lain (FLEXAS 100815): keduanya tidak tampil', () => {
    const html = form('100815')
    expect(html).not.toContain(`>${Z.zone}<`)
    expect(html).not.toContain(`>${Z.konstruksi}<`)
  })
  it('Zone terisi -> teks baca-saja + ikon "Switch to edit mode"', () => {
    const html = form('100829', 'Daerah Zona 5')
    expect(html).toContain('<div class="nbf-inward__teks">Daerah Zona 5</div>')
    expect(html).toContain(`title="${Z.ubahZona}"`)
  })
  it('Zone dipilih / Coverage diganti -> CallAct "rate"; ganti coverage mengosongkan Zone + Construction (SetShowZone_DT)', () => {
    expect(SUMBER).toContain("onPilih={(s) => ubahHitung({ ...c, zone: s.label }, 'percent', 'rate')}")
    expect(SUMBER).toContain("coverageNote: row.nama, zone: '', constructionEarthquake: '' },")
    expect(SUMBER).toContain("onClick={() => ubah({ ...c, zone: '', constructionEarthquake: '' })}")
  })
})

describe('FormCoverage - koma desimal diterima (Net Rate dkk.)', () => {
  it('isian angka coverage & ‰ Total Net Rate mengubah koma menjadi titik sebelum hitung ulang', () => {
    expect(SUMBER).toContain('(MEDAN_ANGKA as readonly string[]).includes(k) ? dariKetikan(v) : v')
    expect(SUMBER).toContain("ubahHitung({ ...c, [k]: dariKetikan(v) }, 'percent')")
    const tab = readFileSync(join(__dirname, 'TabCoverage.tsx'), 'utf8')
    expect(tab).toContain('onUbah(dariKetikan(e.target.value))')
  })
})

describe('FormCoverage - data objek untuk CountPremiRate_ACT (tiket 52)', () => {
  it('callAct "rate" membawa provinsi dan OccupationId objek', () => {
    expect(SUMBER).toContain("...(callAct ? { callAct, province: lokasi?.province ?? '', occupationIds: lokasi?.occupationIds ?? [] } : {}),")
    const tab = readFileSync(join(__dirname, 'TabCoverage.tsx'), 'utf8')
    expect(tab).toContain('lokasi={{ province: x.province, occupationIds: x.occupations.map((k) => k.occupationId) }}')
  })
})

describe('Zone / 4.2 Construction - daftar muncul saat kotak difokus', () => {
  it('kedua isian memakai muatSaatFokus; IsianSaran memuat daftar saat fokus', () => {
    expect(SUMBER.match(/muatSaatFokus\n/g)).toHaveLength(2)
    const isian = readFileSync(join(__dirname, 'IsianSaran.tsx'), 'utf8')
    expect(isian).toContain('if (!muatSaatFokus) return')
  })
})
