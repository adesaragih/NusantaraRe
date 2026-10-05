-- 883 mundur - kota yang DITAMBAH 883 (penanda CREATE_OP = 'MIGRASI-883') dibuang, lalu CITYINPUT.PROVINCENAME
-- (isinya ikut hilang; isi awal dapat dibuat ulang dari RW oleh 883). ⚠️ Kota itu ikut terbuang walau sesudahnya
-- diubah lewat menu (pengubah tercatat di UPDATE_OP, bukan CREATE_OP), dan DISTRICTINPUT.CITYID yang merujuknya tertinggal.
DELETE FROM {skema}.CITYINPUT WHERE CREATE_OP = 'MIGRASI-883'
/
ALTER TABLE {skema}.CITYINPUT DROP (PROVINCENAME)
/
