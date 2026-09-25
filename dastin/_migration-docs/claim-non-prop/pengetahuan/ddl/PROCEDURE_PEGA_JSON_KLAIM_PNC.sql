-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE EDITIONABLE PROCEDURE "POOLDATA"."PEGA_JSON_KLAIM_PNC" 
(
    No_Klaim   IN VARCHAR2,
    PolisNO    IN VARCHAR2,
    PegaID     IN VARCHAR2,
    DataPega   IN CLOB,
    ErrMsg     OUT VARCHAR2,
    StsSimpan  OUT NUMBER
)
AS
BEGIN
    -- UPDATE dulu
    UPDATE POOLDATA.JSON_KLAIM
    SET 
        DATA_JSON    = DataPega
    WHERE IDPEGA = PegaID;

    -- jika tidak ada row yang diupdate, lakukan INSERT
    IF SQL%ROWCOUNT = 0 THEN
        INSERT INTO POOLDATA.JSON_KLAIM
        (
            MNK_NO_KLAIM,
            DATA_JSON,
            IDPEGA,
            TGL_INPUT,
            TGL_KONVERSI,
            NOPOLIS,
            IDPROD,
            STS_KONVERSI
        )
        VALUES
        (
            No_Klaim,
            DataPega,
            PegaID,
            SYSDATE,
            NULL,
            PolisNO,
            1,
            NULL
        );
    END IF;

    -- sukses
    ErrMsg := 'Data sudah disimpan';
    StsSimpan := 1;

EXCEPTION
    WHEN OTHERS THEN
        -- rollback jika error
        ROLLBACK;
        ErrMsg := 'PEGA_JSON_KLAIM Error : ' || SQLERRM;
        StsSimpan := 0;
END;