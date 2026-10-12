import { useState, type FormEvent } from 'react'

import { TOMBOL } from '../labels'
import { PanelPindahTahap } from '../components/PanelPindahTahap'
import { PanelTutupKlaim } from '../components/PanelTutupKlaim'
import { RincianPeserta } from '../components/RincianPeserta'

import {
  ambilKlaimLife,
  bolehCabutPeserta,
  cabutPeserta,
  dampakHapusKlaim,
  hapusKlaim,
  type DampakHapus,
  type Klaim,
} from '../api'
import { kodeStatusGalat, pesanGalat } from '../../../../inti/frontend/klien'
import { PitaGalat } from '../components/PitaGalat'

// ============================================================================
// pages/KlaimLife.tsx — halaman "buka satu klaim Life" (tiket 01 AC-2).
//
// Cara membaca berkas ini bila baru mengenal React:
//   - Komponen = fungsi yang MENGEMBALIKAN tampilan. Bagian yang mirip HTML di
//     dalam `return (...)` disebut JSX; karena berkas ini TypeScript, namanya TSX.
//   - `useState` = kotak penyimpan nilai. Setiap kali nilainya diganti lewat
//     fungsi set-nya (misalnya `setKlaim`), React menggambar ulang tampilan.
//   - `<string>` di belakang useState memberi tahu TypeScript isi kotak itu apa.
//   - `{ ... }` di dalam JSX = "sisipkan nilai JavaScript di sini".
//
// ⛔ Nilai uang tidak pernah disentuh aritmetika dan tidak pernah dibungkus
// Number(). Ia datang sebagai teks desimal dan ditampilkan apa adanya
// (ADR-U-0003, ADR-U-0016). Penjelasannya di modul/claimlife/frontend/api.ts bab 3.
// ============================================================================
export default function KlaimLife() {
  // Empat kotak keadaan halaman ini.
  const [id, setId] = useState<string>('') // yang sedang diketik pengguna
  const [klaim, setKlaim] = useState<Klaim | null>(null) // klaim yang berhasil dibaca
  const [galat, setGalat] = useState<string | null>(null) // pesan bila gagal
  const [memuat, setMemuat] = useState<boolean>(false) // sedang menunggu server?
  const [mencabut, setMencabut] = useState<string | null>(null) // peserta yang sedang dicabut
  // Nomor klaim yang BARU SAJA ditutup - dipakai kalimat penutupnya.
  const [tutupSelesai, setTutupSelesai] = useState<string | null>(null)

  // ⛔ 10-10-2026: rincian per peserta (tanggal, total, dokumen, diagnosa, adjustment, Reject, Komite) dan Close Claim
  // DIPINDAH ke components/RincianPeserta.tsx dan components/PanelTutupKlaim.tsx - dipakai bersama layar Outstanding.

  const [dampak, setDampak] = useState<DampakHapus | null>(null) // isi popup konfirmasi
  const [menghapus, setMenghapus] = useState<boolean>(false) // permintaan hapus berjalan

  // Membuka popup: MEMBACA dampak, tidak menghapus apa pun.
  async function bukaKonfirmasiHapus() {
    if (!klaim) return
    setGalat(null)
    try {
      setDampak(await dampakHapusKlaim(klaim.id))
    } catch {
      setGalat('Failed to calculate the deletion impact.')
    }
  }

  // ⛔ Batal hanya menutup popup. Tidak ada yang perlu dibatalkan, sebab
  // membuka popup tidak menulis apa pun.
  function batalHapus() {
    setDampak(null)
    setGalat(null)
  }

  async function lanjutkanHapus() {
    if (!klaim) return
    setGalat(null)
    setMenghapus(true)
    try {
      await hapusKlaim(klaim.id)
      setKlaim(null)
      setDampak(null)
    } catch (err: unknown) {
      const kode = kodeStatusGalat(err)
      setGalat(
        // ⛔ 405, bukan 501 (GILIRAN-11): backend menjawab 405 + `Allow`
        // sejak ADR-U-0031 ditegakkan; layar dulu menunggu 501 dan jatuh ke
        // "Gagal menghapus klaim." - kalimat yang menyembunyikan sebabnya.
        kode === 405 || kode === 501
          ? (pesanGalat(err) ??
            'Claim deletion is a marker, not a physical delete (ADR-U-0031). The marker column has not been decided by the work owner.')
          : kode === 409
            ? // Kalimat server lebih dulu: 409 kini juga berarti kasus tertutup.
              (pesanGalat(err) ?? 'The claim has been sent to the committee and cannot be deleted.')
            : kode === 403
              ? 'Only ReasLifeAdmin can delete a claim.'
              : 'Failed to delete the claim.',
      )
      setDampak(null)
    } finally {
      setMenghapus(false)
    }
  }

  // OQ-M6 (GILIRAN-17): tombol `DELETE` b17865 - tanpa konfirmasi (b18021).
  async function cabut(pesertaID: string) {
    if (!klaim) return
    setGalat(null)
    setMencabut(pesertaID)
    try {
      await cabutPeserta(klaim.id, pesertaID)
    } catch (err: unknown) {
      setGalat(pesanGalat(err) ?? 'Failed to remove the insured.')
      setMencabut(null)
      return
    }
    // Pencabutan SUDAH berhasil; gagal membaca ulang bukan gagal mencabut.
    try {
      setKlaim(await ambilKlaimLife(klaim.id))
    } catch {
      setGalat('The insured was removed, but the claim could not be reloaded; press Open.')
    } finally {
      setMencabut(null)
    }
  }

  // Membaca ulang klaim yang SEDANG terbuka.
  //
  // ⚠️ Dipakai sesudah perpindahan tahap: tahapnya berubah, dan tombol
  // yang ditawarkan layar ikut berubah. Tanpa pembacaan ulang, layar tetap
  // menawarkan tombol tahap LAMA - dan tiap kliknya ditolak 409.
  async function muatUlang(): Promise<void> {
    if (klaim === null) return
    try {
      setKlaim(await ambilKlaimLife(klaim.id))
    } catch {
      setGalat('Failed to reload the claim after the move.')
    }
  }

  // Dipanggil saat tombol "Buka" ditekan (form dikirim).
  async function cari(e: FormEvent<HTMLFormElement>) {
    e.preventDefault() // jangan biarkan browser memuat ulang halaman
    setGalat(null)
    setKlaim(null)
    // Kalimat penutup klaim SEBELUMNYA dibersihkan: membiarkannya membuat
    // orang membaca "Klaim X ditutup" di atas detail klaim Y.
    setTutupSelesai(null)
    setMemuat(true)
    try {
      setKlaim(await ambilKlaimLife(id.trim()))
    } catch (err: unknown) {
      setGalat(kodeStatusGalat(err) === 404 ? 'The claim does not exist.' : 'Failed to load the claim.')
    } finally {
      setMemuat(false) // apa pun hasilnya, berhenti "memuat"
    }
  }

  return (
    <section>
      <h2>Claim Life</h2>

      <form onSubmit={cari}>
        <label htmlFor="idKlaim">Case ID</label>{' '}
        <input
          id="idKlaim"
          value={id}
          onChange={(e) => setId(e.target.value)}
          placeholder="e.g. CLMLF-000001"
        />{' '}
        <button type="submit" disabled={memuat || id.trim() === ''}>
          {memuat ? 'Loading…' : 'Open'}
        </button>
      </form>

      {/* PitaGalat (11-10-2026): tempat galat yang sama dengan layar Claim Life lain. */}
      <PitaGalat galat={galat ?? ''} onTutup={() => setGalat(null)} />

      {/* ⛔ Kalimat penutup, sesudah `closeContainer` melepas klaimnya dari
          layar. Tanpa ini penutupan terlihat seperti layar yang kosong tiba-
          tiba, dan orang akan menekan Buka lagi untuk memeriksa. */}
      {tutupSelesai !== null && (
        <p role="status">
          Claim {tutupSelesai} closed. The case is completed and removed from the inbox;
          nothing can be changed on it anymore.
        </p>
      )}

      {klaim && (
        <article>
          <h3>{klaim.nomorKlaim || klaim.id}</h3>
          <dl>
            <dt>Policy No</dt>
            <dd>{klaim.nomorPolis || '—'}</dd>
            <dt>Business Name</dt>
            <dd>{klaim.namaBisnis || '—'}</dd>
            <dt>Number of Rows</dt>
            <dd>{klaim.cacahBaris}</dd>
            {/* Status klaim adalah TURUNAN dari baris-barisnya, bukan kolom
                tersimpan. Yang ditampilkan kata, bukan angka - dan bukan pula
                terjemahan `kodeStatus`, yang hanya cerminan baris terakhir. */}
            <dt>Claim Status</dt>
            <dd>{klaim.statusTurunan}</dd>
          </dl>

          {/* `.map` = ulangi blok di bawah untuk SETIAP peserta.
              `key` wajib ada supaya React tahu baris mana yang berubah. */}
          <p>
            <button type="button" onClick={() => void bukaKonfirmasiHapus()}>
              Delete Claim…
            </button>
          </p>

          {/* Popup konfirmasi: angka dulu, keputusan kemudian. Tiap jenis
              disebut sendiri - satu angka total menyembunyikan tingkat mana
              yang ternyata lebih besar dari dugaan. */}
          {dampak && (
            <aside>
              <h4>Delete claim {klaim.nomorKlaim || klaim.id}?</h4>
              <p>The following will also be deleted:</p>
              <ul>
                <li>Claim header: {dampak.header}</li>
                <li>Insured: {dampak.peserta}</li>
                <li>Adjustment rows: {dampak.adjustment}</li>
                <li>Spreading: {dampak.spreading}</li>
                <li>Spreading retro: {dampak.spreadingRetro}</li>
                <li>Documents: {dampak.dokumen}</li>
                <li>Work rows: {dampak.barisWork}</li>
              </ul>
              <p>Total {dampak.total} rows.</p>
              {/* ⛔ Baris datar warisan BUKAN bagian daftar di atas dan tidak
                  dijumlahkan ke totalnya. Nasibnya saat klaim dihapus belum
                  diputuskan work owner; menaruhnya di bawah judul "yang akan
                  ikut terhapus" berarti menjawab pertanyaan itu diam-diam. */}
              <p>
                ⚠️ There are also {dampak.barisDatarWarisan} legacy flat rows with the same CASEID.{' '}
                <strong>Their handling has not been decided</strong>: whether they are deleted too or kept
                because downstream systems already read them.
              </p>
              <button
                type="button"
                disabled={menghapus}
                onClick={() => void lanjutkanHapus()}
              >
                {menghapus ? 'Deleting…' : 'Yes, delete'}
              </button>{' '}
              <button type="button" onClick={batalHapus}>
                Cancel
              </button>
            </aside>
          )}

          {klaim.peserta.map((p) => (
            <section key={p.id}>
              <h4>
                Insured {p.nomorSertifikat || p.id}{' '}
                <small>({p.baris.length} rows)</small>
                {bolehCabutPeserta(klaim) && (
                  <>
                    {' '}
                    <button
                      type="button"
                      aria-label={`${TOMBOL.cabutPeserta} insured ${p.nomorSertifikat || p.id}`}
                      disabled={mencabut === p.id}
                      onClick={() => void cabut(p.id)}
                    >
                      {mencabut === p.id ? 'Removing…' : TOMBOL.cabutPeserta}
                    </button>
                  </>
                )}
              </h4>
              {/* Isi rincian peserta - satu komponen bersama dengan baris terbuka grid Participant Details layar
                  Outstanding (10-10-2026). */}
              <RincianPeserta klaim={klaim} p={p} onMuatUlang={muatUlang} />
            </section>
          ))}
        </article>
      )}

      {klaim !== null && (
        <section className="os__aksi">
          {/* ⛔ Tiap tahap menawarkan PERSIS perpindahan yang section-nya
              punya - lihat PanelPindahTahap. Menyalin tombol antarlayar
              membuka jalur yang di sistem lama tidak ada, dan backend
              menolaknya 409. */}
          <PanelPindahTahap
            klaimID={klaim.id}
            tahap={klaim.tahap}
            sesudahPindah={() => {
              void muatUlang()
            }}
          />

          {/* Close Claim - components/PanelTutupKlaim.tsx (dipindah 10-10-2026, dipakai bersama layar Outstanding). */}
          <PanelTutupKlaim
            klaim={klaim}
            sesudahTutup={(nomor) => {
              // Padanan `closeContainer` b1129: klaim dilepas dari layar, bukan dimuat ulang.
              setKlaim(null)
              setTutupSelesai(nomor)
            }}
          />
        </section>
      )}
    </section>
  )
}
