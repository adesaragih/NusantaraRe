// Menu modul Endorsment Fac In - irisan I0 (keputusan work owner 11-10-2026, W-EDM-01). Menu DATAR (keputusan work
// owner 30-09-2026): satu modul satu tombol; label tombol datang dari `M_NAV_MENU.LABEL` barisnya (900), dinyalakan
// slot menu `966_menu_endorsmentfacin.sql`.
//
// `PENDAFTARAN_MENU` dibaca perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul di berkas
// bersama.

import type { MenuModul } from '../../../inti/frontend/modul'

/** Nama modul - SAMA dengan `const Nama` di `modul/endorsmentfacin/backend/modul.go` dan nama foldernya. */
export const NAMA_EDMFACIN = 'endorsmentfacin'

/** Nama modul - nama folder korpus VERBATIM (`Endorsment Fac In`) = `M_NAV_MENU.LABEL` barisnya (900). */
export const KELOMPOK_EDMFACIN = 'Endorsment Fac In'

/**
 * `endorsmentfacin-portal` - harness Pega `SFAPortalEndorsement` (portal endorsement);
 * `endorsmentfacin-addendum` - form *Create Addendum* (`Section/WorkPrimaryDetails.xml`, kelas
 * `ASM-SFAGIS-Work-Endorsement`), dibuka dari tombol Create Addendum portal.
 */
export const HALAMAN_EDMFACIN = ['endorsmentfacin-portal', 'endorsmentfacin-addendum'] as const
export type HalamanEdmFacIn = (typeof HALAMAN_EDMFACIN)[number]

/** Halaman yang dibuka tombol "Endorsment Fac In" di sidebar = portal endorsement. */
export const HALAMAN_AWAL_EDMFACIN: HalamanEdmFacIn = 'endorsmentfacin-portal'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanEdmFacIn> = {
  nama: NAMA_EDMFACIN,
  kelompok: KELOMPOK_EDMFACIN,
  halaman: HALAMAN_EDMFACIN,
  halamanAwal: HALAMAN_AWAL_EDMFACIN,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    endorsmentfacin: HalamanEdmFacIn
  }
}
