// Uji exportXlsx (r189) — tdd. Menguji BENTUK BYTE arsip ZIP/XLSX secara
// murni (tanpa DOM/Blob — unduhXlsx efek DOM tidak diuji di sini, sama
// pola dengan exportCsv.test.ts).

import { expect, test } from "vitest";

import { bangunZip, crc32, kolomKe, rakitSheetXml, rakitXlsxBytes, xmlEscape } from "./exportXlsx";

test("crc32 — nilai dikenal (string kosong = 0)", () => {
  expectEqual(crc32(new Uint8Array()), 0);
});

test("crc32 — 'ABC' cocok nilai referensi standar", () => {
  // Nilai referensi CRC-32 untuk "ABC" (dihitung independen, algoritma
  // IEEE 802.3 standar yang dipakai ZIP/PNG/gzip).
  const bytes = new TextEncoder().encode("ABC");
  expectEqual(crc32(bytes), 0xa3830348);
});

test("kolomKe — pemetaan indeks ke huruf kolom Excel", () => {
  expectEqual(kolomKe(0), "A");
  expectEqual(kolomKe(25), "Z");
  expectEqual(kolomKe(26), "AA");
  expectEqual(kolomKe(27), "AB");
  expectEqual(kolomKe(51), "AZ");
  expectEqual(kolomKe(52), "BA");
});

test("xmlEscape — membungkus karakter reserved XML", () => {
  expectEqual(xmlEscape("A & B"), "A &amp; B");
  expectEqual(xmlEscape("<tag>"), "&lt;tag&gt;");
  expectEqual(xmlEscape('kata "kutip"'), "kata &quot;kutip&quot;");
});

test("xmlEscape — membuang karakter kontrol yang XML 1.0 tolak", () => {
  // eslint-disable-next-line no-control-regex
  const rusak = "A\x01B\x08C";
  expectEqual(xmlEscape(rusak), "ABC");
});

test("xmlEscape — tab/LF/CR TETAP dipertahankan (XML 1.0 mengizinkannya)", () => {
  expectEqual(xmlEscape("A\tB\nC"), "A\tB\nC");
});

test("rakitSheetXml — header dari label, sel numerik pakai <v>, sel teks pakai inlineStr", () => {
  const xml = rakitSheetXml(
    [
      { kunci: "kode", label: "Kode" },
      { kunci: "nilai", label: "Nilai" },
    ],
    [{ kode: "017", nilai: 12345 }],
  );
  expectMatch(xml, /<row r="1">.*Kode.*Nilai.*<\/row>/);
  // Kode "017" harus INLINE STRING (bukan <v>) — nol di depan TIDAK
  // boleh hilang seperti yang terjadi bila Excel membaca sebagai angka.
  expectMatch(xml, /<c r="A2" t="inlineStr"><is><t[^>]*>017<\/t><\/is><\/c>/);
  // Nilai numerik harus <v> ASLI — Excel dapat menghitungnya (SUM dst).
  expectMatch(xml, /<c r="B2"><v>12345<\/v><\/c>/);
});

test("rakitSheetXml — nol baris tetap menghasilkan header sah", () => {
  const xml = rakitSheetXml([{ kunci: "a", label: "A" }], []);
  expectMatch(xml, /<sheetData><row r="1">/);
  expectNoMatch(xml, /<row r="2">/);
});

test("bangunZip — arsip diawali signature local file header PK\\x03\\x04", () => {
  const zip = bangunZip([{ nama: "test.txt", data: new TextEncoder().encode("hi") }]);
  expectEqual(zip[0], 0x50); // 'P'
  expectEqual(zip[1], 0x4b); // 'K'
  expectEqual(zip[2], 0x03);
  expectEqual(zip[3], 0x04);
});

test("bangunZip — diakhiri signature end-of-central-directory PK\\x05\\x06", () => {
  const zip = bangunZip([{ nama: "a.txt", data: new TextEncoder().encode("x") }]);
  const n = zip.length;
  // EOCD minimal 22 byte di akhir arsip (tanpa comment).
  const eocd = zip.slice(n - 22, n - 18);
  expectEqual(eocd[0], 0x50);
  expectEqual(eocd[1], 0x4b);
  expectEqual(eocd[2], 0x05);
  expectEqual(eocd[3], 0x06);
});

test("bangunZip — pulang-pergi: local file header CRC cocok crc32 data", () => {
  const data = new TextEncoder().encode("konten uji pulang-pergi");
  const zip = bangunZip([{ nama: "f.txt", data }]);
  // CRC-32 di local file header berada pada offset 14 (4 byte, little-endian)
  // setelah signature(4)+version(2)+flag(2)+method(2)+time(2)+date(2).
  // tsconfig kami menyalakan noUncheckedIndexedAccess; indeksnya dijamin
  // ada oleh panjang header, dan itu dinyatakan alih-alih dipaksa.
  const bita = (i: number): number => zip[i] ?? 0;
  const crcDiArsip =
    bita(14) | (bita(15) << 8) | (bita(16) << 16) | (bita(17) << 24);
  expectEqual(crcDiArsip >>> 0, crc32(data));
});

test("rakitXlsxBytes — menghasilkan arsip ZIP non-kosong yang diawali PK", () => {
  const bytes = rakitXlsxBytes([{ kunci: "x", label: "X" }], [{ x: 1 }]);
  expectOk(bytes.length > 0);
  expectEqual(bytes[0], 0x50);
  expectEqual(bytes[1], 0x4b);
});

test("rakitXlsxBytes — memuat SELURUH 7 bagian wajib XLSX minimal", () => {
  const bytes = rakitXlsxBytes([{ kunci: "x", label: "X" }], [{ x: 1 }]);
  const teks = new TextDecoder("latin1").decode(bytes);
  for (const nama of [
    "[Content_Types].xml",
    "_rels/.rels",
    "xl/workbook.xml",
    "xl/_rels/workbook.xml.rels",
    "xl/worksheets/sheet1.xml",
    "docProps/core.xml",
    "docProps/app.xml",
  ]) {
    expectOk(teks.includes(nama), `bagian ${nama} tidak ditemukan di arsip`);
  }
});


// Jembatan node:assert -> vitest. Uji di-port APA ADANYA; yang berubah
// hanya pustaka penegasnya, supaya isinya dapat dibandingkan baris demi
// baris dengan berkas aslinya.
function expectEqual(a: unknown, b: unknown): void {
  expect(a).toEqual(b);
}
function expectOk(a: unknown, pesan?: string): void {
  expect(a, pesan).toBeTruthy();
}
function expectMatch(teks: string, pola: RegExp): void {
  expect(teks).toMatch(pola);
}
function expectNoMatch(teks: string, pola: RegExp): void {
  expect(teks).not.toMatch(pola);
}
