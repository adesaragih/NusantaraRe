// Rute modul Plan untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './planlife.css'
import type { HalamanPL } from './menu'
import PlanLife from './pages/PlanLife'

export function RutePL({ halaman }: PropsRute<HalamanPL>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="planlife">{halaman === 'planlife-daftar' && <PlanLife />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanPL> = RutePL
