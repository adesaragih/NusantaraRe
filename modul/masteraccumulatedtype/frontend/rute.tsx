// Rute modul Accumulated Type - satu halaman: layar master generik inti atas prefix rute modul ini.

import HalamanMaster from '../../../inti/frontend/master/HalamanMaster'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { PREFIX_MASTERACCUMULATEDTYPE } from './labels'
import type { HalamanAccumulatedType as Halaman } from './menu'

export function RuteAccumulatedType({ halaman }: PropsRute<Halaman>) {
  if (halaman !== 'masteraccumulatedtype-daftar') return null
  return <HalamanMaster prefix={PREFIX_MASTERACCUMULATEDTYPE} />
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<Halaman> = RuteAccumulatedType
