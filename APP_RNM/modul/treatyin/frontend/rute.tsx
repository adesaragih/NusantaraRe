// Rute modul Treaty In — tiket 14 dan 15.

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanTreatyIn } from './menu'
import AcuanTreatyIn from './pages/AcuanTreatyIn'

export function RuteTreatyIn({ halaman }: PropsRute<HalamanTreatyIn>) {
  return <>{halaman === 'treatyin-acuan' && <AcuanTreatyIn />}</>
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanTreatyIn> = RuteTreatyIn
