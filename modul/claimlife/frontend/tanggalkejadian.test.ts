import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { ApiFailure, pesanGalat } from '../../../inti/frontend/klien'
import { ubahTanggalKejadian } from './api'

// Uji klien tanggal kejadian — tombol `Edit Date` layar Detail.
//
// ⚠️ Berkas ini lahir bersama pemanggil pertamanya. Rute backendnya ada
// sejak tiket 06, tetapi sampai kelompok Detail & Tutup nol pemanggil di
// React — dan rute tanpa pemanggil tidak berbunyi di uji mana pun: backend
// hijau, layar hijau, jalurnya tetap tidak dapat dijalankan orang.

let permintaan: { url: string; init: RequestInit } | null = null

beforeEach(() => {
  permintaan = null
  vi.stubGlobal(
    'fetch',
    vi.fn((url: string, init: RequestInit) => {
      permintaan = { url, init }
      // 204 WAJIB berbadan null menurut spesifikasi Fetch; Response
      // menolak 204 yang berbadan teks, sekalipun teks kosong.
      return Promise.resolve(new Response(null, { status: 204 }))
    }),
  )
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('ubahTanggalKejadian', () => {
  it('memakai PUT dan jalur peserta, bukan jalur klaim', async () => {
    await ubahTanggalKejadian('K-1', 'P-9', '2026-03-01')
    expect(permintaan?.init.method).toBe('PUT')
    expect(permintaan?.url).toContain('/api/klaim-life/K-1/peserta/P-9/tanggal-kejadian')
  })

  it('pengenal berisi karakter jalur di-encode', async () => {
    await ubahTanggalKejadian('K/1', 'P 9', '2026-03-01')
    expect(permintaan?.url).toContain('K%2F1')
    expect(permintaan?.url).toContain('P%209')
  })

  it('badannya membawa tanggalKejadian apa adanya', async () => {
    await ubahTanggalKejadian('K-1', 'P-9', '2026-03-01')
    expect(JSON.parse(String(permintaan?.init.body))).toEqual({
      tanggalKejadian: '2026-03-01',
    })
  })

  it('204 bukan kegagalan', async () => {
    await expect(ubahTanggalKejadian('K-1', 'P-9', '2026-03-01')).resolves.toBeUndefined()
  })

  it('422 di luar jendela valuasi meneruskan kalimat backend', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(() =>
        Promise.resolve(
          new Response(JSON.stringify({ galat: 'Invalid DOL', pesertaId: 'P-9' }), {
            status: 422,
            headers: { 'Content-Type': 'application/json' },
          }),
        ),
      ),
    )
    const e = await ubahTanggalKejadian('K-1', 'P-9', '2020-01-01').catch(
      (x: unknown) => x,
    )
    expect(e).toBeInstanceOf(ApiFailure)
    expect((e as ApiFailure).status).toBe(422)
    // Kalimat `ValidasiDOL_Act` b410 sampai ke pemakai, bukan teks bawaan.
    expect(pesanGalat(e)).toBe('Invalid DOL')
  })
})
