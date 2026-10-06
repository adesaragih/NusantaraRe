// Popup "Rate Detail" (b11446; tombol Detail b11398 -> `setIDUsedBy_Act` + harness `InboxRIRate`): baris RATE_LIFE satu
// ringkasan, baca saja, 50 per halaman. Jalur tulis rate = Upload CSV.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilRate, type Halaman, type Rate, type Ringkasan } from '../api'
import { jumlahHalaman, nomorAwal } from '../aturan'
import { RR } from '../labels'

export default function RateDetail({ ringkasan, onTutup }: { ringkasan: Ringkasan; onTutup: () => void }) {
  const [halaman, setHalaman] = useState(1)
  const [data, setData] = useState<Halaman<Rate> | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let batal = false
    setData(null)
    ambilRate(ringkasan.id, halaman).then(
      (d) => {
        if (batal) return
        setData(d)
        setGalat(null)
      },
      (g: unknown) => {
        if (!batal) setGalat(g)
      },
    )
    return () => {
      batal = true
    }
  }, [ringkasan.id, halaman])

  const dari = data === null ? 1 : jumlahHalaman(data.total, data.ukuran)

  return (
    <Modal judul={`${RR.judulDetail} — ${ringkasan.id} ${ringkasan.usedby}`} onTutup={onTutup} labelBatal={RR.tutup} lebar>
      {galat !== null && <Gagal galat={galat} />}
      {data === null && galat === null && <Memuat pesan={RR.memuat} />}
      {data !== null && data.daftar.length === 0 && <Kosong pesan={RR.kosongDetail} />}
      {data !== null && data.daftar.length > 0 && (
        <>
          <div className="riratelife__gulir">
            <table className="inbox__tabel riratelife__tabel">
              <thead>
                <tr>
                  <th>No</th>
                  <th>{RR.usedby}</th>
                  <th>{RR.contract}</th>
                  <th>{RR.gender}</th>
                  <th>{RR.age}</th>
                  <th className="riratelife__angka">{RR.rate}</th>
                </tr>
              </thead>
              <tbody>
                {data.daftar.map((r, i) => (
                  <tr key={r.id} className="inbox__baris">
                    <td>{nomorAwal(data.halaman, data.ukuran) + i}</td>
                    <td>{r.usedby}</td>
                    <td>{r.contract}</td>
                    <td>{r.gender}</td>
                    <td>{r.age}</td>
                    <td className="riratelife__angka">{r.rate}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="riratelife__halaman">
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman <= 1} onClick={() => setHalaman((h) => h - 1)}>
              {RR.sebelum}
            </button>
            <span className="muted">{RR.halaman(data.halaman, dari)}</span>
            <button type="button" className="btn btn--ghost btn--sm" disabled={halaman >= dari} onClick={() => setHalaman((h) => h + 1)}>
              {RR.sesudah}
            </button>
          </div>
        </>
      )}
    </Modal>
  )
}
