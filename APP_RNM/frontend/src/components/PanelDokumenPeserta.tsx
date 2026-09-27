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

import { DOKUMEN } from '../assets/labels'
import type { Dokumen } from '../services/api'
import { BelumTersedia } from './ui/dasar'

/** Satu baris daftar, sesudah diputuskan apa yang tampil. */
export interface BarisDokumen {
  id: number
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

export function PanelDokumenPeserta({ dokumen }: { dokumen: Dokumen[] }) {
  const baris = barisDokumen(dokumen)
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
                  {b.terunggah ? (
                    <span title={PRANALA_MENUNGGU}>{PRANALA_MENUNGGU}</span>
                  ) : (
                    <span>Belum terunggah</span>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* ⛔ Ketiganya DINYATAKAN, bukan dihilangkan — lihat kepala berkas. */}
      <p className="dokumen__aksi">
        <BelumTersedia apa={DOKUMEN.tambahLampiran} />{' '}
        <BelumTersedia apa={DOKUMEN.lihatOfficeOnline} />{' '}
        <BelumTersedia apa={DOKUMEN.hapus} />
      </p>
    </section>
  )
}
