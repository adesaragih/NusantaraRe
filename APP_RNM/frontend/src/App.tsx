import { useEffect, useState } from 'react'

import Beranda from './Beranda'
import { Shell } from './inti/components/Shell'
import { BelumTersedia } from './inti/components/ui/dasar'
import { ambilMenu, ambilModulAktif } from './inti/klien'
import { modulDipasang, type KeadaanMenuTabel } from './inti/lib/daftarMenu'
import { pelakuStub } from './inti/store/sesi'
import { ENTRI_MENU, halamanAktif, MODUL_FRONTEND, type Halaman } from './modul/daftar'

// App = identitas + Shell.
//
// ⛔ TIDAK ADA FORM LOGIN `[perintah work owner 27-09-2026]`. Aplikasi membuka
// Inbox langsung; identitasnya datang dari env saat menyala (`store/sesi`).
//
// ⛔ Bila stub-nya mati, layar MENYATAKANNYA - bukan pecah, dan bukan pula
// diam-diam menampilkan daftar kosong yang terbaca "tidak ada pekerjaan".
export default function App() {
  const masuk = pelakuStub()
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

  if (masuk === null) {
    return (
      // ⚠️ `BelumTersedia` hanya menerima `apa`: parameter `sebab`-nya
      // sengaja dibuang di REFERENSI_UI ronde 70 karena diterima lalu
      // diabaikan. Keterangan teknisnya karena itu berdiri di sini.
      <section>
        <BelumTersedia apa="Identitas pelaku" />
        <p className="polis__catatan" role="note">
          Selama IAM belum terpasang, identitas dibentuk dari env Vite dan
          menuntut <code>VITE_AUTH_STUB=true</code>. Tanpa itu setiap
          permintaan ditolak backend, dan layar ini tidak dapat
          menampilkan pekerjaan siapa pun.
        </p>
      </section>
    )
  }

  return (
    <Shell masuk={masuk} halaman={halaman} onPindah={setHalaman} menu={ENTRI_MENU} menuTabel={menuTabel}>
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
