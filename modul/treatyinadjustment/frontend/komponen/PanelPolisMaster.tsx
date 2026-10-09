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

/**
 * `ringkas` — bentuk Pega (9 Oktober 2026): panel duduk DI KANAN kepala,
 * judul kecil tanpa kartu, grid berkolom nomor baris (`1`) di kiri.
 */
export default function PanelPolisMaster({ idMaster, ringkas = false }: { idMaster: string; ringkas?: boolean }) {
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
  if (ringkas) {
    return (
      <section className="tria__polis" aria-label={POLIS_MASTER.judul}>
        <h4 className="tria__polis-judul">{POLIS_MASTER.judul}</h4>
        {galat !== null && <Gagal galat={galat} />}
        {baris === null && galat === null && <Memuat />}
        {baris !== null && (
          <div className="table-wrap">
            <table className="tria__tabel tria__tabel--pega">
              <thead>
                <tr>
                  <th scope="col" className="tria__polis-no" aria-label="No" />
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
                    <td colSpan={POLIS_MASTER.kolom.length + 1}>{POLIS_MASTER.tanpaBaris}</td>
                  </tr>
                )}
                {baris.map((b, i) => (
                  <tr key={`${b.nomorPolis}|${b.pegaID}|${String(i)}`}>
                    <td className="tria__polis-no">{i + 1}</td>
                    <td>{b.nomorPolis}</td>
                    <td>{b.pegaID}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    )
  }
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
