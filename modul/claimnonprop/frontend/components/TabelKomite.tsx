// Tabel komite di bawah inbox Claim Non Prop (perintah work owner 09-10-2026, pola Claim Prop: menu Komite Claim Non Prop
// disembunyikan, kasusnya tampil di bawah inbox). Isinya kasus komite yang sedang menunggu workbasket / akun pelaku
// (worklist KomiteRouter, rute pinjaman `GET /api/komite-claim-non-prop/kasus`); tanpa switch, tidak ikut tab. Klik
// baris = layar komite DI TEMPAT (`onBukaModul`); Back / Submit kembali ke inbox ini dan daftarnya dibaca ulang.
// Disalin dari `modul/claimprop/frontend/components/TabelKomite.tsx` (bukan impor), tanpa kolom nilai / mata uang
// (akseptasi Non Prop bernilai per layer).

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { daftarKomite, MODUL_KOMITE, type BarisKomite } from '../api'
import { tampilTanggal } from '../ketikTanggal'
import { CNP } from '../labels'

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
    if (onBukaModul?.(MODUL_KOMITE, id) !== true) setPesan(CNP.komiteTakTerpasang)
  }

  return (
    <section className="claimnonprop__komite" aria-labelledby="claimnonprop-judul-komite">
      <h3 id="claimnonprop-judul-komite" className="inbox__judul">
        {CNP.judulKomite}
      </h3>
      {pesan && (
        <div className="alert alert--error" role="alert">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={CNP.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={CNP.kosongKomite} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel claimnonprop__tabel">
          <thead>
            <tr>
              <th>{CNP.kolomID}</th>
              <th>{CNP.komiteKolomTgl}</th>
              <th>{CNP.komiteKolomKlaim}</th>
              <th>{CNP.komiteKolomTingkat}</th>
              <th>{CNP.komiteKolomJabatan}</th>
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
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
