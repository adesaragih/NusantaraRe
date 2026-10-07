// Tab Spreading kasus FIRE (tiket 48) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { describe, expect, it } from 'vitest'

import { SPREADING as S } from '../labels'
import { TEKS_SPREADING } from '../labels'
import { lebihDari100, templateLebih, treatyKembar } from './TabSpreading'

const KORPUS = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\'
const SUMBER = readFileSync(join(__dirname, 'TabSpreading.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe.skipIf(!existsSync(KORPUS + 'InputDtlSpreadingCoverage_FacIn.xml'))('label tab Spreading = korpus', () => {
  const ada = (berkas: string, teks: readonly string[], tag = 'pyValue') => {
    const xml = readFileSync(KORPUS + berkas, 'utf-8')
    return teks.filter((t) => !xml.includes(`<${tag}>${t}</${tag}>`))
  }
  it('% Share RNM, Copy Spreading, Type Treaty / % Share, Copy To All Spreading (InputInwardFacultativeDtl)', () => {
    const xml = readFileSync(KORPUS + 'InputInwardFacultativeDtl.xml', 'utf-8')
    expect(xml).toContain(`<pyLabelFieldValue>${S.percentShare}</pyLabelFieldValue>`)
    for (const t of [S.judulCopy, ...S.kolomTemplate, S.salinSemua]) expect(xml).toContain(`>${t}<`)
  })
  it('kolom lokasi / item / total / coverage / spreading / ringkasan', () => {
    expect(ada('CoverageSpreadingList.xml', S.kolomLokasi.slice(1))).toEqual([])
    expect(ada('PropertyItemListCoverageSpreading.xml', S.kolomItem)).toEqual([])
    expect(ada('PropertyItemListCoverageSpreading.xml', S.kolomTotal.slice(1))).toEqual([])
    expect(ada('InputCoverageSpreadingFire.xml', S.kolomCoverage)).toEqual([])
    expect(ada('InputDtlSpreadingCoverage_FacIn.xml', S.kolomSpreading)).toEqual([])
    expect(ada('SummarySpreading_Section.xml', S.kolomRingkasanTreaty)).toEqual([])
  })
  it('detail coverage (SpreadingItem)', () => {
    const xml = readFileSync(KORPUS + 'SpreadingItem.xml', 'utf-8')
    for (const t of [S.detailCoverage.tsi, S.detailCoverage.tsiLiability, S.detailCoverage.tsiNusantaraRe, S.detailCoverage.premiNusantaraRe]) {
      expect(xml).toContain(`>${t}<`)
    }
  })
})

describe('TabSpreading - aturan', () => {
  it('lebihDari100 eksak (tanpa float)', () => {
    expect([lebihDari100('100'), lebihDari100('100.0000'), lebihDari100('100.0001'), lebihDari100('99.9999'), lebihDari100('150')]).toEqual([
      false,
      false,
      true,
      false,
      true,
    ])
  })
  it('Σ % Share template Copy Spreading <= 100 (cekSpreadingFactIn_Act)', () => {
    const t = (s: string) => ({ treatyType: 'UJI', treatyName: 'UJI', sharePercentage: s })
    expect(templateLebih([t('60'), t('40')])).toBe(false)
    expect(templateLebih([t('60'), t('40.01')])).toBe(true)
    expect(templateLebih([])).toBe(false)
  })
  it('% Share RNM -> hitung-share berjeda; Copy To All -> salin; Save -> simpan; hitungan di backend (tanpa float)', () => {
    expect(SUMBER).toContain('hitung(() => hitungShareSpreading(caseId, nilai.trim()))')
    expect(SUMBER).toContain('hitung(() => salinSpreading(caseId, share.trim(), v.template))')
    expect(SUMBER).toContain('await simpanSpreading(caseId, { ...v, percentShare: share.trim() })')
    expect(SUMBER).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|parseInt/)
  })
})

describe('Proteksi Copy Spreading (CalcultePersentageSpeading_Act)', () => {
  const t = (treatyType: string, treatyName: string) => ({ treatyType, treatyName, sharePercentage: '50' })
  const AKT = 'D:\\migrasi\\RNM\\NB FacIn\\Activity\\CalcultePersentageSpeading_Act.xml'

  it.skipIf(!existsSync(AKT))('pesan = korpus verbatim', () => {
    const xml = readFileSync(AKT, 'utf-8')
    expect(xml).toContain(`<PropertiesValue>"${TEKS_SPREADING.treatySama}"</PropertiesValue>`)
  })

  it('syarat QS hanya di backend: "QS" dicari di REINSURANCETYPE.NOTE (RDBList GetTreatyName), bukan nama dropdown', () => {
    expect(SUMBER).not.toMatch(/includes\('QS'\)/)
    expect(SUMBER).toContain('hitung(() => salinSpreading(caseId, share.trim(), v.template))')
  })

  it('Type Treaty kembar ditolak; baris belum dipilih diabaikan', () => {
    expect(treatyKembar([t('UJI1', 'A'), t('UJI1', 'A')])).toBe(true)
    expect(treatyKembar([t('UJI1', 'A'), t('UJI2', 'B')])).toBe(false)
    expect(treatyKembar([t('', ''), t('', '')])).toBe(false)
    expect(treatyKembar([t('UJI1', 'A'), t(' UJI1 ', 'A')])).toBe(true)
  })

  it('Copy To All ditahan selama pesan tampil', () => {
    expect(SUMBER).toContain('disabled={lebih || kembar || v.template.length === 0 || menghitung}')
  })
})
