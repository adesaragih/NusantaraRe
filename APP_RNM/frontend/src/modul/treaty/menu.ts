// Menu modul Treaty Contract Out - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `ENTRI_MENU` di `lib/daftarMenu.ts`.

import { MODUL } from '../../inti/labels'
import type { ButirMenuModul } from '../../inti/lib/daftarMenu'
import { MENU_TCO } from './labels'

/** Nama modul - SAMA dengan `const Nama` di `modul/treaty/modul.go`. */
export const NAMA_TREATY = 'treaty'

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
  { modul: 'tco-tahun', label: MENU_TCO.treatyContractOut, kelompok: MODUL.treatyContractOut },
]
