// Halaman Register klaim Life — tiket 02.
//
// Alurnya satu arah: cari peserta per nomor premium list → pilih → kirim.
//
// ⚠️ Uang tidak muncul di halaman ini sama sekali. Pendaftaran hanya memilih
// peserta; angka klaim lahir di tiket 03. Bila kelak ada, ia tetap TEKS
// (ADR-U-0003) — tidak pernah number JavaScript, yang membulatkan diam-diam.
import { useState } from 'react'

import { REGISTER } from '../assets/labels'
import { PanelDataPolis } from '../components/PanelDataPolis'
import {
  cariPesertaLife,
  daftarKlaimLife,
  type CalonPeserta,
  type HasilDaftar,
} from '../services/api'

/** pesanGalat mengambil pesan dari backend apa adanya bila ada. */
function pesanGalat(e: unknown): string {
  const jawaban = (e as { response?: { data?: { galat?: string } } })?.response
  return jawaban?.data?.galat ?? 'Gagal menghubungi server.'
}

export default function RegisterKlaim() {
  const [pl, setPl] = useState('')
  // Dua kotak `Find Insured`. Keduanya `pyReadOnly = false` dan
  // `pxTextInput`: b16277 `Certificate No` -> CARI2, b16039
  // `Name of Insured` -> CARI3.
  const [cariSertifikat, setCariSertifikat] = useState('')
  const [cariNama, setCariNama] = useState('')
  const [type, setType] = useState('QP')
  const [kodeBisnis, setKodeBisnis] = useState('')
  const [peserta, setPeserta] = useState<CalonPeserta[]>([])
  const [dipilih, setDipilih] = useState<string[]>([])
  const [hasil, setHasil] = useState<HasilDaftar | null>(null)
  const [galat, setGalat] = useState('')
  const [sibuk, setSibuk] = useState(false)

  // Panel data polis menampilkan polis dari peserta PERTAMA yang
  // dicentang: di Pega satu klaim menunjuk satu polis, dan peserta
  // yang dipilih berasal dari premium list yang sama.
  const pertamaDipilih = peserta.find((p) => dipilih.includes(kunci(p)))

  async function cari() {
    setGalat('')
    setHasil(null)
    setSibuk(true)
    try {
      setPeserta(
        await cariPesertaLife(pl, { sertifikat: cariSertifikat, nama: cariNama }),
      )
      setDipilih([])
    } catch (e) {
      setPeserta([])
      setGalat(pesanGalat(e))
    } finally {
      setSibuk(false)
    }
  }

  // ⛔ Kunci pilihan memakai polis DAN sertifikat, bukan sertifikat saja.
  // Satu premium list dapat memuat beberapa polis, dan nomor sertifikat hanya
  // unik di dalam polisnya - dua peserta bersertifikat "006" dari polis
  // berbeda akan tercentang berbarengan bila kuncinya sertifikat saja.
  function kunci(p: CalonPeserta): string {
    return `${p.nomorPolis}|${p.nomorSertifikat}`
  }

  function pilih(k: string) {
    setDipilih((lama) => (lama.includes(k) ? lama.filter((s) => s !== k) : [...lama, k]))
  }

  async function kirim() {
    setGalat('')
    setSibuk(true)
    try {
      const terpilih = peserta.filter((p) => dipilih.includes(kunci(p)))
      // ⚠️ Polis dan mata uang diambil dari peserta terpilih PERTAMA. Itu
      // hanya benar bila seluruh pilihan sepolis dan semata-uang; layar
      // menolak campuran lebih dulu, sebab mendaftarkan klaim bermata uang
      // campur akan menyimpan satu mata uang untuk semuanya tanpa ada yang
      // menyadarinya.
      const polis = new Set(terpilih.map((p) => p.nomorPolis))
      const uang = new Set(terpilih.map((p) => p.mataUang))
      if (polis.size > 1 || uang.size > 1) {
        setGalat('Peserta terpilih berbeda polis atau mata uang. Pilih yang sepolis dan semata-uang.')
        return
      }
      setHasil(
        await daftarKlaimLife({
          nomorPremiList: pl,
          nomorPolis: terpilih[0]?.nomorPolis ?? '',
          type,
          kodeBisnis,
          mataUang: terpilih[0]?.mataUang ?? '',
          sertifikat: terpilih.map((p) => p.nomorSertifikat),
        }),
      )
    } catch (e) {
      setGalat(pesanGalat(e))
    } finally {
      setSibuk(false)
    }
  }

  return (
    <section>
      <h2>Register Klaim Life</h2>

      <label>
        Nomor premium list
        <input value={pl} onChange={(e) => setPl(e.target.value)} placeholder="PL-..." />
      </label>
      {/* ⛔ Kedua kotak penyaring ini BOLEH kosong, dan kosong berarti
          "jangan saring". Label VERBATIM dari sectionnya: b16277 dan
          b16039; keduanya `pyReadOnly = false`. */}
      <label>
        {REGISTER.sertifikat}
        <input
          value={cariSertifikat}
          onChange={(e) => setCariSertifikat(e.target.value)}
        />
      </label>
      <label>
        {REGISTER.namaTertanggung}
        {/* Huruf besarnya dikerjakan server (b405 `@toUpperCase`), bukan di
            sini: mengubah ketikan orang saat ia mengetik membuat kotaknya
            terasa rusak, dan aturan pencarian tetap satu tempat. */}
        <input value={cariNama} onChange={(e) => setCariNama(e.target.value)} />
      </label>

      {/* Label VERBATIM `InputRegisterClaimLife.xml:16553` `pxButton` ->
          `LoadDataPesertaSpesifik_Act` b16576. Tombol `Find Insured` b7057
          adalah pembuka panelnya, bukan penjalan pencarian. */}
      <button onClick={cari} disabled={sibuk || pl.trim() === ''}>
        {REGISTER.cari}
      </button>

      {peserta.length > 0 && (
        <table>
          <thead>
            <tr>
              {/* Label VERBATIM: Select Insured 20008, Name of Insured
                  16064, Choose Policy No 3776. */}
              <th>{REGISTER.pilihTertanggung}</th>
              <th>Sertifikat</th>
              <th>{REGISTER.namaTertanggung}</th>
              <th>{REGISTER.pilihPolis}</th>
            </tr>
          </thead>
          <tbody>
            {peserta.map((p) => (
              <tr key={kunci(p)}>
                <td>
                  <input
                    type="checkbox"
                    checked={dipilih.includes(kunci(p))}
                    onChange={() => pilih(kunci(p))}
                  />
                </td>
                {/* Nomor sertifikat TEKS: "006" bukan 6 (ADR-U-0022). */}
                <td>{p.nomorSertifikat}</td>
                <td>{p.namaTertanggung}</td>
                <td>{p.nomorPolis}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* ⛔ Himpunan medan PERSIS section: sebelas medan
          `.PolicyDataLife.*`. Yang belum punya sumber di backend
          DINYATAKAN di dalam panel, bukan dihilangkan - paritas yang
          tampak lengkap padahal tidak adalah paritas yang tidak akan
          pernah dicari lagi. */}
      <PanelDataPolis
        nomorPolis={pertamaDipilih?.nomorPolis}
        namaTertanggung={pertamaDipilih?.namaTertanggung}
        mataUang={pertamaDipilih?.mataUang}
        type={type}
      />

      {/* Label VERBATIM `InputRegisterClaimLife.xml:9104`
          `.PolicyDataLife.Type` `pxDropdown`. */}
      <label>
        {REGISTER.type}
        <input value={type} onChange={(e) => setType(e.target.value)} />
      </label>
      <label>
        Kode bisnis
        <input value={kodeBisnis} onChange={(e) => setKodeBisnis(e.target.value)} />
      </label>
      <button onClick={kirim} disabled={sibuk || dipilih.length === 0}>
        Daftarkan klaim
      </button>

      {hasil && (
        <p>
          Klaim terdaftar: <strong>{hasil.id}</strong>
          {hasil.nomorKlaim !== '' && <> — nomor {hasil.nomorKlaim}</>}
        </p>
      )}
      {/* Pesan backend diteruskan apa adanya: bila nomor klaim belum dapat
          dibentuk, pesannya menyebut keputusan mana yang ditunggu. */}
      {galat !== '' && <p role="alert">{galat}</p>}
    </section>
  )
}
