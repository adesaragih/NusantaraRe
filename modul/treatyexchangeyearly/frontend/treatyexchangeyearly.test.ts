import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilDaftar, ambilPilihan, tambah, ubah, type Kurs } from './api'
import {
  gantiTahun,
  isianDari,
  keInputTanggal,
  opsiMataUang,
  periksaIsian,
  periodeBawaan,
  QUARTER,
  rapikanIsian,
} from './aturan'
import { MENU_TEY, TEY } from './labels'

const kurs: Kurs = {
  kunci: 'AAA1',
  id: '10010',
  treatyYear: '2019',
  idCurrency: '10001',
  currency: 'USD',
  currencyName: 'UJI DOLLAR',
  startDate: '20190801T00000.000 GMT',
  endDate: '20200630T000000.000 GMT',
  mulai: '01-08-2019',
  akhir: '30-06-2020',
  toIdr: '14500.00',
  toUsd: '0',
  quarter: '0',
  userId: '',
  dateIu: '',
  dateIn: '',
  diubah: '',
}
const sah = {
  kunci: '',
  treatyYear: '2026',
  idCurrency: '10001',
  startDate: '2026-07-01',
  endDate: '2027-06-30',
  toIdr: '16500.5',
  toUsd: '',
  quarter: '0',
}

describe('label Treaty Exchange Yearly', () => {
  it('nama menu = M_NAV_MENU.LABEL migrasi inti 919, golongan MASTER TREATY, langsung menyala', () => {
    const sql = readFileSync(
      `${__dirname}/../../../inti/backend/migrations/919_m_nav_menu_treatyexchangeyearly.sql`,
      'utf8',
    )
    expect(sql).toContain(
      `'treatyexchangeyearly', '${MENU_TEY.kelompok}', 'MASTER TREATY', 'treatyexchangeyearly', 6, '1'`,
    )
  })
})

describe('aturan Treaty Exchange Yearly', () => {
  it('tanggal Pega ke isian; juga bentuk warisan rusak dan ber-jam', () => {
    expect(keInputTanggal('20190801T00000.000 GMT')).toBe('2019-08-01')
    expect(keInputTanggal('20250701T075400.000 GMT')).toBe('2025-07-01')
    expect(keInputTanggal('lain')).toBe('')
  })

  it('isian: Add Quarter 0; Edit membawa kunci dan tanggal', () => {
    expect(isianDari(null).quarter).toBe('0')
    expect(isianDari(kurs)).toEqual({
      kunci: 'AAA1',
      treatyYear: '2019',
      idCurrency: '10001',
      startDate: '2019-08-01',
      endDate: '2020-06-30',
      toIdr: '14500.00',
      toUsd: '0',
      quarter: '0',
    })
    expect(QUARTER).toEqual(['0', '1', '2', '3', '4'])
  })

  it('periode bawaan 1 Juli - 30 Juni; ganti tahun Add mengikuti selama tanggal belum diubah', () => {
    expect(periodeBawaan('2026')).toEqual({ startDate: '2026-07-01', endDate: '2027-06-30' })
    expect(periodeBawaan('26')).toBeNull()
    const a = gantiTahun(isianDari(null), '2026', true)
    expect([a.startDate, a.endDate]).toEqual(['2026-07-01', '2027-06-30'])
    const b = gantiTahun(a, '2027', true)
    expect([b.startDate, b.endDate]).toEqual(['2027-07-01', '2028-06-30'])
    const ubahan = gantiTahun({ ...a, endDate: '2027-03-31' }, '2027', true)
    expect(ubahan.endDate).toBe('2027-03-31')
    expect(gantiTahun(isianDari(kurs), '2020', false).startDate).toBe('2019-08-01')
  })

  it('periksa: tahun 4 angka, currency, tanggal berurutan, kurs angka bertitik', () => {
    expect(periksaIsian(sah)).toBeNull()
    expect(periksaIsian({ ...sah, treatyYear: '26' })).toBe(TEY.galatTahun)
    expect(periksaIsian({ ...sah, idCurrency: '' })).toBe(TEY.galatMataUang)
    expect(periksaIsian({ ...sah, endDate: '' })).toBe(TEY.galatTanggal)
    expect(periksaIsian({ ...sah, endDate: '2026-06-30' })).toBe(TEY.galatUrutan)
    expect(periksaIsian({ ...sah, toIdr: '16.500,50' })).toBe(TEY.galatIdr)
    expect(periksaIsian({ ...sah, toUsd: 'x' })).toBe(TEY.galatUsd)
    expect(rapikanIsian({ ...sah, toIdr: ' 1 ', quarter: ' ' })).toMatchObject({ toIdr: '1', quarter: '0' })
  })

  it('opsi Currency; mata uang lama yang hilang dari view tetap terbaca', () => {
    const m = [{ id: '10001', kode: 'USD', nama: 'UJI DOLLAR' }]
    expect(opsiMataUang(m, '10001')).toEqual([{ value: '10001', label: 'USD - UJI DOLLAR' }])
    expect(opsiMataUang(m, '10099', 'XYZ')).toContainEqual({ value: '10099', label: TEY.nilaiLama('XYZ') })
    expect(TEY.labelQuarter('0')).toBe('0 - Yearly')
  })
})

describe('klien Treaty Exchange Yearly', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan handlers/rute.go; Add tanpa kunci, Edit berkunci; tanpa hapus', async () => {
    const panggil: string[] = []
    const badan: unknown[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')}`)
        if (typeof init.body === 'string') badan.push(JSON.parse(init.body))
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar(' usd ', '2025')
    await ambilDaftar('', '')
    await ambilPilihan()
    await tambah({ ...sah, kunci: 'X' })
    await ubah({ ...sah, kunci: 'AAA1' })
    expect(panggil).toEqual([
      'GET /api/treaty-exchange-yearly?q=usd&tahun=2025',
      'GET /api/treaty-exchange-yearly',
      'GET /api/treaty-exchange-yearly/pilihan',
      'POST /api/treaty-exchange-yearly',
      'PUT /api/treaty-exchange-yearly',
    ])
    expect(badan).toMatchObject([{ kunci: '' }, { kunci: 'AAA1' }])
  })
})

describe('tanpa popup View (perintah work owner 05-10-2026: "yang view nya popup ... ga usah di buat, toh bisa liat langsung dari list tabelnya")', () => {
  it('tidak ada mode lihat, tombol View, atau komponen Lihat*', () => {
    const berkas = (dir: string): string[] =>
      readdirSync(dir).flatMap((n) => {
        const p = join(dir, n)
        if (statSync(p).isDirectory()) return berkas(p)
        return /\.tsx?$/.test(n) && !/\.test\.tsx?$/.test(n) ? [p] : []
      })
    for (const p of berkas(__dirname)) {
      const isi = readFileSync(p, 'utf8')
      expect(isi, p).not.toMatch(/jenis: 'lihat'|[A-Z]+\.view\b|judulLihat|function Lihat|from '\.\/Lihat/)
      expect(p, p).not.toMatch(/Lihat[A-Z][a-z]*\.tsx$/)
    }
  })
})
