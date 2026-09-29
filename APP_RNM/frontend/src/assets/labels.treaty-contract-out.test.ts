// Penjaga bukti label Treaty Contract Out — tiket 02 →.
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN. Berkas ini
// membuka korpus XML pada baris yang tiap label sebut dan memastikan teksnya
// memang ada di sana. Korpus READ-ONLY — hanya dibaca.
//
// ⚠️ Bila korpus tidak terjangkau, test DILEWATI dengan pesan — bukan gagal.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { JENIS_REASURANSI_TCO, LAMPIRAN_TCO, MENU_TCO, TAHUN_TCO } from './labels.treaty-contract-out'

const KORPUS = 'D:\\XML\\RNM_BRD\\Treaty Contract Out'
const adaKorpus = existsSync(KORPUS)

/** Membaca satu baris (1-based) dari sebuah berkas korpus. */
function baris(relatif: string, nomor: number): string {
  const isi = readFileSync(join(KORPUS, relatif), 'utf8')
  return isi.split('\n')[nomor - 1] ?? ''
}

const GRID = 'Section\\InputTreatyContract.xml'
const FORM = 'Section\\InputDtlTreatyContact.xml'
const LAMPIRAN = 'Section\\GridTreatyArrangementAttachment.xml'

describe.skipIf(!adaKorpus)('label Treaty Contract Out berbukti barisnya', () => {
  const kasus: Array<[string, string, number, string]> = [
    ['Harness\\InboxTreatyContract.xml', 'pyLabel', 151, MENU_TCO.inboxTreatyContract],
    ['Harness\\InboxTreatyContractReinsType.xml', 'pyLabel', 151, MENU_TCO.inboxTreatyContractReinsType],
    ['Harness\\InboxTreatyContractDescription.xml', 'pyLabel', 359, MENU_TCO.inboxTreatyContractDescription],
    ['Section\\InputTreatyContractReinsType.xml', 'pyLabelFieldValue', 2652, JENIS_REASURANSI_TCO.reinsType],
    [GRID, 'pyLabelFieldValue', 7925, JENIS_REASURANSI_TCO.reinsuranceType],
    [GRID, 'pyLabelPreview', 7948, JENIS_REASURANSI_TCO.reinsuranceType],
    // tiket 03 - layar tahun treaty
    ['Section\\GridTreatyContract.xml', 'pyValue', 1057, TAHUN_TCO.judul],
    [FORM, 'pyValue', 5437, TAHUN_TCO.inputNewData],
    [GRID, 'pyValue', 17378, TAHUN_TCO.kolomUnderwritingYear],
    [GRID, 'pyValue', 17531, TAHUN_TCO.kolomTransactionYear],
    [GRID, 'pyValue', 17684, TAHUN_TCO.kolomStartDate],
    [GRID, 'pyValue', 17837, TAHUN_TCO.kolomEndDate],
    [GRID, 'pyValue', 17990, TAHUN_TCO.kolomTreatyGroup],
    [GRID, 'pyValue', 18136, TAHUN_TCO.kolomReinsuranceType],
    [GRID, 'pyLabel', 16387, TAHUN_TCO.add],
    [GRID, 'pyLabel', 19939, TAHUN_TCO.edit],
    [GRID, 'pyLabel', 20778, TAHUN_TCO.reinsType],
    [GRID, 'pyLabel', 22196, TAHUN_TCO.listDescription],
    [FORM, 'pyLabelFieldValue', 6379, TAHUN_TCO.formId],
    [FORM, 'pyLabelFieldValue', 6560, TAHUN_TCO.formTreatyGroup],
    [FORM, 'pyLabelFieldValue', 6800, TAHUN_TCO.formReinsuranceType],
    [FORM, 'pyLabelFieldValue', 7532, TAHUN_TCO.formStartDate],
    [FORM, 'pyLabelFieldValue', 7816, TAHUN_TCO.formEndDate],
    [FORM, 'pyLabelFieldValue', 8004, TAHUN_TCO.formUnderwritingYear],
    [FORM, 'pyLabelFieldValue', 8284, TAHUN_TCO.formTransactionYear],
    [FORM, 'pyLabelFieldValue', 9097, TAHUN_TCO.formModifiedDate],
    [FORM, 'pyLabelFieldValue', 9282, TAHUN_TCO.formUsername],
    [FORM, 'pyLabel', 10332, TAHUN_TCO.save],
    [FORM, 'pyLabel', 10622, TAHUN_TCO.cancel],
    // tiket 12 - panel lampiran
    [GRID, 'pyValue', 11721, LAMPIRAN_TCO.attachmentFor],
    [LAMPIRAN, 'pyValue', 1785, LAMPIRAN_TCO.forTreatyContractOut],
    [LAMPIRAN, 'pyLabel', 578, LAMPIRAN_TCO.addAttachment],
    [LAMPIRAN, 'pyLabel', 1023, LAMPIRAN_TCO.refresh],
    [LAMPIRAN, 'pyLabel', 2659, LAMPIRAN_TCO.downloadAll],
    [LAMPIRAN, 'pyValue', 3032, LAMPIRAN_TCO.kolomFileName],
    [LAMPIRAN, 'pyValue', 3170, LAMPIRAN_TCO.kolomType],
    [LAMPIRAN, 'pyLabel', 3897, LAMPIRAN_TCO.delete],
    ['Activity\\TreatyOutSaveAttachment.xml', 'PropertiesValue', 376, `"${LAMPIRAN_TCO.tanpaBerkas}"`],
  ]

  it.each(kasus)('%s baris %i memuat <%s>', (jalur, tag, nomor, teks) => {
    expect(baris(jalur, nomor).trim()).toBe(`<${tag}>${teks}</${tag}>`)
  })

  it('OQ-TCO-05: label tahun memang bersilang antara grid dan form', () => {
    // Grid: sel .UnderwritingYear (b18964) lalu .TreatyYear (b19125) - urutan
    // judul Underwriting Year (b17378) lalu Transaction Year (b17531).
    expect(baris(GRID, 18964).trim()).toBe('<pyValue>.UnderwritingYear</pyValue>')
    expect(baris(GRID, 19125).trim()).toBe('<pyValue>.TreatyYear</pyValue>')
    // Form: label Underwriting Year terikat TreatyYear, Transaction Year
    // terikat UnderwritingYear.
    expect(baris(FORM, 8035).trim()).toBe('<pyValue>InputTreatyYear.TreatyYear</pyValue>')
    expect(baris(FORM, 8315).trim()).toBe('<pyValue>InputTreatyYear.UnderwritingYear</pyValue>')
    // Dan kolom Reinsurance Type grid menampilkan .Proportion.
    expect(baris(GRID, 19724).trim()).toBe('<pyValue>.Proportion</pyValue>')
    expect(baris(FORM, 6829).trim()).toBe('<pyValue>InputTreatyYear.Proportion</pyValue>')
  })

  it('tombol Copy b20459 ada di korpus dan SENGAJA tidak dibawa (AC 72)', () => {
    expect(baris(GRID, 20459).trim()).toBe('<pyLabel>Copy</pyLabel>')
    expect(Object.values(TAHUN_TCO)).not.toContain('Copy')
  })

  it('tiket 12: panel lampiran disertakan form tahun, aksinya terbukti', () => {
    expect(baris(GRID, 13074).trim()).toBe('<pyInclude>GridTreatyArrangementAttachment</pyInclude>')
    expect(baris(LAMPIRAN, 1041).trim()).toBe('<pyActivity>LoadAttachmentTreatyOut</pyActivity>')
    expect(baris(LAMPIRAN, 2677).trim()).toBe('<pyActivity>TreatyOutDownloadAll_Act</pyActivity>')
    expect(baris(LAMPIRAN, 3488).trim()).toBe('<pyActivity>TreatyOutDownloadOne</pyActivity>')
    expect(baris(LAMPIRAN, 3915).trim()).toBe('<pyActivity>DeleteAttachmentTreaty</pyActivity>')
    // `Download` b2391 ada di korpus dan sengaja tidak dibawa sebagai tombol kedua.
    expect(baris(LAMPIRAN, 2391).trim()).toBe('<pyLabel>Download</pyLabel>')
    expect(Object.values(LAMPIRAN_TCO)).not.toContain('Download')
  })

  it('nama kelompok adalah nama folder korpus', () => {
    expect(existsSync(join('D:\\XML\\RNM_BRD', MENU_TCO.kelompok))).toBe(true)
  })

  it('tiga harness portal memang berkelas Data-Portal', () => {
    for (const h of ['InboxTreatyContract', 'InboxTreatyContractReinsType', 'InboxTreatyContractDescription']) {
      const isi = readFileSync(join(KORPUS, 'Harness', `${h}.xml`), 'utf8')
      expect(isi, h).toContain('<pyClassName>Data-Portal</pyClassName>')
    }
  })
})

describe('label yang tidak bergantung korpus', () => {
  it('nol kata Old / testing di label (penyimpangan sadar 8)', () => {
    for (const nilai of [...Object.values(MENU_TCO), ...Object.values(JENIS_REASURANSI_TCO), ...Object.values(TAHUN_TCO)]) {
      expect(nilai).not.toMatch(/\bOld\b|testing/i)
    }
  })
  it('yang tidak ada di korpus ditandai begitu', () => {
    const sumber = readFileSync(join(__dirname, 'labels.treaty-contract-out.ts'), 'utf8')
    for (const kunci of ['masterKosong', 'belumDipilih', 'kosong:', 'menungguTiket', 'catatanLabelBersilang']) {
      const i = sumber.indexOf(kunci)
      expect(i, kunci).toBeGreaterThan(0)
      expect(sumber.slice(sumber.lastIndexOf('/**', i), i)).toContain('[tidak ada di korpus]')
    }
  })
})
