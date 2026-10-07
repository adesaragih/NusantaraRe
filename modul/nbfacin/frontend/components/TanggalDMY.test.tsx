// Kotak tanggal dd/mm/yyyy (permintaan work owner 03-10-2026) - fungsi format + render.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import TanggalDMY, { keKabel, keTampil, rapikanKetikan } from './TanggalDMY'

describe('format dd/mm/yyyy', () => {
  it('ketikan dirapikan: hanya angka, garis miring otomatis, paling panjang 10', () => {
    expect(rapikanKetikan('1')).toBe('1')
    expect(rapikanKetikan('103')).toBe('10/3')
    expect(rapikanKetikan('10032026')).toBe('10/03/2026')
    expect(rapikanKetikan('10/03/2026')).toBe('10/03/2026')
    expect(rapikanKetikan('10-03-2026999')).toBe('10/03/2026')
    expect(rapikanKetikan('ab10cd')).toBe('10')
  })

  it('dd/mm/yyyy -> kabel DD-MM-YYYY hanya bila lengkap dan tanggalnya ada', () => {
    expect(keKabel('10/03/2026')).toBe('10-03-2026')
    expect(keKabel('29/02/2024')).toBe('29-02-2024')
    expect(keKabel('29/02/2026')).toBe('')
    expect(keKabel('31/04/2026')).toBe('')
    expect(keKabel('10/3/2026')).toBe('')
    expect(keKabel('')).toBe('')
  })

  it('kabel -> tampil dd/mm/yyyy; bentuk lain kosong', () => {
    expect(keTampil('10-03-2026')).toBe('10/03/2026')
    expect(keTampil('2026-03-10')).toBe('')
    expect(keTampil('')).toBe('')
  })
})

describe('render', () => {
  const h = (value: string) =>
    renderToStaticMarkup(
      <TanggalDMY label="UJI Tanggal" value={value} onChange={() => {}} required labelKalender="Pilih tanggal" pesanFormat="salah" />,
    )

  it('kotak terlihat = teks dd/mm/yyyy, bukan input date bawaan', () => {
    const html = h('10-03-2026')
    expect(html).toMatch(/type="text" inputMode="numeric" placeholder="dd\/mm\/yyyy" value="10\/03\/2026"/)
    // Satu-satunya input date adalah pemilih tersembunyi yang dibuka tombol kalender.
    expect(html.match(/type="date"/g)).toHaveLength(1)
    expect(html).toMatch(/class="nbf-tanggal__pemilih" type="date"[^>]*value="2026-03-10"/)
    expect(html).toContain('aria-label="Pilih tanggal"')
  })

  it('tanda wajib dan tanpa galat saat kosong', () => {
    const html = h('')
    expect(html).toContain('<span class="field__req">*</span>')
    expect(html).not.toContain('field__error')
  })
})
