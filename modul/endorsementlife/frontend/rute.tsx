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
import BuatEndorsement from './pages/BuatEndorsement'
import InboxEndorsementLife from './pages/InboxEndorsementLife'
import InputEDMLife from './pages/InputEDMLife'

/** Layar yang tampil di halaman menu: kotak masuk, harness buat endorsement, atau kasus. */
type Layar = { jenis: 'inbox' } | { jenis: 'buat' } | { jenis: 'kasus'; id: string }

export function RuteEDM({ halaman }: PropsRute<HalamanEDM>) {
  const [layar, setLayar] = useState<Layar>({ jenis: 'inbox' })
  if (halaman !== 'edm-inbox') return null
  const keInbox = () => setLayar({ jenis: 'inbox' })
  const keKasus = (id: string) => setLayar({ jenis: 'kasus', id })
  return (
    <BahasaUI.Provider value="en">
      {layar.jenis === 'inbox' && <InboxEndorsementLife onBuka={keKasus} onBuat={() => setLayar({ jenis: 'buat' })} />}
      {layar.jenis === 'buat' && <BuatEndorsement onTutup={keInbox} onDibuat={keKasus} />}
      {layar.jenis === 'kasus' && <InputEDMLife kasusId={layar.id} onTutup={keInbox} />}
    </BahasaUI.Provider>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanEDM> = RuteEDM
