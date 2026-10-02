// Paritas layar portal Opportunity dengan harness Pega `SFAPortalOpportunities` - tiket 25.
// Halaman DIRENDER (react-dom/server), bukan dibaca sebagai teks sumber.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { KEPALA_PORTAL, KOLOM_PORTAL, SARING_PORTAL, TEKS_PORTAL } from '../labels'
import PortalOpportunity from './PortalOpportunity'

const HTML = renderToStaticMarkup(<PortalOpportunity onBuat={() => {}} />)

/** Teks setiap `<th>` berurutan. */
const judulKolom = [...HTML.matchAll(/<th[^>]*>([^<]*)<\/th>/g)].map((m) => m[1])

describe('PortalOpportunity = unsur yang TAMPIL di Pega', () => {
  it('akar halaman memasang kelas modul (gaya terisolasi)', () => {
    expect(HTML.startsWith('<div class="nbfacin">')).toBe(true)
  })

  it('judul dan tombol Create opportunity, berurutan', () => {
    const judul = HTML.indexOf(`>${KEPALA_PORTAL.judul.label}</h4>`)
    const tombol = HTML.indexOf(`>${KEPALA_PORTAL.buat.label}</button>`)
    expect(judul).toBeGreaterThan(-1)
    expect(tombol).toBeGreaterThan(judul)
  })

  it('kotak saring: placeholder korpus, label sebagai nama aksesibel (pyIncludeLabel=false)', () => {
    expect(HTML).toContain(`placeholder="${SARING_PORTAL.placeholder.label}"`)
    expect(HTML).toContain(`aria-label="${SARING_PORTAL.label.label}"`)
    expect(HTML).not.toContain(`>${SARING_PORTAL.label.label}<`)
  })

  it('judul kolom = grid GetListOpportunityF, berurutan, kolom 3 dan 7 kosong', () => {
    expect(judulKolom).toEqual(KOLOM_PORTAL.map((k) => k.label))
    expect(judulKolom).toEqual(['Offer No', 'Name', '', 'Group Business', 'Insured Name', 'Marketing', '', 'Status'])
  })

  it('unsur tersembunyi permanen di Pega TIDAK dirender', () => {
    for (const t of ['Stage view', 'List view', '>All<', 'Individual', 'Corporate', 'Phase :', 'Proposal', 'Closed', 'Export', 'Refresh']) {
      expect(HTML).not.toContain(t)
    }
    // `Create Opportunity` (O besar, sel 71) tersembunyi; yang tampil `Create opportunity` (sel 72).
    expect(HTML).not.toContain('Create Opportunity')
  })

  it('daftar tanpa sumber data = BelumTersedia, bukan tabel kosong (B-2)', () => {
    expect(HTML).toContain(`${TEKS_PORTAL.daftar} tidak dapat dimuat saat ini.`)
    expect(HTML).not.toContain('<tbody')
  })

  it('kotak saring, ikon hapus, dan Filter nonaktif selama daftar tidak dapat dimuat', () => {
    const saring = HTML.slice(HTML.indexOf('class="nbf-saring"'), HTML.indexOf('class="table-wrap"'))
    expect(saring.match(/disabled=""/g)).toHaveLength(3)
    expect(saring).toContain(`>${SARING_PORTAL.tombol.label}</button>`)
  })

  it('tombol Create opportunity hidup (membuka form, tiket 26)', () => {
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn btn--primary">${KEPALA_PORTAL.buat.label}</button>`))
  })
})
