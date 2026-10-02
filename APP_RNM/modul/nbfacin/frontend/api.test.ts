import { describe, expect, it } from 'vitest'

import { badanPremiCargo, cariAccount, daftarClassOfBusiness } from './api'

describe('badanPremiCargo', () => {
  it('angka tetap teks, lini MARINE CARGO, bukan master policy', () => {
    const b = badanPremiCargo({ mataUang: ' IDR ', rate: '0.125', tsi: '1000000.10' })
    expect(b).toEqual({ liniBisnis: 'MARINE CARGO', mataUang: 'IDR', tsi: '1000000.10', rate: '0.125', masterPolicy: false })
    expect(typeof b.tsi).toBe('string')
    expect(typeof b.rate).toBe('string')
  })

  it('tidak membulatkan atau menormalkan angka', () => {
    // 0.1 + 0.2 di float = 0.30000000000000004; teks harus utuh.
    expect(badanPremiCargo({ mataUang: 'IDR', rate: '0.30', tsi: '0000123' }).rate).toBe('0.30')
    expect(badanPremiCargo({ mataUang: 'IDR', rate: '1', tsi: '0000123' }).tsi).toBe('0000123')
  })
})

describe('cariAccount (tiket 26 C-8 / tiket 27)', () => {
  it('GET /api/nbfacin/account dengan cari (dipangkas) dan halaman', async () => {
    const asli = globalThis.fetch
    let url = ''
    let metode = ''
    globalThis.fetch = (async (u: RequestInfo | URL, init?: RequestInit) => {
      url = String(u)
      metode = init?.method ?? ''
      return new Response(JSON.stringify({ baris: [], total: 0, halaman: 2, ukuran: 20 }), { status: 200 })
    }) as typeof fetch
    try {
      const h = await cariAccount('  UJI a  ', 2)
      expect(metode).toBe('GET')
      expect(url).toMatch(/\/api\/nbfacin\/account\?/)
      const q = new URL(url, 'http://x').searchParams
      expect(q.get('cari')).toBe('UJI a')
      expect(q.get('halaman')).toBe('2')
      expect(h).toEqual({ baris: [], total: 0, halaman: 2, ukuran: 20 })
    } finally {
      globalThis.fetch = asli
    }
  })
})

describe('daftarClassOfBusiness (tiket 26 C-11 / tiket 28)', () => {
  it('GET /api/nbfacin/class-of-business dengan groupBusinessId', async () => {
    const asli = globalThis.fetch
    let url = ''
    globalThis.fetch = (async (u: RequestInfo | URL) => {
      url = String(u)
      return new Response(JSON.stringify({ baris: [{ id: 'UJI-1', note: 'UJI COB' }] }), { status: 200 })
    }) as typeof fetch
    try {
      const h = await daftarClassOfBusiness('UJI-GRUP')
      expect(url).toMatch(/\/api\/nbfacin\/class-of-business\?/)
      expect(new URL(url, 'http://x').searchParams.get('groupBusinessId')).toBe('UJI-GRUP')
      expect(h.baris).toEqual([{ id: 'UJI-1', note: 'UJI COB' }])
    } finally {
      globalThis.fetch = asli
    }
  })
})
