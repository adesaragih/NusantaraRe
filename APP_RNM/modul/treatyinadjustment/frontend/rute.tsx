// Rute modul Treaty In Adjustment — tiket 01 dan 05.

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanTreatyInAdjustment } from './menu'
import RantaiVersi from './pages/RantaiVersi'

export function RuteTreatyInAdjustment({ halaman }: PropsRute<HalamanTreatyInAdjustment>) {
  return <>{halaman === 'treatyinadjustment-rantai-versi' && <RantaiVersi />}</>
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanTreatyInAdjustment> = RuteTreatyInAdjustment
