import { readFileSync } from 'node:fs'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilDaftar, ambilRincian, ambilRingkasan, cariMasterTreaty, hapus, pratinjau, simpan, type IrisanRingkasan, type Kelompok } from './api'
import {
  barisPesan,
  barisTotal,
  formatAngka,
  formatBulat,
  formatPersen,
  jumlahHalaman,
  KEPALA_TEMPLATE,
  KOLOM_GRID,
  kunciDari,
  tampilanRingkasan,
} from './aturan'

const kelompok: Kelompok = {
  tanggalInput: '04-10-2026',
  cedingCode: 'UJI-C1',
  cedingName: 'UJI CEDING SATU',
  treatyType: 'QS',
  asAt: '31-12-2024',
  uwYear: '2024',
  jumlahBaris: 3,
  inputTerakhir: '04-10-2026 09:30',
}

const IRISAN: IrisanRingkasan[] = [
  { cedingCode: 'UJI-C1', cedingName: 'UJI CEDING SATU', treatyType: 'QS', coverage: 'EQVET', rnmValueInUsd: '100.10' },
  { cedingCode: 'UJI-C1', cedingName: 'UJI CEDING SATU', treatyType: 'QS', coverage: 'FLOOD', rnmValueInUsd: '50.05' },
  { cedingCode: 'UJI-C1', cedingName: 'UJI CEDING SATU', treatyType: 'OR', coverage: 'EQVET', rnmValueInUsd: '49.85' },
  { cedingCode: 'UJI-C2', cedingName: 'UJI CEDING DUA', treatyType: 'SURPLUS', coverage: 'EQVET', rnmValueInUsd: '300' },
]

describe('aturan Aggregate', () => {
  it('KOLOM_GRID = models.KolomGrid backend: nama, urutan, angka', () => {
    const go = readFileSync(`${__dirname}/../backend/models/aggregate.go`, 'utf8')
    const blok = /var KolomGrid = \[\]Kolom\{([\s\S]*?)\n\}/.exec(go)?.[1] ?? ''
    const dariGo = [...blok.matchAll(/\{Nama: "([A-Z_]+)", Jenis: (\w+)[^}]*\}/g)].map((m) => ({
      nama: m[1]!,
      angka: m[2] === 'Angka',
    }))
    expect(dariGo).toHaveLength(40)
    expect(KOLOM_GRID).toEqual(dariGo)
  })

  it('kepala template = 38 kolom CSV, pemisah titik koma', () => {
    const kolom = KEPALA_TEMPLATE.split(';')
    expect(kolom).toHaveLength(38)
    expect([kolom[0], kolom[8], kolom[37]]).toEqual(['ASSESMENT ZONE', 'TO USD', 'REMARK'])
  })

  it('formatAngka: format Indonesia, 2 angka di belakang koma, bulat setengah ke atas', () => {
    expect(formatAngka('1234567.5')).toBe('1.234.567,50')
    expect(formatAngka('-1000')).toBe('-1.000,00')
    expect(formatAngka('2166.528600000000140824359')).toBe('2.166,53')
    expect(formatAngka('35747721.9')).toBe('35.747.721,90')
    expect(formatAngka('37535107.995')).toBe('37.535.108,00')
    expect(formatAngka('999.995')).toBe('1.000,00')
    expect(formatAngka('0.004')).toBe('0,00')
    expect(formatAngka('-0.004')).toBe('0,00')
    expect(formatAngka('0.00006060606060606061')).toBe('0,00')
    expect(formatAngka('.5')).toBe('0,50')
    expect(formatAngka('10')).toBe('10,00')
    expect(formatAngka('')).toBe('')
    expect(formatAngka(undefined)).toBe('')
    expect(formatAngka('abc')).toBe('abc')
  })

  it('formatBulat dan formatPersen: format Indonesia', () => {
    expect(formatBulat(13450)).toBe('13.450')
    expect(formatBulat(3)).toBe('3')
    expect(formatPersen(75)).toBe('75,00%')
    expect(formatPersen(100 / 3)).toBe('33,33%')
  })

  it('baris Total, kunci daftar, halaman, pesan Save', () => {
    expect(barisTotal({ ASSESMENT_ZONE: 'Total :' })).toBe(true)
    expect(barisTotal({ ASSESMENT_ZONE: '1.1' })).toBe(false)
    expect(kunciDari(kelompok)).toEqual({
      tanggalInput: '04-10-2026', cedingCode: 'UJI-C1', cedingName: 'UJI CEDING SATU', treatyType: 'QS', asAt: '31-12-2024', uwYear: '2024',
    })
    expect([jumlahHalaman(0, 50), jumlahHalaman(50, 50), jumlahHalaman(51, 50)]).toEqual([1, 1, 2])
    expect(barisPesan('Assesment Zone in list 1 Not Found\n\nTo USD in list 2 cannot be empty\n')).toEqual([
      'Assesment Zone in list 1 Not Found',
      'To USD in list 2 cannot be empty',
    ])
  })

  it('chart tingkat Ceding: batang per ceding terbesar dulu, ditumpuk per Treaty Type, jumlah eksak', () => {
    const v = tampilanRingkasan(IRISAN, {})
    expect([v.tingkat, v.total, v.bisaBuka, v.seri]).toEqual(['ceding', '500.00', true, ['OR', 'QS', 'SURPLUS']])
    expect(v.batang.map((b) => [b.kunci, b.label, b.nilai, b.persen])).toEqual([
      ['UJI-C2', 'UJI CEDING DUA', '300', 60],
      ['UJI-C1', 'UJI CEDING SATU', '200.00', 40],
    ])
    expect(v.batang[0]!.lebar).toBe(100)
    expect(v.batang[1]!.lebar).toBeCloseTo(66.667, 2)
    expect(v.batang[1]!.bagian.map((s) => [s.kunci, s.nilai, s.seri])).toEqual([
      ['OR', '49.85', 0],
      ['QS', '150.15', 1],
    ])
    expect(v.batang[1]!.bagian[1]!.lebar).toBeCloseTo(50.05, 2)
    expect(v.batang[0]!.bagian.map((s) => s.seri)).toEqual([2])
  })

  it('chart dibuka: ceding -> Treaty Type ditumpuk per Coverage -> Coverage tanpa tumpukan', () => {
    const t = tampilanRingkasan(IRISAN, { ceding: 'UJI-C1' })
    expect([t.tingkat, t.total, t.namaCeding, t.bisaBuka, t.seri]).toEqual(['treatyType', '200.00', 'UJI CEDING SATU', true, ['EQVET', 'FLOOD']])
    expect(t.batang.map((b) => [b.kunci, b.nilai, b.bagian.map((s) => `${s.kunci}=${s.nilai}`)])).toEqual([
      ['QS', '150.15', ['EQVET=100.10', 'FLOOD=50.05']],
      ['OR', '49.85', ['EQVET=49.85']],
    ])
    const c = tampilanRingkasan(IRISAN, { ceding: 'UJI-C1', treatyType: 'QS' })
    expect([c.tingkat, c.total, c.bisaBuka, c.seri]).toEqual(['coverage', '150.15', false, []])
    expect(c.batang.map((b) => [b.kunci, b.nilai, b.bagian.length])).toEqual([
      ['EQVET', '100.10', 1],
      ['FLOOD', '50.05', 1],
    ])
    expect(tampilanRingkasan([], {}).batang).toEqual([])
  })
})

describe('klien Aggregate', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan handlers/rute.go', async () => {
    const panggil: [string, string, string][] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push([String(init.method), url, typeof init.body === 'string' ? init.body : ''])
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    await ambilDaftar('uji', 2)
    await ambilRingkasan()
    await ambilRincian(kunciDari(kelompok))
    await hapus(kunciDari(kelompok))
    await cariMasterTreaty('uji')
    await pratinjau('a;b', ['UJI-V1'])
    await simpan([{ ASSESMENT_ZONE: '1.1' }])
    expect(panggil.map(([m, u]) => `${m} ${u.replace(/^https?:\/\/[^/]+/, '')}`)).toEqual([
      'GET /api/aggregate?q=uji&halaman=2',
      'GET /api/aggregate/ringkasan',
      'GET /api/aggregate/rincian?tanggalInput=04-10-2026&cedingCode=UJI-C1&cedingName=UJI+CEDING+SATU&treatyType=QS&asAt=31-12-2024&uwYear=2024',
      'POST /api/aggregate/hapus',
      'GET /api/aggregate/master-treaty?q=uji',
      'POST /api/aggregate/pratinjau',
      'POST /api/aggregate/simpan',
    ])
    expect(JSON.parse(panggil[5]![2])).toEqual({ csv: 'a;b', masterTreaty: ['UJI-V1'] })
    expect(JSON.parse(panggil[6]![2])).toEqual({ baris: [{ ASSESMENT_ZONE: '1.1' }] })
  })
})
