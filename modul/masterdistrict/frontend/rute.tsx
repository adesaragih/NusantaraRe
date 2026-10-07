// Rute modul District - satu halaman: layar master generik inti atas prefix rute modul ini.

import HalamanMaster from '../../../inti/frontend/master/HalamanMaster'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import { PREFIX_MASTERDISTRICT } from './labels'
import type { HalamanDistrict as Halaman } from './menu'

export function RuteDistrict({ halaman }: PropsRute<Halaman>) {
  if (halaman !== 'masterdistrict-daftar') return null
  return <HalamanMaster prefix={PREFIX_MASTERDISTRICT} />
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<Halaman> = RuteDistrict
