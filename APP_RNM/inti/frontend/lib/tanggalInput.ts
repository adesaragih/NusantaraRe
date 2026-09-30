/**
 * Konversi tanggal di BATAS kontrol input.
 *
 * `<input type="date">` hanya mengerti satu bentuk: `YYYY-MM-DD`. Aplikasi
 * ini memakai bentuk lain, dan bentuk-bentuk itu TIDAK berubah:
 *
 *	kabel API    DD-MM-YYYY   (httpx.Date, NFR-14)
 *	tampilan     DD-MM-YYYY   (lib/format.ts, NFR-14)
 *	simpanan DB  YYYYMMDD / stempel Pega — urusan backend, tidak disentuh
 *
 * Yang berubah HANYA kontrol input-nya. Konversinya hidup di sini — satu
 * tempat, teruji — bukan disalin ke setiap form. Konversi tanggal yang
 * disalin akan berbeda suatu hari, dan yang berbeda adalah TANGGAL: bidang
 * yang selisih satu harinya tidak terlihat sampai seseorang merekonsiliasi.
 *
 * # Kenapa TIDAK memakai `new Date()`
 *
 * `new Date("2020-01-01")` di JavaScript diperlakukan sebagai UTC, lalu
 * `getDate()` membacanya kembali dalam zona waktu lokal. Di WIB (UTC+7) itu
 * kebetulan aman, tetapi di zona negatif ia menggeser satu hari — dan
 * pergeseran pada 1 Januari menggeser TAHUN.
 *
 * Karena itu seluruh konversi di berkas ini murni MANIPULASI TEKS. Tidak ada
 * satu pun objek `Date` yang dibuat, sehingga tidak ada zona waktu yang
 * dapat ikut campur.
 */

/** Memeriksa kewajaran tanggal, termasuk hari yang melebihi panjang bulan. */
function tanggalWajar(th: number, bl: number, hr: number): boolean {
  if (th < 1900 || th > 9999) return false;
  if (bl < 1 || bl > 12) return false;
  if (hr < 1) return false;
  const kabisat = (th % 4 === 0 && th % 100 !== 0) || th % 400 === 0;
  const panjang = [31, kabisat ? 29 : 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];
  // ⚠️ `bl` sudah dijaga 1..12 di atas, jadi indeksnya pasti ada. Nilai
  // cadangan 0 membuat tanggal apa pun DITOLAK bila penjaga itu kelak
  // dilonggarkan - gagal tertutup, bukan gagal terbuka.
  return hr <= (panjang[bl - 1] ?? 0);
}

const p2 = (n: string) => n.padStart(2, "0");

/**
 * keInputTanggal mengubah nilai simpanan/kabel menjadi `YYYY-MM-DD`.
 *
 * Bentuk masuknya SENGAJA longgar. Data existing bercampur — `formatDate`
 * sudah menghadapi `20200101` dan `20201031T170000.000 GMT` di layar.
 * Menolak keduanya di sini membuat tanggal yang TERSIMPAN DENGAN BENAR
 * tampil kosong di form, dan pengguna yang menyimpannya kembali akan
 * menghapus tanggal yang sebenarnya ada.
 *
 * Nilai yang tidak terbaca menghasilkan string kosong, bukan tanggal
 * karangan: input kosong terlihat sebagai "belum diisi", sementara tanggal
 * karangan terlihat sebagai data.
 */
export function keInputTanggal(v: string | null | undefined): string {
  const s = (v ?? "").trim();
  if (s === "") return "";

  let th = "";
  let bl = "";
  let hr = "";

  // DD-MM-YYYY — bentuk kabel dan tampilan.
  let m = /^(\d{2})-(\d{2})-(\d{4})$/.exec(s);
  if (m) {
    [hr, bl, th] = [m[1] ?? "", m[2] ?? "", m[3] ?? ""];
  }

  // YYYY-MM-DD — sudah benar; tetap divalidasi supaya "2026-13-45" tidak
  // lolos hanya karena bentuknya kebetulan cocok.
  if (!th) {
    m = /^(\d{4})-(\d{2})-(\d{2})(?:[T ].*)?$/.exec(s);
    if (m) [th, bl, hr] = [m[1] ?? "", m[2] ?? "", m[3] ?? ""];
  }

  // YYYYMMDD padat, dengan atau tanpa stempel jam Pega.
  if (!th) {
    m = /^(\d{4})(\d{2})(\d{2})(?:T\d{6}.*)?$/.exec(s);
    if (m) [th, bl, hr] = [m[1] ?? "", m[2] ?? "", m[3] ?? ""];
  }

  if (!th) return "";
  if (!tanggalWajar(Number(th), Number(bl), Number(hr))) return "";
  // Waktu NOL Go ("0001-01-01") sudah tersaring `tanggalWajar` lewat batas
  // tahun 1900 — ambang yang sama dengan lib/format.ts, supaya form dan
  // tabel tidak berbeda pendapat tentang apa yang disebut "tidak ada nilai".
  return `${th}-${p2(bl)}-${p2(hr)}`;
}

/**
 * dariInputTanggal mengubah `YYYY-MM-DD` menjadi bentuk KABEL `DD-MM-YYYY`.
 *
 * Bentuk keluarnya hanya SATU, karena itulah yang kontrak API terima
 * (httpx.Date, NFR-14). Melebarkannya berarti mengirim bentuk yang backend
 * tidak janjikan akan dibaca.
 *
 * Aturan wajib-isi dan rentang TIDAK digandakan di sini: `offer.ParseDate`
 * di backend tetap otoritasnya. Yang ditolak di sini hanya nilai yang bukan
 * tanggal sama sekali — supaya teks sampah tidak menyeberang sebagai
 * tanggal.
 */
export function dariInputTanggal(v: string | null | undefined): string {
  const s = (v ?? "").trim();
  if (s === "") return "";
  const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(s);
  if (!m) return "";
  if (!tanggalWajar(Number(m[1]), Number(m[2]), Number(m[3]))) return "";
  return `${m[3]}-${m[2]}-${m[1]}`;
}
