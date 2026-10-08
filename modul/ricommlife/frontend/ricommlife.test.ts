import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambil, ambilDaftar, ambilDetail, hapus, pratinjauUnggah, simpanUnggah, tambah, tambahDetail, ubah, ubahDetail } from './api'
import {
  berkasSah,
  DIGIT_YEAR,
  gantiUrut,
  isianDariKomisi,
  ISIAN_KOMISI_KOSONG,
  jumlahHalaman,
  KOLOM_URUT,
  MAKS_BYTES_CSV,
  nomorAwal,
  periksaIsianKomisi,
  periksaNama,
  SARINGAN_AWAL,
  tampilDesimal,
  tandaUrut,
  UKURAN_HALAMAN,
} from './aturan'
import { MENU_RC, RC } from './labels'

const XML_LABEL = [
  // InboxSummaryRIComm.xml
  'R/I COMM SUMMARY',
  'USEDBY',
  'R/I COMM NAME',
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
  'Format excel : USEDBY, CONTRACT, YEAR, COMM',
  // InboxRIComm.xml
  'R/I COMM DETAIL',
  'Clear Field',
  'CONTRACT',
  'YEAR',
  'COMM',
  '0',
]

describe('label R/I Comm Life', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 925, golongan MASTER TREATY URUTAN 10, langsung menyala', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/925_m_nav_menu_ricommlife.sql`, 'utf8')
    expect(sql).toContain(`'ricommlife', '${MENU_RC.kelompok}', 'MASTER TREATY', 'ricommlife', 10, '1'`)
  })

  it('label layar VERBATIM dari kedua XML', () => {
    const nilai = Object.values(RC).filter((v) => typeof v === 'string')
    for (const l of XML_LABEL) expect(nilai, l).toContain(l)
  })
})

describe('aturan R/I Comm Life', () => {
  it('urut: kolom baru menaik, kolom sama balik arah; halaman kembali 1; bawaan backend (ID menaik)', () => {
    const a = gantiUrut({ ...SARINGAN_AWAL, halaman: 3 }, 'usedby')
    expect(a).toMatchObject({ urut: 'usedby', turun: false, halaman: 1 })
    expect(gantiUrut(a, 'usedby').turun).toBe(true)
    expect(tandaUrut(a, 'usedby')).toBe(' ▲')
    expect(tandaUrut(a, 'id')).toBe('')
    expect(KOLOM_URUT).toEqual(['id', 'usedby', 'operatorid', 'modifieddate'])
    expect(SARINGAN_AWAL.urut).toBe('')
  })

  it('halaman 50 baris (pyPageSize ringkasan b10081, detail b9716)', () => {
    expect(UKURAN_HALAMAN).toBe(50)
    expect(jumlahHalaman(0)).toBe(1)
    expect(jumlahHalaman(51)).toBe(2)
    expect(nomorAwal(2)).toBe(51)
  })

  it('nama wajib; berkas .csv sampai 4 MB', () => {
    expect(periksaNama('  ')).toBe(RC.galatNama)
    expect(periksaNama('UJI COMM')).toBeNull()
    expect(berkasSah('uji.CSV', 10)).toBe(true)
    expect(berkasSah('uji.xlsx', 10)).toBe(false)
    expect(berkasSah('uji.csv', MAKS_BYTES_CSV + 1)).toBe(false)
  })

  it('detail: CONTRACT, YEAR, COMM wajib; YEAR tepat 4 angka; COMM tidak negatif tanpa batas 100; Edit mengisi form berkoma desimal', () => {
    expect(DIGIT_YEAR).toBe(4)
    expect(periksaIsianKomisi(ISIAN_KOMISI_KOSONG)).toBe(RC.galatContract)
    expect(periksaIsianKomisi({ contract: '1', year: '', comm: '' })).toBe(RC.galatYear)
    expect(periksaIsianKomisi({ contract: '1', year: '2026', comm: ' ' })).toBe(RC.galatComm)
    expect(periksaIsianKomisi({ contract: '0', year: '2026', comm: '0,5' })).toBeNull()
    expect(periksaIsianKomisi({ contract: '0', year: '26', comm: '1' })).toBe(RC.galatYearEmpat)
    expect(periksaIsianKomisi({ contract: '0', year: '0026', comm: '1' })).toBe(RC.galatYearEmpat)
    expect(periksaIsianKomisi({ contract: '0', year: '2026', comm: '-1' })).toBe(RC.galatCommNegatif)
    expect(periksaIsianKomisi({ contract: '0', year: '2026', comm: '250' })).toBeNull()
    expect(isianDariKomisi({ id: '1000044', idUsedBy: '1000003', usedby: 'UJI', contract: ' 2', year: '2021 ', comm: '12.5' })).toEqual({
      contract: '2',
      year: '2021',
      comm: '12,5',
    })
    expect(tampilDesimal('0.125')).toBe('0,125')
    expect(tampilDesimal('7')).toBe('7')
  })

  it('View only: form, Upload, Edit, Delete hanya bila bolehUbah; Detail selalu', () => {
    const halaman = readFileSync(`${__dirname}/pages/RICommLife.tsx`, 'utf8')
    expect(halaman).toContain('useBolehUbah(NAMA_RC)')
    for (const label of ['RC.edit', 'RC.hapus']) {
      const i = halaman.indexOf(`{${label}}`)
      expect(halaman.lastIndexOf('{bolehUbah && (', i), label).toBeGreaterThan(halaman.lastIndexOf('</button>', i))
    }
    expect(halaman.match(/\{bolehUbah && \(/g)?.length).toBe(4)
    expect(halaman).toContain('bolehUbah={bolehUbah}')
  })

  it('detail View only: form dan Edit hanya bila bolehUbah; nama R/I dari ringkasan, bukan isian', () => {
    const detail = readFileSync(`${__dirname}/components/KomisiDetail.tsx`, 'utf8')
    expect(detail.match(/\{bolehUbah && \(/g)?.length).toBe(2)
    expect(detail).toContain('value={ringkasan.usedby} readOnly disabled')
    expect(detail).not.toMatch(/usedby:\s*isi/)
  })

  it('pesan Delete menyebut jumlah rincian yang ikut terhapus', () => {
    expect(RC.tanyaHapus('UJI', 3)).toBe('Delete "UJI" and its 3 detail rows? This cannot be undone.')
    expect(RC.tanyaHapus('UJI', 1)).toContain('1 detail row?')
  })
})

describe('klien R/I Comm Life', () => {
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
    await ambil('1000003')
    await tambah('UJI')
    await ubah('1000003', 'UJI 2')
    await hapus('1000003')
    await ambilDetail('1000003', 1)
    await ambilDetail('1000003', 2)
    await tambahDetail('1000003', { contract: '1', year: '2026', comm: '0,5' })
    await ubahDetail('1000003', '1000044', { contract: '2', year: '2025', comm: '1' })
    await pratinjauUnggah('a;b')
    await simpanUnggah('a;b')
    expect(panggil).toEqual([
      'GET /api/ri-comm-life',
      'GET /api/ri-comm-life?id=10&usedby=uji&urut=usedby&arah=desc&halaman=2',
      'GET /api/ri-comm-life/1000003',
      'POST /api/ri-comm-life {"usedby":"UJI"}',
      'PUT /api/ri-comm-life/1000003 {"usedby":"UJI 2"}',
      'DELETE /api/ri-comm-life/1000003',
      'GET /api/ri-comm-life/1000003/detail',
      'GET /api/ri-comm-life/1000003/detail?halaman=2',
      'POST /api/ri-comm-life/1000003/detail {"contract":"1","year":"2026","comm":"0,5"}',
      'PUT /api/ri-comm-life/1000003/detail/1000044 {"contract":"2","year":"2025","comm":"1"}',
      'POST /api/ri-comm-life/unggah/pratinjau {"csv":"a;b"}',
      'POST /api/ri-comm-life/unggah {"csv":"a;b"}',
    ])
  })
})
