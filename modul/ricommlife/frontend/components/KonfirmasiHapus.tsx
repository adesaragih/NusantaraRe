// Dialog Delete (b9680 -> `DeleteSummaryDetail` b9704, DeleteID=.ID b9719): menyebut jumlah baris R/I COMM DETAIL yang
// ikut terhapus (dibaca ulang dari server sebelum tombol aktif).

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambil, hapus, type Ringkasan } from '../api'
import { RC } from '../labels'

export default function KonfirmasiHapus({
  ringkasan,
  onTutup,
  onTerhapus,
}: {
  ringkasan: Ringkasan
  onTutup: () => void
  onTerhapus: (id: string, komisiTerhapus: number) => void
}) {
  const [jumlah, setJumlah] = useState<number | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    let batal = false
    ambil(ringkasan.id).then(
      (r) => {
        if (!batal) setJumlah(r.jumlahKomisi ?? 0)
      },
      (g: unknown) => {
        if (!batal) setGalat(g)
      },
    )
    return () => {
      batal = true
    }
  }, [ringkasan.id])

  const kirim = () => {
    if (sibuk || jumlah === null) return
    setSibuk(true)
    hapus(ringkasan.id).then(
      (h) => onTerhapus(h.id, h.komisiTerhapus),
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
      },
    )
  }

  return (
    <Modal
      judul={RC.judulHapus(ringkasan.id)}
      onTutup={onTutup}
      onKirim={kirim}
      labelBatal={RC.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk || jumlah === null}>
          {sibuk ? RC.menghapus : RC.hapus}
        </button>
      }
    >
      {galat !== null && <Gagal galat={galat} />}
      {jumlah === null && galat === null && <Memuat pesan={RC.menghitung} />}
      {jumlah !== null && <p>{RC.tanyaHapus(ringkasan.usedby, jumlah)}</p>}
    </Modal>
  )
}
