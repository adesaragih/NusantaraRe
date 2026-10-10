// Penampil `View Office Online` - `AttachmentDetailBdx` b2562 -> `DownloadAttachmentBdx` langkah 2 b388 (PRE
// `Param.ViewOffice`): `LinkDocument.url = "<penampil>?src=" + @encodeURL(LinkDocument.url)`, URL bertanda tangan
// dari langkah 1 (`GetUrlGoogleStorage_Act`).
//
// ⛔ SATU-SATUNYA alamat literal modul ini, dan ia milik XML-nya (b388) - pola Master Product Name Life (keputusan
// work owner 03-10-2026 "izinkan ditulis di kode"); dijaga `bordereaux.test.ts`.
// ⛔ Dibuka lewat FORM GET ber-`action` tetap ke penampil ini, ke bingkai `BINGKAI_PENAMPIL` di popup - bukan
// `window.open` / `location.href` (penjaga lintas modul `unduhdokumen.test.ts`). Pega membuka jendela popup
// (`openUrlInWindow` b2217). Form GET menyandikan nilai seperti `@encodeURL`.

/** Alamat penampil kantor (b388, tanpa kueri) - `action` form. */
export const PENAMPIL_OFFICE = 'https://view.officeapps.live.com/op/view.aspx'

/** Nama parameter kueri penampil (b388 `?src=`) - `name` input tersembunyi form. */
export const PARAM_PENAMPIL = 'src'

/** Nama bingkai penampil di popup lampiran - `target` form dan `name` iframe. */
export const BINGKAI_PENAMPIL = 'bordereaux-penampil-lampiran'

/** `View Office Online` hanya untuk ekstensi ini (`AttachmentDetailBdx` b2876). */
export const EKSTENSI_OFFICE: readonly string[] = ['xls', 'xlsx', 'doc', 'docx', 'ppt', 'pptx']

export function bisaViewOffice(ekstensi: string): boolean {
  return EKSTENSI_OFFICE.includes(ekstensi.toLowerCase())
}

/**
 * `View` - pdf dan gambar raster yang ditampilkan langsung di popup penampil (permintaan work owner 08-10-2026, seperti
 * Product Name Life; tidak ada di XML Bordereaux). Isi datang dari rute unduh beridentitas, tipe objek URL-nya dari
 * ekstensi - bukan dari jawaban server, jadi berkas lain (svg, html) tidak pernah dirender.
 */
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

/** Tipe isi objek URL `View` - dari ekstensi. */
export function mimeViewOnline(ekstensi: string): string {
  const e = ekstensi.toLowerCase()
  return Object.hasOwn(VIEW_ONLINE, e) ? VIEW_ONLINE[e]!.mime : ''
}
