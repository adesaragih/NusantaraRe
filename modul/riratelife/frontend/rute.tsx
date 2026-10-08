// Rute modul R/I Rate Life untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './riratelife.css'
import type { HalamanRR } from './menu'
import RIRateLife from './pages/RIRateLife'

export function RuteRR({ halaman }: PropsRute<HalamanRR>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="riratelife">{halaman === 'riratelife-daftar' && <RIRateLife />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanRR> = RuteRR
