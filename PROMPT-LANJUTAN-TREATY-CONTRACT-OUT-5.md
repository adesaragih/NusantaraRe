# PROMPT — LANJUTAN 5 MODUL **Treaty Contract Out** *(folder `OUTPUT_HASIL_RNM`, cabang `main` @ `626bb9c` atau lebih baru)*: **OQ-TCO-24 ditutup dari `ALL_SOURCE`**

> Hanya konteks Treaty Contract Out. Commit dengan jalur eksplisit. **Jangan** menyentuh berkas frontend yang sedang disunting sesi lain di
> `main` *(`index.html`, `App.tsx`, `PagarGalat.tsx`, `labels.ts`, `styles.css`, `Shell.tsx`, `dasar.tsx`, `hooks/useTema.ts`, `lib/tema*.ts`)*.

## 0. VERIFIKASI LANJUTAN 4 *(asisten, 29-09-2026)*

`708a351`, `8426880`, `80d1985`, `626bb9c` ada; Go di `main` **914 PASS · 0 FAIL**, dengan tag `db` **53 SKIP**; vet bersih. Data DEV:
`TREATYYEAR` berbentuk lain **0** baris *(`NOT REGEXP_LIKE(STARTDATE/ENDDATE,'^[0-9]{8}$')`)* — penolakan bentuk tidak mengganggu data yang ada.

## 1. OQ-TCO-24 — badan `PEGA_M_ATTACHMENT` dibaca asisten *(`ALL_SOURCE`, 35 baris, baris komentar dibuang)*

| Baris | Isi |
| --- | --- |
| 1 | `PEGA_M_ATTACHMENT(IDPega IN VARCHAR2, DataPega IN CLOB, ErrMsg OUT VARCHAR2, StsSimpan OUT NUMBER)` |
| 7 | `SELECT count(id) INTO id_count FROM pooldata.M_ATTACHMENTTREATY WHERE ID = IDPega` — **`id_count` variabel lokal, bukan tabel** |
| 14–16 | bila ada: `UPDATE pooldata.M_ATTACHMENTTREATY SET DATA_JSON = DataPega WHERE ID = IDPega; COMMIT` |
| 24–26 | bila tidak: `INSERT INTO pooldata.M_ATTACHMENTTREATY (ID, DATA_JSON, DATEINPUT) VALUES (IDPega, DataPega, TO_DATE(to_char(sysdate,'dd/MM/yyyy'),'dd/MM/yyyy'))` *(tanpa `COMMIT` di cabang ini)* |

Katalog: `M_ATTACHMENTTREATY` **dan** `M_ATTACHMENTTREATY_2` keduanya ada di `POOLDATA`.

**Kesimpulan:** `M_ATTACHMENTTREATY` adalah **simpanan JSON halaman** *(`DATA_JSON` CLOB per `ID`)*, satu-satunya penulisnya prosedur ini
lewat `TreatyOutSaveAttachment` b1674. `M_ATTACHMENTTREATY_2` adalah tabel **relasional** yang dibaca dan dihapus seluruh RDB modul.

**Keputusan `[asisten dari bukti; veto work owner]`:** aplikasi **tetap menulis dan membaca `M_ATTACHMENTTREATY_2`**; **tidak** menulis
simpanan JSON `M_ATTACHMENTTREATY` — konsisten dengan keputusan proyek bahwa simpanan JSON halaman Pega tidak diteruskan *(spec modul:
"tabel JSON dibuang"; PremiumList: revisi penyimpanan JSON dibuang)*. Pembaca hilir yang membaca JSON lama tidak ada di korpus modul ini.

## 2. KERJAKAN — satu commit

- `dba-procedures.md`: badan prosedur §1 *(struktur, tanpa baris komentar)*; OQ-TCO-24 **ditutup** di register OQ dengan kesimpulan dan
  keputusan §1; ralat bertanggal di tiket 12 dan STRUKTUR *(`ID_COUNT` bukan tabel)*.
- Kode: komentar di repository lampiran merujuk OQ-TCO-24; **perilaku tidak berubah**.

Commit `docs: treaty-contract-out — OQ-TCO-24 ditutup, M_ATTACHMENTTREATY simpanan JSON tidak diteruskan`.

## 3. MASIH TERBUKA

- **OQ-TCO-22** *(folder penyimpanan, durasi, nama berkas)* — work owner.
- **Paket 3 lanjutan 4** *(halaman yatim `tco-kontrak`/`tco-klausul` di `App.tsx`)* — menunggu `App.tsx` bebas dari suntingan sesi lain.

---

*Disusun 29 September 2026 dari `ALL_SOURCE` `PEGA_M_ATTACHMENT` (DEV, baca-saja), katalog tabel lampiran, dan verifikasi `626bb9c`.*
