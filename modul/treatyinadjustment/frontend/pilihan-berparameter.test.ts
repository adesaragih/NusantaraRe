// Isi dropdown panel New yang dahulu KOSONG (kotak teks) — audit 8 Oktober
// 2026: "pastikan isi dari dropdown … sama seperti yang ada di treaty in".
//
//   • `associated` rincian Layers Non-Prop (Layer · Part Of · Cover ·
//     Currency Relation · Note Reinstatement) — daftar opsi-limits Treaty In;
//   • RD BERPARAMETER — Class of Business (`pTreatyGroupId`) dan Spreading
//     Type (`TreatyGroupID`, `StartDate`) — rute Treaty In yang sama;
//   • `pySetValueOnSelect` — pasangan kode (`CurrencyID`, `TreatyGroupID`,
//     `ClassOfBusinessID`) ikut diisi saat sebuah pilihan dipilih.

import { describe, expect, it, vi } from 'vitest'

import type { SumberPilihan } from './ekspor/jenis'
import { opsiKepala } from './komponen/SisiPenyesuaian'
import { ikutTerpilih, nilaiParam, opsiUntuk, type DataOpsi, type Opsi, type PemuatOpsi } from './komponen/pilihan'

const LAYER = [{ value: 'Layer', label: 'Layer' }]
const DATA: DataOpsi = {
  limits: {
    jenisTreaty: [],
    kelompokTreaty: [],
    mataUang: [{ id: '10026', nama: 'IDR' }],
    jenisLayer: LAYER,
    cover: [{ value: 'Risk', label: 'Risk' }],
    relasiMataUang: [{ value: 'AND', label: 'AND' }],
    catatanReinstatement: [{ value: 'Free', label: 'Free' }],
  },
}

const KELAS: SumberPilihan = {
  sumber: 'reportdefinition',
  rd: 'BrowseTreatyBusinessWOType_RD',
  nilai: 'BIZNAME',
  param: { pTreatyGroupId: '.TreatyGroupID' },
  setel: [{ target: 'ClassOfBusinessID', dari: 'BizCode' }],
}
const SPREAD_XOL: SumberPilihan = {
  sumber: 'reportdefinition',
  rd: 'BrowseTreatyArrangement_ParentReinsMasterTrt',
  nilai: 'ReinsTypeName',
  tampil: 'ReinsTypeName',
  param: { TreatyGroupID: '.TreatyGroupList(1).TreatyGroupID', StartDate: 'TreatyIn.Commencement', TreatyDescID: '"10001"' },
}

describe('associated rincian Layers — daftar Treaty In', () => {
  it('Layer dan Part Of memakai jenisLayer (seperti TabLimitsNonProp), lalu Cover, Currency Relation, Note', () => {
    for (const k of ['LayerType', 'LayerPartType']) expect(opsiUntuk({ sumber: 'associated' }, k, DATA)).toEqual(LAYER)
    expect(opsiUntuk({ sumber: 'associated' }, 'Cover', DATA)?.[0]?.value).toBe('Risk')
    expect(opsiUntuk({ sumber: 'associated' }, 'CurrencyRelation', DATA)?.[0]?.value).toBe('AND')
    expect(opsiUntuk({ sumber: 'associated' }, 'ReinstatementNote', DATA)?.[0]?.value).toBe('Free')
  })
})

describe('parameter RD dibaca atas halaman sel/medan', () => {
  const tempat = {
    halaman: { Commencement: '20260701', '.TreatyGroupID': '7' },
    sisi: { medan: { TreatyGroupID: '7' }, larik: { TreatyGroupList: [{ TreatyGroupID: '9' }] } },
  }
  it('harfiah, akar, skalar bertitik, baris larik ke-n', () => {
    expect(nilaiParam('"10001"', tempat)).toBe('10001')
    expect(nilaiParam('TreatyIn.Commencement', tempat)).toBe('20260701')
    expect(nilaiParam('.TreatyGroupID', tempat)).toBe('7')
    expect(nilaiParam('.TreatyGroupList(1).TreatyGroupID', tempat)).toBe('9')
    expect(nilaiParam('TempSprd.TreatyGroupID', tempat)).toBe('')
  })

  it('daftar berparameter: dimuat sekali per kunci, lalu terpakai', async () => {
    const muat = vi.fn()
    const pemuat: PemuatOpsi = {
      kelasBisnis: vi.fn(() => Promise.resolve([{ id: 'B01', nama: 'FIRE' }])),
      indukSpreading: vi.fn(() => Promise.resolve([{ reinsTypeId: '10002', reinsTypeName: 'QS 40%' }])),
    }
    // Belum dimuat: daftar KOSONG (dropdown), dan permintaan muat dikirim.
    expect(opsiUntuk(KELAS, 'ClassOfBusiness', { ...DATA, muat }, tempat, pemuat)).toEqual([])
    expect(muat).toHaveBeenCalledWith('kelas-bisnis|7', expect.any(Function))
    const [, ambil] = muat.mock.calls[0] as [string, () => Promise<Opsi[]>]
    expect(await ambil()).toEqual([{ value: 'FIRE', label: 'FIRE', id: 'B01' }])
    expect(pemuat.kelasBisnis).toHaveBeenCalledWith('7')
    // Spreading Non-Prop: grup dari `.TreatyGroupList(1)`, nilai = NAMA.
    opsiUntuk(SPREAD_XOL, 'SpreadingTypeXOL', { ...DATA, muat }, tempat, pemuat)
    expect(muat).toHaveBeenLastCalledWith('spreading|9|20260701', expect.any(Function))
    const muatan = { 'spreading|9|20260701': [{ value: '10002', label: 'QS 40%', id: '10002' }] }
    expect(opsiUntuk(SPREAD_XOL, 'SpreadingTypeXOL', { ...DATA, muatan, muat }, tempat, pemuat)).toEqual([
      { value: 'QS 40%', label: 'QS 40%', id: '10002' },
    ])
  })
})

describe('pySetValueOnSelect — pasangan kode ikut diisi', () => {
  const opsi: Opsi[] = [{ value: 'FIRE', label: 'FIRE', id: 'B01' }]
  it('pilihan di daftar → target = id; nilai di luar daftar → target dikosongkan', () => {
    expect(ikutTerpilih(KELAS, opsi, 'FIRE')).toEqual({ ClassOfBusinessID: 'B01' })
    expect(ikutTerpilih(KELAS, opsi, 'KETIK BEBAS')).toEqual({ ClassOfBusinessID: '' })
    expect(ikutTerpilih({ sumber: 'associated' }, opsi, 'FIRE')).toEqual({})
  })
})

describe('dropdown kepala — daftar opsi-kepala Treaty In', () => {
  it('memakai label rute Treaty In bila sudah dimuat, cadangan nilai mentah bila belum', () => {
    const kepala = {
      bordereaux: [{ value: 'reporting', label: 'Reporting' }],
      caraPembukuan: [{ value: 'underwriting', label: 'Underwriting Year' }],
      caraPembukuanNonProp: [],
      periodePelaporan: [],
    }
    expect(opsiKepala('AccountingMode', kepala)).toEqual([{ value: 'underwriting', label: 'Underwriting Year' }])
    expect(opsiKepala('Bordeaux', undefined).map((o) => o.value)).toEqual(['reporting', 'nonreporting'])
  })
})
