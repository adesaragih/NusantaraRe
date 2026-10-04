import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { bolehUbahTanggalKlaim, ubahTanggalKlaim, type Klaim } from './api'

// Uji klien tiga tanggal klaim — tombol `Save` `EditDateClaimLife_Section`
// b1910 → `UpdateDateClaimLife_Act` (sensus §3.1, 28-09-2026).

let permintaan: { url: string; init: RequestInit } | null = null

beforeEach(() => {
  permintaan = null
  vi.stubGlobal(
    'fetch',
    vi.fn((url: string, init: RequestInit) => {
      permintaan = { url, init }
      return Promise.resolve(new Response(null, { status: 204 }))
    }),
  )
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('ubahTanggalKlaim', () => {
  it('memakai PUT di sarang peserta, satu permintaan untuk ketiga tanggal', async () => {
    await ubahTanggalKlaim('K/1', 'P 9', {
      tanggalTerimaKlaim: '2026-03-02',
      tanggalDokumenLengkap: '',
      tanggalKonfirmasi: '2026-03-04',
    })
    expect(permintaan?.init.method).toBe('PUT')
    expect(permintaan?.url).toContain('/api/klaim-life/K%2F1/peserta/P%209/tanggal-klaim')
    // Kosong dikirim kosong — backend menulisnya NULL, seperti isian Pega
    // yang dikosongkan.
    expect(JSON.parse(String(permintaan?.init.body))).toEqual({
      tanggalTerimaKlaim: '2026-03-02',
      tanggalDokumenLengkap: '',
      tanggalKonfirmasi: '2026-03-04',
    })
  })
})

function klaim(tahap: string, statusWork = '', kodeBaris: string[] = []): Klaim {
  return {
    tahap,
    statusWork,
    peserta: [{ baris: kodeBaris.map((kodeStatus) => ({ kodeStatus })) }],
  } as unknown as Klaim
}

describe('bolehUbahTanggalKlaim', () => {
  it('hanya Outstanding Claim — pemegang Admin DAN grid peserta terbuka', () => {
    // b1000/b1313/b1550/b1788 `pyPosition!='ReasLifeAdmin'` → baca-saja.
    expect(bolehUbahTanggalKlaim(klaim('Outstanding Claim'))).toBe(true)
    expect(bolehUbahTanggalKlaim(klaim('Input Register'))).toBe(false)
    expect(bolehUbahTanggalKlaim(klaim('Medical Check'))).toBe(false)
    expect(bolehUbahTanggalKlaim(klaim('Claim Analis'))).toBe(false)
  })

  it('kasus tertutup tidak dapat diubah', () => {
    expect(bolehUbahTanggalKlaim(klaim('Outstanding Claim', 'Resolved-Completed'))).toBe(false)
    expect(bolehUbahTanggalKlaim(null)).toBe(false)
  })

  it('terkunci sesudah Save to RNM pertama — ada baris berstatus (OQ-M1)', () => {
    // b1000 `CLAIM_NO!=''`: padanannya baris adjustment yang STS_REJECT-nya terisi.
    expect(bolehUbahTanggalKlaim(klaim('Outstanding Claim', '', ['', '']))).toBe(true)
    expect(bolehUbahTanggalKlaim(klaim('Outstanding Claim', '', ['', '0']))).toBe(false)
    expect(bolehUbahTanggalKlaim(klaim('Outstanding Claim', '', ['2']))).toBe(false)
  })
})
