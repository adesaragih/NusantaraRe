// Rute modul Accounts. Pembungkus `.accounts` (`display: contents`) adalah akar CSS modul ini - setiap pemilih
// `accounts.css` berada di bawahnya (`gaya.test.ts`).

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './accounts.css'
import type { HalamanACC } from './menu'
import Accounts from './pages/Accounts'

export function RuteACC({ halaman }: PropsRute<HalamanACC>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="accounts">{halaman === 'accounts-daftar' && <Accounts />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanACC> = RuteACC
