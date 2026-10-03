// Rute modul Marketing Officer. Pembungkus `.marketingofficer` (`display: contents`) adalah akar CSS modul ini -
// setiap pemilih `marketingofficer.css` berada di bawahnya (`gaya.test.ts`).

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './marketingofficer.css'
import type { HalamanMO } from './menu'
import MarketingOfficer from './pages/MarketingOfficer'

export function RuteMO({ halaman }: PropsRute<HalamanMO>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="marketingofficer">{halaman === 'marketingofficer-daftar' && <MarketingOfficer />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanMO> = RuteMO
