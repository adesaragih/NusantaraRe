// Dialog Delete (b9624 -> `DeleteSummaryDetail` b9648, DeleteID=.ID b9664): menyebut jumlah baris R/I RISK DETAIL yang
// ikut terhapus (dibaca ulang dari server sebelum tombol aktif).

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambil, hapus, type Ringkasan } from '../api'
import { RK } from '../labels'

export default function KonfirmasiHapus({
  ringkasan,
  onTutup,
  onTerhapus,
}: {
  ringkasan: Ringkasan
  onTutup: () => void
  onTerhapus: (id: string, rincianTerhapus: number) => void
}) {
  const [jumlah, setJumlah] = useState<number | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    let batal = false
    ambil(ringkasan.id).then(
      (r) => {
        if (!batal) setJumlah(r.jumlahRincian ?? 0)
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
      (h) => onTerhapus(h.id, h.rincianTerhapus),
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
      },
    )
  }

  return (
    <Modal
      judul={RK.judulHapus(ringkasan.id)}
      onTutup={onTutup}
      onKirim={kirim}
      labelBatal={RK.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk || jumlah === null}>
          {sibuk ? RK.menghapus : RK.hapus}
        </button>
      }
    >
      {galat !== null && <Gagal galat={galat} />}
      {jumlah === null && galat === null && <Memuat pesan={RK.menghitung} />}
      {jumlah !== null && <p>{RK.tanyaHapus(ringkasan.usedby, jumlah)}</p>}
    </Modal>
  )
}
