// DAFTAR modul frontend - satu-satunya berkas yang mengenal semua modul
// sekaligus. Padanan daftar Go `inti/backend/daftar`.
//
// Modul di `modul/<nama>/frontend/` tidak pernah mengimpor modul lain, dan
// `inti/frontend/` tidak mengenal modul mana pun; berkas ini yang merakit menu
// dan rute tiap modul, dan `App.tsx` yang memasang modul yang AKTIF
// (`GET /api/modul-aktif`).
//
// # Struktur tim satu folder per modul (30-09-2026): daftar dari FOLDER
//
// ⛔ Nol baris per modul di berkas ini. `import.meta.glob` (Vite) mengumpulkan
// setiap `modul/*/frontend/menu.ts` (ekspor `PENDAFTARAN_MENU`) dan
// `modul/*/frontend/rute.tsx` (ekspor `RUTE_MODUL`) saat dibangun. Menambah
// modul = folder baru dengan dua berkas itu; berkas ini tidak disentuh. Union
// halaman terbentuk dari `HalamanModul`, yang diperluas setiap `menu.ts`.
//
// Nama modulnya SAMA dengan nama folder dan `const Nama` di
// `modul/<nama>/backend/modul.go` - dijaga `rakitModulFrontend` di bawah dan
// `daftar.modulAktif.test.ts`.

import { HALAMAN_BERANDA, modulDipasang, type EntriMenu } from '../inti/frontend/lib/daftarMenu'
import type { HalamanTerdaftar, MenuModul, ModulFrontend, RuteModul } from '../inti/frontend/modul'

/**
 * Halaman yang aplikasi dapat tampilkan, sebagai UNION — bukan `string`.
 *
 * ⛔ Nama halaman yang salah ketik membuka layar kosong, dan tidak satu pun
 * uji TEKS akan melihatnya: `'inbok'` tetap cocok dengan setiap pola yang
 * memeriksa bentuk. Union memindahkan penjagaannya ke kompiler — lebih
 * awal, dan tanpa pagar tambahan.
 */
export type Halaman = typeof HALAMAN_BERANDA | HalamanTerdaftar

/** Berkas `menu.ts` dan `rute.tsx` satu modul, berkunci jalur glob. */
export type BerkasMenu = Record<string, { PENDAFTARAN_MENU: MenuModul }>
export type BerkasRute = Record<string, { RUTE_MODUL: RuteModul<Halaman> }>

/** Nama folder modul dari jalur glob `../modul/<nama>/frontend/<berkas>`. */
function folderModul(jalur: string): string {
  const m = /\/modul\/([^/]+)\/frontend\/[^/]+$/.exec(jalur)
  if (m === null) throw new Error(`daftar modul: jalur di luar modul/<nama>/frontend: ${jalur}`)
  return m[1] ?? ''
}

/**
 * Merakit modul frontend dari berkas `menu.ts` dan `rute.tsx` setiap folder,
 * berurutan menurut NAMA folder - urutan yang sama dengan daftar Go.
 *
 * ⛔ Folder yang hanya punya salah satunya, atau yang `PENDAFTARAN_MENU.nama`-nya
 * berbeda dari nama foldernya, DITOLAK saat aplikasi dimuat - bukan modul yang
 * diam-diam tidak tampil.
 */
export function rakitModulFrontend(menu: BerkasMenu, rute: BerkasRute): ModulFrontend<Halaman>[] {
  const perFolder = new Map<string, { menu?: MenuModul; rute?: RuteModul<Halaman> }>()
  for (const [jalur, isi] of Object.entries(menu)) {
    perFolder.set(folderModul(jalur), { ...perFolder.get(folderModul(jalur)), menu: isi.PENDAFTARAN_MENU })
  }
  for (const [jalur, isi] of Object.entries(rute)) {
    perFolder.set(folderModul(jalur), { ...perFolder.get(folderModul(jalur)), rute: isi.RUTE_MODUL })
  }
  return [...perFolder.entries()]
    .sort(([a], [b]) => (a < b ? -1 : a > b ? 1 : 0))
    .map(([folder, { menu: m, rute: r }]) => {
      if (m === undefined || r === undefined) {
        throw new Error(
          `daftar modul: modul/${folder}/frontend punya ${m === undefined ? 'rute.tsx tanpa menu.ts' : 'menu.ts tanpa rute.tsx'}`,
        )
      }
      if (m.nama !== folder) {
        throw new Error(`daftar modul: modul/${folder}/frontend/menu.ts menyebut nama "${m.nama}"; nama modul = nama foldernya`)
      }
      return { nama: m.nama, kelompok: m.kelompok, halaman: m.halaman, menu: m.menu, Rute: r }
    })
}

/**
 * Modul terdaftar, berurutan menurut NAMA - urutan yang sama dengan daftar Go
 * (`inti/backend/daftar`) dan `GET /api/modul-aktif`.
 */
export const MODUL_FRONTEND: readonly ModulFrontend<Halaman>[] = rakitModulFrontend(
  import.meta.glob<{ PENDAFTARAN_MENU: MenuModul }>('../modul/*/frontend/menu.ts', { eager: true }),
  import.meta.glob<{ RUTE_MODUL: RuteModul<Halaman> }>('../modul/*/frontend/rute.tsx', { eager: true }),
)

/**
 * Entri sidebar yang benar-benar dapat dibuka.
 *
 * ⛔ LIMA butir modul (tco5: Treaty Contract Out satu butir), ditambah Beranda. Empat belas kelompok lain berdiri
 * di sidebar TANPA butir — dan karena itu tidak ada di sini pula. Entri yang
 * berdiri di daftar ini tetapi tidak di sidebar dapat dibuka lewat palet
 * walau menunya tidak terlihat; itu persis cacat yang REFERENSI_UI bayar
 * sekali dan tuliskan pelajarannya.
 *
 * Sejak menu dari tabel M_NAV_MENU (30-09-2026) sidebar dan palet mengikuti
 * urutan TABEL; daftar ini hanya PEMOTONGNYA (`susunMenu`), Beranda pertama.
 */
export const ENTRI_MENU: readonly EntriMenu<Halaman>[] = [
  { modul: HALAMAN_BERANDA, label: 'Beranda', kelompok: 'Beranda', pemilik: null },
  ...MODUL_FRONTEND.flatMap((m) => m.menu.map((b) => ({ ...b, pemilik: m.nama }))),
]

/**
 * Modul backend pemilik tiap halaman - nama yang sama dengan `MODUL_AKTIF` dan
 * `GET /api/modul-aktif` (refactor bentuk B). `null` = milik aplikasi, selalu
 * tampil (Beranda).
 */
export const MODUL_BACKEND: Readonly<Record<Halaman, string | null>> = Object.fromEntries([
  [HALAMAN_BERANDA, null],
  ...MODUL_FRONTEND.flatMap((m) => m.halaman.map((h) => [h, m.nama])),
]) as Record<Halaman, string | null>

/** Apakah halaman ini tampil, menurut daftar modul aktif dari backend. */
export function halamanAktif(halaman: Halaman, aktif: readonly string[] | null): boolean {
  return modulDipasang(MODUL_BACKEND[halaman], aktif)
}
