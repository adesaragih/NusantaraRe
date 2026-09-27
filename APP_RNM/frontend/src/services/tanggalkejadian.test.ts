import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { ApiFailure, pesanGalat, ubahTanggalKejadian } from './api'

// Uji klien tanggal kejadian — tombol `Edit Date` layar Detail.
//
// ⚠️ Rutenya ada di backend sejak tiket 06, tetapi nol pemanggil di React.
// Rute tanpa pemanggil tidak berbunyi di uji mana pun: backend hijau, layar
// hijau, dan jalurnya tetap tidak pernah dapat dijalankan orang.

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

  it('tanggal kejadian milik PESERTA, jadi pengenal peserta wajib ada di jalur', () => {
    // ⛔ `ValidasiDOL_Act` berkelas Int-LIFE_PREMIUM_DETAIL dan menempelkan
    // galatnya pada `.DATE_OF_LOSS` PESERTA. Jalur yang hanya menyebut klaim
    // akan menyetel tanggal untuk semua peserta sekaligus.
    expect(ubahTanggalKejadian).toHaveLength(3)
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
