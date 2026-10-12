-- Jalur mundur 029 - kembali ke VARCHAR2(8).
--
-- ⚠️ DAPAT GAGAL, dan itu disengaja (pola 018_down): bila sudah ada nilai lebih dari 8 byte (mis. "Non-Catastrophe",
-- "Catastrophe"), Oracle menolak dengan ORA-01441 "cannot decrease column length because some value is too big" - TANPA
-- memotong data. Kosongkan / ubah nilainya dulu bila jalur mundur memang diperlukan.
ALTER TABLE {skema}.T_GENERAL_CLAIM MODIFY (
  STS_KATASTROFE VARCHAR2(8),
  IS_KPR         VARCHAR2(8)
)
/
