// Daftar kotak masuk Beranda NB Treaty In (keputusan work owner 06-10-2026: klik workbasket / jenis di Beranda
// menampilkan berkasnya tanpa masuk menu NB Treaty In; kolom = portal + kolom yang ada datanya). Fixture UJI-.

import { describe, expect, it } from 'vitest'

import type { RingkasanKasus } from './api'
import { keDaftarBeranda, sejak } from './beranda'
import { KOLOM_BERANDA } from './labels'
import { PENDAFTARAN_MENU } from './menu'

const baris: RingkasanKasus = {
  id: 'UJI-NB-1',
  proportionalType: 'Proportional',
  businessName: 'UJI BISNIS',
  insuredName: 'UJI TERTANGGUNG',
  marketingName: 'UJI MO',
  nbStatus: 'UJI STATUS',
  statusWork: 'Input Realitation',
  positionNote: 'ReasTreatyInAdmin',
  noPolis: '',
  tglCreate: '2026-10-03 09:00:00',
  createOpName: 'UJI PEMBUAT',
  cedingCoName: 'UJI CEDING',
  startDate: '2026-10-01 00:00:00',
  tglUpdate: '2026-09-11 08:00:00',
  productionDate: '',
}

describe('daftar kotak masuk Beranda NB Treaty In', () => {
  it('kolom berurutan: portal + Ceding Company, Inception Date, Time Since Last Update', () => {
    const d = keDaftarBeranda([baris], new Date(2026, 9, 6, 8, 0, 0))
    expect(d.kolom.map((k) => k.label)).toEqual([
      'Offer No',
      'Type',
      'Insured Name',
      'Group Business',
      'Ceding Company',
      'Inception Date',
      'Marketing',
      'Status Inbox',
      'User Create',
      'Tanggal Create',
      'Time Since Last Update',
    ])
    expect(Object.values(KOLOM_BERANDA)).toEqual(d.kolom.map((k) => k.label))
  })

  it('sel: nilai DB apa adanya, tanggal dd-mm-yyyy, lama sejak update', () => {
    const d = keDaftarBeranda([baris], new Date(2026, 9, 6, 8, 0, 0))
    expect(d.baris).toHaveLength(1)
    expect(d.baris[0]?.id).toBe('UJI-NB-1')
    expect(d.kolom.map((k) => d.baris[0]?.sel[k.kunci])).toEqual([
      'UJI-NB-1',
      'Proportional',
      'UJI TERTANGGUNG',
      'UJI BISNIS',
      'UJI CEDING',
      '01-10-2026',
      'UJI MO',
      'UJI STATUS',
      'UJI PEMBUAT',
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
    expect(sejak('2024-10-06 12:00:00', kini)).toBe('2y ago')
    expect(sejak('', kini)).toBe('')
  })

  it('terdaftar di PENDAFTARAN_MENU sebagai penyedia daftar Beranda', () => {
    expect(typeof PENDAFTARAN_MENU.daftarBeranda).toBe('function')
  })
})

// WO 07-10-2026 "TAMPILAN NYA HANYA NB-XXX AJA, BERLAKU NB DAN EDM TREATY": ID kasus salinan Copy Old = IDPEGA Pega
// utuh (`<kelas> <pyID>`); layar hanya menampilkan pyID, kunci buka / kirim tetap ID utuh.
describe('Beranda - berkas salinan Copy Old', () => {
  it('sel nomor tampil pyID; kunci baris (pembuka berkas) tetap ID utuh', () => {
    const d = keDaftarBeranda([{ ...baris, id: 'ASM-FW-GISFW-WORK UJI-NB-9' }], new Date(2026, 9, 6, 8, 0, 0))
    expect(d.baris[0]?.id).toBe('ASM-FW-GISFW-WORK UJI-NB-9')
    expect(d.baris[0]?.sel.id).toBe('UJI-NB-9')
  })
})
