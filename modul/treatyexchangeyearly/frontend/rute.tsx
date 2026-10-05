// Rute modul Treaty Exchange Yearly untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './treatyexchangeyearly.css'
import type { HalamanTEY } from './menu'
import TreatyExchangeYearly from './pages/TreatyExchangeYearly'

export function RuteTEY({ halaman }: PropsRute<HalamanTEY>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="treatyexchangeyearly">
        {halaman === 'treatyexchangeyearly-daftar' && <TreatyExchangeYearly />}
      </div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanTEY> = RuteTEY
