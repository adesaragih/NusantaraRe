// Paritas layar Inward Facultative tahap 1 (tiket 30) - dirender react-dom/server, data uji sintetis.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { PERIODE as P, PILIHAN_PERIODE, SHOW_DETAIL, TEKS_INWARD, TOMBOL_KAKI_INWARD as KAKI } from '../labels'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import InwardFacultative, { endSebelumBegin, hariIniKabel, tambahSatuTahun, type KasusBaru } from './InwardFacultative'

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
    // Templat Kelola User: nomor case = judul `h2.inbox__judul` di kepala halaman.
    expect(HTML.indexOf('<h2 class="inbox__judul">NB-1</h2>')).toBeGreaterThan(-1)
    expect(HTML.indexOf('<h2 class="inbox__judul">NB-1</h2>')).toBeLessThan(HTML.indexOf(`>${P.judul.label}</h4>`))
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

  it('tombol fitur lain nonaktif; Cancel dan Save for later hidup, Submit nonaktif', () => {
    for (const t of [P.uploadQuotation, P.uploadRISlip, P.search, P.downloadTemplateCsv, P.uploadCsv, P.viewUpload, P.saveData, P.insertAccumulation]) {
      expect(HTML).toMatch(new RegExp(`<button[^>]*disabled=""[^>]*>${t.label.replace(/[/.]/g, '\\$&')}</button>`))
    }
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn btn--ghost">${KAKI.batal.label}</button>`))
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn btn--ghost">${KAKI.simpan.label}</button>`))
    expect(HTML).toMatch(new RegExp(`<button[^>]*disabled=""[^>]*>${KAKI.submit.label}</button>`))
  })

  it('blok SUMMARY dihapus (work owner 03-10-2026); Show Detail belum dicentang -> tab tidak tampil', () => {
    expect(HTML).not.toContain('SUMMARY')
    expect(HTML).not.toContain('Object Name')
    expect(HTML).not.toContain(TEKS_INWARD.kosong)
    expect(HTML).toContain(`/> ${SHOW_DETAIL}`)
    expect(HTML).not.toContain('role="tablist"')
  })

  it('hariIniKabel = DD-MM-YYYY', () => {
    expect(hariIniKabel(new Date(2026, 9, 3))).toBe('03-10-2026')
  })

  it('tahap 2: case dimuat saat dibuka, Marketing Name dari backend, Save for later = PUT general (tiket 31)', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toContain('ambilKasus(kasus.caseId).then(')
    expect(SUMBER).toContain('}, [kasus.caseId])')
    expect(SUMBER).toContain('setOpsiMarketing(h.baris.map((b) => ({ value: b.id, label: b.nama })))')
    expect(SUMBER).toContain('opsi={opsiMarketing}')
    const awal = SUMBER.indexOf('async function simpan()')
    const simpan = SUMBER.slice(awal, SUMBER.indexOf('\n  }\n', awal))
    expect(simpan).toMatch(/simpanGeneral\(kasus\.caseId, \{[\s\S]*marketingId: marketing,[\s\S]*\}\)/)
    expect(simpan).not.toMatch(/sourceOfBusiness:|cedingCoName|cedingList|groupName|oldPolicyNumber/)
    // Offering date tetap hari ini bila server belum punya nilainya (InwardFacultative_PreDT).
    expect(SUMBER).toContain('setPenawaran(g.offeringDate || hariIniKabel())')
  })

  it('Change SOB hidup membuka popup; Choose mengisi kode + nama; kode ikut Save for later (tiket 33, E-1)', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn btn--sm">${P.changeSob.label}</button>`))
    expect(HTML).not.toContain('role="dialog"')
    expect(SUMBER).toContain('onClick={() => setPilihSob(true)}')
    expect(SUMBER).toMatch(/onPilih=\{\(b\) => \{\s*setSob\(\{ id: b\.id, nama: b\.name \}\)/)
    expect(SUMBER).toContain('<Tampil label={P.sourceOfBusiness.label} nilai={sob.nama} />')
    expect(SUMBER).toContain('sourceOfBusinessId: sob.id,')
    expect(SUMBER).toContain('setSob({ id: g.sourceOfBusinessId, nama: g.sourceOfBusiness })')
  })

  it('Change Ceding Co hidup membuka popup daftar; Submit mengganti daftar; nama digabung `;`; kode ikut Save for later (tiket 34)', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn btn--sm">${P.changeCedingCo.label}</button>`))
    expect(SUMBER).toContain('onClick={() => setUbahCeding(true)}')
    expect(SUMBER).toMatch(/<PopupCedingCoList\s*awal=\{ceding\}/)
    expect(SUMBER).toMatch(/onSubmit=\{\(d\) => \{\s*setCeding\(d\)/)
    expect(SUMBER).toContain("nilai={ceding.map((c) => c.name).join(';')}")
    expect(SUMBER).toContain('cedingIds: ceding.map((c) => c.id),')
    expect(SUMBER).toContain('setCeding(g.cedingList ?? [])')
  })

  it('End date = Begin date + 1 tahun (29-02 -> 28-02); masukan tidak sah -> kosong', () => {
    expect(tambahSatuTahun('15-03-2026')).toBe('15-03-2027')
    expect(tambahSatuTahun('31-12-2026')).toBe('31-12-2027')
    expect(tambahSatuTahun('29-02-2028')).toBe('28-02-2029')
    expect(tambahSatuTahun('28-02-2027')).toBe('28-02-2028')
    expect(tambahSatuTahun('')).toBe('')
    expect(tambahSatuTahun('2026-03-15')).toBe('')
  })

  it('Begin date sah mengisi End date; ketikan setengah jadi tidak; End date tetap dapat diubah sendiri', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toMatch(/setMulai\(k\)[\s\S]{0,160}if \(k !== ''\) setSelesai\(tambahSatuTahun\(k\)\)/)
    expect(SUMBER).toMatch(/label=\{P\.endDate\.label\}\s*value=\{selesai\}\s*onChange=\{setSelesai\}/)
  })

  it('End date < Begin date = peringatan di End date dan Save for later tertahan; sama = sah (F-1)', () => {
    expect(endSebelumBegin('15-03-2026', '14-03-2026')).toBe(true)
    expect(endSebelumBegin('15-03-2026', '01-01-2026')).toBe(true)
    expect(endSebelumBegin('31-12-2026', '01-01-2027')).toBe(false)
    expect(endSebelumBegin('15-03-2026', '15-03-2026')).toBe(false)
    expect(endSebelumBegin('', '14-03-2026')).toBe(false)
    expect(endSebelumBegin('15-03-2026', '')).toBe(false)
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toContain('error={tanggalSalah ? TEKS_INWARD.endSebelumBegin : undefined}')
    expect(SUMBER).toMatch(/async function simpan\(\) \{[\s\S]{0,120}if \(tanggalSalah\) return/)
    expect(SUMBER).toContain('disabled={menyimpan || tanggalSalah}')
  })

  it('tab Object kasus FIRE = TabObject; selain FIRE tetap BelumTersedia (G-5)', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toContain("const kasusFire = (op.groupBusiness ?? '').trim().toUpperCase() === 'FIRE'")
    expect(SUMBER).toContain("{tab === 'Object' && kasusFire ? <TabObject caseId={kasus.caseId} insuredName={insured} /> : <BelumTersedia")
  })
})
