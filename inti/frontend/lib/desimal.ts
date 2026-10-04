/**
 * Aritmetika desimal EKSAK di atas string.
 *
 * Uang dan persentase datang dari backend sebagai STRING (money.Amount /
 * money.Factor) justru supaya presisinya tidak pernah lewat float64.
 * Menjumlahkannya dengan `Number` membuang jaminan itu di lapis terakhir:
 * total berdesimal panjang kehilangan digit belakangnya tanpa suara — dan
 * digit belakang itulah yang menentukan apakah total TEPAT 100 (BR-01
 * share, BR-02 cicilan).
 *
 * Karena itu penjumlahan di sini bekerja pada DIGIT lewat BigInt, bukan pada
 * float. Berkas ini tidak memuat satu pun `Number(...)` atau `parseFloat`.
 */

/** Bentuk desimal yang diterima: titik sebagai pemisah desimal, sesuai kabel. */
const BENTUK = /^[+-]?(?:\d+(?:\.\d*)?|\.\d+)$/;

/** Apakah teks ini nilai desimal yang dapat dihitung. */
export function desimalSah(teks: string): boolean {
  return BENTUK.test(teks.trim());
}

function panjangPecahan(teks: string): number {
  const i = teks.indexOf(".");
  return i < 0 ? 0 : teks.length - i - 1;
}

/** Ubah "12.34" pada skala 4 menjadi BigInt(123400). */
function keBulat(teks: string, skala: number): bigint {
  const negatif = teks.startsWith("-");
  const angka = teks.replace(/^[+-]/, "");
  const titik = angka.indexOf(".");
  const bulat = titik < 0 ? angka : angka.slice(0, titik);
  const pecahan = titik < 0 ? "" : angka.slice(titik + 1);
  const n = BigInt((bulat || "0") + pecahan.padEnd(skala, "0"));
  return negatif ? -n : n;
}

function keTeks(n: bigint, skala: number): string {
  const negatif = n < 0n;
  const digit = (negatif ? -n : n).toString().padStart(skala + 1, "0");
  const potong = digit.length - skala;
  const bulat = digit.slice(0, potong);
  const pecahan = skala === 0 ? "" : digit.slice(potong);
  return (negatif ? "-" : "") + bulat + (pecahan === "" ? "" : "." + pecahan);
}

/** Hasil penjumlahan: totalnya, beserta isian yang TIDAK dapat dibaca. */
export interface Jumlah {
  total: string;
  takTerbaca: string[];
}

/**
 * Jumlahkan nilai desimal secara eksak.
 *
 * Isian yang bukan desimal TIDAK dianggap nol — ia dikembalikan lewat
 * `takTerbaca`, supaya layar dapat mengatakan totalnya belum dapat
 * dipastikan alih-alih menampilkan angka yang menelan isian rusak. Isian
 * KOSONG dilewati: kosong berarti belum diisi, bukan nilai yang rusak.
 */
export function jumlahDesimal(nilai: string[]): Jumlah {
  const bersih: string[] = [];
  const takTerbaca: string[] = [];
  for (const v of nilai) {
    const t = v.trim();
    if (t === "") continue;
    if (desimalSah(t)) bersih.push(t);
    else takTerbaca.push(v);
  }
  const skala = bersih.reduce((m, t) => Math.max(m, panjangPecahan(t)), 0);
  let jumlah = 0n;
  for (const t of bersih) jumlah += keBulat(t, skala);
  return { total: keTeks(jumlah, skala), takTerbaca };
}

/**
 * Bandingkan dua nilai desimal secara EKSAK, per digit.
 *
 * "100" dan "100.000000000" dinyatakan sama; "99.999999999" TIDAK. Itulah
 * perbandingan yang dilakukan backend, dan menyamakannya di sini membuat
 * penanda di layar sejalan dengan penerimaan/penolakan di server.
 */
export function samaDengan(a: string, b: string): boolean {
  if (!desimalSah(a) || !desimalSah(b)) return false;
  const skala = Math.max(panjangPecahan(a), panjangPecahan(b));
  return keBulat(a.trim(), skala) === keBulat(b.trim(), skala);
}

/**
 * Geser titik desimal `langkah` posisi — perkalian/pembagian dengan 10^n
 * yang EKSAK.
 *
 * Dipakai untuk berpindah antara pecahan-dari-satu dan persentase: `.6` yang
 * tersimpan tampil 60%, dan 60% yang diketik tersimpan sebagai 0.6.
 *
 * # Kenapa bukan `Number(v) * 100`
 *
 * Perpindahan ini menyentuh nilai yang DIPERSIST, bukan sekadar tampilan.
 * Jalur float64 yang dipakai sebelumnya merusaknya dalam tiga cara sekaligus
 * (audit ronde 69): "0.00000000005" dibulatkan menjadi 0, "1e-7" berubah
 * menjadi "1e-9" yang backend tolak, dan "0x10" terbaca sebagai 16 lalu
 * tersimpan sebagai 0,16. Menggeser titik pada DIGIT tidak punya satu pun
 * dari ketiga mode gagal itu.
 *
 * Teks yang bukan desimal dikembalikan APA ADANYA — layar tidak menebak, dan
 * backend yang menolaknya dengan pesan yang tepat.
 */
export function geserTitik(teks: string, langkah: number): string {
  const t = teks.trim();
  if (!desimalSah(t) || langkah === 0) return teks;

  const negatif = t.startsWith("-");
  const angka = t.replace(/^[+-]/, "");
  const titik = angka.indexOf(".");
  const bulat = titik < 0 ? angka : angka.slice(0, titik);
  const pecahan = titik < 0 ? "" : angka.slice(titik + 1);

  // Seluruh digit dirangkai; yang berpindah hanyalah letak titiknya. Nol
  // ditambahkan di sisi yang kekurangan digit — itu satu-satunya cara
  // pergeseran dapat melampaui panjang angka tanpa kehilangan digit.
  const digit = bulat + pecahan;
  const letak = bulat.length + langkah;

  let kiri: string;
  let kanan: string;
  if (letak <= 0) {
    kiri = "0";
    kanan = "0".repeat(-letak) + digit;
  } else if (letak >= digit.length) {
    kiri = digit + "0".repeat(letak - digit.length);
    kanan = "";
  } else {
    kiri = digit.slice(0, letak);
    kanan = digit.slice(letak);
  }

  const bulatBaru = kiri.replace(/^0+(?=\d)/, "") || "0";
  const pecahanBaru = kanan.replace(/0+$/, "");
  const hasil = bulatBaru + (pecahanBaru === "" ? "" : "." + pecahanBaru);
  // Nol tidak bertanda: "-0" adalah nol yang terbaca seperti nilai negatif.
  return (negatif && /[1-9]/.test(hasil) ? "-" : "") + hasil;
}

/**
 * Pecahan-dari-satu <-> persentase, EKSAK: yang berpindah hanya titik
 * desimalnya (`geserTitik`). Nilai ini DIPERSIST, jadi jalur float64 yang
 * dipakai sebelumnya bukan sekadar tidak rapi — ia mengubah data.
 *
 * Dipakai bersama oleh tampilan grid (`App.tsx` `selDisplay`) dan formulir
 * (`components/RecordForm.tsx` `bentukIsian`) — satu sumber, bukan disalin.
 */
export function pecahanDariPersen(v: string): string {
  return v.trim() === "" ? v : geserTitik(v, -2);
}

export function persenDariPecahan(v: string | undefined): string | undefined {
  return v === undefined || v.trim() === "" ? v : geserTitik(v, 2);
}

/**
 * Buang nol di ekor pecahan, menyisakan paling sedikit `minimal` digit.
 *
 * Ini BUKAN pembulatan: nol di ekor tidak membawa informasi, sehingga
 * membuangnya tidak mengubah nilai. Gunanya supaya total share yang
 * backend kirim sembilan desimal ("100.000000000") tetap terbaca sebagai
 * 100,000 tanpa perlu dibulatkan — pembulatan justru dapat menampilkan
 * "100,000" untuk 99,999999999 dan menyembunyikan pelanggaran BR-01.
 */
export function pangkasNolEkor(teks: string, minimal: number): string {
  if (!desimalSah(teks)) return teks;
  const t = teks.trim();
  if (!t.includes(".")) return t;
  // ⚠️ Dipotong lewat letak titik, bukan `split`, supaya tsc melihat kedua
  // potongan sebagai `string`: `t.includes(".")` di atas sudah menjamin
  // keduanya ada, tetapi destructuring array tidak menyampaikannya.
  const titik = t.indexOf(".");
  const bulat = t.slice(0, titik);
  const pecahan = t.slice(titik + 1);
  let akhir = pecahan.length;
  while (akhir > minimal && pecahan[akhir - 1] === "0") akhir--;
  const sisa = pecahan.slice(0, akhir);
  return sisa === "" ? bulat : bulat + "." + sisa;
}
