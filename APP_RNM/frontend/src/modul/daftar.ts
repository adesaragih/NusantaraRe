// DAFTAR modul frontend - satu-satunya berkas yang mengenal semua modul
// sekaligus. Padanan `modul/daftar.go` di backend (refactor bentuk B,
// 30-09-2026).
//
// Modul di `modul/<nama>/` tidak pernah mengimpor modul lain, dan `inti/`
// tidak mengenal modul mana pun; berkas ini yang merakit menu dan rute dari
// `menu.ts` / `rute.tsx` tiap modul, dan `App.tsx` yang memasang modul yang
// AKTIF (`GET /api/modul-aktif`).
//
// Menambah modul baru = `modul/<nama>/menu.ts` + `rute.tsx`, lalu satu baris
// di `MODUL_FRONTEND`. Nama modulnya SAMA dengan `const Nama` di
// `modul/<nama>/modul.go` - dijaga `daftar.modulAktif.test.ts`.

import { HALAMAN_BERANDA, modulDipasang, type EntriMenu } from '../inti/lib/daftarMenu'
import type { ModulFrontend } from '../inti/modul'
import { HALAMAN_CLAIMLIFE, MENU_CLAIMLIFE, NAMA_CLAIMLIFE, type HalamanClaimLife } from './claim-life/menu'
import { RuteClaimLife } from './claim-life/rute'
import { HALAMAN_KOMITE, MENU_KOMITE, NAMA_KOMITE, type HalamanKomite } from './komite-claim-life/menu'
import { RuteKomite } from './komite-claim-life/rute'
import {
  HALAMAN_PREMIUMLIST,
  MENU_PREMIUMLIST,
  NAMA_PREMIUMLIST,
  type HalamanPremiumList,
} from './premiumlist-life/menu'
import { RutePremiumList } from './premiumlist-life/rute'
import { HALAMAN_TREATY, MENU_TREATY, NAMA_TREATY, type HalamanTreaty } from './treaty-contract-out/menu'
import { RuteTreaty } from './treaty-contract-out/rute'

/**
 * Halaman yang aplikasi dapat tampilkan, sebagai UNION — bukan `string`.
 *
 * ⛔ Nama halaman yang salah ketik membuka layar kosong, dan tidak satu pun
 * uji TEKS akan melihatnya: `'inbok'` tetap cocok dengan setiap pola yang
 * memeriksa bentuk. Union memindahkan penjagaannya ke kompiler — lebih
 * awal, dan tanpa pagar tambahan.
 */
export type Halaman = typeof HALAMAN_BERANDA | HalamanClaimLife | HalamanPremiumList | HalamanKomite | HalamanTreaty

/**
 * Modul terdaftar, berurutan seperti sidebar - dan seperti `modul.Rakit` di
 * backend, urutan `GET /api/modul-aktif`.
 */
export const MODUL_FRONTEND: readonly ModulFrontend<Halaman>[] = [
  { nama: NAMA_CLAIMLIFE, halaman: HALAMAN_CLAIMLIFE, menu: MENU_CLAIMLIFE, Rute: RuteClaimLife },
  { nama: NAMA_PREMIUMLIST, halaman: HALAMAN_PREMIUMLIST, menu: MENU_PREMIUMLIST, Rute: RutePremiumList },
  { nama: NAMA_KOMITE, halaman: HALAMAN_KOMITE, menu: MENU_KOMITE, Rute: RuteKomite },
  { nama: NAMA_TREATY, halaman: HALAMAN_TREATY, menu: MENU_TREATY, Rute: RuteTreaty },
]

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
