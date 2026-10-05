// Rute modul Adjuster Consultant untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './adjusterconsultant.css'
import type { HalamanADJ } from './menu'
import AdjusterConsultant from './pages/AdjusterConsultant'

export function RuteADJ({ halaman }: PropsRute<HalamanADJ>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="adjusterconsultant">{halaman === 'adjusterconsultant-daftar' && <AdjusterConsultant />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanADJ> = RuteADJ
