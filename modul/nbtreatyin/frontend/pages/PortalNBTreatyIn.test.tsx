// W6 audit silang putaran 3: portal = unsur yang TAMPIL di `Section/SFAPortal_OpportunitiesList`
// (dibaca ulang 04-10-2026). Halaman DIRENDER (react-dom/server); harapan diketik dari XML.
//
//   grid `GetListOpportunity` 6 kolom: LABEL "Offer No" (`.TextNoQuotation`), "Name"
//   (`.Name` pxLink openWorkByHandle - `.Name` ditulis NOL rule, RALAT tiket 11:
//   tautan dipindah ke Offer No, kolom tidak dirender), "Group Business", "Insured Name",
//   "Marketing", "Status". Tidak ada kolom "Position" / "No Polis".
//   baris saringan: `.FilterTermForOpportunity` + ikon pengosong C[1.2] (pxIcon
//   `pyImage webwb/pyiconclearfield.png`, click -> setValue "" -> postValue, TANPA refresh)
//   + tombol Filter.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { RingkasanKasus } from '../api'
import { KOLOM_PORTAL, PORTAL } from '../labels'
import PortalNBTreatyIn, { TabelPortal } from './PortalNBTreatyIn'

const baris: RingkasanKasus = {
  id: 'NB-1',
  businessName: 'UJI BISNIS',
  insuredName: 'UJI TERTANGGUNG',
  marketingName: 'UJI MO',
  nbStatus: 'UJI STATUS',
  positionNote: 'ReasTreatyInAdmin',
  noPolis: 'UJI-NOPOL',
  statusWork: 'Input Realitation',
  tglCreate: '2026-10-03 09:00:00',
}

describe('portal NB Treaty In = SFAPortal_OpportunitiesList', () => {
  const html = renderToStaticMarkup(<TabelPortal baris={[baris]} onBuka={() => {}} />)
  const judul = [...html.matchAll(/<th[^>]*>([^<]*)<\/th>/g)].map((m) => m[1])

  it('judul kolom grid GetListOpportunity, berurutan, tanpa Position / No Polis', () => {
    expect(judul).toEqual(['Offer No', 'Group Business', 'Insured Name', 'Marketing', 'Status'])
    expect(Object.values(KOLOM_PORTAL)).toEqual(judul)
    expect(html).not.toContain('Position')
    expect(html).not.toContain('No Polis')
    expect(html).not.toContain('UJI-NOPOL')
  })

  it('isi sel: .TextNoQuotation (tautan pembuka berkas), BusinessName, InsuredName, MarketingName, NBStatus', () => {
    const sel = [...html.matchAll(/<td[^>]*>(.*?)<\/td>/g)].map((m) => (m[1] ?? '').replace(/<[^>]+>/g, ''))
    expect(sel).toEqual(['NB-1', 'UJI BISNIS', 'UJI TERTANGGUNG', 'UJI MO', 'UJI STATUS'])
    expect(html).toContain('<button type="button" class="nbti__tautan">NB-1</button>')
  })

  it('baris saringan: kotak, ikon pengosong (pyiconclearfield), tombol Filter', () => {
    const portal = renderToStaticMarkup(<PortalNBTreatyIn onBuka={() => {}} />)
    expect(portal).toContain(`placeholder="${PORTAL.placeholder}"`)
    expect(portal).toContain(`aria-label="${PORTAL.bersihkan}"`)
    expect(portal.indexOf(`aria-label="${PORTAL.bersihkan}"`)).toBeLessThan(portal.indexOf('>Filter</button>'))
  })
})
