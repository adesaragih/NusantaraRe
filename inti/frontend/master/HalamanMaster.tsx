// Halaman satu master (inti) - dipasang modul master masing-masing (`modul/master<nama>/frontend/rute.tsx`) dengan
// prefix rutenya sendiri (keputusan work owner 04-10-2026: satu menu per master di grup MASTER, "8 modul terpisah").
// Template Kelola User: kepala, keterangan, lalu daftar.

import { useEffect, useMemo, useState } from 'react'

import { Gagal, Memuat } from '../components/ui/dasar'
import { klienMaster, type MetaMaster } from './api'
import DaftarMaster from './DaftarMaster'
import { TEKS_MASTER } from './labels'

export default function HalamanMaster({ prefix }: { prefix: string }) {
  const klien = useMemo(() => klienMaster(prefix), [prefix])
  const [meta, setMeta] = useState<MetaMaster | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let batal = false
    setMeta(null)
    setGalat(null)
    klien.meta().then(
      (m) => {
        if (!batal) setMeta(m)
      },
      (err: unknown) => {
        if (!batal) setGalat(err)
      },
    )
    return () => {
      batal = true
    }
  }, [klien])

  return (
    <section className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{meta?.judul ?? ''}</h2>
      </header>
      <p className="muted">{TEKS_MASTER.sub}</p>
      {meta === null && galat === null && <Memuat pesan={TEKS_MASTER.memuat} />}
      <Gagal galat={galat} />
      {meta && <DaftarMaster klien={klien} meta={meta} />}
    </section>
  )
}
