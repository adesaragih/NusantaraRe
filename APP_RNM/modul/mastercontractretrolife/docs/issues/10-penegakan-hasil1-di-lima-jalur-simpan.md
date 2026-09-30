# 10: Penegakan `HASIL1` di kelima jalur simpan — gagal terang-terangan, HTML dibersihkan

**Status:** ready-for-agent

**Blocked by:** 01 (tahun), 02 (kontrak), 05 (reinsurer), 06 (security reinsurer), 07 (business) —
kelima jalur simpan harus ada agar konformansi dapat diuji

## Hasil & nilai pengguna

Sebagai **admin master**, saya **diberi tahu bila penyimpanan ditolak** basis data — dengan pesan
yang **terbaca manusia**, bukan potongan markup — sehingga saya tidak pernah mengira data tersimpan
padahal tidak. *(User story 32–33 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Pemeriksaan hasil di **kelima** jalur simpan; **pembersihan HTML** di satu tempat |
| `internal/services` | Pemetaan hasil basis data menjadi galat domain |
| `internal/handlers` | Galat sampai ke pemanggil sebagai kegagalan, bukan sukses |
| `frontend/` | Pesan galat tampil bersih di dekat form yang bersangkutan |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Procedure |
| --- | --- | --- |
| `SaveMasterTreatyYear_Life_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `ASM!SAVEMASTERTREATYYEAR_LIFE_SQL` / `RULE-CONNECT-SQL` | `INSERTTREATYYEAR_LIFE` |
| `SaveMasterTreatyContract_Life_SQL` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` / `ASM!SAVEMASTERTREATYCONTRACT_LIFE_SQL` / `RULE-CONNECT-SQL` | `INSERTTREATYCONTRACT_LIFE` |
| `SaveMasterTreatyReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYREINSURER_LIFE_SQL` / `RULE-CONNECT-SQL` | `INSERTREINSURER_LIFE` |
| `SaveMasterTreatySecurityReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYSECURITYREINSURER_LIFE_SQL` / `RULE-CONNECT-SQL` | `INSERTSECURITYREINSURER_LIFE` |
| `SaveMasterTreatyBusiness_Life_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!SAVEMASTERTREATYBUSINESS_LIFE_SQL` / `RULE-CONNECT-SQL` | `INSERTBUSINESS_LIFE` |

`[terverifikasi]` Kelimanya mengembalikan `{OutputData.HASIL1 out}`. Sensus **27 Activity** modul
ini: **hanya `DeleteRowBusiness.xml`** yang menyebut `HASIL1`; **tidak satu pun activity `Save*`
membacanya**. **Galat procedure jatuh diam-diam.**

`[data DBA]` **Semantiknya:**

| Nilai `o_message` | Arti |
| --- | --- |
| **kosong / `NULL`** | **sukses** |
| **berisi teks** | **gagal** |

⚠️ `[data DBA]` Isinya **pesan bergaya UI Pega yang mengandung HTML** —
`<span style="color:red">…</span>` — ditambah `SQLERRM`.

## ADR terkait

**ADR-0007** (jejak audit kegagalan), **ADR-0015** (kegagalan ditangani eksplisit, tidak ditelan).

## Acceptance criteria

- [ ] ⚠️ **`o_message` DIPERIKSA** setelah **setiap** pemanggilan procedure di **kelima** jalur
      simpan: **kosong/NULL = sukses**, **berisi teks = gagal**. *(AC 46 spec;
      `[keputusan work owner]` — penyimpangan sadar 6)*
- [ ] ⚠️ Penyimpanan yang ditolak basis data **tidak pernah tampak berhasil** — API mengembalikan
      kegagalan, layar menampilkan galat. *(AC 46 spec)*
- [ ] ⚠️ **HTML dari pesan galat basis data tidak bocor** ke API maupun layar; pesan yang sampai ke
      pengguna **bersih dan terbaca**. Test yang menemukan tag markup pada respons API **gagal**.
      *(AC 47 spec; `[keputusan work owner]`)*
- [ ] **Uji konformansi lintas jalur**: test yang menemukan **jalur simpan mana pun** tanpa
      pemeriksaan hasil **gagal**. Kelima jalur diperiksa, bukan sebagian.
- [ ] Pembersihan HTML berada di **satu tempat**, bukan diulang lima kali.
- [ ] Pesan asli dari basis data **tetap tercatat** di log/jejak audit untuk penelusuran — yang
      dibersihkan hanya yang **sampai ke pengguna**.
- [ ] Kegagalan pada **penerapan massal** (tiket 08) diperiksa **per baris**, memakai jalur
      pemeriksaan yang sama. *(AC 31 spec)*
- [ ] Galat basis data yang **tidak dikenali** tetap dilaporkan sebagai kegagalan — tidak ada jalur
      yang menganggap galat tak dikenal sebagai sukses.

## Blocker

**Tidak ada.**

## Catatan

⚠️ **Tiket ini menutup celah, bukan membangun dari nol.** Tiket 01, 02, 05, 06, dan 07 masing-masing
sudah mengikat "`o_message` diperiksa, kegagalan ditampilkan" sebagai AC — karena jalur simpan yang
menelan galat bukanlah slice yang lengkap. Yang ditambahkan **di sini**: **pembersihan HTML yang
benar**, **pencatatan pesan asli**, dan **uji konformansi** yang menjamin tidak ada jalur tertinggal.

⚠️ **Nama menyesatkan** `[data DBA]`: parameter bernama `HASIL1` terbaca seperti kode hasil berangka;
ia sebenarnya **pesan galat teks**. Di lapisan domain, beri nama yang jujur.

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — galat procedure hanya dapat dipicu dengan
jujur pada basis data sungguhan.

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
| *"**`o_message` DIPERIKSA** setelah **setiap** pemanggilan procedure"* | procedure **tidak dipanggil** (keputusan o) → tidak ada `o_message` (K8). Yang ditegakkan: setiap galat Go — Oracle, validasi, keadaan data — sampai ke layar berkata-kata; galat tak terduga 500 dicatat di log server; nol galat ditelan. Uji konformansi: kelima jalur simpan menolak dengan badan `{"galat": …}` tanpa markup |
