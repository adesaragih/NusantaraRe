// Layar tabel acuan Treaty In — tiket 15.
//
// ⛔ INI SELURUH LAYAR MODUL INI HARI INI, dan itu disengaja. Papan tiket
// menyatakan `L-4` ("tidak ada spesifikasi layar di mana pun") dan melarang
// mengarang layar untuk memenuhi bentuk irisan tegak: layar yang dikarang
// dipakai sebagai kriteria selesai oleh orang yang tidak tahu ia karangan.
// Halaman ini TIDAK mengarang alur kontrak - ia memperlihatkan isi keenam tabel
// acuan yang tiket 15 buat, tidak lebih.

import { useEffect, useState } from 'react'

import { Gagal, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { ambilAcuan, HIMPUNAN_ACUAN, type Acuan, type HimpunanAcuan } from '../api'
import { ACUAN_TREATYIN, LABEL_HIMPUNAN, MENU_TREATYIN } from '../labels'

export default function AcuanTreatyIn() {
  const [himpunan, setHimpunan] = useState<HimpunanAcuan>('mata-uang')
  const [baris, setBaris] = useState<Acuan[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let dibuang = false
    setBaris(null)
    setGalat(null)
    ambilAcuan(himpunan)
      .then((h) => {
        if (!dibuang) setBaris(h)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [himpunan])

  const bersusun = himpunan === 'jenis-reasuransi'

  return (
    <div className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{MENU_TREATYIN.acuan}</h2>
      </header>
      <p>{ACUAN_TREATYIN.keterangan}</p>
      <StripTab
        tab={HIMPUNAN_ACUAN}
        aktif={himpunan}
        onPilih={setHimpunan}
        label={(t) => LABEL_HIMPUNAN[t]}
      />
      {galat !== null && <Gagal galat={galat} />}
      {baris === null && galat === null && <Memuat />}
      {baris !== null && baris.length === 0 && <p>{ACUAN_TREATYIN.kosong}</p>}
      {baris !== null && baris.length > 0 && (
        <table>
          <thead>
            <tr>
              <th scope="col">{ACUAN_TREATYIN.kolomKode}</th>
              <th scope="col">{ACUAN_TREATYIN.kolomNama}</th>
              <th scope="col">{ACUAN_TREATYIN.kolomAktif}</th>
              {bersusun && <th scope="col">{ACUAN_TREATYIN.kolomInduk}</th>}
            </tr>
          </thead>
          <tbody>
            {baris.map((b) => (
              <tr key={b.id}>
                <td>{b.kode}</td>
                <td>{b.nama}</td>
                <td>{b.aktif}</td>
                {bersusun && <td>{b.idInduk ?? ''}</td>}
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}
