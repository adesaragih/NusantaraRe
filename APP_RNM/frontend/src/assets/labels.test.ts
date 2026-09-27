// Penjaga bukti label — F0.1.
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN. Berkas ini
// membuka korpus XML yang tiap label sebut, pada baris yang tiap label sebut,
// dan memastikan teksnya memang ada di sana.
//
// ⚠️ Korpus `D:\XML\RNM_BRD\` READ-ONLY - dibaca sebagai bukti, tidak pernah
// ditulis. Test ini hanya membaca.
//
// ⚠️ Bila korpusnya tidak terjangkau (mesin lain, CI tanpa korpus), test ini
// DILEWATI dengan pesan - bukan gagal. Menggagalkannya membuat orang
// mematikannya, dan penjaga yang dimatikan tidak menjaga apa pun.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { LAYAR, PERAN, REGISTER, TAHAP, TOMBOL } from './labels'

const KORPUS = 'D:\\XML\\RNM_BRD\\Claim Life'
const adaKorpus = existsSync(KORPUS)

/** Membaca satu baris (1-based) dari sebuah berkas korpus. */
function baris(relatif: string, nomor: number): string {
  const isi = readFileSync(join(KORPUS, relatif), 'utf8')
  return isi.split('\n')[nomor - 1] ?? ''
}

/** Seluruh isi sebuah berkas korpus. */
function berkas(relatif: string): string {
  return readFileSync(join(KORPUS, relatif), 'utf8')
}

describe.skipIf(!adaKorpus)('label membawa bukti XML yang benar', () => {
  it('judul tahap ada sebagai <pyTaskName> di Register_Flow', () => {
    const alur = berkas('Flow\\Register_Flow.xml')
    for (const nilai of Object.values(TAHAP)) {
      expect(alur, `tahap ${nilai}`).toContain(`<pyTaskName>${nilai}</pyTaskName>`)
    }
  })

  // Tiap tombol: berkas, baris, dan teks yang harus ada di baris itu.
  const tombol: Array<[string, string, number, string]> = [
    ['simpanAdjustment', 'Section\\ClaimLifeDetailGCNM.xml', 22590, TOMBOL.simpanAdjustment],
    ['simpanKeOutstanding', 'Section\\AdjustmentDetail_Section.xml', 16169, TOMBOL.simpanKeOutstanding],
    ['tolakOutstanding', 'Section\\AdjustmentDetail_Section.xml', 15101, TOMBOL.tolakOutstanding],
    ['kirimBalikKeMedis', 'Section\\InputAkseptasiClaimLife.xml', 20467, TOMBOL.kirimBalikKeMedis],
    ['kirimBalikKeAdmin', 'Section\\InputAkseptasiClaimLife.xml', 20221, TOMBOL.kirimBalikKeAdmin],
    ['tutupKlaim', 'Section\\CloseClaim_Section.xml', 1028, TOMBOL.tutupKlaim],
  ]

  it.each(tombol)('tombol %s ada di %s baris %i', (_nama, jalur, nomor, teks) => {
    expect(baris(jalur, nomor)).toContain(teks)
  })

  it('tiap flow action yang dinamai punya berkasnya', () => {
    // Kunci LAYAR -> nama berkas FlowAction. Dipetakan eksplisit: nama
    // berkasnya TIDAK dapat diturunkan dari judulnya.
    const berkasLayar: Record<keyof typeof LAYAR, string> = {
      register: 'InputRegisterClaimLife',
      outstanding: 'OSClaimLife',
      adjustmentDetail: 'Adjustment_Detail',
      retro: 'RetroClaimLife',
      lampiran: 'AttachDocumentLife',
      unggahCSV: 'UploadCSV_ClaimLife',
      hapusLampiran: 'ConfirmDeleteAttachment',
      medis: 'MedicalCheck',
      kirimKeMedis: 'SendtoMedical',
      kirimKeAdmin: 'SendtoAdmin',
      akseptasi: 'AkseptasiClaimLife',
      tolakOS: 'RejectOSClaimLife',
      tutupKlaim: 'CloseClaim',
      ubahTanggal: 'ShowEditClaimLife',
      detailKlaim: 'ViewClaimDetailLifeGCNM',
      detailPolis: 'PL_DetailAction_ViewPolis',
    }
    // ⛔ Enam belas, persis cacah flow action korpus. Bila daftarnya berubah,
    // sensus PARITAS-LAYAR-DAN-AKSI.md ikut usang.
    expect(Object.keys(berkasLayar)).toHaveLength(16)
    for (const nama of Object.values(berkasLayar)) {
      expect(
        existsSync(join(KORPUS, 'FlowAction', `${nama}.xml`)),
        `FlowAction\\${nama}.xml`,
      ).toBe(true)
    }
  })
})

describe('label yang tidak bergantung korpus', () => {
  it('ketiga peran memakai pengenal ADR-U-0002, bukan sebutan layar', () => {
    // ⛔ Nilainya dikirim sebagai `X-Peran`; mengubahnya mengubah WEWENANG.
    expect(PERAN.admin).toBe('ReasLifeAdmin')
    expect(PERAN.medis).toBe('ReasLifeMedicalAdvisor')
    expect(PERAN.spv).toBe('ReasLifeSPV')
  })

  it('judul tahap tidak diterjemahkan di tempat aslinya', () => {
    // Terjemahan berdiri DI SAMPING (TAHAP_ID), bukan menggantikan.
    expect(TAHAP.medicalCheck).toBe('Medical Check')
    expect(TAHAP.claimAnalis).toBe('Claim Analis')
  })
})

describe.skipIf(!adaKorpus)('label layar Register berbukti barisnya', () => {
  // Tiap baris: label, berkas, nomor baris yang `labels.ts` sebut.
  const medan: Array<[string, number]> = [
    ['Choose Policy No', 3776],
    ['Find Insured', 7057],
    ['Type', 9104],
    ['Marketing Officer', 9890],
    ['Ceding', 11541],
    ['Policy Holder', 11736],
    ['Class of Business', 12171],
    ['Date Received Email', 13384],
    ['Response Date', 13590],
    ['Confirmation Date', 13795],
    ['Status', 14002],
    ['Updated Status', 14406],
    ['Realization Date', 14601],
    ['Name of Insured', 16064],
    ['Select Insured', 20008],
  ]

  it.each(medan)('%s ada di InputRegisterClaimLife.xml baris %i', (label, nomor) => {
    expect(baris('Section/InputRegisterClaimLife.xml', nomor)).toContain(
      `<pyLabelPreview>${label}</pyLabelPreview>`,
    )
  })

  it('REGISTER memuat kelima belas label itu, tidak kurang', () => {
    const nilai = Object.values(REGISTER)
    for (const [label] of medan) {
      expect(nilai).toContain(label)
    }
    // ⛔ Cacahnya dikunci: medan yang DIHILANGKAN dari layar sama merusaknya
    // dengan medan yang dikarang, dan yang pertama tidak berbunyi.
    expect(nilai).toHaveLength(medan.length)
  })
})
