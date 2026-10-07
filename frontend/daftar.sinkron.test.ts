import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { AKAR_APLIKASI, folderKorpusBelumDimigrasi } from '../inti/frontend/uji/sumber'
import { menuTabelDariMigrasi } from '../inti/frontend/uji/menuBersih'
import { FOLDER_KORPUS, LABEL_MENU } from './katalogKorpus'
import { daftarPalet, saringPalet, susunMenu } from '../inti/frontend/lib/daftarMenu'
import { ENTRI_MENU, MODUL_FRONTEND } from './daftar'

// Sinkron sidebar ↔ palet — DUA ARAH, butir bg.
//
// ⛔ Palet tidak membaca DOM sidebar: ia membaca `MenuTersusun.entri`, yang
// `susunMenu` isi dari tombol yang SAMA dengan sidebar. Berkas ini menagih
// KEDUA arah — entri yang ada di palet tetapi tidak di sidebar dapat dibuka
// lewat Ctrl+K walau menunya tidak terlihat, dan itu persis cacat yang
// REFERENSI_UI bayar sekali.
//
// Menu datar (keputusan work owner 30-09-2026): palet = Beranda + satu entri
// per modul yang tampil dan dimigrasi, urutan sama dengan sidebar.

const SHELL = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'components', 'Shell.tsx'), 'utf8')
const DAFTAR_MENU = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'lib', 'daftarMenu.ts'), 'utf8')
const APP = readFileSync(join(__dirname, 'App.tsx'), 'utf8')

const PALET = susunMenu(menuTabelDariMigrasi(null), ENTRI_MENU)
const DAFTAR = daftarPalet(PALET.entri)

describe('sidebar ↔ palet, dua arah', () => {
  it('palet = Beranda + satu entri per tombol modul yang dapat dibuka, urutan sidebar', () => {
    const dariSidebar = PALET.golongan.flatMap((g) => g.modul.flatMap((t) => (t.halaman === null ? [] : [[t.halaman, t.label]])))
    expect(PALET.entri.map((e) => [e.modul, e.label])).toEqual([['beranda', 'Home'], ...dariSidebar])
    expect(dariSidebar).toHaveLength(MODUL_FRONTEND.length)
    // Label entri palet = LABEL tabel = nama folder korpus.
    const sah = new Set<string>(Object.values(LABEL_MENU))
    for (const e of PALET.entri.slice(1)) expect(sah.has(e.label), e.label).toBe(true)
  })

  it('setiap tombol modul yang Shell render ada di palet', () => {
    // Shell menurunkan butirnya DARI `ENTRI_MENU`, jadi arah ini terjaga
    // oleh konstruksi. Uji ini menahannya tetap begitu: bila seseorang
    // menuliskan butir literal di Shell, `butirKelompok` tidak lagi
    // satu-satunya sumber dan penyimpangan menjadi mungkin lagi.
    // ⛔ Refactor bentuk B paket 7: `ENTRI_MENU` kini tiba lewat prop `menu`
    // (App -> Shell -> PaletMenu).
    // ⛔ Sejak menu dari tabel M_NAV_MENU (30-09-2026): Shell memotong menu
    // `GET /api/menu` dengan `menu` (`susunMenu`), lalu SIDEBAR merender
    // `tersusun.golongan` dan PALET menerima `tersusun.entri` - dua keluaran
    // SATU pemotongan, dan `susunMenu` mengisi `entri` dari tombol yang sama
    // yang ia masukkan ke golongan.
    expect(APP).toContain('menu={ENTRI_MENU}')
    expect(SHELL).toContain('susunMenu(menuTabel.menu, menu)')
    expect(SHELL).toContain('tersusun?.golongan.map((g) =>')
    expect(SHELL).toContain('menu={tersusun?.entri ?? berandaSaja}')
    const susun = DAFTAR_MENU.slice(DAFTAR_MENU.indexOf('export function susunMenu'))
    const masuk = 'modul.push({ kode: m.kode, label: m.label, halaman: r.modul, halamanModul: r.halamanModul ?? [r.modul] })'
    expect(susun).toContain(masuk)
    expect(susun.slice(susun.indexOf(masuk), susun.indexOf(masuk) + 400)).toContain('hasil.entri.push(')
  })
})

describe('saringPalet', () => {
  it('mencocokkan POTONGAN, bukan awalan', () => {
    // Pemakai mengingat kata TENGAH sesering kata depan.
    const hasil = saringPalet(DAFTAR, 'contract')
    expect(hasil.map((h) => h.label)).toContain(FOLDER_KORPUS.treatyContractOut)
  })

  it('setiap KATA harus cocok, urutannya bebas', () => {
    const hasil = saringPalet(DAFTAR, 'life komite')
    expect(hasil.map((h) => h.label)).toEqual([FOLDER_KORPUS.komiteClaimLife])
    // Dan kata yang tidak ada menolak barisnya.
    expect(saringPalet(DAFTAR, 'claim borderaux')).toHaveLength(0)
  })

  it('golongan ikut dicari', () => {
    // `KLAIM` bukan kata di label modulnya; ia golongannya.
    const hasil = saringPalet(DAFTAR, 'klaim')
    expect(hasil.map((h) => h.label).sort()).toEqual([FOLDER_KORPUS.claimLife, FOLDER_KORPUS.komiteClaimLife].sort())
  })

  it('kueri kosong mengembalikan seluruhnya, urutan sidebar', () => {
    const hasil = saringPalet(DAFTAR, '   ')
    expect(hasil).toHaveLength(PALET.entri.length)
    expect(hasil[0]?.label).toBe('Home')
  })

  it('huruf besar-kecil diabaikan', () => {
    expect(saringPalet(DAFTAR, 'PREMIUMLIST').length).toBeGreaterThan(0)
  })
})

describe('nol menu dikarang', () => {
  it('palet TIDAK memuat modul yang belum dimigrasi', () => {
    // ⛔ Modul yang belum dimigrasi berdiri di sidebar sebagai tombol NONAKTIF.
    // Bila salah satunya muncul di palet, ia dapat dibuka lewat Ctrl+K — layar
    // yang tidak ada. Daftarnya pernyataan `Status` di MODUL.md setiap modul
    // (dulu angka 16 di sini, yang harus disunting setiap modul baru).
    const dapatDibuka = new Set(PALET.entri.map((e) => e.label))
    const tanpaButir = Object.values(LABEL_MENU).filter((n) => !dapatDibuka.has(n))
    expect([...tanpaButir].sort()).toEqual(folderKorpusBelumDimigrasi())
    expect(tanpaButir.length).toBeGreaterThan(0)
    for (const nama of tanpaButir) {
      // ⛔ DIPERSEMPIT KE MAKSUDNYA 28-09-2026 (sesi Treaty Contract Out).
      // Dulu: hasil pencarian nama kelompok harus KOSONG. Pencocokan palet
      // memakai POTONGAN kata, sehingga kueri `NB Treaty In` menemukan
      // butir `InboxTreatyContract` (`nb` dan `in` ada di dalam "inbox",
      // `treaty` di kelompoknya) - butir yang SAH, milik kelompok yang
      // sudah dimigrasi. Yang dijaga tetap sama: tidak satu pun hasil
      // BERKELOMPOK modul yang belum dimigrasi.
      for (const hasil of saringPalet(DAFTAR, nama)) {
        expect(hasil.label).not.toBe(nama)
      }
    }
  })
})
