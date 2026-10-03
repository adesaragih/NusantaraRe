// Paritas layar Inward Facultative tahap 1 (tiket 30) - dirender react-dom/server, data uji sintetis.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { KOLOM_RINGKASAN, PERIODE as P, PILIHAN_PERIODE, SHOW_DETAIL, TEKS_INWARD, TOMBOL_KAKI_INWARD as KAKI } from '../labels'
import InwardFacultative, { hariIniKabel, type KasusBaru } from './InwardFacultative'

const KASUS: KasusBaru = {
  caseId: 'NB-1',
  insuredName: 'UJI INSURED',
  isian: {
    estimatedClosingDate: '10-03-2026', businessProspectName: 'UJI', accountId: 'UJI-A', insuredId: 'UJI-I',
    groupBusinessId: 'UJI-G', groupBusiness: 'UJI GRUP', classOfBusiness: 'UJI COB', typeOfInward: 'Facultative',
    typeOfFacultative: 'Facultative In', phase: 'Proposal', stage: 'Opportunity', opportunitySource: '',
    businessStatus: 'New Business', description: '',
  },
}

const HTML = renderToStaticMarkup(<InwardFacultative kasus={KASUS} onBatal={() => {}} />)

/** Label medan berurutan (label/span field__label), dengan penanda wajib. */
const label = [...HTML.matchAll(/<(?:label|span) class="field__label"[^>]*>([^<]*)(<span class="field__req">\*<\/span>)?/g)].map((m) => ({
  teks: m[1],
  wajib: m[2] !== undefined,
}))

describe('Inward Facultative tahap 1 = tangkapan layar kasus FIRE + XML Periode', () => {
  it('judul = nomor case, lalu blok General', () => {
    expect(HTML.indexOf('>NB-1</h3>')).toBeGreaterThan(-1)
    expect(HTML.indexOf('>NB-1</h3>')).toBeLessThan(HTML.indexOf(`>${P.judul.label}</h4>`))
  })

  it('urutan label: kolom kiri lalu kanan', () => {
    expect(label.map((l) => l.teks)).toEqual([
      P.reffNumber.label, P.businessStatus.label, P.insuredName.label, P.qqName.label, P.beginDate.label,
      P.offeringDate.label, P.policyType.label, P.riskScoring.label,
      P.classOfBusiness.label, P.typeFacultative.label, P.sourceOfBusiness.label, P.cedingCoName.label,
      P.groupName.label, P.endDate.label, P.followingPolicyNumber.label, P.oldPolicyNumber.label,
      P.marketingName.label, P.day.label,
    ])
  })

  it('tanda wajib: Begin date, Offering date, End date, Marketing Name', () => {
    expect(label.filter((l) => l.wajib).map((l) => l.teks)).toEqual([P.beginDate.label, P.offeringDate.label, P.endDate.label, P.marketingName.label])
  })

  it('nilai awal dari Create opportunity; Class of business = Group Business; Offering date = hari ini', () => {
    expect(HTML).toContain('<div class="nbf-inward__teks">New Business</div>')
    expect(HTML).toContain('<div class="nbf-inward__teks">UJI INSURED</div>')
    expect(HTML).toContain('<div class="nbf-inward__teks">UJI GRUP</div>')
    expect(HTML).not.toContain('UJI COB')
    expect(HTML).toMatch(/<option value="Facultative In" selected="">Facultative In<\/option>/)
    expect(HTML).toContain(`value="${hariIniKabel().replace(/-/g, '/')}"`)
  })

  it('radio Policy Type dan Day dari tangkapan layar, belum terpilih', () => {
    for (const p of [...PILIHAN_PERIODE.policyType, ...PILIHAN_PERIODE.day]) expect(HTML).toContain(`value="${p}"/> ${p}`)
    expect(HTML).not.toMatch(/type="radio"[^>]*checked/)
  })

  it('tombol fitur lain nonaktif; Cancel hidup, Save for later dan Submit nonaktif', () => {
    for (const t of [P.uploadQuotation, P.uploadRISlip, P.changeSob, P.changeCedingCo, P.search, P.downloadTemplateCsv, P.uploadCsv, P.viewUpload, P.saveData, P.insertAccumulation]) {
      expect(HTML).toMatch(new RegExp(`<button[^>]*disabled=""[^>]*>${t.label.replace(/[/.]/g, '\\$&')}</button>`))
    }
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn btn--ghost">${KAKI.batal.label}</button>`))
    expect(HTML).toMatch(new RegExp(`<button[^>]*disabled=""[^>]*>${KAKI.simpan.label}</button>`))
    expect(HTML).toMatch(new RegExp(`<button[^>]*disabled=""[^>]*>${KAKI.submit.label}</button>`))
  })

  it('ringkasan: kolom Object Name / Location, kosong; Show Detail belum dicentang -> tab tidak tampil', () => {
    for (const k of KOLOM_RINGKASAN) expect(HTML).toContain(`>${k}</th>`)
    expect(HTML).toContain(TEKS_INWARD.kosong)
    expect(HTML).toContain(`/> ${SHOW_DETAIL}`)
    expect(HTML).not.toContain('role="tablist"')
  })

  it('hariIniKabel = DD-MM-YYYY', () => {
    expect(hariIniKabel(new Date(2026, 9, 3))).toBe('03-10-2026')
  })
})
