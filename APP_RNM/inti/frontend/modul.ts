// Bentuk satu modul frontend - refactor bentuk B (30-09-2026).
//
// Padanan `inti.Modul` di backend (`inti/backend/modul.go`): setiap modul
// menyebut NAMA-nya (sama dengan `const Nama` di
// `modul/<nama>/backend/modul.go` dan `MODUL_AKTIF`), halaman-halamannya, butir
// halaman AWAL-nya (`menu.ts`, dibuka tombol modul di sidebar - menu datar,
// keputusan work owner 30-09-2026), dan komponen rutenya (`rute.tsx`).
// `App.tsx` hanya memasang yang AKTIF. Berkas ini tidak mengenal modul mana pun.
//
// # Struktur tim satu folder per modul (30-09-2026)
//
// Daftar modul frontend (`frontend/daftar.ts`) tidak lagi mengimpor tiap modul
// dengan tangan: ia memakai `import.meta.glob` atas
// `modul/*/frontend/menu.ts` dan `rute.tsx`. Karena itu setiap modul
// mengekspor nama yang SAMA - `PENDAFTARAN_MENU` (bentuk `MenuModul`) dan
// `RUTE_MODUL` (bentuk `RuteModul`) - dan memperluas `HalamanModul` di bawah.

import type { ComponentType } from 'react'

import type { Sesi } from './store/sesi'

/**
 * Halaman setiap modul, diperluas SETIAP MODUL di `menu.ts`-nya sendiri lewat
 * declaration merging - bukan didaftar di berkas ini:
 *
 *     declare module '../../../inti/frontend/modul' {
 *       interface HalamanModul {
 *         claimlife: HalamanClaimLife
 *       }
 *     }
 *
 * ⛔ Nama halaman yang salah ketik membuka layar kosong, dan tidak satu pun uji
 * TEKS akan melihatnya. Union `HalamanTerdaftar` menjaga itu di kompiler -
 * dan kini terbentuk dari folder modul, tanpa satu baris per modul di berkas
 * bersama.
 */
export interface HalamanModul {}

/** Setiap halaman modul terdaftar, sebagai UNION - bukan `string`. */
export type HalamanTerdaftar = HalamanModul[keyof HalamanModul]

/** Bentuk ekspor `PENDAFTARAN_MENU` di `modul/<nama>/frontend/menu.ts`. */
export interface MenuModul<H extends string = HalamanTerdaftar> {
  /** Nama modul backend - SAMA dengan nama folder dan `const Nama` Go. */
  nama: string
  /**
   * Nama modul - nama folder korpus VERBATIM (kartu Beranda). Label tombol
   * sidebar datang dari `M_NAV_MENU.LABEL`; penjaga dua arah menuntut keduanya
   * sama.
   */
  kelompok: string
  /** Seluruh halaman modul ini, termasuk yang dibuka DARI DALAM kasus. */
  halaman: readonly H[]
  /**
   * Halaman yang dibuka tombol modul di sidebar dan palet - satu modul, satu
   * menu (keputusan work owner 30-09-2026). Halaman lain modul itu dibuka
   * dari dalam halaman ini.
   */
  halamanAwal: H
}

/** Bentuk ekspor `RUTE_MODUL` di `modul/<nama>/frontend/rute.tsx`. */
export type RuteModul<H extends string = HalamanTerdaftar> = ComponentType<PropsRute<H>>

/**
 * Masukan rute satu modul.
 *
 * `halaman` adalah halaman yang sedang tampil - boleh milik modul lain; rute
 * modul hanya merender halamannya sendiri. `onPindah` memindahkan ke halaman
 * modul itu (atau ke Beranda lewat Shell).
 */
export interface PropsRute<H extends string> {
  halaman: string
  masuk: Sesi
  onPindah: (h: H) => void
  /**
   * Bertambah SETIAP KALI pemakai memilih menu (sidebar atau palet), termasuk
   * menu halaman yang sedang tampil. Rute yang ingin kembali ke layar awalnya
   * saat menunya diklik ulang membandingkannya dengan nilai sebelumnya —
   * `halaman` saja tidak cukup: memilih halaman yang sama tidak mengubahnya.
   *
   * Opsional dan aditif (01-10-2026, permintaan PremiumList Life): rute yang
   * tidak membacanya tetap seperti sebelumnya.
   */
  ketukMenu?: number
}

/** Satu modul frontend terdaftar. */
export interface ModulFrontend<H extends string> {
  /** Nama modul backend - `MODUL_AKTIF`, `GET /api/modul-aktif`. */
  nama: string
  /** Nama modul - nama folder korpus VERBATIM (lihat `MenuModul.kelompok`). */
  kelompok: string
  /** Seluruh halaman modul ini, termasuk yang dibuka DARI DALAM kasus. */
  halaman: readonly H[]
  /** Halaman yang dibuka tombol modul (lihat `MenuModul.halamanAwal`). */
  halamanAwal: H
  /**
   * Komponen rute. ⛔ Ia TETAP terpasang selama modulnya aktif, supaya
   * keadaannya (kasus yang sedang dibuka) bertahan saat pemakai pindah
   * halaman - persis seperti ketika keadaan itu hidup di `App.tsx`.
   */
  Rute: ComponentType<PropsRute<H>>
}
