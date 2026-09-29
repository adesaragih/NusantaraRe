// Uji layar Outstanding — A3 kelompok Outstanding, butir aw.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { OUTSTANDING, TOMBOL_OS } from '../labels'
import { TAHAP_JALUR } from '../api'

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
    [TOMBOL_OS.simpanRNM, 21102],
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
  it('perpindahan memakai PanelPindahTahap tahap Outstanding, bukan salinan', () => {
    // ⛔ GILIRAN-11: salinan lokal melewatkan konfirmasi `Send Back to Admin?`.
    // Isi daftarnya (dua tombol, tujuan, konfirmasi) dikunci
    // `PanelPindahTahap.test.ts` — satu sumber.
    expect(SUMBER).toContain('<PanelPindahTahap')
    expect(SUMBER).toContain('tahap={TAHAP.outstandingClaim}')
    expect(SUMBER).not.toContain('const PERPINDAHAN')
  })

  it('klaim yang gagal dimuat dinyatakan, bukan ditelan', () => {
    // ⛔ Dulu satu `catch` menelan galat klaim DAN polis: layar kosong tanpa sebab.
    expect(SUMBER).toContain('setGalatMuat(e)')
    expect(SUMBER).toContain('<Gagal galat={galatMuat} />')
  })

  it('tombol hanya bagi kasus yang SEBENARNYA di Outstanding Claim', () => {
    // ⛔ Kotak masuk membuka layar ini untuk tahap mana pun (App.tsx).
    expect(SUMBER).toContain('setTahapKasus(k.tahap)')
    expect(SUMBER.match(/tahapKasus === TAHAP\.outstandingClaim &&/g) ?? []).toHaveLength(2)
  })

  it('Save to RNM memanggil simpanKeRNM', () => {
    expect(SUMBER).toContain('simpanKeRNM(klaimID)')
    expect(SUMBER).toContain('TOMBOL_OS.simpanRNM')
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
