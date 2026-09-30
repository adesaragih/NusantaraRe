/**
 * Daftar menu yang dapat dicari palet Ctrl+K — butir **bg**.
 *
 * Pola `REFERENSI_UI/frontend/src/lib/daftarMenu.ts`, disederhanakan: di
 * sana separuh menu datang dari daftar entitas master yang dimuat backend;
 * di sini seluruhnya tetap, sebab modul kami tidak punya registry dinamis.
 *
 * # Refactor bentuk B (30-09-2026): berkas ini tidak mengenal modul
 *
 * Isi menu dirakit dari `modul/<nama>/menu.ts` oleh daftar modul aplikasi
 * (`modul/daftar.ts`, `ENTRI_MENU`) dan diteruskan `App.tsx` ke Shell. Yang
 * tinggal di sini hanya BENTUK entri dan aturan yang sama untuk semua modul:
 * penyaring modul aktif, penurun kelompok sidebar, dan pencari palet.
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
 * Kedua daftar dijaga tetap sama oleh `daftar.sinkron.test.ts`, DUA
 * ARAH — dan berkas itu menuliskan harga pilihan ini apa adanya.
 */

/** Halaman kerangka aplikasi - bukan milik modul mana pun, selalu ada. */
export const HALAMAN_BERANDA = 'beranda'

/**
 * Satu butir menu milik sebuah modul, seperti yang ditulis `menu.ts`-nya.
 *
 * `H` adalah union halaman modul itu, bukan `string`: nama halaman yang salah
 * ketik membuka layar kosong, dan tidak satu pun uji TEKS akan melihatnya.
 */
export interface ButirMenuModul<H extends string = string> {
  /** Halaman tujuan — argumen yang sidebar pakai. */
  modul: H
  /** Nama yang TAMPIL. Selalu dari label, tidak pernah diketik ulang. */
  label: string
  /** Nama kelompok sidebar tempat ia berada — ditampilkan sebagai konteks. */
  kelompok: string
  /**
   * true = butir SATU-SATUNYA kelompoknya tampil DATAR di sidebar: satu
   * tombol langsung, tanpa judul kelompok yang dilipat dan tanpa anak.
   * Kelompok berbutir lebih dari satu mengabaikannya. Bawaan: bertingkat.
   */
  datar?: boolean
}

/** Satu entri menu yang dapat dicari, beserta modul backend pemiliknya. */
export interface EntriMenu<H extends string = string> extends ButirMenuModul<H> {
  /**
   * Modul backend pemilik - nama yang sama dengan `MODUL_AKTIF` dan
   * `GET /api/modul-aktif`. `null` = milik aplikasi, selalu tampil (Beranda).
   */
  pemilik: string | null
}

/**
 * Apakah modul `pemilik` dipasang, menurut daftar modul aktif dari backend.
 *
 * ⛔ `null` - daftar belum terbaca atau backend tak terjangkau - berarti SEMUA
 * tampil: persis perilaku sebelum `MODUL_AKTIF` ada. Menyembunyikan menu
 * karena permintaannya gagal akan membuat aplikasi tampak kosong padahal yang
 * rusak hanya satu pembacaan.
 */
export function modulDipasang(pemilik: string | null, aktif: readonly string[] | null): boolean {
  return aktif === null || pemilik === null || aktif.includes(pemilik)
}

/** Satu butir sidebar - bentuk yang Shell turunkan dari menu. */
export interface ButirSidebar<H extends string = string> {
  halaman: H
  label: string
  pemilik: string | null
  /** Lihat `ButirMenuModul.datar`; hanya ada bila true. */
  datar?: true
}

/** Satu kelompok sidebar beserta butirnya. */
export interface KelompokSidebar<H extends string = string> {
  nama: string
  butir: readonly ButirSidebar<H>[]
}

/** Butir satu kelompok, DITURUNKAN dari menu - tidak pernah diketik ulang. */
export function butirKelompok<H extends string>(menu: readonly EntriMenu<H>[], nama: string): ButirSidebar<H>[] {
  return menu.filter((e) => e.kelompok === nama).map((e) => ({
    halaman: e.modul,
    label: e.label,
    pemilik: e.pemilik,
    ...(e.datar === true ? { datar: true as const } : {}),
  }))
}

/** Butir yang dirender DATAR: kelompok beranggota tepat satu butir bertanda `datar`. */
export function butirDatar<H extends string>(butir: readonly ButirSidebar<H>[]): ButirSidebar<H> | undefined {
  const [b] = butir
  return butir.length === 1 && b?.datar === true ? b : undefined
}

/**
 * Kelompok sidebar yang tampil, masing-masing dengan butir modul aktifnya saja.
 *
 * ⛔ Kelompok yang memang TANPA butir (belum dimigrasi) tetap berdiri dan
 * menyebut sebabnya. Kelompok yang butirnya HABIS tersaring - modulnya
 * NONAKTIF lewat `MODUL_AKTIF` - hilang sama sekali: modulnya ada, hanya
 * tidak dipasang di proses ini, jadi "belum dimigrasi" akan berbohong.
 */
export function kelompokTampil<H extends string>(
  kelompok: readonly KelompokSidebar<H>[],
  aktif: readonly string[] | null,
): { k: KelompokSidebar<H>; butir: readonly ButirSidebar<H>[] }[] {
  return kelompok
    .map((k) => ({ k, butir: k.butir.filter((b) => modulDipasang(b.pemilik, aktif)) }))
    .filter(({ k, butir }) => k.butir.length === 0 || butir.length > 0)
}

/** Satu baris hasil palet. */
export interface HasilPalet<H extends string = string> {
  kunci: string
  label: string
  kelompok: string
  modul: H
}

/** Daftar yang dapat dicari palet - hanya menu modul yang aktif. */
export function daftarPalet<H extends string>(
  menu: readonly EntriMenu<H>[],
  aktif: readonly string[] | null = null,
): HasilPalet<H>[] {
  return menu.filter((e) => modulDipasang(e.pemilik, aktif)).map((e) => ({
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
export function saringPalet<H extends string>(daftar: readonly HasilPalet<H>[], kueri: string): HasilPalet<H>[] {
  const kata = kueri.trim().toLowerCase().split(/\s+/).filter(Boolean)
  if (kata.length === 0) return [...daftar]
  return daftar.filter((h) => {
    const jerami = (h.label + ' ' + h.kelompok).toLowerCase()
    return kata.every((k) => jerami.includes(k))
  })
}
