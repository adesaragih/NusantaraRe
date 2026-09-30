// Shell aplikasi — F0.3, pola `REFERENSI_UI/frontend/src/App.tsx`.
//
// Sidebar (terlipat / laci di ponsel) + topbar + menu profil + palet Ctrl+K,
// dengan `PagarGalat` membungkus isinya.
//
// ⛔ SEJAK 30-09-2026 MENU DARI TABEL `M_NAV_MENU` (brief menu, permintaan work
// owner): sidebar dan palet dirakit dari `GET /api/menu` - golongan TREATY,
// FACULTATIVE, KLAIM, MASTER sebagai kepala bagian, lalu kelompok modul, lalu
// butir - DIPOTONG dengan rute yang terdaftar di `modul/daftar.ts`
// (`susunMenu`). Bila `GET /api/menu` gagal, sidebar menampilkan galatnya.
// Catatan di bawah (bg) tetap berlaku untuk ISI-nya; urutan dan pengelompokan
// kini milik tabel.
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

import { useCallback, useEffect, useMemo, useRef, useState, type ReactNode } from 'react'

import {
  BERANDA,
  KERANGKA,
  KETERANGAN_BELUM_DIMIGRASI,
  MENU,
  PERAN_ID,
  PRODUK,
} from '../labels'
import {
  butirDatar,
  entriAplikasi,
  HALAMAN_BERANDA,
  susunMenu,
  type EntriMenu,
  type KeadaanMenuTabel,
} from '../lib/daftarMenu'
import { singkatanUnik } from '../lib/singkatan'
import { PagarGalat } from '../PagarGalat'
import { type Sesi } from '../store/sesi'
import { KelompokMenu } from './KelompokMenu'
import { PaletMenu } from './PaletMenu'
import { useTema } from '../hooks/useTema'
import {
  Gagal,
  IkonBulan,
  IkonCari,
  IkonChevron,
  IkonMatahari,
  IkonMenu,
  IkonRumah,
  IkonTutup,
  Memuat,
} from './ui/dasar'

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

export interface ShellProps<H extends string> {
  masuk: Sesi
  /** Halaman yang sedang tampil - termasuk yang dibuka dari dalam kasus. */
  halaman: string
  onPindah: (h: H | typeof HALAMAN_BERANDA) => void
  children: ReactNode
  /**
   * Rute menu frontend (`ENTRI_MENU`) - PEMOTONG pohon tabel: hanya butir
   * tabel yang berute di sini yang tampil, di sidebar maupun palet.
   */
  menu: readonly EntriMenu<H>[]
  /** Pembacaan `GET /api/menu` (M_NAV_MENU): `null` = memuat. */
  menuTabel: KeadaanMenuTabel
}

export function Shell<H extends string>({
  masuk,
  halaman,
  onPindah,
  children,
  menu,
  menuTabel,
}: ShellProps<H>) {
  // Pohon tabel dipotong rute frontend - SATU hasil untuk sidebar dan palet.
  const tersusun = useMemo(
    () => (menuTabel !== null && 'menu' in menuTabel ? susunMenu(menuTabel.menu, menu) : null),
    [menuTabel, menu],
  )
  // Palet sebelum/tanpa menu tabel: Beranda saja.
  const berandaSaja = useMemo(() => entriAplikasi(menu), [menu])
  // Lencana dua huruf per kelompok — pengganti ikon saat panel terciut.
  // ⚠️ BUKAN ikon Lucide per kelompok: kelompok berbagi lima keluarga (Claim,
  // Komite Claim, Endorsement, Master, NB…), jadi ikon per keluarga membuat
  // empat Claim dan empat Komite tampil kembar di panel ikon. `singkatanUnik`
  // menjamin tak ada dua lencana yang sama - dihitung atas kelompok yang TAMPIL
  // (menu dari tabel, 30-09-2026), jadi daftarnya tidak lagi diketik di sini.
  const lencana = useMemo(
    () => singkatanUnik(tersusun?.golongan.flatMap((g) => g.kelompok.map((k) => k.nama)) ?? []),
    [tersusun],
  )
  // Baris tabel yang tidak tampil DICATAT, supaya tabel yang mendahului
  // kodenya (atau salah ketik KODE, atau butir di bawah kelompok yang belum
  // dimigrasi) terlihat di konsol.
  useEffect(() => {
    if (tersusun === null) return
    if (tersusun.tanpaRute.length > 0) {
      console.warn(
        `menu: ${tersusun.tanpaRute.length} butir M_NAV_MENU tanpa rute frontend, tidak tampil: ${tersusun.tanpaRute.join(', ')}`,
      )
    }
    if (tersusun.terlipat.length > 0) {
      console.warn(
        `menu: ${tersusun.terlipat.length} butir M_NAV_MENU di bawah kelompok DIMIGRASI='0', tidak tampil: ${tersusun.terlipat.join(', ')}`,
      )
    }
  }, [tersusun])
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
    (h: H | typeof HALAMAN_BERANDA) => {
      onPindah(h)
      // Laci ponsel menutup sesudah memilih - bila tidak, menu menutupi
      // halaman yang baru saja dibuka.
      setLaciBuka(false)
      setPaletBuka(false)
    },
    [onPindah],
  )

  const judulAktif =
    menu.find((e) => e.modul === halaman)?.label ?? MENU.inbox

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

  // Logo resmi perusahaan (diserahkan work owner 30-09-2026) menggantikan
  // kotak berinisial; gambarnya dipasang di CSS (`.shell__merek-tanda`),
  // bukan `src={...}`, supaya penjaga "nol src dinamis" tetap utuh.
  // Dekoratif: nama produk ditulis di sebelahnya.
  const tandaMerek = (
    <span className="shell__merek-tanda" aria-hidden="true" />
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
                  halaman === HALAMAN_BERANDA ? ' shell__butir--aktif' : ''
                }`}
                aria-current={halaman === HALAMAN_BERANDA ? 'page' : undefined}
                title={terciut ? BERANDA.judul : undefined}
                onClick={() => {
                  pilih(HALAMAN_BERANDA)
                }}
              >
                <IkonRumah />
                <span className="shell__label">{BERANDA.judul}</span>
              </button>
            </li>

            {menuTabel === null && (
              <li className="shell__menu-keadaan">
                <Memuat pesan={KERANGKA.memuatMenu} />
              </li>
            )}
            {menuTabel !== null && 'galat' in menuTabel && (
              /* GET /api/menu gagal: galatnya TAMPIL di sidebar - bukan menu
                 kosong diam-diam (brief menu 30-09-2026). Beranda tetap. */
              <li className="shell__menu-keadaan">
                <Gagal galat={menuTabel.galat} />
              </li>
            )}
            {tersusun !== null && tersusun.golongan.length === 0 && (
              /* Tabel menjawab, tetapi tak satu baris pun dapat tampil (kosong,
                 seluruhnya STATUS_AKTIF '0', atau tak satu pun berute):
                 DIKATAKAN, bukan sidebar yang hanya berisi Beranda tanpa sebab. */
              <li className="shell__menu-keadaan">
                <p role="status">{KERANGKA.menuKosong}</p>
              </li>
            )}
            {/* Kepala bagian = GROUPMENU (TREATY, FACULTATIVE, KLAIM, MASTER),
                urutan dari backend. Modul NONAKTIF (MODUL_AKTIF) tidak dikirim
                butirnya, jadi kelompoknya tidak tampil (`susunMenu`). */}
            {tersusun?.golongan.map((g) => (
              <li key={g.kode} className="shell__golongan">
                <p className="shell__golongan-judul" aria-hidden="true">
                  {g.kode}
                </p>
                <ul className="shell__daftar" aria-label={g.kode}>
                  {g.kelompok.map((k) => {
                    const datar = butirDatar(k.butir)
                    return (
                      <li key={k.kode}>
                        {!k.dimigrasi ? (
                          /* Kelompok `DIMIGRASI = '0'` tetap BERDIRI dan menyebut
                             sebabnya. Menyembunyikannya membuat aplikasi tampak
                             lengkap padahal enam belas modul belum ada. Ia bukan
                             tombol: tidak ada yang dapat dibuka di dalamnya. */
                          <div
                            className="kelompok kelompok--kosong"
                            title={terciut ? `${k.nama} — ${KETERANGAN_BELUM_DIMIGRASI}` : undefined}
                          >
                            <span className="kelompok__lencana" aria-hidden="true">
                              {lencana.get(k.nama)}
                            </span>
                            <span className="kelompok__teks">
                              {k.nama}
                              <span className="shell__belum">{KETERANGAN_BELUM_DIMIGRASI}</span>
                            </span>
                          </div>
                        ) : datar !== undefined ? (
                          /* Kelompok beranggota SATU butir bertanda `datar`: satu tombol
                             langsung, tanpa judul kelompok yang dilipat dan tanpa anak
                             (`ButirMenuModul.datar`). */
                          <button
                            type="button"
                            className={`shell__butir shell__butir--datar${halaman === datar.halaman ? ' shell__butir--aktif' : ''}`}
                            aria-current={halaman === datar.halaman ? 'page' : undefined}
                            title={datar.label}
                            onClick={() => {
                              pilih(datar.halaman)
                            }}
                          >
                            <span className="kelompok__lencana" aria-hidden="true">
                              {lencana.get(k.nama)}
                            </span>
                            <span className="shell__label">{datar.label}</span>
                          </button>
                        ) : (
                          <KelompokMenu
                            nama={k.nama}
                            lencana={lencana.get(k.nama) ?? ''}
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
                    )
                  })}
                </ul>
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
              pilih(HALAMAN_BERANDA)
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
          menu={tersusun?.entri ?? berandaSaja}
        />
      )}
    </div>
  )
}
