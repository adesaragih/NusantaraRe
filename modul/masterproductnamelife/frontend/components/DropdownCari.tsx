// Dropdown bercari - mesin bersama `DropdownMaster` (ketujuh pemilih master, keputusan work owner 02-10-2026) dan
// `Plan Name` grid PLAN LIST (keputusan work owner 03-10-2026 "tolong ubah jadi model dropdown").
//
// Nilai HANYA dari daftar sumber (tidak dapat diketik). Medan sendiri membuka daftar; `Search` di dalam panel menyaring
// lewat sumber (`cari`, server); daftar paling banyak BATAS_DROPDOWN baris dengan potongan dinyatakan. Kolom daftar
// dan isinya ditentukan pemakai. Keyboard: Enter/Spasi/panah bawah membuka, panah/PageUp/PageDown menggeser, Enter
// memilih, Escape/klik di luar menutup tanpa memilih. Mode lihat: teks nilai, tidak dapat dibuka.

import { useEffect, useId, useRef, useState, type KeyboardEvent } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { geserAktif } from '../bentuk'
import { LAIN_MPNL, PEMILIH_MPNL } from '../labels'

/** Jeda ketik `Search` sebelum sumber dibaca. */
const JEDA_MS = 250
/** Lebar panel paling besar (px) - untuk memilih rata kiri atau kanan saat dibuka. */
const LEBAR_PANEL = 440
const LEBAR_PANEL_LEBAR = 760
/** Langkah `PageUp` / `PageDown`. */
const LANGKAH_HALAMAN = 10

/** Jawaban sumber: baris yang dirender dan apakah sumber memuat lebih (`potongPilihan`). */
export interface HasilCari<T> {
  tampil: T[]
  lebih: boolean
}

export default function DropdownCari<T>({
  labelAria,
  nilai,
  lihat,
  cari,
  kunci,
  kolom,
  terpilih,
  onPilih,
  lebar,
}: {
  /** Label VERBATIM medan / kepala kolom - untuk pembaca layar (label tampilnya milik `Medan` / kepala grid). */
  labelAria: string
  /** Teks nilai yang tampil di medan. */
  nilai: string
  /** Mode lihat (`IsView == 'true'`): teks, tidak dapat dibuka. */
  lihat: boolean
  /** Sumber daftar menurut kata `Search` (kosong = semua). */
  cari: (kata: string) => Promise<HasilCari<T>>
  kunci: (t: T) => string
  /** Kepala kolom, isi per baris, dan (opsional) `grid-template-columns` bila lebih dari dua kolom - dalam `rem`. */
  kolom: { judul: readonly string[]; isi: (t: T) => readonly string[]; lebar?: string }
  /** Baris yang sedang terpilih - jadi baris aktif saat daftar dibuka. */
  terpilih: (t: T) => boolean
  onPilih: (t: T) => void
  /** Panel lebar (daftar berkolom banyak). */
  lebar?: boolean
}) {
  const id = useId()
  const idDaftar = `${id}-daftar`
  const idButir = (i: number): string => `${id}-butir-${i}`

  const akar = useRef<HTMLDivElement>(null)
  const pemicu = useRef<HTMLDivElement>(null)
  // Sumber dan penanda baris terpilih lewat ref: fungsi baru setiap render tidak memicu pembacaan ulang.
  const cariRef = useRef(cari)
  cariRef.current = cari
  const terpilihRef = useRef(terpilih)
  terpilihRef.current = terpilih

  const [buka, setBuka] = useState(false)
  const [rataKanan, setRataKanan] = useState(false)
  const [kata, setKata] = useState('')
  const [daftar, setDaftar] = useState<T[] | null>(null)
  const [lebih, setLebih] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [aktif, setAktif] = useState(-1)

  // Setiap perubahan kata membatalkan permintaan sebelumnya: jawaban lambat tidak menimpa hasil yang lebih baru.
  useEffect(() => {
    if (!buka) return
    let batal = false
    const jam = setTimeout(
      () => {
        cariRef
          .current(kata)
          .then(({ tampil, lebih: terpotong }) => {
            if (batal) return
            setDaftar(tampil)
            setLebih(terpotong)
            setGalat(null)
            const kini = tampil.findIndex((t) => terpilihRef.current(t))
            setAktif(tampil.length === 0 ? -1 : Math.max(0, kini))
          })
          .catch((e: unknown) => {
            if (!batal) setGalat(e)
          })
      },
      kata === '' ? 0 : JEDA_MS,
    )
    return () => {
      batal = true
      clearTimeout(jam)
    }
  }, [buka, kata])

  // Klik di luar dropdown menutupnya.
  useEffect(() => {
    if (!buka) return
    const tutupDiLuar = (e: MouseEvent): void => {
      if (e.target instanceof Node && !akar.current?.contains(e.target)) setBuka(false)
    }
    document.addEventListener('mousedown', tutupDiLuar)
    return () => {
      document.removeEventListener('mousedown', tutupDiLuar)
    }
  }, [buka])

  useEffect(() => {
    // `idButir` hanya bergantung pada `id` yang tetap sepanjang umur komponen.
    if (buka && aktif >= 0) document.getElementById(idButir(aktif))?.scrollIntoView({ block: 'nearest' })
  }, [buka, aktif])

  if (lihat) return <span>{nilai}</span>

  function bukaDaftar(): void {
    const r = pemicu.current?.getBoundingClientRect()
    const tepi = akar.current?.closest('.mpnl')?.getBoundingClientRect().right ?? window.innerWidth
    setRataKanan(r !== undefined && r.left + (lebar ? LEBAR_PANEL_LEBAR : LEBAR_PANEL) > tepi)
    setKata('')
    setDaftar(null)
    setLebih(false)
    setGalat(null)
    setAktif(-1)
    setBuka(true)
  }

  function tutup(): void {
    setBuka(false)
    pemicu.current?.focus()
  }

  function pilih(t: T): void {
    onPilih(t)
    tutup()
  }

  function tombolPemicu(e: KeyboardEvent<HTMLDivElement>): void {
    if (e.key === 'Enter' || e.key === ' ' || e.key === 'ArrowDown') {
      e.preventDefault()
      if (!buka) bukaDaftar()
    } else if (e.key === 'Escape' && buka) {
      e.preventDefault()
      setBuka(false)
    }
  }

  function tombolCari(e: KeyboardEvent<HTMLInputElement>): void {
    const n = daftar?.length ?? 0
    const langkah = { ArrowDown: 1, ArrowUp: -1, PageDown: LANGKAH_HALAMAN, PageUp: -LANGKAH_HALAMAN }[e.key]
    if (langkah !== undefined) {
      e.preventDefault()
      setAktif((a) => geserAktif(a, langkah, n))
    } else if (e.key === 'Enter') {
      e.preventDefault()
      const t = daftar?.[aktif]
      if (t !== undefined) pilih(t)
    } else if (e.key === 'Escape') {
      e.preventDefault()
      tutup()
    } else if (e.key === 'Tab') {
      setBuka(false)
    }
  }

  const kelas = ['mpnl-dropdown mpnl-dropdown--sel', lebar ? 'mpnl-dropdown--lebar' : '', rataKanan ? 'mpnl-dropdown--kanan' : '']
    .filter((k) => k !== '')
    .join(' ')
  const gayaKolom = kolom.lebar === undefined ? undefined : { gridTemplateColumns: kolom.lebar }

  return (
    <div className={kelas} ref={akar}>
      <div
        ref={pemicu}
        className="field__input mpnl-dropdown__pemicu"
        role="combobox"
        tabIndex={0}
        aria-haspopup="listbox"
        aria-expanded={buka}
        aria-controls={idDaftar}
        aria-label={labelAria}
        onClick={() => {
          if (buka) setBuka(false)
          else bukaDaftar()
        }}
        onKeyDown={tombolPemicu}
      >
        <span className="mpnl-dropdown__nilai">{nilai}</span>
        <span className="mpnl-dropdown__panah" aria-hidden="true">
          ▾
        </span>
      </div>
      {buka && (
        <div className="mpnl-dropdown__panel">
          <input
            className="field__input"
            type="search"
            placeholder={PEMILIH_MPNL.search}
            aria-label={PEMILIH_MPNL.search}
            aria-controls={idDaftar}
            aria-activedescendant={aktif >= 0 && daftar !== null ? idButir(aktif) : undefined}
            autoFocus
            value={kata}
            onChange={(e) => {
              setKata(e.target.value)
            }}
            onKeyDown={tombolCari}
          />
          {daftar === null && galat === null && <Memuat />}
          {galat !== null && <Gagal galat={galat} />}
          {daftar !== null && daftar.length === 0 && galat === null && <Kosong pesan={LAIN_MPNL.kosong} />}
          {daftar !== null && daftar.length > 0 && (
            <>
              <ul className="mpnl-dropdown__daftar" id={idDaftar} role="listbox" aria-label={labelAria}>
                {/* Kepala di DALAM daftar (menempel saat digulir): lebarnya sama dengan baris - kolom `fr` tetap sejajar
                    walau daftar bergulir (scrollbar mengurangi lebar daftar, bukan kepala di luarnya). */}
                <li className="mpnl-dropdown__kepala" role="presentation" aria-hidden="true" style={gayaKolom}>
                  {kolom.judul.map((j) => (
                    <span key={j}>{j}</span>
                  ))}
                </li>
                {daftar.map((t, i) => (
                  <li
                    key={`${kunci(t)}-${i}`}
                    id={idButir(i)}
                    className={i === aktif ? 'mpnl-dropdown__butir mpnl-dropdown__butir--aktif' : 'mpnl-dropdown__butir'}
                    style={gayaKolom}
                    role="option"
                    aria-selected={i === aktif}
                    onMouseEnter={() => {
                      setAktif(i)
                    }}
                    onMouseDown={(e) => {
                      // Mousedown, bukan click: fokus tidak sempat pindah dan menutup daftar lebih dulu.
                      e.preventDefault()
                      pilih(t)
                    }}
                  >
                    {kolom.isi(t).map((v, k) => (
                      <span key={k} className={k === 0 ? 'mpnl-dropdown__id' : undefined}>
                        {v}
                      </span>
                    ))}
                  </li>
                ))}
              </ul>
              {lebih && <p className="mpnl-dropdown__catatan">{LAIN_MPNL.dropdownTerpotong}</p>}
            </>
          )}
        </div>
      )}
    </div>
  )
}
