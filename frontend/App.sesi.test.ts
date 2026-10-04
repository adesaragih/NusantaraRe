import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

// Layar login tanpa kedip (laporan work owner 01-10-2026: "halaman login suka
// kedip kedip"). Penjaga statik atas App.tsx - dua sebab yang ditemukan:
//
//   1. Pemeriksaan sesi dulu merender paragraf polos di halaman putih, lalu
//      berganti ke kartu login bergradien: kedip di SETIAP muat ulang. Latar
//      login sebagai gantinya ganti berkedip UNGU sebelum Shell bagi yang sudah
//      login (teramati di Chrome headless, 01-10-2026) - layarnya kini NETRAL
//      dan teksnya tertunda.
//   2. `GET /api/menu` dikirim saat App menyala, SEBELUM login. Tanpa sesi ia
//      dijawab 401, dan 401 memicu `PERISTIWA_SESI_BERAKHIR` di tengah
//      pemeriksaan sesi.
const APP = readFileSync(join(__dirname, 'App.tsx'), 'utf8')
const SHELL = readFileSync(join(__dirname, '..', 'inti', 'frontend', 'components', 'Shell.tsx'), 'utf8')
const CSS = readFileSync(join(__dirname, '..', 'inti', 'frontend', 'styles.css'), 'utf8')

describe('App: pemeriksaan sesi tanpa kedip', () => {
  it('layar periksa sesi netral: tanpa latar login, teks tertunda', () => {
    expect(APP).toContain('<main className="periksa-sesi" aria-busy="true">')
    expect(APP).not.toContain('<main className="halaman-masuk" aria-busy="true">')
    expect(APP).not.toContain('className="polis__catatan"')
    expect(CSS).toMatch(/\.periksa-sesi__status \{[^}]*opacity: 0;[^}]*animation: tampak-tertunda 0\.2s ease 0\.6s forwards;/)
  })

  it('"Memuat menu" sidebar tertunda - menu dibaca sesudah sesi', () => {
    expect(SHELL).toContain('<li className="shell__menu-keadaan shell__menu-keadaan--tunda">')
    expect(CSS).toMatch(/\.shell__menu-keadaan--tunda \{[^}]*animation: tampak-tertunda 0\.2s ease 0\.6s forwards;/)
  })

  it('menu dibaca sesudah identitas diketahui, ulang bila akunnya berganti', () => {
    // Akun yang wajib ganti sandi TIDAK membaca menu (401 -> kembali ke login).
    expect(APP).toContain("const kunciMenu = stub ? 'stub' : profil && !profil.wajibGantiSandi ? profil.akunId : null")
    expect(APP).toContain('if (kunciMenu === null) return')
    expect(APP).toContain('}, [kunciMenu, versiMenu])')
  })

  // Kelola User (01-10-2026): modul dipasang menurut menu akun, dan halaman
  // Kelola User hanya bagi pemegangnya.
  it('modul dan Kelola User dipasang menurut menu akun', () => {
    expect(APP).toContain('MODUL_FRONTEND.filter((m) => modulDipasang(m.nama, modulBoleh))')
    expect(APP).toContain('modulAktif={modulBoleh}')
    expect(APP).toContain('{halaman === HALAMAN_KELOLA_USER && bolehKelola && <KelolaUser')
    expect(APP).not.toContain('modulDipasang(m.nama, modulAktif)')
  })
})
