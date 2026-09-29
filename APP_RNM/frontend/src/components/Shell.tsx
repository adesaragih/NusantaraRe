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
//
// ⚠️ TAMPILAN mengikuti `D:/XML/RNM_BRD/workpage-template.html` (29-09-2026,
// keputusan pemilik proyek: palet dan huruf persis template). Tiga lebar:
// ponsel (<768px) panel jadi laci off-canvas; tablet (768–1023px) panel mulai
// terciut jadi ikon; desktop panel terbentang. SATU tombol menu melayani
// ketiganya — ciut/bentang di layar lebar, buka laci di ponsel.

import { useCallback, useEffect, useRef, useState, type ReactNode } from 'react'

import {
  BERANDA,
  KERANGKA,
  KETERANGAN_BELUM_DIMIGRASI,
  MENU,
  MODUL,
} from '../assets/labels'
import { PERAN_ID, PRODUK } from '../assets/labels.claimlife'
import { ENTRI_MENU, type ModulTetap } from '../lib/daftarMenu'
import { singkatanUnik } from '../lib/singkatan'
import { PagarGalat } from '../PagarGalat'
import { type Sesi } from '../store/sesi'
import { KelompokMenu } from './KelompokMenu'
import { PaletMenu } from './PaletMenu'
import { useTema } from '../hooks/useTema'
import {
  IkonBulan,
  IkonCari,
  IkonChevron,
  IkonMatahari,
  IkonMenu,
  IkonRumah,
  IkonTutup,
} from './ui/dasar'

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
  // Kelompok ke-18 (ralat: folder korpus 20) — sesi Treaty Contract Out, aditif.
  MODUL.treatyContractOut,
].map((nama) => ({ nama, butir: butirKelompok(nama) }))

/**
 * Lencana dua huruf per kelompok — pengganti ikon saat panel terciut.
 *
 * ⚠️ BUKAN ikon Lucide per kelompok: delapan belas kelompok berbagi lima
 * keluarga (Claim, Komite Claim, Endorsement, Master, NB…), jadi ikon per
 * keluarga membuat empat Claim dan empat Komite tampil kembar di panel ikon.
 * `singkatanUnik` menjamin tak ada dua lencana yang sama.
 */
const LENCANA = singkatanUnik(KELOMPOK.map((k) => k.nama))

/** Huruf awal dua kata pertama — "Nusantara Re" → "NR", "UJI-ADMIN" → "UA". */
function inisial(teks: string): string {
  const kata = teks.split(/[^A-Za-z0-9]+/).filter(Boolean)
  const dua =
    kata.length >= 2
      ? kata
          .slice(0, 2)
          .map((k) => k.charAt(0))
          .join('')
      : teks.replace(/[^A-Za-z0-9]/g, '').slice(0, 2)
  return dua.toUpperCase()
}

/** Batas lebar yang sama dengan template: di bawahnya panel menjadi laci. */
const KUERI_LEBAR = '(min-width: 768px)'

function layarLebarKini(): boolean {
  return typeof window === 'undefined' || typeof window.matchMedia !== 'function'
    ? true
    : window.matchMedia(KUERI_LEBAR).matches
}

/** Mengikuti apakah layar selebar tablet ke atas. */
function useLayarLebar(): boolean {
  const [lebar, setLebar] = useState(layarLebarKini)
  useEffect(() => {
    if (typeof window.matchMedia !== 'function') return
    const mq = window.matchMedia(KUERI_LEBAR)
    const ubah = (): void => {
      setLebar(mq.matches)
    }
    mq.addEventListener('change', ubah)
    return () => {
      mq.removeEventListener('change', ubah)
    }
  }, [])
  return lebar
}

export interface ShellProps {
  masuk: Sesi
  halaman: Halaman
  onPindah: (h: Halaman) => void
  children: ReactNode
}

export function Shell({ masuk, halaman, onPindah, children }: ShellProps) {
  const lebar = useLayarLebar()
  const { tema, balik: balikTema } = useTema()
  // Tablet (768–1023px) mulai dengan panel terciut, seperti template:
  // ruang kerja lebih lega tanpa menyembunyikan menu.
  const [terlipat, setTerlipat] = useState(
    () => typeof window !== 'undefined' && window.innerWidth >= 768 && window.innerWidth < 1024,
  )
  const [laciBuka, setLaciBuka] = useState(false)
  const [profilBuka, setProfilBuka] = useState(false)
  const [paletBuka, setPaletBuka] = useState(false)
  const profilRef = useRef<HTMLDivElement>(null)
  const kananRef = useRef<HTMLDivElement>(null)
  const tombolMenuRef = useRef<HTMLButtonElement>(null)
  const tombolTutupRef = useRef<HTMLButtonElement>(null)
  const laciPernahBuka = useRef(false)
  const fokusKembali = useRef(true)

  /** Panel HANYA terciut di layar lebar; di ponsel laci selalu penuh. */
  const terciut = lebar && terlipat
  const laciTerbuka = !lebar && laciBuka

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

  // Layar dilebarkan ke ukuran tablet: laci ditutup TANPA memindah fokus.
  useEffect(() => {
    if (lebar) {
      fokusKembali.current = false
      setLaciBuka(false)
    }
  }, [lebar])

  // Laci ponsel: fokus pindah ke tombol Tutup saat terbuka, KEMBALI ke tombol
  // menu saat tertutup; isi di belakangnya `inert` supaya Tab tidak lolos ke
  // halaman yang tertutup tirai, dan halaman tidak ikut tergulir.
  useEffect(() => {
    if (kananRef.current !== null) kananRef.current.inert = laciTerbuka
    if (laciTerbuka) {
      laciPernahBuka.current = true
      tombolTutupRef.current?.focus()
      const lama = document.body.style.overflow
      document.body.style.overflow = 'hidden'
      return () => {
        document.body.style.overflow = lama
      }
    }
    if (laciPernahBuka.current && fokusKembali.current) tombolMenuRef.current?.focus()
    fokusKembali.current = true
    return undefined
  }, [laciTerbuka])

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

  const labelMenu = lebar
    ? terlipat
      ? KERANGKA.perluasMenu
      : KERANGKA.ciutkanMenu
    : KERANGKA.bukaMenu

  const peranUtama = masuk.peran[0]
  const sebutanPeran =
    peranUtama === undefined
      ? ''
      : PERAN_ID[peranUtama] + (masuk.peran.length > 1 ? ` +${masuk.peran.length - 1}` : '')

  const tandaMerek = (
    <span className="shell__merek-tanda" aria-hidden="true">
      {inisial(PRODUK.nama)}
    </span>
  )

  return (
    <div
      className={`shell${terciut ? ' shell--terlipat' : ''}${laciTerbuka ? ' shell--laci' : ''}`}
    >
      <a className="shell__lewati" href="#isi">
        {KERANGKA.lewati}
      </a>

      <aside className="shell__sidebar" id="shell-sidebar" aria-label={KERANGKA.menuSamping}>
        <div className="shell__merek">
          {tandaMerek}
          <span className="shell__merek-teks">
            <span className="shell__merek-nama">{PRODUK.nama}</span>
            <span className="shell__merek-sub">{PRODUK.sub}</span>
          </span>
          {/* Tombol tutup: hanya tampil di ponsel. */}
          <button
            ref={tombolTutupRef}
            type="button"
            className="shell__ikon shell__tutup"
            aria-label={KERANGKA.tutupMenu}
            onClick={() => {
              setLaciBuka(false)
            }}
          >
            <IkonTutup ukuran={20} />
          </button>
        </div>

        <nav className="shell__nav" aria-label={KERANGKA.navUtama}>
          <ul className="shell__daftar">
            {/* Beranda berdiri SENDIRI di atas kelompok - ia bukan modul.  */}
            <li>
              <button
                type="button"
                className={`shell__butir shell__butir--beranda${
                  halaman === 'beranda' ? ' shell__butir--aktif' : ''
                }`}
                aria-current={halaman === 'beranda' ? 'page' : undefined}
                title={terciut ? BERANDA.judul : undefined}
                onClick={() => {
                  pilih('beranda')
                }}
              >
                <IkonRumah />
                <span className="shell__label">{BERANDA.judul}</span>
              </button>
            </li>

            {KELOMPOK.map((k) => (
              <li key={k.nama}>
                {k.butir.length === 0 ? (
                  /* Kelompok tanpa butir tetap BERDIRI dan menyebut sebabnya.
                     Menyembunyikannya membuat aplikasi tampak lengkap padahal
                     empat belas modul belum ada. Ia bukan tombol: tidak ada
                     yang dapat dibuka di dalamnya. */
                  <div
                    className="kelompok kelompok--kosong"
                    title={terciut ? `${k.nama} — ${KETERANGAN_BELUM_DIMIGRASI}` : undefined}
                  >
                    <span className="kelompok__lencana" aria-hidden="true">
                      {LENCANA.get(k.nama)}
                    </span>
                    <span className="kelompok__teks">
                      {k.nama}
                      <span className="shell__belum">{KETERANGAN_BELUM_DIMIGRASI}</span>
                    </span>
                  </div>
                ) : (
                  <KelompokMenu
                    nama={k.nama}
                    lencana={LENCANA.get(k.nama) ?? ''}
                    memuatAktif={k.butir.some((b) => b.halaman === halaman)}
                    terciut={terciut}
                    onBentang={() => {
                      setTerlipat(false)
                    }}
                  >
                    {k.butir.map((b) => (
                      <button
                        key={b.halaman}
                        type="button"
                        className={`shell__butir shell__butir--anak${
                          halaman === b.halaman ? ' shell__butir--aktif' : ''
                        }`}
                        aria-current={halaman === b.halaman ? 'page' : undefined}
                        // Nama menu VERBATIM korpus bisa lebih panjang dari
                        // panel dan terpotong elipsis — tooltip memuat utuhnya.
                        title={b.label}
                        onClick={() => {
                          pilih(b.halaman)
                        }}
                      >
                        <span className="shell__label">{b.label}</span>
                      </button>
                    ))}
                  </KelompokMenu>
                )}
              </li>
            ))}
          </ul>
        </nav>
      </aside>

      {/* Tirai laci ponsel: klik untuk menutup. */}
      <div
        className="shell__tirai"
        aria-hidden="true"
        onClick={() => {
          setLaciBuka(false)
        }}
      />

      <div className="shell__kanan" ref={kananRef}>
        <header className="shell__topbar">
          <button
            ref={tombolMenuRef}
            type="button"
            className="shell__ikon"
            aria-controls="shell-sidebar"
            aria-expanded={lebar ? !terlipat : laciBuka}
            aria-label={labelMenu}
            title={labelMenu}
            onClick={() => {
              if (lebar) setTerlipat((v) => !v)
              else setLaciBuka(true)
            }}
          >
            <IkonMenu />
          </button>

          {/* Nama merek: hanya tampil di ponsel, saat panel tersembunyi. */}
          <button
            type="button"
            className="shell__topbar-merek"
            aria-label={`${PRODUK.nama}, ${KERANGKA.keBeranda}`}
            onClick={() => {
              pilih('beranda')
            }}
          >
            {tandaMerek}
            <span aria-hidden="true">{PRODUK.nama}</span>
          </button>

          {/* Judul halaman untuk pembaca layar; judul yang TERLIHAT adalah
              kepala tiap halaman, seperti `page-head` di template. */}
          <h1 className="sr-only">{judulAktif}</h1>

          <button
            type="button"
            className="shell__cari"
            aria-label={KERANGKA.cariMenu}
            onClick={() => {
              setPaletBuka(true)
            }}
          >
            <IkonCari ukuran={20} />
            <span className="shell__cari-teks" aria-hidden="true">
              {KERANGKA.cariMenu}
            </span>
            <kbd aria-hidden="true">{KERANGKA.pintasCari}</kbd>
          </button>

          {/* `.topbar-actions` template: tombol tema + profil, didorong ke kanan. */}
          <div className="shell__aksi">
            {/* Sakelar tema. `aria-pressed` + nama TETAP "Mode gelap": pembaca
                layar mendengar "Mode gelap, tertekan/tidak". Ikonnya tema TUJUAN. */}
            <button
              type="button"
              className="shell__ikon"
              aria-pressed={tema === 'gelap'}
              aria-label={KERANGKA.modeGelap}
              title={tema === 'gelap' ? KERANGKA.modeTerang : KERANGKA.modeGelap}
              onClick={balikTema}
            >
              {tema === 'gelap' ? <IkonMatahari /> : <IkonBulan />}
            </button>
            <div className="shell__profil" ref={profilRef}>
              <button
                type="button"
                className="shell__profil-pemicu"
                aria-expanded={profilBuka}
                aria-controls="shell-profil"
                aria-label={`${KERANGKA.profil} ${masuk.akunID}`}
                onClick={() => {
                  setProfilBuka((v) => !v)
                }}
              >
                <span className="shell__avatar" aria-hidden="true">
                  {inisial(masuk.akunID)}
                </span>
                <span className="shell__profil-teks" aria-hidden="true">
                  <strong>{masuk.akunID}</strong>
                  <span>{sebutanPeran}</span>
                </span>
                <span className="shell__profil-panah" aria-hidden="true">
                  <IkonChevron />
                </span>
              </button>
              {profilBuka && (
                <div className="shell__profil-menu" id="shell-profil">
                  <p className="shell__profil-akun">{masuk.akunID}</p>
                  <ul className="shell__profil-peran">
                    {masuk.peran.map((p) => (
                      <li key={p}>
                        {PERAN_ID[p]}
                        <span>{p}</span>
                      </li>
                    ))}
                  </ul>
                  {/* ⛔ Tombol Keluar DIBUANG: tanpa masuk tidak ada keluar.
                      Identitas datang dari env saat aplikasi menyala. */}
                  <p className="shell__profil-stub">{KERANGKA.modeStub}</p>
                </div>
              )}
            </div>
          </div>
        </header>

        <main className="shell__isi" id="isi" tabIndex={-1}>
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
