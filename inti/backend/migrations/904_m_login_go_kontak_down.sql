-- Jalur mundur 904 - keempat kolom kontak dibuang beserta isinya. Jalur mundur ini untuk skema uji.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029).
ALTER TABLE {skema}.M_LOGIN_GO DROP (EMAIL, PHONE_NUMBER, EMPLOYEE_ID, JOB_POSITION)
/
