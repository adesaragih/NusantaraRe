// Rute modul Cause Of Loss Life untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './causeoflosslife.css'
import type { HalamanCOL } from './menu'
import CauseOfLossLife from './pages/CauseOfLossLife'

export function RuteCOL({ halaman }: PropsRute<HalamanCOL>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="causeoflosslife">{halaman === 'causeoflosslife-daftar' && <CauseOfLossLife />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanCOL> = RuteCOL
