// Penjaga bukti label NB FacIn - tiket 21; portal Opportunity tiket 25.
//
// ⛔ Komentar bukti yang tidak pernah diperiksa adalah HIASAN. Berkas ini membuka
// section korpus dan memastikan setiap sel yang disebut `labels.ts` memang berlabel
// dan berproperti seperti tertulis. Korpus READ-ONLY - hanya dibaca.
//
// ⚠️ Bila korpus tidak terjangkau, test DILEWATI dengan pesan - bukan gagal.

import { existsSync, readFileSync } from 'node:fs'

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
