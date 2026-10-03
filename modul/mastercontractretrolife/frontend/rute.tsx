// Rute modul Master Contract Retro Life - paket 10.
//
// Label korpus berbahasa Inggris; `BahasaUI` "en" membuat teks bawaan komponen bersama (Memuat,
// penomoran halaman, panel backend mati, pilihan kosong, tutup popup) ikut Inggris di layar ini saja.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanMCRL } from './menu'
import MasterContractRetroLife from './pages/MasterContractRetroLife'

export function RuteMCRL({ halaman }: PropsRute<HalamanMCRL>) {
  return <BahasaUI.Provider value="en">{halaman === 'mcrl-tahun' && <MasterContractRetroLife />}</BahasaUI.Provider>
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanMCRL> = RuteMCRL
