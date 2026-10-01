import { useCallback, useEffect, useState } from 'react'

import Beranda from './Beranda'
import GantiSandi from '../inti/frontend/components/GantiSandi'
import Login from '../inti/frontend/components/Login'
import { Shell } from '../inti/frontend/components/Shell'
import { ApiFailure, ambilMenu, ambilModulAktif, ambilSesiSaya, keluarLogin, type ProfilLogin } from '../inti/frontend/klien'
import { LOGIN } from '../inti/frontend/labels'
import { modulDipasang, type KeadaanMenuTabel } from '../inti/frontend/lib/daftarMenu'
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
  const [menuTabel, setMenuTabel] = useState<KeadaanMenuTabel>(null)
  useEffect(() => {
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
  }, [])
  // Daftar modul aktif tiba SESUDAH pemakai sempat membuka halaman modul yang
  // ternyata nonaktif (semua menu tampil selama daftarnya `null`): rute modul
  // itu dilepas, jadi halamannya kembali ke Beranda alih-alih layar kosong
  // (temuan /code-review). Tanpa MODUL_AKTIF tidak pernah terjadi.
  useEffect(() => {
    if (!halamanAktif(halaman, modulAktif)) setHalaman('beranda')
  }, [halaman, modulAktif])

  if (!stub && galatSesi !== null) {
    return (
      <div className="login">
        <div className="login__card">
          <div className="alert alert--error" role="alert">
            {galatSesi instanceof ApiFailure ? galatSesi.message : LOGIN.gagal}
          </div>
          <button type="button" className="btn btn--primary" onClick={periksaSesi}>
            {LOGIN.cobaLagi}
          </button>
        </div>
      </div>
    )
  }
  if (!stub && profil === undefined) {
    return (
      <p role="status" className="polis__catatan">
        {LOGIN.memuatSesi}
      </p>
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
      {halaman === 'beranda' && <Beranda masuk={masuk} onBuka={setHalaman} modulAktif={modulAktif} />}
      {/*
        Refactor bentuk B (30-09-2026): setiap modul AKTIF merender halamannya
        sendiri (`modul/<nama>/rute.tsx`) dan menyimpan keadaannya sendiri -
        kasus, polis, atau kasus komite yang sedang dibuka. Rutenya TETAP
        terpasang selama modulnya aktif, jadi keadaan itu bertahan saat pemakai
        pindah halaman, persis seperti ketika ia hidup di sini.
      */}
      {MODUL_FRONTEND.filter((m) => modulDipasang(m.nama, modulAktif)).map((m) => (
        <m.Rute key={m.nama} halaman={halaman} masuk={masuk} onPindah={setHalaman} />
      ))}
    </Shell>
  )
}
