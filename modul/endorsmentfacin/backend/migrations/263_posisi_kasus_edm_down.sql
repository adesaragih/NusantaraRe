-- Jalur mundur 263 - kasus Endorsment Fac In kembali ke bentuk lama (POSITION = workbasket, POSITION_NOTE kosong).
UPDATE {skema}.T_WORK_POLIS SET POSITION = POSITION_NOTE, POSITION_NOTE = NULL, TGL_UPDATE = SYSDATE
WHERE LINI = 'FACEDM' AND POSITION = 'Policy' AND POSITION_NOTE = 'ReasFacInMarketing'
/
