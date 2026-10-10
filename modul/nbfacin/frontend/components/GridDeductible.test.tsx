// Paritas grid / form Deductible (tiket 45, C3) - section Pega `addDeductible`. Data uji sintetis.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { FORM_DEDUCTIBLE as F, KOLOM_TIME_EXCESS, TEKS_INWARD } from '../labels'
import GridDeductible, { adaGalatDeductible, deductibleBaru } from './GridDeductible'

const grid = (daftar = [deductibleBaru('USD')]) =>
  renderToStaticMarkup(<GridDeductible daftar={daftar} currency="USD" mataUang={['USD', 'IDR']} ubah={() => {}} />)

describe('GridDeductible', () => {
  it('kosong = No items; deductible baru bermata uang item (R-2)', () => {
    expect(grid([])).toContain(TEKS_INWARD.kosong)
    expect(deductibleBaru('IDR').currency).toBe('IDR')
  })

  it('baris: label pilihan ditampilkan (bukan kode); Condition Other = teks isian; Amount gaya Indonesia', () => {
    const html = grid([
      { ...deductibleBaru('USD'), typeDeductible: '1', pctDeductible: '10', minMax: '1', amount: '2500000', condition: '5', inputCondition: 'UJI KONDISI', timeExcess: '7' },
    ])
    expect(html).toContain(
      '<td>% Of Claim</td><td>Min</td><td class="nbf-angka">—</td><td>—</td><td>USD</td><td class="nbf-angka">2.500.000,0000</td><td>UJI KONDISI</td>',
    )
  })

  it('MinMax = 3 (Or): kolom Pct / Type Deductible 2 terisi di baris grid (CoverageItem.xml); MinMax lain "—"', () => {
    const atas = { ...deductibleBaru('USD'), typeDeductible: '3', pctDeductible: '10', pctDeductible2: '12', typeDeductible2: '1' }
    expect(grid([{ ...atas, minMax: '3' }])).toContain('<td>Or</td><td class="nbf-angka">12</td><td>% Of Claim</td><td>USD</td>')
    expect(grid([{ ...atas, minMax: '1' }])).toContain('<td>Min</td><td class="nbf-angka">—</td><td>—</td><td>USD</td>')
  })

  it('tata letak rapi (work owner 08-10-2026): Tambah di kepala blok + jumlah; sel kosong "—"; Hapus bergaya bahaya', () => {
    const html = grid([deductibleBaru('USD'), deductibleBaru('USD')])
    expect(html).toMatch(/<div class="nbf-dedgrid__kepala"><h5 class="nbf-objek__judul">Deductible<span class="nbf-dedgrid__jumlah">2<\/span><\/h5><button type="button" class="btn nbf-tema__tombol nbf-tema__tombol--utama btn--sm nbf-dedgrid__tambah">/)
    expect(html).toContain('<td class="nbf-angka">—</td><td>—</td><td>—</td><td class="nbf-angka">—</td><td>—</td><td>USD</td>')
    expect(html).toContain('class="btn btn--ghost btn--sm nbf-dedgrid__hapus"')
  })

  it('kepala kolom kosong seperti CoverageItem.xml, hanya Time Excess (Days) (work owner 08-10-2026)', () => {
    const kepala = grid().match(/<thead>[\s\S]*<\/thead>/)?.[0] ?? ''
    expect(kepala.match(/<th[^>]*>[^<]+<\/th>/g)).toEqual([`<th scope="col" class="nbf-angka">${KOLOM_TIME_EXCESS}</th>`])
  })

  it('angka tak sah ditolak; kosong / desimal sah', () => {
    expect(adaGalatDeductible([{ ...deductibleBaru('USD'), pctDeductible: '1,5' }])).toBe(true)
    expect(adaGalatDeductible([{ ...deductibleBaru('USD'), pctDeductible: '1.5', amount: '1000' }])).toBe(false)
    expect(adaGalatDeductible(undefined)).toBe(false)
  })

  it('label form mengikuti addDeductible', () => {
    expect([F.typeDeductible.label, F.minMax.label, F.timeExcess.label]).toEqual(['Type Deductible', 'MinMax', 'Time Excess (In Days)'])
  })
})
