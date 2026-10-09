import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilDaftar, ambilPilihanBenefit, ambilPilihanBusiness, tambah, ubah } from './api'
import {
  gantiUrut,
  ISIAN_KOSONG,
  isianDariPlan,
  jumlahHalaman,
  periksaIsian,
  SARINGAN_AWAL,
  saringBenefit,
  saringBusiness,
  tandaUrut,
  UKURAN_HALAMAN,
} from './aturan'
import { MENU_PL, PL } from './labels'

const HALAMAN = readFileSync(`${__dirname}/pages/PlanLife.tsx`, 'utf8')
const KODE = HALAMAN.split('\n')
  .filter((b) => !b.trim().startsWith('//'))
  .join('\n')

/** Label VERBATIM dari `InboxProductType.xml` (nomor baris di labels.ts). */
const XML_LABEL = ['Plan', 'Plan Name', 'Business', 'Benefit', 'New', 'Save', 'Edit', 'OLDID', 'Note', 'ID']

describe('label Plan', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 949 (K6), MASTER TREATY URUTAN 13, langsung menyala', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/949_m_nav_menu_planlife.sql`, 'utf8')
    expect(MENU_PL.kelompok).toBe('Plan')
    expect(sql).toContain(`'planlife', '${MENU_PL.kelompok}', 'MASTER TREATY', 'planlife', 13, '1'`)
  })

  it('label layar VERBATIM dari XML', () => {
    const nilai = Object.values(PL).filter((v) => typeof v === 'string')
    for (const l of XML_LABEL) expect(nilai, l).toContain(l)
  })
})

describe('aturan Plan', () => {
  it('grid 10 baris (b4401); urut Plan Name / Benefit saja; bawaan backend', () => {
    expect(UKURAN_HALAMAN).toBe(10)
    expect(SARINGAN_AWAL).toEqual({ urut: '', turun: false, halaman: 1 })
    const a = gantiUrut({ ...SARINGAN_AWAL, halaman: 3 }, 'covername')
    expect(a).toEqual({ urut: 'covername', turun: false, halaman: 1 })
    expect(gantiUrut(a, 'covername').turun).toBe(true)
    expect(tandaUrut(a, 'covername')).toBe(' ▲')
    expect(tandaUrut(a, 'benefit')).toBe('')
    expect(jumlahHalaman(11)).toBe(2)
  })

  it('K5: ketiga medan wajib, <= 200 byte', () => {
    expect(periksaIsian(ISIAN_KOSONG)).toBe(PL.galatWajib('Plan Name'))
    expect(periksaIsian({ coverName: 'UJI', business: ' ', benefit: 'UJI' })).toBe(PL.galatWajib('Business'))
    expect(periksaIsian({ coverName: 'UJI', business: 'UJI', benefit: '' })).toBe(PL.galatWajib('Benefit'))
    expect(periksaIsian({ coverName: 'A'.repeat(201), business: 'UJI', benefit: 'UJI' })).toBe(PL.galatPanjang('Plan Name', 200))
    expect(periksaIsian({ coverName: 'UJI', business: 'UJI', benefit: 'UJI' })).toBeNull()
  })

  it('autocomplete: cari memuat tanpa beda huruf pada Note / Benefit (pyUseForSearch b1230 / b1550)', () => {
    const biz = [
      { id: '9001', oldId: 'L1', note: 'UJI KREDIT' },
      { id: '9002', oldId: 'L2', note: 'UJI Jiwa' },
    ]
    expect(saringBusiness(biz, 'jiwa').map((b) => b.id)).toEqual(['9002'])
    expect(saringBusiness(biz, 'L1')).toEqual([])
    expect(saringBusiness(biz, '')).toHaveLength(2)
    expect(saringBenefit([{ id: '100001', benefit: 'UJI RAWAT INAP' }], 'rawat')).toHaveLength(1)
    expect(isianDariPlan({ coverName: 'A', business: 'B', benefit: 'C' })).toEqual({ coverName: 'A', business: 'B', benefit: 'C' })
  })
})

describe('layar InboxProductType', () => {
  it('Plan Name teks; Business dan Benefit autocomplete teks bebas (b1125 / b1478) dengan kolom OLDID / Note / ID', () => {
    expect(KODE).toContain('{PL.oldId}')
    expect(KODE).toContain('{PL.note}')
    expect(KODE).toContain('setIsi((s) => ({ ...s, business: b.note }))')
    expect(KODE).toContain('setIsi((s) => ({ ...s, benefit: b.benefit }))')
    expect(KODE).not.toMatch(/businessId|benefitId/)
  })

  it('Save selalu; New HANYA saat Edit (DATASHOW = IsEdit b6827); Edit mengisi form', () => {
    expect(KODE).toContain('{sibuk ? PL.menyimpan : PL.save}')
    expect(KODE).toMatch(/\{ubahID !== '' && \(\s*<button[^>]*onClick=\{kosongkanForm\}>\s*\{PL\.baru\}/)
    expect(KODE).toContain("ubahID === '' ? tambah(kirim) : ubah(ubahID, kirim)")
    expect(KODE).toMatch(/setUbahID\(p\.id\)\s*setIsi\(isianDariPlan\(p\)\)/)
  })

  it('grid Plan Name, Business, Benefit, Edit; tanpa kolom ID, saring, Delete, Upload, detail', () => {
    const kepala = KODE.slice(KODE.indexOf('<thead>'), KODE.indexOf('</thead>'))
    expect(kepala.indexOf('{PL.planName}')).toBeLessThan(kepala.indexOf('{PL.business}'))
    expect(kepala.indexOf('{PL.business}')).toBeLessThan(kepala.indexOf('{PL.benefit}'))
    expect(kepala).not.toContain('{PL.id}')
    expect(kepala).not.toMatch(/gantiUrut\(s, 'business'\)/)
    for (const kata of ['Delete', 'hapus', 'Upload', 'unggah', 'Detail', 'type="search"']) expect(KODE, kata).not.toContain(kata)
  })

  it('View only: form dan Edit hanya bila bolehUbah', () => {
    expect(KODE).toContain('useBolehUbah(NAMA_PL)')
    expect(KODE.match(/\{bolehUbah && \(/g)?.length).toBe(2)
    expect(KODE).toContain('{bolehUbah && <th')
  })
})

describe('klien Plan', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan handlers/rute.go (tanpa DELETE)', async () => {
    const panggil: string[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')} ${String(init.body ?? '')}`.trim())
        return Promise.resolve(new Response('[]', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar(SARINGAN_AWAL)
    await ambilDaftar({ urut: 'benefit', turun: true, halaman: 2 })
    await ambilPilihanBusiness()
    await ambilPilihanBenefit()
    await tambah({ coverName: 'UJI', business: 'UJI KREDIT', benefit: 'UJI RAWAT INAP' })
    await ubah('100001', { coverName: 'UJI 2', business: 'UJI KREDIT', benefit: 'UJI RAWAT INAP' })
    expect(panggil).toEqual([
      'GET /api/plan-life',
      'GET /api/plan-life?urut=benefit&arah=desc&halaman=2',
      'GET /api/plan-life/pilihan/business',
      'GET /api/plan-life/pilihan/benefit',
      'POST /api/plan-life {"coverName":"UJI","business":"UJI KREDIT","benefit":"UJI RAWAT INAP"}',
      'PUT /api/plan-life/100001 {"coverName":"UJI 2","business":"UJI KREDIT","benefit":"UJI RAWAT INAP"}',
    ])
  })
})
