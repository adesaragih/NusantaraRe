// Popup "View Upload" (`InboxSummaryRIRisk` b3676 -> harness `ViewCSVResult_RIRisk` b3695): hasil urai server tanpa menulis - galat per baris,
// ringkasan tujuan (lama / baru), dan baris sah.

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { pratinjauUnggah, type HasilUnggah as Hasil } from '../api'
import { tampilDesimal } from '../aturan'
import { RK } from '../labels'

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
    <Modal judul={RK.lihatUnggah} onTutup={onTutup} labelBatal={RK.tutup} lebar>
      {galat !== null && <Gagal galat={galat} />}
      {hasil === null && galat === null && <Memuat pesan={RK.memuat} />}
      {hasil !== null && (
        <div className="ririsklife__form">
          <div className={hasil.sah ? 'alert alert--ok' : 'alert alert--warn'} role="status">
            {RK.barisSah(hasil.baris.length)} · {RK.barisGalat(hasil.galat.length)}.{' '}
            {hasil.sah ? RK.siapSimpan : RK.tidakSiap}
          </div>
          {hasil.galat.length > 0 && (
            <div className="ririsklife__gulir">
              <table className="inbox__tabel ririsklife__tabel">
                <thead>
                  <tr>
                    <th>{RK.barisKe}</th>
                    <th>{RK.pesan}</th>
                  </tr>
                </thead>
                <tbody>
                  {hasil.galat.map((g, i) => (
                    <tr key={`${g.baris}-${i}`} className="inbox__baris">
                      <td>{g.baris}</td>
                      <td className="ririsklife__galat">{g.pesan}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {hasil.ringkasan.length > 0 && (
            <p className="muted ririsklife__catatan">
              {RK.tujuan}:{' '}
              {hasil.ringkasan
                .map((r) => `${r.usedby} (${r.baru ? RK.ringkasanBaru : `ID ${r.id}`}, ${r.jumlah})`)
                .join('; ')}
            </p>
          )}
          {hasil.baris.length > 0 && (
            <div className="ririsklife__gulir">
              <table className="inbox__tabel ririsklife__tabel">
                <thead>
                  <tr>
                    <th>{RK.barisKe}</th>
                    <th>{RK.kolomUsedby}</th>
                    <th>{RK.contract}</th>
                    <th>{RK.year}</th>
                    <th>{RK.month}</th>
                    <th className="ririsklife__angka">{RK.risk}</th>
                  </tr>
                </thead>
                <tbody>
                  {hasil.baris.map((b) => (
                    <tr key={b.baris} className="inbox__baris">
                      <td>{b.baris}</td>
                      <td>{b.usedby}</td>
                      <td>{b.contract}</td>
                      <td>{b.year}</td>
                      <td>{b.month}</td>
                      <td className="ririsklife__angka">{tampilDesimal(b.risk)}</td>
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
