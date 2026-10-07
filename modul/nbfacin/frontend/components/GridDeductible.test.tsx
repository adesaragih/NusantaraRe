// Paritas grid / form Deductible (tiket 45, C3) - section Pega `addDeductible`. Data uji sintetis.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { FORM_DEDUCTIBLE as F, TEKS_INWARD } from '../labels'
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
    expect(html).toContain('<td>% Of Claim</td><td>Min</td><td>USD</td><td class="nbf-angka">2.500.000</td><td>UJI KONDISI</td>')
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
