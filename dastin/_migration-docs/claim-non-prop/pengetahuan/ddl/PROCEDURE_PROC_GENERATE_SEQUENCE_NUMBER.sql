-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE EDITIONABLE PROCEDURE "POOLDATA"."PROC_GENERATE_SEQUENCE_NUMBER" (
    p_class       IN  VARCHAR2,
    p_jenis       IN  VARCHAR2,
    p_proddate    IN  DATE,
    p_bulan       OUT VARCHAR2,
    p_seq_number  OUT VARCHAR2
)
AS
    v_now           DATE;
    v_periode_date  DATE;
    v_tahun         NUMBER(4);
    v_mm_yyyy       VARCHAR2(7);
    v_day_closing   NUMBER; 
    v_seq           NUMBER;
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

        -- Jika ada → increment
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
            -- Jika belum ada → insert pertama
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
    p_bulan      := v_mm_yyyy;

END;