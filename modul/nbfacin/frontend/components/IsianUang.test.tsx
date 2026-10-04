// IsianUang - format Indonesia saat diketik, kabel desimal bertitik, tanpa float.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import IsianUang, { keKabelUang, keTampilUang, rapikanUang } from './IsianUang'

describe('IsianUang', () => {
  it('ketikan dirapikan: titik ribuan otomatis (permintaan work owner: 20.000.000)', () => {
    expect(rapikanUang('20000000')).toBe('20.000.000')
    expect(rapikanUang('20.000.0001')).toBe('200.000.001')
    expect(rapikanUang('1234,5')).toBe('1.234,5')
    expect(rapikanUang('1234,')).toBe('1.234,')
    expect(rapikanUang('0007')).toBe('7')
    expect(rapikanUang('abc')).toBe('')
    expect(rapikanUang(',5')).toBe('0,5')
    expect(rapikanUang('1,123456789')).toBe('1,12345678')
    expect(rapikanUang('-5')).toBe('5')
  })

  it('kabel tanpa pemisah ribuan, desimal bertitik; pulang-pergi', () => {
    expect(keKabelUang('20.000.000')).toBe('20000000')
    expect(keKabelUang('1.234,5')).toBe('1234.5')
    expect(keKabelUang('1.234,')).toBe('1234')
    expect(keKabelUang('')).toBe('')
    expect(keTampilUang('20000000.50')).toBe('20.000.000,50')
    expect(keTampilUang('99999999999999999999.12345678')).toBe('99.999.999.999.999.999.999,12345678')
    expect(keKabelUang(keTampilUang('3000000000'))).toBe('3000000000')
  })

  it('render: nilai awal tampil berformat; tanpa Number()/parseFloat', () => {
    const html = renderToStaticMarkup(<IsianUang label="UJI" value="20000000" onChange={() => {}} />)
    expect(html).toContain('value="20.000.000"')
    expect(html).toContain('inputMode="decimal"')
    const sumber = readFileSync(join(__dirname, 'IsianUang.tsx'), 'utf8')
    expect(sumber).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|parseInt/)
  })
})
