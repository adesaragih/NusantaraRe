// Isian kolom rujukan (mis. Nation pada Province): mengetik mencari master rujukan yang AKTIF (`GET
// /api/masterdata/{master}?q=&status=aktif`) dan memilih saran mengisi nilainya (ID, atau Code untuk Group Of CZone).
// Struktur `div.field > label, input, ul` sama dengan medan dasar.

import { useEffect, useRef, useState } from 'react'

import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import { cariMaster, type BarisMaster } from '../api'
import { KOLOM_NAMA, TEKS } from '../labels'

/** Jeda sebelum mencari. */
export const JEDA_CARI_MS = 400

export default function PilihRujukan({
  label,
  value,
  onChange,
  master,
  nilai,
  required,
  error,
}: {
  label: string
  value: string
  onChange: (v: string) => void
  /** Kunci master rujukan. */
  master: string
  /** Kunci kolom baris rujukan yang disimpan. */
  nilai: string
  required?: boolean
  error?: string
}) {
  const [saran, setSaran] = useState<BarisMaster[]>([])
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
      cariMaster(master, value, 'aktif', 1).then(
        (h) => {
          if (n === nomor.current) setSaran(h.baris)
        },
        (err: unknown) => {
          if (n === nomor.current) setGalat(err)
        },
      )
    }, JEDA_CARI_MS)
    return () => window.clearTimeout(jadwal)
  }, [value, ketik, master])

  const nama = KOLOM_NAMA[master] ?? 'note'
  return (
    <div className="field">
      <label className="field__label">
        {label}
        {required && <span className="field__req">*</span>}
      </label>
      <input
        className={'field__input' + (error ? ' field__input--error' : '')}
        type="text"
        value={value}
        placeholder={TEKS.cariRujukan}
        onChange={(e) => {
          setKetik(true)
          setGalat(null)
          onChange(e.target.value)
        }}
      />
      {error && <div className="field__error">{error}</div>}
      {saran.length > 0 && (
        <ul className="md-rujukan" role="listbox" aria-label={label}>
          {saran.map((b) => {
            const v = String(b[nilai] ?? '')
            return (
              <li key={v}>
                <button
                  type="button"
                  className="btn btn--ghost btn--sm"
                  onClick={() => {
                    setKetik(false)
                    setSaran([])
                    onChange(v)
                  }}
                >
                  {[v, String(b[nama] ?? '')].filter(Boolean).join(' · ')}
                </button>
              </li>
            )
          })}
        </ul>
      )}
      {galat !== null && <Gagal galat={galat} />}
    </div>
  )
}
