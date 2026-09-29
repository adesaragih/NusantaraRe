import { beforeEach, describe, expect, it, vi } from 'vitest'

import { bolehCabutPeserta, cabutPeserta, type Klaim } from './api'

// Cabut peserta — OQ-M6 (GILIRAN-17): tombol `DELETE` `InputOSClaimLife`
// b17865, tampil bila `CLAIM_NO == ''` (b18082), tanpa konfirmasi (b18021).

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

function klaim(tahap: string, statusWork = '', kodeBaris: string[] = []): Klaim {
  return {
    tahap,
    statusWork,
    peserta: [{ baris: kodeBaris.map((kodeStatus) => ({ kodeStatus })) }],
  } as unknown as Klaim
}

describe('cabutPeserta', () => {
  it('POST ke sarang peserta — penandaan, bukan DELETE', async () => {
    await cabutPeserta('K/1', 'P 9')
    expect(permintaan?.init.method).toBe('POST')
    expect(permintaan?.url).toContain('/api/klaim-life/K%2F1/peserta/P%209/cabut')
  })
})

describe('bolehCabutPeserta', () => {
  it('hanya Outstanding, sebelum Save to RNM, kasus terbuka', () => {
    expect(bolehCabutPeserta(klaim('Outstanding Claim'))).toBe(true)
    expect(bolehCabutPeserta(klaim('Outstanding Claim', '', ['0']))).toBe(false)
    expect(bolehCabutPeserta(klaim('Input Register'))).toBe(false)
    expect(bolehCabutPeserta(klaim('Outstanding Claim', 'Resolved-Completed'))).toBe(false)
    expect(bolehCabutPeserta(null)).toBe(false)
  })
})
