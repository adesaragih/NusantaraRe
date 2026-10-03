// Isian angka di sel grid (spreading, angsuran) - berformat selama tidak
// difokus bila sel Pega ber-`pyShowReadonlyFormatting=true` (K14, `sajian.ts`),
// mentah saat diketik. Nilai yang dikirim ke backend selalu mentah.

import { useState } from 'react'

import { sajikan, type Sajian } from '../sajian'

export default function InputAngka({
  value,
  sajian,
  label,
  onChange,
  onBlur,
}: {
  value: string
  sajian: Sajian
  label: string
  onChange: (v: string) => void
  onBlur: () => void
}) {
  const [fokus, setFokus] = useState(false)
  const berformat = !fokus && sajian !== 'tanggal' && sajian.formatSaatSunting === true
  return (
    <input
      className="field__input"
      aria-label={label}
      value={berformat ? sajikan(value, sajian) : value}
      onFocus={() => setFokus(true)}
      onChange={(e) => onChange(e.target.value)}
      onBlur={() => {
        setFokus(false)
        onBlur()
      }}
    />
  )
}
