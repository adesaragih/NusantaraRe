// Menu modul Treaty Contract Out - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `ENTRI_MENU` di `lib/daftarMenu.ts`.
//
// Struktur tim satu folder per modul (30-09-2026): nama modul tinggal DI SINI,
// bukan di `inti/frontend/labels.ts`; dan `PENDAFTARAN_MENU` di bawah dibaca
// perakit `frontend/daftar.ts` lewat `import.meta.glob` - nol baris per modul
// di berkas bersama.
//
// Menu DATAR (keputusan work owner 30-09-2026): satu modul satu menu - tombol
// "Treaty Contract Out" membuka halaman awal di bawah; butir `tco-tahun`
// dicabut. Sebelumnya (tco5) modul ini sudah satu tombol datar; kini SEMUA
// modul begitu.

import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_TCO } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `modul/treatycontractout/backend/modul.go`. */
export const NAMA_TREATY = 'treatycontractout'

/** Nama modul - nama folder korpus VERBATIM (`labels.ts` `MENU_TCO.kelompok`). */
export const KELOMPOK_TREATY = MENU_TCO.kelompok

/**
 * Halaman modul ini.
 *
 * ⚠️ tco5: `tco-kontrak` / `tco-klausul` dibuka tombol `ReinsType` /
 * `List Description` layar tahun treaty. Nilainya tetap di sini karena
 * rutenya ada (`rute.tsx`).
 */
export const HALAMAN_TREATY = ['tco-tahun', 'tco-kontrak', 'tco-klausul'] as const
export type HalamanTreaty = (typeof HALAMAN_TREATY)[number]

/** Halaman yang dibuka tombol "Treaty Contract Out" - layar tahun treaty. */
export const HALAMAN_AWAL_TREATY: HalamanTreaty = 'tco-tahun'

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanTreaty> = {
  nama: NAMA_TREATY,
  kelompok: KELOMPOK_TREATY,
  halaman: HALAMAN_TREATY,
  halamanAwal: HALAMAN_AWAL_TREATY,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    treatycontractout: HalamanTreaty
  }
}
