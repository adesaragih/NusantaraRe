// Halaman Master Data (rencana 01, M-1): satu menu, satu tab per master (urutan `GET /api/masterdata`: Nation,
// Province, City, District, CZone, Accumulated Type, Accumulation, Object Item Type). Template Kelola User.

import { useEffect, useState } from 'react'

import { Gagal, Memuat, StripTab } from '../../../../inti/frontend/components/ui/dasar'
import { daftarMaster, type MetaMaster } from '../api'
import DaftarMaster from '../components/DaftarMaster'
import { TEKS } from '../labels'
import { KELOMPOK_MASTERDATA } from '../menu'

export default function HalamanMasterData() {
  const [master, setMaster] = useState<MetaMaster[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [aktif, setAktif] = useState('')

  useEffect(() => {
    let batal = false
    daftarMaster().then(
      (h) => {
        if (batal) return
        setMaster(h.tabel)
        setAktif((a) => a || (h.tabel[0]?.kunci ?? ''))
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [])

  const meta = master?.find((m) => m.kunci === aktif)
  return (
    <section className="inbox masterdata">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{KELOMPOK_MASTERDATA}</h2>
      </header>
      <p className="muted">{TEKS.sub}</p>
      {master === null && galat === null && <Memuat pesan={TEKS.memuat} />}
      <Gagal galat={galat} />
      {master !== null && master.length > 0 && (
        <StripTab
          tab={master.map((m) => m.kunci)}
          aktif={aktif}
          onPilih={setAktif}
          label={(k) => master.find((m) => m.kunci === k)?.judul ?? k}
        />
      )}
      {meta && <DaftarMaster key={meta.kunci} meta={meta} />}
    </section>
  )
}
