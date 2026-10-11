import { describe, expect, it } from 'vitest'

import type { LayarModul, MenuModul, RuteModul } from '../inti/frontend/modul'
import { rakitModulFrontend, type BerkasLayar, type BerkasMenu, type BerkasRute, type Halaman } from './daftar'

// Perakit daftar modul frontend - struktur tim satu folder per modul
// (30-09-2026). Daftarnya terbentuk dari FOLDER (`import.meta.glob`), jadi yang
// dijaga di sini adalah pasangan `menu.ts` + `rute.tsx` per folder - dengan
// modul tiruan, tanpa bergantung pada modul yang kebetulan ada.

const ruteKosong: RuteModul<Halaman> = () => null

function menu(nama: string, halamanAwal = 'awal'): { PENDAFTARAN_MENU: MenuModul } {
  // `MenuModul` bawaan berhalaman `HalamanTerdaftar`; modul tiruan memakai
  // halaman tiruan, jadi bentuknya dilonggarkan di satu tempat ini.
  return { PENDAFTARAN_MENU: { nama, kelompok: `Kelompok ${nama}`, halaman: ['awal'], halamanAwal } as unknown as MenuModul }
}

function layar(nama: string): { PENDAFTARAN_LAYAR: LayarModul } {
  return { PENDAFTARAN_LAYAR: { nama, kelompok: `Kelompok ${nama}`, halaman: ['awal'], halamanAwal: 'awal' } as unknown as LayarModul }
}

describe('rakitModulFrontend', () => {
  it('memasangkan menu.ts dan rute.tsx per folder, berurutan menurut nama folder', () => {
    const m: BerkasMenu = { '../modul/beta/frontend/menu.ts': menu('beta'), '../modul/alfa/frontend/menu.ts': menu('alfa') }
    const r: BerkasRute = {
      '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong },
      '../modul/beta/frontend/rute.tsx': { RUTE_MODUL: ruteKosong },
    }
    const hasil = rakitModulFrontend(m, r)
    expect(hasil.map((x) => x.nama)).toEqual(['alfa', 'beta'])
    expect(hasil[0]?.kelompok).toBe('Kelompok alfa')
    expect(hasil[0]?.halamanAwal).toBe('awal')
    expect(hasil[0]?.Rute).toBe(ruteKosong)
  })

  it('penghitung kotak masuk Beranda (antreanBeranda) opsional diteruskan apa adanya (06-10-2026)', () => {
    const hitung = async () => [{ workbasket: 'UJI-WB', nama: 'UJI WB', jumlah: 2 }]
    const berantrean = menu('alfa')
    berantrean.PENDAFTARAN_MENU = { ...berantrean.PENDAFTARAN_MENU, antreanBeranda: hitung }
    const hasil = rakitModulFrontend(
      { '../modul/alfa/frontend/menu.ts': berantrean, '../modul/beta/frontend/menu.ts': menu('beta') },
      { '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong }, '../modul/beta/frontend/rute.tsx': { RUTE_MODUL: ruteKosong } },
    )
    expect(hasil[0]?.antreanBeranda).toBe(hitung)
    expect(hasil[1]).not.toHaveProperty('antreanBeranda')
  })

  it('penyedia daftar kotak masuk Beranda (daftarBeranda) opsional diteruskan apa adanya (06-10-2026)', () => {
    const daftar = async () => ({ kolom: [], baris: [] })
    const berdaftar = menu('alfa')
    berdaftar.PENDAFTARAN_MENU = { ...berdaftar.PENDAFTARAN_MENU, daftarBeranda: daftar }
    const hasil = rakitModulFrontend(
      { '../modul/alfa/frontend/menu.ts': berdaftar, '../modul/beta/frontend/menu.ts': menu('beta') },
      { '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong }, '../modul/beta/frontend/rute.tsx': { RUTE_MODUL: ruteKosong } },
    )
    expect(hasil[0]?.daftarBeranda).toBe(daftar)
    expect(hasil[1]).not.toHaveProperty('daftarBeranda')
  })

  it('HALAMAN_AWAL yang bukan halaman modul itu DITOLAK (menu datar 30-09-2026)', () => {
    expect(() =>
      rakitModulFrontend(
        { '../modul/alfa/frontend/menu.ts': menu('alfa', 'lain') },
        { '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong } },
      ),
    ).toThrow('halaman awal "lain" bukan halaman modul itu')
  })

  it('folder yang hanya punya salah satu berkasnya DITOLAK dengan menyebut foldernya', () => {
    expect(() => rakitModulFrontend({ '../modul/alfa/frontend/menu.ts': menu('alfa') }, {})).toThrow(
      'modul/alfa/frontend punya menu.ts tanpa rute.tsx',
    )
    expect(() => rakitModulFrontend({}, { '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong } })).toThrow(
      'modul/alfa/frontend punya rute.tsx tanpa menu.ts',
    )
  })

  it('nama modul yang berbeda dari nama foldernya DITOLAK', () => {
    expect(() =>
      rakitModulFrontend(
        { '../modul/alfa/frontend/menu.ts': menu('lain') },
        { '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong } },
      ),
    ).toThrow('menyebut nama "lain"')
  })

  // Perintah work owner 09-10-2026: menu komite dihapus, "anggap menu itu tidak pernah ada" - komite digabung ke menu
  // klaim masing-masing. Modul TANPA MENU memakai `layar.ts` (PENDAFTARAN_LAYAR) di tempat `menu.ts`.
  it('layar.ts + rute.tsx = modul tanpa menu (tanpaMenu), berurutan bersama modul bermenu', () => {
    const r: BerkasRute = {
      '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong },
      '../modul/beta/frontend/rute.tsx': { RUTE_MODUL: ruteKosong },
    }
    const l: BerkasLayar = { '../modul/beta/frontend/layar.ts': layar('beta') }
    const hasil = rakitModulFrontend({ '../modul/alfa/frontend/menu.ts': menu('alfa') }, r, l)
    expect(hasil.map((x) => [x.nama, x.tanpaMenu === true])).toEqual([
      ['alfa', false],
      ['beta', true],
    ])
    expect(hasil[1]?.halamanAwal).toBe('awal')
    expect(hasil[1]?.Rute).toBe(ruteKosong)
  })

  it('menu.ts DAN layar.ts di satu folder DITOLAK; layar.ts tanpa rute.tsx DITOLAK', () => {
    expect(() =>
      rakitModulFrontend(
        { '../modul/alfa/frontend/menu.ts': menu('alfa') },
        { '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong } },
        { '../modul/alfa/frontend/layar.ts': layar('alfa') },
      ),
    ).toThrow('modul/alfa/frontend punya menu.ts dan layar.ts')
    expect(() => rakitModulFrontend({}, {}, { '../modul/alfa/frontend/layar.ts': layar('alfa') })).toThrow(
      'modul/alfa/frontend punya layar.ts tanpa rute.tsx',
    )
    expect(() => rakitModulFrontend({}, { '../modul/alfa/frontend/rute.tsx': { RUTE_MODUL: ruteKosong } }, {})).toThrow(
      'modul/alfa/frontend punya rute.tsx tanpa menu.ts',
    )
  })
})
