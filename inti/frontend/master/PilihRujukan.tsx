// Isian kolom rujukan master (mis. Nation pada Province): dropdown yang dapat dicari (`PilihSaring` inti) atas
// `GET P/rujukan/{kunci}` - baris master yang dirujuk, aktif saja. Nilai yang disimpan = kolom `nilai` baris rujukan,
// saran menampilkan `nilai · nama`.

import { useRef, useState } from 'react'

import { PilihSaring, type OpsiSaring } from '../components/ui/pilihSaring'
import type { KlienMaster, RujukanMaster } from './api'

/** Opsi saran dari satu baris rujukan. */
export function opsiRujukan(r: RujukanMaster, b: Record<string, unknown>): OpsiSaring {
  const nilai = String(b[r.nilai] ?? '')
  const nama = String(b[r.nama] ?? '')
  return { value: nilai, label: nama || nilai, keterangan: nama ? nilai : undefined }
}

export default function PilihRujukan({
  klien,
  rujukan,
  label,
  value,
  onChange,
  required,
  error,
}: {
  klien: KlienMaster
  rujukan: RujukanMaster
  label: string
  value: string
  onChange: (v: string) => void
  required?: boolean
  error?: string
}) {
  const [opsi, setOpsi] = useState<OpsiSaring[]>([])
  const [teks, setTeks] = useState(value)
  const [memuat, setMemuat] = useState(false)
  // Jawaban lama dibuang: hanya permintaan terakhir yang mengisi opsi.
  const nomor = useRef(0)

  function cari(kata: string) {
    const n = ++nomor.current
    setMemuat(true)
    klien.rujukan(rujukan.kunci, kata, 1).then(
      (h) => {
        if (n !== nomor.current) return
        setOpsi(h.baris.map((b) => opsiRujukan(rujukan, b)))
        setMemuat(false)
      },
      () => {
        if (n === nomor.current) setMemuat(false)
      },
    )
  }

  return (
    <div>
      <PilihSaring
        label={label}
        value={value}
        teksTerpilih={teks}
        opsi={opsi}
        onCari={cari}
        onPilih={(o) => {
          setTeks(o.label)
          onChange(o.value)
        }}
        required={required}
        memuat={memuat}
      />
      {error && <div className="field__error">{error}</div>}
    </div>
  )
}
