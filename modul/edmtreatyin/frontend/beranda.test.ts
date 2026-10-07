// Daftar kotak masuk Beranda EDM Treaty In (pola `modul/nbtreatyin/frontend/beranda.test.ts`). Fixture UJI-.

import { describe, expect, it } from 'vitest'

import type { RingkasanKasus } from './api'
import { keDaftarBeranda, sejak } from './beranda'
import { KOLOM_BERANDA } from './labels'
import { PENDAFTARAN_MENU } from './menu'

const baris: RingkasanKasus = {
  id: 'UJI-EDMT-1',
  noOffer: 'UJI-OFFER',
  noPolis: 'UJI-POLIS',
  edmNo: 'UJI-POLIS/E01',
  sobName: 'UJI SOB',
  cedingCoName: 'UJI CEDING',
  edmType: '3',
  proportionalType: 'Proportional',
  marketingName: 'UJI MO',
  nbStatus: 'UJI STATUS',
  statusWork: 'Acceptance by Sec Treaty',
  positionNote: 'ReasTreatyInSecHead',
  tglCreate: '2026-10-03 09:00:00',
  startDate: '2026-10-01 00:00:00',
  tglUpdate: '2026-09-11 08:00:00',
}

describe('daftar kotak masuk Beranda EDM Treaty In', () => {
  it('kolom berurutan: sepuluh kolom portal + Tanggal Create, Time Since Last Update', () => {
    const d = keDaftarBeranda([baris], new Date(2026, 9, 6, 8, 0, 0))
    expect(d.kolom.map((k) => k.label)).toEqual([
      'EDM Number',
      'Offer No',
      'Policy Number',
      'EDM Number',
      'SOB',
      'Ceding',
      'EDM Type',
      'Proportional Type',
      'Marketing',
      'Status',
      'Tanggal Create',
      'Time Since Last Update',
    ])
    expect(Object.values(KOLOM_BERANDA)).toEqual(d.kolom.map((k) => k.label))
  })

  it('sel: nilai DB apa adanya, EDM Type berteks DT TreatyEDMListType, tanggal dd-mm-yyyy, lama sejak update', () => {
    const d = keDaftarBeranda([baris], new Date(2026, 9, 6, 8, 0, 0))
    expect(d.baris).toHaveLength(1)
    expect(d.baris[0]?.id).toBe('UJI-EDMT-1')
    expect(d.kolom.map((k) => d.baris[0]?.sel[k.kunci])).toEqual([
      'UJI-EDMT-1',
      'UJI-OFFER',
      'UJI-POLIS',
      'UJI-POLIS/E01',
      'UJI SOB',
      'UJI CEDING',
      'Adjustment Premium',
      'Proportional',
      'UJI MO',
      'UJI STATUS',
      '03-10-2026',
      '25d ago',
    ])
  })

  it('lama sejak update: menit, jam, hari, bulan, tahun + bulan', () => {
    const kini = new Date(2026, 9, 6, 12, 0, 0)
    expect(sejak('2026-10-06 11:59:40', kini)).toBe('just now')
    expect(sejak('2026-10-06 11:15:00', kini)).toBe('45m ago')
    expect(sejak('2026-10-06 07:00:00', kini)).toBe('5h ago')
    expect(sejak('2026-09-11 12:00:00', kini)).toBe('25d ago')
    expect(sejak('2025-11-06 12:00:00', kini)).toBe('11mo ago')
    expect(sejak('2025-06-06 12:00:00', kini)).toBe('1y 4mo ago')
    expect(sejak('', kini)).toBe('')
  })

  it('menu: satu halaman portal, kelompok = folder korpus VERBATIM, penyedia Beranda terdaftar', () => {
    expect(PENDAFTARAN_MENU.nama).toBe('edmtreatyin')
    expect(PENDAFTARAN_MENU.kelompok).toBe('EDM Treaty In')
    expect(PENDAFTARAN_MENU.halaman).toEqual(['edmtreatyin-portal'])
    expect(PENDAFTARAN_MENU.halamanAwal).toBe('edmtreatyin-portal')
    // keputusan work owner 07-10-2026: EDM ikut kotak masuk Beranda
    expect(typeof PENDAFTARAN_MENU.antreanBeranda).toBe('function')
    expect(typeof PENDAFTARAN_MENU.daftarBeranda).toBe('function')
  })
})
