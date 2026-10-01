// Autocomplete `pxAutoComplete` - medan `Ceding` b4107, `SOB` b4495, `R/I Risk Name` b7430, `Policy Holder`
// b17129, `Currency` b28173, dan `Plan Name` b33124 (grid `PLAN LIST`).
//
// Pega: mengetik menyaring RD sumbernya (`Contains`), memilih satu baris menyalin ID + nama (`set*_DT`).
// Di sini: mengetik mengubah nama dan MENGOSONGKAN pasangan ID-nya (pemanggil) - server menolak nama tanpa ID
// dengan kalimat "must be chosen from the master list"; memilih saran mengisi keduanya.

import { useEffect, useRef, useState } from 'react'

import { Gagal } from '../../../../inti/frontend/components/ui/dasar'

export interface PropsSaran<T> {
  /** Label medan VERBATIM; kosong untuk sel grid (label = kepala kolom). */
  label?: string
  /** Label untuk pembaca layar bila `label` kosong. */
  labelAria?: string
  nilai: string
  onKetik: (teks: string) => void
  cari: (kata: string) => Promise<readonly T[]>
  teks: (t: T) => string
  kunci: (t: T) => string
  onPilih: (t: T) => void
  readOnly?: boolean
}

/** Jeda ketik sebelum RD dibaca. */
const JEDA_MS = 250
/** Saran yang ditampilkan paling banyak. */
const SARAN_MAKS = 20

export default function Saran<T>({ label, labelAria, nilai, onKetik, cari, teks, kunci, onPilih, readOnly }: PropsSaran<T>) {
  const [buka, setBuka] = useState(false)
  const [daftar, setDaftar] = useState<readonly T[]>([])
  const [galat, setGalat] = useState<unknown>(null)
  const cariRef = useRef(cari)
  cariRef.current = cari

  useEffect(() => {
    if (!buka) return
    let batal = false
    const jam = setTimeout(() => {
      cariRef
        .current(nilai)
        .then((d) => {
          if (!batal) {
            setDaftar(d.slice(0, SARAN_MAKS))
            setGalat(null)
          }
        })
        .catch((e: unknown) => {
          if (!batal) setGalat(e)
        })
    }, JEDA_MS)
    return () => {
      batal = true
      clearTimeout(jam)
    }
  }, [buka, nilai])

  const input = (
    <input
      className="field__input"
      value={nilai}
      readOnly={readOnly}
      aria-label={label === undefined || label === '' ? labelAria : undefined}
      onChange={(e) => {
        onKetik(e.target.value)
        setBuka(true)
      }}
      onFocus={() => {
        if (!readOnly) setBuka(true)
      }}
      onBlur={() => {
        // Klik saran (mousedown) terjadi sebelum blur menutup daftar.
        setTimeout(() => {
          setBuka(false)
        }, 150)
      }}
    />
  )

  return (
    <div className={label ? 'field mpnl-saran' : 'mpnl-saran'}>
      {label ? <label className="field__label">{label}</label> : null}
      {input}
      {buka && daftar.length > 0 && (
        <ul className="mpnl-saran__daftar" role="listbox">
          {daftar.map((t) => (
            <li
              key={kunci(t)}
              className="mpnl-saran__butir"
              role="option"
              aria-selected={false}
              onMouseDown={(e) => {
                e.preventDefault()
                onPilih(t)
                setBuka(false)
              }}
            >
              {teks(t)}
            </li>
          ))}
        </ul>
      )}
      {buka && galat !== null && <Gagal galat={galat} />}
    </div>
  )
}
