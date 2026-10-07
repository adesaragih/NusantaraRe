// Rute modul Province - satu halaman: layar master generik inti atas prefix rute modul ini.

import HalamanMaster from '../../../inti/frontend/master/HalamanMaster'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { PREFIX_MASTERPROVINCE } from './labels'
import type { HalamanProvince as Halaman } from './menu'

export function RuteProvince({ halaman }: PropsRute<Halaman>) {
  if (halaman !== 'masterprovince-daftar') return null
  return <HalamanMaster prefix={PREFIX_MASTERPROVINCE} />
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<Halaman> = RuteProvince
