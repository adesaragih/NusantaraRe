// Halaman login - M_LOGIN_GO (keputusan work owner 01-10-2026).
//
// Tampilan: kartu login `REFERENSI_UI` ronde 277 - satu kartu di tengah, logo,
// latar "Pipeline" (`lib/pipeline.ts`), CSS `.login*` di styles.css. Teksnya
// `LOGIN` (bahasa Indonesia, seperti bingkai aplikasi).
//
// ⛔ Sandi hanya hidup di state komponen ini sepanjang satu percobaan, lalu
// dikosongkan - tidak pernah ke penyimpanan peramban atau log. Cookie sesi
// HttpOnly dipasang backend; layar tidak menyentuhnya.

import { useEffect, useRef, useState } from 'react'

import { ApiFailure, masukLogin, type ProfilLogin } from '../klien'
import { LOGIN } from '../labels'
import { mulaiPipeline } from '../lib/pipeline'
import { Field, FieldSandi } from './ui/dasar'

/** Pesan gagal dari STATUS jawaban - teks backend tidak ditampilkan apa adanya. */
export function pesanGagalLogin(galat: unknown): string {
  if (galat instanceof ApiFailure) {
    if (galat.status === 401) return LOGIN.salah
    if (galat.status === 423) return LOGIN.terkunci
    if (galat.status === 503) return LOGIN.tidakTersedia
  }
  return LOGIN.gagal
}

/** Latar bergerak - dihentikan saat layar dilepas, tidak dimulai bila pemakai meminta gerak minimal. */
export function useLatarPipeline() {
  const kanvasRef = useRef<HTMLCanvasElement | null>(null)
  useEffect(() => {
    const k = kanvasRef.current
    if (!k || window.matchMedia('(prefers-reduced-motion: reduce)').matches) return
    return mulaiPipeline(k)
  }, [])
  return kanvasRef
}

export default function Login({ onMasuk }: { onMasuk: (p: ProfilLogin) => void }) {
  const [akun, setAkun] = useState('')
  const [sandi, setSandi] = useState('')
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<string | null>(null)
  const kanvasRef = useLatarPipeline()

  return (
    <div className="login">
      <canvas ref={kanvasRef} className="login__latar" aria-hidden="true" role="presentation" />
      <form
        className="login__card"
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
        <div className="login__brand">
          {/* Gambar LATAR (pola `.shell__merek-tanda`): pagar unduhdokumen menolak
              atribut sumber gambar di layar. */}
          <span className="login__logo login__logo--latar" role="img" aria-label="Nusantara Re" />
        </div>
        <h1 className="login__judul">{LOGIN.judul}</h1>
        <p className="login__sub">{LOGIN.sub}</p>
        {galat !== null && (
          <div className="alert alert--error" role="alert">
            {galat}
          </div>
        )}
        <Field label={LOGIN.akun} value={akun} onChange={setAkun} required autoFocus />
        <FieldSandi label={LOGIN.sandi} value={sandi} onChange={setSandi} required />
        <button type="submit" className="btn btn--primary" disabled={sibuk || akun.trim() === '' || sandi === ''}>
          {sibuk ? LOGIN.memproses : LOGIN.masuk}
        </button>
        <p className="login__bantuan">{LOGIN.bantuan}</p>
      </form>
    </div>
  )
}
