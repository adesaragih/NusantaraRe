// Rute modul Aggregate. Pembungkus `.aggregate` (`display: contents`) adalah akar CSS modul ini - setiap pemilih
// `aggregate.css` berada di bawahnya (`gaya.test.ts`).

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './aggregate.css'
import type { HalamanAG } from './menu'
import Aggregate from './pages/Aggregate'

export function RuteAG({ halaman }: PropsRute<HalamanAG>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="aggregate">{halaman === 'aggregate-daftar' && <Aggregate />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanAG> = RuteAG
