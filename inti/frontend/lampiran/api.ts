// Klien lampiran "Reas" - rute `inti/backend/dokumenpolis.Pasang` di bawah jalur dasar modul pemakai
// (mis. `/api/nb-treaty-in/kasus/<id>/lampiran`). Kategori lewat kueri: NOTE memuat garis miring (`R/I SLIP`).

import { BATAS_WAKTU_MS, kegagalanDari, minta, mintaFormulir, rakitURL, unduhBerkasBeridentitas } from '../klien'
import { headerIdentitas } from '../store/sesi'

/** Satu baris grid `AttachmentGridReas`. */
export interface KategoriReas {
  nama: string
  cacah: number
}

/** Jawaban grid: kategori dan apakah Upload / Delete boleh (kasus belum Resolve). */
export interface GridReas {
  daftar: KategoriReas[]
  bolehUbah: boolean
}

/** Satu dokumen `DOCUMENT_POLIS`. Kunci objek penyimpanan TIDAK dikirim backend. */
export interface DokumenReas {
  id: string
  namaFile: string
  /** `MIME` - ekstensi huruf kecil. */
  ekstensi: string
  /** `KATEGORI_2` - kolom Note. */
  kategori: string
  /** `TANGGAL` - kolom Upload Date, `DD-MM-YYYY HH:mm`. */
  tanggal: string
  pengunggah: string
  /** Objek penyimpanan tercatat (View Office Online, b2951). */
  adaObjek: boolean
}

function jalurDokumen(dasar: string, d: DokumenReas): string {
  return `${dasar}/dokumen/${encodeURIComponent(d.id)}`
}

/** Grid kategori (`SetCategoryAttach`). */
export function ambilGrid(dasar: string): Promise<GridReas> {
  return minta<GridReas>(dasar)
}

/** Popup View File (`InputParamUploadReas_act`). */
export function ambilDokumen(dasar: string, kategori: string): Promise<{ daftar: DokumenReas[] }> {
  return minta(`${dasar}/dokumen`, { kueri: { kategori } })
}

/** Upload File - satu berkas per permintaan (`InsertDocument_Act`). */
export function unggahDokumen(dasar: string, kategori: string, berkas: File): Promise<DokumenReas> {
  const isi = new FormData()
  isi.append('berkas', berkas)
  return mintaFormulir(`${dasar}/dokumen?kategori=${encodeURIComponent(kategori)}`, isi)
}

/** Tautan File (`DownloadDocumentPolis`) - fetch beridentitas, bukan pranala. */
export function unduhDokumen(dasar: string, d: DokumenReas): Promise<void> {
  return unduhBerkasBeridentitas(`${jalurDokumen(dasar, d)}/isi`, d.namaFile)
}

/** Isi satu dokumen sebagai Blob - untuk `View` pdf / gambar di popup penampil. */
export async function ambilIsiDokumen(dasar: string, d: DokumenReas): Promise<Blob> {
  const kendali = new AbortController()
  const jam = setTimeout(() => {
    kendali.abort()
  }, BATAS_WAKTU_MS)
  try {
    const jawab = await fetch(rakitURL(`${jalurDokumen(dasar, d)}/isi`), {
      method: 'GET',
      headers: { ...headerIdentitas() },
      signal: kendali.signal,
    })
    if (!jawab.ok) throw kegagalanDari(jawab.status, await jawab.text())
    return await jawab.blob()
  } finally {
    clearTimeout(jam)
  }
}

/** View Office Online - URL bertanda tangan untuk penampil (`DownloadDocumentPolis` ViewOffice). */
export function tautanOffice(dasar: string, d: DokumenReas): Promise<{ url: string }> {
  return minta(`${jalurDokumen(dasar, d)}/office`)
}

/** Delete (`DeleteDocumentPolis_Act`). */
export function hapusDokumen(dasar: string, d: DokumenReas): Promise<{ ok: boolean }> {
  return minta(`${jalurDokumen(dasar, d)}/hapus`, { metode: 'POST' })
}
