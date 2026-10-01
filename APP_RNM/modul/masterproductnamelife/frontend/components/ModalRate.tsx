// `View Rate` b34067 (grid `PLAN LIST`, vis `OTHER .RIRATE!=''`) → `SetParamRate` b34299 + `localAction ViewRate`
// b34332: FlowAction `ViewRate` (submit `Submit` b18, `Cancel` b20), section `ViewRate` - judul `Outward List`
// b843, grid `ID` · `USEDBY` · `GENDER` · `CONTRACT` · `AGE` · `RATE` (RD atas `RATE_LIFE`).
//
// ⏸️ OQ-MPNL-03: sumber rate (view atas JSON) belum disetujui - server menjawab 503 berkalimat dan dialog
// menampilkannya; tombol dan kolomnya tetap ada seperti XML.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilRate } from '../api'
import { LAIN_MPNL, RATE_MPNL } from '../labels'

const KOLOM = [
  ['ID', RATE_MPNL.kolomId],
  ['USEDBY', RATE_MPNL.usedby],
  ['GENDER', RATE_MPNL.gender],
  ['CONTRACT', RATE_MPNL.contract],
  ['AGE', RATE_MPNL.age],
  ['RATE', RATE_MPNL.rate],
] as const

export default function ModalRate({ riRateId, onTutup }: { riRateId: string; onTutup: () => void }) {
  const [baris, setBaris] = useState<Record<string, string>[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let batal = false
    ambilRate(riRateId)
      .then((d) => {
        if (!batal) setBaris(d.daftar)
      })
      .catch((e: unknown) => {
        if (!batal) setGalat(e)
      })
    return () => {
      batal = true
    }
  }, [riRateId])

  return (
    <Modal
      judul={RATE_MPNL.judul}
      onTutup={onTutup}
      labelBatal={RATE_MPNL.cancel}
      lebar
      aksi={
        <button type="button" className="btn btn--primary" onClick={onTutup}>
          {RATE_MPNL.submit}
        </button>
      }
    >
      {baris === null && galat === null && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {baris !== null && baris.length === 0 && <Kosong pesan={LAIN_MPNL.kosong} />}
      {baris !== null && baris.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              {KOLOM.map(([k, l]) => (
                <th key={k}>{l}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {baris.map((b, i) => (
              <tr key={i} className="inbox__baris">
                {KOLOM.map(([k]) => (
                  <td key={k}>{b[k] ?? ''}</td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Modal>
  )
}
