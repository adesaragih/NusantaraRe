// Menu modul Treaty Contract Out - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `ENTRI_MENU` di `lib/daftarMenu.ts`.
//
// Struktur tim satu folder per modul (30-09-2026): nama kelompok sidebar modul
// ini tinggal DI SINI, bukan di `inti/frontend/labels.ts`; dan
// `PENDAFTARAN_MENU` di bawah dibaca perakit `frontend/daftar.ts` lewat
// `import.meta.glob` - nol baris per modul di berkas bersama.

import type { ButirMenuModul } from '../../../inti/frontend/lib/daftarMenu'
import type { MenuModul } from '../../../inti/frontend/modul'
import { MENU_TCO } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `modul/treatycontractout/backend/modul.go`. */
export const NAMA_TREATY = 'treatycontractout'

/** Nama kelompok sidebar - nama folder korpus VERBATIM (`labels.ts` `MENU_TCO.kelompok`). */
export const KELOMPOK_TREATY = MENU_TCO.kelompok

/**
 * Halaman modul ini.
 *
 * ⚠️ tco5: `tco-kontrak` / `tco-klausul` TIDAK lagi punya butir menu — kedua
 * layar dibuka tombol `ReinsType` / `List Description` layar tahun treaty.
 * Nilainya tetap di sini karena rutenya ada (`rute.tsx`).
 */
export const HALAMAN_TREATY = ['tco-tahun', 'tco-kontrak', 'tco-klausul'] as const
export type HalamanTreaty = (typeof HALAMAN_TREATY)[number]

export const MENU_TREATY: readonly ButirMenuModul<HalamanTreaty>[] = [
  // Treaty Contract Out — tco5 [keputusan work owner 29-09-2026]: SATU butir
  // "Treaty Contract Out" (asal harness InboxTreatyContract). Butir ReinsType
  // dan Description dibuang: di Pega keduanya popup form kontrak (b20778, b22196).
  // [keputusan work owner 30-09-2026] DATAR: satu tombol langsung, tanpa model
  // kelompok-beranak "Treaty Contract Out ▸ Treaty Contract Out".
  { modul: 'tco-tahun', label: MENU_TCO.treatyContractOut, kelompok: KELOMPOK_TREATY, datar: true },
]

/** Menu modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const PENDAFTARAN_MENU: MenuModul<HalamanTreaty> = {
  nama: NAMA_TREATY,
  kelompok: KELOMPOK_TREATY,
  halaman: HALAMAN_TREATY,
  menu: MENU_TREATY,
}

declare module '../../../inti/frontend/modul' {
  interface HalamanModul {
    treatycontractout: HalamanTreaty
  }
}
