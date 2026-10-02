// Dropdown master - keputusan work owner 02-10-2026 ("perubahan pada tampilan untuk semua Choose ubah jadi
// dropdown saja"): menggantikan tombol `Choose*` + popup FlowAction `Choose*` (PARITAS §4) untuk ketujuh master -
// Ceding, SOB, R/I Risk, Cause Of Loss, Policy Holder, Currency, dan R/I Rate baris `PLAN LIST`.
//
// Yang TETAP seperti XML: nilai HANYA dari daftar master - medannya tidak dapat diketik (`pyReadOnly` true b4040,
// b4428, b7362, b10626, b17062, b28105); kata `Search` dihurufbesarkan server (`SearchPolicyHolder_act` 1 b236) dan
// dicocokkan "Contains"; memilih menyalin ID + nama (`set*_DT`); kolom daftar `ID` / `Name` (`RIRate Name` untuk
// R/I Rate); mode lihat (`IsView`) tidak dapat dibuka - nama tampil sebagai teks (di form, baris `Medan` yang
// menampilkannya; dropdown hanya dirender di mode sunting).
// Yang BERUBAH: daftar dibuka dari medannya sendiri, memuat paling banyak BATAS_DROPDOWN baris (master `CLIENT`
// ratusan ribu baris) - potongan dinyatakan, sisanya dicapai lewat `Search` di dalam dropdown.

import { useEffect, useId, useRef, useState, type KeyboardEvent } from 'react'

import { Gagal, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { cariMaster, type JenisMaster, type NilaiMaster } from '../api'
import { BATAS_DROPDOWN, geserAktif, potongPilihan } from '../bentuk'
import { LAIN_MPNL, PEMILIH_MPNL } from '../labels'

/** Jeda ketik `Search` sebelum RD dibaca. */
const JEDA_MS = 250
/** Lebar panel daftar paling besar (px) - untuk memilih rata kiri atau kanan saat dibuka. */
const LEBAR_PANEL = 440
/** Langkah `PageUp` / `PageDown`. */
const LANGKAH_HALAMAN = 10

export default function DropdownMaster({
  labelAria,
  jenis,
  nilai,
  kolomNama = PEMILIH_MPNL.kolomName,
  lihat,
  onPilih,
}: {
  /** Label VERBATIM medan / kepala kolom - untuk pembaca layar (label tampilnya milik `Medan` / kepala grid). */
  labelAria: string
  jenis: JenisMaster
  /** Nama yang tampil (`.Ceding`, `.SOBName`, ...). */
  nilai: string
  /** Kepala kolom nama (`Name`; `RIRate Name` untuk R/I Rate). */
  kolomNama?: string
  /** Mode lihat (`IsView == 'true'`): baca-saja, tidak dapat dibuka. */
  lihat: boolean
  /** Penerima `set*_DT`: menyalin ID + nama. */
  onPilih: (v: NilaiMaster) => void
}) {
  const id = useId()
  const idDaftar = `${id}-daftar`
  const idButir = (i: number): string => `${id}-butir-${i}`

  const akar = useRef<HTMLDivElement>(null)
  const pemicu = useRef<HTMLDivElement>(null)
  const nilaiRef = useRef(nilai)
  nilaiRef.current = nilai

  const [buka, setBuka] = useState(false)
  const [rataKanan, setRataKanan] = useState(false)
  const [kata, setKata] = useState('')
  const [daftar, setDaftar] = useState<NilaiMaster[] | null>(null)
  const [lebih, setLebih] = useState(false)
  const [galat, setGalat] = useState<unknown>(null)
  const [aktif, setAktif] = useState(-1)

  // Setiap perubahan kata membatalkan permintaan sebelumnya: jawaban lambat tidak menimpa hasil yang lebih baru.
  useEffect(() => {
    if (!buka) return
    let batal = false
    const jam = setTimeout(
      () => {
        cariMaster(jenis, kata, BATAS_DROPDOWN + 1)
          .then((d) => {
            if (batal) return
            const { tampil, lebih: terpotong } = potongPilihan(d.daftar)
            setDaftar(tampil)
            setLebih(terpotong)
            setGalat(null)
            const kini = tampil.findIndex((v) => v.nama === nilaiRef.current)
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
  }, [buka, jenis, kata])

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
    setRataKanan(r !== undefined && r.left + LEBAR_PANEL > tepi)
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

  function pilih(v: NilaiMaster): void {
    onPilih(v)
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
      const v = daftar?.[aktif]
      if (v !== undefined) pilih(v)
    } else if (e.key === 'Escape') {
      e.preventDefault()
      tutup()
    } else if (e.key === 'Tab') {
      setBuka(false)
    }
  }

  const kelas = rataKanan ? 'mpnl-dropdown mpnl-dropdown--sel mpnl-dropdown--kanan' : 'mpnl-dropdown mpnl-dropdown--sel'

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
              <div className="mpnl-dropdown__kepala" aria-hidden="true">
                <span>{PEMILIH_MPNL.kolomId}</span>
                <span>{kolomNama}</span>
              </div>
              <ul className="mpnl-dropdown__daftar" id={idDaftar} role="listbox" aria-label={labelAria}>
                {daftar.map((v, i) => (
                  <li
                    key={`${v.id}-${i}`}
                    id={idButir(i)}
                    className={i === aktif ? 'mpnl-dropdown__butir mpnl-dropdown__butir--aktif' : 'mpnl-dropdown__butir'}
                    role="option"
                    aria-selected={i === aktif}
                    onMouseEnter={() => {
                      setAktif(i)
                    }}
                    onMouseDown={(e) => {
                      // Mousedown, bukan click: fokus tidak sempat pindah dan menutup daftar lebih dulu.
                      e.preventDefault()
                      pilih(v)
                    }}
                  >
                    <span className="mpnl-dropdown__id">{v.id}</span>
                    <span>{v.nama}</span>
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
