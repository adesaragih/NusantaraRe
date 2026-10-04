import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { BINGKAI_PENAMPIL, PARAM_PENAMPIL, PENAMPIL_OFFICE } from './penampilOffice'

const PANEL = readFileSync(join(__dirname, 'components', 'PanelLampiran.tsx'), 'utf8')

describe('View Office Online - penampil kantor (DownloadAttProdName_Act 7 b1103)', () => {
  it('alamat penampil dan parameter VERBATIM dari XML', () => {
    expect(PENAMPIL_OFFICE.endsWith('view.officeapps.live.com/op/view.aspx')).toBe(true)
    expect(PARAM_PENAMPIL).toBe('src')
  })

  it('form GET menyandikan URL bertanda tangan sebagai satu nilai kueri (@encodeURL)', () => {
    const u = 'UJI objek/Contract Doc.pptx?X-Goog-Signature=a&b=c'
    const kueri = new URLSearchParams([[PARAM_PENAMPIL, u]]).toString()
    expect(kueri).toBe('src=UJI+objek%2FContract+Doc.pptx%3FX-Goog-Signature%3Da%26b%3Dc')
    expect(new URLSearchParams(kueri).get(PARAM_PENAMPIL)).toBe(u)
  })

  // Permintaan work owner 03-10-2026: "dari popup atau windows baru (bukan tab baru)" - popup di dalam aplikasi; jendela
  // peramban baru menuntut `window.open`, yang dilarang penjaga lintas-modul `unduhdokumen.test.ts`.
  it('View Office Online membuka penampil DI POPUP: form GET tersembunyi berbingkai, bukan tab baru', () => {
    expect(BINGKAI_PENAMPIL).not.toBe('_blank')
    expect(PANEL).toContain('<form ref={formOffice} method="get" action={PENAMPIL_OFFICE} target={BINGKAI_PENAMPIL} hidden>')
    expect(PANEL).toContain('<input ref={urlOffice} type="hidden" name={PARAM_PENAMPIL} />')
    expect(PANEL).toMatch(/<iframe[^>]*name=\{BINGKAI_PENAMPIL\}/)
    expect(PANEL).not.toMatch(/window\.open|location\.|target="_blank"/)
  })

  it('View pdf / gambar di popup: hanya objek URL lokal yang dipasang sebagai sumber', () => {
    const pasang = [...PANEL.matchAll(/\.src\s*=(?!=)\s*([^\n;]+)/g)].map((m) => m[1]!.trim())
    expect(pasang.length).toBeGreaterThan(0)
    expect(pasang.every((v) => v === 'objekURL')).toBe(true)
    expect(PANEL).toContain('const objekURL = URL.createObjectURL(')
    expect(PANEL).toContain('URL.revokeObjectURL(')
  })
})
