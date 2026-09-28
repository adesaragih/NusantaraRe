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

import { JENIS_REASURANSI_TCO, MENU_TCO } from './labels.treaty-contract-out'

const KORPUS = 'D:\\XML\\RNM_BRD\\Treaty Contract Out'
const adaKorpus = existsSync(KORPUS)

/** Membaca satu baris (1-based) dari sebuah berkas korpus. */
function baris(relatif: string, nomor: number): string {
  const isi = readFileSync(join(KORPUS, relatif), 'utf8')
  return isi.split('\n')[nomor - 1] ?? ''
}

describe.skipIf(!adaKorpus)('label Treaty Contract Out berbukti barisnya', () => {
  const kasus: Array<[string, string, number, string]> = [
    ['Harness\\InboxTreatyContract.xml', 'pyLabel', 151, MENU_TCO.inboxTreatyContract],
    ['Harness\\InboxTreatyContractReinsType.xml', 'pyLabel', 151, MENU_TCO.inboxTreatyContractReinsType],
    ['Harness\\InboxTreatyContractDescription.xml', 'pyLabel', 359, MENU_TCO.inboxTreatyContractDescription],
    ['Section\\InputTreatyContractReinsType.xml', 'pyLabelFieldValue', 2652, JENIS_REASURANSI_TCO.reinsType],
    ['Section\\InputTreatyContract.xml', 'pyLabelFieldValue', 7925, JENIS_REASURANSI_TCO.reinsuranceType],
    ['Section\\InputTreatyContract.xml', 'pyLabelPreview', 7948, JENIS_REASURANSI_TCO.reinsuranceType],
  ]

  it.each(kasus)('%s baris %i memuat <%s>', (jalur, tag, nomor, teks) => {
    expect(baris(jalur, nomor).trim()).toBe(`<${tag}>${teks}</${tag}>`)
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
    for (const nilai of [...Object.values(MENU_TCO), ...Object.values(JENIS_REASURANSI_TCO)]) {
      expect(nilai).not.toMatch(/\bOld\b|testing/i)
    }
  })
  it('yang tidak ada di korpus ditandai begitu', () => {
    const sumber = readFileSync(join(__dirname, 'labels.treaty-contract-out.ts'), 'utf8')
    const blok = sumber.slice(sumber.indexOf('masterKosong'))
    expect(sumber.slice(0, sumber.indexOf('masterKosong'))).toContain('[tidak ada di korpus]')
    expect(blok).toContain('[tidak ada di korpus]')
  })
})
