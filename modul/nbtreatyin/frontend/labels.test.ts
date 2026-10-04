// Audit silang putaran 3 (bab 3.2 butir 4): teks yang TAMPIL di layar sama
// dengan LABEL / pyLabel XML (korpus `Section/*`, `FlowAction/*`, dibaca ulang
// 04-10-2026). Harapan diketik dari XML, bukan dari labels.ts.

import { describe, expect, it } from 'vitest'

import {
  BAGIAN,
  JUDUL,
  KOLOM_ANGSURAN,
  KOLOM_BISNIS,
  KOLOM_PORTAL,
  KOLOM_SPREADING,
  KOLOM_USULAN,
  KONFIRMASI_TOLAK,
  PORTAL,
} from './labels'
import { KOLOM_RINCI } from './nonprop'

describe('label layar = LABEL XML', () => {
  it('grid spreading: SpreadingRiskList S2 / DetailPolicyTreatyIn / DetailDeptHeadTreatyIn_UW', () => {
    expect([
      KOLOM_SPREADING.treatyType,
      KOLOM_SPREADING.share,
      KOLOM_SPREADING.premium,
      KOLOM_SPREADING.claimPct,
      KOLOM_SPREADING.claim,
      KOLOM_SPREADING.totalShare,
    ]).toEqual(['Type Treaty', '% Share', 'Premium', '% Share', 'Claim', 'Total'])
  })

  it('grid angsuran ListInstallment: judul kolom dan judul wadah (pyIncludeHeader, pyTitle)', () => {
    expect([KOLOM_ANGSURAN.no, KOLOM_ANGSURAN.dueDate, KOLOM_ANGSURAN.pct, KOLOM_ANGSURAN.premium, KOLOM_ANGSURAN.total]).toEqual([
      'No',
      'Due Date',
      '% Installment',
      'Premium Nusantara Re',
      'Total Payment',
    ])
    expect(BAGIAN.angsuran).toBe('Installment Data Information')
  })

  it('ListSuggest: Date, PIC, Approval, Suggest', () => {
    expect([KOLOM_USULAN.tanggal, KOLOM_USULAN.operator, KOLOM_USULAN.putusan, KOLOM_USULAN.catatan]).toEqual([
      'Date',
      'PIC',
      'Approval',
      'Suggest',
    ])
  })

  it('modal tolak: FlowAction PolicyTreatyInDeclineConfirm pyLabel + LABEL section', () => {
    expect(JUDUL.tolak).toBe('Confirm Decline NB')
    expect(KONFIRMASI_TOLAK).toBe('Are you sure you want to decline this NB')
  })

  it('modal nomor polis: FlowAction ShowPolicyNoTreaty pyLabel', () => {
    expect(JUDUL.nomorPolis).toBe('Show PolicyNo')
  })

  it('rincian angsuran NonProp (Section InstallmentList): kolom ke-3 tanpa judul', () => {
    expect(KOLOM_RINCI.map((k) => k.label)).toContain('')
    expect(KOLOM_RINCI.map((k) => k.label)).not.toContain('Currency')
  })

  it('portal SFAPortal_OpportunitiesList / Header: judul kolom, judul halaman, placeholder saring', () => {
    expect([KOLOM_PORTAL.id, KOLOM_PORTAL.bisnis, KOLOM_PORTAL.tertanggung, KOLOM_PORTAL.marketing, KOLOM_PORTAL.status]).toEqual([
      'Offer No',
      'Group Business',
      'Insured Name',
      'Marketing',
      'Status',
    ])
    expect(PORTAL.judul).toBe('Opportunity')
    expect(PORTAL.placeholder).toBe('NB-1234 or Name')
  })

  it('popup pilih bisnis: Section BusinessAndSOBList grid aktif (BrowseTreatyJoinEDM) baris 1 LABEL / baris 2 sel, pyWindowName', () => {
    // judul baris 1 sel 322-346 dan properti baris 2 sel 349-373, sel demi sel (sel 321/348 = tombol Choose)
    expect(KOLOM_BISNIS.map((k) => [k.judul, k.kolom])).toEqual([
      ['Treaty Offer ID', 'TREATYID'],
      ['Contract Name', 'TREATYCONTRACTNAME'],
      ['Class of Business', 'CLASSOFBUSINESS'],
      ['Source Of Business', 'SOB'],
      ['Insured Name', 'CEDING'],
      ['Proportion Type', 'PROPORTIONTYPE'],
      ['Treaty Type', 'TREATYTYPE'],
      ['Treaty Group', 'TREATYGROUP'],
      ['Treaty Year', 'TREATYYEAR'],
      ['Currency', 'LIMITCURRENCY'],
      ['Limit', 'LIMITVALUE'],
      ['Currency', 'RETENTIONCURRENCY'],
      ['Retention', 'RETENTIONVALUE'],
      ['Currency', 'EPICURRENCY'],
      ['EPI', 'EPIVALUE'],
      ['Layer', 'LAYERTYPE'],
      ['', 'LAYER'],
      ['Part of', 'LAYERPARTTYPE'],
      ['', 'LAYERPART'],
      ['Currency', 'MDPCURRENCY'],
      ['MDP', 'MDPVALUE'],
      ['Currency', 'NETPREMICURRENCY'],
      ['Net Premium', 'NETPREMIVALUE'],
      ['Currency', 'SHARECURRENCY'],
      ['Share RNM Value', 'SHAREVALUE'],
    ])
    // showHarness tombol Choose Business (DetailPolicyTreatyIn) pyWindowName
    expect(JUDUL.pilihBisnis).toBe('Business And SOB List')
  })
})
