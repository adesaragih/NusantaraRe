-- 841 - kunci utama T_M_ACCOUNT (keputusan work owner 04-10-2026: "tambahkan pk nya"). Katalog DEV 04-10-2026:
-- 18.658 baris, 18.658 ID berbeda, nol ID kosong, nol ID yang hanya beda huruf/spasi - PK dapat dipasang tanpa
-- membersihkan data. Indeks unik ikut dibuat Oracle.
ALTER TABLE {skema}.T_M_ACCOUNT ADD CONSTRAINT PK_T_M_ACCOUNT PRIMARY KEY (ID)
/
