// Tampilan polis NB Fac In HANYA-BACA (unsur tata `polisNB`, perintah work owner 10-10-2026 "IKUTI MODUL NB FACIN.
// TAMPILAN READ ONLY"): salinan rupa layar Inward Facultative `modul/nbfacin/frontend/pages/InwardFacultative.tsx` -
// blok General dua kolom lalu tab detail LANGSUNG (pemisah kosong dan centang Show Detail nbfacin dibuang - perintah
// work owner 11-10-2026 "INI DI HAPUS, BUAT LEBIH SIMPEL!"). Tab dan isinya per lini seperti nbfacin (`tabTampil` +
// dispatch L345-584; work owner 11-10-2026 "tampilan view polis yang belum tergenerate tolong diperbaiki. cek tampilan
// nya di nb facin!"): FIRE (When `IsFire` server, `kendali` "fire") > keluarga Aneka > Marine Cargo (`polisNBAneka.ts`
// `liniNB`); Clauses / Scoring Risk hanya FIRE, Loss Record / Correspondence hanya akun IT. Lini lain = pesan nbfacin
// `cobBelum`. Medan isian, tombol, dan bilah CSV nbfacin tidak dibawa (hanya-baca).

import { useEffect, useState } from 'react'

import { BelumTersedia } from '../../../../inti/frontend/components/ui/dasar'
import type { Halaman } from '../api'
import { PILIHAN_DAY, PILIHAN_POLICY_TYPE, TAB_NB, TEKS_NB, UMUM_NB, angkaNB, umumNB, type TabNB } from '../polisNB'
import { liniNB, type LiniNB } from '../polisNBAneka'
import { ambilPolisNBIT, type PolisNBIT } from '../polisNBIT'
import { TabCoverageAnekaNB, TabObjekAnekaNB, TabSpreadingAnekaNB } from './PolisNBAneka'
import { TabCorrespondenceNB, TabLossRecordNB } from './PolisNBIT'
import { TabCedantLengkapNB, TabDeductionNB, TabPaymentNB } from './PolisNBBayar'
import { TabCoverageCargoNB, TabObjekCargoNB, TabSpreadingCargoNB } from './PolisNBCargo'
import { LabelNB, RadioNB, TeksNB } from './PolisNBDasar'
import { TabClausesNB } from './PolisNBKlausa'
import { LencanaSkorNB, TabScoringNB } from './PolisNBScoring'
import { TabCoverageNB, TabObjekNB, TabSpreadingNB } from './PolisNBTab'
import { TabGeser } from './TabGeser'

/** `Tata.kendali` unsur `polisNB` bila When `IsFire` benar (`models.KPolisFire`). */
export const KENDALI_POLIS_FIRE = 'fire'

/** Tab tampil (`tabTampil` nbfacin): Clauses / Scoring Risk hanya FIRE; Loss Record / Correspondence hanya akun IT. */
export function tabTampilNB(fire: boolean, aksesIT: boolean): TabNB[] {
  return TAB_NB.filter(
    (t) =>
      (t !== 'Loss Record' || aksesIT) &&
      (t !== 'Correspondence' || aksesIT) &&
      (t !== 'Clauses' || fire) &&
      (t !== 'Scoring Risk' || fire),
  )
}

function IsiTabNB({ tab, h, lini, it }: { tab: TabNB; h: Halaman; lini: LiniNB; it: PolisNBIT | null }) {
  // nbfacin: Correspondence (akun IT) untuk setiap COB; Loss Record (akun IT) hanya FIRE
  if (tab === 'Correspondence' && it) return <TabCorrespondenceNB daftar={it.korespondensi} />
  if (tab === 'Loss Record' && it && lini.fire) return <TabLossRecordNB lr={it.lossRecord} />
  if (lini.fire || lini.aneka || lini.cargo) {
    // tab bersama FIRE / Aneka / Marine Cargo (nbfacin: komponen sama ber-prop aneka / cargo)
    switch (tab) {
      case 'Inw Fac Cedant Panels':
        return <TabCedantLengkapNB h={h} lini={lini} />
      case 'Premium Deduction / Brokerage Fee':
        return <TabDeductionNB h={h} lini={lini} />
      case 'Payment':
        return <TabPaymentNB h={h} />
    }
  }
  if (lini.fire) {
    switch (tab) {
      case 'Object':
        return <TabObjekNB h={h} />
      case 'Coverage':
        return <TabCoverageNB h={h} />
      case 'Spreading':
        return <TabSpreadingNB h={h} />
      case 'Clauses':
        return <TabClausesNB h={h} />
      case 'Scoring Risk':
        return <TabScoringNB h={h} />
    }
  } else if (lini.aneka) {
    switch (tab) {
      case 'Object':
        return <TabObjekAnekaNB h={h} lini={lini} />
      case 'Coverage':
        return <TabCoverageAnekaNB h={h} lini={lini} />
      case 'Spreading':
        return <TabSpreadingAnekaNB h={h} />
    }
  } else if (lini.cargo) {
    switch (tab) {
      case 'Object':
        return <TabObjekCargoNB h={h} />
      case 'Coverage':
        return <TabCoverageCargoNB h={h} />
      case 'Spreading':
        return <TabSpreadingCargoNB h={h} />
    }
  } else {
    // nbfacin: COB di luar FIRE / Aneka / Marine Cargo = pesan `cobBelum` (layar Pega lini lain belum diport)
    return (
      <div className="alert alert--info">
        {TEKS_NB.cobBelum(tab, h.nilai['OfferFacIn.QuotationData.BusinessName'] ?? '')}
      </div>
    )
  }
  // nbfacin: tab tanpa isi di lini selain FIRE = pesan `cobBelum`, di FIRE = belum tersedia
  return lini.fire ? (
    <BelumTersedia apa={`${TEKS_NB.isiTab} ${tab}`} />
  ) : (
    <div className="alert alert--info">
      {TEKS_NB.cobBelum(tab, h.nilai['OfferFacIn.QuotationData.BusinessName'] ?? '')}
    </div>
  )
}

export default function PolisNB({ h, fire }: { h: Halaman; fire: boolean }) {
  const lini = liniNB(h, fire)
  const kasusId = h.nilai['pyID'] ?? ''
  // Tab Loss Record / Correspondence hanya bila server mengizinkan (JOB_POSITION "IT", 403 selain itu) - lunak seperti
  // nbfacin `AksesTabIT`: ditolak / gagal = kedua tab tidak tampil.
  const [it, setIT] = useState<PolisNBIT | null>(null)
  useEffect(() => {
    let batal = false
    setIT(null)
    if (kasusId === '') return
    ambilPolisNBIT(kasusId).then(
      (v) => {
        if (!batal) setIT(v)
      },
      () => {
        if (!batal) setIT(null)
      },
    )
    return () => {
      batal = true
    }
  }, [kasusId])
  const tabs = tabTampilNB(fire, it !== null)
  const [pilih, setTab] = useState<TabNB>(tabs[0] ?? TAB_NB[0])
  const tab = tabs.includes(pilih) ? pilih : (tabs[0] ?? TAB_NB[0])
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
            <TeksNB
              kiri
              label={L.riskScoring}
              nilai={
                u.riskScoring === '' ? (
                  ''
                ) : (
                  <>
                    {angkaNB(u.riskScoring)} <LencanaSkorNB catatan={u.catatanSkor} umum />
                  </>
                )
              }
            />
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
          <TabGeser tab={tabs} aktif={tab} onPilih={setTab} />
          <div key={tab} role="tabpanel" className="claimfacin__tab-isi">
            <IsiTabNB tab={tab} h={h} lini={lini} it={it} />
          </div>
        </div>
      </section>
    </div>
  )
}
