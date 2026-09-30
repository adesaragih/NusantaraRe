import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { AKAR_APLIKASI } from '../inti/frontend/uji/sumber'
import { FOLDER_KORPUS } from './katalogKorpus'
import { kartuModul } from './Beranda'
import { ambilModulAktif } from '../inti/frontend/klien'
import { daftarPalet, susunMenu, type MenuTabel } from '../inti/frontend/lib/daftarMenu'
import { ENTRI_MENU, halamanAktif, MODUL_BACKEND, MODUL_FRONTEND } from './daftar'

// MODUL_AKTIF di frontend - refactor bentuk B paket 6.
//
// ⛔ Yang dijaga: menu modul yang TIDAK dipasang backend hilang dari sidebar
// DAN palet (dua arah tetap sinkron), `null` (daftar belum/gagal terbaca)
// menampilkan SEMUA seperti sebelum MODUL_AKTIF ada, dan nama modul di
// frontend adalah nama yang sama dengan `const Nama` di `modul/*/modul.go`.

const SRC = __dirname
const AKAR_MODUL = join(AKAR_APLIKASI, 'modul')

/**
 * Pohon `GET /api/menu` seperti yang backend kirim untuk `aktif`: satu kelompok
 * per MODUL, butir modul NONAKTIF tidak dikirim (`inti/menu` `Susun`).
 * Golongan tidak diuji di sini - satu golongan tiruan.
 */
function tabel(aktif: readonly string[]): MenuTabel {
  return {
    golongan: [
      {
        kode: 'KLAIM',
        kelompok: Object.values(FOLDER_KORPUS).map((nama) => {
          const milik = ENTRI_MENU.filter((e) => e.kelompok === nama)
          const modul = milik[0]?.pemilik ?? nama.toLowerCase().replace(/ /g, '')
          return {
            kode: modul,
            label: nama,
            modul,
            dimigrasi: milik.length > 0,
            butir: milik
              .filter((e) => e.pemilik !== null && aktif.includes(e.pemilik))
              .map((e) => ({ kode: e.modul, label: e.label, modul })),
          }
        }),
      },
    ],
  }
}

describe('menu modul nonaktif hilang', () => {
  // Saringannya kini di BACKEND (GET /api/menu tidak mengirim butir modul di
  // luar MODUL_AKTIF - `inti/menu/menu_test.go`); yang dijaga di sini adalah
  // sisi frontend-nya: kelompok yang butirnya tidak dikirim HILANG, yang
  // memang belum dimigrasi tetap BERDIRI, dan palet = sidebar.
  it('sidebar: kelompok modul nonaktif hilang, sisanya utuh', () => {
    const s = susunMenu(tabel(['claimlife']), ENTRI_MENU)
    const nama = s.golongan.flatMap((g) => g.kelompok.map((k) => k.nama))
    for (const lain of [FOLDER_KORPUS.premiumListLife, FOLDER_KORPUS.komiteClaimLife, FOLDER_KORPUS.treatyContractOut]) {
      expect(nama).not.toContain(lain)
    }
    const claimLife = s.golongan[0]?.kelompok.find((k) => k.nama === FOLDER_KORPUS.claimLife)
    expect(claimLife?.butir.map((b) => b.halaman)).toEqual(['inbox', 'register'])
    // Kelompok yang memang belum dimigrasi tetap berdiri - ia bukan modul nonaktif.
    const belum = s.golongan.flatMap((g) => g.kelompok.filter((k) => !k.dimigrasi))
    expect(belum).toHaveLength(16)
    expect(s.tanpaRute).toEqual([])
  })

  it('palet: Beranda dan butir sidebar yang sama', () => {
    const s = susunMenu(tabel(['premiumlistlife', 'treatycontractout']), ENTRI_MENU)
    const dariSidebar = s.golongan.flatMap((g) => g.kelompok.flatMap((k) => k.butir.map((b) => b.halaman)))
    expect(daftarPalet(s.entri).map((h) => h.modul)).toEqual(['beranda', ...dariSidebar])
    expect(dariSidebar).toEqual(['premiumlist', 'tco-tahun'])
  })

  it('Beranda: kartu modul nonaktif hilang, yang belum dimigrasi tetap', () => {
    // Tombol kartu Beranda MEMBUKA modul - ia menu juga.
    const kartu = kartuModul(['komiteclaimlife'])
    expect(kartu.filter((k) => k.tujuan !== null).map((k) => k.nama)).toEqual([FOLDER_KORPUS.komiteClaimLife])
    expect(kartu.filter((k) => k.tujuan === null)).toHaveLength(16)
    expect(kartuModul(null)).toEqual(kartuModul())
    // Cacah antrean Claim Life tidak diminta bila modul itu nonaktif.
    const beranda = readFileSync(join(SRC, 'Beranda.tsx'), 'utf8')
    expect(beranda).toContain('const claimLifeAktif = modulDipasang(NAMA_CLAIMLIFE, modulAktif)')
    expect(beranda).toContain('if (!claimLifeAktif) {')
    expect(readFileSync(join(SRC, 'App.tsx'), 'utf8')).toContain('<Beranda masuk={masuk} onBuka={setHalaman} modulAktif={modulAktif} />')
  })

  it('null = semua rute terpasang, persis seperti sebelum MODUL_AKTIF', () => {
    for (const e of ENTRI_MENU) expect(halamanAktif(e.modul, null), e.modul).toBe(true)
    // Beranda milik aplikasi: tampil walau tak satu modul pun disebut.
    expect(halamanAktif('beranda', [])).toBe(true)
  })
})

describe('nama modul sama dengan backend', () => {
  it('MODUL_BACKEND menyebut tepat modul `modul/*/backend/modul.go`', () => {
    // Satu folder per modul (30-09-2026): modul TERDAFTAR = folder yang punya
    // `backend/modul.go`. Folder kerangka (hanya MODUL.md dan docs/) bukan modul
    // terdaftar - ia belum dimigrasi.
    const dariGo = readdirSync(AKAR_MODUL, { withFileTypes: true })
      .filter((d) => d.isDirectory() && existsSync(join(AKAR_MODUL, d.name, 'backend', 'modul.go')))
      .map((d) => {
        const sumber = readFileSync(join(AKAR_MODUL, d.name, 'backend', 'modul.go'), 'utf8')
        const cocok = /^const Nama = "([a-z]+)"$/m.exec(sumber)
        expect(cocok, `modul/${d.name}/backend/modul.go tanpa const Nama`).not.toBeNull()
        // Tabel nama modul (30-09-2026): nama folder backend = `const Nama`.
        expect(cocok?.[1], `modul/${d.name}: folder dan const Nama berbeda`).toBe(d.name)
        return cocok?.[1] ?? ''
      })
    expect(dariGo.length).toBeGreaterThanOrEqual(4)
    // Dan setiap modul terdaftar punya `frontend/` di foldernya sendiri - satu
    // nama untuk backend, frontend, dan dokumen.
    const folderFrontend = readdirSync(AKAR_MODUL, { withFileTypes: true })
      .filter((d) => d.isDirectory() && existsSync(join(AKAR_MODUL, d.name, 'frontend')))
      .map((d) => d.name)
    expect(new Set(folderFrontend)).toEqual(new Set(dariGo))
    const dariFrontend = Object.values(MODUL_BACKEND).filter((m): m is string => m !== null)
    expect(new Set(dariFrontend)).toEqual(new Set(dariGo))
    // Daftar modul frontend = daftar modul backend (refactor bentuk B paket 7).
    expect(new Set(MODUL_FRONTEND.map((m) => m.nama))).toEqual(new Set(dariGo))
    // ⛔ Nama halaman UNIK lintas modul: `MODUL_BACKEND` dirakit dengan
    // `Object.fromEntries`, dan halaman kembar diam-diam dimiliki modul yang
    // terdaftar belakangan - dua rute merender halaman yang sama (/code-review).
    const semuaHalaman = MODUL_FRONTEND.flatMap((m) => [...m.halaman])
    expect(new Set(semuaHalaman).size).toBe(semuaHalaman.length)
    expect(semuaHalaman).not.toContain('beranda')
  })

  it('setiap entri menu punya pemilik yang terdaftar', () => {
    for (const e of ENTRI_MENU) {
      expect(Object.keys(MODUL_BACKEND)).toContain(e.modul)
      expect(e.pemilik).toBe(MODUL_BACKEND[e.modul])
    }
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
    jawab('{"modul":["claimlife","komiteclaimlife"]}')
    await expect(ambilModulAktif()).resolves.toEqual(['claimlife', 'komiteclaimlife'])
    const panggil = vi.mocked(fetch).mock.calls[0]
    expect(String(panggil?.[0])).toContain('/api/modul-aktif')
    const rakit = readFileSync(join(AKAR_APLIKASI, 'cmd', 'api', 'rakit.go'), 'utf8')
    expect(rakit).toContain('mux.HandleFunc("GET /api/modul-aktif"')
  })

  it('badan yang tak dikenal menjadi null (semua tampil), bukan daftar kosong', async () => {
    for (const badan of ['{}', '{"modul":"claimlife"}', '{"modul":[1]}', '']) {
      jawab(badan)
      await expect(ambilModulAktif()).resolves.toBeNull()
    }
  })

  it('App membacanya untuk rute dan Beranda; menunya disaring backend', () => {
    const app = readFileSync(join(SRC, 'App.tsx'), 'utf8')
    expect(app).toContain('ambilModulAktif()')
    expect(app).toContain('modulAktif={modulAktif}')
    // Sejak menu dari tabel (30-09-2026) Shell tidak lagi menyaring modul
    // aktif: GET /api/menu tidak mengirim butir modul nonaktif (cmd/api
    // meneruskan daftar modul aktif ke rute menu).
    const shell = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'components', 'Shell.tsx'), 'utf8')
    expect(shell).not.toContain('modulAktif')
    const rakit = readFileSync(join(AKAR_APLIKASI, 'cmd', 'api', 'rakit.go'), 'utf8')
    expect(rakit).toContain('mux.HandleFunc("GET /api/menu", ruteMenu(dasar, aktif, stubPelaku))')
    // Dan rute modul nonaktif tidak dipasang (paket 7: App merakit dari modul),
    // dan halaman modul yang ternyata nonaktif kembali ke Beranda.
    expect(app).toContain('MODUL_FRONTEND.filter((m) => modulDipasang(m.nama, modulAktif))')
    expect(app).toContain("if (!halamanAktif(halaman, modulAktif)) setHalaman('beranda')")
  })
})
