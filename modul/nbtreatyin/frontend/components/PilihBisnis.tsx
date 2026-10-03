// Popup `Harness/BusinessAndSOBList` - grid RD `BrowseTreatyInDetail` dengan
// tombol "Choose" per baris (`SetValue_Act(ID=.ID)`).

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { daftarBisnis, type BarisKontrak } from '../api'
import { JUDUL, KOLOM_BISNIS, TOMBOL } from '../labels'

export default function PilihBisnis({ onPilih, onTutup }: { onPilih: (id: string) => void; onTutup: () => void }) {
  const [cari, setCari] = useState('')
  const [kueri, setKueri] = useState('')
  const [baris, setBaris] = useState<BarisKontrak[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let dibuang = false
    setBaris(null)
    setGalat(null)
    daftarBisnis(kueri)
      .then((b) => {
        if (!dibuang) setBaris(b)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [kueri])

  return (
    <Modal judul={JUDUL.pilihBisnis} onTutup={onTutup} penuh>
      <form
        className="nbti__saring"
        onSubmit={(e) => {
          e.preventDefault()
          setKueri(cari.trim())
        }}
      >
        <input className="field__input" value={cari} placeholder="TREATYID" onChange={(e) => setCari(e.target.value)} />
        <button type="submit" className="btn">
          {TOMBOL.filter}
        </button>
      </form>
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && baris.length === 0 && <Kosong pesan="—" />}
      {baris !== null && baris.length > 0 && (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col" />
                {KOLOM_BISNIS.map((k) => (
                  <th scope="col" key={k}>
                    {k}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {baris.map((b, i) => (
                <tr key={`${b.ID ?? ''}-${i}`}>
                  <td>
                    <button type="button" className="btn btn--primary" onClick={() => onPilih(b.ID ?? '')}>
                      {TOMBOL.choose}
                    </button>
                  </td>
                  {KOLOM_BISNIS.map((k) => (
                    <td key={k}>{b[k] ?? ''}</td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Modal>
  )
}
