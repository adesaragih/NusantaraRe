// Rute modul Bordereaux. Pembungkus `.bordereaux` (`display: contents`) adalah akar CSS modul ini - setiap pemilih
// `bordereaux.css` berada di bawahnya (`gaya.test.ts`).

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './bordereaux.css'
import type { HalamanBDX } from './menu'
import Bordereaux from './pages/Bordereaux'

export function RuteBDX({ halaman }: PropsRute<HalamanBDX>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="bordereaux">{halaman === 'bordereaux-daftar' && <Bordereaux />}</div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanBDX> = RuteBDX
