# 56: Tab Correspondence (kasus Inward Facultative)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 06-10-2026 "lanjut ke tab
> Correspondence". Survei korpus oleh subagen (pengurai XML); butir kunci diperiksa ulang sesi 0f: `SendCorrespondence`
> (To / CC, langkah set `.IsEditable`), kontainer lampiran `1=2`, label verbatim `Correspondence` / `CorrespondenceContent`.

**Status:** frontend selesai 06-10-2026 (uji hijau); backend + SMTP + migrasi tabel baru diserahkan ke sesi c3.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\`

- `Section\InputInwardFacultativeDtl.xml` kontainer S95 "Correspondence": selalu tampil (tanpa `IsIT` / When), tanpa
  pra-muat, isi section `Correspondence` (kelas Work). Sama di `InputInwardFacultativeDtl_IsUW.xml` dan `Harness\ViewPolis.xml`.
- `Section\Correspondence.xml` (`ASM-FW-GISFW-WORK!CORRESPONDENCE`, memo "added sent status"): grid `.OfferFacIn.CorrespondenceList`
  (kelas `ASM-FW-GISFW-Data-Correspondence`) — Date · Sender · Subject · Recipient · Status (`.IsEditable`, dropdown
  baca-saja); tombol Add (`DataTransform\CorrespondenceAdd.xml`: IsEditable 1, Sender / SenderMail = operator login,
  Date = sekarang); Remove per baris bila `.IsEditable = 1`; baris dibuka = flow action `CorrespondenceContent`.
- `Section\CorrespondenceContent.xml` (memo "Added send email button, hidden attachment module bcs literally unusable"):
  Sender, Sender Mail Address (baca-saja); grid Recipients: (`.RecipientList`: Recipient autocomplete operator
  `pyGetListOfOperators` → `.RecipientMail` = email operator; Add `RecipientListAdd`; Remove); Subject; Message (rich text);
  kontainer lampiran SELALU tersembunyi (`pyContainerVisibleWhen` = `1=2`); tombol Save (`SaveFacIn_Act` = Obj-Save) dan
  Send Email (`SendCorrespondence`) bila `.IsEditable = 1`.
- `Activity\SendCorrespondence.xml`: `SendEmailWithAttachments` dengan **To = `local.Email` = `.SenderMail` (pengirim
  sendiri)**, **CC = semua `.RecipientMail`** (dipisah koma), HTML; SMTP host / akun ditulis di langkah (konfigurasi);
  lalu `.IsEditable = 0` untuk korespondensi dan semua penerimanya (status "terkirim", terkunci). Tanpa Obj-Save.
- `DataTransform\SetCorrespondenceReadOnly.xml`: semua korespondensi `.IsEditable = "0"` saat flow action
  `InwardFacultative` (Submit) dijalankan (`AddToListSuggestOfferFacIn_DT` langkah 4).
- Penyimpanan Pega: tertanam di halaman kerja (`pyWorkPage.OfferFacIn.CorrespondenceList`) + salinan JSON
  `POOLDATA.JSON_OFFER`; tidak ada tabel khusus di `DDL\`. Sampel `DDL\CONTOH` / `P-5`: 113 `CorrespondenceList`, SEMUA kosong.

## Cacat / keanehan Pega `[terverifikasi]`

1. Email dikirim ke pengirim sendiri (To), penerima hanya CC.
2. Kolom grid "Recipient" (`.Recipient` tingkat korespondensi) tidak pernah diisi rule mana pun → kosong `[dugaan]`.
3. Send Email tidak menyimpan (tanpa Obj-Save) — status "terkirim" hilang bila tidak di-Save `[dugaan]`.

## Batas pengetahuan / ketergantungan

- Pengiriman email nyata: `inti/backend/outbox` punya efek email, tetapi `ErrEmailBelumDisetujui` — menghubungkan layanan
  email nyata menuntut persetujuan manusia (CLAUDE.md §10, ADR-0008 / ADR-0015, gerbang lingkungan ADR-0005).
- Daftar penerima: operator Pega → padanan sistem baru = akun `M_LOGIN_GO` (NAME, EMAIL).
- Nilai tampilan Status (`.IsEditable` daftar nilai properti) tidak ada di korpus.

## Keputusan work owner 06-10-2026 ("iya kirim email sungguhan. saya setuju dengan rekomendasi kamu, buatin dan jangan halu")

1. **Email dikirim sungguhan** — persetujuan manusia (CLAUDE.md §10) untuk menghubungkan SMTP nyata bagi korespondensi
   nbfacin. Alamat / akun SMTP dari env var (ADR-0004), nol literal di kode. SMTP belum dikonfigurasi → tidak terkirim,
   status tetap "Not Sent", isi tetap tersimpan.
2. **Penerima dibetulkan** — To = penerima, Cc + Reply-To = pengirim; kolom Recipient grid = nama penerima.
3. **Daftar penerima** = akun sistem baru ber-email (`M_LOGIN_GO`); alamat bebas boleh diketik.
4. **Lampiran tidak dibuat** (Pega `1=2`).
5. **Tabel baru** (migrasi ditulis c3, dijalankan work owner).
6. Gambar layar Pega diterima (grid, editor; Status "Not Sent"; tanggal `06/10/26 17:28`).
7. **Akses** (work owner 06-10-2026 "ini visibility menu nya juga cuma it ya") — tab hanya untuk akun ber-`JOB_POSITION` = "IT",
   aturan sama dengan Loss Record (tiket 55): `KasusNB.aksesKorespondensi`, tab disaring; backend 401 / 403 di semua endpoint
   korespondensi (diminta ke c3). Pega: S95 selalu tampil.
8. **Jawaban atas K56-9** (work owner 06-10-2026 "setuju" atas saran sesi 0f): (1) TLS WAJIB - tanpa TLS tidak dikirim,
   kecuali dimatikan eksplisit lewat env var; (2) wewenang kirim tidak dibatasi lebih jauh (tab sudah hanya IT); (3) tanpa
   gerbang lingkungan - DEV cukup tidak mengisi `NBFACIN_SMTP_*`. Diteruskan ke c3.

## Frontend (sesi 0f)

- `components/TabCorrespondence.tsx` (+ `TabCorrespondence.test.tsx`), `components/EditorTeksKaya.tsx` (contentEditable +
  execCommand, tanpa pustaka baru), `components/htmlKorespondensi.ts` (pembersih HTML daftar-izin, sama dengan backend);
  `api.ts` kontrak; `labels.ts` `KORESPONDENSI` / `TEKS_KORESPONDENSI`; `nbfacin.css` `.nbf-kor*` / `.nbf-teks-kaya*`;
  `pages/InwardFacultative.tsx` (tab untuk semua case). Diperiksa visual lewat Chrome headless (data sintetis).

## Keputusan agent

- **K56-1** Add langsung membuat draf di server (tanggal + pengirim dari akun login, bukan dari layar); Remove langsung
  menghapus draf (Pega: baris clipboard sampai Save).
- **K56-2** Send Email menyimpan isi dulu lalu mengirim (Pega: tanpa Obj-Save).
- **K56-3** Toolbar Message: Bold, Italic, Underline, daftar berpoin / bernomor, hapus format, Undo, Redo. Font, Size, warna
  teks / latar Pega TIDAK disediakan — memerlukan atribut gaya yang dibuang pembersih HTML (keamanan).
- **K56-4** Status "Sent" `[dugaan]` (hanya "Not Sent" terlihat di gambar layar; daftar nilai `.IsEditable` tidak di korpus).
- **K56-6 (sesi 0f, menjawab pertanyaan c3 no. 3)** Penghitung ukuran Message di bawah editor (`bitaPesan` = bita UTF-8 HTML
  bersih, sama dengan hitungan backend); > 4000 bita → peringatan merah, Save dan Send Email nonaktif. Batas mengikuti
  MESSAGE VARCHAR2(4000) (penjaga inti melarang CLOB di migrasi modul).
- **K56-5** Penguncian semua korespondensi saat Submit (`SetCorrespondenceReadOnly`) menyusul ketika Submit dibangun.

## Backend (sesi c3, 06-10-2026)

**Status:** backend selesai; migrasi 207 DITULIS, belum dijalankan. Kontrak = `frontend/api.ts` (`Korespondensi`,
`TampilanKorespondensi`, `IsiKorespondensi`, `saranPenerimaKorespondensi`). Persetujuan manusia untuk SMTP nyata = keputusan
work owner 06-10-2026 (diteruskan sesi 9d, dicatat di sini); tidak ada email yang dikirim selama pengembangan / uji.

- **K56-1 Tabel** — migrasi `207_korespondensi.sql`: `T_CORRESPONDENCELIST` (induk T_GENERAL_POLIS, SEQ_NO = urutan Add) +
  `T_RECIPIENTLIST`. ⚠️ MESSAGE = **VARCHAR2(4000)**, bukan CLOB: penjaga inti `migrasi_test.go` melarang CLOB / BLOB / JSON
  di migrasi kecuali pengecualian bernama milik tim inti (outbox 015, templat 912). Pesan bersih > 4000 bita → 400.
  `[pertanyaan terbuka]` work owner / tim inti: bila email lebih panjang dibutuhkan, tim inti menambah pengecualian CLOB
  untuk kolom ini dan nbfacin menulis migrasi baru.
- **K56-2 Endpoint** — kasus apa pun (Pega S95 selalu tampil). Add = draf (tanggal sekarang, pengirim = NAME / EMAIL
  M_LOGIN_GO akun login; akun tak ada → kosong); 401 tanpa sesi (semua tulis). PUT = Save draf (409 terkirim, 404 nomor bukan
  milik kasus / bukan angka). DELETE = draf saja → 204. Kirim = simpan dulu (seperti PUT), lalu kirim.
- **K56-3 Pesan** — `BersihkanHTMLKorespondensi` (`services/htmlkorespondensi.go`) = port satu-satu `bersihkanHtml`
  frontend, dijalankan sebelum simpan dan kirim. Instrumen: 22 masukan (termasuk img onerror, javascript:, `<SCRIPT>`,
  svg onload, iframe, tag bersarang, script tanpa penutup) — jawabannya dihasilkan pembersih frontend sendiri (Node 24)
  dan port Go sama persis (`TestBersihkanHTMLKorespondensiParitas`).
- **K56-4 Kirim** — To = recipientMail unik (tanpa beda huruf; baris tanpa email dilewati), Cc + Reply-To = senderMail sah
  (Cc dilewati bila pengirim sudah di To); From = `NBFACIN_SMTP_FROM`; Subject encoded-word UTF-8; badan text/html UTF-8
  base64; Subject panjang / non-ASCII dilipat (encoded-word ≤ 45 bita per potong, tiap baris header ≤ 998 — RFC 5322).
  Alamat = alamat polos sah (`net/mail`, tanpa nama tampilan, tanpa CR / LF). Urutan galat: 401 → 404 / 409 (status baris
  menang atas isian tidak sah) → 400 isian → 400 penerima (tidak ada yang ber-email / alamat tidak sah, pesan per alamat) → 503 SMTP belum
  dikonfigurasi → 502 SMTP gagal. Pada 400 penerima / 503 / 502 isi TETAP tersimpan, status tetap draf. Berhasil →
  `IS_EDITABLE` "0" + SENT_AT (korespondensi dan penerimanya). Jejak: satu baris log berhasil / gagal (kasus, nomor, jumlah
  alamat, akun / galat) — tanpa isi pesan, tanpa kredensial; rincian galat server SMTP tidak dikirim ke layar.
  Email terkirim tetapi status gagal ditulis → **500 "email terkirim tetapi status gagal disimpan - jangan kirim ulang"**
  (galat tersendiri); penulisan status memakai `context.WithoutCancel` sehingga klien yang menutup koneksi sesudah email
  keluar tidak membuat baris tetap tampil draf.
- **K56-5 SMTP** (`repository/surel.go`) — konfigurasi HANYA env lewat `inti/backend/config` (penjaga inti: env dibaca di
  satu tempat): `NBFACIN_SMTP_HOST`, `NBFACIN_SMTP_PORT` (kosong = 587; 465 = TLS implisit), `NBFACIN_SMTP_USER`,
  `NBFACIN_SMTP_PASSWORD`, `NBFACIN_SMTP_FROM`, `NBFACIN_SMTP_TANPA_TLS`. Host atau From kosong → tidak dipasang (503).
  **TLS wajib** (K56-10): server tanpa STARTTLS → tidak ada yang dikirim (502), kecuali `NBFACIN_SMTP_TANPA_TLS=true` dan
  akun kosong (sandi tidak pernah lewat kanal terbuka; `smtp.PlainAuth` juga menolaknya); TLS ≥ 1.2;
  batas sambung 15 dtk, percakapan 60 dtk. `config.SMTP.String()` tidak memuat host / akun / sandi (aman dicatat saat
  modul menyala). Uji memakai pengirim tiruan dan server SMTP tiruan di loopback — tidak pernah server nyata.
  ⚠️ Tidak ada gerbang lingkungan — **keputusan work owner 06-10-2026** (jawaban K56-9 no. 3): DEV cukup tidak mengisi
  `NBFACIN_SMTP_*`; DEV yang mengisinya AKAN mengirim email sungguhan ke alamat penerima.
- **K56-6 Penerima** — `GET /api/nbfacin/korespondensi/penerima/saran?q=`: M_LOGIN_GO `IS_ACTIVE = '1'` dan EMAIL terisi;
  NAME / EMAIL mengandung q (tidak peka huruf); urut NAME; batas 50; id = LOGIN_ID, label = NAME, ekstra = EMAIL.
- **K56-7 Belum dibangun** — `SetCorrespondenceReadOnly` (kunci semua korespondensi saat Submit flow InwardFacultative):
  sambungkan saat Submit dibangun. Kirim ganda serentak (dua klik bersamaan sebelum status tersimpan) dapat mengirim dua
  kali — kunci baris sebelum kirim tidak dibuat (percakapan SMTP di dalam transaksi ditolak) `[dugaan]` risiko kecil.
- **K56-8 Urutan deploy** — 207 sebelum tab dipakai (GET kasus TIDAK membacanya; tanpa 207 hanya endpoint korespondensi
  yang gagal). Isi `NBFACIN_SMTP_*` di lingkungan yang boleh mengirim email; tanpa itu Send = 503 dan Save tetap jalan.
- **K56-9 Tinjauan dua sumbu (standar + spec), 06-10-2026** — dibetulkan: Subject terlipat (sebelumnya satu baris tanpa
  batas — 1000 bita ASCII melampaui 998); TLS implisit 465 menghormati pembatalan permintaan (`tls.Dialer.DialContext`);
  status-gagal-sesudah-kirim tersendiri (lihat K56-4); 409 / 404 lebih dulu dari 400; Message dibersihkan ULANG saat GET
  (baris yang ditulis di luar layanan — loader / SQL langsung — tidak sampai ke layar mentah; pembersih idempoten, diuji);
  `NBFACIN_SMTP_FROM` wajib alamat polos — selain itu `ErrKonfigurasi` saat start (perubahan **inti** `config.go`, diuji
  `smtp_test.go`; perlu diketahui tim inti). Mutasi atas empat perbaikan layanan masing-masing menggagalkan ujinya.
  Belum diubah — dicatat apa adanya:
  - Batas 4000 bita dihitung atas HTML **bersih** di server. **Terjawab** (0f, 06-10-2026): frontend menambah penghitung
    4000 bita (Save / Send nonaktif bila lewat); server tetap menjadi penentu (400).
  - senderMail akun tidak sah / kosong → Cc dan Reply-To dilewati tanpa pesan; email tetap terkirim.
  - Akun TIDAK diisi dan server tidak menawarkan STARTTLS → dulu email dikirim tanpa enkripsi. **Terjawab** (work owner
    06-10-2026 "setuju"): TLS wajib, pengecualian hanya `NBFACIN_SMTP_TANPA_TLS=true` — dibangun di K56-10.
  - EHLO memakai nama bawaan `net/smtp` ("localhost"); sebagian server menolak / menandai. `[dugaan]` cukup untuk relay
    internal; bila ditolak, nama EHLO dijadikan env (perubahan inti).
  - Wewenang: dulu siapa pun yang login dapat mengirim sampai 100 alamat bebas per email. **Terjawab** (work owner
    06-10-2026): tab kini hanya JOB_POSITION "IT" (K56-10); tanpa batasan peran / domain penerima tambahan.
  - Galat server SMTP yang dicatat di log dapat memuat alamat penerima (teks tolakan RCPT dari server) — log server
    diperlakukan sebagai data pribadi.
  - Kirim ganda serentak — lihat K56-7.
- **K56-10 Akses IT + TLS wajib (c3, 06-10-2026)** — keputusan work owner "ini visibility menu nya juga cuma it ya" dan
  "setuju" atas jawaban K56-9 (diteruskan sesi 9d; perubahan authz / sistem eksternal — CLAUDE.md §10).
  - **Akses**: aturan SAMA dengan Loss Record (K55-1) — `M_LOGIN_GO.JOB_POSITION` "IT" (dipangkas, tanpa beda huruf),
    satu konstanta `jabatanTabIT` untuk kedua tab. GET kasus dan jawaban PUT …/general: `aksesKorespondensi` dari SATU
    pembacaan lunak JOB_POSITION bersama `aksesLossRecord` (`AksesTabIT`; galat baca → false, layar kasus tidak gagal).
    Keenam endpoint korespondensi (GET / POST …/korespondensi, PUT / DELETE …/{no}, POST …/{no}/kirim, GET …/penerima/saran):
    tanpa identitas 401, bukan "IT" / akun tak ada 403, galat baca jabatan diteruskan (500) — di layanan
    (`WajibAksesKorespondensi`) SEBELUM status baris (draf terkirim / nomor asing tetap 403, bukan 409 / 404), dan untuk
    PUT / kirim juga di handler SEBELUM badan diurai (badan bukan JSON → 403, bukan 400). Jalur ditolak tidak menulis dan
    tidak mengirim email (diuji).
  - **TLS wajib**: port bukan 465 dan server tanpa STARTTLS → sambungan diputus SEBELUM MAIL FROM; jawaban **502 "server
    SMTP tidak mendukung TLS - email tidak dikirim; isi tersimpan, status tetap Not Sent"** (galat tersendiri, bukan 502
    umum). Pengecualian hanya `NBFACIN_SMTP_TANPA_TLS=true` (`strconv.ParseBool`; bukan boolean → gagal start) DAN akun
    kosong. Bawaan = wajib TLS. Diuji dengan server tiruan loopback tanpa STARTTLS: bawaan → ditolak; TANPA_TLS → terkirim;
    TANPA_TLS + akun → tetap ditolak. Jalur STARTTLS sungguhan tidak diuji (butuh sertifikat) `belum terverifikasi`.
  - Perubahan **inti**: `config.SMTP.TanpaTLS` + pembacaan `_TANPA_TLS` (uji `smtp_test.go`) — perlu diketahui tim inti.
  - Instrumen: mutasi (cek handler sebelum urai, cek layanan di lima metode, `aksesKorespondensi` selalu false, TLS tidak
    ditegakkan, TANPA_TLS mengabaikan akun) — ketujuh mutasi akses dan dua mutasi TLS menggagalkan uji. Ikut dibetulkan:
    uji 400 `TestKirimKorespondensi` sebelumnya bergantung pada urutan iterasi map (bisa gagal acak) — kini diperiksa per
    butir.

## Ditegaskan work owner 07-10-2026

"Tab Loss Record dan Correspondence itu cuma untuk IT aja , selain itu gausa munculin" - tetap: tab hanya tampil dan hanya
dapat diubah oleh akun ber-JOB_POSITION "IT"; akun lain tidak melihat tabnya. Saat case keluar dari tahap marketing, tab ini
ikut baca-saja (syarat TAHAP saja, tanpa workbasket marketing - keputusan c3 K57-20 dipertahankan).
