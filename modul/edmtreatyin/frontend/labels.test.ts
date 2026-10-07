// Teks yang TAMPIL di layar = LABEL / pyCaption / pyButtonLabel korpus `EDM Treaty In` (dibaca 06-10-2026). Harapan
// diketik dari XML, bukan dari labels.ts (pola `modul/nbtreatyin/frontend/labels.test.ts`).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import {
  LABEL_JENIS_EDM,
  labelJenisEDM,
  BAGIAN,
  BUAT,
  JUDUL,
  JUDUL_POSISI,
  KOLOM_ANGSURAN,
  KOLOM_BISNIS,
  KOLOM_PORTAL,
  KOLOM_SPREADING,
  KOLOM_USULAN,
  KONFIRMASI_TOLAK,
  NOMOR_DIAKSEP,
  PORTAL,
  TAB_DATA,
  TOMBOL,
} from './labels'
import { BARIS_PER_HALAMAN_GRID, BARIS_PER_HALAMAN_USULAN } from './paginasi'
import { KOLOM_RINCI_ANGSURAN, KOLOM_XOL, KOLOM_XOL_RINCI } from './xol'

describe('label layar = LABEL XML', () => {
  it('portal SFAPortal_Endorsement_Treaty: judul, tombol, saring, 10 kolom grid (kolom 1 dan 4 kembar)', () => {
    expect(JUDUL.portal).toBe('Addendum Treaty')
    expect([TOMBOL.create, TOMBOL.filter, TOMBOL.refresh, TOMBOL.refreshTooltip]).toEqual([
      'Create New Addendum Treaty',
      'Filter',
      'Refresh',
      'Refresh EDM Grid',
    ])
    expect([PORTAL.placeholder, PORTAL.filter]).toEqual([
      'Search EDM no, master ID, policy no, insured, business, ceding, marketing...',
      'Filter Term for Endorsement',
    ])
    expect(KOLOM_PORTAL.map((k) => k.judul)).toEqual([
      'EDM Number',
      'Offer No',
      'Policy Number',
      'EDM Number',
      'SOB',
      'Ceding',
      'EDM Type',
      'Proportional Type',
      'Marketing',
      'Status',
    ])
    expect(BARIS_PER_HALAMAN_GRID).toBe(50)
  })

  it('layar Create TreatyCreateEdm: judul harness, tiga medan, tombol', () => {
    expect(JUDUL.buat).toBe('Create EDM')
    expect([BUAT.noPolis, BUAT.noMaster, BUAT.jenis]).toEqual([
      'No Polis Treaty',
      'No Master Treaty',
      'Source of Change',
    ])
    expect(TOMBOL.buat).toBe('Create')
  })

  it('judul layar kasus = nama assignment flow InputAddendumTreatyIn', () => {
    expect(JUDUL_POSISI).toEqual({
      ReasTreatyInAdmin: 'Input Realitation',
      ReasTreatyInSecHead: 'Acceptance by Sec Treaty',
      ReasTreatyInDeptHead: 'Acceptance by Dept. Head',
    })
  })

  it('tab layout group DetailPolicyTreatyInAddGeneralEditable dan judul wadah ber-pyIncludeHeader', () => {
    expect(TAB_DATA).toEqual(['Old Data', 'New Data', 'Value Difference'])
    expect(Object.values(BAGIAN)).toEqual([
      'Installment Data Information',
      'OGP',
      'ONP',
      'Previous Premium',
      'Current Premium',
      'Total Difference',
    ])
  })

  it('tombol layar kasus', () => {
    // Add grid spreading dibuang (keputusan work owner 07-10-2026)
    expect([TOMBOL.chooseBusiness, TOMBOL.save, TOMBOL.submit, TOMBOL.hitungSelisih]).toEqual([
      'Choose Business',
      'Save',
      'Submit',
      'Calculate Value Difference',
    ])
  })

  it('grid spreading dan angsuran section Prop', () => {
    expect(Object.values(KOLOM_SPREADING)).toEqual(['Type Treaty', '% Share', 'Premium', '% Share', 'Claim', 'Total'])
    expect([
      KOLOM_ANGSURAN.no,
      KOLOM_ANGSURAN.dueDate,
      KOLOM_ANGSURAN.pct,
      KOLOM_ANGSURAN.premium,
      KOLOM_ANGSURAN.total,
    ]).toEqual(['No', 'Due Date', '% Installment', 'Premium Nusantara Re', 'Total Payment'])
  })

  it('AddPremi / DetailPolicyAddPremiDetail / InstallmentList', () => {
    expect(KOLOM_XOL.map((k) => k.judul)).toEqual([
      '',
      'Currency',
      'Gross Premium',
      'Deduction',
      'PPN',
      'PPh',
      'Net Premium',
      'Net Premium After PPN',
      'Net Premium After PPh',
      'Net Premium After Tax',
    ])
    expect(KOLOM_XOL[0]?.teks).toBe('Total For Currency')
    expect(KOLOM_XOL_RINCI).toHaveLength(16)
    expect(KOLOM_XOL_RINCI.filter((k) => k.judul !== '').map((k) => k.judul)).toEqual([
      'Gross Premium',
      'Deduction',
      'PPN',
      'PPh',
      'Net Premium',
      'Net Premium After PPN',
      'Net Premium After PPh',
      'Net Premium After Tax',
    ])
    expect(KOLOM_XOL_RINCI[2]?.teks).toBe('Part of')
    expect(KOLOM_RINCI_ANGSURAN.map((k) => k.judul)).toEqual([
      'Payment Date',
      'Percentage',
      '',
      'Balance Before Tax',
      'Balance Before Withholding Tax (PPH 2.2)',
      'Balance Due To',
    ])
  })

  it('popup BusinessAndSOBListEDM: judul jendela dan 8 kolom sesudah Choose', () => {
    expect(JUDUL.pilihBisnis).toBe('Business And SOB List')
    expect(TOMBOL.choose).toBe('Choose')
    expect(KOLOM_BISNIS.map((k) => k.judul)).toEqual([
      'ID Revision',
      'Previous ID',
      'Treaty Contract Name',
      'Proportion Type',
      'SOB',
      'Ceding',
      'Start Date',
      'End Date',
    ])
  })

  it('ListSuggestEDM, PolicyTreatyInDeclineConfirm, ShowPolicyNoTreaty_SC', () => {
    expect(Object.values(KOLOM_USULAN)).toEqual(['Date', 'PIC', 'Approval', 'Suggest'])
    expect(BARIS_PER_HALAMAN_USULAN).toBe(5)
    expect(JUDUL.tolak).toBe('Confirm Decline NB')
    expect(KONFIRMASI_TOLAK).toBe('Are you sure you want to decline this NB')
    expect([TOMBOL.yes, TOMBOL.no]).toEqual(['Yes', 'No'])
    expect(JUDUL.nomorPolis).toBe('Show PolicyNo')
    expect(NOMOR_DIAKSEP).toBe('telah diaksep menjadi')
    expect(TOMBOL.ok).toBe('OK')
  })
})

// DT `TreatyEDMListType` (screenshot Pega, work owner 07-10-2026) - harapan diketik dari screenshot; dan SAMA dengan
// peta backend `models.LabelJenisEDM` (yang mengisi `acuan.jenisEdm` dropdown Source of Change dan header).
describe('label jenis EDM', () => {
  const screenshot = { '1': 'Internal', '2': 'External', '3': 'Adjustment Premium', '4': 'Cancel Input' }

  it('teks VERBATIM screenshot; kode di luar daftar tampil apa adanya', () => {
    expect(LABEL_JENIS_EDM).toEqual(screenshot)
    expect(labelJenisEDM('9')).toBe('9')
  })

  it('sama dengan backend models.LabelJenisEDM', () => {
    const go = readFileSync(join(__dirname, '..', 'backend', 'models', 'edm_buat.go'), 'utf8')
    const baris = go.split('\n').find((b) => b.startsWith('var LabelJenisEDM = map[string]string{')) ?? ''
    const isi = Object.fromEntries([...baris.matchAll(/"(\d)": "([^"]+)"/g)].map((m) => [m[1], m[2]]))
    expect(isi).toEqual(screenshot)
  })
})
