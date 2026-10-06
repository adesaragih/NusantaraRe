// Rute modul Business Group untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './businessgroup.css'
import type { HalamanBG } from './menu'
import BusinessGroup from './pages/BusinessGroup'

export function RuteBG({ halaman }: PropsRute<HalamanBG>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="businessgroup">{halaman === 'businessgroup-daftar' && <BusinessGroup />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanBG> = RuteBG
