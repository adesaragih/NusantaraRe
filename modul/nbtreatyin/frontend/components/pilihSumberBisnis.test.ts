// Uji syarat tampil pemilih Source Of Business - VERBATIM XML:
//   tombol `Select Source Of Business` (`Section/DetailPolicyTreatyIn`) pyVisible `.ClaimType = 'XOL Retro'`
//   tombol `Choose` (`Section/SourceHierarki`) pyVisible `.ChildCount = 0`
// dan F4 (IKUTI XML): hasil `SearchHierarkiSourceBizAgent_PostDT` dipegang state layar
// lalu ikut terkirim pada Save/Submit - klik tidak menyimpan.

import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  PREFIX_NBTREATYIN,
  kirimKasus,
  pilihSumberBisnis,
  simpanKasus,
  type Halaman,
  type HasilSumberBisnis,
} from '../api'
import { JALUR_SUMBER_BISNIS, pegangSumberBisnis, tampilChoose, tampilTombolSOB } from './PilihSumberBisnis'

const hal = (nilai: Record<string, string>): Halaman => ({ nilai, daftar: {} })

describe('pemilih Source Of Business', () => {
  it('tombol hanya tampil bila ClaimType persis "XOL Retro"', () => {
    expect(tampilTombolSOB(hal({ 'PolicyTreatyIn.ClaimType': 'XOL Retro' }))).toBe(true)
    expect(tampilTombolSOB(hal({ 'PolicyTreatyIn.ClaimType': 'XOL' }))).toBe(false)
    expect(tampilTombolSOB(hal({}))).toBe(false)
  })

  it('Choose hanya untuk simpul tanpa anak (ChildCount = 0; kosong = 0)', () => {
    expect(tampilChoose({ id: 'UJI-1', clientName: 'UJI', leader0: '', childCount: '0', clientId: '' })).toBe(true)
    expect(tampilChoose({ id: 'UJI-2', clientName: 'UJI', leader0: '', childCount: '', clientId: '' })).toBe(true)
    expect(tampilChoose({ id: 'UJI-3', clientName: 'UJI', leader0: '', childCount: '2', clientId: '' })).toBe(false)
  })
})

/** Satu panggilan `fetch` yang direkam. */
interface Rekam {
  url: string
  metode: string
  badan: { halaman?: Halaman; idAgen?: string } | undefined
}

function rekamFetch(jawab: unknown): Rekam[] {
  const rekam: Rekam[] = []
  vi.stubGlobal(
    'fetch',
    vi.fn(async (url: string, init: RequestInit) => {
      rekam.push({
        url,
        metode: init.method ?? 'GET',
        badan: init.body === undefined ? undefined : JSON.parse(String(init.body)),
      })
      return new Response(JSON.stringify(jawab), { status: 200 })
    }),
  )
  return rekam
}

afterEach(() => {
  vi.unstubAllGlobals()
})

// PostDT atas baris ChildCount 0: ID, ClientName, Leader0 (kosong - RD efektif NB
// `LEADER0 IS NULL`), Leader1 bukan kolom RD -> "".
const HASIL_AG1: HasilSumberBisnis = {
  sourceOfBusiness: 'UJI-AG-1',
  sobName: 'UJI SUMBER SATU',
  sobLeader0: '',
  sobLeader1: '',
}
// `@if(.ChildCount > 0, "", ...)`: simpul beranak mengosongkan keempatnya.
const HASIL_BERANAK: HasilSumberBisnis = { sourceOfBusiness: '', sobName: '', sobLeader0: '', sobLeader1: '' }

const layarXOLRetro = (): Halaman =>
  hal({
    'PolicyTreatyIn.ClaimType': 'XOL Retro',
    'PolicyTreatyIn.SOBName': 'UJI SOB KONTRAK',
    'Quotation.SourceOfBusiness': 'UJI-LAMA',
    'PolicyTreatyIn.QuotationData.SourceOfBusiness': 'UJI-LAMA',
  })

describe('F4: pilihan Source Of Business dipegang layar', () => {
  it('klik baris = POST pilih-sumber-bisnis {idAgen, halaman}; jawabannya keempat medan PostDT', async () => {
    const rekam = rekamFetch(HASIL_AG1)
    const h = layarXOLRetro()
    const hasil = await pilihSumberBisnis('UJI-NB 1', 'UJI-AG-1', h)
    expect(rekam).toEqual([
      {
        url: `${PREFIX_NBTREATYIN}/kasus/UJI-NB%201/pilih-sumber-bisnis`,
        metode: 'POST',
        badan: { idAgen: 'UJI-AG-1', halaman: h },
      },
    ])
    expect(hasil).toEqual(HASIL_AG1)
  })

  it('pegangSumberBisnis hanya menulis empat jalur PostDT', () => {
    expect(Object.values(JALUR_SUMBER_BISNIS)).toEqual([
      'Quotation.SourceOfBusiness',
      'Quotation.SobName',
      'Quotation.SobLeader0',
      'Quotation.SobLeader1',
    ])
    const h = layarXOLRetro()
    const selundupan = { ...HASIL_AG1, 'PolicyTreatyIn.SOBName': 'UJI-SELUNDUPAN' } as HasilSumberBisnis
    const dipegang = pegangSumberBisnis(h, selundupan)
    expect(dipegang.nilai).toEqual({
      ...h.nilai,
      'Quotation.SourceOfBusiness': 'UJI-AG-1',
      'Quotation.SobName': 'UJI SUMBER SATU',
      'Quotation.SobLeader0': '',
      'Quotation.SobLeader1': '',
    })
    // PostDT tidak menyentuh medan layar SOBName maupun salinan QuotationData.
    expect(dipegang.nilai['PolicyTreatyIn.SOBName']).toBe('UJI SOB KONTRAK')
    expect(dipegang.nilai['PolicyTreatyIn.QuotationData.SourceOfBusiness']).toBe('UJI-LAMA')
    expect(h.nilai['Quotation.SourceOfBusiness']).toBe('UJI-LAMA') // halaman asal tidak diubah
  })

  it('simpul beranak: jawaban kosong menimpa pilihan yang dipegang', () => {
    const dipegang = pegangSumberBisnis(pegangSumberBisnis(layarXOLRetro(), HASIL_AG1), HASIL_BERANAK)
    for (const j of Object.values(JALUR_SUMBER_BISNIS)) expect(dipegang.nilai[j]).toBe('')
  })

  it('pilihan yang dipegang ikut terkirim pada Save dan Submit', async () => {
    const rekam = rekamFetch(HASIL_AG1)
    const h = layarXOLRetro()
    const dipegang = pegangSumberBisnis(h, await pilihSumberBisnis('UJI-NB-1', 'UJI-AG-1', h))
    await simpanKasus('UJI-NB-1', dipegang)
    await kirimKasus('UJI-NB-1', dipegang)
    expect(rekam.map((r) => `${r.metode} ${r.url}`)).toEqual([
      `POST ${PREFIX_NBTREATYIN}/kasus/UJI-NB-1/pilih-sumber-bisnis`,
      `PUT ${PREFIX_NBTREATYIN}/kasus/UJI-NB-1`,
      `POST ${PREFIX_NBTREATYIN}/kasus/UJI-NB-1/kirim`,
    ])
    for (const r of rekam.slice(1)) {
      expect(r.badan?.halaman?.nilai['Quotation.SourceOfBusiness']).toBe('UJI-AG-1')
      expect(r.badan?.halaman?.nilai['Quotation.SobName']).toBe('UJI SUMBER SATU')
    }
  })
})
