// Tab Inw Fac Cedant Panels kasus FIRE (tiket 49) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

import { CEDANT as C, OPSI_SHARE_CEDANT_TYPE, TEKS_CEDANT as T } from '../labels'
import { adaShareLebih, galatCedant, totalShareLebih, kaliPersen, mataUangCedant, shareSah, tambahCedant, teksShareOfCeding } from './TabCedant'

const KORPUS = 'D:\\migrasi\\RNM\\'
const SUMBER = readFileSync(join(__dirname, 'TabCedant.tsx'), 'utf8').replace(/\r\n/g, '\n')
const b = (cedingCoName: string, shareCeding: string) => ({ cedingCo: cedingCoName === '' ? '' : 'UJI-ID', cedingCoName, shareCeding })

describe.skipIf(!existsSync(KORPUS + 'NB FacIn\\Section\\InputInwardFacultativeDtl.xml'))('label = korpus', () => {
  it('blok atas, grid, tombol (InputInwardFacultativeDtl tab Inw Fac Cedant Panels)', () => {
    const xml = readFileSync(KORPUS + 'NB FacIn\\Section\\InputInwardFacultativeDtl.xml', 'utf-8')
    expect(xml).toContain('<pyTitle>Inw Fac Cedant Panels</pyTitle>')
    for (const t of [C.sob, C.shareCedantType, C.percentShare, C.totalTsi, C.totalPremi]) {
      expect(xml).toContain(`<pyLabelFieldValue>${t}</pyLabelFieldValue>`)
    }
    for (const t of [C.kolomCeding, C.kolomShare, C.tambah, C.hapus]) expect(xml).toContain(`>${t}<`)
  })
  it('pilihan radio (DDL ShareCedantType.xml)', () => {
    const ddl = readFileSync(KORPUS + 'DDL\\ShareCedantType.xml', 'utf-8')
    for (const o of OPSI_SHARE_CEDANT_TYPE) {
      expect(ddl).toContain(`<pyLocalizedValue>${o.label}</pyLocalizedValue>`)
      expect(ddl).toContain(`<pyStandardValue>${o.value}</pyStandardValue>`)
    }
  })
  it('pesan proteksi verbatim (ProtectShareCedant_Act)', () => {
    const xml = readFileSync(KORPUS + 'NB FacIn\\Activity\\ProtectShareCedant_Act.xml', 'utf-8')
    for (const t of [T.shareTidakSah, T.cedingKosong, T.listKosong]) expect(xml).toContain(`<PropertiesValue>"${t}"</PropertiesValue>`)
    const set = readFileSync(KORPUS + 'NB FacIn\\Activity\\SetTSIPremiCedant_Act.xml', 'utf-8')
    expect(set).toContain(`<PropertiesValue>"${T.nilaiLebih}"</PropertiesValue>`)
    expect(set).toContain('<pyStepsPreCondParamsWhen>.ShareCeding&gt;100</pyStepsPreCondParamsWhen>')
  })
})

describe('TabCedant - aturan', () => {
  it('Add menyalin Ceding Co tab General; grid berisi -> hanya yang belum ada (kode); tanpa baris kosong (C-6)', () => {
    const c = (cedingCo: string, cedingCoName: string, shareCeding = '') => ({ cedingCo, cedingCoName, shareCeding })
    const umum = [c('UJI1', 'UJI A'), c('UJI2', 'UJI B')]
    expect(tambahCedant([], umum)).toEqual([c('UJI1', 'UJI A'), c('UJI2', 'UJI B')])
    expect(tambahCedant([c('UJI1', 'UJI A', '40')], umum)).toEqual([c('UJI1', 'UJI A', '40'), c('UJI2', 'UJI B')])
    expect(tambahCedant(umum, umum)).toEqual(umum)
    expect(tambahCedant([], [])).toEqual([])
    expect(SUMBER).not.toContain('PopupPilihAgent')
  })

  it('% Share > 0 dan <= 100, eksak', () => {
    expect(['50', '100', '100.0000', '0.0001'].map(shareSah)).toEqual([true, true, true, true])
    expect(['', '0', '0.0000', '-5', '100.0001', '150', 'x'].map(shareSah)).toEqual([false, false, false, false, false, false, false])
  })

  it('proteksi: list kosong; % Share tidak sah; ceding kosong', () => {
    expect(galatCedant([])).toEqual([T.listKosong])
    expect(galatCedant([b('UJI A', '60'), b('UJI B', '40')])).toEqual([])
    expect(galatCedant([b('UJI A', '0'), b('', '40')])).toEqual([T.shareTidakSah, T.cedingKosong])
  })

  it('proteksi layar hanya bila backend menyatakan wajib; Choose Ceding mengisi ID + nama', () => {
    expect(SUMBER).toContain('const pesan = v.wajib ? galatCedant(v.cedant) : []')
    expect(SUMBER).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|parseInt/)
  })
})

describe('TabCedant - % Share > 100 selalu ditolak (SetTSIPremiCedant_Act)', () => {
  it('eksak; teks bukan angka tidak membuat galat ini (dan tidak melempar)', () => {
    expect(adaShareLebih([b('UJI A', '100'), b('UJI B', '100.0001')])).toBe(true)
    expect(adaShareLebih([b('UJI A', '100.0000'), b('UJI B', '')])).toBe(false)
    expect(adaShareLebih([b('UJI A', 'x')])).toBe(false)
  })
})

describe('TabCedant - detail baris (CedingCedant + SetTSIPremiCedant_Act)', () => {
  it.skipIf(!existsSync(KORPUS + 'NB FacIn\\Section\\CedingCedant.xml'))('label detail = korpus', () => {
    const sec = readFileSync(KORPUS + 'NB FacIn\\Section\\CedingCedant.xml', 'utf-8')
    for (const k of C.kolomMataUang) expect(sec).toContain(`<pyValue>${k}</pyValue>`)
    expect(sec).toContain(`<pyLabelFieldValue>${C.dari.toUpperCase()}</pyLabelFieldValue>`)
    const akt = readFileSync(KORPUS + 'NB FacIn\\Activity\\SetTSIPremiCedant_Act.xml', 'utf-8')
    // Langkah 4 (CurrencyList) yang berjalan; langkah 5 (TotalSpreadAll) ber-label "//" = mati.
    expect(akt).toContain('<pyStepsObjectName>pyWorkPage.OfferFacIn.CurrencyList</pyStepsObjectName>')
    expect(akt).toContain('<pyStepsBlockName>//</pyStepsBlockName>')
    expect(akt).toContain('<PropertiesValue>@Math.divide((Primary.ShareCeding * .TSI),100,20)</PropertiesValue>')
  })

  it('"of" = SetShareOfCeding: Gross 100%, Share RNM PercentShare%, belum dipilih -> tersimpan', () => {
    expect(teksShareOfCeding('0', '10', '')).toBe('100%')
    expect(teksShareOfCeding('1', '10.00000000', '')).toBe('10%')
    expect(teksShareOfCeding('1', '2.5', '')).toBe('2.5%')
    expect(teksShareOfCeding('', '10', '25%')).toBe('25%')
  })

  it('kaliPersen eksak: a x b / 100', () => {
    expect(kaliPersen('100', '50000000000')).toBe('50000000000.00')
    expect(kaliPersen('33.3333', '3')).toBe('0.999999')
    expect(kaliPersen('', '1')).toBe('')
  })

  it('CurrencyList per tipe; % Share tidak sah atau tipe kosong -> kosong', () => {
    const d = [{ currency: 'IDR', tsiRnm: '5000000000', premiumRnm: '939100000', tsiGross: '50000000000', premiumGross: '9391000000' }]
    expect(mataUangCedant('100', '0', d)).toEqual([{ currency: 'IDR', tsi: '50000000000.00', premium: '9391000000.00' }])
    expect(mataUangCedant('50', '1', d)).toEqual([{ currency: 'IDR', tsi: '2500000000.00', premium: '469550000.00' }])
    expect(mataUangCedant('', '0', d)).toEqual([])
    expect(mataUangCedant('100', '', d)).toEqual([])
  })
})

describe('TabCedant - Σ % Share cedant <= 100 selalu (C-7)', () => {
  it('eksak; isian bukan angka diabaikan; Save ditahan', () => {
    expect(totalShareLebih([b('UJI A', '50'), b('UJI B', '60')])).toBe(true)
    expect(totalShareLebih([b('UJI A', '50'), b('UJI B', '50')])).toBe(false)
    expect(totalShareLebih([b('UJI A', '50'), b('UJI B', '50.0001')])).toBe(true)
    expect(totalShareLebih([b('UJI A', 'x'), b('UJI B', '')])).toBe(false)
    expect(SUMBER).toContain('if (typeKosong || pesan.length > 0 || shareLebih || totalLebih) return')
    // Tiket 64 K2: Share Cedant Type juga wajib di tahap Binding (wajibType), pesan pilihan work owner 08-10-2026.
    expect(SUMBER).toContain('const typeWajib = v.wajib || v.wajibType === true')
    expect(SUMBER).toContain("const typeKosong = typeWajib && v.shareCedantType === ''")
  })
})

describe('TabCedant - total per mata uang (Dtl S62 IsFire, TotalSpreadAll)', () => {
  it.skipIf(!existsSync(KORPUS + 'NB FacIn\\Section\\InputInwardFacultativeDtl.xml'))('grid TotalSpreadAll ada di tab Cedant', () => {
    const xml = readFileSync(KORPUS + 'NB FacIn\\Section\\InputInwardFacultativeDtl.xml', 'utf-8').split('\n')
    const tab = xml.slice(60509, 74900).join('\n')
    expect(tab).toContain('<pyPageListProperty>TotalSpreadAll.pxResults</pyPageListProperty>')
    for (const k of ['Total TSI Top Risk Spreaded', 'Total TSI Spreaded', 'Total Limit of Liability', 'Total Premium Spreaded']) {
      expect(tab).toContain(`<pyValue>${k}</pyValue>`)
    }
  })
  it('layar memakai ringkasan Spreading per mata uang saja', () => {
    expect(SUMBER).toContain('{sp !== null && <RingkasanSpreading v={sp} hanyaMataUang tanpaTopRisk={aneka || cargo} />}')
  })
})
