// Paritas sub-tab Occupation (tiket 40) - section Pega `OccupationList` / `OccupationItemFacIn_Section` / popup
// `ChooseOccupation` dan `ChooseClassofContraction`. Data uji sintetis.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { GRID_OKUPASI, TEKS_INWARD } from '../labels'
import SubTabOkupasi, { kategoriDari, okupasiBaru, pilihKonstruksi, pilihOkupasi } from './SubTabOkupasi'

const SUMBER = readFileSync(join(__dirname, 'SubTabOkupasi.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('SubTabOkupasi', () => {
  it('grid kosong: kolom berurutan + Tambah, "No items"', () => {
    const html = renderToStaticMarkup(<SubTabOkupasi caseId="NB-1" occupations={[]} ubah={() => {}} />)
    const kolom = [...html.matchAll(/<th[^>]*>([\s\S]*?)<\/th>|<th[^>]*\/>/g)].map((m) => (m[1] ?? '').replace(/<[^>]+>/g, ''))
    expect(kolom).toEqual(['', ...GRID_OKUPASI.map((k) => k.label), 'Tambah'])
    expect(html).toContain(TEKS_INWARD.kosong)
  })

  it('baris tampil Occupation ID, Name, Class Of Construction', () => {
    const html = renderToStaticMarkup(
      <SubTabOkupasi
        caseId="NB-1"
        occupations={[{ occupationId: 'UJI-1', occupationName: 'UJI NAMA', category: 'I', constructionClass: 'UJI KELAS', pctLimit: '' }]}
        ubah={() => {}}
      />,
    )
    expect(html).toContain('<td>UJI-1</td><td>UJI NAMA</td><td>UJI KELAS</td>')
  })

  it('SetDataOccupation: 03/02/01 -> III/II/I, lainnya kosong; Class of Construction tidak disentuh', () => {
    expect(['03', '02', '01', '04', undefined].map(kategoriDari)).toEqual(['III', 'II', 'I', '', ''])
    const o = pilihOkupasi({ ...okupasiBaru(), constructionClass: 'LAMA' }, { oldId: 'UJI-9', name: 'UJI N', kdRiskExposure: '02' })
    expect(o).toEqual({ occupationId: 'UJI-9', occupationName: 'UJI N', category: 'II', constructionClass: 'LAMA', pctLimit: '' })
  })

  it('SetDataClassofConstraction: Description + PctLimit (teks apa adanya)', () => {
    expect(pilihKonstruksi(okupasiBaru(), { description: 'UJI D', pctLimit: '62,5' })).toMatchObject({ constructionClass: 'UJI D', pctLimit: '62,5' })
  })

  it('Choose Class of Construction butuh Category; popup bergantian (tidak bertumpuk); Occupation dimuat saat dibuka', () => {
    // Ditahan hanya bila Occupation belum dipilih; Category kosong tetap membuka popup (seperti Pega).
    expect(SUMBER).toMatch(/if \(o\.occupationId === ''\) \{\s*setPesan\(TEKS_OKUPASI\.pilihOkupasiDulu\)/)
    expect(SUMBER).not.toMatch(/if \(o\.category === ''\)/)
    expect(SUMBER).toContain("{popup?.jenis === 'okupasi' && (")
    expect(SUMBER).toContain("{popup?.jenis === 'konstruksi' && (")
    expect(SUMBER).toContain("kotak === '' ? 0 : JEDA_MS")
  })
})
