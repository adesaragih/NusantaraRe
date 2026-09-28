// Pencarian diagnosa — A3 kelompok Medis.
//
// Meniru `Harness/Diagnose_Harness.xml` → `Section/Diagnose_Section.xml`,
// dibuka tombol `Find Disease` (`ClaimLifeDetailGCNM.xml` b5061 →
// `showHarness` b5081).
//
// ⛔ DUA kotak pencarian, bukan satu, dan itu bukan pilihan tata letak:
//
//   `Diagnose_Section.xml` b519 `CARI1` → b1645 `Param.ICD_Code`
//   `Diagnose_Section.xml` b801 `CARI2` → b1651 `Param.Disease`
//
// dan `BrowseDiseaseLife_RD.xml` b535 menyambungnya **`A AND B`**.
// Menggabungkannya menjadi satu kotak mengubah maknanya — satu kotak hanya
// dapat berarti OR atau "cari di mana saja", dan bedanya baru terlihat pada
// dua kata kunci sekaligus.
//
// ⛔ Keduanya dinaikkan ke HURUF BESAR di backend (`SearchDiagnose_act.xml`
// b255, b302). Tidak dilakukan di sini: satu aturan, satu tempat.
//
// ⛔ BATASNYA DARI RULE. `DISEASE_LIFE` berisi 97.586 baris; `pyMaxRecords`
// b659 menyebut 500 dan `pyPageSize` b514 menyebut 50. Layar ini meminta 50
// dan backend menjepitnya lagi — layar bukan penjaga.
//
// ✅ Tombol `Choose` (b2509 → `SetDisease` b2528) KINI TERPASANG — butir
// **bd**, 27-09-2026. Kepala berkas ini dulu berbunyi *"belum terpasang, dan
// sebabnya bukan kemalasan: ... satu lawan banyak, OQ-K"*. Premisnya keliru:
// yang dibandingkan adalah RepeatGrid b3923 dengan kolom TUNGGAL
// `T_CLAIMLF_PREMIUMLIST_DETAIL.DISEASE` warisan migrasi 003, seolah tabel
// itu satu-satunya tempat yang mungkin. Jawabannya ada dua baris di bawah
// tempat pembacaan itu berhenti: `Add` b4690 → `addRow` b4700 dan `Delete`
// b6160 → `deleteRow` b6170. Grid **dapat** berarti tampilan satu baris;
// grid ber-`Add` DAN ber-`Delete` **tidak dapat**. Tabelnya kini ada
// (`T_CLAIMLF_DIAGNOSE`, migrasi 018), dan `Choose` menulis ke BARIS yang
// memanggilnya.
//
// ⛔ Karena itu pencarian ini berdiri DI DALAM baris grid, bukan di
// sampingnya: `Find Disease` b5061 adalah sel 37 pada baris data (b5005),
// sehingga barisnya sendirilah konteks `Choose`. Tanpa itu, "diagnosa yang
// mana" harus ditebak.

import { useState } from 'react'

import { DETAIL } from '../../assets/labels.claimlife'
import {
  cariPenyakit,
  pesanGalat,
  UKURAN_HALAMAN_PENYAKIT,
  type Penyakit,
} from '../../services/api'

/**
 * Label VERBATIM dari korpus.
 *
 * ⚠️ `buka` MENGAMBIL dari DETAIL, tidak menuliskannya lagi. Label yang
 * sama di dua tempat adalah dua tempat untuk bergeser, dan yang bergeser
 * tidak akan berbunyi: keduanya benar menurut dirinya sendiri. Itu bentuk
 * cacat yang sudah empat kali terjadi di modul ini.
 */
export const LABEL_CARI_DIAGNOSA = {
  /** `ClaimLifeDetailGCNM.xml` b5061 `pyLabel` — satu sumber: DETAIL. */
  buka: DETAIL.cariPenyakit,
  /** `Diagnose_Section.xml` b2509 `pyLabel`. */
  pilih: 'Choose',
} as const

/**
 * Menyusun kalimat ringkas tentang hasil pencarian.
 *
 * ⚠️ Dipisah supaya dapat diuji tanpa DOM, dan supaya satu hal terjaga:
 * ketika hasilnya menyentuh batas, pemakai HARUS diberi tahu. Daftar yang
 * terpotong diam-diam terbaca sebagai daftar yang lengkap, dan orang akan
 * menyimpulkan diagnosanya tidak ada.
 */
export function ringkasanHasil(jumlah: number, batas: number): string {
  if (jumlah === 0) return 'Tidak ada diagnosa yang cocok.'
  if (jumlah >= batas) {
    return `${jumlah} diagnosa ditampilkan — daftarnya TERPOTONG pada batas ${batas}. Persempit kata kuncinya.`
  }
  return `${jumlah} diagnosa ditemukan.`
}

export function CariDiagnosa({
  onPilih,
  sibuk: sibukLuar = false,
  batas = UKURAN_HALAMAN_PENYAKIT,
}: {
  /**
   * Dipanggil ketika `Choose` ditekan — padanan `SetDisease` b2528.
   *
   * ⚠️ WAJIB, bukan opsional. Pencarian tanpa tempat menaruh hasilnya
   * adalah layar yang menyibukkan orang tanpa mengubah apa pun — dan
   * "prop opsional yang selalu diisi" hanyalah cabang mati yang menunggu
   * seseorang lupa mengisinya.
   */
  onPilih: (p: Penyakit) => void | Promise<void>
  /** Ada perubahan lain yang sedang berjalan pada baris ini. */
  sibuk?: boolean
  batas?: number
}) {
  // ⛔ TERTUTUP sampai ditekan, dan itu bukan kosmetik. Tombolnya berdiri
  // PER PESERTA - `Find Disease` b5061 ada di `ClaimLifeDetailGCNM`, section
  // berkelas `Int-LIFE_PREMIUM_DETAIL` - sehingga klaim grup berpeserta 500
  // akan merender 500 formulir pencarian sekaligus bila selalu terbuka,
  // masing-masing dengan keadaannya sendiri. Di Pega pun b5061 adalah TOMBOL
  // yang MEMBUKA harness (`showHarness` b5081), bukan formulir yang selalu
  // tampak. Jadi menutupnya justru lebih setia, bukan kurang.
  const [terbuka, setTerbuka] = useState(false)
  const [kodeIcd, setKodeIcd] = useState('')
  const [nama, setNama] = useState('')
  const [hasil, setHasil] = useState<Penyakit[] | null>(null)
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<string | null>(null)

  async function cari(): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    try {
      setHasil(await cariPenyakit(kodeIcd, nama, batas))
    } catch (e) {
      setHasil(null)
      setGalat(pesanGalat(e) ?? 'Pencarian diagnosa gagal.')
    } finally {
      setSibuk(false)
    }
  }

  if (!terbuka) {
    return (
      <p className="diagnosa">
        <button type="button" onClick={() => setTerbuka(true)}>
          {LABEL_CARI_DIAGNOSA.buka}
        </button>
      </p>
    )
  }

  return (
    <section className="diagnosa">
      <h4 className="diagnosa__judul">{LABEL_CARI_DIAGNOSA.buka}</h4>

      <p className="diagnosa__kotak">
        <label>
          ICD Code{' '}
          <input
            type="text"
            value={kodeIcd}
            onChange={(e) => setKodeIcd(e.target.value)}
          />
        </label>{' '}
        <label>
          Disease{' '}
          <input type="text" value={nama} onChange={(e) => setNama(e.target.value)} />
        </label>{' '}
        <button type="button" disabled={sibuk} onClick={() => void cari()}>
          {sibuk ? 'Mencari…' : 'Cari'}
        </button>{' '}
        <button
          type="button"
          disabled={sibuk}
          onClick={() => {
            setTerbuka(false)
          }}
        >
          Tutup
        </button>
      </p>

      {galat !== null && <p role="alert">{galat}</p>}

      {hasil !== null && (
        <>
          <p role="status">{ringkasanHasil(hasil.length, batas)}</p>
          {hasil.length > 0 && (
            <table className="diagnosa__tabel">
              <thead>
                <tr>
                  <th scope="col">ID</th>
                  <th scope="col">Disease</th>
                  <th scope="col">ICD Code</th>
                  <th scope="col">{LABEL_CARI_DIAGNOSA.pilih}</th>
                </tr>
              </thead>
              <tbody>
                {hasil.map((p) => (
                  <tr key={`${p.nomor}:${p.kodeIcd}`}>
                    <td>{p.nomor}</td>
                    <td>{p.nama}</td>
                    <td>{p.kodeIcd}</td>
                    <td>
                      <button
                        type="button"
                        disabled={sibuk || sibukLuar}
                        onClick={() => {
                          void (async () => {
                            await onPilih(p)
                            // Harness Pega tertutup sesudah `Choose`
                            // (`SetDisease` b389 menyimpan lalu kembali).
                            setTerbuka(false)
                          })()
                        }}
                      >
                        {LABEL_CARI_DIAGNOSA.pilih}
                      </button>
                    </td>
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
