// Asal: salinan `modul/nbtreatyin/frontend/ambil.test.ts` (06-10-2026), kelas `nbti__` -> `edmt__`.
// Efek ambil-dengan-batal: jawaban permintaan usang dibuang.

import { describe, expect, it } from 'vitest'

import { ambilBatal } from './ambil'

describe('ambilBatal', () => {
  it('meneruskan hasil bila belum dibatalkan', async () => {
    let hasil = ''
    ambilBatal(
      () => Promise.resolve('UJI-1'),
      (x) => (hasil = x),
      () => undefined,
    )
    await Promise.resolve()
    await Promise.resolve()
    expect(hasil).toBe('UJI-1')
  })

  it('membuang hasil dan galat sesudah dibatalkan', async () => {
    let hasil = ''
    let galat: unknown = null
    const batal1 = ambilBatal(
      () => Promise.resolve('UJI-usang'),
      (x) => (hasil = x),
      (e) => (galat = e),
    )
    const batal2 = ambilBatal(
      () => Promise.reject(new Error('UJI-galat')),
      (x) => (hasil = x),
      (e) => (galat = e),
    )
    batal1()
    batal2()
    await new Promise((r) => setTimeout(r, 0))
    expect(hasil).toBe('')
    expect(galat).toBeNull()
  })
})
