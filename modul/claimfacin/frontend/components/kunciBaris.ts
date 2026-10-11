// Disalin dari `modul/claimnonprop/frontend/components/kunciBaris.ts` (pola, bukan impor; asal Claim Prop): keputusan
// work owner yang disebut di bawah adalah keputusan layar Claim Prop yang ditiru Claim Fac In (prompt Claim Fac In tahap 1).
// Key baris popup: HARUS unik per baris. Key kembar membuat React membiarkan baris lama tertinggal saat daftar berganti
// (temuan work owner Claim Prop 08-10-2026: filter menampilkan baris lain). Grid hasil Choose Polis (`GetPolisForClaim_SQL`)
// memakai SELECT DISTINCT atas seluruh kolomnya - nomor polis yang sama dapat muncul berkali-kali (Prodke / endorsemen
// berbeda), jadi gabungan seluruh kolom unik.

import type { BarisPolisCari } from '../api'

/** Pemisah kolom key: karakter kendali, tidak pernah ada di data polis. */
const PEMISAH = String.fromCharCode(1)

export function kunciPolis(b: BarisPolisCari): string {
  return [
    b.policyNo,
    b.customerName,
    b.sourceOfBusinessName,
    b.cedingCoName,
    b.qq,
    b.startDateTime,
    b.endDateTime,
    b.prodke,
    b.businessName,
  ].join(PEMISAH)
}
