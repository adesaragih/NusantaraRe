/**
 * Menu sidebar dan palet Ctrl+K — butir **bg**, menu dari tabel `M_NAV_MENU`.
 *
 * Pola `REFERENSI_UI/frontend/src/lib/daftarMenu.ts`, disederhanakan: di
 * sana separuh menu datang dari daftar entitas master yang dimuat backend;
 * di sini seluruhnya datang dari `GET /api/menu`.
 *
 * # Refactor bentuk B (30-09-2026): berkas ini tidak mengenal modul
 *
 * Halaman awal setiap modul dirakit dari `modul/<nama>/menu.ts` oleh daftar
 * modul aplikasi (`frontend/daftar.ts`, `ENTRI_MENU`) dan diteruskan `App.tsx`
 * ke Shell. Yang tinggal di sini hanya BENTUK entri dan aturan yang sama untuk
 * semua modul: penyaring modul aktif, pemotong menu tabel, dan pencari palet.
 *
 * # Menu DATAR — keputusan work owner 30-09-2026
 *
 * `PROMPT-MENU-DATAR-PER-GROUPMENU.md`: *"menu jangan ada model seperti child.
 * Buat grouping menu antar GROUPMENU dari tabel M_NAV_MENU ... 1 modul 1
 * menu."* `GET /api/menu` mengirim golongan → modul (satu tingkat), dan
 * `susunMenu` di bawah MEMOTONGNYA dengan modul yang benar-benar terdaftar di
 * `frontend/daftar.ts`: baris dimigrasi tanpa modul frontend tidak tampil
 * (dicatat di konsol), modul frontend tanpa baris tabel tidak tampil pula.
 * Klik tombol modul membuka halaman awalnya (`HALAMAN_AWAL_<X>` di `menu.ts`).
 * `frontend/daftar.menuTabel.test.ts` menjaga keduanya dua arah terhadap
 * hasil bersih migrasi 900 + 901 + slot menu.
 *
 * # Palet membaca daftar yang SAMA dengan sidebar
 *
 * `susunMenu` menghasilkan tombol sidebar DAN entri palet dari satu
 * pemotongan, jadi keduanya tidak dapat menyimpang (`daftar.sinkron.test.ts`).
 */

/** Halaman kerangka aplikasi - bukan milik modul mana pun, selalu ada. */
export const HALAMAN_BERANDA = 'beranda'

/**
 * Halaman Kelola User (keputusan work owner 01-10-2026) - milik aplikasi
 * seperti Beranda, tetapi tampil HANYA bagi akun yang memegang menunya.
 */
export const HALAMAN_KELOLA_USER = 'kelolauser'

/**
 * KODE menu Kelola User - `menu.KodeKelolaUser` di backend. Bukan baris
 * `M_NAV_MENU`: `GET /api/menu` mengirimnya di golongan ADMIN hanya bagi
 * pemegangnya.
 */
export const KODE_MENU_KELOLA_USER = 'kelolauser'

/**
 * Modul yang boleh dipasang untuk akun yang login: modul aktif (MODUL_AKTIF)
 * yang menunya ia pegang (`M_LOGIN_GO_MENU`). Padanan gerbang 403 backend -
 * layar modul yang tidak boleh dibuka tidak dipasang, jadi tidak satu pun
 * permintaannya berangkat.
 *
 * `menuAkun` null = tanpa saringan akun (mode stub). `aktif` null = daftar
 * modul aktif belum terbaca: `semua` dipakai, seperti `modulDipasang`.
 */
export function modulUntukAkun(
  aktif: readonly string[] | null,
  menuAkun: readonly string[] | null,
  semua: readonly string[],
): readonly string[] | null {
  if (menuAkun === null) return aktif
  return (aktif ?? semua).filter((n) => menuAkun.includes(n))
}

/**
 * Satu entri menu yang dapat dibuka - Beranda, atau halaman AWAL satu modul.
 *
 * `H` adalah union halaman, bukan `string`: nama halaman yang salah ketik
 * membuka layar kosong, dan tidak satu pun uji TEKS akan melihatnya.
 */
export interface EntriMenu<H extends string = string> {
  /** Halaman tujuan — argumen yang sidebar pakai. */
  modul: H
  /**
   * Nama yang tampil bila menu tabel belum terbaca (Beranda, kartu Beranda).
   * Tombol sidebar dan palet memakai `M_NAV_MENU.LABEL`, bukan ini.
   */
  label: string
  /**
   * Nama modul - nama folder korpus VERBATIM (`Beranda` untuk Beranda); kunci
   * kartu Beranda. Penjaga dua arah menuntut ia sama dengan LABEL barisnya.
   */
  kelompok: string
  /**
   * Modul backend pemilik - nama yang sama dengan `MODUL_AKTIF` dan
   * `GET /api/modul-aktif`. `null` = milik aplikasi, selalu tampil (Beranda).
   */
  pemilik: string | null
  /**
   * GROUPMENU tombol ini (TREATY, …) - konteks baris palet. Hanya entri yang
   * disusun dari menu tabel (`susunMenu`) yang membawanya.
   */
  golongan?: string
  /**
   * Seluruh halaman modul ini - tombolnya menyala selama salah satunya tampil
   * (mis. Outstanding milik Claim Life). Kosong untuk Beranda.
   */
  halamanModul?: readonly H[]
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

/** Satu tombol sidebar - SATU modul. */
export interface TombolModul<H extends string = string> {
  /** `M_NAV_MENU.KODE` - nama modul backend. */
  kode: string
  /** `M_NAV_MENU.LABEL` - nama folder korpus VERBATIM. */
  label: string
  /** Halaman awal modul; `null` = belum dimigrasi (tombol nonaktif). */
  halaman: H | null
  /** Seluruh halaman modul - tombolnya aktif selama salah satunya tampil. */
  halamanModul: readonly H[]
}

/** Satu kepala bagian sidebar - GROUPMENU (TREATY, FACULTATIVE, KLAIM, MASTER). */
export interface GolonganSidebar<H extends string = string> {
  kode: string
  modul: readonly TombolModul<H>[]
}

// ---------------------------------------------------------------------------
// Bentuk `GET /api/menu` - backend `inti/backend/menu` (`menu.Menu`).
// ---------------------------------------------------------------------------

export interface ModulMenuTabel {
  /** Nama modul backend (`M_NAV_MENU.KODE` = `MODUL`). */
  kode: string
  label: string
  modul: string
  urutan: number
  /** false = belum dimigrasi (`DIMIGRASI = '0'`). */
  dimigrasi: boolean
}

export interface GolonganMenuTabel {
  kode: string
  modul: readonly ModulMenuTabel[]
}

export interface MenuTabel {
  golongan: readonly GolonganMenuTabel[]
}

/** Keadaan pembacaan `GET /api/menu`: `null` = sedang dimuat. */
export type KeadaanMenuTabel = { menu: MenuTabel } | { galat: unknown } | null

const teks = (x: unknown): x is string => typeof x === 'string'
const objek = (x: unknown): x is Record<string, unknown> => typeof x === 'object' && x !== null

/**
 * Apakah `x` berbentuk `GET /api/menu`.
 *
 * ⛔ Bentuk yang tidak dikenal DITOLAK (pemanggilnya menampilkan galat), bukan
 * dijadikan menu kosong: sidebar yang kosong diam-diam terbaca "aplikasi tanpa
 * menu". Bentuk pohon lama (golongan → kelompok → butir) ikut ditolak.
 */
export function bentukMenuTabel(x: unknown): x is MenuTabel {
  const modulSah = (m: unknown): boolean =>
    objek(m) && teks(m.kode) && teks(m.label) && teks(m.modul) && typeof m.urutan === 'number' &&
    typeof m.dimigrasi === 'boolean'
  const golonganSah = (g: unknown): boolean =>
    objek(g) && teks(g.kode) && Array.isArray(g.modul) && g.modul.every(modulSah)
  return objek(x) && Array.isArray(x.golongan) && x.golongan.every(golonganSah)
}

/** Hasil `susunMenu`: yang sidebar render, yang palet cari, dan yang dicatat. */
export interface MenuTersusun<H extends string = string> {
  /** Golongan yang TAMPIL, masing-masing dengan tombol modulnya. */
  golongan: GolonganSidebar<H>[]
  /** Daftar palet: Beranda, lalu setiap modul dimigrasi yang tampil - urutan sidebar. */
  entri: EntriMenu<H>[]
  /** Baris dimigrasi tanpa modul frontend terdaftar, `<kode>` - untuk konsol. */
  tanpaRute: string[]
  /** Modul frontend terdaftar yang barisnya `DIMIGRASI = '0'` (tombol nonaktif) - untuk konsol. */
  nonaktifBerute: string[]
}

/** Entri milik APLIKASI (Beranda, pemilik `null`) - bukan baris tabel, selalu ada. */
export function entriAplikasi<H extends string>(rute: readonly EntriMenu<H>[]): EntriMenu<H>[] {
  return rute.filter((e) => e.pemilik === null)
}

/**
 * Memotong menu `GET /api/menu` dengan modul frontend yang terdaftar.
 *
 * Aturannya (brief menu datar 30-09-2026 §1, §3):
 *   - satu TOMBOL per modul, di bawah kepala GROUPMENU; label = LABEL tabel
 *   - modul `DIMIGRASI = '0'` tampil NONAKTIF ("belum dimigrasi"), tanpa
 *     halaman - walau frontend (keliru) mendaftarkannya: itu dicatat di
 *     `nonaktifBerute`
 *   - modul dimigrasi TANPA modul frontend terdaftar tidak tampil, dicatat di
 *     `tanpaRute`; modul frontend tanpa baris tabel tidak tampil pula
 *   - modul dimigrasi di luar MODUL_AKTIF tidak dikirim backend - tidak tampil
 *   - golongan tanpa tombol hilang
 *   - palet: Beranda lalu setiap tombol yang dapat dibuka, urutan sidebar
 *
 * Beranda bukan baris tabel: ia diambil dari `rute` (pemilik `null`).
 */
export function susunMenu<H extends string>(tabel: MenuTabel, rute: readonly EntriMenu<H>[]): MenuTersusun<H> {
  const hasil: MenuTersusun<H> = { golongan: [], entri: entriAplikasi(rute), tanpaRute: [], nonaktifBerute: [] }
  for (const g of tabel.golongan) {
    const modul: TombolModul<H>[] = []
    for (const m of g.modul) {
      const r = rute.find((e) => e.pemilik !== null && e.pemilik === m.kode)
      if (!m.dimigrasi) {
        if (r !== undefined) hasil.nonaktifBerute.push(m.kode)
        modul.push({ kode: m.kode, label: m.label, halaman: null, halamanModul: [] })
        continue
      }
      if (r === undefined) {
        hasil.tanpaRute.push(m.kode)
        continue
      }
      modul.push({ kode: m.kode, label: m.label, halaman: r.modul, halamanModul: r.halamanModul ?? [r.modul] })
      // Palet = sidebar: entri dari tombol yang SAMA. `kelompok` tetap nama
      // modul (satu makna); golongannya dibawa sendiri untuk konteks palet.
      hasil.entri.push({ modul: r.modul, label: m.label, kelompok: m.label, golongan: g.kode, pemilik: m.kode })
    }
    if (modul.length > 0) hasil.golongan.push({ kode: g.kode, modul })
  }
  return hasil
}

/** Satu baris hasil palet. */
export interface HasilPalet<H extends string = string> {
  kunci: string
  label: string
  kelompok: string
  modul: H
}

/**
 * Daftar yang dapat dicari palet.
 *
 * `menu` adalah `MenuTersusun.entri` - modul nonaktif sudah tidak dikirim
 * backend, dan modul yang belum dimigrasi tidak punya entri.
 */
export function daftarPalet<H extends string>(menu: readonly EntriMenu<H>[]): HasilPalet<H>[] {
  return menu.map((e) => ({
    kunci: 'modul:' + e.modul,
    label: e.label,
    // Konteks baris palet: GROUPMENU tombolnya; Beranda (tanpa golongan) memakai namanya sendiri.
    kelompok: e.golongan ?? e.kelompok,
    modul: e.modul,
  }))
}

/**
 * Saringan palet.
 *
 * # Mencocokkan POTONGAN, bukan awalan
 *
 * Nama modul di sini berkata-kata banyak (`Komite Claim Life`,
 * `Treaty Contract Out`). Pemakai mengingat kata TENGAH sesering kata depan,
 * jadi pencocokan awalan akan gagal untuk ketikan yang paling wajar.
 *
 * # Setiap KATA harus cocok, dan urutannya bebas
 *
 * `'life claim'` menemukan `Claim Life` — mengetik dua kata yang diingat,
 * dalam urutan apa pun, adalah cara orang mencari. Pencocokan seluruh-frasa
 * akan menolaknya.
 *
 * # Huruf besar-kecil diabaikan, golongan ikut dicari
 *
 * Mengetik `klaim` menemukan setiap modul di golongan KLAIM.
 */
export function saringPalet<H extends string>(daftar: readonly HasilPalet<H>[], kueri: string): HasilPalet<H>[] {
  const kata = kueri.trim().toLowerCase().split(/\s+/).filter(Boolean)
  if (kata.length === 0) return [...daftar]
  return daftar.filter((h) => {
    const jerami = (h.label + ' ' + h.kelompok).toLowerCase()
    return kata.every((k) => jerami.includes(k))
  })
}
