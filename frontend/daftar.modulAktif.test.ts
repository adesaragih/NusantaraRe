import { existsSync, readdirSync, readFileSync } from 'node:fs'
import { join } from 'node:path'

import { afterEach, describe, expect, it, vi } from 'vitest'

import { AKAR_APLIKASI, folderKorpusBelumDimigrasi } from '../inti/frontend/uji/sumber'
import { menuTabelDariMigrasi } from '../inti/frontend/uji/menuBersih'
import { FOLDER_KORPUS, FOLDER_TANPA_MENU } from './katalogKorpus'
import { modulBerkotakMasuk } from './Beranda'
import { ambilModulAktif } from '../inti/frontend/klien'
import { daftarPalet, KODE_MENU_KELOLA_USER, KODE_MENU_TEMPLATE_MANAGER, susunMenu } from '../inti/frontend/lib/daftarMenu'
import { ENTRI_MENU, halamanAktif, MODUL_BACKEND, MODUL_FRONTEND } from './daftar'

// MODUL_AKTIF di frontend - refactor bentuk B paket 6.
//
// ⛔ Yang dijaga: menu modul yang TIDAK dipasang backend hilang dari sidebar
// DAN palet (dua arah tetap sinkron), `null` (daftar belum/gagal terbaca)
// menampilkan SEMUA seperti sebelum MODUL_AKTIF ada, dan nama modul di
// frontend adalah nama yang sama dengan `const Nama` di `modul/*/modul.go`.

const SRC = __dirname
const AKAR_MODUL = join(AKAR_APLIKASI, 'modul')

describe('menu modul nonaktif hilang', () => {
  // Saringannya di BACKEND (GET /api/menu tidak mengirim modul dimigrasi di
  // luar MODUL_AKTIF - `inti/backend/menu/menu_test.go`); menu tabel di sini
  // disusun dari HASIL BERSIH migrasi menu seperti backend menyusunnya
  // (`menuTabelDariMigrasi`). Yang dijaga: tombol modul nonaktif HILANG, modul
  // yang belum dimigrasi tetap BERDIRI nonaktif, dan palet = sidebar.
  it('sidebar: tombol modul nonaktif hilang, sisanya utuh', () => {
    const s = susunMenu(menuTabelDariMigrasi(['claimlife']), ENTRI_MENU)
    const tombol = s.golongan.flatMap((g) => g.modul)
    for (const lain of [FOLDER_KORPUS.premiumListLife, FOLDER_KORPUS.komiteClaimLife, FOLDER_KORPUS.treatyContractOut]) {
      expect(tombol.map((t) => t.label)).not.toContain(lain)
    }
    expect(tombol.find((t) => t.label === FOLDER_KORPUS.claimLife)?.halaman).toBe('inbox')
    // Modul yang memang belum dimigrasi tetap berdiri - ia bukan modul nonaktif.
    // Daftarnya pernyataan `Status` di MODUL.md setiap modul, bukan angka di sini.
    const belum = tombol.filter((t) => t.halaman === null)
    expect(belum.map((t) => t.label).sort()).toEqual(folderKorpusBelumDimigrasi().filter((f) => !FOLDER_TANPA_MENU.includes(f)))
    expect(s.tanpaRute).toEqual([])
    expect(s.nonaktifBerute).toEqual([])
  })

  it('palet: Beranda dan tombol sidebar yang sama', () => {
    const s = susunMenu(menuTabelDariMigrasi(['premiumlistlife', 'treatycontractout']), ENTRI_MENU)
    const dariSidebar = s.golongan.flatMap((g) => g.modul.flatMap((t) => (t.halaman === null ? [] : [t.halaman])))
    expect(daftarPalet(s.entri).map((h) => h.modul)).toEqual(['beranda', ...dariSidebar])
    expect(dariSidebar).toEqual(['premiumlist', 'tco-tahun'])
  })

  it('Beranda: kotak masuk hanya dari modul aktif', () => {
    // Panel kotak masuk MEMBUKA modul dan berkasnya - ia menu juga. (Tabel modul dan kartu tahap Claim Life
    // dibuang 06-10-2026, perintah work owner "buang aja, ga perlu".)
    expect(modulBerkotakMasuk(MODUL_FRONTEND, ['komiteclaimlife'])).toEqual([])
    expect(modulBerkotakMasuk(MODUL_FRONTEND, null)).toEqual(modulBerkotakMasuk(MODUL_FRONTEND))
    // Sejak Kelola User (01-10-2026) Beranda menerima modul aktif yang menunya
    // DIPEGANG akun (`modulBoleh` = `modulUntukAkun(modulAktif, ...)`).
    // (06-10-2026: elemen <Beranda> kini berbaris banyak - prop onBukaKasus kotak masuk; yang dijaga tetap modulBoleh)
    expect(readFileSync(join(SRC, 'App.tsx'), 'utf8')).toMatch(/<Beranda\s[\s\S]*?modulAktif=\{modulBoleh\}\s*\/>/)
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
      // Menu aplikasi Kelola User dan Template Manager: `pemilik` adalah KODE menunya, bukan modul
      // backend - ia tidak tunduk pada MODUL_AKTIF (`daftar.kelolauser.test.ts`).
      const menuAplikasi = e.pemilik === KODE_MENU_KELOLA_USER || e.pemilik === KODE_MENU_TEMPLATE_MANAGER
      expect(menuAplikasi ? null : e.pemilik).toBe(MODUL_BACKEND[e.modul])
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
    // Modul aktif DISARING menu akun sebelum dipakai (Kelola User 01-10-2026).
    expect(app).toContain('modulUntukAkun(modulAktif, menuAkun, MODUL_FRONTEND.map((m) => m.nama))')
    expect(app).toContain('modulAktif={modulBoleh}')
    // Sejak menu dari tabel (30-09-2026) Shell tidak lagi menyaring modul
    // aktif: GET /api/menu tidak mengirim modul dimigrasi yang nonaktif
    // (cmd/api meneruskan daftar modul aktif ke rute menu).
    const shell = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'components', 'Shell.tsx'), 'utf8')
    expect(shell).not.toContain('modulAktif')
    const rakit = readFileSync(join(AKAR_APLIKASI, 'cmd', 'api', 'rakit.go'), 'utf8')
    expect(rakit).toContain('mux.HandleFunc("GET /api/menu", ruteMenu(dasar, aktif, stubPelaku))')
    // Dan rute modul nonaktif tidak dipasang (paket 7: App merakit dari modul),
    // dan halaman modul yang ternyata nonaktif kembali ke Beranda.
    expect(app).toContain('MODUL_FRONTEND.filter((m) => modulDipasang(m.nama, modulBoleh))')
    expect(app).toContain('if (!halamanAktif(halaman, modulBoleh) || (halaman === HALAMAN_KELOLA_USER && !bolehKelola)) {')
  })
})
