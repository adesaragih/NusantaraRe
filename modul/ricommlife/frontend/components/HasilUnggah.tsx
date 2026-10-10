// Popup "View Upload" (b3706 -> `ViewCSVResult_RIComm` b3725): hasil urai server tanpa menulis - galat per baris,
// ringkasan tujuan (lama / baru), dan baris sah.

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { pratinjauUnggah, type HasilUnggah as Hasil } from '../api'
import { tampilDesimal } from '../aturan'
import { RC } from '../labels'

export default function HasilUnggah({ isi, onTutup }: { isi: string; onTutup: () => void }) {
  const [hasil, setHasil] = useState<Hasil | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let batal = false
    pratinjauUnggah(isi).then(
      (h) => {
        if (!batal) setHasil(h)
      },
      (g: unknown) => {
        if (!batal) setGalat(g)
      },
    )
    return () => {
      batal = true
    }
  }, [isi])

  return (
    <Modal judul={RC.lihatUnggah} onTutup={onTutup} labelBatal={RC.tutup} lebar>
      {galat !== null && <Gagal galat={galat} />}
      {hasil === null && galat === null && <Memuat pesan={RC.memuat} />}
      {hasil !== null && (
        <div className="ricommlife__form">
          <div className={hasil.sah ? 'alert alert--ok' : 'alert alert--warn'} role="status">
            {RC.barisSah(hasil.baris.length)} · {RC.barisGalat(hasil.galat.length)}.{' '}
            {hasil.sah ? RC.siapSimpan : RC.tidakSiap}
          </div>
          {hasil.galat.length > 0 && (
            <div className="ricommlife__gulir">
              <table className="inbox__tabel ricommlife__tabel">
                <thead>
                  <tr>
                    <th>{RC.barisKe}</th>
                    <th>{RC.pesan}</th>
                  </tr>
                </thead>
                <tbody>
                  {hasil.galat.map((g, i) => (
                    <tr key={`${g.baris}-${i}`} className="inbox__baris">
                      <td>{g.baris}</td>
                      <td className="ricommlife__galat">{g.pesan}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {hasil.ringkasan.length > 0 && (
            <p className="muted ricommlife__catatan">
              {RC.tujuan}:{' '}
              {hasil.ringkasan
                .map((r) => `${r.usedby} (${r.baru ? RC.ringkasanBaru : `ID ${r.id}`}, ${r.jumlah})`)
                .join('; ')}
            </p>
          )}
          {hasil.baris.length > 0 && (
            <div className="ricommlife__gulir">
              <table className="inbox__tabel ricommlife__tabel">
                <thead>
                  <tr>
                    <th>{RC.barisKe}</th>
                    <th>{RC.usedby}</th>
                    <th>{RC.contract}</th>
                    <th>{RC.year}</th>
                    <th className="ricommlife__angka">{RC.comm}</th>
                  </tr>
                </thead>
                <tbody>
                  {hasil.baris.map((b) => (
                    <tr key={b.baris} className="inbox__baris">
                      <td>{b.baris}</td>
                      <td>{b.usedby}</td>
                      <td>{b.contract}</td>
                      <td>{b.year}</td>
                      <td className="ricommlife__angka">{tampilDesimal(b.comm)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}
    </Modal>
  )
}
