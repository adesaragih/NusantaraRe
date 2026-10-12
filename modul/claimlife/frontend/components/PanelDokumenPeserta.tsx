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
// ⛔ RALAT 28-09-2026 (sensus remark GILIRAN-12): bab ini dulu menyebut
// "penyimpangan sadar pada URL" - baris ber-URL kosong dibuang Pega
// (`LoadDocumentLife_ACT` 1.4.2, b1310). Langkah 1.4 beserta anaknya
// TER-REMARK (`//` b913); yang berjalan 1.3 b797, SELURUH baris tanpa URL.
// Daftar yang menampilkan semua baris karena itu PARITAS, bukan penyimpangan.
// Penanda "menunggu" tetap membedakan baris yang efek outbox-nya belum
// selesai (`tStorageId` kosong).

import { useEffect, useState } from 'react'

import { DOKUMEN } from '../labels'
import { simpanBlob } from '../../../../inti/frontend/lib/simpanBlob'
import { ambilIsiDokumenKlaim, ambilKategoriDokumen, hapusDokumen, unggahDokumen, type Dokumen } from '../api'
import { Modal } from '../../../../inti/frontend/components/ui/dasar'
import { pesanGalat } from '../../../../inti/frontend/klien'

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

/**
 * Kalimat baris yang efek penyimpanannya belum selesai (T_STORAGE_ID kosong). ⚠️ 11-10-2026: Claim Life BELUM punya pekerja
 * latar outbox, jadi baris ini tidak berpindah sendiri ke "View Office Online" sampai pekerjanya dinyalakan (lihat PARITAS).
 */
export const PRANALA_MENUNGGU = 'Processing — waiting for storage'

/** Tanggal dokumen dd/mm/yyyy HH:mm (11-10-2026, bukan ISO mentah); kosong / tak terbaca apa adanya. */
export function tanggalDokumen(teks: string): string {
  if (teks === '') return '—'
  const d = new Date(teks)
  if (Number.isNaN(d.getTime())) return teks
  const p = (n: number) => String(n).padStart(2, '0')
  return `${p(d.getDate())}/${p(d.getMonth() + 1)}/${d.getFullYear()} ${p(d.getHours())}:${p(d.getMinutes())}`
}

/** Kalimat modal konfirmasi Delete dokumen (keputusan work owner 11-10-2026: "buat document bisa dihapus"). */
export const KONFIRMASI_HAPUS_DOKUMEN = {
  judul: 'Delete this document?',
  peringatan: 'This cannot be undone.',
  batal: 'Cancel',
  ya: 'Yes, Delete',
  proses: 'Deleting…',
  lihat: 'View',
} as const

/** Lebar T_CLAIMLF_DOCUMENT.KATEGORI_2 (byte) — sama dengan backend services.LebarKategoriDokumen. */
export const LEBAR_KATEGORI_BYTE = 64

/** Nama Document Claim yang tidak muat di kolom kategori (ditolak server 422; tidak dipotong). */
export function kategoriTerlaluPanjang(k: string): boolean {
  return new TextEncoder().encode(k.trim()).length > LEBAR_KATEGORI_BYTE
}

/** Label panel (11-10-2026: ke Inggris). */
export const LABEL_DOKUMEN = {
  judul: 'Claim Documents',
  kosong: 'No claim documents for this insured yet.',
  namaFile: 'File name',
  kategori: 'Category',
  tanggal: 'Date',
  berkas: 'File',
  pilihKategori: 'Select category…',
  pilihBerkas: 'Choose file',
  kategoriKosong: 'No claim document types are set for this product in Master Product Name Life.',
  terlaluPanjang: 'Name too long to store (the category column must be widened first)',
} as const

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
  const [berkas, setBerkas] = useState<File | null>(null)
  // Dokumen yang sedang dikonfirmasi untuk dihapus (modal).
  const [hapus, setHapus] = useState<BarisDokumen | null>(null)
  const [sibuk, setSibuk] = useState(false)
  const [galat, setGalat] = useState<string | null>(null)
  // Dropdown "Category" (11-10-2026): daftar dari server — sumber yang SAMA dengan pemeriksa unggahan, jadi pilihan
  // pasti diterima. null = sedang dimuat.
  const [daftarKategori, setDaftarKategori] = useState<string[] | null>(null)
  useEffect(() => {
    let batal = false
    ambilKategoriDokumen(klaimID)
      .then((d) => {
        if (!batal) setDaftarKategori(d)
      })
      .catch((e: unknown) => {
        if (!batal) {
          setDaftarKategori([])
          setGalat(pesanGalat(e) ?? 'Document categories could not be loaded.')
        }
      })
    return () => {
      batal = true
    }
  }, [klaimID])

  async function jalankan(kerja: () => Promise<void>): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    try {
      await kerja()
      await onBerubah()
    } catch (e) {
      setGalat(pesanGalat(e) ?? 'The document change failed.')
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
      // Berbatas klaim; dokumen yang belum selesai disimpan dilayani dari salinan lokal server.
      simpanBlob(await ambilIsiDokumenKlaim(klaimID, b.id), b.namaFile)
    } catch (e) {
      setGalat(pesanGalat(e) ?? 'The document could not be downloaded.')
    }
  }

  return (
    <section className="dokumen">
      <h4 className="dokumen__judul">{LABEL_DOKUMEN.judul}</h4>

      {baris.length === 0 ? (
        <p className="dokumen__kosong">{LABEL_DOKUMEN.kosong}</p>
      ) : (
        <div className="inbox__gulir">{/* gulir horizontal, bukan terpotong (11-10-2026) */}
<table className="dokumen__tabel">
          <thead>
            <tr>
              <th scope="col">{LABEL_DOKUMEN.namaFile}</th>
              <th scope="col">{LABEL_DOKUMEN.kategori}</th>
              <th scope="col">{LABEL_DOKUMEN.tanggal}</th>
              <th scope="col">{LABEL_DOKUMEN.berkas}</th>
            </tr>
          </thead>
          <tbody>
            {baris.map((b) => (
              <tr key={b.id}>
                <td>{b.namaFile}</td>
                <td>{b.kategori}</td>
                <td>{tanggalDokumen(b.tanggal)}</td>
                <td>
                  {/* ✅ Penanda "URL menunggu penyambungan" DICABUT
                      27-09-2026 (butir be): unduhannya kini sungguhan, lewat
                      rute kita sendiri. Yang TETAP dibedakan adalah baris
                      yang efek outbox-nya belum selesai — `tStorageId`
                      kosong berarti SEDANG diproses, bukan hilang. */}
                  {/* View + Delete untuk SETIAP baris, termasuk yang masih diproses (11-10-2026). */}
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    onClick={() => {
                      void unduh(b)
                    }}
                  >
                    {KONFIRMASI_HAPUS_DOKUMEN.lihat}
                  </button>{' '}
                  <button
                    type="button"
                    className="btn btn--ghost btn--sm"
                    disabled={sibuk}
                    onClick={() => {
                      setHapus(b)
                    }}
                  >
                    {DOKUMEN.hapus}
                  </button>
                  {!b.terunggah && (
                    <small className="dokumen__status" title={PRANALA_MENUNGGU}>
                      {' '}
                      {PRANALA_MENUNGGU}
                    </small>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
</div>
      )}

      {galat !== null && (
        <p role="alert" className="alert alert--error">
          {galat}
        </p>
      )}

      {/* `Add attachment` b1245 (11-10-2026): Category = DROPDOWN dari server (dulu teks bebas → hampir selalu 422
          "kategori tidak dikenal", dan kotak berkas NONAKTIF sampai teks diisi — itulah "tidak bisa upload"). Pilih kategori
          dan berkas, lalu tombol "Add attachment" mengunggah. */}
      <div className="dokumen__aksi">
        <label>
          {LABEL_DOKUMEN.kategori}{' '}
          <select
            className="field__input"
            value={kategori}
            disabled={sibuk || daftarKategori === null || daftarKategori.length === 0}
            onChange={(e) => setKategori(e.target.value)}
          >
            <option value="">{LABEL_DOKUMEN.pilihKategori}</option>
            {/* Document Claim produk, urutan master. Nama > 64 byte belum dapat disimpan (kolom KATEGORI_2) — dinyatakan. */}
            {(daftarKategori ?? []).map((k) =>
              kategoriTerlaluPanjang(k) ? (
                <option key={k} value={k} disabled title={LABEL_DOKUMEN.terlaluPanjang}>
                  {k} — {LABEL_DOKUMEN.terlaluPanjang}
                </option>
              ) : (
                <option key={k} value={k}>
                  {k}
                </option>
              ),
            )}
          </select>
        </label>{' '}
        <label>
          {LABEL_DOKUMEN.pilihBerkas}{' '}
          <input
            type="file"
            disabled={sibuk}
            onChange={(e) => {
              setBerkas(e.target.files?.[0] ?? null)
            }}
          />
        </label>{' '}
        <button
          type="button"
          className="btn btn--primary btn--sm"
          disabled={sibuk || kategori === '' || berkas === null}
          onClick={() => {
            const f = berkas
            if (f === null) return
            void jalankan(async () => {
              await unggahDokumen(klaimID, pesertaID, f, kategori)
              setBerkas(null)
            })
          }}
        >
          {sibuk ? 'Uploading…' : DOKUMEN.tambahLampiran}
        </button>
        {daftarKategori !== null && daftarKategori.length === 0 && (
          <span className="os__galat-isian">{LABEL_DOKUMEN.kategoriKosong}</span>
        )}
      </div>
      {hapus !== null && (
        <Modal
          judul={KONFIRMASI_HAPUS_DOKUMEN.judul}
          onTutup={() => {
            if (!sibuk) setHapus(null)
          }}
          tanpaTutup={sibuk}
          labelBatal={KONFIRMASI_HAPUS_DOKUMEN.batal}
          aksi={
            <button
              type="button"
              className="btn btn--danger"
              disabled={sibuk}
              onClick={() => {
                const target = hapus
                void jalankan(async () => {
                  await hapusDokumen(klaimID, target.id)
                  setHapus(null)
                })
              }}
            >
              {sibuk ? KONFIRMASI_HAPUS_DOKUMEN.proses : KONFIRMASI_HAPUS_DOKUMEN.ya}
            </button>
          }
        >
          <p>
            <strong>{hapus.namaFile}</strong>
          </p>
          <p role="note" className="alert alert--warn">
            {KONFIRMASI_HAPUS_DOKUMEN.peringatan}
          </p>
        </Modal>
      )}
    </section>
  )
}
