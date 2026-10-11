// Rincian satu peserta - `pyEditAction` `PL_DetailAction` (`InputEDMLife.xml` b16187, b20938) → FlowAction
// `PL_DetailAction` (b170) → `Section/PL_Detail_Sec.xml` (medan per wadah `pyWorkPage.Type`), dan grid retro
// `pyEditAction` `RetroLife` (`PL_Detail_Sec.xml` b11653) → `RetroDetailLife`. Seluruhnya baca-saja.
//
// Tata letak dirapikan (keputusan work owner 08-10-2026, tampilan saja): medan dikelompokkan bagian berjudul
// (`BAGIAN_RINCI`: Participant, Period, Coverage, Premium); lalu bagian Retro Spreading - setiap lapisan (Treaty
// Type) sub-judul berisi Retroceded Share lapisan, lalu tabel reinsurernya. Semua angka lewat `selAngka`.

import { useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilRincian, type RincianPesertaEDM } from '../api'
import { RINCI_PESERTA, RINCI_RETRO, RINCI_SPREADING } from '../kolomKorpus'
import { RINCI_EDM, UMUM_EDM } from '../labels'
import { kelompokRinci, MEDAN_RINCI_LEBAR, selAngka } from '../tampilan'
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
  const [labelTreaty, labelBagian] = [RINCI_SPREADING[0]?.label ?? '', RINCI_SPREADING[1]?.label ?? '']
  return (
    <div className="edm-rinci">
      <div className="edm-rinci__semua">
        {kelompokRinci(RINCI_PESERTA[t]).map((b) => (
          <section key={b.judul} className="edm-rinci__bagian">
            <h4 className="edm-subjudul">{b.judul}</h4>
            <MedanKorpus
              kolom={b.kolom}
              nilai={isi.peserta.nilai}
              className="edm-kepala edm-rinci__grid"
              lebar={MEDAN_RINCI_LEBAR}
              kelasLebar="edm-rinci__lebar"
            />
          </section>
        ))}
      </div>
      {t !== 'QP' && isi.spreading.length > 0 && (
        <section className="edm-rinci__spreading">
          <h4 className="edm-subjudul">{RINCI_EDM.bagianSpreading}</h4>
          {isi.spreading.map((s) => (
            <div key={s.id} className="edm-rinci__lapisan">
              <h5 className="edm-rinci__lapisanjudul">
                <span>
                  {labelTreaty}: <strong>{s.treatyTypeName}</strong>
                </span>
                <span>
                  {labelBagian}: <strong className="edm-angka">{selAngka(s.retrocadedShare)}</strong>
                </span>
              </h5>
              {s.retro.length > 0 && (
                <TabelKorpus
                  kolom={RINCI_RETRO}
                  baris={s.retro.map((r) => ({ ...r }))}
                  kunci={(b, i) => `${s.id}/${b.reinsurerName ?? ''}/${i}`}
                />
              )}
            </div>
          ))}
        </section>
      )}
    </div>
  )
}
