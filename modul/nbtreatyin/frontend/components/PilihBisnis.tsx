// Popup `Harness/BusinessAndSOBList` - grid RD `BrowseTreatyInDetail` dengan
// tombol "Choose" per baris (`SetValue_Act(ID=.ID)`). Kolom `*VALUE` = pxCurrency
// tanpa `pyDecimalPlaces` -> sajian pola inti (K14, `sajian.ts`).

import { useState } from 'react'

import { Gagal, Kosong, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { useAmbil } from '../ambil'
import { daftarBisnis } from '../api'
import { JUDUL, KOLOM_BISNIS, TOMBOL } from '../labels'
import { sajikan } from '../sajian'

export default function PilihBisnis({ onPilih, onTutup }: { onPilih: (id: string) => void; onTutup: () => void }) {
  const [cari, setCari] = useState('')
  const [kueri, setKueri] = useState('')
  const { data: baris, galat } = useAmbil(() => daftarBisnis(kueri), [kueri])

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
                    <td key={k}>{k.endsWith('VALUE') ? sajikan(b[k] ?? '', {}) : (b[k] ?? '')}</td>
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
