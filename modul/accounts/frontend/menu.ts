import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_ACC } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `backend/modul.go` dan nama foldernya. */
export const NAMA_ACC = 'accounts'
export const HALAMAN_ACC = ['accounts-daftar'] as const
export type HalamanACC = (typeof HALAMAN_ACC)[number]

/** Halaman yang dibuka tombol "Accounts" di sidebar. */
export const HALAMAN_AWAL_ACC: HalamanACC = 'accounts-daftar'

export const PENDAFTARAN_MENU: MenuModul<HalamanACC> = {
  nama: NAMA_ACC,
  kelompok: MENU_ACC.kelompok,
  halaman: HALAMAN_ACC,
  halamanAwal: HALAMAN_AWAL_ACC,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    accounts: HalamanACC
  }
}
