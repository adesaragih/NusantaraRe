// Rute modul NB FacIn - tiket 21.

import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanNbFacIn } from './menu'
import CoverageCargo from './pages/CoverageCargo'

export function RuteNbFacIn({ halaman }: PropsRute<HalamanNbFacIn>) {
  return <>{halaman === 'nbfacin-coverage-cargo' && <CoverageCargo />}</>
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanNbFacIn> = RuteNbFacIn
