import { describe, expect, it } from 'vitest'

import { bentukMenuTabel, daftarPalet, entriAplikasi, HALAMAN_BERANDA, susunMenu, type EntriMenu, type MenuTabel } from './daftarMenu'

// Menu DATAR dari tabel M_NAV_MENU (keputusan work owner 30-09-2026): menu
// `GET /api/menu` (golongan → modul) DIPOTONG dengan modul frontend yang
// benar-benar terdaftar. Modul di sini tiruan - `inti/` tidak mengenal modul;
// uji dua arah dengan `frontend/daftar.ts` ada di
// `frontend/daftar.menuTabel.test.ts`.
//
// Menggantikan uji pohon golongan → kelompok → butir (label butir, penanda
// `datar`, kelompok terlipat).

type H = 'beranda' | 'inbox' | 'register' | 'outstanding' | 'tco-tahun' | 'lain'

const RUTE: readonly EntriMenu<H>[] = [
  { modul: HALAMAN_BERANDA, label: 'Beranda', kelompok: 'Beranda', pemilik: null },
  { modul: 'inbox', label: 'Claim Life', kelompok: 'Claim Life', pemilik: 'claimlife', halamanModul: ['inbox', 'register', 'outstanding'] },
  { modul: 'tco-tahun', label: 'Treaty Contract Out', kelompok: 'Treaty Contract Out', pemilik: 'treatycontractout', halamanModul: ['tco-tahun'] },
  // Modul frontend terdaftar TANPA baris tabel - tidak tampil.
  { modul: 'lain', label: 'Lain', kelompok: 'Lain', pemilik: 'modullain', halamanModul: ['lain'] },
]

const TABEL: MenuTabel = {
  golongan: [
    {
      kode: 'FACULTATIVE',
      modul: [
        { kode: 'nbfacin', label: 'NB FacIn', modul: 'nbfacin', urutan: 1, dimigrasi: false },
        // DIMIGRASI '0' yang (keliru) terdaftar di frontend: tetap nonaktif, dicatat.
        { kode: 'treatycontractout', label: 'Treaty Contract Out', modul: 'treatycontractout', urutan: 2, dimigrasi: false },
      ],
    },
    {
      kode: 'KLAIM',
      modul: [
        { kode: 'claimfacin', label: 'Claim Fac In', modul: 'claimfacin', urutan: 1, dimigrasi: false },
        { kode: 'claimlife', label: 'Claim Life', modul: 'claimlife', urutan: 2, dimigrasi: true },
        // Dimigrasi TANPA modul frontend - tidak tampil, dicatat.
        { kode: 'komiteclaimlife', label: 'Komite Claim Life', modul: 'komiteclaimlife', urutan: 6, dimigrasi: true },
      ],
    },
    { kode: 'MASTER', modul: [{ kode: 'tanpamodul', label: 'Tanpa Modul', modul: 'tanpamodul', urutan: 1, dimigrasi: true }] },
  ],
}

describe('susunMenu: satu tombol per modul di bawah GROUPMENU', () => {
  const s = susunMenu(TABEL, RUTE)

  it('golongan dan modul berurutan seperti tabel; golongan tanpa tombol hilang', () => {
    expect(s.golongan.map((g) => g.kode)).toEqual(['FACULTATIVE', 'KLAIM'])
    expect(s.golongan[1]?.modul.map((m) => m.label)).toEqual(['Claim Fac In', 'Claim Life'])
  })

  it('tombol modul dimigrasi membuka HALAMAN AWAL-nya, label dari tabel', () => {
    const cl = s.golongan[1]?.modul[1]
    expect(cl).toEqual({ kode: 'claimlife', label: 'Claim Life', halaman: 'inbox', halamanModul: ['inbox', 'register', 'outstanding'] })
  })

  it("DIMIGRASI '0': tombol nonaktif tanpa halaman - walau modul frontend terdaftar (dicatat)", () => {
    const nonaktif = s.golongan.flatMap((g) => g.modul.filter((m) => m.halaman === null)).map((m) => m.kode)
    expect(nonaktif).toEqual(['nbfacin', 'treatycontractout', 'claimfacin'])
    expect(s.nonaktifBerute).toEqual(['treatycontractout'])
  })

  it('baris dimigrasi tanpa modul frontend tidak tampil dan dicatat; modul frontend tanpa baris tidak tampil', () => {
    const kode = s.golongan.flatMap((g) => g.modul.map((m) => m.kode))
    expect(kode).not.toContain('komiteclaimlife')
    expect(kode).not.toContain('modullain')
    expect(s.tanpaRute).toEqual(['komiteclaimlife', 'tanpamodul'])
  })

  it('palet: Beranda lalu setiap tombol yang dapat dibuka, urutan sidebar, golongan sebagai konteks', () => {
    // `kelompok` SATU makna - nama modul; golongannya dibawa terpisah.
    expect(s.entri.map((e) => [e.modul, e.label, e.kelompok, e.golongan])).toEqual([
      ['beranda', 'Beranda', 'Beranda', undefined],
      ['inbox', 'Claim Life', 'Claim Life', 'KLAIM'],
    ])
    expect(daftarPalet(s.entri).map((h) => h.kelompok)).toEqual(['Beranda', 'KLAIM'])
  })

  it('entriAplikasi: hanya milik aplikasi (Beranda)', () => {
    expect(entriAplikasi(RUTE).map((e) => e.modul)).toEqual(['beranda'])
  })
})

describe('bentukMenuTabel', () => {
  it('menerima bentuk GET /api/menu', () => {
    expect(bentukMenuTabel(TABEL)).toBe(true)
    expect(bentukMenuTabel({ golongan: [] })).toBe(true)
  })
  it('menolak bentuk lain - tidak pernah dijadikan menu kosong', () => {
    for (const x of [
      null,
      {},
      { golongan: {} },
      { golongan: [{ kode: 'KLAIM' }] },
      // Pohon lama (golongan → kelompok → butir) DITOLAK.
      { golongan: [{ kode: 'KLAIM', kelompok: [{ kode: 'a', label: 'A', modul: 'a', dimigrasi: true, butir: [] }] }] },
      { golongan: [{ kode: 'KLAIM', modul: [{ kode: 'a', label: 'A', modul: 'a', urutan: 1, dimigrasi: '1' }] }] },
      { golongan: [{ kode: 'KLAIM', modul: [{ kode: 'a', label: 'A', modul: 'a', urutan: '1', dimigrasi: true }] }] },
    ]) {
      expect(bentukMenuTabel(x), JSON.stringify(x)).toBe(false)
    }
  })
})
