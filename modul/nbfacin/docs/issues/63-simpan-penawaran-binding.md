# 63: SAVE JSON_OFFER + FACINOFFER → marketing tahap Binding (jalur akhir Accept Underwriting)

> ⚠️ **Disusun agent sesi `nusantarare-c3` — bukan hasil `/to-tickets`.** Pemicu: K62-1 (tiket 62: "63 SAVE JSON_OFFER + FACINOFFER →
> marketing tahap Binding"). Riset korpus: scratchpad sesi 0f `riset-C-json.md` dan scratchpad sesi c3 `riset63\laporan.md`
> (subagen; tabel 50 kolom, SQL persis). Keputusan dikumpulkan sesi `nusantarare-9d` dari work owner 08-10-2026 (AskUserQuestion).

**Status:** backend selesai sesi c3 08-10-2026 (belum di-commit). Tanpa migrasi, tanpa CALL, tanpa endpoint baru — `POST
…/submit` tiket 61 dipakai apa adanya. Frontend tidak berubah (Accept di form UW tiket 61; case Binding baca-saja).

## Keputusan 08-10-2026

- **D1** `[keputusan work owner]` "B: baris + JSON minimal": baris `POOLDATA.JSON_OFFER` dibuat (ID, IDPEGA, NOOFFER, PRODKE 0,
  TGL_INPUT); JSONDATA = dokumen minimal berisi kunci yang dibaca SQL korpus (`IDNewBisnis`). Sejalan **ADR-0028** (penyusun
  muatan JSON tidak dimigrasikan) — **penyimpangan sadar** dari `@ASM.GetPageJSONString()`, dengan syarat: tidak ada sistem
  hilir yang membaca isi JSONDATA (pertanyaan terbuka: pemilik hilir).
  ⚠️ **RALAT K62-4** ("JSON_OFFER dibuat" = muatan penuh): ralat ini kesalahan sesi 9d — ADR-0028 tidak disebut saat work owner
  ditanya K62-4. Dicatat atas permintaan 9d.
- **D2** `[keputusan agent, ikut korpus + ADR-0043]`: prosedur `PEGA_M_JSON_OFFER` di-port ke Go (sequence NEXTVAL boleh, CALL
  tidak); upsert berkunci IDPEGA; satu transaksi dengan Submit (gagal → Submit batal); sequence tidak ada saat runtime → galat
  jelas.
- **D3** `[keputusan work owner]` "Ikut Pega + perbaiki cacat": 50 kolom seperti Pega, RATE berkoma (a); tanggal ditulis benar
  (b); literal "NULL" → SQL NULL (c); INSERT gagal → Submit batal (d); GROUPNAME = GROUP_NAME (e); GROUP_TYPE dari IsGroup.
- **D4** `[keputusan agent]`: POSITION ReasFacInMarketing, FLAG 2, IsCedingConfirm Binding, NB_STATUS pola marketing; email
  "Waiting for Bind" hanya di produksi (env) sesudah commit, BCC individu tidak diport; case Binding baca-saja sampai tiket 64.
- **D5** `[keputusan work owner]`: penanda group = When IsGroup kiriman work owner (`D:\migrasi\RNM\DDL\IsGroup.xml`) — "saya
  ada kirim di folder, when IsGroup. itu penanda yang masuk group". Literal login / kode marketing **tidak** ditulis di kode,
  tes, atau dokumen (CLAUDE.md §4.10): env `NBFACIN_GROUP_LOGIN` / `NBFACIN_GROUP_MARKETINGCODE`. Decision16 IsGroup → POLICY
  (tiket 64) ditahan.
- **D6** `[keputusan agent]`: `UpdateJSONOffer_ACT` dan tombol "Save to Json Offer" TIDAK diport.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\`

- `NB FacIn\Flow\InputInwardFacultativeOffer.xml`: Decision23 ELSE → Utility16 "SAVE JSON_OFFER" (Transition35, tanpa SET) →
  Decision16 (Transition59) → WHEN IsGroup → SubProcess1 POLICY (Transition60: PositionNote "ReasFacInAdmin", FlagOnGoingPolicy 1)
  / ELSE → Transition38 (PositionNote "ReasFacInMarketing", `.IsCedingConfirm = "Binding"`, `.FlagOnGoingPolicy = 2`, NBStatus =
  prefix + " IS IN " + MarketingName + "'S INBOX") → Utility6 SendEmailPolicy → Assignment9 "MARKETING (BINDING)" (ToWorkBasket
  ReasFacInMarketing). Diperiksa ulang sesi c3 (baris 2965-3032).
- `DDL\PEGA_M_JSON_OFFER.txt` (prosedur utuh): COUNT IDPEGA → INSERT (ID = '1' || JSON_OFFER_SEQ.NEXTVAL; NOOFFER 'F' =
  'RNM-OFR.' || MM.YYYY || '.' || LPAD(RNM_JSON_OFFER_SEQ.NEXTVAL, 5, '0'); PRODKE 0) / UPDATE JSONDATA; COMMIT sendiri.
- `NB FacIn\RDBList\InsertOfferProduction_Sql.xml`: DELETE ber-IDPEGA = pyID, INSERT 50 kolom per baris SpreadingList (cabang
  FIRE `Activity\SaveOfferProduction_Act.xml` langkah 12). Kolom dihitung dua cara (pisah koma = 50, token = 50); daftar kolom
  kode = daftar INSERT korpus (urutan sama, A\B = B\A = ∅). DDL FACINOFFER 52 kolom; TSITOP_MENJADI / TSITOP_SELISIH tidak diisi.
- `DDL\IsGroup.xml` (ASM-FW-GISFW-WORK!ISGROUP, "A OR B OR C OR D OR E"): A `.OfferFacIn.QuotationData.IsGroup = "Group"`;
  B/C/D `OperatorID.pyUserIdentifier` = 3 literal login; E `QuotationData.MarketingCode` = 1 literal.
  `NB FacIn\Activity\InputOfferFacInEngineer_preACT.xml` langkah 23-25: IsGroup "Group" bila `@equalsIgnoreCase(pxCreateOperator,
  …)` salah satu 3 login, selain itu "NonGroup" ("SemiGroup" → "NonGroup"). Himpunan login langkah 23 = himpunan B-D (dibandingkan
  tanpa mencetak nilai: 3 = 3, irisan 3).
- `NB FacIn\Activity\SendEmailPolicy.xml`: langkah 8 Subject "Waiting for Bind ID Pega : " + pyID + " - " + InsuredName
  (FlagOnGoingPolicy 2); To = email pembuat case (bukan "@nusantarare.com" / kosong → CC); CC = MO + atasan (`GetMKTandLeader_SQL`);
  langkah 14 kirim hanya bila IsPEGAPROD.

## Backend sesi c3 (08-10-2026)

- **K63-B1 Jalur akhir** — `services/submituw.go` `actFlow`: Decision23 Else (`acceptance.Next` Selesai) dan posisi tanpa jabatan
  tabel limit (UW Financial; Pega: login tanpa baris → SQL kosong → Else, keputusan NB-11) → `tujuanBinding` (`services/binding.go`).
  Tangga tidak pasti (`galatTanggaElse`: seri, bentuk B, BusinessOldId tak diketahui, …) tetap DITAHAN — pesan baru
  "Jenjang approval berikutnya tidak dapat ditentukan dari tabel limit; hubungi administrator." `[keputusan agent]` (tiket 62
  menahannya dengan pesan "jalur akhir belum tersedia" yang kini tidak benar lagi).
- **K63-B2 Group (D5)** — `services/grup.go`: A = `T_WORK_POLIS.CREATE_OP` di `NBFACIN_GROUP_LOGIN` (tidak peka huruf, seperti
  equalsIgnoreCase langkah 23; dihitung saat itu, tanpa migrasi); B-D = akun yang memproses di daftar yang sama (`[keputusan agent]`
  tidak peka huruf; Pega "=" `[dugaan]` peka huruf); E = `MARKETINGOFFICER.CLIENTID` lewat MOID (`GetMOID_Sql clientid =
  MarketingCode`) di `NBFACIN_GROUP_MARKETINGCODE`, persis. Akibat: (i) tiket 62 ActFlow langkah 3 "exit kalau group" berlaku
  (`Pengguna.AnggotaGrup`) → Else; (ii) Decision16 IsGroup → DITAHAN "Accept meneruskan case ke jalur Policy (case group) yang belum
  tersedia." tanpa tulisan; (iii) FACINOFFER.GROUP_TYPE = "1" bila A, selain itu "2" — praktis selalu "2", sebab case group
  ditahan sebelum menulis (di Pega Utility16 jalan SEBELUM Decision16).
- **K63-B3 Tulisan (satu transaksi Submit)** — `repository/binding.go` `tulisBinding`, dipanggil `TulisSubmit` sesudah posisi:
  `T_WORK_POLIS.FLAG_ONGOING_POLICY = 2`; JSON_OFFER: `SELECT COUNT(*), MAX(NOOFFER) … WHERE IDPEGA` → INSERT (NOOFFER dulu, lalu ID,
  urutan prosedur; LPAD memotong ke 5 karakter — ditiru) / UPDATE JSONDATA; FACINOFFER: DELETE IDPEGA = ID case, INSERT per baris
  spreading. JSONDATA = `{"IDNewBisnis":"<ID case>"}`. Bulan nomor offer = bulan WIB jam Submit `[keputusan agent]` (prosedur:
  SYSDATE server basis data). ORA-02289 → `ErrSequenceJSONOffer`.
- **K63-B4 Pemetaan FACINOFFER** — `services/binding.go` `barisFacInOffer`: objek (`BacaObjek`) dan spreading (`bacaSpreading`)
  dipasangkan per indeks lokasi / item / coverage + OLDID (tidak sejalan → galat, Submit batal). ROADNAME = RoadType + RoadName
  (tanpa spasi, seperti sumber); BUILDINGNO = "NO. " + nomor; OCCUPATIONCODE = id baru di DEPAN (`.OccupationId + " || " +
  CARI29`); BEGIN/END = jam dinding WIB 24-jam (D3 b); TSI_MENJADI = TSI_SELISIH = TSISpreaded, PREMI_* = PremiumSpreaded;
  MARKETINGOFFICERCODE = MOID tersimpan `[keputusan agent]` (Pega mencari ulang MARKETINGOFFICER.id lewat CLIENTID + MOSTATUS '1');
  BUSINESSCODE = BUSINESS.ID Class of Business case (butir 89); TOP_RISK "true" / "false"; IsCedingConfirm tidak punya kolom —
  status Binding = FLAG 2 (kotak masuk).
- **K63-B5 Email** — sesudah commit, hanya bila `IS_PEGA_PROD=true` dan `NBFACIN_SMTP_*` terpasang; galat hanya dicatat (Pega
  tidak membaca hasil kirim). Penerima dari `M_LOGIN_GO.EMAIL` (pembuat case = CREATE_OP; marketing + atasan lewat
  `MARKETINGOFFICER.AKSES_LOGIN`). BCC (2 alamat individu ter-hardcode) tidak diport.
- **K63-B6 Config** — `inti/backend/config`: `NBFACIN_GROUP_LOGIN`, `NBFACIN_GROUP_MARKETINGCODE` (dipisah koma; `String()`
  hanya jumlah); `.env.example` ditambah. Catatan: `NBFACIN_SMTP_*` (tiket 56) tidak tercatat di `.env.example` — penjaga
  `TestSetiapEnvTerdokumentasi` hanya menangkap `os.Getenv("LITERAL")`, sedangkan SMTP dibaca lewat awalan; dilaporkan, tidak diubah.

**Uji** — `TestKasusGroup`, `TestAcceptGroupDitahan`, `TestBarisFacInOffer`, `TestSubmitBinding` (DEV tanpa email / produksi
satu email; batas bulan WIB), `TestSurelBinding`, `TestAcceptJenjangLanjut` (17 kasus: Else → Binding, tak pasti → ditahan;
alasan tiap kasus diperiksa dari log), `TestSQLBinding`, `TestNoOfferBaru`, `TestGalatSequence`, `TestGrupNBFacIn` (config).
gofmt / vet bersih; `go test ./modul/nbfacin/... ./inti/backend/config/` hijau; rules dengan `KORPUS_RNM_AKAR` hijau; penjaga
hanya kegagalan lama premiumlistlife; index git 0. **Mutasi 27 / 28 mati**; hidup: pemetaan scan `MarketingCode` di
`KonteksSubmit` (jalur Oracle, hanya teruji dengan basis data — seperti medan scan lainnya). `tulisBinding` (cabang INSERT /
UPDATE JSON_OFFER) juga hanya teruji teks SQL-nya; uji basis data belum ada.

## Deploy

- Tanpa migrasi. Biner baru MENULIS tabel lama `POOLDATA.JSON_OFFER` dan `POOLDATA.FACINOFFER` (INSERT / UPDATE / DELETE) dan
  memakai `POOLDATA.JSON_OFFER_SEQ` / `RNM_JSON_OFFER_SEQ` — akun aplikasi butuh hak itu (DBA).
- Env baru (opsional): `NBFACIN_GROUP_LOGIN`, `NBFACIN_GROUP_MARKETINGCODE` — kosong = tidak ada case group.
- Email "Waiting for Bind" hanya terkirim bila `IS_PEGA_PROD=true` dan `NBFACIN_SMTP_*` terisi.

## Terbuka

- `[pertanyaan terbuka — DBA]` keberadaan dan hak akses `JSON_OFFER_SEQ` / `RNM_JSON_OFFER_SEQ` (DDL sequence tidak ada di korpus).
- `[pertanyaan terbuka — Product]` nilai `Quotation.BusinessFac` NB FIRE = 'F'? (cabang NOOFFER 'F' dipakai; 'T' / 'L' tidak).
- `[pertanyaan terbuka — pemilik hilir]` tidak ada sistem yang membaca isi `JSON_OFFER.JSONDATA` (syarat D1); kolom CEDING …
  UW_NOTE JSON_OFFER tidak ditulis rule mana pun di korpus — pengisinya tidak diketahui.
- `[pertanyaan terbuka — pemilik export Pega]` badan email `NotifEmailPolicy_HTML` tidak ada di korpus (hanya dirujuk); badan
  ringkas `[keputusan agent]`.
- `[pertanyaan terbuka — IAM]` padanan JABATAN UW Financial (tiket 62); kini Accept UW Financial → Decision23 Else → Binding.
- ⚠️ **Temuan keamanan untuk pemilik Pega:** `NB FacIn\Activity\SetNBStatus_Act.xml` langkah 9 (nonaktif) memuat host SMTP dan
  kredensial dalam teks biasa. Nilainya tidak disalin ke artefak mana pun.

## Catatan silang — prakondisi `.pxListSubscript>1` (11-10-2026, tiket 74 P5)

`[terverifikasi]` Langkah INSERT `SaveOfferProduction_Act` (FACINOFFER) di cabang ini berprakondisi `.pxListSubscript>1` (T=2 lanjut / F=3 lewati), sama dengan
setiap cabang lain (FIRE / Aneka / Marine). Port sistem baru menulis SEMUA baris spreading. **Keputusan work owner 11-10-2026
(K74-23, `docs/issues/74-marine-cargo.md`):** tetap begitu — kode tiket ini tidak diubah.

## Catatan silang — IDPEGA = case id polos (11-10-2026)

Tulisan IDPEGA tiket ini kini memakai case id polos; pembaca menerima dua bentuk (polos / berawalan). Rincian: `docs/issues/69-jalur-policy.md`
bagian "IDPEGA = case id polos".
