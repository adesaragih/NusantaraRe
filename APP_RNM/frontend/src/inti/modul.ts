// Bentuk satu modul frontend - refactor bentuk B (30-09-2026).
//
// Padanan `inti.Modul` di backend (`inti/modul.go`): setiap modul menyebut
// NAMA-nya (sama dengan `const Nama` di `modul/<nama>/modul.go` dan
// `MODUL_AKTIF`), halaman-halamannya, butir menunya (`menu.ts`), dan
// komponen rutenya (`rute.tsx`). `modul/daftar.ts` mendaftarkannya; `App.tsx`
// hanya memasang yang AKTIF. Berkas ini tidak mengenal modul mana pun.

import type { ComponentType } from 'react'

import type { ButirMenuModul } from './lib/daftarMenu'
import type { Sesi } from './store/sesi'

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
}

/** Satu modul frontend terdaftar. */
export interface ModulFrontend<H extends string> {
  /** Nama modul backend - `MODUL_AKTIF`, `GET /api/modul-aktif`. */
  nama: string
  /** Seluruh halaman modul ini, termasuk yang dibuka DARI DALAM kasus. */
  halaman: readonly H[]
  /** Butir menu sidebar dan palet, berurutan seperti tampil. */
  menu: readonly ButirMenuModul<H>[]
  /**
   * Komponen rute. ⛔ Ia TETAP terpasang selama modulnya aktif, supaya
   * keadaannya (kasus yang sedang dibuka) bertahan saat pemakai pindah
   * halaman - persis seperti ketika keadaan itu hidup di `App.tsx`.
   */
  Rute: ComponentType<PropsRute<H>>
}
