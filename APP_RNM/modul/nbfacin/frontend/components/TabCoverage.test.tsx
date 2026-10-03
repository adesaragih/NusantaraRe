// Paritas tab Coverage FIRE (tiket 43, C1). Data uji sintetis; uang = teks desimal.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { FORM_COV as F, LABEL_COVERAGE_BASIS, TEKS_COVERAGE } from '../labels'
import FormCoverage, { adaGalatCoverage, angkaSah, coverageBaru } from './FormCoverage'
import { bolehNetRate, totalPremiItem } from './TabCoverage'

const SUMBER_TAB = readFileSync(join(__dirname, 'TabCoverage.tsx'), 'utf8').replace(/\r\n/g, '\n')
const SUMBER_FORM = readFileSync(join(__dirname, 'FormCoverage.tsx'), 'utf8').replace(/\r\n/g, '\n')
const form = (c = coverageBaru()) => renderToStaticMarkup(<FormCoverage caseId="NB-1" c={c} tsiItem="1000000" ubah={() => {}} />)
const label = (html: string) => [...html.matchAll(/class="field__label"[^>]*>([^<]*)/g)].map((m) => m[1])

describe('FormCoverage', () => {
  it('coverage baru: basis Sum Insured (1), % Indemnity 100', () => {
    const c = coverageBaru({ id: 'ID1', oldId: '100815', nama: 'UJI' })
    expect([c.coverageBasis, c.indemnityPercentage, c.coverage, c.oldId, c.coverageNote]).toEqual(['1', '100', 'ID1', '100815', 'UJI'])
  })

  it('basis 1: medan umum; tanpa medan khusus First Loss / EML / Sub Limit', () => {
    const l = label(form())
    expect(l).toContain(LABEL_COVERAGE_BASIS)
    for (const x of [F.coverage, F.day, F.tsi, F.rate, F.discountPercentage, F.netRate, F.indemnityPercentage, F.lostLimit, F.premium]) {
      expect(l).toContain(x.label)
    }
    for (const x of [F.firstLoss, F.firstScale, F.indemnity, F.emlPml, F.sublimit]) expect(l).not.toContain(x.label)
  })

  it('basis 2 menampilkan First Loss / First Scale / Indemnity; basis 3 EML/PML; basis 4 Sub Limit; basis 5 pesan', () => {
    const l2 = label(form({ ...coverageBaru(), coverageBasis: '2' }))
    for (const x of [F.firstLoss, F.firstScale, F.indemnity, F.tsiLiability]) expect(l2).toContain(x.label)
    expect(label(form({ ...coverageBaru(), coverageBasis: '3' }))).toContain(F.emlPml.label)
    expect(label(form({ ...coverageBaru(), coverageBasis: '4' }))).toContain(F.sublimit.label)
    expect(form({ ...coverageBaru(), coverageBasis: '5' })).toContain(TEKS_COVERAGE.layeringBelum)
  })

  it('TSI tampil = TSI item (format Indonesia); Gross Rate bertanda wajib', () => {
    const html = form()
    expect(html).toContain('>1.000.000</div>')
    expect(html).toMatch(new RegExp(`${F.rate.label.replace('‰', '‰')}<span class="field__req">\\*</span>`))
  })

  it('hitung lewat backend: Gross Premium = mode amount, lainnya percent; isian tak sah tidak dihitung', () => {
    expect(SUMBER_FORM).toContain("uang('premium', F.premium.label, 'amount')")
    expect(SUMBER_FORM).toMatch(/if \(MEDAN_ANGKA\.some\(\(k\) => !angkaSah\(baru\[k\]\)\) \|\| baru\.coverageBasis === '5'\) return/)
    expect(SUMBER_FORM).toMatch(/hitungCoverage\(caseId, \{\s*coverage: baru,\s*tsi: tsiItem,\s*mode,\s*modeDiskon: modeDiskon\.current,\s*isAdjustable: item\?\.isAdjustable,/)
    // DiscountStatus: % Discount = percent, Discount = amount.
    expect(SUMBER_FORM).toContain("angkaDiskon('discountPercentage', F.discountPercentage.label, 'percent')")
    expect(SUMBER_FORM).toContain("angkaDiskon('discount', F.discount.label, 'amount')")
    expect(SUMBER_FORM).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|parseInt/)
  })

  it('angka sah: kosong / desimal bertitik; koma ditolak', () => {
    expect(['', '1.5', '0'].every(angkaSah)).toBe(true)
    expect(angkaSah('1,5')).toBe(false)
    expect(adaGalatCoverage([{ ...coverageBaru(), rate: '2,1' }])).toBe(true)
  })
})

describe('TabCoverage', () => {
  it('Total Gross Premium item = Σ premi coverage eksak; tanpa coverage = nilai server', () => {
    const c = (p: string) => ({ ...coverageBaru(), premium: p })
    expect(totalPremiItem([c('0.1'), c('0.2'), c('99999999999999999999.12345678')], '5')).toBe('99999999999999999999.42345678')
    expect(totalPremiItem([], '7.5')).toBe('7.5')
    expect(totalPremiItem(undefined, undefined)).toBe('')
  })

  it('Tambah: kosong -> lima coverage otomatis; ada -> satu coverage baru; Save = PUT objek; tanpa float', () => {
    expect(SUMBER_TAB).toMatch(/const h = await coverageOtomatis\(\)\s*ubahCoverage\(o, i, h\.baris\.map\(\(b\) => coverageBaru\(b\)\)\)/)
    expect(SUMBER_TAB).toContain('ubahCoverage(o, i, [...ada, coverageBaru()])')
    expect(SUMBER_TAB).toContain('const h = await simpanObjek(caseId, objek)')
    expect(SUMBER_TAB).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|parseInt/)
  })
})

describe('FormCoverage - medan uang berformat ribuan', () => {
  it('Gross Premium, Discount, Limit of Liability = IsianUang; rate / persen tetap isian biasa', () => {
    expect(SUMBER_FORM).toContain("uang('limitOfLiability', F.limitOfLiability.label)")
    expect(SUMBER_FORM).toContain("return k === 'discount' ? <IsianUang {...props} /> : <Field {...props} />")
    const html = form({ ...coverageBaru(), premium: '2500000.5' })
    expect(html).toContain('value="2.500.000,5"')
  })
})

describe('TabCoverage - net rate (tiket 44, CekNetRate_ACT)', () => {
  const c = (oldId: string) => ({ ...coverageBaru(), oldId })
  it('‰ Total Net Rate dapat diisi hanya bila kelima OLDID ada (cocok harfiah)', () => {
    expect(bolehNetRate(['FLEXAS', '4.1A CC', '4.3', '4.2 PRGBI', 'OTHERS'].map(c))).toBe(true)
    expect(bolehNetRate(['FLEXAS', '4.1A CC', '4.3', '4.2 PRGBI'].map(c))).toBe(false)
    expect(bolehNetRate(['flexas', '4.1A CC', '4.3', '4.2 PRGBI', 'OTHERS'].map(c))).toBe(false)
    expect(bolehNetRate(undefined)).toBe(false)
  })

  it('perubahan net rate -> POST hitung-net-rate sesudah jeda; jawaban lama dibuang', () => {
    expect(SUMBER_TAB).toContain('hitungNetRate(caseId, {')
    expect(SUMBER_TAB).toContain('if (n === nomorNet.current) ubahCoverage(o, i, h.coverages)')
    expect(SUMBER_TAB).toContain('{bolehNetRate(covs) ? (')
  })
})
