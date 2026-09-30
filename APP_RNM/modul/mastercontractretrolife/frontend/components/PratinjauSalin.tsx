// Pratinjau + konfirmasi `Copy to all Reinstype` - tiket 08, penyimpangan sadar 5.
//
// Pega (`SaveBusinessToAllLife_Act` b12252) langsung menulis ke setiap kontrak LAIN setahun yang
// berjenis BERBEDA (langkah 3.1 b646 `.REINSTYPEID == Param.REINSTYPEID` T=3 = dilewati; RALAT R2).
// Di sini pratinjau lebih dulu: sasaran dan jenis dasar pemilihan disebut; nol sasaran = diberi tahu
// TANPA tombol konfirmasi; `Cancel` tidak menulis apa pun; `Yes` mengirim sasaran yang DILIHAT.

import { useEffect, useState } from 'react'

import { pratinjauSalinSemua, salinSemua, type Business, type PratinjauSalin as Pratinjau } from '../api'
import { BUSINESS_MCRL, KONTRAK_MCRL, SALIN_MCRL } from '../labels'
import { sel } from '../tampilan'
import { Gagal, Modal } from '../../../../inti/frontend/components/ui/dasar'

/** ID sasaran yang dikirim `Yes` - persis yang tampil di pratinjau. */
export function sasaranDariPratinjau(p: Pratinjau): string[] {
  return p.sasaran.map((k) => k.id)
}

export default function PratinjauSalin({
  business,
  onSelesai,
  onBatal,
}: {
  business: Business
  onSelesai: (pesan: string) => void
  onBatal: () => void
}) {
  const [pratinjau, setPratinjau] = useState<Pratinjau | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    let hidup = true
    pratinjauSalinSemua(business.id)
      .then((p) => {
        if (hidup) setPratinjau(p)
      })
      .catch((e: unknown) => {
        if (hidup) setGalat(e)
      })
    return () => {
      hidup = false
    }
  }, [business.id])

  async function ya(): Promise<void> {
    if (pratinjau === null || sibuk) return
    setSibuk(true)
    setGalat(null)
    try {
      const h = await salinSemua(business.id, sasaranDariPratinjau(pratinjau))
      onSelesai(h.pesan)
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  const adaSasaran = pratinjau !== null && pratinjau.sasaran.length > 0

  return (
    <Modal
      judul={BUSINESS_MCRL.copyToAll}
      onTutup={onBatal}
      labelBatal={BUSINESS_MCRL.cancel}
      lebar
      aksi={
        adaSasaran ? (
          <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => void ya()}>
            {SALIN_MCRL.ya}
          </button>
        ) : undefined
      }
    >
      <p>
        <strong>{sel(business.bizName)}</strong>
      </p>
      {pratinjau === null && galat === null && <p role="status">{SALIN_MCRL.memuat}</p>}
      {galat !== null && <Gagal galat={galat} />}
      {pratinjau !== null && (
        <>
          <p>
            {SALIN_MCRL.dasar} <strong>{sel(pratinjau.business.reinsTypeName || pratinjau.reinsTypeId)}</strong>
          </p>
          {adaSasaran ? (
            <>
              <p>
                {SALIN_MCRL.akanDitulis} <strong>{pratinjau.sasaran.length}</strong>
              </p>
              <table className="inbox__tabel">
                <thead>
                  <tr>
                    <th>{KONTRAK_MCRL.kolomId}</th>
                    <th>{KONTRAK_MCRL.kolomReinsType}</th>
                  </tr>
                </thead>
                <tbody>
                  {pratinjau.sasaran.map((k) => (
                    <tr key={k.id} className="inbox__baris">
                      <td>{sel(k.id)}</td>
                      <td>{sel(k.reinsTypeName)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </>
          ) : (
            <p role="status">{SALIN_MCRL.nol}</p>
          )}
        </>
      )}
    </Modal>
  )
}
