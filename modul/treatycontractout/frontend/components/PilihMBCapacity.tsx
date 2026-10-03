// Pemilih form MB Capacity — 10017 LimitMB, Treaty Contract Out [keputusan work owner 02-10-2026: "ikuti
// xml-nya aja"].
//
// `[terverifikasi]` `Section/GridTreatyArrangementLIMITMB.xml`:
// - Occupation: Dropdown `ID_Occupation` berteks kosong `Choose` (b1696); berubah → `SetOccupationLimitMB`
//   (b1669) yang mengisi nama untuk ID 01–04. Opsinya dari server (`GET /klausul-pilihan/occupation-limitmb`,
//   `models.PilihanOccupationLimitMB`); nama yang tersimpan diisi server.
// - TerritorialLimit: Autocomplete atas `BrowseTreatyGroup_RD` (b3485), tampil dan DISIMPAN `.TreatyGroupName`
//   (`pyPropertyTarget` associated property). Daftarnya master grup treaty (`GET /grup-treaty`, tiket 03),
//   disaring di sini; ID tampil sebagai keterangan.

import { useEffect, useState } from 'react'

import { Gagal, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import { PilihSaring, type OpsiSaring } from '../../../../inti/frontend/components/ui/pilihSaring'
import { ambilGrupTreaty, cariPilihanKlausul, type GrupTreaty, type PilihanKlausul } from '../api'
import { KLAUSUL_TCO } from '../labels'

/** Opsi dropdown Occupation: ID nilai, nama tampil. */
export function opsiOccupationMB(daftar: readonly PilihanKlausul[]): Opsi[] {
  return daftar.map((p) => ({ value: p.id, label: p.nama }))
}

/** Opsi autocomplete grup treaty untuk kata itu: nilai = nama (yang disimpan), nama kembar tampil sekali. */
export function saringGrupTreaty(daftar: readonly GrupTreaty[], kata: string): OpsiSaring[] {
  const k = kata.trim().toLowerCase()
  const lihat = new Set<string>()
  const hasil: OpsiSaring[] = []
  for (const g of daftar) {
    if (lihat.has(g.treatyGroupName)) continue
    if (k !== '' && !g.treatyGroupName.toLowerCase().includes(k) && !g.id.includes(k)) continue
    lihat.add(g.treatyGroupName)
    hasil.push({ value: g.treatyGroupName, label: g.treatyGroupName, keterangan: g.id })
  }
  return hasil
}

export function PilihOccupationMB({
  label,
  value,
  onChange,
}: {
  label: string
  value: string
  onChange: (id: string) => void
}) {
  const [daftar, setDaftar] = useState<PilihanKlausul[]>([])
  const [galat, setGalat] = useState<unknown>(null)
  useEffect(() => {
    cariPilihanKlausul('occupation-limitmb', '').then(setDaftar, setGalat)
  }, [])
  if (galat !== null) return <Gagal galat={galat} />
  return <Pilih label={label} value={value} onChange={onChange} opsi={opsiOccupationMB(daftar)} kosong={KLAUSUL_TCO.choose} />
}

export function PilihGrupTreatyNama({
  label,
  value,
  onChange,
}: {
  label: string
  /** `TreatyGroupName` terpilih (medan `TerritorialLimit`). */
  value: string
  onChange: (nama: string) => void
}) {
  const [daftar, setDaftar] = useState<GrupTreaty[] | null>(null)
  const [kata, setKata] = useState('')
  const [galat, setGalat] = useState<unknown>(null)
  useEffect(() => {
    ambilGrupTreaty().then((j) => setDaftar(j.daftar ?? []), setGalat)
  }, [])
  if (galat !== null) return <Gagal galat={galat} />
  return (
    <PilihSaring
      label={label}
      value={value}
      teksTerpilih={value}
      opsi={saringGrupTreaty(daftar ?? [], kata)}
      onCari={setKata}
      onPilih={(o) => onChange(o.value)}
      memuat={daftar === null}
    />
  )
}
