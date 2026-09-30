// Pemilih jenis reasuransi yang DAPAT DIFILTER — ReinsType Treaty Limit.
//
// `[terverifikasi]` Kedua grid Treaty Limit memakai `pxAutoComplete`, bukan
// dropdown biasa:
//   - induk `GridTreatyArrangementTreatyLimit.xml` b3025: RD
//     `BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull` (daftar induk tiket 02),
//     dicari dan ditampilkan `.Note` (b3119-b3122), `.ID` tersembunyi (b3133)
//   - anak `GridTreatyArrTreatyLimitList.xml` b2892: `ReinsTypeList` dari
//     `TreatyContractSetReinsTypeList` (tak diekspor) — dicari `.CARI2` (nama,
//     b3004), `.CARI1` (ID) ikut TAMPIL (b3017). Isinya porsi + induknya
//     [keputusan work owner 30-09-2026], disaring server.
//
// ⛔ Saringan di sini HANYA teks ketikan atas nama — aturan daftar (porsi,
// blacklist, Flag) tetap milik server.

import { useEffect, useState } from 'react'

import {
  ambilJenisReasuransiAnakTreatyLimit,
  ambilJenisReasuransiTreaty,
  type JenisReasuransiTreaty,
} from '../api'
import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import { PilihSaring, type OpsiSaring } from '../../../../inti/frontend/components/ui/pilihSaring'

/**
 * Opsi untuk kata `kata`: dicocokkan pada `.Note` saja (pyUseForSearch), tanpa
 * peduli huruf besar-kecil; `denganID` = ID tampil sebagai keterangan (anak).
 */
export function opsiSaringJenisReasuransi(
  daftar: readonly JenisReasuransiTreaty[],
  kata: string,
  denganID: boolean,
): OpsiSaring[] {
  const cari = kata.trim().toLowerCase()
  return daftar
    .filter((j) => cari === '' || j.note.toLowerCase().includes(cari))
    .map((j) => (denganID ? { value: j.id, label: j.note, keterangan: j.id } : { value: j.id, label: j.note }))
}

export default function PilihJenisReasuransiSaring({
  label,
  value,
  onChange,
  required,
  anakTreatyLimitDari,
}: {
  label: string
  value: string
  onChange: (v: string) => void
  required?: boolean
  /** Diisi = pilihan anak Treaty Limit di bawah induk ini; kosong = daftar induk. */
  anakTreatyLimitDari?: string
}) {
  const [daftar, setDaftar] = useState<JenisReasuransiTreaty[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [kata, setKata] = useState('')

  useEffect(() => {
    let hidup = true
    setDaftar(null)
    setGalat(null)
    const baca =
      anakTreatyLimitDari === undefined
        ? ambilJenisReasuransiTreaty()
        : ambilJenisReasuransiAnakTreatyLimit(anakTreatyLimitDari)
    baca.then(
      (h) => {
        if (hidup) setDaftar(h.daftar)
      },
      (e: unknown) => {
        if (hidup) setGalat(e)
      },
    )
    return () => {
      hidup = false
    }
  }, [anakTreatyLimitDari])

  if (galat !== null) return <Gagal galat={galat} />
  const terpilih = daftar?.find((j) => j.id === value)
  return (
    <PilihSaring
      label={label}
      value={value}
      teksTerpilih={terpilih?.note ?? value}
      opsi={daftar === null ? [] : opsiSaringJenisReasuransi(daftar, kata, anakTreatyLimitDari !== undefined)}
      memuat={daftar === null}
      onCari={setKata}
      onPilih={(o) => {
        onChange(o.value)
      }}
      required={required}
    />
  )
}
