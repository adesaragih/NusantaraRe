// Lampiran klaim induk di layar komite (perintah work owner 11-10-2026): grid hanya kategori yang dikirim server (yang
// berisi), tombol attach (Add attachment, Upload File) NONAKTIF, View File aktif.

import { createElement } from 'react'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { LAMPIRAN_KCFI } from './labels'
import PanelLampiran, { GridKategori } from './pages/PanelLampiran'

const tombol = (html: string, label: string) => {
  const m = new RegExp(`<button[^>]*aria-label="${label}"[^>]*>`).exec(html)
  expect(m, label).not.toBeNull()
  return m![0]
}

describe('lampiran klaim induk di layar komite', () => {
  it('grid: satu baris per kategori berisi; Upload File nonaktif, View File aktif', () => {
    const html = renderToStaticMarkup(
      createElement(GridKategori, {
        kategori: [
          { id: 'ADU', label: 'Approval', countAttach: 1 },
          { id: 'File', label: 'File', countAttach: 2 },
        ],
        onLihat: () => {},
      }),
    )
    expect(html.match(/<tr>/g)).toHaveLength(3) // kepala + 2 kategori
    expect(html).toContain('<td>Approval</td><td>1</td>')
    for (const k of ['Approval', 'File']) {
      expect(tombol(html, `${LAMPIRAN_KCFI.unggah} ${k}`)).toContain('disabled')
      expect(tombol(html, `${LAMPIRAN_KCFI.lihat} ${k}`)).not.toContain('disabled')
    }
  })

  it('panel: Add attachment nonaktif, Refresh aktif', () => {
    const html = renderToStaticMarkup(createElement(PanelLampiran, { id: 'KMT-UJI001' }))
    expect(html).toContain(`<h3 class="panel__title komiteclaimfacin__kartu-judul">${LAMPIRAN_KCFI.judul}</h3>`)
    expect(html).toMatch(new RegExp(`<button[^>]*disabled[^>]*>${LAMPIRAN_KCFI.tambah}</button>`))
    expect(html).toMatch(new RegExp(`<button type="button" class="btn btn--ghost btn--sm">${LAMPIRAN_KCFI.muatUlang}`))
  })
})
