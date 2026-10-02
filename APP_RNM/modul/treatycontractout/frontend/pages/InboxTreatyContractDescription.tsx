// Layar klausul dari menu — tiket 08 Treaty Contract Out.
//
// Butir menu `InboxTreatyContractDescription` (`Harness/InboxTreatyContractDescription.xml`
// b359). Di Pega harness ini dibuka BERKONTEKS dari tombol `List Description`
// baris tahun treaty (b22196); dibuka langsung dari portal ia tidak punya
// tahun. Seperti layar kontrak (tiket 04), pemakai memilih tahunnya lebih dulu
// (`[tidak ada di korpus]`), lalu panel yang SAMA dengan tombol itu tampil.

import { useEffect, useState } from 'react'

import { JUDUL_TAMPIL_TCO, KLAUSUL_TCO } from '../labels'
import PanelKlausulTahun from '../components/PanelKlausulTahun'
import { Gagal, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { ambilTahunTreaty, type TahunTreaty } from '../api'
import { labelTahun } from './InboxTreatyContractReinsType'

export default function InboxTreatyContractDescription() {
  const [tahun, setTahun] = useState<TahunTreaty[] | null>(null)
  const [pilih, setPilih] = useState('')
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    // `BrowseTreatyYear_RD` pyMaxRecords 500 b757 - batas yang sama dengan tiket 04.
    ambilTahunTreaty(1, 500)
      .then((h) => {
        setTahun(h.baris)
      })
      .catch((e: unknown) => {
        setGalat(e)
      })
  }, [])

  const terpilih = tahun?.find((t) => t.id === pilih) ?? null

  return (
    <div className="inbox tco">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{JUDUL_TAMPIL_TCO.deskripsi}</h2>
      </header>
      {galat !== null && <Gagal galat={galat} />}
      {tahun === null && galat === null && <Memuat />}
      {tahun !== null && (
        <div className="form-grid">
          <Pilih
            label={KLAUSUL_TCO.pilihTahun}
            value={pilih}
            onChange={setPilih}
            opsi={tahun.map((t) => ({ value: t.id, label: labelTahun(t) }))}
          />
        </div>
      )}
      {terpilih !== null && <PanelKlausulTahun key={terpilih.id} tahun={terpilih} />}
    </div>
  )
}
