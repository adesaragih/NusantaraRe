import { useCallback, useEffect, useMemo, useState } from 'react'

import Beranda from './Beranda'
import GantiSandi from '../inti/frontend/components/GantiSandi'
import Login, { KerangkaMasuk, pesanGagalLogin } from '../inti/frontend/components/Login'
import { Shell } from '../inti/frontend/components/Shell'
import { ApiFailure, ambilMenu, ambilModulAktif, ambilSesiSaya, keluarLogin, type ProfilLogin } from '../inti/frontend/klien'
import KelolaUser from '../inti/frontend/kelolauser/KelolaUser'
import { LOGIN } from '../inti/frontend/labels'
import {
  HALAMAN_KELOLA_USER,
  KODE_MENU_KELOLA_USER,
  modulDipasang,
  modulUntukAkun,
  type KeadaanMenuTabel,
} from '../inti/frontend/lib/daftarMenu'
import { bolehMasukStub, PERISTIWA_SESI_BERAKHIR, pelakuStub, sesiDariProfil, type Sesi } from '../inti/frontend/store/sesi'
import { ENTRI_MENU, halamanAktif, MODUL_FRONTEND, type Halaman } from './daftar'

// App = identitas + Shell.
//
// ⭐ LOGIN SUNGGUHAN (M_LOGIN_GO, keputusan work owner 01-10-2026): saat
// menyala App bertanya `GET /api/auth/saya`. 401 = layar login; akun yang
// wajib ganti sandi = layar ganti sandi; selebihnya Shell. 401 di tengah
// pemakaian (`PERISTIWA_SESI_BERAKHIR`) kembali ke layar login.
//
// Mode stub (`VITE_AUTH_STUB=true`) tetap: identitas dari env, tanpa layar
// login - untuk pengembangan.
export default function App() {
  const stub = bolehMasukStub()
  // `undefined` = sedang diperiksa; `null` = belum login; selebihnya profil.
  const [profil, setProfil] = useState<ProfilLogin | null | undefined>(stub ? null : undefined)
  const [galatSesi, setGalatSesi] = useState<unknown>(null)
  const [gantiSandi, setGantiSandi] = useState(false)
  const periksaSesi = useCallback(() => {
    setGalatSesi(null)
    setProfil(undefined)
    ambilSesiSaya().then(setProfil, (g: unknown) => {
      if (g instanceof ApiFailure && g.status === 401) setProfil(null)
      else setGalatSesi(g)
    })
  }, [])
  useEffect(() => {
    if (stub) return
    periksaSesi()
    const berakhir = () => {
      setProfil(null)
    }
    window.addEventListener(PERISTIWA_SESI_BERAKHIR, berakhir)
    return () => {
      window.removeEventListener(PERISTIWA_SESI_BERAKHIR, berakhir)
    }
  }, [stub, periksaSesi])
  const masuk: Sesi | null = stub ? pelakuStub() : profil ? sesiDariProfil(profil) : null
  // ⛔ Beranda layar AWAL sejak butir bg: ia yang menyebut modul mana
  // yang sudah ada dan mana yang belum. Membuka langsung ke Inbox membuat
  // aplikasi tampak hanya punya satu modul.
  const [halaman, setHalaman] = useState<Halaman>('beranda')
  // Penghitung pilihan menu - `PropsRute.ketukMenu`. Memilih menu halaman yang
  // SEDANG tampil tidak mengubah `halaman`, jadi rute modul tidak akan tahu
  // menunya diklik ulang tanpa sinyal kedua ini.
  const [ketukMenu, setKetukMenu] = useState(0)
  const pilihDariMenu = useCallback((h: Halaman) => {
    setHalaman(h)
    setKetukMenu((k) => k + 1)
  }, [])
  // Modul yang dipasang backend (MODUL_AKTIF, refactor bentuk B). `null` =
  // belum terbaca atau gagal dibaca: SEMUA menu tampil, persis seperti
  // sebelum MODUL_AKTIF ada - satu pembacaan yang gagal tidak mengosongkan
  // aplikasi.
  const [modulAktif, setModulAktif] = useState<readonly string[] | null>(null)
  useEffect(() => {
    let batal = false
    ambilModulAktif().then(
      (m) => {
        if (!batal) setModulAktif(m)
      },
      () => {
        // Tetap `null`: layar modul menampilkan galat backend-nya sendiri.
      },
    )
    return () => {
      batal = true
    }
  }, [])
  // Menu dari tabel M_NAV_MENU (brief menu 30-09-2026). `null` = sedang
  // dimuat; gagal = sidebar MENAMPILKAN galatnya, bukan menu kosong.
  //
  // ⛔ Dibaca SESUDAH identitas diketahui, dan dibaca ULANG saat akunnya
  // berganti (01-10-2026): menu disaring per akun, dan `GET /api/menu`
  // tanpa sesi dijawab 401 - yang, bila dikirim sebelum login, memicu
  // `PERISTIWA_SESI_BERAKHIR` di tengah pemeriksaan sesi dan membuat layar
  // login berkedip.
  //
  // ⛔ Akun yang WAJIB ganti sandi bukan pemegang menu (backend tidak menaruh
  // menunya di context): menunya dibaca SESUDAH sandi diganti. Membacanya di
  // layar ganti sandi dijawab 401, dan 401 itu melempar akun baru kembali ke
  // layar login berulang-ulang (temuan /code-review).
  const kunciMenu = stub ? 'stub' : profil && !profil.wajibGantiSandi ? profil.akunId : null
  // Naik saat admin mengubah akunnya SENDIRI di Kelola User: menu dibaca ulang.
  const [versiMenu, setVersiMenu] = useState(0)
  const [menuTabel, setMenuTabel] = useState<KeadaanMenuTabel>(null)
  useEffect(() => {
    setMenuTabel(null)
    if (kunciMenu === null) return
    let batal = false
    ambilMenu().then(
      (menu) => {
        if (!batal) setMenuTabel({ menu })
      },
      (galat: unknown) => {
        if (!batal) setMenuTabel({ galat })
      },
    )
    return () => {
      batal = true
    }
  }, [kunciMenu, versiMenu])
  // Menu akun (M_LOGIN_GO_MENU, Kelola User 01-10-2026): modul yang menunya
  // tidak dipegang TIDAK DIPASANG - layarnya tidak ada dan tidak satu pun
  // permintaannya berangkat (backend menjawabnya 403). `null` = mode stub,
  // tanpa saringan akun - juga bila backend (versi lama) tidak mengirim medan
  // `menu`: backend itu memang belum menyaring per akun.
  const menuAkun = !stub && profil && Array.isArray(profil.menu) ? profil.menu : null
  const modulBoleh = useMemo(
    () => modulUntukAkun(modulAktif, menuAkun, MODUL_FRONTEND.map((m) => m.nama)),
    [modulAktif, menuAkun],
  )
  const bolehKelola = menuAkun !== null && menuAkun.includes(KODE_MENU_KELOLA_USER)
  // Admin mengubah akunnya sendiri: profil (menu) dan sidebar dibaca ulang.
  const segarkanDiri = useCallback(() => {
    ambilSesiSaya().then(setProfil, () => {
      // 401 sudah memicu PERISTIWA_SESI_BERAKHIR; galat lain menunggu muat ulang.
    })
    setVersiMenu((v) => v + 1)
  }, [])
  // Daftar modul aktif tiba SESUDAH pemakai sempat membuka halaman modul yang
  // ternyata nonaktif (semua menu tampil selama daftarnya `null`): rute modul
  // itu dilepas, jadi halamannya kembali ke Beranda alih-alih layar kosong
  // (temuan /code-review). Tanpa MODUL_AKTIF tidak pernah terjadi. Sama bila
  // menu akunnya dicabut, dan bila Kelola User tidak (lagi) dipegang.
  useEffect(() => {
    if (!halamanAktif(halaman, modulBoleh) || (halaman === HALAMAN_KELOLA_USER && !bolehKelola)) {
      setHalaman('beranda')
    }
  }, [halaman, modulBoleh, bolehKelola])

  if (!stub && galatSesi !== null) {
    return (
      <KerangkaMasuk judul={LOGIN.masuk} sub={pesanGagalLogin(galatSesi)}>
        <div className="halaman-masuk__form">
          {/* Kalimat backend tetap tampil KECIL di bawahnya - ia yang menyebut
              sebab sebenarnya (mis. migrasi yang belum dijalankan) bagi IT. */}
          {galatSesi instanceof ApiFailure && galatSesi.detail.message !== undefined && (
            <p className="halaman-masuk__bantuan" role="alert">
              {galatSesi.detail.message}
            </p>
          )}
          <button type="button" className="halaman-masuk__tombol" onClick={periksaSesi}>
            {LOGIN.cobaLagi}
          </button>
        </div>
      </KerangkaMasuk>
    )
  }
  if (!stub && profil === undefined) {
    // ⛔ Layar NETRAL - latar halaman, tanpa kartu, teksnya baru tampak bila
    // pemeriksaan lebih dari 600 ms (laporan work owner 01-10-2026 "halaman
    // login kedip"). Tujuannya belum diketahui: kartu login (ungu) ATAU Shell
    // (terang). Layar antara yang meniru salah satunya berkedip bagi yang lain -
    // dulu teks polos sebelum kartu login, lalu latar ungu sebelum Shell.
    return (
      <main className="periksa-sesi" aria-busy="true">
        <p role="status" className="periksa-sesi__status">
          {LOGIN.memuatSesi}
        </p>
      </main>
    )
  }
  if (!stub && profil === null) {
    return <Login onMasuk={setProfil} />
  }
  if (!stub && profil && (profil.wajibGantiSandi || gantiSandi)) {
    return (
      <GantiSandi
        wajib={profil.wajibGantiSandi}
        onSelesai={(p) => {
          setGantiSandi(false)
          setProfil(p)
        }}
        onBatal={() => {
          setGantiSandi(false)
        }}
      />
    )
  }
  if (masuk === null) {
    // Mode stub tanpa identitas tidak mungkin (bolehMasukStub benar), dan
    // login sungguhan sudah dijawab di atas.
    return null
  }

  return (
    <Shell
      masuk={masuk}
      halaman={halaman}
      onPindah={setHalaman}
      menu={ENTRI_MENU}
      menuTabel={menuTabel}
      onKeluar={
        stub
          ? undefined
          : () => {
              // Sesi dibuang di layar walau logout ke server gagal - pemakai
              // yang menekan Keluar harus keluar.
              keluarLogin().finally(() => {
                setProfil(null)
              })
            }
      }
      onGantiSandi={
        stub
          ? undefined
          : () => {
              setGantiSandi(true)
            }
      }
    >
      {halaman === 'beranda' && <Beranda masuk={masuk} onBuka={setHalaman} modulAktif={modulBoleh} />}
      {halaman === HALAMAN_KELOLA_USER && bolehKelola && <KelolaUser akunSaya={masuk.akunID} onDiriBerubah={segarkanDiri} />}
      {/*
        Refactor bentuk B (30-09-2026): setiap modul AKTIF merender halamannya
        sendiri (`modul/<nama>/rute.tsx`) dan menyimpan keadaannya sendiri -
        kasus, polis, atau kasus komite yang sedang dibuka. Rutenya TETAP
        terpasang selama modulnya aktif, jadi keadaan itu bertahan saat pemakai
        pindah halaman, persis seperti ketika ia hidup di sini.
      */}
      {MODUL_FRONTEND.filter((m) => modulDipasang(m.nama, modulBoleh)).map((m) => (
        <m.Rute key={m.nama} halaman={halaman} masuk={masuk} onPindah={setHalaman} />
      ))}
    </Shell>
  )
}
