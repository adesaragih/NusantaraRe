// Rute modul Treaty Description untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './treatydescription.css'
import type { HalamanTD } from './menu'
import TreatyDescription from './pages/TreatyDescription'

export function RuteTD({ halaman }: PropsRute<HalamanTD>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="treatydescription">{halaman === 'treatydescription-daftar' && <TreatyDescription />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanTD> = RuteTD
