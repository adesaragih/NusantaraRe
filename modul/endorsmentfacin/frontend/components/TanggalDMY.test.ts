// Pembantu kotak tanggal dd/mm/yyyy (salinan pola NB FacIn) - kabel DD-MM-YYYY.

import { describe, expect, it } from 'vitest'

import { keKabel, keTampil, rapikanKetikan } from './TanggalDMY'

describe('TanggalDMY', () => {
  it('ketikan dirapikan: garis miring disisipkan, paling panjang dd/mm/yyyy', () => {
    expect(rapikanKetikan('0')).toBe('0')
    expect(rapikanKetikan('0102')).toBe('01/02')
    expect(rapikanKetikan('01022026999')).toBe('01/02/2026')
    expect(rapikanKetikan('01-02-2026')).toBe('01/02/2026')
  })

  it('kabel hanya untuk tanggal lengkap yang ada', () => {
    expect(keKabel('01/02/2026')).toBe('01-02-2026')
    expect(keKabel('31/02/2026')).toBe('')
    expect(keKabel('01/02/20')).toBe('')
    expect(keTampil('01-02-2026')).toBe('01/02/2026')
    expect(keTampil('')).toBe('')
  })
})
