// Shell aplikasi — F0.3, pola `REFERENSI_UI/frontend/src/App.tsx`.
//
// Sidebar (terlipat / laci di ponsel) + topbar + menu profil + palet Ctrl+K,
// dengan `PagarGalat` membungkus isinya.
//
// ⛔ TUJUH BELAS KELOMPOK, EMPAT BUTIR — butir **bg**, 28-09-2026.
//
// Kelompoknya nama FOLDER korpus `D:/XML/RNM_BRD/` apa adanya. Butirnya
// hanya untuk tiga modul yang punya bukti XML: Claim Life (dua), PremiumList
// Life (satu), Komite Claim Life (satu). **Empat belas kelompok lain berdiri
// terlipat, TANPA butir**, berketerangan `belum dimigrasi`.
//
// ⛔ Nol butir dikarang. Menu yang tidak ada di sistem lama adalah menu yang
// dikarang, dan penjaganya ada di `Shell.test.ts`.
//
// ⚠️ Kelompok kosong tetap BERDIRI, tidak disembunyikan. Aplikasi yang
// menampilkan tiga modul dari tujuh belas tampak lengkap padahal tidak — dan
// layar yang tampak lengkap padahal tidak adalah layar yang tidak akan dicari
// lagi (pelajaran butir av).
//
// ⚠️ Perilaku yang ditiru persis referensi: menu profil tertutup oleh klik di
// luar DAN oleh Esc; lebar tablet melipat panel; laci ponsel menutup sesudah
// sebuah menu dipilih.

import { useCallback, useEffect, useRef, useState, type ReactNode } from 'react'

import { BERANDA, KETERANGAN_BELUM_DIMIGRASI, MENU, MODUL } from '../assets/labels'
import { PERAN_ID, PRODUK } from '../assets/labels.claimlife'
import { ENTRI_MENU, type ModulTetap } from '../lib/daftarMenu'
import { PagarGalat } from '../PagarGalat'
import { type Sesi } from '../store/sesi'
import { KelompokMenu } from './KelompokMenu'
import { PaletMenu } from './PaletMenu'
import { IkonCari, IkonPanel, IkonPengguna } from './ui/dasar'

/** Halaman yang Shell dapat tampilkan. */
// ⚠️ `outstanding` dan `detail` BUKAN butir menu: di Pega keduanya
// dibuka DARI DALAM kasus (flow action), bukan dari navigasi.
export type Halaman = ModulTetap | 'outstanding' | 'detail'

/** Satu butir menu. */
interface ButirMenu {
  halaman: Halaman
  label: string
}

/** Satu kelompok sidebar beserta butirnya. */
interface KelompokSidebar {
  nama: string
  butir: readonly ButirMenu[]
}

/**
 * Ketujuh belas kelompok, berurutan seperti folder korpus.
 *
 * ⛔ Butirnya DITURUNKAN dari `ENTRI_MENU`, bukan diketik ulang. Dua daftar
 * yang masing-masing menyebut menu yang sama adalah dua daftar yang akan
 * menyimpang — dan yang menyimpang tidak akan berbunyi: palet membuka menu
 * yang sidebar tidak punya, atau sebaliknya. `daftarMenu.sinkron.test.ts`
 * menjaganya dua arah; penurunan ini membuat penjagaan itu hampir tak perlu.
 */
function butirKelompok(nama: string): readonly ButirMenu[] {
  return ENTRI_MENU.filter((e) => e.kelompok === nama).map((e) => ({
    halaman: e.modul,
    label: e.label,
  }))
}

const KELOMPOK: readonly KelompokSidebar[] = [
  MODUL.claimFacIn,
  MODUL.claimLife,
  MODUL.claimNonProp,
  MODUL.claimProp,
  MODUL.edmTreatyIn,
  MODUL.endorsementLife,
  MODUL.endorsmentFacIn,
  MODUL.komiteClaimFacIn,
  MODUL.komiteClaimLife,
  MODUL.komiteClaimNonProp,
  MODUL.komiteClaimProp,
  MODUL.masterContractRetroLife,
  MODUL.masterProductNameLife,
  MODUL.nbFacIn,
  MODUL.nbTreatyIn,
  MODUL.premiumListLife,
  MODUL.rnwFacIn,
].map((nama) => ({ nama, butir: butirKelompok(nama) }))

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

  const judulAktif =
    ENTRI_MENU.find((e) => e.modul === halaman)?.label ?? MENU.inbox

  return (
    <div className={`shell${terlipat ? ' shell--terlipat' : ''}${laciBuka ? ' shell--laci' : ''}`}>
      <aside className="shell__sidebar" aria-label="Navigasi utama">
        <div className="shell__merek">
          <span className="shell__merek-nama">{PRODUK.nama}</span>
          <span className="shell__merek-sub">{PRODUK.sub}</span>
        </div>

        <nav className="shell__nav">
          {/* Beranda berdiri SENDIRI di atas kelompok - ia bukan modul.  */}
          <button
            type="button"
            className={`shell__butir shell__butir--beranda${
              halaman === 'beranda' ? ' shell__butir--aktif' : ''
            }`}
            aria-current={halaman === 'beranda' ? 'page' : undefined}
            onClick={() => {
              pilih('beranda')
            }}
          >
            {BERANDA.judul}
          </button>

          {KELOMPOK.map((k) => (
            <KelompokMenu
              key={k.nama}
              nama={k.nama}
              memuatAktif={k.butir.some((b) => b.halaman === halaman)}
            >
              {k.butir.length === 0 ? (
                /* Kelompok tanpa butir tetap BERDIRI dan menyebut sebabnya.
                   Menyembunyikannya membuat aplikasi tampak lengkap padahal
                   empat belas modul belum ada. */
                <p className="shell__belum">{KETERANGAN_BELUM_DIMIGRASI}</p>
              ) : (
                k.butir.map((b) => (
                  <button
                    key={b.halaman}
                    type="button"
                    className={`shell__butir${
                      halaman === b.halaman ? ' shell__butir--aktif' : ''
                    }`}
                    aria-current={halaman === b.halaman ? 'page' : undefined}
                    onClick={() => {
                      pilih(b.halaman)
                    }}
                  >
                    {b.label}
                  </button>
                ))
              )}
            </KelompokMenu>
          ))}
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
        <PaletMenu
          onTutup={() => {
            setPaletBuka(false)
          }}
          onPilih={pilih}
        />
      )}
    </div>
  )
}
