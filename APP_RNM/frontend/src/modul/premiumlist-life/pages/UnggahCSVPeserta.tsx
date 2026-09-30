// Layar unggah CSV peserta — tiket 04.
//
// Meniru `UploadCSV_LifePremium` → `UploadCSVLifePremium_Act` dan harness
// `ViewCSVResult_LifePremium*`, yang di Pega menampilkan hasil unggahan
// SEBELUM disimpan.
//
// ⛔ DUA LANGKAH, DAN LANGKAH PERTAMA TIDAK MENYIMPAN APA PUN. Pemakai
// meninjau dulu — baris mana lolos, baris mana ditolak, dan kolom apa yang
// salah — lalu memutuskan. Satu tombol yang "unggah dan simpan sekaligus"
// menghilangkan tepat yang AC tiket ini minta.
//
// ⛔ TOMBOL SIMPAN MATI SELAMA MASIH ADA PENOLAKAN. Tombol yang hidup lalu
// selalu dijawab 409 mengajari orang mengabaikan galat.
//
// ⛔ SETIAP PENOLAKAN MENYEBUT NOMOR BARIS DAN NAMA KOLOM. Penolakan tanpa
// keduanya memaksa orang mencocokkan pesan dengan berkas ratusan baris
// dengan mata.
//
// ⚠️ DUA KALIMAT PER PENOLAKAN, dan keduanya perlu: pesan VERBATIM yang
// dikenali pemakai lama, dan sebab yang tepat. `SUM INSURED HARUS ADA` yang
// muncul untuk nilai `1,000` berbohong tentang penyebabnya.

import { useRef, useState } from 'react'

import { UNGGAH_CSV } from '../labels'
import { Gagal } from '../../../inti/components/ui/dasar'
import {
  simpanUnggahPolis,
  tinjauUnggahPolis,
  type HasilTinjauUnggah,
  type PenolakanUnggah,
} from '../api'

/** Ringkasan satu tinjauan, dalam satu kalimat. */
export function ringkasanTinjau(h: HasilTinjauUnggah): string {
  if (h.lolos) {
    return `${String(h.cacahBaris)} baris terbaca, semuanya lolos.`
  }
  const baris = new Set(h.ditolak.map((p) => p.baris)).size
  return (
    `${String(h.cacahBaris)} baris terbaca; ${String(baris)} baris ditolak ` +
    `dengan ${String(h.cacahDitolak)} alasan. Tidak ada yang tersimpan.`
  )
}

/** Kunci baris tabel penolakan — satu baris dapat punya banyak alasan. */
export function kunciPenolakan(p: PenolakanUnggah, urut: number): string {
  return `${String(p.baris)}-${p.kolom}-${String(urut)}`
}

export default function UnggahCSVPeserta({
  polisID,
  onTersimpan,
}: {
  polisID: string
  /** Dipanggil sesudah penyimpanan berhasil, supaya grid disegarkan. */
  onTersimpan: () => void
}) {
  const [berkas, setBerkas] = useState<File | null>(null)
  const [tinjau, setTinjau] = useState<HasilTinjauUnggah | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [kabar, setKabar] = useState('')
  const pilih = useRef<HTMLInputElement>(null)

  function pilihBerkas(f: File | null): void {
    setBerkas(f)
    // ⛔ Tinjauan LAMA dibuang begitu berkasnya berganti. Membiarkannya
    // membuat layar menampilkan hasil berkas yang sudah tidak dipilih —
    // dan tombol simpan akan merujuk berkas yang berbeda dari yang terlihat.
    setTinjau(null)
    setKabar('')
    setGalat(null)
  }

  async function jalankanTinjau(): Promise<void> {
    if (berkas === null || sibuk) return
    setSibuk(true)
    setGalat(null)
    setKabar('')
    try {
      setTinjau(await tinjauUnggahPolis(polisID, berkas))
    } catch (e) {
      setGalat(e)
      setTinjau(null)
    } finally {
      setSibuk(false)
    }
  }

  async function jalankanSimpan(): Promise<void> {
    if (berkas === null || sibuk || tinjau === null || !tinjau.lolos) return
    setSibuk(true)
    setGalat(null)
    try {
      const h = await simpanUnggahPolis(polisID, berkas)
      setKabar(
        `${String(h.cacahDisimpan)} peserta tersimpan` +
          (h.cacahDihapus > 0
            ? `, menggantikan ${String(h.cacahDihapus)} baris sebelumnya.`
            : '.'),
      )
      setTinjau(null)
      setBerkas(null)
      if (pilih.current !== null) pilih.current.value = ''
      onTersimpan()
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  const bolehSimpan = tinjau !== null && tinjau.lolos && !sibuk

  return (
    <section className="unggah-csv">
      <h3 className="unggah-csv__judul">{UNGGAH_CSV.judul}</h3>

      <p className="unggah-csv__aturan" role="note">
        {UNGGAH_CSV.aturanPemisah}
      </p>

      <p className="unggah-csv__pilih">
        <input
          ref={pilih}
          type="file"
          accept=".csv,text/csv"
          aria-label={UNGGAH_CSV.pilihBerkas}
          onChange={(e) => {
            pilihBerkas(e.target.files?.[0] ?? null)
          }}
        />{' '}
        <button
          type="button"
          disabled={berkas === null || sibuk}
          onClick={() => {
            void jalankanTinjau()
          }}
        >
          {UNGGAH_CSV.tinjau}
        </button>{' '}
        <button
          type="button"
          className="unggah-csv__simpan"
          disabled={!bolehSimpan}
          onClick={() => {
            void jalankanSimpan()
          }}
        >
          {UNGGAH_CSV.simpan}
        </button>
      </p>

      {galat !== null && <Gagal galat={galat} />}
      {kabar !== '' && <p role="status">{kabar}</p>}

      {tinjau !== null && (
        <>
          <p className="unggah-csv__ringkas" role="status">
            {ringkasanTinjau(tinjau)}
          </p>
          {/* ⛔ Sebab tombol simpan mati DIKATAKAN, bukan dibiarkan ditebak. */}
          {!tinjau.lolos && <p role="note">{UNGGAH_CSV.perbaikiDulu}</p>}
          {tinjau.ditolak.length > 0 && (
            <table className="unggah-csv__tabel">
              <thead>
                <tr>
                  <th>{UNGGAH_CSV.kolomBaris}</th>
                  <th>{UNGGAH_CSV.kolomKolom}</th>
                  <th>{UNGGAH_CSV.kolomPesan}</th>
                  <th>{UNGGAH_CSV.kolomSebab}</th>
                </tr>
              </thead>
              <tbody>
                {tinjau.ditolak.map((p, i) => (
                  <tr key={kunciPenolakan(p, i)}>
                    {/* Baris 0 berarti penolakan atas BERKASNYA, bukan baris. */}
                    <td>{p.baris === 0 ? '—' : p.baris}</td>
                    <td>{p.kolom === '' ? '—' : p.kolom}</td>
                    <td>{p.pesan}</td>
                    <td>{p.sebab}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </>
      )}
    </section>
  )
}
