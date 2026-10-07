import { describe, expect, it } from 'vitest'

import { nilai, setel, type Halaman } from './api'
import { buatPenjagaIsian } from './penjagaIsian'

// Tinjauan kode 06-10-2026: jawaban server mengganti seluruh halaman kerja.

function tertunda<T>() {
  let selesai!: (x: T) => void
  const janji = new Promise<T>((r) => (selesai = r))
  return { janji, selesai }
}

const kosong: Halaman = { nilai: {}, daftar: {} }
const halaman = (isi: Record<string, string>): Halaman =>
  Object.entries(isi).reduce((h, [j, v]) => setel(h, j, v), kosong)

describe('penjaga isian layar kasus', () => {
  it('isian yang diketik selama permintaan berjalan diterapkan ulang di atas jawaban', async () => {
    const p = buatPenjagaIsian(halaman({ 'PolicyTreatyIn.Suggest': '' }))
    const jawab = tertunda<Halaman>()
    const jalan = p.kirim(
      () => jawab.janji,
      (x) => x,
    )
    await Promise.resolve()
    p.ubah((h) => setel(h, 'PolicyTreatyIn.Suggest', 'UJI-catatan'))
    jawab.selesai(halaman({ 'PolicyTreatyIn.Suggest': '', 'PolicyTreatyIn.DueTo': 'UJI-DUE' }))
    await jalan
    const h = p.kini() as Halaman
    expect(nilai(h, 'PolicyTreatyIn.Suggest')).toBe('UJI-catatan')
    expect(nilai(h, 'PolicyTreatyIn.DueTo')).toBe('UJI-DUE')
  })

  it('permintaan berurutan: yang kedua berangkat SESUDAH jawaban pertama, dengan halaman terbaru', async () => {
    const p = buatPenjagaIsian(halaman({ 'PolicyTreatyIn.FlagPPH': '0' }))
    const satu = tertunda<Halaman>()
    const dikirim: string[] = []
    const a = p.kirim(
      (h) => {
        dikirim.push(nilai(h, 'PolicyTreatyIn.FlagPPH'))
        return satu.janji
      },
      (x) => x,
    )
    await Promise.resolve()
    p.ubah((h) => setel(h, 'PolicyTreatyIn.FlagPPH', '1'))
    const b = p.kirim(
      (h) => {
        dikirim.push(nilai(h, 'PolicyTreatyIn.FlagPPH') + '|' + nilai(h, 'PolicyTreatyIn.PPHValue'))
        return Promise.resolve(setel(h, 'PolicyTreatyIn.PPHValue', '2'))
      },
      (x) => x,
    )
    satu.selesai(halaman({ 'PolicyTreatyIn.FlagPPH': '0', 'PolicyTreatyIn.PPHValue': '0' }))
    await Promise.all([a, b])
    // kedua berangkat membawa centang baru di atas jawaban pertama; jawaban pertama tidak menimpa yang kedua
    expect(dikirim).toEqual(['0', '1|0'])
    expect(nilai(p.kini() as Halaman, 'PolicyTreatyIn.FlagPPH')).toBe('1')
    expect(nilai(p.kini() as Halaman, 'PolicyTreatyIn.PPHValue')).toBe('2')
  })

  it('permintaan gagal tidak menahan antrean dan tidak membuang isian', async () => {
    const p = buatPenjagaIsian(halaman({ 'PolicyTreatyIn.Suggest': '' }))
    const gagal = p.kirim(
      () => Promise.reject(new Error('UJI-gagal')),
      (x: Halaman) => x,
    )
    p.ubah((h) => setel(h, 'PolicyTreatyIn.Suggest', 'UJI-tetap'))
    await expect(gagal).rejects.toThrow('UJI-gagal')
    await p.kirim(
      (h) => Promise.resolve(h),
      () => null,
    )
    expect(nilai(p.kini() as Halaman, 'PolicyTreatyIn.Suggest')).toBe('UJI-tetap')
  })
})
