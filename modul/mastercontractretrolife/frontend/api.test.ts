// Permintaan HTTP modul - jalur, metode, dan badan SAMA dengan rute `backend/handlers`.

import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  PREFIX_MCRL,
  ambilRate,
  cariMasterReinsurer,
  hapus,
  rencanaSimpan,
  salinSemua,
  simpanKontrak,
  simpanTahun,
} from './api'

/** Satu panggilan `fetch` yang direkam. */
interface Rekam {
  url: string
  metode: string
  badan: unknown
}

function rekamFetch(jawab: unknown = {}): Rekam[] {
  const rekam: Rekam[] = []
  vi.stubGlobal(
    'fetch',
    vi.fn(async (url: string, init: RequestInit) => {
      rekam.push({ url, metode: init.method ?? 'GET', badan: init.body === undefined ? undefined : JSON.parse(String(init.body)) })
      return new Response(JSON.stringify(jawab), { status: 200 })
    }),
  )
  return rekam
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('rencana simpan: baru = POST ke induk, ubah = PUT ke /{id} (ADR-0006)', () => {
  it('tahun (akar)', () => {
    expect(rencanaSimpan('tahun', null, '')).toEqual({ metode: 'POST', jalur: `${PREFIX_MCRL}/tahun` })
    expect(rencanaSimpan('tahun', null, 'UJI-T1')).toEqual({ metode: 'PUT', jalur: `${PREFIX_MCRL}/tahun/UJI-T1` })
  })

  it('anak berinduk', () => {
    expect(rencanaSimpan('reinsurer', { segmen: 'kontrak', id: 'UJI-K1' }, '')).toEqual({
      metode: 'POST',
      jalur: `${PREFIX_MCRL}/kontrak/UJI-K1/reinsurer`,
    })
    expect(rencanaSimpan('reinsurer', { segmen: 'kontrak', id: 'UJI-K1' }, 'UJI-R1')).toEqual({
      metode: 'PUT',
      jalur: `${PREFIX_MCRL}/reinsurer/UJI-R1`,
    })
  })

  it('id disandi di jalur', () => {
    expect(rencanaSimpan('kontrak', { segmen: 'tahun', id: 'UJI T/1' }, '').jalur).toBe(`${PREFIX_MCRL}/tahun/UJI%20T%2F1/kontrak`)
  })
})

describe('permintaan', () => {
  it('simpan tahun baru: POST tanpa id di jalur', async () => {
    const r = rekamFetch()
    await simpanTahun({ id: '', treatyYear: '2026', underwritingYear: '2026', startDate: '2026-01-01', endDate: '2026-12-31' })
    expect(r).toEqual([
      { url: `${PREFIX_MCRL}/tahun`, metode: 'POST', badan: { id: '', treatyYear: '2026', underwritingYear: '2026', startDate: '2026-01-01', endDate: '2026-12-31' } },
    ])
  })

  it('ubah kontrak: PUT /kontrak/{id}', async () => {
    const r = rekamFetch()
    await simpanKontrak('UJI-T1', { id: 'UJI-K1', reinsTypeId: '10196', idr: '1', usd: '', bIdr: '0', bUsd: '0' })
    expect(r[0]?.url).toBe(`${PREFIX_MCRL}/kontrak/UJI-K1`)
    expect(r[0]?.metode).toBe('PUT')
  })

  it('hapus: DELETE berbadan dampak yang dilihat di popup', async () => {
    const r = rekamFetch({ pesan: 'Data Berhasil di Hapus', terhapus: { security: 0, reinsurer: 0, business: 0 } })
    const h = await hapus('kontrak', 'UJI-K1', { security: 1, reinsurer: 2, business: 3 })
    expect(r).toEqual([{ url: `${PREFIX_MCRL}/kontrak/UJI-K1`, metode: 'DELETE', badan: { dampak: { security: 1, reinsurer: 2, business: 3 } } }])
    expect(h.pesan).toBe('Data Berhasil di Hapus')
  })

  it('Copy to all Reinstype: POST berbadan sasaran', async () => {
    const r = rekamFetch({ pesan: 'Copied to all reins types.', jumlah: 1, baru: [] })
    await salinSemua('UJI-B1', ['UJI-K2'])
    expect(r).toEqual([{ url: `${PREFIX_MCRL}/business/UJI-B1/salin-semua`, metode: 'POST', badan: { sasaran: ['UJI-K2'] } }])
  })

  it('autocomplete dan rate lewat kueri', async () => {
    const r = rekamFetch({ daftar: [], total: 0 })
    await cariMasterReinsurer('UJI RE')
    await ambilRate('UJI-RT')
    expect(r.map((x) => x.url)).toEqual([`${PREFIX_MCRL}/master-reinsurer?cari=UJI+RE`, `${PREFIX_MCRL}/rate?idusedby=UJI-RT`])
  })
})
