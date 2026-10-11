// `Section/ConfirmSection.xml` (wadah `InputEDMLife.xml` b35518 `.IsJsonPolis=1`): radio `Status` b496,
// `Comment` b829, riwayat b1856, lalu `Submit` - b37494 tampil bila `EmailTypePL = 1` (Confirm),
// b38109 bila `2`/`7` (Decline). Keduanya berlabel sama; satu tombol tampil sesuai pilihan.
//
// Sejak 08-10-2026 (keputusan work owner): tombol berlabel "Confirm" (utama), selalu tampil tetapi nonaktif selama
// Status belum dipilih; klik membuka popup konfirmasi (`KonfirmasiPutusan`), Yes menjalankan `kirim` yang sama.

import { useState } from 'react'

import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import { putuskanKasus, type HasilPutusanEDM, type RiwayatEDM } from '../api'
import { KONFIRMASI_PUTUSAN_EDM, OPSI_KEPUTUSAN, PUTUSAN_EDM, UMUM_EDM } from '../labels'
import { sel, selWaktu } from '../tampilan'
import KonfirmasiPutusan from './KonfirmasiPutusan'

export default function KeputusanEDM({
  kasusId,
  riwayat,
  onSelesai,
}: {
  kasusId: string
  riwayat: readonly RiwayatEDM[]
  /** Hasil + Status yang dipilih ('1' Accept) - layar memilih popup hasil (09-10-2026). */
  onSelesai: (h: HasilPutusanEDM, status: string) => void
}) {
  const [status, setStatus] = useState('')
  const [komentar, setKomentar] = useState('')
  const [mengirim, setMengirim] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [tanya, setTanya] = useState(false)

  const kirim = async () => {
    setMengirim(true)
    setGalat(null)
    try {
      onSelesai(await putuskanKasus(kasusId, { status, comment: komentar }), status)
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
      <div className="edm-aksi">
        <button
          type="button"
          className="btn btn--primary btn--sm"
          disabled={status === '' || mengirim}
          onClick={() => setTanya(true)}
        >
          {mengirim ? UMUM_EDM.memutuskan : KONFIRMASI_PUTUSAN_EDM.confirm}
        </button>
      </div>
      {tanya && (
        <KonfirmasiPutusan
          status={status}
          komentar={komentar}
          onTutup={() => setTanya(false)}
          onLanjut={() => {
            setTanya(false)
            void kirim()
          }}
        />
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
