// Menu modul NB FacIn - tiket 21, halaman depan tiket 25/26. Menu DATAR (keputusan work owner 30-09-2026): satu
// modul satu tombol; label tombol datang dari `M_NAV_MENU.LABEL` barisnya.

import type { MenuModul } from '../../../inti/frontend/modul'
import { kotakMasukNbFacIn } from './api'
import { daftarBeranda } from './beranda'
import { VIEW_POLICY } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `modul/nbfacin/backend/modul.go`. */
export const NAMA_NBFACIN = 'nbfacin'

/** Nama modul - nama folder korpus VERBATIM (`NB FacIn`) = `M_NAV_MENU.LABEL` barisnya (900). */
export const KELOMPOK_NBFACIN = 'NB FacIn'

// 'nbfacin-view-policy' (tiket 71): menu View Policy - tombol sidebar di bawah Home (`menuAtas`), tanpa baris M_NAV_MENU baru.
export const HALAMAN_NBFACIN = ['nbfacin-portal', 'nbfacin-opportunity', 'nbfacin-inward', 'nbfacin-coverage-cargo', 'nbfacin-view-policy'] as const
export type HalamanNbFacIn = (typeof HALAMAN_NBFACIN)[number]

/**
 * Halaman yang dibuka tombol "NB FacIn" di sidebar = portal Opportunity (tiket 25), meniru harness Pega
 * `SFAPortalOpportunities`. Form Opportunity (tiket 26) dibuka dari tombol `Create opportunity`; coverage
 * kargo (tiket 21) bukan lagi halaman awal - work owner 02-10-2026: halaman depan lama salah tempat.
 */
export const HALAMAN_AWAL_NBFACIN: HalamanNbFacIn = 'nbfacin-portal'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanNbFacIn> = {
  nama: NAMA_NBFACIN,
  kelompok: KELOMPOK_NBFACIN,
  halaman: HALAMAN_NBFACIN,
  halamanAwal: HALAMAN_AWAL_NBFACIN,
  // Kotak masuk Beranda per workbasket yang dipegang akun (work owner 07-10-2026; pola NB Treaty In 06-10-2026).
  antreanBeranda: kotakMasukNbFacIn,
  // Daftar case yang menunggu, ID membuka layar case (`PropsRute.bukaKasus`).
  daftarBeranda,
  // Tiket 71 (work owner 09-10-2026 "View Policy dibuatin di menu baru. dibawah menu home").
  menuAtas: [{ label: VIEW_POLICY.judul, halaman: 'nbfacin-view-policy' }],
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    nbfacin: HalamanNbFacIn
  }
}
