// Dialog Delete (b12238 -> `DeleteSummaryDetail`, DeleteID=.ID): menyebut jumlah baris rate yang ikut terhapus
// (dibaca ulang dari server sebelum tombol aktif).

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambil, hapus, type Ringkasan } from '../api'
import { RR } from '../labels'

export default function KonfirmasiHapus({
  ringkasan,
  onTutup,
  onTerhapus,
}: {
  ringkasan: Ringkasan
  onTutup: () => void
  onTerhapus: (id: string, rateTerhapus: number) => void
}) {
  const [jumlah, setJumlah] = useState<number | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  useEffect(() => {
    let batal = false
    ambil(ringkasan.id).then(
      (r) => {
        if (!batal) setJumlah(r.jumlahRate ?? 0)
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
      (h) => onTerhapus(h.id, h.rateTerhapus),
      (g: unknown) => {
        setSibuk(false)
        setGalat(g)
      },
    )
  }

  return (
    <Modal
      judul={RR.judulHapus(ringkasan.id)}
      onTutup={onTutup}
      onKirim={kirim}
      labelBatal={RR.batal}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk || jumlah === null}>
          {sibuk ? RR.menghapus : RR.hapus}
        </button>
      }
    >
      {galat !== null && <Gagal galat={galat} />}
      {jumlah === null && galat === null && <Memuat pesan={RR.menghitung} />}
      {jumlah !== null && <p>{RR.tanyaHapus(ringkasan.usedby, jumlah)}</p>}
    </Modal>
  )
}
