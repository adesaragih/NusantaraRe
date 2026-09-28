import { readFileSync } from 'node:fs'
import { join } from 'node:path'

import { describe, expect, it } from 'vitest'

import { MENU, MENU_MODUL, MODUL } from '../assets/labels'
import { daftarPalet, ENTRI_MENU, saringPalet } from './daftarMenu'

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

const SHELL = readFileSync(join(__dirname, '..', 'components', 'Shell.tsx'), 'utf8')

describe('sidebar ↔ palet, dua arah', () => {
  it('setiap entri palet punya kelompok yang Shell render', () => {
    const kelompokSah = new Set<string>(Object.values(MODUL))
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
    expect(SHELL).toContain('butirKelompok')
    expect(SHELL).toContain("ENTRI_MENU.filter((e) => e.kelompok === nama)")
  })

  it('label palet datang dari labels.ts, tidak diketik ulang', () => {
    const label = ENTRI_MENU.map((e) => e.label)
    expect(label).toContain(MENU.inbox)
    expect(label).toContain(MENU.register)
    expect(label).toContain(MENU_MODUL.premiumList)
    expect(label).toContain(MENU_MODUL.inboxKomite)
  })
})

describe('saringPalet', () => {
  it('mencocokkan POTONGAN, bukan awalan', () => {
    // Pemakai mengingat kata TENGAH sesering kata depan.
    const hasil = saringPalet(daftarPalet(), 'claim')
    expect(hasil.map((h) => h.label)).toContain(MENU.inbox)
  })

  it('setiap KATA harus cocok, urutannya bebas', () => {
    const hasil = saringPalet(daftarPalet(), 'life claim')
    expect(hasil.map((h) => h.label)).toContain(MENU.inbox)
    // Dan kata yang tidak ada menolak barisnya.
    expect(saringPalet(daftarPalet(), 'claim borderaux')).toHaveLength(0)
  })

  it('kelompok ikut dicari', () => {
    // `PremiumList` label butirnya; `PremiumList Life` nama kelompoknya.
    const hasil = saringPalet(daftarPalet(), 'premiumlist life')
    expect(hasil.map((h) => h.label)).toContain(MENU_MODUL.premiumList)
  })

  it('kueri kosong mengembalikan seluruhnya, urutan sidebar', () => {
    const hasil = saringPalet(daftarPalet(), '   ')
    expect(hasil).toHaveLength(ENTRI_MENU.length)
    expect(hasil[0]?.label).toBe('Beranda')
  })

  it('huruf besar-kecil diabaikan', () => {
    expect(saringPalet(daftarPalet(), 'INBOX').length).toBeGreaterThan(0)
  })
})

describe('nol menu dikarang', () => {
  it('palet TIDAK memuat modul yang belum dimigrasi', () => {
    // ⛔ Empat belas kelompok berdiri di sidebar TANPA butir. Bila salah
    // satunya muncul di palet, ia dapat dibuka lewat Ctrl+K — layar yang
    // tidak ada.
    const berbutir = new Set(ENTRI_MENU.map((e) => e.kelompok))
    const tanpaButir = Object.values(MODUL).filter((n) => !berbutir.has(n))
    expect(tanpaButir).toHaveLength(14)
    for (const nama of tanpaButir) {
      // ⛔ DIPERSEMPIT KE MAKSUDNYA 28-09-2026 (sesi Treaty Contract Out).
      // Dulu: hasil pencarian nama kelompok harus KOSONG. Pencocokan palet
      // memakai POTONGAN kata, sehingga kueri `NB Treaty In` menemukan
      // butir `InboxTreatyContract` (`nb` dan `in` ada di dalam "inbox",
      // `treaty` di kelompoknya) - butir yang SAH, milik kelompok yang
      // sudah dimigrasi. Yang dijaga tetap sama: tidak satu pun hasil
      // BERKELOMPOK modul yang belum dimigrasi.
      for (const hasil of saringPalet(daftarPalet(), nama)) {
        expect(hasil.kelompok).not.toBe(nama)
      }
    }
  })
})
