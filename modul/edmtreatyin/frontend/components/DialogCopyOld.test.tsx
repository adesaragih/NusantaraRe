// Popup Copy Old (perintah work owner 07-10-2026 "SAMA SEPERTI MASTER PRODUCTNAME LIFE, KHUSUS BUAT SUPERUSER") dan
// tombolnya di kepala portal - render statis (react-dom/server). Fixture UJI-.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { DokumenLama } from '../api'
import { COPY_OLD } from '../labels'
import { KepalaPortal } from '../pages/PortalEDMTreatyIn'
import { TabelLama } from './DialogCopyOld'

const dok = (id: string, boleh: boolean): DokumenLama => ({
  id,
  noPolis: 'UJI-POL-1',
  edmNo: 'UJI-POL-1/E01',
  prodKe: 1,
  edmType: '3',
  sobName: 'UJI SOB',
  cedingCoName: 'UJI CEDING',
  tglProd: '2017-10-02 08:00:00',
  bolehDisalin: boleh,
  alasan: boleh ? [] : ['UJI previous generation missing'],
})

const kosong = () => {}

describe('Copy Old', () => {
  it('tombol Copy Old hanya bila superadmin, di samping Create', () => {
    const tanpa = renderToStaticMarkup(<KepalaPortal copyOld={false} onCopyOld={kosong} onBuat={kosong} />)
    expect(tanpa).not.toContain(COPY_OLD.tombol)
    const html = renderToStaticMarkup(<KepalaPortal copyOld onCopyOld={kosong} onBuat={kosong} />)
    const tombol = [...html.matchAll(/<button[^>]*>([^<]*)<\/button>/g)].map((m) => m[1])
    expect(tombol).toEqual([COPY_OLD.tombol, 'Create New Addendum Treaty'])
  })

  it('baris yang tidak dapat disalin: centang terkunci, alasan tampil; baris lain dapat dicentang', () => {
    const html = renderToStaticMarkup(
      <TabelLama
        tampil={[dok('EDMT-1', true), dok('EDMT-2', false)]}
        terpilih={new Set(['EDMT-1'])}
        proses={false}
        onCentang={kosong}
        onCentangSemua={kosong}
      />,
    )
    expect(html).toMatch(/<input type="checkbox" aria-label="Select EDMT-1" checked=""\/>/)
    expect(html).toMatch(/<input type="checkbox" aria-label="Select EDMT-2" disabled=""\/>/)
    expect(html).toContain('Cannot be copied: UJI previous generation missing')
    // label EDM Type DT TreatyEDMListType, Production Date format sistem
    expect(html).toContain('<td>Adjustment Premium</td>')
    expect(html).toContain('<td>02-10-2017</td>')
    // kepala: centang semua tercentang bila semua baris yang dapat disalin terpilih
    expect(html).toMatch(/aria-label="Select all" checked=""/)
  })
})

// WO 07-10-2026 "TAMPILAN NYA HANYA NB-XXX AJA, BERLAKU NB DAN EDM TREATY": ID kasus salinan Copy Old = IDPEGA Pega
// utuh (`<kelas> <pyID>`); layar hanya menampilkan pyID, kunci buka / kirim tetap ID utuh.
describe('Copy Old - nomor dokumen Pega', () => {
  it('kolom nomor dan label centang tampil pyID; centang tetap berkunci ID utuh', () => {
    const id = 'ASM-FW-GISFW-WORK UJI-EDMT-9'
    const html = renderToStaticMarkup(
      <TabelLama
        tampil={[dok(id, true)]}
        terpilih={new Set([id])}
        proses={false}
        onCentang={kosong}
        onCentangSemua={kosong}
      />,
    )
    expect(html).toContain('<td>UJI-EDMT-9</td>')
    expect(html).toContain('aria-label="Select UJI-EDMT-9" checked=""')
    expect(html).not.toContain('ASM-FW-GISFW-WORK')
  })
})
