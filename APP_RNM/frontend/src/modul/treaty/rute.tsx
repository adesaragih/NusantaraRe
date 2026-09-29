// Rute modul Treaty Contract Out - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `App.tsx`.

import type { PropsRute } from '../../inti/modul'
import type { HalamanTreaty } from './menu'
import InboxTreatyContract from './pages/InboxTreatyContract'
import InboxTreatyContractDescription from './pages/InboxTreatyContractDescription'
import InboxTreatyContractReinsType from './pages/InboxTreatyContractReinsType'

export function RuteTreaty({ halaman }: PropsRute<HalamanTreaty>) {
  return (
    <>
      {/* Treaty Contract Out tiket 03 — layar tahun treaty (harness InboxTreatyContract). */}
      {halaman === 'tco-tahun' && <InboxTreatyContract />}
      {/* Tiket 04 — editor kontrak dari menu (harness InboxTreatyContractReinsType). */}
      {halaman === 'tco-kontrak' && <InboxTreatyContractReinsType />}
      {halaman === 'tco-klausul' && <InboxTreatyContractDescription />}
    </>
  )
}
