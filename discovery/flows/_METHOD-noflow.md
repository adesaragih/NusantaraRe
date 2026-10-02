# Konvensi Telusur D2 — Modul **TANPA** rule `Flow`

Pelengkap **`_METHOD.md`**, yang tetap berlaku penuh. Dokumen ini hanya menambahkan **cara
menemukan titik masuk dan merekonstruksi mesin status** ketika tidak ada graf `Flow`.
Ditetapkan 2026-09-13 untuk STEP D2 Tahap 5 (5 modul, OQ-005).

Seluruh aturan `_METHOD.md` tetap mengikat: batas fase §0, empat aturan anti-halusinasi §1,
format 7 bagian §3, label wajib §4, dan larangan memakai **nama** sebagai bukti perilaku (§2.2).

---

## 1. Mengapa metode terpisah diperlukan

`_METHOD.md` §2.1 memulai telusur dari `Flow` → `<pyStartActivity>`. Lima modul tidak punya rule
`Flow` sama sekali (OQ-005) `[terverifikasi]`:

| Modul | File `.xml` | `Harness` | `FlowAction` | `Section` | `Activity` |
| --- | ---: | ---: | ---: | ---: | ---: |
| Treaty In | 329 | 3 | 32 | 50 | 149 |
| Treaty In Adjustment | 379 | 6 | 43 | 68 | 155 |
| Treaty Contract Out | 303 | 3 | 2 | 49 | 168 |
| Master Product Name Life | 114 | 1 | 11 | 12 | 37 |
| Master Contract Retro Life | 66 | 4 | 2 | 8 | 27 |

Perintah audit:
```
for m in "Treaty In" "Treaty In Adjustment" "Treaty Contract Out" \
         "Master Product Name Life" "Master Contract Retro Life"; do
  echo "$m $(find "$m" -name '*.xml' | wc -l)"
  for d in "$m"/*/; do echo "   $(basename "$d") $(ls "$d" | wc -l)"; done
done
```

Tanpa `Flow`, **tidak ada `<pyFrom>`/`<pyTo>`, tidak ada `<pyShapeType>`, tidak ada
`<pyConditionType>`** — tiga tag yang menjadi tulang punggung §2.2 `_METHOD.md`. Graf harus
direkonstruksi dari sumber lain, dan **caranya harus ditetapkan lebih dulu, bukan diimprovisasi**.

---

## 2. Titik masuk berbukti

### 2.1 Urutan pencarian yang mengikat

1. **`Harness/`** — layar/portal. Identitas dibaca dari `<pxInsName>`. Class `DATA-PORTAL`
   menandai layar portal; class aplikasi (mis. `ASM-FW-GISFW-INT-PRODUCT_LIFE`) menandai layar
   yang menempel pada satu entitas data.
2. **`Section/` yang di-`<pyInclude>` oleh Harness** — struktur layar. `_METHOD.md` §2.4 tetap
   berlaku: `<pySection>` di aksi *Refresh* **bukan** dependency.
3. **Section panel tombol** — dalam praktiknya inilah titik masuk perilaku yang sesungguhnya
   (§3.1).
4. **`FlowAction/`** — aksi yang dapat dipanggil; sebagian dipakai sebagai *modal/sub-layar*,
   bukan sebagai langkah proses.
5. Dari situ turun ke `Activity` / `DataTransform` / `When` / `RDBList` **persis seperti
   `_METHOD.md` §2.3**.

### 2.2 Harness besar: grep, jangan baca

Harness terbesar korpus ada di kelompok ini: `Treaty In/Harness/InputTreatyInOffer.xml`
**8.225.095 byte** `[terverifikasi]` (`ls -l`). **Dilarang membacanya utuh** (`_METHOD.md` §2.5).
Yang diambil dengan grep:

```
grep -o "<pxInsName>[^<]*"  "<harness>.xml" | head -1      # identitas
grep -oE "<pySection>[^<]*" "<harness>.xml" | sed 's/<[^>]*>//' | sort | uniq -c | sort -rn
```

### 2.3 Tabel titik masuk yang dipakai di Tahap 5 `[terverifikasi]`

| Modul | Harness titik masuk | `pxInsName` |
| --- | --- | --- |
| Treaty In | `Harness/InputTreatyInOffer.xml` | `DATA-PORTAL!INPUTTREATYINOFFER` |
| Treaty In Adjustment | `Harness/InputTreatyInAdjustment.xml` | `DATA-PORTAL!INPUTTREATYINADJUSTMENT` |
| Treaty Contract Out | `Harness/InboxTreatyContract.xml` | `DATA-PORTAL!INBOXTREATYCONTRACT` |
| Master Product Name Life | `Harness/InwardProductName.xml` | `ASM-FW-GISFW-INT-PRODUCT_LIFE!INWARDPRODUCTNAME` |
| Master Contract Retro Life | `Harness/InboxRetroLifeReinsurersList.xml` | `DATA-PORTAL!INBOXRETROLIFEREINSURERSLIST` |

---

## 3. Merekonstruksi mesin status tanpa graf

### 3.1 Tombol = transisi. Panel tombol = daftar transisi yang mungkin.

Pola yang terbukti di Tahap 5: satu `Section` memuat seluruh tombol aksi, dan **tiap tombol
mengikat satu rule** lewat `<pyName>`. Bukti: `Treaty In/Section/TreatyInActionButtons.xml`
(`DATA-PORTAL!TREATYINACTIONBUTTONS`, 369.789 byte).

```
grep -oE "<pyLabel>[^<]{1,45}" "<section>.xml" | sed 's/<pyLabel>//' | sort -u    # label tombol
grep -oE "<pyName>[^<]{1,45}"  "<section>.xml" | sed 's/<pyName>//'  | sort -u    # rule terikat
```

Aturan pembacaan:
- **Label tombol adalah label** (`_METHOD.md` §2.2) — yang membuktikan perilaku adalah rule di
  `<pyName>`.
- Tombol yang tidak punya rule terikat **tidak dicatat sebagai transisi**.

### 3.2 Kondisi visibilitas tombol = guard transisi

Ini pengganti `<pyConditionType>` pada connector. Ekspresi visibilitas menyebut properti status
dan keanggotaan workbasket secara literal:

```
grep -ohE "<py[A-Za-z]*>[^<]{0,110}(\.Position|StatusAkseptasi|RevisionState|ViewState)[^<]{0,70}" \
  "<section>.xml" | sed 's/&amp;#61;/=/g;s/&amp;quot;/"/g;s/&amp;amp;/\&/g' | sort -u
```

**Kewajiban:** guard dicatat **apa adanya**, termasuk yang tampak selalu salah
(mis. `FALSE && …`, `… && 1=2`). Tombol ber-guard demikian dicatat sebagai **tidak dapat muncul
menurut ekspresi yang terbaca** — bukan sebagai "dihapus" atau "mati", karena Pega dapat
menimpanya dari tempat lain yang tidak terlihat. Ini paralel dengan OQ-023 (shape tanpa connector).

### 3.3 Transisi status yang sesungguhnya: **Data Transform berjenjang**

Temuan metodologis terpenting Tahap 5 `[terverifikasi]`: pada modul tanpa `Flow`, tangga
persetujuan tidak hilang — ia pindah ke **satu `RULE-OBJ-MODEL` (Data Transform) bersarang**
dengan aksi `WHEN` / `OTHERWISE_WHEN` / `SET`.

Bukti: `Treaty In/DataTransform/Akseptasi_DT.xml` (`DATA-PORTAL!AKSEPTASI_DT`) —
10 `WHEN`, 18 `OTHERWISE_WHEN`, 64 `SET`.

```
grep -oE "<pyActionName>[^<]*" "<dt>.xml" | sed 's/<[^>]*>//' | sort | uniq -c
```

**Cara membacanya (mengikat).** `rowdata` di ekspor Pega **bersarang**, dan tag penutup anak
muncul sebelum tag penutup induk. Pembacaan dengan pencocokan baris datar **akan salah**. Wajib
memakai penghitung kedalaman:

```awk
function val(l,t,  s,e,v){s=index(l,"<" t ">"); if(s==0)return ""; v=substr(l,s+length(t)+2);
                          e=index(v,"<"); if(e>0)v=substr(v,1,e-1); return v}
{
  if ($0 ~ /<rowdata/) { depth++; A[depth]=""; N[depth]=""; V[depth]="" }
  if (depth>=1) {
    if (index($0,"<pyActionName>"))                      A[depth]=val($0,"pyActionName")
    if (index($0,"<pyPropertiesName>")  && N[depth]=="") N[depth]=val($0,"pyPropertiesName")
    if (index($0,"<pyPropertiesValue>") && V[depth]=="") V[depth]=val($0,"pyPropertiesValue")
  }
  if ($0 ~ /<\/rowdata>/) {
    a=A[depth]; n=N[depth]; v=V[depth]; ind=sprintf("%*s",(depth-1)*2,"")
    if (a=="WHEN"||a=="OTHERWISE_WHEN") printf "%s%s: %s\n", ind, a, n
    else if (a=="SET")                  printf "%s  SET %s := %s\n", ind, n, v
    depth--
  }
}
```

**Urutan keluaran:** blok `SET` muncul **sebelum** baris `WHEN` miliknya (karena anak ditutup
lebih dulu). Jadi setiap baris `WHEN:` adalah **label bagi blok `SET` tepat di atasnya**.
Kekeliruan membaca arah ini akan membalik seluruh mesin status.

### 3.4 Bentuk state machine yang dihasilkan

Untuk tiap cabang, catat **tripel** yang di-`SET`, bukan hanya statusnya:

```
(kondisi) -> { properti status , properti posisi/antrean , properti pemilik berikutnya }
```

Contoh bentuk (nilai diisi dari korpus, bukan dari dugaan):

```
WHEN <guard tingkat-1>
   WHEN <guard posisi saat ini>
        WHEN <aksi yang dipilih pengguna>
             SET <properti status>   := "<nilai literal>"
             SET <properti posisi>   := "<nama workbasket literal>"
             SET <properti pemilik>  := <sumber nilai>
```

**Dilarang** menuliskan "disetujui", "ditolak", atau "naik ke atasan" sebagai perilaku bila korpus
hanya menunjukkan perpindahan nilai properti. Arti nilai status ditulis
`arti belum terverifikasi` (OQ-020), sama seperti `ProposalAcceptStatus` di Tahap 4 — di mana
pemetaan nilai → hasil ternyata **tidak ada di korpus** (OQ-043).

### 3.5 Yang tetap tidak dapat direkonstruksi

Tanpa `Flow`, hal-hal berikut **tidak punya sumber bukti sama sekali** dan wajib dicatat sebagai
batas pengetahuan, bukan diisi dugaan:

| Tidak terbaca | Alasan |
| --- | --- |
| Urutan wajib antar layar | tidak ada connector; urutan hanya tersirat dari guard |
| Assignment ke workbasket sebagai mekanisme Pega | tidak ada `<pyImplementation>`/`<pyRouteTo>`; yang ada hanya properti data bernama posisi |
| SLA / timer / eskalasi | tidak ada shape yang membawanya |
| Titik akhir proses | tidak ada shape `End` |
| Siapa yang boleh menekan tombol | hanya terbaca sejauh ekspresi visibilitas menyebutnya |

---

## 4. Ringkas status konflik & open-question untuk 5 modul ini

`[terverifikasi]` Entri register `../inventory/_oq011-konflik-isi.md` yang menyentuh tiap modul —
`Treaty In` **98**, `Treaty In Adjustment` **56**, `Master Product Name Life` **22**,
`Treaty Contract Out` **15**, `Master Contract Retro Life` **5**
(`grep -cF "<modul>/" ../inventory/_oq011-konflik-isi.md`).
`Treaty In` ↔ `Treaty In Adjustment` berbagi **323 file bernama sama**; setelah normalisasi 18 tag
**277 identik (85,8 %)**, **46 berbeda**, 6 hanya di `Treaty In`, 56 hanya di `Adjustment` — jadi
telusur difokuskan pada 46 + 56 itu (OQ-010, OQ-011).
OQ yang mengikat di sini: **OQ-005** (titik masuk — dijawab dokumen ini), **OQ-010** (dua modul
atau satu), **OQ-022** (`Treaty Contract Out` bukan outward), **OQ-011** (per-varian),
**OQ-002/OQ-012** (SP & kolom JSON), **OQ-021/OQ-027** (guard identitas), **OQ-024/OQ-007**
(workbasket & otorisasi). `IsPEGAPROD` **tidak ada di kelima modul** (`find … -iname IsPEGAPROD.xml`
→ 0), sehingga OQ-029 tidak berlaku di Tahap 5.
