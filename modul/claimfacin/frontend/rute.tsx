// Rute modul Claim Fac In untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './claimfacin.css'
import type { HalamanCFI } from './menu'
import ClaimFacIn from './pages/ClaimFacIn'

export function RuteCFI({ halaman, masuk, onBukaModul, bukaKasus, onBeranda }: PropsRute<HalamanCFI>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="claimfacin">
        {halaman === 'claimfacin-daftar' && (
          <ClaimFacIn pelaku={masuk.akunID} onBukaModul={onBukaModul} bukaKasus={bukaKasus} onBeranda={onBeranda} />
        )}
      </div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanCFI> = RuteCFI
