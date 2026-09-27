// Layar masuk — STUB, F0.2 (brief lanjutan 6 §2 aturan 5).
//
// ⛔ INI BUKAN AUTENTIKASI. Tidak ada sandi yang diminta, dikirim, maupun
// disimpan; yang terjadi hanyalah pemakai MENGAKU sebagai sebuah akun dan
// memilih peran, lalu pengakuan itu dikirim sebagai header `X-Pelaku` /
// `X-Peran`. Backend memercayainya hanya ketika stub-nya menyala
// (`handlers.pelakuDari(r, stubAktif)`); tanpa itu ia mengembalikan pelaku
// KOSONG dan setiap jalur beridentitas menolak.
//
// ⛔ Gerbangnya `VITE_AUTH_STUB === 'true'`, dan ia GAGAL TERTUTUP: env yang
// tidak disetel berarti layar ini menolak memberi jalan masuk. IAM sungguhan
// adalah tiket 07 / ADR-U-0030.
//
// ⚠️ Pita peringatan di layar bukan hiasan: satu-satunya hal yang membedakan
// layar ini dari layar masuk sungguhan adalah pengetahuan pemakainya.

import { useState, type FormEvent } from 'react'

import { PERAN, PERAN_ID, PRODUK } from '../assets/labels'
import type { KodePeran } from '../assets/labels'
import { Field, IkonMasuk } from '../components/ui/dasar'
import { bolehMasukStub, sesi, PERAN_TERSEDIA, type Sesi } from '../store/sesi'

/** Awalan akun yang disarankan - menegaskan bahwa ini akun UJI. */
const AWALAN_UJI = 'UJI-'

export interface MasukProps {
  /** Dipanggil sesudah sesi tersimpan. */
  onMasuk: (s: Sesi) => void
}

export default function Masuk({ onMasuk }: MasukProps) {
  const [akun, setAkun] = useState(`${AWALAN_UJI}ADMIN`)
  const [peran, setPeran] = useState<KodePeran[]>([PERAN.admin])
  const [galat, setGalat] = useState<string | null>(null)

  const stubMenyala = bolehMasukStub()

  function ubahPeran(p: KodePeran, dipilih: boolean): void {
    setPeran((lama) => (dipilih ? [...lama, p] : lama.filter((x) => x !== p)))
  }

  function kirim(e: FormEvent): void {
    e.preventDefault()
    const bersih = akun.trim()
    if (bersih === '') {
      setGalat('Nama akun wajib diisi.')
      return
    }
    // ⛔ Nol peran ditolak DI SINI juga, bukan hanya di backend: sesi tanpa
    // peran akan membuat setiap permintaan ditolak 403 tanpa pemakai tahu
    // sebabnya.
    if (peran.length === 0) {
      setGalat('Pilih sekurangnya satu peran.')
      return
    }
    const s: Sesi = { akunID: bersih, peran }
    if (!sesi.simpan(s)) {
      setGalat(
        'Peramban menolak menyimpan sesi (mode privat atau data situs ' +
          'diblokir). Masuk tidak dapat dilanjutkan.',
      )
      return
    }
    setGalat(null)
    onMasuk(s)
  }

  if (!stubMenyala) {
    return (
      <main className="masuk">
        <div className="masuk__kartu">
          <h1 className="masuk__judul">{PRODUK.nama}</h1>
          <p className="masuk__galat" role="alert">
            Masuk lewat stub dimatikan (<code>VITE_AUTH_STUB</code> bukan{' '}
            <code>true</code>). Belum ada sumber identitas lain, jadi aplikasi
            tidak dapat dimasuki dari layar ini.
          </p>
        </div>
      </main>
    )
  }

  return (
    <main className="masuk">
      <form className="masuk__kartu" onSubmit={kirim}>
        <h1 className="masuk__judul">{PRODUK.nama}</h1>
        <p className="masuk__sub">{PRODUK.sub}</p>

        {/* ⛔ Pita ini TIDAK boleh dihapus selama stub menyala. */}
        <p className="masuk__pita" role="status">
          <strong>Mode stub.</strong> Identitas tidak diverifikasi dan tidak
          ada sandi. Akun serta peran yang dipilih di sini dikirim apa adanya
          sebagai header permintaan.
        </p>

        <Field
          label="Akun"
          value={akun}
          onChange={setAkun}
          placeholder={`${AWALAN_UJI}ADMIN`}
        />

        <fieldset className="masuk__peran">
          <legend>Peran</legend>
          {/* ⚠️ Kotak centang, bukan pilihan tunggal: `pelakuDari` memecah
              `X-Peran` pada koma, jadi satu pelaku memang boleh memegang
              lebih dari satu peran. */}
          {PERAN_TERSEDIA.map((p) => (
            <label key={p} className="masuk__peran-baris">
              <input
                type="checkbox"
                checked={peran.includes(p)}
                onChange={(e) => {
                  ubahPeran(p, e.target.checked)
                }}
              />
              <span>{PERAN_ID[p]}</span>
              <code className="masuk__peran-kode">{p}</code>
            </label>
          ))}
        </fieldset>

        {galat !== null && (
          <p className="masuk__galat" role="alert">
            {galat}
          </p>
        )}

        <button type="submit" className="masuk__tombol">
          <IkonMasuk />
          Masuk
        </button>
      </form>
    </main>
  )
}
