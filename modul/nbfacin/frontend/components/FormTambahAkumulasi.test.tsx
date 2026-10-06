// Add New accumulation (tiket 46) - data uji sintetis.

import { existsSync, readFileSync } from 'node:fs'
import { join } from 'node:path'
import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'

import { TAMBAH_AKUMULASI as T, TEKS_TAMBAH_AKUMULASI as TEKS } from '../labels'
import FormTambahAkumulasi, { keBadan, periksaTambah, TAMBAH_KOSONG, UKURAN_HALAMAN_ZIP } from './FormTambahAkumulasi'

const SUMBER_POPUP = readFileSync(join(__dirname, 'PopupAkumulasi.tsx'), 'utf8').replace(/\r\n/g, '\n')
const KORPUS = 'D:\\migrasi\\RNM\\NB FacIn\\Section\\'

/** Isi `rowdata` sel ber-`pyCellId` tertentu (rowdata terdalam). */
function blokSel(xml: string, sel: string): string[] {
  const hasil: string[] = []
  const tumpukan: { mulai: number; cocok: boolean }[] = []
  for (const m of xml.matchAll(/<rowdata\b[^>]*>|<\/rowdata>|<pyCellId>(\d+)<\/pyCellId>/g)) {
    const i = m.index ?? 0
    if (m[0].startsWith('<rowdata')) tumpukan.push({ mulai: i, cocok: false })
    else if (m[0] === '</rowdata>') {
      const atas = tumpukan.pop()
      if (atas?.cocok) hasil.push(xml.slice(atas.mulai, i + m[0].length))
    } else if (m[1] === sel) {
      const atas = tumpukan.at(-1)
      if (atas) atas.cocok = true
    }
  }
  return hasil
}

describe.skipIf(!existsSync(KORPUS + 'InputAccumulationCov.xml'))('label Add New = korpus', () => {
  const form = readFileSync(KORPUS + 'InputAccumulationCov.xml', 'utf-8')
  const cari = readFileSync(KORPUS + 'SearchRiskAccumCov.xml', 'utf-8')
  const ada = (xml: string, u: { sel: string; tag: string; label: string }) =>
    blokSel(xml, u.sel).some((b) => b.includes(`<${u.tag}>${u.label}</${u.tag}>`))
  it('tombol Add New (SearchRiskAccumCov sel 39)', () => {
    expect(ada(cari, T.tambahBaru)).toBe(true)
  })
  it.each(
    [T.accumulationType, T.province, T.scopeArea, T.crestaZone, T.primaryZip, T.pilihZip, T.keyword, T.note, T.simpan, T.tutup].map(
      (u) => [u.label, u] as const,
    ),
  )('%s (InputAccumulationCov)', (_, u) => {
    expect(ada(form, u)).toBe(true)
  })
  it('kolom Choose Zip Code (ChooseZipCodeDtl)', () => {
    const zip = readFileSync(KORPUS + 'ChooseZipCodeDtl.xml', 'utf-8')
    for (const k of T.kolomZip) expect(zip).toContain(`<pyValue>${k}</pyValue>`)
  })
})

describe('FormTambahAkumulasi', () => {
  it('wajib: tujuh isian bertanda * (Accumulation Type, Province, Key Word, Description, Scope Area, Cresta Zone, Zip)', () => {
    expect(Object.keys(periksaTambah(TAMBAH_KOSONG)).sort()).toEqual(
      ['accumulationType', 'cZone', 'keyword', 'note', 'provinceName', 'scopeArea', 'zipCode'].sort(),
    )
  })
  it('Accumulation Type / Province yang hanya diketik (belum dipilih dari saran) ditolak', () => {
    const g = periksaTambah({ ...TAMBAH_KOSONG, accumulationType: 'UJI', provinceName: 'UJI' })
    expect(g.accumulationType).toBe(TEKS.pilihDariSaran)
    expect(g.provinceName).toBe(TEKS.pilihDariSaran)
    const ok = periksaTambah({ ...TAMBAH_KOSONG, accumulationType: 'UJI', accumulation: 'T1', provinceName: 'UJI', provinceId: 'P1' })
    expect(ok.accumulationType).toBeUndefined()
    expect(ok.provinceName).toBeUndefined()
  })
  it('badan kiriman tanpa teks tampil, dipangkas', () => {
    const b = keBadan({
      ...TAMBAH_KOSONG,
      accumulation: 'T1',
      accumulationType: ' UJI TIPE ',
      accumulationName: 'NAMA TAMPIL',
      note: ' jalan uji ',
      provinceId: 'P1',
      provinceName: 'PROV UJI',
      zipCode: '12345',
      negara: 'UJI',
    })
    expect(b).toMatchObject({ accumulation: 'T1', accumulationType: 'UJI TIPE', note: 'jalan uji', provinceId: 'P1', zipCode: '12345', negara: 'UJI' })
    expect(b).not.toHaveProperty('accumulationName')
    expect(b).not.toHaveProperty('provinceName')
  })
  it('render: Primary Zip Code baca-saja "---" + tombol Choose Zip Code; Save / Close', () => {
    const html = renderToStaticMarkup(<FormTambahAkumulasi zipSaring="" onTutup={() => {}} onTersimpan={() => {}} />)
    expect(html).toContain(TEKS.judul)
    expect(html).toMatch(/value="---"[^>]*readOnly|readOnly[^>]*value="---"/i)
    for (const t of [T.pilihZip.label, T.simpan.label, T.tutup.label, T.note.label]) expect(html).toContain(t)
  })
  it('popup: Add New menggantikan kartu cari; Save berhasil -> coverage memakai accumulation baru', () => {
    expect(SUMBER_POPUP).toContain('{TAMBAH_AKUMULASI.tambahBaru.label}')
    expect(SUMBER_POPUP).toMatch(/onTersimpan=\{\(id, note\) => onPilih\(id, note\)\}/)
    expect(TEKS.tidakLengkap).toBe("Postal code, Nation, CZone and Accumulation Description can't be null!")
  })
})

describe('Choose Zip Code - 15 per halaman, menurut Province terpilih (04-10-2026)', () => {
  const SUMBER = readFileSync(join(__dirname, 'FormTambahAkumulasi.tsx'), 'utf8').replace(/\r\n/g, '\n')
  it('ukuran 15; data dari server per halaman menurut Province form', () => {
    expect(UKURAN_HALAMAN_ZIP).toBe(15)
    expect(SUMBER).toContain('cariZipAkumulasi(provinceName, q, halaman)')
    expect(SUMBER).toContain('provinceName={s.provinceName}')
  })
  it('Choose Zip Code = popup di atas popup (portal ke body, akar nbfacin, Escape capture hanya menutup popup zip)', () => {
    expect(SUMBER).toContain('return createPortal(')
    expect(SUMBER).toContain('document.body')
    expect(SUMBER).toContain('<div className="nbfacin">')
    expect(SUMBER).toContain("window.addEventListener('keydown', onKey, true)")
    expect(SUMBER).toContain('e.stopPropagation()')
  })
  it('memilih zip SELALU mengganti Cresta Zone (RW.CZONE zip itu); awal hanya mengisi bila ada', () => {
    expect(SUMBER).toContain('isiCzone(b.zipCode, true)')
    expect(SUMBER).toContain('isiCzone(zipSaring, false)')
    expect(SUMBER).toContain("if (ganti || h.cZone !== '') setS((x) => ({ ...x, cZone: h.cZone, cZoneId: h.cZoneId }))")
  })
  it('kata cari / province berubah -> halaman 1; pager bila total > ukuran', () => {
    expect(SUMBER).toMatch(/useEffect\(\(\) => \{\s*setHalaman\(1\)\s*\}, \[q, provinceName\]\)/)
    expect(SUMBER).toContain('data.total > data.ukuran')
  })
})
