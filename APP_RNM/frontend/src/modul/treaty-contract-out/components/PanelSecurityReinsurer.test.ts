// Uji panel security — tiket 06 Treaty Contract Out.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import {
  formSecurityDari,
  formSecurityKosong,
  kataCariSecurity,
  keMasukSecurity,
  labelPilihanSecurity,
  pilihanDariTeks,
} from './PanelSecurityReinsurer'

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
    expect(KODE).toContain('cariReinsurerMaster(kataCariSecurity(teks))')
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
  const master = [
    { id: '10001', clientName: 'UJI REAS SATU', clientId: '' },
    { id: '10002', clientName: 'UJI REAS SATU', clientId: '' },
    { id: '10003', clientName: '', clientId: '' },
  ]
  it('isian "Cari security" dibuang; satu isian bertautan datalist', () => {
    expect(KODE).not.toMatch(/cariSecurity|Cari security/)
    expect(KODE).toContain('list={idDaftar}')
    expect(KODE).toContain('<datalist id={idDaftar}>')
  })
  it('label pilihan menyebut ID — nama kembar tetap dapat dibedakan', () => {
    expect(labelPilihanSecurity(master[0]!)).toBe('UJI REAS SATU (10001)')
    expect(labelPilihanSecurity(master[1]!)).not.toBe(labelPilihanSecurity(master[0]!))
    expect(labelPilihanSecurity(master[2]!)).toBe('10003')
  })
  it('hanya teks yang TEPAT label pilihan yang memilih security', () => {
    expect(pilihanDariTeks('UJI REAS SATU (10002)', master)?.id).toBe('10002')
    expect(pilihanDariTeks('UJI REAS', master)).toBeUndefined()
  })
  it('kata cari ke server tanpa ekor " (ID)"', () => {
    expect(kataCariSecurity('UJI REAS SATU (10001)')).toBe('UJI REAS SATU')
    expect(kataCariSecurity('UJI REAS SATU (100')).toBe('UJI REAS SATU')
    expect(kataCariSecurity('  asuransi ')).toBe('asuransi')
  })
})

