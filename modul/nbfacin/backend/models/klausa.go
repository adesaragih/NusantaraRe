package models

// ArgumenKlausa - satu argumen klausa (`.ArgumentList`; master: M_ARGCLAUSEFIRE JSON ArgumentNumber / Description /
// DefaultValue, RDBList SearchClauseArgFireSQL). Semua teks (tiket 47).
type ArgumenKlausa struct {
	Number, Description, Value string
}

// KlausaKasus - satu baris `pyWorkPage.ClauseList` (kelas Data-Clause; SearchClauseFireSQL_PostAct). ArgumentCount teks
// angka seperti di Pega / frontend.
type KlausaKasus struct {
	Code, Title, Description, Language, LanguageID, Content, ContentTemp, ArgumentCount string
	Arguments                                                                           []ArgumenKlausa
}

// HasilKlausa - satu hasil popup Choose Clause (RetrieveClauseSQL, kelas Int-CLAUSE): ID, Title / Text menurut bahasa,
// Info; Language = label bahasa yang dicari; ArgumentCount = RetrieveArgumentNumberClauseSQL.
type HasilKlausa struct {
	ID, Title, Info, Text, Language, ArgumentCount string
}
