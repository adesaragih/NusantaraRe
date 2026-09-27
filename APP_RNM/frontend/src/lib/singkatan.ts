/**
 * Lencana dua huruf untuk daftar menu sidebar — bebas tabrakan.
 *
 * # Kenapa ia menerima DAFTAR, bukan satu nama
 *
 * Sampai ronde 248, `App.tsx` memakai `singkatan(nama)` yang membaca DUA
 * KATA PERTAMA saja. Fungsi yang hanya melihat satu nama tidak akan pernah
 * tahu nama lain memilih huruf yang sama, dan hasilnya TIGA tabrakan pada
 * daftar menu yang sungguhan:
 *
 *	TE -> Treaty Exchange · Treaty Exchange Yearly · Treaty Excess of Loss
 *	TO -> Table of Limit  · Treaty Out Non Proportional
 *	RT -> Reinsurance Type · Realisasi Treaty In
 *
 * Pemilik proyek melaporkan yang pertama ("DUA TE kembar"); dua sisanya
 * hanya terlihat dari sapuan seluruh daftar. Karena tabrakan adalah sifat
 * KUMPULAN dan bukan sifat satu nama, seam-nya ikut berpindah ke kumpulan.
 *
 * # Yang TIDAK dilakukan fungsi ini
 *
 * Ia tidak menyentuh nama menu sama sekali — tidak memotong, tidak
 * menormalkan, tidak menerjemahkan. Nama menu wajib sama dengan existing
 * (pagar §4.1 ronde 248), dan lencana hanyalah hiasan di sebelahnya. PEGA
 * sendiri TIDAK punya lencana (`MasterNavTreaty.xml`: `pyImageSource=none`
 * 18/18), jadi bentuk lencana adalah SD-14 murni — tidak ada angka maupun
 * keputusan bisnis yang berubah karenanya.
 */

/** Huruf/angka pertama sebuah kata, atau "" bila kata itu tanpa keduanya. */
function hurufAwal(kata: string): string {
  const m = /[A-Za-z0-9]/.exec(kata);
  return m ? m[0].toUpperCase() : "";
}

/**
 * Calon lencana untuk satu nama, URUT dari yang paling disukai.
 *
 * Urutannya sengaja menaruh bentuk yang paling terbaca di depan, sehingga
 * nama yang TIDAK bertabrakan tetap mendapat lencana yang sama seperti
 * sebelum ronde 248 — daftar yang sudah dihafal orang tidak ikut berubah
 * hanya karena mekanismenya diganti.
 */
function calon(nama: string): string[] {
  const kata = nama
    .trim()
    .split(/\s+/)
    .map(hurufAwal)
    .filter(Boolean);
  const bersih = nama.replace(/[^A-Za-z0-9]/g, "").toUpperCase();
  const out: string[] = [];
  const tambah = (s: string) => {
    if (/^[A-Z0-9]{2}$/.test(s) && !out.includes(s)) out.push(s);
  };

  if (kata.length === 0) {
    // Nama tanpa huruf maupun angka sama sekali — tidak ada yang dapat
    // disingkat darinya, dan penyelesaiannya memang diserahkan ke gelung
    // abjad di tingkat 5 (lewat `awal` yang jatuh ke "X").
    //
    // KOREKSI TV2-299 (ronde 250): cabang ini dulu berbunyi
    // `tambah("--")`, yang TIDAK PERNAH menambah apa pun — `tambah`
    // menyaring /^[A-Z0-9]{2}$/ dan "--" gagal saringan itu. Kode mati
    // yang tampak seperti penanganan kasus; menghapusnya tidak mengubah
    // satu pun hasil, dan komentar inilah penanganannya yang sebenarnya.
  } else if (kata.length === 1) {
    tambah(bersih.slice(0, 2));
  } else {
    // ⚠️ Cabang ini hanya berjalan bila `kata.length >= 2`, jadi kata
    // pertama dan kedua PASTI ada - tetapi tsc tidak melihatnya dari
    // `else`. Keduanya diambil sekali ke variabel bertipe pasti.
    const k0 = kata[0] ?? "";
    const k1 = kata[1] ?? "";
    const kAkhir = kata[kata.length - 1] ?? "";
    // 1. dua kata pertama — bentuk LAMA, dipertahankan sebagai pilihan utama
    tambah(k0 + k1);
    // 2. kata pertama + kata TERAKHIR — "Treaty Exchange Yearly" -> TY
    tambah(k0 + kAkhir);
    // 3. kata pertama + tiap kata di tengah, kiri ke kanan
    for (let i = 2; i < kata.length - 1; i++) tambah(k0 + (kata[i] ?? ""));
  }
  // 4. huruf pertama + tiap huruf berikutnya dari nama itu sendiri
  const awal = kata[0] ?? bersih[0] ?? "X";
  for (const h of bersih.slice(1)) tambah(awal + h);
  // 5. abjad lalu angka, semuanya berawalan huruf pertama nama.
  //    Tingkat 4 menghasilkan bagian dari himpunan yang sama; ia tetap
  //    ada karena URUTANNYA berbeda — huruf yang benar-benar muncul di
  //    nama itu lebih dulu ditawarkan daripada huruf sembarang.
  //
  //    Ini BUKAN jalan terakhir. Kalau ke-36 kombinasi ini pun terpakai,
  //    `lencanaSisa` menyapu seluruh ruang 36x36 (TV2-299).
  for (const h of AKSARA) tambah(awal + h);
  return out;
}

/**
 * Memetakan setiap nama ke lencana dua huruf yang UNIK dalam daftar itu.
 *
 * Deterministik: daftar yang sama selalu menghasilkan pemetaan yang sama,
 * karena nama diproses menurut urutannya dan calon dipilih menurut urutan
 * tetap. Lencana yang berubah tiap render adalah lencana yang tidak dapat
 * dihafal.
 *
 * Nama KEMBAR dalam satu daftar berbagi satu lencana — ia nama yang sama,
 * bukan dua menu. Kuncinya adalah nama itu apa adanya.
 */
export function singkatanUnik(daftar: string[]): Map<string, string> {
  const hasil = new Map<string, string>();
  const dipakai = new Set<string>();
  for (const nama of daftar) {
    if (hasil.has(nama)) continue;
    let pilih = "";
    for (const c of calon(nama)) {
      if (!dipakai.has(c)) {
        pilih = c;
        break;
      }
    }
    if (!pilih) pilih = lencanaSisa(dipakai, nama);
    hasil.set(nama, pilih);
    dipakai.add(pilih);
  }
  return hasil;
}

/** Seluruh aksara yang boleh muncul di lencana, urut tetap. */
const AKSARA = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

/**
 * Lencana apa pun yang belum terpakai, dicari di SELURUH ruang dua
 * aksara (36 x 36 = 1.296 kemungkinan).
 *
 * # KOREKSI TV2-299 (ronde 250)
 *
 * Jalan terakhir di sini dulu berbunyi `pilih = "--"` — dan ia
 * membatalkan justru jaminan yang fungsi ini ada untuk memberi. Dua
 * nama yang sama-sama kehabisan calon akan sama-sama menerima `"--"`,
 * yaitu lencana KEMBAR: persis cacat yang dilaporkan pemilik proyek
 * ("DUA TE kembar") dan yang seluruh berkas ini dibangun untuk
 * menutupnya. Ia juga melanggar bentuk yang dijanjikan
 * (`/^[A-Z0-9]{2}$/`), karena ia satu-satunya jalur yang tidak melewati
 * `tambah`.
 *
 * Terjangkaunya nyata, bukan teoretis: `calon()` hanya menawarkan
 * kombinasi yang berawalan huruf pertama nama itu — paling banyak 36.
 * Tiga puluh tujuh menu yang berawalan huruf sama sudah cukup.
 * Terbukti dengan 40 nama "Treaty Nomor N" (uji ronde 250).
 *
 * Ruang penuh 1.296 tidak dapat dihabiskan oleh daftar menu mana pun
 * yang masuk akal. Bila toh habis, fungsi ini MELEMPAR galat alih-alih
 * mengembalikan lencana kembar — kontraknya keunikan, dan keunikan yang
 * gagal diam-diam adalah keadaan yang justru tidak boleh ada di sini.
 *
 * # Pilihan yang DITOLAK, dan kenapa (temuan code-review, ronde 251)
 *
 * Melempar berarti seluruh layar berhenti demi sebuah HIASAN. Itu
 * keberatan yang sah, dan jawabannya bukan "tidak mungkin terjadi"
 * melainkan: galatnya tertangkap `PagarGalat` di akar (`main.tsx`),
 * jadi yang muncul halaman galat berpesan, bukan layar putih.
 *
 * Alternatif KETIGA yang dipertimbangkan dan tidak dipakai: kembalikan
 * "tanpa lencana" dan biarkan layar hidup. Ia ditolak karena
 * memindahkan kegagalan ke tempat yang tidak ada yang melihatnya —
 * 22 menu tampil, dua di antaranya tanpa lencana, dan tidak ada yang
 * tahu sistem penamaannya sudah kehabisan ruang. Bila daftar menu suatu
 * hari benar-benar melewati 1.296, itu keadaan yang HARUS dilaporkan,
 * bukan dirapikan.
 */
function lencanaSisa(dipakai: Set<string>, nama: string): string {
  for (const a of AKSARA) {
    for (const b of AKSARA) {
      const s = a + b;
      if (!dipakai.has(s)) return s;
    }
  }
  throw new Error(
    `singkatanUnik: seluruh ${AKSARA.length * AKSARA.length} lencana dua ` +
      `aksara sudah terpakai; tidak ada yang tersisa untuk ${JSON.stringify(nama)}. ` +
      "Mengembalikan lencana kembar akan menyembunyikan keadaan ini, " +
      "dan lencana kembar adalah cacat yang fungsi ini ada untuk menutup.",
  );
}
