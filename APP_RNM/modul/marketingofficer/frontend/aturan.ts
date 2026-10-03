// Aturan layar Marketing Officer - fungsi murni, diuji tanpa DOM. Aturan sebenarnya dijaga backend
// (`backend/services`); yang di sini hanya mencegah kiriman yang pasti ditolak dan menyusun tampilan.

import type { BarisMO, Cabang, Isian, MarketingOfficer, OpsiLeader, Pilihan } from './api'
import { MO } from './labels'

/** `CLIENTID2` baris leader. */
export const NILAI_LEADER = 'LEADER'

export function isianKosong(): Isian {
  return { aksesLogin: '', leader: false, leaderId: '', branchParent: '', branchDetailId: '', aktif: true }
}

/** Isian form ubah dari baris tersimpan. */
export function isianDari(b: MarketingOfficer): Isian {
  const leader = b.clientId2 === NILAI_LEADER
  return {
    aksesLogin: b.aksesLogin,
    leader,
    leaderId: leader ? '' : b.clientId2,
    branchParent: b.branchParent,
    branchDetailId: b.branchDetailId,
    aktif: b.moStatus === '1',
  }
}

/** Galat pertama yang pasti ditolak backend; `null` = boleh dikirim. Baris lama tanpa akun boleh tetap tanpa akun. */
export function periksa(isi: Isian, baru: boolean): string | null {
  if (baru && isi.aksesLogin.trim() === '') return MO.galatAkun
  if (!isi.leader && isi.leaderId.trim() === '') return MO.galatLeader
  return null
}

/** Sub Branch milik Branch terpilih (`BrowseBranchDetail_RD` disaring `ParentID`); Branch kosong = semua. */
export function subBranchUntuk(p: Pick<Pilihan, 'subBranch'>, branch: string): Cabang[] {
  return branch === '' ? p.subBranch : p.subBranch.filter((c) => c.induk === branch)
}

/** Pilihan Leader tanpa baris itu sendiri. */
export function leaderUntuk(p: Pick<Pilihan, 'leader'>, id: string | null): OpsiLeader[] {
  return p.leader.filter((l) => l.id !== id)
}

/** Ganti Branch: Sub Branch yang bukan miliknya dikosongkan. */
export function setelBranch(isi: Isian, p: Pick<Pilihan, 'subBranch'>, branch: string): Isian {
  const masih = subBranchUntuk(p, branch).some((c) => c.id === isi.branchDetailId)
  return { ...isi, branchParent: branch, branchDetailId: masih ? isi.branchDetailId : '' }
}

/** Ganti Sub Branch: Branch mengikuti induknya. */
export function setelSubBranch(isi: Isian, p: Pick<Pilihan, 'subBranch'>, sub: string): Isian {
  const c = p.subBranch.find((x) => x.id === sub)
  return { ...isi, branchDetailId: sub, branchParent: c === undefined ? isi.branchParent : c.induk }
}

export type Saringan = 'semua' | 'aktif' | 'nonaktif'

/** Saring daftar: status, lalu setiap kata harus ada di salah satu kolom yang tampil. */
export function saring(d: readonly BarisMO[], kueri: string, s: Saringan): BarisMO[] {
  const kata = kueri.trim().toLowerCase().split(/\s+/).filter(Boolean)
  return d.filter((b) => {
    if (s === 'aktif' && b.moStatus !== '1') return false
    if (s === 'nonaktif' && b.moStatus === '1') return false
    const jerami = [b.id, b.clientName, b.clientId, b.aksesLogin, b.moLeader, b.branchDetailName, b.teamGroup]
      .join(' ')
      .toLowerCase()
    return kata.every((k) => jerami.includes(k))
  })
}

/** Keterangan akun di daftar: kosong bila akun aktif ada di M_LOGIN_GO. */
export function tandaAkun(b: Pick<BarisMO, 'statusAkun'>): string {
  switch (b.statusAkun) {
    case '':
      return MO.tanpaAkun
    case 'tidak-ada':
      return MO.akunTidakAda
    case 'nonaktif':
      return MO.akunNonaktif
    default:
      return ''
  }
}
