-- Jalur mundur untuk 431_sequences_tab_treatyin.sql
--
-- ⚠️ Dibongkar SEBELUM 430 (nomor turun), jadi tabelnya masih ada ketika
-- sequence-nya hilang. Itu keadaan yang tidak dapat dipakai menulis - dan
-- memang begitu maksudnya: membongkar mundur satu langkah lalu berhenti
-- harus terasa, bukan diam.
DROP SEQUENCE {skema}.SEQ_MTI_COMMENT
/
DROP SEQUENCE {skema}.SEQ_MTI_INSTALLMENTITEM
/
DROP SEQUENCE {skema}.SEQ_MTI_INSTALLMENT
/
DROP SEQUENCE {skema}.SEQ_MTI_RETENTION
/
DROP SEQUENCE {skema}.SEQ_MTI_EGNPI
/
DROP SEQUENCE {skema}.SEQ_MTI_ACCUMULATION
/
DROP SEQUENCE {skema}.SEQ_MTI_PORTFOLIO
/
DROP SEQUENCE {skema}.SEQ_MTI_REPORTINGPERIOD
/
