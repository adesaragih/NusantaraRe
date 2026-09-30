// Layar kontrak treaty dari menu — tiket 04 Treaty Contract Out.
//
// Butir menu `InboxTreatyContractReinsType` (`Harness/InboxTreatyContractReinsType.xml`
// b151). Di Pega harness ini dibuka sebagai popup BERKONTEKS dari tombol
// `ReinsType` baris tahun treaty; dibuka langsung dari portal ia tidak punya
// tahun. Di sini pemakai memilih tahunnya lebih dulu (`[tidak ada di korpus]`),
// lalu editor yang SAMA dengan popup itu tampil.

import { useEffect, useState } from 'react'

import { KONTRAK_TCO } from '../labels'
import PanelKontrakTahun from '../components/PanelKontrakTahun'
import { Gagal, Memuat, Pilih } from '../../../../inti/frontend/components/ui/dasar'
import { ambilTahunTreaty, type TahunTreaty } from '../api'
import { labelProporsi } from '../proporsi'

/** Label pilihan tahun: tahun treaty · grup · jenis. */
export function labelTahun(t: TahunTreaty): string {
  return [t.treatyYear, t.treatyGroupName, labelProporsi(t.proportion)].filter((v) => v.trim() !== '').join(' · ') || t.id
}

export default function InboxTreatyContractReinsType() {
  const [tahun, setTahun] = useState<TahunTreaty[] | null>(null)
  const [pilih, setPilih] = useState('')
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    // `BrowseTreatyYear_RD` pyMaxRecords 500 b757 - batas yang sama.
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
    <div className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{KONTRAK_TCO.judul}</h2>
      </header>
      {galat !== null && <Gagal galat={galat} />}
      {tahun === null && galat === null && <Memuat />}
      {tahun !== null && (
        <div className="form-grid">
          <Pilih
            label={KONTRAK_TCO.pilihTahun}
            value={pilih}
            onChange={setPilih}
            opsi={tahun.map((t) => ({ value: t.id, label: labelTahun(t) }))}
          />
        </div>
      )}
      {terpilih !== null && <PanelKontrakTahun key={terpilih.id} tahun={terpilih} />}
    </div>
  )
}
