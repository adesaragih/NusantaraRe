// Halaman awal Komite Claim Prop - daftar kerja penyetuju: worklist assignment "KomiteRouter" (`KomiteTreaty_Flow`,
// router `KomiteRouter` S6.1: kasus yang baris tangga pertamanya yang menunggu milik pelaku). Nol tombol "buat baru"
// (`pyCanCreateWorkObject=false`: kasus lahir dari Claim Prop). Klik baris = buka layar komite.

import { useEffect, useState } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { daftarKerja, type BarisKerja } from '../api'
import { KCP } from '../labels'
import { tampil } from '../nilai'
import KasusKomite from './KasusKomite'

/** Sel kosong ditandai. */
function sel(v: string) {
  return v.trim() === '' ? <span className="muted">—</span> : v
}

export default function DaftarKerja({
  onLihatBerkas,
}: {
  /** `PropsRute.onLihatBerkas` - tombol View more details. */
  onLihatBerkas?: (modul: string, id: string) => boolean
}) {
  const [daftar, setDaftar] = useState<BarisKerja[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [buka, setBuka] = useState<string | null>(null)
  const [segar, setSegar] = useState(0)

  useEffect(() => {
    if (buka) return
    let aktif = true
    setDaftar(null)
    daftarKerja().then(
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
  }, [buka, segar])

  if (buka) {
    return (
      <KasusKomite
        id={buka}
        onLihatBerkas={onLihatBerkas}
        onKembali={() => {
          setBuka(null)
          setSegar((s) => s + 1)
        }}
      />
    )
  }

  return (
    <section className="inbox komiteclaimprop__akar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{KCP.judul}</h2>
      </header>
      {galat !== null && <Gagal galat={galat} />}
      {daftar === null && galat === null && <Memuat pesan={KCP.memuat} />}
      {daftar !== null && daftar.length === 0 && <Kosong pesan={KCP.kosong} />}
      {daftar !== null && daftar.length > 0 && (
        <table className="inbox__tabel komiteclaimprop__tabel">
          <thead>
            <tr>
              <th>{KCP.kolomKasus}</th>
              <th>{KCP.kolomTgl}</th>
              <th>{KCP.kolomStatus}</th>
              <th>{KCP.kolomKlaim}</th>
              <th>{KCP.kolomTingkat}</th>
              <th>{KCP.kolomJabatan}</th>
              <th className="komiteclaimprop__angka">{KCP.kolomNilai}</th>
              <th>{KCP.kolomMataUang}</th>
              <th>{KCP.kolomStatusBaris}</th>
            </tr>
          </thead>
          <tbody>
            {daftar.map((b) => (
              <tr key={b.kasusId} className="inbox__baris" onClick={() => setBuka(b.kasusId)}>
                <td>{b.kasusId}</td>
                <td>{tampil(b.tglUpdate, 'tanggalJam')}</td>
                <td>{sel(b.statusWork)}</td>
                <td>{b.noKlaim ? `${b.noKlaim} / ${b.klaimId}` : b.klaimId}</td>
                <td>{`${String(b.tingkat)} / ${String(b.komiteLoop)}`}</td>
                <td>{sel(b.jabatan)}</td>
                <td className="komiteclaimprop__angka">{sel(tampil(b.nilai, 'angka'))}</td>
                <td>{sel(b.mataUang)}</td>
                <td>{sel(b.statusBaris)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
