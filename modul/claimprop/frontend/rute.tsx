// Rute modul Claim Prop untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './claimprop.css'
import type { HalamanCP } from './menu'
import ClaimProp from './pages/ClaimProp'

export function RuteCP({ halaman, masuk, onLihatBerkas, onBukaModul, bukaKasus, onBeranda }: PropsRute<HalamanCP>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="claimprop">
        {halaman === 'claimprop-daftar' && (
          <ClaimProp
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

export const RUTE_MODUL: RuteModul<HalamanCP> = RuteCP
