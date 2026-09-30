# 06: Security reinsurer — retrosesi atas retrosesi dan eksposur berjenjang

**Status:** ready-for-agent

**Blocked by:** 05 (security reinsurer lahir di bawah reinsurer)

## Hasil & nilai pengguna

Sebagai **admin master**, saya menambahkan **security reinsurer** di bawah seorang reinsurer; dan
sebagai **underwriter** saya melihat **eksposur efektifnya terhadap treaty** — bukan angka mentah
yang menyesatkan. *(User story 18–20 di spec)*

⚠️ Inilah tiket dengan risiko salah-baca terbesar di modul ini.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas security reinsurer; rujukan ke reinsurer induk |
| `internal/repository` | Pemanggilan procedure penulis security reinsurer |
| `internal/services` | **Perhitungan eksposur berjenjang** — share anak × share induk |
| `internal/handlers` | Endpoint CRUD; eksposur efektif pada respons |
| `frontend/` | Grid security reinsurer; **eksposur efektif ditampilkan berdampingan** dengan share mentah |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatySecurityReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYSECURITYREINSURER_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveMasterTreatySecurityReinsurer_Life_SQL.xml` | `POOLDATA.INSERTSECURITYREINSURER_LIFE` |
| `SaveSecurityReinsurerLife_Act` | `ASM-FW-GISFW-…` / `SAVESECURITYREINSURERLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveSecurityReinsurerLife_Act.xml` | orkestrator simpan |
| `SetSecurityReinsurerLife_Act` | `ASM-FW-GISFW-…` / `SETSECURITYREINSURERLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetSecurityReinsurerLife_Act.xml` | isi form |
| `BrowseSecurityReinsurer_Life_RD` | `ASM-FW-GISFW-INT-TREATYSECURITYREINSURER_LIFE` / `BROWSESECURITYREINSURER_LIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseSecurityReinsurer_Life_RD.xml` | daftar |
| `InboxSecurityReinsurerLife` | `DATA-PORTAL` / `INBOXSECURITYREINSURERLIFE` / `RULE-HTML-HARNESS` | `Master Contract Retro Life/Harness/InboxSecurityReinsurerLife.xml` (444.874 byte) | layar grid |

`[terverifikasi]` `INSERTSECURITYREINSURER_LIFE` menerima **`p_TREATYREINSURERID`** — penunjuk baris
**reinsurer induk**. `[data DBA]` Kini ditegakkan **FK** di basis data
(`TREATYSECURITYREINSURER_LIFE.TREATYREINSURERID` → `TREATYREINSURER_LIFE.ID`).

`[fakta bisnis — work owner]` **`PCTSHARE` security reinsurer adalah persentase DARI SHARE REINSURER
INDUKNYA**, **bukan** dari keseluruhan treaty. Ini **retrosesi atas retrosesi**.

> Reinsurer A memperoleh **40%** treaty. Security Reinsurer X di bawah A memperoleh **10%**.
> **Eksposur X terhadap treaty = 10% × 40% = 4%.**

`[data DBA]` `PCTSHARE` bertipe **`NUMBER`** — desimal, bukan teks.

## ADR terkait

**ADR-0003** (persentase non-float), **ADR-0006**, **ADR-0007**,
**ADR-0001** (angka eksposur ini dikonsumsi konteks hilir).

## Acceptance criteria

- [ ] Security reinsurer lahir **di bawah** satu reinsurer; tanpa reinsurer induk **ditolak**.
      *(AC 21 spec)*
- [ ] ⚠️ `PCTSHARE` security reinsurer dibaca sebagai **persentase dari share induknya**, **bukan**
      dari treaty. *(AC 22 spec; `[fakta bisnis — work owner]`)*
- [ ] **Eksposur efektif** terhadap treaty dihitung **share anak × share induk**, dan **itulah** angka
      yang ditampilkan sebagai eksposur. Test memuat kasus **`10% × 40% = 4%`**. *(AC 23 spec)*
- [ ] Mengubah share **induk** mengubah eksposur efektif **seluruh anaknya** — tanpa menyentuh baris
      anak. *(AC 24 spec)*
- [ ] Layar menampilkan **share mentah dan eksposur efektif berdampingan**, dengan label yang
      membedakan keduanya — sehingga `10%` tidak pernah terbaca sebagai eksposur terhadap treaty.
- [ ] `PCTSHARE` di luar rentang **0–100 ditolak**. *(sejalan AC 16 spec)*
- [ ] Share diperlakukan sebagai **desimal presisi arbitrer**; **tidak** melewati `float`.
      *(AC 49 spec; **ADR-0003**)*
- [ ] Menyimpan security reinsurer yang sudah ada = **upsert**, bukan baris kedua; identitas baru
      dibuat basis data. *(AC 45, 4 spec)*
- [ ] `o_message` diperiksa; kegagalan ditampilkan. *(AC 46 spec — HTML dan konformansi di
      **tiket 10**)*

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Jebakan terbesar modul ini.** Menampilkan `10%` tanpa konteks induknya **salah besar** — ia
tampak empat kali lebih besar dari eksposur sebenarnya. Angka ini merambat ke perhitungan klaim dan
eksposur di **Claim Life** dan **Komite Claim Life**.

⚠️ **Anti-dobel logis belum diputuskan** — sama seperti tiket 05. `[data DBA]` Basis data tidak
mencegah dua security reinsurer yang sama di bawah satu induk.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — kaskade FK dan perilaku upsert hanya
berperilaku benar pada basis data sungguhan.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Ralat bertanggal 30-09-2026 — sesi implementasi (paket 0)

> Sumber: `RALAT-DEV-30-09-2026.md` (K1–K8 katalog DEV, R1–R12 pembacaan ulang XML) dan `PARITAS-LAYAR-DAN-AKSI.md`. Kalimat di atas **tidak dihapus**; yang berlaku adalah ralat ini.

| Kalimat lama | Ralat |
| --- | --- |
| *"Layar menampilkan **share mentah dan eksposur efektif berdampingan**"* | kolom eksposur tidak ada di grid Pega — ditambahkan dengan label tiket, OQ-MCRL-08 (R11) |
| *"ditegakkan **FK** di basis data"* | nol FK di DEV (K1); induk diperiksa Go |
| — | `Add` Pega mengosongkan halaman reinsurer, bukan form security — sistem baru mengosongkan form security (OQ-MCRL-12); `(%) SHARE` 0..100 (R3) |
