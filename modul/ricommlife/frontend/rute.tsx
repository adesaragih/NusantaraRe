// Rute modul R/I Comm Life untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './ricommlife.css'
import type { HalamanRC } from './menu'
import RICommLife from './pages/RICommLife'

export function RuteRC({ halaman }: PropsRute<HalamanRC>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="ricommlife">{halaman === 'ricommlife-daftar' && <RICommLife />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanRC> = RuteRC
