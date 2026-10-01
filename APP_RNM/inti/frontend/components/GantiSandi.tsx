// Ganti sandi - M_LOGIN_GO (keputusan work owner 01-10-2026).
//
// Dua jalan masuk: WAJIB (akun baru atau sandi direset - layar penuh, tanpa
// Batal, sebelum aplikasi dapat dipakai) dan PILIHAN (dari menu profil).
// Backend menegakkan aturannya; di sini hanya pemeriksaan awal supaya pemakai
// tidak menunggu jawaban server untuk salah yang jelas. Tampilan: kerangka
// halaman login (`KerangkaMasuk`).

import { useState } from 'react'

import { ApiFailure, gantiSandiLogin, type ProfilLogin } from '../klien'
import { LOGIN, PANJANG_MIN_SANDI } from '../labels'
import { IsianSandi, KerangkaMasuk } from './Login'

/** Pemeriksaan awal; `null` = boleh dikirim. Panjang dihitung KARAKTER, seperti backend. */
export function periksaSandiBaru(baru: string, ulang: string): string | null {
  if ([...baru].length < PANJANG_MIN_SANDI) return LOGIN.minimal
  if (baru !== ulang) return LOGIN.tidakSama
  return null
}

/**
 * Galat ganti sandi: 400 membawa kalimat backend (`login/rute.go`) yang
 * DIPETAKAN ke label berbahasa Inggris - layar login berbahasa Inggris
 * (permintaan work owner 01-10-2026); kalimat tak dikenal = pesan umum.
 * 401 = sesi berakhir.
 */
export function pesanGagalGanti(galat: unknown): string {
  if (galat instanceof ApiFailure && galat.status === 400) {
    const pesan = galat.detail.message ?? ''
    if (pesan.includes('sandi lama salah')) return LOGIN.lamaSalah
    if (pesan.includes('sandi minimal')) return LOGIN.minimal
    if (pesan.includes('sandi maksimal')) return LOGIN.terlaluPanjang
    if (pesan.includes('sama dengan sandi lama')) return LOGIN.samaDenganLama
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

  const kosongkan = () => {
    setLama('')
    setBaru('')
    setUlang('')
  }

  return (
    <KerangkaMasuk judul={LOGIN.judulGanti} sub={wajib ? LOGIN.wajibGanti : LOGIN.minimal}>
      <form
        className="halaman-masuk__form"
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
        {galat !== null && (
          <div className="alert alert--error" role="alert">
            {galat}
          </div>
        )}
        <IsianSandi
          id="ganti-sandi-lama"
          label={LOGIN.sandiLama}
          placeholder={LOGIN.sandiLama}
          value={lama}
          onChange={setLama}
          autoComplete="current-password"
          autoFocus
        />
        <IsianSandi
          id="ganti-sandi-baru"
          label={LOGIN.sandiBaru}
          placeholder={LOGIN.sandiBaru}
          value={baru}
          onChange={setBaru}
          autoComplete="new-password"
        />
        <IsianSandi
          id="ganti-sandi-ulang"
          label={LOGIN.ulangiSandi}
          placeholder={LOGIN.ulangiSandi}
          value={ulang}
          onChange={setUlang}
          autoComplete="new-password"
        />
        <button className="halaman-masuk__tombol" type="submit" disabled={sibuk || lama === '' || baru === ''}>
          {sibuk ? LOGIN.memproses : LOGIN.simpanSandi}
        </button>
        {!wajib && onBatal && (
          <button type="button" className="halaman-masuk__tombol-kedua" onClick={onBatal} disabled={sibuk}>
            {LOGIN.batal}
          </button>
        )}
      </form>
    </KerangkaMasuk>
  )
}
