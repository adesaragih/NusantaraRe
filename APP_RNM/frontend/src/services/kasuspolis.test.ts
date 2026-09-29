// Tombol portal PremiumList membuat kasus — GILIRAN-13 paket 1 (butir bn).

import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { buatKasusPolis, FLAG_POLIS } from './api'

interface Tangkapan {
  url: string
  init: RequestInit
}

let tertangkap: Tangkapan[] = []

beforeEach(() => {
  tertangkap = []
  vi.stubEnv('VITE_AUTH_STUB', 'true')
  vi.stubEnv('VITE_STUB_PELAKU', 'UJI-INPUTOR')
  vi.stubEnv('VITE_STUB_PERAN', '')
  vi.stubGlobal('fetch', (url: string, init: RequestInit) => {
    tertangkap.push({ url, init })
    return Promise.resolve(
      new Response(
        JSON.stringify({ caseId: 'NBLF-1', statusWork: 'Input Offer Life', position: 'Offer', flag: '1' }),
        { status: 201, headers: { 'Content-Type': 'application/json' } },
      ),
    )
  })
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.unstubAllEnvs()
})

describe('buatKasusPolis', () => {
  it('nilai bendera VERBATIM FlagPolicy tombol portal', () => {
    // Section/PremiumList.xml b3310/b3597 ("0") dan b3958/b4233 ("1").
    expect(FLAG_POLIS.inputOffer).toBe('0')
    expect(FLAG_POLIS.inputPremium).toBe('1')
  })

  it('POST ke /api/polis-life membawa bendera dan identitas', async () => {
    const hasil = await buatKasusPolis(FLAG_POLIS.inputPremium)
    expect(tertangkap).toHaveLength(1)
    const t = tertangkap[0]!
    expect(t.url).toBe('/api/polis-life')
    expect(t.init.method).toBe('POST')
    expect(JSON.parse(String(t.init.body))).toEqual({ flag: '1' })
    expect((t.init.headers as Record<string, string>)['X-Pelaku']).toBe('UJI-INPUTOR')
    expect(hasil.caseId).toBe('NBLF-1')
    expect(hasil.statusWork).toBe('Input Offer Life')
  })
})
