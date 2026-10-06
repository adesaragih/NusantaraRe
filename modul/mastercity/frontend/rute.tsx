// Rute modul City - satu halaman: layar master generik inti atas prefix rute modul ini.

import HalamanMaster from '../../../inti/frontend/master/HalamanMaster'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { PREFIX_MASTERCITY } from './labels'
import type { HalamanCity as Halaman } from './menu'

export function RuteCity({ halaman }: PropsRute<Halaman>) {
  if (halaman !== 'mastercity-daftar') return null
  return <HalamanMaster prefix={PREFIX_MASTERCITY} />
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<Halaman> = RuteCity
