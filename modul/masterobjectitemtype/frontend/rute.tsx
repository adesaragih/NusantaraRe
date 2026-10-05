// Rute modul Object Item Type - satu halaman: layar master generik inti atas prefix rute modul ini.

import HalamanMaster from '../../../inti/frontend/master/HalamanMaster'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { PREFIX_MASTEROBJECTITEMTYPE } from './labels'
import type { HalamanObjectItemType as Halaman } from './menu'

export function RuteObjectItemType({ halaman }: PropsRute<Halaman>) {
  if (halaman !== 'masterobjectitemtype-daftar') return null
  return <HalamanMaster prefix={PREFIX_MASTEROBJECTITEMTYPE} />
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<Halaman> = RuteObjectItemType
