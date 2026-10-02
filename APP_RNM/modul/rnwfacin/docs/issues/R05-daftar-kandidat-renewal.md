# R05: Daftar kandidat renewal dan portalnya

**What to build:** Underwriter membuka daftar polis yang mendekati jatuh tempo dan memilih mana yang
akan diperpanjang — tanpa perlu tahu nomor polisnya lebih dulu. Daftar yang sama muncul di portal.

`[terverifikasi]` `ReportDefinition\RenewalList_RD` — **11 kolom**:

```
.pyID · .pzInsKey · .pxCreateDateTime · .pxCreateOperator
.Quotation.OldPolicyNo · .Quotation.InsuredName · .Quotation.MarketingName
.OfferFacIn.PolicyData.StartDateTime · .OfferFacIn.PolicyData.EndDateTime
.NBStatus · .NBStatusNew
```

Berkunci **`OldPolicyNo`** dan tanggal berakhirnya polis — itulah yang membuatnya daftar *kandidat*,
bukan sekadar daftar kasus.

`[terverifikasi]` `Section\SFAPortal_Renewal` mengikat 9 properti dengan bentuk yang sama.

### ⚠️ `Work-Renewal` adalah irisan pelaporan, bukan kelas kerja

`[terverifikasi]` `RenewalList_RD` berkelas `ASM-FW-GISFW-Work-Renewal`, **tetapi** `pyWorkClass` pada
flow renewal adalah `ASM-FW-GISFW-Work` — **sama dengan NB** — dan `Work-Renewal` tidak dipakai rule
lain mana pun.

⛔ **Jangan membuat kelas kerja atau antrean tersendiri untuk renewal.** Kasus renewal hidup di kelas
kerja yang sama dengan NB; `Work-Renewal` hanya cara melaporkannya.

**Blocked by:** R01

**Status:** ready-for-human — layanan diport 01-10-2026; premis "daftar kandidat" dikoreksi (A40, dikonfirmasi work owner — butir 56 nbfacin); portal dan layar belum · ⏸ pengerjaan rnwfacin dihentikan work owner 01-10-2026 (`PROMPT-LANJUT-IMPLEMENT-NB-SAJA.md`)

- [ ] Daftar kandidat menampilkan kesebelas kolom, berkunci nomor polis lama dan tanggal berakhir
- [ ] Portal menampilkan daftar yang sama
- [ ] Memilih satu baris mengarah ke alur masuk renewal (R01) dengan nomor polis terisi
- [x] **Tidak ada kelas kerja atau antrean baru** untuk renewal — kasus tetap di kelas kerja NB
- [x] Penyaringan status dan kelompok tim direproduksi sesuai perilaku lama
- [x] ⚠️ Kolom nama tertanggung dan nama marketing **ditampilkan sesuai perilaku lama**, tetapi **nilainya tidak pernah masuk fixture, test, atau artefak mana pun**

## Comments

### 2026-10-01 — layanan diport (agent); premis dikoreksi

**Kode:** `backend/services/daftar.go` (`ServiceDaftar.Daftar`, `Saringan`, `HasilDaftar`), model
`backend/models/daftar.go` (`BarisDaftarRenewal`, `StatusKerjaSelesai`/`Ditolak`), sumber
`repository.SumberDaftarRenewal` (tabel menunggu DBA). Tes `services/daftar_test.go`, nilai sintetis berawalan `UJI-`; uji mutasi 13/14 (yang lolos ekuivalen: pesan galat pembuat kosong tidak memuat apa pun).

⚠️ **A40 — `RenewalList_RD` bukan daftar kandidat jatuh tempo.** `[terverifikasi]` filternya (`pyContent/pyFilters`,
`B AND A AND C AND D`): pembuat = `Param.UserNameID`, status kerja bukan `Resolved-Completed`/`Resolved-Rejected`,
`.Quotation.TeamGroup = Param.TeamGroup`. **Tidak ada** filter tanggal berakhir atau `OldPolicyNo` — keduanya kolom tampil.
Urut `pxCreateDateTime` DESC lalu `pyID` DESC; 500 baris (`pyMaxRecords`). *(Ralat 01-10 malam: ringkasnya = kasus renewal milik pembuatnya di team group yang sama, **kecuali** Resolved-Completed dan Resolved-Rejected — status Resolved-* lain tetap lolos.)* Baris = kasus renewal yang **sudah ada**
(`pyID`, `pzInsKey`), jadi kriteria "memilih baris mengarah ke alur masuk R01" kemungkinan keliru — `[dugaan]` memilih
baris membuka kasus itu. Kriteria 1, 2, 3 (tampilan, portal, navigasi) = layar, belum (butir 55).

`[pertanyaan terbuka]` padanan status kerja Pega di sistem baru; perilaku Pega saat parameter kosong (`pyUseNullIfEmpty`
tidak diisi — layanan menolak dengan `ErrSaringanKosong`); arti `NBStatus`/`NBStatusNew`. `Param.UserNameID` disediakan
pemanggil dari sesi — layanan tidak menyentuh autentikasi.

### 2026-10-01 — code review dua sumbu (agent)

Spec: A40 dikonfirmasi reviewer (`RenewalList_RD.xml` `pyUIFilters` L606–686, `pyFilters` L1063–1143; parameter rule hanya
`TeamGroup` dan `UserNameID`). Diperbaiki: baris berstatus kerja kosong kini **ditolak** (`ErrStatusKerjaKosong`,
`[dugaan]` SQL `<>` Oracle membuang '' = NULL) alih-alih diloloskan; `ErrSaringanKosong` dinyatakan sebagai
**penyimpangan sadar** (`[dugaan]` Pega melewati filter berparameter kosong — tanpa filter A daftar memuat kasus pengguna
lain). Standar: pesan galat tidak lagi memuat login pembuat (CLAUDE.md §4.10); fixture berawalan `UJI-` (README-BACA-DULU
§4, juga `masuk_test.go` dan `uji/lintasmodul/renewal_tangga_test.go` R01); angka sensus disertai cara hitungnya.
Dicatat, tidak diubah: sumber mengembalikan semua baris dan batas 500 diterapkan di memori — pemuat Oracle kelak boleh
mendorong filter ke SQL selama tes ini tetap hijau.
