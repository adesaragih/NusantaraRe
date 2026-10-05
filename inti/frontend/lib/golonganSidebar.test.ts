import { describe, expect, it } from 'vitest'

import {
  alihkanGolongan,
  bacaGolonganTertutup,
  bukaGolongan,
  KUNCI_GOLONGAN_TERTUTUP,
  simpanGolonganTertutup,
} from './golonganSidebar'

// Golongan sidebar yang dapat dibuka-tutup (permintaan work owner 05-10-2026).

function gudang(awal: Record<string, string> = {}) {
  const isi = { ...awal }
  return {
    isi,
    getItem: (k: string) => isi[k] ?? null,
    setItem: (k: string, v: string) => {
      isi[k] = v
    },
  }
}

describe('golongan sidebar dibuka-tutup', () => {
  it('bawaan: semua golongan terbuka', () => {
    expect([...bacaGolonganTertutup(gudang())]).toEqual([])
    expect([...bacaGolonganTertutup(null)]).toEqual([])
  })

  it('klik kepala golongan menutup lalu membuka lagi', () => {
    const tutup = alihkanGolongan(new Set(), 'MASTER')
    expect([...tutup]).toEqual(['MASTER'])
    expect([...alihkanGolongan(tutup, 'MASTER')]).toEqual([])
    expect([...alihkanGolongan(tutup, 'KLAIM')].sort()).toEqual(['KLAIM', 'MASTER'])
  })

  it('pilihan diingat peramban: simpan lalu baca menghasilkan himpunan yang sama', () => {
    const g = gudang()
    simpanGolonganTertutup(new Set(['MASTER TREATY', 'KLAIM']), g)
    expect(g.isi[KUNCI_GOLONGAN_TERTUTUP]).toBe('["KLAIM","MASTER TREATY"]')
    expect([...bacaGolonganTertutup(g)].sort()).toEqual(['KLAIM', 'MASTER TREATY'])
  })

  it('isi penyimpanan rusak atau bukan daftar teks = semua terbuka', () => {
    for (const rusak of ['{', '"MASTER"', '{"a":1}', '[1, null]']) {
      expect([...bacaGolonganTertutup(gudang({ [KUNCI_GOLONGAN_TERTUTUP]: rusak }))], rusak).toEqual([])
    }
    expect([...bacaGolonganTertutup(gudang({ [KUNCI_GOLONGAN_TERTUTUP]: '["MASTER", 3]' }))]).toEqual(['MASTER'])
  })

  it('penyimpanan yang melempar galat tidak merusak sidebar', () => {
    const melempar = {
      getItem: () => {
        throw new Error('diblokir')
      },
      setItem: () => {
        throw new Error('penuh')
      },
    }
    expect([...bacaGolonganTertutup(melempar)]).toEqual([])
    expect(() => {
      simpanGolonganTertutup(new Set(['MASTER']), melempar)
    }).not.toThrow()
  })

  it('golongan halaman yang tampil dibuka; yang sudah terbuka tidak membuat himpunan baru', () => {
    const tertutup = new Set(['MASTER', 'KLAIM'])
    expect([...bukaGolongan(tertutup, 'MASTER')]).toEqual(['KLAIM'])
    expect(bukaGolongan(tertutup, 'TREATY')).toBe(tertutup)
  })
})
