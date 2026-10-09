// Rute modul Benefit untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './benefitlife.css'
import type { HalamanBN } from './menu'
import BenefitLife from './pages/BenefitLife'

export function RuteBN({ halaman }: PropsRute<HalamanBN>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="benefitlife">{halaman === 'benefitlife-daftar' && <BenefitLife />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanBN> = RuteBN
