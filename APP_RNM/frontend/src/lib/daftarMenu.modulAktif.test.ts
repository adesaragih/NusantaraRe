import { readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { MODUL } from '../assets/labels'
import { ambilModulAktif } from '../services/api'
import {
  daftarPalet,
  ENTRI_MENU,
  halamanAktif,
  kelompokTampil,
  MODUL_BACKEND,
  type KelompokSidebar,
} from './daftarMenu'

// MODUL_AKTIF di frontend - refactor bentuk B paket 6.
//
// ⛔ Yang dijaga: menu modul yang TIDAK dipasang backend hilang dari sidebar
// DAN palet (dua arah tetap sinkron), `null` (daftar belum/gagal terbaca)
// menampilkan SEMUA seperti sebelum MODUL_AKTIF ada, dan nama modul di
// frontend adalah nama yang sama dengan `const Nama` di `modul/*/modul.go`.

const SRC = join(__dirname, '..')
const AKAR_MODUL_GO = join(SRC, '..', '..', 'modul')

/** Kelompok sidebar, diturunkan persis seperti `Shell.tsx` menurunkannya. */
const KELOMPOK: readonly KelompokSidebar[] = Object.values(MODUL).map((nama) => ({
  nama,
  butir: ENTRI_MENU.filter((e) => e.kelompok === nama).map((e) => ({ halaman: e.modul, label: e.label })),
}))

/** Nama kelompok yang memuat butir milik modul backend `nama`. */
function kelompokMilik(nama: string): string[] {
  return KELOMPOK.filter((k) => k.butir.some((b) => MODUL_BACKEND[b.halaman] === nama)).map((k) => k.nama)
}

describe('menu modul nonaktif hilang', () => {
  it('sidebar: kelompok modul nonaktif hilang, sisanya utuh', () => {
    const tampil = kelompokTampil(KELOMPOK, ['claimlife'])
    const nama = tampil.map(({ k }) => k.nama)
    for (const lain of ['premiumlist', 'komite', 'treaty']) {
      const milik = kelompokMilik(lain)
      expect(milik.length).toBeGreaterThan(0)
      for (const n of milik) expect(nama).not.toContain(n)
    }
    const claimLife = tampil.find(({ k }) => k.nama === MODUL.claimLife)
    expect(claimLife?.butir.map((b) => b.halaman)).toEqual(['inbox', 'register'])
    // Kelompok yang memang belum dimigrasi tetap berdiri - ia bukan modul nonaktif.
    const tanpaButir = KELOMPOK.filter((k) => k.butir.length === 0).map((k) => k.nama)
    expect(tanpaButir).toHaveLength(14)
    for (const n of tanpaButir) expect(nama).toContain(n)
  })

  it('palet: hanya Beranda dan menu modul aktif', () => {
    const halaman = daftarPalet(['claimlife']).map((h) => h.modul)
    expect(halaman).toEqual(['beranda', 'inbox', 'register'])
    expect(daftarPalet(['treaty', 'komite']).map((h) => h.modul)).toEqual(
      ENTRI_MENU.filter((e) => ['beranda', 'komite', 'tco-tahun', 'tco-kontrak', 'tco-klausul'].includes(e.modul)).map(
        (e) => e.modul,
      ),
    )
  })

  it('sidebar dan palet menyaring dengan aturan yang sama', () => {
    const aktif = ['premiumlist', 'treaty']
    const dariSidebar = kelompokTampil(KELOMPOK, aktif).flatMap(({ butir }) => butir.map((b) => b.halaman))
    const dariPalet = daftarPalet(aktif)
      .map((h) => h.modul)
      .filter((m) => m !== 'beranda')
    expect(new Set(dariSidebar)).toEqual(new Set(dariPalet))
  })

  it('null = semua tampil, persis seperti sebelum MODUL_AKTIF', () => {
    expect(kelompokTampil(KELOMPOK, null).map(({ k, butir }) => ({ nama: k.nama, butir }))).toEqual(
      KELOMPOK.map((k) => ({ nama: k.nama, butir: k.butir })),
    )
    expect(daftarPalet(null)).toEqual(daftarPalet())
    expect(daftarPalet(null)).toHaveLength(ENTRI_MENU.length)
    // Beranda milik aplikasi: tampil walau tak satu modul pun disebut.
    expect(halamanAktif('beranda', [])).toBe(true)
  })
})

describe('nama modul sama dengan backend', () => {
  it('MODUL_BACKEND menyebut tepat modul `modul/*/modul.go`', () => {
    const dariGo = readdirSync(AKAR_MODUL_GO, { withFileTypes: true })
      .filter((d) => d.isDirectory())
      .map((d) => {
        const sumber = readFileSync(join(AKAR_MODUL_GO, d.name, 'modul.go'), 'utf8')
        const cocok = /^const Nama = "([a-z]+)"$/m.exec(sumber)
        expect(cocok, `modul/${d.name}/modul.go tanpa const Nama`).not.toBeNull()
        return cocok?.[1] ?? ''
      })
    expect(dariGo.length).toBeGreaterThanOrEqual(4)
    const dariFrontend = Object.values(MODUL_BACKEND).filter((m): m is string => m !== null)
    expect(new Set(dariFrontend)).toEqual(new Set(dariGo))
  })

  it('setiap entri menu punya pemilik yang terdaftar', () => {
    for (const e of ENTRI_MENU) expect(Object.keys(MODUL_BACKEND)).toContain(e.modul)
  })
})

describe('GET /api/modul-aktif', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  function jawab(badan: string): void {
    vi.stubGlobal(
      'fetch',
      vi.fn(() => Promise.resolve(new Response(badan, { status: 200 }))),
    )
  }

  it('membaca daftar modul dari jalur yang dipasang cmd/api', async () => {
    jawab('{"modul":["claimlife","komite"]}')
    await expect(ambilModulAktif()).resolves.toEqual(['claimlife', 'komite'])
    const panggil = vi.mocked(fetch).mock.calls[0]
    expect(String(panggil?.[0])).toContain('/api/modul-aktif')
    const rakit = readFileSync(join(SRC, '..', '..', 'cmd', 'api', 'rakit.go'), 'utf8')
    expect(rakit).toContain('mux.HandleFunc("GET /api/modul-aktif"')
  })

  it('badan yang tak dikenal menjadi null (semua tampil), bukan daftar kosong', async () => {
    for (const badan of ['{}', '{"modul":"claimlife"}', '{"modul":[1]}', '']) {
      jawab(badan)
      await expect(ambilModulAktif()).resolves.toBeNull()
    }
  })

  it('App membacanya dan Shell meneruskannya ke sidebar dan palet', () => {
    const app = readFileSync(join(SRC, 'App.tsx'), 'utf8')
    expect(app).toContain('ambilModulAktif()')
    expect(app).toContain('modulAktif={modulAktif}')
    const shell = readFileSync(join(SRC, 'components', 'Shell.tsx'), 'utf8')
    expect(shell).toContain('kelompokTampil(KELOMPOK, modulAktif)')
    expect(shell).toContain('modulAktif={modulAktif}')
    const palet = readFileSync(join(SRC, 'components', 'PaletMenu.tsx'), 'utf8')
    expect(palet).toContain('daftarPalet(modulAktif)')
  })
})
