// Rute modul Master Data - satu halaman menu (`masterdata-daftar`): daftar bertab per master, tambah / ubah, aktif /
// nonaktif (rencana 01 langkah 3, sesi 0f).

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanMasterData as Halaman } from './menu'
import HalamanMasterData from './pages/HalamanMasterData'
import './masterdata.css'

export function RuteMasterData({ halaman }: PropsRute<Halaman>) {
  if (halaman !== 'masterdata-daftar') return null
  return <HalamanMasterData />
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<Halaman> = RuteMasterData
