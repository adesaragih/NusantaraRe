// DI-PORT APA ADANYA dari `REFERENSI_UI/frontend/src/lib/exportXlsx.ts`
// (butir bg, 28-09-2026). Nol impor - ia mandiri, dan itulah sebab ia
// satu-satunya dari enam berkas pola grid yang dibawa.
//
// ⛔ LIMA YANG TIDAK DI-PORT, dan sebabnya: `RecordForm`,
// `BilahSaringRegistry`, `saringRegistry` bergantung pada tipe
// `MasterEntity`/`MasterField` - registry master yang aplikasi ini tidak
// punya. `CariSebaris` dan `PilihCari` bergantung pada
// `@tanstack/react-query` dan endpoint lookup `?q=` yang juga belum ada.
// Membawanya berarti membawa kode untuk daftar yang tidak ada - kode
// mati, yang brief larang.
//
// ⚠️ Dua penyesuaian, keduanya karena `noUncheckedIndexedAccess`
// menyala di tsconfig kami dan tidak di sana. Nol perubahan perilaku.
/**
 * Ekspor tabel ke `.xlsx` ASLI (OOXML) — TANPA menambah dependensi
 * (proyek ini melarangnya, lihat komponen/ui/dasar.tsx §IKON).
 *
 * # Kenapa `.xlsx` sungguhan bisa dibuat tanpa library
 *
 * `.xlsx` adalah arsip ZIP berisi beberapa berkas XML. ZIP mendukung
 * metode "STORED" (tanpa kompresi, metode 0) yang SAH menurut spesifikasi
 * — Excel membukanya sama seperti ZIP terkompresi, hanya berkasnya lebih
 * besar. Karena STORED tidak butuh algoritma kompresi (DEFLATE), seluruh
 * arsip dapat ditulis dengan operasi byte murni: tanpa JSZip, tanpa
 * exceljs, tanpa dependensi apa pun.
 *
 * Struktur minimal yang Excel terima (dan diverifikasi dokumentasi
 * OOXML/ECMA-376): `[Content_Types].xml`, `_rels/.rels`,
 * `xl/workbook.xml`, `xl/_rels/workbook.xml.rels`,
 * `xl/worksheets/sheet1.xml`, plus `docProps/core.xml` + `docProps/app.xml`
 * untuk kompatibilitas maksimal (tanpanya sebagian versi Excel lama
 * menampilkan dialog "perbaiki berkas" — jinak, tapi dihindari di sini).
 *
 * Sel numerik (`typeof nilai === "number"`) ditulis sebagai `<v>` numerik
 * asli — Excel dapat menghitungnya (SUM, dst), bukan sekadar teks
 * berformat. Sel lain (kode, tanggal-teks, label) ditulis `inlineStr`
 * supaya angka berawalan nol (mis. kode cabang "017") TIDAK kehilangan
 * nol-nya — kelas kesalahan yang justru sering terjadi pada CSV/Excel.
 */

// ===========================================================================
// CRC-32 — dibutuhkan header ZIP setiap berkas.
// ===========================================================================

const CRC_TABLE: Uint32Array = (() => {
  const table = new Uint32Array(256);
  for (let n = 0; n < 256; n++) {
    let c = n;
    for (let k = 0; k < 8; k++) {
      c = c & 1 ? 0xedb88320 ^ (c >>> 1) : c >>> 1;
    }
    table[n] = c >>> 0;
  }
  return table;
})();

/** crc32 menghitung checksum standar ZIP/PNG atas suatu buffer byte. */
export function crc32(buf: Uint8Array): number {
  let crc = 0xffffffff;
  for (let i = 0; i < buf.length; i++) {
    // tsconfig kami menyalakan noUncheckedIndexedAccess; indeks tabel
    // dijamin 0..255 oleh mask `& 0xff`, dan itu dinyatakan.
    crc = (CRC_TABLE[(crc ^ (buf[i] ?? 0)) & 0xff] ?? 0) ^ (crc >>> 8);
  }
  return (crc ^ 0xffffffff) >>> 0;
}

// ===========================================================================
// Penulis byte kecil — menghindari penggabungan string besar berulang.
// ===========================================================================

class PenulisByte {
  private potongan: Uint8Array[] = [];
  private panjang = 0;

  tambah(b: Uint8Array): void {
    this.potongan.push(b);
    this.panjang += b.length;
  }
  u16(v: number): void {
    this.tambah(Uint8Array.of(v & 0xff, (v >>> 8) & 0xff));
  }
  u32(v: number): void {
    this.tambah(
      Uint8Array.of(v & 0xff, (v >>> 8) & 0xff, (v >>> 16) & 0xff, (v >>> 24) & 0xff),
    );
  }
  teks(s: string): void {
    this.tambah(new TextEncoder().encode(s));
  }
  hasil(): Uint8Array {
    const out = new Uint8Array(this.panjang);
    let o = 0;
    for (const p of this.potongan) {
      out.set(p, o);
      o += p.length;
    }
    return out;
  }
}

// ===========================================================================
// Penulis ZIP — metode STORED (tanpa kompresi), nol dependensi.
// ===========================================================================

interface EntriZip {
  nama: string;
  data: Uint8Array;
}

/** TANGGAL_ZIP_BAWAAN = 1980-01-01 dalam format tanggal DOS ZIP (bit-packed:
 *  hari 1-31 | bulan 1-12 | tahun sejak 1980). Nilai apa pun yang SAH
 *  cukup — pembaca ZIP tidak memvalidasi terhadap tanggal sebenarnya. */
const TANGGAL_ZIP_BAWAAN = 0x21;

/**
 * bangunZip merakit arsip ZIP dari daftar entri — local file header per
 * berkas, lalu central directory, lalu end-of-central-directory record.
 * Urutan dan bentuk field mengikuti PKZIP APPNOTE.TXT §4.3.
 */
export function bangunZip(entri: EntriZip[]): Uint8Array {
  const bagianLokal: Uint8Array[] = [];
  const bagianPusat: Uint8Array[] = [];
  let ofsBerjalan = 0;
  const infoPusat: { namaBytes: Uint8Array; crc: number; ukuran: number; ofs: number }[] = [];

  for (const e of entri) {
    const namaBytes = new TextEncoder().encode(e.nama);
    const crc = crc32(e.data);

    const w = new PenulisByte();
    w.u32(0x04034b50); // signature local file header
    w.u16(20); // version needed to extract
    w.u16(0); // general purpose flag
    w.u16(0); // compression method = STORED
    w.u16(0); // last mod time
    w.u16(TANGGAL_ZIP_BAWAAN); // last mod date
    w.u32(crc);
    w.u32(e.data.length); // compressed size == uncompressed (STORED)
    w.u32(e.data.length);
    w.u16(namaBytes.length);
    w.u16(0); // extra field length
    w.tambah(namaBytes);
    w.tambah(e.data);

    const lokal = w.hasil();
    bagianLokal.push(lokal);
    infoPusat.push({ namaBytes, crc, ukuran: e.data.length, ofs: ofsBerjalan });
    ofsBerjalan += lokal.length;
  }

  const ofsPusatMulai = ofsBerjalan;
  for (const info of infoPusat) {
    const w = new PenulisByte();
    w.u32(0x02014b50); // signature central directory header
    w.u16(20); // version made by
    w.u16(20); // version needed
    w.u16(0); // flag
    w.u16(0); // compression = STORED
    w.u16(0); // time
    w.u16(TANGGAL_ZIP_BAWAAN); // date
    w.u32(info.crc);
    w.u32(info.ukuran);
    w.u32(info.ukuran);
    w.u16(info.namaBytes.length);
    w.u16(0); // extra length
    w.u16(0); // comment length
    w.u16(0); // disk number start
    w.u16(0); // internal file attributes
    w.u32(0); // external file attributes
    w.u32(info.ofs); // relative offset of local header
    w.tambah(info.namaBytes);
    bagianPusat.push(w.hasil());
  }
  const ukuranPusat = bagianPusat.reduce((a, c) => a + c.length, 0);

  const eocd = new PenulisByte();
  eocd.u32(0x06054b50); // signature end of central directory
  eocd.u16(0); // disk number
  eocd.u16(0); // disk with start of CD
  eocd.u16(entri.length); // entries on this disk
  eocd.u16(entri.length); // total entries
  eocd.u32(ukuranPusat);
  eocd.u32(ofsPusatMulai);
  eocd.u16(0); // comment length

  const total = new PenulisByte();
  for (const p of bagianLokal) total.tambah(p);
  for (const p of bagianPusat) total.tambah(p);
  total.tambah(eocd.hasil());
  return total.hasil();
}

// ===========================================================================
// XML XLSX minimal.
// ===========================================================================

/** xmlEscape membungkus karakter reserved XML, dan membuang karakter
 *  kontrol yang XML 1.0 TIDAK IZINKAN sama sekali (selain tab/LF/CR) —
 *  membiarkannya akan menghasilkan berkas yang Excel tolak buka. */
export function xmlEscape(s: string): string {
  return s
    // eslint-disable-next-line no-control-regex
    .replace(/[\x00-\x08\x0B\x0C\x0E-\x1F]/g, "")
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

/** kolomKe merakit huruf kolom Excel dari indeks 0-based (0 -> A, 25 -> Z,
 *  26 -> AA, dst) — algoritma bijection basis-26 standar. */
export function kolomKe(idxNol: number): string {
  let idx = idxNol + 1;
  let s = "";
  while (idx > 0) {
    const sisa = (idx - 1) % 26;
    s = String.fromCharCode(65 + sisa) + s;
    idx = Math.floor((idx - 1) / 26);
  }
  return s;
}

/** Satu kolom ekspor: kunci field di baris data + label header. */
export interface KolomEksporXlsx<T> {
  kunci: keyof T;
  label: string;
}

/** selXml merender SATU sel — numerik (typeof number) sebagai <v> asli
 *  (Excel dapat menghitungnya), lainnya sebagai inlineStr (kode berawalan
 *  nol TIDAK kehilangan nolnya, beda dari sel numerik). */
function selXml(ref: string, nilai: unknown): string {
  if (typeof nilai === "number" && Number.isFinite(nilai)) {
    return `<c r="${ref}"><v>${nilai}</v></c>`;
  }
  const teks = nilai == null ? "" : String(nilai);
  return `<c r="${ref}" t="inlineStr"><is><t xml:space="preserve">${xmlEscape(teks)}</t></is></c>`;
}

/**
 * rakitSheetXml merender seluruh isi lembar (header + baris data) sebagai
 * XML `sheet1.xml` — dipisah dari rakitXlsxBytes supaya dapat diuji tanpa
 * merakit ZIP-nya.
 */
export function rakitSheetXml<T>(kolom: KolomEksporXlsx<T>[], baris: T[]): string {
  const baris1 = kolom.map((k, i) => selXml(kolomKe(i) + "1", k.label)).join("");
  const barisData = baris
    .map((b, r) => {
      const nomorBaris = r + 2; // baris 1 = header
      const sel = kolom
        .map((k, i) => selXml(kolomKe(i) + String(nomorBaris), b[k.kunci]))
        .join("");
      return `<row r="${nomorBaris}">${sel}</row>`;
    })
    .join("");
  return (
    '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>' +
    '<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">' +
    "<sheetData>" +
    `<row r="1">${baris1}</row>` +
    barisData +
    "</sheetData>" +
    "</worksheet>"
  );
}

const CONTENT_TYPES_XML =
  '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>' +
  '<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">' +
  '<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>' +
  '<Default Extension="xml" ContentType="application/xml"/>' +
  '<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>' +
  '<Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>' +
  '<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>' +
  '<Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>' +
  "</Types>";

const ROOT_RELS_XML =
  '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>' +
  '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">' +
  '<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>' +
  '<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>' +
  '<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>' +
  "</Relationships>";

const WORKBOOK_XML =
  '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>' +
  '<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">' +
  '<sheets><sheet name="Sheet1" sheetId="1" r:id="rId1"/></sheets>' +
  "</workbook>";

const WORKBOOK_RELS_XML =
  '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>' +
  '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">' +
  '<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>' +
  "</Relationships>";

const CORE_XML =
  '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>' +
  '<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/">' +
  "<dc:creator>Treaty V2</dc:creator>" +
  "<cp:lastModifiedBy>Treaty V2</cp:lastModifiedBy>" +
  "</cp:coreProperties>";

const APP_XML =
  '<?xml version="1.0" encoding="UTF-8" standalone="yes"?>' +
  '<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties">' +
  "<Application>Treaty V2</Application>" +
  "</Properties>";

/**
 * rakitXlsxBytes merakit arsip `.xlsx` LENGKAP sebagai bytes — dipisah
 * dari unduhXlsx (efek DOM) supaya dapat diuji murni.
 */
export function rakitXlsxBytes<T>(kolom: KolomEksporXlsx<T>[], baris: T[]): Uint8Array {
  const enc = new TextEncoder();
  return bangunZip([
    { nama: "[Content_Types].xml", data: enc.encode(CONTENT_TYPES_XML) },
    { nama: "_rels/.rels", data: enc.encode(ROOT_RELS_XML) },
    { nama: "docProps/core.xml", data: enc.encode(CORE_XML) },
    { nama: "docProps/app.xml", data: enc.encode(APP_XML) },
    { nama: "xl/workbook.xml", data: enc.encode(WORKBOOK_XML) },
    { nama: "xl/_rels/workbook.xml.rels", data: enc.encode(WORKBOOK_RELS_XML) },
    { nama: "xl/worksheets/sheet1.xml", data: enc.encode(rakitSheetXml(kolom, baris)) },
  ]);
}

/** unduhXlsx memicu unduhan `.xlsx` lewat elemen <a> sementara — pola
 *  sama dengan unduhCsv (lib/exportCsv.ts). */
export function unduhXlsx<T>(namaBerkas: string, kolom: KolomEksporXlsx<T>[], baris: T[]): void {
  const bytes = rakitXlsxBytes(kolom, baris);
  const blob = new Blob([bytes.buffer as ArrayBuffer], {
    type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = namaBerkas.endsWith(".xlsx") ? namaBerkas : namaBerkas + ".xlsx";
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}
