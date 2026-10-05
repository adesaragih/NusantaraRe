// Rute modul Master Contract Retro Life - paket 10.
//
// Label korpus berbahasa Inggris; `BahasaUI` "en" membuat teks bawaan komponen bersama (Memuat,
// penomoran halaman, panel backend mati, pilihan kosong, tutup popup) ikut Inggris di layar ini saja.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanMCRL } from './menu'
import MasterContractRetroLife from './pages/MasterContractRetroLife'
import { catatOperator } from './tampilan'

// Setiap pilihan menu (`ketukMenu` naik) memasang ULANG halaman lewat `key`: form,
// panel, popup, dan halaman tabel kembali ke keadaan awal (permintaan work owner
// 03-10-2026: "setiap kali menunya ditekan tampilan langsung balik ke awal").
export function RuteMCRL({ halaman, ketukMenu, masuk }: PropsRute<HalamanMCRL>) {
  // Inputor = akun yang login, di form Add maupun Edit (04-10-2026).
  catatOperator(masuk.akunID)
  return (
    <BahasaUI.Provider value="en">
      {halaman === 'mcrl-tahun' && <MasterContractRetroLife key={ketukMenu ?? 0} />}
    </BahasaUI.Provider>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanMCRL> = RuteMCRL
