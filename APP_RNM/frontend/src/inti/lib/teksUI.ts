/**
 * Teks BAWAAN komponen bersama (`components/ui/dasar.tsx`), per bahasa.
 *
 * Bawaannya `id` — setiap modul yang tidak memilih bahasa tampil persis
 * seperti sebelum berkas ini ada. Modul memilih `en` dengan membungkus
 * layarnya di `BahasaUI.Provider` (`components/ui/bahasaUI.tsx`); Treaty
 * Contract Out melakukannya [keputusan work owner 30-09-2026: "untuk bahasa
 * pake bahasa inggris, jangan indo"].
 *
 * Modul ini MURNI (tanpa React) supaya dapat diuji langsung.
 */

export type Bahasa = 'id' | 'en'

export interface TeksUI {
  memuat: string
  pilihKosong: string
  muatUlang: string
  tutup: string
  tutupEsc: string
  tidakAdaBaris: string
  menampilkan: (awal: number, akhir: number, total: number) => string
  halaman: (kini: number) => string
  halamanDari: (kini: number, total: number) => string
  sebelumnya: string
  berikutnya: string
  /** `Pilih`: nilai tersimpan yang tidak ada di daftar pilihan. */
  tidakDiDaftar: string
  /** `PilihSaring`: petunjuk kotak, daftar kosong, tombol panah. */
  ketikUntukMenyaring: string
  tidakCocok: string
  bukaDaftar: string
}

export const TEKS_UI: Readonly<Record<Bahasa, TeksUI>> = {
  id: {
    memuat: 'Memuat...',
    pilihKosong: '-- pilih --',
    muatUlang: 'Muat ulang',
    tutup: 'Tutup',
    tutupEsc: 'Tutup (Esc)',
    tidakAdaBaris: 'Tidak ada baris',
    menampilkan: (awal, akhir, total) => `Menampilkan ${awal}–${akhir} dari ${total}`,
    halaman: (kini) => `Halaman ${kini}`,
    halamanDari: (kini, total) => `Halaman ${kini} dari ${total}`,
    sebelumnya: 'Sebelumnya',
    berikutnya: 'Berikutnya',
    tidakDiDaftar: '(tidak ada di daftar referensi)',
    ketikUntukMenyaring: 'Ketik untuk menyaring',
    tidakCocok: 'Tidak ada yang cocok',
    bukaDaftar: 'Buka daftar',
  },
  en: {
    memuat: 'Loading...',
    pilihKosong: '-- select --',
    muatUlang: 'Reload',
    tutup: 'Close',
    tutupEsc: 'Close (Esc)',
    tidakAdaBaris: 'No rows',
    menampilkan: (awal, akhir, total) => `Showing ${awal}–${akhir} of ${total}`,
    halaman: (kini) => `Page ${kini}`,
    halamanDari: (kini, total) => `Page ${kini} of ${total}`,
    sebelumnya: 'Previous',
    berikutnya: 'Next',
    tidakDiDaftar: '(not in the reference list)',
    ketikUntukMenyaring: 'Type to filter',
    tidakCocok: 'No matches',
    bukaDaftar: 'Open list',
  },
}
