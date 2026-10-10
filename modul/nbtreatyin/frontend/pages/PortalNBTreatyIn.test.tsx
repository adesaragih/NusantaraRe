// W6 audit silang putaran 3: portal = unsur yang TAMPIL di `Section/SFAPortal_OpportunitiesList`
// (dibaca ulang 04-10-2026). Halaman DIRENDER (react-dom/server); harapan diketik dari XML.
//
//   grid `GetListOpportunity` 6 kolom: LABEL "Offer No" (`.TextNoQuotation`), "Name"
//   (`.Name` pxLink openWorkByHandle - `.Name` ditulis NOL rule, RALAT tiket 11:
//   tautan dipindah ke Offer No, kolom tidak dirender), "Treaty Business" (perintah work owner 10-10-2026, dulu "Group Business"), "Insured Name",
//   "Marketing", "Status". Tidak ada kolom "Position" / "No Polis".
//   baris saringan: `.FilterTermForOpportunity` + ikon pengosong C[1.2] (pxIcon
//   `pyImage webwb/pyiconclearfield.png`, click -> setValue "" -> postValue, TANPA refresh)
//   + tombol Filter.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import type { RingkasanKasus } from '../api'
import { KOLOM_PORTAL, PORTAL, STATUS_PORTAL } from '../labels'
import PortalNBTreatyIn, { TabelPortal } from './PortalNBTreatyIn'

const baris: RingkasanKasus = {
  id: 'UJI-NB-1',
  treatyGroupName: 'UJI BISNIS',
  insuredName: 'UJI TERTANGGUNG',
  marketingName: 'UJI MO',
  nbStatus: 'UJI STATUS',
  positionNote: 'ReasTreatyInAdmin',
  noPolis: 'UJI-NOPOL',
  statusWork: 'Input Realitation',
  tglCreate: '2026-10-03 09:00:00',
  createOpName: 'UJI PEMBUAT',
  proportionalType: 'NonProportional',
  cedingCoName: '',
  startDate: '',
  tglUpdate: '',
  productionDate: '',
}

describe('portal NB Treaty In = SFAPortal_OpportunitiesList', () => {
  const html = renderToStaticMarkup(<TabelPortal baris={[baris]} onBuka={() => {}} />)
  const judul = [...html.matchAll(/<th[^>]*>([^<]*)<\/th>/g)].map((m) => m[1])

  it('judul kolom grid GetListOpportunity, berurutan, tanpa Position / No Polis', () => {
    // + User Create / Tanggal Create (keputusan work owner 06-10-2026; bukan kolom GetListOpportunity)
    expect(judul).toEqual([
      'Offer No',
      'Type',
      'Treaty Business',
      'Insured Name',
      'Marketing',
      'Status',
      'User Create',
      'Tanggal Create',
    ])
    expect(Object.values(KOLOM_PORTAL)).toEqual(judul)
    expect(html).not.toContain('Position')
    expect(html).not.toContain('No Polis')
    expect(html).not.toContain('UJI-NOPOL')
  })

  // Perintah work owner 07-10-2026: "KALO DAH RESOLVE TAMBAHIN KOLOM NOPOLISNYA", "SEKALIAN KELUARIN TANGGAL
  // PRODUKSINYA AJA DD-MM-YYYY" - hanya tab Resolved; In Progress tetap kolom di atas.
  it('tab Resolved: Policy Number sesudah Offer No, Production Date dd-mm-yyyy di akhir', () => {
    const b = { ...baris, statusWork: 'Resolved-Completed', productionDate: '2026-10-07 10:00:00' }
    const h = renderToStaticMarkup(<TabelPortal baris={[b]} onBuka={() => {}} selesai />)
    expect([...h.matchAll(/<th[^>]*>([^<]*)<\/th>/g)].map((m) => m[1])).toEqual([
      'Offer No',
      'Policy Number',
      'Type',
      'Treaty Business',
      'Insured Name',
      'Marketing',
      'Status',
      'User Create',
      'Tanggal Create',
      'Production Date',
    ])
    const sel = [...h.matchAll(/<td[^>]*>(.*?)<\/td>/g)].map((m) => (m[1] ?? '').replace(/<[^>]+>/g, ''))
    expect(sel[1]).toBe('UJI-NOPOL')
    expect(sel.at(-1)).toBe('07-10-2026')
  })

  it('isi sel: .TextNoQuotation (tautan pembuka berkas), BusinessName, InsuredName, MarketingName, NBStatus', () => {
    const sel = [...html.matchAll(/<td[^>]*>(.*?)<\/td>/g)].map((m) => (m[1] ?? '').replace(/<[^>]+>/g, ''))
    // Type = T_POLIS_QUOTATION.PROPORTIONAL_TYPE apa adanya, tanpa disingkat (perintah work owner 06-10-2026)
    expect(sel).toEqual([
      'UJI-NB-1',
      'NonProportional',
      'UJI BISNIS',
      'UJI TERTANGGUNG',
      'UJI MO',
      'UJI STATUS',
      'UJI PEMBUAT',
      '03-10-2026',
    ])

    expect(html).toContain('<button type="button" class="nbti__tautan">UJI-NB-1</button>')
  })

  it('baris saringan: kotak, ikon pengosong (pyiconclearfield), tombol Filter', () => {
    const portal = renderToStaticMarkup(<PortalNBTreatyIn onBuka={() => {}} />)
    expect(portal).toContain(`placeholder="${PORTAL.placeholder}"`)
    expect(portal).toContain(`aria-label="${PORTAL.bersihkan}"`)
    expect(portal.indexOf(`aria-label="${PORTAL.bersihkan}"`)).toBeLessThan(portal.indexOf('>Filter</button>'))
  })

  // switch di atas daftar (keputusan work owner 06-10-2026: "lihat yang lagi proses atau resolve, default ke proses")
  it('switch In Progress / Resolved di atas saringan, bawaan In Progress', () => {
    const tab = (html: string) =>
      [...html.matchAll(/<button[^>]*role="tab"[^>]*aria-selected="(true|false)"[^>]*>([^<]*)<\/button>/g)].map(
        (m) => `${m[2]}:${m[1]}`,
      )
    const portal = renderToStaticMarkup(<PortalNBTreatyIn onBuka={() => {}} />)
    expect(tab(portal)).toEqual(['In Progress:true', 'Resolved:false'])
    expect(portal.indexOf('role="tablist"')).toBeLessThan(portal.indexOf('role="search"'))
    const resolved = renderToStaticMarkup(
      <PortalNBTreatyIn onBuka={() => {}} status={STATUS_PORTAL[1]} onStatus={() => {}} />,
    )
    expect(tab(resolved)).toEqual(['In Progress:false', 'Resolved:true'])
  })

  // perintah work owner 06-10-2026: "statusnya yang resolve pake status resolve di work_polis aja"
  it('kolom Status: berkas Resolved memakai STATUS_WORK, selain itu NBStatus', () => {
    const status = (b: RingkasanKasus) =>
      [
        ...renderToStaticMarkup(<TabelPortal baris={[b]} onBuka={() => {}} />).matchAll(
          /<span class="nbti__status">([^<]*)<\/span>/g,
        ),
      ].map((m) => m[1])
    expect(status(baris)).toEqual(['UJI STATUS'])
    expect(status({ ...baris, statusWork: 'Resolved-Completed', nbStatus: 'UJI NB IS IN INBOX' })).toEqual([
      'Resolved-Completed',
    ])
    expect(status({ ...baris, statusWork: 'Resolved-Rejected', nbStatus: 'UJI DECLINED' })).toEqual([
      'Resolved-Rejected',
    ])
  })
})

// WO 07-10-2026 "TAMPILAN NYA HANYA NB-XXX AJA, BERLAKU NB DAN EDM TREATY": ID kasus salinan Copy Old = IDPEGA Pega
// utuh (`<kelas> <pyID>`); layar hanya menampilkan pyID, kunci buka / kirim tetap ID utuh.
describe('portal NB Treaty In - berkas salinan Copy Old', () => {
  it('Offer No tampil pyID saja', () => {
    const html = renderToStaticMarkup(
      <TabelPortal baris={[{ ...baris, id: 'ASM-FW-GISFW-WORK UJI-NB-9' }]} onBuka={() => {}} selesai />,
    )
    expect(html).toContain('<button type="button" class="nbti__tautan">UJI-NB-9</button>')
    expect(html).not.toContain('ASM-FW-GISFW-WORK')
  })
})
