// Rute modul Reinsurance Type untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './reinsurancetype.css'
import type { HalamanRT } from './menu'
import ReinsuranceType from './pages/ReinsuranceType'

export function RuteRT({ halaman }: PropsRute<HalamanRT>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="reinsurancetype">{halaman === 'reinsurancetype-daftar' && <ReinsuranceType />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanRT> = RuteRT
