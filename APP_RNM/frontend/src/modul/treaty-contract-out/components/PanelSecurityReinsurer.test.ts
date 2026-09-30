// Uji panel security — tiket 06 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { formSecurityDari, formSecurityKosong, keMasukSecurity } from './PanelSecurityReinsurer'

const KODE = readFileSync(join(__dirname, 'PanelSecurityReinsurer.tsx'), 'utf8')
  .split('\n')
  .filter((b) => !b.trimStart().startsWith('//'))
  .join('\n')

const BARIS = {
  id: '1000011', thnTreaty: '2026', reasId: '1000007', reasSecurity: 'UJI-R2', clientName: 'UJI REAS DUA',
  pctShare: '33.33333333', topId: '', tpTreaty: '', userId: '',
}

describe('form security', () => {
  it('kosong (InputNewSecurityReinsurer) dan dari baris (ShowEditSecurityReinsurer)', () => {
    expect(formSecurityKosong()).toEqual({ id: '', reasSecurity: '', clientName: '', pctShare: '' })
    expect(formSecurityDari(BARIS)).toEqual({ id: '1000011', reasSecurity: 'UJI-R2', clientName: 'UJI REAS DUA', pctShare: '33.33333333' })
  })
  it('badan simpan membawa ID baris dan TIDAK membawa nama tampil', () => {
    const m = keMasukSecurity({ id: '1000011', reasSecurity: ' UJI-R3 ', clientName: 'KARANGAN', pctShare: ' 20,5 ' })
    expect(m).toEqual({ id: '1000011', reasSecurity: 'UJI-R3', pctShare: '20,5' })
    expect(Object.keys(m)).not.toContain('clientName')
  })
})

describe('kabel', () => {
  it('Security ID hanya dibaca; Security Name dari master yang sama dengan reinsurer', () => {
    expect((KODE.match(/readOnly/g) ?? []).length).toBe(1)
    expect(KODE).toContain('const master = useCariReinsurerMaster(setGalat)')
  })
  it('hapus langsung menurut ID, tanpa konfirmasi (Pega: aksi refresh)', () => {
    expect(KODE).toContain('hapusSecurity(tahunID, kontrakID, reinsurerID, s.id)')
    expect(KODE).not.toMatch(/confirm\(/)
  })
  it('persen tidak menjadi angka JavaScript', () => {
    // `(?<![A-Za-z])`: `formatNumber(` (pemformat TEKS bersama, bekerja pada
    // digit tanpa float - `inti/lib/format.ts`) bukan `Number(` JavaScript.
    expect(KODE).not.toMatch(/(?<![A-Za-z])Number\(|parseFloat|toFixed/)
  })
})

describe('Security Name: dropdown yang dapat difilter (keputusan work owner 30-09-2026)', () => {
  it('isian "Cari security" dibuang; satu PilihSaring', () => {
    expect(KODE).not.toMatch(/cariSecurity|Cari security|<datalist/)
    expect(KODE).toContain('<PilihSaring')
    expect(KODE).toContain('onCari={master.cari}')
  })
  it('membuka form lain membatalkan jawaban cari yang masih di jalan', () => {
    const buka = KODE.slice(KODE.indexOf('function buka'))
    expect(buka.slice(0, buka.indexOf('\n  }\n'))).toContain('master.reset()')
  })
})
