import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambil, ambilDaftar, ambilDetail, hapus, pratinjauUnggah, simpanUnggah, tambah, tambahDetail, ubah, ubahDetail } from './api'
import {
  berkasSah,
  DIGIT_YEAR_MONTH,
  gantiUrut,
  isianDariRincian,
  ISIAN_RINCIAN_KOSONG,
  jumlahHalaman,
  KOLOM_URUT,
  MAKS_BYTES_CSV,
  nomorAwal,
  periksaIsianRincian,
  periksaNama,
  SARINGAN_AWAL,
  tampilDesimal,
  tandaUrut,
  UKURAN_HALAMAN,
  UKURAN_HALAMAN_RINCIAN,
} from './aturan'
import { MENU_RK, RK } from './labels'

const XML_LABEL = [
  // InboxSummaryRIRisk.xml
  'R/I RISK SUMMARY',
  'ID',
  'R/I RISK NAME',
  'MODIFY OPERATOR',
  'MODIFY DATE',
  'Upload CSV',
  'View Upload',
  'Simpan Upload',
  'Edit',
  'Detail',
  'Delete',
  'Format excel : USEDBY, CONTRACT, YEAR, MONTH, RISK',
  // InboxRIRisk.xml
  'R/I RISK DETAIL',
  'Clear Field',
  'Save',
  'Cancel',
  'CONTRACT',
  'YEAR',
  'MONTH',
  'RISK (PERMIL)',
  'EDIT',
  '0',
]

describe('label R/I Risk', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 941 (K4), golongan MASTER TREATY URUTAN 11, langsung menyala', () => {
    const sql = readFileSync(`${__dirname}/../../../inti/backend/migrations/941_m_nav_menu_ririsklife.sql`, 'utf8')
    expect(MENU_RK.kelompok).toBe('R/I Risk')
    expect(sql).toContain(`'ririsklife', '${MENU_RK.kelompok}', 'MASTER TREATY', 'ririsklife', 11, '1'`)
  })

  it('label layar VERBATIM dari kedua XML', () => {
    const nilai = Object.values(RK).filter((v) => typeof v === 'string')
    for (const l of XML_LABEL) expect(nilai, l).toContain(l)
  })
})

describe('aturan R/I Risk', () => {
  it('urut: kolom baru menaik, kolom sama balik arah; halaman kembali 1; bawaan backend (ID menaik)', () => {
    const a = gantiUrut({ ...SARINGAN_AWAL, halaman: 3 }, 'usedby')
    expect(a).toMatchObject({ urut: 'usedby', turun: false, halaman: 1 })
    expect(gantiUrut(a, 'usedby').turun).toBe(true)
    expect(tandaUrut(a, 'usedby')).toBe(' ▲')
    expect(tandaUrut(a, 'id')).toBe('')
    expect(KOLOM_URUT).toEqual(['id', 'usedby', 'operatorid', 'modifieddate'])
    expect(SARINGAN_AWAL.urut).toBe('')
  })

  it('ringkasan 50 baris (b10030), detail 200 baris (pyPageSizeOther b10169)', () => {
    expect(UKURAN_HALAMAN).toBe(50)
    expect(UKURAN_HALAMAN_RINCIAN).toBe(200)
    expect(jumlahHalaman(0)).toBe(1)
    expect(jumlahHalaman(51)).toBe(2)
    expect(jumlahHalaman(201, UKURAN_HALAMAN_RINCIAN)).toBe(2)
    expect(nomorAwal(2)).toBe(51)
  })

  it('nama ringkasan wajib; berkas .csv sampai 4 MB', () => {
    expect(periksaNama('  ')).toBe(RK.galatNama)
    expect(periksaNama('UJI RISK')).toBeNull()
    expect(berkasSah('uji.CSV', 10)).toBe(true)
    expect(berkasSah('uji.xlsx', 10)).toBe(false)
    expect(berkasSah('uji.csv', MAKS_BYTES_CSV + 1)).toBe(false)
  })

  it('detail: CONTRACT dan RISK wajib; YEAR / MONTH boleh kosong, bulat <= 4 angka; RISK tidak negatif; EDIT mengisi form berkoma desimal', () => {
    expect(DIGIT_YEAR_MONTH).toBe(4)
    expect(periksaIsianRincian(ISIAN_RINCIAN_KOSONG)).toBe(RK.galatContract)
    expect(periksaIsianRincian({ contract: '1', year: '', month: '', risk: ' ' })).toBe(RK.galatRisk)
    expect(periksaIsianRincian({ contract: '0', year: '', month: '', risk: '921,9' })).toBeNull()
    expect(periksaIsianRincian({ contract: '0', year: '', month: '180', risk: '580.894351210924' })).toBeNull()
    expect(periksaIsianRincian({ contract: '0', year: '12345', month: '', risk: '1' })).toBe(RK.galatAngka(RK.year, 4))
    expect(periksaIsianRincian({ contract: '0', year: '', month: '1,5', risk: '1' })).toBe(RK.galatAngka(RK.month, 4))
    expect(periksaIsianRincian({ contract: '1.5', year: '', month: '', risk: '1' })).toBe(RK.galatAngka(RK.contract, 10))
    expect(periksaIsianRincian({ contract: '0', year: '1', month: '', risk: '-1' })).toBe(RK.galatRiskNegatif)
    expect(
      isianDariRincian({ id: '131720', idUsedBy: '1000003', usedby: 'UJI', contract: ' 2', year: '1 ', month: '', risk: '921.9' }),
    ).toEqual({ contract: '2', year: '1', month: '', risk: '921,9' })
    expect(tampilDesimal('580.894351210924')).toBe('580,894351210924')
    expect(tampilDesimal('7')).toBe('7')
  })

  it('ringkasan: Save dan Cancel (saat Edit) DITAMPILKAN - keputusan WO 08-10-2026; Edit menyimpan lewat ubah(); label format tetap tidak', () => {
    const halaman = readFileSync(`${__dirname}/pages/RIRiskLife.tsx`, 'utf8')
    expect(halaman).toContain('useBolehUbah(NAMA_RK)')
    expect(halaman).toContain('{sibuk ? RK.menyimpan : RK.save}')
    expect(halaman).toMatch(/\{ubahID !== '' && \(\s*<button[^>]*onClick=\{kosongkanForm\}>\s*\{RK\.batal\}/)
    expect(halaman).toContain("ubahID === '' ? tambah(nama.trim()) : ubah(ubahID, nama.trim())")
    expect(halaman).not.toContain('{RK.format}')
    for (const label of ['RK.edit', 'RK.hapus']) {
      const i = halaman.indexOf(`{${label}}`)
      expect(halaman.lastIndexOf('{bolehUbah && (', i), label).toBeGreaterThan(halaman.lastIndexOf('</button>', i))
    }
    expect(halaman.match(/\{bolehUbah && \(/g)?.length).toBe(4)
    expect(halaman).toContain('bolehUbah={bolehUbah}')
    expect(readFileSync(`${__dirname}/components/UnggahCSV.tsx`, 'utf8')).not.toContain('{RK.format}')
  })

  it('detail View only: form dan EDIT hanya bila bolehUbah; nama R/I dari ringkasan, bukan isian; kolom MONTH', () => {
    const detail = readFileSync(`${__dirname}/components/RincianDetail.tsx`, 'utf8')
    expect(detail.match(/\{bolehUbah && \(/g)?.length).toBe(2)
    expect(detail).toContain('value={ringkasan.usedby} readOnly disabled')
    expect(detail).not.toMatch(/usedby:\s*isi/)
    expect(detail).toContain('{RK.month}')
    expect(detail).toContain('{RK.editRincian}')
  })

  it('pesan Delete menyebut jumlah rincian yang ikut terhapus', () => {
    expect(RK.tanyaHapus('UJI', 3)).toBe('Delete "UJI" and its 3 detail rows? This cannot be undone.')
    expect(RK.tanyaHapus('UJI', 1)).toContain('1 detail row?')
  })
})

describe('klien R/I Risk', () => {
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
    await tambahDetail('1000003', { contract: '1', year: '1', month: '', risk: '921,9' })
    await ubahDetail('1000003', '131720', { contract: '2', year: '', month: '12', risk: '1' })
    await pratinjauUnggah('a;b')
    await simpanUnggah('a;b')
    expect(panggil).toEqual([
      'GET /api/ri-risk-life',
      'GET /api/ri-risk-life?id=10&usedby=uji&urut=usedby&arah=desc&halaman=2',
      'GET /api/ri-risk-life/1000003',
      'POST /api/ri-risk-life {"usedby":"UJI"}',
      'PUT /api/ri-risk-life/1000003 {"usedby":"UJI 2"}',
      'DELETE /api/ri-risk-life/1000003',
      'GET /api/ri-risk-life/1000003/detail',
      'GET /api/ri-risk-life/1000003/detail?halaman=2',
      'POST /api/ri-risk-life/1000003/detail {"contract":"1","year":"1","month":"","risk":"921,9"}',
      'PUT /api/ri-risk-life/1000003/detail/131720 {"contract":"2","year":"","month":"12","risk":"1"}',
      'POST /api/ri-risk-life/unggah/pratinjau {"csv":"a;b"}',
      'POST /api/ri-risk-life/unggah {"csv":"a;b"}',
    ])
  })
})
