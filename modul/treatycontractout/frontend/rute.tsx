// Rute modul Treaty Contract Out - refactor bentuk B (30-09-2026): dipindah apa
// adanya dari `App.tsx`.
//
// [keputusan work owner 30-09-2026] Layar Treaty Contract Out berbahasa
// INGGRIS: `BahasaUI` "en" membuat teks bawaan komponen bersama (Memuat,
// penomoran halaman, panel backend mati, pilihan kosong, tutup popup) ikut
// Inggris di layar ini saja; label modul di `labels.ts`.

import { BahasaUI } from '../../../inti/frontend/components/ui/bahasaUI'
import type { PropsRute, RuteModul } from '../../../inti/frontend/modul'
import type { HalamanTreaty } from './menu'
import InboxTreatyContract from './pages/InboxTreatyContract'
import InboxTreatyContractDescription from './pages/InboxTreatyContractDescription'
import InboxTreatyContractReinsType from './pages/InboxTreatyContractReinsType'
import './tco.css'

export function RuteTreaty({ halaman }: PropsRute<HalamanTreaty>) {
  return (
    <BahasaUI.Provider value="en">
      {/* Treaty Contract Out tiket 03 — layar tahun treaty (harness InboxTreatyContract). */}
      {halaman === 'tco-tahun' && <InboxTreatyContract />}
      {/* Tiket 04 — editor kontrak dari menu (harness InboxTreatyContractReinsType). */}
      {halaman === 'tco-kontrak' && <InboxTreatyContractReinsType />}
      {halaman === 'tco-klausul' && <InboxTreatyContractDescription />}
    </BahasaUI.Provider>
  )
}

/** Rute modul ini untuk perakit `frontend/daftar.ts` - nama ekspor sama di setiap modul. */
export const RUTE_MODUL: RuteModul<HalamanTreaty> = RuteTreaty
