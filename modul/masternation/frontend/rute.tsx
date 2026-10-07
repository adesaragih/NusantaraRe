// Rute modul Nation - satu halaman: layar master generik inti atas prefix rute modul ini.

import HalamanMaster from '../../../inti/frontend/master/HalamanMaster'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { PREFIX_MASTERNATION } from './labels'
import type { HalamanNation as Halaman } from './menu'

export function RuteNation({ halaman }: PropsRute<Halaman>) {
  if (halaman !== 'masternation-daftar') return null
  return <HalamanMaster prefix={PREFIX_MASTERNATION} />
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<Halaman> = RuteNation
