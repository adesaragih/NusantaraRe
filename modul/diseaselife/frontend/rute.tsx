// Rute modul Disease Life untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './diseaselife.css'
import type { HalamanDSL } from './menu'
import DiseaseLife from './pages/DiseaseLife'

export function RuteDSL({ halaman }: PropsRute<HalamanDSL>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="diseaselife">{halaman === 'diseaselife-daftar' && <DiseaseLife />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanDSL> = RuteDSL
