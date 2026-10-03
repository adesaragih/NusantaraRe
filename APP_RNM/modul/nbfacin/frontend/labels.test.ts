// Penjaga bukti label NB FacIn - tiket 21; portal Opportunity tiket 25.
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN. Berkas ini membuka
// section korpus dan memastikan setiap sel yang disebut `labels.ts` memang berlabel
// dan berproperti seperti tertulis. Korpus READ-ONLY - hanya dibaca.
//
// ⚠️ Bila korpus tidak terjangkau, test DILEWATI dengan pesan - bukan gagal.

import { existsSync, readdirSync, readFileSync } from 'node:fs'

import { describe, expect, it } from 'vitest'

import {
  KEPALA_PORTAL,
  BANGUNAN_KOSONG,
  GRID_OBJEK,
  KOLOM_PORTAL,
  MEDAN_COVERAGE_CARGO,
  OBJECT_ADDRESS,
  OBJECT_TYPE_LAINNYA,
  OPSI_FLOOR_TYPE,
  OPSI_OBJECT_TYPE,
  OPSI_ROOF_TYPE,
  OPSI_WALL_TYPE,
  PERIODE,
  POPUP_CEDING,
  POPUP_RISK,
  POPUP_TAMBAH_RISK,
  SISI_SEKITAR,
  MEDAN_SISI,
  LAIN_SEKITAR,
  TEKS_SEKITAR,
  GRID_ITEM,
  FORM_ITEM,
  ITEM_KOSONG,
  TEKS_ITEM,
  GRID_OKUPASI,
  FORM_OKUPASI,
  POPUP_OKUPASI,
  POPUP_KONSTRUKSI,
  GRID_FEA,
  FORM_FEA,
  GRID_KERUGIAN,
  FORM_KERUGIAN,
  LOSS_RATIO,
  GRID_KLAIM_INTERNAL,
  OPSI_CONSTRUCTION,
  CONSTRUCTION_KOSONG,
  OPSI_OWNERSHIP,
  OPSI_HOUSEKEEPING,
  OPSI_FLOOD_STATUS,
  OPSI_FLOOD_AREA,
  OPSI_CONDITION,
  OPSI_MINMAX,
  OPSI_KONDISI_DEDUCTIBLE,
  POPUP_AKUMULASI as P_AKUM,
  OPSI_TYPE_DEDUCTIBLE,
  OPSI_TYPE_DEDUCTIBLE2,
  FORM_DEDUCTIBLE,
  JUDUL_DEDUCTIBLE,
  KOLOM_TIME_EXCESS,
  DEDUCTIBLE_KOSONG,
  CONDITION_KOSONG,
  OPSI_PCT_ADJUST,
  OPSI_FIRE_BRIGADE,
  OPSI_SOP_SAFETY,
  OPSI_SOP_RISIKO,
  OPSI_REMARKS,
  LABEL_REMARKS,
  GRID_COV_OBJEK,
  GRID_COV_ITEM,
  GRID_COV_TOTAL,
  GRID_COVERAGE,
  FORM_COV,
  LABEL_COVERAGE_BASIS,
  OPSI_COVERAGE_BASIS,
  SIMPAN_COVERAGE,
  OPSI_TITLE_RISK,
  POPUP_SOB,
  SARING_PORTAL,
  SIMPAN_OBJEK,
  SUBTAB_OBJEK,
  SHOW_DETAIL,
  TAB_DETAIL,
  TEKS_OBJEK,
  TOMBOL_COVERAGE_CARGO,
  TOMBOL_KAKI_INWARD,
} from './labels'

const SECTION = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\InputCoverageCargo_FacIn.xml'
const ada = existsSync(SECTION)
const KEPALA = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\SFAPortalOpportunitiesHeader.xml'
const DAFTAR = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\SFAPortal_OpportunitiesList.xml'
const adaPortal = existsSync(KEPALA) && existsSync(DAFTAR)

/**
 * Isi elemen sel ber-`pyCellId` tertentu: `rowdata` terdalam yang terbuka saat
 * `<pyCellId>` itu ditemukan, sampai penutup PASANGANNYA (rowdata bersarang - mis.
 * aksi tombol - ikut di dalamnya). Pemindai kedalaman, bukan regex tak-rakus.
 */
function blokSel(xml: string, sel: string): string[] {
  const hasil: string[] = []
  const pola = /<rowdata\b[^>]*>|<\/rowdata>|<pyCellId>(\d+)<\/pyCellId>/g
  const tumpukan: { mulai: number; cocok: boolean }[] = []
  for (const m of xml.matchAll(pola)) {
    const i = m.index ?? 0
    if (m[0].startsWith('<rowdata')) {
      tumpukan.push({ mulai: i, cocok: false })
    } else if (m[0] === '</rowdata>') {
      const atas = tumpukan.pop()
      if (atas?.cocok) hasil.push(xml.slice(atas.mulai, i + m[0].length))
    } else if (m[1] === sel) {
      const atas = tumpukan.at(-1)
      if (atas) atas.cocok = true
    }
  }
  return hasil
}

describe.skipIf(!ada)('label NB FacIn = korpus InputCoverageCargo_FacIn', () => {
  const xml = ada ? readFileSync(SECTION, 'utf-8') : ''

  it.each(Object.entries(MEDAN_COVERAGE_CARGO))('medan %s', (_, m) => {
    const blok = blokSel(xml, m.sel)
    expect(blok.some((b) => b.includes(`<pyLabelFor>${m.label}</pyLabelFor>`) && b.includes(`<pyValue>${m.properti}</pyValue>`))).toBe(true)
  })

  it.each(Object.entries(TOMBOL_COVERAGE_CARGO))('tombol %s', (_, t) => {
    expect(blokSel(xml, t.sel).some((b) => b.includes(`<pyLabel>${t.label}</pyLabel>`))).toBe(true)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(xml, '11').some((b) => b.includes('<pyLabelFor>Rate (‰)</pyLabelFor>'))).toBe(false)
  })
})

if (!ada) {
  it('korpus tidak terjangkau - bukti label dilewati', () => {
    console.warn(`korpus tidak ada di ${SECTION}; uji bukti label NB FacIn dilewati`)
  })
}

describe.skipIf(!adaPortal)('label portal Opportunity = korpus SFAPortalOpportunitiesHeader + SFAPortal_OpportunitiesList', () => {
  const kepala = adaPortal ? readFileSync(KEPALA, 'utf-8') : ''
  const daftar = adaPortal ? readFileSync(DAFTAR, 'utf-8') : ''
  /** Sel itu SATU di berkasnya, dan teksnya ada di tag yang disebut. */
  const ditemukan = (xml: string, u: { sel: string; tag: string; label: string }) => {
    const blok = blokSel(xml, u.sel)
    return blok.length === 1 && blok[0]!.includes(`<${u.tag}>${u.label}</${u.tag}>`)
  }

  it.each(Object.entries(KEPALA_PORTAL))('kepala %s', (_, u) => {
    expect(ditemukan(kepala, u)).toBe(true)
  })

  it.each(Object.entries(SARING_PORTAL))('saring %s', (_, u) => {
    expect(ditemukan(daftar, u)).toBe(true)
  })

  it.each(KOLOM_PORTAL.map((k) => [k.sel, k] as const))('judul kolom sel %s', (_, k) => {
    const blok = blokSel(daftar, k.sel)
    expect(blok).toHaveLength(1)
    if (k.label === '') {
      // Kolom berjudul kosong: sel judulnya tanpa `pyValue` berisi.
      expect(blok[0]).not.toMatch(/<pyValue>[^<]+<\/pyValue>/)
    } else {
      expect(blok[0]).toContain(`<pyValue>${k.label}</pyValue>`)
    }
  })

  it('kolom berurutan = sel 89..96 (grid GetListOpportunityF)', () => {
    expect(KOLOM_PORTAL.map((k) => k.sel)).toEqual(['89', '90', '91', '92', '93', '94', '95', '96'])
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(ditemukan(kepala, { sel: '72', tag: 'pyLabel', label: 'Create Opportunity' })).toBe(false)
    expect(ditemukan(daftar, { sel: '89', tag: 'pyValue', label: 'Offer No.' })).toBe(false)
  })
})

if (!adaPortal) {
  it('korpus portal tidak terjangkau - bukti label portal dilewati', () => {
    console.warn(`korpus tidak ada di ${KEPALA} / ${DAFTAR}; uji bukti label portal dilewati`)
  })
}

const NBFACIN = 'D:\\migrasi\\RNM\\NB FacIn\\'
const BERKAS_INWARD = {
  periode: NBFACIN + 'Section\\Periode.xml',
  layar: NBFACIN + 'Section\\InputInwardFacultative.xml',
  dtl: NBFACIN + 'Section\\InputInwardFacultativeDtl.xml',
  aksi: NBFACIN + 'FlowAction\\InwardFacultative.xml',
}
const adaInward = Object.values(BERKAS_INWARD).every((b) => existsSync(b))

describe.skipIf(!adaInward)('label Inward Facultative = korpus (Periode, InputInwardFacultative, Dtl, flow action)', () => {
  const baca = (b: string) => (adaInward ? readFileSync(b, 'utf-8') : '')
  const periode = baca(BERKAS_INWARD.periode)

  it.each(Object.entries(PERIODE))('Periode %s', (_, u) => {
    if (u.sel === '') {
      // Judul layout / label tanpa sel: cukup ada di berkas dengan tag itu.
      expect(periode).toContain(`<${u.tag}>${u.label}</${u.tag}>`)
    } else {
      const blok = blokSel(periode, u.sel)
      expect(blok.some((b) => b.includes(`<${u.tag}>${u.label}</${u.tag}>`))).toBe(true)
    }
  })

  it('Show Detail = caption checkbox InputInwardFacultative', () => {
    expect(baca(BERKAS_INWARD.layar)).toContain(`<pyCheckboxCaption>${SHOW_DETAIL}</pyCheckboxCaption>`)
  })

  it.each(TAB_DETAIL.map((t) => [t]))('judul tab %s = pyTitle InputInwardFacultativeDtl', (t) => {
    expect(baca(BERKAS_INWARD.dtl)).toContain(`<pyTitle>${t}</pyTitle>`)
  })

  it.each(Object.entries(TOMBOL_KAKI_INWARD))('tombol kaki %s = flow action InwardFacultative', (_, u) => {
    expect(baca(BERKAS_INWARD.aksi)).toContain(`<${u.tag}>${u.label}</${u.tag}>`)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(periode, '9').some((b) => b.includes('<pyLabelFieldValue>Reff. Number</pyLabelFieldValue>'))).toBe(false)
    expect(baca(BERKAS_INWARD.dtl)).not.toContain('<pyTitle>Objects</pyTitle>')
  })
})

const HARNESS_SOB = NBFACIN + 'Harness\\SOB.xml'
const adaSob = existsSync(HARNESS_SOB) && existsSync(BERKAS_INWARD.periode)

describe.skipIf(!adaSob)('popup Change SOB = korpus (Periode sel 52 → harness SOB) - tiket 33', () => {
  const sob = adaSob ? readFileSync(HARNESS_SOB, 'utf-8') : ''
  const periode = adaSob ? readFileSync(BERKAS_INWARD.periode, 'utf-8') : ''

  it('tombol Change SOB membuka harness SOB sebagai popup berjudul Change SOB', () => {
    const blok = blokSel(periode, PERIODE.changeSob.sel).find((b) => b.includes('<pyHarnessName>SOB</pyHarnessName>')) ?? ''
    expect(blok).toContain(`<pyWindowName>${POPUP_SOB.judul}</pyWindowName>`)
    expect(blok).toContain('<pyTarget>popup</pyTarget>')
  })

  it('kotak Search dan tombol Choose', () => {
    expect(sob).toContain(`<${POPUP_SOB.cari.tag}>${POPUP_SOB.cari.label}</${POPUP_SOB.cari.tag}>`)
    expect(sob).toContain(`<${POPUP_SOB.pilih.tag}>${POPUP_SOB.pilih.label}</${POPUP_SOB.pilih.tag}>`)
  })

  it('grid kelas Data-Agent dengan kolom .ID dan .ClientName', () => {
    expect(sob).toContain('<pySmartPromptClass>ASM-FW-GISFW-Data-Agent</pySmartPromptClass>')
    expect(sob).toContain('<pyValue>.ID</pyValue>')
    expect(sob).toContain('<pyValue>.ClientName</pyValue>')
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(sob).not.toContain('<pyLabel>Select</pyLabel>')
    expect(blokSel(periode, '53').some((b) => b.includes('<pyHarnessName>SOB</pyHarnessName>'))).toBe(false)
  })
})

const SECTION_CEDING = NBFACIN + 'Section\\ShowCedingCoList.xml'
const HARNESS_CEDING = NBFACIN + 'Harness\\ShowCedingCoList.xml'
const PENCARI_CEDING = NBFACIN + 'Section\\CedingCoHierarki.xml'
const adaCeding = [SECTION_CEDING, HARNESS_CEDING, PENCARI_CEDING, BERKAS_INWARD.periode].every((b) => existsSync(b))

describe.skipIf(!adaCeding)('popup Change Ceding Co = korpus (Periode sel 53 → ShowCedingCoList → CedingCompany) - tiket 34', () => {
  const baca = (b: string) => (adaCeding ? readFileSync(b, 'utf-8') : '')
  const daftar = baca(SECTION_CEDING)
  const harness = baca(HARNESS_CEDING)
  const pencari = baca(PENCARI_CEDING)
  const periode = baca(BERKAS_INWARD.periode)

  it('tombol Change Ceding Co membuka harness ShowCedingCoList sebagai popup; judulnya pyLabel harness', () => {
    const blok = blokSel(periode, PERIODE.changeCedingCo.sel).find((b) => b.includes('<pyHarnessName>ShowCedingCoList</pyHarnessName>')) ?? ''
    expect(blok).toContain('<pyTarget>popup</pyTarget>')
    expect(harness).toContain(`<pyLabel>${POPUP_CEDING.judul}</pyLabel>`)
  })

  it.each([['kolom', POPUP_CEDING.kolom], ['tambah', POPUP_CEDING.tambah], ['hapus', POPUP_CEDING.hapus]] as const)('sel %s', (_, u) => {
    expect(blokSel(daftar, u.sel).some((b) => b.includes(`<${u.tag}>${u.label}</${u.tag}>`))).toBe(true)
  })

  it('kepala Ceding Co wajib; Add / Select Ceding membuka CedingCompany berjudul Ceding Company; grid atas CedingCoList', () => {
    expect(blokSel(daftar, '14').some((b) => b.includes('<pyRequired>true</pyRequired>'))).toBe(true)
    const tambah = blokSel(daftar, POPUP_CEDING.tambah.sel).find((b) => b.includes('<pyHarnessName>CedingCompany</pyHarnessName>')) ?? ''
    expect(tambah).toContain(`<pyWindowName>${POPUP_CEDING.judulCari}</pyWindowName>`)
    expect(tambah).toContain('<pyActivity>AddCedingList_act</pyActivity>')
    expect(daftar).toContain('<pyPageListProperty>pyWorkPage.Quotation.CedingCoList</pyPageListProperty>')
  })

  it('Submit = tombol harness', () => {
    expect(harness).toContain(`<${POPUP_CEDING.submit.tag}>${POPUP_CEDING.submit.label}</${POPUP_CEDING.submit.tag}>`)
    expect(harness).toContain('<pyActivity>SetCedingCo_Act</pyActivity>')
  })

  it('popup pencari CedingCoHierarki: kolom ID / Client ID / Name, Choose, RD BrowseAgentNonLife_RD (sama dengan SOB)', () => {
    for (const k of POPUP_SOB.kolom) expect(pencari).toContain(`<pyValue>${k}</pyValue>`)
    expect(pencari).toContain(`<pyLabel>${POPUP_SOB.pilih.label}</pyLabel>`)
    expect(pencari).toContain('<pyGridRDName>BrowseAgentNonLife_RD</pyGridRDName>')
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(daftar, '14').some((b) => b.includes('<pyValue>Ceding Company</pyValue>'))).toBe(false)
    expect(harness).not.toContain('<pyLabel>Ceding Co Lists</pyLabel>')
  })
})

const BERKAS_OBJEK = {
  daftar: NBFACIN + 'Section\\ObjectList.xml',
  subtab: NBFACIN + 'Section\\Property.xml',
  alamat: NBFACIN + 'Section\\ObjectDetails.xml',
  lantai: NBFACIN + 'Activity\\SetErrorMessageFloorNumber_Act.xml',
  unggah: NBFACIN + 'Activity\\InsertUploadFire_act.xml',
  dtl: NBFACIN + 'Section\\InputInwardFacultativeDtl.xml',
}
const adaObjek = Object.values(BERKAS_OBJEK).every((b) => existsSync(b))

describe.skipIf(!adaObjek)('tab Object FIRE = korpus (ObjectList, Property, ObjectDetails) - tiket 35', () => {
  const baca = (b: string) => (adaObjek ? readFileSync(b, 'utf-8') : '')
  const daftar = baca(BERKAS_OBJEK.daftar)
  const alamat = baca(BERKAS_OBJEK.alamat)

  it.each(GRID_OBJEK.kolom.map((k) => [k.label, k] as const))('kolom grid %s', (_, k) => {
    expect(blokSel(daftar, k.sel).some((b) => b.includes(`<${k.tag}>${k.label}</${k.tag}>`))).toBe(true)
  })

  it('grid atas .LocationList, master-detail, 10 per halaman', () => {
    expect(daftar).toContain('<pyPageListProperty>.LocationList</pyPageListProperty>')
    expect(daftar).toContain('<pyRowEditing>masterDetail</pyRowEditing>')
    expect(daftar).toContain(`<pyPageSize>${GRID_OBJEK.ukuran}</pyPageSize>`)
  })

  it.each(SUBTAB_OBJEK.map((t) => [t]))('sub-tab %s = pyTitle Property', (t) => {
    expect(baca(BERKAS_OBJEK.subtab)).toContain(`<pyTitle>${t}</pyTitle>`)
  })

  it.each(Object.entries(OBJECT_ADDRESS))('Object Address %s', (_, u) => {
    if (u.sel === '') expect(alamat).toContain(`<${u.tag}>${u.label}</${u.tag}>`)
    else expect(blokSel(alamat, u.sel).some((b) => b.includes(`<${u.tag}>${u.label}</${u.tag}>`))).toBe(true)
  })

  it('Object Type wajib; Object Name tampil hanya bila Others', () => {
    expect(blokSel(alamat, OBJECT_ADDRESS.objectType.sel).some((b) => b.includes('<pyRequired>true</pyRequired>'))).toBe(true)
    expect(alamat).toContain(`<pyCondition>.Property.ObjectType = '${OBJECT_TYPE_LAINNYA}'</pyCondition>`)
  })

  it('Save = Dtl sel 25; pesan lantai minus = SetErrorMessageFloorNumber_Act', () => {
    expect(blokSel(baca(BERKAS_OBJEK.dtl), SIMPAN_OBJEK.sel).some((b) => b.includes(`<${SIMPAN_OBJEK.tag}>${SIMPAN_OBJEK.label}</${SIMPAN_OBJEK.tag}>`))).toBe(true)
    expect(baca(BERKAS_OBJEK.lantai)).toContain(`<Message>${TEKS_OBJEK.lantaiMinus}</Message>`)
  })

  it('pilihan Object Type = HIMPUNAN di ekspresi InsertUploadFire_act (selain itu "Others"); urutan dari tangkapan layar', () => {
    const unggah = baca(BERKAS_OBJEK.unggah)
    const nilai = [...(/OBJECT_TYPE!=("[^,]*),"Others"/.exec(unggah)?.[1] ?? '').matchAll(/"([^"]+)"/g)].map((m) => m[1])
    expect(nilai).toHaveLength(9)
    expect([...nilai, 'Others'].sort()).toEqual([...OPSI_OBJECT_TYPE].sort())
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(alamat, '5').some((b) => b.includes('<pyLabelFieldValue>Object Type</pyLabelFieldValue>'))).toBe(false)
    expect(blokSel(daftar, '14').some((b) => b.includes('<pyValue>No</pyValue>'))).toBe(false)
  })
})

const DDL = 'D:\\migrasi\\RNM\\DDL\\'
const BERKAS_BANGUNAN = { RoofType: OPSI_ROOF_TYPE, WallType: OPSI_WALL_TYPE, FloorType: OPSI_FLOOR_TYPE }
const adaBangunan = Object.keys(BERKAS_BANGUNAN).every((n) => existsSync(`${DDL}${n}.xml`))

/** Baris PromptList aturan properti: [pyStandardValue, pyLocalizedValue], berurutan. */
function daftarPrompt(xml: string): [string, string][] {
  const isi = /<pyPromptTableList[^>]*>([\s\S]*?)<\/pyPromptTableList>/.exec(xml)?.[1] ?? ''
  return isi
    .split(/<rowdata REPEATINGINDEX="/)
    .slice(1)
    .map((b) => [/<pyStandardValue>([^<]*)</.exec(b)?.[1] ?? '', /<pyLocalizedValue>([^<]*)</.exec(b)?.[1] ?? ''])
}

describe.skipIf(!adaBangunan)('Roof / Wall / Floor Type = aturan properti Pega (DDL\\*.xml) - tiket 35', () => {
  it.each(Object.entries(BERKAS_BANGUNAN))('%s: baris pertama = pilihan kosong, sisanya = value/label berurutan', (nama, opsi) => {
    const baris = daftarPrompt(readFileSync(`${DDL}${nama}.xml`, 'utf-8'))
    expect(baris[0]).toEqual(['', BANGUNAN_KOSONG])
    expect(baris.slice(1)).toEqual(opsi.map((o) => [o.value, o.label]))
  })

  it('uji ini menggigit: pasangan yang tertukar tidak cocok', () => {
    const baris = daftarPrompt(readFileSync(`${DDL}FloorType.xml`, 'utf-8'))
    expect(baris).not.toContainEqual(['KELAS I', 'Keramik'])
  })
})

const BERKAS_RISK = {
  saring: NBFACIN + 'Section\\ChooseRiskAddress.xml',
  hasil: NBFACIN + 'Section\\ChooseRiskAddress_ResultList.xml',
  alamat: NBFACIN + 'Section\\ObjectDetails.xml',
}
const adaRisk = Object.values(BERKAS_RISK).every((b) => existsSync(b))

describe.skipIf(!adaRisk)('popup Choose Risk Address = korpus (ChooseRiskAddress, _ResultList) - tiket 36', () => {
  const baca = (b: string) => (adaRisk ? readFileSync(b, 'utf-8') : '')
  const saring = baca(BERKAS_RISK.saring)
  const hasil = baca(BERKAS_RISK.hasil)

  it('Choose Risk Address membuka harness ChooseRiskAddress berjudul Choose Risk Location', () => {
    const blok = blokSel(baca(BERKAS_RISK.alamat), '16').find((b) => b.includes('<pyHarnessName>ChooseRiskAddress</pyHarnessName>')) ?? ''
    expect(blok).toContain(`<pyWindowName>${POPUP_RISK.judul}</pyWindowName>`)
    expect(blok).toContain('<pyTarget>popup</pyTarget>')
  })

  it.each([...Object.values(POPUP_RISK.saring), POPUP_RISK.cari, POPUP_RISK.tambah].map((u) => [u.label, u] as const))('saringan/tombol %s', (_, u) => {
    expect(blokSel(saring, u.sel).some((b) => b.includes(`<${u.tag}>${u.label}</${u.tag}>`))).toBe(true)
  })

  it.each(POPUP_RISK.kolom.map((k) => [k.label, k] as const))('kolom grid %s', (_, k) => {
    expect(blokSel(hasil, k.sel).some((b) => b.includes(`<${k.tag}>${k.label}</${k.tag}>`))).toBe(true)
  })

  it('grid 10 per halaman bernomor; Pilih = SetRiskIdDT_FacIn', () => {
    expect(hasil).toContain(`<pyPageSize>${POPUP_RISK.ukuran}</pyPageSize>`)
    expect(hasil).toContain('<pyPageMode>Numeric</pyPageMode>')
    expect(hasil).toContain('<pyName>SetRiskIdDT_FacIn</pyName>')
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(saring, '79').some((b) => b.includes('<pyLabelFieldValue>Zip code</pyLabelFieldValue>'))).toBe(false)
  })
})

const BERKAS_TAMBAH = {
  isian: NBFACIN + 'Section\\InputRiskAddress.xml',
  akumulasi: NBFACIN + 'RDBList\\SearchAccumulationbypersetase_SQL.xml',
  pencari: NBFACIN + 'Section\\ChooseRiskAddress.xml',
}
const adaTambah = Object.values(BERKAS_TAMBAH).every((b) => existsSync(b))

describe.skipIf(!adaTambah)('popup Add alamat risiko = korpus (InputRiskAddress) - tiket 37', () => {
  const baca = (b: string) => (adaTambah ? readFileSync(b, 'utf-8') : '')
  const isian = baca(BERKAS_TAMBAH.isian)
  const medan = Object.entries(POPUP_TAMBAH_RISK).filter(([k]) => k !== 'judul') as [string, { sel: string; tag: string; label: string }][]

  it.each(medan)('medan/tombol %s', (_, u) => {
    expect(blokSel(isian, u.sel).some((b) => b.includes(`<${u.tag}>${u.label}</${u.tag}>`))).toBe(true)
  })

  it('judul = tombol Add pembuka harness ChooseRiskLocation', () => {
    const blok = blokSel(baca(BERKAS_TAMBAH.pencari), '74').find((b) => b.includes('<pyHarnessName>ChooseRiskLocation</pyHarnessName>')) ?? ''
    expect(blok).toContain(`<pyLabel>${POPUP_TAMBAH_RISK.judul}</pyLabel>`)
  })

  it('syarat tampil berantai Province <- Country, City <- Province, District <- City, Territory <- District', () => {
    for (const k of ['NationName', 'ProvinceName', 'CityName', 'DistrictName']) {
      expect(isian).toContain(`<pyCondition>InputRiskAddress.${k} != ''</pyCondition>`)
    }
  })

  it('pilihan Title = urutan REGEXP di SearchAccumulationbypersetase_SQL', () => {
    const re = /\((DESA\|[^)]*)\)/.exec(baca(BERKAS_TAMBAH.akumulasi))?.[1] ?? ''
    expect(re.split('|').map((x) => x.replace('\\.', '.'))).toEqual([...OPSI_TITLE_RISK])
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(isian, '13').some((b) => b.includes('<pyLabelFieldValue>Zip code</pyLabelFieldValue>'))).toBe(false)
  })
})

const BERKAS_SEKITAR = {
  bagian: NBFACIN + 'Section\\RiskAround.xml',
  minus: NBFACIN + 'Activity\\NegativeIsNotAllowed.xml',
}
const adaSekitar = Object.values(BERKAS_SEKITAR).every((b) => existsSync(b))

describe.skipIf(!adaSekitar)('Surrounding Risk = korpus (RiskAround) - tiket 38', () => {
  const baca = (b: string) => (adaSekitar ? readFileSync(b, 'utf-8') : '')
  const xml = baca(BERKAS_SEKITAR.bagian)
  const ada = (sel: string, tag: string, label: string) => blokSel(xml, sel).some((b) => b.includes(`<${tag}>${label}</${tag}>`))

  it.each(SISI_SEKITAR.map((s) => [s.judul, s] as const))('sisi %s: judul + empat medan di sel yang benar', (_, s) => {
    expect(xml).toContain(`<pyTitle>${s.judul}</pyTitle>`)
    expect(ada(s.occupation, 'pyLabelFieldValue', MEDAN_SISI.occupation)).toBe(true)
    expect(ada(s.construction, 'pyLabelFieldValue', MEDAN_SISI.construction)).toBe(true)
    expect(ada(s.distance, 'pyLabelFieldValue', MEDAN_SISI.distance)).toBe(true)
    expect(ada(s.note, 'pyLabelFieldValue', MEDAN_SISI.note)).toBe(true)
    const nama = s.judul
    expect(blokSel(xml, s.occupation).some((b) => b.includes(`<pyValue>.Property.SurroundingRisk.${nama}Occupation</pyValue>`))).toBe(true)
    expect(blokSel(xml, s.distance).some((b) => b.includes(`<pyValue>.Property.SurroundingRisk.${nama}Distance</pyValue>`))).toBe(true)
  })

  it('Occupation = data page D_BrowseOccupationFacInFIRE (TYPE FIRE)', () => {
    expect(blokSel(xml, '9').some((b) => b.includes('D_BrowseOccupationFacInFIRE') && b.includes('FIRE'))).toBe(true)
  })

  it.each(Object.entries(LAIN_SEKITAR))('Other Description %s', (_, u) => {
    if (u.sel === '') expect(xml).toContain(`<${u.tag}>${u.label}</${u.tag}>`)
    else expect(ada(u.sel, u.tag, u.label)).toBe(true)
  })

  it('Flood Area tampil hanya bila Flood Area Status = 0; pesan Distance minus = NegativeIsNotAllowed', () => {
    expect(xml).toContain('<pyCondition>.Property.SurroundingRisk.FloodAreaStatus==0</pyCondition>')
    expect(baca(BERKAS_SEKITAR.minus)).toContain(TEKS_SEKITAR.jarakMinus)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(ada('13', 'pyLabelFieldValue', 'Distance')).toBe(false)
  })
})

const BERKAS_SEKITAR_DDL = {
  FrontConstruction: { opsi: OPSI_CONSTRUCTION, kosong: CONSTRUCTION_KOSONG as string | null },
  Ownership: { opsi: OPSI_OWNERSHIP, kosong: null },
  HousekeepingStatus: { opsi: OPSI_HOUSEKEEPING, kosong: null },
  FloodAreaStatus: { opsi: OPSI_FLOOD_STATUS, kosong: null },
  FloodArea: { opsi: OPSI_FLOOD_AREA, kosong: CONSTRUCTION_KOSONG as string | null },
  PrivateFireBrigade: { opsi: OPSI_FIRE_BRIGADE, kosong: null },
  TeamSOPSafety: { opsi: OPSI_SOP_SAFETY, kosong: null },
  TeamSOPRiskManagement: { opsi: OPSI_SOP_RISIKO, kosong: null },
  Remarks: { opsi: OPSI_REMARKS, kosong: null },
}
const adaSekitarDDL = Object.keys(BERKAS_SEKITAR_DDL).every((n) => existsSync(`${DDL}${n}.xml`))

describe.skipIf(!adaSekitarDDL)('dropdown Surrounding Risk = aturan properti Pega (DDL\\*.xml) - tiket 38', () => {
  it.each(Object.entries(BERKAS_SEKITAR_DDL))('%s: pasangan value/label per rowdata, berurutan', (nama, u) => {
    const baris = daftarPrompt(readFileSync(`${DDL}${nama}.xml`, 'utf-8'))
    const isi = u.kosong === null ? baris : baris.slice(1)
    if (u.kosong !== null) expect(baris[0]).toEqual(['', u.kosong])
    expect(isi).toEqual(u.opsi.map((o) => [o.value, o.label]))
  })

  it('Flood Area Status "Yes" bernilai "0" = syarat tampil Flood Area', () => {
    expect(OPSI_FLOOD_STATUS.find((o) => o.label === 'Yes')?.value).toBe('0')
  })
})

const BERKAS_ITEM = {
  daftar: NBFACIN + 'Section\\PropertyItemList.xml',
  form: NBFACIN + 'Section\\PropertyItemFacIn_Section.xml',
  pct: NBFACIN + 'Activity\\ValidateAdjustPct.xml',
}
const adaItem = Object.values(BERKAS_ITEM).every((b) => existsSync(b))

describe.skipIf(!adaItem)('Object Item = korpus (PropertyItemList, PropertyItemFacIn_Section) - tiket 39', () => {
  const baca = (b: string) => (adaItem ? readFileSync(b, 'utf-8') : '')
  const daftar = baca(BERKAS_ITEM.daftar)
  const form = baca(BERKAS_ITEM.form)

  it.each([...GRID_ITEM.kolom, ...GRID_ITEM.total].map((k) => [k.label, k] as const))('kolom grid %s', (_, k) => {
    expect(blokSel(daftar, k.sel).some((b) => b.includes(`<pyValue>${k.label}</pyValue>`))).toBe(true)
  })

  it('grid item atas .Property.PropertyItemList (master-detail); total atas .Property.TotalTSIList', () => {
    expect(daftar).toContain('<pyPageListProperty>.Property.PropertyItemList</pyPageListProperty>')
    expect(daftar).toContain('<pyPageListProperty>.Property.TotalTSIList</pyPageListProperty>')
  })

  it.each(Object.entries(FORM_ITEM))('form %s', (_, u) => {
    expect(blokSel(form, u.sel).some((b) => b.includes(`<${u.tag}>${u.label}</${u.tag}>`))).toBe(true)
  })

  it('pilihan kosong "Choose"; Adjustment Pct. angka tampil bila Adjustable; pesan ValidateAdjustPct', () => {
    expect(form).toContain(`<pyNoSelectionText>${ITEM_KOSONG}</pyNoSelectionText>`)
    expect(blokSel(form, FORM_ITEM.pctAdjust.sel).some((b) => b.includes('<pyCondition>IsAdjustableFlag</pyCondition>'))).toBe(true)
    expect(baca(BERKAS_ITEM.pct)).toContain(TEKS_ITEM.pctAdjust)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(form, '17').some((b) => b.includes('<pyLabelFieldValue>TSI Object Item</pyLabelFieldValue>'))).toBe(false)
  })
})

const BERKAS_OKUPASI = {
  daftar: NBFACIN + 'Section\\OccupationList.xml',
  form: NBFACIN + 'Section\\OccupationItemFacIn_Section.xml',
  cari: NBFACIN + 'Section\\ChooseOccupation.xml',
  konstruksi: NBFACIN + 'Section\\ChooseClassofContraction.xml',
  data: NBFACIN + 'Activity\\SetDataOccupation.xml',
}
const adaOkupasi = Object.values(BERKAS_OKUPASI).every((b) => existsSync(b))

describe.skipIf(!adaOkupasi)('Occupation = korpus (OccupationList, OccupationItemFacIn_Section, Choose*) - tiket 40', () => {
  const baca = (b: string) => (adaOkupasi ? readFileSync(b, 'utf-8') : '')
  const daftar = baca(BERKAS_OKUPASI.daftar)
  const form = baca(BERKAS_OKUPASI.form)
  const cari = baca(BERKAS_OKUPASI.cari)
  const konstruksi = baca(BERKAS_OKUPASI.konstruksi)
  const ada = (xml: string, sel: string, tag: string, label: string) => blokSel(xml, sel).some((b) => b.includes(`<${tag}>${label}</${tag}>`))

  it.each(GRID_OKUPASI.map((k) => [k.label, k] as const))('kolom grid %s', (_, k) => {
    expect(ada(daftar, k.sel, 'pyValue', k.label)).toBe(true)
  })

  it('grid atas .Property.OccupationList', () => {
    expect(daftar).toContain('<pyPageListProperty>.Property.OccupationList</pyPageListProperty>')
  })

  it.each(Object.entries(FORM_OKUPASI))('form %s', (_, u) => {
    expect(ada(form, u.sel, u.tag, u.label)).toBe(true)
  })

  it('Class of Construction wajib', () => {
    expect(blokSel(form, FORM_OKUPASI.konstruksi.sel).some((b) => b.includes('<pyRequired>true</pyRequired>'))).toBe(true)
  })

  it('popup Occupation: kotak cari, kolom ID / Name, Choose; RD BrowseOccupationFacInFIRE_RD', () => {
    expect(ada(cari, POPUP_OKUPASI.cari.sel, POPUP_OKUPASI.cari.tag, POPUP_OKUPASI.cari.label)).toBe(true)
    for (const k of POPUP_OKUPASI.kolom) expect(ada(cari, k.sel, 'pyValue', k.label)).toBe(true)
    expect(ada(cari, POPUP_OKUPASI.pilih.sel, POPUP_OKUPASI.pilih.tag, POPUP_OKUPASI.pilih.label)).toBe(true)
    expect(cari).toContain('BrowseOccupationFacInFIRE_RD')
  })

  it('popup Class of Construction: kolom Description, Choose; RD BrowseTableOfLimit_RD', () => {
    expect(ada(konstruksi, POPUP_KONSTRUKSI.kolom.sel, 'pyValue', POPUP_KONSTRUKSI.kolom.label)).toBe(true)
    expect(ada(konstruksi, POPUP_KONSTRUKSI.pilih.sel, POPUP_KONSTRUKSI.pilih.tag, POPUP_KONSTRUKSI.pilih.label)).toBe(true)
    expect(konstruksi).toContain('BrowseTableOfLimit_RD')
  })

  it('SetDataOccupation: KDRiskExposure 03/02/01 -> III/II/I', () => {
    expect(baca(BERKAS_OKUPASI.data)).toMatch(/@if\(\.TableOfLimit\.Category="03","III",@if\(\.TableOfLimit\.Category="02","II",@if\(\.TableOfLimit\.Category="01","I",""\)\)\)/)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(ada(form, '9', 'pyLabelFieldValue', 'Class Of Construction')).toBe(false)
  })
})

const BERKAS_FEA = {
  daftar: NBFACIN + 'Section\\FEAList.xml',
  form: NBFACIN + 'Section\\InputFEA_IsUW.xml',
  aksi: NBFACIN + 'FlowAction\\InputFEA.xml',
}
const adaFEA = Object.values(BERKAS_FEA).every((b) => existsSync(b))
/** Label berisi `&` tertulis `&amp;` di XML. */
const xmlTeks = (s: string) => s.replace(/&/g, '&amp;')

describe.skipIf(!adaFEA)('FEA = korpus (FEAList, InputFEA_IsUW, flow action InputFEA) - tiket 41', () => {
  const baca = (b: string) => (adaFEA ? readFileSync(b, 'utf-8') : '')
  const daftar = baca(BERKAS_FEA.daftar)
  const form = baca(BERKAS_FEA.form)

  it('FEA = Fire Extinguisher Availability (pyLabel flow action InputFEA)', () => {
    expect(baca(BERKAS_FEA.aksi)).toContain('<pyLabel>Input Fire Extinguisher Availability</pyLabel>')
  })

  it.each(GRID_FEA.map((k) => [k.label, k] as const))('kolom grid %s', (_, k) => {
    expect(blokSel(daftar, k.sel).some((b) => b.includes(`<pyValue>${xmlTeks(k.label)}</pyValue>`))).toBe(true)
  })

  it('grid atas .FEAList (master-detail, flow action InputFEA)', () => {
    expect(daftar).toContain('<pyPageListProperty>.FEAList</pyPageListProperty>')
    expect(daftar).toContain('<pyRowEditing>masterDetail</pyRowEditing>')
  })

  it.each(Object.entries(FORM_FEA))('form %s', (_, u) => {
    expect(blokSel(form, u.sel).some((b) => b.includes(`<pyLabelFieldValue>${xmlTeks(u.label)}</pyLabelFieldValue>`))).toBe(true)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(blokSel(form, '12').some((b) => b.includes('<pyLabelFieldValue>APAR</pyLabelFieldValue>'))).toBe(false)
  })
})

const BERKAS_KERUGIAN = {
  daftar: NBFACIN + 'Section\\CauseOfLoss_FacIn.xml',
  form: NBFACIN + 'Section\\InputCauseOfLoss_FacIn.xml',
  rasio: NBFACIN + 'Section\\InputOfferFacInLossRatio.xml',
  internal: NBFACIN + 'Section\\CauseOfLossClaim_FacIn.xml',
}
const adaKerugian = Object.values(BERKAS_KERUGIAN).every((b) => existsSync(b))

describe.skipIf(!adaKerugian)('Loss Record + Loss Record Internal = korpus - tiket 42', () => {
  const baca = (b: string) => (adaKerugian ? readFileSync(b, 'utf-8') : '')
  const daftar = baca(BERKAS_KERUGIAN.daftar)
  const form = baca(BERKAS_KERUGIAN.form)
  const rasio = baca(BERKAS_KERUGIAN.rasio)
  const internal = baca(BERKAS_KERUGIAN.internal)
  const ada = (xml: string, sel: string, tag: string, label: string) => blokSel(xml, sel).some((b) => b.includes(`<${tag}>${label}</${tag}>`))

  it.each(GRID_KERUGIAN.map((k) => [k.label, k] as const))('kolom Loss Record %s', (_, k) => {
    expect(ada(daftar, k.sel, 'pyValue', k.label)).toBe(true)
  })

  it.each(Object.entries(FORM_KERUGIAN))('form %s', (_, u) => {
    expect(ada(form, u.sel, 'pyLabelFieldValue', u.label)).toBe(true)
  })

  it('Loss Detail wajib; grid atas .Property.ListCauseOfLoss', () => {
    expect(blokSel(form, FORM_KERUGIAN.detail.sel).some((b) => b.includes('<pyRequired>true</pyRequired>'))).toBe(true)
    expect(daftar).toContain('<pyPageListProperty>.Property.ListCauseOfLoss</pyPageListProperty>')
  })

  it.each(LOSS_RATIO.map((l) => [l.label, l] as const))('loss ratio %s', (_, l) => {
    expect(ada(rasio, l.sel, 'pyLabelFieldValue', l.label)).toBe(true)
  })

  it.each(GRID_KLAIM_INTERNAL.map((k) => [k.label, k] as const))('kolom Loss Record Internal %s', (_, k) => {
    expect(ada(internal, k.sel, 'pyValue', k.label)).toBe(true)
  })

  it('Loss Record Internal: tombol Get tersembunyi (1=0), grid baca-saja atas .Property.ListCauseOfLossClaim', () => {
    expect(internal).toContain('<pyPageListProperty>.Property.ListCauseOfLossClaim</pyPageListProperty>')
    expect(blokSel(internal, '1').some((b) => b.includes('<pyCondition>1=0</pyCondition>'))).toBe(true)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(ada(form, '7', 'pyLabelFieldValue', 'Total Claim')).toBe(false)
  })
})

describe.skipIf(!existsSync(`${DDL}PctAdjust2.xml`))('PctAdjust2 = LocalList aturan properti (DDL\\PctAdjust2.xml)', () => {
  it('pyTableOption LocalList, satu nilai "100"', () => {
    const xml = readFileSync(`${DDL}PctAdjust2.xml`, 'utf-8')
    expect(xml).toContain('<pyTableOption>LocalList</pyTableOption>')
    const isi = /<pyListValues[^>]*>([\s\S]*?)<\/pyListValues>/.exec(xml)?.[1] ?? ''
    const nilai = [...isi.matchAll(/<pyLabel>([^<]*)<\/pyLabel>/g)].map((m) => m[1])
    expect(nilai).toEqual(OPSI_PCT_ADJUST.map((o) => o.value))
  })
})

describe.skipIf(!existsSync(`${DDL}Remarks.xml`))('label Remarks = pyLabel aturan properti (DDL\\Remarks.xml)', () => {
  it('pyLabel aturan CAUSEOFLOSS!REMARKS = LABEL_REMARKS', () => {
    const xml = readFileSync(`${DDL}Remarks.xml`, 'utf-8')
    expect(xml).toContain('<pxInsName>ASM-FW-GISFW-DATA-CAUSEOFLOSS!REMARKS</pxInsName>')
    expect(xml).toContain(`<pyLabel>${LABEL_REMARKS}</pyLabel>`)
  })
})

const BERKAS_COV = {
  dtl: NBFACIN + 'Section\\InputInwardFacultativeDtl.xml',
  objek: NBFACIN + 'Section\\CoverageList.xml',
  item: NBFACIN + 'Section\\PropertyItemListCoverage.xml',
  grid: NBFACIN + 'Section\\InputCoverageFire.xml',
  form: NBFACIN + 'Section\\CoverageItem.xml',
  premi: NBFACIN + 'Activity\\CountPremi_ACT.xml',
  otomatis: NBFACIN + 'Activity\\AddCoverageAutoFire.xml',
}
const adaCov = Object.values(BERKAS_COV).every((b) => existsSync(b))

describe.skipIf(!adaCov)('tab Coverage FIRE = korpus (CoverageList -> PropertyItemListCoverage -> InputCoverageFire -> CoverageItem) - tiket 43', () => {
  const baca = (b: string) => (adaCov ? readFileSync(b, 'utf-8') : '')
  const ada = (xml: string, sel: string, tag: string, label: string) => blokSel(xml, sel).some((b) => b.includes(`<${tag}>${label}</${tag}>`))
  const objek = baca(BERKAS_COV.objek)
  const item = baca(BERKAS_COV.item)
  const grid = baca(BERKAS_COV.grid)
  const form = baca(BERKAS_COV.form)

  it('Dtl: tab Coverage + tombol Save', () => {
    const dtl = baca(BERKAS_COV.dtl)
    expect(dtl).toContain('<pyTitle>Coverage</pyTitle>')
    expect(dtl).toContain(`<pyLabel>${SIMPAN_COVERAGE}</pyLabel>`)
  })

  it.each(GRID_COV_OBJEK.map((k) => [k.label, k] as const))('kolom objek %s', (_, k) => {
    expect(ada(objek, k.sel, 'pyValue', k.label)).toBe(true)
  })
  it.each([...GRID_COV_ITEM, ...GRID_COV_TOTAL].map((k) => [k.label, k] as const))('kolom item / total %s', (_, k) => {
    expect(ada(item, k.sel, 'pyValue', k.label)).toBe(true)
  })
  it.each(GRID_COVERAGE.map((k) => [k.label, k] as const))('kolom coverage %s', (_, k) => {
    expect(ada(grid, k.sel, 'pyValue', k.label)).toBe(true)
  })
  it.each(Object.entries(FORM_COV))('form %s', (_, u) => {
    expect(ada(form, u.sel, u.tag, u.label)).toBe(true)
  })

  it('‰ Gross Rate wajib; grid berantai sesuai PageList', () => {
    expect(blokSel(form, FORM_COV.rate.sel).some((b) => b.includes('<pyRequired>true</pyRequired>'))).toBe(true)
    expect(item).toContain('<pyPageListProperty>.Property.PropertyItemList</pyPageListProperty>')
    expect(item).toContain('<pyPageListProperty>.Property.TotalTSIPremiGrossList</pyPageListProperty>')
    expect(grid).toContain('<pyPageListProperty>.CoverageList</pyPageListProperty>')
  })

  it('rumus premi Sum Insured di CountPremi_ACT (rate permil, dasar 1e9)', () => {
    expect(baca(BERKAS_COV.premi)).toContain('@Math.divide((.TSI*.Rate*Local.Prorate*.IndemnityPercentage*Local.LossLimit),1000000000,20)')
  })

  it('AddCoverageAutoFire: lima coverage awal', () => {
    const t = baca(BERKAS_COV.otomatis)
    for (const k of ['100815', '100825', '100828', '100829', '100840']) expect(t).toContain(`<PropertiesValue>"${k}"</PropertiesValue>`)
  })

  it('uji ini menggigit: label salah tidak ditemukan', () => {
    expect(ada(form, '29', 'pyLabelFieldValue', 'Gross Rate')).toBe(false)
  })
})

describe.skipIf(!existsSync(`${DDL}CoverageBasis.xml`))('Coverage Basis = aturan properti (DDL\\CoverageBasis.xml)', () => {
  it('label + lima pilihan berurutan', () => {
    const xml = readFileSync(`${DDL}CoverageBasis.xml`, 'utf-8')
    expect(xml).toContain(`<pyLabel>${LABEL_COVERAGE_BASIS}</pyLabel>`)
    expect(daftarPrompt(xml).filter(([v]) => v !== '')).toEqual(OPSI_COVERAGE_BASIS.map((o) => [o.value, o.label]))
  })
})

/**
 * Cari berkas aturan properti di `DDL\` menurut IDENTITAS rule (`<pxInsName>`), bukan nama berkas - dua aturan berbeda
 * dapat dikirim dengan nama berkas yang sama (03-10-2026: `Condition.xml` Deductible menimpa `Condition.xml`
 * PropertyItem). Kosong bila tidak ada.
 */
function aturanDDL(pxInsName: string): string {
  if (!existsSync(DDL)) return ''
  for (const f of readdirSync(DDL)) {
    if (!f.toLowerCase().endsWith('.xml')) continue
    const xml = readFileSync(`${DDL}${f}`, 'utf-8')
    if (xml.includes(`<pxInsName>${pxInsName}</pxInsName>`)) return xml
  }
  return ''
}

describe('Condition Object Item = ASM-FW-GISFW-DATA-PROPERTYITEM!CONDITION (dicari menurut pxInsName)', () => {
  const xml = aturanDDL('ASM-FW-GISFW-DATA-PROPERTYITEM!CONDITION')
  it.skipIf(xml === '')('pasangan value/label per rowdata; baris pertama = Please Select', () => {
    const baris = daftarPrompt(xml)
    expect(baris[0]).toEqual(['', CONDITION_KOSONG])
    expect(baris.slice(1)).toEqual(OPSI_CONDITION.map((o) => [o.value, o.label]))
  })
})

describe('Deductible MinMax / Condition = aturan DATA-DEDUCTIBLE (tahap C3)', () => {
  const minMax = aturanDDL('ASM-FW-GISFW-DATA-DEDUCTIBLE!MINMAX')
  const kondisi = aturanDDL('ASM-FW-GISFW-DATA-DEDUCTIBLE!CONDITION')
  it.skipIf(minMax === '')('MinMax = Min / Max / Or', () => {
    expect(daftarPrompt(minMax)).toEqual(OPSI_MINMAX.map((o) => [o.value, o.label]))
  })
  it.skipIf(kondisi === '')('Condition deductible = lima pilihan', () => {
    expect(daftarPrompt(kondisi)).toEqual(OPSI_KONDISI_DEDUCTIBLE.map((o) => [o.value, o.label]))
  })
})

describe.skipIf(!existsSync(NBFACIN + 'Section\\addDeductible.xml'))('Deductible = korpus (addDeductible, CoverageItem) - tiket 45', () => {
  const form = readFileSync(NBFACIN + 'Section\\addDeductible.xml', 'utf-8')
  const cov = readFileSync(NBFACIN + 'Section\\CoverageItem.xml', 'utf-8')
  it.each(Object.entries(FORM_DEDUCTIBLE))('form %s', (_, u) => {
    expect(blokSel(form, u.sel).some((b) => b.includes(`<pyLabelFieldValue>${u.label}</pyLabelFieldValue>`))).toBe(true)
  })
  it('syarat tampil: Pct / MinMax bila Type != 7; Type2 / Pct2 bila MinMax = 3; Condition teks bila Condition = 5', () => {
    expect(blokSel(form, '6').some((b) => b.includes('<pyCondition>.TypeDeductible != 7</pyCondition>'))).toBe(true)
    expect(blokSel(form, '9').some((b) => b.includes('<pyCondition>.MinMax = 3</pyCondition>'))).toBe(true)
    expect(blokSel(form, '13').some((b) => b.includes('<pyCondition>.Condition = 5</pyCondition>'))).toBe(true)
    expect(form).toContain(`<pyNoSelectionText>${DEDUCTIBLE_KOSONG}</pyNoSelectionText>`)
  })
  it('grid CoverageItem: judul Deductible, kolom Time Excess (Days), PageList .DeductibleList', () => {
    expect(cov).toContain(`>${JUDUL_DEDUCTIBLE}<`)
    expect(cov).toContain(`<pyValue>${KOLOM_TIME_EXCESS}</pyValue>`)
    expect(cov).toContain('<pyPageListProperty>.DeductibleList</pyPageListProperty>')
  })
})

describe('Type Deductible (1/2) = aturan DATA-DEDUCTIBLE (dicari menurut pxInsName)', () => {
  const t1 = aturanDDL('ASM-FW-GISFW-DATA-DEDUCTIBLE!TYPEDEDUCTIBLE')
  const t2 = aturanDDL('ASM-FW-GISFW-DATA-DEDUCTIBLE!TYPEDEDUCTIBLE2')
  it.skipIf(t1 === '')('TypeDeductible', () => {
    expect(daftarPrompt(t1)).toEqual(OPSI_TYPE_DEDUCTIBLE.map((o) => [o.value, o.label]))
  })
  it.skipIf(t2 === '')('TypeDeductible2', () => {
    expect(daftarPrompt(t2)).toEqual(OPSI_TYPE_DEDUCTIBLE2.map((o) => [o.value, o.label]))
  })
})

describe.skipIf(!existsSync(NBFACIN + 'Section\\CoverageItem.xml'))('Form Coverage - kontrol & syarat (tiket 46)', () => {
  const form = readFileSync(NBFACIN + 'Section\\CoverageItem.xml', 'utf-8')
  const sel = (n: string) => blokSel(form, n).join('\n')
  it('Days = radio mendatar', () => {
    expect(sel('26')).toContain('<pyFormat>pxRadioButtons</pyFormat>')
    expect(sel('26')).toContain('<pyOrientation>horizontal</pyOrientation>')
  })
  it('Indemnity / Indemnity Unit / % Loss Limit ALWAYS (syarat sisa diabaikan); First Loss / EML OTHER', () => {
    for (const n of ['28', '40', '44']) expect(sel(n)).toContain('<pyVisible>ALWAYS</pyVisible>')
    for (const n of ['30', '42', '45']) expect(sel(n)).toContain('<pyVisible>OTHER</pyVisible>')
  })
  it('Copy Accumulation hanya coverage pertama; Choose Accumulation Code -> SearchAccumAct', () => {
    expect(sel('19')).toContain('<pyCondition>.pxListSubscript==1</pyCondition>')
    expect(sel('19')).toContain('<pyActivity>CopyAccumulationCode_Act</pyActivity>')
    expect(sel('18')).toContain('<pyActivity>SearchAccumAct</pyActivity>')
  })
})

describe.skipIf(!existsSync(NBFACIN + 'Section\\SearchRiskAccumCov.xml'))('Popup Choose Accumulation = SearchRiskAccumCov (tiket 46)', () => {
  const xml = readFileSync(NBFACIN + 'Section\\SearchRiskAccumCov.xml', 'utf-8')
  const cov = readFileSync(NBFACIN + 'Section\\CoverageItem.xml', 'utf-8')
  const ada = (sel: string, tag: string, label: string) => blokSel(xml, sel).some((b) => b.includes(`<${tag}>${label}</${tag}>`))
  it.each(
    [P_AKUM.cari, P_AKUM.accumulationCode, P_AKUM.road, P_AKUM.zipCode, P_AKUM.czone, P_AKUM.filter, P_AKUM.bersih, ...P_AKUM.kolom, P_AKUM.pilih].map(
      (u) => [u.label, u] as const,
    ),
  )('%s', (_, u) => {
    expect(ada(u.sel, u.tag, u.label)).toBe(true)
  })
  it('judul jendela + RD grid sel 90', () => {
    expect(cov).toContain(`<pyWindowName>${P_AKUM.judul}</pyWindowName>`)
    expect(xml).toContain('SearchRiskAccumulation_RD')
  })
})
