import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { cariPesertaLife } from './api'

// Uji BENTUK URL pencarian peserta — `Find Insured` layar Register.
//
// Yang dijaga bukan hasilnya melainkan parameter yang benar-benar terkirim:
// kotak yang dibiarkan kosong TIDAK boleh menjadi `sertifikat=` kosong di URL.
// Di backend, parameter kosong dan parameter tidak terkirim memang berarti
// hal yang sama — tetapi menyamakannya di SATU tempat saja berarti sisi
// satunya bebas berubah tanpa ada yang menyadarinya.

/** URL terakhir yang diminta fetch. */
let terakhir = ''

beforeEach(() => {
  terakhir = ''
  vi.stubGlobal(
    'fetch',
    vi.fn((url: string) => {
      terakhir = url
      return Promise.resolve(
        new Response('[]', {
          status: 200,
          headers: { 'Content-Type': 'application/json' },
        }),
      )
    }),
  )
})

afterEach(() => {
  vi.unstubAllGlobals()
})

/** Parameter kueri dari URL terakhir. */
function kueri(): URLSearchParams {
  return new URL(terakhir, 'http://uji.invalid').searchParams
}

describe('cariPesertaLife menyusun kueri', () => {
  it('tanpa penyaring: hanya pl dan n', async () => {
    await cariPesertaLife('UJI-PL-1')
    expect(kueri().get('pl')).toBe('UJI-PL-1')
    expect(kueri().get('n')).toBe('50')
    // ⛔ Bukan string kosong — TIDAK ADA sama sekali.
    expect(kueri().has('sertifikat')).toBe(false)
    expect(kueri().has('nama')).toBe(false)
  })

  it('kotak berisi spasi saja dianggap kosong', async () => {
    await cariPesertaLife('UJI-PL-1', { sertifikat: '   ', nama: '  ' })
    expect(kueri().has('sertifikat')).toBe(false)
    expect(kueri().has('nama')).toBe(false)
  })

  it('kedua penyaring terkirim, spasi tepi dibuang', async () => {
    await cariPesertaLife('UJI-PL-1', { sertifikat: '  UJI-006 ', nama: ' budi ' })
    expect(kueri().get('sertifikat')).toBe('UJI-006')
    expect(kueri().get('nama')).toBe('budi')
  })

  it('huruf besar nama BUKAN urusan layar', async () => {
    // b405 `@toUpperCase` dikerjakan server. Bila layar ikut meng-uppercase,
    // aturannya punya dua rumah, dan yang satu akan bergeser diam-diam.
    await cariPesertaLife('UJI-PL-1', { nama: 'budi' })
    expect(kueri().get('nama')).toBe('budi')
  })

  it('satu penyaring saja boleh', async () => {
    await cariPesertaLife('UJI-PL-1', { nama: 'ani' })
    expect(kueri().has('sertifikat')).toBe(false)
    expect(kueri().get('nama')).toBe('ani')
  })

  it('nomor premium list selalu ikut, sebab tabelnya 66,8 juta baris', async () => {
    await cariPesertaLife('UJI-PL-9', { sertifikat: 'X' })
    expect(kueri().get('pl')).toBe('UJI-PL-9')
  })

  it('batas hasil dapat diatur dan selalu terkirim', async () => {
    await cariPesertaLife('UJI-PL-1', {}, 10)
    expect(kueri().get('n')).toBe('10')
  })

  it('daftar kosong dari Go (null) menjadi array kosong', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(() =>
        Promise.resolve(
          new Response('null', {
            status: 200,
            headers: { 'Content-Type': 'application/json' },
          }),
        ),
      ),
    )
    await expect(cariPesertaLife('UJI-PL-1')).resolves.toEqual([])
  })
})
