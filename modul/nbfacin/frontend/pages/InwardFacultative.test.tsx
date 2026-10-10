// Paritas layar Inward Facultative tahap 1 (tiket 30) - dirender react-dom/server, data uji sintetis.

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { PERIODE as P, PILIHAN_PERIODE, SHOW_DETAIL, TEKS_INWARD, TOMBOL_KAKI_INWARD as KAKI } from '../labels'
import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import InwardFacultative, { endSebelumBegin, hariIniKabel, kodePolicyType, tambahSatuTahun, type KasusBaru } from './InwardFacultative'

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
    // Type facultative baca-saja (work owner 05-10-2026: "tidak bisa dipilih lgi, krn awalnya sudah dipilih").
    expect(HTML).toContain('<div class="nbf-inward__teks">Facultative In</div>')
    expect(HTML).not.toContain('<option value="Facultative In"')
    expect(HTML).toContain(`value="${hariIniKabel().replace(/-/g, '/')}"`)
  })

  it('radio Policy Type dan Day dari tangkapan layar, belum terpilih', () => {
    for (const p of PILIHAN_PERIODE.policyType) expect(HTML).toContain(`value="${p.value}"/> ${p.label}`)
    for (const p of PILIHAN_PERIODE.day) expect(HTML).toContain(`value="${p}"/> ${p}`)
    expect(HTML).not.toMatch(/type="radio"[^>]*checked/)
  })

  it('mode hanya lihat (View Details / View Policy, work owner 09-10-2026): Show Detail terbuka, tanpa kaki Cancel / Save / Submit', () => {
    const LIHAT = renderToStaticMarkup(<InwardFacultative kasus={KASUS} onBatal={() => {}} hanyaLihat />)
    expect(LIHAT).not.toContain(`>${KAKI.simpan.label}</button>`)
    expect(LIHAT).not.toContain(`>${KAKI.submit.label}</button>`)
    expect(LIHAT).toMatch(/<input type="checkbox" checked=""\/> Show Detail/)
    expect(HTML).toContain(`>${KAKI.submit.label}</button>`)
  })

  it('Class of Business selain FIRE (mis. Liability): tab berpesan "belum tersedia", bukan galat pemuatan (work owner 09-10-2026)', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toContain("dariServer !== null && !kasusFire ? (")
    expect(TEKS_INWARD.cobBelum('Object', 'LIABILITY')).toBe('Tab Object untuk Class of Business LIABILITY belum tersedia di sistem baru (saat ini baru FIRE).')
  })

  it('Policy Type tersimpan standard value DDL PolicyType.xml (work owner 09-10-2026); label lama dipetakan ke kode', () => {
    expect(PILIHAN_PERIODE.policyType.map((p) => `${p.value}:${p.label}`)).toEqual(['0:Individual Policy', '1:Master Policy'])
    expect(kodePolicyType('Master Policy')).toBe('1')
    expect(kodePolicyType('Individual Policy')).toBe('0')
    expect(kodePolicyType('1')).toBe('1')
    expect(kodePolicyType('')).toBe('')
    const DDL = 'D:\\migrasi\\RNM\\DDL\\PolicyType.xml'
    if (existsSync(DDL)) {
      const xml = readFileSync(DDL, 'utf-8')
      for (const p of PILIHAN_PERIODE.policyType)
        expect(xml).toMatch(new RegExp(`<pyLocalizedValue>${p.label}</pyLocalizedValue>\\s*<pyStandardValue>${p.value}</pyStandardValue>`))
    }
  })

  it('tombol fitur lain nonaktif; Cancel dan Save for later hidup, Submit nonaktif', () => {
    // Bilah CSV objek kini aktif (tiket 66) - lihat BilahCsvObjek.test.tsx.
    for (const t of [P.search]) {
      expect(HTML).toMatch(new RegExp(`<button[^>]*disabled=""[^>]*>${t.label.replace(/[/.]/g, '\\$&')}</button>`))
    }
    // Upload AI hanya akun IT (work owner 08-10-2026) - render tanpa data server = bukan IT.
    for (const t of [P.uploadQuotation, P.uploadRISlip, P.insertAccumulation]) expect(HTML).not.toContain(`>${t.label}</button>`)
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn nbf-tema__tombol btn--ghost">${KAKI.batal.label}</button>`))
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn nbf-tema__tombol btn--ghost">${KAKI.simpan.label}</button>`))
    expect(HTML).toMatch(new RegExp(`<button[^>]*disabled=""[^>]*>${KAKI.submit.label}</button>`))
  })

  it('blok SUMMARY dihapus (work owner 03-10-2026); Show Detail belum dicentang -> tab tidak tampil', () => {
    expect(HTML).not.toContain('SUMMARY')
    expect(HTML).not.toContain('Object Name')
    expect(HTML).not.toContain(TEKS_INWARD.kosong)
    expect(HTML).toContain(`/> ${SHOW_DETAIL}`)
    expect(HTML).not.toContain('role="tablist"')
  })

  it('Risk Scoring = skor tersimpan + kelas berwarna seperti tab Scoring; Save Scoring memperbarui General (tiket 54)', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toContain('setSkorRisiko(k.riskScoring)')
    expect(SUMBER).not.toContain('<span className="nbf-score__skor" data-nada={nadaKelas(skorRisiko.noteFinalScore)}>')
    expect(SUMBER).toContain('<span className="nbf-score__pil nbf-score__pil--umum" data-nada={nadaKelas(skorRisiko.noteFinalScore)}>')
    expect(SUMBER).toContain('<TabScoring caseId={kasus.caseId} saatTersimpan={setSkorRisiko} />')
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
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn nbf-tema__tombol btn--sm">${P.changeSob.label}</button>`))
    expect(HTML).not.toContain('role="dialog"')
    expect(SUMBER).toContain('onClick={() => setPilihSob(true)}')
    expect(SUMBER).toMatch(/onPilih=\{\(b\) => \{\s*setSob\(\{ id: b\.id, nama: b\.name \}\)/)
    expect(SUMBER).toContain('<Tampil label={P.sourceOfBusiness.label} nilai={sob.nama} />')
    expect(SUMBER).toContain('sourceOfBusinessId: sob.id,')
    expect(SUMBER).toContain('setSob({ id: g.sourceOfBusinessId, nama: g.sourceOfBusiness })')
  })

  it('Change Ceding Co hidup membuka popup daftar; Submit mengganti daftar; nama digabung `;`; kode ikut Save for later (tiket 34)', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(HTML).toMatch(new RegExp(`<button type="button" class="btn nbf-tema__tombol btn--sm">${P.changeCedingCo.label}</button>`))
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
    expect(SUMBER).toMatch(/label=\{P\.endDate\.label\}\s*value=\{selesai\}\s*disabled=\{kunciPeriode\}\s*onChange=\{setSelesai\}/)
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
    expect(SUMBER).toContain('disabled={menyimpan || tanggalSalah || bacaSaja}')
  })

  it('tahap Binding (tiket 64 K1): Begin / End date terkunci bila kunciPeriode', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toContain('const kunciPeriode = dariServer?.kunciPeriode === true')
    expect(SUMBER).toMatch(/label=\{P\.beginDate\.label\}\s*value=\{mulai\}\s*disabled=\{kunciPeriode\}/)
    expect(SUMBER).toMatch(/label=\{P\.endDate\.label\}\s*value=\{selesai\}\s*disabled=\{kunciPeriode\}/)
  })

  it('tahap Underwriting (tiket 59): baca-saja kecuali tab Spreading; grid baris hanya bila ubahBarisSpreading', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toContain("const spreadingTerbuka = !hanyaLihat && bacaSaja && dariServer?.bolehUbahSpreading === true")
    expect(SUMBER).toContain("const detailBacaSaja = bacaSaja && !(spreadingTerbuka && tab === 'Spreading')")
    expect(SUMBER).toContain("<BacaSajaKasus key={`${detailBacaSaja ? `bacasaja${kecualiDetail ?? ''}` : 'ubah'}-${versiDetail}`} aktif={detailBacaSaja} kecuali={kecualiDetail} sembunyi={hanyaLihat}>")
    expect(SUMBER).toContain('<TabSpreading caseId={kasus.caseId} tahapUw={dariServer?.ubahBarisSpreading === true} />')
    // General + Save for later tetap baca-saja di tahap UW.
    expect(SUMBER).toContain('<BacaSajaKasus aktif={bacaSaja} sembunyi={hanyaLihat}>')
    expect(SUMBER).toContain('disabled={menyimpan || tanggalSalah || bacaSaja}')
  })

  it('bilah CSV: empat tombol CSV tidak tampil saat baca-saja (UW dst.); Insert Accumulation hanya IT (work owner 08-10-2026)', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toContain('{!hanyaLihat && (!bacaSaja || akunIT) && (')
    // Tahap ubah: bilah CSV objek (tiket 66) + Insert Accumulation hanya IT; baca-saja: hanya Insert Accumulation (bila IT).
    expect(SUMBER).toMatch(/\{!bacaSaja \? \(\s*<BilahCsvObjek\s+caseId=\{kasus\.caseId\}\s+akhir=\{akunIT && <TombolNonaktif label=\{P\.insertAccumulation\.label\} \/>\}/)
    expect(SUMBER).toMatch(/\) : \(\s*<div className="nbf-opp__tombol">\s*<TombolNonaktif label=\{P\.insertAccumulation\.label\} \/>/)
  })

  it('tombol Upload QUOTATION/PLACING SLIP AI dan Upload R/I SLIP AI hanya untuk akun IT (work owner 08-10-2026)', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toContain('const akunIT = dariServer?.aksesLossRecord === true')
    expect(SUMBER).toMatch(/\{akunIT && \(\s*<>\s*<div className="nbf-opp__tombol">\s*<TombolNonaktif label=\{P\.uploadQuotation\.label\} \/>[\s\S]{0,160}P\.uploadRISlip\.label/)
  })

  it('tahap Underwriting (tiket 60): Accumulation di tab Coverage terbuka lewat pengecualian data-ubah-akumulasi', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(SUMBER).toContain('const akumulasiTerbuka = !hanyaLihat && bacaSaja && dariServer?.ubahAkumulasi === true')
    expect(SUMBER).toContain("const kecualiDetail = akumulasiTerbuka && tab === 'Coverage' ? PILIH_UBAH_AKUMULASI : undefined")
    const FORM = readFileSync(join(__dirname, '..', 'components', 'FormCoverage.tsx'), 'utf8').replace(/\r\n/g, '\n')
    // Choose + Copy Accumulation + popup bertanda; field coverage lain TIDAK (tetap baca-saja di UW - CoverageItem_IsUW Read-only).
    expect(FORM.match(/data-ubah-akumulasi=""/g)).toHaveLength(3)
    const TAB = readFileSync(join(__dirname, '..', 'components', 'TabCoverage.tsx'), 'utf8').replace(/\r\n/g, '\n')
    expect(TAB).toContain("data-ubah-akumulasi={hanyaAkumulasi ? '' : undefined}")
  })

  it('tab Object kasus FIRE = TabObject; FIRE = Group Business "FIRE" (perilaku lama, tak boleh hilang) ATAU lini backend (tiket 73 P0), Clauses / Scoring Risk hanya FIRE (G-5)', () => {
    const SUMBER = readFileSync(join(__dirname, 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')
    // Regresi 09-10-2026: case FIRE tampil "belum tersedia" karena layar FIRE hanya menurut `lini.fire` backend.
    expect(SUMBER).toContain("const kasusFire = (op.groupBusiness ?? '').trim().toUpperCase() === 'FIRE' || dariServer?.lini?.fire === true")
    expect(SUMBER).not.toContain('dariServer !== null ? dariServer.lini?.fire === true')
    expect(SUMBER).toContain("(t !== 'Clauses' || kasusFire) &&")
    expect(SUMBER).toContain("(t !== 'Scoring Risk' || kasusFire)")
    expect(SUMBER).toMatch(/\{tab === 'Object' && kasusFire \? \(\s*<TabObject caseId=\{kasus\.caseId\} insuredName=\{insured\} \/>/)
    // Tab Coverage FIRE (tiket 43).
    expect(SUMBER).toMatch(/: tab === 'Coverage' && kasusFire \? \(\s*<TabCoverage caseId=\{kasus\.caseId\} hanyaAkumulasi=\{akumulasiTerbuka\} \/>/)
  })
})
