// Rute modul NB FacIn - tiket 21; halaman depan tiket 25/26; Inward Facultative tiket 30.

import { useState } from 'react'

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanNbFacIn } from './menu'
import CoverageCargo from './pages/CoverageCargo'
import FormOpportunity from './pages/FormOpportunity'
import InwardFacultative, { type KasusBaru } from './pages/InwardFacultative'
import PortalOpportunity from './pages/PortalOpportunity'
import './nbfacin.css'

export function RuteNbFacIn({ halaman, masuk, onPindah }: PropsRute<HalamanNbFacIn>) {
  // Case yang baru dibuat Create opportunity - dibawa ke layar Inward Facultative (assignment pertama).
  // Komponen rute tetap terpasang selama modul aktif, jadi keadaan ini bertahan saat pindah halaman.
  const [kasus, setKasus] = useState<KasusBaru | null>(null)
  // Pesan konfirmasi sesudah Submit berhasil - tampil di portal (tiket 57).
  const [pesanPortal, setPesanPortal] = useState('')
  return (
    <>
      {halaman === 'nbfacin-portal' && (
        <PortalOpportunity
          onBuat={() => {
            setPesanPortal('')
            onPindah('nbfacin-opportunity')
          }}
          onBuka={(caseId) => {
            setPesanPortal('')
            setKasus({ caseId })
            onPindah('nbfacin-inward')
          }}
          pesan={pesanPortal}
          onTutupPesan={() => setPesanPortal('')}
        />
      )}
      {halaman === 'nbfacin-opportunity' && (
        <FormOpportunity
          pemilik={masuk.nama ?? masuk.akunID}
          onDibuat={(k) => {
            setKasus(k)
            onPindah('nbfacin-inward')
          }}
        />
      )}
      {halaman === 'nbfacin-inward' && kasus && (
        <InwardFacultative
          kasus={kasus}
          onBatal={() => onPindah('nbfacin-portal')}
          onSelesai={(pesan) => {
            // Submit berhasil: case ditutup, kembali ke portal (Pega `finishAssignment`).
            setPesanPortal(pesan)
            setKasus(null)
            onPindah('nbfacin-portal')
          }}
        />
      )}
      {halaman === 'nbfacin-coverage-cargo' && <CoverageCargo />}
    </>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanNbFacIn> = RuteNbFacIn
