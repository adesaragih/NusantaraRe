// Isian teks dengan saran sambil mengetik - padanan `pxAutoComplete` Pega bersumber Report Definition (popup Choose
// Accumulation, tiket 46). Mengetik mengubah teks dan memanggil `cari` sesudah jeda; memilih saran memanggil `onPilih`
// (nilai yang disalin Pega saat memilih: ID, NationInitial, ZipCode, …). Struktur = `div.field > label, input, ul`
// supaya tata letak berlabel kiri (`nbf-labelkiri`) menaruh daftar saran di bawah kotak.

import { useEffect, useRef, useState } from 'react'

import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import type { SaranAkumulasi } from '../api'

/** Jeda sebelum mencari (sama dengan popup lain modul ini). */
export const JEDA_SARAN_MS = 400

/** Saran yang labelnya sama persis dengan teks (tanpa beda huruf besar/kecil, spasi tepi diabaikan); null bila tidak ada. */
export function saranSama(teks: string, baris: SaranAkumulasi[]): SaranAkumulasi | null {
  const t = teks.trim().toUpperCase()
  return t === '' ? null : (baris.find((s) => s.label.trim().toUpperCase() === t) ?? null)
}

export default function IsianSaran({
  label,
  value,
  onKetik,
  cari,
  onPilih,
  tampil = (s) => s.label,
  muatSaatFokus = false,
  wajibDariDaftar = false,
}: {
  label: string
  value: string
  /** Teks diketik (bukan dipilih) - pemanggil mengosongkan nilai hasil pilihan sebelumnya. */
  onKetik: (v: string) => void
  cari: (q: string) => Promise<{ baris: SaranAkumulasi[] }>
  onPilih: (s: SaranAkumulasi) => void
  /** Teks satu baris saran. */
  tampil?: (s: SaranAkumulasi) => string
  /**
   * Daftar langsung dimuat saat kotak difokus, sebelum mengetik (dropdown pxAutoComplete Pega; work owner 06-10-2026:
   * "di dropdown kebawah juga harus muncul list nya, jgn saat mulai di ketik huruf saja"); ditutup saat kotak ditinggal.
   */
  muatSaatFokus?: boolean
  /**
   * Nilai wajib dari daftar (work owner 06-10-2026, Class Of Business Loss Record: "kalau yang diketik tidak ada dari list,
   * langsung dihps aja ketikannya"): saat kotak ditinggal sesudah diketik, teks yang sama persis dengan satu saran dipilih
   * (`onPilih`), selain itu dikosongkan (`onKetik('')`).
   */
  wajibDariDaftar?: boolean
}) {
  const [saran, setSaran] = useState<SaranAkumulasi[]>([])
  const [ketik, setKetik] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const nomor = useRef(0)
  // Diketik sejak pilihan terakhir (wajibDariDaftar).
  const diketik = useRef(false)

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
    <div className="field nbf-saran">
      <label className="field__label">{label}</label>
      <input
        className="field__input"
        type="text"
        value={value}
        onChange={(e) => {
          diketik.current = true
          setKetik(true)
          setGalat(null)
          onKetik(e.target.value)
        }}
        onFocus={() => {
          if (!muatSaatFokus) return
          setGalat(null)
          setKetik(true)
        }}
        onBlur={() => {
          // Jeda singkat supaya klik pada saran sempat diproses sebelum daftar ditutup.
          if (muatSaatFokus) window.setTimeout(() => setKetik(false), 200)
          if (!wajibDariDaftar) return
          const teks = value
          window.setTimeout(() => {
            if (!diketik.current) return
            if (teks.trim() === '') {
              diketik.current = false
              return
            }
            cari(teks).then(
              (h) => {
                if (!diketik.current) return
                diketik.current = false
                setKetik(false)
                setSaran([])
                const s = saranSama(teks, h.baris)
                if (s) onPilih(s)
                else onKetik('')
              },
              (err: unknown) => setGalat(err),
            )
          }, 200)
        }}
      />
      {saran.length > 0 ? (
        <ul className="nbf-saran__daftar" role="listbox" aria-label={label}>
          {saran.map((s, i) => (
            <li key={`${i}-${s.id}`}>
              <button
                type="button"
                className="btn btn--ghost btn--sm"
                onClick={() => {
                  diketik.current = false
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
