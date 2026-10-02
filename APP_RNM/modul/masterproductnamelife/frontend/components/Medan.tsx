// Satu baris form label-kiri / nilai-kanan - tata letak `InboxProductName` (foto layar Pega work owner 02-10-2026:
// label di kiri, nilai di kanan, angka rata kanan, nilai kosong bertanda).
//
// Mode lihat (`ProductName.IsView == 'true'`): TEKS, bukan isian baca-saja - angka berpemisah ribuan (`tampilAngka`),
// tanggal DD/MM/YYYY (`tampilTanggal`), kosong `—`. Mode sunting: isian bawaan (teks / angka rata kanan / tanggal),
// atau `children` (dropdown master, pilihan, area) bila diberikan. Tanpa `onUbah` dan tanpa `children` = medan yang
// SELALU baca-saja (`Product Code`).

import { useId, type ReactNode } from 'react'

import { useTeksUI } from '../../../../inti/frontend/components/ui/bahasaUI'
import { tampilAngka, tampilTanggal } from '../bentuk'

export type JenisMedan = 'teks' | 'angka' | 'tanggal'

/** Tanda nilai kosong di mode lihat (Pega: `——`). */
export const TANDA_KOSONG = '—'

export default function Medan({
  label,
  nilai,
  lihat,
  onUbah,
  jenis = 'teks',
  required,
  kode,
  tampil,
  catatan,
  children,
}: {
  /** Label medan VERBATIM. */
  label: string
  nilai: string
  lihat: boolean
  onUbah?: (v: string) => void
  /** Bentuk tampilan dan isian - menurut tipe kolom flat (`mpnl_flat.go`, dijaga `FormProduk.test.ts`). */
  jenis?: JenisMedan
  /** `pyRequired` true di XML: tanda `*`. */
  required?: boolean
  /** Identitas (`Product Code`): tampil berbingkai seperti kotak baca-saja Pega. */
  kode?: boolean
  /** Teks mode lihat bila berbeda dari nilai (label pilihan pembayaran). */
  tampil?: string
  /** Keterangan di kanan nilai (`OF SUM REASURED` b24880). */
  catatan?: ReactNode
  /** Kontrol mode sunting pengganti isian bawaan. */
  children?: ReactNode
}) {
  const id = useId()
  const judul = (
    <>
      {label}
      {required && (
        <span className="mpnl-wajib" aria-hidden="true">
          *
        </span>
      )}
    </>
  )
  const teks = tampil ?? (jenis === 'angka' ? tampilAngka(nilai) : jenis === 'tanggal' ? tampilTanggal(nilai) : nilai)
  const kosong = teks.trim() === ''
  const kelasNilai = ['mpnl-nilai', jenis === 'angka' ? 'mpnl-nilai--angka' : '', kode ? 'mpnl-kode' : '', kosong ? 'mpnl-nilai--kosong' : '']
    .filter((k) => k !== '')
    .join(' ')
  const nilaiTeks = <span className={kelasNilai}>{kosong ? TANDA_KOSONG : teks}</span>

  if (lihat) {
    return (
      <div className="mpnl-medan mpnl-medan--lihat">
        <span className="mpnl-medan__label">{judul}</span>
        <div className="mpnl-medan__isi">
          {nilaiTeks}
          {catatan}
        </div>
      </div>
    )
  }

  const isianBawaan = children === undefined && onUbah !== undefined
  // `pxDateTime` - isian tanggal untuk `YYYY-MM-DD`; teks lama yang bukan tanggal tetap isian teks (tidak dibuang).
  const tipe = jenis === 'tanggal' && (nilai === '' || /^\d{4}-\d{2}-\d{2}$/.test(nilai)) ? 'date' : 'text'
  return (
    <div className="mpnl-medan">
      <label className="mpnl-medan__label" htmlFor={isianBawaan ? id : undefined}>
        {judul}
      </label>
      <div className="mpnl-medan__isi">
        {children ??
          (onUbah === undefined ? (
            nilaiTeks
          ) : (
            <input
              id={id}
              className={jenis === 'angka' ? 'field__input mpnl-isian--angka' : 'field__input'}
              type={tipe}
              inputMode={jenis === 'angka' ? 'decimal' : undefined}
              aria-required={required}
              value={nilai}
              onChange={(e) => {
                onUbah(e.target.value)
              }}
            />
          ))}
        {catatan}
      </div>
    </div>
  )
}

/** Pilihan statis di dalam `Medan` (`Premium Payment Method` b25611) - nilai di luar daftar tetap tampil, tidak dibuang. */
export function PilihanMedan({
  labelAria,
  value,
  opsi,
  onChange,
}: {
  labelAria: string
  value: string
  opsi: readonly { value: string; label: string }[]
  onChange: (v: string) => void
}) {
  const teksUI = useTeksUI()
  const asing = value !== '' && !opsi.some((o) => o.value === value)
  return (
    <select
      className="field__input"
      aria-label={labelAria}
      value={value}
      onChange={(e) => {
        onChange(e.target.value)
      }}
    >
      <option value="">{teksUI.pilihKosong}</option>
      {asing && (
        <option value={value}>
          {value} {teksUI.tidakDiDaftar}
        </option>
      )}
      {opsi.map((o) => (
        <option key={o.value} value={o.value}>
          {o.label}
        </option>
      ))}
    </select>
  )
}
