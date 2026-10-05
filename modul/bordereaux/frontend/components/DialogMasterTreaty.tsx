// Popup Choose Master Treaty - padanan flow action `ChooseMasterID` / section `ShowMasterTreatyIn_Sec`: pilih Cedant
// (autocomplete `BrowseAgentNonLife_RD`), lalu grid Master Treaty cedant itu (`BrowseTREATY_IN`) dengan tombol
// Choose (`ChooseIDMaster`). Kotak "Adjustment Treaty In" Pega tidak dibawa: tombol Choose grid EDM-nya tanpa aksi.

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { cariCedant, cariTreaty, type Cedant, type MasterTreaty } from '../api'
import { BDX } from '../labels'

export default function DialogMasterTreaty({ onPilih, onTutup }: { onPilih: (t: MasterTreaty) => void; onTutup: () => void }) {
  const [kata, setKata] = useState('')
  const [cedant, setCedant] = useState<Cedant[]>([])
  const [terpilih, setTerpilih] = useState<Cedant | null>(null)
  const [treaty, setTreaty] = useState<MasterTreaty[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    if (terpilih !== null) return
    let hidup = true
    const jam = setTimeout(() => {
      cariCedant(kata.trim()).then(
        (r) => {
          if (hidup) setCedant(r.daftar)
        },
        (g: unknown) => {
          if (hidup) setGalat(g)
        },
      )
    }, 250)
    return () => {
      hidup = false
      clearTimeout(jam)
    }
  }, [kata, terpilih])

  useEffect(() => {
    if (terpilih === null) return
    let hidup = true
    setTreaty(null)
    cariTreaty(terpilih.id).then(
      (r) => {
        if (hidup) setTreaty(r.daftar)
      },
      (g: unknown) => {
        if (hidup) setGalat(g)
      },
    )
    return () => {
      hidup = false
    }
  }, [terpilih])

  return (
    <Modal judul={BDX.chooseMasterTreaty} onTutup={onTutup} labelBatal={BDX.batal} penuh>
      <div className="bordereaux__master">
        <label className="field bordereaux__cedant">
          <span className="field__label">{BDX.cedant}</span>
          <input
            className="field__input"
            type="search"
            placeholder={BDX.cariCedant}
            value={terpilih === null ? kata : terpilih.nama}
            onChange={(e) => {
              setTerpilih(null)
              setTreaty(null)
              setKata(e.target.value)
            }}
          />
        </label>
        {terpilih === null && cedant.length > 0 && (
          <ul className="bordereaux__saran" role="listbox" aria-label={BDX.cedant}>
            {cedant.map((c) => (
              <li key={c.id}>
                <button type="button" className="bordereaux__saran-butir" onClick={() => setTerpilih(c)}>
                  {c.nama}
                </button>
              </li>
            ))}
          </ul>
        )}
        {galat !== null && <Gagal galat={galat} />}
        {terpilih !== null && treaty === null && galat === null && <Memuat pesan={BDX.mencari} />}
        {treaty !== null && treaty.length === 0 && <p className="muted">{BDX.treatyKosong}</p>}
        {treaty !== null && treaty.length > 0 && (
          <div className="bordereaux__gulir">
            <table className="inbox__tabel bordereaux__grid">
              <thead>
                <tr>
                  <th>{BDX.id}</th>
                  <th>{BDX.contractName}</th>
                  <th>{BDX.reinsuranceType}</th>
                  <th>{BDX.sourceOfBusiness}</th>
                  <th>{BDX.ceding}</th>
                  <th className="table__actions">{BDX.aksi}</th>
                </tr>
              </thead>
              <tbody>
                {treaty.map((t) => (
                  <tr key={t.id} className="inbox__baris">
                    <td>{t.id}</td>
                    <td>{t.contractName}</td>
                    <td>{t.reinsType}</td>
                    <td>{t.sobName}</td>
                    <td>{t.cedingName}</td>
                    <td className="table__actions">
                      <button type="button" className="btn btn--primary btn--sm" onClick={() => onPilih(t)}>
                        {BDX.choose}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </Modal>
  )
}
