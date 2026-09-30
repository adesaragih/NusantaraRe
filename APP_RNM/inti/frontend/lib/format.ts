/**
 * Pemformat TAMPILAN — satu-satunya tempat nilai berubah menjadi teks layar.
 *
 * Modul ini murni: tidak ada fetch, tidak ada state, tidak ada React. Itu
 * disengaja. Angka yang tampil di layar adalah angka UANG, dan satu-satunya
 * cara memastikan ia tampil sama dengan sistem existing adalah dengan
 * mengujinya — yang mustahil selama pemformatnya bersarang di dalam klien
 * HTTP. Ujinya ada di format.test.ts.
 *
 * Nilai TIDAK PERNAH berubah di sini; yang berubah hanya bentuk bacanya.
 */

/**
 * Format tanggal SERAGAM: DD-MM-YYYY, tanpa jam (NFR-14).
 *
 * Sistem lama mencampur YYYYMMDD, RRRRMMDD, dan DDMMRRRR — bahkan antar
 * insert induk dan anak pada satu operasi. Itu sumber bug yang nyata.
 */
export function formatDate(v: string | Date | null | undefined): string {
  if (!v) return "";
  const d = typeof v === "string" ? parseDate(v) : v;
  if (!d || Number.isNaN(d.getTime())) return "";
  // Waktu NOL Go ("0001-01-01...") = "tidak ada nilai", bukan tanggal —
  // pasca-ADR-007 posisi alur tidak menyimpan waktu, dan menampilkannya
  // sebagai 01-01-0001 menyamarkan ketiadaan sebagai data.
  if (d.getFullYear() < 1900) return "";
  const p = (n: number) => String(n).padStart(2, "0");
  return `${p(d.getDate())}-${p(d.getMonth() + 1)}-${d.getFullYear()}`;
}

function parseDate(s: string): Date | null {
  const m = /^(\d{2})-(\d{2})-(\d{4})$/.exec(s);
  if (m) return new Date(Number(m[3]), Number(m[2]) - 1, Number(m[1]));
  // Dua bentuk nilai existing yang `new Date()` TIDAK kenali dan sebelumnya
  // jatuh diam-diam menjadi "-" (ronde 54):
  //   "20201031T170000.000 GMT"  — stempel Pega (TANGGALSTATEMENT dkk)
  //   "20200101"                 — tanggal padat (Commencement/Termination)
  const pega = /^(\d{4})(\d{2})(\d{2})(?:T\d{6})?/.exec(s.trim());
  if (pega)
    return new Date(Number(pega[1]), Number(pega[2]) - 1, Number(pega[3]));
  const d = new Date(s);
  return Number.isNaN(d.getTime()) ? null : d;
}

/**
 * Jumlah desimal TIDAK dibatasi — padanan `pyDecimalPlaces = -999` di Pega.
 *
 * Nilainya tampil dengan sebanyak digit yang backend kirim: tidak dibulatkan,
 * tidak dipangkas, tidak dipadankan dengan nol. Dipakai untuk `To IDR`,
 * `Limit (%)`, dan kolom `Rp`/`Usd` Arrangement (UI_EXISTING_REFERENSI §9.2).
 */
export const DESIMAL_TAK_DIBATASI = -1;

/**
 * Angka bergaya Indonesia: **titik ribuan, koma desimal** — 1.234.567,89.
 *
 * Keputusan pemilik proyek. Sumber Pega TIDAK memuat setelan pemisah ribuan
 * (UI_EXISTING_REFERENSI §9.1: rule kontrol pxNumber/pxCurrency dan locale
 * operator tidak ikut diekspor); satu-satunya petunjuk adalah judul grid
 * literal `Risk with Sum Insured less than USD 100.000.000` — titik sebagai
 * pemisah ribuan.
 *
 * `desimal` adalah BATAS ATAS, bukan panjang tetap. Nilainya ditiru per jenis:
 * 6 uang realisasi · 3 share/komisi · 2 limit/TSI · 1 Pct arrangement ·
 * DESIMAL_TAK_DIBATASI untuk Rp/Usd arrangement. NFR-13 merekonsiliasi hasil
 * TREATY V2 terhadap Pega, sehingga selisih tampilan akan terbaca sebagai
 * selisih data.
 *
 * # Nol di EKOR tidak ditampilkan (ronde 69)
 *
 * Versi sebelumnya memadankan pecahan ke panjang tetap, sehingga 2484250
 * tampil "2.484.250,000000" pada batas 6 dan setiap nilai bulat memperoleh
 * ekor nol. Pemilik proyek menolaknya dengan menunjuk sistem existing:
 *
 *	"DI PEGA 2484250 MAKA DI SISTEM BARU JANGAN 2484250.000000000"
 *
 * Yang dipangkas hanya nol yang TIDAK BERMAKNA — tidak satu digit berarti pun
 * hilang, dan nilai tersimpan tidak disentuh. Membulatkan akan menghilangkan
 * digit berarti; memangkas nol ekor tidak.
 *
 * # Mengapa string, bukan Number
 *
 * Uang TETAP string di API (money.Amount) dan pemformatan hanya lapis
 * tampilan. Melewatkannya melalui `Number` mengembalikan presisi ganda yang
 * baru saja dihindari backend: nilai 6 desimal berukuran besar kehilangan
 * digit terakhirnya secara diam-diam. Pembulatan di bawah karena itu bekerja
 * pada DIGIT, bukan pada float.
 *
 * Teks yang bukan angka dikembalikan APA ADANYA — lebih baik menampilkan
 * nilai mentah daripada berpura-pura ia nol.
 */
export function formatNumber(
  nilai: string | number | null | undefined,
  desimal: number,
): string {
  if (nilai === null || nilai === undefined) return "";
  const mentah = String(nilai).trim();
  if (mentah === "") return "";

  // Pemisah desimal yang DITERIMA hanya titik — itulah yang backend kirim.
  // Menerima koma juga berarti menebak: "1,234" bisa berarti 1,234 atau
  // 1234, dan menebak salah pada nilai uang.
  const m = /^([+-]?)(\d*)(?:\.(\d*))?$/.exec(mentah);
  if (!m) return mentah;
  const pecahanMentah = m[3] ?? "";
  if (m[2] === "" && pecahanMentah === "") return mentah;

  // ⚠️ Grup regex dibaca lewat `??`: grup yang TIDAK cocok memang tidak
  // ada, dan itu berbeda dari grup yang cocok dengan teks kosong.
  const tanda = m[1] === "-" ? "-" : "";
  const bulatMentah = m[2] ?? "";
  let bulat = bulatMentah === "" ? "0" : bulatMentah;
  let pecahan = pecahanMentah;
  if (desimal >= 0) [bulat, pecahan] = bulatkan(bulat, pecahan, desimal);

  // Nol di EKOR dibuang: ia tidak menambah satu pun informasi, dan dalam
  // konvensi Indonesia "2.484.250,000000" justru lebih sulit dibaca daripada
  // "2.484.250". Yang di DEPAN koma tidak disentuh.
  pecahan = pecahan.replace(/0+$/, "");

  // Nol di depan dibuang — ini angka, bukan pengenal (kolom pengenal tidak
  // ditandai isNumber justru karena nol di depannya bermakna).
  bulat = bulat.replace(/^0+(?=\d)/, "");
  const grup = bulat.replace(/\B(?=(\d{3})+(?!\d))/g, ".");
  return tanda + grup + (pecahan === "" ? "" : "," + pecahan);
}

/**
 * Isi satu sel tabel yang jenis kolomnya TIDAK dinyatakan backend.
 *
 * Tiga layar menampilkan kolom apa adanya dari sumbernya — Laporan Realisasi
 * (view existing), Limit Treaty In, dan grid XOL. Kolomnya dinamis, jadi tidak
 * ada penanda isNumber/decimals seperti pada grid master.
 *
 * Aturannya sengaja sempit: hanya nilai yang MENGANDUNG TITIK DESIMAL yang
 * diformat. Itu yang membedakan "2484250.000000000" (uang, perlu dirapikan)
 * dari "2026" (tahun) dan "1000809" (pengenal) — dan tanpa pembeda itu
 * UWYEAR akan tampil "2.026", persis cara kolom tahun rusak di layar.
 *
 * Nilai bergaya koma-legacy ("0,027") tidak lolos pemeriksaan dan tampil apa
 * adanya: menafsirkannya berarti menebak apakah komanya ribuan atau desimal.
 */
export function formatSel(nilai: string | null | undefined): string {
  const t = (nilai ?? "").trim();
  if (!/^[+-]?\d*\.\d+$/.test(t)) return t;
  return formatNumber(t, DESIMAL_TAK_DIBATASI);
}

/**
 * Tampilkan sebuah PERSENTASE — angka yang sama, ditambah tanda %.
 *
 * Pemangkasan nol ekor dikerjakan formatNumber (lihat catatan di sana), jadi
 * "30,000000" sudah tiba di sini sebagai "30". Yang tersisa hanya menempelkan
 * tandanya.
 *
 * Jumlah desimalnya TIDAK dikurangi paksa: memaksa dua desimal menampilkan
 * 33,333% sebagai "33,33%", dan tiga share 33,333 berjumlah TEPAT 100
 * sementara tiga share 33,33 tidak (SD-05, BR-01).
 */
export function formatPersen(
  nilai: string | number | null | undefined,
  desimalMaks: number,
): string {
  const teks = formatNumber(nilai, desimalMaks);
  return teks === "" ? "" : teks + "%";
}

/**
 * Pembulatan setengah-ke-atas pada digit, tanpa float.
 *
 * Mengembalikan [bagian bulat, bagian pecahan] dengan pecahan PALING BANYAK
 * `desimal` digit. Limpahan dari pembulatan merambat ke bagian bulat: 9,999
 * pada 2 desimal menjadi 10 — bukan 9,100.
 *
 * Pecahan yang lebih pendek dari batas dibiarkan apa adanya; memadankannya
 * dengan nol adalah persis kebiasaan yang ronde 69 buang.
 */
function bulatkan(
  bulat: string,
  pecahan: string,
  desimal: number,
): [string, string] {
  if (pecahan.length <= desimal) return [bulat, pecahan];

  const simpan = pecahan.slice(0, desimal);
  if (pecahan.charCodeAt(desimal) - 48 < 5) return [bulat, simpan];

  const digit = (bulat + simpan).split("");
  let i = digit.length - 1;
  let sisa = 1;
  while (i >= 0 && sisa === 1) {
    // ⚠️ `i >= 0` sudah dijaga gelungnya, tetapi tsc tidak dapat
    // membuktikan indeksnya di dalam rentang; nilainya diambil sekali
    // ke variabel yang bertipe pasti.
    const dig = digit[i] ?? "0";
    const d = dig.charCodeAt(0) - 48 + 1;
    digit[i] = String(d % 10);
    sisa = d === 10 ? 1 : 0;
    i--;
  }
  const s = (sisa === 1 ? "1" : "") + digit.join("");
  const potong = s.length - desimal;
  return [s.slice(0, potong) || "0", s.slice(potong)];
}

/**
 * Petunjuk isian — menyebut ARTI nilainya, bukan konvensi pemisah desimalnya.
 *
 * Pemilik proyek menolak "angka, desimal pakai titik" pada field persentase
 * (ronde 72): field itu persen, dan yang perlu diketahui pengguna adalah
 * satuannya. Keduanya dipusatkan di sini karena sebelumnya tersalin di tiga
 * berkas dengan bunyi yang sama.
 */
export const PETUNJUK_ANGKA = "Isi dengan angka";
export const PETUNJUK_PERSEN = "Isi dengan angka";
