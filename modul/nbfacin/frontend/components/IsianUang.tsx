// Kotak isian uang berformat Indonesia saat diketik - permintaan work owner 03-10-2026 ("TSI Object Item (All Unit)
// biar angka nya ada pemisah nya saat diketik. misal jadi 20.000.000").
//
// Yang TAMPIL dan diketik: titik ribuan, koma desimal (gaya `formatNumber` inti). Nilai KELUAR = teks desimal bertitik
// tanpa pemisah ribuan ("20000000.5") - bentuk kabel API (ADR-0003/0034). Bekerja pada DIGIT (teks), tanpa float.
// Titik yang diketik pengguna dianggap pemisah ribuan (dibuang); koma pertama = pemisah desimal; paling banyak 8 desimal
// (NUMBER(38,8), ADR-0016). Tanpa tanda minus (uang di layar ini >= 0).

import { useEffect, useState } from 'react'

/** Maksimum digit desimal (ADR-0016). */
const MAKS_DESIMAL = 8

/** Kelompokkan digit bulat dengan titik ribuan. */
function kelompok(bulat: string): string {
  return bulat.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
}

/** Ketikan bebas -> teks tampil (titik ribuan, koma desimal; koma di ujung dipertahankan selama mengetik). */
export function rapikanUang(s: string): string {
  const koma = s.indexOf(',')
  const kiri = (koma < 0 ? s : s.slice(0, koma)).replace(/\D/g, '')
  const kanan = koma < 0 ? null : s.slice(koma + 1).replace(/\D/g, '').slice(0, MAKS_DESIMAL)
  if (kiri === '' && kanan === null) return ''
  const bulat = kiri.replace(/^0+(?=\d)/, '') || '0'
  return kelompok(bulat) + (kanan === null ? '' : `,${kanan}`)
}

/** Teks tampil -> kabel ("20.000.000,5" -> "20000000.5"); kosong -> "". */
export function keKabelUang(tampil: string): string {
  if (tampil === '') return ''
  const [bulat = '', pecahan = ''] = tampil.split(',')
  const b = bulat.replace(/\./g, '')
  return pecahan === '' ? b : `${b}.${pecahan}`
}

/** Kabel -> teks tampil ("20000000.50" -> "20.000.000,50"); bukan desimal sah -> apa adanya. */
export function keTampilUang(kabel: string): string {
  const m = /^(\d+)(?:\.(\d*))?$/.exec(kabel.trim())
  if (!m) return kabel
  return kelompok(m[1]!.replace(/^0+(?=\d)/, '')) + (m[2] ? `,${m[2]}` : '')
}

export default function IsianUang({
  label,
  value,
  onChange,
  error,
  required,
}: {
  label: string
  /** Kabel: teks desimal bertitik, atau kosong. */
  value: string
  onChange: (kabel: string) => void
  error?: string
  required?: boolean
}) {
  const [teks, setTeks] = useState(keTampilUang(value))

  // Nilai dari luar berubah (mis. dimuat ulang dari server) -> tampilan mengikuti; ketikan sendiri tidak diganggu.
  useEffect(() => {
    if (keKabelUang(teks) !== value) setTeks(keTampilUang(value))
  }, [value])

  return (
    <div className="field">
      <label className="field__label">
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      <input
        className={'field__input nbf-angka' + (error ? ' field__input--error' : '')}
        type="text"
        inputMode="decimal"
        value={teks}
        onChange={(e) => {
          const rapi = rapikanUang(e.target.value)
          setTeks(rapi)
          onChange(keKabelUang(rapi))
        }}
      />
      {error && <div className="field__error">{error}</div>}
    </div>
  )
}
