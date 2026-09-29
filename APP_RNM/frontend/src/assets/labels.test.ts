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

import {
  DETAIL,
  EDIT_DATE,
  KONFIRMASI_BALIK,
  LAYAR,
  PERAN,
  REGISTER,
  REJECT_OS,
  TAHAP,
  TOMBOL,
  TOMBOL_KOMITE,
} from './labels.claimlife'

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
    ['cabutPeserta', 'Section\\InputOSClaimLife.xml', 17865, TOMBOL.cabutPeserta],
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
    ['System Reinsurance', 9301],
    ['Premium Method', 9694],
    ['Marketing Officer', 9890],
    ['WPC', 10149],
    ['Retro Name', 10335],
    ['Security Reinsurer', 10617],
    ['Ceding', 11541],
    ['Policy Holder', 11736],
    ['SOB', 11930],
    ['Class of Business', 12171],
    ['Product Name ID', 12450],
    ['Product Name', 12648],
    ['Date Received Email', 13384],
    ['Response Date', 13590],
    ['Confirmation Date', 13795],
    ['Status', 14002],
    ['Confirmation Reserved', 14198],
    ['Updated Status', 14406],
    ['Realization Date', 14601],
    ['Underwriter Note', 14806],
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
    // Butir bo: tombol grid `.AdjustmentList` b17126. (`Delete` b19120 tidak
    // dirender - OQ-N7 ditutup 29-09-2026, ADR-U-0031.)
    ['Add', 17937, 'pyLabel'],
    // Butir bk: sel read-only `.MAXCLAIM_RECEIVED` b12131.
    ['MAX CLAIM RECEIVED', 12124, 'pyLabelPreview'],
    // ⛔ ENAM total, bukan lima. Yang ini sempat luput karena
    // pencacahannya memakai rujukan `CheckTotalAdjustmentClaim`, dan ia
    // satu-satunya total yang TIDAK punya aksi refresh.
    ['Total Ceding Retention', 20629, 'pyLabelPreview'],
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

  it('tombol Komite VERBATIM dari ClaimComite, bukan karangan', () => {
    // ⛔ RALAT: layar sempat memakai `Send ke Komite`, campuran
    // Indonesia-Inggris yang tidak ada di korpus mana pun. Label tombol
    // diambil dari XML - pemakai sistem lama mencari kalimat yang sama, dan
    // pengujinya membandingkan kedua layar kata demi kata.
    expect(baris('Section\\ClaimComite.xml', 7033)).toBe(
      `<pyLabel>${TOMBOL_KOMITE.serahkan}</pyLabel>`,
    )
    expect(baris('Section\\ClaimComite.xml', 7715)).toBe(
      `<pyLabel>${TOMBOL_KOMITE.batal}</pyLabel>`,
    )
  })

  // Dialog Edit Date - label di berkas section-nya sendiri.
  const medanEditDate: Array<[string, number, string]> = [
    [EDIT_DATE.terimaKlaim, 1069, 'pyLabelPreview'],
    [EDIT_DATE.dokumenLengkap, 1381, 'pyLabelPreview'],
    [EDIT_DATE.konfirmasi, 1619, 'pyLabelPreview'],
    [EDIT_DATE.simpan, 1910, 'pyLabel'],
  ]

  it('konfirmasi jalur balik VERBATIM dari kedua section-nya', () => {
    expect(baris('Section/SendtoAdmin_Section.xml', 566)).toBe(
      `<pyValue>${KONFIRMASI_BALIK.keAdmin}</pyValue>`,
    )
    expect(baris('Section/SendtoAdmin_Section.xml', 1229)).toBe(
      `<pyLabel>${KONFIRMASI_BALIK.kirim}</pyLabel>`,
    )
    expect(baris('Section/SendtoMedical_Section.xml', 577)).toBe(
      `<pyValue>${KONFIRMASI_BALIK.keMedis}</pyValue>`,
    )
    expect(baris('Section/SendtoMedical_Section.xml', 1267)).toBe(
      `<pyLabel>${KONFIRMASI_BALIK.kirim}</pyLabel>`,
    )
    // Tombol batal modal: `pyCancelLabel` flow action-nya, bukan karangan.
    expect(baris('FlowAction/SendtoAdmin.xml', 19)).toBe(
      `<pyCancelLabel>${KONFIRMASI_BALIK.batal}</pyCancelLabel>`,
    )
    expect(baris('FlowAction/SendtoMedical.xml', 21)).toBe(
      `<pyCancelLabel>${KONFIRMASI_BALIK.batal}</pyCancelLabel>`,
    )
  })

  it('tombol jalur balik membuka local action, tombol maju tidak', () => {
    // ⛔ Yang menentukan ada-tidaknya konfirmasi: `pyLocalAction` sesudah
    // label tombolnya. Dibaca di tiga section, bukan ditebak dari namanya.
    const la = (n: string) => `<pyLocalAction>${n}</pyLocalAction>`
    expect(baris('Section/InputOSClaimLife.xml', 21433)).toBe(la('SendtoAdmin'))
    expect(baris('Section/MedicalCheckClaimLife.xml', 20285)).toBe(la('SendtoAdmin'))
    expect(baris('Section/InputAkseptasiClaimLife.xml', 20250)).toBe(la('SendtoAdmin'))
    expect(baris('Section/InputAkseptasiClaimLife.xml', 20496)).toBe(la('SendtoMedical'))
    const act = '<pyActivity>SendtoAdmin_Act1</pyActivity>'
    expect(baris('Section/InputOSClaimLife.xml', 21863)).toBe(act)
    expect(baris('Section/MedicalCheckClaimLife.xml', 21174)).toBe(act)
  })

  // Dialog Reject Outstanding (OQ-M5, GILIRAN-17) - label di section-nya
  // sendiri; batalnya `pyCancelLabel` flow action-nya.
  it('dialog Reject Outstanding VERBATIM dari RejectOSClaimLife_Sec', () => {
    const s = 'Section/RejectOSClaimLife_Sec.xml'
    expect(baris(s, 783)).toBe(`<pyLabelPreview>${REJECT_OS.tanggal}</pyLabelPreview>`)
    expect(baris(s, 969)).toBe(`<pyLabelPreview>${REJECT_OS.pic}</pyLabelPreview>`)
    expect(baris(s, 1680)).toBe(`<pyLabelPreview>${REJECT_OS.alasan}</pyLabelPreview>`)
    expect(baris(s, 3098)).toBe(`<pyLabel>${REJECT_OS.kirim}</pyLabel>`)
    expect(baris('FlowAction/RejectOSClaimLife.xml', 19)).toBe(
      `<pyCancelLabel>${REJECT_OS.batal}</pyCancelLabel>`,
    )
  })

  it.each(medanEditDate)('%s ada di EditDateClaimLife_Section baris %i sebagai <%s>',
    (label, nomor, tag) => {
      expect(baris('Section/EditDateClaimLife_Section.xml', nomor)).toBe(
        `<${tag}>${label}</${tag}>`,
      )
    })

  it('DETAIL memuat setiap label itu, tidak kurang dan tidak lebih', () => {
    const nilai = Object.values(DETAIL)
    for (const [label] of [...medanDetail, ...medanTutup]) {
      expect(nilai).toContain(label)
    }
    expect(nilai).toHaveLength(medanDetail.length + medanTutup.length)
  })

  it('rule CheckTotalAdjustmentClaim dirujuk section tetapi NOL berkasnya', () => {
    // ⚠️ DIPERTAHANKAN, TETAPI ARTINYA DIRALAT 27-09-2026.
    //
    // Fakta di bawah tetap benar: sepuluh rujukan, nol berkas rule. Yang
    // keliru adalah KESIMPULAN yang pernah digantungkan padanya - bahwa
    // angka totalnya karena itu tidak dapat ditiru. Yang hilang hanya
    // pemanggil REFRESH; nilainya dihitung activity lain yang ADA (lihat
    // uji berikutnya). Rujukan menggantung ini tetap OQ-H, dengan
    // pertanyaan yang lebih sempit: kenapa rule refresh-nya tidak ikut.
    const section = berkas('Section/ClaimLifeDetailGCNM.xml')
    const rujukan = section.match(/CheckTotalAdjustmentClaim/g) ?? []
    expect(rujukan).toHaveLength(10)
    expect(existsSync(join(KORPUS, 'Activity', 'CheckTotalAdjustmentClaim.xml'))).toBe(
      false,
    )
  })

  it('keenam total DIHITUNG oleh dua activity yang ada di korpus', () => {
    // ⛔ Uji yang menutup ralat itu. Ia memeriksa rumusnya di SUMBER,
    // bukan di kode kita - sehingga ia akan gagal bila kelak seseorang
    // menyalin rumus yang berbeda dari yang XML tulis.
    //
    // `SavePesertaClaim` langkah 8.1 b4221 dan `SaveOutStandingLife_Act`
    // langkah 23.1 b10841 sama-sama mengulang `.AdjustmentList` dan
    // menjumlahkan keenam kolom ke penampung `local.Total*`.
    for (const nama of ['SavePesertaClaim', 'SaveOutStandingLife_Act']) {
      const act = berkas(`Activity/${nama}.xml`)
      for (const kolom of [
        'CEDING_RETENTION',
        'SHARE_NUSANTARA_RE',
        'SUM_INSURED',
        'SUM_REASURED',
        'SHARE_RETRO',
        'CLAIM_AMOUNT',
      ]) {
        expect(
          act,
          `${nama} harus menjumlahkan ${kolom} ke penampungnya`,
        ).toContain(`.${kolom} + local.Total`)
      }
      // Dan penampungnya di-reset ke literal 0, sehingga peserta tanpa
      // baris adjustment bertotal 0 - bukan kosong.
      expect(act).toContain('<PropertiesName>local.TotalClaimAmount</PropertiesName>')
    }
  })

  it('langkah penjumlah itu BERPRASYARAT KOSONG - baris ditolak ikut', () => {
    // ⛔ Ini pembenar satu-satunya untuk "seluruh baris, termasuk yang
    // ditolak". Bila kelak ternyata ada prasyarat `STS_REJECT`, uji ini
    // gagal dan keputusan bc harus ditinjau ulang - bukan kodenya diam-diam
    // disesuaikan.
    const act = berkas('Activity/SavePesertaClaim.xml')
    const mulai = act.indexOf('RH_1.pySteps(8).pySteps(1)')
    const akhir = act.indexOf('RH_1.pySteps(8).pySteps(2)')
    expect(mulai).toBeGreaterThan(0)
    expect(akhir).toBeGreaterThan(mulai)
    const langkah = act.slice(mulai, akhir)
    expect(langkah).toContain('<pyStepsObjectName>.AdjustmentList</pyStepsObjectName>')
    expect(langkah).toContain('<pyStepsPreCondParamsWhen/>')
    expect(langkah).not.toContain('STS_REJECT')
  })
})
