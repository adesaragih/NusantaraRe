import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { AKAR_APLIKASI, folderKorpusBelumDimigrasi } from '../inti/frontend/uji/sumber'
import { MENU } from '../inti/frontend/labels'
import { LABEL_MENU_KOMITE } from '../modul/komiteclaimlife/frontend/labels'
import { LABEL_MENU_PREMIUMLIST } from '../modul/premiumlistlife/frontend/labels'
import { FOLDER_KORPUS } from './katalogKorpus'
import { daftarPalet, saringPalet } from '../inti/frontend/lib/daftarMenu'
import { ENTRI_MENU } from './daftar'

// Sinkron sidebar ↔ palet — DUA ARAH, butir bg.
//
// ⛔ HARGANYA DITULIS APA ADANYA. Palet tidak boleh membaca DOM sidebar:
// `KelompokMenu` melepas anak kelompok yang TERLIPAT dari DOM, dan sejak bg
// empat belas kelompok memang terlipat — palet yang membaca DOM akan
// menjawab "tidak ada" untuk menu yang ADA.
//
// Harga dari pilihan itu: dua daftar yang dapat menyimpang. Berkas ini yang
// membayarnya, dan ia menagih KEDUA arah — entri yang ada di palet tetapi
// tidak di sidebar dapat dibuka lewat Ctrl+K walau menunya tidak terlihat,
// dan itu persis cacat yang REFERENSI_UI bayar sekali.

const SHELL = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'components', 'Shell.tsx'), 'utf8')
const DAFTAR_MENU = readFileSync(join(AKAR_APLIKASI, 'inti', 'frontend', 'lib', 'daftarMenu.ts'), 'utf8')
const APP = readFileSync(join(__dirname, 'App.tsx'), 'utf8')

describe('sidebar ↔ palet, dua arah', () => {
  it('setiap entri palet punya kelompok yang Shell render', () => {
    const kelompokSah = new Set<string>(Object.values(FOLDER_KORPUS))
    kelompokSah.add('Beranda')
    for (const e of ENTRI_MENU) {
      expect(kelompokSah.has(e.kelompok)).toBe(true)
    }
  })

  it('setiap butir modul yang Shell render ada di palet', () => {
    // Shell menurunkan butirnya DARI `ENTRI_MENU`, jadi arah ini terjaga
    // oleh konstruksi. Uji ini menahannya tetap begitu: bila seseorang
    // menuliskan butir literal di Shell, `butirKelompok` tidak lagi
    // satu-satunya sumber dan penyimpangan menjadi mungkin lagi.
    // ⛔ Refactor bentuk B paket 7: `ENTRI_MENU` kini tiba lewat prop `menu`
    // (App -> Shell -> PaletMenu).
    // ⛔ Sejak menu dari tabel M_NAV_MENU (30-09-2026): Shell memotong pohon
    // `GET /api/menu` dengan `menu` (`susunMenu`), lalu SIDEBAR merender
    // `tersusun.golongan` dan PALET menerima `tersusun.entri` - dua keluaran
    // SATU pemotongan, dan `susunMenu` mengisi `entri` dari butir yang
    // sama yang ia masukkan ke kelompok.
    expect(APP).toContain('menu={ENTRI_MENU}')
    expect(SHELL).toContain('susunMenu(menuTabel.menu, menu)')
    expect(SHELL).toContain('tersusun?.golongan.map((g) =>')
    expect(SHELL).toContain('menu={tersusun?.entri ?? berandaSaja}')
    const susun = DAFTAR_MENU.slice(DAFTAR_MENU.indexOf('export function susunMenu'))
    const masuk = "kelompok.push({ kode: k.kode, nama: k.label, dimigrasi: true, butir })"
    expect(susun).toContain(masuk)
    expect(susun.slice(susun.indexOf(masuk))).toContain('of butir) hasil.entri.push(')
  })

  it('label palet datang dari labels.ts, tidak diketik ulang', () => {
    const label = ENTRI_MENU.map((e) => e.label)
    expect(label).toContain(MENU.inbox)
    expect(label).toContain(MENU.register)
    expect(label).toContain(LABEL_MENU_PREMIUMLIST.premiumList)
    expect(label).toContain(LABEL_MENU_KOMITE.inboxKomite)
  })
})

describe('saringPalet', () => {
  it('mencocokkan POTONGAN, bukan awalan', () => {
    // Pemakai mengingat kata TENGAH sesering kata depan.
    const hasil = saringPalet(daftarPalet(ENTRI_MENU), 'claim')
    expect(hasil.map((h) => h.label)).toContain(MENU.inbox)
  })

  it('setiap KATA harus cocok, urutannya bebas', () => {
    const hasil = saringPalet(daftarPalet(ENTRI_MENU), 'life claim')
    expect(hasil.map((h) => h.label)).toContain(MENU.inbox)
    // Dan kata yang tidak ada menolak barisnya.
    expect(saringPalet(daftarPalet(ENTRI_MENU), 'claim borderaux')).toHaveLength(0)
  })

  it('kelompok ikut dicari', () => {
    // `PremiumList` label butirnya; `PremiumList Life` nama kelompoknya.
    const hasil = saringPalet(daftarPalet(ENTRI_MENU), 'premiumlist life')
    expect(hasil.map((h) => h.label)).toContain(LABEL_MENU_PREMIUMLIST.premiumList)
  })

  it('kueri kosong mengembalikan seluruhnya, urutan sidebar', () => {
    const hasil = saringPalet(daftarPalet(ENTRI_MENU), '   ')
    expect(hasil).toHaveLength(ENTRI_MENU.length)
    expect(hasil[0]?.label).toBe('Beranda')
  })

  it('huruf besar-kecil diabaikan', () => {
    expect(saringPalet(daftarPalet(ENTRI_MENU), 'INBOX').length).toBeGreaterThan(0)
  })
})

describe('nol menu dikarang', () => {
  it('palet TIDAK memuat modul yang belum dimigrasi', () => {
    // ⛔ Kelompok yang belum dimigrasi berdiri di sidebar TANPA butir. Bila
    // salah satunya muncul di palet, ia dapat dibuka lewat Ctrl+K — layar
    // yang tidak ada. Daftarnya pernyataan `Status` di MODUL.md setiap modul
    // (dulu angka 16 di sini, yang harus disunting setiap modul baru).
    const berbutir = new Set(ENTRI_MENU.map((e) => e.kelompok))
    const tanpaButir = Object.values(FOLDER_KORPUS).filter((n) => !berbutir.has(n))
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
      for (const hasil of saringPalet(daftarPalet(ENTRI_MENU), nama)) {
        expect(hasil.kelompok).not.toBe(nama)
      }
    }
  })
})
