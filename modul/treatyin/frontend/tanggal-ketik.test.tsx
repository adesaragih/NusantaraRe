// Uji kotak TANGGAL YANG BISA DIKETIK — permintaan pemakai 7 Oktober 2026.
// ⛔ Contoh yang SAMA ada di modul treatyin dan treatyinadjustment: kedua
// salinan `TanggalKetik.tsx` harus berperilaku sama.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { KotakTanggalKetik, keTeksTanggal, uraiTanggalKetik } from './components/TanggalKetik'

describe('uraiTanggalKetik — ketikan → bentuk kabel DD-MM-YYYY', () => {
  it.each([
    ['18/01/2024', '18-01-2024'],
    ['18-01-2024', '18-01-2024'],
    ['18.01.2024', '18-01-2024'],
    ['18 01 2024', '18-01-2024'],
    ['1/8/2024', '01-08-2024'],
    ['18012024', '18-01-2024'],
    ['180124', '18-01-2024'],
    ['18/01/24', '18-01-2024'],
    ['20240118', '18-01-2024'],
    ['2024-01-18', '18-01-2024'],
    ['  29/02/2024 ', '29-02-2024'],
    ['', ''],
  ])('%s → %s', (ketik, mau) => {
    expect(uraiTanggalKetik(ketik)).toBe(mau)
  })

  it.each(['31/02/2024', '29/02/2025', '13/13/2024', '00/01/2024', 'besok', '18/01', '1801202', '01/01/1899'])(
    '%s bukan tanggal → null (nilai lama tetap)',
    (ketik) => {
      expect(uraiTanggalKetik(ketik)).toBeNull()
    },
  )
})

describe('kotak', () => {
  it('nilai tersimpan (YYYYMMDD / kabel / stempel Pega) tampil DD/MM/YYYY', () => {
    expect(keTeksTanggal('20240118')).toBe('18/01/2024')
    expect(keTeksTanggal('18-01-2024')).toBe('18/01/2024')
    expect(keTeksTanggal('20240118T170000.000 GMT')).toBe('18/01/2024')
    expect(keTeksTanggal('')).toBe('')
  })

  it('kotak TEKS yang dapat diketik + ikon kalender + kalender bawaan tak terlihat', () => {
    const html = renderToStaticMarkup(<KotakTanggalKetik label="Due Date" value="20240118" onChange={() => undefined} />)
    expect(html).toContain('type="text"')
    expect(html).toContain('value="18/01/2024"')
    expect(html).toContain('aria-label="Pick from the calendar — Due Date"')
    expect(html).toMatch(/type="date"[^>]*value="2024-01-18"|value="2024-01-18"[^>]*type="date"/)
    // Kosong = kosong, bukan `dd/mm/yyyy` yang terbaca seperti nilai.
    expect(renderToStaticMarkup(<KotakTanggalKetik label="X" value="" onChange={() => undefined} />)).toContain('type="text" autoComplete="off" aria-label="X" title=')
  })
})
