import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { ApiFailure, pesanGalat, simpanKeRNM } from './api'

// Uji klien `Save to RNM` — `InputOSClaimLife.xml` b21102 → SaveOutStandingLife_Act.

let permintaan: { url: string; init: RequestInit } | null = null

function jawab(status: number, badan: unknown): void {
  vi.stubGlobal(
    'fetch',
    vi.fn((url: string, init: RequestInit) => {
      permintaan = { url, init }
      return Promise.resolve(
        new Response(JSON.stringify(badan), {
          status,
          headers: { 'Content-Type': 'application/json' },
        }),
      )
    }),
  )
}

beforeEach(() => {
  permintaan = null
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('simpanKeRNM', () => {
  it('memakai POST di jalur klaim, tanpa badan', async () => {
    jawab(200, { nomorKlaim: 'UJI-1', nomorBaru: false, barisDitandai: 0, arasapas: 'dilewati' })
    const hasil = await simpanKeRNM('K/1')
    expect(permintaan?.init.method).toBe('POST')
    expect(permintaan?.url).toContain('/api/klaim-life/K%2F1/outstanding')
    expect(hasil.nomorKlaim).toBe('UJI-1')
  })

  it('422 meneruskan kalimat XML apa adanya', async () => {
    jawab(422, { galat: 'DOB cannnot be blank No 1', langkah: '11.12' })
    const e = await simpanKeRNM('K-1').catch((x: unknown) => x)
    expect(e).toBeInstanceOf(ApiFailure)
    expect((e as ApiFailure).status).toBe(422)
    expect(pesanGalat(e)).toBe('DOB cannnot be blank No 1')
  })
})
