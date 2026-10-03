// Kotak masuk Endorsement Life - `Section/InboxEndorsementLife.xml` (harness portal `InboxEndorsementLife`).
//
// Grid b8284 dari RD `InboxEDMLife`: kasus yang belum `Resolved-Completed` dan belum `Resolved-Rejected`,
// terbaru dahulu. Tautan `Case ID` (b10202 → `openAssignment` b10349) membuka layar kasus.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilInbox, type BarisInbox, type HalamanEDM } from '../api'
import { INBOX_EDM, UMUM_EDM } from '../labels'
import { UKURAN_HALAMAN_EDM, labelEdmType, sel, selWaktu } from '../tampilan'
import '../endorsementlife.css'

/** Kepala grid VERBATIM, urut korpus (b8418 … b10052). */
const KOLOM: readonly string[] = [
  INBOX_EDM.kolomCaseId,
  INBOX_EDM.kolomEndorsementNo,
  INBOX_EDM.kolomType,
  INBOX_EDM.kolomEdmType,
  INBOX_EDM.kolomPolicyNo,
  INBOX_EDM.kolomSob,
  INBOX_EDM.kolomCeding,
  INBOX_EDM.kolomPolicyHolder,
  INBOX_EDM.kolomMarketingName,
  INBOX_EDM.kolomCreateDate,
  INBOX_EDM.kolomCreateOperator,
  INBOX_EDM.kolomEdmTypeBatal,
  INBOX_EDM.kolomStatus,
]

export default function InboxEndorsementLife({ onBuka, onBuat }: { onBuka: (id: string) => void; onBuat: () => void }) {
  const [halaman, setHalaman] = useState(1)
  const [isi, setIsi] = useState<HalamanEDM<BarisInbox> | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  const muat = useCallback(async (h: number) => {
    setGalat(null)
    try {
      setIsi(await ambilInbox(h))
    } catch (e) {
      setGalat(e)
    }
  }, [])

  useEffect(() => {
    void muat(halaman)
  }, [muat, halaman])

  return (
    <section className="inbox edm-inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{INBOX_EDM.judul}</h2>
        {/* `Create Addendum` b6620 → `showHarness` b6984 `EndorsmentLife_harnes` (`pyTarget=current`). */}
        <button type="button" className="btn btn--primary btn--sm" onClick={onBuat}>
          {INBOX_EDM.createAddendum}
        </button>
      </header>
      {galat !== null && <Gagal galat={galat} />}
      {isi === null && galat === null && <Memuat pesan={UMUM_EDM.memuat} />}
      {isi !== null && isi.baris.length === 0 && <Kosong pesan={UMUM_EDM.kosong} />}
      {isi !== null && isi.baris.length > 0 && (
        <div className="edm-gulir">
          <table className="inbox__tabel">
            <thead>
              <tr>
                {KOLOM.map((k) => (
                  <th key={k}>{k}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {isi.baris.map((b) => (
                <tr key={b.caseId}>
                  <td>
                    <button type="button" className="btn btn--ghost btn--sm edm-tautan" onClick={() => onBuka(b.caseId)}>
                      {b.caseId}
                    </button>
                  </td>
                  <td>{sel(b.endorsementNo)}</td>
                  <td>{sel(b.type)}</td>
                  <td>{labelEdmType(b.edmType)}</td>
                  <td>{sel(b.policyNo)}</td>
                  <td>{sel(b.sob)}</td>
                  <td>{sel(b.ceding)}</td>
                  <td>{sel(b.policyHolder)}</td>
                  <td>{sel(b.marketingName)}</td>
                  <td>{selWaktu(b.createDate)}</td>
                  <td>{sel(b.createOperator)}</td>
                  {/* `A.EdmTypeBatal` b12243 `VIS A.EdmType==3` - opsinya belum ditetapkan (OQ-EDM-005). */}
                  <td>{b.edmType === '3' ? sel('') : ''}</td>
                  <td>{sel(b.status)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
      {isi !== null && isi.total > 0 && (
        <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN_EDM} total={isi.total} onPindah={setHalaman} />
      )}
    </section>
  )
}
