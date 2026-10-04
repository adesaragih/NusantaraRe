// Layar A — daftar kontrak Treaty In.
//
// ⛔ DISALIN dari `Section/InputTreatyInOffer.xml` ekspor 2026-09, bukan
// dirancang di sini. Kesembilan kolom, urutannya, dan susunan tombol aksinya
// datang dari sana; jejak per medan ada di `../labels.ts`.
//
// ⚠️ `L-4` — *"tidak ada spesifikasi layar di mana pun"* — **DICABUT 3 Oktober
// 2026**: spesifikasinya ada, 50 Section + 3 Harness. Larangan mengarang layar
// tetap berlaku, dan kini tidak ada alasan untuk mengarang.
//
// ⭐ ISINYA NYATA sejak 3 Oktober 2026 — keputusan pemilik proses: layar ini
// membaca `POOLDATA.TREATY_IN`, tabel sistem lama, **1.854 baris**.
//
// ⚠️ Berkas ini pernah berbunyi "Tabelnya KOSONG hari ini (tiket 59 belum
// jalan)". Itu benar ketika ia membaca `KONTRAK`/`VERSI_KONTRAK` — model
// baru, yang memang nol baris. Layar ini **tidak menunggu tiket 59**; yang
// menunggu tiket 59 adalah daftar model baru, dan itu daftar yang berbeda.
//
// ⛔ BACA SAJA. Nol penulisan ke `TREATY_IN` dari jalur mana pun.

import { useEffect, useState } from 'react'

import { Gagal, Halaman, Kosong, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { ambilDaftarWarisan, type BarisDaftarWarisan } from '../api'
import { DAFTAR_KONTRAK, KOLOM_DAFTAR } from '../labels'

/**
 * Tombol aksi sebuah baris, dan susunannya bergantung KEADAAN barisnya.
 *
 * ⛔ Jangan diratakan. Baris berstatus `Resolve Complete` tidak lagi dapat
 * disunting, dan satu-satunya petunjuk di layar bahwa ia terkunci adalah
 * hilangnya tombol `Edit`. Menyeragamkan keempat tombol menghapus petunjuk itu
 * dan membuat pemakai menekan `Edit` yang akan ditolak di belakang.
 */
export function aksiUntuk(keadaan: string): readonly string[] {
  if (keadaan === DAFTAR_KONTRAK.keadaanTerkunci) {
    return [DAFTAR_KONTRAK.lihat, DAFTAR_KONTRAK.salin, DAFTAR_KONTRAK.revisi]
  }
  return [DAFTAR_KONTRAK.edit, DAFTAR_KONTRAK.lihat]
}

/**
 * Kolom layar -> medan baris warisan.
 *
 * ⛔ Pemetaannya berdiri SENDIRI, bukan tersebar di JSX: kesembilan kolom
 * datang dari ekspor Pega dan kesembilan medannya dari `TREATY_IN`, dan
 * pemetaan di antaranya adalah hal yang paling mudah salah diam-diam.
 * Berdiri sendiri, ia dapat diuji tanpa merender.
 */
export const MEDAN_WARISAN: Record<(typeof KOLOM_DAFTAR)[number]['kunci'], keyof BarisDaftarWarisan> = {
  id: 'id',
  namaKontrak: 'namaKontrak',
  sifatProporsi: 'sifatProporsi',
  idAsalBisnis: 'asalBisnis',
  idCedant: 'cedant',
  tanggalMulai: 'tanggalMulai',
  tanggalBerakhir: 'tanggalBerakhir',
  posisiKe: 'posisiKe',
  keadaanSiklusHidup: 'statusAkseptasi',
}

/**
 * Isi satu sel.
 *
 * ⚠️ Sel kosong ditandai `—` hanya karena NILAINYA kosong — bukan karena
 * kolomnya tidak ada. Bentuk lamanya menyangkal `Position To` sebagai
 * "belum punya rumah di model baru"; di tabel warisan `POSITIONUSERNAME`
 * ADA, terisi pada 32 dari 1.854 baris.
 */
function sel(b: BarisDaftarWarisan, kunci: (typeof KOLOM_DAFTAR)[number]['kunci']) {
  const v = b[MEDAN_WARISAN[kunci]]
  if (v === '') {
    return <span className="trin__redup">—</span>
  }
  return v
}

/**
 * Baris per halaman.
 *
 * ⚠️ Angka ini KEPUTUSAN KITA, bukan salinan: ekspor memakai penomoran Pega
 * (`pyGridPaginator`) yang ukurannya tidak tertulis di Section mana pun.
 * Dinyatakan di sini supaya tidak terbaca sebagai angka yang dibaca dari
 * sistem lama.
 */
export const UKURAN_HALAMAN = 25

export interface DaftarKontrakProps {
  /**
   * Membuka satu kontrak di form.
   *
   * ⛔ Pengenalnya TEKS. `TREATY_IN.ID` adalah `VARCHAR2(100)`, dan ke-1.854
   * nilainya kebetulan berupa angka tujuh digit — kebetulan bukan jaminan.
   */
  onBuka: (id: string) => void
  /** Membuka form kosong — tombol `Add`. */
  onTambah: () => void
}

export default function DaftarKontrakTreatyIn({ onBuka, onTambah }: DaftarKontrakProps) {
  const [baris, setBaris] = useState<BarisDaftarWarisan[] | null>(null)
  const [total, setTotal] = useState(0)
  const [galat, setGalat] = useState<unknown>(null)
  const [saringTampil, setSaringTampil] = useState(false)
  // Isi penyaring per kolom. `pyColumnFilteringDropDown` bernilai `true`
  // pada 1.049 sel dan `false` hanya pada 6 di seluruh ekspor, jadi
  // kesembilan kolom mendapat penyaringnya.
  const [saring, setSaring] = useState<Record<string, string>>({})
  const [halaman, setHalaman] = useState(1)

  // ⛔ SATU HALAMAN PER PERMINTAAN, dan itu syarat — bukan penghematan.
  // Tabelnya 1.854 baris; menariknya sekaligus ke peramban untuk menampilkan
  // 25 adalah ongkos yang tidak perlu dibayar siapa pun, dan ia tumbuh
  // bersama tabelnya.
  useEffect(() => {
    let dibuang = false
    setBaris(null)
    setGalat(null)
    ambilDaftarWarisan(halaman)
      .then((h) => {
        if (dibuang) return
        setBaris(h.baris)
        setTotal(h.total)
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [halaman])

  // ⚠️ Penyaring bekerja pada HALAMAN YANG SEDANG TAMPIL, bukan pada 1.854
  // baris — dan itu dinyatakan di layar, bukan disembunyikan. Menyaring
  // seluruh tabel menuntut parameter saring di `GET /kontrak-warisan`;
  // membangunnya sekarang berarti membangun yang belum ada yang memintanya.
  // Begitu rute itu menerima penyaring, blok ini yang pindah.
  const tersaring = (baris ?? []).filter((b) =>
    KOLOM_DAFTAR.every((k) => {
      const cari = (saring[k.kunci] ?? '').trim().toLowerCase()
      return cari === '' || b[MEDAN_WARISAN[k.kunci]].toLowerCase().includes(cari)
    }),
  )
  const tampil = baris === null ? null : tersaring

  return (
    <div className="inbox">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{DAFTAR_KONTRAK.judul}</h2>
        <button type="button" className="btn btn--primary" onClick={onTambah}>
          {DAFTAR_KONTRAK.tambah}
        </button>
        {/* Tombol ada di ekspor; yang dialihkannya BELUM dibangun, dan itu
            dinyatakan di layar alih-alih tombol yang diam-diam tidak
            melakukan apa pun. */}
        <div className="trin__kepala-kanan">
          <button
            type="button"
            className="btn btn--ghost"
            aria-expanded={saringTampil}
            onClick={() => {
              setSaringTampil((v) => !v)
            }}
          >
            {DAFTAR_KONTRAK.saring}
          </button>
          {/* Penomoran halaman di KANAN ATAS, seperti rujukan. Ia hanya
              muncul ketika ada baris: penomoran di atas tabel kosong
              menyatakan halaman yang tidak ada. */}
          {/* ⛔ `total` datang dari SERVER, bukan dari panjang halaman.
              Dengan 25 baris per halaman, `tampil.length` akan membuat
              penomoran mengira selalu ada satu halaman — dan tombol
              "berikutnya" tidak akan pernah muncul. */}
          {total > 0 && (
            <Halaman halaman={halaman} ukuran={UKURAN_HALAMAN} total={total} onPindah={setHalaman} />
          )}
        </div>
      </header>

      {galat !== null && <Gagal galat={galat} />}

      {/* ⛔ KEPALA KOLOM SELALU DIRENDER, termasuk saat nol baris.
          Sebelumnya seluruh <table> disembunyikan ketika kosong, sehingga
          kesembilan kolom yang disalin dari ekspor TIDAK TERLIHAT SAMA SEKALI
          — layar kosong tidak memperlihatkan rancangannya sendiri, dan
          pembacanya tidak punya cara tahu kolom apa yang akan datang. */}
      <section className="panel">
        <div className="table-wrap trin__tabel-daftar-bungkus">
          <table className="trin__tabel">
            <thead>
              <tr>
                {KOLOM_DAFTAR.map((k) => (
                  <th key={k.kunci} scope="col" className={`trin__k--${k.kunci}`}>
                    <span className="trin__kepala-kolom">
                      {k.label}
                      {/* Ikon saring per kolom — ekspor memberi SETIAP kolom
                          `pyColumnFilteringDropDown = true`. Menekannya
                          membuka baris saring, sama seperti tombol di kepala. */}
                      <button
                        type="button"
                        className="trin__ikon-saring"
                        aria-label={`Saring ${k.label}`}
                        aria-pressed={(saring[k.kunci] ?? '') !== ''}
                        onClick={() => {
                          setSaringTampil(true)
                        }}
                      >
                        ▾
                      </button>
                    </span>
                  </th>
                ))}
                <th scope="col">{DAFTAR_KONTRAK.aksi}</th>
              </tr>
              {saringTampil && (
                <tr className="trin__baris-saring">
                  {KOLOM_DAFTAR.map((k) => (
                    <td key={k.kunci}>
                      <input
                        type="text"
                        aria-label={`Saring ${k.label}`}
                        value={saring[k.kunci] ?? ''}
                        onChange={(e) => {
                          setHalaman(1)
                          setSaring((v) => ({ ...v, [k.kunci]: e.target.value }))
                        }}
                      />
                    </td>
                  ))}
                  <td />
                </tr>
              )}
            </thead>
            <tbody>
              {baris === null && galat === null && (
                <tr>
                  <td colSpan={KOLOM_DAFTAR.length + 1}>
                    <Memuat />
                  </td>
                </tr>
              )}
              {tampil !== null && tampil.length === 0 && (
                <tr>
                  <td colSpan={KOLOM_DAFTAR.length + 1}>
                    <Kosong pesan={DAFTAR_KONTRAK.kosong} petunjuk={DAFTAR_KONTRAK.kosongPetunjuk} />
                  </td>
                </tr>
              )}
              {/* ⛔ NOL pemotongan di sini: servernya yang sudah memotong.
                  Memotong dua kali membuat halaman kedua menampilkan
                  halaman kosong. */}
              {tampil?.map((b) => (
                <tr key={b.id}>
                  {KOLOM_DAFTAR.map((k) => (
                    <td key={k.kunci} className={`trin__k--${k.kunci}`}>
                      {sel(b, k.kunci)}
                    </td>
                  ))}
                  <td className="trin__aksi">
                    {aksiUntuk(b.statusAkseptasi).map((a) => (
                      <button
                        key={a}
                        type="button"
                        className="btn btn--ghost btn--sm"
                        onClick={() => {
                          onBuka(b.id)
                        }}
                      >
                        {a}
                      </button>
                    ))}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      {/* ⛔ KETERANGAN ASAL-USUL DILIPAT, BUKAN DIBUANG.
          Keempat kalimat ini benar dan berguna — dari mana barisnya dibaca,
          kenapa "Position To" sering kosong, kenapa Ceding berupa nama di
          sini dan pengenal di model baru. Tetapi sebagai empat baris prosa
          padat di bawah tabel, ia terbaca sebagai catatan rilis dan memakan
          layar yang isinya data. Dilipat: yang mencarinya menemukannya,
          yang tidak tidak terganggu. */}
      <details className="trin__kaki-lipat">
        <summary>{DAFTAR_KONTRAK.kakiJudul}</summary>
        <p className="trin__kaki">{DAFTAR_KONTRAK.catatanSumber}</p>
        <p className="trin__kaki">
          {DAFTAR_KONTRAK.keterangan} {DAFTAR_KONTRAK.catatanPosisiKe}{' '}
          {DAFTAR_KONTRAK.catatanNama}
        </p>
      </details>
    </div>
  )
}
