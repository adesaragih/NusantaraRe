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

import {
  JENIS_REASURANSI_TCO,
  KLAUSUL_TCO,
  KONTRAK_TCO,
  HAPUS_TCO,
  KURS_TCO,
  LABEL_MEDAN_KHUSUS,
  JUDUL_GRID_KLAUSUL,
  KOLOM_GRID_KLAUSUL,
  LABEL_MEDAN_KLAUSUL,
  MENU_TCO,
  REINSURER_TCO,
  SECURITY_TCO,
  BUSINESS_TCO,
  TAHUN_TCO,
} from './labels'

const KORPUS = 'D:\\XML\\RNM_BRD\\Treaty Contract Out'
const adaKorpus = existsSync(KORPUS)

/** Membaca satu baris (1-based) dari sebuah berkas korpus. */
function baris(relatif: string, nomor: number): string {
  const isi = readFileSync(join(KORPUS, relatif), 'utf8')
  return isi.split('\n')[nomor - 1] ?? ''
}

const GRID = 'Section\\InputTreatyContract.xml'
const FORM = 'Section\\InputDtlTreatyContact.xml'
const KONTRAK = 'Section\\InputTreatyContractReinsType.xml'
const REAS = 'Section\\ViewDetailTreatyReinsurerGrid1.xml'
const BIZ = 'Section\\ViewDetailTreatyBusinessGrid.xml'

describe.skipIf(!adaKorpus)('label Treaty Contract Out berbukti barisnya', () => {
  const kasus: Array<[string, string, number, string]> = [
    ['Harness\\InboxTreatyContract.xml', 'pyLabel', 151, MENU_TCO.inboxTreatyContract],
    ['Harness\\InboxTreatyContractReinsType.xml', 'pyLabel', 151, MENU_TCO.inboxTreatyContractReinsType],
    ['Harness\\InboxTreatyContractDescription.xml', 'pyLabel', 359, MENU_TCO.inboxTreatyContractDescription],
    ['Section\\InputTreatyContractReinsType.xml', 'pyLabelFieldValue', 2652, JENIS_REASURANSI_TCO.reinsType],
    [GRID, 'pyLabelFieldValue', 7925, JENIS_REASURANSI_TCO.reinsuranceType],
    [GRID, 'pyLabelPreview', 7948, JENIS_REASURANSI_TCO.reinsuranceType],
    // tiket 03 - layar tahun treaty
    // Ditulis huruf kapital di awal kata (keputusan work owner 03-10-2026); korpus menulisnya kapital semua.
    ['Section\\GridTreatyContract.xml', 'pyValue', 1057, TAHUN_TCO.judul.toUpperCase()],
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
    // tiket 04 - editor kontrak
    ['Harness\\InboxTreatyContractReinsType.xml', 'pyValue', 1670, KONTRAK_TCO.judul],
    [KONTRAK, 'pyLabelFieldValue', 1145, KONTRAK_TCO.headerUnderwritingYear],
    [KONTRAK, 'pyLabelFieldValue', 1358, KONTRAK_TCO.headerReinsType],
    [KONTRAK, 'pyLabelFieldValue', 2905, KONTRAK_TCO.formStartDate],
    [KONTRAK, 'pyLabelFieldValue', 3244, KONTRAK_TCO.formEndDate],
    [KONTRAK, 'pyLabel', 3618, KONTRAK_TCO.save],
    [KONTRAK, 'pyLabelFieldValue', 4363, KONTRAK_TCO.formModifiedDate],
    [KONTRAK, 'pyLabelFieldValue', 4547, KONTRAK_TCO.formUsername],
    [KONTRAK, 'pyLabel', 5343, KONTRAK_TCO.undo],
    [KONTRAK, 'pyLabelFieldValue', 6400, KONTRAK_TCO.information],
    [KONTRAK, 'pyLabel', 8528, KONTRAK_TCO.add],
    [KONTRAK, 'pyValue', 9164, KONTRAK_TCO.kolomReinsType],
    [KONTRAK, 'pyValue', 9304, KONTRAK_TCO.kolomTreatyStart],
    [KONTRAK, 'pyValue', 9444, KONTRAK_TCO.kolomTreatyEnd],
    [KONTRAK, 'pyLabel', 10519, KONTRAK_TCO.edit],
    [KONTRAK, 'pyLabel', 10842, KONTRAK_TCO.businessList],
    [KONTRAK, 'pyLabel', 11308, KONTRAK_TCO.reinsurerList],
    [KONTRAK, 'pyLabel', 11809, KONTRAK_TCO.delete],
    // tiket 05 - panel reinsurer
    [REAS, 'pyLabel', 1730, REINSURER_TCO.add],
    [REAS, 'pyValue', 2418, REINSURER_TCO.kolomReinsId],
    [REAS, 'pyValue', 2560, REINSURER_TCO.kolomReinsurer],
    [REAS, 'pyValue', 2702, REINSURER_TCO.kolomShare],
    [REAS, 'pyValue', 2844, REINSURER_TCO.kolomComm],
    [REAS, 'pyValue', 2990, REINSURER_TCO.kolomRating],
    [REAS, 'pyValue', 3138, REINSURER_TCO.kolomOperatorName],
    [REAS, 'pyLabel', 4491, REINSURER_TCO.edit],
    [REAS, 'pyLabel', 4936, REINSURER_TCO.delete],
    [REAS, 'pyLabel', 5277, REINSURER_TCO.securityReinsurer],
    [REAS, 'pyLabelFieldValue', 7842, REINSURER_TCO.formId],
    [REAS, 'pyLabelFieldValue', 8042, REINSURER_TCO.formReinsId],
    [REAS, 'pyLabelFieldValue', 8226, REINSURER_TCO.formReinsurer],
    [REAS, 'pyLabelFieldValue', 8522, REINSURER_TCO.formShare],
    [REAS, 'pyLabelFieldValue', 8800, REINSURER_TCO.formComm],
    [REAS, 'pyLabelFieldValue', 9076, REINSURER_TCO.formRating],
    [REAS, 'pyLabelFieldValue', 11100, REINSURER_TCO.formOperatorName],
    [REAS, 'pyLabel', 11405, REINSURER_TCO.save],
    [REAS, 'pyLabelFieldValue', 12131, REINSURER_TCO.error],
    [REAS, 'pyLabelFieldValue', 12868, 'Informasi'],
    // tiket 07 - panel business
    [BIZ, 'pyLabelFieldValue', 1988, BUSINESS_TCO.judul],
    [BIZ, 'pyLabel', 2785, BUSINESS_TCO.add],
    [BIZ, 'pyValue', 3422, BUSINESS_TCO.kolomTreatyGroup],
    [BIZ, 'pyValue', 3566, BUSINESS_TCO.kolomBusinessId],
    [BIZ, 'pyValue', 3710, BUSINESS_TCO.kolomBusinessName],
    [BIZ, 'pyLabel', 4547, BUSINESS_TCO.edit],
    [BIZ, 'pyLabel', 4826, BUSINESS_TCO.delete],
    [BIZ, 'pyLabelFieldValue', 6241, BUSINESS_TCO.formBusinessName],
    [BIZ, 'pyLabelFieldValue', 6499, BUSINESS_TCO.formActive],
    [BIZ, 'pyLabelFieldValue', 6680, BUSINESS_TCO.formBusinessCode],
    [BIZ, 'pyLabel', 6966, BUSINESS_TCO.save],
    [BIZ, 'pyLabelFieldValue', 9319, BUSINESS_TCO.information],
    [BIZ, 'pyLabel', 10889, BUSINESS_TCO.closeList],
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

  it('tiket 04: label kepala ReinsType memang memuat nama grup, dan aksi tombolnya terbukti', () => {
    expect(baris(KONTRAK, 1386).trim()).toBe('<pyValue>InputTreatyContractReinsType.TreatyGroupName</pyValue>')
    expect(baris(KONTRAK, 3642).trim()).toBe('<pyActivity>SaveTreatyContract_Act</pyActivity>')
    expect(baris(KONTRAK, 5366).trim()).toBe('<pyActivity>UndoOperation</pyActivity>')
    expect(baris(KONTRAK, 8552).trim()).toBe('<pyActivity>NewInputTreatyContract_Act</pyActivity>')
    expect(baris(KONTRAK, 10543).trim()).toBe('<pyActivity>SetUbahTreatyContract</pyActivity>')
    expect(baris(KONTRAK, 9992).trim()).toBe('<pyValue>.ReinsTypeName</pyValue>')
    // Medan ID hanya berlabel bawaan kontrol - label "ID" kita tandai bukan dari korpus.
    expect(baris(KONTRAK, 2478).trim()).toBe('<pyLabelFieldValue>Formatted Text</pyLabelFieldValue>')
  })

  it('tiket 05: Total Share (ber-entitas di XML), aksi tombol, dan Tambah yang tidak dibawa', () => {
    expect(baris(REAS, 6186).trim().replace(/&gt;/g, '>')).toBe(`<pyValue>${REINSURER_TCO.totalShare}</pyValue>`)
    expect(baris(REAS, 1754).trim()).toBe('<pyActivity>NewTreatyReinsurerDetail_Act</pyActivity>')
    expect(baris(REAS, 4543).trim()).toBe('<pyActivity>SetUbahTreatyReinsurerList_Act</pyActivity>')
    expect(baris(REAS, 11429).trim()).toBe('<pyActivity>SaveTreatyReinsurerDetail1_Act</pyActivity>')
    expect(baris(KONTRAK, 11325).trim()).toBe('<pyActivity>BrowseTreatyReinsurerList_Act</pyActivity>')
    expect(baris(REAS, 1996).trim()).toBe('<pyLabel>Tambah</pyLabel>')
    expect(Object.values(REINSURER_TCO)).not.toContain('Tambah')
  })

  it('tiket 07: aksi tombol business terbukti; kolom Business ID = ID baris', () => {
    expect(baris(BIZ, 2809).trim()).toBe('<pyActivity>NewTreatyBusinessDetail_Act</pyActivity>')
    expect(baris(BIZ, 4571).trim()).toBe('<pyActivity>SetUbahTreatyBusinessList_Act</pyActivity>')
    expect(baris(BIZ, 4850).trim()).toBe('<pyActivity>DeleteRowBusiness</pyActivity>')
    expect(baris(BIZ, 6993).trim()).toBe('<pyActivity>SaveTreatyBusinessDetail_Act</pyActivity>')
    expect(baris(BIZ, 10914).trim()).toBe('<pyActivity>CancelActivity</pyActivity>')
    expect(baris(BIZ, 4228).trim()).toBe('<pyValue>.ID</pyValue>')
    expect(baris(KONTRAK, 10859).trim()).toBe('<pyActivity>BrowseTreatyBusinessList_Act</pyActivity>')
  })

  it('tiket 08: layar klausul — kepala, dua grid, tombol, dan label medan per jenis VERBATIM', () => {
    const H = join('Harness', 'InboxTreatyContractDescription.xml')
    const lf = (v: string) => `<pyLabelFieldValue>${v}</pyLabelFieldValue>`
    const nilai = (v: string) => `<pyValue>${v}</pyValue>`
    const tombol = (v: string) => `<pyLabel>${v}</pyLabel>`
    expect(baris(H, 359).trim()).toBe(tombol(MENU_TCO.inboxTreatyContractDescription))
    const kepala: [number, string][] = [
      [2003, KLAUSUL_TCO.headerTreatyGroupId], [2489, KLAUSUL_TCO.headerUnderwritingYear],
      [2665, KLAUSUL_TCO.headerTransactionYear], [2839, KLAUSUL_TCO.headerStartDate], [3025, KLAUSUL_TCO.headerEndDate],
      [3209, KLAUSUL_TCO.headerTreatyDescription], [3382, KLAUSUL_TCO.headerProportionType],
    ]
    for (const [n, v] of kepala) expect(baris(H, n).trim(), String(n)).toBe(lf(v))
    // b4880 "For Non XOL" diganti "Treaty Desc" dan grid XOL b7971 dibuang [keputusan work owner 30-09-2026].
    expect(baris(H, 4880).trim()).toBe(nilai('For Non XOL'))
    expect(KLAUSUL_TCO.gridTreatyDesc).toBe('Treaty Desc')
    for (const [n, v] of [[5440, KLAUSUL_TCO.kolomId],
      [5549, KLAUSUL_TCO.kolomDescriptionName]] as [number, string][]) {
      expect(baris(H, n).trim(), String(n)).toBe(nilai(v))
    }
    expect(baris(H, 6059).trim()).toBe(tombol(KLAUSUL_TCO.show))
    expect(baris(GRID, 22196).trim()).toBe(tombol(TAHUN_TCO.listDescription))
    const epi = join('Section', 'GridTreatyArrangementEpi.xml')
    for (const [n, v] of [[2777, LABEL_MEDAN_KLAUSUL.ReinsTypeID], [3215, LABEL_MEDAN_KLAUSUL.Line], [3372, LABEL_MEDAN_KLAUSUL.Rp],
      [3798, LABEL_MEDAN_KLAUSUL.Usd], [4487, KLAUSUL_TCO.formModifiedDate]] as [number, string][]) {
      expect(baris(epi, n).trim(), String(n)).toBe(lf(v))
    }
    for (const [n, v] of [[5501, KLAUSUL_TCO.save], [8980, KLAUSUL_TCO.add], [10917, KLAUSUL_TCO.edit],
      [11200, KLAUSUL_TCO.showChild]] as [number, string][]) {
      expect(baris(epi, n).trim(), String(n)).toBe(tombol(v))
    }
    const sec = (b: string) => join('Section', `GridTreatyArrangement${b}.xml`)
    const anak = join('Section', 'GridTreatyArrTreatyEpiList.xml')
    expect(baris(anak, 3029).trim()).toBe(lf(LABEL_MEDAN_KLAUSUL.Pct))
    expect(baris(anak, 8657).trim()).toBe(tombol(KLAUSUL_TCO.closeChild))
    const medan: [string, number, string][] = [
      ['BordereAux', 2702, LABEL_MEDAN_KLAUSUL.Method], ['TerrLimit', 2742, LABEL_MEDAN_KLAUSUL.TerritorialLimit],
      ['ProfitCommision', 3223, LABEL_MEDAN_KLAUSUL.PctMe], ['ProfitCommision', 3481, LABEL_MEDAN_KLAUSUL.Ydcf],
      ['Coins', 3394, LABEL_MEDAN_KLAUSUL.CoIns_Min], ['Coins', 4193, LABEL_MEDAN_KLAUSUL.CoIns_Max],
      ['Coins', 4880, LABEL_MEDAN_KLAUSUL.TreatyLimit],
      ['MinLOL', 1540, LABEL_MEDAN_KHUSUS.MinLOL.Pct], ['MInLOLMB', 1548, LABEL_MEDAN_KHUSUS.MinLOLMB.Pct],
      ['MaxCoinsPanel', 1525, LABEL_MEDAN_KHUSUS.MaxCoinsPanel.CoIns_Max],
      ['ExclutionTreatyOccupation', 2007, LABEL_MEDAN_KLAUSUL.ID_Occupation],
      ['ExclutionTreatyOccupation', 2296, `.${LABEL_MEDAN_KLAUSUL.Occupation}`],
      ['ExclutionTreatyOccupation', 2745, LABEL_MEDAN_KHUSUS['ExclutionTreaty/Occupation'].Line],
      ['ExclutionTreatyOccupation', 2935, LABEL_MEDAN_KHUSUS['ExclutionTreaty/Occupation'].Usd],
      ['ExclutionTreatyOccupation', 3222, LABEL_MEDAN_KHUSUS['ExclutionTreaty/Occupation'].Rp],
      ['ExclutionTreatyClausule', 2030, LABEL_MEDAN_KLAUSUL.ID_Clause],
      ['ExclutionTreatyClausule', 2320, `.${LABEL_MEDAN_KLAUSUL.Clause}`],
      ['ExclutionTreatyPeriode', 500, 'Max Periode (Month)'], // label kita: 'Max Period (Month)' (Inggris)
    ]
    for (const [b, n, v] of medan) expect(baris(sec(b), n).trim(), `${b} ${n}`).toBe(lf(v))
    expect(baris(sec('ExclutionTreatyObject'), 566).trim().replace(/&gt;/g, '>')).toBe(lf(LABEL_MEDAN_KHUSUS['ExclutionTreaty/Object'].Pct))
    expect(baris(sec('ExclutionTreaty'), 955).trim()).toBe(nilai(KLAUSUL_TCO.exclusionTreaty))
    // 10014 Co-Ins Scale: judul bagian, dua judul grid, kolom, Add/Edit kedua grid.
    const judul = (v: string) => `<pyTitle>${v}</pyTitle>`
    for (const [n, v] of [[917, nilai(KLAUSUL_TCO.coInsScale)],
      [8708, judul(JUDUL_GRID_KLAUSUL['CoinsPanel/Less Than'])], [12971, judul(JUDUL_GRID_KLAUSUL['CoinsPanel/More Than'])],
      [10180, nilai(KLAUSUL_TCO.coInsuranceShare)], [14445, nilai(KLAUSUL_TCO.coInsuranceShare)],
      [10340, nilai(LABEL_MEDAN_KLAUSUL.TreatyLimit)], [14605, nilai(LABEL_MEDAN_KLAUSUL.TreatyLimit)],
      [10554, tombol(KLAUSUL_TCO.add)], [14815, tombol(KLAUSUL_TCO.add)], [11254, tombol(KLAUSUL_TCO.edit)],
      [9982, '<pyFieldValueForNoRows>pzRDLNoResults</pyFieldValueForNoRows>']] as [number, string][]) {
      expect(baris(sec('Coins'), n).trim(), `Coins ${n}`).toBe(v)
    }
    // 10017 MB Capacity: judul, label form, kepala dan sel grid, tombol, teks Choose.
    const mb = LABEL_MEDAN_KHUSUS.LimitMB
    const esc = (v: string) => v.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
    const kolomMB: Record<string, string> = Object.fromEntries(KOLOM_GRID_KLAUSUL.LimitMB)
    for (const [n, v] of [[872, nilai(JUDUL_GRID_KLAUSUL.LimitMB)], [1586, lf(mb.ID_Occupation)], [1928, lf(mb.Pct)],
      [2169, lf(mb.PctMe)], [2410, lf(mb.Rp)], [2651, lf(mb.Usd)], [2892, lf(mb.MoreRp)], [3133, lf(mb.MoreUsd)],
      [3374, lf(`.${mb.TerritorialLimit}`)], [1696, `<pyNoSelectionText>${KLAUSUL_TCO.choose}</pyNoSelectionText>`],
      [8045, nilai(kolomMB.Occupation ?? '')], [8205, nilai(esc(kolomMB.Pct ?? ''))], [8361, nilai(esc(kolomMB.PctMe ?? ''))],
      [8517, nilai(esc(kolomMB.Usd ?? ''))], [8673, nilai(esc(kolomMB.MoreUsd ?? ''))], [8829, nilai(kolomMB.TerritorialLimit ?? '')],
      [9281, nilai('.Occupation')], [9451, nilai('.Pct')], [9640, nilai('.PctMe')], [9829, nilai('.Usd')],
      [10018, nilai('.MoreUsd')], [10207, nilai('.TerritorialLimit')],
      [5082, tombol(KLAUSUL_TCO.save)], [5361, tombol(KLAUSUL_TCO.cancel)], [9041, tombol(KLAUSUL_TCO.add)],
      [10428, tombol(KLAUSUL_TCO.edit)]] as [number, string][]) {
      expect(baris(sec('LIMITMB'), n).trim(), `LIMITMB ${n}`).toBe(v)
    }
  })

  it('tiket 06: grid dan form security VERBATIM, aksi tombolnya terbukti', () => {
    const nilai = (v: string) => `<pyValue>${v}</pyValue>`
    const lf = (v: string) => `<pyLabelFieldValue>${v}</pyLabelFieldValue>`
    const tombol = (v: string) => `<pyLabel>${v}</pyLabel>`
    const kasus: [number, string][] = [
      [15459, tombol(SECURITY_TCO.add)], [16088, nilai(SECURITY_TCO.kolomReasSecurity)],
      [16228, nilai(SECURITY_TCO.kolomSecurityName)], [16368, nilai(SECURITY_TCO.kolomPercentShare)],
      [17252, tombol(SECURITY_TCO.edit)], [17559, tombol(SECURITY_TCO.delete)],
      [19468, lf(SECURITY_TCO.formSecurityId)], [19648, lf(SECURITY_TCO.formSecurityName)],
      [19888, lf(SECURITY_TCO.formShare)], [20246, tombol(SECURITY_TCO.save)],
      [20980, lf(SECURITY_TCO.error)], [21717, lf('Informasi')],
      [16724, nilai('.REAS_SECURITY')], [16878, nilai('.CLIENTNAME')], [17013, nilai('.PCT_SHARE')],
      [17276, '<pyActivity>ShowEditSecurityReinsurer</pyActivity>'],
      [17583, '<pyActivity>DeleteSecurityReinsurer</pyActivity>'],
      [20270, '<pyActivity>SaveSecurityReinsurer_Act</pyActivity>'],
      [19711, '<pySourceName>BrowseAgentReinsSOA_RD</pySourceName>'],
    ]
    for (const [n, v] of kasus) expect(baris(KONTRAK, n).trim(), String(n)).toBe(v)
    expect(baris(REAS, 5277).trim()).toBe(tombol(REINSURER_TCO.securityReinsurer))
  })

  it('tiket 11: jalur kurs yang dipakai - testingKurs di kedua grid; XOL mengirim TreatyYear, bukan StartDate', () => {
    const H = join('Harness', 'InboxTreatyContractDescription.xml')
    expect(baris(H, 6138).trim()).toBe('<pyActivity>testingKurs</pyActivity>')
    expect(baris(H, 6153).trim()).toBe('<pyValue>InputTreatyArrangementDesc.StartDate</pyValue>')
    expect(baris(H, 9256).trim()).toBe('<pyActivity>testingKurs</pyActivity>')
    expect(baris(H, 9271).trim()).toBe('<pyValue>InputTreatyArrangementDesc.TreatyYear</pyValue>')
    expect(baris(join('Activity', 'testingKurs.xml'), 273).trim()).toBe('<PropertiesValue>Param.StartDate</PropertiesValue>')
    expect(baris(join('Section', 'NitipKurs.xml'), 512).trim()).toBe('<pyValue>InputTreatyArrangement.Kurs</pyValue>')
    expect(baris(join('Activity', 'NewTreatyArrEpi.xml'), 870).trim()).toBe(
      '<PropertiesValue>"Tidak ada Nilai Kurs di Tahun : "+ InputTreatyArrangementDesc.TreatyYear</PropertiesValue>')
    expect(Object.values(KURS_TCO).join(' ')).not.toMatch(/testing/i)
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
    const sumber = readFileSync(join(__dirname, 'labels.ts'), 'utf8')
    for (const kunci of ['masterKosong', 'belumDipilih', 'kosong:', 'menungguTiket']) {
      const i = sumber.indexOf(kunci)
      expect(i, kunci).toBeGreaterThan(0)
      expect(sumber.slice(sumber.lastIndexOf('/**', i), i)).toContain('[tidak ada di korpus]')
    }
  })
})

describe('bahasa Inggris (keputusan work owner 30-09-2026)', () => {
  it('label korpus berbahasa Indonesia diterjemahkan', () => {
    expect(REINSURER_TCO.informasi).toBe('Information')
    expect(SECURITY_TCO.informasi).toBe('Information')
    expect(LABEL_MEDAN_KLAUSUL.Layer).toBe('Max Period (Month)')
  })
  it('nol kata Indonesia di nilai label Treaty', () => {
    const semua = [MENU_TCO, JENIS_REASURANSI_TCO, TAHUN_TCO, KONTRAK_TCO, REINSURER_TCO, BUSINESS_TCO,
      KLAUSUL_TCO, LABEL_MEDAN_KLAUSUL, SECURITY_TCO, KURS_TCO, HAPUS_TCO].flatMap((o) => Object.values(o).map(String))
    const indo = /\b(tidak|belum|pilih|tersimpan|tutup|kombinasi|hapus|berlaku|kurs|bisnis|klausul|lampiran|tahun|ikut|ulangi|periksa|informasi|aktif|nonaktif|ya|batal|cari|sampai|terkirim|tertunda|gagal|baris|milik|bukan|dengan|yang)\b/i
    for (const v of semua) expect(v, v).not.toMatch(indo)
  })
})

