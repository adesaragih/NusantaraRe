import { MENU, MENU_MODUL, MODUL } from '../assets/labels'
import { MENU_TCO } from '../assets/labels.treaty-contract-out'

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
// ⚠️ tco5: `tco-kontrak` / `tco-klausul` TIDAK lagi punya butir menu — kedua
// layar dibuka tombol `ReinsType` / `List Description` layar tahun treaty.
// Nilainya tetap di union karena rutenya ada di `App.tsx`, yang memuat suntingan
// work owner yang belum di-commit (brief lanjutan 3: jangan disentuh).
export type ModulTetap = 'beranda' | 'inbox' | 'register' | 'premiumlist' | 'komite' | 'tco-tahun' | 'tco-kontrak' | 'tco-klausul'

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
 * ⛔ LIMA butir modul (tco5: Treaty Contract Out satu butir), ditambah Beranda. Empat belas kelompok lain berdiri
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
  // Treaty Contract Out — tco5 [keputusan work owner 29-09-2026]: SATU butir
  // "Treaty Contract Out" (asal harness InboxTreatyContract). Butir ReinsType
  // dan Description dibuang: di Pega keduanya popup form kontrak (b20778, b22196).
  { modul: 'tco-tahun', label: MENU_TCO.treatyContractOut, kelompok: MODUL.treatyContractOut },
]

/**
 * Modul backend pemilik tiap halaman - nama yang sama dengan `MODUL_AKTIF` dan
 * `GET /api/modul-aktif` (refactor bentuk B). `null` = milik aplikasi, selalu
 * tampil (Beranda).
 */
export const MODUL_BACKEND: Readonly<Record<ModulTetap, string | null>> = {
  beranda: null,
  inbox: 'claimlife',
  register: 'claimlife',
  premiumlist: 'premiumlist',
  komite: 'komite',
  'tco-tahun': 'treaty',
  'tco-kontrak': 'treaty',
  'tco-klausul': 'treaty',
}

/**
 * Apakah halaman ini tampil, menurut daftar modul aktif dari backend.
 *
 * ⛔ `null` - daftar belum terbaca atau backend tak terjangkau - berarti SEMUA
 * tampil: persis perilaku sebelum `MODUL_AKTIF` ada. Menyembunyikan menu
 * karena permintaannya gagal akan membuat aplikasi tampak kosong padahal yang
 * rusak hanya satu pembacaan.
 */
export function halamanAktif(halaman: ModulTetap, aktif: readonly string[] | null): boolean {
  const pemilik = MODUL_BACKEND[halaman]
  return aktif === null || pemilik === null || aktif.includes(pemilik)
}

/** Satu butir sidebar - bentuk yang Shell turunkan dari `ENTRI_MENU`. */
export interface ButirSidebar {
  halaman: ModulTetap
  label: string
}

/** Satu kelompok sidebar beserta butirnya. */
export interface KelompokSidebar {
  nama: string
  butir: readonly ButirSidebar[]
}

/**
 * Kelompok sidebar yang tampil, masing-masing dengan butir modul aktifnya saja.
 *
 * ⛔ Kelompok yang memang TANPA butir (belum dimigrasi) tetap berdiri dan
 * menyebut sebabnya. Kelompok yang butirnya HABIS tersaring - modulnya
 * NONAKTIF lewat `MODUL_AKTIF` - hilang sama sekali: modulnya ada, hanya
 * tidak dipasang di proses ini, jadi "belum dimigrasi" akan berbohong.
 */
export function kelompokTampil(
  kelompok: readonly KelompokSidebar[],
  aktif: readonly string[] | null,
): { k: KelompokSidebar; butir: readonly ButirSidebar[] }[] {
  return kelompok
    .map((k) => ({ k, butir: k.butir.filter((b) => halamanAktif(b.halaman, aktif)) }))
    .filter(({ k, butir }) => k.butir.length === 0 || butir.length > 0)
}

/** Satu baris hasil palet. */
export interface HasilPalet {
  kunci: string
  label: string
  kelompok: string
  modul: ModulTetap
}

/** Daftar yang dapat dicari palet - hanya menu modul yang aktif. */
export function daftarPalet(aktif: readonly string[] | null = null): HasilPalet[] {
  return ENTRI_MENU.filter((e) => halamanAktif(e.modul, aktif)).map((e) => ({
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
