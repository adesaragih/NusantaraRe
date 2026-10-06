import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambil, ambilDaftar, ambilRate, hapus, pratinjauUnggah, simpanUnggah, tambah, ubah } from './api'
import {
  berkasSah,
  gantiUrut,
  jumlahHalaman,
  KOLOM_URUT,
  MAKS_BYTES_CSV,
  nomorAwal,
  periksaNama,
  SARINGAN_AWAL,
  tandaUrut,
  UKURAN_HALAMAN,
} from './aturan'
import { MENU_RR, RR } from './labels'

const XML_LABEL = [
  'R/I RATE SUMMARY',
  'R/I RATE NAME',
  'MODIFY OPERATOR',
  'MODIFY DATE',
  'Save',
  'Cancel',
  'Upload CSV',
  'View Upload',
  'Simpan Upload',
  'Edit',
  'Detail',
  'Delete',
  'Rate Detail',
  'Format excel : USEDBY, CONTRACT, GENDER, AGE, RATE',
]

describe('label R/I Rate Life', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 922, golongan MASTER TREATY URUTAN 9, langsung menyala', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/922_m_nav_menu_riratelife.sql`, 'utf8')
    expect(sql).toContain(`'riratelife', '${MENU_RR.kelompok}', 'MASTER TREATY', 'riratelife', 9, '1'`)
  })

  it('label layar VERBATIM dari RIRate.xml', () => {
    const nilai = Object.values(RR).filter((v) => typeof v === 'string')
    for (const l of XML_LABEL) expect(nilai, l).toContain(l)
  })
})

describe('aturan R/I Rate Life', () => {
  it('urut: kolom baru menaik, kolom sama balik arah; halaman kembali 1', () => {
    const a = gantiUrut({ ...SARINGAN_AWAL, halaman: 3 }, 'usedby')
    expect(a).toMatchObject({ urut: 'usedby', turun: false, halaman: 1 })
    expect(gantiUrut(a, 'usedby').turun).toBe(true)
    expect(gantiUrut(gantiUrut(a, 'usedby'), 'id')).toMatchObject({ urut: 'id', turun: false })
    expect(tandaUrut(a, 'usedby')).toBe(' ▲')
    expect(tandaUrut(a, 'id')).toBe('')
    expect(KOLOM_URUT).toEqual(['id', 'usedby', 'operatorid', 'modifieddate'])
  })

  it('halaman 50 baris (pyPageSize b12645)', () => {
    expect(UKURAN_HALAMAN).toBe(50)
    expect(jumlahHalaman(0)).toBe(1)
    expect(jumlahHalaman(50)).toBe(1)
    expect(jumlahHalaman(51)).toBe(2)
    expect(nomorAwal(2)).toBe(51)
  })

  it('nama wajib; berkas .csv sampai 4 MB', () => {
    expect(periksaNama('  ')).toBe(RR.galatNama)
    expect(periksaNama('UJI RATE')).toBeNull()
    expect(berkasSah('uji.CSV', 10)).toBe(true)
    expect(berkasSah('uji.xlsx', 10)).toBe(false)
    expect(berkasSah('uji.csv', 0)).toBe(false)
    expect(berkasSah('uji.csv', MAKS_BYTES_CSV + 1)).toBe(false)
  })

  it('View only: form, Upload, Edit, Delete hanya bila bolehUbah; Detail selalu', () => {
    const halaman = readFileSync(`${__dirname}/pages/RIRateLife.tsx`, 'utf8')
    expect(halaman).toContain('useBolehUbah(NAMA_RR)')
    for (const label of ['RR.edit', 'RR.hapus']) {
      const i = halaman.indexOf(`{${label}}`)
      expect(halaman.lastIndexOf('{bolehUbah && (', i), label).toBeGreaterThan(halaman.lastIndexOf('</button>', i))
    }
    expect(halaman.match(/\{bolehUbah && \(/g)?.length).toBe(4)
  })

  it('pesan Delete menyebut jumlah rate yang ikut terhapus', () => {
    expect(RR.tanyaHapus('UJI', 3)).toBe('Delete "UJI" and its 3 rate rows? This cannot be undone.')
    expect(RR.tanyaHapus('UJI', 1)).toContain('1 rate row?')
  })
})

describe('klien R/I Rate Life', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan handlers/rute.go', async () => {
    const panggil: string[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')} ${String(init.body ?? '')}`.trim())
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar(SARINGAN_AWAL)
    await ambilDaftar({ id: ' 10 ', usedby: 'uji', urut: 'usedby', turun: true, halaman: 2 })
    await ambil('101')
    await tambah('UJI')
    await ubah('101', 'UJI 2')
    await hapus('101')
    await ambilRate('101', 1)
    await ambilRate('101', 2)
    await pratinjauUnggah('a;b')
    await simpanUnggah('a;b')
    expect(panggil).toEqual([
      'GET /api/ri-rate-life',
      'GET /api/ri-rate-life?id=10&usedby=uji&urut=usedby&arah=desc&halaman=2',
      'GET /api/ri-rate-life/101',
      'POST /api/ri-rate-life {"usedby":"UJI"}',
      'PUT /api/ri-rate-life/101 {"usedby":"UJI 2"}',
      'DELETE /api/ri-rate-life/101',
      'GET /api/ri-rate-life/101/rate',
      'GET /api/ri-rate-life/101/rate?halaman=2',
      'POST /api/ri-rate-life/unggah/pratinjau {"csv":"a;b"}',
      'POST /api/ri-rate-life/unggah {"csv":"a;b"}',
    ])
  })
})
