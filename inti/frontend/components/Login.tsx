// Halaman login - M_LOGIN_GO (keputusan work owner 01-10-2026).
//
// Tampilan: `loginbaru.html` lampiran work owner 01-10-2026 - kartu kaca dua
// kolom (ilustrasi meja kerja di kiri pada layar lebar, form di kanan), logo
// "Nusantara Re | Login aplikasi" bertulisan hitam. CSS `.halaman-masuk*` di
// styles.css. Teks `LOGIN`.
//
// ⛔ Sandi hanya hidup di state komponen ini sepanjang satu percobaan, lalu
// dikosongkan - tidak pernah ke penyimpanan peramban atau log. Cookie sesi
// HttpOnly dipasang backend; layar tidak menyentuhnya.

import { useState, type ReactNode } from 'react'

import { ApiFailure, masukLogin, type ProfilLogin } from '../klien'
import { LOGIN } from '../labels'
import IlustrasiLogin from './IlustrasiLogin'

/** Pesan gagal dari STATUS jawaban - teks backend tidak ditampilkan apa adanya. */
export function pesanGagalLogin(galat: unknown): string {
  if (galat instanceof ApiFailure) {
    if (galat.status === 401) return LOGIN.salah
    if (galat.status === 423) return LOGIN.terkunci
    if (galat.status === 503) return LOGIN.tidakTersedia
  }
  return LOGIN.gagal
}

/** Kerangka bersama Login dan Ganti sandi: kartu kaca, ilustrasi, logo. */
export function KerangkaMasuk({ judul, sub, children }: { judul: string; sub: string; children: ReactNode }) {
  return (
    <main className="halaman-masuk">
      <div className="halaman-masuk__kartu">
        <div className="halaman-masuk__panel">
          <IlustrasiLogin />
        </div>
        <section className="halaman-masuk__sisi-form">
          <div className="halaman-masuk__kotak-form">
            {/* Gambar LATAR (pola `.shell__merek-tanda`): pagar unduhdokumen menolak
                atribut sumber gambar di layar. */}
            <span className="halaman-masuk__logo" role="img" aria-label={LOGIN.logo} />
            <h1 className="halaman-masuk__judul">{judul}</h1>
            <p className="halaman-masuk__sub">{sub}</p>
            {children}
          </div>
        </section>
      </div>
    </main>
  )
}

/** Isian sandi dengan tombol tampilkan/sembunyikan - ikon Lucide (lisensi ISC). */
export function IsianSandi({
  id,
  label,
  placeholder,
  value,
  onChange,
  autoComplete,
  autoFocus,
}: {
  id: string
  label: string
  placeholder: string
  value: string
  onChange: (v: string) => void
  autoComplete: 'current-password' | 'new-password'
  autoFocus?: boolean
}) {
  const [terlihat, setTerlihat] = useState(false)
  return (
    <>
      <label className="sr-only" htmlFor={id}>
        {label}
      </label>
      <div className="halaman-masuk__sandi">
        <input
          className="halaman-masuk__input"
          id={id}
          name={id}
          type={terlihat ? 'text' : 'password'}
          placeholder={placeholder}
          autoComplete={autoComplete}
          autoFocus={autoFocus}
          required
          value={value}
          onChange={(e) => {
            onChange(e.target.value)
          }}
        />
        <button
          className="halaman-masuk__toggle"
          type="button"
          aria-label={terlihat ? LOGIN.sembunyikanSandi : LOGIN.tampilkanSandi}
          aria-pressed={terlihat}
          aria-controls={id}
          onClick={() => {
            setTerlihat((v) => !v)
          }}
        >
          {terlihat ? (
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
              <path d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0" />
              <circle cx="12" cy="12" r="3" />
            </svg>
          ) : (
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
              <path d="M10.733 5.076a10.744 10.744 0 0 1 11.205 6.575 1 1 0 0 1 0 .696 10.747 10.747 0 0 1-1.444 2.49" />
              <path d="M14.084 14.158a3 3 0 0 1-4.242-4.242" />
              <path d="M17.479 17.499a10.75 10.75 0 0 1-15.417-5.151 1 1 0 0 1 0-.696 10.75 10.75 0 0 1 4.446-5.143" />
              <path d="m2 2 20 20" />
            </svg>
          )}
        </button>
      </div>
    </>
  )
}

export default function Login({ onMasuk }: { onMasuk: (p: ProfilLogin) => void }) {
  const [akun, setAkun] = useState('')
  const [sandi, setSandi] = useState('')
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<string | null>(null)

  return (
    <KerangkaMasuk judul={LOGIN.judul} sub={LOGIN.sub}>
      <form
        className="halaman-masuk__form"
        onSubmit={(e) => {
          e.preventDefault()
          if (sibuk) return
          setSibuk(true)
          setGalat(null)
          masukLogin(akun, sandi).then(
            (p) => {
              setSandi('')
              setSibuk(false)
              onMasuk(p)
            },
            (g: unknown) => {
              setSandi('')
              setSibuk(false)
              setGalat(pesanGagalLogin(g))
            },
          )
        }}
      >
        {galat !== null && (
          <div className="alert alert--error" role="alert">
            {galat}
          </div>
        )}
        <label className="sr-only" htmlFor="masuk-akun">
          {LOGIN.akun}
        </label>
        <input
          className="halaman-masuk__input"
          id="masuk-akun"
          name="username"
          type="text"
          placeholder={LOGIN.isianAkun}
          autoComplete="username"
          autoFocus
          required
          value={akun}
          onChange={(e) => {
            setAkun(e.target.value)
          }}
        />
        <IsianSandi
          id="masuk-sandi"
          label={LOGIN.sandi}
          placeholder={LOGIN.isianSandi}
          value={sandi}
          onChange={setSandi}
          autoComplete="current-password"
        />
        <button className="halaman-masuk__tombol" type="submit" disabled={sibuk || akun.trim() === '' || sandi === ''}>
          {sibuk ? LOGIN.memproses : LOGIN.masuk}
        </button>
        <p className="halaman-masuk__bantuan">{LOGIN.bantuan}</p>
      </form>
    </KerangkaMasuk>
  )
}
