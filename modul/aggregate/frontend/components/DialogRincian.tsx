// Rincian satu baris daftar (klik ganda `GridDasbordAgg` -> `DetailAggregate_Act`): seluruh baris AGGREGATE
// berkunci sama, hanya dibaca.

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilRincian, type Baris, type Kunci } from '../api'
import { AG } from '../labels'
import GridAggregate from './GridAggregate'

export default function DialogRincian({ kunci, onTutup }: { kunci: Kunci; onTutup: () => void }) {
  const [baris, setBaris] = useState<Baris[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    ambilRincian(kunci).then(
      (r) => {
        if (hidup) setBaris(r.baris)
      },
      (g: unknown) => {
        if (hidup) setGalat(g)
      },
    )
    return () => {
      hidup = false
    }
  }, [kunci])

  return (
    <Modal judul={`${AG.judulRincian} — ${kunci.cedingName || kunci.cedingCode} · ${kunci.asAt}`} onTutup={onTutup} labelBatal={AG.tutup} penuh>
      <div className="aggregate__rincian">
        {galat !== null && <Gagal galat={galat} />}
        {baris === null && galat === null && <Memuat pesan={AG.memuatRincian} />}
        {baris !== null && <GridAggregate baris={baris} sertaID />}
      </div>
    </Modal>
  )
}
