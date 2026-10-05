// Rute modul Treaty Group untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './treatygroup.css'
import type { HalamanTG } from './menu'
import TreatyGroup from './pages/TreatyGroup'

export function RuteTG({ halaman }: PropsRute<HalamanTG>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="treatygroup">{halaman === 'treatygroup-daftar' && <TreatyGroup />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanTG> = RuteTG
