// Rute modul Master Data - satu halaman menu (`masterdata-daftar`).
//
// ⚠️ KERANGKA (sesi c3, 04-10-2026): halaman placeholder supaya modul terpasang dan uji perakit hijau; layar daftar /
// form / aktif-nonaktif dibangun sesi 0f (rencana langkah 3) dan menggantikan `HalamanMasterData`.

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanMasterData as Halaman } from './menu'
import HalamanMasterData from './pages/HalamanMasterData'

export function RuteMasterData({ halaman }: PropsRute<Halaman>) {
  if (halaman !== 'masterdata-daftar') return null
  return <HalamanMasterData />
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<Halaman> = RuteMasterData
