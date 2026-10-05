// Rute modul Treaty Group OJK untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './treatygroupojk.css'
import type { HalamanTGO } from './menu'
import TreatyGroupOjk from './pages/TreatyGroupOjk'

export function RuteTGO({ halaman }: PropsRute<HalamanTGO>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="treatygroupojk">{halaman === 'treatygroupojk-daftar' && <TreatyGroupOjk />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanTGO> = RuteTGO
