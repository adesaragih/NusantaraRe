// ⛔ Laporan pemakai 8 Oktober 2026: Currency 100% Limit yang SUDAH tersimpan
// tampil "Choose" lagi sesudah Save, dan Save berikutnya menulis data basi.
//
// Sebabnya `useProperti` membekukan nilai awal tab saat tab PERTAMA dirender.
// Muat ulang sesudah Save mengosongkan penampung — tab yang tidak dilahirkan
// ulang lalu menyemai ulang data SEBELUM Save. Uji ini menjaga dua hal:
// fieldset tab berkunci `generasi`, dan `generasi` naik SESUDAH data baru
// diterapkan (bukan sebelum, saat `warisan` masih yang lama).

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

const FORM = readFileSync(join(__dirname, 'pages', 'FormKontrakTreatyIn.tsx'), 'utf8').replace(/\r\n/g, '\n')

describe('muat ulang sesudah Save melahirkan ulang tab', () => {
  it('key fieldset tab memuat generasi data', () => {
    expect(FORM).toContain("key={`${idKontrak}|${warisan === null ? '-' : 'isi'}|${mode}|${generasi}`}")
  })

  it('generasi naik SESUDAH setWarisan, pengosongan, dan penyemaian penampung', () => {
    const naik = FORM.indexOf('setGenerasi((n) => n + 1)')
    expect(naik).toBeGreaterThan(-1)
    const blok = FORM.slice(FORM.lastIndexOf('setWarisan(k)', naik), naik)
    expect(blok).toContain('penampung.kosongkan()')
    expect(blok).toContain('k.penampung ??')
    expect(blok).toContain('k.penampungLarik ??')
  })
})
