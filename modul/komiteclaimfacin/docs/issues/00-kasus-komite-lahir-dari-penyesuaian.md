# 00: Kasus komite lahir dari penyesuaian — jenjang dibekukan

**Status:** ready-for-agent
**Blocked by:** `claim-facin\issues\13` *(kontrak muatan)*
**Menutup:** AC 1 · 2 · 3 · 4 · 5 · 6 · 7 · 8 · 96 · 97 · 98 · 99 · 100 *(13 AC)* — US 1–5

## Hasil & nilai pengguna

Hari ini Penyesuaian di atas kewenangan penilai **belum dapat berubah menjadi kasus komite**, dan tidak ada yang menetapkan **berapa jenjang** harus menyetujui.

Sesudah tiket ini, Penilai menyerahkan penyesuaian, dan **kasus komite lahir** dengan **jumlah jenjang yang ditetapkan susunan jenjang komite**. ⭐ Tiap jenjang **dibekukan bersama jabatan dan akun pemegangnya** — pergantian pemegang jabatan sesudah itu **tidak memindahkan kasus**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Hari ini Penyesuaian di atas kewenangan penilai
> **belum dapat berubah menjadi kasus komite**, dan tidak ada yang menetapkan **berapa jenjang** harus menyetujui."* dan
> *"⭐ Tiap jenjang **dibekukan bersama jabatan dan akun pemegangnya** — pergantian pemegang jabatan sesudah itu **tidak
> memindahkan kasus**."* →
>
> - **Kelahiran kasus `KMT-` TT2 sudah dibangun Claim Fac In tahap 1** (`CreateKMTNo_Act`, commit `8c3b2e71`,
>   `modul/claimfacin/docs/PARITAS.md` §5): `T_WORK_CLAIM` `KMT-` (`COVER_KEY` = kasus klaim, `LINI` FACIN, `TAHAP`
>   `Komite_Flow`, `POSITION` = `OPERATOR_ID` tingkat 1), `T_GENERAL_KOMITE`, dan tangga awal dari roster calon
>   (`RosterKomiteCalon` = `SetListKomite_act`: non-retro `LIMIT_BOTTOM ≤ 1`, retro `ValueAdjustment` dibatasi `LIMIT_TOP`
>   Technic Div. Head). Modul komite **tidak** melahirkan kasus TT2; ia membaca kasus yang sudah lahir. "Blocked by
>   `claim-facin\issues\13`" sudah terpenuhi.
> - **"Jenjang dibekukan" diganti KCF-02** (prompt tahap 2 §3): perluasan tangga ikut `ApprovalKomite_Act` L4–L8. Bila
>   `KomiteCount = 1` dan tangga baru satu baris, anggota roster FACIN aktif DEGREE > 1 ber-`LIMIT_BOTTOM` < total
>   adjustment ditambahkan (urut DEGREE), calon yang sama dengan anggota tingkat 1 dibuang (L6.2), lalu `KOMITE_LOOP`
>   dihitung ulang (L8). Perluasan **dihitung saat tingkat 1 membuka kasus** (tampil di List of Committee, tanpa menulis
>   saat GET) dan **disimpan saat tingkat 1 Submit**. Kasus Fac Retro melewati perluasan (L3); tangganya sudah dibentuk
>   sisi klaim.
> - **Pemegang tidak dibekukan.** `KOMITE_OPERATORID` / `POSITION` berisi workbasket (KCF-01), dan yang boleh memutus
>   adalah anggota workbasket tingkat berjalan saat Submit (pola Komite Prop 09-10-2026). Pergantian anggota workbasket
>   memang memindahkan siapa yang dapat memutus.
> - Dari baris `Menutup`: AC 5 dan AC 6 **gugur**; AC 2 / 4 / 7 / 8 / 96 / 97 dibaca menurut RALAT bertanggal sama di
>   `spec.md`; AC 1 / 3 / 98 / 99 / 100 tetap.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Pembentukan kasus | rule pembuat nomor komite menyemai **jumlah jenjang** dan **giliran mulai**, lalu membuat kasus anak |
| Jalur satu jenjang | rule kirim tutup klaim dan kirim tolak klaim — ⭐ keduanya menyemai **satu jenjang** |
| Penguncian induk | modul komite **mengunci kasus induk** sebelum mengerjakan apa pun |
| Empat jalur masuk | penyesuaian · tolak · tutup klaim · survey — ⛔ **jalur survey tidak dialihkan**, rule-nya tidak ada di korpus mana pun |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"Pembentukan kasus | rule pembuat nomor komite
> menyemai **jumlah jenjang** dan **giliran mulai**, lalu membuat kasus anak"* dan *"Jalur satu jenjang | rule kirim
> tutup klaim dan kirim tolak klaim — ⭐ keduanya menyemai **satu jenjang**"* →
>
> - Rule pembuat nomor komite = `CreateKMTNo_Act` (korpus Claim Fac In), **sudah dibangun claimfacin tahap 1** (lihat
>   RALAT di atas).
> - **TT3 / TT4 → KCF-03.** `SendRejectClaimToKomite2` (TT3) dan `SendCloseClaimToKomite` (TT4) memang menyemai satu
>   anggota, tetapi anggotanya akun orang tertulis mati; sistem baru memakai **`ReasClaimDeptHead`** (prompt §5 #7).
>   Kasusnya **tanpa adjustment**, jadi migrasi komiteclaimfacin (640–679) `MODIFY T_GENERAL_KOMITE.ADJUSTMENT_ID` menjadi
>   boleh kosong dan menambah kolom jenis penyerahan **`TRANSFER_TYPE`** (2 / 3 / 4 = `TransferType` XML), izin work
>   owner. Tombol **Yes** pada `SureRejectClaim` / `PreventRejectClaim` di claimfacin dinyalakan (tahap 1 nonaktif,
>   OQ-CFI-27). Hanya Fac In; Prop / Non Prop menyusul.
> - Empat jalur masuk cocok dengan `KomitePostAct`: S2 TT2 → `KomitePost_Adjustment`, S4 TT3 → `KomitePost_Reject`, S5
>   TT4 → `KomitePost_CloseClaim`; S3 TT1 (survey) ter-remark.

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K13** — ⭐ Jenjang dibekukan bersama **akun pemegangnya**, bukan hanya jabatannya
- **K6** — ⭐ Komite lini FAC dan lini PROP memakai **tabel yang sama** — ⛔ nol tabel baru

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"**K13** — ⭐ Jenjang dibekukan bersama **akun
> pemegangnya**, bukan hanya jabatannya"* dan *"**K6** — ⭐ Komite lini FAC dan lini PROP memakai **tabel yang sama** —
> ⛔ nol tabel baru"* →
>
> - **K13 gugur** (prompt tahap 2 §3 "Gugur"): digantikan **KCF-01** (roster `EMAILKOMITE` STS_KLAIM FACIN → workbasket
>   per DEGREE, JABATAN dan LIMIT tetap) dan **KCF-02** (perluasan tangga).
> - **"K6" di sini bukan K6 spec modul ini** (spec ID-9 = surel galat kasir). Ia nomor `STRUKTUR-TABEL-CLAIM-FACIN.md`
>   §5d (`modul/claimfacin/docs/`), sejalan §0 K4. Nol tabel baru tetap berlaku, tetapi `T_GENERAL_KOMITE` kini **diubah**
>   (KCF-03: `MODIFY ADJUSTMENT_ID` + kolom `TRANSFER_TYPE`), tidak lagi dipakai "apa adanya".

## Yang harus diuji

- [ ] Kasus komite lahir dengan **jumlah jenjang benar** dari susunan jenjang
- [ ] Giliran mulai pada **jenjang pertama**
- [ ] ⭐ Jenjang **dibekukan**: pergantian pemegang jabatan sesudah pembentukan **tidak memindahkan kasus**
- [ ] ⛔ Jabatan **tanpa pemegang aktif** ⇒ pembentukan **ditolak dengan galat yang menyebut jabatannya**
- [ ] ⭐ Jalur **tutup klaim** dan **tolak klaim** membentuk komite **satu jenjang**
- [ ] ⭐ Seluruh data kutipan yang dibutuhkan penggolongan **tersalin dari kasus induk**
- [ ] ⛔ Jalur **survey** **tidak dibangun**

> ⛔ **RALAT 10-10-2026.** Butir lamanya dikutip utuh, tidak dihapus: *"Kasus komite lahir dengan **jumlah jenjang
> benar** dari susunan jenjang"*, *"⭐ Jenjang **dibekukan**: pergantian pemegang jabatan sesudah pembentukan **tidak
> memindahkan kasus**"*, *"⛔ Jabatan **tanpa pemegang aktif** ⇒ pembentukan **ditolak dengan galat yang menyebut
> jabatannya**"* dan *"⭐ Seluruh data kutipan yang dibutuhkan penggolongan **tersalin dari kasus induk**"* →
>
> - Jumlah jenjang saat lahir diuji claimfacin. Yang diuji di sini **perluasan** KCF-02: tangga satu baris + total
>   adjustment ⇒ baris DEGREE > 1 ber-`LIMIT_BOTTOM` < total; **pita SPV B** (`ApprovalKomite_Act` L5, KCF-01: pemutus
>   tingkat 1 anggota `ReasClaimSPVB` dan 30.000.000 < total ≤ 57.750.000) ⇒ hanya satu jenjang atas (DEGREE > 1
>   pertama); Fac Retro ⇒ tanpa perluasan (L3).
> - "Dibekukan" **gugur** (lihat RALAT di kepala tiket).
> - "Tanpa pemegang aktif ⇒ ditolak" **tidak dibangun**: nol bukti XML. `ApprovalKomite_Act` L4 hanya menyaring roster
>   `STS_AKTIF = "1"`; workbasket tanpa anggota aktif berarti kasus menunggu.
> - **Komite tidak menyalin data kutipan.** Halaman klaim induk, termasuk `OfferFacIn`, dibaca lewat kontrak
>   `kontrak.KlaimFacInKomite` (`BacaKlaimFacIn`, `inti/backend/kontrak/klaimfacin.go`, disediakan claimfacin).
>   `SetValueKomite` S12 hanya menyalin `BusinessType` / `GroupPanel` ke `pyWorkPage.Quotation`; sistem baru tidak
>   butuh salinan itu.
> - "Giliran mulai pada jenjang pertama", "satu jenjang" (kini `ReasClaimDeptHead`), dan "survey tidak dibangun" tetap.

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **39** | ⚠️ **Satu jabatan dipegang LEBIH DARI SATU orang — siapa menerima giliran?** ⭐ Usul asisten: tolak dengan galat yang jelas, ⛔ jangan memilih diam-diam | ⚠️ **MENAHAN pembentukan kasus** bila keadaan itu terjadi |
| **6** | Isi daftar jabatan dan susunan jenjang | ⛔⛔ **MENAHAN** penentuan jumlah jenjang |

> ⛔ **RALAT 10-10-2026.** Sel lamanya dikutip utuh, tidak dihapus: *"⚠️ **MENAHAN pembentukan kasus** bila keadaan itu
> terjadi"* dan *"⛔⛔ **MENAHAN** penentuan jumlah jenjang"* → **keduanya tidak menahan lagi.** Butir **39** tertutup
> pola Komite Prop 09-10-2026: penyetuju = anggota workbasket tingkat berjalan, siapa pun anggotanya boleh memutus, surel
> ke semua anggota. Butir **6** tertutup ADR-0030 (`OUTPUT_HASIL_RNM/docs/bersama/adr/0030-aturan-peran-ditetapkan-sekali-lintas-modul.md`)
> + KCF-01: roster FACIN DEGREE 1–5 (Claim Supervisor / Claim Dept. Head / Technic Div. Head / Operational Director /
> Technical Director) → `ReasClaimSPVA` (cadangan `ReasClaimSPVB`) / `ReasClaimDeptHead` / `ReasClaimTechDivHead` /
> `ReasClaimOpsDir` / `ReasClaimTechDir`.

## Seam & verifikasi

**Seam:** lapisan layanan komite — ⭐ **satu pintu masuk**, bukan layar.
2. Serahkan satu penyesuaian — ⭐ kasus komite lahir, jumlah jenjang benar.
3. Ganti pemegang jabatan **sesudah** pembentukan — ⭐ kasus **tidak berpindah**.
4. Kosongkan pemegang salah satu jabatan — ⛔ pembentukan **ditolak**, galatnya **menyebut jabatannya**.
5. Serahkan lewat jalur tutup klaim — ⭐ komitenya **satu jenjang**.

> ⛔ **RALAT 10-10-2026.** Langkah lamanya dikutip utuh, tidak dihapus: *"3. Ganti pemegang jabatan **sesudah**
> pembentukan — ⭐ kasus **tidak berpindah**."* dan *"4. Kosongkan pemegang salah satu jabatan — ⛔ pembentukan
> **ditolak**, galatnya **menyebut jabatannya**."* → langkah 3 dan 4 **gugur** (KCF-01, nol bukti XML). Penggantinya:
> buka kasus sebagai anggota `ReasClaimSPVA` dan `ReasClaimSPVB` ⇒ List of Committee menampilkan perluasan tanpa baris
> tersimpan; Submit tingkat 1 ⇒ tangga dan `KOMITE_LOOP` tersimpan. Langkah 5 = TT4 Close Without Payment, satu tingkat
> `ReasClaimDeptHead` (KCF-03).
