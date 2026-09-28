// Daftar dokumen pendukung SATU PESERTA — A3 kelompok 1.
//
// Meniru `Section/DocumentLife.xml`, dibaca 27-09-2026 bersama activity
// pemuatnya `Activity/LoadDocumentLife_ACT.xml`.
//
// ⛔ MILIK PESERTA, bukan klaim. `T_CLAIMLF_DOCUMENT.PREMIUM_LIST_DETAIL_ID`
// menunjuk baris peserta, dan `LoadDocumentLife_ACT` berkelas
// `Int-LIFE_PREMIUM_DETAIL` — saringan browse-nya membaca `.DOCUMENT` pada
// halaman peserta yang sedang berjalan. Menaruh daftar ini di tingkat klaim
// akan mencampur dokumen milik peserta yang berbeda, dan itu kesalahan yang
// sama dengan panel total sebelum diralat.
//
// ⛔ APA YANG BELUM ADA, DAN KENAPA DISEBUT. Dari empat tombol section itu,
// hanya daftarnya yang dibangun kelompok ini:
//
//   b611  `Refresh`             -> `pyAction refresh` b619
//   b1245 `Add attachment`      -> localAction `AttachDocumentLife` b1273
//   b3502 `View Office Online`  -> runActivity `DownloadDocumentClaim` b3519
//   b4288 `Delete`              -> localAction `ConfirmDeleteAttachment` b4317
//                                  DAN `closeContainer` b4441
//
// Ketiganya selain Refresh menyentuh penyimpanan luar (Google Storage,
// ADR-U-0010) dan outbox `T_LOG_SERVICE_RNM` — itu kelompok Dokumen.
// Tombolnya DINYATAKAN di sini, bukan dihilangkan: layar yang tampak lengkap
// padahal tidak adalah layar yang tidak akan dicari lagi (butir av).
//
// ⛔ PENYIMPANGAN SADAR pada URL. `LoadDocumentLife_ACT` langkah 1.4.2
// berprasyarat `DataImage.URLImage==""` dengan `WhenTrue=3` — artinya baris
// yang URL penyimpanannya kosong TIDAK ikut ditambahkan ke daftar. Di sini
// panggilan ke Google Storage belum dilakukan, sehingga URL-nya SELALU
// kosong; meniru gerbang itu akan membuat daftar SELALU kosong dan layar
// berkata "tidak ada dokumen" untuk peserta yang dokumennya lengkap.
// Barisnya karena itu tetap tampil, dengan penanda bahwa pranalanya menunggu
// penyambungan. Arah selisihnya disengaja: daftar yang menyebut apa yang
// belum tersambung lebih jujur daripada daftar yang diam-diam kosong.

import { useState } from 'react'

import { DOKUMEN } from '../../assets/labels.claimlife'
import { simpanBlob } from '../../lib/simpanBlob'
import {
  ambilIsiDokumen,
  hapusDokumen,
  pesanGalat,
  unggahDokumen,
  type Dokumen,
} from '../../services/api'

/** Satu baris daftar, sesudah diputuskan apa yang tampil. */
export interface BarisDokumen {
  /** TEKS — lihat `Dokumen.id`: 17 angka di luar jangkauan aman JS. */
  id: string
  /** Nama berkas. ⛔ Nama BERKAS, tidak pernah nama orang. */
  namaFile: string
  /** Kategori yang `DocumentLife.xml` tampilkan kepada manusia. */
  kategori: string
  /**
   * Teks tanggal, SUDAH dinormalkan.
   *
   * ⛔ Backend mengirim `null` untuk kolom DATE yang NULL (`*time.Time`),
   * bukan `""`. Keduanya dinormalkan menjadi `''` DI SINI, satu tempat -
   * sehingga komponen dan ujinya cukup mengenal satu bentuk "kosong".
   */
  tanggal: string
  /**
   * Apakah berkasnya sudah terunggah ke penyimpanan luar.
   *
   * ⚠️ Ini BUKAN "punya URL". URL-nya belum pernah diminta sama sekali —
   * yang ada hanya penunjuk `T_STORAGE_ID`. Membedakan keduanya penting:
   * "belum diunggah" dan "sudah diunggah tetapi pranalanya belum tersambung"
   * adalah dua keadaan berbeda, dan hanya yang kedua yang menunggu kita.
   */
  terunggah: boolean
}

/** Kalimat penanda pranala yang belum tersambung. */
export const PRANALA_MENUNGGU = 'URL menunggu penyambungan penyimpanan'

/**
 * Menyusun baris daftar dari jawaban backend.
 *
 * ⚠️ Dipisah dari komponennya supaya dapat diuji tanpa DOM.
 */
export function barisDokumen(dokumen: Dokumen[]): BarisDokumen[] {
  return dokumen.map((d) => ({
    id: d.id,
    namaFile: d.namaFile,
    // `DocumentLife.xml` menampilkan KATEGORI_2 kepada manusia; kategori1
    // adalah penggolong yang dipakai saringan browse-nya.
    kategori: d.kategori2,
    // ⛔ `?? ''` bukan kosmetik: `null` dan `''` sama-sama berarti "kolomnya
    // kosong", dan tanpa ini penanda — tidak pernah menyala untuk `null`.
    tanggal: d.tanggal ?? '',
    terunggah: d.tStorageId !== '',
  }))
}

export function PanelDokumenPeserta({
  klaimID,
  pesertaID,
  dokumen,
  onBerubah,
}: {
  klaimID: string
  pesertaID: string
  dokumen: Dokumen[]
  /** Dipanggil sesudah setiap perubahan supaya pemanggil membaca ulang. */
  onBerubah: () => void | Promise<void>
}) {
  const baris = barisDokumen(dokumen)
  const [kategori, setKategori] = useState('')
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<string | null>(null)

  async function jalankan(kerja: () => Promise<void>): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    try {
      await kerja()
      await onBerubah()
    } catch (e) {
      setGalat(pesanGalat(e) ?? 'Perubahan dokumen gagal.')
    } finally {
      setSibuk(false)
    }
  }

  /**
   * `View Office Online` b3502 - unduhan LEWAT klien, supaya identitas ikut.
   *
   * ⚠️ Tidak memanggil `onBerubah`: mengunduh tidak mengubah apa pun, dan
   * membaca ulang daftar sesudahnya hanya menambah satu permintaan.
   */
  async function unduh(b: BarisDokumen): Promise<void> {
    setGalat(null)
    try {
      simpanBlob(await ambilIsiDokumen(b.id), b.namaFile)
    } catch (e) {
      setGalat(pesanGalat(e) ?? 'Gagal mengunduh dokumen.')
    }
  }

  return (
    <section className="dokumen">
      <h4 className="dokumen__judul">Dokumen pendukung</h4>

      {baris.length === 0 ? (
        <p className="dokumen__kosong">Belum ada dokumen pada peserta ini.</p>
      ) : (
        <table className="dokumen__tabel">
          <thead>
            <tr>
              <th scope="col">Nama berkas</th>
              <th scope="col">Kategori</th>
              <th scope="col">Tanggal</th>
              <th scope="col">Berkas</th>
            </tr>
          </thead>
          <tbody>
            {baris.map((b) => (
              <tr key={b.id}>
                <td>{b.namaFile}</td>
                <td>{b.kategori}</td>
                <td>{b.tanggal === '' ? '—' : b.tanggal}</td>
                <td>
                  {/* ✅ Penanda "URL menunggu penyambungan" DICABUT
                      27-09-2026 (butir be): unduhannya kini sungguhan, lewat
                      rute kita sendiri. Yang TETAP dibedakan adalah baris
                      yang efek outbox-nya belum selesai — `tStorageId`
                      kosong berarti SEDANG diproses, bukan hilang. */}
                  {b.terunggah ? (
                    <>
                      <button
                        type="button"
                        onClick={() => {
                          void unduh(b)
                        }}
                      >
                        {DOKUMEN.lihatOfficeOnline}
                      </button>{' '}
                      <button
                        type="button"
                        disabled={sibuk}
                        onClick={() => {
                          void jalankan(() => hapusDokumen(klaimID, b.id))
                        }}
                      >
                        {DOKUMEN.hapus}
                      </button>
                    </>
                  ) : (
                    <span title={PRANALA_MENUNGGU}>{PRANALA_MENUNGGU}</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {galat !== null && <p role="alert">{galat}</p>}

      {/* `Add attachment` b1245 — kategori WAJIB: unggahan menolak kategori
          di luar daftar butir ar1. (Gerbang kelengkapan per kategori di Save
          ke Outstanding ter-remark di XML dan dibuang — butir bl.) */}
      <p className="dokumen__aksi">
        <label>
          Kategori{' '}
          <input
            type="text"
            value={kategori}
            onChange={(e) => setKategori(e.target.value)}
          />
        </label>{' '}
        <label>
          {DOKUMEN.tambahLampiran}{' '}
          <input
            type="file"
            disabled={sibuk || kategori.trim() === ''}
            onChange={(e) => {
              const f = e.target.files?.[0]
              // Kotak berkas DIKOSONGKAN lagi: tanpa itu, memilih berkas
              // yang SAMA dua kali tidak memicu `change` sama sekali, dan
              // unggahan kedua tampak diabaikan tanpa sebab.
              e.target.value = ''
              if (f === undefined) return
              void jalankan(async () => {
                await unggahDokumen(klaimID, pesertaID, f, kategori.trim())
              })
            }}
          />
        </label>
      </p>
    </section>
  )
}
