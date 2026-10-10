// Popup "View Upload" (b3467 -> `ViewCSVResult_RIRate` b3486): hasil urai server tanpa menulis - galat per baris,
// ringkasan tujuan (lama / baru), dan baris sah.

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { pratinjauUnggah, type HasilUnggah as Hasil } from '../api'
import { RR } from '../labels'

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
    <Modal judul={RR.lihatUnggah} onTutup={onTutup} labelBatal={RR.tutup} lebar>
      {galat !== null && <Gagal galat={galat} />}
      {hasil === null && galat === null && <Memuat pesan={RR.memuat} />}
      {hasil !== null && (
        <div className="riratelife__form">
          <div className={hasil.sah ? 'alert alert--ok' : 'alert alert--warn'} role="status">
            {RR.barisSah(hasil.baris.length)} · {RR.barisGalat(hasil.galat.length)}.{' '}
            {hasil.sah ? RR.siapSimpan : RR.tidakSiap}
          </div>
          {hasil.galat.length > 0 && (
            <div className="riratelife__gulir">
              <table className="inbox__tabel riratelife__tabel">
                <thead>
                  <tr>
                    <th>{RR.barisKe}</th>
                    <th>{RR.pesan}</th>
                  </tr>
                </thead>
                <tbody>
                  {hasil.galat.map((g, i) => (
                    <tr key={`${g.baris}-${i}`} className="inbox__baris">
                      <td>{g.baris}</td>
                      <td className="riratelife__galat">{g.pesan}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {hasil.ringkasan.length > 0 && (
            <p className="muted riratelife__catatan">
              {RR.tujuan}:{' '}
              {hasil.ringkasan
                .map((r) => `${r.usedby} (${r.baru ? RR.ringkasanBaru : `ID ${r.id}`}, ${r.jumlah})`)
                .join('; ')}
            </p>
          )}
          {hasil.baris.length > 0 && (
            <div className="riratelife__gulir">
              <table className="inbox__tabel riratelife__tabel">
                <thead>
                  <tr>
                    <th>{RR.barisKe}</th>
                    <th>{RR.usedby}</th>
                    <th>{RR.contract}</th>
                    <th>{RR.gender}</th>
                    <th>{RR.age}</th>
                    <th className="riratelife__angka">{RR.rate}</th>
                  </tr>
                </thead>
                <tbody>
                  {hasil.baris.map((b) => (
                    <tr key={b.baris} className="inbox__baris">
                      <td>{b.baris}</td>
                      <td>{b.usedby}</td>
                      <td>{b.contract}</td>
                      <td>{b.gender}</td>
                      <td>{b.age}</td>
                      <td className="riratelife__angka">{b.rate}</td>
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
