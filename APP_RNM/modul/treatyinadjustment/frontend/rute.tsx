// Rute modul Treaty In Adjustment — tiket 01 dan 05.

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanTreatyInAdjustment } from './menu'
import RantaiVersi from './pages/RantaiVersi'
import './treatyinadjustment.css'

export function RuteTreatyInAdjustment({ halaman }: PropsRute<HalamanTreatyInAdjustment>) {
  return (
    // Akar gaya modul: semua aturan `treatyinadjustment.css` diawali `.treatyinadjustment`
    // (`display: contents`).
    <div className="treatyinadjustment">
      {halaman === 'treatyinadjustment-rantai-versi' && <RantaiVersi />}
    </div>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanTreatyInAdjustment> = RuteTreatyInAdjustment
