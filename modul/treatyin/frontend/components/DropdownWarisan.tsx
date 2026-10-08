// Pemilih bercari layar Treaty In.
//
// ⛔ DIPINDAHKAN dari `pages/FormKontrakTreatyIn.tsx` 5 Oktober 2026
// (pemindahan MURNI), lalu DIBERI PAPAN TIK pada langkah berikutnya —
// dua langkah terpisah, supaya uji yang merah dapat dibedakan antara
// "pindahnya salah" dan "tambahannya salah".
//
// ⭐ POLANYA DITIRU dari `modul/masterproductnamelife/frontend/components/
// DropdownCari.tsx` — dibaca, bukan diimpor: impor lintas modul dilarang.
// Yang ditiru: Enter/Spasi/panah membuka, panah dan PageUp/PageDown
// menggeser, Enter memilih, Escape menutup tanpa memilih, baris aktif
// digulirkan ke dalam pandangan, dan `aria-activedescendant` menyebut
// baris itu kepada pembaca layar.
//
// ⚠️ SATU BEDA YANG DISENGAJA, dan ia bukan kelalaian: `DropdownCari`
// MENOLAK ketikan bebas (nilainya hanya dari daftar). Pemilih ini HARUS
// dapat diketik — permintaan pemilik proses 5 Oktober 2026, sebab 131
// cedant lebih cepat disaring dengan tiga huruf daripada digulir.
//
// ⚠️ BEDA KEDUA, dan ia milik kita: NAMA KEMBAR diberi pengenalnya.
// Modul contoh tidak menghadapi soal itu; kita menghadapinya pada 36
// nama cedant.

import { useEffect, useId, useRef, useState, type KeyboardEvent } from 'react'

import { Gagal } from '../../../../inti/frontend/components/ui/dasar'
import type { PilihanWarisan } from '../api'
import { FORM_KONTRAK } from '../labels'
import { saringTerdekat } from '../saring'

/** Langkah `PageUp` / `PageDown` — sama dengan modul contoh. */
const LANGKAH_HALAMAN = 10

/**
 * Pemilih "Choose Ceding" / "Choose Source of Business" — DAFTAR PILIHAN
 * YANG DAPAT DIKETIK.
 *
 * ⛔ Bentuknya dropdown, tetapi kotaknya menerima ketikan dan daftarnya
 * MENYARING saat diketik — permintaan pemilik proses 5 Oktober 2026. Daftar
 * pilihan biasa tidak cukup di sini: cedant ada 131 baris dan asal bisnis 96,
 * dan menggulir 131 baris untuk menemukan satu nama lebih lambat daripada
 * mengetik tiga huruf.
 *
 * ⚠️ NAMA KEMBAR DIBERI PENGENALNYA. Terukur: 36 nama cedant dipakai lebih
 * dari satu pengenal. Dua baris bernama sama tanpa pembeda membuat yang
 * memilih menebak, dan tebakannya menulis pengenal yang salah ke kontrak —
 * salah yang baru terbaca saat rekonsiliasi. Hanya yang kembar yang diberi;
 * membubuhkannya pada semua baris membuat daftar sulit dibaca tanpa menolong
 * siapa pun.
 */
export default function DropdownWarisan({
  label,
  ambil,
  nilai,
  onPilih,
}: {
  label: string
  ambil: () => Promise<PilihanWarisan[]>
  nilai: string
  onPilih: (p: PilihanWarisan) => void
}) {
  const [daftar, setDaftar] = useState<PilihanWarisan[]>([])
  const [galat, setGalat] = useState<unknown>(null)
  const [ketik, setKetik] = useState('')
  const [buka, setBuka] = useState(false)
  // ⭐ Baris AKTIF — yang panah geser dan yang Enter pilih. `-1` berarti nol
  // baris aktif, dan Enter saat itu tidak memilih apa pun.
  const [aktif, setAktif] = useState(-1)
  const id = useId()
  const idBaris = (i: number): string => `${id}-baris-${i}`
  const kotak = useRef<HTMLInputElement>(null)

  useEffect(() => {
    // ⛔ Penjaga `hidup`: layar dapat ditinggalkan sebelum jawabannya tiba.
    let hidup = true
    ambil().then(
      (d) => {
        if (hidup) setDaftar(d)
      },
      (e: unknown) => {
        if (hidup) setGalat(e)
      },
    )
    return () => {
      hidup = false
    }
  }, [ambil])

  // ⚠️ Baris aktif digulirkan ke dalam pandangan — tanpa ini, panah-bawah
  // pada daftar 131 baris menggeser penanda ke luar layar dan yang memakai
  // papan tik kehilangan jejaknya.
  useEffect(() => {
    if (buka && aktif >= 0) {
      document.getElementById(idBaris(aktif))?.scrollIntoView({ block: 'nearest' })
    }
  }, [buka, aktif, idBaris])

  if (galat !== null) return <Gagal galat={galat} />

  /** Label satu baris — pengenal dibubuhkan HANYA pada nama kembar. */
  const teks = (b: PilihanWarisan) => (b.kembar ? `${b.nama} — ${b.id}` : b.nama)

  // ⛔ Hanya yang PALING MIRIP yang tampil — permintaan pemilik proses
  // 5 Oktober 2026. Substring biasa memunculkan yang sekadar KEBETULAN
  // mengandung ketikan, dan pada 131 cedant yang namanya hampir semua
  // berawalan "ASURANSI", itu berarti separuh daftar ikut muncul setiap kali.
  const q = ketik.trim().toLowerCase()
  const cocok = saringTerdekat(daftar, q)

  /**
   * ⛔ DIJEPIT di kedua ujung, tidak melingkar.
   *
   * Daftar yang melingkar membuat panah-bawah di baris terakhir melompat ke
   * baris pertama, dan yang menekan terus tidak pernah tahu ia sudah di
   * ujung. Sama dengan `geserAktif` di modul contoh.
   */
  const geser = (langkah: number): void => {
    setAktif((a) => {
      if (cocok.length === 0) return -1
      return Math.min(cocok.length - 1, Math.max(0, a + langkah))
    })
  }

  const pilih = (b: PilihanWarisan): void => {
    onPilih(b)
    setBuka(false)
    setAktif(-1)
  }

  function papanTik(e: KeyboardEvent<HTMLInputElement>): void {
    if (!buka && (e.key === 'ArrowDown' || e.key === 'Enter')) {
      e.preventDefault()
      setBuka(true)
      setAktif(-1)
      return
    }
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      geser(1)
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      geser(-1)
    } else if (e.key === 'PageDown') {
      e.preventDefault()
      geser(LANGKAH_HALAMAN)
    } else if (e.key === 'PageUp') {
      e.preventDefault()
      geser(-LANGKAH_HALAMAN)
    } else if (e.key === 'Enter') {
      // ⛔ `preventDefault` SELALU saat daftarnya terbuka: Enter di dalam
      // form mengirim form, dan form ini belum punya jalur simpan.
      e.preventDefault()
      // ⭐ Tanpa baris tersorot, Enter memilih yang PALING MIRIP (baris
      // pertama saringan) — bila ada ketikan. Keluhan pemakai 7 Oktober 2026:
      // ketik `QUOTA` lalu Enter → tidak ada yang terpilih, kotaknya kembali
      // kosong, dan seluruh isian yang bergantung pada Treaty Type (QS %,
      // persen Retention/Cession) ikut lenyap — terbaca "semua ke-reset".
      // Autocomplete Pega menyorot hasil pertama; Enter memilihnya.
      const b = cocok[aktif] ?? (q !== '' ? cocok[0] : undefined)
      if (b !== undefined) pilih(b)
    } else if (e.key === 'Escape') {
      // ⚠️ Menutup TANPA memilih, dan nilainya kembali seperti semula —
      // Escape yang meninggalkan pilihan setengah jalan lebih buruk
      // daripada Escape yang tidak berbuat apa-apa.
      e.preventDefault()
      setBuka(false)
      setAktif(-1)
    } else if (e.key === 'Tab') {
      setBuka(false)
    }
  }

  return (
    <div className="field trin__ketikpilih">
      <label className="field__label">{label}</label>
      <input
        className="field__input"
        type="text"
        ref={kotak}
        role="combobox"
        aria-expanded={buka}
        aria-controls={`${id}-daftar`}
        aria-autocomplete="list"
        aria-activedescendant={buka && aktif >= 0 ? idBaris(aktif) : undefined}
        aria-label={label}
        autoComplete="off"
        value={buka ? ketik : nilai}
        placeholder={nilai === '' ? '' : undefined}
        onKeyDown={papanTik}
        onFocus={() => {
          // ⛔ KETIKAN TIDAK DIHAPUS bila daftarnya SUDAH terbuka.
          //
          // Bentuk sebelumnya mengosongkannya pada SETIAP fokus. Akibatnya
          // nyata dan terlihat di tangkapan layar pemilik proses 6 Oktober
          // 2026: diketik `zurich`, daftar menyaring ke satu baris, lalu satu
          // peristiwa fokus susulan — render ulang, klik kembali ke kotak,
          // pengembalian fokus dari daftar — mengembalikan seluruh 131 baris
          // sementara kotaknya masih memperlihatkan `zurich`.
          //
          // ⚠⚠ Dan itu terbaca sebagai SARINGAN YANG RUSAK, padahal
          // saringannya benar. `saringTerdekat('zurich')` terbukti
          // mengembalikan tepat satu baris (`saring.test.ts`).
          //
          // ⭐ Membuka dari keadaan TERTUTUP tetap mengosongkannya — di sana
          // kotak memperlihatkan nilai TERSIMPAN, bukan ketikan, dan membawa
          // ketikan lama ke pembukaan berikutnya menyaring tanpa ada yang
          // mengetiknya.
          if (!buka) {
            setKetik('')
            setBuka(true)
            setAktif(-1)
          }
        }}
        // ⛔ Penutupan DITUNDA satu putaran: `onBlur` menyala sebelum `onClick`
        // baris, dan menutup seketika membuat pilihan tidak pernah sampai.
        onBlur={() => {
          window.setTimeout(() => {
            setBuka(false)
          }, 150)
        }}
        onChange={(e) => {
          setKetik(e.target.value)
          setBuka(true)
          // ⚠️ Ketikan baru MENGGANTI baris aktif: baris ke-3 dari daftar
          // LAMA bukan baris ke-3 dari daftar baru, dan membiarkannya
          // membuat Enter memilih yang tidak dilihat siapa pun.
          // ⭐ Ada ketikan → baris PERTAMA daftar BARU tersorot (yang paling
          // mirip), jadi yang Enter pilih selalu terlihat — seperti
          // autocomplete Pega. Kotak dikosongkan → nol baris tersorot.
          setAktif(e.target.value.trim() === '' ? -1 : 0)
        }}
      />
      {buka && (
        <ul
          id={`${id}-daftar`}
          className="trin__ketikpilih-daftar"
          role="listbox"
          aria-label={label}
        >
          {cocok.length === 0 && (
            <li className="trin__ketikpilih-kosong">{FORM_KONTRAK.pilihKosong}</li>
          )}
          {cocok.map((b, i) => (
            <li key={b.id}>
              <button
                type="button"
                id={idBaris(i)}
                role="option"
                aria-selected={i === aktif}
                className={
                  'trin__ketikpilih-baris' + (i === aktif ? ' trin__ketikpilih-baris--aktif' : '')
                }
                onMouseEnter={() => {
                  setAktif(i)
                }}
                onMouseDown={() => {
                  // ⛔ `onMouseDown`, bukan `onClick`: ia mendahului `onBlur`.
                  pilih(b)
                }}
              >
                {teks(b)}
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
