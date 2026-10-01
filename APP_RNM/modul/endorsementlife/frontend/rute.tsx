// Rute modul Endorsement Life.
//
// Label korpus berbahasa Inggris; `BahasaUI` "en" membuat teks bawaan komponen bersama (Memuat,
// penomoran halaman, panel backend mati) ikut Inggris di layar ini saja.
//
// ⛔ Satu halaman menu (`edm-inbox`); kasus yang sedang dibuka disimpan di keadaan rute, supaya ia
// bertahan saat pemakai pindah halaman lalu kembali (pola PremiumList Life).

import { useState } from 'react'

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanEDM } from './menu'
import InboxEndorsementLife from './pages/InboxEndorsementLife'
import InputEDMLife from './pages/InputEDMLife'

export function RuteEDM({ halaman }: PropsRute<HalamanEDM>) {
  const [kasusId, setKasusId] = useState('')
  if (halaman !== 'edm-inbox') return null
  return (
    <BahasaUI.Provider value="en">
      {kasusId === '' && <InboxEndorsementLife onBuka={setKasusId} />}
      {kasusId !== '' && <InputEDMLife kasusId={kasusId} onTutup={() => setKasusId('')} />}
    </BahasaUI.Provider>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanEDM> = RuteEDM
