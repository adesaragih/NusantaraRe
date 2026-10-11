// Paritas layar portal Opportunity dengan harness Pega `SFAPortalOpportunities` - tiket 25; daftar case NB tiket 32.
// Tata letak = halaman depan NB Treaty In (work owner 11-10-2026). Halaman DIRENDER (react-dom/server).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { KEPALA_PORTAL, KOLOM_PORTAL, KOLOM_PORTAL_SELESAI, SARING_PORTAL, STATUS_PORTAL } from '../labels'
import PortalOpportunity from './PortalOpportunity'

const HTML = renderToStaticMarkup(<PortalOpportunity onBuat={() => {}} onBuka={() => {}} />)
const HTML_SELESAI = renderToStaticMarkup(<PortalOpportunity onBuat={() => {}} onBuka={() => {}} status="Resolved" onStatus={() => {}} />)
const SUMBER = readFileSync(join(__dirname, 'PortalOpportunity.tsx'), 'utf8').replace(/\r\n/g, '\n')
const RUTE = readFileSync(join(__dirname, '..', 'rute.tsx'), 'utf8').replace(/\r\n/g, '\n')

/** Teks setiap `<th>` berurutan. */
const judulKolom = (h: string) => [...h.matchAll(/<th[^>]*>([^<]*)<\/th>/g)].map((m) => m[1])

describe('PortalOpportunity = unsur yang TAMPIL di Pega', () => {
  it('akar halaman memasang kelas modul (gaya terisolasi)', () => {
    expect(HTML.startsWith('<div class="nbfacin">')).toBe(true)
  })

  it('kotak saring: placeholder korpus, label sebagai nama aksesibel (pyIncludeLabel=false)', () => {
    expect(HTML).toContain(`placeholder="${SARING_PORTAL.placeholder.label}"`)
    expect(HTML).toContain(`aria-label="${SARING_PORTAL.label.label}"`)
    expect(HTML).not.toContain(`>${SARING_PORTAL.label.label}<`)
  })

  it('judul kolom In Progress = grid GetListOpportunityF, berurutan, kolom 3 kosong', () => {
    expect(judulKolom(HTML)).toEqual(KOLOM_PORTAL.map((k) => k.label))
    expect(judulKolom(HTML)).toEqual(['Offer No', 'Name', '', 'Group Business', 'Insured Name', 'Marketing', 'NB Status', 'Status'])
  })

  it('tab Resolved: Policy Number sesudah Offer No, Production Date di akhir (pola NB Treaty In)', () => {
    expect(judulKolom(HTML_SELESAI)).toEqual([
      'Offer No', KOLOM_PORTAL_SELESAI.noPolis, 'Name', '', 'Group Business', 'Insured Name', 'Marketing', 'NB Status', 'Status',
      KOLOM_PORTAL_SELESAI.tglProduksi,
    ])
  })

  it('unsur tersembunyi permanen di Pega TIDAK dirender; Copy Old NB Treaty In tidak dibawa', () => {
    for (const t of ['Stage view', 'List view', '>All<', 'Individual', 'Corporate', 'Phase :', 'Proposal', 'Closed', 'Export', 'Refresh', 'Copy Old']) {
      expect(HTML).not.toContain(t)
    }
    expect(HTML).not.toContain('Create Opportunity')
  })

  it('daftar case NB dimuat saat dibuka dengan saringan switch; belum ada baris sebelum jawaban datang', () => {
    expect(SUMBER).toMatch(/useEffect\(\(\) => \{\s*void muat\('', 1\)/)
    expect(SUMBER).toContain('daftarCaseNB(cari, halaman, sudahSelesai)')
    expect(HTML).not.toContain('<tbody')
    expect(HTML).not.toContain('tidak dapat dimuat')
  })

  it('jawaban lama dibuang; kosong = pesan; galat apa adanya', () => {
    expect(SUMBER).toContain('if (nomor === nomorPermintaan.current) setHasil(h)')
    expect(SUMBER).toContain('<Kosong pesan={TEKS_PORTAL.tanpaCase} />')
    expect(SUMBER).toContain('<Gagal galat={galat} />')
  })
})

describe('PortalOpportunity = tata letak NB Treaty In (work owner 11-10-2026)', () => {
  it('kepala: judul kiri, Create opportunity di kanan kepala', () => {
    const kepala = HTML.slice(HTML.indexOf('<header'), HTML.indexOf('</header>'))
    expect(kepala).toContain(`<h2 class="inbox__judul nbf-tema__judul">${KEPALA_PORTAL.judul.label}</h2>`)
    expect(kepala).toContain(
      `<div class="nbf-portal__aksi"><button type="button" class="btn btn--primary nbf-tema__tombol nbf-tema__tombol--utama">${KEPALA_PORTAL.buat.label}</button></div>`,
    )
  })

  it('switch In Progress / Resolved sesudah kepala, bawaan In Progress; posisinya dipegang rute', () => {
    expect(STATUS_PORTAL).toEqual(['In Progress', 'Resolved'])
    const kepala = HTML.indexOf('</header>')
    for (const t of STATUS_PORTAL) expect(HTML.indexOf(`>${t}<`)).toBeGreaterThan(kepala)
    expect(SUMBER).toContain('void muat(kunciCari.current, 1, s === STATUS_PORTAL[1])')
    expect(RUTE).toContain('status={statusPortal}')
    expect(RUTE).toContain('onStatus={setStatusPortal}')
  })

  it('saringan SATU kapsul: ikon, isian, pengosong (tersembunyi saat kosong), Filter; Enter / Escape / ✕', () => {
    const kapsul = HTML.slice(HTML.indexOf('<form class="nbf-saring" role="search">'), HTML.indexOf('</form>'))
    expect(kapsul).toContain('<svg class="nbf-saring__ikon"')
    expect(kapsul).toContain('class="nbf-saring__isian"')
    expect(kapsul).toContain('class="nbf-saring__hapus" aria-label="Hapus isian saring" disabled=""')
    expect(kapsul).toContain(`<button type="submit" class="nbf-saring__tombol">${SARING_PORTAL.tombol.label}</button>`)
    expect(SUMBER).toContain('void muat(kotak, 1)')
    expect(SUMBER).toMatch(/if \(e\.key === 'Escape'\) \{\s*setKotak\(''\)\s*void muat\('', 1\)/)
    expect(SUMBER).toContain('onPindah={(h) => void muat(kunciCari.current, h)}')
  })

  it('baris: nomor case = tautan pembuka case, status = chip, baris dapat diklik', () => {
    const baris = SUMBER.slice(SUMBER.indexOf('hasil.baris.map((b) =>'), SUMBER.indexOf('</tbody>'))
    expect(baris).toContain('className="nbf-tautan nbf-portal__nomor"')
    expect(baris).toContain('e.stopPropagation()')
    expect(SUMBER).toContain('<tr key={b.caseId} className="inbox__baris" onClick={() => onBuka(b.caseId)}>')
    expect(baris).toContain('<span className="nbf-portal__status">{b.nbStatus}</span>')
    expect(baris).toContain('<span className="nbf-portal__status">{b.status}</span>')
    expect(baris).toContain("{selesai && <td>{b.noPolis ?? ''}</td>}")
    expect(baris).toContain("{selesai && <td>{b.tglProduksi ?? ''}</td>}")
  })
})
