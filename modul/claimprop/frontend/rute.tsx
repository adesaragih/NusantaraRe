// Rute modul Claim Prop untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './claimprop.css'
import type { HalamanCP } from './menu'
import ClaimProp from './pages/ClaimProp'

export function RuteCP({ halaman, masuk }: PropsRute<HalamanCP>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="claimprop">{halaman === 'claimprop-daftar' && <ClaimProp pelaku={masuk.akunID} />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanCP> = RuteCP
