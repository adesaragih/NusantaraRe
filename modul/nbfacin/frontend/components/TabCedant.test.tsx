// Tab Inw Fac Cedant Panels kasus FIRE (tiket 49) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

import { CEDANT as C, OPSI_SHARE_CEDANT_TYPE, TEKS_CEDANT as T } from '../labels'
import { adaShareLebih, galatCedant, shareSah, tambahCedant } from './TabCedant'

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
  it('Choose Ceding (CedingCedant.xml) dan pilihan radio (DDL ShareCedantType.xml)', () => {
    const sec = readFileSync(KORPUS + 'NB FacIn\\Section\\CedingCedant.xml', 'utf-8')
    expect(sec).toContain(`<pyLabel>${C.pilihCeding}</pyLabel>`)
    expect(sec).toContain(`<pyWindowName>${C.judulPilih}</pyWindowName>`)
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
  it('Add: grid kosong menyalin Ceding Co General; tanpa daftar -> satu baris kosong; grid berisi -> baris kosong', () => {
    const umum = [b('UJI A', ''), b('UJI B', '')]
    expect(tambahCedant([], umum)).toEqual([b('UJI A', ''), b('UJI B', '')])
    expect(tambahCedant([], [])).toEqual([b('', '')])
    expect(tambahCedant([b('UJI A', '40')], umum)).toEqual([b('UJI A', '40'), b('', '')])
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
    expect(SUMBER).toContain('ubahBaris(pilihUntuk, { cedingCo: a.id, cedingCoName: a.name })')
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
