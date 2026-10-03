// Penampil `View Office Online` b69291 - `DownloadAttProdName_Act` langkah 7 b1080 (PRE `Param.ViewOffice`):
// `LinkDocument.url = "<penampil>?src=" + @encodeURL(LinkDocument.url)` (b1103), URL bertanda tangan dari langkah 6 b953.
//
// ⛔ SATU-SATUNYA alamat literal modul ini: keputusan work owner 03-10-2026 "izinkan ditulis di kode" (OQ-MPNL-11
// dibalik) - pengecualian bernama `alamatDiizinkan` di `TestMPNLNolAlamatLayanan`. Berkas ini hanya boleh memuat
// alamat ini.
// ⛔ Dibuka lewat FORM GET berakhir tetap ke penampil ini (`target="_blank"`, input `src`) - bukan `window.open` /
// `location.href` (penjaga lintas-modul `unduhdokumen.test.ts`: tidak ada navigasi yang dapat menuju backend). Form
// GET menyandikan nilai seperti `@encodeURL` (`application/x-www-form-urlencoded`).

/** Alamat penampil kantor (b1103, tanpa kueri) - `action` form. */
export const PENAMPIL_OFFICE = 'https://view.officeapps.live.com/op/view.aspx'

/** Nama parameter kueri penampil (b1103 `?src=`) - `name` input tersembunyi form. */
export const PARAM_PENAMPIL = 'src'
