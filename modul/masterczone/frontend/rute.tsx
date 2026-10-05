// Rute modul CZone - satu halaman: layar master generik inti atas prefix rute modul ini.

import HalamanMaster from '../../../inti/frontend/master/HalamanMaster'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { PREFIX_MASTERCZONE } from './labels'
import type { HalamanCZone as Halaman } from './menu'

export function RuteCZone({ halaman }: PropsRute<Halaman>) {
  if (halaman !== 'masterczone-daftar') return null
  return <HalamanMaster prefix={PREFIX_MASTERCZONE} />
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<Halaman> = RuteCZone
