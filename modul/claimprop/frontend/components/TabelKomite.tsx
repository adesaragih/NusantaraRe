// Tabel komite di bawah inbox Claim Prop (keputusan work owner 09-10-2026: menu Komite Claim Prop dibuang, "semua komite
// muncul di inbox klaim prop dengan tambahan tabel di bawah"). Isinya kasus komite yang sedang menunggu workbasket /
// akun pelaku (worklist KomiteRouter, rute pinjaman `GET /api/komite-claim-prop/kasus`); tanpa switch, tidak ikut tab.
// Klik baris = layar komite DI TEMPAT (`onBukaModul`, "jangan pop up, langsung buka komitenya"); Back / Submit kembali
// ke inbox ini dan daftarnya dibaca ulang.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { daftarKomite, MODUL_KOMITE, type BarisKomite } from '../api'
import { tampilTanggal } from '../ketikTanggal'
import { CP } from '../labels'
import { tampilAngka } from '../nilai'

/** Sel kosong ditandai. */
function sel(v: string) {
  return v.trim() === '' ? <span className="muted">—</span> : v
}

export default function TabelKomite({ onBukaModul }: { onBukaModul?: (modul: string, id: string) => boolean }) {
  const [daftar, setDaftar] = useState<BarisKomite[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [pesan, setPesan] = useState<string | null>(null)

  useEffect(() => {
    let aktif = true
    daftarKomite().then(
      (d) => {
        if (!aktif) return
        setDaftar(d)
        setGalat(null)
      },
      (g: unknown) => {
        if (aktif) setGalat(g)
      },
    )
    return () => {
      aktif = false
    }
  }, [])

  const buka = (id: string) => {
    setPesan(null)
    if (onBukaModul?.(MODUL_KOMITE, id) !== true) setPesan(CP.komiteTakTerpasang)
  }

  return (
    <section className="claimprop__komite" aria-labelledby="claimprop-judul-komite">
      <h3 id="claimprop-judul-komite" className="inbox__judul">
        {CP.judulKomite}
      </h3>
      {pesan && (
        <div className="alert alert--error" role="alert">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={CP.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={CP.kosongKomite} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel claimprop__tabel">
          <thead>
            <tr>
              <th>{CP.kolomID}</th>
              <th>{CP.komiteKolomTgl}</th>
              <th>{CP.komiteKolomKlaim}</th>
              <th>{CP.komiteKolomTingkat}</th>
              <th>{CP.komiteKolomJabatan}</th>
              <th>
                <span className="claimprop__angka">{CP.komiteKolomNilai}</span>
              </th>
              <th>{CP.komiteKolomMataUang}</th>
            </tr>
          </thead>
          <tbody>
            {daftar.map((b) => (
              <tr key={b.kasusId} className="inbox__baris" onClick={() => buka(b.kasusId)}>
                <td>{b.kasusId}</td>
                <td>{tampilTanggal(b.tglUpdate, 'tanggal-waktu')}</td>
                <td>{b.noKlaim ? `${b.noKlaim} / ${b.klaimId}` : b.klaimId}</td>
                <td>{`${String(b.tingkat)} / ${String(b.komiteLoop)}`}</td>
                <td>{sel(b.jabatan)}</td>
                <td>
                  <span className="claimprop__angka">{sel(tampilAngka(b.nilai))}</span>
                </td>
                <td>{sel(b.mataUang)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
