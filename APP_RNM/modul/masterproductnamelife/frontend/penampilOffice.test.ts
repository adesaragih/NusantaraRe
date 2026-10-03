import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { PARAM_PENAMPIL, PENAMPIL_OFFICE } from './penampilOffice'

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

  it('link View Office Online yang ada langsung membuka penampil (satu klik, form GET tersembunyi, tanpa jendela tambahan)', () => {
    expect(PANEL).toContain('<form ref={formOffice} method="get" action={PENAMPIL_OFFICE} target="_blank" hidden>')
    expect(PANEL).toContain('<input ref={urlOffice} type="hidden" name={PARAM_PENAMPIL} />')
    expect(PANEL).not.toMatch(/window\.open|location\.|setOffice|<Modal\s+judul=\{LAMPIRAN_MPNL\.viewOffice\}/)
  })
})
