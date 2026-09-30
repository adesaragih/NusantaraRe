# Sumber penomoran klaim — `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` dan tabel pendukungnya

`[data DBA — dibaca sendiri 26 September 2026 dari katalog instance pengembangan, belum dikonfirmasi DBA]`

Dibaca lewat `ALL_SOURCE` dan `ALL_TAB_COLUMNS` pada instance **pengembangan** (Oracle 12.2.0.1),
skema `POOLDATA`. Yang **tidak** dibaca: nilai `NO_SEQ` (isi counter) dan nilai `TANGGAL_CLOSING`.
Dokumen ini memenuhi brief ronde 4 §5 butir 4–6. Ia **tidak** memutuskan apa pun tentang §2 o1–o3;
ia hanya membuat keputusan itu dapat diambil tanpa menebak.

## Sumber procedure — 93 baris, disalin apa adanya

Dua karakter komentar di baris 65 dan 78 tiba sebagai karakter pengganti (masalah pengkodean saat
dibaca), ditulis di sini sebagai `->`. Sisanya verbatim.

```sql
PROCEDURE	   PROC_GENERATE_SEQUENCE_NUMBER (
    p_class	  IN  VARCHAR2,
    p_jenis	  IN  VARCHAR2,
    p_proddate	  IN  DATE,
    p_bulan	  OUT VARCHAR2,
    p_seq_number  OUT VARCHAR2
)
AS
    v_now	    DATE;
    v_periode_date  DATE;
    v_tahun	    NUMBER(4);
    v_mm_yyyy	    VARCHAR2(7);
    v_day_closing   NUMBER;
    v_seq	    NUMBER;
BEGIN
    ------------------------------------------------------------------
    -- Current datetime Jakarta
    ------------------------------------------------------------------
    IF p_proddate IS NULL THEN
	v_now := CAST(SYSTIMESTAMP AT TIME ZONE 'Asia/Jakarta' AS DATE);
    ELSE
	v_now := p_proddate;
    END IF;

    ------------------------------------------------------------------
    -- Ambil tanggal closing
    ------------------------------------------------------------------
    SELECT TO_NUMBER(tanggal)
    INTO v_day_closing
    FROM POOLDATA.TANGGAL_CLOSING
    WHERE ROWNUM = 1;

    ------------------------------------------------------------------
    -- Tentukan periode
    ------------------------------------------------------------------
    IF TRUNC(v_now) <= TO_DATE('02/01/2026','DD/MM/YYYY') THEN
	v_mm_yyyy := '12.2025';
	v_tahun   := 2025;
    ELSE
	v_periode_date := ADD_MONTHS(
	    v_now,
	    CASE
		WHEN TO_NUMBER(TO_CHAR(v_now,'DD')) > v_day_closing
		THEN 1
		ELSE 0
	    END
	);

	v_mm_yyyy := TO_CHAR(v_periode_date,'MM.YYYY');
	v_tahun   := EXTRACT(YEAR FROM v_periode_date);
    END IF;

    ------------------------------------------------------------------
    -- LOCK row (kalau ada)
    ------------------------------------------------------------------
    BEGIN
	SELECT no_seq
	INTO v_seq
	FROM POOLDATA.GENERATE_SEQUENCE_NUMBER
	WHERE class = p_class
	  AND jenis = p_jenis
	  AND tahun = v_tahun
	FOR UPDATE;

	-- Jika ada -> increment
	v_seq := v_seq + 1;

	UPDATE POOLDATA.GENERATE_SEQUENCE_NUMBER
	SET no_seq  = v_seq,
	    mm_yyyy = v_mm_yyyy,
	    tanggal = v_now
	WHERE class = p_class
	  AND jenis = p_jenis
	  AND tahun = v_tahun;

    EXCEPTION
	WHEN NO_DATA_FOUND THEN
	    -- Jika belum ada -> insert pertama
	    v_seq := 1;

	    INSERT INTO POOLDATA.GENERATE_SEQUENCE_NUMBER
		(class, jenis, tahun, no_seq, tanggal, mm_yyyy)
	    VALUES
		(p_class, p_jenis, v_tahun, v_seq, v_now, v_mm_yyyy);
    END;

    ------------------------------------------------------------------
    -- Output
    ------------------------------------------------------------------
    p_seq_number := LPAD(v_seq, 5, '0');
    p_bulan	 := v_mm_yyyy;

END;
```

## Yang terbaca dari sumbernya — fakta, bukan keputusan

1. **Tidak ada `COMMIT` di dalam procedure.** `COMMIT` yang terlihat di rule Pega
   `GetSequenceNumber_SQL` berada di blok pemanggil, bukan di procedure. Ini bukti untuk **OQ-013**
   (batas transaksi di sisi database); penutupannya milik work owner / DBA.
2. Kunci baris: `SELECT … FOR UPDATE` atas `(CLASS, JENIS, TAHUN)`; bila belum ada, `INSERT`
   pertama dengan `NO_SEQ = 1`. Dua pemanggil serentak diserialkan oleh kunci baris itu.
3. Periode: hari tutup buku `TANGGAL_CLOSING.TANGGAL` (teks `VARCHAR2(10)`, dibaca `TO_NUMBER`,
   `ROWNUM = 1`); tanggal di atas hari itu digeser satu bulan. Cutover: `TRUNC(now) <= 02/01/2026`
   memaksa periode `12.2025` dan tahun 2025.
4. Waktu: `SYSTIMESTAMP AT TIME ZONE 'Asia/Jakarta'`, kecuali `p_proddate` diberikan.
5. Keluaran: `p_seq_number = LPAD(no_seq, 5, '0')` dan `p_bulan = MM.YYYY`. **Format nomor lengkap
   (`<prefix>K<kode bisnis>.MM.YYYY.<5 digit>`) dirakit pemanggil**, bukan procedure — jadi separuh
   logika penomoran ada di rule Pega, bukan di sini.
6. `TAHUN` disimpan sebagai `VARCHAR2(5)` tetapi dibandingkan dengan `v_tahun NUMBER(4)` —
   konversi implisit Oracle. Kalau disalin ke Go, bandingkan sebagai teks yang dibentuk sama.

## DDL pendukung — katalog

`GENERATE_SEQUENCE_NUMBER` (25 baris pada tanggal pembacaan; nilai `NO_SEQ` tidak dibaca):

| # | Kolom | Tipe | Null |
| ---: | --- | --- | :-: |
| 1 | `CLASS` | `VARCHAR2(200)` | N |
| 2 | `JENIS` | `VARCHAR2(50)` | N |
| 3 | `TAHUN` | `VARCHAR2(5)` | N |
| 4 | `NO_SEQ` | `NUMBER` | Y |
| 5 | `TANGGAL` | `DATE` | Y |
| 6 | `MM_YYYY` | `VARCHAR2(10)` | Y |

Pasangan `(CLASS, JENIS)` yang menyangkut Life dan Claim, tanpa nilai counter:

| `CLASS` | `JENIS` |
| --- | --- |
| `ASM-FW-GCNMFW-Work-ClaimLife` | `RNML-K` |
| `ASM-FW-GCNMFW-Work-KomiteLife` | `RNML-A` |
| `ASM-FW-GISFW-Work-LIFE` | `RNML-QR/QP/TP/TR` |
| `ASM-FW-GCNMFW-Work-ClaimTreaty` | `RNM-K`, `RNM-A`, `RNM-P` |
| `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp` | `RNM-K`, `RNM-A` |

`KODE_PRODUKSI` (2 baris): `KODE VARCHAR2(5)`, `TYPE VARCHAR2(10)`; baris `LIFE` bernilai
`RNML-`. `TANGGAL_CLOSING` (1 baris): satu kolom `TANGGAL VARCHAR2(10)` — hari tutup buku
disimpan sebagai **teks**, nilainya tidak dibaca.

## Yang dokumen ini TIDAK putuskan

⛔ Apakah logika di atas disalin ke Go (keputusan o, brief ronde 4 §2 o1–o3), siapa yang memegang
`GENERATE_SEQUENCE_NUMBER` selama koeksistensi, dan bagaimana AC 2/3/10 tiket 02 ditulis ulang.
Semuanya milik work owner dan DBA.
