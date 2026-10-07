package models

// Untuk apa berkas ini: penanda jalur NonProporsional di layar EDM. Asal: salinan sebagian
// `modul/nbtreatyin/backend/models/nonprop_layar.go` - bagian NB lain (subsection `DetailPoliciesNonProportional`,
// spreading NonProp yang dapat disunting, pajak NonProp preACT 16) tidak ada di section EDM.
//
// Bukti XML EDM: `Section/DetailPolicyTreatyInAddendum` S11 (`.IsNewPolicyNonProp = 1`) menyertakan
// `DetailPolicyTreatyInAddPremi` (grid XOL Previous / Current / Total Difference - semuanya hanya-baca lewat panel
// `DetailPolicyAddPremiDetail`) dan Installment; S17 (`.IsNewPolicyNonProp = 0 || ''`) menyertakan tab Old / New /
// Value Difference.

// PolisNonPropBaru = `.IsNewPolicyNonProp = 1`.
func PolisNonPropBaru(h *Halaman) bool { return samaDenganSatu(h.Ambil(pt + "IsNewPolicyNonProp")) }
