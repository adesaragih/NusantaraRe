// Satu daftar cari di menu View Policy (tiket 71; tab Facultative / Treaty / Life, work owner 11-10-2026) - pola section Pega
// `DATA-PORTAL!VIEWOPENPOLICYMARINEINBOX` (`D:\migrasi\RNM\DDL\ViewOpenPolicyMarineInbox.xml`): View Option (bila daftar
// punya) + Filter + Refresh, grid berhalaman, tautan View Details. Pencarian berjalan saat mengetik / ganti View Option /
// Enter; filter kosong = daftar kosong (`GetListViewPolicy_act` langkah 7). Kolom dan isi sel disuplai pemanggil.

import { useEffect, useRef, useState, type ReactNode } from 'react'

import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { FOLLOWING as F, TEKS_FOLLOWING, VIEW_POLICY as V } from '../labels'
import { JEDA_CARI_MS } from './PopupPilihAgent'

/** Satu halaman hasil cari (bentuk kabel semua daftar View Policy). */
export interface HalamanCari<T> {
  baris: T[]
  total: number
  halaman: number
  ukuran: number
}

export default function DaftarCariPolis<T>({
  judul,
  opsi,
  placeholder,
  kolom,
  cari,
  kunci,
  sel,
  onLihat,
  bisaLihat = () => true,
}: {
  /** Judul daftar ("List Polis" / "List Offer"). */
  judul: string
  /** Pilihan View Option; tanpa ini daftar mencari dari isi teks saja (Life: tanpa View Option, work owner 11-10-2026). */
  opsi?: readonly { value: string; label: string }[]
  placeholder: string
  kolom: readonly string[]
  cari: (opsi: string, filter: string, halaman: number) => Promise<HalamanCari<T>>
  kunci: (b: T, i: number) => string
  sel: (b: T) => ReactNode[]
  /** View Details; tanpa ini kolom tautan tidak ada (detail daftar itu belum diport). */
  onLihat?: (b: T) => void
  bisaLihat?: (b: T) => boolean
}) {
  const perluOpsi = opsi !== undefined
  const [pilihan, setPilihan] = useState('')
  const [filter, setFilter] = useState('')
  const [hasil, setHasil] = useState<HalamanCari<T> | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [memuat, setMemuat] = useState(false)
  const nomorPermintaan = useRef(0)

  async function muat(o: string, f: string, halaman: number) {
    const nomor = ++nomorPermintaan.current
    // GetListViewPolicy_act langkah 7: filter kosong -> keluar tanpa mencari.
    if ((perluOpsi && o === '') || f.trim() === '') {
      setHasil(null)
      setMemuat(false)
      return
    }
    setMemuat(true)
    setGalat(null)
    try {
      const h = await cari(o, f, halaman)
      if (nomor === nomorPermintaan.current) setHasil(h)
    } catch (err) {
      if (nomor === nomorPermintaan.current) {
        setHasil(null)
        setGalat(err)
      }
    } finally {
      if (nomor === nomorPermintaan.current) setMemuat(false)
    }
  }

  useEffect(() => {
    const jadwal = window.setTimeout(() => void muat(pilihan, filter, 1), JEDA_CARI_MS)
    return () => window.clearTimeout(jadwal)
  }, [pilihan, filter])

  return (
    <div className="nbf-vp__daftar">
      <h3 className="nbf-vp__subjudul">{judul}</h3>
      <form
        className="toolbar"
        onSubmit={(e) => {
          e.preventDefault()
          void muat(pilihan, filter, 1)
        }}
      >
        {perluOpsi && (
          <select
            className="field__input nbf-tema__isian nbf-vp__opsi"
            aria-label={`${judul} - ${F.viewOption}`}
            value={pilihan}
            onChange={(e) => setPilihan(e.target.value)}
          >
            {/* Pilihan kosong = "Choose" saja seperti Pega (work owner 11-10-2026); nama "View Option" tetap di aria-label. */}
            <option value="">{F.pilihKosong}</option>
            {opsi.map((o) => (
              <option key={o.value} value={o.value}>
                {o.label}
              </option>
            ))}
          </select>
        )}
        <input
          className="field__input nbf-tema__isian nbf-saring__kotak"
          type="search"
          aria-label={`${judul} - ${F.filter}`}
          placeholder={placeholder}
          value={filter}
          disabled={perluOpsi && pilihan === ''}
          onChange={(e) => setFilter(e.target.value)}
        />
        <button type="button" className="btn btn--sm nbf-tema__tombol" onClick={() => setFilter('')}>
          {F.refresh}
        </button>
      </form>
      {hasil && hasil.baris.length > 0 && (
        <div className="nbf-tema__pager">
          <Halaman halaman={hasil.halaman} ukuran={hasil.ukuran} total={hasil.total} onPindah={(h) => void muat(pilihan, filter, h)} />
        </div>
      )}
      <table className="inbox__tabel nbf-tema__tabel">
        <thead>
          <tr>
            {kolom.map((k) => (
              <th key={k} scope="col">
                {k}
              </th>
            ))}
            {onLihat && <th scope="col" />}
          </tr>
        </thead>
        {hasil && hasil.baris.length > 0 && (
          <tbody>
            {hasil.baris.map((b, i) => {
              const lihat = onLihat !== undefined && bisaLihat(b)
              return (
                <tr key={kunci(b, i)} className={lihat ? 'inbox__baris' : undefined} onClick={lihat ? () => onLihat(b) : undefined}>
                  {sel(b).map((isi, j) => (
                    <td key={j}>{isi}</td>
                  ))}
                  {onLihat && (
                    <td>
                      <button
                        type="button"
                        className="btn btn--sm nbf-tema__tombol"
                        disabled={!lihat}
                        onClick={(e) => {
                          // Baris sudah membuka detail; cegah panggilan ganda.
                          e.stopPropagation()
                          onLihat(b)
                        }}
                      >
                        {F.lihat}
                      </button>
                    </td>
                  )}
                </tr>
              )
            })}
          </tbody>
        )}
      </table>
      {memuat && <Memuat />}
      {!memuat && hasil === null && <Kosong pesan={perluOpsi ? TEKS_FOLLOWING.isiFilter : V.teks.isiFilterSaja} />}
      {!memuat && hasil !== null && hasil.baris.length === 0 && <Kosong pesan={TEKS_FOLLOWING.kosong} />}
      <Gagal galat={galat} />
    </div>
  )
}
