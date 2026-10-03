import { describe, expect, it } from 'vitest'

import {
  ambilKasus,
  ambilObjek,
  badanPremiCargo,
  buatOpportunity,
  cariAccount,
  cariSOB,
  daftarCaseNB,
  daftarClassOfBusiness,
  daftarMarketing,
  simpanGeneral,
  simpanObjek,
} from './api'

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

describe('buatOpportunity (tiket 29)', () => {
  it('POST /api/nbfacin/opportunity dengan isian apa adanya, jawaban caseId', async () => {
    const asli = globalThis.fetch
    let url = ''
    let metode = ''
    let badan = ''
    globalThis.fetch = (async (u: RequestInfo | URL, init?: RequestInit) => {
      url = String(u)
      metode = init?.method ?? ''
      badan = String(init?.body ?? '')
      return new Response(JSON.stringify({ caseId: 'NB-1' }), { status: 201 })
    }) as typeof fetch
    const isian = {
      estimatedClosingDate: '10-03-2026', businessProspectName: 'UJI', accountId: 'UJI-A', insuredId: 'UJI-I',
      groupBusinessId: 'UJI-G', groupBusiness: 'UJI GRUP', classOfBusiness: 'UJI COB', typeOfInward: 'Facultative',
      typeOfFacultative: 'Facultative In', phase: 'Proposal', stage: 'Opportunity', opportunitySource: '',
      businessStatus: 'New Business', description: '',
    }
    try {
      const h = await buatOpportunity(isian)
      expect(metode).toBe('POST')
      expect(url).toMatch(/\/api\/nbfacin\/opportunity$/)
      expect(JSON.parse(badan)).toEqual(isian)
      expect(h).toEqual({ caseId: 'NB-1' })
    } finally {
      globalThis.fetch = asli
    }
  })
})

describe('kasus NB - baca, simpan General, marketing officer (tiket 31)', () => {
  it('rute dan metode sesuai kontrak; caseId di-encode', async () => {
    const asli = globalThis.fetch
    const panggil: { url: string; metode: string; badan: string }[] = []
    globalThis.fetch = (async (u: RequestInfo | URL, init?: RequestInit) => {
      panggil.push({ url: String(u), metode: init?.method ?? '', badan: String(init?.body ?? '') })
      return new Response(JSON.stringify({ baris: [] }), { status: 200 })
    }) as typeof fetch
    try {
      await ambilKasus('NB-1')
      await simpanGeneral('NB-1', {
        reffNumber: 'UJI', qqName: '', beginDate: '01-10-2026', offeringDate: '03-10-2026', endDate: '',
        policyType: '', marketingId: 'UJI-M', day: '', typeFacultative: 'Facultative In', sourceOfBusinessId: 'UJI-S', cedingIds: ['UJI-C1', 'UJI-C2'],
      })
      await daftarMarketing()
      expect(panggil.map((p) => [p.metode, p.url.replace(/^.*(\/api\/)/, '$1')])).toEqual([
        ['GET', '/api/nbfacin/kasus/NB-1'],
        ['PUT', '/api/nbfacin/kasus/NB-1/general'],
        ['GET', '/api/nbfacin/marketing-officer'],
      ])
      expect(JSON.parse(panggil[1]!.badan).marketingId).toBe('UJI-M')
    } finally {
      globalThis.fetch = asli
    }
  })
})

describe('daftarCaseNB (tiket 32)', () => {
  it('GET /api/nbfacin/opportunity dengan cari (dipangkas) dan halaman', async () => {
    const asli = globalThis.fetch
    let url = ''
    globalThis.fetch = (async (u: RequestInfo | URL) => {
      url = String(u)
      return new Response(JSON.stringify({ baris: [], total: 0, halaman: 1, ukuran: 15 }), { status: 200 })
    }) as typeof fetch
    try {
      await daftarCaseNB('  NB-1  ', 1)
      expect(url).toMatch(/\/api\/nbfacin\/opportunity\?/)
      const q = new URL(url, 'http://x').searchParams
      expect(q.get('cari')).toBe('NB-1')
      expect(q.get('halaman')).toBe('1')
    } finally {
      globalThis.fetch = asli
    }
  })
})

describe('cariSOB (tiket 33)', () => {
  it('GET /api/nbfacin/sob dengan cari (dipangkas) dan halaman', async () => {
    const asli = globalThis.fetch
    let url = ''
    globalThis.fetch = (async (u: RequestInfo | URL) => {
      url = String(u)
      return new Response(JSON.stringify({ baris: [], total: 0, halaman: 2, ukuran: 20 }), { status: 200 })
    }) as typeof fetch
    try {
      await cariSOB('  uji  ', 2)
      expect(url).toMatch(/\/api\/nbfacin\/sob\?/)
      const q = new URL(url, 'http://x').searchParams
      expect(q.get('cari')).toBe('uji')
      expect(q.get('halaman')).toBe('2')
    } finally {
      globalThis.fetch = asli
    }
  })
})

describe('tab Object (tiket 35)', () => {
  it('GET dan PUT /api/nbfacin/kasus/{caseId}/objek; PUT membawa { baris }', async () => {
    const asli = globalThis.fetch
    const panggil: { url: string; metode: string; badan: string }[] = []
    globalThis.fetch = (async (u: RequestInfo | URL, init?: RequestInit) => {
      panggil.push({ url: String(u), metode: init?.method ?? '', badan: String(init?.body ?? '') })
      return new Response(JSON.stringify({ baris: [] }), { status: 200 })
    }) as typeof fetch
    try {
      await ambilObjek('NB-1')
      await simpanObjek('NB-1', [])
      expect(panggil.map((p) => [p.metode, p.url.replace(/^.*(\/api\/)/, '$1')])).toEqual([
        ['GET', '/api/nbfacin/kasus/NB-1/objek'],
        ['PUT', '/api/nbfacin/kasus/NB-1/objek'],
      ])
      expect(JSON.parse(panggil[1]!.badan)).toEqual({ baris: [] })
    } finally {
      globalThis.fetch = asli
    }
  })
})
