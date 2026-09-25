-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE EDITIONABLE PROCEDURE "POOLDATA"."PEGA_JSON_OS_AKSEP_KLAIMTNP" (Currency IN VARCHAR2, LayerT IN VARCHAR2, NoClaim IN VARCHAR2,IDMaster IN VARCHAR2, PolisNO IN VARCHAR2, PegaID IN VARCHAR2, DataPega IN CLOB, STATUS_RJT IN VARCHAR2,KONVERSI IN VARCHAR2, ErrMsg OUT VARCHAR2, StsSimpan OUT NUMBER)
AS
id_count INTEGER;

BEGIN
    BEGIN
        SELECT count(1) INTO id_count FROM OS_AKSEPTASI_KLAIM a WHERE CASEID = PegaID AND a.data_json.TypeLoss= LayerT and a.data_json.Currency = Currency ;
        EXCEPTION
        WHEN OTHERS THEN
                ErrMsg := 'SELECT OS_AKSEPTASI_KLAIM COUNT Error : ' || sqlerrm;
                ROLLBACK;
                RETURN;
    END;
    
    BEGIN
        INSERT INTO OS_AKSEPTASI_KLAIM (CASEID,NOCLAIM,MASTERID, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, STS_KONVERSI) VALUES (PegaID, NoClaim,IDMaster, DataPega,TO_DATE(to_char(sysdate,'dd/MM/yyyy'), 'dd/MM/yyyy'), PolisNO, STATUS_RJT, KONVERSI);
        EXCEPTION
            WHEN OTHERS THEN
                ErrMsg := 'OS_AKSEPTASI_KLAIM Error : ' || sqlerrm;
                StsSimpan := 0;
                ROLLBACK;
                RETURN;
        END;
    
--    IF id_count = 0 THEN
--        BEGIN
--            INSERT INTO OS_AKSEPTASI_KLAIM (CASEID,NOCLAIM, DATA_JSON, TANGGAL, NOPOLIS, STS_REJECT, STS_KONVERSI) VALUES (PegaID, NoClaim, DataPega,TO_DATE(to_char(sysdate,'dd/MM/yyyy'), 'dd/MM/yyyy'), PolisNO, STATUS_RJT, KONVERSI);
--        EXCEPTION
--            WHEN OTHERS THEN
--                ErrMsg := 'OS_AKSEPTASI_KLAIM Error : ' || sqlerrm;
--                StsSimpan := 0;
--                ROLLBACK;
--                RETURN;
--        END;
--    ELSE 
--        BEGIN
--            UPDATE OS_AKSEPTASI_KLAIM a SET DATA_JSON = DataPega WHERE CASEID = PegaID and a.data_json.TypeLoss= LayerT and a.data_json.Currency = Currency ;
--        EXCEPTION
--            WHEN OTHERS THEN
--                ErrMsg := 'OS_AKSEPTASI_KLAIM Error : ' || sqlerrm;
--                StsSimpan := 0;
--                ROLLBACK;
--                RETURN;
--        END; 
--    END IF;
    ErrMsg := 'Data sudah di simpan';
    StsSimpan := 1;

EXCEPTION
    WHEN OTHERS THEN
        ErrMsg := 'SELECT OS_AKSEPTASI_KLAIM Error : ' || sqlerrm;
        StsSimpan := 0;
        ROLLBACK;
        RETURN;
END;