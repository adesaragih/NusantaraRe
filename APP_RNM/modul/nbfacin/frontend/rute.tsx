// Rute modul NB FacIn - tiket 21; halaman depan tiket 25/26.

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanNbFacIn } from './menu'
import CoverageCargo from './pages/CoverageCargo'
import FormOpportunity from './pages/FormOpportunity'
import PortalOpportunity from './pages/PortalOpportunity'
import './nbfacin.css'

export function RuteNbFacIn({ halaman, masuk, onPindah }: PropsRute<HalamanNbFacIn>) {
  return (
    <>
      {halaman === 'nbfacin-portal' && <PortalOpportunity onBuat={() => onPindah('nbfacin-opportunity')} />}
      {halaman === 'nbfacin-opportunity' && <FormOpportunity pemilik={masuk.nama ?? masuk.akunID} />}
      {halaman === 'nbfacin-coverage-cargo' && <CoverageCargo />}
    </>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanNbFacIn> = RuteNbFacIn
