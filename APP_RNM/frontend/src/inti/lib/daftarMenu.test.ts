import { describe, expect, it } from 'vitest'

import { bentukMenuTabel, HALAMAN_BERANDA, susunMenu, type EntriMenu, type MenuTabel } from './daftarMenu'

// Menu dari tabel M_NAV_MENU (brief menu 30-09-2026): pohon `GET /api/menu`
// DIPOTONG dengan rute frontend yang benar-benar terdaftar. Rute di sini
// tiruan - `inti/` tidak mengenal modul; uji dua arah dengan `modul/daftar.ts`
// ada di `modul/daftar.menuTabel.test.ts`.

type H = 'beranda' | 'inbox' | 'register' | 'tco-tahun' | 'rute-tanpa-baris'

const RUTE: readonly EntriMenu<H>[] = [
  { modul: HALAMAN_BERANDA, label: 'Beranda', kelompok: 'Beranda', pemilik: null },
  { modul: 'inbox', label: 'Inbox Claim Life', kelompok: 'Claim Life', pemilik: 'claimlife' },
  { modul: 'register', label: 'Register', kelompok: 'Claim Life', pemilik: 'claimlife' },
  { modul: 'tco-tahun', label: 'Treaty Contract Out', kelompok: 'Treaty Contract Out', pemilik: 'treatycontractout', datar: true },
  { modul: 'rute-tanpa-baris', label: 'Tanpa Baris', kelompok: 'Claim Life', pemilik: 'claimlife' },
]

const TABEL: MenuTabel = {
  golongan: [
    {
      kode: 'FACULTATIVE',
      kelompok: [{ kode: 'nbfacin', label: 'NB FacIn', modul: 'nbfacin', dimigrasi: false, butir: [] }],
    },
    {
      kode: 'KLAIM',
      kelompok: [
        { kode: 'claimfacin', label: 'Claim Fac In', modul: 'claimfacin', dimigrasi: false, butir: [] },
        {
          kode: 'claimlife',
          label: 'Claim Life',
          modul: 'claimlife',
          dimigrasi: true,
          butir: [
            { kode: 'inbox', label: 'Inbox Claim Life', modul: 'claimlife' },
            { kode: 'baris-tanpa-rute', label: 'Tanpa Rute', modul: 'claimlife' },
            { kode: 'register', label: 'Register', modul: 'claimlife' },
          ],
        },
        // Modul dimigrasi yang butirnya tidak dikirim (MODUL_AKTIF) - hilang.
        { kode: 'komiteclaimlife', label: 'Komite Claim Life', modul: 'komiteclaimlife', dimigrasi: true, butir: [] },
      ],
    },
    {
      kode: 'MASTER',
      kelompok: [
        {
          kode: 'treatycontractout',
          label: 'Treaty Contract Out',
          modul: 'treatycontractout',
          dimigrasi: true,
          butir: [{ kode: 'tco-tahun', label: 'Treaty Contract Out', modul: 'treatycontractout' }],
        },
      ],
    },
  ],
}

describe('susunMenu: pohon tabel dipotong rute frontend', () => {
  const s = susunMenu(TABEL, RUTE)

  it('golongan dan kelompok berurutan seperti tabel', () => {
    expect(s.golongan.map((g) => g.kode)).toEqual(['FACULTATIVE', 'KLAIM', 'MASTER'])
    expect(s.golongan[1]?.kelompok.map((k) => k.nama)).toEqual(['Claim Fac In', 'Claim Life'])
  })

  it('baris tabel tanpa rute frontend tidak tampil dan dicatat', () => {
    const claimLife = s.golongan[1]?.kelompok[1]
    expect(claimLife?.butir.map((b) => b.halaman)).toEqual(['inbox', 'register'])
    expect(s.tanpaRute).toEqual(['claimlife/baris-tanpa-rute'])
  })

  it('butir frontend tanpa baris tabel tidak tampil', () => {
    const semua = s.golongan.flatMap((g) => g.kelompok.flatMap((k) => k.butir.map((b) => b.halaman)))
    expect(semua).not.toContain('rute-tanpa-baris')
    expect(s.entri.map((e) => e.modul)).not.toContain('rute-tanpa-baris')
  })

  it('kelompok belum dimigrasi tetap berdiri; kelompok dimigrasi tanpa butir hilang', () => {
    const kelompok = s.golongan.flatMap((g) => g.kelompok)
    expect(kelompok.filter((k) => !k.dimigrasi).map((k) => k.kode)).toEqual(['nbfacin', 'claimfacin'])
    expect(kelompok.map((k) => k.kode)).not.toContain('komiteclaimlife')
  })

  it('golongan tanpa kelompok tampil hilang', () => {
    const t = susunMenu(
      { golongan: [{ kode: 'TREATY', kelompok: [{ kode: 'x', label: 'X', modul: 'x', dimigrasi: true, butir: [] }] }] },
      RUTE,
    )
    expect(t.golongan).toEqual([])
  })

  it('LABEL dari tabel; penanda datar dari rute frontend', () => {
    const tco = s.golongan[2]?.kelompok[0]?.butir[0]
    expect(tco).toEqual({ halaman: 'tco-tahun', label: 'Treaty Contract Out', pemilik: 'treatycontractout', datar: true })
    expect(s.golongan[1]?.kelompok[1]?.butir.some((b) => 'datar' in b)).toBe(false)
  })

  it('palet: Beranda lalu butir yang tampil, urutan sidebar', () => {
    expect(s.entri.map((e) => [e.modul, e.kelompok])).toEqual([
      ['beranda', 'Beranda'],
      ['inbox', 'Claim Life'],
      ['register', 'Claim Life'],
      ['tco-tahun', 'Treaty Contract Out'],
    ])
  })
})

describe('bentukMenuTabel', () => {
  it('menerima bentuk GET /api/menu', () => {
    expect(bentukMenuTabel(TABEL)).toBe(true)
    expect(bentukMenuTabel({ golongan: [] })).toBe(true)
  })
  it('menolak bentuk lain - tidak pernah dijadikan menu kosong', () => {
    for (const x of [null, {}, { golongan: {} }, { golongan: [{ kode: 'KLAIM' }] },
      { golongan: [{ kode: 'KLAIM', kelompok: [{ kode: 'a', label: 'A', modul: 'a', dimigrasi: '1', butir: [] }] }] },
      { golongan: [{ kode: 'KLAIM', kelompok: [{ kode: 'a', label: 'A', modul: 'a', dimigrasi: true, butir: [{ kode: 1 }] }] }] }]) {
      expect(bentukMenuTabel(x), JSON.stringify(x)).toBe(false)
    }
  })
})
