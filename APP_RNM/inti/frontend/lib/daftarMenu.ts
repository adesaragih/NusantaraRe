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
 * (`frontend/daftar.ts`, `ENTRI_MENU`) dan diteruskan `App.tsx` ke Shell. Yang
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
 *
 * # Sejak 30-09-2026: sidebar dan palet dari tabel `M_NAV_MENU`
 *
 * Permintaan work owner (`PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`): menu dibuat
 * dari tabel, supaya kelak dapat disaring per akun. `GET /api/menu` mengirim
 * pohon GROUPMENU → kelompok → butir; `susunMenu` di bawah MEMOTONGNYA dengan
 * rute yang benar-benar terdaftar di `frontend/daftar.ts`. Keduanya harus
 * sepakat: baris tabel tanpa rute tidak tampil (dicatat di konsol), rute
 * tanpa baris tabel juga tidak - dan `frontend/daftar.menuTabel.test.ts`
 * menjaganya dua arah terhadap isi awal migrasi 900.
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
  datar?: true
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
  /** `M_NAV_MENU.KODE` - nama modul backend. */
  kode: string
  /** Nama tampil - `M_NAV_MENU.LABEL`, nama folder korpus VERBATIM. */
  nama: string
  /** false = "belum dimigrasi" (`M_NAV_MENU.DIMIGRASI = '0'`). */
  dimigrasi: boolean
  butir: readonly ButirSidebar<H>[]
}

/** Satu kepala bagian sidebar - GROUPMENU (TREATY, FACULTATIVE, KLAIM, MASTER). */
export interface GolonganSidebar<H extends string = string> {
  kode: string
  kelompok: readonly KelompokSidebar<H>[]
}

// ---------------------------------------------------------------------------
// Bentuk `GET /api/menu` - backend `inti/menu` (`menu.Menu`).
// ---------------------------------------------------------------------------

export interface ButirMenuTabel {
  /** Kunci halaman frontend. */
  kode: string
  label: string
  /** Modul backend pemilik. */
  modul: string
}

export interface KelompokMenuTabel {
  kode: string
  label: string
  modul: string
  dimigrasi: boolean
  butir: readonly ButirMenuTabel[]
}

export interface GolonganMenuTabel {
  kode: string
  kelompok: readonly KelompokMenuTabel[]
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
 * menu".
 */
export function bentukMenuTabel(x: unknown): x is MenuTabel {
  const butirSah = (b: unknown): boolean => objek(b) && teks(b.kode) && teks(b.label) && teks(b.modul)
  const kelompokSah = (k: unknown): boolean =>
    objek(k) && teks(k.kode) && teks(k.label) && teks(k.modul) && typeof k.dimigrasi === 'boolean' &&
    Array.isArray(k.butir) && k.butir.every(butirSah)
  const golonganSah = (g: unknown): boolean =>
    objek(g) && teks(g.kode) && Array.isArray(g.kelompok) && g.kelompok.every(kelompokSah)
  return objek(x) && Array.isArray(x.golongan) && x.golongan.every(golonganSah)
}

/** Hasil `susunMenu`: yang sidebar render, yang palet cari, dan yang dicatat. */
export interface MenuTersusun<H extends string = string> {
  /** Golongan yang TAMPIL, masing-masing dengan kelompok yang tampil. */
  golongan: GolonganSidebar<H>[]
  /** Daftar palet: Beranda, lalu setiap butir yang tampil - urutan sidebar. */
  entri: EntriMenu<H>[]
  /** Butir tabel tanpa rute frontend, `<kelompok>/<butir>` - untuk konsol. */
  tanpaRute: string[]
  /** Butir tabel di bawah kelompok `DIMIGRASI = '0'`, tidak tampil - untuk konsol. */
  terlipat: string[]
}

/** Entri milik APLIKASI (Beranda, pemilik `null`) - bukan baris tabel, selalu ada. */
export function entriAplikasi<H extends string>(rute: readonly EntriMenu<H>[]): EntriMenu<H>[] {
  return rute.filter((e) => e.pemilik === null)
}

/**
 * Memotong pohon `GET /api/menu` dengan rute frontend yang terdaftar.
 *
 * Aturannya (brief menu 30-09-2026 §3):
 *   - butir tabel TANPA rute di `rute` tidak tampil, dan dicatat di `tanpaRute`
 *   - rute TANPA butir tabel tidak tampil (yang diulang hanya baris tabel)
 *   - LABEL dari tabel; penanda `datar` dari rute frontend (cara tampil)
 *   - kelompok `DIMIGRASI = '0'` tetap BERDIRI, terlipat "belum dimigrasi" -
 *     tanpa butir, walau tabel (keliru) memberinya butir berute: butir itu
 *     dicatat di `terlipat`
 *   - kelompok dimigrasi yang butirnya habis - modulnya nonaktif (MODUL_AKTIF,
 *     backend tidak mengirim butirnya) atau tak satu pun berute - hilang:
 *     "belum dimigrasi" akan berbohong
 *   - golongan tanpa kelompok tampil hilang
 *
 * Beranda bukan baris tabel: ia diambil dari `rute` (pemilik `null`).
 */
export function susunMenu<H extends string>(tabel: MenuTabel, rute: readonly EntriMenu<H>[]): MenuTersusun<H> {
  const hasil: MenuTersusun<H> = { golongan: [], entri: entriAplikasi(rute), tanpaRute: [], terlipat: [] }
  for (const g of tabel.golongan) {
    const kelompok: KelompokSidebar<H>[] = []
    for (const k of g.kelompok) {
      if (!k.dimigrasi) {
        for (const b of k.butir) hasil.terlipat.push(`${k.kode}/${b.kode}`)
        kelompok.push({ kode: k.kode, nama: k.label, dimigrasi: false, butir: [] })
        continue
      }
      const butir: ButirSidebar<H>[] = []
      for (const b of k.butir) {
        const r = rute.find((e) => e.pemilik !== null && e.modul === b.kode)
        if (r === undefined) {
          hasil.tanpaRute.push(`${k.kode}/${b.kode}`)
          continue
        }
        butir.push({ halaman: r.modul, label: b.label, pemilik: b.modul, ...(r.datar === true ? { datar: true as const } : {}) })
      }
      if (butir.length === 0) continue
      kelompok.push({ kode: k.kode, nama: k.label, dimigrasi: true, butir })
      // Palet = sidebar: entri dari butir yang SAMA, datar ikut apa adanya.
      for (const { halaman, ...sisa } of butir) hasil.entri.push({ modul: halaman, kelompok: k.label, ...sisa })
    }
    if (kelompok.length > 0) hasil.golongan.push({ kode: g.kode, kelompok })
  }
  return hasil
}

/** Butir yang dirender DATAR: kelompok beranggota tepat satu butir bertanda `datar`. */
export function butirDatar<H extends string>(butir: readonly ButirSidebar<H>[]): ButirSidebar<H> | undefined {
  const [b] = butir
  return butir.length === 1 && b?.datar === true ? b : undefined
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
 * Sejak menu dari tabel (30-09-2026) `menu` adalah `MenuTersusun.entri` -
 * butir modul nonaktif sudah tidak dikirim backend, jadi saringan modul aktif
 * tidak lagi di sini.
 */
export function daftarPalet<H extends string>(menu: readonly EntriMenu<H>[]): HasilPalet<H>[] {
  return menu.map((e) => ({
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
