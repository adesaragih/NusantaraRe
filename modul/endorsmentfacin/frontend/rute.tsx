// Rute modul Endorsment Fac In untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. Irisan I0-I3
// (keputusan work owner 11-10-2026, W-EDM-01…12): portal endorsement -> form Create Addendum -> kasus `EDM-` lahir,
// kembali ke portal dengan pesan.

import { useState } from 'react'

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import './endorsmentfacin.css'
import { PERIKSA } from './labels'
import type { HalamanEdmFacIn } from './menu'
import FormAddendum from './pages/FormAddendum'
import PortalEndorsement from './pages/PortalEndorsement'

export function RuteEdmFacIn({ halaman, onPindah }: PropsRute<HalamanEdmFacIn>) {
  // Pesan sesudah kasus lahir - tampil di portal (rute tetap terpasang selama modul aktif).
  const [pesanPortal, setPesanPortal] = useState('')
  return (
    <BahasaUI.Provider value="en">
      <div className="endorsmentfacin">
        {halaman === 'endorsmentfacin-portal' && (
          <PortalEndorsement
            pesan={pesanPortal}
            onTutupPesan={() => setPesanPortal('')}
            onBuat={() => {
              setPesanPortal('')
              onPindah('endorsmentfacin-addendum')
            }}
          />
        )}
        {halaman === 'endorsmentfacin-addendum' && (
          <FormAddendum
            onBatal={() => onPindah('endorsmentfacin-portal')}
            onDibuat={(id) => {
              setPesanPortal(PERIKSA.dibuat(id))
              onPindah('endorsmentfacin-portal')
            }}
          />
        )}
      </div>
    </BahasaUI.Provider>
  )
}

export const RUTE_MODUL: RuteModul<HalamanEdmFacIn> = RuteEdmFacIn
