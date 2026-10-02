// Permintaan HTTP modul - jalur, metode, dan badan SAMA dengan rute `backend/handlers`.

import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  PREFIX_MPNL,
  ambilDaftarProduk,
  ambilLampiran,
  ambilProduk,
  ambilRate,
  cariMaster,
  cariPlan,
  hapusLampiran,
  lihatOffice,
  simpanProduk,
  ulangiLampiran,
} from './api'
import { produkBaru } from './bentuk'

/** Satu panggilan `fetch` yang direkam. */
interface Rekam {
  url: string
  metode: string
  badan: unknown
}

function rekamFetch(jawab: unknown = {}, status = 200): Rekam[] {
  const rekam: Rekam[] = []
  vi.stubGlobal(
    'fetch',
    vi.fn(async (url: string, init: RequestInit) => {
      rekam.push({
        url,
        metode: init.method ?? 'GET',
        badan: typeof init.body === 'string' ? JSON.parse(init.body) : init.body,
      })
      return new Response(JSON.stringify(jawab), { status })
    }),
  )
  return rekam
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('rute baca (paket 1)', () => {
  it('prefix sama dengan handlers.Prefix', () => {
    expect(PREFIX_MPNL).toBe('/api/master-product-name-life')
  })

  it('grid daftar dan satu produk; id disandi di jalur', async () => {
    const rekam = rekamFetch({ daftar: [], total: 0 })
    await ambilDaftarProduk()
    await ambilProduk('UJI 1/2')
    expect(rekam.map((r) => [r.metode, r.url])).toEqual([
      ['GET', `${PREFIX_MPNL}/produk`],
      ['GET', `${PREFIX_MPNL}/produk/UJI%201%2F2`],
    ])
  })
})

describe('simpan (paket 3-9)', () => {
  it('produk baru = POST tanpa id; ubah = PUT /{id}; penanda permintaan ikut badan', async () => {
    const rekam = rekamFetch({})
    await simpanProduk(produkBaru())
    await simpanProduk({ ...produkBaru(), id: '100044', hitungOutward: true })
    await simpanProduk({ ...produkBaru(), salinanDari: '100007' })
    expect(rekam.map((r) => [r.metode, r.url])).toEqual([
      ['POST', `${PREFIX_MPNL}/produk`],
      ['PUT', `${PREFIX_MPNL}/produk/100044`],
      ['POST', `${PREFIX_MPNL}/produk`],
    ])
    expect((rekam[0]?.badan as { id: string }).id).toBe('')
    expect((rekam[1]?.badan as { hitungOutward: boolean }).hitungOutward).toBe(true)
    expect((rekam[2]?.badan as { salinanDari: string }).salinanDari).toBe('100007')
  })

  it('galat server tampil dengan kalimatnya (amplop {galat})', async () => {
    rekamFetch({ galat: 'Product Name Empty; Ceding Empty' }, 422)
    await expect(simpanProduk(produkBaru())).rejects.toThrow()
  })
})

describe('pemilih dan lampiran', () => {
  it('jalur master, plan, rate, lampiran', async () => {
    const rekam = rekamFetch({ daftar: [], total: 0 })
    await cariMaster('pemegang-polis', 'uji')
    await cariPlan('UJI PLAN')
    await ambilRate('R1')
    await ambilLampiran('100044')
    await ulangiLampiran('100044', 'L1')
    await lihatOffice('100044', 'L1')
    await hapusLampiran('100044', 'L1')
    expect(rekam.map((r) => [r.metode, r.url])).toEqual([
      ['GET', `${PREFIX_MPNL}/master/pemegang-polis?cari=uji`],
      ['GET', `${PREFIX_MPNL}/master-plan?cari=UJI+PLAN`],
      ['GET', `${PREFIX_MPNL}/rate?riRateId=R1`],
      ['GET', `${PREFIX_MPNL}/produk/100044/lampiran`],
      ['POST', `${PREFIX_MPNL}/produk/100044/lampiran/L1/ulangi`],
      ['GET', `${PREFIX_MPNL}/produk/100044/lampiran/L1/office`],
      ['DELETE', `${PREFIX_MPNL}/produk/100044/lampiran/L1`],
    ])
  })
})
