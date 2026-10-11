// Panel "Lampiran" tanpa pil tab tunggal di atasnya (perintah work owner 11-10-2026: "ITU DI KLAIM PROP DI HAPUS AJA!");
// judul panel tetap.

import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { CP } from '../labels'
import PanelLampiran from './PanelLampiran'

describe('panel lampiran', () => {
  it('tanpa strip tab, judul panel tetap', () => {
    const html = renderToStaticMarkup(
      createElement(PanelLampiran, { id: 'CLMP-UJI001', hanyaLihat: true, bolehSimpan: false, onSimpan: () => {} }),
    )
    expect(html).not.toContain('role="tablist"')
    expect(html).toContain(`<h3 class="panel__title">${CP.lampiran}</h3>`)
  })
})
