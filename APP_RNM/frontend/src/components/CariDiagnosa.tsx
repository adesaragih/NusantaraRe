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
// ⛔ Tombol `Choose` (b2509 → `SetDisease` b2528) BELUM terpasang, dan
// sebabnya bukan kemalasan: `SetDisease` menulis `.DISEASE`/`.ICDCODE` pada
// halaman berkelas `Data-DiagnoseLife`, dan `.DiagnoseList`
// (`ClaimLifeDetailGCNM.xml` b3923) adalah **RepeatGrid** — banyak diagnosa
// per peserta. Tetapi `T_CLAIMLF_PREMIUMLIST_DETAIL` hanya punya `DISEASE`
// dan `ICD_CODE` TUNGGAL (migrasi 003). Satu lawan banyak, dan memilih
// salah satunya tanpa keputusan berarti membuang diagnosa orang. OQ-K.

import { useState } from 'react'

import { DETAIL } from '../assets/labels'
import { cariPenyakit, pesanGalat, type Penyakit } from '../services/api'
import { BelumTersedia } from './ui/dasar'

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

export function CariDiagnosa({ batas = 50 }: { batas?: number }) {
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
                      {/* ⛔ Lihat kepala berkas: satu lawan banyak, OQ-K. */}
                      <BelumTersedia apa={LABEL_CARI_DIAGNOSA.pilih} />
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
