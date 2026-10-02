// `View Rate` b34113 (grid `PLAN LIST`, vis `OTHER .RIRATE!=''`) → `SetParamRate` b34310 + `localAction ViewRate`
// b34354: FlowAction `ViewRate` (submit `Submit` b18, `Cancel` b20), section `ViewRate` - judul `Outward List`
// b843, grid `ID` · `USEDBY` · `GENDER` · `CONTRACT` · `AGE` · `RATE` (RD atas `RATE_LIFE`).
//
// K1 keputusan work owner 01-10-2026 (OQ-MPNL-03): view `RATE_LIFE` dibaca saja, disaring `RIRATEID` baris
// plan. ⚠️ Penyimpangan sadar: grid `ViewRate.xml` b1024 menyaring `ParamID.OUTWARDRATEID` yang tidak pernah
// diisi rule mana pun (`SetParamRate` b259 mengisi `ParamID.RIRATEID` dari halaman Retro Life). View tak
// terbaca = 503 berkalimat yang menyebut view-nya, dan kalimat itu yang tampil.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilRate, type BarisRate } from '../api'
import { LAIN_MPNL, RATE_MPNL } from '../labels'

const KOLOM = [
  ['id', RATE_MPNL.kolomId],
  ['usedBy', RATE_MPNL.usedby],
  ['gender', RATE_MPNL.gender],
  ['contract', RATE_MPNL.contract],
  ['age', RATE_MPNL.age],
  ['rate', RATE_MPNL.rate],
] as const satisfies readonly (readonly [keyof BarisRate, string])[]

export default function ModalRate({ riRateId, onTutup }: { riRateId: string; onTutup: () => void }) {
  const [baris, setBaris] = useState<BarisRate[] | null>(null)
  const [terpotong, setTerpotong] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let batal = false
    ambilRate(riRateId)
      .then((d) => {
        if (!batal) {
          setBaris(d.daftar)
          setTerpotong(d.terpotong)
        }
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
      {terpotong && <p className="mpnl-catatan-medan">{LAIN_MPNL.terpotong}</p>}
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
                  <td key={k}>{b[k]}</td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Modal>
  )
}
