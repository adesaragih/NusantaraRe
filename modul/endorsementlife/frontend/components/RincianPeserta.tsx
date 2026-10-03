// Rincian satu peserta - `pyEditAction` `PL_DetailAction` (`InputEDMLife.xml` b16187, b20938) → FlowAction
// `PL_DetailAction` (b170) → `Section/PL_Detail_Sec.xml` (medan per wadah `pyWorkPage.Type`), dan grid retro
// `pyEditAction` `RetroLife` (`PL_Detail_Sec.xml` b11653) → `RetroDetailLife`. Seluruhnya baca-saja.

import { useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilRincian, type RincianPesertaEDM } from '../api'
import { RINCI_PESERTA, RINCI_RETRO, RINCI_SPREADING } from '../kolomKorpus'
import { UMUM_EDM } from '../labels'
import { MedanKorpus, TabelKorpus } from './TabelKorpus'
import { tipePopup } from './PolisLama'

export default function RincianPeserta({ kasusId, pesertaId, tipe }: { kasusId: string; pesertaId: string; tipe: string }) {
  const [isi, setIsi] = useState<RincianPesertaEDM | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let batal = false
    ambilRincian(kasusId, pesertaId)
      .then((r) => {
        if (!batal) setIsi(r)
      })
      .catch((e: unknown) => {
        if (!batal) setGalat(e)
      })
    return () => {
      batal = true
    }
  }, [kasusId, pesertaId])

  if (galat !== null) return <Gagal galat={galat} />
  if (isi === null) return <Memuat pesan={UMUM_EDM.memuat} />
  // Wadah `QP` di `PL_Detail_Sec` tanpa grid retro (b12819 … b21495) - seperti korpus.
  const t = tipePopup(tipe)
  return (
    <div className="edm-rinci">
      <MedanKorpus kolom={RINCI_PESERTA[t]} nilai={isi.peserta.nilai} />
      {t !== 'QP' && isi.spreading.length > 0 && (
        <>
          <TabelKorpus
            kolom={RINCI_SPREADING}
            baris={isi.spreading.map((s) => ({ treatyTypeName: s.treatyTypeName, retrocadedShare: s.retrocadedShare }))}
            kunci={(_, i) => isi.spreading[i]?.id ?? String(i)}
          />
          {isi.spreading.map((s) =>
            s.retro.length === 0 ? null : (
              <TabelKorpus
                key={s.id}
                kolom={RINCI_RETRO}
                baris={s.retro.map((r) => ({ ...r }))}
                kunci={(b, i) => `${s.id}/${b.reinsurerName ?? ''}/${i}`}
              />
            ),
          )}
        </>
      )}
    </div>
  )
}
