import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { ApiFailure, pesanGalat } from '../../../inti/frontend/klien'
import { GagalSimpanRNM, pindahTahap, simpanKeRNM, simpanMedis } from './api'

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

describe('simpanKeRNM: seluruh kekurangan kolom wajib (11-10-2026)', () => {
  it('422 dengan kekurangan -> GagalSimpanRNM membawa daftar lengkap, pesan server tetap', async () => {
    const kekurangan = [
      { nomorPeserta: 1, pesertaId: 'P1', sertifikat: 'UJI-1', nama: 'UJI A', kolom: 'Date of Loss', pesan: 'x' },
      { nomorPeserta: 2, pesertaId: 'P2', sertifikat: 'UJI-2', nama: 'UJI B', kolom: 'Claim RNM', pesan: 'y' },
    ]
    jawab(422, { galat: 'x\ny', langkah: 'wajib', kekurangan })
    const e = await simpanKeRNM('K-1').catch((x: unknown) => x)
    expect(e).toBeInstanceOf(GagalSimpanRNM)
    expect(e).toBeInstanceOf(ApiFailure)
    expect((e as GagalSimpanRNM).kekurangan).toEqual(kekurangan)
    expect(pesanGalat(e)).toBe('x\ny')
  })
  it('422 tanpa kekurangan (klaim ganda) -> ApiFailure biasa', async () => {
    jawab(422, { galat: 'Person number 1 has already been accepted.', langkah: '11.4', kekurangan: [] })
    const e = await simpanKeRNM('K-1').catch((x: unknown) => x)
    expect(e).toBeInstanceOf(ApiFailure)
    expect(e).not.toBeInstanceOf(GagalSimpanRNM)
  })
})

describe('pindahTahap: Send to Medical Check berkekurangan (11-10-2026)', () => {
  it('422 berdaftar -> GagalSimpanRNM; 409 biasa -> ApiFailure', async () => {
    const kekurangan = [{ nomorPeserta: 1, pesertaId: 'P1', sertifikat: 'S1', nama: 'UJI A', kolom: 'Document', pesan: 'x' }]
    jawab(422, { galat: 'x', langkah: 'dokumen-medis', kekurangan })
    const e = await pindahTahap('K-1', 'medical-check').catch((x: unknown) => x)
    expect(e).toBeInstanceOf(GagalSimpanRNM)
    expect((e as GagalSimpanRNM).kekurangan).toEqual(kekurangan)
    expect(permintaan?.init.method).toBe('POST')
    expect(permintaan?.url).toContain('/api/klaim-life/K-1/tahap/medical-check')
    jawab(409, { galat: 'menunggu Save to RNM' })
    const e2 = await pindahTahap('K-1', 'medical-check').catch((x: unknown) => x)
    expect(e2).toBeInstanceOf(ApiFailure)
    expect(e2).not.toBeInstanceOf(GagalSimpanRNM)
  })
})

describe('simpanMedis: tombol Save layar Medical Check (11-10-2026)', () => {
  it('PUT /medis membawa Recommendation per peserta; 403 diteruskan sebagai ApiFailure', async () => {
    jawab(200, {})
    await simpanMedis('K/1', { rekomendasi: { P1: 'UJI REKOMENDASI' } })
    expect(permintaan?.init.method).toBe('PUT')
    expect(permintaan?.url).toContain('/api/klaim-life/K%2F1/medis')
    expect(JSON.parse(String(permintaan?.init.body))).toEqual({ rekomendasi: { P1: 'UJI REKOMENDASI' } })
    jawab(403, { galat: 'Only the Medical Advisor can save the Medical Check.' })
    const e = await simpanMedis('K-1', { rekomendasi: { P1: 'x' } }).catch((x: unknown) => x)
    expect(e).toBeInstanceOf(ApiFailure)
    expect(pesanGalat(e)).toBe('Only the Medical Advisor can save the Medical Check.')
  })
})
