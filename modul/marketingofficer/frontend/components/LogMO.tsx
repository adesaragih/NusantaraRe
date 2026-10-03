// Log perubahan satu MO dari MARKETINGOFFICER_LOG (permintaan work owner 03-10-2026). Terbaru di atas; setiap
// butir = satu UPDATE: kapan, oleh siapa, kolom apa berubah dari apa menjadi apa.

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilLog, type BarisMO, type Riwayat } from '../api'
import { labelRuas, teksNilaiLog } from '../aturan'
import { MO } from '../labels'

export default function LogMO({ baris, onTutup }: { baris: BarisMO; onTutup: () => void }) {
  const [riwayat, setRiwayat] = useState<Riwayat | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    ambilLog(baris.id).then(
      (r) => {
        if (hidup) setRiwayat(r)
      },
      (g: unknown) => {
        if (hidup) setGalat(g)
      },
    )
    return () => {
      hidup = false
    }
  }, [baris.id])

  return (
    <Modal judul={MO.judulLog(baris.clientName, baris.id)} onTutup={onTutup} labelBatal={MO.tutup} lebar>
      <p className="muted marketingofficer__catatan">{MO.catatanLog}</p>
      {galat !== null && <Gagal galat={galat} />}
      {riwayat === null && galat === null && <Memuat pesan={MO.memuatLog} />}
      {riwayat !== null && riwayat.perubahan.length === 0 && <p className="muted">{MO.logKosong}</p>}
      {riwayat !== null && riwayat.perubahan.length > 0 && (
        <ol className="marketingofficer__log">
          {riwayat.perubahan.map((p, i) => (
            <li key={`${i}-${p.waktu}`} className="marketingofficer__log-butir">
              <div className="marketingofficer__log-kepala">
                <strong>{p.waktu === '' ? MO.waktuTakDiketahui : p.waktu}</strong>
                {p.oleh !== '' && <span className="muted">{` ${MO.oleh} ${p.oleh}`}</span>}
                {p.perkiraan && <span className="marketingofficer__tanda">{MO.perkiraan}</span>}
              </div>
              {p.ruas.length === 0 ? (
                <p className="muted marketingofficer__catatan">{MO.tanpaPerubahan}</p>
              ) : (
                <table className="marketingofficer__log-tabel">
                  <tbody>
                    {p.ruas.map((r) => (
                      <tr key={r.kolom}>
                        <th scope="row">{labelRuas(r)}</th>
                        <td>{teksNilaiLog(r.kolom, r.sebelum)}</td>
                        <td aria-hidden="true">→</td>
                        <td>{teksNilaiLog(r.kolom, r.sesudah)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </li>
          ))}
        </ol>
      )}
    </Modal>
  )
}
