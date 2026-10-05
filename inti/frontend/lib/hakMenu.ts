// Hak menu per akun (`M_LOGIN_GO_MENU.HAK`, migrasi 914, keputusan work owner 04-10-2026): menu yang hanya boleh
// DILIHAT menyembunyikan tombol tulis layarnya. Backend tetap menolak tulis (gerbang `cmd/api`), jadi ini kenyamanan
// layar, bukan pengamannya.
//
// Disediakan `App.tsx` lewat `KonteksHakMenu`; modul membaca `useBolehUbah('<kode menu>')`. Tanpa penyedia (mode
// stub, uji) = penuh.

import { createContext, useContext } from 'react'

export interface HakMenu {
  /** KODE menu ber-hak LIHAT - berlaku juga bagi superadmin (keputusan work owner 05-10-2026). */
  lihat: readonly string[]
}

export const HAK_PENUH: HakMenu = { lihat: [] }

export const KonteksHakMenu = createContext<HakMenu>(HAK_PENUH)

/** Boleh menulis lewat menu `kode`? */
export function bolehUbahMenu(h: HakMenu, kode: string): boolean {
  return !h.lihat.includes(kode)
}

/** Hook layar modul: false = menu ini View only bagi akun yang login. */
export function useBolehUbah(kode: string): boolean {
  return bolehUbahMenu(useContext(KonteksHakMenu), kode)
}
