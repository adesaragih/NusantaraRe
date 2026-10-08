// Tombol `Copy` daftar kontrak — `Section/InputTreatyInOffer.xml` cell 993 →
// `SetTreatyIn_Act(ID=.ID)` → `Activity/TreatyInCopy.xml` (8 Oktober 2026).
//
// ⭐ Copy SENDIRI tidak menulis: `TreatyInCopy` hanya `Property-Set` dan
// `RDB-List GetCurrentDate`. Ia membuka DRAF kontrak baru (`ID =
// "UnknownId"`, `OLDID` = sumber), dan draf itu masuk tabel lewat Save.
// Uji ini menjaga keempat sambungannya: daftar → rute → form → api.

import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { ambilDraftSalinan, simpanSalinan } from './api'
import { bolehRevisi } from './pages/DaftarKontrakTreatyIn'

const baca = (...p: string[]) => readFileSync(join(__dirname, ...p), 'utf8').replace(/\r\n/g, '\n')
const DAFTAR = baca('pages', 'DaftarKontrakTreatyIn.tsx')
const RUTE = baca('rute.tsx')
const FORM = baca('pages', 'FormKontrakTreatyIn.tsx')

describe('tombol Copy — TreatyInCopy', () => {
  it('syarat tampil cell 993 SAMA dengan cell 994 (Revision)', () => {
    // `(WB = 'ReasTreatyInAdmin' && .Position = '' && .StatusAkseptasi = 'Resolve Complete')`
    expect(DAFTAR).toContain(
      '.filter((a) => (a !== DAFTAR_KONTRAK.revisi && a !== DAFTAR_KONTRAK.salin) || bolehRevisi(b, workbasket))',
    )
    const tuntas = { statusAkseptasi: 'Resolve Complete', posisi: '' }
    expect(bolehRevisi(tuntas, ['ReasTreatyInAdmin'])).toBe(true)
    expect(bolehRevisi(tuntas, ['ReasTreatyInSecHead'])).toBe(false)
  })

  it('Copy membuka DRAF lewat onSalin — event `click`, nol fungsi tulis di daftar', () => {
    expect(DAFTAR).toMatch(/onClick=\{\(\) => \{\s*if \(a === DAFTAR_KONTRAK\.salin && onSalin !== undefined\) onSalin\(b\.id\)/)
    expect(DAFTAR).not.toContain('simpanSalinan')
    expect(DAFTAR).not.toContain('ambilDraftSalinan')
  })

  it('rute: draf = pengenal KOSONG + sumber, mode ubah (ViewState = 0, IsEditData = 0)', () => {
    expect(RUTE).toMatch(/onSalin=\{\(id\) => \{[\s\S]*?setSalinDari\(id\)\s*setMode\('ubah'\)\s*setDibuka\(''\)/)
    expect(RUTE).toContain('salinDari={salinDari ?? undefined}')
    // Sesudah Save, pengenal BARU dibuka dan drafnya dilupakan.
    expect(RUTE).toMatch(/onTersimpan=\{\(id\) => \{\s*setSalinDari\(null\)\s*setDibuka\(id\)/)
  })

  it('form: draf dimuat dari sumbernya, dan HANYA Save/Submit/Decline menulis', () => {
    expect(FORM).toContain("const sumberSalin = idKontrak === '' ? (salinDari ?? '') : ''")
    expect(FORM).toContain("sumberSalin !== '' ? ambilDraftSalinan(sumberSalin) : ambilKontrakWarisan(idKontrak)")
    expect(FORM).toContain("sumberSalin !== '' ? simpanSalinan(isiSalinan('')) : simpanKontrak(isiTombol())")
    expect(FORM).toContain('simpanSalinan(isiSalinan(aksi, tambahan))')
    // Dua pemanggilan saja — keduanya di dalam `tekanTulis`.
    expect(FORM.match(/simpanSalinan\(/g)?.length).toBe(2)
  })
})

describe('api Copy', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  const tangkap = () => {
    const panggil: { url: string; metode: string; badan: unknown }[] = []
    vi.stubGlobal(
      'fetch',
      vi.fn((url: string, init: RequestInit) => {
        panggil.push({
          url: String(url),
          metode: init.method ?? 'GET',
          badan: typeof init.body === 'string' ? JSON.parse(init.body) : undefined,
        })
        return Promise.resolve(new Response('{"id":"1001857"}', { status: 200 }))
      }),
    )
    return panggil
  }

  it('draf: GET tanpa badan — nol tulisan', async () => {
    const p = tangkap()
    await ambilDraftSalinan('1001001')
    expect(p).toHaveLength(1)
    expect(p[0]?.url).toContain('/api/treaty-in/kontrak-warisan/1001001/salin')
    expect(p[0]?.metode).toBe('GET')
    expect(p[0]?.badan).toBeUndefined()
  })

  it('Save draf: POST /kontrak/salin membawa sumber dan isi layar', async () => {
    const p = tangkap()
    const h = await simpanSalinan({ idSumber: '1001001', dokumen: { TreatyContractName: 'X' }, aksi: '' })
    expect(h.id).toBe('1001857')
    expect(p[0]?.url).toContain('/api/treaty-in/kontrak/salin')
    expect(p[0]?.metode).toBe('POST')
    expect(p[0]?.badan).toEqual({ idSumber: '1001001', dokumen: { TreatyContractName: 'X' }, aksi: '' })
  })
})
