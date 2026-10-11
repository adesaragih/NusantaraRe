// Halaman Register klaim Life — tiket 02.
//
// Alurnya: cari peserta per nomor premium list → pilih → (cari PL lain → pilih lagi) → kirim.
//
// ⛔ BANYAK PL SEKALIGUS (keputusan work owner 10-10-2026): pilihan BERTAHAN antar pencarian dan dikirim per
// premium list; server membuat satu case per PL number + mata uang, semua atau tidak sama sekali. Menyimpang dari
// Pega (satu polis per register) - PARITAS 10-10-2026.
//
// ⚠️ Uang tidak muncul di halaman ini sama sekali. Pendaftaran hanya memilih
// peserta; angka klaim lahir di tiket 03. Bila kelak ada, ia tetap TEKS
// (ADR-U-0003) — tidak pernah number JavaScript, yang membulatkan diam-diam.
import { useEffect, useRef, useState } from 'react'

import { REGISTER, REGISTER_BANYAK } from '../labels'
import { PanelDataPolis, tanggalTampil } from '../components/PanelDataPolis'
import KonfirmasiDaftar from '../components/KonfirmasiDaftar'
import {
  cariPesertaLife,
  daftarKlaimLife,
  ambilDataPolisIdPega,
  type CalonPeserta,
  type HasilDaftar,
  type PolicyDataLife,
} from '../api'
import { pesanGalat } from '../../../../inti/frontend/klien'
import {
  alihPilihan,
  kalimatJumlahCase,
  kelompokkanPilihan,
  kelompokPermintaan,
  kunciPilihan,
  sudahDipilih,
  bacaPolisPeserta,
  type PesertaTerpilih,
} from '../pilihanPeserta'
import { PitaGalat } from '../components/PitaGalat'

/**
 * kalimatGalat memilih kalimat yang ditampilkan pita merah.
 *
 * ⛔ Layar ini sempat punya PEMBACA GALATNYA SENDIRI, berbentuk axios
 * (`e.response.data.galat`) - dan axios sudah dibuang di F0.2. Bentuk itu
 * karena itu tidak pernah cocok dengan apa pun, sehingga SETIAP penolakan
 * backend tampil sebagai "Gagal menghubungi server."
 *
 * Itu bukan sekadar pesan yang hilang, melainkan pesan yang MENYESATKAN: ia
 * menuduh jaringan padahal backend menjawab, dan menjawab dengan sebab yang
 * tepat. Orang yang membacanya akan memeriksa koneksi, bukan datanya.
 *
 * Kini satu jalan saja: `pesanGalat` bersama di `inti/klien`.
 */
function kalimatGalat(e: unknown): string {
  return pesanGalat(e) ?? 'Permintaan ditolak backend.'
}

/** Case ID hasil Register, urut jawaban server (dikirim ke inbox lewat state rute, 10-10-2026). */
export function idTerdaftar(h: HasilDaftar): string[] {
  return h.kasus.map((k) => k.id)
}

export default function RegisterKlaim({
  onTerdaftar,
}: {
  /**
   * Register BERHASIL (keputusan work owner 10-10-2026 "langsung ke inbox claim life yang paling depan"): rute pindah ke
   * Inbox (tab depan) dan menampilkan pemberitahuan berisi Case ID-nya. Gagal = tetap di sini.
   */
  onTerdaftar?: (ids: string[]) => void
} = {}) {
  const [pl, setPl] = useState('')
  // Dua kotak `Find Insured`. Keduanya `pyReadOnly = false` dan
  // `pxTextInput`.
  //
  // ⚠️ Barisnya disebut BERSAMA TAG-nya, sebab satu label dapat berdiri di
  // lebih dari satu baris dengan tag berbeda, dan menyebut angkanya saja
  // membuat dua kutipan yang sah tampak bertentangan:
  //
  //   `Certificate No`   b16277 <pyLabelFieldValue>  -> CARI2
  //                      (tidak punya <pyLabelPreview> sama sekali)
  //   `Name of Insured`  b16039 <pyLabelFieldValue>  -> CARI3
  //                      b16064 <pyLabelPreview>     <- yang dikunci uji
  const [cariSertifikat, setCariSertifikat] = useState('')
  const [cariNama, setCariNama] = useState('')
  const [peserta, setPeserta] = useState<CalonPeserta[]>([])
  // ⛔ Pilihan BERTAHAN antar pencarian (10-10-2026): disimpan sebagai salinan barisnya, bukan indeks ke hasil
  // pencarian yang sedang tampil - mencari PL lain tidak menghapusnya.
  const [dipilih, setDipilih] = useState<PesertaTerpilih[]>([])
  const [galat, setGalat] = useState('')
  const [sibuk, setSibuk] = useState(false)
  // Modal konfirmasi Register Claim (10-10-2026) + penjaga klik ganda: satu konfirmasi = satu permintaan.
  const [konfirmasi, setKonfirmasi] = useState(false)
  const sedangKirim = useRef(false)

  // Panel data polis menampilkan polis dari peserta TERAKHIR yang dicentang (pilihan dapat berasal dari beberapa
  // premium list sejak 10-10-2026; panelnya menunjukkan yang sedang dikerjakan).
  const pertamaDipilih = dipilih.length > 0 ? dipilih[dipilih.length - 1] : undefined
  const kelompok = kelompokkanPilihan(dipilih)

  // ⛔ BUTIR av. Data polis DIBACA dari modul PremiumList Life begitu nomor
  // polisnya diketahui - di Pega ia terisi sendiri saat `Choose Policy No`,
  // dan medannya `pyReadOnly`. Membiarkan pemakai mengetiknya ulang membuat
  // dua salinan data polis yang akan berbeda pada hari pertama.
  //
  // ⚠️ 404 BUKAN KERUSAKAN: klaim dapat didaftarkan atas polis yang belum ada
  // di modul itu. Panelnya menyatakan keadaan itu alih-alih pecah.
  const [polis, setPolis] = useState<PolicyDataLife | null>(null)
  const [polisTidakDitemukan, setPolisTidakDitemukan] = useState<{ idPega: string } | undefined>(undefined)
  // ⛔ 10-10-2026: header polis = baris T_PREMIUM_LIST ber-IDPEGA peserta yang TERAKHIR dicentang (keputusan work
  // owner "ambil ke tabel T_PREMIUM_LIST berdasarkan idpega peserta nya saja") - terisi otomatis saat dicentang;
  // tanpa centang panel kosong.
  const adaTerpilih = pertamaDipilih !== undefined
  const idPegaTerpilih = pertamaDipilih?.idPega ?? ''
  useEffect(() => {
    if (!adaTerpilih) {
      setPolis(null)
      setPolisTidakDitemukan(undefined)
      return
    }
    let hidup = true
    void (async () => {
      const h = await bacaPolisPeserta(ambilDataPolisIdPega, idPegaTerpilih)
      if (!hidup) return
      setPolis(h.polis)
      setPolisTidakDitemukan(h.tidakDitemukan ? { idPega: h.idPega } : undefined)
    })()
    return () => {
      hidup = false
    }
  }, [adaTerpilih, idPegaTerpilih])

  async function cari() {
    setGalat('')
    setSibuk(true)
    try {
      // ⛔ Pilihan TIDAK dikosongkan: hasil baru hanya menandai baris yang sudah terpilih.
      setPeserta(
        await cariPesertaLife(pl, { sertifikat: cariSertifikat, nama: cariNama }),
      )
    } catch (e) {
      setPeserta([])
      setGalat(kalimatGalat(e))
    } finally {
      setSibuk(false)
    }
  }

  // ⛔ Kunci pilihan = PL + polis + sertifikat (pilihanPeserta.ts): nomor sertifikat hanya unik di dalam polisnya.
  function pilih(p: CalonPeserta) {
    setDipilih((lama) => alihPilihan(lama, p))
  }

  // ⛔ Hanya dipanggil dari tombol "Register Claim" di MODAL konfirmasi (10-10-2026), tidak dari layar.
  async function kirim() {
    if (sedangKirim.current) return
    sedangKirim.current = true
    setGalat('')
    setSibuk(true)
    try {
      // ⛔ Tidak ada lagi penolakan pilihan campuran polis / mata uang: setiap PL + mata uang menjadi case sendiri, dan
      // server yang memastikan polis tiap kelompok cocok dengan PL-nya.
      // Type dan kode bisnis dibaca server dari polis PL tiap case (10-10-2026).
      const h = await daftarKlaimLife({ kelompok: kelompokPermintaan(dipilih) })
      // Berhasil: pilihan dikosongkan dan rute pindah ke Inbox (tabel hasil di layar ini dihapus 10-10-2026).
      setDipilih([])
      onTerdaftar?.(idTerdaftar(h))
    } catch (e) {
      // Galat tampil sebagai alert di LAYAR (modal ditutup): kalimat server bisa panjang (beberapa PL/sertifikat),
      // dan pilihan tetap utuh untuk diperbaiki.
      setGalat(kalimatGalat(e))
    } finally {
      setSibuk(false)
      setKonfirmasi(false)
      sedangKirim.current = false
    }
  }

  return (
    <section className="register">
      <h2>{REGISTER_BANYAK.judulHalaman}</h2>
      {/* Pesan backend apa adanya, di ATAS isi (PitaGalat, 11-10-2026 "rapi seperti premium list"). */}
      <PitaGalat galat={galat} onTutup={() => setGalat('')} />

      {/* Kartu pencarian (tampilan 10-10-2026): tiga kotak + Search sejajar. ⛔ Kedua kotak penyaring BOLEH kosong =
          "jangan saring" (b16277, b16039). Huruf besar nama dikerjakan server (b405). */}
      <form
        className="register__kartu register__cari"
        onSubmit={(e) => {
          e.preventDefault()
          if (!sibuk && pl.trim() !== '') void cari()
        }}
      >
        <label>
          {REGISTER_BANYAK.labelPL}
          <input value={pl} onChange={(e) => setPl(e.target.value)} placeholder="RNML-..." />
        </label>
        <label>
          {REGISTER_BANYAK.kolomSertifikat}
          <input value={cariSertifikat} onChange={(e) => setCariSertifikat(e.target.value)} />
        </label>
        <label>
          {REGISTER_BANYAK.kolomNama}
          <input value={cariNama} onChange={(e) => setCariNama(e.target.value)} />
        </label>
        {/* Label VERBATIM `InputRegisterClaimLife.xml:16553` -> `LoadDataPesertaSpesifik_Act` b16576. */}
        <button type="submit" className="btn btn--primary" disabled={sibuk || pl.trim() === ''}>
          {REGISTER.cari}
        </button>
      </form>

      {peserta.length > 0 && (
        <div className="register__kartu register__gulir">
          <table className="inbox__tabel">
            <thead>
              <tr>
                <th>{REGISTER_BANYAK.kolomPilih}</th>
                <th>{REGISTER_BANYAK.kolomSertifikat}</th>
                <th>{REGISTER_BANYAK.kolomNama}</th>
                <th>{REGISTER_BANYAK.kolomTanggalLahir}</th>
              </tr>
            </thead>
            <tbody>
              {peserta.map((p) => {
                const aktif = sudahDipilih(dipilih, p)
                return (
                  // Baris dapat diklik untuk mencentang; baris terpilih tersorot.
                  <tr
                    key={kunciPilihan(p)}
                    className={aktif ? 'register__baris register__baris--pilih' : 'register__baris'}
                    onClick={() => pilih(p)}
                  >
                    <td>
                      <input
                        type="checkbox"
                        aria-label={`${REGISTER_BANYAK.kolomPilih} ${p.nomorSertifikat}`}
                        checked={aktif}
                        onClick={(e) => e.stopPropagation()}
                        onChange={() => pilih(p)}
                      />
                    </td>
                    {/* Nomor sertifikat TEKS: "006" bukan 6 (ADR-U-0022). */}
                    <td>{p.nomorSertifikat}</td>
                    <td>{p.namaTertanggung}</td>
                    {/* DOB dd/mm/yyyy menggantikan Policy No ceding (10-10-2026); nomor polis tetap di data. */}
                    <td>{tanggalTampil(p.tanggalLahir)}</td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      )}

      <PanelDataPolis
        nomorPolis={pertamaDipilih?.nomorPolis}
        namaTertanggung={pertamaDipilih?.namaTertanggung}
        mataUang={pertamaDipilih?.mataUang}
        polis={polis}
        polisTidakDitemukan={polisTidakDitemukan}
      />

      {/* ⛔ Isian Type dan kode bisnis DIHAPUS (keputusan work owner 10-10-2026): server membacanya dari polis PL. */}
      <DaftarTerpilih
        kelompok={kelompok}
        onHapus={(p) => setDipilih((lama) => alihPilihan(lama, p))}
        onKosongkan={() => setDipilih([])}
        onDaftar={() => setKonfirmasi(true)}
        sibuk={sibuk}
      />
      {konfirmasi && (
        <KonfirmasiDaftar
          kelompok={kelompok}
          sibuk={sibuk}
          onTutup={() => {
            if (!sibuk) setKonfirmasi(false)
          }}
          onLanjut={() => void kirim()}
        />
      )}

    </section>
  )
}

/**
 * Blok "Selected insured" (10-10-2026): pilihan dikelompokkan per PL number + mata uang - satu kelompok = satu case -
 * dengan tombol hapus per peserta dan kalimat jumlah case.
 */
export function DaftarTerpilih({
  kelompok,
  onHapus,
  onKosongkan,
  onDaftar,
  sibuk = false,
}: {
  kelompok: ReturnType<typeof kelompokkanPilihan>
  onHapus: (p: PesertaTerpilih) => void
  onKosongkan: () => void
  /** Tombol utama "Register Claim" (tampilan 10-10-2026: sejajar di kanan, bersama Clear selection). */
  onDaftar?: () => void
  sibuk?: boolean
}) {
  if (kelompok.length === 0) {
    return (
      <section className="register__kartu" aria-label={REGISTER_BANYAK.terpilih}>
        <h3>{REGISTER_BANYAK.terpilih}</h3>
        <p className="register__ringan">{REGISTER_BANYAK.kosongTerpilih}</p>
      </section>
    )
  }
  return (
    <section className="register__kartu" aria-label={REGISTER_BANYAK.terpilih}>
      <h3>{REGISTER_BANYAK.terpilih}</h3>
      {kelompok.map((k) => (
        <div key={`${k.nomorPremiList}|${k.mataUang}`} className="register__kelompok">
          <h4>
            {REGISTER_BANYAK.judulKelompokPL.split('{pl}').join(k.nomorPremiList)}
            <span className="register__chip">{k.mataUang}</span>
            <span className="register__jumlah">
              {REGISTER_BANYAK.jumlahPesertaKelompok.split('{n}').join(String(k.peserta.length))}
            </span>
          </h4>
          <div className="register__gulir">
            <table className="inbox__tabel">
              <thead>
                <tr>
                  <th>{REGISTER_BANYAK.kolomNo}</th>
                  <th>{REGISTER_BANYAK.kolomSertifikat}</th>
                  <th>{REGISTER_BANYAK.kolomNama}</th>
                  <th>{REGISTER_BANYAK.kolomTanggalLahir}</th>
                  <th aria-label={REGISTER_BANYAK.hapus} />
                </tr>
              </thead>
              <tbody>
                {k.peserta.map((p, i) => (
                  <tr key={kunciPilihan(p)}>
                    <td>{i + 1}</td>
                    <td>{p.nomorSertifikat}</td>
                    <td>{p.namaTertanggung}</td>
                    <td>{tanggalTampil(p.tanggalLahir)}</td>
                    <td>
                      <button type="button" className="btn btn--ghost btn--sm" onClick={() => onHapus(p)}>
                        {REGISTER_BANYAK.hapus}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      ))}
      <div className="register__aksi">
        <p className="register__ringkas">{kalimatJumlahCase(kelompok.length, REGISTER_BANYAK.jumlahCase)}</p>
        <button type="button" className="btn btn--ghost" onClick={onKosongkan} disabled={sibuk}>
          {REGISTER_BANYAK.kosongkan}
        </button>
        {onDaftar !== undefined && (
          <button type="button" className="btn btn--primary" onClick={onDaftar} disabled={sibuk}>
            {REGISTER_BANYAK.daftarkan}
          </button>
        )}
      </div>
    </section>
  )
}
