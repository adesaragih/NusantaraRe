// Dropdown ID Occupation / ID Clause — 10013 Exclusion Treaty, dapat difilter.
//
// `[terverifikasi]` XML: `ID_Occupation` input teks (`GridTreatyArrangementExclutionTreatyOccupation.xml`
// b2038) yang diisi Autocomplete nama `Occupation` (RD `BrowseOccupationFIRE_RD`, tampil `.Name`, target
// `.ID_Occupation` b2470); Clause sama (`…ExclutionTreatyClausule.xml` tampil `.Info`, target b2501).
//
// [keputusan work owner 02-10-2026] Kotak "Search" terpisah dibuang; dropdown terisi begitu dibuka
// (100 baris pertama berurut nama, `GET /klausul-pilihan/{master}?cari=`), ketikan di dalamnya
// menyaring di SERVER sehingga seluruh master terjangkau. Memilih ID ikut mengisi nama untuk tampilan;
// nama yang tersimpan tetap diambil server dari master.

import { useEffect, useRef, useState } from 'react'

import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import { PilihSaring, type OpsiSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import { cariPilihanKlausul, type PilihanKlausul } from '../api'

/** Opsi dropdown: nama tampil, ID sebagai keterangan. */
export function opsiMasterKlausul(daftar: readonly PilihanKlausul[]): OpsiSaring[] {
  return daftar.map((p) => ({ value: p.id, label: p.nama, keterangan: p.id }))
}

export default function PilihMasterKlausul({
  master,
  label,
  value,
  nama,
  onPilih,
  required,
}: {
  master: 'occupation' | 'clause'
  label: string
  /** ID terpilih. */
  value: string
  /** Nama terpilih (medan `Occupation` / `Clause` form). */
  nama: string
  onPilih: (id: string, nama: string) => void
  required?: boolean
}) {
  const [kata, setKata] = useState('')
  const [daftar, setDaftar] = useState<PilihanKlausul[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  // Hanya jawaban pencarian TERAKHIR yang dipakai (ketikan cepat).
  const urutan = useRef(0)

  useEffect(() => {
    const ke = ++urutan.current
    setDaftar(null)
    setGalat(null)
    cariPilihanKlausul(master, kata.trim()).then(
      (d) => {
        if (ke === urutan.current) setDaftar(d)
      },
      (e: unknown) => {
        if (ke === urutan.current) setGalat(e)
      },
    )
  }, [master, kata])

  if (galat !== null) return <Gagal galat={galat} />
  return (
    <PilihSaring
      label={label}
      value={value}
      teksTerpilih={nama !== '' ? nama : value}
      opsi={opsiMasterKlausul(daftar ?? [])}
      memuat={daftar === null}
      onCari={setKata}
      onPilih={(o) => {
        onPilih(o.value, o.label)
      }}
      required={required}
    />
  )
}
