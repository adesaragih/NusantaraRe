// Tampilan polis NB Fac In HANYA-BACA (unsur tata `polisNB`, perintah work owner 10-10-2026 "IKUTI MODUL NB FACIN.
// TAMPILAN READ ONLY"): salinan rupa layar Inward Facultative `modul/nbfacin/frontend/pages/InwardFacultative.tsx` -
// blok General dua kolom lalu sepuluh tab detail LANGSUNG (pemisah kosong dan centang Show Detail nbfacin dibuang -
// perintah work owner 11-10-2026 "INI DI HAPUS, BUAT LEBIH SIMPEL!"). Seperti nbfacin, isi tab hanya untuk kasus
// FIRE (When `IsFire` dievaluasi server: `kendali` "fire") dan hanya Object / Coverage / Spreading / Inw Fac Cedant
// Panels; Clauses (`pyWorkPage.ClauseList` kasus NB, bukan bagian halaman polis `OfferFacIn`) dan tab lain = pesan
// inti `BelumTersedia`. Medan isian, tombol, dan bilah CSV nbfacin tidak dibawa (hanya-baca).

import { useState } from 'react'

import { BelumTersedia } from '../../../../inti/frontend/components/ui/dasar'
import type { Halaman } from '../api'
import {
  PILIHAN_DAY,
  PILIHAN_POLICY_TYPE,
  TAB_BERISI_NB,
  TAB_NB,
  TEKS_NB,
  UMUM_NB,
  umumNB,
  type TabNB,
} from '../polisNB'
import { LabelNB, RadioNB, TeksNB } from './PolisNBDasar'
import { TabCedantNB, TabCoverageNB, TabObjekNB, TabSpreadingNB } from './PolisNBTab'
import { TabGeser } from './TabGeser'

/** `Tata.kendali` unsur `polisNB` bila When `IsFire` benar (`models.KPolisFire`). */
export const KENDALI_POLIS_FIRE = 'fire'

function IsiTabNB({ tab, h }: { tab: TabNB; h: Halaman }) {
  switch (tab) {
    case 'Object':
      return <TabObjekNB h={h} />
    case 'Coverage':
      return <TabCoverageNB h={h} />
    case 'Spreading':
      return <TabSpreadingNB h={h} />
    case 'Inw Fac Cedant Panels':
      return <TabCedantNB h={h} />
    default:
      return <BelumTersedia apa={`${TEKS_NB.isiTab} ${tab}`} />
  }
}

export default function PolisNB({ h, fire }: { h: Halaman; fire: boolean }) {
  const [tab, setTab] = useState<TabNB>(TAB_NB[0])
  const u = umumNB(h)
  const L = UMUM_NB
  return (
    <div className="claimfacin__nb">
      <section className="claimfacin__nb-kartu">
        <h4 className="claimfacin__nb-judul">{L.judul}</h4>
        <div className="claimfacin__nb-kolom">
          <div className="claimfacin__nb-tumpuk">
            <TeksNB kiri label={L.reffNumber} nilai={u.reff} />
            <TeksNB kiri label={L.businessStatus} nilai={u.statusBisnis} />
            <TeksNB kiri label={L.insuredName} nilai={u.insured} />
            <TeksNB kiri label={L.qqName} nilai={u.qq} />
            <TeksNB kiri label={L.beginDate} nilai={u.mulai} />
            <TeksNB kiri label={L.offeringDate} nilai={u.penawaran} />
            <RadioNB kiri label={L.policyType} pilihan={PILIHAN_POLICY_TYPE} nilai={u.jenisPolis} />
            <TeksNB kiri label={L.riskScoring} nilai={u.riskScoring} />
          </div>
          <div className="claimfacin__nb-tumpuk">
            <TeksNB kiri label={L.classOfBusiness} nilai={u.kelas} />
            <TeksNB kiri label={L.typeFacultative} nilai={u.typeFac} />
            <TeksNB kiri label={L.sourceOfBusiness} nilai={u.sob} />
            <TeksNB kiri label={L.cedingCoName} nilai={u.ceding} />
            <TeksNB kiri label={L.groupName} nilai={u.grup} />
            <TeksNB kiri label={L.endDate} nilai={u.akhir} />
            <LabelNB label={L.followingPolicyNumber} />
            <TeksNB kiri label={L.oldPolicyNumber} nilai={u.lama} />
            <TeksNB kiri label={L.marketingName} nilai={u.marketing} />
            <RadioNB kiri label={L.day} pilihan={PILIHAN_DAY} nilai={u.hari} />
          </div>
        </div>
      </section>
      <section className="claimfacin__nb-kartu">
        <div className="claimfacin__tab">
          <TabGeser tab={TAB_NB} aktif={tab} onPilih={setTab} />
          <div key={tab} role="tabpanel" className="claimfacin__tab-isi">
            {fire && TAB_BERISI_NB.includes(tab) ? (
              <IsiTabNB tab={tab} h={h} />
            ) : (
              <BelumTersedia apa={`${TEKS_NB.isiTab} ${tab}`} />
            )}
          </div>
        </div>
      </section>
    </div>
  )
}
