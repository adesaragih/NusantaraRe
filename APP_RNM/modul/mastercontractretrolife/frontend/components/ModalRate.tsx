// Section `ViewRate` (`Rate List`) - dibuka tombol `View Rate` form (`localAction ViewRate`, pre-processing
// `SetParamRate`: `ParamID.RIRATEID ← InputBusinessLife.RIRATEID`) dan tombol `View Rate` baris
// (`SetParamRateTable` `RIRATEID ← .RIRATEID` + `localAction ViewRateTable`). Grid RD `BrowseRateLife_RD`
// param `idusedby = ParamID.RIRATEID` (`ViewRate.xml` b1061); tanpa penomoran.
//
// ⚠️ OQ-MCRL-13: sumber tabel rate menunggu persetujuan work owner - server menjawab 503 berkalimat,
// dan kalimat itu yang tampil (bukan daftar kosong yang membuat pengguna mengira tidak ada rate).

import { useEffect, useState } from 'react'

import { ambilRate, type BarisRate } from '../api'
import { RATE_MCRL, UMUM_MCRL } from '../labels'
import { sel } from '../tampilan'
import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'

export default function ModalRate({ idusedby, onTutup }: { idusedby: string; onTutup: () => void }) {
  const [daftar, setDaftar] = useState<BarisRate[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    ambilRate(idusedby)
      .then((d) => {
        if (hidup) setDaftar(d.daftar)
      })
      .catch((e: unknown) => {
        if (hidup) setGalat(e)
      })
    return () => {
      hidup = false
    }
  }, [idusedby])

  return (
    <Modal
      judul={RATE_MCRL.judul}
      onTutup={onTutup}
      labelBatal={RATE_MCRL.cancel}
      lebar
      aksi={
        // FlowAction tanpa post-processing: `Submit` sekadar menutup dialog, seperti `Cancel`.
        <button type="button" className="btn btn--primary" onClick={onTutup}>
          {RATE_MCRL.submit}
        </button>
      }
    >
      {daftar === null && galat === null && <Memuat />}
      {galat !== null && <Gagal galat={galat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={UMUM_MCRL.kosong} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              <th>{RATE_MCRL.kolomId}</th>
              <th>{RATE_MCRL.kolomUsedBy}</th>
              <th>{RATE_MCRL.kolomGender}</th>
              <th>{RATE_MCRL.kolomContract}</th>
              <th>{RATE_MCRL.kolomAge}</th>
              <th>{RATE_MCRL.kolomRate}</th>
            </tr>
          </thead>
          <tbody>
            {daftar.map((r) => (
              <tr key={r.id} className="inbox__baris">
                <td>{sel(r.id)}</td>
                <td>{sel(r.usedBy)}</td>
                <td>{sel(r.gender)}</td>
                <td>{sel(r.contract)}</td>
                <td>{sel(r.age)}</td>
                <td>{sel(r.rate)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Modal>
  )
}
