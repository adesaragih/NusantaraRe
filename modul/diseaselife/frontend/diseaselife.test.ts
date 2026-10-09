import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilDaftar, tambah, ubah } from './api'
import {
  BATAS_DISEASE,
  BATAS_ICD,
  hurufBesar,
  jumlahHalaman,
  periksaIsian,
  pilihUrut,
  SARINGAN_AWAL,
  tandaArah,
  UKURAN_HALAMAN,
} from './aturan'
import { DSL, MENU_DSL } from './labels'

const HALAMAN = readFileSync(`${__dirname}/pages/DiseaseLife.tsx`, 'utf8')
/** Kode saja - komentar kepala berkas menyebut apa yang SENGAJA tidak dibuat. */
const KODE = HALAMAN.split('\n')
  .filter((b) => !b.trim().startsWith('//'))
  .join('\n')

/** Label VERBATIM dari `InboxDisease.xml` (nomor baris di labels.ts). */
const XML_LABEL = ['DISEASE', 'Number / ID', 'ID', 'ICD Code', 'Disease', 'Edit', 'Save', 'Cancel']

describe('label Disease Life', () => {
  it('nama menu = M_NAV_MENU.LABEL slot menu modul 951 (D4), golongan MASTER TREATY URUTAN 15, langsung menyala', () => {
    const sql = readFileSync(`${__dirname}/../backend/migrations/951_menu_diseaselife.sql`, 'utf8')
    expect(MENU_DSL.kelompok).toBe('Disease Life')
    expect(sql).toContain(`'diseaselife', '${MENU_DSL.kelompok}', 'MASTER TREATY', 'diseaselife', 15, '1'`)
  })

  it('label layar VERBATIM dari XML', () => {
    const nilai = Object.values(DSL).filter((v) => typeof v === 'string')
    for (const l of XML_LABEL) expect(nilai, l).toContain(l)
  })
})

describe('aturan Disease Life', () => {
  it('grid 10 baris (b5169), ID menurun bawaan (b5012); kepala ID / ICD Code mengurut dan kembali ke halaman 1', () => {
    expect(UKURAN_HALAMAN).toBe(10)
    expect(SARINGAN_AWAL).toEqual({ icdCode: '', disease: '', urut: 'id', naik: false, halaman: 1 })
    expect(tandaArah(SARINGAN_AWAL, 'id')).toBe(' ▼')
    expect(tandaArah(SARINGAN_AWAL, 'icd')).toBe('')
    const a = pilihUrut({ ...SARINGAN_AWAL, halaman: 3 }, 'id')
    expect(a).toMatchObject({ urut: 'id', naik: true, halaman: 1 })
    const b = pilihUrut(a, 'icd')
    expect(b).toMatchObject({ urut: 'icd', naik: true, halaman: 1 })
    expect(pilihUrut(b, 'icd')).toMatchObject({ urut: 'icd', naik: false })
    expect(pilihUrut(b, 'id')).toMatchObject({ urut: 'id', naik: false })
    expect(jumlahHalaman(0)).toBe(1)
    expect(jumlahHalaman(97586)).toBe(9759)
  })

  it('ICD Code (D3) dan Disease (b1465) wajib, huruf besar (SetUpperCase_DT b1265 / b1538), batas lebar kolom', () => {
    expect(periksaIsian('  ', 'uji')).toBe(DSL.galatICD)
    expect(periksaIsian('uji01', ' ')).toBe(DSL.galatDisease)
    expect(periksaIsian('uji01', 'uji kolera')).toBeNull()
    expect(periksaIsian('A'.repeat(BATAS_ICD), 'A'.repeat(BATAS_DISEASE))).toBeNull()
    expect(periksaIsian('A'.repeat(BATAS_ICD + 1), 'UJI')).toBe(DSL.galatPanjang('ICD Code', 100))
    expect(periksaIsian('UJI', '€'.repeat(334))).toBe(DSL.galatPanjang('Disease', 1000))
    expect(hurufBesar('uji Kolera')).toBe('UJI KOLERA')
  })
})

describe('layar InboxDisease', () => {
  it('Number / ID tidak dapat diisi (disabled b1070); ICD Code input (pxTextInput b1232) dan Disease textarea (pxTextArea b1502), keduanya huruf besar', () => {
    expect(HALAMAN).toMatch(/\{DSL\.nomorId\}<\/span>\s*<input[^>]*value=\{ubahID\}[^>]*readOnly disabled/)
    expect(HALAMAN).toMatch(/<input\s+className="field__input"\s+type="text"\s+value=\{icdCode\}/)
    expect(HALAMAN).toContain('<textarea')
    expect(HALAMAN).toContain('onBlur={() => setIcdCode((v) => hurufBesar(v))}')
    expect(HALAMAN).toContain('onBlur={() => setDisease((v) => hurufBesar(v))}')
  })

  it('Save selalu; Cancel HANYA saat Edit (DATASHOW = IsEdit b2525); Edit mengisi form (EditList_DT b4852)', () => {
    expect(HALAMAN).toContain('{sibuk ? DSL.menyimpan : DSL.save}')
    expect(HALAMAN).toMatch(/\{ubahID !== '' && \(\s*<button[^>]*onClick=\{kosongkanForm\}>\s*\{DSL\.batal\}/)
    expect(HALAMAN).toContain("ubahID === '' ? tambah(icd, nama) : ubah(ubahID, icd, nama)")
    expect(HALAMAN).toMatch(/setUbahID\(p\.id\)\s*setIcdCode\(p\.icdCode\)\s*setDisease\(p\.disease\)/)
  })

  it('grid ID, ICD Code (dapat diurutkan), Disease (tidak), kolom Edit; saring ICD Code / Disease; nol Delete / Upload / detail', () => {
    const kepala = KODE.slice(KODE.indexOf('<thead>'), KODE.indexOf('</thead>'))
    expect(kepala.indexOf('{DSL.id}')).toBeLessThan(kepala.indexOf('{DSL.icdCode}'))
    expect(kepala.indexOf('{DSL.icdCode}')).toBeLessThan(kepala.indexOf('{DSL.disease}'))
    expect(kepala).toContain("pilihUrut(s, 'id')")
    expect(kepala).toContain("pilihUrut(s, 'icd')")
    expect(kepala).toContain('<th>{DSL.disease}</th>')
    expect(KODE).toContain('placeholder={DSL.saringICD}')
    expect(KODE).toContain('placeholder={DSL.saringDisease}')
    for (const kata of ['Delete', 'hapus', 'Upload', 'unggah', 'Detail', 'rincian']) {
      expect(KODE, kata).not.toContain(kata)
    }
  })

  it('⛔ 97.586 baris: layar hanya meminta SATU halaman bersaring ke server - nol pemuatan seluruh tabel', () => {
    expect(KODE).toContain('ambilDaftar(s)')
    expect(KODE).not.toMatch(/ukuran\s*:\s*\d{3,}|semua|all=/i)
  })

  it('View only: form dan Edit hanya bila bolehUbah', () => {
    expect(HALAMAN).toContain('useBolehUbah(NAMA_DSL)')
    expect(HALAMAN.match(/\{bolehUbah && \(/g)?.length).toBe(2)
    expect(HALAMAN).toContain('{bolehUbah && <th')
  })
})

describe('klien Disease Life', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute, saring, urut, halaman sama dengan handlers/rute.go (tanpa DELETE)', async () => {
    const panggil: string[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')} ${String(init.body ?? '')}`.trim())
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar(SARINGAN_AWAL)
    await ambilDaftar({ icdCode: ' a00 ', disease: 'kolera', urut: 'icd', naik: true, halaman: 2 })
    await tambah('UJI01', 'UJI KOLERA')
    await ubah('100002', 'UJI02', 'UJI DEMAM')
    expect(panggil).toEqual([
      'GET /api/disease-life',
      'GET /api/disease-life?icd=a00&disease=kolera&urut=icd&arah=asc&halaman=2',
      'POST /api/disease-life {"icdCode":"UJI01","disease":"UJI KOLERA"}',
      'PUT /api/disease-life/100002 {"icdCode":"UJI02","disease":"UJI DEMAM"}',
    ])
  })
})
