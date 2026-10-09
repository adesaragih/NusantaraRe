// Jalur REVISI tab `Information & Submit` — `Section/TreatyInfoSubmit.xml`
// cell 9 (Comment), 20 (Submit biasa), 21 (Submit revisi), 22 (Decline
// offer). Tombol Revision daftar kontrak menyetel `RevisionState=1`.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import TabInfoSubmit from './components/TabInfoSubmit'

const tidakApa = () => undefined
const tombol = (html: string) => [...html.matchAll(/<button[^>]*>([^<]*)<\/button>/g)].map((m) => m[1])

describe('⭐ jalur revisi — tepat satu Submit, tanpa Decline offer', () => {
  it('mode lihat kontrak revisi: Submit (cell 21) tampil, Decline tidak', () => {
    const html = renderToStaticMarkup(
      <TabInfoSubmit mode="lihat" statusAkseptasi="" revisi onKirim={tidakApa} onTolak={tidakApa} />,
    )
    expect(tombol(html)).toEqual(['Submit'])
  })

  it('mode ubah kontrak revisi: tetap SATU Submit — ViewState disetel ulang ke 1', () => {
    const html = renderToStaticMarkup(
      <TabInfoSubmit mode="ubah" statusAkseptasi="" revisi onKirim={tidakApa} onTolak={tidakApa} />,
    )
    expect(tombol(html)).toEqual(['Submit'])
  })

  it('bukan revisi: mode lihat tanpa tombol, mode ubah Submit + Decline offer', () => {
    expect(tombol(renderToStaticMarkup(<TabInfoSubmit mode="lihat" statusAkseptasi="" onKirim={tidakApa} />))).toEqual([])
    expect(
      tombol(renderToStaticMarkup(<TabInfoSubmit mode="ubah" statusAkseptasi="" onKirim={tidakApa} onTolak={tidakApa} />)),
    ).toEqual(['Submit', 'Decline offer'])
  })

  it('Additional Information tetap terkunci di revisi; Comment tidak', () => {
    const html = renderToStaticMarkup(<TabInfoSubmit mode="lihat" statusAkseptasi="" revisi onKirim={tidakApa} />)
    const i = html.indexOf('<fieldset')
    expect(html.slice(i, i + 80)).toContain('disabled')
    // Comment di LUAR fieldset itu.
    expect(html.indexOf('Comment')).toBeGreaterThan(html.indexOf('</fieldset>'))
  })

  it('form tidak mematikan tab ini untuk kontrak revisi', () => {
    const form = readFileSync(join(__dirname, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8')
    expect(form).toContain('disabled={!bisaUbah && !(revisi && tabTampil === TAB_REVISI)}')
    expect(form).toContain("const revisi = tersimpanPenampung.RevisionState === '1'")
    expect(form).toContain('revisi={revisi}')
  })
})
