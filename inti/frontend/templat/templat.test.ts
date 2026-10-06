import { afterEach, describe, expect, it, vi } from 'vitest'

import { TEMPLATE_MANAGER as TM } from '../labels'
import { aktifkan, ambilDaftar, ambilRiwayat, periksa, unggah, type Slot } from './api'
import { kelompokSlot, saringSlot, ukuranBerkas, versiBerikut } from './aturan'

const slot = (kode: string, grup: string, nama: string, berkas?: string): Slot => ({
  kode,
  menu: grup.toLowerCase(),
  grup,
  nama,
  dipakaiDi: '',
  ekstensi: '.csv',
  pemisah: ';',
  jumlahKolom: 3,
  namaUnduhan: `${grup} - ${nama}.csv`,
  aktif:
    berkas === undefined
      ? null
      : { versi: 1, namaBerkas: berkas, ukuran: 10, jumlahKolom: 3, catatan: '', aktif: true, diunggahOleh: 'UJI-IT', tglUnggah: '04-10-2026 09:30' },
})

const DAFTAR = [
  slot('uji.premi.fire', 'Bordereaux', 'PREMIUM · FIRE', 'uji_fire_rev.csv'),
  slot('uji.agg', 'Aggregate', 'Aggregate'),
  slot('uji.claim.fire', 'Bordereaux', 'CLAIM · FIRE'),
]

describe('aturan Template Manager', () => {
  it('dikelompokkan per menu, urutan kemunculan', () => {
    expect(kelompokSlot(DAFTAR).map((g) => [g.grup, g.slot.length])).toEqual([
      ['Bordereaux', 2],
      ['Aggregate', 1],
    ])
  })

  it('saringan teks dan menu', () => {
    expect(saringSlot(DAFTAR, 'claim', '').map((s) => s.kode)).toEqual(['uji.claim.fire'])
    expect(saringSlot(DAFTAR, 'rev', '').map((s) => s.kode)).toEqual(['uji.premi.fire'])
    expect(saringSlot(DAFTAR, '', 'Aggregate').map((s) => s.kode)).toEqual(['uji.agg'])
    expect(saringSlot(DAFTAR, 'fire', 'Aggregate')).toEqual([])
  })

  it('ukuran, versi berikut, syarat', () => {
    expect([ukuranBerkas(832), ukuranBerkas(13830)]).toEqual(['832 B', '13.5 KB'])
    expect([versiBerikut([]), versiBerikut([{ versi: 2 }, { versi: 1 }])]).toEqual([1, 3])
    expect(TM.syaratTeks('.csv', ';', 51)).toBe('.csv · separator ; · 51 columns')
    expect(TM.syaratTeks('.xlsx', '', 0)).toBe('.xlsx')
  })
})

describe('klien Template Manager', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('rute dan metode sama dengan inti/backend/templat/rute', async () => {
    const panggil: string[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push(`${String(init.method ?? 'GET')} ${url.replace(/^https?:\/\/[^/]+/, '')} ${init.body instanceof FormData ? 'form' : ''}`.trim())
        return Promise.resolve(new Response('{}', { status: 200, headers: { 'Content-Type': 'application/json' } }))
      }),
    )
    const f = new File(['A;B;C\n'], 'uji.csv')
    await ambilDaftar()
    await ambilRiwayat('uji.premi.fire')
    await periksa('uji.premi.fire', f)
    await unggah('uji.premi.fire', f, 'catatan')
    await aktifkan('uji.premi.fire', 0)
    expect(panggil).toEqual([
      'GET /api/templat',
      'GET /api/templat/uji.premi.fire/riwayat',
      'POST /api/templat/uji.premi.fire/periksa form',
      'POST /api/templat/uji.premi.fire/unggah form',
      'POST /api/templat/uji.premi.fire/aktifkan',
    ])
  })
})
