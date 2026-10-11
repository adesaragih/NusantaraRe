// Isian angka tampil 4 desimal (work owner 06-10-2026) - data uji sintetis.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import IsianAngka4, { dariKetikan, untukDiketik } from './IsianAngka4'

describe('IsianAngka4', () => {
  it('tidak difokus: tampil 4 desimal gaya Indonesia; nilai asli tidak diubah', () => {
    const html = renderToStaticMarkup(<IsianAngka4 value="100.00000000000000000000" onChange={() => {}} ariaLabel="% Installment" />)
    expect(html).toContain('value="100,0000"')
  })
  it('difokus: 4 desimal bertitik untuk diketik; koma ketikan dibaca titik', () => {
    expect(untukDiketik('33.33333333333333333333')).toBe('33.3333')
    expect(untukDiketik('164342.50000000')).toBe('164342.5000')
    expect(untukDiketik('')).toBe('')
    expect(dariKetikan('10,5')).toBe('10.5')
    expect(dariKetikan(' 25 ')).toBe('25')
  })
})
