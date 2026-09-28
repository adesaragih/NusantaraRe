// Pemilih jenis reasuransi — tiket 02 Treaty Contract Out.
//
// SATU komponen untuk layar kontrak (`InputTreatyContractReinsType.xml`
// b2652 `ReinsType`) maupun seluruh grid klausul (RD dominan dipakai 11
// grid). Daftarnya datang dari `GET /api/treaty-contract-out/jenis-reasuransi`
// — SUDAH tersaring di server (12 awalan blacklist, Flag active, Type 1/2/3).
//
// ⛔ Layar TIDAK menyaring ulang dan TIDAK menulis daftar: pilihan yang ditulis
// di React pasti karangan (`Pilih`, ui/dasar.tsx).
//
// ⛔ Master kosong / tidak terbaca DINYATAKAN (503 dari server, ADR-0015) —
// bukan dropdown kosong yang terbaca "belum ada jenis".

import { useEffect, useState } from 'react'

import { JENIS_REASURANSI_TCO } from '../../assets/labels.treaty-contract-out'
import { ambilJenisReasuransiTreaty, type JenisReasuransiTreaty } from '../../services/api'
import { Gagal, Memuat, Pilih, type Opsi } from '../ui/dasar'

/** Pilihan `Pilih` dari daftar server: nilai = `.ID` (TEKS), label = `.Note`. */
export function opsiJenisReasuransi(daftar: readonly JenisReasuransiTreaty[]): Opsi[] {
  return daftar.map((j) => ({ value: j.id, label: j.note }))
}

export default function PilihJenisReasuransi({
  value,
  onChange,
  label = JENIS_REASURANSI_TCO.reinsType,
  required,
  error,
}: {
  value: string
  onChange: (v: string) => void
  /** `ReinsType` (form kontrak) atau `Reinsurance Type` (form tahun). */
  label?: string
  required?: boolean
  error?: string
}) {
  const [daftar, setDaftar] = useState<JenisReasuransiTreaty[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)

  useEffect(() => {
    let hidup = true
    void (async () => {
      try {
        const hasil = await ambilJenisReasuransiTreaty()
        if (hidup) setDaftar(hasil.daftar)
      } catch (e) {
        if (hidup) setGalat(e)
      }
    })()
    return () => {
      hidup = false
    }
  }, [])

  if (galat !== null) return <Gagal galat={galat} />
  if (daftar === null) return <Memuat />
  return (
    <Pilih
      label={label}
      value={value}
      onChange={onChange}
      opsi={opsiJenisReasuransi(daftar)}
      required={required}
      error={error}
      kosong={JENIS_REASURANSI_TCO.belumDipilih}
    />
  )
}
