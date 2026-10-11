// Rute modul Claim Non Prop untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './claimnonprop.css'
import type { HalamanCNP } from './menu'
import ClaimNonProp from './pages/ClaimNonProp'

export function RuteCNP({ halaman, masuk, onLihatBerkas, onBukaModul, bukaKasus, onBeranda }: PropsRute<HalamanCNP>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="claimnonprop">
        {halaman === 'claimnonprop-daftar' && (
          <ClaimNonProp
            pelaku={masuk.akunID}
            onLihatBerkas={onLihatBerkas}
            onBukaModul={onBukaModul}
            bukaKasus={bukaKasus}
            onBeranda={onBeranda}
          />
        )}
      </div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanCNP> = RuteCNP
