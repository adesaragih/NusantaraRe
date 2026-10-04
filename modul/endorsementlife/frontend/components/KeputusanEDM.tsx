// `Section/ConfirmSection.xml` (wadah `InputEDMLife.xml` b35518 `.IsJsonPolis=1`): radio `Status` b496,
// `Comment` b829, riwayat b1856, lalu `Submit` - b37494 tampil bila `EmailTypePL = 1` (Confirm),
// b38109 bila `2`/`7` (Decline). Keduanya berlabel sama; satu tombol tampil sesuai pilihan.

import { useState } from 'react'

import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import { putuskanKasus, type HasilPutusanEDM, type RiwayatEDM } from '../api'
import { OPSI_KEPUTUSAN, PUTUSAN_EDM, UMUM_EDM } from '../labels'
import { sel, selWaktu } from '../tampilan'

export default function KeputusanEDM({
  kasusId,
  riwayat,
  onSelesai,
}: {
  kasusId: string
  riwayat: readonly RiwayatEDM[]
  onSelesai: (h: HasilPutusanEDM) => void
}) {
  const [status, setStatus] = useState('')
  const [komentar, setKomentar] = useState('')
  const [mengirim, setMengirim] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)

  const kirim = async () => {
    setMengirim(true)
    setGalat(null)
    try {
      onSelesai(await putuskanKasus(kasusId, { status, comment: komentar }))
    } catch (e) {
      setGalat(e)
    } finally {
      setMengirim(false)
    }
  }

  return (
    <section className="edm-putusan">
      <fieldset className="edm-putusan__status">
        <legend>{PUTUSAN_EDM.status}</legend>
        {Object.entries(OPSI_KEPUTUSAN).map(([nilai, label]) => (
          <label key={nilai}>
            <input type="radio" name="edm-status" value={nilai} checked={status === nilai} onChange={() => setStatus(nilai)} />
            {label}
          </label>
        ))}
      </fieldset>
      <label className="edm-putusan__komentar">
        <span>{PUTUSAN_EDM.comment}</span>
        <textarea maxLength={255} rows={3} value={komentar} onChange={(e) => setKomentar(e.target.value)} />
      </label>
      {status !== '' && (
        <div className="edm-aksi">
          <button type="button" className="btn btn--sm" disabled={mengirim} onClick={() => void kirim()}>
            {mengirim ? UMUM_EDM.memutuskan : PUTUSAN_EDM.submit}
          </button>
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {riwayat.length > 0 && (
        <div className="edm-gulir">
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{PUTUSAN_EDM.kolomDate}</th>
                <th>{PUTUSAN_EDM.kolomPic}</th>
                <th>{PUTUSAN_EDM.kolomStatus}</th>
                <th>{PUTUSAN_EDM.kolomComment}</th>
              </tr>
            </thead>
            <tbody>
              {riwayat.map((r) => (
                <tr key={r.no}>
                  <td>{selWaktu(r.date)}</td>
                  <td>{sel(r.pic)}</td>
                  <td>{sel(r.status)}</td>
                  <td>{sel(r.comment)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  )
}
