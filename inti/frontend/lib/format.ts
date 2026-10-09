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

// ---------------------------------------------------------------------
// ISIAN ANGKA — pemisah ribuan SAAT MENGETIK
// ---------------------------------------------------------------------
//
// ⛔ MENGAPA INI BUKAN `formatNumber` YANG DIPANGGIL TIAP KETUKAN.
// `modul/treatyinadjustment/frontend/komponen/KerangkaTab.tsx` pernah
// menuliskan sebabnya apa adanya:
//
//   "Nilai MENTAH ke kotak isian, tanpa format angka: memformat di tiap
//    ketukan membuat koma desimal mustahil diketik."
//
// Dan itu benar untuk `formatNumber`: ia MEMBULATKAN dan MEMBUANG NOL EKOR.
// Mengetik `12,` menjadi `12`; koma yang baru saja diketik hilang sebelum
// digit berikutnya sempat masuk. Mengetik `1,05` mustahil: `1,0` menjadi
// `1`.
//
// ⭐ Pemecahannya BUKAN membatalkan pemformatan, melainkan memformat HANYA
// BAGIAN BULAT dan membiarkan ekor desimal persis seperti yang diketik —
// termasuk koma sendirian dan nol di ekor.

/** Pemisah konvensi Indonesia: `.` ribuan, `,` desimal. */
const RIBUAN = ".";
const DESIMAL = ",";

/** Sisipkan titik tiap tiga digit dari kanan. */
function kelompokkan(bulat: string): string {
  let keluar = "";
  for (let i = 0; i < bulat.length; i++) {
    if (i > 0 && (bulat.length - i) % 3 === 0) keluar += RIBUAN;
    keluar += bulat[i];
  }
  return keluar;
}

/**
 * Bentuk TAMPIL saat mengetik: `1000000,5` → `1.000.000,5`.
 *
 * ⭐ EKOR DESIMAL LEWAT APA ADANYA. `12,` tetap `12,` dan `1,00` tetap
 * `1,00` — keduanya keadaan SAH di tengah pengetikan, dan keduanya hilang
 * bila `formatNumber` yang dipanggil.
 *
 * ⚠️ `desimalMaks` MEMOTONG ekor, tidak membulatkannya: yang mengetik digit
 * kesembilan di belakang koma sedang salah tekan, dan membulatkan diam-diam
 * mengubah angkanya. Pemotongan terlihat seketika di layar.
 *
 * Menerima titik MAUPUN koma sebagai pemisah desimal yang diketik — papan
 * tik angka banyak yang hanya punya titik.
 */
export function formatKetik(mentah: string, desimalMaks: number): string {
  const t = String(mentah ?? "").trim();
  if (t === "") return "";
  const minus = t.startsWith("-") ? "-" : "";
  const tanpaTanda = t.replace(/^[+-]/, "");

  // ⛔ TITIK SELALU PEMISAH RIBUAN, KOMA SELALU PEMISAH DESIMAL.
  //
  // Bentuk pertama menerima keduanya sebagai pemisah desimal, dan itu CACAT
  // yang dilaporkan pemilik proses 7 Oktober 2026 (*"mentok … tidak bisa
  // meng input lebih"*): kotak terkendali mengumpankan KELUARANNYA SENDIRI
  // kembali, jadi sesudah `1.000` ketukan berikutnya tiba sebagai `1.0005`.
  // Titik pemisah ribuan itu terbaca sebagai pemisah desimal, ekornya
  // dipotong dua digit, dan angkanya kembali ke `1,00` — setiap ketukan,
  // selamanya.
  //
  // ⭐ Titik yang benar-benar DIKETIK diterjemahkan menjadi koma di batas
  // masukan oleh `normalisasiKetikan`, yang tahu huruf mana yang baru
  // disisipkan. Di sini nol tebakan.
  // ⛔ Kolom BULAT: yang diketik sesudah koma DIBUANG, bukan disambung.
  // Menyambungnya mengubah `1000,5` menjadi `10005` — keliru sepuluh kali
  // lipat, tanpa satu pun tanda di layar.
  if (desimalMaks === 0) {
    const batas = tanpaTanda.indexOf(DESIMAL);
    const sumber = batas >= 0 ? tanpaTanda.slice(0, batas) : tanpaTanda;
    const angka = sumber.replace(/\D/g, "").replace(/^0+(?=\d)/, "");
    return angka === "" ? "" : minus + kelompokkan(angka);
  }

  const iKoma = tanpaTanda.lastIndexOf(DESIMAL);
  const adaPemisah = iKoma >= 0;

  const bulatMentah = (adaPemisah ? tanpaTanda.slice(0, iKoma) : tanpaTanda).replace(/\D/g, "");
  let ekor = adaPemisah ? tanpaTanda.slice(iKoma + 1).replace(/\D/g, "") : "";
  if (desimalMaks >= 0 && ekor.length > desimalMaks) ekor = ekor.slice(0, desimalMaks);

  // Nol di depan dibuang, tetapi `0` sendirian dipertahankan — yang mengetik
  // `0,5` melewati keadaan `0`, dan membuangnya memakan angkanya.
  const bulat = bulatMentah.replace(/^0+(?=\d)/, "");
  const kiri = bulat === "" ? (adaPemisah ? "0" : "") : kelompokkan(bulat);
  if (!adaPemisah) return minus + kiri;
  return minus + kiri + DESIMAL + ekor;
}

/**
 * Terjemahkan TITIK YANG BARU DIKETIK menjadi koma.
 *
 * ⛔ Mengapa membandingkan dengan teks sebelumnya, bukan menebak dari
 * bentuknya: `1.0005` dapat berarti "satu koma nol nol nol lima" ATAU
 * "seribu, lalu digit kelima diketik". Bentuknya SAMA; yang membedakan
 * hanya apa yang berubah. Jadi yang diperiksa adalah huruf yang disisipkan.
 *
 * ⭐ Gunanya papan tik angka: banyak yang hanya punya titik, dan yang
 * mengetiknya bermaksud desimal.
 */
export function normalisasiKetikan(baru: string, lama: string): string {
  if (baru.length !== lama.length + 1) return baru;
  let i = 0;
  while (i < lama.length && baru[i] === lama[i]) i++;
  return baru[i] === "." ? baru.slice(0, i) + DESIMAL + baru.slice(i + 1) : baru;
}

/**
 * Bentuk KABEL dari bentuk tampil: `1.000.000,5` → `1000000.5`.
 *
 * ⛔ Backend hanya menerima titik sebagai pemisah desimal; `formatNumber`
 * menyatakannya di kepalanya. Kebalikan `formatKetik`, dan keduanya WAJIB
 * tetap sepasang — uji `format.test.ts` mengadu bolak-baliknya.
 */
export function keKabelAngka(tampil: string): string {
  const t = String(tampil ?? "").trim();
  if (t === "") return "";
  const minus = t.startsWith("-") ? "-" : "";
  const tanpaTanda = t.replace(/^[+-]/, "");
  const iKoma = tanpaTanda.lastIndexOf(DESIMAL);
  if (iKoma < 0) return minus + tanpaTanda.replace(/\D/g, "");
  const bulat = tanpaTanda.slice(0, iKoma).replace(/\D/g, "");
  const ekor = tanpaTanda.slice(iKoma + 1).replace(/\D/g, "");
  // ⚠️ Koma TANPA digit di belakangnya bukan angka yang sah di kabel —
  // `12,` dikirim sebagai `12`. Yang mengetik belum selesai, dan layar
  // tetap memperlihatkan komanya.
  return ekor === "" ? minus + bulat : minus + bulat + "." + ekor;
}
