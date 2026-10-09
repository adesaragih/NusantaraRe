// Rute modul R/I Risk untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './ririsklife.css'
import type { HalamanRK } from './menu'
import RIRiskLife from './pages/RIRiskLife'

export function RuteRK({ halaman }: PropsRute<HalamanRK>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="ririsklife">{halaman === 'ririsklife-daftar' && <RIRiskLife />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanRK> = RuteRK
