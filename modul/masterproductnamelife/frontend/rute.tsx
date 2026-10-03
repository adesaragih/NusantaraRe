// Rute modul Master Product Name Life - paket 10.
//
// Label korpus berbahasa Inggris; `BahasaUI` "en" membuat teks bawaan komponen bersama (Memuat, penomoran
// halaman, panel backend mati, pilihan kosong, tutup popup) ikut Inggris di layar ini saja.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './masterproductnamelife.css'
import type { HalamanMPNL } from './menu'
import MasterProductNameLife from './pages/MasterProductNameLife'

export function RuteMPNL({ halaman }: PropsRute<HalamanMPNL>) {
  return <BahasaUI.Provider value="en">{halaman === 'mpnl-produk' && <MasterProductNameLife />}</BahasaUI.Provider>
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanMPNL> = RuteMPNL
