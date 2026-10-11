// Portal = unsur yang TAMPIL di `Section/SFAPortal_Endorsement_Treaty` (dibaca 06-10-2026). Halaman DIRENDER
// (react-dom/server, pola uji portal NB); `fetch` distub - render statis tidak menjalankan efek, jadi tidak memanggilnya.

import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import type { RingkasanKasus } from '../api'
import { KOLOM_PORTAL } from '../labels'
import BuatEDM, { tampilTombolBuat } from '../components/BuatEDM'
import PortalEDMTreatyIn, { TabelPortal } from './PortalEDMTreatyIn'

const baris: RingkasanKasus = {
  id: 'UJI-EDMT-1',
  noOffer: 'UJI-OFFER',
  noPolis: 'UJI-POLIS',
  edmNo: 'UJI-POLIS/E01',
  sobName: 'UJI SOB',
  cedingCoName: 'UJI CEDING',
  edmType: '3',
  proportionalType: 'NonProportional',
  marketingName: 'UJI MO',
  nbStatus: 'UJI STATUS',
  statusWork: 'Input Realitation',
  positionNote: 'ReasTreatyInAdmin',
  tglCreate: '2026-10-03 09:00:00',
  startDate: '',
  tglUpdate: '',
}

const ambil = vi.fn(() => Promise.reject(new Error('UJI: fetch tidak boleh dipanggil render statis')))
beforeEach(() => vi.stubGlobal('fetch', ambil))
afterEach(() => vi.unstubAllGlobals())

describe('portal EDM Treaty In = SFAPortal_Endorsement_Treaty', () => {
  // perintah work owner 07-10-2026 (screenshot tab Resolved): "KALO DAH RESOLVE STATUS NYA PAKE STATUS RESOLVE" - sama
  // dengan portal NB (STATUS_WORK, bukan NBStatus terakhir)
  it('kolom Status: berkas Resolved memakai STATUS_WORK, selain itu NBStatus', () => {
    const status = (b: RingkasanKasus) =>
      [
        ...renderToStaticMarkup(<TabelPortal baris={[b]} onBuka={() => {}} selesai />).matchAll(
          /<span class="edmt__status">([^<]*)<\/span>/g,
        ),
      ].map((m) => m[1])
    expect(status(baris)).toEqual(['UJI STATUS'])
    expect(status({ ...baris, statusWork: 'Resolved-Completed', nbStatus: 'UJI EDMT IS IN DEPT HEAD INBOX' })).toEqual([
      'Resolved-Completed',
    ])
    expect(status({ ...baris, statusWork: 'Resolved-Rejected', nbStatus: 'UJI DECLINED' })).toEqual([
      'Resolved-Rejected',
    ])
  })

  it('grid InboxEDM_RD2: 10 judul kolom VERBATIM berurutan (EDM Number kembar)', () => {
    const html = renderToStaticMarkup(<TabelPortal baris={[baris]} onBuka={() => {}} />)
    const judul = [...html.matchAll(/<th[^>]*>([^<]*)<\/th>/g)].map((m) => m[1])
    expect(judul).toEqual(KOLOM_PORTAL.map((k) => k.judul))
    expect(judul).toEqual([
      'EDM Number',
      'Offer No',
      'Policy Number',
      'EDM Number',
      'SOB',
      'Ceding',
      'EDM Type',
      'Proportional Type',
      'Marketing',
      'Status',
    ])
  })

  it('tab Resolved: kolom tambahan Production Date DD-MM-YYYY di ujung, Policy Number tetap kolom XML; In Progress tanpa', () => {
    const selesai = { ...baris, statusWork: 'Resolved-Completed', positionNote: '', tglProd: '2026-10-07 11:30:00' }
    const html = renderToStaticMarkup(<TabelPortal baris={[selesai]} onBuka={() => {}} selesai />)
    const kepala = [...html.matchAll(/<th[^>]*>(.*?)<\/th>/g)].map((m) => m[1])
    expect(kepala.at(-1)).toBe('Production Date')
    expect(kepala).toContain('Policy Number')
    expect(html).toContain('data-label="Production Date">07-10-2026</td>')
    // berkas selesai dibuka hanya-baca: tautan aktif walau posisinya kosong
    expect(html).not.toMatch(/<button[^>]*disabled[^>]*>UJI-EDMT-1<\/button>/)
    const proses = renderToStaticMarkup(<TabelPortal baris={[baris]} onBuka={() => {}} />)
    expect(proses).not.toContain('Production Date')
  })

  it('isi sel: nilai DB apa adanya, EDM Type berteks DT TreatyEDMListType; status chip NBStatus', () => {
    const html = renderToStaticMarkup(<TabelPortal baris={[baris]} onBuka={() => {}} />)
    const sel = [...html.matchAll(/<td[^>]*>(.*?)<\/td>/g)].map((m) => (m[1] ?? '').replace(/<[^>]+>/g, ''))
    expect(sel).toEqual([
      'UJI-EDMT-1',
      'UJI-OFFER',
      'UJI-POLIS',
      'UJI-POLIS/E01',
      'UJI SOB',
      'UJI CEDING',
      'Adjustment Premium',
      'NonProportional',
      'UJI MO',
      'UJI STATUS',
    ])
    expect(html).toContain('<span class="edmt__status">UJI STATUS</span>')
  })

  it('tautan EDM Number aktif hanya bila PositionNote = ReasTreatyInAdmin', () => {
    const aktif = renderToStaticMarkup(<TabelPortal baris={[baris]} onBuka={() => {}} />)
    expect(aktif).toContain('<button type="button" class="edmt__tautan">UJI-EDMT-1</button>')
    const mati = renderToStaticMarkup(
      <TabelPortal baris={[{ ...baris, positionNote: 'ReasTreatyInSecHead' }]} onBuka={() => {}} />,
    )
    expect(mati).toContain('<button type="button" class="edmt__tautan" disabled="">UJI-EDMT-1</button>')
  })

  it('50 baris per halaman (pyPageSize 50, Numeric)', () => {
    const banyak = Array.from({ length: 51 }, (_, i) => ({ ...baris, id: `UJI-EDMT-${i + 1}` }))
    const html = renderToStaticMarkup(<TabelPortal baris={banyak} onBuka={() => {}} />)
    expect((html.match(/<tr>/g) ?? []).length).toBe(51) // 1 baris judul + 50 baris
    expect(html).toContain('aria-current="page">1</button>')
    expect(html).toContain('>2</button>')
  })

  it('kepala: judul "Addendum Treaty", tombol Create, switch In Progress / Resolved (aturan portal NB), saring (nomor EDM / polis / insured ...) + Filter, tombol Refresh', () => {
    const html = renderToStaticMarkup(<PortalEDMTreatyIn onBuka={() => {}} />)
    expect(html).toContain('<h2 class="inbox__judul">Addendum Treaty</h2>')
    expect(html).toContain('>Create New Addendum Treaty</button>')
    expect(html).toContain('placeholder="Search EDM no, master ID, policy no, insured, business, ceding, marketing..."')
    expect(html).toContain('aria-label="Filter Term for Endorsement"')
    expect(html).toContain('>Filter</button>')
    expect(html).toContain('title="Refresh EDM Grid">Refresh</button>')
    // keputusan work owner 07-10-2026: aturan portal NB berlaku - switch bawaan In Progress
    expect(html).toContain('aria-selected="true" class="tabs__item tabs__item--aktif">In Progress</button>')
    expect(html).toContain('>Resolved</button>')
    expect(ambil).not.toHaveBeenCalled()
  })
})

describe('layar Create TreatyCreateEdm', () => {
  it('tiga medan VERBATIM; No Master Treaty tidak dapat diisi; tombol Create tersembunyi selama No Polis kosong', () => {
    const html = renderToStaticMarkup(
      <BuatEDM opsiJenis={[{ value: '3', label: '3' }]} onBuka={() => {}} onTutup={() => {}} />,
    )
    expect(html).toContain('Create EDM')
    expect(html).toContain('No Polis Treaty')
    expect(html).toContain('No Master Treaty')
    expect(html).toContain('Source of Change')
    expect(html).toMatch(/<input[^>]*readOnly=""/i)
    expect(html).not.toContain('>Create</button>')
  })

  it('syarat tombol Create: TrtERR kosong, No Polis dan No Master terisi (Source of Change tidak disyaratkan)', () => {
    expect(tampilTombolBuat('', 'UJI-POLIS', 'UJI-MASTER')).toBe(true)
    expect(
      tampilTombolBuat("There's EDM with this policy no that haven't finish yet!", 'UJI-POLIS', 'UJI-MASTER'),
    ).toBe(false)
    expect(tampilTombolBuat('', '', 'UJI-MASTER')).toBe(false)
    expect(tampilTombolBuat('', 'UJI-POLIS', '')).toBe(false)
  })
})

// WO 07-10-2026 "TAMPILAN NYA HANYA NB-XXX AJA, BERLAKU NB DAN EDM TREATY": ID kasus salinan Copy Old = IDPEGA Pega
// utuh (`<kelas> <pyID>`); layar hanya menampilkan pyID, kunci buka / kirim tetap ID utuh.
describe('portal EDM Treaty In - berkas salinan Copy Old', () => {
  it('EDM Number tampil pyID saja', () => {
    const html = renderToStaticMarkup(
      <TabelPortal baris={[{ ...baris, id: 'ASM-FW-GISFW-WORK UJI-EDMT-9' }]} onBuka={() => {}} selesai />,
    )
    expect(html).toContain('>UJI-EDMT-9</button>')
    expect(html).not.toContain('ASM-FW-GISFW-WORK')
  })
})
