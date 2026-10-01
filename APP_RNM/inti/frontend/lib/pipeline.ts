// Disalin UTUH dari REFERENSI_UI/frontend/src/lib/pipeline.ts (latar layar login
// ronde 277, "Pipeline") untuk halaman login M_LOGIN_GO - keputusan work owner
// 01-10-2026. Satu penyesuaian: setiap pembacaan `pipa[i]` diberi `?? 0`, sebab
// tsconfig APP_RNM menyalakan noUncheckedIndexedAccess. Selebihnya sama persis.

/**
 * Latar bergerak layar Login — "Pipeline".
 *
 * Diturunkan dari AmbientCanvasBackgrounds karya **Sean Free** (Codrops,
 * 2018) — https://tympanus.net/codrops/ , demo `index5.html` /
 * `js/pipeline.js` + `js/util.js`.
 *
 * Lisensinya mengizinkan pemakaian bila DIBANGUN DI ATASNYA pada web app
 * (README.md sumber): *"This resource can be used freely if integrated or
 * build upon in personal or commercial projects such as websites, web apps
 * …"*. Kredit ini dipasang karena biayanya nol dan README yang sama meminta
 * *"a visible mention and link to the original work"* untuk karya turunan.
 *
 * ---------------------------------------------------------------------------
 * APA YANG BERBEDA DARI SUMBERNYA, DAN KENAPA
 * ---------------------------------------------------------------------------
 *
 * Ini BUKAN salinan. Empat perbedaan, seluruhnya dituntut §0.1:
 *
 * 1. **Modul, bukan skrip global.** Sumbernya `'use strict'` dengan ~16
 *    `const` tingkat atas (`TO_RAD`, `rand`, `lerp`, …) yang mencemari
 *    scope global begitu dimuat. Di sini SELURUH keadaan hidup di dalam
 *    closure `mulaiPipeline`, dan yang keluar hanya fungsi penghentinya.
 *
 * 2. **Dapat DIHENTIKAN.** Sumbernya memanggil
 *    `window.requestAnimationFrame(draw)` tanpa menyimpan id-nya sama
 *    sekali — loop-nya tidak dapat dihentikan siapa pun. Layar Login sering
 *    ditinggal terbuka berjam-jam; loop yang tidak pernah berhenti adalah
 *    CPU yang terbakar tanpa ada yang melihat layarnya (§0.1.1).
 *
 * 3. **Jejaknya MEMUDAR.** Ini perbedaan yang paling penting, dan ia bukan
 *    selera.
 *
 *    Sumbernya menggambar pipa ke `canvas.a` dan **TIDAK PERNAH
 *    membersihkannya**. Setiap sapuan menambah alpha 0.125 di atas yang
 *    sebelumnya, selamanya. Pada demo Codrops itu tidak terlihat karena
 *    tidak ada yang membiarkannya terbuka berjam-jam — di layar Login,
 *    latar belakangnya akan TERUS MENERANG tanpa batas.
 *
 *    Akibatnya terukur, bukan diduga: `.login__foot` adalah teks PUTIH di
 *    LUAR kartu, jadi kontrasnya diukur terhadap titik paling terang di
 *    layar. Latar yang menerang tanpa batas berarti rasio kontras yang
 *    MENURUN tanpa batas — dan ronde 264 baru saja memperbaiki teks itu
 *    dari 3,18:1 (di bawah AA).
 *
 *    Karena itu `canvas.a` dipudarkan tiap bingkai. Jejaknya tetap ada,
 *    tetapi terangnya BERBATAS — dan batas itulah yang membuat §0.1.3 dapat
 *    dijawab dengan angka.
 *
 * 4. **Warnanya dari palet aplikasi**, dibaca dari token CSS saat jalan
 *    (`--navy`), bukan `hsla(150,80%,1%,1)` nyaris hitam milik sumbernya.
 *    Nol warna baru diperkenalkan berkas ini (§8.8).
 *
 * ---------------------------------------------------------------------------
 * CACAT SUMBER YANG SENGAJA TIDAK DIPERBAIKI
 * ---------------------------------------------------------------------------
 *
 * `checkBounds` pada sumbernya menerima `x`/`y` sebagai nilai dan
 * menugaskan ulang parameter lokalnya — SESUDAH keduanya ditulis kembali ke
 * `pipeProps`. Jadi pembungkusan tepi layarnya tidak pernah benar-benar
 * terjadi; pipa hanyut keluar layar lalu mati oleh TTL-nya.
 *
 * Dipertahankan apa adanya. Memperbaikinya mengubah GERAK yang pemilik
 * proyek setujui dari tangkapan layar, dan "membetulkan" hal yang tidak
 * dikeluhkan siapa pun bukan bagian tahap ini. Dicatat supaya pembaca
 * berikutnya tahu ia disengaja, bukan terlewat.
 */

/** Sebuah animasi yang sedang berjalan. Panggil untuk menghentikannya. */
export type Penghenti = () => void;

/**
 * Parameter bentuk, disalin apa adanya dari sumbernya kecuali yang
 * disebutkan. Dipisah dari kode supaya nilainya dapat dibaca — dan diuji —
 * tanpa menjalankan kanvas.
 */
/**
 * Nama token palet. Nilainya TIDAK ditulis di modul ini.
 *
 * Pola yang sama dipakai `closingPlasma.ts` dan `webglLiquid.ts`, dan
 * alasannya sama: warna yang ditulis di dalam modul adalah warna yang
 * tidak seorang pun putuskan. Token dideklarasikan pada `.login`, jadi
 * pembacanya WAJIB `getComputedStyle(kanvas)` — kanvasnya anak `.login`.
 * `getComputedStyle(document.documentElement)` TIDAK akan melihatnya.
 */
export const TOKEN_PIPA = {
  /**
   * ⚠️ `--plasma-a` -> `--login-latar` — RONDE 279 B.1.
   *
   * Keputusan pemilik proyek 24-09-2026: latar Login `#1F2555`. Ia
   * TIDAK dipasang dengan mengubah nilai `--plasma-a`, dan itu
   * disengaja: `closingPlasma.ts` masih membaca token itu, dan modul
   * yang sedang DITAHAN (pola ronde 270 D.9) harus tetap dapat
   * dijalankan apa adanya bila latar ini kelak ditolak.
   *
   * Token baru berarti kedua modul membaca warnanya masing-masing, dan
   * tidak ada yang berubah diam-diam untuk modul yang tidak diminta
   * berubah.
   */
  latar: "--login-latar",
  /**
   * ⚠️ `--cair-tengah` -> `--cair-sorot` — RONDE 281 B.2.
   *
   * Keputusan pemilik proyek 24-09-2026 (pesan malam): *"JEJAK-JEJAKNYA
   * DITERANGIN LAGI"*. Ini menjawab pilihan A/B/C yang ronde 279 B
   * ajukan atas lima jalur samar, dan ia menggeser plafon "tidak norak"
   * DENGAN KEPUTUSAN — bukan diam-diam.
   *
   * Tuasnya WARNA, bukan alpha: ronde 279 mengukur `alphaSapuan`
   * sepanjang rentang wajar dan kontras jejak hanya bergerak
   * 1,303:1 -> 1,541:1. Dengan warna, jalur mayoritas naik
   * **1,406:1 -> 5,381:1** (titik tetap kanvas memudar).
   */
  pipa: "--cair-sorot",
  /**
   * ⚠️ `--cair-sorot` -> `--cair-kepala` — RONDE 281 B.2.
   *
   * Jalur aksen (2 dari 7) kini PUTIH, bukan biru muda: ia yang paling
   * menonjol, dan `peluangSorot` tetap berarti apa adanya. Titik tetap
   * 6,919:1. Nol hex baru — putih sudah hidup sebagai `--bg`.
   */
  sorot: "--cair-kepala",
  /**
   * KEPALA jalur yang sedang bergerak — RONDE 281 B.1.
   *
   * *"SAAT BERGERAK WARNANYA LEBIH TERLIHAT PUTIHNYA"*. Ia digambar
   * sebagai sapuan TERPISAH di posisi terkini, ber-`alphaKepala` yang
   * lebih tinggi daripada `alphaSapuan`.
   *
   * Kenapa sapuan terpisah, bukan sekadar mewarnai sapuan biasa: di
   * kanvas yang memudar, SETIAP titik jejak pernah menjadi kepala.
   * Mewarnai sapuannya putih akan memutihkan seluruh jalur, bukan
   * kepalanya. Yang membedakan kepala hanyalah bahwa ia BELUM memudar —
   * jadi yang harus ditambah di sana adalah TERANG, bukan warna lain.
   */
  kepala: "--cair-kepala",
} as const;
// ⚠️ Warna keempat rancangan ber-ACC, `--pipa-aksen` (#2A4B78), SENGAJA
// tidak ada di sini. Ia dipakai VIGNETTE, dan vignette hidup di CSS
// (`.login::after`) karena ia gradien diam — menggambarnya di kanvas
// berarti membayarnya setiap bingkai, berlawanan dengan seluruh sebab
// ronde ini ada.
//
// Versi pertama memuatnya, dan itu keliru: peta bernama TOKEN_PIPA yang
// mencantumkan token yang TIDAK PERNAH DIBACA modul ini berbohong kepada
// pembacanya — dan pagar yang menyapu `Object.entries(TOKEN_PIPA)` akan
// lulus atas entri itu secara hampa, karena namanya memang ada di berkas
// ini: di deklarasinya sendiri.

export const PIPA = {
  /**
   * ⚠️ 30 -> 7 RONDE 277 (B.3) — keputusan pemilik proyek 24-09-2026.
   *
   * Rancangan ber-ACC (`KEPUTUSAN_PENDING.md`): *"<=7 jalur"*, dan
   * kalimat yang menjelaskannya *"LEBIH TERLIHAT TAPI JANGAN NORAK"*.
   *
   * Keduanya terdengar berlawanan dan sebenarnya satu: 30 jalur tipis
   * saling menimpa menjadi kabut, dan menaikkan opasitas kabut membuatnya
   * NORAK tanpa membuat satu jalur pun lebih terbaca. Tujuh jalur yang
   * masing-masing lebih pekat terbaca sebagai GARIS.
   *
   * Ia juga yang membayar kenaikan opasitas di bawah: peluang dua pipa
   * bersilangan turun kira-kira sebanding (30/7)^2 ~ 18x, dan kasus dua
   * sapuan itulah yang membatasi opasitas sejak ronde 269.
   */
  jumlah: 7,
  belokan: 8,
  peluangBelok: 58,
  lajuDasar: 0.5,
  lajuRentang: 1,
  umurDasar: 100,
  umurRentang: 300,
  /**
   * Lebar garis, 3..8 px (`lebarDasar + acak(lebarRentang)`).
   *
   * §1.3 ronde 269 — pemilik proyek: *"PIPELINE NYA KURANG TERLIHAT, BUAT
   * LEBIH TERLIHAT DIKIT TAPI TETAP TERKESAN CLEAN DAN PROFESIONAL"*.
   * Naik dari 2..6 (rata 4) menjadi 3..8 (rata 5,5) = **+37,5% luas
   * tersapu**.
   *
   * # Kenapa LEBAR yang dinaikkan, bukan alpha
   *
   * Lebar garis TIDAK MUNCUL di rumus titik tetap `C*` sama sekali — ia
   * menaikkan LUAS yang tersapu, bukan PUNCAK kecerahan sebuah piksel.
   * Jadi seluruh tabel kontras tidak bergerak satu digit pun karena
   * perubahan ini.
   *
   * Itu membuatnya tuas termurah yang tersedia, dan ronde 268 sudah
   * membuktikan tuas termahal (`alphaSapuan`) berbahaya: percobaan
   * pilihan-rasa di sana MELEWATI plafonnya. `alphaSapuan` karena itu
   * TIDAK disentuh ronde ini — lihat catatannya di bawah.
   */
  lebarDasar: 3,
  lebarRentang: 5,
  /**
   * ⚠️ MENGGANTIKAN `ronaDasar: 180` / `ronaRentang: 60` — RONDE 277 B.3.
   *
   * Yang lama memilih warna sendiri lewat `hsla(rona,75%,50%)`, dengan
   * alasan yang waktu itu benar: rentang 180..240 kebetulan sudah berada
   * di keluarga warna aplikasi, jadi tak ada warna yang perlu dikarang.
   *
   * Keputusan pemilik proyek 24-09-2026 menetapkan paletnya secara
   * EKSPLISIT — `#0B1120` / `#134d93` / `#8cecff` / `#2A4B78` — dan palet
   * yang ditetapkan mengalahkan palet yang kebetulan cocok. Warnanya kini
   * dibaca dari token (`TOKEN_PIPA`), jadi modul ini tidak lagi memilih
   * warna sama sekali.
   *
   * # Kenapa SOROT harus JARANG
   *
   * `--cair-sorot` berluminans 0,72786 — EMPAT KALI plafon AA. Ronde 272
   * sudah menghadapinya dan menjawabnya dengan memindahkan seluruh teks
   * ke DALAM kartu, bukan dengan menggeser warnanya. Jawaban itu masih
   * berlaku (pagar `webglLiquid.ronde272.test.ts` menjaganya), tetapi
   * "tidak menyakiti teks" bukan berarti "boleh membanjiri layar".
   *
   * Rancangan ber-ACC menyebutnya *"glow tipis"*. Di sini itu berarti:
   * satu dari tujuh jalur membawa warna sorot, sisanya `--cair-tengah`.
   */
  /**
   * ⚠️ 1/7 -> 2/7 — RONDE 279 B.2, dan ia tuas yang BERBEDA dari alpha.
   *
   * Pemilik proyek 24-09-2026: efek pipa *"LEBIH TERLIHAT LAGI"* — kedua
   * kalinya, sesudah ronde 277 sudah menaikkan alpha 0,03 -> 0,06.
   *
   * # Kenapa alpha SAJA tidak akan menjawabnya
   *
   * Latar barunya `#1F2555` (L=0,02270) 3,9x lebih terang daripada
   * `#0B1120` (L=0,00576), dan terang itu MENELAN `--cair-tengah`:
   *
   *	pipa #134D93 vs latar LAMA   2,251:1
   *	pipa #134D93 vs latar BARU   1,727:1   <- permintaan pertama
   *	                                          MELAWAN yang kedua
   *
   * Pada titik tetap kanvas, menaikkan alpha nyaris tidak menggerakkan
   * warna itu — a=0,06 memberi 1,303:1 dan a=0,20 hanya 1,541:1. Yang
   * benar-benar terbaca adalah jalur ber-`--cair-sorot`: 3,944:1 pada
   * a=0,06 dan 5,381:1 pada a=0,10.
   *
   * Jadi yang dinaikkan PORSI jalur sorot, bukan hanya opasitasnya.
   *
   * # Kenapa 2/7 dan bukan lebih
   *
   * Alasan ronde 277 tetap berlaku apa adanya: `--cair-sorot` berluminans
   * 0,72786, empat kali plafon AA. Ia boleh menjadi AKSEN, bukan
   * permukaan. Dua dari tujuh masih minoritas; pagar
   * `login.ronde279.test.ts` mematoknya <= 3/7 supaya "lebih terlihat"
   * tidak pelan-pelan menjadi "seluruhnya sorot".
   */
  peluangSorot: 2 / 7,
  /**
   * Alpha sapuan. Sumbernya 0.125.
   *
   * Diturunkan karena jejaknya kini memudar (lihat perbedaan 3): dengan
   * pemudaran, terang MANTAP-nya punya titik tetap tertutup
   *
   *     C* = (navy*f*(1-a) + pipa*a) / (1 - (1-f)(1-a))
   *
   * untuk piksel terburuk — yang tersapu SETIAP bingkai pada alpha puncak.
   *
   * Pasangan ini DIHITUNG, bukan dipilih. Percobaan pertama (a=0.055,
   * f=0.045) MELEBIHI plafonnya: L=0.18253 lawan `--blue` 0.16571, yaitu
   * membatalkan perbaikan ronde 264 tanpa satu pun gejala di layar.
   *
   *	Terukur 22-09-2026, rona 180 (cyan — paling terang di rentang):
   *	  a=0.03 f=0.05  ->  RGB mantap (14, 84, 141)   L = 0.09025
   *	  plafon `--blue`                               L = 0.16571
   *	  putih di atasnya                              7,49:1
   *	Rona 190..240 seluruhnya LEBIH gelap (8,61:1 .. 14,15:1).
   *
   * Dipilih pasangan dengan JEJAK TERPANJANG yang masih di bawah plafon —
   * pudar 0.05 berarti jejak meluruh ~20 bingkai (~0,33 detik pada 60fps).
   */
  /**
   * ⚠️ 0.03 -> 0.06 RONDE 277 (B.3) — dan PLAFONNYA berganti, bukan
   *   dilonggarkan. Seluruh tabel di atas dihitung terhadap `--navy`
   *   (#03045e) dengan warna HSL; keduanya sudah tidak berlaku.
   *
   * Rancangan ber-ACC: *"sorot #8cecff **opasitas naik**"*.
   *
   * # Plafon lama SUDAH MATI sebelum ronde ini, dan itu yang menentukan
   *
   * Plafon `--blue` (L=0,16571) diwarisi dari ronde 264, ketika teks
   * MASIH berdiri di atas latar Login. Ronde 272 mencabut keadaan itu:
   * `.login__foot` pindah ke DALAM kartu (§6.1 jalan (a)), kartunya buram
   * penuh, dan pemilik proyek menetapkan `--cair-sorot` (L=0,72786 —
   * 4,4x plafon `--blue`) sebagai warna latar yang sah.
   *
   * Jadi mempertahankan plafon `--blue` di sini berarti menegakkan syarat
   * yang premisnya sudah tidak ada — sementara latar di sebelahnya sudah
   * empat kali melewatinya dengan persetujuan tertulis.
   *
   * # Plafon yang MENGGANTIKANNYA, dan ia bukan sekadar lebih longgar
   *
   * Latar tidak boleh menjadi LEBIH TERANG daripada warna paling terang
   * yang pemilik proyek setujui, yaitu `--cair-sorot` itu sendiri.
   * Akumulasi yang melewatinya berarti animasi ini mengarang terang yang
   * tidak ada di paletnya.
   *
   *	Terukur 24-09-2026, latar #0B1120, sapuan #8cecff:
   *	  a=0.06 f=0.06  1 sapuan   L = 0,1940
   *	  a=0.06 f=0.06  2 sapuan   L = 0,3403
   *	  plafon `--cair-sorot`     L = 0,72786
   *
   * Kasus dua sapuan masih 2,1x di bawah plafon. Yang membelinya
   * `jumlah: 30 -> 7`: persilangan turun ~18x, dan kasus dua sapuan
   * itulah yang sejak ronde 269 membatasi opasitas.
   *
   * `lajuPudar` SENGAJA tidak ikut naik: menaikkannya memperpendek jejak,
   * dan jejak itulah yang membuat tujuh jalur terbaca sebagai GARIS
   * alih-alih tujuh titik.
   */
  /**
   * ⚠️ 0,06 -> 0,10 — RONDE 279 B.2, dan PLAFONNYA TIDAK DIGESER.
   *
   * Rencana B.2 mengizinkan menggeser ambang "tidak norak" asalkan
   * disertai angka. Izin itu TIDAK dipakai: 0,10 adalah nilai TERBESAR
   * yang masih berada di dalam plafon yang sudah ada.
   *
   *	a        bobot1    bobot2    plafon 0,65 / 0,80
   *	0,06     0,5155    0,6871    (ronde 277)
   *	0,10     0,6494    0,7963    LULUS, margin 0,0006
   *	0,12     0,6944    0,8292    MELEWATI
   *
   * Marginnya tipis dan itu dinyatakan: ia berarti tuas ini sudah
   * mentok, dan permintaan "lebih terlihat" berikutnya harus dijawab
   * tuas LAIN (porsi sorot, lebar) atau dengan keputusan baru tentang
   * plafonnya — bukan dengan menggeser plafon diam-diam karena kebetulan
   * saya yang diuntungkan.
   *
   * Akumulasi tetap jauh di bawah plafon warna: dua sapuan mencapai
   * L=0,48264 lawan `--cair-sorot` 0,72786 (1,51x di bawah).
   */
  alphaSapuan: 0.1,
  /**
   * Alpha sapuan KEPALA — RONDE 281 B.1.
   *
   * Keputusan pemilik proyek 24-09-2026 (pesan malam): *"SAAT BERGERAK
   * WARNANYA LEBIH TERLIHAT PUTIHNYA"*.
   *
   * # Kenapa ini BOLEH melewati `alphaSapuan`, sementara `alphaSapuan`
   *   sendiri tetap dipatok
   *
   * Peringatan pada `alphaSapuan` di atas berbunyi: plafonnya mentok,
   * dan permintaan "lebih terlihat" berikutnya harus dijawab tuas LAIN
   * atau keputusan baru. Ronde ini memakai KEDUANYA, dan itu disengaja:
   * tuas lain (sapuan kepala tersendiri), dengan keputusan pemilik yang
   * dikutip di atas.
   *
   * Yang membuatnya AMAN sementara menaikkan `alphaSapuan` tidak:
   * `alphaSapuan` berlaku di SEPANJANG jejak, sehingga menaikkannya
   * menaikkan terang seluruh permukaan — persis "norak" yang rancangan
   * tolak. Sapuan kepala hanya menyentuh SATU titik per pipa per
   * bingkai, dan titik itu langsung mulai memudar begitu pipa bergerak.
   *
   * # Angkanya
   *
   *	alpha 0,10 -> kepala 1,788:1 (n=2)   tidak terlihat sebagai kepala
   *	alpha 0,20 -> kepala 3,055:1
   *	alpha 0,30 -> kepala 4,684:1         DIPILIH
   *	alpha 0,45 -> kepala 7,488:1         mulai menjadi titik menyilaukan
   *
   * Pada 0,3 kepala berjarak 3,12x dari jejak sorot di bawahnya — cukup
   * untuk terbaca sebagai kepala, belum menjadi lampu.
   */
  alphaKepala: 0.3,
  /**
   * Jari-jari kepala sebagai NISBAH terhadap lebar sapuan — ronde 281 B.1.
   *
   * Lebih kecil dari sapuannya supaya kepala terbaca sebagai TITIK di
   * dalam jalur, bukan sebagai jalur kedua yang lebih putih. Nisbah,
   * bukan piksel: lebar tiap pipa diacak saat lahir, dan angka tetap
   * akan membuat kepala melebihi badannya pada pipa yang tipis.
   */
  lebarKepala: 0.55,
  /**
   * Alpha latar yang ditimpakan tiap bingkai untuk memudarkan jejak.
   * Lihat perhitungan pada `alphaSapuan` — keduanya satu pasangan dan
   * tidak boleh disetel sendiri-sendiri tanpa mengukur ulang.
   *
   * # KOREKSI ronde 269 — 0.05 -> 0.06, dan sebabnya lubang di pagar
   *   ronde 268 yang dipasang oleh ronde 268 sendiri
   *
   * Rumus titik tetap di atas memodelkan piksel yang tersapu SATU kali per
   * bingkai. Laporan ronde 268 menyatakan plafonnya seolah berlaku umum —
   * ia berlaku untuk sapuan tunggal. Di tempat dua pipa BERSILANGAN dan
   * berjalan beriringan, sapuannya dua kali per bingkai, dan alpha
   * efektifnya `1-(1-a)^2`, bukan `a`.
   *
   *	Terukur 22-09-2026, rona 180, a=0.03:
   *	  f=0.05  1 sapuan   L = 0.09025   -45,5% thd plafon   LULUS
   *	  f=0.05  2 sapuan   L = 0.17812    +7,5% thd plafon   GAGAL
   *	  f=0.06  1 sapuan   L = 0.07442   -55,1%              LULUS
   *	  f=0.06  2 sapuan   L = 0.15173    -8,4%              LULUS
   *
   * Jadi keadaan SEBELUM ronde ini sudah melewati plafon di titik
   * silang; pagarnya hijau karena ia hanya menguji baris pertama.
   * Menaikkan lebar garis memperbesar daerah persilangan itu, sehingga
   * mengirim kenaikan lebar tanpa menutup lubangnya adalah kelalaian.
   *
   * Ongkosnya: jejak meluruh ~17 bingkai alih-alih ~20 (0,28 dtk lawan
   * 0,33 dtk pada 60fps) — sebagian termakan, dan kenaikan lebar +37,5%
   * lebih dari menggantinya.
   *
   * `alphaSapuan` SENGAJA tidak ikut disentuh: menurunkannya akan
   * MEREDUPKAN, berlawanan dengan yang diminta; menaikkannya melewati
   * plafon pada dua sapuan (a maks untuk 2 sapuan = 0.027, di BAWAH nilai
   * sekarang).
   */
  lajuPudar: 0.06,
} as const;

export type RGB = [number, number, number];

/**
 * `#rgb` / `#rrggbb` -> RGB, atau `null` bila tokennya kosong.
 *
 * ⚠️ SENGAJA tidak diimpor dari `closingPlasma.ts` (yang mengekspornya)
 * maupun `webglLiquid.ts` (yang meneruskannya). Keduanya modul yang
 * DITAHAN — jalan pulang yang boleh dihapus kapan saja pemilik proyek
 * memutuskan. Modul yang HIDUP tidak boleh bergantung pada modul yang
 * sedang menunggu penghapusan; kalau tidak, "menahan" berubah menjadi
 * "tidak dapat dihapus". Enam baris adalah harga kemerdekaan itu.
 */
export function bacaHex(nilai: string): RGB | null {
  const t = nilai.trim().replace("#", "");
  if (t.length !== 3 && t.length !== 6) return null;
  const p =
    t.length === 3
      ? t.split("").map((c) => c + c)
      : [t.slice(0, 2), t.slice(2, 4), t.slice(4, 6)];
  const v = p.map((h) => parseInt(h, 16));
  return v.some((n) => Number.isNaN(n)) ? null : (v as RGB);
}

/**
 * Warna sapuan sebuah pipa: `--cair-sorot` bila ia jalur sorot, selain
 * itu `--cair-tengah`.
 *
 * Fungsi MURNI supaya keputusannya dapat diuji tanpa kanvas — dan supaya
 * pagar "sorot JARANG" mengikat pilihannya, bukan ejaan `hsla()`.
 */
export function warnaSapuan(sorot: boolean, pipa: RGB, terang: RGB): RGB {
  return sorot ? terang : pipa;
}

/**
 * mulaiPipeline memasang animasi pada `kanvas` dan mengembalikan fungsi
 * penghentinya.
 *
 * Pemanggil WAJIB memanggil penghentinya saat komponen dilepas. Tanpa itu
 * loop-nya hidup terus walau layarnya sudah tidak ada.
 */
export function mulaiPipeline(kanvas: HTMLCanvasElement): Penghenti {
  const { PI, cos, sin, abs, round, random } = Math;
  const TAU = 2 * PI;
  const SETENGAH_PI = 0.5 * PI;
  const TO_RAD = PI / 180;
  const acak = (n: number) => n * random();
  const pudarMasukKeluar = (t: number, m: number) => {
    const hm = 0.5 * m;
    return abs(((t + hm) % m) - hm) / hm;
  };

  const ctx = kanvas.getContext("2d");
  if (!ctx) return () => {};

  // Palet dibaca dari token, bukan ditulis di sini.
  //
  // ⚠️ `getComputedStyle(kanvas)` — BUKAN `document.documentElement`.
  // Keempat token dideklarasikan pada `.login`, dan kanvas ini anaknya.
  // Membacanya dari akar dokumen mengembalikan string KOSONG, dan modul
  // akan diam total tanpa satu pun galat (lihat penolakan di bawah).
  const gaya = getComputedStyle(kanvas);
  const bacaLatar = bacaHex(gaya.getPropertyValue(TOKEN_PIPA.latar));
  const bacaPipa = bacaHex(gaya.getPropertyValue(TOKEN_PIPA.pipa));
  const bacaSorot = bacaHex(gaya.getPropertyValue(TOKEN_PIPA.sorot));
  const bacaKepala = bacaHex(gaya.getPropertyValue(TOKEN_PIPA.kepala));
  // Tanpa palet, TIDAK menggambar. Mengarang warna pengganti berarti
  // layar Login tampil dengan warna yang tidak seorang pun putuskan —
  // pola yang sama dipakai `webglLiquid.ts` dan `closingPlasma.ts`.
  if (!bacaLatar || !bacaPipa || !bacaSorot || !bacaKepala) return () => {};
  // Diikat ulang BERTIPE, bukan dipaksa dengan `!`. Penyempitan `null`
  // tidak bertahan melintasi batas closure, dan `!` akan membungkam
  // kompiler alih-alih menjawabnya — tanda seru itu tepat sekali gagal
  // ketika seseorang kelak memindahkan penolakan di atas.
  const WARNA_LATAR: RGB = bacaLatar;
  const WARNA_PIPA: RGB = bacaPipa;
  const WARNA_SOROT: RGB = bacaSorot;
  const WARNA_KEPALA: RGB = bacaKepala;
  const latar = `rgb(${WARNA_LATAR[0]},${WARNA_LATAR[1]},${WARNA_LATAR[2]})`;

  const PROP = 8;
  const panjang = PIPA.jumlah * PROP;
  let pipa = new Float32Array(panjang);
  let tick = 0;
  let idFrame = 0;
  let berhenti = false;

  function ukurUlang() {
    // devicePixelRatio dibatasi 2: di atas itu biaya piksel naik kuadratik
    // sementara bedanya tidak terlihat pada latar yang sengaja kabur.
    const dpr = Math.min(window.devicePixelRatio || 1, 2);
    kanvas.width = Math.floor(kanvas.clientWidth * dpr);
    kanvas.height = Math.floor(kanvas.clientHeight * dpr);
  }

  function awaliPipa(i: number) {
    pipa.set(
      [
        acak(kanvas.width),
        0.5 * kanvas.height,
        round(acak(1)) ? SETENGAH_PI : TAU - SETENGAH_PI,
        PIPA.lajuDasar + acak(PIPA.lajuRentang),
        0,
        PIPA.umurDasar + acak(PIPA.umurRentang),
        PIPA.lebarDasar + acak(PIPA.lebarRentang),
        // Slot 7 dulu menyimpan RONA; kini bendera SOROT (1 = jalur
        // terang). Nilainya ditetapkan sekali saat pipa lahir, bukan tiap
        // bingkai — jalur yang berkedip antara dua warna terbaca sebagai
        // kerlip, dan itu persis "norak" yang rancangannya tolak.
        random() < PIPA.peluangSorot ? 1 : 0,
      ],
      i,
    );
  }

  function awaliSemua() {
    pipa = new Float32Array(panjang);
    for (let i = 0; i < panjang; i += PROP) awaliPipa(i);
  }

  function gambarSatu(i: number) {
    const x = pipa[i] ?? 0;
    const y = pipa[i + 1] ?? 0;
    const umur = pipa[i + 4] ?? 0;
    const ttl = pipa[i + 5] ?? 0;
    const lebar = pipa[i + 6] ?? 0;
    const sorot = pipa[i + 7] === 1;

    const selubung = pudarMasukKeluar(umur, ttl);
    const a = selubung * PIPA.alphaSapuan;
    const [r, g, b] = warnaSapuan(sorot, WARNA_PIPA, WARNA_SOROT);
    ctx!.save();
    ctx!.strokeStyle = `rgba(${r},${g},${b},${a})`;
    ctx!.beginPath();
    ctx!.arc(x, y, lebar, 0, TAU);
    ctx!.stroke();
    ctx!.closePath();

    /* KEPALA — ronde 281 B.1, keputusan pemilik proyek 24-09-2026:
       "SAAT BERGERAK WARNANYA LEBIH TERLIHAT PUTIHNYA".

       Sapuan KEDUA di titik yang sama, lebih putih dan lebih pekat,
       berjari-jari lebih kecil. Ia ikut memudar seperti sapuan pertama —
       jadi yang terlihat adalah titik terang di ujung yang bergerak,
       dengan ekor putih pendek yang cepat luruh di atas jejak berwarna.

       Memakai selubung yang SAMA (`pudarMasukKeluar`): pipa yang sedang
       lahir atau sedang mati tidak boleh tiba-tiba punya kepala terang
       sementara badannya belum/tidak lagi ada.

       Biayanya satu `arc` tambahan per pipa per bingkai: 7 -> 14 sapuan.
       Ia TIDAK menyentuh loop medan per-piksel yang menjadi biaya
       dominan latar (lihat lib/bebanLatar.ronde277.test.ts). */
    const [kr, kg, kb] = WARNA_KEPALA;
    ctx!.strokeStyle = `rgba(${kr},${kg},${kb},${selubung * PIPA.alphaKepala})`;
    ctx!.beginPath();
    ctx!.arc(x, y, lebar * PIPA.lebarKepala, 0, TAU);
    ctx!.stroke();
    ctx!.closePath();
    ctx!.restore();
  }

  function majukan(i: number) {
    let x = pipa[i] ?? 0;
    let y = pipa[i + 1] ?? 0;
    let arah = pipa[i + 2] ?? 0;
    const laju = pipa[i + 3] ?? 0;
    let umur = pipa[i + 4] ?? 0;
    const ttl = pipa[i + 5] ?? 0;

    gambarSatu(i);

    umur++;
    x += cos(arah) * laju;
    y += sin(arah) * laju;

    const pembagi = round(acak(PIPA.peluangBelok)) || 1;
    const belok = !(tick % pembagi) && (!(round(x) % 6) || !(round(y) % 6));
    if (belok) {
      arah += (360 / PIPA.belokan) * TO_RAD * (round(acak(1)) ? -1 : 1);
    }

    pipa[i] = x;
    pipa[i + 1] = y;
    pipa[i + 2] = arah;
    pipa[i + 4] = umur;

    if (umur > ttl) awaliPipa(i);
  }

  function bingkai() {
    if (berhenti) return;
    tick++;

    // Pemudaran: latar ditimpakan tipis SEBELUM pipa baru digambar.
    // Inilah yang membatasi terangnya — lihat perbedaan 3 di kepala berkas.
    ctx!.save();
    ctx!.globalAlpha = PIPA.lajuPudar;
    ctx!.fillStyle = latar;
    ctx!.fillRect(0, 0, kanvas.width, kanvas.height);
    ctx!.restore();

    for (let i = 0; i < panjang; i += PROP) majukan(i);

    idFrame = window.requestAnimationFrame(bingkai);
  }

  function jalan() {
    if (berhenti || idFrame) return;
    idFrame = window.requestAnimationFrame(bingkai);
  }

  function jeda() {
    if (idFrame) {
      window.cancelAnimationFrame(idFrame);
      idFrame = 0;
    }
  }

  // Tab tersembunyi: hentikan. Tanpa ini, tab Login yang ditinggalkan di
  // latar belakang tetap menggambar 60 bingkai per detik selamanya (§0.1.1).
  function saatTampilBerubah() {
    if (document.hidden) jeda();
    else jalan();
  }

  ukurUlang();
  // Latar dicat PENUH sekali di awal supaya bingkai pertama tidak
  // memperlihatkan kanvas transparan sekejap.
  ctx.fillStyle = latar;
  ctx.fillRect(0, 0, kanvas.width, kanvas.height);
  awaliSemua();

  function saatUkuranBerubah() {
    ukurUlang();
    ctx!.fillStyle = latar;
    ctx!.fillRect(0, 0, kanvas.width, kanvas.height);
    awaliSemua();
  }

  window.addEventListener("resize", saatUkuranBerubah);
  document.addEventListener("visibilitychange", saatTampilBerubah);
  if (!document.hidden) jalan();

  return () => {
    berhenti = true;
    jeda();
    window.removeEventListener("resize", saatUkuranBerubah);
    document.removeEventListener("visibilitychange", saatTampilBerubah);
  };
}
