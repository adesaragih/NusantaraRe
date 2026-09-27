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

import { DETAIL, LAYAR, PERAN, REGISTER, TAHAP, TOMBOL } from './labels'

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

  // Dua label kotak pencarian `Find Insured`. Keduanya BUKAN
  // `pyLabelPreview`, dan itu sebabnya ia berdiri terpisah dari `medan`:
  // memaksanya masuk daftar di atas berarti menguji tag yang tidak ada dan
  // lalu melonggarkan ujinya sampai lulus. Diperiksa: nol
  // `<pyLabelPreview>Certificate No` di seluruh berkas.
  const medanCari: Array<[string, number, string]> = [
    ['Certificate No', 16277, 'pyLabelFieldValue'],
    ['Search', 16553, 'pyLabel'],
  ]

  it.each(medanCari)('%s ada di baris %i sebagai <%s>', (label, nomor, tag) => {
    expect(baris('Section/InputRegisterClaimLife.xml', nomor)).toBe(
      `<${tag}>${label}</${tag}>`,
    )
  })

  it('REGISTER memuat ketujuh belas label itu, tidak kurang', () => {
    const nilai = Object.values(REGISTER)
    for (const [label] of medan) {
      expect(nilai).toContain(label)
    }
    for (const [label] of medanCari) {
      expect(nilai).toContain(label)
    }
    // ⛔ Cacahnya dikunci: medan yang DIHILANGKAN dari layar sama merusaknya
    // dengan medan yang dikarang, dan yang pertama tidak berbunyi.
    expect(nilai).toHaveLength(medan.length + medanCari.length)
  })
})

describe.skipIf(!adaKorpus)('label layar Detail berbukti barisnya', () => {
  // Tiap label disebut BERSAMA TAG-nya: satu teks dapat berdiri di lebih dari
  // satu baris dengan tag berbeda (`Name of Insured` b16039/b16064 di layar
  // Register), dan menyebut angkanya saja membuat dua kutipan yang sah tampak
  // bertentangan.
  const medanDetail: Array<[string, number, string]> = [
    ['Find Disease', 5061, 'pyLabel'],
    ['Edit Date', 14115, 'pyLabel'],
    ['Save Adjustment', 22641, 'pyLabel'],
    ['Total Share Nusantara Re', 20914, 'pyLabelPreview'],
    ['Total Sum Insured', 21201, 'pyLabelPreview'],
    ['Total Sum Reasured', 21488, 'pyLabelPreview'],
    ['Total Share Retro', 21775, 'pyLabelPreview'],
    ['Total Claim Amount', 22063, 'pyLabelPreview'],
  ]

  it.each(medanDetail)('%s ada di ClaimLifeDetailGCNM baris %i sebagai <%s>',
    (label, nomor, tag) => {
      expect(baris('Section/ClaimLifeDetailGCNM.xml', nomor)).toBe(
        `<${tag}>${label}</${tag}>`,
      )
    })

  // Dua label ini milik `CloseClaim_Section.xml`, bukan section Detail, jadi
  // buktinya dicari di berkas yang benar - bukan di berkas yang kebetulan
  // memuat teks yang sama.
  const medanTutup: Array<[string, number, string]> = [
    ['Close Claim', 1081, 'pyLabel'],
    ['Are you sure want to Close Claim?', 499, 'pyValue'],
  ]

  it.each(medanTutup)('%s ada di CloseClaim_Section baris %i sebagai <%s>',
    (label, nomor, tag) => {
      expect(baris('Section/CloseClaim_Section.xml', nomor)).toBe(
        `<${tag}>${label}</${tag}>`,
      )
    })

  it('tombol Close Claim menjalankan DUA aksi, bukan satu', () => {
    // ⛔ Ronde pertama melaporkan "tombol itu tidak punya aksi lain selain
    // gerbangnya". Keliru: satu klik menjalankan `refresh` -> activity DAN
    // `closeContainer`. Uji ini menahan klaim itu supaya tidak diulang.
    const tutup = berkas('Section/CloseClaim_Section.xml')
    expect(tutup).toContain('<pyActivity>ProtectCloseClaim_act</pyActivity>')
    expect(tutup).toContain('<pyAction>closeContainer</pyAction>')
  })

  it('DETAIL memuat kesepuluh label itu, tidak kurang', () => {
    const nilai = Object.values(DETAIL)
    for (const [label] of [...medanDetail, ...medanTutup]) {
      expect(nilai).toContain(label)
    }
    expect(nilai).toHaveLength(medanDetail.length + medanTutup.length)
  })

  it('rule CheckTotalAdjustmentClaim dirujuk section tetapi NOL berkasnya', () => {
    // ⛔ Inti keputusan panel total. Uji ini memeriksa KEDUA sisinya:
    // rujukannya memang ada (jadi medannya nyata dan harus tampil), dan
    // rule-nya memang tidak ada (jadi angkanya tidak boleh dikarang).
    //
    // Bila kelak ekspornya dilengkapi, uji inilah yang gagal - dan gagalnya
    // adalah kabar baik: aturannya sudah dapat ditiru.
    const section = berkas('Section/ClaimLifeDetailGCNM.xml')
    const rujukan = section.match(/CheckTotalAdjustmentClaim/g) ?? []
    expect(rujukan).toHaveLength(10)
    expect(existsSync(join(KORPUS, 'Activity', 'CheckTotalAdjustmentClaim.xml'))).toBe(
      false,
    )
  })
})
