import { MENU, MENU_MODUL, MODUL } from '../assets/labels'

/**
 * Daftar menu yang dapat dicari palet Ctrl+K — butir **bg**.
 *
 * Pola `REFERENSI_UI/frontend/src/lib/daftarMenu.ts`, disederhanakan: di
 * sana separuh menu datang dari daftar entitas master yang dimuat backend;
 * di sini seluruhnya tetap, sebab modul kami tidak punya registry dinamis.
 *
 * # Kenapa daftar ini TERPISAH dari sidebar
 *
 * Diukur, bukan diduga:
 *
 *   `components/KelompokMenu.tsx` merender anaknya HANYA ketika kelompoknya
 *   terbuka — kelompok yang terlipat melepas anaknya dari DOM.
 *
 * Jadi palet TIDAK boleh membaca DOM sidebar: ia hanya akan menemukan entri
 * di kelompok yang kebetulan terbuka, lalu menjawab "tidak ada" untuk menu
 * yang ADA. Dan sejak butir **bg** empat belas kelompok memang terlipat.
 *
 * Kedua daftar dijaga tetap sama oleh `daftarMenu.sinkron.test.ts`, DUA
 * ARAH — dan berkas itu menuliskan harga pilihan ini apa adanya.
 */

/**
 * Halaman yang sidebar buka, sebagai UNION — bukan `string`.
 *
 * ⛔ Nama halaman yang salah ketik membuka layar kosong, dan tidak satu pun
 * uji TEKS akan melihatnya: `'inbok'` tetap cocok dengan setiap pola yang
 * memeriksa bentuk. Union memindahkan penjagaannya ke kompiler — lebih
 * awal, dan tanpa pagar tambahan.
 */
export type ModulTetap = 'beranda' | 'inbox' | 'register' | 'premiumlist' | 'komite'

/** Satu entri menu yang dapat dicari. */
export interface EntriMenu {
  /** Halaman tujuan — argumen yang sidebar pakai. */
  modul: ModulTetap
  /** Nama yang TAMPIL. Selalu dari label, tidak pernah diketik ulang. */
  label: string
  /** Nama kelompok sidebar tempat ia berada — ditampilkan sebagai konteks. */
  kelompok: string
}

/**
 * Entri sidebar yang benar-benar dapat dibuka.
 *
 * ⛔ EMPAT butir modul, ditambah Beranda. Empat belas kelompok lain berdiri
 * di sidebar TANPA butir — dan karena itu tidak ada di sini pula. Entri yang
 * berdiri di daftar ini tetapi tidak di sidebar dapat dibuka lewat palet
 * walau menunya tidak terlihat; itu persis cacat yang REFERENSI_UI bayar
 * sekali dan tuliskan pelajarannya.
 *
 * Urutannya SAMA dengan urutan render sidebar, supaya apa yang pemakai lihat
 * pertama kali cocok dengan yang sudah ia hafal letaknya.
 */
export const ENTRI_MENU: readonly EntriMenu[] = [
  { modul: 'beranda', label: 'Beranda', kelompok: 'Beranda' },
  { modul: 'inbox', label: MENU.inbox, kelompok: MODUL.claimLife },
  { modul: 'register', label: MENU.register, kelompok: MODUL.claimLife },
  { modul: 'premiumlist', label: MENU_MODUL.premiumList, kelompok: MODUL.premiumListLife },
  { modul: 'komite', label: MENU_MODUL.inboxKomite, kelompok: MODUL.komiteClaimLife },
]

/** Satu baris hasil palet. */
export interface HasilPalet {
  kunci: string
  label: string
  kelompok: string
  modul: ModulTetap
}

/** Daftar yang dapat dicari palet. */
export function daftarPalet(): HasilPalet[] {
  return ENTRI_MENU.map((e) => ({
    kunci: 'modul:' + e.modul,
    label: e.label,
    kelompok: e.kelompok,
    modul: e.modul,
  }))
}

/**
 * Saringan palet.
 *
 * # Mencocokkan POTONGAN, bukan awalan
 *
 * Nama menu di sini berkata-kata banyak (`Inbox Claim Life`,
 * `Inbox Komite`). Pemakai mengingat kata TENGAH sesering kata depan, jadi
 * pencocokan awalan akan gagal untuk ketikan yang paling wajar.
 *
 * # Setiap KATA harus cocok, dan urutannya bebas
 *
 * `'life claim'` menemukan `Inbox Claim Life` — mengetik dua kata yang
 * diingat, dalam urutan apa pun, adalah cara orang mencari. Pencocokan
 * seluruh-frasa akan menolaknya.
 *
 * # Huruf besar-kecil diabaikan, kelompok ikut dicari
 *
 * Mengetik `premiumlist` menemukan butir di kelompok `PremiumList Life`
 * walau label butirnya sendiri hanya `PremiumList`.
 */
export function saringPalet(daftar: readonly HasilPalet[], kueri: string): HasilPalet[] {
  const kata = kueri.trim().toLowerCase().split(/\s+/).filter(Boolean)
  if (kata.length === 0) return [...daftar]
  return daftar.filter((h) => {
    const jerami = (h.label + ' ' + h.kelompok).toLowerCase()
    return kata.every((k) => jerami.includes(k))
  })
}
