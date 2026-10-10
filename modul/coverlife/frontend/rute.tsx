// Rute modul Cover Life untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './coverlife.css'
import type { HalamanCVL } from './menu'
import CoverLife from './pages/CoverLife'

export function RuteCVL({ halaman }: PropsRute<HalamanCVL>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="coverlife">{halaman === 'coverlife-daftar' && <CoverLife />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanCVL> = RuteCVL
