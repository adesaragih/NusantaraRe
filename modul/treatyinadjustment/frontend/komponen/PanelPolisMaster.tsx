// Panel `Existing Policy for Master ID` — `Section/InputTreatyInAdjustment.xml`
// @111283: grid `PolisList.pxResults`, dua kolom (sel 31 `Policy No`, sel 32
// `Pega ID`), di antara kepala dan panel Old/New. Wadahnya
// `pyIsVisibilityOption = ALWAYS` di bawah `DATASHOW = 1`.
//
// ⭐ Seam SENDIRI (`/kontrak-warisan/{id}/polis`): polis yang gagal dibaca
// tidak mengosongkan layar detail, sama seperti Attachment dan History.

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Panel } from '../../../../inti/frontend/components/ui/dasar'
import { ambilPolisMaster, type BarisPolisMaster } from '../api'
import { POLIS_MASTER } from '../labelsPenyesuaian'

/**
 * Pengenal yang Activity `FetchTreatyExistingProduction` cari:
 * `@if(TreatyIn.EDMState="", TreatyIn.ID, TreatyIn.OLDID)`.
 */
export function idMasterPolis(medan: Readonly<Record<string, string>>, id: string, idAsal: string): string {
  return (medan.EDMState ?? '') === '' ? id : idAsal
}

export default function PanelPolisMaster({ idMaster }: { idMaster: string }) {
  const [baris, setBaris] = useState<BarisPolisMaster[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  useEffect(() => {
    if (idMaster === '') {
      setBaris([])
      return
    }
    let dibuang = false
    setBaris(null)
    setGalat(null)
    ambilPolisMaster(idMaster)
      .then((d) => {
        if (!dibuang) setBaris(d)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [idMaster])
  return (
    <Panel judul={POLIS_MASTER.judul}>
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && (
        <div className="table-wrap">
          <table className="tria__tabel">
            <thead>
              <tr>
                {POLIS_MASTER.kolom.map((k) => (
                  <th key={k} scope="col">
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {baris.length === 0 && (
                <tr>
                  <td colSpan={POLIS_MASTER.kolom.length}>{POLIS_MASTER.tanpaBaris}</td>
                </tr>
              )}
              {/* ⛔ Pengenal — tidak pernah diformat. Baris yang sama dapat
                  muncul lebih dari sekali: DISTINCT rule-nya atas empat kolom
                  (termasuk kuartal), panel menampilkan dua. */}
              {baris.map((b, i) => (
                <tr key={`${b.nomorPolis}|${b.pegaID}|${String(i)}`}>
                  <td>{b.nomorPolis}</td>
                  <td>{b.pegaID}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Panel>
  )
}
