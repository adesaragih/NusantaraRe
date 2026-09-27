import { beforeEach, describe, expect, it, vi } from 'vitest'

import {
  BATAS_BARIS_PENYAKIT,
  UKURAN_HALAMAN_PENYAKIT,
  cariPenyakit,
} from './api'

// Uji sisi klien pencarian diagnosa — kelompok Medis.
//
// ⛔ Tabelnya 97.586 baris. Yang dijaga berkas ini: tidak ada permintaan yang
// berangkat tanpa batas, dan batas yang berlebihan dijepit sebelum dikirim.

function jawab(isi: unknown): Response {
  return new Response(JSON.stringify(isi), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  })
}

function urlDari(palsu: { mock: { calls: unknown[][] } }): URL {
  const panggilan = palsu.mock.calls[0]
  expect(panggilan).toBeDefined()
  return new URL(String(panggilan?.[0]), 'http://uji.invalid')
}

describe('cariPenyakit', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('selalu mengirim batas, bahkan tanpa kata kunci', () => {
    // ⛔ Justru pencarian KOSONG yang paling harus berbatas: di Pega
    // `Contains ""` cocok dengan seluruh 97.586 baris, dan yang menahannya
    // hanya `pyMaxRecords`.
    const palsu = vi
      .spyOn(globalThis, 'fetch')
      .mockImplementation(() => Promise.resolve(jawab([])))

    return cariPenyakit('', '').then(() => {
      const u = urlDari(palsu)
      expect(u.pathname).toContain('/api/penyakit-life')
      expect(u.searchParams.get('batas')).toBe(String(UKURAN_HALAMAN_PENYAKIT))
    })
  })

  it('batas berlebihan dijepit ke pyMaxRecords, bukan diteruskan', async () => {
    const palsu = vi
      .spyOn(globalThis, 'fetch')
      .mockImplementation(() => Promise.resolve(jawab([])))

    await cariPenyakit('', '', 1_000_000)
    expect(urlDari(palsu).searchParams.get('batas')).toBe(
      String(BATAS_BARIS_PENYAKIT),
    )
  })

  it('batas nol atau negatif tidak pernah berangkat sebagai nol', async () => {
    // `batas=0` di backend berarti "pakai ukuran halaman", tetapi mengirim
    // nol dari sini menggantungkan artinya pada backend. Dijepit ke 1 di
    // sisi ini, dan backend tetap menjepitnya lagi.
    const palsu = vi
      .spyOn(globalThis, 'fetch')
      .mockImplementation(() => Promise.resolve(jawab([])))

    await cariPenyakit('', '', -5)
    const batas = Number(urlDari(palsu).searchParams.get('batas'))
    expect(batas).toBeGreaterThan(0)
  })

  it('kedua kata kunci dikirim sebagai parameter terpisah', async () => {
    // ⛔ `icd` dan `nama` TERPISAH, sebab backend menyambungnya dengan AND
    // (b535 `A AND B`). Menggabungkannya menjadi satu kotak pencarian akan
    // mengubah maknanya, dan bedanya baru terlihat pada dua kata kunci.
    const palsu = vi
      .spyOn(globalThis, 'fetch')
      .mockImplementation(() => Promise.resolve(jawab([])))

    await cariPenyakit('A00', 'diabetes')
    const u = urlDari(palsu)
    expect(u.searchParams.get('icd')).toBe('A00')
    expect(u.searchParams.get('nama')).toBe('diabetes')
  })

  it('angkanya dari rule, bukan dari selera kami', () => {
    expect(BATAS_BARIS_PENYAKIT).toBe(500)
    expect(UKURAN_HALAMAN_PENYAKIT).toBe(50)
  })

  it('hasilnya dibawa apa adanya, nomor tetap TEKS', async () => {
    // ADR-U-0022: pengenal adalah teks. "007" bukan "7", dan Number() akan
    // menyamakan keduanya.
    vi.spyOn(globalThis, 'fetch').mockImplementation(() =>
      Promise.resolve(jawab([{ nomor: '007', nama: 'UJI-PENYAKIT', kodeIcd: 'A00' }])),
    )
    const hasil = await cariPenyakit('A00', '')
    expect(hasil).toHaveLength(1)
    expect(hasil[0]?.nomor).toBe('007')
    expect(typeof hasil[0]?.nomor).toBe('string')
  })
})
