// Rute modul Komite Claim Prop untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './komiteclaimprop.css'
import type { HalamanKCP } from './menu'
import DaftarKerja from './pages/DaftarKerja'

export function RuteKCP({ halaman, onLihatBerkas }: PropsRute<HalamanKCP>) {
  return (
    <BahasaUI.Provider value="en">
      <div className="komiteclaimprop">
        {halaman === 'komiteclaimprop-daftar' && <DaftarKerja onLihatBerkas={onLihatBerkas} />}
      </div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanKCP> = RuteKCP
