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
      setGalat('Gagal menghitung dampak penghapusan.')
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
            'Penghapusan klaim adalah penanda, bukan hapus fisik (ADR-U-0031). Kolom penandanya belum diputuskan work owner.')
          : kode === 409
            ? // Kalimat server lebih dulu: 409 kini juga berarti kasus tertutup.
              (pesanGalat(err) ?? 'Klaim sudah diserahkan ke Komite dan tidak dapat dihapus.')
            : kode === 403
              ? 'Hanya ReasLifeAdmin yang dapat menghapus klaim.'
              : 'Gagal menghapus klaim.',
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
      setGalat(pesanGalat(err) ?? 'Gagal mencabut peserta.')
      setMencabut(null)
      return
    }
    // Pencabutan SUDAH berhasil; gagal membaca ulang bukan gagal mencabut.
    try {
      setKlaim(await ambilKlaimLife(klaim.id))
    } catch {
      setGalat('Peserta sudah dicabut, tetapi klaim gagal dibaca ulang; tekan Buka.')
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
      setGalat('Gagal membaca ulang klaim sesudah perpindahan.')
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
      setGalat(kodeStatusGalat(err) === 404 ? 'Klaim tidak ada.' : 'Gagal membaca klaim.')
    } finally {
      setMemuat(false) // apa pun hasilnya, berhenti "memuat"
    }
  }

  return (
    <section>
      <h2>Klaim Life</h2>

      <form onSubmit={cari}>
        <label htmlFor="idKlaim">Pengenal klaim</label>{' '}
        <input
          id="idKlaim"
          value={id}
          onChange={(e) => setId(e.target.value)}
          placeholder="mis. UJI-KLAIM-1"
        />{' '}
        <button type="submit" disabled={memuat || id.trim() === ''}>
          {memuat ? 'Memuat…' : 'Buka'}
        </button>
      </form>

      {/* PitaGalat (11-10-2026): tempat galat yang sama dengan layar Claim Life lain. */}
      <PitaGalat galat={galat ?? ''} onTutup={() => setGalat(null)} />

      {/* ⛔ Kalimat penutup, sesudah `closeContainer` melepas klaimnya dari
          layar. Tanpa ini penutupan terlihat seperti layar yang kosong tiba-
          tiba, dan orang akan menekan Buka lagi untuk memeriksa. */}
      {tutupSelesai !== null && (
        <p role="status">
          Klaim {tutupSelesai} ditutup. Kasusnya selesai dan hilang dari kotak
          masuk; tidak ada lagi yang dapat diubah padanya.
        </p>
      )}

      {klaim && (
        <article>
          <h3>{klaim.nomorKlaim || klaim.id}</h3>
          <dl>
            <dt>Nomor polis</dt>
            <dd>{klaim.nomorPolis || '—'}</dd>
            <dt>Nama bisnis</dt>
            <dd>{klaim.namaBisnis || '—'}</dd>
            <dt>Jumlah baris</dt>
            <dd>{klaim.cacahBaris}</dd>
            {/* Status klaim adalah TURUNAN dari baris-barisnya, bukan kolom
                tersimpan. Yang ditampilkan kata, bukan angka - dan bukan pula
                terjemahan `kodeStatus`, yang hanya cerminan baris terakhir. */}
            <dt>Status klaim</dt>
            <dd>{klaim.statusTurunan}</dd>
          </dl>

          {/* `.map` = ulangi blok di bawah untuk SETIAP peserta.
              `key` wajib ada supaya React tahu baris mana yang berubah. */}
          <p>
            <button type="button" onClick={() => void bukaKonfirmasiHapus()}>
              Hapus klaim…
            </button>
          </p>

          {/* Popup konfirmasi: angka dulu, keputusan kemudian. Tiap jenis
              disebut sendiri - satu angka total menyembunyikan tingkat mana
              yang ternyata lebih besar dari dugaan. */}
          {dampak && (
            <aside>
              <h4>Hapus klaim {klaim.nomorKlaim || klaim.id}?</h4>
              <p>Yang akan ikut terhapus:</p>
              <ul>
                <li>Header klaim: {dampak.header}</li>
                <li>Peserta: {dampak.peserta}</li>
                <li>Baris adjustment: {dampak.adjustment}</li>
                <li>Spreading: {dampak.spreading}</li>
                <li>Spreading retro: {dampak.spreadingRetro}</li>
                <li>Dokumen: {dampak.dokumen}</li>
                <li>Baris work: {dampak.barisWork}</li>
              </ul>
              <p>Total {dampak.total} baris.</p>
              {/* ⛔ Baris datar warisan BUKAN bagian daftar di atas dan tidak
                  dijumlahkan ke totalnya. Nasibnya saat klaim dihapus belum
                  diputuskan work owner; menaruhnya di bawah judul "yang akan
                  ikut terhapus" berarti menjawab pertanyaan itu diam-diam. */}
              <p>
                ⚠️ Selain itu terdapat {dampak.barisDatarWarisan} baris datar warisan
                ber-CASEID sama. <strong>Nasibnya belum diputuskan</strong>: apakah ikut
                terhapus atau ditinggal karena hilir sudah membacanya.
              </p>
              <button
                type="button"
                disabled={menghapus}
                onClick={() => void lanjutkanHapus()}
              >
                {menghapus ? 'Menghapus…' : 'Ya, hapus'}
              </button>{' '}
              <button type="button" onClick={batalHapus}>
                Batal
              </button>
            </aside>
          )}

          {klaim.peserta.map((p) => (
            <section key={p.id}>
              <h4>
                Peserta {p.nomorSertifikat || p.id}{' '}
                <small>({p.baris.length} baris)</small>
                {bolehCabutPeserta(klaim) && (
                  <>
                    {' '}
                    <button
                      type="button"
                      aria-label={`${TOMBOL.cabutPeserta} peserta ${p.nomorSertifikat || p.id}`}
                      disabled={mencabut === p.id}
                      onClick={() => void cabut(p.id)}
                    >
                      {mencabut === p.id ? 'Mencabut…' : TOMBOL.cabutPeserta}
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
