// Penjaga bukti label NB FacIn - tiket 21; portal Opportunity tiket 25.
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN. Berkas ini membuka
// section korpus dan memastikan setiap sel yang disebut `labels.ts` memang berlabel
// dan berproperti seperti tertulis. Korpus READ-ONLY - hanya dibaca.
//
// ⚠️ Bila korpus tidak terjangkau, test DILEWATI dengan pesan - bukan gagal.

import { existsSync, readFileSync } from 'node:fs'

import { describe, expect, it } from 'vitest'

import {
  KEPALA_PORTAL,
  KOLOM_PORTAL,
  KOLOM_RINGKASAN,
  MEDAN_COVERAGE_CARGO,
  PERIODE,
  SARING_PORTAL,
  SHOW_DETAIL,
  TAB_DETAIL,
  TOMBOL_COVERAGE_CARGO,
  TOMBOL_KAKI_INWARD,
} from './labels'

const SECTION = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\InputCoverageCargo_FacIn.xml'
const ada = existsSync(SECTION)
const KEPALA = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\SFAPortalOpportunitiesHeader.xml'
const DAFTAR = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\SFAPortal_OpportunitiesList.xml'
const adaPortal = existsSync(KEPALA) && existsSync(DAFTAR)

/**
 * Isi elemen sel ber-`pyCellId` tertentu: `rowdata` terdalam yang terbuka saat
 * `<pyCellId>` itu ditemukan, sampai penutup PASANGANNYA (rowdata bersarang - mis.
 * aksi tombol - ikut di dalamnya). Pemindai kedalaman, bukan regex tak-rakus.
 */
function blokSel(xml: string, sel: string): string[] {
  const hasil: string[] = []
  const pola = /<rowdata\b[^>]*>|<\/rowdata>|<pyCellId>(\d+)<\/pyCellId>/g
  const tumpukan: { mulai: number; cocok: boolean }[] = []
  for (const m of xml.matchAll(pola)) {
    const i = m.index ?? 0
    if (m[0].startsWith('<rowdata')) {
      tumpukan.push({ mulai: i, cocok: false })
    } else if (m[0] === '</rowdata>') {
      const atas = tumpukan.pop()
      if (atas?.cocok) hasil.push(xml.slice(atas.mulai, i + m[0].length))
    } else if (m[1] === sel) {
      const atas = tumpukan.at(-1)
      if (atas) atas.cocok = true
    }
  }
  return hasil
}

describe.skipIf(!ada)('label NB FacIn = korpus InputCoverageCargo_FacIn', () => {
  const xml = ada ? readFileSync(SECTION, 'utf-8') : ''

  it.each(Object.entries(MEDAN_COVERAGE_CARGO))('medan %s', (_, m) => {
    const blok = blokSel(xml, m.sel)
    expect(blok.some((b) => b.includes(`<pyLabelFor>${m.label}</pyLabelFor>`) && b.includes(`<pyValue>${m.properti}</pyValue>`))).toBe(true)
  })

  it.each(Object.entries(TOMBOL_COVERAGE_CARGO))('tombol %s', (_, t) => {
    expect(blokSel(xml, t.sel).some((b) => b.includes(`<pyLabel>${t.label}</pyLabel>`))).toBe(true)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(xml, '11').some((b) => b.includes('<pyLabelFor>Rate (‰)</pyLabelFor>'))).toBe(false)
  })
})

if (!ada) {
  it('korpus tidak terjangkau - bukti label dilewati', () => {
    console.warn(`korpus tidak ada di ${SECTION}; uji bukti label NB FacIn dilewati`)
  })
}

describe.skipIf(!adaPortal)('label portal Opportunity = korpus SFAPortalOpportunitiesHeader + SFAPortal_OpportunitiesList', () => {
  const kepala = adaPortal ? readFileSync(KEPALA, 'utf-8') : ''
  const daftar = adaPortal ? readFileSync(DAFTAR, 'utf-8') : ''
  /** Sel itu SATU di berkasnya, dan teksnya ada di tag yang disebut. */
  const ditemukan = (xml: string, u: { sel: string; tag: string; label: string }) => {
    const blok = blokSel(xml, u.sel)
    return blok.length === 1 && blok[0]!.includes(`<${u.tag}>${u.label}</${u.tag}>`)
  }

  it.each(Object.entries(KEPALA_PORTAL))('kepala %s', (_, u) => {
    expect(ditemukan(kepala, u)).toBe(true)
  })

  it.each(Object.entries(SARING_PORTAL))('saring %s', (_, u) => {
    expect(ditemukan(daftar, u)).toBe(true)
  })

  it.each(KOLOM_PORTAL.map((k) => [k.sel, k] as const))('judul kolom sel %s', (_, k) => {
    const blok = blokSel(daftar, k.sel)
    expect(blok).toHaveLength(1)
    if (k.label === '') {
      // Kolom berjudul kosong: sel judulnya tanpa `pyValue` berisi.
      expect(blok[0]).not.toMatch(/<pyValue>[^<]+<\/pyValue>/)
    } else {
      expect(blok[0]).toContain(`<pyValue>${k.label}</pyValue>`)
    }
  })

  it('kolom berurutan = sel 89..96 (grid GetListOpportunityF)', () => {
    expect(KOLOM_PORTAL.map((k) => k.sel)).toEqual(['89', '90', '91', '92', '93', '94', '95', '96'])
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(ditemukan(kepala, { sel: '72', tag: 'pyLabel', label: 'Create Opportunity' })).toBe(false)
    expect(ditemukan(daftar, { sel: '89', tag: 'pyValue', label: 'Offer No.' })).toBe(false)
  })
})

if (!adaPortal) {
  it('korpus portal tidak terjangkau - bukti label portal dilewati', () => {
    console.warn(`korpus tidak ada di ${KEPALA} / ${DAFTAR}; uji bukti label portal dilewati`)
  })
}

const NBFACIN = 'D:\\migrasi\\RNM\\NB FacIn\\'
const BERKAS_INWARD = {
  periode: NBFACIN + 'Section\\Periode.xml',
  layar: NBFACIN + 'Section\\InputInwardFacultative.xml',
  dtl: NBFACIN + 'Section\\InputInwardFacultativeDtl.xml',
  ringkasan: NBFACIN + 'Section\\FireSummarySection.xml',
  aksi: NBFACIN + 'FlowAction\\InwardFacultative.xml',
}
const adaInward = Object.values(BERKAS_INWARD).every((b) => existsSync(b))

describe.skipIf(!adaInward)('label Inward Facultative = korpus (Periode, InputInwardFacultative, Dtl, FireSummarySection, flow action)', () => {
  const baca = (b: string) => (adaInward ? readFileSync(b, 'utf-8') : '')
  const periode = baca(BERKAS_INWARD.periode)

  it.each(Object.entries(PERIODE))('Periode %s', (_, u) => {
    if (u.sel === '') {
      // Judul layout / label tanpa sel: cukup ada di berkas dengan tag itu.
      expect(periode).toContain(`<${u.tag}>${u.label}</${u.tag}>`)
    } else {
      const blok = blokSel(periode, u.sel)
      expect(blok.some((b) => b.includes(`<${u.tag}>${u.label}</${u.tag}>`))).toBe(true)
    }
  })

  it('Show Detail = caption checkbox InputInwardFacultative', () => {
    expect(baca(BERKAS_INWARD.layar)).toContain(`<pyCheckboxCaption>${SHOW_DETAIL}</pyCheckboxCaption>`)
  })

  it.each(TAB_DETAIL.map((t) => [t]))('judul tab %s = pyTitle InputInwardFacultativeDtl', (t) => {
    expect(baca(BERKAS_INWARD.dtl)).toContain(`<pyTitle>${t}</pyTitle>`)
  })

  it.each(KOLOM_RINGKASAN.map((k) => [k]))('kolom ringkasan %s = FireSummarySection', (k) => {
    expect(baca(BERKAS_INWARD.ringkasan)).toContain(`<pyValue>${k}</pyValue>`)
  })

  it.each(Object.entries(TOMBOL_KAKI_INWARD))('tombol kaki %s = flow action InwardFacultative', (_, u) => {
    expect(baca(BERKAS_INWARD.aksi)).toContain(`<${u.tag}>${u.label}</${u.tag}>`)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(periode, '9').some((b) => b.includes('<pyLabelFieldValue>Reff. Number</pyLabelFieldValue>'))).toBe(false)
    expect(baca(BERKAS_INWARD.dtl)).not.toContain('<pyTitle>Objects</pyTitle>')
  })
})
