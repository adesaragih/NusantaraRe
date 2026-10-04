// Isian teks dengan saran sambil mengetik - padanan `pxAutoComplete` Pega bersumber Report Definition (popup Choose
// Accumulation, tiket 46). Mengetik mengubah teks dan memanggil `cari` sesudah jeda; memilih saran memanggil `onPilih`
// (nilai yang disalin Pega saat memilih: ID, NationInitial, ZipCode, …). Struktur = `div.field > label, input, ul`
// supaya tata letak berlabel kiri (`nbf-labelkiri`) menaruh daftar saran di bawah kotak.

import { useEffect, useRef, useState } from 'react'

import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import type { SaranAkumulasi } from '../api'

/** Jeda sebelum mencari (sama dengan popup lain modul ini). */
export const JEDA_SARAN_MS = 400

export default function IsianSaran({
  label,
  value,
  onKetik,
  cari,
  onPilih,
  tampil = (s) => s.label,
}: {
  label: string
  value: string
  /** Teks diketik (bukan dipilih) - pemanggil mengosongkan nilai hasil pilihan sebelumnya. */
  onKetik: (v: string) => void
  cari: (q: string) => Promise<{ baris: SaranAkumulasi[] }>
  onPilih: (s: SaranAkumulasi) => void
  /** Teks satu baris saran. */
  tampil?: (s: SaranAkumulasi) => string
}) {
  const [saran, setSaran] = useState<SaranAkumulasi[]>([])
  const [ketik, setKetik] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const nomor = useRef(0)

  useEffect(() => {
    if (!ketik) {
      setSaran([])
      return
    }
    const n = ++nomor.current
    const jadwal = window.setTimeout(() => {
      cari(value).then(
        (h) => {
          if (n === nomor.current) setSaran(h.baris)
        },
        (err: unknown) => {
          if (n === nomor.current) setGalat(err)
        },
      )
    }, JEDA_SARAN_MS)
    return () => window.clearTimeout(jadwal)
    // `cari` sengaja tidak menjadi dependensi: pemanggil membuat fungsi baru tiap render.
  }, [value, ketik])

  return (
    <div className="field">
      <label className="field__label">{label}</label>
      <input
        className="field__input"
        type="text"
        value={value}
        onChange={(e) => {
          setKetik(true)
          setGalat(null)
          onKetik(e.target.value)
        }}
      />
      {saran.length > 0 ? (
        <ul className="nbf-tambah-risk__saran" role="listbox" aria-label={label}>
          {saran.map((s, i) => (
            <li key={`${i}-${s.id}`}>
              <button
                type="button"
                className="btn btn--ghost btn--sm"
                onClick={() => {
                  setKetik(false)
                  setSaran([])
                  onPilih(s)
                }}
              >
                {tampil(s)}
              </button>
            </li>
          ))}
        </ul>
      ) : (
        galat !== null && <Gagal galat={galat} />
      )}
    </div>
  )
}
