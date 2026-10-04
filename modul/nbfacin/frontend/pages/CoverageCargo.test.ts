// Paritas layar CoverageCargo dengan section Pega - tiket 21. Tanpa perender React di
// repositori, uji ini membaca SUMBER halaman (pola modul lain): urutan medan = urutan
// sel Pega, sel 16 tidak tampil, tombol Pega nonaktif, hanya tiga medan yang dapat diisi.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const SUMBER = readFileSync(join(__dirname, 'CoverageCargo.tsx'), 'utf8').replace(/\r\n/g, '\n')
const awalJSX = SUMBER.indexOf('return (\n    <Panel')
if (awalJSX < 0) throw new Error('blok JSX CoverageCargo tidak ditemukan - uji ini perlu diperbarui')
const JSX = SUMBER.slice(awalJSX)

describe('CoverageCargo = blok pertama InputCoverageCargo_FacIn', () => {
  it('urutan medan dan tombol = urutan sel Pega (3,4,5,10,11,12,15,17,18,21,22,23,24)', () => {
    const urutan = [
      'MEDAN.keteranganJaminan', 'MEDAN.coverageInitial', 'TOMBOL.pilihCoverage', 'MEDAN.mataUang', 'MEDAN.rate',
      'MEDAN.limitOfLiability', 'MEDAN.currencyMaster', 'TOMBOL.ambilCurrencyMaster', 'MEDAN.diskonPersen',
      'MEDAN.tsi', 'MEDAN.premi', 'MEDAN.minPremi', 'MEDAN.diskon',
    ]
    // Dengan titik akhir: `MEDAN.diskon.` tidak boleh cocok dengan `MEDAN.diskonPersen`.
    const letak = urutan.map((u) => JSX.indexOf(u + '.'))
    expect(letak.every((l) => l >= 0)).toBe(true)
    expect([...letak].sort((a, b) => a - b)).toEqual(letak)
  })

  it('empat kelompok layout Pega = empat form-grid', () => {
    expect(JSX.match(/className="form-grid"/g)).toHaveLength(4)
  })

  it('hanya Name, Rate (%), TSI dapat diisi (butir 61); medan lain read-only', () => {
    // TSI memakai IsianUang (berformat ribuan saat diketik, permintaan work owner 03-10-2026).
    const dapatDiisi = [...JSX.matchAll(/<(?:Field|IsianUang) label=\{MEDAN\.(\w+)\.label\}[^/]*?onChange=\{ubah\(/g)].map((m) => m[1])
    expect(dapatDiisi).toEqual(['mataUang', 'rate', 'tsi'])
    const readOnly = [...JSX.matchAll(/<Field label=\{MEDAN\.(\w+)\.label\}[^/]*?readOnly/g)].map((m) => m[1])
    expect(readOnly).toEqual(['coverageInitial', 'limitOfLiability', 'currencyMaster', 'diskonPersen', 'premi', 'minPremi', 'diskon'])
  })

  it('tombol Pega yang belum diport nonaktif, dengan keterangan terlihat', () => {
    expect(SUMBER).toMatch(/<button[^>]*disabled[^>]*aria-describedby/)
    expect(SUMBER).toContain('<small>({TEKS.belumDiport})</small>')
  })

  it('frontend tidak menghitung uang: tanpa Number/parseFloat atas premi', () => {
    expect(SUMBER).not.toMatch(/parseFloat|Number\(|toFixed/)
  })
})
