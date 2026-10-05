// Rute modul Accumulation - satu halaman: layar master generik inti atas prefix rute modul ini.

import HalamanMaster from '../../../inti/frontend/master/HalamanMaster'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { PREFIX_MASTERACCUMULATION } from './labels'
import type { HalamanAccumulation as Halaman } from './menu'

export function RuteAccumulation({ halaman }: PropsRute<Halaman>) {
  if (halaman !== 'masteraccumulation-daftar') return null
  return <HalamanMaster prefix={PREFIX_MASTERACCUMULATION} />
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<Halaman> = RuteAccumulation
