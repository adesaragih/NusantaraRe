// Shell aplikasi — F0.3, pola `REFERENSI_UI/frontend/src/App.tsx`.
//
// Sidebar (terlipat / laci di ponsel) + topbar + menu profil + palet Ctrl+K,
// dengan `PagarGalat` membungkus isinya.
//
// ⛔ MENUNYA HANYA YANG BERBUKTI KORPUS — dua butir, satu kelompok. Lihat
// `assets/labels.ts` `MENU`. Sebelas modul lain, butir `BelumTersedia`, dan
// layar-layar yang di Pega dibuka DARI DALAM kasus (16 flow action + 3 popup)
// TIDAK berdiri di sini. Menu yang tidak ada di sistem lama adalah menu yang
// dikarang, dan penjaganya ada di `Shell.test.ts`.
//
// ⚠️ Perilaku yang ditiru persis referensi: menu profil tertutup oleh klik di
// luar DAN oleh Esc; lebar tablet melipat panel; laci ponsel menutup sesudah
// sebuah menu dipilih.

import { useCallback, useEffect, useRef, useState, type ReactNode } from 'react'

import { MENU, PERAN_ID, PRODUK } from '../assets/labels'
import { PagarGalat } from '../PagarGalat'
import { type Sesi } from '../store/sesi'
import { KelompokMenu } from './KelompokMenu'
import { IkonCari, IkonPanel, IkonPengguna, IkonTutup } from './ui/dasar'

/** Halaman yang Shell dapat tampilkan. */
export type Halaman = 'inbox' | 'register' | 'detail'

/** Satu butir menu. */
interface ButirMenu {
  halaman: Halaman
  label: string
}

/**
 * Butir sidebar — DUA, dan keduanya berbukti.
 *
 * `detail` sengaja TIDAK di sini: di Pega ia flow action
 * (`ViewClaimDetailLifeGCNM`) yang dibuka dari dalam kasus, bukan menu.
 */
const BUTIR: readonly ButirMenu[] = [
  { halaman: 'inbox', label: MENU.inbox },
  { halaman: 'register', label: MENU.register },
]

export interface ShellProps {
  masuk: Sesi
  halaman: Halaman
  onPindah: (h: Halaman) => void
  children: ReactNode
}

export function Shell({ masuk, halaman, onPindah, children }: ShellProps) {
  const [terlipat, setTerlipat] = useState(false)
  const [laciBuka, setLaciBuka] = useState(false)
  const [profilBuka, setProfilBuka] = useState(false)
  const [paletBuka, setPaletBuka] = useState(false)
  const profilRef = useRef<HTMLDivElement>(null)

  // Ctrl+K / Cmd+K membuka palet; Esc menutup palet dan menu profil.
  useEffect(() => {
    function tekan(e: KeyboardEvent): void {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault()
        setPaletBuka((b) => !b)
        return
      }
      if (e.key === 'Escape') {
        setPaletBuka(false)
        setProfilBuka(false)
        setLaciBuka(false)
      }
    }
    window.addEventListener('keydown', tekan)
    return () => {
      window.removeEventListener('keydown', tekan)
    }
  }, [])

  // ⛔ Klik DI LUAR menutup menu profil. Tanpa ini menu menggantung terbuka
  // sesudah pemakai mengklik ke tempat lain, dan menutupi isi layar.
  useEffect(() => {
    if (!profilBuka) return
    function klik(e: MouseEvent): void {
      if (profilRef.current !== null && !profilRef.current.contains(e.target as Node)) {
        setProfilBuka(false)
      }
    }
    document.addEventListener('mousedown', klik)
    return () => {
      document.removeEventListener('mousedown', klik)
    }
  }, [profilBuka])

  const pilih = useCallback(
    (h: Halaman) => {
      onPindah(h)
      // Laci ponsel menutup sesudah memilih - bila tidak, menu menutupi
      // halaman yang baru saja dibuka.
      setLaciBuka(false)
      setPaletBuka(false)
    },
    [onPindah],
  )

  const judulAktif = BUTIR.find((b) => b.halaman === halaman)?.label ?? MENU.inbox

  return (
    <div className={`shell${terlipat ? ' shell--terlipat' : ''}${laciBuka ? ' shell--laci' : ''}`}>
      <aside className="shell__sidebar" aria-label="Navigasi utama">
        <div className="shell__merek">
          <span className="shell__merek-nama">{PRODUK.nama}</span>
          <span className="shell__merek-sub">{PRODUK.sub}</span>
        </div>

        <nav className="shell__nav">
          <KelompokMenu nama={MENU.kelompokClaimLife} memuatAktif>
            {BUTIR.map((b) => (
              <button
                key={b.halaman}
                type="button"
                className={`shell__butir${halaman === b.halaman ? ' shell__butir--aktif' : ''}`}
                aria-current={halaman === b.halaman ? 'page' : undefined}
                onClick={() => {
                  pilih(b.halaman)
                }}
              >
                {b.label}
              </button>
            ))}
          </KelompokMenu>
        </nav>
      </aside>

      {/* Tirai laci ponsel. */}
      {laciBuka && (
        <button
          type="button"
          className="shell__tirai"
          aria-label="Tutup menu"
          onClick={() => {
            setLaciBuka(false)
          }}
        />
      )}

      <div className="shell__kanan">
        <header className="shell__topbar">
          <button
            type="button"
            className="shell__ikon"
            aria-label={terlipat ? 'Bentangkan panel' : 'Lipat panel'}
            onClick={() => {
              setTerlipat((v) => !v)
            }}
          >
            <IkonPanel />
          </button>
          <button
            type="button"
            className="shell__ikon shell__ikon--laci"
            aria-label="Buka menu"
            onClick={() => {
              setLaciBuka(true)
            }}
          >
            <IkonPanel />
          </button>

          <h1 className="shell__judul">{judulAktif}</h1>

          <button
            type="button"
            className="shell__cari"
            onClick={() => {
              setPaletBuka(true)
            }}
          >
            <IkonCari />
            <span>Cari menu</span>
            <kbd>Ctrl K</kbd>
          </button>

          <div className="shell__profil" ref={profilRef}>
            <button
              type="button"
              className="shell__ikon"
              aria-haspopup="menu"
              aria-expanded={profilBuka}
              onClick={() => {
                setProfilBuka((v) => !v)
              }}
            >
              <IkonPengguna />
            </button>
            {profilBuka && (
              <div className="shell__profil-menu" role="menu">
                <p className="shell__profil-akun">{masuk.akunID}</p>
                <ul className="shell__profil-peran">
                  {masuk.peran.map((p) => (
                    <li key={p}>{PERAN_ID[p]}</li>
                  ))}
                </ul>
                {/* ⛔ Tombol Keluar DIBUANG: tanpa masuk tidak ada keluar.
                    Identitas datang dari env saat aplikasi menyala. */}
                <p className="shell__profil-stub">mode stub</p>
              </div>
            )}
          </div>
        </header>

        <main className="shell__isi">
          <PagarGalat>{children}</PagarGalat>
        </main>
      </div>

      {paletBuka && (
        <PaletMenuSederhana
          butir={BUTIR}
          onTutup={() => {
            setPaletBuka(false)
          }}
          onPilih={pilih}
        />
      )}
    </div>
  )
}

/**
 * Palet menu Ctrl+K.
 *
 * ⚠️ Versi SEDERHANA dari `components/PaletMenu.tsx` referensi: yang di sana
 * mencari entitas master dinamis, sedangkan menu kita dua butir tetap.
 * Menyalin pencari entitas berarti menyalin kode untuk daftar yang tidak ada.
 * Perilaku yang ditiru: fokus masuk ke kotak, panah naik/turun, Enter memilih,
 * Esc menutup, dan fokus KEMBALI ke elemen sebelumnya.
 */
function PaletMenuSederhana({
  butir,
  onTutup,
  onPilih,
}: {
  butir: readonly ButirMenu[]
  onTutup: () => void
  onPilih: (h: Halaman) => void
}) {
  const [kueri, setKueri] = useState('')
  const [sorot, setSorot] = useState(0)
  const kotak = useRef<HTMLInputElement>(null)
  const fokusSebelum = useRef<Element | null>(null)

  useEffect(() => {
    fokusSebelum.current = document.activeElement
    kotak.current?.focus()
    return () => {
      // Tanpa ini fokus tertinggal di <body>: Tab berikutnya melompat ke awal
      // halaman, dan pemakai papan ketik kehilangan tempatnya.
      if (fokusSebelum.current instanceof HTMLElement) fokusSebelum.current.focus()
    }
  }, [])

  const cocok = butir.filter((b) => b.label.toLowerCase().includes(kueri.trim().toLowerCase()))

  return (
    <div className="palet" role="dialog" aria-modal="true" aria-label="Cari menu">
      <div className="palet__kotak">
        <input
          ref={kotak}
          className="palet__isian"
          value={kueri}
          placeholder="Cari menu…"
          onChange={(e) => {
            setKueri(e.target.value)
            setSorot(0)
          }}
          onKeyDown={(e) => {
            if (e.key === 'ArrowDown') {
              e.preventDefault()
              setSorot((s) => Math.min(s + 1, cocok.length - 1))
            } else if (e.key === 'ArrowUp') {
              e.preventDefault()
              setSorot((s) => Math.max(s - 1, 0))
            } else if (e.key === 'Enter') {
              const p = cocok[sorot]
              if (p !== undefined) onPilih(p.halaman)
            } else if (e.key === 'Escape') {
              onTutup()
            }
          }}
        />
        <ul className="palet__hasil">
          {cocok.map((b, i) => (
            <li key={b.halaman}>
              <button
                type="button"
                className={`palet__hasil-butir${i === sorot ? ' palet__hasil-butir--sorot' : ''}`}
                onClick={() => {
                  onPilih(b.halaman)
                }}
              >
                {b.label}
              </button>
            </li>
          ))}
          {cocok.length === 0 && <li className="palet__kosong">Tidak ada menu yang cocok.</li>}
        </ul>
        <button type="button" className="palet__tutup" aria-label="Tutup" onClick={onTutup}>
          <IkonTutup />
        </button>
      </div>
    </div>
  )
}
