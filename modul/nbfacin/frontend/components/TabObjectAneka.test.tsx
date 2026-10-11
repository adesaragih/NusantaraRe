// Paritas tab Object Aneka (tiket 73 P1) - grid lokasi `InputDtlObject_FacIn.xml` L36463, form lokasi
// `InputDtlObjectLocation_FacIn.xml`, detail okupasi `InputOkupasiAneka_FacIn.xml`, popup HARNESS `ChooseOccupation`.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { GRID_GOLF_ANEKA, GRID_LOKASI_ANEKA, GRID_OKUPASI_ANEKA, POPUP_OKUPASI_ANEKA } from '../labels'
import { PopupOkupasiAneka } from './TabObjectAneka'

const SUMBER = readFileSync(join(__dirname, 'TabObjectAneka.tsx'), 'utf8').replace(/\r\n/g, '\n')
const API = readFileSync(join(__dirname, '..', 'api.ts'), 'utf8').replace(/\r\n/g, '\n')
const HALAMAN = readFileSync(join(__dirname, '..', 'pages', 'InwardFacultative.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('TabObjectAneka', () => {
  it('label grid = label Pega (No / Top Risk / Location; Occupation Code / Name; Golf Object Name / Currency / TSI - TSI Old hanya EDM, work owner 09-10-2026)', () => {
    expect(GRID_LOKASI_ANEKA).toEqual(['No', 'Top Risk', 'Location'])
    expect(GRID_OKUPASI_ANEKA).toEqual(['Occupation Code', 'Occupation Name'])
    expect(GRID_GOLF_ANEKA).toEqual(['Object Name', 'Currency', 'TSI'])
  })

  it('popup Choose Occupation = HARNESS ChooseOccupation (Search, ID / Occupation Name, Choose)', () => {
    const html = renderToStaticMarkup(<PopupOkupasiAneka caseId="NB-1" onTutup={() => {}} onPilih={() => {}} />)
    expect(html).toContain(POPUP_OKUPASI_ANEKA.judul)
    expect(html).toContain('>Search<')
    for (const k of POPUP_OKUPASI_ANEKA.kolom) expect(html).toContain(`>${k}<`)
  })

  it('Choose Occupation: id → occupationId (OCCUPATION.ID, BUKAN OLDID), name → occupationName (GetOccupation_act)', () => {
    expect(SUMBER).toContain('ganti(pilih, { ...k, occupationId: b.id, occupationName: b.name })')
    expect(API).toContain('`/api/nbfacin/kasus/${encodeURIComponent(caseId)}/occupation-aneka`, { kueri: { q } }')
  })

  it('Save = PUT objek-aneka; 422 proteksi ditampilkan sebagai daftar pesan (pesanProteksiDari), galat lain lewat Gagal', () => {
    expect(SUMBER).toContain('const h = await simpanObjekAneka(caseId, daftar, proRate)')
    expect(SUMBER).toContain('const p = pesanProteksiDari(err)')
    expect(API).toContain("minta<DaftarAneka>(jalurAneka(caseId), { metode: 'PUT', badan: { lokasi, proRatePercent } })")
  })

  it('form lokasi: Top Risk hanya bila tampilanAneka.topRisk; grid Golf bila golf; grid okupasi kecuali Golf murni', () => {
    expect(SUMBER).toContain('{tampilan.topRisk && (')
    expect(SUMBER).toContain('{tampilan.golf && (')
    expect(SUMBER).toContain('{tampilOkupasi(aneka, tampilan) && (')
  })

  it('halaman case: tab Object Aneka di belakang FIRE (FIRE selalu menang - layar FIRE tidak berubah)', () => {
    expect(HALAMAN).toContain(
      'const kasusAneka = !kasusFire && lini !== undefined && tampilanAneka !== undefined && (lini.anekaNotBonding || lini.golf || lini.bonding)',
    )
    expect(HALAMAN).toMatch(/\{tab === 'Object' && kasusFire \? \([\s\S]{0,400}?<TabObject caseId=\{kasus\.caseId\} insuredName=\{insured\} \/>/)
    expect(HALAMAN).toContain("tab === 'Object' && kasusAneka && tampilanAneka !== undefined ? (")
    expect(HALAMAN).toContain('<TabObjectAneka caseId={kasus.caseId} tampilan={tampilanAneka} aneka={lini?.aneka === true} />')
  })
})

describe('TabCoverageAneka (tiket 73 P2)', () => {
  const COV = readFileSync(join(__dirname, 'TabCoverageAneka.tsx'), 'utf8').replace(/\r\n/g, '\n')

  it('halaman case: tab Coverage Aneka di belakang FIRE', () => {
    expect(HALAMAN).toContain("tab === 'Coverage' && kasusAneka && tampilanAneka !== undefined ? (")
    expect(HALAMAN).toMatch(/: tab === 'Coverage' && kasusFire \? \(\s*<TabCoverage /)
  })

  it('Save mengirim % Pro Rate case; tab Object mengirim balik nilai GET (PUT menulis kosong = NULL)', () => {
    expect(COV).toContain('const h = await simpanObjekAneka(caseId, lokasi, proRateRef.current)')
    expect(SUMBER).toContain('const h = await simpanObjekAneka(caseId, daftar, proRate)')
    expect(SUMBER).toContain("setProRate(h.proRatePercent ?? '')")
  })

  it('hitung membawa % Pro Rate case dan memasang nilai barunya; jawaban usang dibuang', () => {
    expect(COV).toContain("hitungCoverageAneka(caseId, { proRatePercent: proRateRef.current, baris: [{ coverage: c, pemicu }] })")
    expect(COV).toContain('if (nomorHitung.current.get(kunci) !== nomor) return')
  })

  it('tombol salin: lokasi pertama saja; Copy Deductible To All hanya objek pertama (pxListSubscript = 1)', () => {
    expect(COV).toContain('{bolehSalin && p.l === 0 && (')
    expect(COV).toMatch(/\{p\.j === 0 && \(\s*<button type="button" className="btn nbf-tema__tombol btn--sm" onClick=\{\(\) => salinDeductible\(p\)\}>\s*\{SALIN_ANEKA\.deductibleSemua\}/)
  })

  it('tidak ada komponen yang didefinisikan di dalam komponen (isian kehilangan fokus)', () => {
    expect(COV).not.toMatch(/^[ 	]+function [A-Z]\w*\(/m)
    expect(SUMBER.match(/^function [A-Z]\w*\(/gm)?.length).toBeGreaterThan(0) // instrumen: pola menangkap fungsi tingkat atas
  })
})
