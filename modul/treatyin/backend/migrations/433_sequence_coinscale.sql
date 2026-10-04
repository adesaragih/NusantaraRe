-- Sequence identitas untuk M_TREATYIN_COINSCALE
--
-- Berdiri TERPISAH dari tabelnya, mengikuti 402, 415, 419, 421, dan 431.
--
-- INV-02: pengenal datang dari SEQUENCE, tidak pernah dari `MASTERID`
-- ditambah urutan. INV-03: `NOCYCLE` — satu nomor tidak pernah terpakai dua
-- kali. Awalan `SEQ_MTI_` konsisten dengan migrasi 431; `SEQ_MTI_COINSCALE`
-- 17 bita, jauh di bawah batas 30.
CREATE SEQUENCE {skema}.SEQ_MTI_COINSCALE START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
