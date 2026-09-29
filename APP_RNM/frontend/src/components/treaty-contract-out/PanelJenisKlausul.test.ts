// Uji panel jenis klausul — tiket 08 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import type { AturanKlausul, JenisKlausul, Klausul } from '../../services/api'
import {
  aturanAnak,
  aturanInduk,
  barisSubjenis,
  formKlausulDari,
  formKlausulKosong,
  jenisBerkurs,
  keMasukKlausul,
  labelMedan,
  rencanaKonversi,
} from './PanelJenisKlausul'

const KODE = readFileSync(join(__dirname, 'PanelJenisKlausul.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//') && !b.trimStart().startsWith('*'))
  .join('\n')

function aturan(p: Partial<AturanKlausul>): AturanKlausul {
  return { jenis: 'EPI', anak: false, subjenis: '', medan: ['ReinsTypeID', 'Line', 'Rp', 'Usd'],
    wajib: ['ReinsTypeID', 'Rp', 'Usd'], turunan: null, ditahan: '', berkurs: false, konversi: '', sumber: 'SaveTreatyArrEPI_Act', ...p }
}

function klausul(p: Partial<Klausul>): Klausul {
  return { id: '10000001', treatyYear: '2026', treatyYearId: '1000001', treatyGroupId: '10001', treatyDescId: '10009',
    treatyDescName: 'UJI EPI', reinsTypeId: '10003', reinsTypeName: 'UJI QS', parentReinsTypeId: '00', subjenis: '',
    medan: {}, kurs: '', userId: 'UJI-ADMIN', tglUpdate: '', ...p }
}

describe('label medan VERBATIM per jenis', () => {
  it('bawaan, lalu penimpaan per jenis dan per subjenis', () => {
    expect(labelMedan(aturan({}), 'ReinsTypeID')).toBe('ReinsType')
    expect(labelMedan(aturan({ jenis: 'MinLOL' }), 'Pct')).toBe('Minimum LOL (%)')
    expect(labelMedan(aturan({ jenis: 'MinLOLMB' }), 'Pct')).toBe('Minimum LOL MB (%)')
    expect(labelMedan(aturan({ jenis: 'ExclutionTreaty', subjenis: 'Occupation' }), 'Usd')).toBe('TSI Less Than (USD)')
    expect(labelMedan(aturan({ jenis: 'ExclutionTreaty', subjenis: 'Object' }), 'Pct')).toBe('TSI BI >')
    expect(labelMedan(aturan({ jenis: 'CoinsPanel' }), 'CoIns_Min')).toBe('From')
  })
})

describe('form klausul', () => {
  it('kosong dan dari baris — hanya medan aturan', () => {
    const a = aturan({})
    expect(formKlausulKosong(a)).toEqual({ id: '', medan: { ReinsTypeID: '', Line: '', Rp: '', Usd: '' } })
    const f = formKlausulDari(a, klausul({ medan: { ReinsTypeID: '10003', Rp: '1000.5', Usd: '70', Kurs: 'X' } }))
    expect(f).toEqual({ id: '10000001', medan: { ReinsTypeID: '10003', Line: '', Rp: '1000.5', Usd: '70' } })
  })
  it('anak: Rp/Usd TURUNAN tidak dikirim; induk disebut', () => {
    const a = aturan({ jenis: 'EpiList', anak: true, medan: ['ReinsTypeID', 'Pct'], turunan: ['Rp', 'Usd'] })
    const m = keMasukKlausul(a, '10009', { id: '', medan: { ReinsTypeID: '10005', Pct: ' 25 ', Rp: '9' } }, '10003')
    expect(m).toEqual({ id: '', descId: '10009', anak: true, subjenis: '', parentReinsTypeId: '10003',
      medan: { ReinsTypeID: '10005', Pct: '25' } })
  })
  it('induk tidak menyebut induk', () => {
    expect(keMasukKlausul(aturan({}), '10009', formKlausulKosong(aturan({})), '00').parentReinsTypeId).toBe('')
  })
  it('baris exclusion disaring per subjenis', () => {
    const d = [klausul({ id: '1', subjenis: 'Occupation' }), klausul({ id: '2', subjenis: 'Clause' })]
    expect(barisSubjenis(d, 'Clause').map((k) => k.id)).toEqual(['2'])
    expect(barisSubjenis(d, '').length).toBe(2)
  })
  it('aturan induk/anak dari jenis', () => {
    const j: JenisKlausul = { id: '10009', descName: 'UJI EPI', isXol: '0', statusAktif: '1', catatan: '',
      aturan: [aturan({}), aturan({ jenis: 'EpiList', anak: true })] }
    expect(aturanInduk(j).map((a) => a.jenis)).toEqual(['EPI'])
    expect(aturanAnak(j)?.jenis).toBe('EpiList')
  })
})

describe('kabel', () => {
  it('Cancel membuang isian panel INI saja (AC 29) — state form per GridAturan', () => {
    expect(KODE).toMatch(/function GridAturan[\s\S]*useState<FormKlausul \| null>/)
    expect(KODE).toContain('KLAUSUL_TCO.cancel')
  })
  it('jenis ditahan tampil dengan alasannya, tanpa form', () => {
    expect(KODE).toContain('aturan.ditahan !== ')
  })
  it('uang/persen tidak menjadi angka JavaScript', () => {
    expect(KODE).not.toMatch(/Number\(|parseFloat|toFixed/)
  })
  it('Show Child dan Close Child dari label', () => {
    expect(KODE).toContain('KLAUSUL_TCO.showChild')
    expect(KODE).toContain('KLAUSUL_TCO.closeChild')
  })
})

describe('kurs (tiket 11)', () => {
  it('rencana konversi mengikuti HitungRpUsd_depan dan CalculateTSIExcludeTreaty', () => {
    expect(rencanaKonversi('RpKeUsd', 'Rp')).toEqual({ dari: 'Rp', ke: 'Usd', skala: '8' })
    expect(rencanaKonversi('RpKeUsd', 'Usd')).toBeNull()
    expect(rencanaKonversi('DuaArah', 'Rp')).toEqual({ dari: 'Rp', ke: 'Usd', skala: '4' })
    expect(rencanaKonversi('DuaArah', 'Usd')).toEqual({ dari: 'Usd', ke: 'Rp', skala: '8' })
    expect(rencanaKonversi('', 'Rp')).toBeNull()
  })
  it('jenis berkurs bila salah satu aturannya berkurs', () => {
    const j: JenisKlausul = { id: '10009', descName: 'UJI EPI', isXol: '0', statusAktif: '1', catatan: '',
      aturan: [aturan({ berkurs: true, konversi: 'RpKeUsd', turunan: ['Usd'] }), aturan({ jenis: 'EpiList', anak: true, berkurs: true })] }
    expect(jenisBerkurs(j)).toBe(true)
    expect(jenisBerkurs({ ...j, aturan: [aturan({})] })).toBe(false)
  })
  it('induk berkurs: Usd TURUNAN tidak dikirim', () => {
    const a = aturan({ berkurs: true, konversi: 'RpKeUsd', turunan: ['Usd'] })
    const m = keMasukKlausul(a, '10009', { id: '', medan: { ReinsTypeID: '10003', Line: '', Rp: '1000', Usd: '0.06451613' } }, '00')
    expect(Object.keys(m.medan)).not.toContain('Usd')
  })
  it('kurs dibaca dari server; tanpa kurs Add nonaktif; konversi di server', () => {
    expect(KODE).toContain('ambilKursTahun(tahunID)')
    expect(KODE).toContain('disabled={aturan.berkurs && !kursAda}')
    expect(KODE).toContain('konversiKurs(tahunID, r.dari, nilai, r.skala)')
  })
})
