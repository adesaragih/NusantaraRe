import { useEffect, useMemo, useRef, useState } from 'react'

import { KERANGKA } from '../assets/labels'
import { IkonCari } from './ui/dasar'
import {
  daftarPalet,
  saringPalet,
  type HasilPalet,
  type ModulTetap,
} from '../lib/daftarMenu'

/**
 * PaletMenu — cari menu dengan papan ketik (butir **bg**).
 *
 * Port `REFERENSI_UI/frontend/src/components/PaletMenu.tsx`, disederhanakan
 * pada satu hal: di sana daftarnya digabung dengan entitas master yang App
 * muat dari backend; di sini seluruh entri tetap, jadi ia tidak menerima
 * apa pun dari luar selain penanganan peristiwanya.
 *
 * ⛔ NOL pencarian data dagang. Palet ini membuka layar yang SAMA yang
 * sidebar buka, lewat pemanggil yang SAMA. Nol endpoint baru, nol data
 * dibaca, nol baris tersimpan. Yang berubah hanya panjang jalan ke sana —
 * dan batas itu dipagari uji.
 *
 * ⛔ Ia membaca `lib/daftarMenu.ts`, BUKAN DOM sidebar. Sejak butir bg
 * empat belas kelompok terlipat, dan `KelompokMenu` melepas anak kelompok
 * yang terlipat dari DOM: palet yang membaca DOM akan menjawab "tidak ada"
 * untuk menu yang ADA.
 *
 * ⚠️ Pola papan ketiknya dipinjam apa adanya dari referensi: Escape menutup,
 * panah ber-`preventDefault` (tanpanya panah memindahkan KARET di kotak,
 * bukan sorotan), Enter bersyarat, sorotan melingkar, dan fokus KEMBALI ke
 * elemen sebelumnya saat ditutup.
 */
export function PaletMenu({
  onTutup,
  onPilih,
  modulAktif = null,
}: {
  onTutup: () => void
  /** Membuka satu hasil. Shell meneruskannya ke pemindah halamannya. */
  onPilih: (modul: ModulTetap) => void
  /** Modul aktif dari backend; `null` = semua (lihat `halamanAktif`). */
  modulAktif?: readonly string[] | null
}) {
  const [kueri, setKueri] = useState('')
  const [sorot, setSorot] = useState(0)
  const kotak = useRef<HTMLInputElement>(null)
  /**
   * Elemen yang tadinya berfokus, supaya fokus KEMBALI saat palet ditutup.
   *
   * Tanpa ini, menutup palet meninggalkan fokus di `<body>`: penekanan Tab
   * berikutnya melompat ke awal halaman, dan pemakai papan ketik kehilangan
   * tempatnya. Itu bukan kemewahan aksesibilitas — ia satu-satunya cara
   * palet tidak merugikan pemakai yang paling terbantu olehnya.
   */
  const fokusSebelum = useRef<Element | null>(null)

  const semua = useMemo(() => daftarPalet(modulAktif), [modulAktif])
  const hasil = useMemo(() => saringPalet(semua, kueri), [semua, kueri])

  // Sorotan kembali ke atas setiap kali kuerinya berubah. Tanpa ini sorotan
  // tertinggal di indeks lama yang kini menunjuk baris lain — Enter membuka
  // menu yang TIDAK sedang dilihat pemakai.
  useEffect(() => {
    setSorot(0)
  }, [kueri])

  useEffect(() => {
    fokusSebelum.current = document.activeElement
    const t = setTimeout(() => kotak.current?.focus(), 0)
    return () => {
      clearTimeout(t)
      if (fokusSebelum.current instanceof HTMLElement) fokusSebelum.current.focus()
    }
  }, [])

  function pilih(h: HasilPalet): void {
    onPilih(h.modul)
    onTutup()
  }

  return (
    <div
      className="palet__tirai"
      // Klik di luar menutup — perilaku yang sama dengan modal lain.
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) onTutup()
      }}
    >
      <div className="palet" role="dialog" aria-modal="true" aria-label={KERANGKA.cariMenu}>
        <label className="palet__kepala">
          <IkonCari ukuran={20} />
          <input
            ref={kotak}
            className="palet__isian"
            type="text"
            value={kueri}
            placeholder="Cari menu…"
            aria-label={KERANGKA.cariMenu}
            onChange={(e) => {
              setKueri(e.target.value)
            }}
            onKeyDown={(e) => {
              if (e.key === 'Escape') {
                e.preventDefault()
                onTutup()
                return
              }
              if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
                e.preventDefault()
                if (hasil.length === 0) return
                const arah = e.key === 'ArrowDown' ? 1 : -1
                setSorot((s) => (s + arah + hasil.length) % hasil.length)
                return
              }
              if (e.key === 'Enter') {
                const h = hasil[sorot]
                if (h !== undefined) {
                  e.preventDefault()
                  pilih(h)
                }
              }
            }}
          />
          <kbd aria-hidden="true">Esc</kbd>
        </label>
        {hasil.length === 0 ? (
          // ⛔ Daftar kosong SENYAP terbaca sebagai layar rusak. Ia menyebut
          // apa yang dicari, supaya pemakai tahu ia mengetik bukan menunggu.
          <p className="palet__kosong">Tidak ada menu yang cocok dengan “{kueri}”.</p>
        ) : (
          <ul className="palet__hasil">
            {hasil.map((h, i) => (
              <li key={h.kunci}>
                <button
                  type="button"
                  className={`palet__hasil-butir${
                    i === sorot ? ' palet__hasil-butir--sorot' : ''
                  }`}
                  // Sorotan mengikuti tetikus juga, supaya papan ketik dan
                  // tetikus tidak menunjuk dua baris yang berbeda.
                  onMouseEnter={() => {
                    setSorot(i)
                  }}
                  onClick={() => {
                    pilih(h)
                  }}
                >
                  <span className="palet__label">{h.label}</span>
                  <span className="palet__kelompok">{h.kelompok}</span>
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}
