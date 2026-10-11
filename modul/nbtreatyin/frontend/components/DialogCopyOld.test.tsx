// Popup Copy Old NB (perintah work owner 07-10-2026) dan tombolnya di kepala portal - render statis. Fixture UJI-.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { DokumenLama } from '../api'
import { COPY_OLD, TOMBOL } from '../labels'
import { KepalaPortal } from '../pages/PortalNBTreatyIn'
import { TabelLama } from './DialogCopyOld'

const dok = (id: string, boleh: boolean): DokumenLama => ({
  id,
  noOffer: 'UJI-M1',
  noPolis: 'UJI-POL-1',
  insuredName: 'UJI TERTANGGUNG',
  businessName: 'UJI GRUP',
  sobName: 'UJI SOB',
  cedingCoName: 'UJI CEDING',
  tglProd: '2017-10-02 08:00:00',
  bolehDisalin: boleh,
  alasan: boleh ? [] : ['the old JSON cannot be read'],
})

const kosong = () => {}

describe('Copy Old NB', () => {
  it('tombol Copy Old hanya bila superadmin, di samping Create', () => {
    const tanpa = renderToStaticMarkup(
      <KepalaPortal copyOld={false} sibuk={false} onCopyOld={kosong} onBuat={kosong} />,
    )
    expect(tanpa).not.toContain(COPY_OLD.tombol)
    const html = renderToStaticMarkup(<KepalaPortal copyOld sibuk={false} onCopyOld={kosong} onBuat={kosong} />)
    const tombol = [...html.matchAll(/<button[^>]*>([^<]*)<\/button>/g)].map((m) => m[1])
    expect(tombol).toEqual([COPY_OLD.tombol, TOMBOL.create])
  })

  it('baris yang tidak dapat disalin: centang terkunci, alasan tampil; kolom NB', () => {
    const html = renderToStaticMarkup(
      <TabelLama
        tampil={[dok('NB-1', true), dok('NB-2', false)]}
        terpilih={new Set(['NB-1'])}
        proses={false}
        onCentang={kosong}
        onCentangSemua={kosong}
      />,
    )
    const kepala = [...html.matchAll(/<th[^>]*>(?:<input[^>]*\/>)?([^<]*)<\/th>/g)]
      .map((m) => m[1])
      .filter((x) => x !== '')
    expect(kepala).toEqual(Object.values(COPY_OLD.kolom))
    expect(html).toMatch(/<input type="checkbox" aria-label="Select NB-1" checked=""\/>/)
    expect(html).toMatch(/<input type="checkbox" aria-label="Select NB-2" disabled=""\/>/)
    expect(html).toContain('Cannot be copied: the old JSON cannot be read')
    expect(html).toContain('<td>UJI TERTANGGUNG</td>')
    expect(html).toContain('<td>02-10-2017</td>')
  })
})

// WO 07-10-2026 "TAMPILAN NYA HANYA NB-XXX AJA, BERLAKU NB DAN EDM TREATY": ID kasus salinan Copy Old = IDPEGA Pega
// utuh (`<kelas> <pyID>`); layar hanya menampilkan pyID, kunci buka / kirim tetap ID utuh.
describe('Copy Old - nomor dokumen Pega', () => {
  it('kolom nomor dan label centang tampil pyID; centang tetap berkunci ID utuh', () => {
    const id = 'ASM-FW-GISFW-WORK UJI-NB-9'
    const html = renderToStaticMarkup(
      <TabelLama
        tampil={[dok(id, true)]}
        terpilih={new Set([id])}
        proses={false}
        onCentang={kosong}
        onCentangSemua={kosong}
      />,
    )
    expect(html).toContain('<td>UJI-NB-9</td>')
    expect(html).toContain('aria-label="Select UJI-NB-9" checked=""')
    expect(html).not.toContain('ASM-FW-GISFW-WORK')
  })
})
