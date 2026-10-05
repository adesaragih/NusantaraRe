// Rute modul Company Detail. Pembungkus `.companydetail` (`display: contents`) adalah akar CSS modul ini - setiap
// pemilih `companydetail.css` berada di bawahnya (`gaya.test.ts`).

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './companydetail.css'
import type { HalamanCD } from './menu'
import CompanyDetail from './pages/CompanyDetail'

export function RuteCD({ halaman }: PropsRute<HalamanCD>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="companydetail">{halaman === 'companydetail-daftar' && <CompanyDetail />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanCD> = RuteCD
