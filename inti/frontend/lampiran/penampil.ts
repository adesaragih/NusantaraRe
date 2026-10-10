// Penampil lampiran - `View Office Online` (`DownloadDocumentPolis` langkah 3: `LinkDocument.url =
// "<penampil>?src=" + @encodeURL(LinkDocument.url)`) dan `View` pdf / gambar (permintaan work owner 08-10-2026, seperti
// Product Name Life).
//
// ⛔ SATU-SATUNYA alamat literal panel ini, dan ia milik XML-nya - pola Master Product Name Life (keputusan work owner
// 03-10-2026 "izinkan ditulis di kode"); dijaga `lampiran.test.ts`. Dibuka lewat FORM GET ber-`action` tetap ke
// penampil ini, ke bingkai di popup - bukan `window.open` / `location.href` (penjaga `unduhdokumen.test.ts`).

/** Alamat penampil kantor (tanpa kueri) - `action` form. */
export const PENAMPIL_OFFICE = 'https://view.officeapps.live.com/op/view.aspx'

/** Nama parameter kueri penampil (`?src=`) - `name` input tersembunyi form. */
export const PARAM_PENAMPIL = 'src'

/** Nama bingkai penampil - `target` form dan `name` iframe. */
export const BINGKAI_PENAMPIL = 'lampiran-reas-penampil'

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

/** Berkas dari pilihan dan seret-lepas digabung tanpa ganda (nama tanpa beda huruf besar + ukuran). */
export function gabungBerkas(lama: readonly File[], baru: readonly File[]): File[] {
  const kunci = (f: File) => `${f.name.toLowerCase()}|${f.size}`
  const ada = new Set(lama.map(kunci))
  const hasil = [...lama]
  for (const f of baru) {
    if (ada.has(kunci(f))) continue
    ada.add(kunci(f))
    hasil.push(f)
  }
  return hasil
}
