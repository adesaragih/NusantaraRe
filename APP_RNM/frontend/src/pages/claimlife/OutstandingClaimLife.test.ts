// Uji layar Outstanding — A3 kelompok Outstanding, butir aw.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { OUTSTANDING, TOMBOL_OS } from '../../assets/labels.claimlife'
import { TAHAP_JALUR } from '../../services/api'

const KORPUS = 'D:/XML/RNM_BRD/Claim Life'
const adaKorpus = existsSync(KORPUS)
const SUMBER = readFileSync(join(__dirname, 'OutstandingClaimLife.tsx'), 'utf8')

/** Membaca satu baris (1-based) dari section Outstanding. */
function barisOS(nomor: number): string {
  const isi = readFileSync(join(KORPUS, 'Section', 'InputOSClaimLife.xml'), 'utf8')
  return isi.split('\n')[nomor - 1] ?? ''
}

describe.skipIf(!adaKorpus)('label Outstanding VERBATIM korpus', () => {
  const medan: Array<[string, number]> = [
    [OUTSTANDING.nomorKlaim, 951],
    [OUTSTANDING.type, 6898],
    [OUTSTANDING.marketing, 7686],
    [OUTSTANDING.ceding, 9407],
    [OUTSTANDING.pemegangPolis, 9602],
    [OUTSTANDING.kelasBisnis, 10036],
    [OUTSTANDING.tanggalEmail, 11250],
    [OUTSTANDING.tanggalRespon, 11456],
    [OUTSTANDING.tanggalKonfirmasi, 11662],
    [OUTSTANDING.status, 11868],
    [OUTSTANDING.statusDiperbarui, 12271],
    [OUTSTANDING.tanggalRealisasi, 12468],
  ]

  it.each(medan)('%s ada di InputOSClaimLife.xml baris %i', (label, nomor) => {
    expect(barisOS(nomor)).toContain(`<pyLabelPreview>${label}</pyLabelPreview>`)
  })

  it.each([
    [TOMBOL_OS.kembaliKeRegister, 21404],
    [TOMBOL_OS.kirimKeMedis, 21839],
  ])('tombol %s ada di baris %i', (label, nomor) => {
    expect(barisOS(nomor)).toContain(`<pyLabel>${label}</pyLabel>`)
  })

  it('Close Claim ada di baris 22750', () => {
    expect(barisOS(22750)).toContain(`<pyLabelPreview>${TOMBOL_OS.tutupKlaim}</pyLabelPreview>`)
  })
})

describe('butir aw — dua tombol perpindahan', () => {
  it('Send Back to Register menuju tahap input-register', () => {
    // `pyLocalAction SendtoAdmin` (21433) → Decision3 → Assignment2.
    const blok = SUMBER.slice(SUMBER.indexOf('const PERPINDAHAN'))
    const baris = blok.slice(0, blok.indexOf(']'))
    expect(baris).toContain(`TAHAP_JALUR.inputRegister`)
    expect(baris).toContain('kembaliKeRegister')
  })

  it('Send to Medical Check menuju tahap medical-check', () => {
    const blok = SUMBER.slice(SUMBER.indexOf('const PERPINDAHAN'))
    const baris = blok.slice(0, blok.indexOf(']'))
    expect(baris).toContain(`TAHAP_JALUR.medicalCheck`)
    expect(baris).toContain('kirimKeMedis')
  })

  it('cacah tombol perpindahan TEPAT dua', () => {
    const blok = SUMBER.slice(SUMBER.indexOf('const PERPINDAHAN'))
    const baris = blok.slice(0, blok.indexOf(']'))
    expect(baris.match(/tujuan: TAHAP_JALUR\./g) ?? []).toHaveLength(2)
  })

  it('cacatnya DICATAT di berkasnya, bukan disembunyikan', () => {
    // ⛔ Kedua activity berprasyarat ReasLifeMedicalAdvisor padahal layar ini
    // dipegang Admin. Yang ditiru MAKSUDnya, dan itu harus terbaca.
    expect(SUMBER).toContain('ReasLifeMedicalAdvisor')
    expect(SUMBER).toContain('OQ-C')
  })
})

describe('tahap sebagai KATA, bukan angka', () => {
  it('keempat jalur berupa kata berhubung', () => {
    expect(TAHAP_JALUR).toEqual({
      inputRegister: 'input-register',
      outstanding: 'outstanding',
      medicalCheck: 'medical-check',
      claimAnalis: 'claim-analis',
    })
  })
})

describe('panel polis TIDAK disalin', () => {
  it('Outstanding memakai PanelDataPolis yang sama dengan Register', () => {
    // Menyalinnya berarti dua bentuk yang harus berubah bersama, dan tidak
    // ada yang memaksanya.
    expect(SUMBER).toContain("import { PanelDataPolis }")
    expect(SUMBER).toContain('<PanelDataPolis')
  })
})
