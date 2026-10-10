// Tabel komite di bawah inbox Claim Fac In (modul Komite Claim Fac In TANPA menu, prompt tahap 2 §2 butir 2; perintah
// work owner 09-10-2026 "komite digabung ke menu klaimnya masing-masing"). Isinya kasus `KMT-` yang sedang menunggu
// workbasket / akun pelaku (worklist KomiteRouter, rute pinjaman `GET /api/komite-claim-fac-in/kasus`); tanpa switch,
// tidak ikut tab. Klik baris = layar komite DI TEMPAT (`onBukaModul`); Back / Submit kembali ke inbox ini dan daftarnya
// dibaca ulang (halaman awal dipasang ulang). Disalin dari `modul/claimnonprop/frontend/components/TabelKomite.tsx`
// (bukan impor), ditambah kolom jenis penyerahan (ADJUSTMENT / REJECT / CLOSE) dan status baris adjustment.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { daftarKomite, MODUL_KOMITE, type BarisKomite } from '../api'
import { tampilTanggal } from '../ketikTanggal'
import { CFI } from '../labels'

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
        setDaftar(d ?? [])
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
    if (onBukaModul?.(MODUL_KOMITE, id) !== true) setPesan(CFI.komiteTakTerpasang)
  }

  return (
    <section className="claimfacin__komite" aria-labelledby="claimfacin-judul-komite">
      <h3 id="claimfacin-judul-komite" className="inbox__judul">
        {CFI.judulKomite}
      </h3>
      {pesan && (
        <div className="alert alert--error" role="alert">
          {pesan}
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={CFI.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={CFI.kosongKomite} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel claimfacin__tabel">
          <thead>
            <tr>
              <th>{CFI.komiteKolomNo}</th>
              <th>{CFI.komiteKolomKlaimId}</th>
              <th>{CFI.komiteKolomNoKlaim}</th>
              <th>{CFI.komiteKolomJenis}</th>
              <th>{CFI.komiteKolomJabatan}</th>
              <th>{CFI.komiteKolomTingkat}</th>
              <th>{CFI.komiteKolomStatus}</th>
              <th>{CFI.komiteKolomTgl}</th>
            </tr>
          </thead>
          <tbody>
            {daftar.map((b) => (
              <tr
                key={b.kasusId}
                className="inbox__baris"
                tabIndex={0}
                onClick={() => buka(b.kasusId)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault()
                    buka(b.kasusId)
                  }
                }}
              >
                <td>{b.kasusId}</td>
                <td>{b.klaimId}</td>
                <td>{sel(b.noKlaim)}</td>
                <td>{sel(b.jenis)}</td>
                <td>{sel(b.jabatan)}</td>
                <td>{`${String(b.tingkat)} / ${String(b.komiteLoop)}`}</td>
                <td>{sel(b.statusBaris)}</td>
                <td>{sel(tampilTanggal(b.tglUpdate, 'tanggal-waktu'))}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
