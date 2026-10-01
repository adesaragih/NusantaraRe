// Ganti sandi - M_LOGIN_GO (keputusan work owner 01-10-2026).
//
// Dua jalan masuk: WAJIB (akun baru atau sandi direset - layar penuh, tanpa
// Batal, sebelum aplikasi dapat dipakai) dan PILIHAN (dari menu profil).
// Backend menegakkan aturannya; di sini hanya pemeriksaan awal supaya pemakai
// tidak menunggu jawaban server untuk salah yang jelas.

import { useState } from 'react'

import { ApiFailure, gantiSandiLogin, type ProfilLogin } from '../klien'
import { LOGIN, PANJANG_MIN_SANDI } from '../labels'
import { useLatarPipeline } from './Login'
import { FieldSandi } from './ui/dasar'

/** Pemeriksaan awal; `null` = boleh dikirim. Panjang dihitung KARAKTER, seperti backend. */
export function periksaSandiBaru(baru: string, ulang: string): string | null {
  if ([...baru].length < PANJANG_MIN_SANDI) return LOGIN.minimal
  if (baru !== ulang) return LOGIN.tidakSama
  return null
}

/** Galat ganti sandi: 400 membawa kalimat backend; 401 = sesi berakhir. */
export function pesanGagalGanti(galat: unknown): string {
  if (galat instanceof ApiFailure && galat.status === 400) {
    return galat.detail.message?.includes('sandi lama salah') ? LOGIN.lamaSalah : (galat.detail.message ?? LOGIN.gagal)
  }
  return LOGIN.gagal
}

export default function GantiSandi({
  wajib,
  onSelesai,
  onBatal,
}: {
  /** Akun wajib ganti sandi: layar penuh, tanpa Batal. */
  wajib: boolean
  onSelesai: (p: ProfilLogin) => void
  onBatal?: () => void
}) {
  const [lama, setLama] = useState('')
  const [baru, setBaru] = useState('')
  const [ulang, setUlang] = useState('')
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<string | null>(null)
  const kanvasRef = useLatarPipeline()

  const kosongkan = () => {
    setLama('')
    setBaru('')
    setUlang('')
  }

  return (
    <div className="login">
      <canvas ref={kanvasRef} className="login__latar" aria-hidden="true" role="presentation" />
      <form
        className="login__card"
        onSubmit={(e) => {
          e.preventDefault()
          if (sibuk) return
          const awal = periksaSandiBaru(baru, ulang)
          if (awal !== null) {
            setGalat(awal)
            return
          }
          setSibuk(true)
          setGalat(null)
          gantiSandiLogin(lama, baru).then(
            (p) => {
              kosongkan()
              setSibuk(false)
              onSelesai(p)
            },
            (g: unknown) => {
              kosongkan()
              setSibuk(false)
              setGalat(pesanGagalGanti(g))
            },
          )
        }}
      >
        <h1 className="login__judul">{LOGIN.judulGanti}</h1>
        {wajib && <p className="login__sub">{LOGIN.wajibGanti}</p>}
        {galat !== null && (
          <div className="alert alert--error" role="alert">
            {galat}
          </div>
        )}
        <FieldSandi label={LOGIN.sandiLama} value={lama} onChange={setLama} required autoFocus />
        <FieldSandi label={LOGIN.sandiBaru} value={baru} onChange={setBaru} required />
        <FieldSandi label={LOGIN.ulangiSandi} value={ulang} onChange={setUlang} required />
        <p className="login__bantuan">{LOGIN.minimal}</p>
        <button type="submit" className="btn btn--primary" disabled={sibuk || lama === '' || baru === ''}>
          {sibuk ? LOGIN.memproses : LOGIN.simpanSandi}
        </button>
        {!wajib && onBatal && (
          <button type="button" className="btn" onClick={onBatal} disabled={sibuk}>
            {LOGIN.batal}
          </button>
        )}
      </form>
    </div>
  )
}
