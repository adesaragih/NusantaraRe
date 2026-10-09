import { useId } from 'react'

import { KotakTanggalKetik } from './TanggalKetik'

/**
 * Medan tanggal BERLABEL — dapat DIKETIK (`DD/MM/YYYY`, `18-01-2024`,
 * `18012024`, …) atau dipilih dari ikon kalender. Lihat `TanggalKetik.tsx`.
 *
 * ⭐ 7 Oktober 2026 — permintaan pemakai: *"tiap inputan tanggal bisa
 * diketik juga biar gampang"*. Sebelumnya pembungkus `FieldTanggal` inti
 * (`<input type="date">`) yang hanya dapat diisi per segmen, dan yang
 * menulis `dd/mm/yyyy` sendiri saat kosong — masalah yang dulu ditambal
 * dengan kelas `trin__tanggal--kosong`. Kotak teks yang kosong memang
 * kosong, jadi tambalan itu tidak diperlukan lagi.
 *
 * API tetap: `value` bentuk apa pun yang `keInputTanggal` inti baca,
 * `onChange` menerima bentuk KABEL `DD-MM-YYYY` — sama dengan `FieldTanggal`.
 */
export default function TanggalRedup({
  label,
  value,
  onChange,
}: {
  label: string
  value: string
  onChange: (v: string) => void
}) {
  const id = useId()
  return (
    <div className="field">
      <label className="field__label" htmlFor={id}>
        {label}
      </label>
      <KotakTanggalKetik id={id} label={label} value={value} onChange={onChange} />
    </div>
  )
}
