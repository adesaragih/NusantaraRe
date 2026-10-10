import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilDaftar, tambah, ubah } from './api'
import {
  balikArah,
  BATAS_BENEFIT,
  hurufBesar,
  jumlahHalaman,
  periksaBenefit,
  SARINGAN_AWAL,
  tandaArah,
  UKURAN_HALAMAN,
} from './aturan'
import { BN, MENU_BN } from './labels'

const HALAMAN = readFileSync(`${__dirname}/pages/BenefitLife.tsx`, 'utf8')

/** Label VERBATIM dari `InboxBenefit.xml` (nomor baris di labels.ts). */
const XML_LABEL = ['INSURANCE BENEFIT', 'Number / ID', 'ID', 'Benefit', 'Edit', 'Save', 'Cancel']

describe('label Benefit', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 945 (K4), golongan MASTER TREATY URUTAN 12, langsung menyala', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/945_m_nav_menu_benefitlife.sql`, 'utf8')
    expect(MENU_BN.kelompok).toBe('Benefit')
    expect(sql).toContain(`'benefitlife', '${MENU_BN.kelompok}', 'MASTER TREATY', 'benefitlife', 12, '1'`)
  })

  it('label layar VERBATIM dari XML', () => {
    const nilai = Object.values(BN).filter((v) => typeof v === 'string')
    for (const l of XML_LABEL) expect(nilai, l).toContain(l)
  })
})

describe('aturan Benefit', () => {
  it('grid 10 baris (b4526), ID menurun bawaan (b4391); kepala ID membalik arah dan kembali ke halaman 1', () => {
    expect(UKURAN_HALAMAN).toBe(10)
    expect(SARINGAN_AWAL).toEqual({ id: '', benefit: '', naik: false, halaman: 1 })
    expect(tandaArah(SARINGAN_AWAL)).toBe(' ▼')
    const a = balikArah({ ...SARINGAN_AWAL, halaman: 3 })
    expect(a).toMatchObject({ naik: true, halaman: 1 })
    expect(tandaArah(a)).toBe(' ▲')
    expect(jumlahHalaman(0)).toBe(1)
    expect(jumlahHalaman(11)).toBe(2)
  })

  it('Benefit wajib (b1137), huruf besar (SetUpperCase_DT b1211), paling panjang 200 byte (batas kolom)', () => {
    expect(periksaBenefit('   ')).toBe(BN.galatBenefit)
    expect(periksaBenefit('uji rawat inap')).toBeNull()
    expect(periksaBenefit('A'.repeat(BATAS_BENEFIT))).toBeNull()
    expect(periksaBenefit('A'.repeat(BATAS_BENEFIT + 1))).toBe(BN.galatPanjang(200))
    expect(periksaBenefit('€'.repeat(67))).toBe(BN.galatPanjang(200))
    expect(hurufBesar('uji Rawat')).toBe('UJI RAWAT')
  })
})

describe('layar InboxBenefit', () => {
  it('Number / ID tidak dapat diisi (disabled b1013); Benefit textarea wajib (pxTextArea b1175) yang menjadi huruf besar', () => {
    expect(HALAMAN).toMatch(/\{BN\.nomorId\}<\/span>\s*<input[^>]*value=\{ubahID\}[^>]*readOnly disabled/)
    expect(HALAMAN).toContain('<textarea')
    expect(HALAMAN).toContain('onBlur={() => setBenefit((v) => hurufBesar(v))}')
  })

  it('Save selalu; Cancel HANYA saat Edit (DATASHOW = IsEdit b2198); Edit mengisi form (EditList_DT b4243)', () => {
    expect(HALAMAN).toContain('{sibuk ? BN.menyimpan : BN.save}')
    expect(HALAMAN).toMatch(/\{ubahID !== '' && \(\s*<button[^>]*onClick=\{kosongkanForm\}>\s*\{BN\.batal\}/)
    expect(HALAMAN).toContain("ubahID === '' ? tambah(isi) : ubah(ubahID, isi)")
    expect(HALAMAN).toMatch(/setUbahID\(b\.id\)\s*setBenefit\(b\.benefit\)/)
  })

  it('grid kolom ID, Benefit, lalu kolom Edit; nol Delete / Upload / detail (XML tidak memuatnya)', () => {
    const kepala = HALAMAN.slice(HALAMAN.indexOf('<thead>'), HALAMAN.indexOf('</thead>'))
    expect(kepala.indexOf('{BN.id}')).toBeLessThan(kepala.indexOf('{BN.benefit}'))
    expect(kepala).not.toContain('No<')
    // Kode saja - komentar kepala berkas menyebut apa yang SENGAJA tidak dibuat.
    const kode = HALAMAN.split('\n')
      .filter((b) => !b.trim().startsWith('//'))
      .join('\n')
    for (const kata of ['Delete', 'hapus', 'Upload', 'unggah', 'Detail', 'rincian']) expect(kode, kata).not.toMatch(new RegExp(kata, 'i'))
  })

  it('View only: form dan Edit hanya bila bolehUbah', () => {
    expect(HALAMAN).toContain('useBolehUbah(NAMA_BN)')
    expect(HALAMAN.match(/\{bolehUbah && \(/g)?.length).toBe(2)
    expect(HALAMAN).toContain('{bolehUbah && <th')
  })
})

describe('klien Benefit', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan handlers/rute.go (tanpa DELETE)', async () => {
    const panggil: string[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')} ${String(init.body ?? '')}`.trim())
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar(SARINGAN_AWAL)
    await ambilDaftar({ id: ' 1000 ', benefit: 'uji', naik: true, halaman: 2 })
    await tambah('UJI')
    await ubah('100001', 'UJI 2')
    expect(panggil).toEqual([
      'GET /api/benefit-life',
      'GET /api/benefit-life?id=1000&benefit=uji&arah=asc&halaman=2',
      'POST /api/benefit-life {"benefit":"UJI"}',
      'PUT /api/benefit-life/100001 {"benefit":"UJI 2"}',
    ])
  })
})
