// Jendela "Choose Ceding" / "Choose Source of Business" — local action
// `TreatyInSearchReinsured` / `TreatyInSearchSoB` (`Section/` dengan nama yang
// sama, korpus Treaty In Adjustment).
//
//   sel 5   kotak kata kunci       InputData.CARI2
//   sel 8   Search                 DataTransform `TreatyInCedingSetToUppercase`:
//                                  CARI2 = @toUpperCase(input)
//   grid    ID · Name              RD `BrowseAgentNusaRe_RD`, saringan
//                                  `.ClientName Contains Param.ClientName`
//   sel 27  Choose (per baris)     DataTransform `TreatyInSetReinsured`
//                                  (type, name = .ClientName, id = .ID), lalu
//                                  menutup jendela
//
// ⭐ Daftar agennya rute Treaty In `/warisan/cedant` — saringan RD yang sama
// (agen aktif, bukan LIFE, tanpa anak, ber-CLIENTID), dibaca dari AGENT.
// Pencocokan `Contains` dijalankan di layar atas daftar itu.
//
// ⛔ Memilih HANYA mengisi keadaan layar (nama + pengenal). Nol tulisan ke
// basis data sebelum Save.

import { useEffect, useState } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { ambilAgenTreatyIn, type PilihanTreatyIn } from '../api'
import { PILIH_AGEN } from '../labelsPenyesuaian'

/** `.ClientName Contains kata` — kata sudah huruf besar (DataTransform Search). */
export function saringAgen(daftar: readonly PilihanTreatyIn[], kata: string): PilihanTreatyIn[] {
  return kata === '' ? [...daftar] : daftar.filter((a) => a.nama.includes(kata))
}

export default function PilihAgen({
  judul,
  onPilih,
  onTutup,
}: {
  /** Label flow action — `TreatyInSearchReinsured` / `TreatyInSearchSoB`. */
  judul: string
  onPilih: (nama: string, id: string) => void
  onTutup: () => void
}) {
  const [daftar, setDaftar] = useState<PilihanTreatyIn[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [ketik, setKetik] = useState('')
  const [kata, setKata] = useState('')
  useEffect(() => {
    let dibuang = false
    ambilAgenTreatyIn()
      .then((d) => {
        if (!dibuang) setDaftar(d)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [])
  const hasil = saringAgen(daftar ?? [], kata)
  return (
    <Modal judul={judul} onTutup={onTutup} lebar>
      <div className="tria__cari">
        <input
          className="field__input"
          type="text"
          aria-label={PILIH_AGEN.kataKunci}
          value={ketik}
          onChange={(e) => {
            setKetik(e.target.value)
          }}
        />
        <button
          type="button"
          className="btn btn--primary btn--sm"
          onClick={() => {
            // `TreatyInCedingSetToUppercase`: CARI2 = @toUpperCase(input).
            const besar = ketik.toUpperCase()
            setKetik(besar)
            setKata(besar)
          }}
        >
          {PILIH_AGEN.cari}
        </button>
      </div>
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat />}
      {daftar !== null && (
        <div className="table-wrap">
          <table className="tria__tabel">
            <thead>
              <tr>
                <th scope="col">{PILIH_AGEN.kolomId}</th>
                <th scope="col">{PILIH_AGEN.kolomNama}</th>
                <th scope="col" aria-label={PILIH_AGEN.pilih} />
              </tr>
            </thead>
            <tbody>
              {hasil.length === 0 && (
                <tr>
                  <td colSpan={3}>{PILIH_AGEN.tanpaBaris}</td>
                </tr>
              )}
              {hasil.map((a) => (
                <tr key={`${a.id}|${a.nama}`}>
                  <td>{a.id}</td>
                  <td>{a.nama}</td>
                  <td>
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        onPilih(a.nama, a.id)
                        onTutup()
                      }}
                    >
                      {PILIH_AGEN.pilih}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </Modal>
  )
}
