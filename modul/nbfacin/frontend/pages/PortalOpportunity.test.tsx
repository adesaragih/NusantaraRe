// Paritas layar portal Opportunity dengan harness Pega `SFAPortalOpportunities` - tiket 25; daftar case NB tiket 32.
// Halaman DIRENDER (react-dom/server), bukan dibaca sebagai teks sumber.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { KEPALA_PORTAL, KOLOM_PORTAL, SARING_PORTAL } from '../labels'
import PortalOpportunity from './PortalOpportunity'

const HTML = renderToStaticMarkup(<PortalOpportunity onBuat={() => {}} onBuka={() => {}} />)
const SUMBER = readFileSync(join(__dirname, 'PortalOpportunity.tsx'), 'utf8').replace(/\r\n/g, '\n')

/** Teks setiap `<th>` berurutan. */
const judulKolom = [...HTML.matchAll(/<th[^>]*>([^<]*)<\/th>/g)].map((m) => m[1])

describe('PortalOpportunity = unsur yang TAMPIL di Pega', () => {
  it('akar halaman memasang kelas modul (gaya terisolasi)', () => {
    expect(HTML.startsWith('<div class="nbfacin">')).toBe(true)
  })

  it('judul dan tombol Create opportunity, berurutan', () => {
    const judul = HTML.indexOf(`<h2 class="inbox__judul">${KEPALA_PORTAL.judul.label}</h2>`)
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

  it('daftar case NB dimuat saat dibuka; belum ada baris sebelum jawaban datang (tiket 32)', () => {
    expect(SUMBER).toMatch(/useEffect\(\(\) => \{\s*void muat\('', 1\)/)
    expect(SUMBER).toContain('daftarCaseNB(cari, halaman)')
    expect(HTML).not.toContain('<tbody')
    expect(HTML).not.toContain('tidak dapat dimuat')
  })

  it('kotak saring dapat diisi; Filter / Enter mencari dari halaman 1; ✕ mengosongkan dan memuat ulang', () => {
    const saring = HTML.slice(HTML.indexOf('class="toolbar"'), HTML.indexOf('class="toolbar__spacer"'))
    expect(saring).not.toContain('disabled')
    expect(saring).toContain(`<button type="submit" class="btn btn--sm">${SARING_PORTAL.tombol.label}</button>`)
    expect(SUMBER).toContain('void muat(kotak, 1)')
    expect(SUMBER).toMatch(/setKotak\(''\)\s*void muat\('', 1\)/)
    expect(SUMBER).toContain('onPindah={(h) => void muat(kunciCari.current, h)}')
  })

  it('baris: Offer No = case id, Name tautan pembuka case, kolom 3 dan 7 kosong, Status (permintaan work owner)', () => {
    const baris = SUMBER.slice(SUMBER.indexOf('hasil.baris.map((b) =>'), SUMBER.indexOf('</tbody>'))
    const sel = [...baris.matchAll(/<td \/>|<td[^>]*>([\s\S]*?)<\/td>/g)].map((m) => (m[1] ?? '').trim())
    expect(sel).toHaveLength(8)
    expect(sel[0]).toBe('{b.caseId}')
    expect(sel[1]).toContain('onBuka(b.caseId)')
    expect(sel[1]).toContain('e.stopPropagation()')
    expect(SUMBER).toContain('<tr key={b.caseId} className="inbox__baris" onClick={() => onBuka(b.caseId)}>')
    expect(sel[2]).toBe('')
    expect(sel.slice(3, 6)).toEqual(['{b.groupBusiness}', '{b.insuredName}', '{b.marketing}'])
    expect(sel[6]).toBe('')
    expect(sel[7]).toBe('{b.status}')
  })

  it('jawaban lama dibuang; kosong = pesan, bukan tabel kosong tanpa keterangan; galat apa adanya', () => {
    expect(SUMBER).toContain('if (nomor === nomorPermintaan.current) setHasil(h)')
    expect(SUMBER).toContain('<Kosong pesan={TEKS_PORTAL.tanpaCase} />')
    expect(SUMBER).toContain('<Gagal galat={galat} />')
  })

  it('tombol Create opportunity hidup (membuka form, tiket 26)', () => {
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn btn--primary">${KEPALA_PORTAL.buat.label}</button>`))
  })

  it('templat = menu Kelola User: kepala inbox, toolbar (saring kiri, Create opportunity kanan), tabel inbox', () => {
    expect(HTML).toContain('<section class="inbox"><header class="inbox__kepala"><h2 class="inbox__judul">')
    const bilah = HTML.slice(HTML.indexOf('<form class="toolbar"'), HTML.indexOf('</form>'))
    expect(bilah.indexOf('type="search"')).toBeLessThan(bilah.indexOf('toolbar__spacer'))
    expect(bilah.indexOf('toolbar__spacer')).toBeLessThan(bilah.indexOf(`>${KEPALA_PORTAL.buat.label}</button>`))
    expect(HTML).toContain('<table class="inbox__tabel">')
  })
})
