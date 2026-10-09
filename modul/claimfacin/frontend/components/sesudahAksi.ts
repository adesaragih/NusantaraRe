// Keputusan layar sesudah satu aksi berhasil - fungsi murni (diuji tanpa DOM).
//
//   modal    `Layar.bukaModal` membuka modalnya; tombol kaki modal yang berhasil tanpa pesan menutupnya; modal yang tidak
//            lagi dikirim server (kasus pindah tahap / selesai) ditutup.
//   pop-up   aksi server tanpa akibat (`halamanSaja` no-op di `services/aksi.go`) yang di Pega membuka harness: layar
//            membuka pop-up isinya SESUDAH server menerima aksi itu (aksi terbuka, isian tersimpan):
//              BukaSebab        -> CauseofLoss_Harness   (GET pilihan/sebab, Choose = GetNameCauseofLoss)
//              BukaKatastrofe   -> CatastrofeList        (GET pilihan/katastrofe, Choose = SetCatastrope,
//                                                          Save = SaveCatasrtope: daftar terbuka lagi dengan baris
//                                                          barunya, Pega InputCatastrope Cancel)
//              BukaOutstanding  -> Outstanding           (GET pilihan/outstanding?indeks=<objek>)
//              BukaRetro(List)  -> RetroList_Harnness    (FacRetroList halaman polis, hanya-baca)
//            SearchPolis (pop-up Choose Polis) -> grid hasil GET pilihan/polis di dalam modal.

import type { Layar } from '../api'
import { CFI } from '../labels'

/** Pop-up layar (bukan modal tata server). */
export type JenisPopup = 'sebab' | 'katastrofe' | 'outstanding' | 'retro'

const POPUP_SESUDAH: Record<string, JenisPopup> = {
  BukaSebab: 'sebab',
  BukaKatastrofe: 'katastrofe',
  SaveCatasrtope: 'katastrofe',
  BukaOutstanding: 'outstanding',
  BukaRetro: 'retro',
  BukaRetroList: 'retro',
}

/** Aksi tombol Search pop-up Choose Polis (SearchPolis_act). */
export const AKSI_CARI_POLIS = 'SearchPolis'

/** Lanjutan layar sesudah aksi berhasil. */
export type Lanjutan = { popup: JenisPopup; indeks: number } | { cariPolis: true } | null

export function lanjutanAksi(aksi: string, indeks: number): Lanjutan {
  const popup = POPUP_SESUDAH[aksi]
  if (popup) return { popup, indeks }
  if (aksi === AKSI_CARI_POLIS) return { cariPolis: true }
  return null
}

/** Modal yang terbuka sesudah aksi berhasil. `tutup` = aksi tombol kaki modal. */
export function modalSesudah(
  kini: string | null,
  l: Pick<Layar, 'bukaModal' | 'modal' | 'pesan'>,
  tutup: boolean,
): string | null {
  if (l.bukaModal) return l.bukaModal
  if (kini === null) return null
  if (tutup && (l.pesan ?? []).length === 0) return null
  return l.modal?.[kini] !== undefined ? kini : null
}

/** Modal Choose Polis (harness ChoosePolis, `services.ModalPilihPolis`). */
export const MODAL_POLIS = 'pilihPolis'

/**
 * Judul modal = pyWindowName tombol pembuka / pyLabel harness dan flow action. Kunci modal = `konteks` setiap aksi di
 * dalamnya (`services/layanan.go` tataKasus): `pla:<o>`, `dla:<o>`, `cedant:<adj>`, `komite:<adj>` membawa objek /
 * adjustment sasarannya.
 */
export function judulModal(m: string): string {
  if (m === MODAL_POLIS) return CFI.popPolis
  if (m === 'protectDOL') return CFI.popProtectDOL
  if (m === 'tutup') return CFI.popTutup
  if (m === 'tolak') return CFI.popTolak
  if (m.startsWith('pla:')) return CFI.popPLA
  if (m.startsWith('dla:')) return CFI.popDLA
  if (m.startsWith('cedant:')) return CFI.popCedant
  if (m.startsWith('komite:')) return CFI.popKomite
  return m
}
