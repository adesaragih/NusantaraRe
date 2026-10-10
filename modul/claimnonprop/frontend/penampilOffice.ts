// Disalin dari `modul/claimprop/frontend/penampilOffice.ts` (pola, bukan impor; perintah work owner 09-10-2026 "untuk attachment
// juga mengikuti dari klaim prop"): keputusan work owner yang disebut di bawah adalah keputusan layar Claim Prop.
// Penampil lampiran klaim - jendela View File mengikuti NB Treaty In (perintah work owner 09-10-2026, popup
// `ReasViewAttachment`): `View` pdf / gambar dan `View Office Online` (`DownloadDocumentPolis` langkah 3:
// `"<penampil>?src=" + @encodeURL(url bertanda tangan)`).
//
// ⛔ SATU-SATUNYA alamat literal modul ini (pola Master Product Name Life `penampilOffice.ts`, keputusan work owner
// 03-10-2026 "izinkan ditulis di kode", dan NB Treaty In `inti/frontend/lampiran/penampil.ts`). Dibuka lewat FORM GET
// ber-`action` tetap ke bingkai di popup - bukan `window.open` / `location.href` (penjaga `unduhdokumen.test.ts`).

/** Alamat penampil kantor (tanpa kueri) - `action` form. */
export const PENAMPIL_OFFICE = 'https://view.officeapps.live.com/op/view.aspx'

/** Nama parameter kueri penampil (`?src=`) - `name` input tersembunyi form. */
export const PARAM_PENAMPIL = 'src'

/** Nama bingkai penampil - `target` form dan `name` iframe. */
export const BINGKAI_PENAMPIL = 'claimnonprop-penampil-lampiran'

/** `View Office Online` hanya untuk ekstensi ini (`ReasViewAttachment` b2951) - dan bila objeknya tercatat. */
const EKSTENSI_OFFICE: readonly string[] = ['xls', 'xlsx', 'doc', 'docx', 'ppt', 'pptx']

export function bisaViewOffice(ekstensi: string): boolean {
  return EKSTENSI_OFFICE.includes(ekstensi.toLowerCase())
}

/** `View` - pdf dan gambar raster; tipe objek URL dari ekstensi, bukan dari jawaban server (svg / html tidak dirender). */
const VIEW_ONLINE: Readonly<Record<string, { jenis: 'pdf' | 'gambar'; mime: string }>> = {
  pdf: { jenis: 'pdf', mime: 'application/pdf' },
  png: { jenis: 'gambar', mime: 'image/png' },
  jpg: { jenis: 'gambar', mime: 'image/jpeg' },
  jpeg: { jenis: 'gambar', mime: 'image/jpeg' },
  gif: { jenis: 'gambar', mime: 'image/gif' },
  bmp: { jenis: 'gambar', mime: 'image/bmp' },
  webp: { jenis: 'gambar', mime: 'image/webp' },
}

/** Jenis tampilan `View` untuk sebuah ekstensi; null = tidak ditawarkan. */
export function jenisViewOnline(ekstensi: string): 'pdf' | 'gambar' | null {
  const e = ekstensi.toLowerCase()
  return Object.hasOwn(VIEW_ONLINE, e) ? VIEW_ONLINE[e]!.jenis : null
}

/** Tipe isi objek URL `View`. */
export function mimeViewOnline(ekstensi: string): string {
  const e = ekstensi.toLowerCase()
  return Object.hasOwn(VIEW_ONLINE, e) ? VIEW_ONLINE[e]!.mime : ''
}
